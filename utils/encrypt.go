package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

/**
对一个字符串进行加密
*/
func Encrypt(content, key string) ([]byte, error) {
	secret := []byte(key)
	text := []byte(content)
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	cipherext := make([]byte, aes.BlockSize+len(b))
	iv := cipherext[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherext[aes.BlockSize:], []byte(b))
	return cipherext, nil
}

/**
对一个字符串进行解密
*/
func Decrypt(content, key string) ([]byte, error) {
	secret := []byte(key)
	text := []byte(content)
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

//MD5 生成32位md5值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// md5 add salt
func MD5BySalt(str string) string {
	salt := time.Now().Unix()
	m5 := md5.New()
	m5.Write([]byte(str))
	m5.Write([]byte(string(salt)))
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}


// Encrypt string to base64 crypto using AES
func AESEncrypt(secret string, text string) (string, bool) {
	key:= []byte(secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", false
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext), true
}

// Decrypt from base64 to decrypted string
func AESDecrypt(secret string, cryptoText string) (string, bool) {
	key:= []byte(secret)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}
	if len(ciphertext) < aes.BlockSize {
		return "", false
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext), true
}
