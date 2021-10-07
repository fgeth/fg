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
	PubKey			ecdsa.PublicKey
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


