package ring

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"strconv"
	"time"
	"github.com/fgeth/fg/bank"
	"github.com/fgeth/fg/node"
	
)

type Ring struct{
	Id		uint64					//This Nodes current location on Ring
	Table 	[]FingerTable			//The 32 Nodes On Ring that the IP address is Known
	Nodes	[]node.RNode		    //Index is location on Ring and the nodes corresponding Public Key 
	Banks	[]bank.PBNode			//Array of All Known and Trusted Banks
	
}
func(r Ring) SaveRing (dirname string) {
	path :=filepath.Join(dirname, "ring")

	_, err := os.Stat(dirname)
    if err !=nil {
		fmt.Println("error ", err)
		err := os.Mkdir(dirname, 0755)
		if err !=nil{
			fmt.Println("failed to make root directory. Can not save Ring.", err)
				
			}else{
				err = os.Mkdir(path, 0755)
				if err !=nil{
					fmt.Println("failed to make ring directory", err)
				}
			}
			
		}else{
			err = os.Mkdir(path, 0755)
			if err !=nil{
				fmt.Println("failed to make ring directory", err)
			}
		}

    
  
	fileName := filepath.Join(path, "ring.json")
	
	file, _ := json.MarshalIndent(r, "", " ")
	
	
	err = ioutil.WriteFile(fileName, file, 0644)
	if err !=nil{
		fmt.Println("failed to save file", err)
	}
}	
func ImportRing(dirname string ) (Ring, error){

	path :=filepath.Join(dirname, "ring")
	var ring Ring
	var errRing Ring
	fileName := filepath.Join(path, "ring.json")
	//fmt.Println("File Name : ", fileName )
	_, e := os.Stat(fileName)
	if e != nil{
		dirname, _ := os.UserHomeDir()
		path :=filepath.Join(dirname, "ring")
		fileName := filepath.Join(path, "ring.json")
		//fmt.Println("File Name : ", fileName )
		_, e1 := os.Stat(fileName)
		
		if e1 != nil{
			return errRing, e1
			
		}else{
			file, _ := ioutil.ReadFile(fileName)
			//fmt.Println("Unmarshalling File : ", fileName )
			err :=json.Unmarshal(file, &ring)
			
			if err != nil {
				fmt.Println("couldn't unmarshal parameters", err)
				return errRing, err

			}
		}
		
		//fmt.Println( e )
	}else{
		file, _ := ioutil.ReadFile(fileName)
		
		err := json.Unmarshal(file, &ring)
		//fmt.Println("Unmarshalling File : ", fileName )
	if err != nil {
        fmt.Println("couldn't unmarshal parameters", err)
			return errRing, err
    }
		return ring, nil
	}

	return errRing, e
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

func (r Ring) RotateFingerTable(n node.PNode, ringId uint64){
	var tmpNode  node.PNode
	ft := len(r.Table)
	if ft <32 || (n.Id >= ringId && n.Id <= ringId +uint64(8)){
		for x:=0; x< ft; x +=1{
		
		if r.Table[x].Id == n.Id{
			tmpNode =r.Table[x].Node
			r.Table[x].Node = n
			n = tmpNode 
			n.Id += uint64(1)
			
		}


		}
		if ft <32{
			r.Table = append(r.Table, FingerTable{Id: n.Id ,Node: n})
		}
	}
}






func(r Ring) FindPeer() node.PNode{
var result node.PNode
var theUrl =""
x:=0;
	for a:=3; a<256; a +=1{
		if a !=10{
			for b:=0; b<256; b +=1{
				if ((a !=172) &&(b <16 || b>31)) ||((a !=192) &&(b!=168)) {
					for c:=0; c<256; c +=1{
						for d:=0; d<256; d +=1{
							if a != 10 || (a != 172 && (b <16 && b>31))|| (a !=192 && b != 168){
								x +=1
								theUrl = "http://"+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(c)+"."+strconv.Itoa(d)
								fmt.Print(theUrl)
								timeout := 3 * time.Second
								url := theUrl +":42069"
								_, err := net.DialTimeout("tcp",url, timeout)
								if err !=nil{
									fmt.Println(" Port 42069 Not Listening on ", theUrl)
								}else{
									x +=1
									go r.CheckPeer42069(theUrl, strconv.Itoa(a),strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d) )
									if len(r.Table) > 0{
										break
									}
									
								}

								timeout = 3 * time.Second
								url = theUrl +":80"
								_, err = net.DialTimeout("tcp",url, timeout)
								if err !=nil{
									fmt.Println(" Port 80 Not Listening on ", theUrl)
								}else{
									
									go r.CheckPeer80(theUrl, strconv.Itoa(a),strconv.Itoa(b), strconv.Itoa(c), strconv.Itoa(d) )
									if len(r.Table) > 0{
										break
									}
									
								}
								if x % 5000 ==0 {
										 time.Sleep(10 * time.Second)
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
var result node.PNode
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
var result node.PNode
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