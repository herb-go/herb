package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func MustListenAndServeHTTP(server *http.Server, config Config, app http.Handler) {
	go func() {
		l := config.MustListen()
		defer l.Close()
		fmt.Println("Listening " + l.Addr().String())
		server.Handler = app
		err := server.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}
func MustServeHTTP(server *http.Server, l net.Listener, app http.Handler) {
	go func() {
		fmt.Println("Listening " + l.Addr().String())
		server.Handler = app
		err := server.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

}
func ShutdownHTTP(Server *http.Server) {
	ShutdownHTTPWithContext(Server, context.Background())

}
func ShutdownHTTPWithTimeout(Server *http.Server, Timeout time.Duration) {
	ctx, fn := context.WithTimeout(context.Background(), Timeout)
	fn()
	ShutdownHTTPWithContext(Server, ctx)

}
func ShutdownHTTPWithContext(Server *http.Server, ctx context.Context) {
	fmt.Println("Quiting...")
	Server.Shutdown(ctx)
	fmt.Println("Quited.")
}
