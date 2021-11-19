package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"meta_api/f"
	"meta_api/g"
	"meta_api/protocal"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var default_meta_API_ip []string

var api_info_map map[string]string
var api_info_str string
var api_info_json_str string
var meta_data_send g.META_PROTOCAL_S

type obj_info_s struct {
	Info     map[string]string
	BaseInfo g.META_API_INFO_S
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

func process_commands(c *gin.Context, data []byte) {

	var meta_json g.META_PROTOCAL_S

	err := f.JSON_decode(data, &meta_json)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(meta_json)
	/*for k, v := range meta_json {
		fmt.Println(k + ":" + v)
	}*/

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
				f.Dump_meta_data(meta_data)
				if f.Is_response(&meta_json) {
					store_obj_info(c, meta_data)

					//try post my info
					api_list := []string{c.ClientIP()}
					for i := 0; i < len(default_meta_API_ip); i++ {
						api_list[i] = "http://" + default_meta_API_ip[i] + ":8081/api"
					}
					f.Post_node_info(api_list, api_info_map)
					return
				} else {
					if meta_data[0] == 0 {
						c.String(http.StatusOK, api_info_str)
						return
					}
					if meta_data[0] == 1 {
						outb, _ := f.JSON_encode(api_info_map)
						c.String(http.StatusOK, string(outb))
						//c.String(http.StatusOK, api_info_json_str)
						return
					}
					if meta_data[0] == 4 {
						store_obj_info(c, meta_data[1:])
						return
					}

					f.Check_if_need_get_info(c, &object_info_ip_id)
				}

				break
			}
		case 4:
			fmt.Println("nodes post self info")
			data := []byte{1}
			meta_data_send := f.Get_meta_protcal(0, 4, 3, 0, 0, data)

			outb, _ := f.JSON_encode(meta_data_send)
			c.String(http.StatusOK, string(outb))

			break

		case 6:
			{

				limit := (uint32)(meta_data[0])
				offset := (uint32)(meta_data[1]) | (uint32)(meta_data[1])<<8 | (uint32)(meta_data[2])<<16 | (uint32)(meta_data[3])<<24
				Get_nodes(c, limit, offset)
				f.Check_if_need_get_info(c, &object_info_ip_id)
			}
			break
		}

		break

	}
}

func api_handle(c *gin.Context) {
	get_meta_api_info := GET_query(c, "get_meta_api_info")
	do := GET_query(c, "do")
	if get_meta_api_info == "text" {
		c.String(http.StatusOK, api_info_str)
		return
	}
	if get_meta_api_info == "json" {
		outb, _ := f.JSON_encode(api_info_map)
		c.String(http.StatusOK, string(outb))
		//c.String(http.StatusOK, api_info_json_str)
		return
	}
	if do == "test_api_out" {
		data := []byte{1}
		meta_data_send := f.Get_meta_protcal(0, 4, 3, 0, 0, data)

		outb, _ := f.JSON_encode(meta_data_send)
		c.String(http.StatusOK, string(outb))
		return
	}

	if do == "get_nodes" {
		limit := GET_query(c, "limit")
		offset := GET_query(c, "offset")

		Get_nodes(c, (uint32)(stoi(limit)), (uint32)(stoi(offset)))
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	process_commands(c, data)

}

func stoi(s string) int {
	vid, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int(vid)
}

func Get_nodes(c *gin.Context, limit uint32, offset uint32) {
	res := []map[string]string{}
	find_count := 0
	count := (uint32)(0)
	for _, obj := range object_info_id_obj {
		count++
		if count < offset {
			continue
		}
		res = append(res, obj.Info)
		find_count++
		if find_count >= (int)(limit) {
			break
		}
	}
	outb, _ := f.JSON_encode(res)
	c.String(http.StatusOK, string(outb))
}

func main() {
	init_param()

	//1.创建路由
	r := gin.Default()
	r.Use(gin.Recovery())
	//2.绑定路由规则，执行的函数
	r.GET("/api", api_handle)
	r.POST("/api", api_handle)
	//3.监听端口，默认8080
	r.Run(":8081")
}

func get_external() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	return string(content)
}

func store_obj_info(c *gin.Context, meta_data []byte) {
	ip := c.ClientIP()
	_, ok := object_info_ip_id[ip]
	if !ok {
		var api_info map[string]string
		err := f.JSON_decode(meta_data, &api_info)
		if err != nil {
			id := api_info["id"]
			object_info_ip_id[ip] = id
			object_info_obj := obj_info_s{}
			object_info_obj.Info = api_info
			object_info_id_obj[id] = object_info_obj
		}
	}

}

func init_param() {
	protocal.CRC16_init_param()
	protocal.CRC32_init_param()

	default_meta_API_ip = []string{"42.194.159.204"}

	ip := f.GetPulicIP()
	fmt.Println("ip:" + ip)

	api_info_map = map[string]string{
		"meta_api_ver":        "1.0",
		"id":                  "meta-api-server-id-" + ip,
		"name":                "meta-api-server-" + ip,
		"meta_api_class_name": "meta-api-server-" + ip,
		"meta_api_class_id":   "meta-api-server-class-1-" + ip,
		"api_info":            "https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal",
		"info_url":            "https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058",
		"api_url":             "http://" + ip + ":8081/api",
		"meta_api_info_url":   "https://thoughts.aliyun.com/share/61954da2c1a410001add844d#title=元宇宙_API_基础信息原语描述",
	}

	//post_my_info
	api_list := []string{}
	for i := 0; i < len(default_meta_API_ip); i++ {
		api_list = append(api_list, "http://"+default_meta_API_ip[i]+":8081/api")
	}
	f.Post_node_info(api_list, api_info_map)

	object_info_ip_id = make(map[string]string)
	object_info_id_obj = make(map[string]obj_info_s)

	//
	_, ok := object_info_ip_id[ip]
	if !ok {

		id := api_info_map["id"]
		object_info_ip_id[ip] = id
		object_info_obj := obj_info_s{}
		object_info_obj.Info = api_info_map
		object_info_id_obj[id] = object_info_obj
		//backup
		object_info_id_obj[id+"-1"] = object_info_obj
	}

	api_info_str = `meta_api_ver:` + api_info_map["meta_api_ver"] +
		`id:` + api_info_map["id"] + `,
name:` + api_info_map["name"] + `,
meta_api_class_name:` + api_info_map["meta_api_class_name"] + `,
meta_api_class_id:` + api_info_map["meta_api_class_id"] + `,
api_info:` + api_info_map["api_info"] + `,
info_url:` + api_info_map["info_url"] + `,	
api_url:` + api_info_map["api_url"] + `,
meta_api_info_url:` + api_info_map["meta_api_info_url"]

	api_info_json_str = `
{"id":"meta-api-server-id-001",
"name":"meta-api-server",
"meta_api_class_name":"meta-api-server",
"meta_api_class_id":"meta-api-server-class-1",
"api_info":"https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal",
"info_url":"https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058",
"api_url":"http://42.194.159.204:8081/api",
"meta_api_ver":1.0}
`
}

func GET_query(c *gin.Context, key string) string {
	return c.DefaultQuery(key, "")
}
