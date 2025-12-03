<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import Input from '$lib/components/input.svelte';
	import Button from '$lib/components/button.svelte';
	import { type Leilao } from '$lib/helpers/models/leilao';
	import { enhance } from '$app/forms';

	let { form } = $props();

	let leilaoData: Leilao = $state({
		id: '',
		description: '',
		start_date: '',
		end_date: ''
	});
</script>

<Card>
	<form class="subcard" method="POST" action="?/createLeilao" use:enhance>
		<h2>Criar Leilão</h2>
		<div class="grid">
			<Input
				type="text"
				name="id"
				bind:value={leilaoData.id}
				label="ID do Leilão"
				--width="20rem"
			/>
			<Button
				text="Aleatório"
				type="button"
				onclick={() => (leilaoData.id = crypto.randomUUID())}
				--width="10rem"
				--color="orangered"
			/>
			<Input
				type="text"
				name="description"
				bind:value={leilaoData.description}
				label="Descrição do leilão"
				--width="20rem"
			/>
			<Button
				type="button"
				text="Aleatório"
				onclick={() => (leilaoData.description = crypto.randomUUID())}
				--width="10rem"
				--color="orangered"
			/>
			<Input
				type="datetime-local"
				name="start_date"
				bind:value={leilaoData.start_date}
				label="Início do leilão"
				--width="20rem"
			/>
			<Button
				type="button"
				text="Agora"
				onclick={() => (leilaoData.start_date = new Date().toISOString().slice(0, 16))}
				--width="10rem"
				--color="orangered"
			/>
			<Input
				type="datetime-local"
				name="end_date"
				bind:value={leilaoData.end_date}
				label="Fim do leilão"
				--width="20rem"
			/>
			<Button
				type="button"
				text="Agora"
				onclick={() => (leilaoData.end_date = new Date().toISOString().slice(0, 16))}
				--width="10rem"
				--color="orangered"
			/>
		</div>
		<Button text="Cadastrar" type="submit" --color="green" />
		{#if form?.error}
			<p>Erro ao criar Leilão: {form.error}</p>
		{/if}
	</form>
</Card>

<style>
	h2 {
		margin-bottom: 2rem;
	}

	.subcard {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2rem;
	}

	.grid {
		display: grid;
		grid-template-columns: 2fr 1fr;
		column-gap: 1rem;
		row-gap: 1rem;
		align-items: end;
	}

	p {
		color: red;
	}
</style>
