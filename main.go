package main

import (
	"fmt"
	"net/http"
	"github.com/fgeth/fg/common"
	"github.com/fgeth/fg/crypto"

)

func main(){
	   prvKey, err := createKey()
	   if err !=nil{
			fmt.Println("Could not generate PrivateKey")
		}else{
			fmt.Println("PrivateKey Created.")
		}
		plainText := "Starting server for testing HTTP POST...\n"
		
		fmt.Printf(plainText)
		//Create New Keccak State
		khState := NewKeccakState()
	
		//Generate Hash
		hash := HashData(khState, []byte(plainText))
		fmt.Printf(hash)
		
		
	   if err := http.ListenAndServe(":69420", nil); err != nil {
	    	log.Fatal(err)
	   }
}