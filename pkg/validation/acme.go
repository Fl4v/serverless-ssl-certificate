package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Checks the domain and returns the certificate ARN
func checkCertificate(awsConfig aws.Config, ctx context.Context, domain string) string {
	client := acm.NewFromConfig(awsConfig)

	certificates, err := client.ListCertificates(ctx, nil)

	if err != nil {
		panic(err)
	}

	var certArn string

	for _, cert := range certificates.CertificateSummaryList {
		if *cert.DomainName == domain {
			certArn = *cert.CertificateArn
		}
	}

	return certArn
}

// Retrieves the domain hosted zone
func getHostedZoneId(awsConfig aws.Config, ctx context.Context, domain string) (Id string, err error) {
	client := route53.NewFromConfig(awsConfig)

	params := &route53.ListHostedZonesInput{}

	zones, err := client.ListHostedZones(ctx, params)

	if err != nil {
		panic(err)
	}

	var zoneId string

	for _, zone := range zones.HostedZones {
		if strings.Contains(*zone.Name, domain) {
			zoneId = *zone.Id
		}
	}

	if zoneId != "" {
		return zoneId, nil
	}

	return zoneId, fmt.Errorf("no zone found for domain: '%s'", domain)
}

// Creates or updates the acme record
// If the env variable HOSTED_ZONE_ID is not set, it will try and look for it based on the domain name
func upsertAcmeRecord(awsConfig aws.Config, ctx context.Context, domain string, txtValue string) string {

	var hostedZoneId string
	var recordName string = "_acme-challenge." + domain
	var value string = fmt.Sprintf("\"%s\"", txtValue)
	var TTL int64 = 30
	var err error

	hostedZoneId = os.Getenv("HOSTED_ZONE_ID")

	if hostedZoneId == "" {
		hostedZoneId, err = getHostedZoneId(awsConfig, ctx, domain)

		if err != nil {
			panic(err)
		}
	}

	resourceRecord := &types.ResourceRecord{
		Value: &value,
	}
	recordSet := &types.ResourceRecordSet{
		Name: &recordName,
		Type: "TXT",
		ResourceRecords: []types.ResourceRecord{
			*resourceRecord,
		},
		TTL: &TTL,
	}
	change := &types.Change{
		Action:            "UPSERT",
		ResourceRecordSet: recordSet,
	}
	changeBatch := &types.ChangeBatch{
		Changes: []types.Change{
			*change,
		},
	}

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: &hostedZoneId,
	}

	client := route53.NewFromConfig(awsConfig)

	resp, err := client.ChangeResourceRecordSets(ctx, params)

	if err != nil {
		panic(err)
	}

	return *resp.ChangeInfo.Id
}

func main() {

	ctx := context.TODO()

	// Load default config from ~/.aws/config
	awsConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic(err)
	}

	domain := os.Getenv("CERTBOT_DOMAIN")
	token := os.Getenv("CERTBOT_VALIDATION")

	upsertAcmeRecord(awsConfig, ctx, domain, token)

	// Wait N seconds for DNS propagation
	time.Sleep(10 * time.Second)
}
