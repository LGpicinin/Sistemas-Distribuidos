using Routes;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring OpenAPI at https://aka.ms/aspnet/openapi
builder.Services.AddOpenApi();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

Notificacao notificacaoRouter = new Notificacao();
Lance lanceRouter = new Lance();
Leilao leilaoRouter = new Leilao(notificacaoRouter);

lanceRouter.SetupRoutes(app);
leilaoRouter.SetupRoutes(app);
notificacaoRouter.SetupRoutes(app);

await notificacaoRouter.ConnectCreateChannel();
await notificacaoRouter.ConsumeEvents();

// app.UseHttpsRedirection();

app.Run();

