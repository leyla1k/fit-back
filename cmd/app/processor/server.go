package processor

import (
	"context"
	"net/http"
)

var stop func(ctx context.Context) error

func Run() error {
	handler := api()

	server := http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	stop = server.Shutdown

	return server.ListenAndServe()
}

func Shutdown(ctx context.Context) {
	stop(ctx)
}
