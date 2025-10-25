<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/header.svelte';
	import Sidebar from '$lib/components/sidebar.svelte';

	let { children, data } = $props();

	let sidebarOpen: boolean = $state(false);
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
