<template>
	<div class="backdrop" @click="emitClose">
		<div class="alert" role="alert">
			<h3>Error:</h3>
			<h5>{{ msg }}</h5>
			<p class="small">(Click anywhere to close)</p>
		</div>
	</div>
</template>

<script>
export default {
	props: ['msg'],
	methods: {
		emitClose() {
			this.$emit('close')
		},
	}
}
</script>

<style scoped>
.backdrop {
    top: 0;
    position: fixed;
    background-color: rgba(0,0,0,0.5);
    width: 100%;
    height: 100%;
}
.alert  {
	width: 600px;
    height: 200px;
    background-color: rgba(169, 1, 1, 0.817);
    color: var(--light-colour);
    border-radius: 10px;
    position: fixed;
    left: 33%;
    top: 20%;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}
p.small {
	font-size: small;
	margin: 30px 0px 0px 0px;
}
</style>
