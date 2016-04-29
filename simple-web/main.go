package main

import (
    "fmt"
    "net/http"
    "strings"
    "log"
    "html/template"
    "time"
    "crypto/md5"
    "io"
    "strconv"
)

func sayHelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Println(r.Form)
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key: ", k)
        fmt.Println("val: ", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello, astaxie!")
}

func login(w http.ResponseWriter, r *http.Request)  {
    fmt.Println("method: ", r.Method)
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, token)
    } else {
        r.ParseForm()

        fmt.Println("---------- FORM POST -------------")
        fmt.Println(r.Form)
        //fmt.Println("gender: ", r.Form["gender"])
        //fmt.Println("gender: ", r.Form.Get("gender"))
        //fmt.Println("username1: ", template.HTMLEscapeString(r.Form.Get("username")))
        //fmt.Println("username: ", r.Form["username"])
        //fmt.Println("password: ", r.Form["password"])

        token := r.Form.Get("token")
        if token != "" {
            // check token validity
        } else {
            // give error if no token
        }
        template.HTMLEscape(w, []byte(r.Form.Get("username"))) // respond to client

    }
}

func main() {
    http.HandleFunc("/", sayHelloName)
    http.HandleFunc("/login", login)
    err := http.ListenAndServe(":9090", nil)

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
