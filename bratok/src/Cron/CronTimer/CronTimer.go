package CronTimer

/*
The module generates timer's impulses and sends them to main process.
*/

import (
	"Cron/CronMessage"
	"strconv"
	"sync"
	"time"
)

type Timer struct {
	chOut chan CronMessage.Mess
	chIn  chan CronMessage.Mess
}

func Start(mSec int, wg sync.WaitGroup, chOut chan CronMessage.Mess) *Timer {

	chIn := CronMessage.Channel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		//myTimer := time.NewTimer(time.Duration(mSec) * time.Second)
		ticker := time.NewTicker(time.Duration(mSec) * time.Second)

		for {
			select {
			case t := <-ticker.C:
				CronMessage.Send(chOut, "", "timer", []byte{}, nil, t)
			case m, isGood := <-chIn:
				if !isGood {
					ticker.Stop()
					return
				}

				if m.Type == "stop" {
					ticker.Stop()
					return
				}

				if m.Type == "reset" {
					if i, err := strconv.Atoi(string(m.Data)); err == nil {
						if i > 0 {
							mSec = i
							ticker = time.NewTicker(time.Duration(mSec) * time.Second)
						}
					}
				}
			}
		}
	}()

	return &Timer{
		chOut: chOut,
		chIn:  chIn,
	}
}
