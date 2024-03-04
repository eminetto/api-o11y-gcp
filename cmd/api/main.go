package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

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
	"github.com/go-chi/httplog"
	telemetrymiddleware "github.com/go-chi/telemetry"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Logger
	//@todo remover logs de requests pq isso deve ser só métrica
	//@todo adicionar o slog para log estruturado de erros
	logger := httplog.NewLogger("auth", httplog.Options{
		JSON: true,
	})

	err := godotenv.Load()
	if err != nil {
		logger.Panic().Msg(err.Error())
	}

	db, err := sql.Open("sqlite3", "ops/db/api.db")
	if err != nil {
		logger.Panic().Msg(err.Error())
	}
	defer db.Close()

	ctx := context.Background()
	otel, err := telemetry.NewJaeger(ctx, "auth")
	if err != nil {
		logger.Panic().Msg(err.Error())
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
		r.With(middleware.IsAuthenticated(ctx, otel)).
			Post("/", feedback.Store(ctx, fService, otel))
	})

	r.Route("/v1/vote", func(r chi.Router) {
		r.With(middleware.IsAuthenticated(ctx, otel)).
			Post("/", vote.Store(ctx, vService, otel))
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
		logger.Panic().Msg(err.Error())
	}
}
