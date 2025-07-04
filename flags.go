package main

import (
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ProfilerEnabled bool
	ChannelOut      string
	ChannelIn       string
	StreamIdIn      int
	StreamIdOut     int
	Messages        int
	Size            int
	LoggingOn       bool
	Timeout         int
	AeronDir        string
	Idle            string
	ClusterDir      string
	ClusterId       int
}

func mustAtoI(s string) int64 {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int64(v)
}

// backoffp,10,20,1000,1000000
// backoff
// busyspin
// sleeping,1000
// yield
func ToIdleStrategy(idle string) idlestrategy.Idler {
	vs := strings.Split(idle, ",")
	switch vs[0] {
	case "backoffp":
		v1 := mustAtoI(vs[1])
		v2 := mustAtoI(vs[2])
		v3 := mustAtoI(vs[3])
		v4 := mustAtoI(vs[4])
		return idlestrategy.NewBackoffIdleStrategy(v1, v2, v3, v4)
	case "":
		fallthrough
	case "backoff":
		return idlestrategy.NewDefaultBackoffIdleStrategy()
	case "busyspin":
		return idlestrategy.Busy{}
	case "sleeping":
		return idlestrategy.Sleeping{
			SleepFor: time.Duration(mustAtoI(vs[1])) * time.Nanosecond,
		}
	case "yield":
		return idlestrategy.Yielding{}
	default:
		panic("Unknown idle strategy: " + idle)
	}
}
