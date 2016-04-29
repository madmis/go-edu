package httputil

import (
    "net/http"
    "fmt"
    "errors"
)

func ValidateRequestMethod(r *http.Request, method string) error {
    if r.Method != method {
        return errors.New(fmt.Sprintf("Only %s allowed", method))
    }

    return nil
}
