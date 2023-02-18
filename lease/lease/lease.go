package lease

import (
	"log"
	"sync"
	"time"
)

type Lease struct {
	Id         uint64
	Name       string
	Wg         sync.WaitGroup
	TTL        time.Duration
	Timer      *time.Timer
	MutexLock  sync.Mutex
	UnlockTime *time.Time
}

func (l *Lease) TryLock(TTL time.Duration) bool {

	if locked := l.MutexLock.TryLock(); locked {
		var timer = time.NewTimer(TTL)
		var unlockTime = time.Now().Add(TTL)

		log.Printf("Lease with id: %d and name: %s locked for duration %s", l.Id, l.Name, TTL)

		l.TTL = TTL
		l.UnlockTime = &unlockTime
		l.Timer = timer
		l.Wg.Add(1)

		go func() {
			<-timer.C
			l.Wg.Done()
			timer.Stop()
			l.MutexLock.Unlock()
			log.Printf("Lease with id: %d and name: %s unlocked", l.Id, l.Name)
		}()

		return true
	}

	log.Printf("Lease with id: %d and name: %s lock failed", l.Id, l.Name)
	return false
}

func (l *Lease) Renew(TTL time.Duration) {
	var unlockTime = time.Now().Add(TTL)

	l.TTL = TTL
	l.UnlockTime = &unlockTime
	l.Timer.Reset(TTL)

	log.Printf("Lease with id: %d and name: %s renewed for duration %s", l.Id, l.Name, TTL)
}

func (l *Lease) Wait() {
	l.Wg.Wait()
}

func New(id uint64, name string) *Lease {
	return &Lease{
		Id:         id,
		Name:       name,
		TTL:        0,
		Timer:      nil,
		MutexLock:  sync.Mutex{},
		UnlockTime: nil,
	}
}
