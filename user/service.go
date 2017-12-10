package user

type UserService interface {
	Has(uid string) (bool, error)
	CheckAvailable(uid string) (bool, error)
}
type AccountsService interface {
	Accounts(uid string) ([]*UserAccount, error)
}
