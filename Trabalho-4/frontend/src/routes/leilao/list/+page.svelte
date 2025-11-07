<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import type { Leilao } from '$lib/helpers/models/leilao.js';

	let { data } = $props();

	const leiloes: Leilao[] = data.leiloes;
</script>

{#if leiloes !== null}
	{#each leiloes as leilao}
		<Card>
			<p>{leilao.id}</p>
			<p>{leilao.description}</p>
			<p>{leilao.start_date}</p>
			<p>{leilao.end_date}</p>
		</Card>
	{/each}
{:else}
	oi
{/if}
