namespace Routes
{
    class Lance
    {
        private static readonly HttpClient httpClient = new HttpClient();
        private static string MSLanceAddress = "http://localhost:8080";

        public void SetupRoutes(WebApplication app)
        {
            app.MapPost("/lance/new", NewLance);
        }

        public async Task NewLance(HttpContext httpContext)
        {
            var content = await Utils.HTTPHelper.getRequestContent(httpContext);

            using var response = await httpClient.PostAsync($"{MSLanceAddress}/new", content);

            httpContext.Response.StatusCode = (int)response.StatusCode;
            httpContext.Response.ContentType = response.Content.Headers.ContentType?.ToString() ?? "application/json";
            var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(respBody);
        }
    }
}