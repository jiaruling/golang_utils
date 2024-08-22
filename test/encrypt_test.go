package test

import (
	"testing"

	"log"

	"github.com/jiaruling/golang_utils/lib"
)

func TestAesEncrypt(t *testing.T) {
	var PwdKey = []byte("sqsEWQWE51wdq5wqSDdqwewdqSCQCQsW")
	enc := lib.NewAesEncrypt(string(PwdKey))
	encrypted, err := enc.EnPwdCode([]byte("1234567890"))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := enc.DePwdCode(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("\n密文:%s\n明文:%s\n", string(encrypted), string(msg))
}

func TestRsaEncrypt(t *testing.T) {
	r := lib.NewRsaEncrypt(2048, ".", ".")
	message := []byte("Hello, RSA encryption!")
	// 公钥加密 & 私钥解密
	ciphertext, err := r.EncryptWithPublicKey(message)
	if err != nil {
		log.Fatal(err)
	}
	plaintext, err := r.DecryptWithPrivateKey(ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("\n明文:%s\n公钥加密后(base64):%s\n私钥解密后:%s\n", string(message), ciphertext, string(plaintext))

	// 私钥签名 & 公钥验签
	sign, err := r.SignWithPrivateKey(message)
	if err != nil {
		log.Fatal(err)
	}
	ok, err := r.VerifyWithPublicKey(message, sign)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("\n明文:%s\n信息签名(base64):%s\n验证结果:%v\n", string(message), sign, ok)
}
