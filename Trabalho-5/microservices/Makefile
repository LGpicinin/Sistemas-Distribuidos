
run: build
	./bin/ms-leilao &
	./bin/ms-lance &
	./bin/ms-pagamento &

build: ms-leilao ms-lance ms-pagamento

ms-leilao: ./ms-leilao/ms-leilao.go
	cd ./ms-leilao; \
	go build -o ../bin/ms-leilao ./ms-leilao.go

ms-lance: ./ms-lance/ms-lance.go
	cd ./ms-lance; \
	go build -o ../bin/ms-lance ./ms-lance.go

ms-pagamento: ./ms-pagamento/ms-pagamento.go
	cd ./ms-pagamento; \
	go build -o ../bin/ms-pagamento ./ms-pagamento.go



