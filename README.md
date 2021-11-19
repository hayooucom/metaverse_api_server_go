# metaverse_api_server_go

#### 介绍
元宇宙接口 API 服务器 实现。 

The implementation of  [元宇宙统一接口标准](https://thoughts.aliyun.com/share/61953ed66a1d11001aecd4f9#title=元宇宙通用通信协议_Metaverse_General_Protocal)

说明：
json 所有接口，字段通常为string类型

元宇宙基础API接口（可通过此服务器获取其他节点信息，遍历元宇宙）：

http://42.194.159.204:8081/api

服务器开源代码：

https://gitee.com/hayoou/metaverse_api_server_go


API测试：

http://42.194.159.204:8081/api?do=get_nodes&limit=10&offset=0


http://42.194.159.204:8081/api?do=search_nodes&object_id=&field_name=&meta_api_class_id=meta-api-server&limit=10&offset=0


安装 golang(https://golang.org/)
 
安装：设置服务器 8081 端口开放

windows:

下载后，按 Shift 加右键，在当前路径打开你的终端并执行

$ go env -w GO111MODULE=on

$ go env -w GOPROXY=https://goproxy.cn,direct

$ go mod tidy

$ go build metaverse_api_server.go 

双击生成的exe文件 或者 直接运行：

$ go run metaverse_api_server.go



linux:

设置goproxy

$ export GO111MODULE=on

$ export GOPROXY=https://goproxy.cn

$ go mod tidy

$ chmod 777 ./buildandrun.sh

$ ./buildandrun.sh