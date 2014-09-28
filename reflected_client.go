// reflected_client is an ugly client for manual testing.
// Requests must be supplied in JSON format on the commandline.
// Responses will be shown as pretty-printed JSON.

package main

import (
	"code.google.com/p/go.net/context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/scottlamb/luxor/client"
	"github.com/scottlamb/luxor/protocol"
	"os"
	"reflect"
)

var baseURL = flag.String("base_url", "http://luxor/", "Base URL for controller")
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
var typeOfController = reflect.TypeOf((*protocol.Controller)(nil)).Elem()
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

func listSubcommands() {
	fmt.Printf("Valid subcommands:\n")
	for i := 0; i < typeOfController.NumMethod(); i++ {
		fmt.Printf("    %s\n", typeOfController.Method(i).Name)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		fmt.Printf("usage: %s SUBCOMMAND [REQUEST]", os.Args[0])
		flag.PrintDefaults()
		listSubcommands()
		os.Exit(1)
	}

	ctx := context.Background()
	controller := reflect.ValueOf(&client.Controller{*baseURL})
	subcommandName := args[0]
	subcommand := controller.MethodByName(subcommandName)
	if !subcommand.IsValid() {
		fmt.Fprintf(os.Stderr, "No such subcommand %q.\n", subcommandName)
		flag.PrintDefaults()
		listSubcommands()
		os.Exit(1)
	}
	subcommandType := subcommand.Type()
	request := reflect.New(subcommandType.In(1).Elem())
	if len(args) == 2 {
		if err := json.Unmarshal([]byte(args[1]), request.Interface()); err != nil {
			panic(err)
		}
	}
	output := subcommand.Call([]reflect.Value{reflect.ValueOf(ctx), request})
	response := output[0].Interface()
	if !output[1].IsNil() {
		err := output[1].Interface().(error)
		panic(err)
	}
	formatted, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", formatted)
}
