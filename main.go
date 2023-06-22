package main

import (
	"flag"
	"fmt"
	"log"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/machinebox/graphql"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// Default values.
	defaultAddr = "127.0.0.1"
	defaultPort = 9001

	pluginName = "etherfi"
)

var (
	addr string
	port int
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default 9091")

	flag.Parse()
}

func main() {
	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func query_gql() {
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://api.studio.thegraph.com/query/41778/etherfi-mainnet/0.0.3")

	// make a request
	req := graphql.NewRequest(`
		query {
				bids(where: { status: "WON", bidderAddress: "0x7C0576343975A1360CEb91238e7B7985B8d71BF4" }) {
					id
				}
			}
		`)
	// run it and capture the response
	var respData map[string]interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		log.Fatal(err)
	}
	bids, ok := respData["bids"]
	if !ok {
		fmt.Println("bids not found in respData")
	} else {
		for _, bid := range bids.([]interface{}) {
			bidMap := bid.(map[string]interface{})
			if idValue, ok := bidMap["id"]; ok {
				fmt.Println("ID:", idValue)
			} else {
				fmt.Println("ID not found in bid")
			}
		}
	}
	/*
		bytes, err := json.MarshalIndent(respData, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bytes))
	*/
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	// TODO: Fill here.
	//	fmt.Println("asdfasdfasdf")
	query_gql()
	ret := sdk.CallResponse{
		FuncName:   "etherfi_func",
		Message:    "YOUR_MESSAGE_CONTENTS",
		Severity:   pluginpb.SEVERITY_UNKNOWN,
		State:      pluginpb.STATE_NONE,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
