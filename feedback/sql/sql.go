package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/eminetto/api-o11y-gcp/feedback"
	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
)

// SQL repo
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
func (r *SQL) Store(ctx context.Context, f *feedback.Feedback) error {
	ctx, span := r.telemetry.Start(ctx, "feedback: sql")
	defer span.End()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer tx.Commit()
	stmt, err := tx.PrepareContext(ctx, `
		insert into feedback (id, email, title, body, created_at) 
		values(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,
		f.ID,
		f.Email,
		f.Title,
		f.Body,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	return nil
}
