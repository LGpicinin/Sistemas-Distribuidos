import type { Leilao } from '$lib/helpers/models/leilao.js';
import { GATEWAY_ADDRESS } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';
import { formatDateFromInput } from '$lib/helpers/utils/formatDateFromInput.js';

export const actions = {
	createLeilao: async ({ request, fetch }) => {
		const formData = await request.formData();
		const leilaoData: Leilao = {
			id: formData.get('id') as string,
			description: formData.get('description') as string,
			start_date: formatDateFromInput(formData.get('start_date') as string),
			end_date: formatDateFromInput(formData.get('end_date') as string)
		};

		try {
			const response = await fetch(`${GATEWAY_ADDRESS}/leilao/create`, {
				method: 'POST',
				body: JSON.stringify(leilaoData),
				headers: {
					'Content-Type': 'application/json'
				}
			});

			if (response.status === 201) return redirect(303, 'leilao/list');
			else throw new Error(`Erro ao criar leilão: ${response.statusText}`);
		} catch (error: unknown) {
			if (error instanceof Error)
				return fail(403, {
					error: error.message
				});
		}
	}
};
