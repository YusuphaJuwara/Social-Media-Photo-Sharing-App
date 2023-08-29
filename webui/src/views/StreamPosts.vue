<style scoped>
</style>

<script>
import Post from "./Post.vue"

export default {
	components: {
		Post
	},
	data() {
		return {
			errormsg: null,
			loading: false,
      userid: '',
      posts: [],

		}
	},
  methods: {
    async getMyStream() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + this.userid + "/posts/stream/");
				this.posts = response.data;

				if (this.posts.length === 0){
					this.$router.push("/create-post");
				}
			} catch (e) {
				this.errormsg = e.toString();
			} finally{
        this.loading = false;
      }
		},

		async onPostDeleted(){
			await this.getMyStream();
		},

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
	created(){
		this.userid = localStorage.getItem('userid')
	},
	async mounted(){
    await this.getMyStream();

	}

}
</script>

<template>
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <LoadingSpinner v-if="loading"></LoadingSpinner>
  
    <Post v-for="post in posts" :postid="post['post-id']" @postDeleted="onPostDeleted"></Post>
  </div>
</template>