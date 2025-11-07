using System.Numerics;
using System.Text.Json;
using RabbitMQ.Client;

namespace Routes
{
    class Notificacao
    {
        private static readonly HttpClient httpClient = new HttpClient();
        // private static string MSNotificacaoAddress = "http://localhost:8090";
        private Dictionary<string, InterestList> InterestLists = new Dictionary<string, InterestList>();

        class InterestList
        {
            public Dictionary<string, int> ClientIds = new Dictionary<string, int>();
        }

        class InterestRequest
        {
            public string ClientId { get; set; }
            public string LeilaoId { get; set; }
        }

        public async void ConnectCreateChannel()
        {
            ConnectionFactory factory = new ConnectionFactory();
            // "guest"/"guest" by default, limited to localhost connections
            // factory.UserName = user;
            // factory.Password = pass;
            // factory.VirtualHost = vhost;
            // factory.HostName = hostName;

            IConnection conn = await factory.CreateConnectionAsync();
        }

        public void SetupRoutes(WebApplication app)
        {
            app.MapGet("/event", SendNotification);
            app.MapPost("/register", RegisterInterest);
            app.MapPost("/cancel", CancelInterest);
        }

        public async Task SendNotification(HttpContext httpContext)
        {
            httpContext.Response.Headers.Append("Content-Type", "text/event-stream");

            while (true)
            {
                await httpContext.Response.WriteAsync("event: teu_pai\n");
                await httpContext.Response.WriteAsync("data: teu_pai\n\n");
                await httpContext.Response.Body.FlushAsync();

                await Task.Delay(Random.Shared.Next(1000, 5000));
            }
        }

        public async Task RegisterInterest(HttpContext httpContext)
        {
            httpContext.Request.EnableBuffering();

            string body;
            using (var reader = new StreamReader(httpContext.Request.Body, System.Text.Encoding.UTF8, leaveOpen: true))
            {
                body = await reader.ReadToEndAsync();
                httpContext.Request.Body.Position = 0;
            }

            InterestRequest? interest =
                JsonSerializer.Deserialize<InterestRequest>(body);

            if (interest == null)
            {
                await httpContext.Response.WriteAsync("RUIM");
            }
            else
            {
                if (InterestLists.ContainsKey(interest.LeilaoId))
                {
                    var interestList = InterestLists[interest.LeilaoId].ClientIds;
                    if (interestList.ContainsKey(interest.ClientId))
                    {
                        await httpContext.Response.WriteAsync("RUIM");
                    }
                    else
                    {
                        interestList.Add(interest.ClientId, 0);
                        await httpContext.Response.WriteAsync("BOM");
                    }
                }
            }
        }

        public async Task CancelInterest(HttpContext httpContext)
        {
            httpContext.Request.EnableBuffering();

            string body;
            using (var reader = new StreamReader(httpContext.Request.Body, System.Text.Encoding.UTF8, leaveOpen: true))
            {
                body = await reader.ReadToEndAsync();
                httpContext.Request.Body.Position = 0;
            }

            InterestRequest? interest =
                JsonSerializer.Deserialize<InterestRequest>(body);

            if (interest == null)
            {
                await httpContext.Response.WriteAsync("RUIM");
            }
            else
            {
                if (InterestLists.ContainsKey(interest.LeilaoId))
                {
                    var interestList = InterestLists[interest.LeilaoId].ClientIds;
                    if (interestList.ContainsKey(interest.ClientId))
                    {
                        interestList.Remove(interest.ClientId);
                        await httpContext.Response.WriteAsync("BOM");
                    }
                    else
                    {
                        await httpContext.Response.WriteAsync("RUIM");
                    }
                }
            }
        }
    }
}