
run: build
	./bin/ms-leilao &
	./bin/ms-lance &
# 	./bin/ms-pagamento &

build: proto ms-leilao ms-lance ms-pagamento

proto: 
	cd common/proto; \
	protoc --go_out=.. \
	--go-grpc_out=.. \
	./*.proto

ms-leilao: ./ms-leilao/ms-leilao.go
	cd ./ms-leilao; \
	go mod tidy; \
	go build -o ../bin/ms-leilao ./ms-leilao.go

ms-lance: ./ms-lance/ms-lance.go
	cd ./ms-lance; \
	go mod tidy; \
	go build -o ../bin/ms-lance ./ms-lance.go

# ms-pagamento: ./ms-pagamento/ms-pagamento.go
# 	cd ./ms-pagamento; \
# 	go mod tidy; \
# 	go build -o ../bin/ms-pagamento ./ms-pagamento.go



