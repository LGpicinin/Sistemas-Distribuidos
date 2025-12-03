<script lang="ts">
	import Input from '$lib/components/input.svelte';
	import Button from '$lib/components/button.svelte';
	import Card from '$lib/components/card.svelte';
	let userIdInput: string;
</script>

<Card>
	<form method="POST" action="?/setUserId">
		<section class="userIdSection">
			<Input label="ID do Usuário" bind:value={userIdInput} name="userId" --width="20rem" />
			<Button
				type="button"
				onclick={() => {
					userIdInput = crypto.randomUUID();
				}}
				text="Criar ID aleatório"
				--color="brown"
			/>
		</section>
		<Button --color="green" type="submit" text="Confirmar" />
	</form>
</Card>

<style>
	form {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1rem;
	}

	.userIdSection {
		display: flex;
		align-items: end;
		gap: 1rem;
	}
</style>
