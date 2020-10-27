package cmdroot

import (
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	return userErrorRegexp.MatchString(err.Error())
}

// WaitSignal block until os.Interrupt, os.Kill, syscall.SIGTERM
func WaitSignal() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		stop := false
		select {
		case <-sigChannel:
			stop = true
		}
		if stop {
			break
		}
	}
}
