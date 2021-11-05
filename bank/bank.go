package bank

import(

	"github.com/fgeth/fg/note"
)

type Bank struct{
	Id					string						//Public Key as Address for banking key
	Collateral			note.Stack					//Stack that is staked to be a Bank
	BankNode			BNode						//Private Banking Node
	PublicNode			PBNode						//Public Banking Node
}





