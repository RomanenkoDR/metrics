package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func CheckReqSign(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("HashSHA256") == "" {
				fmt.Println("Skip hash calculation")
				next.ServeHTTP(w, r)
				return
			}

			var sign []byte

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Calcucalate hash for request body
			h := hmac.New(sha256.New, []byte(key))
			h.Write(body)
			sha256sum := h.Sum(nil)

			// Read sign from request
			sign, err = hex.DecodeString(r.Header.Get("HashSHA256"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if !hmac.Equal(sign, sha256sum) {
				http.Error(w, "Corrupted sign on request. ", http.StatusBadRequest)
				return
			}

			// Add calculated hash to headers
			w.Header().Set("HashSHA256", hex.EncodeToString(sha256sum))

			// Restore request body for further processing
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			next.ServeHTTP(w, r)
		})
	}
}
