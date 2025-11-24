<script lang="ts">
	import Button from './button.svelte';
	import type { Notification } from '$lib/helpers/models/notification';

	let { message }: { message: Notification } = $props();
	let clickPagamento = $state(false);

	const lanceTypes: Record<string, string> = {
		lance_validado: 'Novo Lance Válido',
		lance_invalidado: 'Novo Lance Inválido',
		leilao_vencedor: 'Lance Vencedor',
		link_pagamento: 'Link para Pagamento',
		status_pagamento: 'Status do Pagamento'
	};

	const realizar_pagamento = async (link: string, message: any) => {
		const response = await fetch(link, {
			method: 'POST',
			body: JSON.stringify(message),
			headers: {
				'Content-Type': 'application/json'
			}
		});

		clickPagamento = true;

		// return fail(403, {
		// 	error: `Erro ao criar leilão: ${response.statusText}`
		// });
	};
</script>

<h4>
	{lanceTypes[message.type]}
</h4>
<div class="info">
	{#if message.type != 'status_pagamento' && message.type != 'link_pagamento'}
		<p><strong>ID do Leilão:</strong> {message.lance!.leilao_id}</p>
		<p><strong>Autor do lance:</strong> {message.lance!.user_id}</p>
		<p><strong>Valor do lance:</strong> R$ {message.lance!.value}.00</p>
	{:else if message.type == 'link_pagamento'}
		{#if clickPagamento == false}
			<Button
				text={'Realizar pagamento'}
				type="disabled"
				onclick={() => realizar_pagamento(message.linkData!.link, message)}
				--width="10rem"
				--color="green"
			/>
		{/if}
	{:else}
		<p>
			<strong>Status do pagamento:</strong>
			{message.statusData!.status == true ? 'Aprovado' : 'Recusado'}
		</p>
		<p><strong>Valor:</strong> {message.statusData!.value}</p>
	{/if}
</div>

<style>
	.info {
		margin-top: 1rem;
		width: 15vw;
	}
</style>
