package code

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"reflect"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	block, err := aes.NewCipher([]byte("12345678901234567890123456789012")) // AES128=16bytes, AES192=24bytes, AES256=32bytes
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name         string
		input        string
		want         string
		wantEncError bool
		wantDecError bool
	}{
		{
			name:  "OK",
			input: "plain text",
			want:  "plain text",
		},
		{
			name:  "OK (black))",
			input: "",
			want:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			encrypted, err := encryptWithBase64(&block, c.input)
			if c.wantEncError && err == nil {
				t.Fatal("Unexpected no error")
			}
			if !c.wantEncError && err != nil {
				t.Fatalf("Unexpected error occured, err=%+v", err)
			}

			decrypted, err := TestdecryptWithBase64(&block, encrypted)
			if c.wantDecError && err == nil {
				t.Fatal("Unexpected no error")
			}
			if !c.wantDecError && err != nil {
				t.Fatalf("Unexpected error occured, err=%+v", err)
			}

			if !reflect.DeepEqual(c.want, decrypted) {
				t.Fatalf("Unexpected not matching: want=%+v, got=%+v", c.want, decrypted)
			}
		})
	}
}
func TestdecryptWithBase64(block *cipher.Block, encrypted string) (string, error) {
	decoded, err := base64.RawStdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	decrypted := Testdecrypt(block, decoded)
	if len(decrypted) < 1 {
		return "", nil
	}

	// Unpadding
	padSize := int(decrypted[len(decrypted)-1])
	return string(decrypted[:len(decrypted)-padSize]), nil
}

func Testdecrypt(block *cipher.Block, encrypted []byte) []byte {
	if len(encrypted) < aes.BlockSize {
		return []byte("")
	}
	iv := encrypted[:aes.BlockSize] // Get Initial Vector form first head block.
	decrypted := make([]byte, len(encrypted[aes.BlockSize:]))
	decrypter := cipher.NewCBCDecrypter(*block, iv)
	decrypter.CryptBlocks(decrypted, encrypted[aes.BlockSize:])
	return decrypted
}
