package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"fmt"
	"time"
)

const KeyNotFoundError = "KeyNotFoundError"
const StateLoadError = "StateLoadError"

const APIStateFilename = "state.json"


type API struct {
	ctx *Context
	address string
	port int
	started_at time.Time
}
var server *API

func (api *API) serve() {
	log.Println("Listening on %s:%s", api.address, api.port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.address, api.port), APIRouter(api.ctx))
}

// NewAPI creates new service object.
func NewAPI(address string, port int) *API {
	return &API{
		ctx: NewContext(true),
		address: address,
		port: port,
		started_at: time.Now(),
	}
}

// init hooks into sys calls and listens for interrupts.
func init() {
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Println("Interrupt signal received. Stopping program!")

		server.ctx.flush()
		os.Exit(0)
	}()
}

// main magic.
func main() {
	server = NewAPI("0.0.0.0", 8080)
	server.serve()
}
