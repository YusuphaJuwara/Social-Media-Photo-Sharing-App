<style scoped>
.error {
  color: red;
}

.hint {
  color: gray;
}

.card-with-bg {
  margin-bottom: 20px;
  background-color:beige;
}

.user-icon {
	width: 500px;
	height: 500px;
  vertical-align: middle;
	border-radius: 50%;
	object-fit: cover;
	display: inline-block;
}

/* .btn-primary[disabled] {
  background-color: rgb(223, 82, 35);
  cursor: not-allowed;
} */
</style>

<script>
import Post from "./Post.vue"

export default {
  components:{
    Post,
  },

	data() {
		return {
      errormsg: null,
      loading: false,
      users: [],
      userids: [],
      postids: [],
      profilePics: [],

      name_hashtag: '',

      startRendering: false,

    }
	},

  methods: {
    async getUserProfile(uid) {
      this.loading = true;
      this.errormsg = null;
      let response = {};
      try {
        response = await this.$axios.get("/users/"+uid);
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
      return response.data;
    },

    async getUserProfiles(){
      this.users = [];
			for (const uid of this.userids) {
				const uDet = await this.getUserProfile(uid);
        this.users.push(uDet)
			}
		},

    async getUserProfilePicture(uid) {
      this.loading = true;
      this.errormsg = null;
      let photo = null;
      try {
        let response = await this.$axios.get("/users/"+uid+"/profile-picture", {responseType: "blob"});
        photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
      } catch (e) {
        this.errormsg = e.toString();
      }finally {
        this.loading = false;
      }
      return photo;
    }, 

    async getUserProfilePictures() {
			this.loading = true;
			this.errormsg = null;
			try {
        this.profilePics = [];
        for (const uid of this.userids) {
          const photo = await this.getUserProfilePicture(uid)
          this.profilePics.push(photo);
        }
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

    async search() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/search-name-or-hashtag?name-hashtag="+this.name_hashtag);

				this.userids = [];
				if (response.data['userIDs'] != null){
					this.userids = response.data['userIDs'];
          await this.getUserProfiles();
          await this.getUserProfilePictures();
				}

				this.postids = [];
				if (response.data['postIDs'] != null){
					this.postids = response.data['postIDs'];
				}

        this.startRendering = true;

			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

    async onPostDeleted(){
			await this.search();
		},

	},

}
</script>

<template>
  <div>
  
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <LoadingSpinner v-if="loading"></LoadingSpinner>

    <div>
      <h1 style="color: blueviolet;" > Search a user by his profile name or a post by its hashtag </h1>
      <div class="row m-3">
        <div class="col-md-8 m-1">
          <textarea id="namehastag" 
            class="form-control" 
            v-model="name_hashtag" 
            placeholder="Profile Name: Enrico204. Hashtag: nicePic">
          </textarea>
        </div>
        <div class="col-md-3 m-1 d-flex justify-content-end">
          <button type="button" class="btn btn-primary" @click="search">Search</button>
        </div>
      </div>
    </div>
  
    <div class="container m-5" style="border: 4px solid red;" v-if="!loading && startRendering===true">
      <div class="card-text" 
      style="color: brown;font-size: xx-large;" 
      v-if="userids.length===0"> 
        No user profile name corresponds to the search term.
      </div>
      <div class="card m-2" style="border: 1px solid rgb(0, 94, 255);" v-if="userids.length>0" v-for="(user, idx) in users" :key="user['user-id']">
        <div class="card-header m-2">
          <LinkToUserProfile
          :profpic="profilePics[idx]" 
          :userprofname="user['profile-name']"
          :uid="user['user-id']">
          </LinkToUserProfile>
        </div>
        <div class="card-body">
          <div class="row">
            <div class="col-md-4 d-flex justify-content-center">Post Count: {{ user['post-count'] }}</div>
            <div class="col-md-4 d-flex justify-content-center">Follower Count: {{ user['follower-count'] }}</div>
            <div class="col-md-4 d-flex justify-content-center">Following Count: {{ user['following-count'] }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="container m-5" style="border: 4px solid red;" v-if="!loading && postids.length===0 && startRendering===true">
      <div class="card-text" style="color: brown;font-size: xx-large;" > 
          No post hashtag corresponds to the search term.
      </div>
    </div>

    <Post v-if="!loading && postids.length>0" v-for="pid in postids" :postid="pid" :key="pid" @postDeleted="onPostDeleted"></Post>
  
  </div>
  </template>

