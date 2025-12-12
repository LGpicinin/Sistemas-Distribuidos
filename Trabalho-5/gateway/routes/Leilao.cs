using System.Text.Json;
using Grpc.Net.Client;
using GrpcLeilao;
using Classes;
using System.Threading.Channels;
using System;
using System.Net.Http;
using System.Threading.Tasks;



namespace Routes
{
    class Leilao
    {
        private static readonly HttpClient httpClient = new HttpClient();
        private static string MSLeilaoAddress = "http://localhost:8090";

        private Notificacao notificacao;
        private LeilaoService.LeilaoServiceClient ms_leilao;

        public Leilao(Notificacao not)
        {
            notificacao = not;
        }

        public async Task ConnectCreateChannel()
        {
            AppContext.SetSwitch(
                "System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);
            // var options = new GrpcChannelOptions();
            // options.Credentials = Grpc.Core.ChannelCredentials.Insecure;
            var channel = GrpcChannel.ForAddress("http://localhost:8090");
            // var channel = GrpcChannel.ForAddress("localhost:8090",  new GrpcChannelOptions
            // {
            //     Credentials = ChannelCredentials.Insecure
            // });
            ms_leilao = new LeilaoService.LeilaoServiceClient(channel);
        }

        

        public void SetupRoutes(WebApplication app)
        {
            app.MapPost("/leilao/create", NewLeilao);
            app.MapGet("/leilao/list", ListLeiloes);
        }

        public async Task NewLeilao(HttpContext httpContext)
        {
            string body = await Utils.HTTPHelper.getRequestBody(httpContext);
            using var content = await Utils.HTTPHelper.getRequestContentFromBody(body, httpContext);

            var leilao = JsonSerializer.Deserialize<LeilaoData>(body);

            var pbLeilao = new LLeilao
            {
                ID = leilao.id,
                Description = leilao.description,
                StartDate = leilao.start_date.ToString(),
                EndDate = leilao.end_date.ToString()
            };

            var response = ms_leilao.Create(pbLeilao);

            // using var response = await httpClient.PostAsync($"{MSLeilaoAddress}/create", content);
            var num_status = 0;
            if (response.Status == "Created")
            {
                num_status = 201;
            } else
            {
                num_status = 400;
            }

            httpContext.Response.StatusCode = num_status;
            if (response.Status.ToString() == "Created")
            {
                var leilaoData = JsonSerializer.Deserialize<LeilaoData>(body);
                notificacao.AddLeilao(leilaoData!.id);
            }
            httpContext.Response.ContentType = "application/json";
            // var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(body);
        }

        public async Task ListLeiloes(HttpContext httpContext)
        {
            CancellationTokenSource source = new CancellationTokenSource();
            CancellationToken token = source.Token;
            var userId = httpContext.Request.Query["userId"];
            // string responseBody = await activeLeiloesResponse.Content.ReadAsStringAsync();
            // var leiloes = JsonSerializer.Deserialize<List<LeilaoData>>(responseBody);
            List<LeilaoDataPlus> leiloesPlus = new List<LeilaoDataPlus>();

            var activeLeiloesResponse = ms_leilao.List(null);
            while (await activeLeiloesResponse.ResponseStream.MoveNext(token))
            // await foreach (var current in call.ResponseStream.ReadAllAsync())
            {
                var current = activeLeiloesResponse.ResponseStream.Current;
                var leilaoData = new LeilaoData(
                    current.ID,
                    current.Description,
                    DateTime.Parse(current.StartDate),
                    DateTime.Parse(current.EndDate)
                );
                var leilaoPlus = new LeilaoDataPlus(leilaoData, false);

                if (!notificacao.InterestLists.ContainsKey(leilaoData.id))
                {
                    notificacao.InterestLists.Add(leilaoData.id, new Notificacao.InterestList());
                }
                var clients = notificacao.InterestLists[leilaoData.id].ClientIds;
                if (clients.ContainsKey(userId!))
                {
                    leilaoPlus.notificar = true;
                }
                leiloesPlus.Add(leilaoPlus);
            }

            var leiloesSerialized = JsonSerializer.Serialize(leiloesPlus);

            httpContext.Response.StatusCode = 200;
            httpContext.Response.ContentType = "application/json";
            await httpContext.Response.WriteAsync(leiloesSerialized);
        }
    }
}