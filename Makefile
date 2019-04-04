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
	curl -H "Content-Type: application/json" \
		 -X POST --data "{ \
			\"specversion\": \"0.2\", \
			\"type\": \"github.com.mchmarny.knative-ws-example.message\", \
			\"source\": \"https://github.com/mchmarny/knative-ws-example\", \
			\"id\": \"6CC459AE-D75D-4556-8C14-CD1ED5D95AE7\", \
			\"time\": \"2019-02-13T17:31:00Z\", \
			\"contenttype\": \"text/plain\", \
			\"data\": \"This is my sample message\" \
		}" \
		http://localhost:8080/?token=${KNOWN_PUBLISHER_TOKEN}

