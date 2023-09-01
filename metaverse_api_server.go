package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"meta_api/f"
	"meta_api/g"
	"meta_api/protocal"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func set_api_info_map() {
	//设置里的 服务器信息 "id"，"name"，"meta_api_class_id"，"meta_api_class_name"，"api_info"，"info_url"等
	//注意 ： 请勿改动 api_url
	api_info_map = map[string]string{
		"meta_api_ver":        "1.0",
		"id":                  "id-" + g.Server_ip,
		"name":                "meta-api-server",
		"meta_api_id":         "meta-api-server-id-" + g.Server_ip,
		"meta_api_class_name": "meta-api-server",
		"meta_api_class_id":   "meta-api-server-class",
		"api_info":            "https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal",
		"info_url":            "https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058",
		"api_url":             "http://" + g.Server_ip + ":" + g.Server_port + "/api",
		"meta_api_info_url":   "https://thoughts.aliyun.com/share/61954da2c1a410001add844d#title=元宇宙_API_基础信息原语描述_Metaverse_API_schema",
	}
}

var default_meta_API_list []string
var api_info_map map[string]string
var api_info_str string
var api_info_json_str string
var meta_data_send g.META_PROTOCAL_S

type obj_info_s struct {
	Info        map[string]string
	Create_time int64
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
	//fmt.Println("process_commands data len:", len(data))
	//fmt.Println("meta_json:")
	//meta_json2, _ := f.JSON_encode(meta_json)
	//fmt.Println(string(meta_json2))
	/*for k, v := range meta_json {
		fmt.Println(k + ":" + v)
	}*/

	meta_data, err := base64.StdEncoding.DecodeString(meta_json.DATA_base64)
	if err != nil {
		fmt.Println(err)
		return
	}
	client_ip := c.ClientIP()
	fmt.Println("process_commands ip:", client_ip, ", CmdSet:", meta_json.CmdSet, ",CmdId:", meta_json.CmdID)

	switch meta_json.CmdSet {
	case 0:
		switch meta_json.CmdID {
		case 0:
			{
				f.Dump_meta_data(meta_data)

				if f.Is_response(&meta_json) {
					fmt.Println("Is_response")

					store_obj_info(c, meta_data)

					//try post my info
					f.Post_node_info(default_meta_API_list, api_info_map)
					return
				} else {

					if meta_data[0] == 0 {
						fmt.Println("get text query")
						c.String(http.StatusOK, api_info_str)
						f.Check_if_need_get_info(client_ip, &object_info_ip_id)
						return
					}
					if meta_data[0] == 1 {
						fmt.Println("get json query")
						outb, _ := f.JSON_encode(api_info_map)
						c.String(http.StatusOK, string(outb))
						//c.String(http.StatusOK, api_info_json_str)
						return
					}
					if meta_data[0] == 2 {
						fmt.Println("get xml query")
						/*
							encoder := xml.NewEncoder(f)
							err = encoder.Encode(info)
							if err != nil {
								fmt.Println("编码错误：", err.Error())
								return
							} else {
								//fmt.Println("编码成功")
								c.String(http.StatusOK, string(outb))
							}
						*/
						//c.String(http.StatusOK, api_info_json_str)
						return
					}
					if meta_data[0] == 4 {
						fmt.Println("store_obj_info")
						store_obj_info(c, meta_data[1:])
						return
					}

				}

				break
			}
		case 4:
			fmt.Println("request connect")
			data := []byte{1}
			meta_data_send := f.Get_meta_protcal(0, 4, 3, 0, 0, data)

			outb, _ := f.JSON_encode(meta_data_send)
			c.String(http.StatusOK, string(outb))

			break

		case 6:
			{
				fmt.Println("连接拓扑 遍历")
				limit := (uint32)(meta_data[0])
				offset := (uint32)(meta_data[1]) | (uint32)(meta_data[1])<<8 | (uint32)(meta_data[2])<<16 | (uint32)(meta_data[3])<<24
				Get_nodes(c, limit, offset)
				f.Check_if_need_get_info(client_ip, &object_info_ip_id)
			}
			break
		case 7:
			{
				fmt.Println("query connect")
				var query_info map[string]string
				err := f.JSON_decode(meta_data, &query_info)
				if err != nil {
					fmt.Println(err)
					return
				}
				object_id := f.Get_map_value(query_info, "object_id")
				field_name := f.Get_map_value(query_info, "field_name")
				meta_api_class_id := f.Get_map_value(query_info, "meta_api_class_id")
				if meta_api_class_id == "" {
					meta_api_class_id = f.Get_map_value(query_info, g.Field_name_map_nor_sort["meta_api_class_id"])
				}
				limit := f.Get_map_value(query_info, "limit")
				offset := f.Get_map_value(query_info, "offset")

				if object_id != "" || field_name != "" || meta_api_class_id != "" {
					Search_nodes(c, object_id, field_name, meta_api_class_id, stoi(limit), stoi(offset))
				}

				f.Check_if_need_get_info(client_ip, &object_info_ip_id)
			}
			break
		}

		break

	}
}

