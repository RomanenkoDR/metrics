package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKey загружает приватный ключ из файла
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Читаем файл с приватным ключом
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения приватного ключа: %v", err)
	}

	// Декодируем PEM-блок
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("неверный формат приватного ключа")
	}

	// Парсим RSA-ключ
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга приватного ключа: %v", err)
	}

	return privateKey, nil
}

// LoadPublicKey загружает публичный ключ из файла
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения публичного ключа: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("неверный формат публичного ключа")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга публичного ключа: %v", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("неверный тип публичного ключа")
	}

	return publicKey, nil
}
