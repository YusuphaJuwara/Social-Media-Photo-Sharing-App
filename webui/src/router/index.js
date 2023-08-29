import {createRouter, createWebHashHistory} from 'vue-router'
import Profile from '../views/Profile.vue'
import ProfileOtherUsers from '../views/ProfileOtherUsers.vue'
import StreamPosts from '../views/StreamPosts.vue'
import PostsNotBanned from '../views/PostsNotBanned.vue'
import AllUsersNotBanned from '../views/AllUsersNotBanned.vue'
import Search from '../views/Search.vue'
import CreatePost from '../views/CreatePost.vue'
import Login from '../views/Login.vue'

import ModalA from '../views/ModalA.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/profile', component: Profile},
		{path: '/:uid/profile', component: ProfileOtherUsers},
		{path: '/stream', component: StreamPosts},
		{path: '/posts', component: PostsNotBanned},
		{path: '/users', component: AllUsersNotBanned},
		{path: '/search', component: Search},
		{path: '/create-post', component: CreatePost},
		{path: '/', component: Login},
		
		{path: '/modala', component: ModalA},
	]
})

export default router
