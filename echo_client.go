package main

import (
	"fmt"
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"log"
	"strconv"
	"time"
)

type EchoClient struct {
	C   *Config
	a   *aeron.Aeron
	pub *aeron.Publication
	sub *aeron.Subscription
}

func (e *EchoClient) clientEcho() {

	to := time.Second * (time.Duration(e.C.Timeout))
	ctx := aeron.NewContext().
		MediaDriverTimeout(to).
		AeronDir(e.C.AeronDir).
		IdleStrategy(ToIdleStrategy(e.C.Idle)).
		ErrorHandler(func(err error) {
			logger.Fatalf("Received error: %v", err)
		})
	var err error
	e.a, err = aeron.Connect(ctx)
	if err != nil {
		logger.Fatalf("Failed to connect to media driver: %s\n", err.Error())
	}
	defer e.a.Close()

	log.Printf("Connected Cnc File: %s\n", ctx.CncFileName())

	go e.keepSending()
	e.listen()
}

func (e *EchoClient) keepSending() {
	pub, err := e.a.AddPublication(e.C.ChannelOut, int32(e.C.StreamIdOut))
	if err != nil {
		logger.Fatal(err)
	}
	defer pub.Close()
	e.pub = pub

	for !pub.IsConnected() {
		log.Printf("Waiting for publication to connect...%s %d", pub.Channel(), pub.StreamID())
		time.Sleep(time.Second)
	}

	log.Printf("Publication created %v", pub)

	// allocatte a 64byte buffer for sending responses
	sendBuffer := atomic.MakeBuffer(make([]byte, 64))

	pubIdle := ToIdleStrategy(e.C.Idle)

	count := 0

	for {
		tNow := time.Now().UnixNano()
		message := fmt.Sprintf("%d", tNow)
		v := []byte(message)
		l := int32(len(v))

		sendBuffer.PutBytesArray(0, &v, 0, l)

		e.send(pubIdle, pub, sendBuffer, 0, l)
		count += 1
		//logger.Infof("msg %d: %s", count, message)
		time.Sleep(time.Second)
	}
}

func (e *EchoClient) listen() {
	sub, err := e.a.AddSubscription(e.C.ChannelIn, int32(e.C.StreamIdIn))
	if err != nil {
		logger.Fatal(err)
	}
	defer sub.Close()
	log.Printf("Subscription created %v", sub)

	e.sub = sub

	subIdle := ToIdleStrategy(e.C.Idle)

	counter := 0

	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
		s := string(buffer.GetBytesArray(offset, length))
		nano, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Printf("Error parsing fragment payload: %s", err)
			return
		}
		tNow := time.Now().UnixNano()

		fmt.Printf("%8.d: Frag offset=%d length=%d delay=%d ns %d us payload: %s\n", counter, offset, length, tNow-nano, (tNow-nano)/1000, s)
	}

	for {
		fragmentsRead := sub.Poll(handler, 1000)
		subIdle.Idle(fragmentsRead)
	}
}

func (e *EchoClient) send(c idlestrategy.Idler, pub *aeron.Publication, buffer *atomic.Buffer, index int32, length int32) {
	for {
		result := pub.Offer(buffer, index, length, nil)
		if result >= 0 {
			return
		} else if result == aeron.BackPressured || result == aeron.AdminAction {
			c.Idle(0)
		} else {
			log.Printf("WARNING: Offer failed - result=%v\n", result)
		}
		time.Sleep(time.Second)
	}

}
