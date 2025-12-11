using GrpcGateway;
using Grpc.Core;
using Routes;
using Classes;

namespace GrpcGateway.Services;


class Gateway : GatewayService.GatewayServiceBase
{
        // required string LeilaoID = 1;
        // required string UserID = 2;
        // required float Value = 3;
        private Notificacao notificacao;

        public Gateway(Notificacao not)
        {
            notificacao = not;
        }


        public override Task<GStatus> PublicaLanceInvalido(GLance request,
            ServerCallContext context)
        {
            // logger.LogInformation("Saying hello to {Name}", request.Name);
        

            var lanceData = new LanceData();
            var lancePlus = new LanceDataType();

            lanceData.leilao_id = request.LeilaoID;
            lanceData.user_id = request.UserID;
            lanceData.value = request.Value;

            lancePlus.lance = lanceData;
            lancePlus.type = "lance_invalidado";
            
            notificacao.ConsumeLanceEvents(lancePlus);
    

            return Task.FromResult(new GStatus 
            {
                Status = "FELICIDADE!!!!!"
            });
        }

        public override Task<GStatus> PublicaLanceValido(GLance request,
            ServerCallContext context)
        {
            // logger.LogInformation("Saying hello to {Name}", request.Name);
        

            var lanceData = new LanceData();
            var lancePlus = new LanceDataType();

            lanceData.leilao_id = request.LeilaoID;
            lanceData.user_id = request.UserID;
            lanceData.value = request.Value;

            lancePlus.lance = lanceData;
            lancePlus.type = "lance_validado";
            
            notificacao.ConsumeLanceEvents(lancePlus);
    

            return Task.FromResult(new GStatus 
            {
                Status = "FELICIDADE!!!!!"
            });
        }

        public override Task<GStatus> PublicaLeilaoVencedor(GLance request,
            ServerCallContext context)
        {
            // logger.LogInformation("Saying hello to {Name}", request.Name);
        

            var lanceData = new LanceData();
            var lancePlus = new LanceDataType();

            lanceData.leilao_id = request.LeilaoID;
            lanceData.user_id = request.UserID;
            lanceData.value = request.Value;

            lancePlus.lance = lanceData;
            lancePlus.type = "leilao_vencedor";
            
            notificacao.ConsumeLanceEvents(lancePlus);
    

            return Task.FromResult(new GStatus 
            {
                Status = "FELICIDADE!!!!!"
            });
        }

        public override Task<GStatus> PublicaLinkPagamento(Link request,
            ServerCallContext context)
        {
            // logger.LogInformation("Saying hello to {Name}", request.Name);
        

            var linkData = new LinkData();
            var linkPlus = new LinkDataType();

            linkData.clientId = request.UserID;
            linkData.link = request.Link_;

            linkPlus.linkData = linkData;
            linkPlus.type = "link_pagamento";
            
            notificacao.ConsumeLinkEvents(linkPlus);
    

            return Task.FromResult(new GStatus 
            {
                Status = "FELICIDADE!!!!!"
            });
        }

        public override Task<GStatus> PublicaStatusPagamento(StatusPayment request,
            ServerCallContext context)
        {

            var statusData = new StatusData();
            var statusPlus = new StatusDataType();

            statusData.clientId = request.UserID;
            statusData.paymentId = request.PaymentId;
            statusData.value = request.Value;
            statusData.status = request.Status;

            statusPlus.statusData = statusData;
            statusPlus.type = "status_pagamento";
            
            notificacao.ConsumeStatusEvents(statusPlus);
    

            return Task.FromResult(new GStatus 
            {
                Status = "FELICIDADE!!!!!"
            });
        }
}
