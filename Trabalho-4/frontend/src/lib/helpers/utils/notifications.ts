import type { Notification } from '../models/notification';
import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store';


export const loadNotifications = (messages: Writable<Notification[]>) => {
	const notifications = sessionStorage.getItem('notifications');

	messages = writable(notifications === null ? [] as Notification[] : JSON.parse(notifications) as Notification[])

	return messages;
};

export const saveNotifications = (messages: Writable<Notification[]>) => {

	messages.subscribe((value) => {
		sessionStorage.setItem('notifications', JSON.stringify(value))
	})

	return messages
};
