REPO_URL = github.com/invincibleinfra/grafana-backup-restore

# Build using golang docker container for reproducibility
# this image also includes things we need like glide
GOLANG_IMAGE=amaysim/golang:1.8.3
CONTAINER_GOPATH=/go
CONTAINER_SOURCE_DIR=$(CONTAINER_GOPATH)/src/$(REPO_URL)

build:
	docker run \
		-e "CGO_ENABLED=0" \
		-v $(CURDIR):$(CONTAINER_SOURCE_DIR) \
		-w $(CONTAINER_SOURCE_DIR) \
		--rm $(GOLANG_IMAGE) \
		/bin/bash -c "glide install && go build -o grafana-backup-restore"
