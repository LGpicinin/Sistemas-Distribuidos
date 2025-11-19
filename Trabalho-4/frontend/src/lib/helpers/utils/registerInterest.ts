import type { Interest } from '../models/interest';

export const registerInterest = async (userId: string, leilao_id: string) => {
	const interest: Interest = {
		UserId: userId,
		LeilaoId: leilao_id
	};

	await fetch(`http://localhost:5059/register`, {
		method: 'POST',
		body: JSON.stringify(interest),
		headers: {
			'Content-Type': 'application/json'
		}
	});
};
