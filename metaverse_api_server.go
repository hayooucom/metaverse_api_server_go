package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var api_info_map map[string]string
var api_info_str string
var api_info_json_str string

/*
func err_handle(err, c *gin.Context) {
	if err {
		c.String(http.StatusOK, err)
		fmt.Println(err)
		return
	}
}*/

func api_handle(c *gin.Context) {
	get_meta_api_info := GET_query(c, "get_meta_api_info")
	if get_meta_api_info == "text" {
		c.String(http.StatusOK, api_info_str)
		return
	}
	if get_meta_api_info == "json" {
		outb, _ := JSON_encode(api_info_map)
		c.String(http.StatusOK, string(outb))
		//c.String(http.StatusOK, api_info_json_str)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	var meat_json META_PROTOCAL_JSON
	if err != nil {
		fmt.Println(err)
		return
	}
	err = JSON_decode(data, &meat_json)
	if err != nil {
		fmt.Println(err)
		return
	}

	meat_data, err := base64.StdEncoding.DecodeString(meat_json.DATA_base64)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch meat_json.CmdSet {
	case 0:

		switch meat_json.CmdID {
		case 0:
			{
				if meat_data[0] == 0 {
					c.String(http.StatusOK, api_info_str)
					return
				}
				if meat_data[0] == 1 {
					c.String(http.StatusOK, api_info_json_str)
					return
				}
				break
			}
		}
		break

	}

}

func main() {
	init_param()
	//1.创建路由
	r := gin.Default()
	//2.绑定路由规则，执行的函数
	r.GET("/api", api_handle)
	//3.监听端口，默认8080
	r.Run(":8081")
}

func init_param() {

	api_info_map = map[string]string{
		"id":                  "meta-api-server-id-001",
		"name":                "meta-api-server",
		"meta_api_class_name": "meta-api-server",
		"meta_api_class_id":   "meta-api-server-class-1",
		"api_info":            "https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal",
		"info_url":            "https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058",
		"api_url":             "http://42.194.159.204/api",
		"meta_api_ver":        "1.0",
	}

	api_info_str = `
id:` + api_info_map["id"] + `,
name:` + api_info_map["name"] + `,
meta_api_class_name:` + api_info_map["meta_api_class_name"] + `,
meta_api_class_id:` + api_info_map["meta_api_class_id"] + `,
api_info:` + api_info_map["id"] + `,
info_url:` + api_info_map["id"] + `,
api_url:` + api_info_map["id"] + `,
meta_api_ver:` + api_info_map["id"]

	api_info_json_str = `
{"id":"meta-api-server-id-001",
"name":"meta-api-server",
"meta_api_class_name":"meta-api-server",
"meta_api_class_id":"meta-api-server-class-1",
"api_info":"https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal",
"info_url":"https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058",
"api_url":"http://42.194.159.204/api",
"meta_api_ver":1.0}
`
}

func GET_query(c *gin.Context, key string) string {
	return c.DefaultQuery(key, "")
}

func JSON_decode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

func JSON_encode(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

type META_PROTOCAL_JSON struct {
	SOF         uint32 `json:"SOF"`
	Version     uint32 `json:"Version"`
	DataType    uint32 `json:"DataType"`
	Length      uint32 `json:"Length"`
	CmdType     uint32 `json:"CmdType"`
	ENC         uint32 `json:"ENC"`
	CmdSet      uint32 `json:"CmdSet"`
	CmdID       uint32 `json:"CmdID"`
	Reseved     uint32 `json:"Reseved"`
	Extend      uint32 `json:"Extend"`
	SEQ         uint32 `json:"SEQ"`
	DATA_base64 string `json:"DATA_base64"`
	CRC_32      uint32 `json:"CRC-32"`
}
