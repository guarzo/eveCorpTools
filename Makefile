APP_NAME=zkillanalytics
GROUP_NAME=flygd


.PHONY: all build acr_build update_aca full_deploy run add_secrets clean help
 
build:
	go build -o zkill

acr_build: increment_version
	@VERSION=$$(cat VERSION) && \
	az acr build --image $(GROUP_NAME)/$(APP_NAME):$${VERSION}  --registry $(GROUP_NAME) --build-arg VERSION=$${VERSION} .

restart_aca:
	@export REVISION=$$(az containerapp revision list --name $(APP_NAME) --resource-group $(GROUP_NAME) --query "[].name" -o tsv) && \
	az containerapp revision restart -n $(APP_NAME) -g $(GROUP_NAME) --revision $${REVISION}

full_deploy: acr_build update_aca

update_aca:
	@VERSION=$$(cat VERSION) && \
	az containerapp update -n $(APP_NAME) -g $(GROUP_NAME) --image $(GROUP_NAME).azurecr.io/$(GROUP_NAME)/$(APP_NAME):$${VERSION}

run:
	go run .
 
clean:
	go clean
	rm -rf data
	rm -rf charts
	rm -f zkill

increment_version:
	@VERSION=$$(cat VERSION) && \
	NEW_VERSION=$$(echo $${VERSION} | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}') && \
	echo $${NEW_VERSION} > VERSION && \
	echo "Updated version to $${NEW_VERSION}"
