package user

import (
	"testing"
)

func TestUserAccount(t *testing.T) {
	account1, err := CaseInsensitiveAcountType.NewAccount("testaccount", "account1")
	if err != nil {
		t.Fatal(err)
	}
	account1equal, err := CaseInsensitiveAcountType.NewAccount("testaccount", "account1")
	if err != nil {
		t.Fatal(err)
	}
	account2, err := CaseInsensitiveAcountType.NewAccount("testaccount", "account2")
	if err != nil {
		t.Fatal(err)
	}
	if !account1.Equal(account1equal) {
		t.Error(account1)
	}
	if account1.Equal(account2) {
		t.Error(account1)
	}
	accounts := UserAccounts{account1}
	if !accounts.Exists(account1equal) {
		t.Error(account1equal)
	}
	if accounts.Exists(account2) {
		t.Error(account2)
	}
	err = accounts.Bind(account1equal)
	if err != ErrAccountBindExists {
		t.Fatal(err)
	}
	err = accounts.Unbind(account2)
	if err != ErrAccountUnbindNotExists {
		t.Fatal(err)
	}
	err = accounts.Bind(account2)
	if err != nil {
		t.Fatal(err)
	}
	if !accounts.Exists(account2) {
		t.Error(account2)
	}
}
