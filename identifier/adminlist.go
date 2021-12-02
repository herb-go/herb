package identifier

import "net/http"

type IDListVerifer []string

func (v IDListVerifer) VerifyID(id string) (bool, error) {
	for _, value := range v {
		if value == id {
			return true, nil
		}
	}
	return false, nil
}

var ForbiddenHanlder = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(403), 403)
})
