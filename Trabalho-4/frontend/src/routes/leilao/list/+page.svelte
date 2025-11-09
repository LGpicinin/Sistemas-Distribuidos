<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import type { LeilaoPlus } from '$lib/helpers/models/leilao.js';
	import type { Interest } from '$lib/helpers/models/interest.js';
	import Button from '$lib/components/button.svelte';
	import { fail, redirect } from '@sveltejs/kit';

	let { data } = $props();

	let leiloes: LeilaoPlus[] = $state(data.leiloes);

	const changeInterest = async (userId: string, leilao: LeilaoPlus, index: number) => {
		const interest: Interest = {
			UserId: userId,
			LeilaoId: leilao.leilao.id
		};

		const route = leilao.notificar ? 'cancel' : 'register';

		const response = await fetch(`http://localhost:5059/${route}?userId=${data.userId}`, {
			method: 'POST',
			body: JSON.stringify(interest),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (response.status === 200) {
			leilao.notificar = !leilao.notificar;
			leiloes[index] = leilao;
		}
	};
</script>

<Card>
	<h2>Leilões Ativos</h2>
	{#if leiloes.length}
		{#each leiloes as leilao, index}
			<Card>
				<h3>Leilão {index + 1}</h3>
				<div class="info">
					<p><strong>ID:</strong> {leilao.leilao.id}</p>
					<p><strong>Descrição:</strong> {leilao.leilao.description}</p>
					<p><strong>Data e hora de início:</strong> {leilao.leilao.start_date}</p>
					<p><strong>Data e hora de término:</strong> {leilao.leilao.end_date}</p>
				</div>

				<Button
					text={leilao.notificar ? 'Cancelar Interesse' : 'Registrar Interesse'}
					type="button"
					onclick={() => changeInterest(data.userId, leilao, index)}
					--width="10rem"
					--color={leilao.notificar ? 'red' : 'green'}
				/>
			</Card>
		{/each}
	{:else}
		<p>Não há leilões ativos.</p>
	{/if}
</Card>

<style>
	.info {
		margin: 1rem;
	}
</style>
