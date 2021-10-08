package common

import (
	"crypto/ecdsa"
	"math/big"
	"hash"
	
)

const (
	//FGE Account Address Length
	AddressLength = 20

	// HashLength is the expected length of the hash
	HashLength = 32
	

)


// Address represents the 20 byte address of an Address
type Address [AddressLength]byte

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte



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

type SignedTx struct {
    Accept			bool				//If node accepts this transaction or rejects the transaction
	R				big.Int
	S				big.Int
	Node			uintptr  			//Able to look up Node and get its public key
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