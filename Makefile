.PHONY: cover start test test-integration

# get newest image to deploy
export image := `aws lightsail get-container-images --service-name canvas --label canvas | jq -r '.containerImages[0].image'`

build:
	docker build -t canvas .

cover:
	go tool cover -html=cover.out

start:
	go run cmd/server/*.go

test:
	go test -coverprofile=cover.out -short ./...

test-integration:
	go test -coverprofile=cover.out -p 1 ./...

deploy:
	aws lightsail push-container-image --service-name deeler --label app --image deeler
	aws lightsail create-container-service-deployment --service-name deeler \
		--containers '{"app":{"image":":deeler.app.1","environment":{"HOST":"","PORT":"8080","LOG_ENV":"production"},"ports":{"8080":"HTTP"}}}' \
		--public-endpoint '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}'
