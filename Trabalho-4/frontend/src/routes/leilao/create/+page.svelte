<script lang="ts">
    import Card from "$lib/components/card.svelte";
    
    let leilaoData : Leilao = {
        id: "",
        description: "",
        start_date: new Date(),
        end_date: new Date()
    };
</script>

<Card>
    <h2>Criar Leilão</h2>
    
    <input type="text" name="leilaoId" id="leilaoId" bind:value={leilaoData.id}>
    <input type="text" name="leilaoDescription" id="leilaoDescription" bind:value={leilaoData.description}>
    <input type="datetime-local" name="leilaoStart" id="leilaoStart" bind:value={leilaoData.start_date}>
    <input type="datetime-local" name="leilaoEnd" id="leilaoEnd" bind:value={leilaoData.end_date}>
    
</Card>

<style>

</style>