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
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-systemd/activation"
	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/journal"
)

func main() {
	journal.Print(journal.PriInfo, "Initialising Commander.")

	showHelp := flag.Bool("h", false, "Show Usage (This menu).")
	isDebug := flag.Bool("d", false, "Enter debug mode.")
	sockFile := flag.String("socket", "/tmp/commander.sock", "Path to unix socket for use in debug mode.")
	flag.Parse()

	if *showHelp {
		fmt.Println("Usage: ", os.Args[0], "[-d] [-socket=<socket_file>]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// If Debug is enabled, then use socket file specified by the user, otherwise
	//  try and retrieve a pre-configured Unix Socket from systemd.
	var listener net.Listener
	if *isDebug {
		var err error
		listener, err = net.Listen("unix", *sockFile)
		if err != nil {
			journal.Print(journal.PriErr, "Unable to activate socket.")
			panic(err)
		}
	} else {
		listeners, err := activation.Listeners(true)
		if err != nil {
			journal.Print(journal.PriErr, "Unable to activate socket.")
			panic(err)
		}

		listener = listeners[0]
	}

	defer listener.Close()

	// Create our Commander object; this will handle all HTTP interactions.
	commander := new(Commander)
	_, err := commander.Init()
	if err != nil {
		journal.Print(journal.PriErr, "Unable to initialise Commander!")
		panic(err)
	}

	// Is systemd watchdog enabled? If so, register a HTTP endpoint which can
	//  handle watchdog functionality.
	isWatchDogEnabled := false
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil && interval == 0 { // erhh... if err is NOT nil, then set watchdog to true...?!?!
		isWatchDogEnabled = true
		http.HandleFunc("/watchdog", func(w http.ResponseWriter, req *http.Request) {
			daemon.SdNotify(false, "WATCHDOG=1")
		})
	}

	// Configure Commander HTTP Handlers
	http.HandleFunc("/listing", commander.Listing)
	http.HandleFunc("/status", commander.ProvideStatus)
	http.HandleFunc("/dispatch", commander.HandleCommand)
	http.Serve(listener, nil)

	journal.Print(journal.PriInfo, "Commander is now listening for requests.")

	// All init is done; so create a go-routine that simply hits our watchdog endpoint
	//  - keeping systemd aware of the health of this process.
	if isWatchDogEnabled {
		go func() {
			for {
				/* Make request to listeners[0] '/watchdog' */
				time.Sleep(interval / 3)
			}
		}()
	}
}
