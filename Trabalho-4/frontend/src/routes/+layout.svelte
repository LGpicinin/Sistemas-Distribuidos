<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/header.svelte';
	import Sidebar from '$lib/components/sidebar.svelte';

	let { children } = $props();

	let sidebarOpen: boolean = $state(false);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="main-wrapper">
	<Header bind:open={sidebarOpen} username="Doglas" />
	<Sidebar bind:isOpen={sidebarOpen} />
	{@render children?.()}
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		min-height: 100vh;
		min-width: 100vw;
	}

	:global(*) {
		padding: 0;
		margin: 0;
		--utfpr-main-color: #ffcc00;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}

	.main-wrapper {
		display: flex;
		justify-content: space-between;
		height: 100%;
		width: 100%;
	}
</style>
