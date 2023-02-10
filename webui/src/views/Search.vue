<style scoped>
.card {
	margin-bottom: 20px;
	border-bottom: 2px solid #333;
}

/* .imgThumbNail {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
} */

/* .user-icon {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 8px;
	display: inline-block;
} */

.user-name {
	color: #333;
	cursor: pointer;
	display: inline-block;
	vertical-align: middle;
}

.highlighted {
	color: blue;
}
</style>

<script>
export default {
	data: function () {
		return {
			token: null,
			userid: null,
			name_hashtag: "",
			highlightProfile: false,
			errormsg: null,
			loading: false,
			userids: [],
			postids: [],
		}
	},
	methods: {
		load() {
			return load
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
			}
			this.loading = false;
		},
		async search() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/search-name-or-hashtag?name-hashtag="+this.name_hashtag);
				this.userids = response.data['userIDs'];
				this.postids = response.data['postIDs'];

			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		
		async home() {
			this.highlightProfile = !this.highlightProfile
			this.$router.push("/"+this.userid+"/stream/");
		},
		async getUserProfilePicture(userid) {
			// this.loading = true;
			this.errormsg = null;
			let photo = ''
			try {
				let response = await this.$axios.get("/users/" + userid + "/profile-picture");
				photo = URL.createObjectURL(new Blob([response.data]));

				// refresh...
				// await this.search();
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
			return photo;
		},
		async getProfileName(userid) {
			// this.loading = true;
			this.errormsg = null;
			let profileName = ''
			try {
				let response = await this.$axios.get("/users/" + userid);
				profileName = response.data['profile-name']

				// refresh...
				// await this.search();
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
			return profileName;
		},
		userProfile(userid) {
			this.highlightProfile = !this.highlightProfile
			this.$router.push("/"+userid+"/profile");
			// this.$route.params.id
		},
		viewPost(postid) {
			this.highlightProfile = !this.highlightProfile
			this.$router.push("/"+postid+"/post");
			// this.$route.params.id
		},
		async getPhotoDetails(postid) {
			// this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/posts/"+postid);
				var photoid =response.data['photo-id'];
				var uidd = response.data['user-id'];

				// refresh...
				// await this.search();
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
			return photoid, uidd;
		},
		async getSinglePhoto(photoid) {
			// this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/photos/" + photoid);
				var photo = URL.createObjectURL(new Blob([response.data]));

				// refresh...
				// await this.getUserPhotos();
			} catch (e) {
				this.errormsg = e.toString();
			}
			// this.loading = false;
			return photo;
		},
  },
	mounted() {
		this.search()
		this.userid = localStorage.getItem('userid')
    this.token = localStorage.getItem('token')
	}
}

</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">{{ getProfileName(userid) }}'s search</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="logOut">
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#log-out" />
						</svg>
						Logout
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="search">
						Refresh search
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="home">
						Stream Home Page
					</button>
				</div>
        <div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="userProfile(userid)">
						Your Profile Page
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div class="row">
			<h1 > Search a user by his profile name or search a post by its hashtag </h1>
			<input type="text" v-model="name_hashtag" class="form-control"/>
			<button type="button" class="btn btn-sm btn-outline-secondary" @click="search"></button>"

			<div class="col-lg-6 col-md-12">
				<h1 > Users </h1>
				<div class="card" v-if="!loading" v-for="idd in userids" :key="idd">
					<div class="card-header">
						<a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
							@click="userProfile(idd)">
							<img :src="getUserProfilePicture(idd)" width="80rem" height="80rem" class="rounded-circle mx-auto" />
							<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ getProfileName(idd) }}
							</h5>
						</a>
					</div>
				</div>
			</div>

			<div class="col-lg-6 col-md-12">
				<h1 > Posts </h1>
				<div class="card" v-if="!loading" v-for="idd in postids" :key="idd">
          <div class="card-header">
            <div class="header-left">
              <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
                @click="viewPost(idd)">
                <img :src="getSinglePhoto(getPhotoDetails(idd)[0])" width="80rem" height="80rem" class="rounded-circle mx-auto" />
                <h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ getProfileName(getPhotoDetails(idd)[1]) }}
                </h5>
              </a>
            </div>
          </div>
        </div>
			</div>

		</div>
	</div>
</template>
