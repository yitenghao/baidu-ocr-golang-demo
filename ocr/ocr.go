package ocr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const tokenurl = "https://aip.baidubce.com/oauth/2.0/token"

// const ocrurl = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic" //基础ocr
const ocrurl = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic" //高精度
// const ocrurl = "https://aip.baidubce.com/rest/2.0/ocr/v1/general"  //含位置信息
// const ocrurl = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate"  //高精度含位置信息
// const ocrurl = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_enhanced"  //含生僻字版

const granttype = "client_credentials"

const clientid = "your clientid"
const clientsecret = "your clientsecret"

var token string

func init() {
	// CreateImgs()
	ok, data := GetToken()
	fmt.Println(ok, data)
	if ok != nil {
		fmt.Println(ok)
		return
	}
	token = data.Token
}

// func main() {
// 	ff, err := ioutil.ReadFile("./imgs/222.png") //快速读文件
// 	if err != nil {
// 		return
// 	}
// 	// bufstore := make([]byte, 4*len(ff))         //数据缓存
// 	// base64.StdEncoding.Encode(bufstore, ff)     // 文件转base64
// 	bufstore := base64.StdEncoding.EncodeToString(ff)
// 	// _ = ioutil.WriteFile("./imgs/base64pic.txt", []byte(bufstore), 0666)
// 	// _ = ioutil.WriteFile("./imgs/222.png", ff, 0666)
// 	b, value := YunOcr(token, bufstore)
// 	if b != nil {
// 		fmt.Println(b)
// 		return
// 	}
// 	fmt.Println(value)
// }
type ReturnData struct {
	Log_id           int64
	Words_result     []map[string]string
	Words_result_num int
}
type ResponseData struct {
	Code  string
	Msg   string
	Token string
}

//OcrImg 调用百度文字识别api
func OcrImg(base string) (data ReturnData, err error) {
	if strings.HasPrefix(base, "data:image") {
		var dt ReturnData
		return dt, errors.New("param format error")
	}
	err, data = YunOcr(token, base)
	return
}

// func CreateImgs() {
// 	err := os.Mkdir("imgs", os.ModePerm) //创建目录
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		os.Chmod("imgs", 0777)
// 	}
// }

func GetToken() (error, ResponseData) {
	response := ResponseData{}
	//读取token文件
	file, err := os.OpenFile("./ocr.tk", os.O_RDWR, 0600)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		file, createrr := os.Create("./ocr.tk")
		if createrr != nil {
			fmt.Println("创建token文件失败")
			return createrr, response
		}
		// t := time.Now()
		// millisecond := t.UnixNano() / 1e6
		// fmt.Println(millisecond)
		rterr, response := authtoken()
		if rterr == nil {
			_, err = file.Write([]byte(response.Token))
		}
		return rterr, response
	}

	//校验token文件的有效期
	fileinfo, _ := file.Stat()
	changetime := fileinfo.ModTime()
	hours := time.Now().Sub(changetime).Hours()
	fmt.Println(hours)
	//大于29天就重新获取 实际有效期是30天
	if int(hours) > 29*24 {
		rterr, response := authtoken()
		if rterr == nil {
			_, err = file.Write([]byte(response.Token))
			if err != nil {
				fmt.Println(err)
				return err, response
			}
		}
		return rterr, response
	}
	buff := make([]byte, 1000)
	n, err := file.Read(buff)
	fmt.Println(string(buff[:n]))
	if err != nil {
		return nil, response
	}
	response.Token = string(buff[:n])
	return err, response
}
func authtoken() (error, ResponseData) {
	data := url.Values{"grant_type": {granttype}, "client_id": {clientid}, "client_secret": {clientsecret}}
	respmsg, msgerr := http.PostForm(tokenurl, data)
	// fmt.Println(respmsg.Request.Form, respmsg.Request.Method, respmsg.Request.URL)
	response := ResponseData{}
	if msgerr != nil {
		response.Code = "400"
		response.Msg = "token请求失败"
		response.Token = ""
		return msgerr, response
	}
	defer respmsg.Body.Close()
	resbody, reserr := ioutil.ReadAll(respmsg.Body)
	if reserr != nil {
		// handle error
		response.Code = "400"
		response.Msg = "token读取失败"
		response.Token = ""
		return reserr, response
	}
	var rtdata struct {
		Access_token string
		Expires_in   int
	}
	err := json.Unmarshal(resbody, &rtdata)
	if err != nil {
		fmt.Println(err)
		response.Code = "400"
		response.Msg = "token解析失败"
		response.Token = ""
		return err, response
	}
	response.Code = ""
	response.Msg = "token读取成功"
	response.Token = rtdata.Access_token
	return nil, response
}

func YunOcr(token, basestr string) (error, ReturnData) {
	sendurl := ocrurl + "?access_token=" + token
	data := url.Values{"image": {basestr}}
	var client = http.DefaultClient
	req, err := http.NewRequest("POST", sendurl, strings.NewReader(data.Encode()))
	// bufstore := make([]byte, len([]byte(basestr))/4+1) //数据缓存
	// base64.StdEncoding.Decode(bufstore, []byte(basestr)) // base64转文件
	// _ = ioutil.WriteFile("./imgs/222.txt", []byte(basestr), 0666)
	var returndata ReturnData
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// fmt.Println(req.Method, req.URL, req.Header)
	rp, err := client.Do(req)
	if err != nil {
		return err, returndata
	}
	defer rp.Body.Close()
	resbody, reserr := ioutil.ReadAll(rp.Body)
	if reserr != nil {
		return reserr, returndata
	}

	err = json.Unmarshal(resbody, &returndata)
	if err != nil {
		fmt.Println("无法解析")
	}
	return err, returndata
}
