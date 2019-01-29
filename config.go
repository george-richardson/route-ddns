package main

type Config struct {
	HostedZones []HostedZone
	Providers   []string
	CycleTime   int
}

type HostedZone struct {
	Id      string
	Records []RecordSet
}

type RecordSet struct {
	TTL        int64
	Name       string
	RecordType string
}
