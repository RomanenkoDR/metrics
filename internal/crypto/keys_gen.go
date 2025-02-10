package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateRSAKeys генерирует пару RSA-ключей и сохраняет их в файлы.
func GenerateRSAKeys(keySize int) error {
	privateKeyPath := "private.pem"
	publicKeyPath := "../agent/public.pem"

	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("ошибка генерации приватного ключа: %v", err)
	}

	// Сохранение приватного ключа
	err = savePrivateKey(privateKey, privateKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка сохранения приватного ключа: %v", err)
	}

	// Сохранение публичного ключа
	err = savePublicKey(&privateKey.PublicKey, publicKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка сохранения публичного ключа: %v", err)
	}

	fmt.Println("RSA-ключи успешно сгенерированы:")
	fmt.Println("Приватный ключ:", privateKeyPath)
	fmt.Println("Публичный ключ:", publicKeyPath)

	return nil
}

// savePrivateKey сохраняет приватный ключ в PEM-формате.
func savePrivateKey(privateKey *rsa.PrivateKey, path string) error {
	privateKeyFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка создания файла приватного ключа: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.Encode(privateKeyFile, privateKeyPEM)
}

// savePublicKey сохраняет публичный ключ в PEM-формате.
func savePublicKey(publicKey *rsa.PublicKey, path string) error {
	publicKeyFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка создания файла публичного ключа: %v", err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга публичного ключа: %v", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	return pem.Encode(publicKeyFile, publicKeyPEM)
}
