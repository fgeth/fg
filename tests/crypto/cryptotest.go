package main

import (
	"crypto/ecdsa"
	"math/big"
	"fmt"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/crypto"

)

const (
	plainText = "This is the secret message"

)

func createKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()

}

func verifyHash(hash common.Hash, r *big.Int, s *big.Int, pubKey *ecdsa.PublicKey) string{
response :=""
if crypto.Verify(hash , r,  s, *pubKey) {
		response = "Hash was signed by prvtKey"
	}else{
		response = "Hash was not signed by prvtKey"
	}
return response

}

func main() {

	//Test to see if we can get an ECDSA PrivateKey generated
	prvKey, err := createKey()
	
	if err !=nil{
		fmt.Println("Could not generate PrivateKey")
	}else{
		fmt.Println("PrivateKey Created.")
	}
	prvKey2, err := createKey()
	
	if err !=nil{
		fmt.Println("Could not generate PrivateKey2")
	}else{
		fmt.Println("PrivateKey2 Created.")
	}
	
	
	//Create New Keccak State
	khState := crypto.NewKeccakState()
	
	//Generate Hash
	hash := crypto.HashData(khState, []byte(plainText))
	
	//Sign Hash returns 2 big.Ints
	r, s, err  := crypto.Sign(hash, *prvKey )
	if err !=nil{
		fmt.Println("Error Signing Hash with prvKey")
	}else{
		fmt.Println("Signed Hash with prvKey")
	}
	
	t, u, err  := crypto.Sign(hash, *prvKey2 )
	if err !=nil{
		fmt.Println("Error Signing Hash with prvKey2")
	}else{
		fmt.Println("Signed Hash with prvKey2")
	}
	
	//Response should be valid that the private key signed the hash
	response := verifyHash(hash , r,  s, &prvKey.PublicKey)
	fmt.Println(response)
	
	//Response should be invalid the private key did not sign the hash
	response = verifyHash(hash , t,  u, &prvKey.PublicKey)
	fmt.Println(response)
	
	

	
	



}