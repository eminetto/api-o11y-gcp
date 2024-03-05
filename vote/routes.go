package vote

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/go-chi/httplog"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

func Store(ctx context.Context, vService UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		ctx, span := otel.Start(ctx, "vote:store")
		defer span.End()
		var v Vote
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		v.Email = r.Context().Value("email").(string)
		var result struct {
			ID uuid.UUID `json:"id"`
		}
		result.ID, err = vService.Store(ctx, &v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			oplog.Error().Msg(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			oplog.Error().Msg(err.Error())
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		return
	}
}
