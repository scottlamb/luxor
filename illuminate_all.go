// illuminate_all is a simple example of using the luxor API.
// It retrieves all themes and illuminates them one-by-one.

package main

import (
	"code.google.com/p/go.net/context"
	"github.com/scottlamb/luxor/client"
	"github.com/scottlamb/luxor/protocol"
)

func main() {
	client := &client.Controller{BaseURL: "http://luxor/"}
	ctx := context.Background()
	themes, err := client.ThemeListGet(ctx, &protocol.ThemeListGetRequest{})
	if err != nil {
		panic(err)
	}
	for _, theme := range themes.ThemeList {
		request := &protocol.IlluminateThemeRequest{ThemeIndex: theme.ThemeIndex, OnOff: 1}
		_, err := client.IlluminateTheme(context.Background(), request)
		if err != nil {
			panic(err)
		}
	}
}
