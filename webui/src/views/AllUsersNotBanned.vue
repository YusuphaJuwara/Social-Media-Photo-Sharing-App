<style scoped>
.card {
	margin-bottom: 20px;
}

.imgThumbNail {
	width: 20px;
	height: 20px;
	border-radius: 50%;
	object-fit: cover;
}

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
			token: localStorage.getItem('token'),
			userid: localStorage.getItem('userid'),
			highlightProfile: false,
			errormsg: null,
			loading: false,
			users: [],
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
		async getAllUsers() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/");
				this.users = response.data;
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
				// await this.getAllUsers();
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

				//delete later
				/*const headersArray = [];
					Object.entries(response.data).forEach(([key, value]) => {
						headersArray.push(`${key}: ${value}`);
					});

				confirm("Profile name: " + profileName+"\n response.data: " + headersArray)
				*/
				//////////////////////////////////7

				// refresh...
				// await this.getAllUsers();
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
  },
	beforeCreate() {
    // Initialize variables here
    this.userid = localStorage.getItem('userid');
    this.token = localStorage.getItem('token');
  },
	mounted() {
		this.getAllUsers()
	}
}

</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">{{ getProfileName(userid) }}'s getAllUsers excluding 'ban' and ...</h1>
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
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="getAllUsers">
						Refresh getAllUsers
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false" class="btn btn-sm btn-outline-primary" @click="home">
						Stream Home Page
					</button>
				</div>
        <div class="btn-group me-2">
					<button type="button" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false" class="btn btn-sm btn-outline-primary" @click="userProfile(userid)">
						Your Profile Page
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div class="card" v-if="!loading" v-for="user in users" :key="user['user-id']">
			<div class="card-header">
				<div class="header-left">
					<a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
						@click="userProfile(user['user-id'])">
						<img class="imgThumbNail" :src="getUserProfilePicture(user['user-id'])" />
						<h5 class="user-name" :class="{ 'highlighted': highlightProfile }">{{ getProfileName(user['user-id']) }}
						</h5>
					</a>
				</div>
			</div>
		</div>
	</div>
</template>
