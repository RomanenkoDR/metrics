package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// DecryptData расшифровывает данные с помощью приватного ключа
func DecryptData(encryptedData []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("ошибка расшифровки данных: %v", err)
	}
	return decryptedData, nil
}
