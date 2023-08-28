# ibsdns-api
InternetBS DNS Updater

## What and why?
updated the repo changing the app to a server application that is accessed using a apikey and a single endpoint. the user should install the app on a server/ computer with static IP as allowed in Internet BS API. Access the endpoint with the configured API key and the app will record your sending ip and update the configured domains to that ip address if different from the last checks.

serveripaddress/update-dns 

with a key passed as a header X-API-KEY

APIkey is user created in config

Strongly suggest that the app is hosted behind a web server like NGINX or Apache with SSL configured and enabled (not secure to snooping otherwise and if someone gets your API they can cause havoc (change your A records))

### Note:
I had to prepay money into Internet.bs to get the API enabled (minimum level in GBP was £8), this is credit against future purchases so using this solution is still free.


## Features

- Reads configuration from a YAML file
- Supports multiple domains
- Can be used behind a reverse proxy
- Checks if the client's IP address has changed before updating the DNS record

## Prerequisites

- Go 1.16 or higher
- A configured DNS server that you have access to (if you want to test DNS functionality)

## Installation

Clone the repository:
git clone https://github.com/xevxx/ibdns.git

Navigate to the project directory:
cd ibdns


Build the application:
go build -o ibdns-api



## Configuration

Before running the server, you'll need to configure it by creating a `config.yaml` file. Here is a sample configuration:
```
Host: "0.0.0.0"
Port: 8080
ApiKey: "your-api-key"
ApiKeyInternetBS: "internetbs supplied api key"
Password: "internetbs supplied password"
Domain: "example.com,example2.com"
```


## Usage

To run the server:
./ibsdns-api



To update the DNS record, make a POST request to the API endpoint:
curl -X POST https://serverip:8080/updateDns

DO USE A WEB SERVER TO ADD SSL

## API Endpoints

- **POST /updateDns**: Update the DNS record. It uses the client's IP address for the update.



# original readme below
## Why?
My buddy Keenan wanted to update home.example.com to point to his IP at home,
however he has a dynamic IP provided to home and so he needed a way to keep it current.
He was insistent on using InternetBS as compared to Route53 as he believed it would be less expensive (free basically)

So I made this for him, though I did it in a way where it can be used by anyone with a InternetBS api key,
and a static host to ssh to and run the update tool from (internetBS limits their IP to a static IP you provide ahead of time)

## HowTo
On server at home (or where the IP will be dynamicly changing) you need to add `grabDynamicIp.sh` to a cronjob, like once every 5 minutes perhaps

```/bin/bash
*/5 * * * * /home/kmosdell/dynamicDNS/grabDynamicIp.sh remoteHost.example.com >> /home/exampleUser/dynamicDNS/dnsLog.txt 2>&1
```

On remotehost add ibsdns binary to `/usr/local/bin/ibsdns` and edit `/opt/ibsdns/config.yaml` for your values.

Ensure passwordless ssh is setup from home to remote host, and that the user can read /opt/ibsdns/config.yaml

## Build
```/bin/bash
export GO111MODULE=on
go mod init
go build
```
