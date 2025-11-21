<script lang="ts">
	import { PUBLIC_GATEWAY_ADDRESS } from '$env/static/public';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/header.svelte';
	import Sidebar from '$lib/components/sidebar.svelte';
	import NotificationBar from '$lib/components/notificationBar.svelte';
	import { loadNotifications, saveNotifications } from '$lib/helpers/utils/notifications.js';
	import { onMount, onDestroy } from 'svelte';
	import { type Notification } from '$lib/helpers/models/notification.js';
	import { writable } from 'svelte/store';

	let { children, data } = $props();
	let sidebarOpen: boolean = $state(false);
	let eventSource: EventSource;
	let messages = $state(writable([] as Notification[]));

	onMount(() => {
		eventSource = new EventSource(`${PUBLIC_GATEWAY_ADDRESS}/event?userId=${data.userId}`);
		messages = loadNotifications(messages);
		eventSource.addEventListener(data.userId, (event) => {
			messages.update((not) => [JSON.parse(event.data), ...not]);
			$inspect($messages).with(console.log);
			messages = saveNotifications(messages);
		});
	});

	onDestroy(() => {
		eventSource?.close();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="main-wrapper">
	<Header bind:open={sidebarOpen} username={data.userId} />
	<main>
		<Sidebar bind:isOpen={sidebarOpen} />
		<section class="content">
			{@render children?.()}
		</section>

		<NotificationBar bind:notifications={$messages} />
	</main>
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
	}

	:global(*) {
		padding: 0;
		margin: 0;
		--utfpr-main-color: #ffcc00;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}

	.main-wrapper {
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		height: 100%;
		width: 100%;
	}

	main {
		display: flex;
		height: 100%;
	}

	.content {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 100%;
	}
</style>
