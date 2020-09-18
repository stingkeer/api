```GET
	api.GetApi().GET(func(mainNetPrivateKey string, tokenAmount *big.Int) Result {
    		return toResult(swap.Deposit(mainNetPrivateKey, tokenAmount))
    	}, "/deposit")

	api.Start(":8080")
}
```
```POST
   api.GetApi().POST(func (param struct {
                         	MainNetPrivateKey string `json:"mainNetPrivateKey"`
                         	Base64PrivateKey  string `json:"base64PrivateKey"`
                         }) Result {
                         	err := gw.MapAccounts(param.MainNetPrivateKey, param.Base64PrivateKey)
                         	return toResult(err)
                          }
   )

   api.Start(":8080")
```