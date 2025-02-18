package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"go.uber.org/zap"
	"os"
)

const (
	PrivateKeyPath = "private.pem"
	PublicKeyPath  = "public.pem"
)

// GenerateAESKey создает новый 32-байтовый AES-ключ
func GenerateAESKey() ([]byte, error) {

	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		logger.Error("Ошибка генерации AES-ключа", zap.Error(err))
		return nil, err
	}
	logger.Info("Сгенерирован новый AES-ключ")
	return aesKey, nil
}

// GenerateKeys проверяет существование ключей и генерирует новые, если их нет
func GenerateKeys() error {
	if _, err := os.Stat(PrivateKeyPath); err == nil {
		if _, err := os.Stat(PublicKeyPath); err == nil {
			logger.Info("Ключи уже существуют, генерация не требуется")
			return nil
		}
	}

	logger.Info("Генерация новой пары RSA-ключей")

	// Генерация приватного ключа (4096 бит)
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.Error("Ошибка генерации ключей", zap.Error(err))
		return err
	}

	// Сохранение приватного ключа
	if err := savePrivateKey(privKey); err != nil {
		return err
	}

	// Сохранение публичного ключа
	if err := savePublicKey(&privKey.PublicKey); err != nil {
		return err
	}

	logger.Info("Ключи успешно сгенерированы и сохранены")
	return nil
}

// savePrivateKey сохраняет приватный ключ в файл
func savePrivateKey(privKey *rsa.PrivateKey) error {
	privKeyFile, err := os.Create(PrivateKeyPath)
	if err != nil {
		logger.Error("Ошибка создания файла приватного ключа", zap.Error(err))
		return err
	}
	defer privKeyFile.Close()

	privKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	if _, err := privKeyFile.Write(privKeyPEM); err != nil {
		logger.Error("Ошибка записи приватного ключа", zap.Error(err))
		return err
	}

	return nil
}

// savePublicKey сохраняет публичный ключ в файл
func savePublicKey(pubKey *rsa.PublicKey) error {
	pubKeyPEM, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		logger.Error("Ошибка кодирования публичного ключа", zap.Error(err))
		return err
	}

	pubKeyFile, err := os.Create(PublicKeyPath)
	if err != nil {
		logger.Error("Ошибка создания файла публичного ключа", zap.Error(err))
		return err
	}
	defer pubKeyFile.Close()

	if _, err := pubKeyFile.Write(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyPEM})); err != nil {
		logger.Error("Ошибка записи публичного ключа", zap.Error(err))
		return err
	}

	return nil
}
