package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/eminetto/api-o11y-gcp/auth"
	"github.com/eminetto/api-o11y-gcp/internal/telemetry"
	"github.com/eminetto/api-o11y-gcp/user"
	mysql_user "github.com/eminetto/api-o11y-gcp/user/mysql"
	"github.com/go-chi/chi"
	"github.com/go-chi/httplog"
	telemetrymiddleware "github.com/go-chi/telemetry"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", dataSourceName)
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

	repo := mysql_user.NewUserMySQL(db, otel)
	uService := user.NewService(repo, otel)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(telemetrymiddleware.Collector(telemetrymiddleware.Config{
		AllowAny: true,
	}, []string{"/v1"})) // path prefix filters basically records generic http request metrics
	r.Post("/v1/auth", auth.UserAuth(ctx, uService, otel))

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
