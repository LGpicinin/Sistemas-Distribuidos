using System.Text.Json;

namespace Routes
{
    class Leilao
    {
        private static readonly HttpClient httpClient = new HttpClient();
        private static string MSLeilaoAddress = "http://localhost:8090";

        private Notificacao notificacao;

        public Leilao(Notificacao not)
        {
            notificacao = not;
        }

        public class LeilaoData
        {
            public string id { get; set; }
            public string description { get; set; }
            public DateTime start_date { get; set; }
            public DateTime end_date { get; set; }

            public LeilaoData(
                string id, string description, DateTime start_date, DateTime end_date
            )
            {
                this.id = id;
                this.description = description;
                this.start_date = start_date;
                this.end_date = end_date;
            }
        }

        public class LeilaoDataPlus
        {
            public LeilaoData leilao { get; set; }
            public bool notificar { get; set; }

            public LeilaoDataPlus(LeilaoData leilao, bool notificar)
            {
                this.leilao = leilao;
                this.notificar = notificar;
            }
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

            using var response = await httpClient.PostAsync($"{MSLeilaoAddress}/create", content);

            httpContext.Response.StatusCode = (int)response.StatusCode;
            if (response.StatusCode.ToString() == "Created")
            {
                var leilao = JsonSerializer.Deserialize<LeilaoData>(body);
                notificacao.AddLeilao(leilao!.id);
            }
            httpContext.Response.ContentType = response.Content.Headers.ContentType?.ToString() ?? "application/json";
            var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(respBody);
        }

        public async Task ListLeiloes(HttpContext httpContext)
        {
            using var activeLeiloesResponse = await httpClient.GetAsync($"{MSLeilaoAddress}/list");
            var userId = httpContext.Request.Query["userId"];

            string responseBody = await activeLeiloesResponse.Content.ReadAsStringAsync();
            var leiloes = JsonSerializer.Deserialize<List<LeilaoData>>(responseBody);
            List<LeilaoDataPlus> leiloesPlus = new List<LeilaoDataPlus>();

            for (int i = 0; i < leiloes?.Count; i++)
            {
                var leilaoPlus = new LeilaoDataPlus(leiloes[i], false);

                if (!notificacao.InterestLists.ContainsKey(leiloes[i].id))
                {
                    notificacao.InterestLists.Add(leiloes[i].id, new Notificacao.InterestList());
                }
                var clients = notificacao.InterestLists[leiloes[i].id].ClientIds;
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