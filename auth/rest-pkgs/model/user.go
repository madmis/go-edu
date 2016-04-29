package model

import (
    "time"
    "regexp"
    "strings"
    "database/sql"
)

const tableCreate = `CREATE TABLE users (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        name VARCHAR(64) NULL,
                        email VARCHAR(64) NULL,
                        password VARCHAR(64) NULL,
                        "created" DATETIME NULL
                     );`

type User struct {
    Id int                      `json:"id"`
    Name string                 `json:"name"`
    Email string                `json:"email"`
    Password string             `json:"-"`
    CreatedAt time.Time         `json:"created_at"`
}

func NewUser(name string, email string) (*User) {
    return &User{
        Name: name,
        Email: email,
        Password: "",
        CreatedAt: time.Now(),
    }
}

func (user *User) Validate() (bool, map[string]string) {
    errors := make(map[string]string)

    re := regexp.MustCompile(".+@.+\\..+")
    matched := re.Match([]byte(user.Email))
    if matched == false {
        errors["email"] = "Please enter a valid email address"
    }

    if strings.TrimSpace(user.Name) == "" {
        errors["name"] = "Please enter name"
    }

    if user.Password == "" {
        errors["password"] = "Password required"
    }

    return len(errors) == 0, errors
}

func (user *User) FindAll(db *sql.DB) ([]*User, error) {
    users := make([]*User, 0)

    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        return users, err
    }
    defer rows.Close()


    for rows.Next() {
        usr := new(User)
        err := rows.Scan(&usr.Id, &usr.Name, &usr.Email, &user.Password, &usr.CreatedAt)
        if err != nil {
            return users, err
        }
        users = append(users, usr)
    }
    if err = rows.Err(); err != nil {
        return users, err
    }

    return users, err
}


func (user *User) Find(db *sql.DB, id int) (error) {
    err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

    return err
}

func (user *User) FindByEmail(db *sql.DB) (error) {
    row := db.QueryRow("SELECT * FROM users WHERE email = ?", user.Email)
    err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

    return err
}

func (user *User) FindByToken(db *sql.DB, token string) (error) {
    row := db.QueryRow("SELECT * FROM users WHERE password = ?", token)
    err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

    return err
}

func (user *User) Create(db *sql.DB) (error) {
    // insert
    stmt, err := db.Prepare("INSERT INTO users(name, email, password, created) values(?, ?, ?, ?)")
    if err != nil {
        return err
    }

    res, err := stmt.Exec(user.Name, user.Email, user.Password, user.CreatedAt)
    if err != nil {
        return err
    }

    id, err := res.LastInsertId()
    if err == nil {
        user.Id = int(id)
    }

    return err
}