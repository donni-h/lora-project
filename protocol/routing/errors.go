package routing

import "log"

func (a *AODV) processError(err error) {
	log.Printf("Error: %s\n", err)
}
