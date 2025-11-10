import type { PageServerLoad } from '../$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const userId = cookies.get('userId');

	return {
		userId: userId!
	};
};
