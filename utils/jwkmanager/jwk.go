package jwkmanager

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"gin-auth-mongo/utils/consts"
	"log"
	mathRand "math/rand"
	"os"
	"sync"

	"github.com/square/go-jose/v3"
)

// global variables for storing keys and locks
var (
	privateJWKs  []jose.JSONWebKey
	publicJWKs   []jose.JSONWebKey
	keyCacheLock sync.RWMutex
)

// load signing keys from file to memory
func LoadSigningKeys(publicKeysFilePath, privateKeysFilePath string) error {
	keyCacheLock.Lock()
	defer keyCacheLock.Unlock()

	// check file exists
	if _, err := os.Stat(publicKeysFilePath); os.IsNotExist(err) {
		return errors.New(publicKeysFilePath + " not found")
	}

	data, err := os.ReadFile(publicKeysFilePath)
	if err != nil {
		return err
	}

	var publicKeys []jose.JSONWebKey
	if err := json.Unmarshal(data, &publicKeys); err != nil {
		return err
	}

	publicJWKs = publicKeys
	log.Println("Successfully loaded ", publicKeysFilePath, " keys into memory")

	if _, err := os.Stat(privateKeysFilePath); os.IsNotExist(err) {
		return errors.New(privateKeysFilePath + " not found")
	}

	data, err = os.ReadFile(privateKeysFilePath)
	if err != nil {
		return err
	}

	var privateKeys []jose.JSONWebKey
	if err := json.Unmarshal(data, &privateKeys); err != nil {
		return err
	}

	privateJWKs = privateKeys
	log.Println("Successfully loaded ", privateKeysFilePath, " keys into memory")

	return nil
}

// get public JWKs
func GetPublicJWKs() ([]jose.JSONWebKey, error) {
	keyCacheLock.RLock()
	defer keyCacheLock.RUnlock()

	if len(publicJWKs) == 0 {
		return nil, errors.New("no public keys found")
	}

	return publicJWKs, nil
}

// get private JWKs
func GetPrivateJWKs() ([]jose.JSONWebKey, error) {
	keyCacheLock.RLock()
	defer keyCacheLock.RUnlock()

	if len(privateJWKs) == 0 {
		return nil, errors.New("no private keys found")
	}
	return privateJWKs, nil
}

// get random JWK
func GetRandomJWK() (jose.JSONWebKey, error) {

	// read lock
	keyCacheLock.RLock()
	defer keyCacheLock.RUnlock()

	if len(privateJWKs) == 0 {
		return jose.JSONWebKey{}, errors.New("no signing keys found")
	}

	return privateJWKs[mathRand.Intn(len(privateJWKs))], nil
}

// get key by kid (for verifying JWT)
func GetKeyByID(kid string) (jose.JSONWebKey, error) {
	keyCacheLock.RLock()
	defer keyCacheLock.RUnlock()

	for _, key := range privateJWKs {
		if key.KeyID == kid {
			return key, nil
		}
	}
	return jose.JSONWebKey{}, errors.New("key not found")
}

// clear key files
func clearKeyFiles(privateFilePath, publicFilePath string) {
	if _, err := os.Stat(privateFilePath); err == nil {
		e := os.Remove(privateFilePath)
		if e != nil {
			log.Fatal("Error removing ", privateFilePath, ":", e)
		}
	}
	if _, err := os.Stat(publicFilePath); err == nil {
		e := os.Remove(publicFilePath)
		if e != nil {
			log.Fatal("Error removing ", publicFilePath, ":", e)
		}
	}
}

// update keys
func UpdateKeys() error {
	privateFilePath := consts.PRIVATE_KEYS_FILE
	publicFilePath := consts.PUBLIC_KEYS_FILE

	// clear old key files
	clearKeyFiles(privateFilePath, publicFilePath)

	// new jwks
	newPrivateJWKs := make([]jose.JSONWebKey, 0)
	newPublicJWKs := make([]jose.JSONWebKey, 0)

	// generate 5 new jwks
	for i := 0; i < 5; i++ {
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}

		hasher := sha256.New()
		hasher.Write(publicKey)
		keyThumbprint := hex.EncodeToString(hasher.Sum(nil))

		publicJWK := jose.JSONWebKey{Key: publicKey, KeyID: keyThumbprint, Algorithm: "Ed25519", Use: "sig"}
		newPublicJWKs = append(newPublicJWKs, publicJWK)

		privateJWK := jose.JSONWebKey{Key: privateKey, KeyID: keyThumbprint, Algorithm: "Ed25519", Use: "sig"}
		newPrivateJWKs = append(newPrivateJWKs, privateJWK)
	}

	jsonData, err := json.Marshal(newPrivateJWKs)
	if err != nil {
		panic(err)
	}
	os.Mkdir(".private", os.ModePerm)
	if err := os.WriteFile(privateFilePath, jsonData, 0644); err != nil {
		panic(err)
	}

	jsonData, err = json.Marshal(newPublicJWKs)
	if err != nil {
		panic(err)
	}
	os.Mkdir(".public", os.ModePerm)
	if err := os.WriteFile(publicFilePath, jsonData, 0644); err != nil {
		panic(err)
	}

	privateJWKs = newPrivateJWKs
	publicJWKs = newPublicJWKs

	err = LoadSigningKeys(publicFilePath, privateFilePath)
	if err != nil {
		return err
	}

	return nil
}
