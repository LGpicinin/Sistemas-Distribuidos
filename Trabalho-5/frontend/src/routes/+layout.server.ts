import { redirect } from '@sveltejs/kit';

export const load = ({ cookies, url }) => {
	const userId = cookies.get('userId');

	if (url.pathname !== '/' && !userId) return redirect(303, '/');

	return {
		userId: userId ?? '',
		isHome: url.pathname === '/home'
	};
};
