<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import type { Leilao } from '$lib/helpers/models/leilao.js';

	let { data } = $props();

	const leiloes: Leilao[] = data.leiloes;
</script>

<Card>
	{#if leiloes !== null}
		{#each leiloes as leilao}
			<Card>
				<p><strong>ID:</strong> {leilao.id}</p>
				<p><strong>Descrição:</strong> {leilao.description}</p>
				<p><strong>Data e hora de início:</strong> {leilao.start_date}</p>
				<p><strong>Data e hora de término:</strong> {leilao.end_date}</p>
			</Card>
		{/each}
	{:else}
		<p>Não há leilões ativos.</p>
	{/if}
</Card>