func website_handle(c *gin.Context) {
	html_str := strings.Replace(api_info_str, "\n", "<br>\n", -1)
	//html_str += "\n<br><a href=\"" + api_info_map[g.Field_name_map_nor_sort["api_url"]] + "?do=get_nodes&count=0&limit=20&offset=0\">" + "connected nodes list</a><br>" +
	//	"\n<br><a href=\"https://thoughts.aliyun.com/share/6195068ebdc2c4001aea0058#title=元宇宙接口标准\">" + "API docs</a><br>"
	//c.String(http.StatusOK, html_str)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":         "Metaverse standard API",
		"body":          template.HTML(html_str),
		"get_nodes_url": api_info_map[g.Field_name_map_nor_sort["api_url"]] + "?do=get_nodes&count=0&limit=255&offset=0",
	})
	return
}

func api_handle(c *gin.Context) {

	api_handle2(c)
	/*
		ci := make(chan int)

		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)

		go func(c1 *gin.Context, ci *chan int) {
			f.Try(func() {
				api_handle2(c1)
			}, func(e interface{}) {
				//print(e)
				fmt.Println("run Handle_api_php error ", e)
				//fmt.Println(string(debug.Stack()))
			})

			*ci <- 1
		}(c, &ci)

		<-ci
	*/
}
func api_handle2(c *gin.Context) {
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
	if get_meta_api_info == "getip" {
		c.String(http.StatusOK, c.ClientIP())
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
	if do == "search_nodes" {
		//http://8.222.174.114:8081/api?do=search_nodes&object_id=meta-api-server-id-8.222.174.114&field_name=&meta_api_class_id=&limit=10&offset=0
		object_id := GET_query(c, "object_id")
		field_name := GET_query(c, "field_name")
		meta_api_class_id := GET_query(c, "meta_api_class_id")
		if meta_api_class_id == "" {
			meta_api_class_id = GET_query(c, g.Field_name_map_nor_sort["meta_api_class_id"])
		}
		limit := GET_query(c, "limit")
		offset := GET_query(c, "offset")
		if object_id != "" || field_name != "" || meta_api_class_id != "" {
			Search_nodes(c, object_id, field_name, meta_api_class_id, stoi(limit), stoi(offset))
		}
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(data) > 0 {
		process_commands(c, data)
	} else {

		c.String(http.StatusOK, api_info_str)
		return
	}

}

func i64tos(i int64) string {
	return fmt.Sprintf("%d", i)
}

func stoi(s string) int {
	vid, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int(vid)
}

func filte_field_name(obj obj_info_s, field_name string, res *[]map[string]string) {
	field_name_list := []string{}

	if field_name != "" {
		field_name_list = strings.Split(field_name, ",")
	}
	//fmt.Println(field_name_list)

	if field_name == "" {
		*res = append(*res, obj.Info)
	} else if len(field_name_list) == 1 {
		field_value, ok2 := obj.Info[field_name]
		if ok2 {
			obj2 := map[string]string{"id": obj.Info["id"], field_name: field_value}
			*res = append(*res, obj2)
		}
	} else {
		obj2 := map[string]string{}
		for _, fn := range field_name_list {

			field_value, ok2 := obj.Info[fn]
			if fn == "" || !ok2 {
				continue
			} else {
				obj2["id"] = obj.Info["id"]
				obj2[fn] = field_value
			}
		}
		if len(obj2) > 0 {
			*res = append(*res, obj2)
		}
	}

}

func Search_nodes(c *gin.Context, object_id string, field_name string, meta_api_class_id string, limit int, offset int) {
	res := []map[string]string{}
	find_count := 0
	count := (0)

	g.Object_info_map_lock.RLock()
	if object_id != "" {
		obj, ok := object_info_id_obj[object_id]
		if ok {
			filte_field_name(obj, field_name, &res)
			find_count++
		}
	}

	if meta_api_class_id != "" {
		for _, obj := range object_info_id_obj {
			count++
			if count < offset {
				continue
			}
			v, ok := obj.Info["mci"]
			if ok && v == meta_api_class_id {
				filte_field_name(obj, field_name, &res)
				find_count++
			}

			if find_count >= (int)(limit) {
				break
			}
		}
	}
	count = 0
	if find_count < (int)(limit) {
		for _, obj := range object_info_id_obj {
			count++
			if count < offset {
				continue
			}
			v, ok := obj.Info["mci"]
			//fmt.Println("Search_nodes 1 " + v + " " + meta_api_class_id)
			if ok && strings.Contains(v, meta_api_class_id) {
				//fmt.Println("Search_nodes 2")
				filte_field_name(obj, field_name, &res)
				find_count++
			}

			if find_count >= (int)(limit) {
				break
			}
		}
	}
	g.Object_info_map_lock.RUnlock()
	outb, _ := f.JSON_encode(res)
	c.String(http.StatusOK, string(outb))
}

func Get_nodes(c *gin.Context, limit uint32, offset uint32) {
	res := []map[string]string{}
	find_count := 0
	count := (uint32)(0)
	g.Object_info_map_lock.RLock()
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
	g.Object_info_map_lock.RUnlock()
	outb, _ := f.JSON_encode(res)
	c.String(http.StatusOK, string(outb))
}

func main() {
	init_param()
	go cron_jobs()
	//1.创建路由
	r := gin.Default()
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("./html/*")
	//2.绑定路由规则，执行的函数
	r.GET("/api", api_handle)
	r.POST("/api", api_handle)
	r.GET("/", website_handle)
	//3.监听端口，默认8080
	r.Run(":" + g.Server_port)
}

func cron_jobs() {
	for {
		select {
		case <-time.After(time.Second * 600):
			{
				fmt.Println("saving object_info_id_obj.json")
				f.Try(func() {
					meta_json2, _ := f.JSON_encode(object_info_id_obj)
					f.Save_file("object_info_id_obj.json", string(meta_json2))
					meta_json3, _ := f.JSON_encode(object_info_ip_id)
					f.Save_file("object_info_ip_id.json", string(meta_json3))
				}, func(e interface{}) {
					fmt.Println(e)
				})
			}
		case <-time.After(time.Second * 2 * 600):
			{
				fmt.Println("remove time out id")
				f.Try(func() {
					for k, v := range object_info_id_obj {
						if v.Create_time > time.Now().Unix()-int64(86400*60) {
							g.Object_info_map_lock.Lock()
							delete(object_info_id_obj, k)
							g.Object_info_map_lock.Unlock()
						}
					}
				}, func(e interface{}) {
					fmt.Println(e)
				})
			}

		case <-time.After(time.Second * 3600):
			{
				fmt.Println("Post_node_info:reg ip")
				f.Try(func() {
					f.Post_node_info(default_meta_API_list, api_info_map)
				}, func(e interface{}) {
					fmt.Println(e)
				})
			}
		}
	}
}

func store_obj_info(c *gin.Context, meta_data []byte) {
	ip := c.ClientIP()
	fmt.Println("store_obj_info ip:" + ip)
	g.Object_info_map_lock.RLock()
	_, ok := object_info_ip_id[ip]
	g.Object_info_map_lock.RUnlock()
	if !ok {
		var api_info map[string]string
		fmt.Println("meta_data:" + string(meta_data))
		err := f.JSON_decode(meta_data, &api_info)
		if err == nil {
			api_info["id"] = g.Server_ip + ":" + api_info["id"]
			id := api_info["id"]
			object_info_obj := obj_info_s{}
			api_info = g.Reset_field_name_map_nor_sort(&api_info)
			api_info["rtsm"] = i64tos(time.Now().UnixNano() / 1e6)
			object_info_obj.Info = api_info
			object_info_obj.Create_time = time.Now().Unix()
			g.Object_info_map_lock.Lock()
			object_info_ip_id[ip] = id
			object_info_id_obj[id] = object_info_obj
			g.Object_info_map_lock.Unlock()
			fmt.Println("ip:" + ip + ",id:" + id + " stored")

		} else {
			fmt.Println("store_obj_info error", err)
		}
	} else {
		fmt.Println("ip already exist in object_info_ip_id ,ip:" + ip)

	}

}

func init_param() {
	protocal.CRC16_init_param()
	protocal.CRC32_init_param()
	g.Get_field_name_map_sort_nor()
	g.Server_port = "8081"

	default_meta_API_list = []string{"http://8.222.174.114:8081/api"}
	var api_list = default_meta_API_list

	g.Server_ip = f.GetPulicIP2()
	if !f.IsPublicIP(g.Server_ip) {
		g.Server_ip = f.HttpGet(api_list[0] + "?do=getip")
		if !f.IsPublicIP(g.Server_ip) {
			g.Server_ip = f.GetPulicIP()
		}
	}

	fmt.Println("Metaverse API server start !\n\nServer info:\n")
	fmt.Println("g.Server_ip:" + g.Server_ip)
	fmt.Println("g.Server_port:" + g.Server_port)

	set_api_info_map()
	/*
	   	api_info_str = `meta_api_ver:` + api_info_map["meta_api_ver"] + `,
	   id:` + api_info_map["id"] + `,
	   name:` + api_info_map["name"] + `,
	   meta_api_class_name:` + api_info_map["meta_api_class_name"] + `,
	   meta_api_class_id:` + api_info_map["meta_api_class_id"] + `,
	   api_info:` + api_info_map["api_info"] + `,
	   info_url:` + api_info_map["info_url"] + `,
	   api_url:` + api_info_map["api_url"] + `,
	   meta_api_info_url:` + api_info_map["meta_api_info_url"]
	*/
	//fmt.Println(api_info_map)
	for k, v := range api_info_map {
		api_info_str += k + ":" + v + ",\n"
	}
	api_info_map = g.Reset_field_name_map_nor_sort(&api_info_map)
	api_info_map["rtsm"] = i64tos(time.Now().UnixNano() / 1e6)

	fmt.Println(api_info_str)

	//post_my_info

	g.Object_info_map_lock = sync.RWMutex{}

	f.Post_node_info(api_list, api_info_map)

	object_info_ip_id = make(map[string]string)
	object_info_id_obj = make(map[string]obj_info_s)

	object_info_ip_id_d := f.File_read("object_info_ip_id.json")
	if len(object_info_ip_id_d) > 0 {
		object_info_ip_id_t := make(map[string]string)
		err := f.JSON_decode(object_info_ip_id_d, &object_info_ip_id_t)
		if err == nil {
			fmt.Println("object_info_ip_id get sucess")
			object_info_ip_id = object_info_ip_id_t
		}
	}
	object_info_id_obj_d := f.File_read("object_info_id_obj.json")
	if len(object_info_id_obj_d) > 0 {
		object_info_id_obj_t := make(map[string]obj_info_s)
		err := f.JSON_decode(object_info_id_obj_d, &object_info_id_obj_t)
		if err == nil {
			fmt.Println("object_info_id_obj get sucess")
			object_info_id_obj = object_info_id_obj_t
		}
	}

	g.Object_info_map_lock.Lock()

	id := api_info_map["id"]

	object_info_obj := obj_info_s{}
	object_info_obj.Info = api_info_map

	object_info_ip_id[g.Server_ip] = id
	object_info_id_obj[id] = object_info_obj
	//backup
	//object_info_id_obj[id+"-1"] = object_info_obj

	g.Object_info_map_lock.Unlock()

}

func GET_query(c *gin.Context, key string) string {
	return c.DefaultQuery(key, "")
}
