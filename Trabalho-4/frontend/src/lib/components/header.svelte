<script lang="ts">
	import { fade } from 'svelte/transition';
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
		style="--rotation: {!open ? '90deg' : '0deg'}"
		transition:fade
	>
		<span>III</span>
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
		rotate: calc(var(--rotation));
		animation: 0.5s fadeIn;
		cursor: pointer;
		font-size: 2rem;
	}
</style>
