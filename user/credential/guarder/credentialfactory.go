package guarder

import (
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

type CredentialFactory interface {
	CreateLoader(r *http.Request) credential.Loader
}
