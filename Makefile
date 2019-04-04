.PHONY: app client service deployment

# DEV
test:
	go test ./... -v

run:
	go run ./cmd/*.go

deps:
	go mod tidy

# BUILD

image:
	gcloud builds submit \
		--project s9-demo \
		--tag gcr.io/s9-demo/tevents

# DEPLOYMENT

deploy:
	kubectl apply -f deployment/service.yaml -n demo

undeploy:
	kubectl delete -f deployment/service.yaml  -n demo

# DEMO

event:
	curl -X POST -H "Content-Type: application/json" -d @sample.json \
		 https://tevents.demo.knative.tech/

local-event:
	curl -XPOST -H "Content-Type: application/json" -d @sample.json \
		 http://localhost:8080/



