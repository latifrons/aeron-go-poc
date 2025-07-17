package main

import (
	"fmt"
	"time"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster/client"
)

type ClusterLatencyCheckClient struct {
	ac                    *client.AeronCluster
	nextSendKeepAliveTime int64
}

func (ctx *ClusterLatencyCheckClient) OnConnect(ac *client.AeronCluster) {
	fmt.Printf("OnConnect - sessionId=%d leaderMemberId=%d leadershipTermId=%d\n",
		ac.ClusterSessionId(), ac.LeaderMemberId(), ac.LeadershipTermId())
	ctx.ac = ac
	ctx.nextSendKeepAliveTime = time.Now().UnixMilli() + time.Second.Milliseconds()
}

func (ctx *ClusterLatencyCheckClient) OnDisconnect(cluster *client.AeronCluster, details string) {
	fmt.Printf("OnDisconnect - sessionId=%d (%s)\n", cluster.ClusterSessionId(), details)
	ctx.ac = nil
}

func (ctx *ClusterLatencyCheckClient) OnMessage(cluster *client.AeronCluster, timestamp int64,
	buffer *atomic.Buffer, offset int32, length int32, header *logbuffer.Header) {
	recvTime := time.Now().UnixNano()
	msgNo := buffer.GetInt32(offset)
	sendTime := buffer.GetInt64(offset + 8)
	latency := recvTime - sendTime
	fmt.Printf("OnMessage - sessionId=%d timestamp=%d pos=%d length=%d msgNo=%d latency=%d\n",
		cluster.ClusterSessionId(), timestamp, header.Position(), length, msgNo, latency)
}

func (ctx *ClusterLatencyCheckClient) OnNewLeader(cluster *client.AeronCluster, leadershipTermId int64, leaderMemberId int32) {
	fmt.Printf("OnNewLeader - sessionId=%d leaderMemberId=%d leadershipTermId=%d\n",
		cluster.ClusterSessionId(), leaderMemberId, leadershipTermId)
}

func (ctx *ClusterLatencyCheckClient) OnError(cluster *client.AeronCluster, details string) {
	fmt.Printf("OnError - sessionId=%d: %s\n", cluster.ClusterSessionId(), details)
}

func (ctx *ClusterLatencyCheckClient) sendKeepAliveIfNecessary() {
	if now := time.Now().UnixMilli(); now > ctx.nextSendKeepAliveTime && ctx.ac != nil && ctx.ac.SendKeepAlive() {
		ctx.nextSendKeepAliveTime += time.Second.Milliseconds()
	}
}

func clusterLatencyCheckClient(c *ClusterClientConfig) {
	ctx := aeron.NewContext().AeronDir(c.AeronDir)

	opts := client.NewOptions()
	//if idleStr := os.Getenv("NO_OP_IDLE"); idleStr != "" {
	opts.IdleStrategy = ToIdleStrategy(c.Idle)
	//}

	// 10002,10102,10202
	opts.IngressChannel = c.IngressChannel
	opts.IngressEndpoints = c.IngressEndpoints
	opts.IngressStreamId = int32(c.IngressStreamId)

	opts.EgressChannel = c.EgressChannel
	opts.EgressStreamId = c.EgressStreamId

	listener := &ClusterLatencyCheckClient{
		//latencies: make([]int64, 1000),
	}
	clusterClient, err := client.NewAeronCluster(ctx, opts, listener)
	if err != nil {
		panic(err)
	}
	clusterClient.Poll()

	for !clusterClient.IsConnected() {
		opts.IdleStrategy.Idle(clusterClient.Poll())
		fmt.Printf("waiting to connect...\n")
		time.Sleep(time.Second)
	}

	sentCt := 0
	sendBuf := atomic.MakeBuffer(make([]byte, 64))

	go listen(opts, clusterClient, listener)

	for {
		sendBuf.PutInt32(0, int32(sentCt))
		sendBuf.PutInt64(8, time.Now().UnixNano())
		for {
			if r := clusterClient.Offer(sendBuf, 0, sendBuf.Capacity()); r >= 0 {
				sentCt++
				break
			}
		}
		//clusterClient.Poll()
		//listener.sendKeepAliveIfNecessary()
		time.Sleep(time.Second)
	}
	clusterClient.Close()
	fmt.Println("done")
	time.Sleep(time.Second)
}

func listen(opts *client.Options, client *client.AeronCluster, listener *ClusterLatencyCheckClient) {
	idleStrategy := opts.IdleStrategy
	for {
		fragmentsRead := client.Poll()
		idleStrategy.Idle(fragmentsRead)
		listener.sendKeepAliveIfNecessary()
	}
}
