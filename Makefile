.DEFAULT_GOAL := mac
.PHONY: google

# TAG=gokarna:5000/$(APP):$(VER)

TAG = europe-west4-docker.pkg.dev/websites-394411/webstekjes/activities:$(VER)
GOKARNA = gokarna:5000/activities
# gokarna : TAG = gokarna:5000/activities:$(VER)

makelinux: linux gokarna pushgokarna
pushgoogle: pull tag push

# gokarna: linux docker push


linux:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

docker:
	# @echo $(TAG)
	docker build --tag $(TAG) . 

gokarna: 
	docker build --tag $(GOKARNA)

pushgokarna:
	docker push $(TAG)

pull: 
	docker pull $(GOKARNA))

tag: 
	docker image tag $(GOKARNA) $(TAG)

push:
	docker push $(TAG)

mac:
	go build

clean:
	rm activities

