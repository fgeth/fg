package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/aes"
    "crypto/cipher"
    //"crypto/sha256"
	"crypto/x509"
    "encoding/pem"
	"encoding/binary"
	//"encoding/hex"
	//"encoding/json"
    //"fmt"
    //"log"
    //"io"
	"io/ioutil"
	//"math"
	"math/big"
	"path/filepath"
	"os"
	"github.com/fgeth/fg/common"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/scrypt"
)
const (
	// StandardScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	ScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	ScryptP = 1
	
	//Crypto Version
	version = 1



	

)


type Key struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	PrivateKey *ecdsa.PrivateKey
}


type encryptedKeyJSONV3 struct {
	Address common.Address      `json:"address"`
	Crypto  CryptoJSON 			`json:"crypto"`
	Id      string     			`json:"id"`
	Version int       			`json:"version"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

type plainKeyJSON struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
	Id         string `json:"id"`
	Version    int    `json:"version"`
}



func Sign(hash common.Hash, prvKey *ecdsa.PrivateKey ) (*big.Int, *big.Int, error){
	r, s, err := ecdsa.Sign(rand.Reader, prvKey, hash[:])
	if err != nil {
		panic(err)
	}
	return r, s, err

}

func Verify(hash common.Hash, r *big.Int, s *big.Int, pubKey *ecdsa.PublicKey) bool{
	return ecdsa.Verify(pubKey, hash[:], r, s) 

}


// GenerateKey generates a new private key.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
}


// NewKeccakState creates a new KeccakState
func NewKeccakState() common.KeccakState {
	return sha3.NewLegacyKeccak256().(common.KeccakState)
}

// HashData hashes the provided data using the KeccakState and returns a 32 byte hash
func HashData(kh common.KeccakState, data []byte) (h common.Hash) {
	kh.Reset()
	kh.Write(data)
	kh.Read(h[:])
	return h
}

func HashToUint64(h common.Hash) uint64{
	data := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	for i:=0; i< len(h); i++{
		data[i] =h[i]	
	}
	theHashAsUint64 := binary.BigEndian.Uint64(data)


}


func Encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
    x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncoded), string(pemEncodedPub)
}

func Decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)

    return privateKey, publicKey
}






func WriteTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}


func StoreKey ( key *ecdsa.PrivateKey, auth string) error{
prvKey, PubKey := Encode(key, &key.PublicKey)

keyjson, err := Encrypt([]byte(auth), []byte(prvKey))
	if err != nil {
		return err
	}
	tmpName, err := WriteTemporaryKeyFile(PubKey, keyjson)
	os.Rename(tmpName, PubKey)
	return err
}

func GetKey(filename, auth string) (*ecdsa.PrivateKey,*ecdsa.PublicKey, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	key, err := Decrypt([]byte(auth), []byte(keyjson))
	if err != nil {
		return nil, nil, err
	}
	prvKey, pubKey := Decode(key,filename)
	// Make sure we're really operating on the requested key (no swap attacks)
	
	return prvKey, pubKey, nil
}
func Encrypt(password, data []byte) ([]byte, error) {
    key, salt, err := DeriveKey(password, nil)
    if err != nil {
        return nil, err
    }
    blockCipher, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err = rand.Read(nonce); err != nil {
        return nil, err
    }
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    ciphertext = append(ciphertext, salt...)
    return ciphertext, nil
}
func Decrypt(key, data []byte) (string, error) {
    salt, data := data[len(data)-32:], data[:len(data)-32]
    key, _, err := DeriveKey(key, salt)
    if err != nil {
        return "", err
    }
    blockCipher, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return "", err
    }
    nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}
func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
    if salt == nil {
        salt = make([]byte, 32)
        if _, err := rand.Read(salt); err != nil {
            return nil, nil, err
        }
    }
    key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
    if err != nil {
        return nil, nil, err
    }
    return key, salt, nil
}
