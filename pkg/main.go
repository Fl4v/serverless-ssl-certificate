package main

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

// User struct
type User struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}
func (u User) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

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

// Constructs the certbot execution command
func certbotExec() string {

	var email string = os.Getenv("EMAIL")
	var domains string = os.Getenv("DOMAINS")

	if email == "" || domains == "" {
		panic(errors.New("email or domains variable not found in environment"))
	}

	EMAIL_OPTIONS := "--email " + email
	DOMAINS_OPTIONS := "--domains " + domains
	OPTIONS := "--manual --preferred-challenges=dns --no-eff-email --agree-tos --manual-auth-hook /app/acme_validation.sh"

	S := " "

	execCommand := "certbot certonly" + S + OPTIONS + S + EMAIL_OPTIONS + S + DOMAINS_OPTIONS

	return execCommand

}

func main() {

	// ctx := context.TODO()

	// // Load default config from ~/.aws/config
	// awsConfig, err := config.LoadDefaultConfig(ctx)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(checkCertificate(awsConfig, ctx, "fl4v.com"))

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}

	user := User{
		Email: os.Getenv("EMAIL"),
		key:   privateKey,
	}

	fmt.Println(user)

	config := lego.NewConfig(&user)

	lego_client, err := lego.NewClient(config)

	if err != nil {
		panic(err)
	}

	custom_validation := domain.upsertRecord()
}
