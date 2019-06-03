package user

import (
	"testing"
)

func TestUserAccount(t *testing.T) {
	account1, err := CaseSensitiveAcountProvider.NewAccount("testaccount", "account1")
	if err != nil {
		t.Fatal(err)
	}
	account1equal, err := CaseInsensitiveAcountProvider.NewAccount("testaccount", "account1")
	if err != nil {
		t.Fatal(err)
	}
	account2, err := CaseInsensitiveAcountProvider.NewAccount("testaccount", "account2")
	if err != nil {
		t.Fatal(err)
	}
	if !account1.Equal(account1equal) {
		t.Error(account1)
	}
	if account1.Equal(account2) {
		t.Error(account1)
	}
	accounts := NewAccounts()
	accounts.Bind(account1)
	if !accounts.Exists(account1equal) {
		t.Error(account1equal)
	}
	if accounts.Exists(account2) {
		t.Error(account2)
	}
	err = accounts.Bind(account1equal)
	if err != ErrAccountBindingExists {
		t.Fatal(err)
	}
	err = accounts.Unbind(account2)
	if err != ErrAccountUnbindingNotExists {
		t.Fatal(err)
	}
	err = accounts.Bind(account2)
	if err != nil {
		t.Fatal(err)
	}
	if !accounts.Exists(account2) {
		t.Error(account2)
	}
	err = accounts.Unbind(account2)
	if err != nil {
		t.Fatal(err)
	}
	if accounts.Exists(account2) {
		t.Error(account2)
	}
}

func TestCIAccountType(t *testing.T) {
	account1cs, err := CaseSensitiveAcountProvider.NewAccount("testaccount", "account1")
	if err != nil {
		t.Fatal(err)
	}
	Account1cs, err := CaseSensitiveAcountProvider.NewAccount("testaccount", "Account1")
	if err != nil {
		t.Fatal(err)
	}
	if account1cs.Equal(Account1cs) {
		t.Error(Account1cs)
	}
	account1ci, err := CaseInsensitiveAcountProvider.NewAccount("testaccount", "account1")
	if err != nil {
		t.Fatal(err)
	}
	Account1ci, err := CaseInsensitiveAcountProvider.NewAccount("testaccount", "Account1")
	if err != nil {
		t.Fatal(err)
	}
	if !account1ci.Equal(Account1ci) {
		t.Error(Account1cs)
	}
}
