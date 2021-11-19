package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"meta_api/protocal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var api_info_map map[string]string
var api_info_str string
var api_info_json_str string
var meta_data_send META_PROTOCAL_S

type obj_info_s struct {
	Info     map[string]string
	BaseInfo META_API_INFO_S
}

var object_info_id_obj map[string]obj_info_s
var object_info_ip_id map[string]string

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
	do := GET_query(c, "do")
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
	if do == "test_api_out" {
		data := []byte{1}
		meta_data_send := get_meta_protcal(0, 4, 3, 0, 0, data)

		outb, _ := JSON_encode(meta_data_send)
		c.String(http.StatusOK, string(outb))
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	var meta_json META_PROTOCAL_S
	if err != nil {
		fmt.Println(err)
		return
	}
	err = JSON_decode(data, &meta_json)
	if err != nil {
		fmt.Println(err)
		return
	}

	meta_data, err := base64.StdEncoding.DecodeString(meta_json.DATA_base64)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch meta_json.CmdSet {
	case 0:

		switch meta_json.CmdID {
		case 0:
			{
				if Is_response(&meta_json) {
					store_obj_info(c, meta_data)
					return
				} else {
					if meta_data[0] == 0 {
						c.String(http.StatusOK, api_info_str)
						return
					}
					if meta_data[0] == 1 {
						outb, _ := JSON_encode(api_info_map)
						c.String(http.StatusOK, string(outb))
						//c.String(http.StatusOK, api_info_json_str)
						return
					}
					if meta_data[0] == 4 {
						store_obj_info(c, meta_data[1:])
						return
					}

					check_if_need_get_info(c)
				}

				break
			}
		case 4:
			data := []byte{1}
			meta_data_send := get_meta_protcal(0, 4, 3, 0, 0, data)

			outb, _ := JSON_encode(meta_data_send)
			c.String(http.StatusOK, string(outb))

			break

		case 6:
			{
				count := (uint32)(0)
				limit := (uint32)(meta_data[0])
				offset := (uint32)(meta_data[1]) | (uint32)(meta_data[1])<<8 | (uint32)(meta_data[2])<<16 | (uint32)(meta_data[3])<<24
				res := []map[string]string{}
				for _, obj := range object_info_id_obj {

					count++
					if count < offset {
						continue
					}
					res = append(res, obj.Info)
					if count >= limit {
						outb, _ := JSON_encode(res)
						c.String(http.StatusOK, string(outb))
						break
					}
				}

				check_if_need_get_info(c)
			}
			break
		}

		break

	}

}

func main() {
	init_param()

	//1.创建路由
	r := gin.Default()
	r.Use(gin.Recovery())
	//2.绑定路由规则，执行的函数
	r.GET("/api", api_handle)
	//3.监听端口，默认8080
	r.Run(":8081")
}

func store_obj_info(c *gin.Context, meta_data []byte) {
	ip := c.ClientIP()
	_, ok := object_info_ip_id[ip]
	if ok {
		var api_info map[string]string
		err := JSON_decode(meta_data, &api_info)
		if err != nil {
			id := api_info["id"]
			object_info_ip_id[ip] = id
			object_info_obj := obj_info_s{}
			object_info_obj.Info = api_info
			object_info_id_obj[id] = object_info_obj
		}
	}

}

func get_api_info(url string, post_data string) string {

	//url := "http://" + ip + ":" + port + "/meta_api"
	//post_data = '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x8bf2df40698ba857dbdff5b460aabfe4913d3595","latest"],"id":1}'//array ("username" : "bob","key" : "12345")

	Header := map[string]string{
		"Content-Type":   "application/json charset=utf-8",
		"Content-Length": Conv_string(len(post_data)),
	}

	output := http_post_with_header(url, post_data, Header)
	return output
}
func Conv_string(i int) string {
	return fmt.Sprintf("%d", i)
}

func check_if_need_get_info(c *gin.Context) {
	ip := c.ClientIP()
	_, ok := object_info_ip_id[ip]
	if !ok {
		//request json of api_info
		data := []byte{byte(1)}
		meta_data_send := get_meta_protcal(0, 0, 3, 0, 0, data)

		outb, _ := JSON_encode(meta_data_send)
		c.String(http.StatusOK, string(outb))
		return

	}
}

