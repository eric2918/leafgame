package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateKey(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	publicKey := privateKey.PublicKey

	privateByte := x509.MarshalPKCS1PrivateKey(privateKey)
	publicByte := x509.MarshalPKCS1PublicKey(&publicKey)

	blockPrivate := pem.Block{Type: "private key", Bytes: privateByte}
	blockPublic := pem.Block{Type: "public key", Bytes: publicByte}

	privateFile, errPri := os.Create("../../private.pem")
	if errPri != nil {
		return errPri
	}
	defer privateFile.Close()
	pem.Encode(privateFile, &blockPrivate)

	publicFile, errPub := os.Create("../../public.crt")
	if errPub != nil {
		return errPub
	}
	defer publicFile.Close()
	pem.Encode(publicFile, &blockPublic)

	return nil
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	bytes := make([]byte, fileInfo.Size())
	file.Read(bytes)

	return bytes, nil
}

func Encrypt(data []byte, filename string) ([]byte, error) {
	bytes, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	res, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func Decrypt(data []byte, filename string) ([]byte, error) {
	bytes, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	res, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
