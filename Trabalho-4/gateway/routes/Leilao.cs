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
            public string ID { get; set; }
            public string Description { get; set; }
            public DateTime StartDate { get; set; }
            public DateTime EndDate { get; set; }
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
            if (httpContext.Response.StatusCode == 201)
            {
                using (var reader = new StreamReader(httpContext.Response.Body, System.Text.Encoding.UTF8, leaveOpen: true))
                {
                    body = await reader.ReadToEndAsync();
                }

                var leilao = JsonSerializer.Deserialize<LeilaoData>(body);
                notificacao.AddLeilao(leilao.ID);
            }
            httpContext.Response.ContentType = response.Content.Headers.ContentType?.ToString() ?? "application/json";
            var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(respBody);
        }

        public async Task ListLeiloes(HttpContext httpContext)
        {
            using var activeLeiloesResponse = await httpClient.GetAsync($"{MSLeilaoAddress}/list");

            httpContext.Response.StatusCode = (int)activeLeiloesResponse.StatusCode;
            httpContext.Response.ContentType = activeLeiloesResponse.Content.Headers.ContentType?.ToString() ?? "application/json";
            var respBody = await activeLeiloesResponse.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(respBody);
        }
    }
}