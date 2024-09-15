package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

//Reference document
//http://www.topgoer.com/%E5%85%B6%E4%BB%96/%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86/%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86.html
//Advanced Encryption Standard (Adevanced Encryption Standard, AES)

// 16, 24, and 32-bit strings correspond to AES-128, AES-192, AES-256 encryption methods respectively
// key cannot be leaked
// var PwdKey = []byte("DIS**#KKKDJJSKDI")
var PwdKey = []byte("alexalexalexalexalexalexalexalex")

// PKCS7 padding mode
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//The function of the Repeat() function is to copy the padding of the slice []byte{byte(padding)}, and then merge it into a new byte slice to return
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// The reverse operation of filling, delete the filling string
func PKCS7UnPadding1(origData []byte) ([]byte, error) {
	//Get data length
	length := len(origData)
	if length == 0 {
		return nil, errors.New("Encrypted string error!")
	} else {
		//Get the fill string length
		unpadding := int(origData[length-1])
		//Intercept the slice, delete the padding bytes, and return the plaintext
		return origData[:(length - unpadding)], nil
	}
}

// Encryption
func AesEcrypt(origData []byte, key []byte) ([]byte, error) {
	//Create an instance of the encryption algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//Get the size of the block
	blockSize := block.BlockSize()
	//Padded the data so that the data length meets the demand
	origData = PKCS7Padding(origData, blockSize)
	//CBC encryption mode in AES encryption method
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//Perform encryption
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// Realize decryption
func AesDeCrypt(cypted []byte, key []byte) (string, error) {
	//Create an instance of the encryption algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//Get the block size
	blockSize := block.BlockSize()
	//Create an encrypted client instance
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//This function can also be used to decrypt
	blockMode.CryptBlocks(origData, cypted)
	//Remove the fill string
	origData, err = PKCS7UnPadding1(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}

// Encrypted base64
func EnPwdCode(pwdStr string) string {
	pwd := []byte(pwdStr)
	result, err := AesEcrypt(pwd, []byte(PwdKey))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(result)
}

// Decrypt
func DePwdCode(pwd string) string {
	temp, _ := hex.DecodeString(pwd)
	//Perform AES decryption
	res, _ := AesDeCrypt(temp, []byte(PwdKey))
	return res
}

/*
func main() {

	//aes encryption
	destring := `{"name":"Lyoha","site":"http://www.alexserv.com"}`
	deStr := EnPwdCode(destring)
	fmt.Println(deStr) //4f4d74c15e0ad4afb323a17927b1176ecb0c95ecbdf8e776ceb093499e3ff4c45157b007ae7dff1688ac2d2bf9fef28644922a1b3bbc6ef5881cb1ed0dff298a

	//aes decrypt
	decodeStr := DePwdCode("2f6c3f5f2b64b47d00fcc9689b5c828fb05267162af5b1d067ec87f2a2e2294d13e1173039cc8994a324bfdbac28f581591e344cc3ee967506f53f7b06d916b6")
	fmt.Println(decodeStr) //{"name":"Novice Tutorial 11","site":"http://www.runoob.com"}
}*/
