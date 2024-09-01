package http

import (
    "net/http"
)

func JSONError(w http.ResponseWriter, r *http.Request, status int, err error) {
    if err != nil {
        w.WriteHeader(status)
        _, err := w.Write([]byte(err.Error()))
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
        }
    }
}
