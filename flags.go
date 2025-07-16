package main

type Config struct {
	ProfilerEnabled  bool
	ChannelOut       string
	ChannelIn        string
	StreamIdIn       int
	StreamIdOut      int
	Messages         int
	Size             int
	LoggingOn        bool
	Timeout          int
	AeronDir         string
	Idle             string
	ClusterDir       string
	ClusterId        int
	IngressChannel   string
	IngressEndpoints string
	IngressStreamId  int
	EgressChannel    string
	EgressStreamId   int
}
