package sql

import (
	"context"
	"database/sql"

	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/eminetto/api-o11y-gcp/user"
	"go.opentelemetry.io/otel/codes"
)

// SQL sql repo
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

// Get an user
func (r *SQL) Get(ctx context.Context, email string) (*user.User, error) {
	ctx, span := r.telemetry.Start(ctx, "sql")
	defer span.End()
	stmt, err := r.db.PrepareContext(ctx, `select id, email, password, first_name, last_name from user where email = ?`)
	if err != nil {
		return nil, err
	}
	var u user.User
	rows, err := stmt.QueryContext(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}
	return &u, nil
}
