/*
	Package pecho provides simple level separated logging library.
	You must make a channel via MkChannel, then call StartDaemon() in a goroutine to get pecho started.
	Afterwards, use channel to send messages.

	The default logging level is defined via environment variable PECHO_LEVEL, of which can be debug, info, warn or crit.

	Example:

	```
	chann := pecho.MkChannel()
	go pecho.StartDaemon(chann)
	chann <- []string{"debug", "Unrecognised log level " + rawUserLevel}
	```
*/

package pecho

import (
	"fmt"
	"os"
)

func MkChannel() chan []string {
	var chann = make(chan []string, 256)
	return chann
}

/*
	Start the logging daemon
	usage: go pecho.StartDaemon(channel)
*/
func StartDaemon(logChan chan []string) {
	rawUserLevel := os.Getenv("PECHO_LEVEL")
	var userLevel int
	switch rawUserLevel {
		case "debug":
			userLevel = 1
		case "info":
			userLevel = 2
		case "warn":
			userLevel = 3
		case "crit":
			userLevel = 4
		default:
			userLevel = 2
			logChan <- []string{"debug", "Unrecognised log level " + rawUserLevel}
	}
	logChan <- []string{"debug", "Started logging daemon"}
	for incoming := range logChan {
		msgLevel := 0
		switch incoming[0] {
			case "debug":
				msgLevel = 1
			case "info":
				msgLevel = 2
			case "warn":
				msgLevel = 3
			case "crit":
				fmt.Println("Critical: " + "incoming[1]")
				os.Exit(1)
				return
		}
		if userLevel <= msgLevel {
			/* SCARY!!!
			This will panic on malformed events...
			We better guard the channel behind a function
			*/
			fmt.Println(
				"[", incoming[0], "]: ",
				incoming[1],
			)
		}
	}
	fmt.Println("The logging daemon has shutdown")
}