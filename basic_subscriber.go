package main

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"log"
	"strconv"
	"time"
)

func basicSubscriber(c *Config) {
	to := time.Second * (time.Duration(c.Timeout))
	ctx := aeron.NewContext().MediaDriverTimeout(to).AeronDir(c.AeronDir).IdleStrategy(ToIdleStrategy(c.Idle))

	a, err := aeron.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer a.Close()

	log.Printf("Connected Cnc File: %s\n", ctx.CncFileName())

	//subscription, err := a.AddSubscription("aeron:ipc", 10)
	subscription, err := a.AddSubscription(c.ChannelIn, int32(c.StreamIdIn))
	if err != nil {
		logger.Fatal(err)
	}
	defer subscription.Close()
	log.Printf("Subscription found %v", subscription)

	counter := 0

	startTime := time.Now()

	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		s := string(buffer.GetBytesArray(offset, length))
		nano, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Printf("Error parsing fragment payload: %s", err)
			return
		}
		tNow := time.Now().UnixNano()

		fmt.Printf("%8.d: Frag offset=%d length=%d delay=%d ns %d us payload: %s\n", counter, offset, length, tNow-nano, (tNow-nano)/1000, s)

		counter++

		if counter == 1 {
			startTime = time.Now()
		}
		if counter%1000 == 0 {
			elapsed := time.Since(startTime)
			log.Printf("Received %d fragments in %s, TPS: %.2f\n", counter, elapsed, float64(counter)/elapsed.Seconds())
		}

		if counter == 1000000 {
			elapsed := time.Since(startTime)
			log.Printf("End %d fragments in %s, TPS: %.2f\n", counter, elapsed, float64(counter)/elapsed.Seconds())
			counter = 0
		}
	}

	//idleStrategy := idlestrategy.Sleeping{SleepFor: time.Millisecond}
	//idleStrategy := idlestrategy.Busy{}

	idle := ToIdleStrategy(c.Idle)
	for {
		fragmentsRead := subscription.Poll(handler, 1000)
		idle.Idle(fragmentsRead)
	}
}
