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

import "C"
import (
	"fmt"
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/spf13/viper"
)

var logger = logging.MustGetLogger("basic_subscriber")

func profileConfig(profile string, clusterId int) (c *Config) {
	switch profile {
	case "clusterClientBenchmark":
		c = &Config{
			ProfilerEnabled:  false,
			ChannelOut:       "",
			ChannelIn:        "",
			StreamIdIn:       0,
			StreamIdOut:      0,
			Messages:         0,
			Size:             0,
			LoggingOn:        false,
			Timeout:          0,
			AeronDir:         fmt.Sprintf("C:\\Users\\LATFIR~1\\AppData\\Local\\Temp\\aeron-latfirons-%d-driver", clusterId),
			Idle:             "",
			ClusterDir:       "",
			ClusterId:        0,
			IngressChannel:   "aeron:udp",
			IngressEndpoints: "0=localhost:9010",
			IngressStreamId:  101,
			EgressChannel:    "aeron:udp?endpoint=localhost:9011",
			EgressStreamId:   102,
		}
	case "clusterServerEcho":
		c = &Config{
			ProfilerEnabled: false,
			ChannelOut:      "",
			ChannelIn:       "",
			StreamIdIn:      0,
			StreamIdOut:     0,
			Messages:        0,
			Size:            0,
			LoggingOn:       false,
			Timeout:         0,
			AeronDir:        fmt.Sprintf("C:\\Users\\LATFIR~1\\AppData\\Local\\Temp\\aeron-latfirons-%d-driver", clusterId),
			Idle:            "",
			ClusterDir:      fmt.Sprintf("E:\\dev\\aeron-cluster-%d", clusterId),
			ClusterId:       clusterId,
		}
	default:
		panic(fmt.Errorf("unknown profile %s", profile))
	}
	return
}

func main() {
	viper.SetEnvPrefix("INJ")
	viper.AutomaticEnv()

	logger.SetLevel(logging.DEBUG)

	cmd := viper.GetString("command")
	//profile := viper.GetString("profile")
	clusterId := viper.GetInt("cluster_id")

	c := profileConfig(cmd, clusterId)

	if !c.LoggingOn {
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
