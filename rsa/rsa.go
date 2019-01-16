package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// GenRsaKey 生成公钥私钥
func GenRsaKey(bits int) (string, string, error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privateBuffer := new(bytes.Buffer)
	err = pem.Encode(privateBuffer, block)
	if err != nil {
		return "", "", err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	publicBuffer := new(bytes.Buffer)
	err = pem.Encode(publicBuffer, block)
	if err != nil {
		return "", "", err
	}
	return publicBuffer.String(), privateBuffer.String(), nil
}

// Encrypt 加密
func Encrypt(publicKey string, origData []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// Decrypt 解密
func Decrypt(privateKey string, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
