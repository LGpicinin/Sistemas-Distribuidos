<script lang="ts">
	let userIdInput: string;
</script>

<div class="card">
	<form method="POST" action="?/setUserId">
		<label for="userId">ID do Usuário:</label>
		<section class="userIdSection">
			<input type="text" name="userId" id="userId" bind:value={userIdInput} />
			<button
				type="button"
				onclick={() => {
					userIdInput = crypto.randomUUID();
				}}>Criar ID aleatório</button
			>
		</section>
		<button type="submit">Confirmar</button>
	</form>
</div>

<style>
	.card {
		width: 50%;
		border: 2px solid black;
		border-radius: 1rem;
		padding: 5rem;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.userIdSection {
		display: flex;
		/* flex-direction: column; */
	}
</style>
