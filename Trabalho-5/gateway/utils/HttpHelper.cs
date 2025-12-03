namespace Utils
{
    class HTTPHelper
    {
        public async static Task<string> getRequestBody(HttpContext httpContext)
        {
            httpContext.Request.EnableBuffering();

            string body;
            using (var reader = new StreamReader(httpContext.Request.Body, System.Text.Encoding.UTF8, leaveOpen: true))
            {
                body = await reader.ReadToEndAsync();
                httpContext.Request.Body.Position = 0;
            }

            return body;
        }

        public async static Task<StringContent> getRequestContentFromBody(string body, HttpContext httpContext)
        {
            var contentType = httpContext.Request.ContentType ?? "application/json";
            var content = new StringContent(body, System.Text.Encoding.UTF8, contentType);

            return content;
        }

        public async static Task<StringContent> getRequestContent(HttpContext httpContext)
        {
            string body = await getRequestBody(httpContext);
            return await getRequestContentFromBody(body, httpContext);
        }
    }
}