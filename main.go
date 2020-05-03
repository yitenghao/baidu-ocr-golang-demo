package main

import (
	"baidu-ocr-golang-demo/ocr"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ResponseData struct {
	Code string
	Msg  string
	Data []string
}

//必须在localhost或https下前端才能获取到摄像头权限
func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/ocr", OcrImg)
	log.Fatal(http.ListenAndServe(":1415", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

//这里可以拓展一些其他参数，如识别等级
type Param struct {
	Base string `json:"base`
}

func OcrImg(w http.ResponseWriter, r *http.Request) {
	var resp ResponseData

	rd := r.Body
	body, err := ioutil.ReadAll(rd)
	if err != nil {
		fmt.Println(err)
		return
	}
	var param Param
	jsonunmerr := json.Unmarshal(body, &param)
	if jsonunmerr != nil {
		fmt.Println(jsonunmerr)
		resp.Code = "400"
		resp.Msg = "参数解析失败"
		bts, _ := json.Marshal(resp)
		w.Write([]byte(bts))
		w.WriteHeader(400)
		return
	}
	imgs := []string{"data:image/png;base64,", "data:image/jpg;base64,", "data:image/jpeg;base64,", "data:image/bmp;base64,"}
	base := param.Base
	for _, imghead := range imgs {
		if strings.HasPrefix(param.Base, imghead) {
			base = param.Base[len(imghead):]
		}
	}
	if base != "" {
		data, err := ocr.OcrImg(base)
		if err == nil {
			var list []string
			for _, item := range data.Words_result {
				list = append(list, item["words"])
			}
			fmt.Println(list)
			resp.Code = "200"
			resp.Msg = "success"
			resp.Data = list
			bts, _ := json.Marshal(resp)

			w.Write(bts)
			return
		}
	}
	resp.Code = "400"
	resp.Msg = "识别失败"
	bts, _ := json.Marshal(resp)
	w.Write([]byte(bts))
	return

}
