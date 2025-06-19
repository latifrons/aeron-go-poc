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
	viper.SetDefault("channel", "aeron:ipc")
	viper.SetDefault("streamId", 10)
	viper.SetDefault("count", 1000000)
	viper.SetDefault("size", 256)
	viper.SetDefault("logging", false)
	viper.SetDefault("timeout", 10)
	viper.SetDefault("dir", "")

	c := &Config{
		ProfilerEnabled: viper.GetBool("profiler"),
		Channel:         viper.GetString("channel"),
		StreamId:        viper.GetInt("streamId"),
		Messages:        viper.GetInt("count"),
		Size:            viper.GetInt("size"),
		LoggingOn:       viper.GetBool("logging"),
		Timeout:         viper.GetInt("timeout"),
		AeronDir:        viper.GetString("dir"),
	}
	cmd := viper.GetString("command")

	//if !c.LoggingOn {
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
	default:
		logger.Fatalf("Unknown example: %s", cmd)
	}

}
