```GET
	api.GET(func(mainNetPrivateKey string, tokenAmount *big.Int) Result {
    		return toResult(swap.Deposit(mainNetPrivateKey, tokenAmount))
    	}, "/deposit")

	api.Start(":8080")
}
```
```POST
   api.POST(func (param struct {
                         	MainNetPrivateKey string `json:"mainNetPrivateKey"`
                         	Base64PrivateKey  string `json:"base64PrivateKey"`
                         }) Result {
                         	err := gw.MapAccounts(param.MainNetPrivateKey, param.Base64PrivateKey)
                         	return toResult(err)
                          }
   )

   api.Start(":8080")
```
upload file
```go
func MulFile(read multipart.Reader) string {
	par, _ := read.NextPart()
	fmt.Println(par.FileName(),par.FormName())
	b, _ := ioutil.ReadAll(par)
	fmt.Println(string(b))
	return string(b)
}
GetApi().POST(MulFile, "/update")
StartServer(":8080")
```