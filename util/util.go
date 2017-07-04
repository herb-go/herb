package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

type NetConfig struct {
	Net  string
	Addr string
}

func (c *NetConfig) Listen() (net.Listener, error) {
	return net.Listen(c.Net, c.Addr)
}
func (c *NetConfig) MustListen() net.Listener {
	l, err := net.Listen(c.Net, c.Addr)
	Must(err)
	return l
}
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
func mustPath(path string, err error) string {
	if err != nil {
		panic(err)
	}
	return path
}

var SrcPath = mustPath(os.Executable())
var Rootpath = path.Join(path.Dir(SrcPath), "../")
var ResouresPath = path.Join(Rootpath, "resources")
var AppdataPath = path.Join(Rootpath, "appdata")
var ConfigPath = path.Join(Rootpath, "config")

func SetConfigPath(paths ...string) {
	ConfigPath = path.Join(paths...)
}
func MustGetWD() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

func joinPath(p string, filepath ...string) string {
	return path.Join(p, path.Join(filepath...))
}
func Resource(filepaths ...string) string {
	return joinPath(ResouresPath, filepaths...)
}
func Config(filepaths ...string) string {
	return joinPath(ConfigPath, filepaths...)
}
func Appdata(filepaths ...string) string {
	return joinPath(AppdataPath, filepaths...)
}

var QuitChan = make(chan int)

func WaitingQuit() {
	<-QuitChan
}

func MustListenAndServeHTTP(server *http.Server, netconfig NetConfig, app http.Handler) {
	go func() {
		l := netconfig.MustListen()
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
	ctx, _ := context.WithTimeout(context.Background(), Timeout)
	ShutdownHTTPWithContext(Server, ctx)

}
func ShutdownHTTPWithContext(Server *http.Server, ctx context.Context) {
	fmt.Println("Quiting...")
	Server.Shutdown(ctx)
	fmt.Println("Quited.")
}
func Quit() {
	defer func() {
		recover()
	}()
	close(QuitChan)
}
