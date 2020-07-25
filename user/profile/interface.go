package profile

type Provider interface {
	GetProfile(id string) (*Profile, error)
}

type Loader interface {
	LoadProfile() (*Profile, error)
}
