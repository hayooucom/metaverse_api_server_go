package g

import "sync"

var Object_info_map_lock sync.RWMutex
var Server_ip string
var Server_port string = "8081"

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

var Sort_field_name_map map[string]string

func Get_sort_field_name_map() {
	Sort_field_name_map = map[string]string{
		"id": "id", "name": "nm", "describe": "desc", "type": "tp", "interact_expired_date": "ied",
		"visible_expired_date": "ved", "class_name": "cn", "meta_api_class_name": "mcn",
		"meta_api_class_id": "mci", "info_url": "iu", "info_json_url": "iju", "api_info": "ai",
		"api_url": "au", "api_url2": "au2", "meta_api_sch_url": "msu", "meta_api_info_url": "miu", "meta_api_ver": "mv", "meta_api_ver_min": "mvmin", "support_connector": "sc", "is_permanent": "isper", "is_in_real_world": "isirw", "is_AI": "isai", "status": "stat", "failure_report_api": "fra", "trusted_connector": "tc", "connector_count": "cc", "info_expire": "ie", "transmit_rate_per_interface_bps": "trpi", "transmit_state_pressure": "tsp", "transmit_state_pressure_percent": "tspp", "transmit_rate_remaining_bandwidth_bps": "trrbb", "can_relay": "cr", "network_access": "na", "network_access_internet": "nai", "network_info": "ni", "debug_info": "di", "debug_state": "ds", "manufactor": "mnf", "create_date": "cd", "can_delete": "cdel", "can_interact": "ci", "text_stream": "ts", "icon_url": "iu", "avatar_hd": "ah", "3D_model_type": "3dt", "3D_model_url": "3mu", "3D_texture_url": "3tu", "3D_model_size": "3ms", "position": "pos", "3D_spin": "3s", "3D_zoom": "3z", "video_output_stream": "vos", "video_texture_zoom": "vtz", "3D_model_controler_stream": "3mcs", "audio_output_stream": "aos", "click_url": "cu", "interact_msg": "im", "motion_trajectory_stream": "mts", "video_input_stream": "vis", "audio_input_stream": "ais", "permanent_object_url": "pou", "unique_id": "ui", "create_api_url": "cau", "pay_api_url": "pau", "interact_api_url": "iau", "meta_api_sch_extend_url_list": "mseul", "meta_api_sch_custum_url_list": "macul", "support_api_stream_type": "spst", "api_interface_data": "aid", "mata_api_public_key": "mapk", "get_meta_api_info": "gmai"}
	for k, v := range Sort_field_name_map {
		Sort_field_name_map[v] = k
	}
}
