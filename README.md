## Usage
```
route-ddns is a ddns client that updates record sets in AWS route53 when an IP change is detected.

Usage:
  route-ddns [flags]

Flags:
      --config string   config file (default is ./route-ddns.yaml)
  -h, --help            help for route-ddns
```

## Example Config
```
hostedZones: 
    - id: Z1WM8KTX7RPOAI
      records: 
        - name: richardsteen.house
          recordType: A
          ttl: 60
        - name: awesome.richardsteen.house
          recordType: A
          ttl: 60
providers: 
  - "http://myip.dnsomatic.com/"
  - "http://icanhazip.com/"
  - "http://bot.whatismyipaddress.com/"
  - "https://api.ipify.org?format=text"
```