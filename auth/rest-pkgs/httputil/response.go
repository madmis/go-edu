package httputil

import (
    "net/http"
    "encoding/json"
)

type Error struct {
    Message string  `json:"message"`
    Key     string  `json:"key"`
}

type Errors struct {
    Message string  `json:"message"`
    Code    int     `json:"code"`
    Errors  []Error `json:"errors"`
}

type Response struct {
    Data   interface{}  `json:"data,omitempty"`
    Errors interface{}  `json:"errors,omitempty"`
}

func NewErrors(message string, code int) (errors *Errors) {
    res := Errors{
        Message: message,
        Code: code,
    }

    res.Errors = make([]Error, 0)

    return &res
}

func (err *Errors) AppendError(message string, key string) {
    item := Error{message, ""}
    if key != "" {
        item.Key = key
    }

    err.Errors = append(err.Errors, item)
}

func SendErrorResponse(w http.ResponseWriter, httpCode int, errors Errors) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(httpCode)

    response := new(Response)
    response.Data = nil
    response.Errors = errors

    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}

func SendSuccessResponse(w http.ResponseWriter, httpCode int, data map[string]interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(httpCode)

    response := new(Response)
    response.Errors = nil
    response.Data = data

    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}
