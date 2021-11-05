package bank

import(
	"crypto/ecdsa"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/note"
)

type Bank struct{
	Id					string						//Public Key as Address for banking key
	PubKey				*ecdsa.PublicKey			//Public Key For Banking Functions
	Collateral			note.Stack					//Stack that is staked to be Bank
	
}


//TODO Complete This
func SubmitPayment(payment note.Stack, theItem item.Buy) (note.Note, error){
var aNote note.Note
	return aNote, nil
}