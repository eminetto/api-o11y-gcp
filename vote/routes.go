package vote

import (
	"encoding/json"
	"net/http"

	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
)

func Store(vService UseCase, otel telemetry.Telemetry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Start(r.Context(), "vote: store")
		defer span.End()
		var v Vote
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
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
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
}
