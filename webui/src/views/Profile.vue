<style>
</style>

<script>
import Post from "./Post.vue"
import ProfileDetails from "./ProfileDetails.vue";

export default {
	components: {
		Post,
		ProfileDetails,
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

				console.log("Profile: user: "+this.user);

      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

		// An event listener to when a post is deleted in Post.vue
		async onPostDeleted(){
			await this.getUserProfile();
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

	<ProfileDetails :userid="userid"></ProfileDetails>

	<Post v-for="pid in user['user-post-ids']" :postid="pid" @postDeleted="onPostDeleted"></Post>

</div>
</template>