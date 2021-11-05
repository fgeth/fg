package bank

import(

)

type BankMessage struct{
	Id			string				//Address of Recieving Bank's Public ECDSA Key
	Pay			string				//Encrypted Stack
	Item		string				//Encrypted Item
	
}

