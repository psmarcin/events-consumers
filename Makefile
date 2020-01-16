NAME=events-consumer
PROJECT_ID=events-consumer
REGION=europe-west1
PKG_DIR=pkg
SERVICE_ACCOUNT=root-481@events-consumer.iam.gserviceaccount.com
INFRA_DIR=./infrastructure

setProject:
	gcloud config set project $(PROJECT_ID)

pack:
	cd $(PKG_DIR)/message/ && zip ./../../dist/message.zip ./* -r
	cd $(PKG_DIR)/get-content/ && zip ./../../dist/get-content.zip ./* -r
	cd $(PKG_DIR)/process-content/ && zip ./../../dist/process-content.zip ./* -r

deploy-production: pack
	cd $(INFRA_DIR) && make infrastructure-apply-prod

deploy-development: pack
	cd $(INFRA_DIR) && make infrastructure-apply-development

destroy: pack
	cd $(INFRA_DIR) && make infrastructure-destroy
