package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/RomanenkoDR/metrics/internal/middleware/logger"
)

// DecryptAESKeyRSA расшифровывает AES-ключ с помощью RSA
func DecryptAESKeyRSA(encAESKey []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	logger.Info("Расшифровка AES-ключа с помощью RSA")
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, privKey, encAESKey, nil)
}
