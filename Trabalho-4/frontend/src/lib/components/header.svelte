<script lang="ts">
	import { fade } from 'svelte/transition';
	import Hamburguer from '$lib/assets/hamburguer.svg';
	import Chevron from '$lib/assets/chevron.svg';

	export let open: boolean;
	export let username: string;
</script>

<header>
	<button
		type="button"
		aria-label="button"
		onclick={() => {
			open = !open;
		}}
	>
		<img transition:fade src={open ? Chevron : Hamburguer} alt="Menu Icon" />
	</button>
	<h1>UTFPR Leilões</h1>
	<span>{username}</span>
</header>

<style>
	header {
		display: flex;
		flex-direction: row;
		padding: 1rem;
		justify-content: space-between;
		align-items: center;
		background-color: var(--utfpr-main-color);
	}

	button {
		background-color: transparent;
		border: none;
		animation: 0.5s fadeIn;
		cursor: pointer;
	}

	img {
		width: 2rem;
	}
</style>
