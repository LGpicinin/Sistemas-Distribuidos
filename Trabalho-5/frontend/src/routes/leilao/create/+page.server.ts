import type { Leilao } from '$lib/helpers/models/leilao.d.ts';
import { GATEWAY_ADDRESS } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';
import { formatDateFromInput } from '$lib/helpers/utils/formatDateFromInput';

export const actions = {
	createLeilao: async ({ request, fetch }) => {
		const formData = await request.formData();
		const leilaoData: Leilao = {
			id: formData.get('id') as string,
			description: formData.get('description') as string,
			start_date: formatDateFromInput(formData.get('start_date') as string),
			end_date: formatDateFromInput(formData.get('end_date') as string)
		};

		if (leilaoData.start_date >= leilaoData.end_date) {
			return fail(403, {
				error: `Favor criar leilão com data de término superior a data de início.`
			});
		}

		const response = await fetch(`${GATEWAY_ADDRESS}/leilao/create`, {
			method: 'POST',
			body: JSON.stringify(leilaoData),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (response.status === 201) return redirect(303, '/leilao/list');

		return fail(403, {
			error: `Erro ao criar leilão: ${response.statusText}`
		});
	}
};
