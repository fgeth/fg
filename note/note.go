package note

import(
		"math/big"
		"time/time"
)

type Note struct{
	Id			string					//Generate a Public Key and turn it into an Address
	Amount		uint32					//Value in Virtual Dollars
	Stack		string					//Id of any stack this Note belongs too
	Hash		string					//Hash of Note amount and Id plus Current Owners Public Key
	PHash		string					//Hash of Note amount and Id plus Past Owners Public Key the OTP
	R			big.Int					//R value of Signature
	S    		big.Int					//S value of Signature
	OTP			string					//Singers Public Key set when transfering the Note to a New Owner
	PR			big.Int					//R value of Previous Signature
	PS    		big.Int					//S value of Previous Signature
	POTP		string					//Previous Singers Public Key set when transfering the Note to a New Owner
}

