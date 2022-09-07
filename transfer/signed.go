package transfer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type SignedRequest struct {
	UUID      string
	Signature string
	Data      interface{}
}

type SignedResponse struct {
	UUID      string
	Signature string
	Data      interface{}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func hash(message string) [32]byte {
	messageBytes := []byte(message)
	return sha256.Sum256(messageBytes)
}

func readKeyFile(keyFilePath string) (*pem.Block, []byte) {
	bytes, err := ioutil.ReadFile(keyFilePath)
	panicOnError(err)

	return pem.Decode(bytes)
}

/***********/
/* Request */
/***********/
func ParseSigned(r *http.Request, data interface{}, publicKeyPath string) (SignedRequest, error) {
	var req SignedRequest

	publicKey := ReadPublicKey(publicKeyPath)

	req.Data = data
	err := json.NewDecoder(r.Body).Decode(&req)

	verified := VerifyMessage([]byte(fmt.Sprintf("%v", req.Data)), req.Signature, publicKey)
	if verified != true {
		err = errors.New("Verification failed")
		req = SignedRequest{}
	}

	return req, err
}

func ReadPublicKey(keyFilePath string) *rsa.PublicKey {
	block, _ := readKeyFile(keyFilePath)
	if block == nil || block.Type != "PUBLIC KEY" {
		panicOnError(errors.New("Key is not readable or is no public key"))
	}

	keyInterface, err := x509.ParsePKCS1PublicKey(block.Bytes)
	panicOnError(err)

	return keyInterface
}

func VerifyMessage(message []byte, signature string, publicKey *rsa.PublicKey) bool {
	sig := hash(signature)

	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, message[:], sig[:])
	return err != nil
}

/************/
/* Response */
/************/
func RespondSigned(w http.ResponseWriter, data interface{}, keyFilePath string) error {
	var resp SignedResponse

	resp.UUID = fmt.Sprintf("%v", uuid.New())

	privateKey := ReadPrivateKey(keyFilePath)
	signature := SignMessage(fmt.Sprintf("%v", data), resp.UUID, privateKey)

	resp.Signature = fmt.Sprintf("%v", signature)

	return json.NewEncoder(w).Encode(&resp)
}

func ReadPrivateKey(keyFilePath string) *rsa.PrivateKey {
	block, _ := readKeyFile(keyFilePath)
	if block == nil { // check for block.Type == private identifier
		panicOnError(errors.New("Key is not readable or is no private key"))
	}

	keyInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	panicOnError(err)

	return keyInterface
}

func SignMessage(message string, uuid string, privateKey *rsa.PrivateKey) []byte {
	hashed := hash(message)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	panicOnError(err)

	return signature
}
