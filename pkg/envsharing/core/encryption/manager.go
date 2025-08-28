package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/devlink/pkg/envsharing/core"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) EncryptEnvFile(envFile *core.EnvFile, shareCode string) (*core.EncryptedData, error) {
	data := []byte(envFile.RawContent)
	return m.Encrypt(data, shareCode)
}

func (m *Manager) DecryptEnvFile(encryptedData *core.EncryptedData, shareCode string) (*core.EnvFile, error) {
	data, err := m.Decrypt(encryptedData, shareCode)
	if err != nil {
		return nil, err
	}

	parser := core.NewParser()
	envFile := parser.ParseContent(string(data), "retrieved")
	return envFile, nil
}

func (m *Manager) Encrypt(data []byte, shareCode string) (*core.EncryptedData, error) {
	key, salt, err := m.deriveKey(shareCode)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)

	return &core.EncryptedData{
		Data:      base64.StdEncoding.EncodeToString(encrypted),
		Nonce:     base64.StdEncoding.EncodeToString(nonce),
		Salt:      base64.StdEncoding.EncodeToString(salt),
		Version:   "1.0",
		Algorithm: "AES-256-GCM",
	}, nil
}

func (m *Manager) Decrypt(encryptedData *core.EncryptedData, shareCode string) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData.Data)
	if err != nil {
		return nil, err
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedData.Nonce)
	if err != nil {
		return nil, err
	}

	salt, err := base64.StdEncoding.DecodeString(encryptedData.Salt)
	if err != nil {
		return nil, err
	}

	key, _, err := m.deriveKeyWithSalt(shareCode, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	ciphertext := encrypted[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (m *Manager) deriveKey(shareCode string) ([]byte, []byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	key, _, err := m.deriveKeyWithSalt(shareCode, salt)
	return key, salt, err
}

func (m *Manager) deriveKeyWithSalt(shareCode string, salt []byte) ([]byte, []byte, error) {
	hash := sha256.New()
	hash.Write([]byte(shareCode))
	hash.Write(salt)
	key := hash.Sum(nil)
	return key, salt, nil
}
