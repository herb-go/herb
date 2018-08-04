package util

import (
	"os"
	"path"
	"syscall"
)

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
var SystemPath = path.Join(RootPath, "system")

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
func System(filepaths ...string) string {
	return joinPath(SystemPath, filepaths...)
}

var QuitChan = make(chan int)

func WaitingQuit() {
	<-QuitChan
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
