package wallet


import(
	"crypto/ecdsa"
	"math/big"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/note"
	"github.com/fgeth/fg/transaction"

)

type Wallet struct {
	Id			string							//Address
	Coins		float64
	Wei			*big.Int
	Dollars		float64
	CoinValue	float64
	Items		Selling							//Items that are for sell
	Buy			Buying							//Items we have bought
	Notes		map[string]note.Note			//Index is note Id 	
	Stack		[]note.Stack
	Debits		Debits
	Auth		string
}


type Selling struct {
	Item			map[string]item.Item								 //Index of Item Id and the Item
	Keys			map[string][]*ecdsa.PrivateKey						//Index is Item Id and array of private keys for the transaction
	
}

type Buying struct {
	Item			map[string]item.Item								 //Index is Item Id and the Item
	Tx				map[string][]transaction.Transaction				//Index is Item Id	and the array of transactions that go with that Item
	Keys			map[string][]*ecdsa.PrivateKey						//Index is Item Id and array of private keys for the transaction

}

type Debits struct {
	Debit	map[string]transaction.BaseTransaction								//Index of debit Tranasaction Hash with transaction
	
}







func removeNote(s []note.Note) []note.Note {
    
    return s[:len(s)-1]
}