package f

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"meta_api/g"
	"meta_api/protocal"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func Is_response(meta_json *g.META_PROTOCAL_S) bool {
	return (meta_json.CmdType & 0x40) != 0
}

func Get_meta_protcal(CmdSet uint32, CmdID uint32, DataType uint32, CmdType uint32, SEQ uint32, data []byte) g.META_PROTOCAL_S {
	var meta_data g.META_PROTOCAL_S
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
	crc16 = protocal.CRC16_update(crc16, crc16data, len(crc16data))
	meta_data.DATA_base64 = base64.StdEncoding.EncodeToString(data)
	crc32 := protocal.CRC32_init()
	crc32data := append(crc16data, data...)
	crc32 = protocal.CRC32_update(crc32, crc32data, len(crc32data))
	meta_data.CRC_32 = crc32
	return meta_data
}

func JSON_decode(data []byte, val interface{}) error {
	return json.Unmarshal(data, val)
}

func JSON_encode(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

func Post_node_info(api_list []string, api_info_map map[string]string) {
	Try(func() {

		my_info_byte, _ := JSON_encode(api_info_map)
		data := []byte{byte(4)}
		data = append(data, my_info_byte...)
		meta_data_send := Get_meta_protcal(0, 0, 3, 0, 0, data)
		outb, _ := JSON_encode(meta_data_send)

		for i := 0; i < len(api_list); i++ {

			Post_api_info(api_list[i], string(outb))
			fmt.Println(string(outb))
		}
	}, func(e interface{}) {

	})

}

func Dump_meta_data(meta_data []byte) {
	count := 0
	fmt.Println("dump_meta_data:")
	for _, v := range meta_data {
		count++
		fmt.Printf("0x%x,", v)
		if count > 22 {
			return
		}
	}
}

func GetPulicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	fmt.Println("GetPulicIP:" + localAddr)
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func Post_api_info(url string, post_data string) string {

	//url := "http://" + ip + ":" + port + "/meta_api"
	//post_data = '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x8bf2df40698ba857dbdff5b460aabfe4913d3595","latest"],"id":1}'//array ("username" : "bob","key" : "12345")

	Header := map[string]string{
		"Content-Type":   "application/json charset=utf-8",
		"Content-Length": Conv_string(len(post_data)),
	}

	output := Http_post_with_header(url, post_data, Header)
	return output
}
func Conv_string(i int) string {
	return fmt.Sprintf("%d", i)
}

func Check_if_need_get_info(c *gin.Context, object_info_ip_id *map[string]string) bool {
	ip := c.ClientIP()
	_, ok := (*object_info_ip_id)[ip]
	if !ok {
		//request json of api_info
		data := []byte{byte(1)}
		meta_data_send := Get_meta_protcal(0, 0, 3, 0, 0, data)

		outb, _ := JSON_encode(meta_data_send)
		c.String(http.StatusOK, string(outb))
		return true

	}
	return false
}

func Http_post_with_header(url string, data string, header map[string]string) string {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		//log.Fatal("Error reading request. ", err)
		return ""
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
		//log.Fatal("Error reading response. ", err)
		return ""
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Fatal("Error reading body. ", err)
		return ""
	}
	return string(body)
}

func HttpGet(url string) string {
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
