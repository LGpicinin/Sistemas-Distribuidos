<div class="card">
	<slot></slot>
</div>

<style>
	.card {
		width: 50%;
		border: 2px solid black;
		border-radius: 1rem;
		padding: 2.5rem;
		margin: 1rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		background-color: var(--bgColor, transparent);
	}
</style>
