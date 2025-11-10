export interface Lance {
	leilao_id: string;
	user_id: string;
	value: string;
}

export interface LanceBody extends Lance {
	value: number;
}
