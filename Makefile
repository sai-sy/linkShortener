# Image names
DEV_IMAGE = app-dev
PROD_IMAGE = app-prod

# Build the dev image
dev-build:
	docker build --target dev -t $(DEV_IMAGE) .

# Run the dev container with hot reload
dev-run:
	docker run --rm -it \
		-p 8080:8080 \
		-v $(PWD):/app \
		$(DEV_IMAGE)

# Convenience target: build + run
dev: dev-build dev-run

