export const formatDateFromInput = (dateString: string) => {
	const pattern = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})$/;
	const match = dateString.match(pattern);
	if (match) {
		const year = match[1];
		const month = match[2];
		const day = match[3];
		const hour = match[4];
		const minute = match[5];

		const date = new Date(0);
		date.setUTCFullYear(Number(year));
		date.setUTCMonth(Number(month) - 1);
		date.setUTCDate(Number(day));
		date.setUTCHours(Number(hour), Number(minute));

		return date.toISOString();
	}

	return new Date().toISOString();
};
