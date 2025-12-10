using System.Numerics;
using System.Text.Json;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;


namespace Routes
{
    class Notificacao
    {
        private string EXCHANGE_NAME = "LEILAO";
        private string QUEUE_LANCE_VALIDADO = "lance_validado";
        private string QUEUE_LANCE_INVALIDADO = "lance_invalidado";
        private string QUEUE_LEILAO_VENCEDOR = "leilao_vencedor";
        private string QUEUE_LINK_PAGAMENTO = "link_pagamento";

        private string QUEUE_STATUS_PAGAMENTO = "status_pagamento";

        private static readonly HttpClient httpClient = new HttpClient();
        // private static string MSNotificacaoAddress = "http://localhost:8090";
        public Dictionary<string, InterestList> InterestLists = new Dictionary<string, InterestList>();
        private Dictionary<string, HttpContext> UserList = new Dictionary<string, HttpContext>();
        // private List<HttpContext> UserList = new List<HttpContext>();
        private ConnectionFactory factory = new ConnectionFactory();

        private IConnection conn = null;

        private IChannel channel = null;

        public class LanceData
        {
            public string leilao_id { get; set; }
            public string user_id { get; set; }
            public float value { get; set; }
        }

        public class StatusData
        {
            public string clientId { get; set; }
            public string paymentId { get; set; }
            public float value { get; set; }

            public bool status { get; set; }
        }

        public class StatusDataType
        {
            public string type { get; set; }
            public StatusData statusData { get; set; }
        }

        public class LinkData
        {
            public string clientId { get; set; }
            public string link { get; set; }
        }

        public class LinkDataType
        {
            public string type { get; set; }
            public LinkData linkData { get; set; }
        }

        public class LanceDataType
        {
            public string type { get; set; }
            public LanceData lance { get; set; }
        }

        public class InterestList
        {
            public Dictionary<string, int> ClientIds = new Dictionary<string, int>();
        }

        class InterestRequest
        {
            public string UserId { get; set; }
            public string LeilaoId { get; set; }
        }

        // public async Task ConnectCreateChannel()
        // {
        //     factory.Uri = new Uri("amqp://guest:guest@localhost:5672/");
        //     conn = await factory.CreateConnectionAsync();

        //     channel = await conn.CreateChannelAsync();

        //     await channel.QueueDeclareAsync(QUEUE_LANCE_VALIDADO, true, false, false, null);
        //     await channel.QueueBindAsync(QUEUE_LANCE_VALIDADO, EXCHANGE_NAME, QUEUE_LANCE_VALIDADO, null);

        //     await channel.QueueDeclareAsync(QUEUE_LEILAO_VENCEDOR, true, false, false, null);
        //     await channel.QueueBindAsync(QUEUE_LEILAO_VENCEDOR, EXCHANGE_NAME, QUEUE_LEILAO_VENCEDOR, null);
        // }

        public void SetupRoutes(WebApplication app)
        {
            app.MapGet("/event", SendNotification);
            app.MapPost("/register", RegisterInterest);
            app.MapPost("/cancel", CancelInterest);
        }

        public async Task ConsumeLanceEvents(LanceDataType lancePlus)
        {
            var lanceSerialized = JsonSerializer.Serialize<LanceDataType>(lancePlus);
            if (InterestLists.ContainsKey(lancePlus.lance.leilao_id))
            {
                var interestList = InterestLists[lancePlus.lance.leilao_id].ClientIds;
                foreach (KeyValuePair<string, int> entry in interestList)
                {
                    HttpContext context = UserList[entry.Key];
                    Console.WriteLine(entry.Key);
                    await context.Response.WriteAsync($"event: {entry.Key}\n");
                    await context.Response.WriteAsync($"data: {lanceSerialized}\n\n");
                    await context.Response.Body.FlushAsync();
                    Console.WriteLine("context");
                }
            }
        }

        public async Task ConsumeStatusEvents(StatusDataType statusPlus)
        {
            var statusSerialized = JsonSerializer.Serialize<StatusDataType>(statusPlus);

            HttpContext context = UserList[statusPlus.statusData.clientId];
            Console.WriteLine(statusPlus.statusData.clientId);
            await context.Response.WriteAsync($"event: {statusPlus.statusData.clientId}\n");
            await context.Response.WriteAsync($"data: {statusSerialized}\n\n");
            await context.Response.Body.FlushAsync();
            Console.WriteLine("context");
        }

        public async Task ConsumeLinkEvents(LinkDataType linkPlus)
        {
            var linkSerialized = JsonSerializer.Serialize<LinkDataType>(linkPlus);

            HttpContext context = UserList[linkPlus.linkData.clientId];
            Console.WriteLine(linkPlus.linkData.clientId);
            await context.Response.WriteAsync($"event: {linkPlus.linkData.clientId}\n");
            await context.Response.WriteAsync($"data: {linkSerialized}\n\n");
            await context.Response.Body.FlushAsync();
            Console.WriteLine("context");
        }

        // recebe requisição sse do usuário
        // salva sessão http em lista para ser usado depois para envio de notificações
        public async Task SendNotification(HttpContext httpContext)
        {
            httpContext.Request.EnableBuffering();

            var userId = httpContext.Request.Query["userId"];

            httpContext.Response.Headers.Append("Content-Type", "text/event-stream");
            if (UserList.ContainsKey(userId))
            {
                UserList[userId] = httpContext;
            }
            else
            {
                UserList.Add(userId, httpContext);
            }
            while (UserList[userId] == httpContext)
            {
                continue;
            }


        }

        public void AddLeilao(string id)
        {
            if (!InterestLists.ContainsKey(id))
            {
                InterestLists.Add(id, new InterestList());
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
                    if (interestList.ContainsKey(interest.UserId))
                    {
                        await httpContext.Response.WriteAsync("RUIM");
                    }
                    else
                    {
                        interestList.Add(interest.UserId, 0);
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
                    if (interestList.ContainsKey(interest.UserId))
                    {
                        interestList.Remove(interest.UserId);
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