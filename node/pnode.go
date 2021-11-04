package node

import(

)


//Public Node with IP
type PNode struct {
	Id				uint64							//Node Ring Location
	Ip				string							//Node Onion Address
	Port			string							//Port that the node is running under
	PKStr			string							//Node PublicKey as string	

	
}
