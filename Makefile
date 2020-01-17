NAME=events-consumer
PROJECT_ID=events-consumer
REGION=europe-west1
PKG_DIR=pkg
DIST_DIR=dist
SERVICE_ACCOUNT=root-481@events-consumer.iam.gserviceaccount.com
INFRA_DIR=./infrastructure
PWD=${shell pwd}

setProject:
	gcloud config set project $(PROJECT_ID)

cleanup-dist:
	rm -rf ${PWD}/${DIST_DIR}/*

pack: cleanup-dist
	cd $(PKG_DIR)/message/ && zip ${PWD}/${DIST_DIR}/message.zip ./* -r -x "*cmd/*" -q
	cd $(PKG_DIR)/get-content/ && zip $(PWD)/${DIST_DIR}/get-content.zip ./* -r -x "*cmd/*" -q
	cd $(PKG_DIR)/process-content/ && zip $(PWD)/${DIST_DIR}/process-content.zip ./* -r -x "*cmd/*" -q

deploy-production: pack
	cd $(INFRA_DIR) && make infrastructure-apply-prod

deploy-development: pack
	cd $(INFRA_DIR) && make infrastructure-apply-development

destroy: pack
	cd $(INFRA_DIR) && make infrastructure-destroy
