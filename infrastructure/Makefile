INFRA_DIR=.
PRODUCTION_CONFIG=$(INFRA_DIR)/production.tfvars
DEVELOPMENT_CONFIG=$(INFRA_DIR)/development.tfvars
GOOGLE_CLOUD_CREDENTIAL_FILE="$(pwd)/credentials.json"

setGoogleCloudCredentials:
	export GOOGLE_CLOUD_KEYFILE_JSON=$(GOOGLE_CLOUD_CREDENTIAL_FILE)

infrastructure-apply-prod: setGoogleCloudCredentials
	terraform apply -var-file="$(PRODUCTION_CONFIG)"

infrastructure-apply-developmemnt: setGoogleCloudCredentials
	terraform apply -var-file="$(DEVELOPMENT_CONFIG)"

infrastructure-destroy: setGoogleCloudCredentials
	terraform destroy
