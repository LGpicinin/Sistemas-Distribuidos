<script lang="ts">
	import Button from '$lib/components/button.svelte';
	import Card from '$lib/components/card.svelte';
	import Input from '$lib/components/input.svelte';
	import { type Lance } from '$lib/helpers/models/lance.js';

	const { data } = $props();

	let lance: Lance = {
		user_id: data.userId!,
		leilao_id: data.leilaoId!,
		value: '0'
	};
</script>

<Card>
	<form action="?/createLance" method="post">
		<h2>Novo Lance</h2>
		<Input
			type="text"
			name="leilao_id"
			bind:value={lance.leilao_id}
			label="ID do leilão"
			disabled={true}
			--width="20rem"
		/>
		<Input
			type="text"
			name="user_id"
			bind:value={lance.user_id}
			label="ID do usuário"
			disabled={true}
			--width="20rem"
		/>
		<Input
			type="number"
			name="value"
			bind:value={lance.value}
			label="Valor do lance"
			--width="20rem"
		/>

		<Button type="submit" text="Realizar" --width="10rem" --color="blue" />
	</form>
</Card>

<style>
	form {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}
</style>
