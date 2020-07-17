package secret

type Public interface {
	PublicBlob(Blob, error)
}
type Pair interface {
	Public
	Secret
}
