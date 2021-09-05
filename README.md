# Serverless SSL Certificate with Certbot

Lambda that automatically creates/updates a SSL Certificate for your domain to [Certificate Manager](https://aws.amazon.com/certificate-manager/).

The main reason for this repo is because I want to learn more about the AWS Cloud ecosystem.

- Get a free Certificate from [Let's Encrypt]()

### Env Variables

| Name | Type | Description |
| ---- | ---- | ---- |
| DOMAINS_NAME | `String` | Comma separated values of the domain(s) |
| HOSTED_ZONE_ID | `String` | Needed for |
| EMAIL | `String` | |
