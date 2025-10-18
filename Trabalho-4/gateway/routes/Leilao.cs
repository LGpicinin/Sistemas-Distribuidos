namespace Routes
{
    class Leilao
    {
        private static readonly HttpClient httpClient = new HttpClient();
        private static string MSLeilaoAddress = "http://localhost:8090";

        public void SetupRoutes(WebApplication app)
        {
            app.MapPost("/leilao/create", NewLeilao);
            app.MapGet("/leilao/list", ListLeiloes);
        }

        public async Task NewLeilao(HttpContext httpContext)
        {

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