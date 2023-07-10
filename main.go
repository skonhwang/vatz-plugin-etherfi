package main

import (
	"flag"
	"fmt"
	"os"

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
	graphUrl    = "https://api.studio.thegraph.com/query/41778/etherfi-mainnet/0.0.3"
	outputPath  = "/mnt/etherfi/sync_client_validator_keys/"

	pluginName = "etherfi"
)

var (
	addr    string
	port    int
	gqlUrl  string
	outPath string
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
	flag.IntVar(&port, "port", defaultPort, "Port number, default 9091")
	flag.StringVar(&gqlUrl, "graphURL", graphUrl, "Ex:https://api.studio.thegraph.com/query/41778/etherfi-mainnet/0.0.3")
	flag.StringVar(&outPath, "output", outputPath, "Dirctory for downloading validator keys")

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

func query_gql() (string, bool) {
	var retString string
	var retBool bool
	// create a client (safe to share across requests)
	//client := graphql.NewClient("https://api.studio.thegraph.com/query/41778/etherfi-mainnet/0.0.3")
	client := graphql.NewClient(gqlUrl)

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
		fmt.Println("Error: ", err)
	}
	bids, ok := respData["bids"]
	if !ok {
		fmt.Println("bids not found in respData")
	} else {
		for _, bid := range bids.([]interface{}) {
			bidMap := bid.(map[string]interface{})
			if idValue, ok := bidMap["id"]; ok {
				//path := fmt.Sprintf("/Users/hwangseungkon/dsrv/2022/etherfi/sync-client-v2/etherfi-sync-clientv2/mnt/etherfi/sync_client_validator_keys/%s", idValue.(string))
				path := fmt.Sprintf("%s/%s", outPath, idValue.(string))
				_, errDir := os.Stat(path)
				if os.IsNotExist(errDir) {
					retString += fmt.Sprintf("%s is new one\n", idValue.(string))
					fmt.Println(retString)
					retBool = true
				} else {
					logstr := fmt.Sprintf("%s is existed!\n", path)
					fmt.Println(logstr)
					retBool = false
				}
			} else {
				retString = "ID not found in bid"
				retBool = false
			}
		}
	}

	fmt.Println(retString)
	return retString, retBool
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	severity := pluginpb.SEVERITY_INFO
	state := pluginpb.STATE_NONE
	// TODO: Fill here.
	str, isNew := query_gql()
	if isNew {
		severity = pluginpb.SEVERITY_CRITICAL
		state = pluginpb.STATE_SUCCESS
		fmt.Println(str)
	} else {
		str = "There is no new one."
		severity = pluginpb.SEVERITY_INFO
		state = pluginpb.STATE_SUCCESS
		fmt.Println(str)
	}
	ret := sdk.CallResponse{
		FuncName:   "etherfi_func",
		Message:    str,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
