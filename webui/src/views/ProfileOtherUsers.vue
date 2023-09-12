<style>
</style>

<script>
import ProfileDetails from "./ProfileDetails.vue";
import Post from "./Post.vue"

export default {
	components: {
		ProfileDetails,
		Post,
	},
	data() {
		return {
			errormsg: null,
			loading: false,
      userid: '',
      uid: '',
			user: {},
			posts: [],

			// Check if the logged in user (userid) can see the posts of the uid
			seePosts: false,

		}
	},
  methods: {
		async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/"+this.uid);
        this.user = response.data;

				if (this.user['user-post-ids'] != null){
					this.posts = this.user['user-post-ids'].slice().reverse();
				} else {
					this.posts = [];
				}

        console.log("ProfileOtherUsers: user: "+this.user);

				// True if follows ...
				const checkFollow = await this.getUserFollows();

				response = await this.$axios.get("/users/" + this.uid + "/private-profile");
      	const pri = response.data ? true : false

				// If the profile is private and the user doesn't follow him, then show only prof details, not posts.
				if (!pri || (pri && checkFollow)){
					this.seePosts = true;
				}
				console.log("\nprivate: "+pri+"\ncheckFollow: "+checkFollow+"\nresponse.data: "+response.data+"\nseePosts = false: "+this.seePosts)

      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

		async getUserFollows() {
      this.loading = true;
      this.errormsg = null;
			let checkFollow = false;
      try {
        const response = await this.$axios.get("/users/"+this.uid+"/follow");

				console.log("getUserFollows() in ProfileOtherUsers: \nresponse.data['followers-array']: "
				+response.data['followers-array']+"\nLoggedin user this.userid: "+this.userid)

				// Check if the logged in user follows the user whose info is been requested.
				if (response.data['followers-array'] != null){
					if (response.data['followers-array'].includes(this.userid)){
						checkFollow = true;
					}
        }
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
			return checkFollow;
    },

		async onPostDeleted(){
			await this.getUserProfile();
		},


	},
	async created(){
		// Logged in user doing the request
		this.userid = localStorage.getItem('userid')

		// User whose profile is been requested
    this.uid = this.$route.params.uid

		await this.getUserProfile();
	},
	async mounted(){
		// await this.getUserProfile();

	}

}
</script>
<template>
<div>

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	<LoadingSpinner v-if="loading"></LoadingSpinner>

	<ProfileDetails :userid="uid"></ProfileDetails>

	<Post v-if="seePosts" v-for="pid in user['user-post-ids']" :postid="pid" :key="pid" @postDeleted="onPostDeleted"></Post>

</div>
</template>