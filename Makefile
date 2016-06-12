#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testzombie .
default: testzombie
	docker build -f Dockerfile -t testzombie .

testzombie:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o testzombie

clean:
	rm testzombie
