NAME=events-consumer
PROJECT_ID=events-consumer
REGION=europe-west1
PKG_DIR=pkg
DIST_DIR=dist
INFRA_DIR=./infrastructure
PWD=${shell pwd}

.PHONY: setProject
setProject:
	gcloud config set project $(PROJECT_ID)

.PHONY: cleanup-dist
cleanup-dist:
	rm -rf ${PWD}/${DIST_DIR}/*

.PHONY: pack
pack: cleanup-dist
	cd $(PKG_DIR)/message/ && zip ${PWD}/${DIST_DIR}/message.zip ./* -r -x "*cmd/*" -q
	cd $(PKG_DIR)/get-content/ && zip $(PWD)/${DIST_DIR}/get-content.zip ./* -r -x "*cmd/*" -q
	cd $(PKG_DIR)/process-content/ && zip $(PWD)/${DIST_DIR}/process-content.zip ./* -r -x "*cmd/*" -q
	cd $(PKG_DIR)/get-jobs/ && zip $(PWD)/${DIST_DIR}/get-jobs.zip ./* -r -x "*cmd/*" -q

.PHONY: deploy-production
deploy-production: pack
	cd $(INFRA_DIR) && make infrastructure-apply-prod

.PHONY: deploy-development
deploy-development: pack
	cd $(INFRA_DIR) && make infrastructure-apply-development

.PHONY: destroy
destroy: pack
	cd $(INFRA_DIR) && make infrastructure-destroy

.PHONY: get-dependencies
get-dependencies:
	find ./pkg/ -maxdepth 1 -type d \( ! -name "pkg" \) -exec bash -c "cd '{}' && go get ." \;

.PHONY: test
test:
	find ./pkg/ -maxdepth 1 -type d \( ! -name "pkg" \) -exec bash -c "cd '{}' && go test ./..." \;

