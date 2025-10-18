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

Lance lanceRouter = new Lance();
Leilao leilaoRouter = new Leilao();

lanceRouter.SetupRoutes(app);
leilaoRouter.SetupRoutes(app);

// app.UseHttpsRedirection();

app.MapGet("/event", async (HttpContext httpContext) =>
{
    httpContext.Response.Headers.Append("Content-Type", "text/event-stream");

    while (true)
    {
        await httpContext.Response.WriteAsync("event: teu_pai\n");
        await httpContext.Response.WriteAsync("data: teu_pai\n\n");
        await httpContext.Response.Body.FlushAsync();

        await Task.Delay(Random.Shared.Next(1000, 5000));
    }
})
.WithName("GetWeatherForecast");

app.Run();

