all: build-local

build-local:
	make clean
	make app-tier-network
	make etcd-server
	make web-app

clean:
	-docker stop etcd-server
	-docker rm etcd-server
	-docker stop url-shortener
	-docker rm url-shortener
	-docker network rm app-tier

app-tier-network:
	docker network create app-tier --driver bridge

etcd-server:
	docker run -d --name etcd-server --network app-tier --publish 2379:2379 --publish 2380:2380 --env ALLOW_NONE_AUTHENTICATION=yes --env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 bitnami/etcd:latest

web-app:
	docker build . -t laupski/url-shortener
	docker run -d --name url-shortener --network app-tier --publish 8000:8000 laupski/url-shortener
