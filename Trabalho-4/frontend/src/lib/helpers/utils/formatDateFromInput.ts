export const formatDateFromInput = (dateString: string) => {
    const pattern = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2})$/
    const match = dateString.match(pattern);
    if (match) {
        let year = match[1];
        let month = match[2];
        let day = match[3];
        let hour = match[4];
        let minute = match[5];

        const date = new Date(0);
        date.setFullYear(Number(year));
        date.setMonth(Number(month));
        date.setDate(Number(day));
        date.setHours(Number(hour), Number(minute));

        return date.toISOString();
    }

    return new Date().toISOString();
}