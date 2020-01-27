# Events Consumer 

Notify about content changes. Use Google Cloud Function, Google PubSub, Google Scheduler and Google Firestore to send notification via Telegram about content changes.

#### Battery included
💻 (code) + 📖 (infrastructure) = ❤️

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) - Infrastructure as a code
* [Go](https://golang.org/dl/) - Function logic
* [Google Cloud CLI](https://cloud.google.com/sdk/docs/quickstarts/) - Interact with GCP


### Installing

A step by step series of examples that tell you how to get a development env running

Run function locally: 
```sh
cd ./pkg/get-content
go run ./cmd/main.go
```

Server will start listing on http://localhost:8080

Make POST request with event payload: 
```sh
curl --location --request POST 'http://localhost:8080' \
--header 'Content-Type: text/plain' \
--data-binary 'request-payload.json'
```

Example payload: 
```json
{
	"data": "eyJjb21tYW5kIjogImN1cmwgaHR0cHM6Ly9nb29nbGUuY29tIiwic2VsZWN0b3IiOiAiYm9keSIsImNvbnRlbnQiOiAib2xkIGNvbnRlbnQiLCJuYW1lIjoiam9iIG5hbWUifQo"
}
```

"data" should be base64 encoded. So JSON example decoded version looks like this:

```json
{
	"data": {
		"command": "curl https://google.com",
		"selector": "body",
		"content": "old content",
		"name":"job name"
	}
}
```



## Running the tests

To run local test go to function directory and run `go test ./...`


### And coding style tests

To run test style formatter run `gofmt -s -w ./pkg/`

## Deployment

We use `terraform` for handling all infrastructure tasks including deployment. 

## Deploy

Set `credential.json` file for authorization in GCP: 
```sh
export GOOGLE_CLOUD_KEYFILE_JSON="./events-consumers/infrastructure/credentials.json"
```


```shell script
make deploy-production
```

### Config files
1. `./infrastructure/production.tfvars`:
```terraform
# https://www.terraform.io/docs/providers/google/index.html
project = ""
region = ""
zone = ""
# app config
get_content_name=""
get_jobs_name=""
telegram_api_key=""
telegram_channel_id=""
```
1. `./infrastructure/development.tfvars`:
```terraform
# https://www.terraform.io/docs/providers/google/index.html
project = ""
region = ""
zone = ""
# app config
get_content_name=""
get_jobs_name=""
telegram_api_key=""
telegram_channel_id=""
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


