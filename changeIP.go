package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	try "gopkg.in/matryer/try.v1"
)

func changeIP(newIP string, cfg Config) error {
	const maxAttempts = 5
	err := try.Do(func(attempt int) (bool, error) {
		err := tryChangeIP(newIP, cfg)
		if err != nil {
			log.Print(fmt.Sprintf("WARNING: Set DNS attempt %v/%v: %v", attempt, maxAttempts, err))
		}
		return attempt < maxAttempts, err
	})
	return err
}

func tryChangeIP(newIP string, cfg Config) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	svc := route53.New(sess)
	for _, zone := range cfg.HostedZones {
		err = processHostedZone(newIP, zone, *svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func processHostedZone(newIP string, zone HostedZone, svc route53.Route53) error {
	// Update all records in Route53 hostedzone
	var inputRecords []*route53.Change
	for _, record := range zone.Records {
		inputRecords = append(inputRecords, &route53.Change{
			Action: aws.String("UPSERT"),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: aws.String(record.Name),
				ResourceRecords: []*route53.ResourceRecord{
					{
						Value: aws.String(newIP),
					},
				},
				TTL:  aws.Int64(record.TTL),
				Type: aws.String(record.RecordType),
			},
		})
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: inputRecords,
			Comment: aws.String("Managed by route-ddns"),
		},
		HostedZoneId: aws.String(zone.Id),
	}
	_, err := svc.ChangeResourceRecordSets(input)
	return err
}
