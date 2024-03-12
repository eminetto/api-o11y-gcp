package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/eminetto/api-o11y-gcp/auth"
	"github.com/eminetto/api-o11y-gcp/feedback"
	sql_feedback "github.com/eminetto/api-o11y-gcp/feedback/sql"
	"github.com/eminetto/api-o11y-gcp/internal/middleware"
	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/eminetto/api-o11y-gcp/user"
	sql_user "github.com/eminetto/api-o11y-gcp/user/sql"
	"github.com/eminetto/api-o11y-gcp/vote"
	sql_vote "github.com/eminetto/api-o11y-gcp/vote/sql"
	"github.com/go-chi/chi"
	telemetrymiddleware "github.com/go-chi/telemetry"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
	}

	db, err := otelsql.Open("sqlite3", "./ops/db/api.db", otelsql.WithAttributes(
		semconv.DBSystemSqlite,
	))
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemSqlite,
	))
	if err != nil {
		logger.Error(err.Error())
	}

	ctx := context.Background()
	otel, err := telemetry.NewGCP(ctx, "api-o11y-gcp")
	if err != nil {
		logger.Error(err.Error())
	}
	defer otel.Shutdown(ctx)

	uRepo := sql_user.NewSQL(db, otel)
	uService := user.NewService(uRepo, otel)

	fRepo := sql_feedback.NewSQL(db, otel)
	fService := feedback.NewService(fRepo, otel)

	vRepo := sql_vote.NewSQL(db, otel)
	vService := vote.NewService(vRepo, otel)

	r := chi.NewRouter()
	r.Use(telemetrymiddleware.Collector(telemetrymiddleware.Config{
		AllowAny: true,
	}, []string{"/v1"})) // path prefix filters basically records generic http request metrics
	r.Post("/v1/auth", auth.UserAuth(ctx, uService, otel))

	r.Route("/v1/feedback", func(r chi.Router) {
		r.With(middleware.IsAuthenticated(ctx, otel, "feedback")).
			Post("/", feedback.Store(fService, otel))
	})

	r.Route("/v1/vote", func(r chi.Router) {
		r.With(middleware.IsAuthenticated(ctx, otel, "vote")).
			Post("/", vote.Store(vService, otel))
	})

	r.Get("/startup", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("startup"))
	})

	r.Get("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("liveness"))
	})

	http.Handle("/", r)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      http.DefaultServeMux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}
}
