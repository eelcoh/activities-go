.DEFAULT_GOAL := local
.PHONY: google

# TAG=gokarna:5000/$(APP):$(VER)

TAG = europe-west4-docker.pkg.dev/websites-394411/webstekjes/activities:$(VER)
GOKARNA = gokarna:5000/activities
# gokarna : TAG = gokarna:5000/activities:$(VER)

gokarna: linux image togokarna
google: pull tag push

# gokarna: linux docker push


linux:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags '-w' -o activities ./cmd/activities

docker:
	# @echo $(TAG)
	docker build --tag $(TAG) . 

image: 
	docker build --tag $(GOKARNA) .

togokarna:
	docker push $(GOKARNA)

pull: 
	docker pull $(GOKARNA)

tag: 
	docker image tag $(GOKARNA) $(TAG)

push:
	docker push $(TAG)

local:
	go build -o activities ./cmd/activities

clean:
	rm -f activities
