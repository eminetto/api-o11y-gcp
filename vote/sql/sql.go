package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/eminetto/api-o11y-gcp/vote"
	"go.opentelemetry.io/otel/codes"
)

// SQL mysql repo
type SQL struct {
	db        *sql.DB
	telemetry telemetry.Telemetry
}

// NewSQL create new repository
func NewSQL(db *sql.DB, telemetry telemetry.Telemetry) *SQL {
	return &SQL{
		db:        db,
		telemetry: telemetry,
	}
}

// Store a feedback
func (r *SQL) Store(ctx context.Context, v *vote.Vote) error {
	ctx, span := r.telemetry.Start(ctx, "vote:mysql")
	defer span.End()
	stmt, err := r.db.Prepare(`
		insert into vote (id, email, talk_name, score, created_at) 
		values(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		v.ID,
		v.Email,
		v.TalkName,
		v.Score,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		span.RecordError(err)
		return err
	}
	err = stmt.Close()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
