<script lang="ts">
	let props = $props();
</script>

<button {...props} class="button">{props.text}</button>

<style>
	.button {
		background-color: var(--bg-color, transparent);
		color: var(--color, black);
		border: 1px solid var(--color, black);
		border-radius: 1rem;
		padding: 0.5rem 1rem;
		cursor: pointer;
		width: var(--width, fit-content);
		height: var(--height, fit-content);
	}
</style>
