

### 1. Quick Start
```
go get go.aew.app/api.v1
```
### 2. Usage
-  helloword
```
api.GET(func() interface{} {
      return "hello word"
}, "/hello")
api.StartService(api.WithListen(":8080"))
```
- login
```
api.GET(func(username,password string) interface{} {
      return "hello word"
}, "/login")
api.StartService(api.WithListen(":8080"))
```
```shell script
curl --location --request GET 'http://127.0.0.1:8011/login?username=alloce&password=123456'
```
- upload file
```
api.POST(func(reader multipart.Reader) interface{} {
	form, _ := reader.ReadForm(0)
	fHeader := form.File["file"][0]
	f, _ := fHeader.Open()
	bytes,_:=ioutil.ReadAll(f)
	return bytes
}, "/fileUpload")
api.StartService(api.WithListen(":8080"))
```
```shell script
curl --location --request GET 'http://127.0.0.1:8011/fileUpload' \
--form 'file=@"/C:/Users/Administrator/Downloads/README.MD"'
```
-  Download file

```
api.GET(func() interface{} {
	f, _ := os.Open("C:/Users/Administrator/api/README.MD")
	return api.NewStream(f)
}, "/download")
api.StartService(api.WithListen(":8080"))
```
### 3. Swagger UI
```shell
go build -tags 'swagger'
```
[http://localhost:8080/ui/index.html#/](http://localhost:8080//ui/index.html#/)
