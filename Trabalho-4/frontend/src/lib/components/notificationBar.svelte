<script lang="ts">
	import { flip } from 'svelte/animate';
	import Card from './card.svelte';
	import Notification from './notification.svelte';

	let { notifications = $bindable() } = $props();
</script>

<section>
	{#each notifications as notification, index (notification)}
		<div>
			<Card --bgColor="lightgray">
				<Notification message={notification} --bgColor="lightgray" />
			</Card>
		</div>
	{/each}
</section>

<style>
	section {
		position: absolute;
		right: 0vw;
		top: 4rem;
		display: flex;
		flex-direction: column;
		width: 15vw;
		height: 100%;
		z-index: 2;
	}
</style>
