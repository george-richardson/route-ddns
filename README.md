# Install

Grab a [release](https://github.com/george-richardson/route-ddns/releases) or run:
```
go get github.com/george-richardson/route-ddns
```

# Usage
```
route-ddns is a ddns client that updates record sets in AWS route53 when an IP change is detected.

Usage:
  route-ddns [flags]

Flags:
      --config string   config file (default is ./route-ddns.yaml)
  -h, --help            help for route-ddns
```

You must provide a valid config file as described below. AWS must be configured with a default profile or through environment variables as described in the [official AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).

## Example Config
```yaml
# A list of hosted zones to update.
hostedZones: 
    - id: Z1WM8KTX7RPOAI    # The ID of the hosted zone.
      records: 
        - name: example.com # The name of the record to update. Will be created if doesn't already exist.
          recordType: A     # The type of record to set.
          ttl: 60           # TTL to set on the record.
        - name: subdomain.example.com
          recordType: A
          ttl: 60
# (Optional) List of public IP resolvers to use (default: https://api.ipify.org?format=text)
providers: 
  - "https://api.ipify.org?format=text"
  - "http://myip.dnsomatic.com/"
  - "http://icanhazip.com/"
  - "http://bot.whatismyipaddress.com/"
# (Optional) How often in seconds IP should be checked and updated in Route53 (default: 300)
cycleTime: 300
```