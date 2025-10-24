<script lang="ts">
	import { fade } from 'svelte/transition';

	export let isOpen;
</script>

{#if isOpen}
	<section transition:fade>aaaaaaaaaaaaaa</section>
{/if}

<style>
	section {
		display: flex;
		flex-direction: column;
		width: 20vw;
		height: 100%;
		background-color: var(--utfpr-main-color);
	}
</style>
