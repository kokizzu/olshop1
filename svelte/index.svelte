<script>
	// should generate index.html
	import { GuestRegister } from "./jsApi.GEN";

	let email = '';
	let password = '';

	async function doRegister() {
		await GuestRegister({
			email: email,
			password: password,
		}, async function(data) {
			console.log('data', data)
		})
	}
</script>

<form on:submit|preventDefault={doRegister}>
	<label for="email">Email</label>
	<input id="email" bind:value={email}>
	<label for="password">Password</label>
	<input type="password" bind:value={password}>

	<input type="submit" value="Register">
</form>