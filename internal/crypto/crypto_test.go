package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes for AES-256
	iv := make([]byte, 16)
	plaintext := []byte("Hello, MegaBasterd!")

	encrypted, err := EncryptAES(plaintext, key, iv)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := DecryptAES(encrypted, key, iv)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match. Got %s, want %s", decrypted, plaintext)
	}
}

func TestAESCTREncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	iv := make([]byte, 16)
	plaintext := []byte("Hello, MegaBasterd CTR mode!")

	encrypted, err := EncryptAESCTR(plaintext, key, iv)
	if err != nil {
		t.Fatalf("CTR encryption failed: %v", err)
	}

	decrypted, err := DecryptAESCTR(encrypted, key, iv)
	if err != nil {
		t.Fatalf("CTR decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("CTR decrypted text doesn't match. Got %s, want %s", decrypted, plaintext)
	}
}

func TestDeriveKey(t *testing.T) {
	password := "testpassword"
	salt := []byte("somesalt12345678")
	iterations := 10000
	keyLen := 32

	key := DeriveKey(password, salt, iterations, keyLen)

	if len(key) != keyLen {
		t.Errorf("Derived key length is %d, want %d", len(key), keyLen)
	}

	// Test consistency
	key2 := DeriveKey(password, salt, iterations, keyLen)
	if !bytes.Equal(key, key2) {
		t.Error("Derived keys are not consistent")
	}
}

func TestGenerateSalt(t *testing.T) {
	length := 16
	salt, err := GenerateSalt(length)
	if err != nil {
		t.Fatalf("Salt generation failed: %v", err)
	}

	if len(salt) != length {
		t.Errorf("Salt length is %d, want %d", len(salt), length)
	}

	// Test uniqueness
	salt2, err := GenerateSalt(length)
	if err != nil {
		t.Fatalf("Second salt generation failed: %v", err)
	}

	if bytes.Equal(salt, salt2) {
		t.Error("Generated salts are identical (should be random)")
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	data := []byte("Test data for base64")

	encoded := Base64Encode(data)
	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Fatalf("Base64 decode failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match. Got %s, want %s", decoded, data)
	}
}

func TestRSAEncryptDecrypt(t *testing.T) {
	// Generate RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("RSA key generation failed: %v", err)
	}

	plaintext := []byte("RSA test message")

	encrypted, err := EncryptRSA(plaintext, &privateKey.PublicKey)
	if err != nil {
		t.Fatalf("RSA encryption failed: %v", err)
	}

	decrypted, err := DecryptRSA(encrypted, privateKey)
	if err != nil {
		t.Fatalf("RSA decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("RSA decrypted text doesn't match. Got %s, want %s", decrypted, plaintext)
	}
}

func TestPKCS7Padding(t *testing.T) {
	data := []byte("test")
	blockSize := 16

	padded := pkcs7Pad(data, blockSize)
	if len(padded)%blockSize != 0 {
		t.Errorf("Padded data length %d is not multiple of block size %d", len(padded), blockSize)
	}

	unpadded, err := pkcs7Unpad(padded, blockSize)
	if err != nil {
		t.Fatalf("Unpadding failed: %v", err)
	}

	if !bytes.Equal(data, unpadded) {
		t.Errorf("Unpadded data doesn't match. Got %s, want %s", unpadded, data)
	}
}
