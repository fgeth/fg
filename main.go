package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/fgeth/fg/crypto"

)

func main(){
	   prvKey, err := crypto.GenerateKey
	   if err !=nil{
			fmt.Println("Could not generate PrivateKey")
		}else{
			fmt.Println("PrivateKey Created.")
		}
		plainText := "Starting server for testing HTTP POST...\n"
		
		fmt.Printf(plainText)
		//Create New Keccak State
		khState := crypto.NewKeccakState()
	
		//Generate Hash
		hash := crypto.HashData(khState, []byte(plainText))
		fmt.Printf(hash)
		
		
	   if err := http.ListenAndServe(":69420", nil); err != nil {
	    	log.Fatal(err)
	   }
}