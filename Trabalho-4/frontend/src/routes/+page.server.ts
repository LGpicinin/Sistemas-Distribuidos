import { redirect } from '@sveltejs/kit';

export const actions = {
	setUserId: async ({ cookies, request }) => {
		const data = await request.formData();
		cookies.set('userId', data.get('userId')?.toString() ?? '', { path: '/' });
		return redirect(303, '/home');
	}
};
