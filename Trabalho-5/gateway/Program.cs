using Routes;
using GrpcGateway;
using GrpcLeilao;
using GrpcLance;
using GrpcPagamento;
using System.Net;
using Microsoft.AspNetCore.Server.Kestrel.Core;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddGrpc();

// Add services to the container.
// Learn more about configuring OpenAPI at https://aka.ms/aspnet/openapi
builder.Services.AddOpenApi();

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAll", policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader();
    });
});

builder.WebHost.ConfigureKestrel((context, serverOptions) =>
{
    var kestrelSection = context.Configuration.GetSection("Kestrel");

    serverOptions.Listen(IPAddress.Loopback, 5060, listenOptions =>
    {
        listenOptions.Protocols = HttpProtocols.Http2;
    });
    serverOptions.Listen(IPAddress.Loopback, 5059, listenOptions =>
    {
        listenOptions.Protocols = HttpProtocols.Http1;
    });
});

Notificacao notificacaoRouter = new Notificacao();
Lance lanceRouter = new Lance();
Routes.Leilao leilaoRouter = new Routes.Leilao(notificacaoRouter);

builder.Services.AddSingleton(notificacaoRouter);

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

// Enable CORS
app.UseCors("AllowAll");

lanceRouter.SetupRoutes(app);
leilaoRouter.SetupRoutes(app);
notificacaoRouter.SetupRoutes(app);

// app.MapGrpcService<LanceService>();
// app.MapGrpcService<LeilaoService>();
app.MapGrpcService<GrpcGateway.Services.Gateway>();

await lanceRouter.ConnectCreateChannel();
await leilaoRouter.ConnectCreateChannel();
// await notificacaoRouter.ConsumeLanceEvents();
// await notificacaoRouter.ConsumeStatusEvents();
// await notificacaoRouter.ConsumeLinkEvents();

// app.UseHttpsRedirection();

app.Run();

