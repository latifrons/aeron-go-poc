/*
Copyright 2016 Stanislav Liberman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/spf13/viper"
	"log"
)

var logger = logging.MustGetLogger("basic_subscriber")

func main() {
	viper.SetEnvPrefix("INJ")
	viper.AutomaticEnv()

	for k, v := range viper.AllSettings() {
		log.Printf("Settings: %v = %v\n", k, v)
	}
	viper.SetDefault("profiler", false)
	viper.SetDefault("channel_out", "aeron:ipc")
	viper.SetDefault("channel_in", "aeron:ipc")
	viper.SetDefault("streamId", 10)
	viper.SetDefault("count", 1000000)
	viper.SetDefault("size", 256)
	viper.SetDefault("logging", true)
	viper.SetDefault("timeout", 10)
	viper.SetDefault("dir", "./aeron")
	viper.SetDefault("idle", "backoff")

	c := &Config{
		ProfilerEnabled: viper.GetBool("profiler"),
		ChannelOut:      viper.GetString("channel_out"),
		ChannelIn:       viper.GetString("channel_in"),
		Messages:        viper.GetInt("count"),
		Size:            viper.GetInt("size"),
		LoggingOn:       viper.GetBool("logging"),
		Timeout:         viper.GetInt("timeout"),
		AeronDir:        viper.GetString("dir"),
		Idle:            viper.GetString("idle"),
	}
	cmd := viper.GetString("command")

	//if !C.LoggingOn {
	//	logging.SetLevel(logging.INFO, "aeron")
	//	logging.SetLevel(logging.INFO, "memmap")
	//	logging.SetLevel(logging.DEBUG, "driver")
	//	logging.SetLevel(logging.INFO, "counters")
	//	logging.SetLevel(logging.INFO, "logbuffers")
	//	logging.SetLevel(logging.INFO, "buffer")
	//}

	// use the first arg to select func
	switch cmd {
	case "basicSubscriber":
		basicSubscriber(c)
		return
	case "basicPublisher":
		basicPublisher(c)
		return
	case "clusterClientBenchmark":
		clusterClientBenchmark(c)
		return
	case "clusterServerEcho":
		clusterServerEcho(c)
		return
	case "echoServer":
		c.StreamIdIn = PingStreamId
		c.StreamIdOut = PongStreamId
		es := EchoServer{
			C: c,
		}
		es.serverEcho()
	case "echoClient":
		c.StreamIdIn = PongStreamId
		c.StreamIdOut = PingStreamId
		es := EchoClient{
			C: c,
		}
		es.clientEcho()
	default:
		logger.Fatalf("Unknown example: %s", cmd)
	}

}
