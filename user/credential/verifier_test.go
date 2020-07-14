package credential

import "testing"

var testVerifier = VerifierFunc(
	func(c *Collection) (string, error) {
		if string(c.Get(TypeAppID)) == "testappid" && string(c.Get(TypeToken)) == "testtoken" {
			return "testappid", nil
		}
		return "", nil
	},
	TypeAppID,
	TypeToken,
)

func TestVerifier(t *testing.T) {
	appid := New().WithType(TypeAppID).WithData([]byte("testappid"))
	token := New().WithType(TypeToken).WithData([]byte("testtoken"))
	id, err := Verify(testVerifier, appid)
	if id != "" || err != nil {
		t.Fatal(id, err)
	}
	id, err = Verify(testVerifier, appid, token)
	if id != "testappid" || err != nil {
		t.Fatal(id, err)
	}
	id, err = Verify(ForbiddenVerifier, appid, token)
	if id != "" || err != nil {
		t.Fatal(id, err)
	}

}
