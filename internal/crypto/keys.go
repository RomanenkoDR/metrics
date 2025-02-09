package crypto

import "crypto/rsa"

// Глобальные переменные для ключей
var PrivateKey *rsa.PrivateKey // Используется сервером
var PublicKey *rsa.PublicKey   // Используется агентом
