<style>
 .error {
    color: red;
  }

  .hint {
    color: gray;
  }
</style>

<script>
// import DProfilePhoto from './DProfilePhoto.vue';
// import DPost from './DPost.vue';
// import DFollowModal from './DFollowModal.vue';

export default {
  // components: {
  //   DProfilePhoto,
  //   DPost,
  //   DFollowModal,
  // },
  data() {
    return {
      isFollowerModalOpen: false,
      profilePhoto: null,
      followerIDs: [],
      pictureError: false,
      pic: null,

      fol: false,
      ban: false,
      banUsers: [],
      user: {},
      
      // userid: this.$route.params.userid,
      userid: null,
      token: null,

      clicked: false,
      hashtags: '',
      caption: '',
      highlightProfile: false,
      showCaption: false,
      showHashtag: false,
      showComments: false,
      showLikes: false,
      errormsg: null,
      loading: false,
      posts: [],
      like: false,
      commentCount: 0,
      likedBy: [],
      messagetext: '',

      username: '',
      isValid: false,
      selectedImage: false,
    }
  },
  computed: {
    checkImage() {
      this.selectedImage = this.$refs.image.files.length > 0;
    },

    isPicValid() {
      if (!this.pic) {
        return false;
      }
      if (this.pic.type !== 'image/png') {
        return false;
      }
      return true;
    },
  },

  methods: {
    async openFollowerModal() {
      this.isFollowerModalOpen = true;
    },

    async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/"+this.userID);
        this.user = response.data;
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },

    async getUserProfilePicture(userid) {
    this.loading = true;
    this.errormsg = null;
    //let photo = ''
    try {
      let response = await this.$axios.get("/users/" + userid + "/profile-picture");
      this.profilePhoto = URL.createObjectURL(new Blob([response.data]));

      // refresh...
      // await this.getUserPhotos();
    } catch (e) {
      this.errormsg = e.toString();
    }
    this.loading = false;
    //return photo;
  }, 

  async getUserFollows(userid) {
      this.loading = true;
      this.errormsg = null;
      try {
        var response = await this.$axios.get("/users/"+userid+"/follow");
        this.followerIDs = response.data['followers-array']
        // await this.getUserPhotos();
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
      //return response.data['followers-array'], response.data['followings-array'];
    },

    async changeUserProfilePicture(userid) {
      this.loading = true;
      this.errormsg = null;
      try {
        const formData = new FormData();
        formData.append("photo", this.$refs.image.files[0]);

        await this.$axios.put("/users/"+userid+"/profile-picture", formData, {
          headers: {
            "content-type": "multipart/form-data"
          }
        });

        this.pic = null;

        await this.getUserProfile(userid);

      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },

    uploadPicture(event) {
      const file = event.target.files[0];
      if (file && file.type === 'image/png') {
        this.pic = file;
        this.pictureError = false;
      } else {
        this.pictureError = true;
        this.$refs.pictureInput.value = '';
      }
    },

//////////////////////////////
  },
  created(){
    this.userid = localStorage.getItem('userid')
    this.token = localStorage.getItem('token')
    // this.userid = this.$route.params.userid

    this.getUserProfile(this.userid)
    this.getUserProfilePicture(this.userid)
    this.getUserFollows(user.ID)

    console.log("\n\n"+
                "userid: "+this.userid+
                "\ntoken: "+this.token+
                "\nuser: "+this.user+
                "\n\n");

  }
  // ,
  // mounted() {
  //   // this.getUserProfile(this.userid)
  //   // this.getUserProfilePicture(this.userid)
  //   // this.getUserFollows(user.ID)
    
  // }
};
</script>

<template>
  <div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<LoadingSpinner v-if="loading"></LoadingSpinner>

    <!--d-profile-photo :photo-id="user.ProfilePhotoID" /-->
    <img :src="profilePhoto" alt="Oops!" width="150rem" height="150rem" class="rounded-circle mx-auto">

    <div>
      <h1>User: {{ user }}</h1>
      <h1>ID: {{ user.ID }}</h1>
      <h1>Profile Name: {{ user['profile-name'] }}</h1>
      <h1>Gender: {{ user['gender']}}</h1>
      <h1>BirthDate: {{ user.BirthDate }}</h1>
      <h1>ProfilePhotoID: {{ user.ProfilePhotoID }}</h1>
      <h1>PostCount: {{ user.PostCount }}</h1>
      <h1>FollowerCount: {{ user.FollowerCount }}</h1>
      <h1>FollowingCount: {{ user.FollowingCount }}</h1>
      <h1>ProfileMessage: {{ user.ProfileMessage }}</h1>
    </div>
      
    <!--d-post :post-ids="user.PostIDs" /-->

    <!--button @click="openFollowerModal">View Followers</button-->
    <!-- <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#staticBackdrop1">Followers: {{ user.FollowerCount }}</button> -->

    <!--d-follow-modal v-if="isFollowerModalOpen" :user-id="user.ID" /-->

    <!--Modals-->
    <!-- Dropdown for followers-->
		<!-- <div class="modal fade" id="staticBackdrop1" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
				aria-labelledby="staticBackdropLabel" aria-hidden="true">
			<div class="modal-dialog">
				<div class="modal-content">

					<div class="modal-header">
						<h1 class="modal-title fs-5" id="staticBackdropLabel">Users who followed you</h1>
						<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
					</div>

					<div class="modal-body" v-for="uid in followerIDs" :key="uid">
						<div>{{ uid }}</div>
					</div>

					<div class="modal-footer">
						<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
					</div>

				</div>
			</div>
		</div> -->
  </div>
</template>