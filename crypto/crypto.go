package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
	"crypto/aes"
    "crypto/cipher"
	"crypto/x509"
    "encoding/pem"
	"encoding/binary"
	"fmt"
	"hash"
	"io/ioutil"
	"math/big"
R	"math/rand"
	"path/filepath"
	"os"
	"strings"
	"time"
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

	// HashLength is the expected length of the hash
	HashLength = 32


	

)
var (
    lowerCharSet   = "abcdedfghijklmnopqrst"
    upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    specialCharSet = "!@#$%&*"
    numberSet      = "0123456789"
    allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)






// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash []byte

type Signer struct {
	PubKey			*ecdsa.PublicKey
	R				big.Int
	S				big.Int
}

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}
// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// HashData hashes the provided data using the KeccakState and returns a 32 byte hash
func HashData(kh KeccakState, data []byte) (h Hash) {
	kh.Reset()
	kh.Write(data)
	kh.Read(h[:])
	return h
}

type SignedTx struct {
    Accept			bool				//If node accepts this transaction or rejects the transaction
	R				big.Int
	S				big.Int
	PubKey			string 			//Able to look up Node and get its public key this is the nodes publickey as string
}



func Sign(hash Hash, prvKey *ecdsa.PrivateKey ) (*big.Int, *big.Int){
	r, s, err := ecdsa.Sign(rand.Reader, prvKey, hash[:])
	if err != nil {
		fmt.Println(err)
	}
	return r, s

}

func Verify(hash Hash, r *big.Int, s *big.Int, pubKey *ecdsa.PublicKey) bool{
	return ecdsa.Verify(pubKey, hash[:], r, s) 

}


// GenerateKey generates a new private key.
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//elliptic.P256() can do 27,000 verifications per second vs 1,700 a second with elliptic.P384()
	//going with elliptic.P256 as each key will only be protecting $10 256 bits is good enough 
	//for example each bitcoin wallet uses the 256 curve and each wallet stores way more than $10
}



func HashToUint64(h Hash) (uint64, uint64, uint64, uint64){
	data0 := []byte{h[0],h[1],h[2],h[3],h[4],h[5],h[6],h[7]}
	data1 := []byte{h[8],h[9],h[10],h[11],h[12],h[13],h[14],h[15]}
	data2 := []byte{h[16],h[17],h[18],h[19],h[20],h[21],h[22],h[23]}
	data3 := []byte{h[24],h[25],h[26],h[27],h[28],h[29],h[30],h[31]}

	uintA := binary.BigEndian.Uint64(data0)
	uintB := binary.BigEndian.Uint64(data1)
	uintC := binary.BigEndian.Uint64(data2)
	uintD := binary.BigEndian.Uint64(data3)
	return uintA, uintB, uintC, uintD	

}

func Uint64ToHash(uintA, uintB, uintC, uintD uint64) Hash{
cn1 := make([]byte, 8)
cn2 := make([]byte, 8)
cn3 := make([]byte, 8)
cn4 := make([]byte, 8)
var h  Hash
binary.BigEndian.PutUint64(cn1, uintA)

binary.BigEndian.PutUint64(cn2, uintB)

binary.BigEndian.PutUint64(cn3, uintC)
binary.BigEndian.PutUint64(cn4, uintD)
for x:=0; x <8; x+=1{
	h[x] =cn1[x]

}
for x:=0; x <8; x+=1{
	h[x+8] =cn2[x]

}
for x:=0; x <8; x+=1{
	h[x+16] =cn3[x]

}
for x:=0; x <8; x+=1{
	h[x+24] =cn4[x]

}
return h
}

func EncodePrv(privateKey *ecdsa.PrivateKey) (string) {
    x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded)
}
func Encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
    x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncoded), string(pemEncodedPub)
}

func EncodePubKey( publicKey *ecdsa.PublicKey) (string) {
    
    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncodedPub)
}

func DecodePubKey( pemEncodedPub string) (*ecdsa.PublicKey) {
   

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*ecdsa.PublicKey)

    return publicKey
}

