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
		slog.Info("HTTP Request", "from", r.RemoteAddr, "path", r.URL.Path)
		if err := f(w, r); err != nil {
			if apiErr, ok := err.(types.ApiError); ok {
				slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
				WriteJSON(w, apiErr.StatusCode, apiErr)
			} else {
				resp := map[string]string{
					"statusCode": "500",
					"msg":        "internal server error",
				}
				slog.Error("INTERNAL error", "err", err.Error(), "path", r.URL.Path)
				WriteJSON(w, http.StatusInternalServerError, resp)
			}
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
