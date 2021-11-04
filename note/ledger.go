package note

import(

)

type Ledger struct{
	Notes  		map[string]Transaction		//Index is Note Id with the last associated transaction

}

