package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateRSAKeys генерирует пару RSA-ключей и сохраняет их в файлы.
func GenerateRSAKeys(bits int) error {
	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("ошибка генерации ключа: %w", err)
	}

	// Создание и запись приватного ключа
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("ошибка создания файла приватного ключа: %w", err)
	}
	defer privateFile.Close()

	err = pem.Encode(privateFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return fmt.Errorf("ошибка сохранения приватного ключа: %w", err)
	}

	// Генерация публичного ключа
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("ошибка кодирования публичного ключа: %w", err)
	}

	// Проверка существования папки перед созданием файла
	publicKeyPath := "public.pem"
	publicDir := filepath.Dir(publicKeyPath)
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		err = os.MkdirAll(publicDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("ошибка создания директории для публичного ключа: %w", err)
		}
	}

	// Создание и запись публичного ключа
	publicFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла публичного ключа: %w", err)
	}
	defer publicFile.Close()

	err = pem.Encode(publicFile, &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
	if err != nil {
		return fmt.Errorf("ошибка сохранения публичного ключа: %w", err)
	}

	fmt.Println("Ключи успешно сгенерированы.")
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
