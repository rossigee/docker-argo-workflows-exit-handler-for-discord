IMAGE_NAME = argo-workflows-exit-handler-for-discord
TAG = v0.3.0
INTERNAL_REPO = harbor.golder.lan/rossigee
PUBLIC_REPO = rossigee

# Build Docker image
build:
	docker build -t $(IMAGE_NAME):$(TAG) .

push-internal:
	docker tag $(IMAGE_NAME):$(TAG) $(INTERNAL_REPO)/$(IMAGE_NAME):$(TAG)
	docker push $(INTERNAL_REPO)/$(IMAGE_NAME):$(TAG)

push-public:
	docker tag $(IMAGE_NAME):$(TAG) $(PUBLIC_REPO)/$(IMAGE_NAME):$(TAG)
	docker push $(PUBLIC_REPO)/$(IMAGE_NAME):$(TAG)

# Build and Push
all: build push-internal push-public

