package ring

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/fgeth/fg/node"
)

type Ring struct{
	Id		uint64					//This Nodes current location on Ring
	Table 	[]FingerTable			//The 32 Nodes On Ring that the IP address is Known
	Nodes	[]node.RNode		    //Index is location on Ring and the nodes corresponding Public Key 
	
}




func(r Ring) RotateKeys(n node.RNode){
	var tmpNode  node.RNode

	
	for x:=0; x < len(r.Nodes); x +=1{
		
		if n.Id < r.Nodes[x].Id {
			tmpNode.Id = r.Nodes[x].Id
			tmpNode.PKStr = r.Nodes[x].PKStr
			r.Nodes[x] = n
			n = tmpNode
		}
	
	}
	r.Nodes = append(r.Nodes, n)
	

}

func (r Ring) RotateFingerTable(n node.Node){
	var tmpNode  node.Node
for x:=0; x< 32; x +=1{
		
		if r.Table[x].Id == n.Id{
			tmpNode =r.Table[x].Node
			r.Table[x].Node = n
			n = tmpNode 
			n.Id += uint64(1)
			
		}


}
}






func(r Ring) FindPeer() node.Node{
var result node.Node
var theUrl =""
	for a:=1; a<256; a +=1{
		if a !=10{
			for b:=0; b<256; b +=1{
				if ((a !=172) &&(b <16 || b>31)) ||((a !=192) &&(b!=168)) {
					for c:=0; c<256; c +=1{
						for d:=0; d<256; d +=1{
							if a != 10 || (a != 172 && (b <16 && b>31))|| (a !=192 && b != 168){
								theUrl = "http://"+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d)+":42069/getPeer"
								fmt.Print(theUrl)
								resp, err := http.Get(theUrl)
									if err != nil {
										theUrl = "http://"+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d)+":80/getPeer"
										fmt.Print(theUrl)
										resp, err := http.Get(theUrl)
											if err != nil {
												fmt.Println("Still Searching for peer", err)
											}else{
												defer resp.Body.Close()
												body, err := ioutil.ReadAll(resp.Body)
											
												
												err = json.Unmarshal([]byte(body), &result)
												if err != nil {
														fmt.Println("Error unmarshaling data from request.")
													}else{
													if result.Ip == strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d){
														return result
													}
												}
											}
									}else{
												defer resp.Body.Close()
												body, err := ioutil.ReadAll(resp.Body)
												
												
												err = json.Unmarshal([]byte(body), &result)
												if err != nil {
														fmt.Println("Error unmarshaling data from request.")
													}else{
													if result.Ip == strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d){
														return result
													}
												}
											}
							}
							
						}
					}
				
				
				}
			}
		}
	}
	return result
}
