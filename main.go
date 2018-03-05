package main

import (
	"fmt"
	"os"

	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	pathname := "/dev/log"
	if len(os.Args) == 2 {
		pathname = os.Args[1]
	}
	err := server.ListenUnixgram(pathname)
	if err != nil {
		panic(err)
	}
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			fmt.Printf("%v [%s/%d] %v\n",
				logParts["timestamp"],
				logParts["app_name"],
				logParts["proc_id"],
				logParts["message"],
			)
		}
	}(channel)

	server.Wait()
}
