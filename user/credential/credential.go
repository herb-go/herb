package credential

type Credential interface {
	Type() CredentialType
	Load() (CredentialData, error)
}

type CredentialType string

var CredentialTypeUID = CredentialType("uid")
var CredentialTypePassword = CredentialType("password")
var CredentialTypeAppID = CredentialType("appid")
var CredentialTypeToken = CredentialType("token")
var CredentialTypeTimestamp = CredentialType("timestamp")
var CredentialTypeSign = CredentialType("sign")
var CredentialTypeSession = CredentialType("session")

type CredentialData []byte

type CredentialDataCollection map[CredentialType]CredentialData

func (d *CredentialDataCollection) Set(t CredentialType, v CredentialData) {
	(*d)[t] = v
}
func (d *CredentialDataCollection) Get(t CredentialType) CredentialData {
	if len((*d)[t]) == 0 {
		return nil
	}
	return (*d)[t]
}

func NewDataCollection() *CredentialDataCollection {
	return &CredentialDataCollection{}
}
