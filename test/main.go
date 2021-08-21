package main

import (
	"context"
	"encoding/json"
	"fmt"

	spvchannels "github.com/libsv/go-spvchannels"
)

func main() {
	fmt.Println("foo")

	cfg := spvchannels.ClientConfig{
		Insecure: true,
		BaseURL:  "localhost:5010",
		Version:  "v1",
		User:     "dev",
		Passwd:   "dev",
		Token:    "",
	}

	client := spvchannels.NewClient(cfg)
	r := spvchannels.GetChannelRequest{
		Method:    "GET",
		AccountId: "1",
		ChannelId: "2vkapEui-Cfb3tY7l9FFviRjpsNGa0Iv4kFEHYoMWJdl4f9PSlvurjOCnTBzH1r_C8VUuvQsn-0NsO0Q2bKGUA",
	}

	reply, err := client.GetChannel(context.Background(), r)
	if err != nil {
		panic("foo")
	}

	bReply, _ := json.Marshal(reply)
	fmt.Println(string(bReply))
}
