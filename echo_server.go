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

type EchoServer struct {
	C   *Config
	a   *aeron.Aeron
	pub *aeron.Publication
	sub *aeron.Subscription
}

func (e *EchoServer) serverEcho() {

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

	e.listenAndEcho()
}

func (e *EchoServer) listenAndEcho() {
	sub, err := e.a.AddSubscription(e.C.ChannelIn, int32(e.C.StreamIdIn))
	if err != nil {
		logger.Fatal(err)
	}
	defer sub.Close()
	log.Printf("Subscription created %+v", sub)

	e.sub = sub

	pub, err := e.a.AddPublication(e.C.ChannelOut, int32(e.C.StreamIdOut))
	if err != nil {
		logger.Fatal(err)
	}
	defer sub.Close()
	log.Printf("Publication created %+v", pub)

	e.pub = pub

	counter := 0

	startTime := time.Now()

	// allocatte a 64byte buffer for sending responses
	sendBuffer := atomic.MakeBuffer(make([]byte, 64))

	pubIdle := ToIdleStrategy(e.C.Idle)
	subIdle := ToIdleStrategy(e.C.Idle)

	handler := func(buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {

		s := string(buffer.GetBytesArray(offset, length))
		nano, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Printf("Error parsing fragment payload: %s", err)
			return
		}
		tNow := time.Now().UnixNano()

		fmt.Printf("%8.d: Frag offset=%d length=%d delay=%d ns %d us payload: %s\n", counter, offset, length, tNow-nano, (tNow-nano)/1000, s)

		responseBody := fmt.Sprintf("%s", s)
		v := []byte(responseBody)
		l := int32(len(v))

		sendBuffer.PutBytesArray(0, &v, 0, l)

		e.send(pubIdle, pub, sendBuffer, 0, l)

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

	for {
		fragmentsRead := sub.Poll(handler, 1000)
		subIdle.Idle(fragmentsRead)
	}
}

//
//func handleDisconnect(image aeron.Image) {
//	log.Printf("Received disconnect for image with correlation ID %d", image.CorrelationID())
//	pub, ok := connectionMap[image.CorrelationID()]
//	if !ok {
//		log.Printf("No disconnect publication found for correlation ID %d", image.CorrelationID())
//		return
//	}
//	log.Printf("Disconnecting image with correlation ID %d", image.CorrelationID())
//	err := pub.Close()
//	if err != nil {
//		log.Printf("Error closing publication: %v", err)
//		return
//	}
//	delete(connectionMap, image.CorrelationID())
//}
//
//func handleConnect(image aeron.Image) {
//	log.Printf("Received connect for image with correlation ID %d", image.CorrelationID())
//	pub, err := aeronObject.AddPublication(
//		cc.ChannelOut+"|response-correlation-id="+strconv.FormatInt(image.CorrelationID(), 10), int32(cc.StreamIdOut))
//	if err != nil {
//		logger.Fatal(err)
//	}
//	defer pub.Close()
//
//	for !pub.IsConnected() {
//		log.Printf("Waiting for publication to connect...%s %d", pub.Channel(), pub.StreamID())
//		time.Sleep(time.Second)
//	}
//
//	log.Printf("Publication created %+v", pub)
//
//	connectionMap[image.CorrelationID()] = pub
//	lastCorrelationID = image.CorrelationID()
//}

func (e *EchoServer) send(c idlestrategy.Idler, pub *aeron.Publication, buffer *atomic.Buffer, index int32, length int32) {
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
