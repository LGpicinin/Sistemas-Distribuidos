
run: build
	./bin/ms-leilao &
	./bin/ms-lance &

build: ms-leilao ms-lance

ms-leilao: ./ms-leilao/ms-leilao.go
	cd ./ms-leilao; \
	go build -o ../bin/ms-leilao ./ms-leilao.go

ms-lance: ./ms-lance/ms-lance.go
	cd ./ms-lance; \
	go build -o ../bin/ms-lance ./ms-lance.go



