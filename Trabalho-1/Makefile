
run: build
	./bin/ms-notificacao &
	./bin/ms-leilao &
	./bin/ms-lance &

build: ms-notificacao client ms-leilao ms-lance

ms-notificacao: ./ms-notificacao/ms-notificacao.go
	cd ./ms-notificacao; \
	go build -o ../bin/ms-notificacao ./ms-notificacao.go

client: ./client/client.go
	cd ./client; \
	go build -o ../bin/client ./client.go

ms-leilao: ./ms-leilao/ms-leilao.go
	cd ./ms-leilao; \
	go build -o ../bin/ms-leilao ./ms-leilao.go

ms-lance: ./ms-lance/ms-lance.go
	cd ./ms-lance; \
	go build -o ../bin/ms-lance ./ms-lance.go



