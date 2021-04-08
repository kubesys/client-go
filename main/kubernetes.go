package kubesys

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "https://119.8.188.235:6443/"
	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjcwajUwUGRjQmtucUwxWmZsa21Jb0dMVDZobmNUZHo3bmp6Mm4zcnpnYncifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudC10b2tlbi04NXFtbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWNsaWVudCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjQ5NGNiZmU5LWExMDctNGFkNi05NTg0LWI5ZDFiMDdhNjY1MSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWNsaWVudCJ9.ByrQzCJCQx9YVPOMKGITXu-amE9hT0rMaVCG-hadUGpaCDwaTrXC264hpPyuEOmTyzm3Z3Ayd_fu39XvDsse7tNFoZRG5zXtcd7fOLD2v4Y-QhmUm8Ut25h7Q3RX5XOUtvRbTKEwJ51I6CbN-Y1HtnUF7rp_k3QKjYTmDTopJQRK8jOMNi8QvEtwXvdpWaovLRsdveFTObjfJ5e2-EjlYA1sMry_JBGmNu6iofU725H8GdXjh2TbMRNQ5jakbuGEBFItPCNCdTCzHSjgMaO7JzzfiMQ142xKjKH63kFwzuGGKXDNyoQ-QJl40CHW9eLINGFXyjNUxY7kDcYEvuVZ3w")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
