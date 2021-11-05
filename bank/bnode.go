package bank

import(
		
	"crypto/ecdsa"

)

type BNode struct{
	Id					string						//Public Key as Address for banking key
	PubKey				*ecdsa.PublicKey			//Public Key For Banking Functions
	PrvKey				*ecdsa.PrivateKey			//Private Key For Banking Functions
}


