package main

import (
	models "common/models"
	pb "common/proto_models"
	utils "common/utils"
	"container/list"
	"context"
	"fmt"
	"log"
	"net"
	"reflect"

	// dto "ms-leilao/DTO"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var activeLeiloes map[string]models.Leilao = make(map[string]models.Leilao)
var leiloesSortedByStart *list.List = list.New()
var leiloesSortedByEnd *list.List = list.New()

var lanceClient pb.LanceServiceClient

type server struct {
	pb.UnimplementedLeilaoServiceServer
}

type createLeilaoHandler struct{}
type listLeilaoHandler struct{}

// função que insere novo leilão em lista ordenada por tempo
func insertionSortOnList(leilaoList *list.List, value models.Leilao, fieldToCompare string) {

	r := reflect.ValueOf(value)
	fieldValue := reflect.Indirect(r).FieldByName(fieldToCompare).Interface().(time.Time)

	if leilaoList.Len() == 0 {
		leilaoList.PushBack(value)
	} else {
		var k *list.Element = nil
		for e := leilaoList.Front(); e != nil; e = e.Next() {
			e_r := reflect.ValueOf(e.Value.(models.Leilao))
			e_value := reflect.Indirect(e_r).FieldByName(fieldToCompare).Interface().(time.Time)

			if e_value.Compare(fieldValue) > 0 {
				k = e
				break
			}
		}
		if k != nil {
			leilaoList.InsertBefore(value, k)
		} else {
			leilaoList.PushBack(value)
		}
	}
}

// função infinita que publica leilão na fila quando ele for iniciado
func publishWhenStarts() {
	for {
		if leiloesSortedByStart.Len() == 0 {
			continue
		}

		first := leiloesSortedByStart.Front()
		firstLeilao := first.Value.(models.Leilao)

		if !firstLeilao.HasStarted() {
			continue
		}

		startDate := firstLeilao.StartDate.Format("2006-01-02T15:04:05.999Z")
		endDate := firstLeilao.EndDate.Format("2006-01-02T15:04:05.999Z")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := lanceClient.PublicaLeilaoIniciado(ctx, &pb.Leilao{
			ID:          &firstLeilao.ID,
			Description: &firstLeilao.Description,
			StartDate:   &startDate,
			EndDate:     &endDate,
		})
		if err != nil {
			log.Printf("[MS-LEILAO] Erro ao publicar leilao iniciado: %v\n", err)
		}
		activeLeiloes[firstLeilao.ID] = firstLeilao
		leiloesSortedByStart.Remove(first)

		log.Printf("[MS-LEILAO] NOVO LEILÃO INICIADO: %s\n\n", firstLeilao.Print())
	}
}

// função infinita que publica leilão na fila quando ele for finalizado
func publishWhenFinishes() {
	for {
		if leiloesSortedByEnd.Len() == 0 {
			continue
		}

		first := leiloesSortedByEnd.Front()
		firstLeilao := first.Value.(models.Leilao)

		if !firstLeilao.HasEnded() {
			continue
		}

		startDate := firstLeilao.StartDate.Format("2006-01-02T15:04:05.999Z")
		endDate := firstLeilao.EndDate.Format("2006-01-02T15:04:05.999Z")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := lanceClient.PublicaLeilaoFinalizado(ctx, &pb.Leilao{
			ID:          &firstLeilao.ID,
			Description: &firstLeilao.Description,
			StartDate:   &startDate,
			EndDate:     &endDate,
		})
		if err != nil {
			log.Printf("[MS-LEILAO] Erro ao publicar leilao iniciado: %v\n", err)
		}
		leiloesSortedByEnd.Remove(first)
		delete(activeLeiloes, firstLeilao.ID)

		log.Printf("[MS-LEILAO] NOVO LEILÃO FINALIZADO %s\n\n", firstLeilao.Print())
	}
}

// recebe requisição http do gateway para criação de novo leilão
// chama função para inserir na lista de ordenada por tempo de início e na de tempo por fim
func (s *server) Create(ctx context.Context, in *pb.LLeilao) (*pb.LStatus, error) {
	startDate, _ := time.Parse("2006-01-02T15:04:05.999Z", in.GetStartDate())
	endDate, _ := time.Parse("2006-01-02T15:04:05.999Z", in.GetEndDate())

	leilao := models.CreateLeilao(in.GetID(), in.GetDescription(), startDate, endDate)

	insertionSortOnList(leiloesSortedByStart, leilao, "StartDate")
	insertionSortOnList(leiloesSortedByEnd, leilao, "EndDate")

	status := http.StatusText(http.StatusCreated)
	return &pb.LStatus{Status: &status, Leilao: in}, nil
}

// recebe requisição http do gateway para listar leilões ativos
// envia resposta por http
func (s *server) List(in *pb.Empty, stream pb.LeilaoService_ListServer) error {
	for _, activeLeilao := range activeLeiloes {
		startDate := activeLeilao.StartDate.Format("2006-01-02T15:04:05.999Z")
		endDate := activeLeilao.EndDate.Format("2006-01-02T15:04:05.999Z")
		if err := stream.Send(&pb.LLeilao{
			ID:          &activeLeilao.ID,
			Description: &activeLeilao.Description,
			StartDate:   &startDate,
			EndDate:     &endDate,
		}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8090")
	utils.FailOnError(err, "Erro ao escutar a porta :8090")

	s := grpc.NewServer()
	pb.RegisterLeilaoServiceServer(s, &server{})
	fmt.Println("Server running on http://localhost:8090")

	lanceConn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.FailOnError(err, "Erro ao conectar ao lance")
	defer lanceConn.Close()

	lanceClient = pb.NewLanceServiceClient(lanceConn)

	go publishWhenStarts()
	go publishWhenFinishes()

	if err = s.Serve(lis); err != nil {
		utils.FailOnError(err, "Erro ao servir")
	}
}
