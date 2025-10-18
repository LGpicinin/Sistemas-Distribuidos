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
            httpContext.Request.EnableBuffering();

            string body;
            using (var reader = new StreamReader(httpContext.Request.Body, System.Text.Encoding.UTF8, leaveOpen: true))
            {
                body = await reader.ReadToEndAsync();
                httpContext.Request.Body.Position = 0;
            }

            var contentType = httpContext.Request.ContentType ?? "application/json";
            using var content = new StringContent(body, System.Text.Encoding.UTF8, contentType);

            using var response = await httpClient.PostAsync($"{MSLanceAddress}/new", content);

            httpContext.Response.StatusCode = (int)response.StatusCode;
            httpContext.Response.ContentType = response.Content.Headers.ContentType?.ToString() ?? "application/json";
            var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(respBody);
        }
    }
}