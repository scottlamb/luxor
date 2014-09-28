Go library for controlling an [FX Luminaire
Luxor ZD](http://www.fxl.com/product/power-and-control/luxor) landscape
lighting system via its wifi module protocol. Provides a typesafe RPC
interface which abstracts away HTTP and JSON handling via methods of the form:

    Method(ctx context.Context, req *MethodRequest) (*MethodResponse, error)

Get started with:

    go get github.com/scottlamb/luxor
    godoc github.com/scottlamb/luxor/protocol
    godoc github.com/scottlamb/luxor/client

See `illuminate_all.go` for a simple example client.
