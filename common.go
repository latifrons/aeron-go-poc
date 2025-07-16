package main

import (
	"github.com/lirm/aeron-go/aeron/idlestrategy"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

type ServerConfig struct {
}

type ClientConfig struct {
}

type ClusterServerConfig struct {
	AeronDir   string
	ClusterDir string
	Idle       string
}

type ClusterClientConfig struct {
	AeronDir         string
	Idle             string
	IngressChannel   string
	IngressEndpoints string
	IngressStreamId  int32
	EgressChannel    string
	EgressStreamId   int32
}

func FromEnvOrDefaultString(key string, defaultValue string) string {
	if v := viper.GetString(key); v != "" {
		return v
	}
	return defaultValue
}

func FromEnvOrDefaultInt(key string, defaultValue int) int {
	if v := viper.GetInt(key); v != 0 {
		return v
	}
	return defaultValue
}

func MustEnv(key string) string {
	v := viper.GetString(key)
	if v == "" {
		panic("Environment variable " + key + " is not set")
	}
	return v
}

func MustEnvInt32(key string) int32 {
	v := viper.GetInt(key)
	if v == 0 {
		panic("Environment variable " + key + " is not set or is zero")
	}
	return int32(v)
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

func mustAtoI(s string) int64 {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int64(v)
}
