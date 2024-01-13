package remote_conn

import "sync"

var (
	ticketDeficit       int
	ticketDeficitLock   sync.Locker
	chConnectionTickets chan struct{}
)

func init() {
	// 200 causes errors
	// 150 sometimes causes errors
	// chConnectionTickets = make(chan struct{}, 150)
	ticketDeficitLock = &sync.Mutex{}
	chConnectionTickets = make(chan struct{})
}

func WaitForConnectionLicense() {
	ticketDeficitLock.Lock()
	ticketDeficit++
	ticketDeficitLock.Unlock()

	<-chConnectionTickets

	ticketDeficitLock.Lock()
	ticketDeficit--
	ticketDeficitLock.Unlock()
}

func ReturnConnLicense() {
	ticketDeficitLock.Lock()

	if ticketDeficit > 0 {
		go func() {
			chConnectionTickets <- struct{}{}
		}()
	}

	ticketDeficitLock.Unlock()
}
