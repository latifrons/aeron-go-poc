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
	"github.com/lirm/aeron-go/aeron/logging"
	"github.com/spf13/viper"
)

var logger = logging.MustGetLogger("basic_subscriber")

//func profileConfig(profile string, clusterId int) (c *Config) {
//	switch profile {
//
//	}
//
//	func
//	profileConfig(profile
//	string, clusterId
//	int) (c * Config)
//	{
//		switch profile {
//		case "clusterClientBenchmark":
//			fallthrough
//		case "clusterLatencyCheckClient":
//			c = &Config{
//				ProfilerEnabled: false,
//				ChannelOut:      "",
//				ChannelIn:       "",
//				StreamIdIn:      0,
//				StreamIdOut:     0,
//				Messages:        0,
//				Size:            0,
//				LoggingOn:       true,
//				Timeout:         0,
//				//AeronDir:         fmt.Sprintf("C:\\Users\\LATFIR~1\\AppData\\Local\\Temp\\aeron-latfirons-%d-driver", clusterId),
//				AeronDir:         fmt.Sprintf("E:\\dev\\aeron-md"),
//				Idle:             "",
//				ClusterDir:       "",
//				ClusterId:        0,
//				IngressChannel:   "aeron:udp",
//				IngressEndpoints: "0=localhost:11010,1=localhost:11110,2=localhost:11210",
//				//IngressEndpoints: "0=localhost:11010",
//				IngressStreamId: 101,
//				EgressChannel:   "aeron:udp?endpoint=localhost:20000",
//				EgressStreamId:  102,
//			}
//		case "clusterServerEcho":
//			c = &Config{
//				ProfilerEnabled: false,
//				ChannelOut:      "",
//				ChannelIn:       "",
//				StreamIdIn:      0,
//				StreamIdOut:     0,
//				Messages:        0,
//				Size:            0,
//				LoggingOn:       true,
//				Timeout:         0,
//				AeronDir:        fmt.Sprintf("C:\\Users\\LATFIR~1\\AppData\\Local\\Temp\\aeron-latfirons-%d-driver", clusterId),
//				Idle:            "",
//				ClusterDir:      fmt.Sprintf("E:\\dev\\aeron-cluster-%d", clusterId),
//				ClusterId:       clusterId,
//			}
//		default:
//			panic(fmt.Errorf("unknown profile %s", profile))
//		}
//		return
//	}

func main() {
	viper.SetEnvPrefix("INJ")
	viper.AutomaticEnv()

	logger.SetLevel(logging.DEBUG)

	cmd := viper.GetString("command")
	//profile := viper.GetString("profile")
	//clusterId := viper.GetInt("cluster_id")

	//if c.LoggingOn {
	logging.SetLevel(logging.DEBUG, "aeron")
	logging.SetLevel(logging.DEBUG, "memmap")
	logging.SetLevel(logging.DEBUG, "driver")
	logging.SetLevel(logging.DEBUG, "counters")
	logging.SetLevel(logging.DEBUG, "logbuffers")
	logging.SetLevel(logging.DEBUG, "buffer")
	logging.SetLevel(logging.DEBUG, "cluster-client")
	//}
	c := &Config{}

	// use the first arg to select func
	switch cmd {
	case "basicSubscriber":
		basicSubscriber(c)
		return
	case "basicPublisher":
		basicPublisher(c)
		return
	case "clusterLatencyCheckClient":
		clusterLatencyCheckClient(&ClusterClientConfig{
			AeronDir:         FromEnvOrDefaultString("aeron.dir", ""),
			Idle:             FromEnvOrDefaultString("aeron.idle", ""),
			IngressChannel:   FromEnvOrDefaultString("aeron.ingressChannel", "aeron:udp"),
			IngressEndpoints: MustEnv("aeron.ingressEndpoints"),
			IngressStreamId:  MustEnvInt32("aeron.ingressStreamId"),
			EgressChannel:    FromEnvOrDefaultString("aeron.egressChannel", "aeron:udp"),
			EgressStreamId:   MustEnvInt32("aeron.egressStreamId"),
		})
		return
	case "clusterBenchmarkClient":
		clusterBenchmarkClient(&ClusterClientConfig{
			AeronDir:         FromEnvOrDefaultString("aeron.dir", ""),
			Idle:             FromEnvOrDefaultString("aeron.idle", ""),
			IngressChannel:   FromEnvOrDefaultString("aeron.ingressChannel", "aeron:udp"),
			IngressEndpoints: MustEnv("aeron.ingressEndpoints"),
			IngressStreamId:  MustEnvInt32("aeron.ingressStreamId"),
			EgressChannel:    FromEnvOrDefaultString("aeron.egressChannel", "aeron:udp"),
			EgressStreamId:   MustEnvInt32("aeron.egressStreamId"),
		})
		return
	case "clusterServerEcho":
		clusterServerEcho(&ClusterServerConfig{
			AeronDir:   FromEnvOrDefaultString("aeron.dir", ""),
			ClusterDir: MustEnv("aeron.cluster.clusterDir"),
			Idle:       FromEnvOrDefaultString("aeron.idle", ""),
		})
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