type META_API_INFO_S struct {
	Meta_api_ver        string
	Id                  string
	Name                string
	Meta_api_class_name string
	Meta_api_class_id   string
	Api_info            string
	Info_url            string
	Api_url             string
	Meta_api_info_url   string
}

func Is_response(meta_json *META_PROTOCAL_S) bool {
	return (meta_json.CmdType & 0x40) != 0
}

func get_meta_protcal(CmdSet uint32, CmdID uint32, DataType uint32, CmdType uint32, SEQ uint32, data []byte) META_PROTOCAL_S {
	var meta_data META_PROTOCAL_S
	meta_data.SOF = 0xEA5A
	meta_data.Version = 16
	meta_data.DataType = DataType
	meta_data.Length = (uint32)(22 + len(data))
	meta_data.CmdType = CmdType
	meta_data.ENC = 0
	meta_data.CmdSet = CmdSet
	meta_data.CmdID = CmdID
	meta_data.Reseved = 0
	meta_data.Extend = 0
	meta_data.SEQ = SEQ
	crc16data := []byte{}
	crc16data = append(crc16data, byte(meta_data.SOF&0xFF), byte(meta_data.SOF>>8), byte(meta_data.Version), byte(meta_data.DataType), byte(meta_data.Length&0xFF), byte(meta_data.Length>>8),
		byte(meta_data.CmdType), byte(meta_data.ENC), byte(meta_data.CmdSet&0xFF), byte(meta_data.CmdSet>>8), byte(meta_data.CmdID&0xFF), byte(meta_data.CmdID>>8),
		byte(meta_data.Reseved&0xFF), byte(meta_data.Reseved>>8), byte(meta_data.SEQ&0xFF), byte(meta_data.SEQ>>8))

	crc16 := protocal.CRC16_init()
	crc16 = protocal.CRC16_update(crc16, crc16data, len(data))
	meta_data.DATA_base64 = base64.StdEncoding.EncodeToString(data)
	crc32 := protocal.CRC32_init()
	crc32data := append(crc16data, data...)
	crc32 = protocal.CRC32_update(crc32, crc32data, len(data))
	meta_data.CRC_32 = crc32
	return meta_data
}
func init_param() {
	protocal.CRC16_init_param()
	protocal.CRC32_init_param()
	api_info_map = map[string]string{
		"meta_api_ver":        "1.0",
		"id":                  "meta-api-server-id-001",
		"name":                "meta-api-server",
		"meta_api_class_name": "meta-api-server",
		"meta_api_class_id":   "meta-api-server-class-1",
		"api_info":            "https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal",
		"info_url":            "https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058",
		"api_url":             "http://42.194.159.204/api",
		"meta_api_info_url":   "https://thoughts.aliyun.com/share/61954da2c1a410001add844d#title=元宇宙_API_基础信息原语描述",
	}

	api_info_str = `meta_api_ver:` + api_info_map["meta_api_ver"] +
		`id:` + api_info_map["id"] + `,
name:` + api_info_map["name"] + `,
meta_api_class_name:` + api_info_map["meta_api_class_name"] + `,
meta_api_class_id:` + api_info_map["meta_api_class_id"] + `,
api_info:` + api_info_map["api_info"] + `,
info_url:` + api_info_map["info_url"] + `,
api_url:` + api_info_map["api_url"] + `,
api_url:` + api_info_map["api_url"]

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

type META_PROTOCAL_S struct {
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
	CRC_16      uint32 `json:"CRC-16"`
	DATA_base64 string `json:"DATA_base64"`
	CRC_32      uint32 `json:"CRC-32"`
}

func httpGet(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: time.Second * 5}

	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Fatal("Error reading body. ", err)
	}
	return string(body)
}

func http_post_with_header(url string, data string, header map[string]string) string {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	// Set headers
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Host", "httpbin.org")

	// Create and Add cookie to request
	//cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
	//req.AddCookie(&cookie)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Validate cookie and headers are attached
	//fmt.Println(req.Cookies())
	//fmt.Println(req.Header)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	return string(body)
}
