package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/brave-intl/bat-go/libs/logging"
	"github.com/brave-intl/bat-go/libs/middleware"
	sentry "github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	chiware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/brave/go-translate/controller"
)

func contextHandler(ctx context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Copy over default net/http server context keys
		if v, ok := ctx.Value(http.ServerContextKey).(*http.Server); ok {
			ctx = context.WithValue(ctx, http.ServerContextKey, v)
		}
		if v, ok := ctx.Value(http.LocalAddrContextKey).(net.Addr); ok {
			ctx = context.WithValue(ctx, http.LocalAddrContextKey, v)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setupRouter(ctx context.Context, logger *zerolog.Logger) (context.Context, *chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(chiware.RequestID)
	r.Use(chiware.RealIP)
	r.Use(chiware.Heartbeat("/"))
	r.Use(chiware.Timeout(60 * time.Second))
	r.Use(middleware.BearerToken)

	if logger != nil {
		// Also handles panic recovery
		r.Use(
			hlog.NewHandler(*logger),
			hlog.UserAgentHandler("user_agent"),
			hlog.RequestIDHandler("req_id", "Request-Id"),
			middleware.RequestLogger(logger))
	}
	r.Get("/metrics", middleware.Metrics())
	tr, err := controller.TranslateRouter(ctx)
	r.Mount("/", tr)

	return ctx, r, err
}

// StartServer starts the translate proxy server on port 8195
func StartServer() {
	serverCtx, logger := logging.SetupLogger(context.Background())

	serverCtx, r, err := setupRouter(serverCtx, logger)
	if err != nil {
		logger.Panic().Err(err).Msg("service setup failed!")
	}
	port := ":8195"

	go func() {
		logger.Info().
			Str("port", ":9090").
			Msg("Starting metrics server")

		err := http.ListenAndServe(":9090", middleware.Metrics())
		if err != nil {
			sentry.CaptureException(err)
			logger.Panic().Err(err).Msg("metrics HTTP server start failed!")
		}
	}()

	logger.Info().
		Str("port", port).
		Msg("Starting API server")

	srv := http.Server{Addr: port, Handler: contextHandler(serverCtx, r)}
	err = srv.ListenAndServe()
	if err != nil {
		sentry.CaptureException(err)
		logger.Panic().Err(err)
	}
}
