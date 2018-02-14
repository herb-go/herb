package util

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

type SiteConfig struct {
	Name    string
	BaseURL string
}
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
var RootPath = path.Join(path.Dir(SrcPath), "../")
var ResouresPath = path.Join(RootPath, "resources")
var AppDataPath = path.Join(RootPath, "appdata")
var ConfigPath = path.Join(RootPath, "config")

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
func AppData(filepaths ...string) string {
	return joinPath(AppDataPath, filepaths...)
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

func Recover(args ...interface{}) {
	if r := recover(); r != nil {
		err := r.(error)
		if _, ok := LoggerIgnoredErrors[err]; ok == false {
			lines := strings.Split(string(debug.Stack()), "\n")
			length := len(lines)
			maxLength := LoggerMaxLength*2 + 7
			if length > maxLength {
				length = maxLength
			}
			var output = make([]string, length-6)
			output[0] = fmt.Sprintf("Panic: %s", err.Error())
			output[0] += "\n" + lines[0]
			copy(output[1:], lines[7:])
			log.Println(strings.Join(output, "\n"))

		}
	}

}

func Quit() {
	defer func() {
		recover()
	}()
	close(QuitChan)
}

var LoggerMaxLength = 5
var LoggerIgnoredErrors = map[error]bool{
	syscall.EPIPE: true,
}

func RecoverMiddleware(logger *log.Logger) func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if logger == nil {
		logger = log.New(os.Stderr, log.Prefix(), log.Flags())
	}
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				if _, ok := LoggerIgnoredErrors[err]; ok == false {
					lines := strings.Split(string(debug.Stack()), "\n")
					length := len(lines)
					maxLength := LoggerMaxLength*2 + 7
					if length > maxLength {
						length = maxLength
					}
					var output = make([]string, length-6)
					output[0] = fmt.Sprintf("Panic: %s - http request %s \"%s\" ", err.Error(), req.Method, req.URL.String())
					output[0] += "\n" + lines[0]
					copy(output[1:], lines[7:])
					logger.Println(strings.Join(output, "\n"))
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next(w, req)
	}
}
