killall metaverse_api_server
go build ./metaverse_api_server.go
 nohup ./metaverse_api_server > /dev/null 2>&1 &
