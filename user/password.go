package user

//PasswordVerifier password verifier interface
type PasswordVerifier interface {
	//VerifyPassword Verify user password.
	//Return verify result and any error if raised
	VerifyPassword(uid string, password string) (bool, error)
}
