package guarder

import "testing"

func initDrivers() {
	UnregisterAllMapper()
	UnregisterAllIdentifier()
	UnregisterAllCredential()
	registerIDTokenHeadersFactory()
	registerBasicAuthFactory()
	registerTokenMapFactory()
	registerTokenFactory()
}

func TestDriver(t *testing.T) {
	defer initDrivers()
	mf := MapperFactories()
	if len(mf) == 0 {
		t.Fatal(mf)
	}
	idf := IdentifierFactories()
	if len(idf) == 0 {
		t.Fatal(idf)
	}
	cf := CredentialFactories()
	if len(cf) == 0 {
		t.Fatal(cf)
	}

	UnregisterAllMapper()
	UnregisterAllIdentifier()
	UnregisterAllCredential()
	mf = MapperFactories()
	if len(mf) != 0 {
		t.Fatal(mf)
	}
	idf = IdentifierFactories()
	if len(idf) != 0 {
		t.Fatal(idf)
	}
	cf = CredentialFactories()
	if len(cf) != 0 {
		t.Fatal(cf)
	}
}
