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
			console.log(event);
			messages = [JSON.parse(event.data), ...messages];
		});
	});

	onDestroy(() => {
		eventSource?.close();
	});
</script>

<Card>
	{#each messages as message}
		<Card>
			<h2>{message.type == "lance_validado" ? "Novo lance" : "Lance Vencedor"} do Leilão {message.lance.leilao_id}</h2>
			<div class="info">
				<p><strong>Valor do lance:</strong> R$ {message.lance.value}.00</p>
			</div>
		</Card>
	{/each}
</Card>
