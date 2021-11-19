package g

import "sync"

var Object_info_map_lock sync.RWMutex

type QUERY_node_S struct {
	Object_id         string `json:"object_id"`
	Field_name        string `json:"field_name"`
	Meta_api_class_id string `json:"meta_api_class_id"`
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
