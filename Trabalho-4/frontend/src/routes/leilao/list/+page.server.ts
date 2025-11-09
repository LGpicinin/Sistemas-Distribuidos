import type { Leilao, LeilaoPlus } from '$lib/helpers/models/leilao.d.ts';
import { GATEWAY_ADDRESS } from '$env/static/private';

export const load = async ({ fetch, cookies }) => {
	const userId = cookies.get('userId');
	const response = await fetch(`${GATEWAY_ADDRESS}/leilao/list?userId=${userId}`);

	// const data = JSON.parse(response)
	const data: LeilaoPlus[] | null = await response.json();
	console.log(data)
	return {
		leiloes: data ?? ([] as LeilaoPlus[])
	};
};
