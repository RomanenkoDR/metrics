package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"go.uber.org/zap"
	"os"
)

// LoadPublicKey загружает публичный RSA-ключ из файла
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	logger.Info("Загрузка публичного ключа", zap.String("path", path))

	keyData, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Ошибка чтения публичного ключа", zap.Error(err))
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		logger.Error("Неверный формат публичного ключа")
		return nil, errors.New("неверный формат публичного ключа")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Error("Ошибка парсинга публичного ключа", zap.Error(err))
		return nil, err
	}

	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		logger.Error("Неверный тип публичного ключа")
		return nil, errors.New("неверный тип публичного ключа")
	}

	return pubKey, nil
}

// LoadPrivateKey загружает приватный RSA-ключ из файла
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	logger.Info("Загрузка приватного ключа", zap.String("path", path))

	keyData, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Ошибка чтения приватного ключа", zap.Error(err))
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.Error("Неверный формат приватного ключа")
		return nil, errors.New("неверный формат приватного ключа")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Error("Ошибка парсинга приватного ключа", zap.Error(err))
		return nil, err
	}

	return privKey, nil
}
