<style scoped>
</style>

<script>
export default {
	data() {
		return {
			errormsg: null,
			loading: false,

		}
	},
  methods: {
    async logOut() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/session");

				// invalidate the data ...
				localStorage.removeItem('token');
				localStorage.removeItem('userid');

				this.$router.push("/login");
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
        this.loading = false;
      }
		},

	},

}
</script>

<template>
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <LoadingSpinner v-if="loading"></LoadingSpinner>

		<div class="d-grid justify-content-center m-5 p-5">
			<div style="color: tomato;font-size: xx-large;"> Do you want to log out?</div>
			<button class="btn btn-primary m-4" type="button" @click="logOut">Log out</button>
		</div>
  
  </div>
</template>