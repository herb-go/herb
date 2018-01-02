package user

type UserAccount struct {
	Keyword string
	Account string
}

func (a *UserAccount) Equal(account *UserAccount) bool {
	return a.Keyword == account.Keyword && a.Account == account.Account
}

type UserAccounts []UserAccount

func (a *UserAccounts) Exists(account *UserAccount) bool {
	for k := range *a {
		if (*a)[k].Equal(account) {
			return true
		}
	}
	return false
}
func (a *UserAccounts) Bind(account *UserAccount) error {
	for k := range *a {
		if (*a)[k].Equal(account) {
			return ErrAccountBindExists
		}
	}
	*a = append(*a, *account)
	return nil
}

func (a *UserAccounts) Unbind(account *UserAccount) error {
	for k := range *a {
		if (*a)[k].Equal(account) {
			(*a) = append((*a)[:k], (*a)[k+1:]...)
			return nil
		}
	}
	return ErrAccountUnbindNotExists
}
