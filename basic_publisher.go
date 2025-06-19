package main

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"log"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
)

func basicPublisher(c *Config) {
	to := time.Second * (time.Duration(c.Timeout))
	ctx := aeron.NewContext().MediaDriverTimeout(to).AeronDir(c.AeronDir)
	a, err := aeron.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer a.Close()

	log.Printf("Connected Cnc File: %s\n", ctx.CncFileName())

	pub, err := a.AddPublication(c.Channel, int32(c.StreamId))
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Close()
	log.Printf("Publication created %v", pub)

	idle := idlestrategy.NewDefaultBackoffIdleStrategy()
	counter := 0

	for i := 0; i < c.Messages; i++ {
		message := fmt.Sprintf("%d", time.Now().UnixNano())

		srcBuffer := atomic.MakeBuffer(([]byte)(message))

		for {
			res := pub.Offer(srcBuffer, 0, int32(len(message)), nil)

			if res > 0 {
				break
			} else if res == aeron.BackPressured || res == aeron.AdminAction {
				idle.Idle(0)
			} else {
				log.Printf("Offer failed: %d", res)
				idle.Idle(0)
			}
		}
		counter++
		fmt.Printf("Published %d messages: %s\n", counter, message)
		time.Sleep(time.Second)
	}
	fmt.Printf("Published %d messages of size %d to channel %s stream %d\n", c.Messages, c.Size, c.Channel, c.StreamId)
}
