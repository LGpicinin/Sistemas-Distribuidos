using Routes;
using GrpcGateway;
using GrpcLeilao;
using GrpcLance;
using GrpcPagamento;

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

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

// Enable CORS
app.UseCors("AllowAll");


Notificacao notificacaoRouter = new Notificacao();
Lance lanceRouter = new Lance();
Routes.Leilao leilaoRouter = new Routes.Leilao(notificacaoRouter);

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

