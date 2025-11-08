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
        private string QUEUE_LEILAO_INICIADO = "leilao_iniciado";
        private string QUEUE_LEILAO_FINALIZADO = "leilao_finalizado";

        private static readonly HttpClient httpClient = new HttpClient();
        // private static string MSNotificacaoAddress = "http://localhost:8090";
        private Dictionary<string, InterestList> InterestLists = new Dictionary<string, InterestList>();
        private Dictionary<string, HttpContext> UserList = new Dictionary<string, HttpContext>();
        private ConnectionFactory factory = new ConnectionFactory();

        private IConnection conn = null;

        private IChannel channel = null;

        public class LanceData
        {
            public string LeilaoID { get; set; }
            public string UserID { get; set; }
            public float Value { get; set; }
        }

        public class LanceDataType
        {
            public string type { get; set; }
            public LanceData lance{ get; set; }
        }

        class InterestList
        {
            public Dictionary<string, int> ClientIds = new Dictionary<string, int>();
        }

        class InterestRequest
        {
            public string ClientId { get; set; }
            public string LeilaoId { get; set; }
        }

        public async Task ConnectCreateChannel()
        {
            factory.Uri = new Uri("amqp://guest:guest@localhost:5672/");
            conn = await factory.CreateConnectionAsync();

            channel = await conn.CreateChannelAsync();

            await channel.QueueDeclareAsync(QUEUE_LANCE_VALIDADO, true, false, false, null);
            await channel.QueueBindAsync(QUEUE_LANCE_VALIDADO, EXCHANGE_NAME, QUEUE_LANCE_VALIDADO, null);

            await channel.QueueDeclareAsync(QUEUE_LEILAO_VENCEDOR, true, false, false, null);
            await channel.QueueBindAsync(QUEUE_LEILAO_VENCEDOR, EXCHANGE_NAME, QUEUE_LEILAO_VENCEDOR, null);
        }

        public void SetupRoutes(WebApplication app)
        {
            app.MapGet("/event", SendNotification);
            app.MapPost("/register", RegisterInterest);
            app.MapPost("/cancel", CancelInterest);
        }

        public async Task ConsumeEvents()
        {
            var consumer = new AsyncEventingBasicConsumer(channel);
            consumer.ReceivedAsync += async (sender, eventArgs) =>
            {
                var routingKey = eventArgs.RoutingKey;
                byte[] body = eventArgs.Body.ToArray();
                string message = Encoding.UTF8.GetString(body);
                var lance = JsonSerializer.Deserialize<LanceData>(message);

                // Acknowledge the message
                await ((AsyncEventingBasicConsumer)sender)
                    .Channel.BasicAckAsync(eventArgs.DeliveryTag, multiple: false);

                var lancePlus = new LanceDataType();
                lancePlus.lance = lance;
                lancePlus.type = routingKey;

                var lanceSerialized = JsonSerializer.Serialize<LanceDataType>(lancePlus);
                if (InterestLists.ContainsKey(lance.LeilaoID))
                {
                    var interestList = InterestLists[lance.LeilaoID].ClientIds;
                    foreach(KeyValuePair<string, int> entry in interestList)
                    {
                        HttpContext httpContext = UserList[entry.Key];
                        await httpContext.Response.WriteAsync(lanceSerialized);
                        await httpContext.Response.Body.FlushAsync();
                    }
                }

            };
            await channel.BasicConsumeAsync("orders", autoAck: false, consumer);
        }

        public async Task SendNotification(HttpContext httpContext)
        {
            httpContext.Request.EnableBuffering();

            string body;
            using (var reader = new StreamReader(httpContext.Request.Body, System.Text.Encoding.UTF8, leaveOpen: true))
            {
                body = await reader.ReadToEndAsync();
                httpContext.Request.Body.Position = 0;
            }

            UserList.Add(body, httpContext);
            // httpContext.Response.Headers.Append("Content-Type", "text/event-stream");

            while (true)
            {
                
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