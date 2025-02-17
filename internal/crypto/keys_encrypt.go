package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
	"go.uber.org/zap"
	"io"
)

// EncryptAESKeyRSA шифрует AES-ключ с помощью RSA
func EncryptAESKeyRSA(aesKey []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	logger.Info("Шифрование AES-ключа с помощью RSA")
	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, pubKey, aesKey, nil)
}

// EncryptData шифрует данные AES-ключом
func EncryptData(plainText []byte, aesKey []byte) ([]byte, error) {
	logger.Info("Шифрование данных с помощью AES")
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		logger.Error("Ошибка создания AES-шифра", zap.Error(err))
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logger.Error("Ошибка генерации IV", zap.Error(err))
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	logger.Info("Данные успешно зашифрованы")
	return cipherText, nil
}

// DecryptData расшифровывает данные AES-ключом
func DecryptData(cipherText []byte, aesKey []byte) ([]byte, error) {
	logger.Info("Расшифровка данных с помощью AES")
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		logger.Error("Ошибка создания AES-шифра", zap.Error(err))
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		logger.Error("Шифрованный текст слишком короткий")
		return nil, errors.New("шифрованный текст слишком короткий")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	logger.Info("Данные успешно расшифрованы")
	return cipherText, nil
}
