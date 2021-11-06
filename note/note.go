package note

import(
		"math/big"
		
)

type Note struct{
	Id			string					//Generate a Public Key and turn it into an Address 						   32 bytes
	Amount		float64					//Value in Virtual Dollars 														8 bytes
	Coins		*big.Int				//Amount of Coins			    											   	8 bytes
	Hash		string					//Hash of Note Coins, Id and Owners Public Key								   32 bytes	
	R			*big.Int				//R value of Signature															8 bytes
	S    		*big.Int				//S value of Signature															8 bytes
	OTP			string					//Singers Public Key set when Destroying the Note							   32 bytes 
	
										//Between 80 and 128 bytes depending on if it been spent or not
}



