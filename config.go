package main

type Config struct {
	HostedZones []HostedZone
	Providers   []string
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
