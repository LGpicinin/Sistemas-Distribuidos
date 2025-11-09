export interface Leilao {
	id: string;
	description: string;
	start_date: string;
	end_date: string;
}

export interface LeilaoPlus {
	leilao: Leilao;
	notificar: boolean;
}
