package bank

import(
	"time"
	"github.com/fgeth/fg/item"
	"github.com/fgeth/fg/note"

)

type Escrow struct{
	Id			string					//Generate a Public Key and turn it into an Address 						   32 bytes
	Item		item.Item				//Item Being Bought											
	SOTP		string					//Seller Publick Key														   32 bytes 	
	BOTP		string					//Buyers Public Key															   32 bytes	
	Shipped		bool					//Set to true when seller ships or transfers item to buyer
	Recieved	bool					//Set to true by buyer when item is recieved or after time expires on escrow
	Expires		time.Time				//Set when Escrow is created set to 2 weeks by default Payment reverts to Buyer if seller never ships or Seller if they have shipped and provided Tracking
	Tracking	string					//Any tracking number provided by seller
	Note		[]note.Note				//Array of Notes to cover Escrow
	Stack		note.Stack				//Stack to cover Escrow

}