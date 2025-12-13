package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	models "common/models"
	pb "common/proto_models"
	utils "common/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type statusPagamentoHandler struct{}

type server struct {
	pb.UnimplementedPagamentoServiceServer
}

var gatewayClient pb.GatewayServiceClient

// função que envia lance ganhador para o sistema de pagamento para gerar um link
// depois, recebe o link gerado e coloca na fila de links
func handleLeilaoGanhador(lance models.Lance) string {

	var payment models.Payment = models.CreatePayment("R$", lance.UserID, lance.Value, "http://localhost:8100/status")

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3333/create-payment", bytes.NewReader(payment.ToByteArray()))

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("request failed: %s", err)
	}
	defer res.Body.Close()

	log.Printf("[MS-PAGAMENTO] Novo pagamento gerado: \n%s\n", payment.Print())

	var link models.JustLink

	_ = json.NewDecoder(res.Body).Decode(&link)

	var link_data models.Link = models.CreateLink(link.Link, lance.UserID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = gatewayClient.PublicaLinkPagamento(ctx, &pb.Link{
		UserID: &link_data.UserID,
		Link:   &link_data.Link,
	})
	if err != nil {
		log.Printf("[MS-LEILAO] Erro ao publicar leilao iniciado: %v\n", err)
	}

	log.Printf("[MS-PAGAMENTO] Novo link gerado: \n%s\n", link_data.Print())
	return http.StatusText(http.StatusOK)
}

// função que escuta fila de lances ganhadores
func (s *server) PublicaLanceVencedor(ctx context.Context, in *pb.LanceVencedor) (*pb.PStatus, error) {
	lance := models.CreateLance(in.GetLeilaoID(), in.GetUserID(), in.GetValue())

	response := handleLeilaoGanhador(lance)

	return &pb.PStatus{Status: &response, Lance: in}, nil
}

// função que recebe status de pagamento via requisição http do sistema de pagamento externo
// coloca na fila de status
// func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error)
func (h *statusPagamentoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var status models.StatusPayment

	_ = json.NewDecoder(r.Body).Decode(&status)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := gatewayClient.PublicaStatusPagamento(ctx, &pb.StatusPayment{
		Value:     &status.Value,
		PaymentId: &status.PaymentId,
		UserID:    &status.UserID,
		Status:    &status.Status,
	})
	if err != nil {
		log.Printf("[MS-LEILAO] Erro ao publicar leilao iniciado: %v\n", err)
	}

	log.Printf("[MS-PAGAMENTO] Novo pagamento gerado: \n%s\n", status.Print())
}

func main() {
	lis, err := net.Listen("tcp", ":8101")
	utils.FailOnError(err, "Erro ao escutar a porta :8101")

	s := grpc.NewServer()
	pb.RegisterPagamentoServiceServer(s, &server{})
	fmt.Println("GRPC Server running on localhost:8101")

	gatewayConn, err := grpc.NewClient("localhost:5060", grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.FailOnError(err, "Erro ao conectar ao gateway")
	defer gatewayConn.Close()

	gatewayClient = pb.NewGatewayServiceClient(gatewayConn)

	mux := http.NewServeMux()
	mux.Handle("/status", &statusPagamentoHandler{})
	fmt.Println("Server running on http://localhost:8100")
	go http.ListenAndServe(":8100", mux)

	if err = s.Serve(lis); err != nil {
		utils.FailOnError(err, "Erro ao servir")
	}
}
