package user

import "strings"

//Account user account struct
type Account struct {
	//User accont keyword
	Keyword string
	//user account name
	Account string
}

//Equal check if an account is euqal to another.
func (a *Account) Equal(account *Account) bool {
	return a.Keyword == account.Keyword && a.Account == account.Account
}

//Accounts type account list
type Accounts []*Account

//Exists check if an account is in account list.
func (a *Accounts) Exists(account *Account) bool {
	for k := range *a {
		if (*a)[k].Equal(account) {
			return true
		}
	}
	return false
}

//Bind add account to accountlist.
//Return any error if raised.
//If account exists in account list,error ErrAccountBindingExists will be raised.
func (a *Accounts) Bind(account *Account) error {
	for k := range *a {
		if (*a)[k].Equal(account) {
			return ErrAccountBindingExists
		}
	}
	*a = append(*a, account)
	return nil
}

//Unbind remove account from accountlist.
//Return any error if raised.
//If account not exists in account list,error ErrAccountUnbindingNotExists will be raised.
func (a *Accounts) Unbind(account *Account) error {
	for k := range *a {
		if (*a)[k].Equal(account) {
			(*a) = append((*a)[:k], (*a)[k+1:]...)
			return nil
		}
	}
	return ErrAccountUnbindingNotExists
}

//AccountProvider account provider interface
type AccountProvider interface {
	//NewAccount create new account with keyword and account
	NewAccount(keyword string, account string) (*Account, error)
}

//PlainAccountProvider plain account provider.
type PlainAccountProvider struct {
	Prefix          string
	CaseInsensitive bool
}

//NewAccount create new account
//is CaseInsensitive is true,account name will be convert to lower
func (s *PlainAccountProvider) NewAccount(keyword string, account string) (*Account, error) {
	if s.CaseInsensitive {
		account = strings.ToLower(account)
	}
	return &Account{
		Keyword: keyword,
		Account: account,
	}, nil
}

//CaseInsensitiveAcountProvider plain account provider which case insensitive
var CaseInsensitiveAcountProvider = &PlainAccountProvider{
	Prefix:          "",
	CaseInsensitive: true,
}

//CaseSensitiveAcountProvider plain account provider which case sensitive
var CaseSensitiveAcountProvider = &PlainAccountProvider{
	Prefix:          "",
	CaseInsensitive: false,
}
