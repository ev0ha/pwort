package encryption

import "golang.org/x/crypto/bcrypt"

func HashPwd(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	return string(bytes), err
}

func CheckPwdHash(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}
