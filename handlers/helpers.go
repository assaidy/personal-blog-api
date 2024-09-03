package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/assaidy/personal-blog-api/types"
)

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

func Make(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("HTTP Request", "path", r.URL.Path, "from", r.RemoteAddr)
		if err := f(w, r); err != nil {
			if apiErr, ok := err.(types.ApiError); ok {
				WriteJSON(w, apiErr.StatusCode, apiErr)
			} else {
				resp := map[string]string{
					"statusCode": "500",
					"msg":        "internal server error",
				}
				WriteJSON(w, http.StatusInternalServerError, resp)
			}
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to encode JSON response", "err", err.Error())
		return err
	}

	return nil
}
