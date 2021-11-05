package note

import(
		"math/big"

)

type Note struct{
	Id			string					//Generate a Public Key and turn it into an Address 						   32 bytes
	Amount		float64					//Value in Virtual Dollars 														4 bytes
	Stack		string					//Id of any stack this Note belongs too										   32 bytes
	Hash		string					//Hash of Note amount and Id plus Current Owners Public Key					   32 bytes	
	PHash		string					//Hash of Note amount and Id plus Past Owners Public Key the OTP			   32 bytes
	R			big.Int					//R value of Signature															8 bytes
	S    		big.Int					//S value of Signature															8 bytes
	OTP			string					//Singers Public Key set when transfering the Note to a New Owner			   32 bytes 	
	PR			big.Int					//R value of Previous Signature													8 bytes
	PS    		big.Int					//S value of Previous Signature												    8 bytes
	POTP		string					//Previous Singers Public Key set when transfering the Note to a New Owner	   32 bytes	
										//Between 196 and 228 bytes depending on if it is part of a stack or not
}



