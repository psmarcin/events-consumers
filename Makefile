NAME=events-consumer
PROJECT_ID=events-consumer
REGION=europe-west1
INFRA_DIR=infrastructure
PKG_DIR=pkg
SERVICE_ACCOUNT=root-481@events-consumer.iam.gserviceaccount.com

setProject:
	gcloud config set project $(PROJECT_ID)

pack:
	cd $(PKG_DIR)/message/ && zip ./../../dist/message.zip ./* -r

deploy: pack
	cd $(INFRA_DIR) && make infrastructure-apply

destroy: pack
	cd $(INFRA_DIR) && make infrastructure-destroy
