<script lang="ts">
	import { fade } from 'svelte/transition';
	import Hamburguer from '$lib/assets/hamburguer.svg';
	import Chevron from '$lib/assets/chevron.svg';
	import Button from './button.svelte';

	export let open: boolean;
	export let username: string;
</script>

<header>
	<div class="left">
		<button
			type="button"
			aria-label="button"
			onclick={() => {
				open = !open;
			}}
		>
			<img transition:fade src={open ? Chevron : Hamburguer} alt="Menu Icon" />
		</button>
	</div>
	<a href="/home"><h1>UTFPR Leilões</h1></a>
	<div class="right">
		<p>{username}</p>
		<form action="/?/deleteUserId" method="post">
			<Button type="submit" text="Sair" />
		</form>
	</div>
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

	.left {
		flex: 1;
	}

	.right {
		flex: 1;
		margin-left: auto;
		display: flex;
		justify-content: flex-end;
		align-items: center;
		gap: 1rem;
	}

	button {
		background-color: transparent;
		border: none;
		animation: 0.5s fadeIn;
		cursor: pointer;
		flex: 1;
	}

	img {
		width: 2rem;
	}

	a {
		text-decoration: none;
		color: black;
	}
</style>
