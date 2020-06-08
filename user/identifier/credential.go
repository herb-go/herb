package identifier

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

type CredentialDataCollection map[CredentialType][]CredentialData

func (d *CredentialDataCollection) Append(t CredentialType, v CredentialData) {
	(*d)[t] = append((*d)[t], v)
}
func (d *CredentialDataCollection) Get(t CredentialType) CredentialData {
	if len((*d)[t]) == 0 {
		return nil
	}
	return (*d)[t][0]
}
func (d *CredentialDataCollection) GetAllByType(t CredentialType) []CredentialData {
	return (*d)[t]
}
func NewDataCollection() *CredentialDataCollection {
	return &CredentialDataCollection{}
}
