package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword recebe uma senha em texto plano e retorna sua versão criptografada usando bcrypt
// A função utiliza o DefaultCost do bcrypt para a complexidade da criptografia
// Retorna a senha criptografada como string e qualquer erro que ocorra durante o processo
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compara uma senha em texto plano com uma senha criptografada
// Retorna true se a senha corresponder ao hash, false caso contrário
// Esta função é comumente usada para verificação de senha durante o login
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}