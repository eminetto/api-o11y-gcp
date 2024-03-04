package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/eminetto/api-o11y-gcp/auth/security"
	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/eminetto/api-o11y-gcp/user"
	"github.com/go-chi/httplog"
	"go.opentelemetry.io/otel/codes"
)

func UserAuth(ctx context.Context, uService user.UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		ctx, span := otel.Start(ctx, "userAuth")
		defer span.End()
		var param struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		err = uService.ValidateUser(ctx, param.Email, param.Password)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			oplog.Error().Msg(err.Error())
			return
		}
		var result struct {
			Token string `json:"token"`
		}
		result.Token, err = security.NewToken(param.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			oplog.Error().Msg(err.Error())
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}

		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return
		}
		return
	}
}
