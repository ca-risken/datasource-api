package code

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func encryptWithBase64(block *cipher.Block, plainText string) (string, error) {
	buf, err := encrypt(block, plainText)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(buf), nil
}

func encrypt(block *cipher.Block, plainText string) ([]byte, error) {
	// PKCS#7 Padding (CBCブロック暗号モードで暗号化したいので、長さが16byteの倍数じゃない場合は末尾をパディングしとく)
	padSize := aes.BlockSize - (len(plainText) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	paddedText := append([]byte(plainText), pad...)

	encrypted := make([]byte, aes.BlockSize+len(paddedText))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return encrypted, err
	}
	encrypter := cipher.NewCBCEncrypter(*block, iv)
	encrypter.CryptBlocks(encrypted[aes.BlockSize:], []byte(paddedText))
	return encrypted, nil
}
