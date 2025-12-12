using Grpc.Net.Client;
using GrpcLance;
using System.Text.Json;
using Classes;
using System;
using System.Net.Http;
using System.Threading.Tasks;

namespace Routes
{
    class Lance
    {
        private static readonly HttpClient httpClient = new HttpClient();
        private static string MSLanceAddress = "http://localhost:8080";

        private LanceService.LanceServiceClient ms_lance;

        public void SetupRoutes(WebApplication app)
        {
            app.MapPost("/lance/new", NewLance);
        }


        public async Task ConnectCreateChannel()
        {
            AppContext.SetSwitch(
                "System.Net.Http.SocketsHttpHandler.Http2UnencryptedSupport", true);
            // var options = new GrpcChannelOptions();
            // options.Credentials = Grpc.Core.ChannelCredentials.Insecure;
            var channel = GrpcChannel.ForAddress("http://localhost:8080");
            // var channel = GrpcChannel.ForAddress("localhost:8090",  new GrpcChannelOptions
            // {
            //     Credentials = ChannelCredentials.Insecure
            // });
            ms_lance = new LanceService.LanceServiceClient(channel);
        }

        public async Task NewLance(HttpContext httpContext)
        {
            var content = await Utils.HTTPHelper.getRequestContent(httpContext);
            string body = await Utils.HTTPHelper.getRequestBody(httpContext);

            var lanceData = JsonSerializer.Deserialize<LanceData>(body);

            var pbLance = new LLance
            {
                LeilaoID = lanceData.leilao_id,
                UserID = lanceData.user_id,
                Value = lanceData.value
            };

            var response = ms_lance.Create(pbLance);

            var num_status = 0;
            if (response.Status_ == "Created")
            {
                num_status = 201;
            } else
            {
                num_status = 400;
            }

            httpContext.Response.StatusCode = num_status;

            httpContext.Response.ContentType = "application/json";
            // var respBody = await response.Content.ReadAsStringAsync();
            await httpContext.Response.WriteAsync(await content.ReadAsStringAsync());
        }
    }
}