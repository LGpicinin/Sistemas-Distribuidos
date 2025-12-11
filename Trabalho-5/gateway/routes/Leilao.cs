using System.Text.Json;
using Grpc.Net.Client;
using GrpcLeilao;
using Classes;



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
            var channel = GrpcChannel.ForAddress("https://localhost:8090");
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

            using var response = await ms_leilao.Create(pbLeilao);

            // using var response = await httpClient.PostAsync($"{MSLeilaoAddress}/create", content);

            httpContext.Response.StatusCode = (int)response.Status;
            if (response.StatusCode.ToString() == "Created")
            {
                var leilaoData = JsonSerializer.Deserialize<LeilaoData>(body);
                notificacao.AddLeilao(leilaoData!.id);
            }
            httpContext.Response.ContentType = "application/json";
            // var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(await content.ReadAsStringAsync());
        }

        public async Task ListLeiloes(HttpContext httpContext)
        {
            using var activeLeiloesResponse = await ms_leilao.List(null);
            var userId = httpContext.Request.Query["userId"];
            
            // string responseBody = await activeLeiloesResponse.Content.ReadAsStringAsync();
            // var leiloes = JsonSerializer.Deserialize<List<LeilaoData>>(responseBody);
            List<LeilaoDataPlus> leiloesPlus = new List<LeilaoDataPlus>();

            for (int i = 0; i < activeLeiloesResponse?.Count; i++)
            {
                var leilaoData = new LeilaoData(
                    activeLeiloesResponse[i].ID,
                    activeLeiloesResponse[i].Description,
                    activeLeiloesResponse[i].StartDate,
                    activeLeiloesResponse[i].EndDate
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

            httpContext.Response.StatusCode = (int)activeLeiloesResponse.StatusCode;
            httpContext.Response.ContentType = activeLeiloesResponse.Content.Headers.ContentType?.ToString() ?? "application/json";
            await httpContext.Response.WriteAsync(leiloesSerialized);
        }
    }
}