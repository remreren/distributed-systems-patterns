package main

import (
	"distributed-systems/lease/lease"
	"log"
	"time"
)

func main() {
	var ls = lease.New(15, "lease1")

	log.SetFlags(log.Ltime | log.Lmicroseconds)

	ls.TryLock(2 * time.Second) // locks for 2 seconds
	ls.TryLock(2 * time.Second) // fails to lock again
	ls.Renew(3 * time.Second)   // renews the lock
	ls.Wait()                   // waits for lock to finish and quits program

	time.Sleep(10 * time.Millisecond)
}
