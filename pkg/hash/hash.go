package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Gunakan bcrypt dengan cost 10 (standar production)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
	return err == nil
}
