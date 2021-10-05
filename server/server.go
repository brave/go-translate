package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brave-intl/bat-go/middleware"
	"github.com/brave/go-translate/controller"
	"github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
)

func setupLogger(ctx context.Context) (context.Context, *logrus.Logger) {
	logger := logrus.New()
	// Redirect output from the standard logging package "log"
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	ctx = lg.WithLoggerContext(ctx, logger)
	return ctx, logger
}

func setupRouter(ctx context.Context, logger *logrus.Logger) (context.Context, *chi.Mux) {
	r := chi.NewRouter()

	r.Use(chiware.RequestID)
	r.Use(chiware.RealIP)
	r.Use(chiware.Heartbeat("/"))
	r.Use(chiware.Timeout(60 * time.Second))
	r.Use(middleware.BearerToken)

	if logger != nil {
		// Also handles panic recovery
		r.Use(middleware.RequestLogger(logger))
	}

	r.Mount("/", controller.TranslateRouter())

	return ctx, r
}

// StartServer starts the translate proxy server on port 8195
func StartServer() {
	serverCtx, logger := setupLogger(context.Background())
	logger.WithFields(logrus.Fields{"prefix": "main"}).Info("Starting server")

	go func() {
		// setup metrics on another non-public port 9090
		err := http.ListenAndServe(":9090", middleware.Metrics())
		if err != nil {
			sentry.CaptureException(err)
			panic(fmt.Sprintf("metrics HTTP server start failed: %s", err.Error()))
		}
	}()

	serverCtx, r := setupRouter(serverCtx, logger)
	port := ":8195"
	fmt.Println("Starting server: http://localhost" + port)
	srv := http.Server{Addr: port, Handler: chi.ServerBaseContext(serverCtx, r)}
	err := srv.ListenAndServe()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}
}
