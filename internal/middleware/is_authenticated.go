package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/eminetto/api-o11y-gcp/auth/security"
	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
)

func IsAuthenticated(ctx context.Context, telemetry telemetry.Telemetry, traceName string) func(next http.Handler) http.Handler {
	return Handler(ctx, telemetry, traceName)
}

func Handler(ctx context.Context, telemetry telemetry.Telemetry, traceName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			newCTX, span := telemetry.Start(ctx, traceName+": isAuthenticated")
			defer span.End()
			errorMessage := "Erro na autenticação" // traduzir pra ingles e remover o Erro do começo
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				err := errors.New("Unauthorized")
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errorMessage)
				return
			}

			t, err := security.ParseToken(tokenString)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errorMessage)
				return
			}
			tData, err := security.GetClaims(t)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errorMessage)
				return
			}
			email := tData["email"].(string)
			type key string

			next.ServeHTTP(rw, r.WithContext(context.WithValue(newCTX, "email", email)))
		}
		return http.HandlerFunc(fn)
	}
}

// RespondWithError return a http error
func respondWithError(w http.ResponseWriter, code int, e string, message string) {
	respondWithJSON(w, code, map[string]string{"code": strconv.Itoa(code), "error": e, "message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
