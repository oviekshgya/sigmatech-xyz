package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

const (
	KEYAES = "41097111528535889696911679134821"
	IVKEY  = "6383462337918274"

	SECRETKEY        = "pYJLL6mIaWNtx6Q4m0QIcq8Svzuv57Qp"
	SECRETREFRESTJWT = "PLSWepO8rqMBPxwQm7qvllxYvoJIdoZ7"
)

type Acceess struct {
	Id          int     `json:"id,omitempty"`
	Exp         float64 `json:"exp,omitempty"`
	AccessUuid  string  `json:"access_uuid,omitempty"`
	RefreshUuid string  `json:"refresh_uuid,omitempty"`
	IsChange    int     `json:"is_change,omitempty"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
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

func ExtractedExt(r *http.Request, tokenStringParam string) (*Acceess, error) {
	var tokenString string
	if tokenStringParam == "" {
		bearToken := r.Header.Get("Authorization")
		const BEARER_SCHEMA = "Bearer "
		//authHeader := c.GetHeader("Authorization")
		tokenString = bearToken[len(BEARER_SCHEMA):]

	} else {
		tokenString = tokenStringParam
	}

	if tokenString == "" {
		fmt.Println("tidak ada token bearer")
		return nil, fmt.Errorf("tidak ada token")
	}

	/*strArr := strings.Split(bearToken, " ")
	if len(strArr) < 2 {
		fmt.Println("ac", strArr, bearToken)
		return nil, fmt.Errorf("invalid Access")
	}*/
	fmt.Println("tokenString", tokenString)
	decrypted := Aes256Decode(tokenString)
	fmt.Println("decrypted", decrypted)

	if decrypted == "" {
		return nil, fmt.Errorf("unknwon access")
	}
	//fmt.Println("This is a decrypted:", string(decrypted))

	token, err := jwt.Parse(string(decrypted), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method")
		}
		return []byte(SECRETKEY), nil
	})
	if err != nil {
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok2 := claims["access_uuid"].(string)
		id := claims["key"].(float64)
		exp := claims["exp"].(float64)
		if ok2 == false {
			return nil, fmt.Errorf("unauthorized %s", accessUuid)
		} else {
			return &Acceess{
				Id:  int(id),
				Exp: exp,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}
