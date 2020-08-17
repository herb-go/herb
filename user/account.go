package user

import "strings"

//NewAccount create new account.
func NewAccount() *Account {
	return &Account{}
}

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

//NewAccounts creatre new accounts
func NewAccounts() *Accounts {
	return &Accounts{}
}

//Exists check if an account is in account list.
func (a *Accounts) Exists(account *Account) bool {
	for k := range *a {
		if (*a)[k].Equal(account) {
			return true
		}
	}
	return false
}

//Data return account data
func (a *Accounts) Data() []*Account {
	return []*Account(*a)
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
func (p *PlainAccountProvider) NewAccount(keyword string, account string) (*Account, error) {
	if p.CaseInsensitive {
		account = strings.ToLower(account)
	}
	a := NewAccount()
	a.Keyword = keyword
	a.Account = p.Prefix + account
	return a, nil
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
