package resetSleep

import (
	"log"
	"sync"
	"testing"
	"time"
)

var globalG *G

func TestResetSleep(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		globalG = GetG()
		log.Println("Begin Sleep")
		time.Sleep(time.Second * 100)
		log.Println("Sleep Done")
		wg.Done()
	}()
	time.Sleep(time.Second)
	Resettimer(globalG.Timer, int64(time.Now().Nanosecond())+time.Second.Nanoseconds())
	wg.Wait()
}