func DecodePrv(pemEncoded string) (*ecdsa.PrivateKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

    return privateKey
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
func Decrypt(password, data []byte) (string, error) {
    salt, data := data[len(data)-32:], data[:len(data)-32]
    key, _, err := DeriveKey(password, salt)
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


func createPass() string{
    R.Seed(time.Now().Unix())
    minSpecialChar := 32
    minNum := 32
    minUpperCase := 32
    passwordLength := 128
    return generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)
	
}

func generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
    var password strings.Builder

    //Set special character
    for i := 0; i < minSpecialChar; i++ {
        random := R.Intn(len(specialCharSet))
        password.WriteString(string(specialCharSet[random]))
    }

    //Set numeric
    for i := 0; i < minNum; i++ {
        random := R.Intn(len(numberSet))
        password.WriteString(string(numberSet[random]))
    }

    //Set uppercase
    for i := 0; i < minUpperCase; i++ {
        random := R.Intn(len(upperCharSet))
        password.WriteString(string(upperCharSet[random]))
    }

    remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
    for i := 0; i < remainingLength; i++ {
        random := R.Intn(len(allCharSet))
        password.WriteString(string(allCharSet[random]))
    }
    inRune := []rune(password.String())
	R.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
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

func GenerateRSAKey() rsa.PrivateKey{
privateKey, err := rsa.GenerateKey(rand.Reader, 4096) 
fmt.Println(err)
return *privateKey
}


func EncodeRSA(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (string, string) {
    x509Encoded:= x509.MarshalPKCS1PrivateKey(privateKey)
    pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncoded), string(pemEncodedPub)
}

func DecodeRSA(pemEncoded string, pemEncodedPub string) (*rsa.PrivateKey, *rsa.PublicKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParsePKCS1PrivateKey(x509Encoded)

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*rsa.PublicKey)

    return privateKey, publicKey
}

func DecodeRSAPvKey(pemEncoded string) (rsa.PrivateKey, rsa.PublicKey) {
    block, _ := pem.Decode([]byte(pemEncoded))
    x509Encoded := block.Bytes
    privateKey, _ := x509.ParsePKCS1PrivateKey(x509Encoded)



    return *privateKey, privateKey.PublicKey
}

func RSAEncrypt(secretMessage string, key rsa.PublicKey) string {
    label := []byte("OAEP Encrypted")
    rng := rand.Reader
    ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
    fmt.Println(err)
    return base64.StdEncoding.EncodeToString(ciphertext)
}

func RSADecrypt(cipherText string, privKey rsa.PrivateKey) string {
    ct, _ := base64.StdEncoding.DecodeString(cipherText)
    label := []byte("OAEP Encrypted")
    rng := rand.Reader
    plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
    fmt.Println(err)
   
    return string(plaintext)
}
func EncodeRSAPubKey( publicKey *rsa.PublicKey) (string) {
    
    x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
    pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

    return string(pemEncodedPub)
}

func DecodeRSAPubKey( pemEncodedPub string) (rsa.PublicKey) {
   

    blockPub, _ := pem.Decode([]byte(pemEncodedPub))
    x509EncodedPub := blockPub.Bytes
    genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
    publicKey := genericPublicKey.(*rsa.PublicKey)

    return *publicKey
}


func StoreRSAKey ( key rsa.PrivateKey, auth string, fileName string) error{
prvKey, PubKey := EncodeRSA(&key, &key.PublicKey)

keyjson, err := Encrypt([]byte(auth), []byte(prvKey))
	if err != nil {
		return err
	}
	fmt.Println(PubKey)
	tmpName, err := WriteTemporaryKeyFile(fileName, keyjson)
	os.Rename(tmpName, fileName)
	return err
}

func GetRSAKey(filename, auth string) (rsa.PrivateKey,rsa.PublicKey, error) {
	// Load the key from the keystore and decrypt its contents
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	key, err := Decrypt([]byte(auth), []byte(keyjson))
	if err != nil {
		fmt.Println(err)
	}
	prvKey, pubKey := DecodeRSAPvKey(key)
	// Make sure we're really operating on the requested key (no swap attacks)
	
	return prvKey, pubKey, nil
}