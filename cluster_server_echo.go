package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/logbuffer"
	"github.com/lirm/aeron-go/cluster"
	"github.com/lirm/aeron-go/cluster/codecs"
)

type ClusterServerEcho struct {
	cluster cluster.Cluster

	messageCount int32
}

func (s *ClusterServerEcho) OnStart(cluster cluster.Cluster, image aeron.Image) {
	s.cluster = cluster
	if image == nil {
		fmt.Printf("OnStart with no image\n")
	} else {
		cnt := image.Poll(func(buf *atomic.Buffer, offset int32, length int32, hdr *logbuffer.Header) {
			if length == 4 && s.messageCount == 0 {
				s.messageCount = buf.GetInt32(offset)
			} else {
				fmt.Printf("WARNING: unexpected snapshot message - pos=%d offset=%d length=%d\n",
					hdr.Position(), offset, length)
			}
		}, 100)
		fmt.Printf("OnStart with image - snapshotMsgCnt=%d messageCount=%d\n", cnt, s.messageCount)
	}
}

func (s *ClusterServerEcho) OnSessionOpen(session cluster.ClientSession, timestamp int64) {
	fmt.Printf("OnSessionOpen - sessionId=%d timestamp=%v\n", session.Id(), timestamp)
}

func (s *ClusterServerEcho) OnSessionClose(
	session cluster.ClientSession,
	timestamp int64,
	reason codecs.CloseReasonEnum,
) {
	fmt.Printf("OnSessionClose - sessionId=%d timestamp=%v reason=%v\n", session.Id(), timestamp, reason)
}

func (s *ClusterServerEcho) OnSessionMessage(
	session cluster.ClientSession,
	timestamp int64,
	buffer *atomic.Buffer,
	offset int32,
	length int32,
	header *logbuffer.Header,
) {
	s.messageCount++
	// Read the incoming timestamp from the message
	clientTimestamp := string(buffer.GetBytesArray(offset, length))

	srcBuffer := atomic.MakeBuffer(([]byte)(clientTimestamp))

	session.Offer(srcBuffer, 0, length, nil)
	fmt.Printf("OnSessionMessage - sessionId=%d time=%d pos=%d len=%d messageCount=%d clientTimestamp=%s\n",
		session.Id(), timestamp, header.Position(), length, s.messageCount, clientTimestamp)

	//
	//for offerCnt := 1; ; offerCnt++ {
	//	result := session.Offer(buffer, offset, length, nil)
	//	if result >= 0 {
	//		return
	//	} else if result == aeron.BackPressured || result == aeron.AdminAction {
	//		s.cluster.IdleStrategy().Idle(0)
	//	} else {
	//		fmt.Printf("WARNING: OnSessionMessage offer failed - sessionId=%d time=%d pos=%d len=%d offerCnt=%d result=%v\n",
	//			session.Id(), timestamp, header.Position(), length, offerCnt, result)
	//	}
	//}
}

func (s *ClusterServerEcho) OnTimerEvent(correlationId, timestamp int64) {
	fmt.Printf("OnTimerEvent - correlationId=%d timestamp=%v\n", correlationId, timestamp)
}

func (s *ClusterServerEcho) OnTakeSnapshot(publication *aeron.Publication) {
	fmt.Printf("OnTakeSnapshot - streamId=%d sessionId=%d messageCount=%d\n",
		publication.StreamID(), publication.SessionID(), s.messageCount)
	buf := atomic.MakeBuffer(make([]byte, 4))
	buf.PutInt32(0, s.messageCount)
	for {
		result := publication.Offer(buf, 0, buf.Capacity(), nil)
		if result >= 0 {
			return
		} else if result == aeron.BackPressured || result == aeron.AdminAction {
			s.cluster.IdleStrategy().Idle(0)
		} else {
			fmt.Printf("WARNING: OnTakeSnapshot offer failed - result=%v\n", result)
		}
	}
}

func (s *ClusterServerEcho) OnRoleChange(role cluster.Role) {
	fmt.Printf("OnRoleChange - role=%v\n", role)
}

func (s *ClusterServerEcho) OnTerminate(cluster cluster.Cluster) {
	fmt.Printf("OnTerminate - role=%v logPos=%d\n", cluster.Role(), cluster.LogPosition())
}

func (s *ClusterServerEcho) OnNewLeadershipTermEvent(
	leadershipTermId int64,
	logPosition int64,
	timestamp int64,
	termBaseLogPosition int64,
	leaderMemberId int32,
	logSessionId int32,
	timeUnit codecs.ClusterTimeUnitEnum,
	appVersion int32,
) {
	fmt.Printf("OnNewLeadershipTermEvent - leaderTermId=%d logPos=%d time=%d termBase=%d leaderId=%d logSessionId=%d timeUnit=%v appVer=%d\n",
		leadershipTermId, logPosition, timestamp, termBaseLogPosition, leaderMemberId, logSessionId, timeUnit, appVersion)
}

func clusterServerEcho(c *Config) {
	aeronDir := strings.ReplaceAll(c.AeronDir, "<id>", strconv.Itoa(c.ClusterId))
	ctx := aeron.NewContext().AeronDir(aeronDir)

	opts := cluster.NewOptions()
	//if idleStr := os.Getenv("NO_OP_IDLE"); idleStr != "" {

	opts.IdleStrategy = ToIdleStrategy(c.Idle)
	opts.ClusterDir = strings.ReplaceAll(c.ClusterDir, "<id>", strconv.Itoa(c.ClusterId))

	service := &ClusterServerEcho{}
	agent, err := cluster.NewClusteredServiceAgent(ctx, opts, service)
	if err != nil {
		panic(err)
	}

	if err := agent.StartAndRun(); err != nil {
		panic(err)
	}
}
