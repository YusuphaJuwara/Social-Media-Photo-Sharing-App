<!-- Delete later; testing 'repeat' -->

<style scoped>
.card-list {
  display: grid;
  grid-gap: 1em;
}

.card-item {
  /* background-color: dodgerblue; */
  padding: 2sp;
}

#app {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
  transition: all 0.2s;
}

body {
  background: #20262E;
  padding: 2sp;
  font-family: Helvetica;
}

ul {
  list-style-type: none;
}
</style>

<script>
export default {
	data() {
		return {
			errormsg: null,
			loading: false,
      userid: '',
			user: {},
      posts: [],
			// post: {},

			profPic: null,
			pics: [],

			numberOfColumns: 4,


		}
	},

	computed: {
    gridStyle() {
      return {
        gridTemplateColumns: `repeat(${this.numberOfColumns}, minmax(300px, 1fr))`
      }
    },
  },

  methods: {
		// addCard() {
    //   this.cards.push('new-card')
    // },
		// load() {
		// 	return load
		// },
		async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/users/"+this.userid);
        this.user = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      } finally {
				this.loading = false;
			}
    },

		async getUserProfilePicture() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get(`/users/"${this.user['user-id']}/profile-picture`, {responseType: "blob"})
				this.profPic = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

		async getSinglePhoto(pUid) {
			this.loading = true;
			this.errormsg = null;
			let photo = '';
			try {
				let response = await this.$axios.get("/photos/" + pUid, {responseType: "blob"});
				photo = URL.createObjectURL(new Blob([response.data]), { type: "image/png" });

				// refresh...
				// await this.getPhoto();
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
			return photo;
		},

		// async getPhotos() {
		// 	this.loading = true;
		// 	this.errormsg = null;
		// 	try {
		// 		let response = await this.$axios.get("/posts/");
		// 		this.posts = response.data;
		// 	} catch (e) {
		// 		this.errormsg = e.toString();
		// 	} finally {
    //     this.loading = false;
    //   }
		// },

		async userPhotos(){
			this.posts = [];
			this.pics = [];
			await this.getUserProfile();
			for (var postid of this.user['user-post-ids']){
				var p = await this.getPhoto(postid);
				this.posts.push(p);

				var photo = await this.getSinglePhoto(p['photo-id']);
				this.pics.push(photo);
			}

		},
		async getPhoto(pid) {
			this.loading = true;
			this.errormsg = null;
			let p = null;
			try {
				let response = await this.$axios.get("/posts/"+pid);
				p = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
			return p;
		},

		async deletePhoto(pid) {
			this.errormsg = null;
			if (confirm('Are you sure you want to delete?')) {
				this.loading = true;
				try {
					await this.$axios.delete("/users/" +this.userid+"/posts/"+pid);

					// refresh ...
					await this.userPhotos();
				} catch (e) {
					this.errormsg = e.toString();
				} finally {
					this.loading = false;
				}
			}
		},


	},
	created(){
		this.userid = localStorage.getItem('userid')
	},
	async mounted(){
		await this.getUserProfile();
		await this.getUserProfilePicture();
		await this.userPhotos();
    // await this.getPhotos();

	}

}
</script>

<template>
<div>

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	<LoadingSpinner v-if="loading"></LoadingSpinner>

	Columns: <input v-model.number="numberOfColumns">
		<!-- <button @click="addCard">Add card</button> -->
	<div class="card m-5" v-if="!loading" id="app" >
		<ul :style="gridStyle" class="card-list">
			<li v-for="(post, idx) in posts" :key="post['post-id']" class="card-item">
				<button class="btn btn-danger" type="button" @click="deletePhoto(post['post-id'])" > Delete Post</button>
				<div class="container m-5 p-5">
					<img class="card-img-top" style="width:100%" :src="pics[idx]" alt="oops, can't load photo!" />
				</div>
			</li>
		</ul>
	</div>

</div>
</template>
