package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// EncryptData шифрует данные с помощью публичного ключа
func EncryptData(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования данных: %v", err)
	}
	return encryptedData, nil
}
