<style scoped>
.card {
	margin-bottom: 20px;
}

.imgThumbNail {
	width: 40px;
	height: 40px;
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
	props: ['profpic', 'userprofname', 'uid'],
	data() {
		return {
			errormsg: null,
			loading: false,
			highlightProfile: false,
      userid: '',
      posts: [],

		}
	},
  methods: {
		async userProfile() {
			this.highlightProfile = !this.highlightProfile

			if (this.userid===this.uid){
				this.$router.push("/profile");
				console.log("userProfile(): this.userid === this.uid. \nuid: "+this.uid+"\nuserid: "+this.userid)
			} else {
				this.$router.push("/"+this.uid+"/profile");
				console.log("userProfile(): this.userid != this.uid. \nuid: "+this.uid+"\nuserid: "+this.userid);
			}
		},

	},

	created(){
		this.userid = localStorage.getItem('userid')
	},

}
</script>

<template>
  <div>
    <a href="javascript:" @mouseover="highlightProfile=true" @mouseout="highlightProfile=false"
        @click="userProfile()">
			<div class="d-flex align-items-center">
				<img class="imgThumbNail me-2 mb-2" :src="profpic" alt="Opps! error" />
				<h5 class="user-name ms-2 mb-2" :class="{ 'highlighted': highlightProfile }">{{ userprofname}}
				</h5>
			</div>
		</a>
  </div>
</template>