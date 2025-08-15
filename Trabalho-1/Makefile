
run-all: build-all
	./bin/ms-notificacao &
	./bin/client &
	./bin/ms-leilao &
	./bin/ms-lance &

build-all: build-ms-notificacao build-client build-ms-leilao build-ms-lance

build-ms-notificacao: ./ms-notificacao/ms-notificacao.go
	cd ./ms-notificacao; \
	go build -o ../bin/ms-notificacao ./ms-notificacao.go

build-client: ./client/client.go
	cd ./client; \
	go build -o ../bin/client ./client.go

build-ms-leilao: ./ms-leilao/ms-leilao.go
	cd ./ms-leilao; \
	go build -o ../bin/ms-leilao ./ms-leilao.go

build-ms-lance: ./ms-lance/ms-lance.go
	cd ./ms-lance; \
	go build -o ../bin/ms-lance ./ms-lance.go



