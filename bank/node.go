package bank

import(
		
	"crypto/ecdsa"

)

type Node struct{
	Id					string						//Public Key as Address for banking key
	PubKey				*ecdsa.PublicKey			//Public Key For Banking Functions
	PrvKey				*ecdsa.PrivateKey			//Private Key For Banking Functions
}


