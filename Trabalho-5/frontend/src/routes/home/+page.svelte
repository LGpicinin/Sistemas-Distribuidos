<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import { PUBLIC_GATEWAY_ADDRESS } from '$env/static/public';
	import { onDestroy, onMount } from 'svelte';
	import Button from '$lib/components/button.svelte';
	import { type Notification as NotificationType } from '$lib/helpers/models/notification.js';
	import { loadNotifications, saveNotifications } from '$lib/helpers/utils/notifications.js';
	import Notification from '$lib/components/notification.svelte';
	import { flip } from 'svelte/animate';
	import { writable } from 'svelte/store';

	const { data } = $props();

	// let messages = $state(writable([] as NotificationType[]));
	// let eventSource: EventSource;

	// onMount(() => {
	// 	eventSource = new EventSource(`${PUBLIC_GATEWAY_ADDRESS}/event?userId=${data.userId}`);
	// 	messages = loadNotifications(messages);
	// 	eventSource.addEventListener(data.userId, (event) => {
	// 		messages.update((not) => [JSON.parse(event.data), ...not])
	// 		// messages = [JSON.parse(event.data), ...messages];
	// 		messages = saveNotifications(messages);
	// 	});
	// });

	// onDestroy(() => {
	// 	eventSource?.close();
	// });
</script>

<Card>
	<h2>Bem vindo ao UTFPR Leilões</h2>
	<!-- {#each $messages as message}
		<Card>
			<Notification {message} />
		</Card>
	{:else}
		<p>Sem notificações para exibir</p>
	{/each} -->
</Card>
