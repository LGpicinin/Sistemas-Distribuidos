using Routes;

var builder = WebApplication.CreateBuilder(args);

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
Leilao leilaoRouter = new Leilao(notificacaoRouter);

lanceRouter.SetupRoutes(app);
leilaoRouter.SetupRoutes(app);
notificacaoRouter.SetupRoutes(app);

await notificacaoRouter.ConnectCreateChannel();
await notificacaoRouter.ConsumeLanceEvents();
await notificacaoRouter.ConsumeStatusEvents();
await notificacaoRouter.ConsumeLinkEvents();

// app.UseHttpsRedirection();

app.Run();

