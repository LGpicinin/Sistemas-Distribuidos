import type { Leilao } from '$lib/helpers/models/leilao.d.ts';
import { GATEWAY_ADDRESS } from '$env/static/private';

export const load = async ({ fetch }) => {
	const response = await fetch(`${GATEWAY_ADDRESS}/leilao/list`);

	const data: Leilao[] | null = await response.json();
	return {
		leiloes: data ?? ([] as Leilao[])
	};
};
