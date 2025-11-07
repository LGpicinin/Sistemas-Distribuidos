<script lang="ts">
	import { fade } from 'svelte/transition';

	export let isOpen;
</script>

{#if isOpen}
	<section transition:fade>
		<ul>
			<span>Leilão</span>
			<li><a href="/leilao/create">Criar</a></li>
			<li><a href="/leilao/list">Listar</a></li>
		</ul>
	</section>
{/if}

<style>
	section {
		display: flex;
		flex-direction: column;
		width: 20vw;
		height: 100%;
		background-color: var(--utfpr-main-color);
		z-index: 2;
	}

	ul {
		padding: 1rem;
	}

	span {
		font-size: 24px;
	}

	li {
		padding-left: 1rem;
		list-style-type: none;
	}

	a {
		color: black;
	}
</style>
