<script lang="ts">
	export let label = '';
	export let type = 'text';
	export let placeholder = '';
	export let name = '';
	export let value: string;
	export let disabled: boolean = false;
</script>

<div>
	{#if label}
		<div>
			<label for={name}>{label}</label>
		</div>
	{/if}
	<input {disabled} {type} {name} id={name} bind:value {placeholder} />
</div>

<style>
	input {
		background-color: transparent;
		border: 1px solid black;
		width: var(--width, fit-content);
		border-radius: 1rem;
		height: 2rem;
		padding: 0.5rem;
	}
</style>
