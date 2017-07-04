package misc

import "net/http"
import "time"

func Time(timezone string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server-Time", time.Now().In(loc).Format("2006-01-02 15:04:05"))
		next(w, r)
	}
}
