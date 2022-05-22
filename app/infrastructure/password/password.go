package password

import "golang.org/x/crypto/bcrypt"

// Verify compares encrypted password with plain password and returns `true` if encrypted plain equal to encrypted hash.
func Verify(encryptedHash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedHash), []byte(plain)) == nil
}

// Encrypt encrypts plain text by bcrypt.
func Encrypt(plain string) (encryptedHash string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	return string(b), err
}
