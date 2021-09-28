package Uwim

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mergermarket/go-pkcs7"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
	"math/rand"
	"strings"
	"time"
)

var (
	SeedLength      int = 32
	PublicKeyLength int = 32

	TransactionRawKey = []byte{139, 111, 224, 92, 142, 122, 138, 224, 138, 118, 30, 229, 209, 155, 193, 186, 180, 234, 69, 249, 75, 71, 195, 105, 20, 61, 211, 13, 104, 253, 72, 5}
	TransactionRawIv  = []byte{22, 129, 2, 139, 42, 15, 11, 131, 158, 197, 170, 43, 114, 14, 178, 167}
)

// Mnemonic
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

//Seed
func SeedFromMnemonic(mnemonic string) ([]byte, error) {
	if len(strings.Split(mnemonic, " ")) != 24 {
		return nil, errors.New("Invalid mnemonic length")
	}

	return pbkdf2.Key([]byte(mnemonic), nil, 2048, SeedLength, sha512.New), nil
}

//SecretKey
func SecretKeyFromSeed(seed []byte) ([]byte, error) {
	if len(seed) != SeedLength {
		return nil, errors.New("Invalid seed length")
	}

	return ed25519.NewKeyFromSeed(seed), nil
}

func SecretKeyFromMnemonic(mnemonic string) ([]byte, error) {
	seed, err := SeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return SecretKeyFromSeed(seed)
}

//PublicKey
func PublicKeyFromSecretKey(secretKey []byte) []byte {
	publicKey := make([]byte, PublicKeyLength)
	copy(publicKey, secretKey[PublicKeyLength:])

	return publicKey
}

func PublicKeyFromMnemonic(mnemonic string) ([]byte, error) {
	seed, err := SeedFromMnemonic(mnemonic)
	if err != nil {
		return seed, nil
	}

	secretKey, err := SecretKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}

	return PublicKeyFromSecretKey(secretKey), nil
}

//Address
func AddressFromPublicKey(hrp string, publicKey []byte) (string, error) {
	address, err := bech32Encode(hrp, publicKey)
	if err != nil {
		return "", err
	}

	return address, nil
}

func AddressFromMnemonic(hrp, mnemonic string) (string, error) {
	publicKey, err := PublicKeyFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	return AddressFromPublicKey(hrp, publicKey)
}

//Signature
func SignMessageWithSecretKey(secretKey []byte, message []byte) []byte {
	return ed25519.Sign(secretKey, message)
}

// TransactionRaw
type Tx struct {
	Type       int64   `json:"type"`
	Nonce      int64   `json:"nonce"`
	HashTx     string  `json:"hashTx"`
	Height     int64   `json:"height"`
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
	TokenLabel string  `json:"tokenLabel"`
	Timestamp  string  `json:"timestamp"`
	Tax        float64 `json:"tax"`
	Signature  []byte  `json:"signature"`
	Comment    Comment `json:"comment"`
}

type Comment struct {
	Title string `json:"title"`
	Data  []byte `json:"data"`
}

type TxRaw struct {
	Type       int64        `json:"type"`
	Nonce      int64        `json:"nonce"`
	From       string       `json:"from"`
	To         string       `json:"to"`
	Amount     float64      `json:"amount"`
	TokenLabel string       `json:"tokenLabel"`
	Signature  string       `json:"signature"`
	Comment    TxRawComment `json:"comment"`
}

type TxRawComment struct {
	Title string `json:"title"`
	Data  string `json:"data"`
}

func GetTransactionRaw(mnemonic, sender, recipient, tokenLabel, commentTitle string, commentData []byte, amount float64, txType int64) (string, error) {

	timestamp, err := timestampUnix()
	if err != nil {
		return "", err
	}

	tx := Tx{
		Type:       txType,
		Nonce:      timestamp + rand.Int63(),
		HashTx:     "",
		Height:     0,
		From:       sender,
		To:         recipient,
		Amount:     amount,
		TokenLabel: tokenLabel,
		Timestamp:  "",
		Tax:        0,
		Signature:  nil,
		Comment: Comment{
			Title: commentTitle,
			Data:  commentData,
		},
	}

	jsonString, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}

	secretKey, err := SecretKeyFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	signature := SignMessageWithSecretKey(secretKey, jsonString)

	txRaw := TxRaw{
		Nonce:      tx.Nonce,
		From:       tx.From,
		To:         tx.To,
		Amount:     tx.Amount,
		TokenLabel: tx.TokenLabel,
		Type:       tx.Type,
		Signature:  base64.StdEncoding.EncodeToString(signature),
		Comment: TxRawComment{
			Title: commentTitle,
			Data:  string(commentData),
		},
	}

	jsonString, err = json.Marshal(txRaw)
	if err != nil {
		return "", err
	}

	txRawEncrypted, err := Encrypt(string(jsonString))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(txRawEncrypted), nil
}

func Encrypt(unencrypted string) ([]byte, error) {
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf(`plainText: "%s" has error`, plainText)
	}

	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return nil, err
	}

	block, err := aes.NewCipher(TransactionRawKey)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))

	mode := cipher.NewCBCEncrypter(block, TransactionRawIv)
	mode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

//Apparel
func bech32Encode(hrp string, data []byte) (string, error) {
	conv, err := ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", err
	}

	encoded, err := Encode(hrp, conv)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

func timestampUnix() (int64, error) {
	timestamp, err := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}

	return timestamp.UnixNano(), nil
}
