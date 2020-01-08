GOOGLE_CLOUD_CREDENTIAL_FILE="$(pwd)/infrastructure/credentials.json"

setGoogleCloudCredentials:
	export GOOGLE_CLOUD_KEYFILE_JSON=$(GOOGLE_CLOUD_CREDENTIAL_FILE)

infrastructure-apply: setGoogleCloudCredentials
	terraform apply

infrastructure-destroy: setGoogleCloudCredentials
	terraform destroy
