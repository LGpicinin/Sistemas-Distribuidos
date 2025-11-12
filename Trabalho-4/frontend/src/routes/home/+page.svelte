<script lang="ts">
	import Card from '$lib/components/card.svelte';
	import { PUBLIC_GATEWAY_ADDRESS } from '$env/static/public';
	import { onDestroy, onMount } from 'svelte';

	const { data } = $props();

	let messages = $state([] as any[]);
	let eventSource: EventSource;

	onMount(() => {
		eventSource = new EventSource(`${PUBLIC_GATEWAY_ADDRESS}/event?userId=${data.userId}`);
		eventSource.addEventListener(data.userId, (event) => {
			console.log(event.data);
			messages = [JSON.parse(event.data), ...messages];
		});
	});

	onDestroy(() => {
		eventSource?.close();
	});

	const lanceTypes: Record<string, string> = {
		lance_validado: 'Novo Lance Válido',
		lance_invalidado: 'Novo Lance Inválido',
		leilao_vencedor: 'Lance Vencedor'
	};
</script>

<Card>
	<h2>Histórico de Notificações</h2>
	{#each messages as message}
		<Card>
			<h4>
				{lanceTypes[message.type]}
			</h4>
			<div class="info">
				<p><strong>ID do Leilão:</strong> {message.lance.leilao_id}</p>
				<p><strong>Autor do lance:</strong> {message.lance.user_id}</p>
				<p><strong>Valor do lance:</strong> R$ {message.lance.value}.00</p>
			</div>
		</Card>
	{:else}
		<p>Sem notificações para exibir</p>
	{/each}
</Card>

<style>
	.info {
		margin-top: 1rem;
	}
</style>
