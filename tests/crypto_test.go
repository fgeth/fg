package main

import (
	"fmt"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/crypto"

)


func createKey() *ecdsa.PrivateKey, error {
	return GenerateKey()

}


func main() {

	//Test to see if we can get an ECDSA PrivateKey generated
	prvKey, err := createKey()
	if err !:=nil{
		fmt.Println("Could not generate PrivateKey")
	}else{
		fmt.Println("PrivateKey Created.")
	}
	
	sign(hash common.Hash, prvKey ecdsa.PrivateKey )
	
	



}