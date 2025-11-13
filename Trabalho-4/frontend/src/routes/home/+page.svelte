<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import { PUBLIC_GATEWAY_ADDRESS } from '$env/static/public';
	import { onDestroy, onMount } from 'svelte';
	import Button from '$lib/components/button.svelte';
	import { type Notification as NotificationType } from '$lib/helpers/models/notification.js';
	import { loadNotifications, saveNotifications } from '$lib/helpers/utils/notifications.js';
	import Notification from '$lib/components/notification.svelte';
	import { flip } from 'svelte/animate';

	const { data } = $props();

	let messages = $state([] as NotificationType[]);
	let eventSource: EventSource;

	onMount(() => {
		eventSource = new EventSource(`${PUBLIC_GATEWAY_ADDRESS}/event?userId=${data.userId}`);
		messages = loadNotifications();
		eventSource.addEventListener(data.userId, (event) => {
			messages = [JSON.parse(event.data), ...messages];
			saveNotifications(messages);
		});
	});

	onDestroy(() => {
		eventSource?.close();
	});
</script>

<Card>
	<h2>Histórico de Notificações</h2>
	{#each messages as message}
		<Card>
			<Notification {message} />
		</Card>
	{:else}
		<p>Sem notificações para exibir</p>
	{/each}
</Card>
