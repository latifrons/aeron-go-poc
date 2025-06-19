package main

type Config struct {
	ProfilerEnabled *bool
	Channel         *string
	StreamId        *int
	Messages        *int
	Size            *int
	LoggingOn       *bool
	Timeout         *int
	AeronDir        *string
}
