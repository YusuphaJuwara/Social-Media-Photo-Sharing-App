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
      uid: '',
			user: {},

		}
	},
  methods: {
		async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/"+this.uid);
        this.user = response.data;

        console.log("ProfileOtherUsers: user: "+this.user);

      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

		async onPostDeleted(){
			await this.getUserProfile();
		},


	},
	created(){
		this.userid = localStorage.getItem('userid')
    this.uid = this.$route.params.uid
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

	<ProfileDetails :userid="uid"></ProfileDetails>

	<Post v-for="pid in user['user-post-ids']" :postid="pid" @postDeleted="onPostDeleted"></Post>

</div>
</template>