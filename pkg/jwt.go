package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

type Acceess struct {
	Id          int
	Exp         float64
	AccessUuid  string
	RefreshUuid string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(id int) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 8760).Unix() //expires after 30 min
	td.TokenUuid = uuid.New().String()
	//key, _ := GetAESEncrypted(fmt.Sprintf("%d", id))

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["key"] = id
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(SECRETKEY))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	td.RtExpires = time.Now().Add(time.Hour * 8760).Unix()
	td.RefreshUuid = fmt.Sprintf("%s++%d", td.TokenUuid, id)

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["key"] = id
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(SECRETREFRESTJWT))
	if err != nil {
		return nil, err
	}
	td.AccessToken, _ = Ase256(td.AccessToken)

	td.RefreshToken, _ = Ase256(td.RefreshToken)

	return td, nil
}

func Ase256(plaintext string) (string, error) {
	bKey := []byte(KEYAES)
	bIV := []byte(IVKEY)
	bPlaintext := PKCS5Padding([]byte(plaintext), aes.BlockSize, len(plaintext))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func Aes256Decode(cipherText string) (decryptedString string) {
	bKey := []byte(KEYAES)
	bIV := []byte(IVKEY)
	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		fmt.Println("err", err.Error())
		return ""
	}
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return string(PKCS5UnPadding256(cipherTextDecoded))
}

func PKCS5UnPadding256(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
