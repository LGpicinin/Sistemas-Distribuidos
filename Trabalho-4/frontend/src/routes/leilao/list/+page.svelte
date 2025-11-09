<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import type { Leilao, LeilaoPlus } from '$lib/helpers/models/leilao.js';
	import type { Interest } from '$lib/helpers/models/interest.js';
	import Button from '$lib/components/button.svelte';
	import { fail, redirect } from '@sveltejs/kit';

	let { data } = $props();

	let leiloes: LeilaoPlus[] = data.leiloes;

	const clickButton = async (userId: string, leilao: LeilaoPlus, index: number) => {
		console.log("entreiii")
		const interest: Interest = {
			UserId: userId,
			LeilaoId: leilao.leilao.id
		};

		const route = leilao.notificar ? "cancel" : "register";

		try {
			const response = await fetch(`http://localhost:5059/leilao/${route}`, {
				method: 'POST',
				body: JSON.stringify(interest),
				headers: {
					'Content-Type': 'application/json'
				}
			});
			leilao.notificar = !leilao.notificar
			leiloes[index] = leilao
			if (response.status === 201) return redirect(303, 'leilao/list');
			else throw new Error(`Erro ao cancelar/registrar interesse: ${response.statusText}`);
		} catch (error: unknown) {
			if (error instanceof Error)
				return fail(403, {
					error: error.message
				});
		}
	}

	
</script>

{#if leiloes !== null}
	{#each leiloes as leilao, index}
		<Card>
			<p>{leilao.leilao.id}</p>
			<p>{leilao.leilao.description}</p>
			<p>{leilao.leilao.start_date}</p>
			<p>{leilao.leilao.end_date}</p>

			<Button
				text={leilao.notificar ? "Cancelar Interesse" : "Registrar Interesse"}
				type="button"
				onclick={() => clickButton(data.userId, leilao, index)}
				--width="10rem"
				--color={leilao.notificar ? "red" : "green"}
			/>
		</Card>
	{/each}
{:else}
	oi
{/if}
