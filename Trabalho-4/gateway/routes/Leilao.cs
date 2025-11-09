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
        }

        public class LeilaoDataPlus
        {
            public LeilaoData leilao { get; set; }
            public bool notificar { get; set; }
        }

        public void SetupRoutes(WebApplication app)
        {
            app.MapPost("/leilao/create", NewLeilao);
            app.MapGet("/leilao/list", ListLeiloes);
        }

        public async Task NewLeilao(HttpContext httpContext)
        {
            httpContext.Request.EnableBuffering();

            string body;
            using (var reader = new StreamReader(httpContext.Request.Body, System.Text.Encoding.UTF8, leaveOpen: true))
            {
                body = await reader.ReadToEndAsync();
                httpContext.Request.Body.Position = 0;
            }

            var contentType = httpContext.Request.ContentType ?? "application/json";
            using var content = new StringContent(body, System.Text.Encoding.UTF8, contentType);

            using var response = await httpClient.PostAsync($"{MSLeilaoAddress}/create", content);

            httpContext.Response.StatusCode = (int)response.StatusCode;
            Console.WriteLine(response.StatusCode.ToString());
            if (response.StatusCode.ToString() == "Created")
            {
                var leilao = JsonSerializer.Deserialize<LeilaoData>(body);
                notificacao.AddLeilao(leilao.id);
            }
            httpContext.Response.ContentType = response.Content.Headers.ContentType?.ToString() ?? "application/json";
            var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(respBody);
        }

        public async Task ListLeiloes(HttpContext httpContext)
        {
            using var activeLeiloesResponse = await httpClient.GetAsync($"{MSLeilaoAddress}/list");

            httpContext.Response.StatusCode = (int)activeLeiloesResponse.StatusCode;

            var userId = httpContext.Request.Query["userId"];
            Console.WriteLine(userId);

            string responseBody = await activeLeiloesResponse.Content.ReadAsStringAsync();


            var leiloes = JsonSerializer.Deserialize<List<LeilaoData>>(responseBody);
            List<LeilaoDataPlus> leiloesPlus = new List<LeilaoDataPlus>();

            for (int i = 0; i < leiloes?.Count; i++)
            {
                var leilaoPlus = new LeilaoDataPlus();
                leilaoPlus.leilao = leiloes[i];
                leilaoPlus.notificar = false;
                if (!notificacao.InterestLists.ContainsKey(leiloes[i].id))
                {
                    Console.WriteLine("vai dar merda");
                    continue;
                }
                var clients = notificacao.InterestLists[leiloes[i].id].ClientIds;
                if (clients.ContainsKey(userId))
                {
                    leilaoPlus.notificar = true;
                }
                leiloesPlus.Add(leilaoPlus);
            }

            var leiloesSerialized = JsonSerializer.Serialize(leiloesPlus);

            httpContext.Response.ContentType = activeLeiloesResponse.Content.Headers.ContentType?.ToString() ?? "application/json";
            await httpContext.Response.WriteAsync(leiloesSerialized);
        }
    }
}