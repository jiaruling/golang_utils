package lib

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

//高级加密标准（Adevanced Encryption Standard ,AES）

type AesEncrypt struct {
	pwdKey []byte
}

func NewAesEncrypt(pwdKey string) *AesEncrypt {
	pwdLen := len(pwdKey)
	if pwdLen == 16 || pwdLen == 24 || pwdLen == 32 {
		return &AesEncrypt{
			pwdKey: []byte(pwdKey),
		}
	} else {
		return nil
	}
}

// 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
// key不能泄露
// var PwdKey = []byte("sqsEWQWE51wdq5wqSDdqwewdqSCQCQsW")

// PKCS7 填充模式
func (a *AesEncrypt) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 填充的反向操作，删除填充字符串
func (a *AesEncrypt) pKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}

// 实现加密
func (a *AesEncrypt) aesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = a.pKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 实现解密
func (a *AesEncrypt) aesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = a.pKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

// 加密base64
func (a *AesEncrypt) EnPwdCode(message []byte) ([]byte, error) {
	result, err := a.aesEcrypt(message, a.pwdKey)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(result)), err
}

// 解密
func (a *AesEncrypt) DePwdCode(message []byte) ([]byte, error) {
	//解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(string(message))
	if err != nil {
		return nil, err
	}
	//执行AES解密
	return a.aesDeCrypt(pwdByte, a.pwdKey)
}

// RSA 公钥私钥对非对称加密
type RsaEncrypt struct {
	bits        int
	privatePath string
	publicPath  string
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

var (
	privateFileName = "/private.pem"
	publicFileName  = "/public.pem"
)

func NewRsaEncrypt(bits int, privatePath, publicPath string) *RsaEncrypt {
	if !(bits == 2048 || bits == 3072 || bits == 4096) {
		bits = 2048
	}
	return &RsaEncrypt{
		bits:        bits,
		privatePath: privatePath,
		publicPath:  publicPath,
	}

}

// 生成RSA密钥对, 并将其保存到文件中
func (r *RsaEncrypt) generateRSAKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, r.bits)
	if err != nil {
		return err
	}
	r.privateKey = privateKey
	r.publicKey = &r.privateKey.PublicKey
	err = r.savePEMKey()
	if err != nil {
		return err
	}
	err = r.savePublicPEMKey()
	if err != nil {
		return err
	}
	return nil
}

func (r *RsaEncrypt) savePEMKey() error {
	outFile, err := os.Create(r.privatePath + privateFileName)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer outFile.Close()

	var privateKeyPEM = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(r.privateKey),
	}
	pem.Encode(outFile, privateKeyPEM)
	return nil
}

func (r *RsaEncrypt) savePublicPEMKey() error {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&r.privateKey.PublicKey)
	if err != nil {
		return err
	}

	var pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(r.publicPath + publicFileName)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer pemfile.Close()

	pem.Encode(pemfile, pemkey)
	return nil
}

// 加载公钥
func (r *RsaEncrypt) loadPublicKey() error {
	if r.publicKey != nil {
		return nil
	}

	pubKeyFile, err := os.ReadFile(r.publicPath + publicFileName)
	if err != nil {
		if os.IsNotExist(err) {
			err = r.generateRSAKeyPair()
			if err != nil {
				return err
			}
		}
		return err
	}

	block, _ := pem.Decode(pubKeyFile)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		r.publicKey = pub
		return nil
	default:
		return fmt.Errorf("unexpected key type")
	}
}

// 使用公钥加密数据
func (r *RsaEncrypt) EncryptWithPublicKey(msg []byte) (string, error) {
	if r.publicKey == nil {
		if err := r.loadPublicKey(); err != nil {
			return "", err
		}
	}
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil // []byte转成base64
}

// 使用公钥验证签名
func (r *RsaEncrypt) VerifyWithPublicKey(message []byte, signatureBase64 string) (bool, error) {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64) // base64转[]byte
	if err != nil {
		return false, err
	}
	if r.publicKey == nil {
		if err := r.loadPublicKey(); err != nil {
			return false, err
		}
	}
	hashed := sha256.Sum256(message)
	err = rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, hashed[:], signature)
	return err == nil, err
}

// 加载私钥
func (r *RsaEncrypt) loadPrivateKey() error {
	if r.privateKey != nil {
		return nil
	}
	privKeyFile, err := os.ReadFile(r.privatePath + privateFileName)
	if err != nil {
		if os.IsNotExist(err) {
			err = r.generateRSAKeyPair()
			if err != nil {
				return err
			}
		}
		return err
	}

	block, _ := pem.Decode(privKeyFile)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	r.privateKey = priv
	return nil
}

// 使用私钥解密数据
func (r *RsaEncrypt) DecryptWithPrivateKey(ciphertextBase64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64) // base64转[]byte
	if err != nil {
		return nil, err
	}
	if r.privateKey == nil {
		if err := r.loadPrivateKey(); err != nil {
			return nil, err
		}
	}
	return rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, ciphertext)
}

// 使用私钥签名数据
func (r *RsaEncrypt) SignWithPrivateKey(message []byte) (string, error) {
	if r.privateKey == nil {
		if err := r.loadPrivateKey(); err != nil {
			return "", err
		}
	}
	hashed := sha256.Sum256(message)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil // []byte转成base64
}
