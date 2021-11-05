package bank

import(
		
	"crypto/ecdsa"
	"crypto/rsa"

)

type PBNode struct{
	Id					string						//Ecdsa Public Key as Address
	PubKey				rsa.PublicKey				//Public Key For Banking Functions
	PbKey				*ecdsa.PublicKey			//Public Key For verifying Banking Functions

}


