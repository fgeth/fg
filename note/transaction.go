package note

import(

)


type Transaction struct {
	Id						string				// Note Id
	BlockNumber				uint64				// BlockNumber when last transaction occured Must wait 10 blocks before spending a Note
	Hash					string				// Hash of Note proving ownership

}