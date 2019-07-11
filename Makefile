.PHONY: app client service deployment

# DEV
test:
	go test ./... -v

run:
	go run ./cmd/*.go

mod:
	go mod tidy
	go mod vendor

# BUILD
image: mod
	gcloud builds submit \
		--project cloudylabs-public \
		--tag gcr.io/cloudylabs-public/tevents:0.1.1

# DEPLOYMENT
service:
	kubectl apply -f deployment/service.yaml -n demo

trigger:
	kubectl apply -f deployment/trigger.yaml -n demo

cleanup:
	kubectl delete -f deployment/service.yaml  -n demo

# DEMO

event:
	curl -X POST -H "Content-Type: application/json" -d @sample.json \
		 https://tweets.demo.knative.tech/

local-event:
	curl -XPOST -H "Content-Type: application/json" -d @sample.json \
		 http://localhost:8080/



