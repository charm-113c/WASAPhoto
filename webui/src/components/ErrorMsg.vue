<script>
export default {
	props: ['msg', 'code']
}
</script>

<template>
	<div class="alert alert-danger" role="alert">
		<p>{{ code }}</p>
		<p>{{ msg }}</p>
	</div>
</template>

<style>
</style>
