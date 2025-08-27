package util

import (
	"encoding/json"
	"testing"
)

func TestNewEncryptionManager(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	if em == nil {
		t.Fatal("NewEncryptionManager() returned nil")
	}

	if em.config == nil {
		t.Error("EncryptionManager should have config")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	testData := []byte("Hello, World! This is a test message.")
	shareCode := "blue-whale-42"

	// Encrypt data
	encryptedData, err := em.Encrypt(testData, shareCode)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Validate encrypted data structure
	if err := em.ValidateEncryptedData(encryptedData); err != nil {
		t.Fatalf("ValidateEncryptedData failed: %v", err)
	}

	// Decrypt data
	decryptedData, err := em.Decrypt(encryptedData, shareCode)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Verify decrypted data matches original
	if string(decryptedData) != string(testData) {
		t.Errorf("Decrypted data doesn't match original. Expected: %s, Got: %s",
			string(testData), string(decryptedData))
	}
}

func TestEncryptDecryptEnvFile(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	// Create test environment file
	envFile := &EnvFile{
		Variables: []EnvVariable{
			{
				Key:         "DATABASE_URL",
				Value:       "postgresql://localhost:5432/mydb",
				IsSensitive: true,
				LineNumber:  1,
			},
			{
				Key:         "API_KEY",
				Value:       "sk-1234567890abcdef",
				IsSensitive: true,
				LineNumber:  2,
			},
			{
				Key:         "NODE_ENV",
				Value:       "development",
				IsSensitive: false,
				LineNumber:  3,
			},
		},
		FilePath:   "test.env",
		TotalLines: 3,
		ValidLines: 3,
	}

	shareCode := "red-dragon-123"

	// Encrypt environment file
	encryptedData, err := em.EncryptEnvFile(envFile, shareCode)
	if err != nil {
		t.Fatalf("EncryptEnvFile failed: %v", err)
	}

	// Validate encrypted data
	if err := em.ValidateEncryptedData(encryptedData); err != nil {
		t.Fatalf("ValidateEncryptedData failed: %v", err)
	}

	// Decrypt environment file
	decryptedEnvFile, err := em.DecryptEnvFile(encryptedData, shareCode)
	if err != nil {
		t.Fatalf("DecryptEnvFile failed: %v", err)
	}

	// Verify decrypted environment file matches original
	if decryptedEnvFile.FilePath != envFile.FilePath {
		t.Errorf("FilePath mismatch. Expected: %s, Got: %s",
			envFile.FilePath, decryptedEnvFile.FilePath)
	}

	if len(decryptedEnvFile.Variables) != len(envFile.Variables) {
		t.Errorf("Variables count mismatch. Expected: %d, Got: %d",
			len(envFile.Variables), len(decryptedEnvFile.Variables))
	}

	// Check specific variables
	for i, expectedVar := range envFile.Variables {
		if i >= len(decryptedEnvFile.Variables) {
			t.Errorf("Missing variable at index %d", i)
			continue
		}

		actualVar := decryptedEnvFile.Variables[i]
		if actualVar.Key != expectedVar.Key {
			t.Errorf("Variable key mismatch at index %d. Expected: %s, Got: %s",
				i, expectedVar.Key, actualVar.Key)
		}

		if actualVar.Value != expectedVar.Value {
			t.Errorf("Variable value mismatch at index %d. Expected: %s, Got: %s",
				i, expectedVar.Value, actualVar.Value)
		}

		if actualVar.IsSensitive != expectedVar.IsSensitive {
			t.Errorf("Variable sensitivity mismatch at index %d. Expected: %t, Got: %t",
				i, expectedVar.IsSensitive, actualVar.IsSensitive)
		}
	}
}

func TestEncryptionUniqueness(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	testData := []byte("Same data, different encryption")
	shareCode := "green-forest-999"

	// Encrypt the same data multiple times
	encrypted1, err := em.Encrypt(testData, shareCode)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := em.Encrypt(testData, shareCode)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Encrypted data should be different due to random nonce and salt
	if encrypted1.Data == encrypted2.Data {
		t.Error("Encrypted data should be different for each encryption")
	}

	if encrypted1.Nonce == encrypted2.Nonce {
		t.Error("Nonce should be different for each encryption")
	}

	if encrypted1.Salt == encrypted2.Salt {
		t.Error("Salt should be different for each encryption")
	}

	// But both should decrypt to the same data
	decrypted1, err := em.Decrypt(encrypted1, shareCode)
	if err != nil {
		t.Fatalf("First decryption failed: %v", err)
	}

	decrypted2, err := em.Decrypt(encrypted2, shareCode)
	if err != nil {
		t.Fatalf("Second decryption failed: %v", err)
	}

	if string(decrypted1) != string(decrypted2) {
		t.Error("Both decryptions should produce the same result")
	}
}

func TestWrongShareCode(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	testData := []byte("Secret data")
	correctShareCode := "blue-whale-42"
	wrongShareCode := "red-dragon-123"

	// Encrypt with correct share code
	encryptedData, err := em.Encrypt(testData, correctShareCode)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to decrypt with wrong share code
	_, err = em.Decrypt(encryptedData, wrongShareCode)
	if err == nil {
		t.Error("Decryption with wrong share code should fail")
	}
}

func TestValidateEncryptedData(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	// Test valid encrypted data
	testData := []byte("Test data")
	shareCode := "blue-whale-42"
	encryptedData, err := em.Encrypt(testData, shareCode)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if err := em.ValidateEncryptedData(encryptedData); err != nil {
		t.Errorf("Valid encrypted data should pass validation: %v", err)
	}

	// Test nil encrypted data
	if err := em.ValidateEncryptedData(nil); err == nil {
		t.Error("Nil encrypted data should fail validation")
	}

	// Test empty data
	emptyData := &EncryptedData{
		Data:      "",
		Nonce:     encryptedData.Nonce,
		Salt:      encryptedData.Salt,
		Version:   encryptedData.Version,
		Algorithm: encryptedData.Algorithm,
	}
	if err := em.ValidateEncryptedData(emptyData); err == nil {
		t.Error("Empty data should fail validation")
	}

	// Test empty nonce
	emptyNonce := &EncryptedData{
		Data:      encryptedData.Data,
		Nonce:     "",
		Salt:      encryptedData.Salt,
		Version:   encryptedData.Version,
		Algorithm: encryptedData.Algorithm,
	}
	if err := em.ValidateEncryptedData(emptyNonce); err == nil {
		t.Error("Empty nonce should fail validation")
	}

	// Test empty salt
	emptySalt := &EncryptedData{
		Data:      encryptedData.Data,
		Nonce:     encryptedData.Nonce,
		Salt:      "",
		Version:   encryptedData.Version,
		Algorithm: encryptedData.Algorithm,
	}
	if err := em.ValidateEncryptedData(emptySalt); err == nil {
		t.Error("Empty salt should fail validation")
	}

	// Test unsupported algorithm
	unsupportedAlgo := &EncryptedData{
		Data:      encryptedData.Data,
		Nonce:     encryptedData.Nonce,
		Salt:      encryptedData.Salt,
		Version:   encryptedData.Version,
		Algorithm: "AES-128-CBC",
	}
	if err := em.ValidateEncryptedData(unsupportedAlgo); err == nil {
		t.Error("Unsupported algorithm should fail validation")
	}
}

func TestEncryptedDataSerialization(t *testing.T) {
	config := DefaultConfig()
	em := NewEncryptionManager(config)

	testData := []byte("Test data for serialization")
	shareCode := "blue-whale-42"

	// Encrypt data
	encryptedData, err := em.Encrypt(testData, shareCode)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(encryptedData)
	if err != nil {
		t.Fatalf("JSON marshaling failed: %v", err)
	}

	// Deserialize from JSON
	var deserializedData EncryptedData
	if err := json.Unmarshal(jsonData, &deserializedData); err != nil {
		t.Fatalf("JSON unmarshaling failed: %v", err)
	}

	// Validate deserialized data
	if err := em.ValidateEncryptedData(&deserializedData); err != nil {
		t.Fatalf("Deserialized data validation failed: %v", err)
	}

	// Decrypt deserialized data
	decryptedData, err := em.Decrypt(&deserializedData, shareCode)
	if err != nil {
		t.Fatalf("Decryption of deserialized data failed: %v", err)
	}

	// Verify decrypted data matches original
	if string(decryptedData) != string(testData) {
		t.Errorf("Decrypted deserialized data doesn't match original. Expected: %s, Got: %s",
			string(testData), string(decryptedData))
	}
}
