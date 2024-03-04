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

func IsAuthenticated(ctx context.Context, telemetry telemetry.Telemetry) func(next http.Handler) http.Handler {
	return Handler(ctx, telemetry)
}

func Handler(ctx context.Context, telemetry telemetry.Telemetry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			_, span := telemetry.Start(ctx, "IsAuthenticated")
			defer span.End()
			errorMessage := "Erro na autenticação"
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				err := errors.New("Unauthorized")
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				respondWithError(rw, http.StatusUnauthorized, err.Error(), errorMessage)
				return
			}

			//@todo não precisa fazer o http get, deve usar o user.ValidateToken
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

			newCTX := context.WithValue(r.Context(), "email", email)
			next.ServeHTTP(rw, r.WithContext(newCTX))
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
