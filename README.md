# baidu-ocr-golang-demo
go语言调百度ocr api
## Config
在ocr.go文件里面配置百度智能云( https://login.bce.baidu.com/ )里面获得的clientid和clientsecret
## Run
第一次执行时会调接口生成token并存在ocr.tk文件里面
以后每次调用会校验有效时间(29天)看是否失效，失效则重新获取
由于前端需要localhost或https下才能获取摄像头权限(首次获取权限浏览器会询问是否给权限，以后不会了)，所以调试时请使用 http://localhost:1415 来测试
