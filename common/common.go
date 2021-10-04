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

type BlockData struct {
	Year				uint					   //The year this data is valid for	
	Balance				map[uint64]big.Int		  //Map index is Block Number and the associated value of the account at that Block Number. History[3] would show the value of the account at the end of Block 2 and to include any of block 3 transactions
	Confirmations		map[uint64][]SignedTx	  //Map index is Block Number and the associated cofirmations of the account Balance at this height	
	EBLY				map[uint]big.Int		  //The Ending balance of past years with the map index equal to the year
}


type SignedTx struct {
	R				big.Int
	S				big.Int
	Node			uintptr  			//Able to look up Node and get its public key
}

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


