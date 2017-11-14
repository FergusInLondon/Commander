// main.go - Fergus In London <fergus@fergus.london> (November 2017)
//
// Simple little systemd enabled service that is capable of registering "commands"
//  that can be called via a simple JSON API over a Unix Socket; this can then allow
//  other client applications (be they native GUI apps, Electron Apps, or even Web
//  Interfaces etc) to defer low level tasks elsewhere.
//
// View the usage for details on how to run without configuring with systemd.
//
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/journal"
)


var (
	showHelp = flag.Bool("h", false, "Show Usage (This menu).")
	isDebug  = flag.Bool("d", false, "Enter debug mode.")
	sockFile = flag.String("socket", "/tmp/commander.sock", "Path to unix socket for use in debug mode.")
)


func main() {
	journal.Print(journal.PriInfo, "Initialising Commander.")

	flag.Parse()
	if *showHelp {
		fmt.Println("Usage: ", os.Args[0], "[-d] [-socket=<socket_file>]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Grab our socket for listening on
	listener, _ := get_api_listener()
	defer listener.Close()

	// Create our Commander object; this will handle all HTTP interactions.
	commander := new(Commander)
	if err := commander.Init(); err != nil {
		journal.Print(journal.PriErr, "Unable to initialise Commander!")
		panic(err)
	}

	// Create a new Router, and pass it to Commander for configuration
	router := mux.NewRouter()
	commander.ConfigureListeners(router)

	// Is systemd watchdog enabled? If so, register a HTTP endpoint which can
	//  handle watchdog functionality.
	isWatchDogEnabled := false
	interval, err := daemon.SdWatchdogEnabled(false)
	if err == nil && interval > 0 {
		isWatchDogEnabled = true
		router.HandleFunc("/watchdog", func(w http.ResponseWriter, req *http.Request) {
			daemon.SdNotify(false, "WATCHDOG=1")
			w.WriteHeader(200)
		})
	}

	// Configuration complete!
	http.Serve(listener, router)
	journal.Print(journal.PriInfo, "Commander is now listening for requests.")

	// All init is done; so create a go-routine that simply hits our watchdog endpoint
	//  - keeping systemd aware of the health of this process.
	if isWatchDogEnabled {
		go func() {
			for {
				call_systemd_healthcheck()
				time.Sleep(interval / 3)
			}
		}()
	}
}


func get_api_listener() (listener net.Listener, err error){
	// If Debug is enabled, then use socket file specified by the user, otherwise
	//  try and retrieve a pre-configured Unix Socket from systemd.
	if *isDebug {
		listener, err = net.Listen("unix", *sockFile)
		if err == nil {
			return
		}
	} else {
		listeners, err := activation.Listeners(true)
		if err == nil {
			return listeners[0], nil
		}
	}

	journal.Print(journal.PriErr, "Unable to activate socket.")
	panic(err)
}


// Thanks to @teknoraver for the example of HTTP over Unix Sockets.
//  - https://gist.github.com/teknoraver/5ffacb8757330715bcbcc90e6d46ac74
func call_systemd_healthcheck() {
	// if this is being called, then systemd integration is enabled - and we can
	//  be quite sure of the address of the unix socket.
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/tmp/commander.sock")
			},
		},
	}

	httpc.Get("http://unix/watchdog")
}
