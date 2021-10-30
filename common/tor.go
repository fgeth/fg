package common


import(
	"crypto/tls"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)



type OnionClient struct {
	client *http.Client
}


func TorDialer(base string) error{
torDialer, err := proxy.SOCKS5("tcp", MyNode.Ip, nil, proxy.Direct)
		transportConfig := &http.Transport{
		Dial:            torDialer.Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		//cookieJar, _ := cookiejar.New(nil)
		//cookieJar.SetCookies(base, Cookies)
		OC.client = &http.Client{
		Transport: transportConfig,
		Timeout: time.Second * 120,
		}
		
		//OC.client = &http.Client{
		//Transport: transportConfig,
		//Jar:       cookieJar,
		//Timeout: time.Second * 120,
		//}
		//OC.client.Get(url)
		return err

}
