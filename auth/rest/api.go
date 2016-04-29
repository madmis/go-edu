package main

import (
    "net/http"
    "log"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/madmis/go-edu/auth/rest-pkgs/model"
    "github.com/madmis/go-edu/auth/rest-pkgs/httputil"
    "golang.org/x/crypto/bcrypt"
    "strings"
    "errors"
)

var db *sql.DB

// Check errors
// For now only throw panic
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

// Init database connection
func init() {
    var err error
    db, err = sql.Open("sqlite3", "./auth.db")
    checkErr(err)

    err = db.Ping()
    checkErr(err)
}

func main() {
    defer db.Close()

    http.HandleFunc("/api/", mainAction)
    http.HandleFunc("/api/register", registerAction)
    http.HandleFunc("/api/login", loginAction)
    http.HandleFunc("/api/profile", profileAction)

    log.Fatal(http.ListenAndServe(":9090", nil))
}

// Get users list action
func mainAction(w http.ResponseWriter, r *http.Request) {

    users, err := model.NewUser("", "").FindAll(db)
    checkErr(err)

    data := map[string]interface{}{
        "users": users,
    }

    httputil.SendSuccessResponse(w, http.StatusOK, data)
}

// Register user action
func registerAction(w http.ResponseWriter, r *http.Request) {
    if err := httputil.ValidateRequestMethod(r, "POST"); err != nil {
        err := httputil.NewErrors(err.Error(), http.StatusMethodNotAllowed)
        httputil.SendErrorResponse(w, http.StatusMethodNotAllowed, *err)

        return
    }

    user := model.NewUser(r.FormValue("name"), r.FormValue("email"))
    user.Password = r.FormValue("password")

    if ok, errors := user.Validate(); !ok {
        err := httputil.NewErrors(
            http.StatusText(http.StatusBadRequest),
            http.StatusBadRequest,
        )

        for k, v := range errors {
            err.AppendError(v, k)
        }

        httputil.SendErrorResponse(w, http.StatusBadRequest, *err)

        return
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    checkErr(err)
    user.Password = string(hashedPassword)

    // query
    err = user.FindByEmail(db)
    if err != sql.ErrNoRows {
        err := httputil.NewErrors("User with same email exists", http.StatusBadRequest)

        httputil.SendErrorResponse(w, http.StatusBadRequest, *err)

        return
    }

    err = user.Create(db)
    checkErr(err)

    userMap := map[string]interface{}{
        "user": *user,
    }
    httputil.SendSuccessResponse(w, http.StatusCreated, userMap)
}

// Get access token action
func loginAction(w http.ResponseWriter, r *http.Request) {
    if err := httputil.ValidateRequestMethod(r, "POST"); err != nil {
        err := httputil.NewErrors(err.Error(), http.StatusMethodNotAllowed)
        httputil.SendErrorResponse(w, http.StatusMethodNotAllowed, *err)

        return
    }

    user := model.NewUser("test", r.FormValue("email"))
    password := r.FormValue("password")
    user.Password = password
    if ok, errors := user.Validate(); !ok {
        err := httputil.NewErrors(
            http.StatusText(http.StatusBadRequest),
            http.StatusBadRequest,
        )

        for k, v := range errors {
            err.AppendError(v, k)
        }

        httputil.SendErrorResponse(w, http.StatusBadRequest, *err)

        return
    }

    err := user.FindByEmail(db)
    if err != nil {
        err := httputil.NewErrors("Bad credentials", http.StatusBadRequest)

        httputil.SendErrorResponse(w, http.StatusBadRequest, *err)

        return
    }

    // Comparing the password with the hash
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        err := httputil.NewErrors("Bad credentials", http.StatusBadRequest)

        httputil.SendErrorResponse(w, http.StatusBadRequest, *err)

        return
    }

    // token - this is user password.
    // this case used only for test application
    // for real application should be used jwt golang implementation
    data := map[string]interface{}{
        "token": user.Password,
    }

    httputil.SendSuccessResponse(w, http.StatusOK, data)
}

// Get user profile action
// This action require authentication (access token)
func profileAction(w http.ResponseWriter, r *http.Request) {
    if origin := r.Header.Get("Origin"); origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    }

    // Stop here for a Preflighted OPTIONS request.
    if r.Method == "OPTIONS" {
        return
    }


    if err := httputil.ValidateRequestMethod(r, "GET"); err != nil {
        err := httputil.NewErrors(err.Error(), http.StatusMethodNotAllowed)
        httputil.SendErrorResponse(w, http.StatusMethodNotAllowed, *err)

        return
    }

    token, err := parseTokenFromRequest(r)
    if err != nil {
        err := httputil.NewErrors(err.Error(), http.StatusBadRequest)
        httputil.SendErrorResponse(w, http.StatusBadRequest, *err)

        return
    }

    user := model.User{}
    err = user.FindByToken(db, token)
    if err != nil {
        if err == sql.ErrNoRows {
            err := httputil.NewErrors("User not found", http.StatusNotFound)
            httputil.SendErrorResponse(w, http.StatusNotFound, *err)
        } else {
            err := httputil.NewErrors(err.Error(), http.StatusInternalServerError)
            httputil.SendErrorResponse(w, http.StatusInternalServerError, *err)
        }

        return
    }

    data := map[string]interface{}{
        "profile": user,
    }

    httputil.SendSuccessResponse(w, http.StatusOK, data)
}


// Try to find the token in an http.Request.
func parseTokenFromRequest(req *http.Request) (token string, err error) {
    // Look for an Authorization header
    if ah := req.Header.Get("Authorization"); ah != "" {
        // Should be a bearer token
        if len(ah) > 6 && strings.ToUpper(ah[0:7]) == "BEARER " {
            return strings.TrimSpace(ah[7:]), nil
        }
    }

    return "", errors.New("No token present in request")
}
