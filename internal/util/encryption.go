package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptedData represents encrypted environment data
type EncryptedData struct {
	Data      string `json:"data"`      // Base64 encoded encrypted data
	Nonce     string `json:"nonce"`     // Base64 encoded nonce
	Salt      string `json:"salt"`      // Base64 encoded salt for key derivation
	Version   string `json:"version"`   // Encryption version for future compatibility
	Algorithm string `json:"algorithm"` // Encryption algorithm used
}

// EncryptionManager handles encryption and decryption operations
type EncryptionManager struct {
	config *Config
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager(config *Config) *EncryptionManager {
	return &EncryptionManager{
		config: config,
	}
}

// Encrypt encrypts environment data using AES-256-GCM
func (em *EncryptionManager) Encrypt(data []byte, shareCode string) (*EncryptedData, error) {
	// Derive encryption key from share code
	key, salt, err := em.deriveKey(shareCode)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Create encrypted data structure
	encryptedData := &EncryptedData{
		Data:      base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:     base64.StdEncoding.EncodeToString(nonce),
		Salt:      base64.StdEncoding.EncodeToString(salt),
		Version:   "1.0",
		Algorithm: "AES-256-GCM",
	}

	return encryptedData, nil
}

// Decrypt decrypts environment data using AES-256-GCM
func (em *EncryptionManager) Decrypt(encryptedData *EncryptedData, shareCode string) ([]byte, error) {
	// Validate algorithm
	if encryptedData.Algorithm != "AES-256-GCM" {
		return nil, fmt.Errorf("unsupported algorithm: %s", encryptedData.Algorithm)
	}

	// Decode base64 data
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedData.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	salt, err := base64.StdEncoding.DecodeString(encryptedData.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Derive key using the same salt
	key, _, err := em.deriveKeyWithSalt(shareCode, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// deriveKey derives a 32-byte key from share code using PBKDF2
func (em *EncryptionManager) deriveKey(shareCode string) ([]byte, []byte, error) {
	// Generate random salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return em.deriveKeyWithSalt(shareCode, salt)
}

// deriveKeyWithSalt derives a key using the provided salt
func (em *EncryptionManager) deriveKeyWithSalt(shareCode string, salt []byte) ([]byte, []byte, error) {
	// Use PBKDF2 with SHA256 for key derivation
	// 100,000 iterations for security
	key := pbkdf2.Key([]byte(shareCode), salt, 100000, 32, sha256.New)
	return key, salt, nil
}

// EncryptEnvFile encrypts an environment file
func (em *EncryptionManager) EncryptEnvFile(envFile *EnvFile, shareCode string) (*EncryptedData, error) {
	// Serialize environment file to JSON
	jsonData, err := json.Marshal(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize env file: %w", err)
	}

	// Encrypt the JSON data
	return em.Encrypt(jsonData, shareCode)
}

// DecryptEnvFile decrypts an environment file
func (em *EncryptionManager) DecryptEnvFile(encryptedData *EncryptedData, shareCode string) (*EnvFile, error) {
	// Decrypt the data
	jsonData, err := em.Decrypt(encryptedData, shareCode)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	// Deserialize to EnvFile
	var envFile EnvFile
	if err := json.Unmarshal(jsonData, &envFile); err != nil {
		return nil, fmt.Errorf("failed to deserialize env file: %w", err)
	}

	return &envFile, nil
}

// ValidateEncryptedData validates the structure of encrypted data
func (em *EncryptionManager) ValidateEncryptedData(encryptedData *EncryptedData) error {
	if encryptedData == nil {
		return fmt.Errorf("encrypted data is nil")
	}

	if encryptedData.Data == "" {
		return fmt.Errorf("encrypted data is empty")
	}

	if encryptedData.Nonce == "" {
		return fmt.Errorf("nonce is empty")
	}

	if encryptedData.Salt == "" {
		return fmt.Errorf("salt is empty")
	}

	if encryptedData.Algorithm != "AES-256-GCM" {
		return fmt.Errorf("unsupported algorithm: %s", encryptedData.Algorithm)
	}

	// Validate base64 encoding
	if _, err := base64.StdEncoding.DecodeString(encryptedData.Data); err != nil {
		return fmt.Errorf("invalid data encoding: %w", err)
	}

	if _, err := base64.StdEncoding.DecodeString(encryptedData.Nonce); err != nil {
		return fmt.Errorf("invalid nonce encoding: %w", err)
	}

	if _, err := base64.StdEncoding.DecodeString(encryptedData.Salt); err != nil {
		return fmt.Errorf("invalid salt encoding: %w", err)
	}

	return nil
}
