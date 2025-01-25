package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/acaloiaro/frm"
	"github.com/acaloiaro/frm/routers/frmchi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

var (
	listenAddress = "0.0.0.0:3000"
	requestLogger *httplog.Logger
	logOptions    *slog.HandlerOptions
	logger        *slog.Logger
)

func initLogging() {
	initAppLogger()
	initRequestLogger()
}

func initAppLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	logOptions = &slog.HandlerOptions{Level: slog.LevelDebug}
	if logLevel == "INFO" {
		logOptions.Level = slog.LevelInfo
	}

	logger = slog.New(slog.NewTextHandler(os.Stdout, logOptions))

}

func initRequestLogger() {
	requestLogger = httplog.NewLogger("web", httplog.Options{
		LogLevel:         logOptions.Level.Level().Level(),
		Concise:          true,
		RequestHeaders:   false,
		MessageFieldName: "message",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"env": "dev",
		},
		QuietDownRoutes: []string{
			"/ping",
			"/js/htmx.js",
			"/js/htmx-response-targets.js",
			"/js/hyperscript.js",
			"/js/choices.min.js",
			"/css/styles.css",
			"/css/choices.min.css",
		},
		QuietDownPeriod: 24 * 60 * 60 * time.Second,
	})
}

func main() {
	initLogging()
	logger.Info("frm dev server started")
	chiRouter := chi.NewRouter()
	chiRouter.Use(httplog.RequestLogger(requestLogger))
	chiRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	const workspaceParamName = "frm_workspace_id"
	f, err := frm.New(frm.Args{
		PostgresURL:         os.Getenv("POSTGRES_URL"),
		PostgresDisableSSL:  true,
		WorkspaceIDUrlParam: workspaceParamName, // name of the chi URL parameter name
		BuilderMountPoint:   fmt.Sprintf("/frm/{%s}/build", workspaceParamName),
		CollectorMountPoint: fmt.Sprintf("/frm/{%s}/collect", workspaceParamName),
		Reciever: func(ctx context.Context, submission frm.FormSubmission) (err error) {
			log.Println("GOT SUBMISSION!", submission)
			return
		},
	})
	if err != nil {
		panic(err)
	}
	err = f.Init(context.Background())
	if err != nil {
		panic(err)
	}
	frmchi.Mount(chiRouter, f)

	s := &http.Server{
		Handler:      chiRouter,
		Addr:         listenAddress,
		ReadTimeout:  time.Duration(10 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
	}
	err = s.ListenAndServe()
	slog.Error("server exited", "error", err)
	os.Exit(1)
}
