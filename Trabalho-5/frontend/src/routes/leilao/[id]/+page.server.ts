import type { PageServerLoad } from './$types';
import type { LanceBody } from '$lib/helpers/models/lance';
import { GATEWAY_ADDRESS } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';
import { registerInterest } from '$lib/helpers/utils/registerInterest';

export const load: PageServerLoad = async ({ cookies, params }) => {
	const userId = cookies.get('userId');
	const leilaoId = params.id;

	return {
		userId,
		leilaoId
	};
};

export const actions = {
	createLance: async ({ request, fetch, params, cookies }) => {
		const formData = await request.formData();
		const leilao_id = params.id;
		const user_id = cookies.get('userId');
		await registerInterest(user_id!, leilao_id);

		const lanceData: LanceBody = {
			user_id: user_id!,
			leilao_id,
			value: Number(formData.get('value') as string)
		};

		const lanceRequest = await fetch(`${GATEWAY_ADDRESS}/lance/new`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(lanceData)
		});

		if (lanceRequest.status === 201) {
			return redirect(303, '/home');
		}

		return fail(403, {
			error: `Erro ao criar lance: ${await lanceRequest.text()}`
		});
	}
};
