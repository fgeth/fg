package node



import (
    "crypto/sha1"
	"crypto/rsa"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"flag"
	"log"
	"bufio"
	"hash"
	"math/big"
	"sync"
	"time"
	"encoding/json"
	"net/http"
	"net/url"
	 "os"

)

type Node struct {
	Id				string							//Node Public Key as a string
	Ip				string							//Node IP or fully qualified domain name that can resolve to an IP
	Port			string							//Port that the node is running under
	PubKey			*ecdsa.PublicKey				//Nodes Public Key
	PrvKey			*ecdsa.PrivateKey				//Nodes Private Key
	Writer			bool							//True if a current Block Node 
	Leader			bool							//True if the current Block Leader
	NumNodes		uint32							//Tracks Number of Block nodes that have submited Txs
	Comms			Comm							//Node RSA Keys
	Items			Selling							//Items that are for sell
	Mtx	 			sync.Mutex
}

type Nodes struct {
	Node 			map[string]Node  				//All active and inactive Nodes.  Easy to get PublicKey using string of public key as map index 

}

type Comm struct{

	RsaPrvKeys		map[rsa.PublicKey][]rsa.PrivateKey	//index is the RSA publicKey

	
}

type Selling struct {
	Item			[]Item							//Array of Items
	Tx				Transaction						//Tranaction shell just has the Debit BaseTransactions
	
}


func (node *Node) genRsa() {
 prvKey := GenerateRSAKey()
 pubKey := prvKey.PublicKey
 secret := RSAEncrypt("Secret", pubKey)
 fmt.Println("Encrypted Message =", secret)
 clearText := RSADecrypt(secret, prvKey)
 fmt.Println("Message =", clearText)
 publicKey := EncodeRSAPubKey(&pubKey)
 fmt.Println("Publick Key =", publicKey)
 pKey := DecodeRSAPubKey(publicKey)
 secret2 := RSAEncrypt("Secret2", pKey)
 fmt.Println("Encrypted Message =", secret2)
 clearText2 := RSADecrypt(secret2, prvKey)
 fmt.Println("Message =", clearText2)
 pk := prvKey
 err:= StoreRSAKey( pk ,"Pass", "Key1")
 fmt.Println(err)
 pvKey, pbKey, err :=GetRSAKey("Key1", "Pass") //rsa.PrivateKey,rsa.PublicKey, error
 secret3 := RSAEncrypt("Secret3", pbKey)
 fmt.Println("Encrypted Message =", secret3)
 clearText3 := RSADecrypt(secret3, pvKey)
 fmt.Println("Message =", clearText3)


}







