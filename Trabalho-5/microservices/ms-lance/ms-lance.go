package main

import (
	models "common/models"
	utils "common/utils"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "common/proto_models"
	"net/http"

	"google.golang.org/grpc"
)

var gatewayClient pb.GatewayServiceClient

type server struct {
	pb.UnimplementedLanceServiceServer
}

var activeLeiloes map[string]models.ActiveLeilao = make(map[string]models.ActiveLeilao)

// função que verifica se o novo lance é valido ou não
// se for, posta na fila de lances validados; caso contrário, posta na de inválidados
func handleLanceCandidate(lanceCandidate models.Lance) string {
	activeLeilao, ok := activeLeiloes[lanceCandidate.LeilaoID]
	if !ok {
		log.Printf("[MS-LANCE] Erro ao acessar leilão ativo: %v\n", lanceCandidate.LeilaoID)
		return http.StatusText(http.StatusNotFound)
	}

	if lanceCandidate.Value <= activeLeilao.LastValidLance.Value {
		log.Printf("[MS-LANCE] Lance não válido: \n%s\n", lanceCandidate.Print())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := gatewayClient.PublicaLanceInvalido(ctx, &pb.GLance{
			LeilaoID: &lanceCandidate.LeilaoID,
			UserID:   &lanceCandidate.UserID,
			Value:    &lanceCandidate.Value,
		})
		if err != nil {
			log.Printf("[MS-LANCE] Erro ao publicar lance inválido: %v\n", err)
		}

		return response.GetStatus()
	}

	activeLeilao.LastValidLance = lanceCandidate

	activeLeiloes[lanceCandidate.LeilaoID] = activeLeilao

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := gatewayClient.PublicaLanceValido(ctx, &pb.GLance{
		LeilaoID: &lanceCandidate.LeilaoID,
		UserID:   &lanceCandidate.UserID,
		Value:    &lanceCandidate.Value,
	})
	if err != nil {
		log.Printf("[MS-LANCE] Erro ao publicar lance válido: %v\n", err)
	}

	log.Printf("[MS-LANCE] Novo lance validado: \n%s\n", lanceCandidate.Print())
	return response.GetStatus()
}

func (s *server) Create(ctx context.Context, in *pb.Lance) (*pb.Status, error) {
	lance := models.CreateLance(in.GetLeilaoID(), in.GetUserID(), in.GetValue())

	resultStatus := handleLanceCandidate(lance)

	return &pb.Status{Status: &resultStatus, Lance: in}, nil
}

// função que salva novo leilão
func handleLeilaoIniciado(leilao models.Leilao) string {
	activeLeiloes[leilao.ID] = models.ActiveLeilao{
		Leilao: leilao,
	}

	log.Printf("[MS-LANCE] NOVO LEILÃO INICIADO: \n%s\n", leilao.Print())
	return http.StatusText(http.StatusOK)
}

// função que escuta fila de leilões iniciados
func (s *server) PublicaLeilaoIniciado(ctx context.Context, in *pb.Leilao) (*pb.LStatus, error) {
	startDate, _ := time.Parse(time.RFC1123Z, in.GetStartDate())
	endDate, _ := time.Parse(time.RFC1123Z, in.GetEndDate())
	leilao := models.CreateLeilao(in.GetID(), in.GetDescription(), startDate, endDate)

	response := handleLeilaoIniciado(leilao)
	return &pb.LStatus{Status: &response, Leilao: in}, nil
}

// função que remove leilão e publica lance vencedor na fila
func handleLeilaoFinalizado(leilao models.Leilao) string {
	var response = ""

	activeLeilao, ok := activeLeiloes[leilao.ID]
	if ok {
		lastLance := activeLeilao.LastValidLance
		if lastLance != (models.Lance{}) {

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			response_, err := gatewayClient.PublicaLeilaoVencedor(ctx, &pb.GLance{
				LeilaoID: &lastLance.LeilaoID,
				UserID:   &lastLance.UserID,
				Value:    &lastLance.Value,
			})
			response = response_.GetStatus()
			if err != nil {
				log.Printf("[MS-LANCE] Erro ao publicar lance válido: %v\n", err)
			}

			log.Printf("[MS-LANCE] NOVO VENCEDOR: \n%s\n", lastLance.Print())
		}

		delete(activeLeiloes, leilao.ID)
		log.Printf("[MS-LANCE] LEILÃO FINALIZADO: \n%s\n", leilao.Print())
	}

	return response
}

// função que escuta fila de leilões finalisados
func (s *server) PublicaLeilaoFinalizado(ctx context.Context, in *pb.Leilao) (*pb.LStatus, error) {
	startDate, _ := time.Parse(time.RFC1123Z, in.GetStartDate())
	endDate, _ := time.Parse(time.RFC1123Z, in.GetEndDate())
	leilao := models.CreateLeilao(in.GetID(), in.GetDescription(), startDate, endDate)

	response := handleLeilaoFinalizado(leilao)
	return &pb.LStatus{Status: &response, Leilao: in}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	utils.FailOnError(err, "Erro ao escutar a porta :8080")

	s := grpc.NewServer()
	pb.RegisterLanceServiceServer(s, &server{})
	fmt.Println("Server running on http://localhost:8080")

	gatewayConn, err := grpc.NewClient("localhost:5060")
	utils.FailOnError(err, "Erro ao conectar ao gateway")
	defer gatewayConn.Close()

	gatewayClient = pb.NewGatewayServiceClient(gatewayConn)

	if err = s.Serve(lis); err != nil {
		utils.FailOnError(err, "Erro ao servir")
	}
}
