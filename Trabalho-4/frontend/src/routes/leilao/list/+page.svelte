<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import type { LeilaoPlus } from '$lib/helpers/models/leilao.js';
	import type { Interest } from '$lib/helpers/models/interest.js';
	import Button from '$lib/components/button.svelte';

	interface LeilaoPlusPlus extends LeilaoPlus {
		href: string;
	}

	let { data } = $props();

	let leiloesPlus: LeilaoPlus[] = $state(data.leiloes);
	let leiloesPlusPlus : LeilaoPlusPlus[] = $state([]);
	let i = 0
	for (i=0; i<leiloesPlus.length; i++){
		leiloesPlusPlus.push({...leiloesPlus[i], href:`/leilao/${leiloesPlus[i].leilao.id}`})
	}

	const changeInterest = async (userId: string, leilao: LeilaoPlusPlus, index: number) => {
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
			leiloesPlusPlus[index] = leilao;
		}
	};
</script>

<Card>
	<h2>Leilões Ativos</h2>
	{#if leiloesPlusPlus.length}
		{#each leiloesPlusPlus as leilao, index}
			<Card>
				<h3>Leilão {index + 1}</h3>
				<div class="info">
					<p><strong>ID:</strong> {leilao.leilao.id}</p>
					<p><strong>Descrição:</strong> {leilao.leilao.description}</p>
					<p><strong>Data e hora de início:</strong> {leilao.leilao.start_date}</p>
					<p><strong>Data e hora de término:</strong> {leilao.leilao.end_date}</p>
				</div>

				<div>
					<Button
						text={leilao.notificar ? 'Cancelar Interesse' : 'Registrar Interesse'}
						type="button"
						onclick={() => changeInterest(data.userId, leilao, index)}
						--width="10rem"
						--color={leilao.notificar ? 'red' : 'green'}
					/>
					<a href={leilao.href} class="link" >Realizar lance</a>
					<!-- <Button
						text="Realizar lance"
						type="button"
						href={leilao.href}
						onclick={() => {
							
						}}
						--width="10rem"
						--color="blue"
					/> -->
				</div>
			</Card>
		{/each}
	{:else}
		<p>Não há leilões ativos.</p>
	{/if}
</Card>

<style>
	.info {
		margin: 1rem;
		width: 100%;
	}
	.link {
		background-color: var(--bg-color, transparent);
		color: var(--color, blue);
		border: 1px solid var(--color, blue);
		border-radius: 1rem;
		padding: 0.5rem 1rem;
		cursor: pointer;
		width: var(--width, fit-content);
		height: var(--height, fit-content);
		text-decoration: none;
		font-size: 13px;
	}
</style>
