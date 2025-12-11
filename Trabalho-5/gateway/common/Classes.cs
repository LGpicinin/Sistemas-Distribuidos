namespace Classes;

public class LanceData
{
    public string leilao_id { get; set; }
    public string user_id { get; set; }
    public float value { get; set; }
}

public class LanceDataType
{
    public string type { get; set; }
    public LanceData lance { get; set; }
}

public class LeilaoData
{
    public string id { get; set; }
    public string description { get; set; }
    public DateTime start_date { get; set; }
    public DateTime end_date { get; set; }
    public LeilaoData(
        string id, string description, DateTime start_date, DateTime end_date
    )
    {
        this.id = id;
        this.description = description;
        this.start_date = start_date;
        this.end_date = end_date;
    }
}

public class LeilaoDataPlus
{
    public LeilaoData leilao { get; set; }
    public bool notificar { get; set; }
    public LeilaoDataPlus(LeilaoData leilao, bool notificar)
    {
        this.leilao = leilao;
        this.notificar = notificar;
    }
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