import type { Notification } from '../models/notification';

export const loadNotifications = () => {
	const notifications = sessionStorage.getItem('notifications');

	if (notifications === null) return [] as Notification[];

	return JSON.parse(notifications) as Notification[];
};

export const saveNotifications = (notifications: Notification[]) => {
	const notificationsString = JSON.stringify(notifications);

	sessionStorage.setItem('notifications', notificationsString);
};
