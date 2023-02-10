import {createRouter, createWebHashHistory} from 'vue-router'
import Profile from '../views/Profile.vue'
import StreamPosts from '../views/StreamPosts.vue'
import PostsNotBanned from '../views/PostsNotBanned.vue'
import AllUsersNotBanned from '../views/AllUsersNotBanned.vue'
import Search from '../views/Search.vue'
import Post from '../views/Post.vue'
import CreatePost from '../views/CreatePost.vue'
import Login from '../views/Login.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/:userid/profile', component: Profile},
		{path: '/:userid/stream/', component: StreamPosts},
		{path: '/:userid/posts/', component: PostsNotBanned},
		{path: '/:userid/users/', component: AllUsersNotBanned},

		{path: '/:postid/post', component: Post},
		{path: '/create-post', component: CreatePost},
		{path: '/login', component: Login},
		{path: '/search', component: Search},

	]
})

export default router
