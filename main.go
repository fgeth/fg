package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/fgeth/fg/crypto"

)

func main(){
	   prvKey := crypto.GenerateKey
	   
		plainText := "Starting server for testing HTTP POST...\n"
		
		fmt.Printf(plainText)
		//Create New Keccak State
		//khState := crypto.NewKeccakState()
	
		//Generate Hash
		//hash := crypto.HashData(khState, []byte(plainText))
		//fmt.Printf(string(hash))
		
		
	   if err := http.ListenAndServe(":69420", nil); err != nil {
	    	log.Fatal(err)
	   }
}