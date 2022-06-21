package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	r.Get("/metrics", middleware.Metrics())

	return ctx, r
}

// StartServer starts the translate proxy server on port 8195
func StartServer() {
	serverCtx, logger := setupLogger(context.Background())
	logger.WithFields(logrus.Fields{"prefix": "main"}).Info("Starting server")
	serverCtx, r := setupRouter(serverCtx, logger)
	port := ":8195"
	fmt.Printf("Starting server: http://localhost%s", port)
	
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	FileServer(r, "/static", filesDir)

	srv := http.Server{Addr: port, Handler: chi.ServerBaseContext(serverCtx, r)}
	err := srv.ListenAndServe()
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))

		w.Header().Set("Access-Control-Allow-Origin", "*")

		fs.ServeHTTP(w, r)
	})
}