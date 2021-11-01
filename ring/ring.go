package ring

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"strconv"
	"time"
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
x:=0;
	for a:=1; a<256; a +=1{
		if a !=10{
			for b:=0; b<256; b +=1{
				if ((a !=172) &&(b <16 || b>31)) ||((a !=192) &&(b!=168)) {
					for c:=0; c<256; c +=1{
						for d:=0; d<256; d +=1{
							if a != 10 || (a != 172 && (b <16 && b>31))|| (a !=192 && b != 168){
								theUrl = "http://"+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d)
								fmt.Print(theUrl)
								timeout := 1 * time.Second
								url := theUrl +":42069"
								_, err := net.DialTimeout("tcp",url, timeout)
								if err !=nil{
									fmt.Println(" Port 42069 Not Listening on ", theUrl)
								}else{
									x +=1
									go r.CheckPeer42069(theUrl, strconv.Itoa(a),strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d) )
									if x % 5000 ==0 {
										 time.Sleep(8 * time.Second)
									}
								}
								timeout = 1 * time.Second
								url = theUrl +":80"
								_, err = net.DialTimeout("tcp",url, timeout)
								if err !=nil{
									fmt.Println(" Port 80 Not Listening on ", theUrl)
								}else{
									x +=1
									go r.CheckPeer80(theUrl, strconv.Itoa(a),strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d) )
									if x % 5000 ==0 {
										 time.Sleep(8 * time.Second)
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


func (r Ring) CheckPeer42069(theUrl, a, b, c, d string){
var result node.Node
var finger FingerTable
var Mtx	sync.Mutex
		url :=theUrl +":42069/getPeer"
		resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Still Searching for peer", err)
				}else{
							defer resp.Body.Close()
							body, err := ioutil.ReadAll(resp.Body)
							
							
							err = json.Unmarshal([]byte(body), &result)
							if err != nil {
									fmt.Println("Error unmarshaling data from request.")
								}else{
								if result.Ip == a+"."+b+"."+c+"."+d{
									
									finger.Id = result.Id
									finger.Node = result
									Mtx.Lock()
									r.Table = append(r.Table, finger)
									Mtx.Unlock()
								}
							}
						}

}

func (r Ring) CheckPeer80(theUrl, a, b, c, d string){
var result node.Node
var finger FingerTable
var Mtx	sync.Mutex
		theUrl = theUrl+":80/getPeer"
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
								if result.Ip == a+"."+b+"."+c+"."+d{
									finger.Id = result.Id
									finger.Node = result
									Mtx.Lock()
									r.Table= append(r.Table, finger)
									Mtx.Unlock()
								}
							}
						}

}