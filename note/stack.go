package note

import(

)

type Stack struct{
	Id			string 				    //Generate a publick key and turn it into an address
	Amount		uint64					//Value in virtual dollars
	Notes		[]Note					//Array of Notes
	Hash		string					//Hash of Stack amount, Notes, Id plus Current Owners Public Key
	PHash		string					//Hash of Stack amount, Notes, Id plus Past Owners Public Key the OTP
	R			big.Int					//R value of Signature
	S    		big.Int					//S value of Signature
	OTP			string					//Singers Public Key set when transfering the Note to a New Owner
	PR			big.Int					//R value of Previous Signature
	PS    		big.Int					//S value of Previous Signature
	POTP		string					//Previous Singers Public Key set when transfering the Note to a New Owner				
}



//TODO Function to Create Stack


//TODO Function to Break Stack