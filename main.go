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
	"flag"
	"github.com/lirm/aeron-go/aeron/logging"
	"os"
)

var logger = logging.MustGetLogger("basic_subscriber")

func main() {

	c := &Config{
		ProfilerEnabled: flag.Bool("profiler", false, "Enable profiler"),
		Channel:         flag.String("channel", "aeron:ipc", "Channel to use for subscription"),
		StreamId:        flag.Int("streamId", 10, "Stream ID"),
		Messages:        flag.Int("count", 1000000, "Number of messages to send"),
		Size:            flag.Int("size", 256, "messages size"),
		LoggingOn:       flag.Bool("logging", false, "Enable logging"),
		Timeout:         flag.Int("timeout", 10, "Timeout in seconds"),
		AeronDir:        flag.String("dir", "", "Aeron directory (default: empty, uses default aeron dir)"),
	}

	cmd := os.Args[1]
	os.Args = os.Args[1:]

	flag.Parse()

	if !*c.LoggingOn {
		logging.SetLevel(logging.INFO, "aeron")
		logging.SetLevel(logging.INFO, "memmap")
		logging.SetLevel(logging.DEBUG, "driver")
		logging.SetLevel(logging.INFO, "counters")
		logging.SetLevel(logging.INFO, "logbuffers")
		logging.SetLevel(logging.INFO, "buffer")
	}

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
