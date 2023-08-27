<style>
</style>

<script>
import Post from "./Post.vue"

export default {
  // props: ['postid'],
	components: {
		Post
	},
	data() {
		return {
			errormsg: null,
			loading: false,
      userid: '',
			user: {},

		}
	},
  methods: {
		async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/"+this.userid);
        this.user = response.data;

				// if (this.user['user-post-ids'].length === 0){
				// 	this.$router.push("/create-post");
				// }
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },


	},
	created(){
		this.userid = localStorage.getItem('userid')
	},
	async mounted(){
		await this.getUserProfile();

	}

}
</script>
<template>
<div>

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	<LoadingSpinner v-if="loading"></LoadingSpinner>

	<Post v-for="pid in user['user-post-ids']" :postid="pid"></Post>

</div>
</template>