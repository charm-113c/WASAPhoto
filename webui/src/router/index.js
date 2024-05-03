import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import StreamView from '../views/StreamView.vue'
import MyProfile from '../views/MyProfile.vue'
import SearchedProfile from '../views/SearchedProfile.vue'
import ImageView from '../views/ImageView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		// {path: '/login', name: 'doLogin', component: LoginView},
		{path: '/users/:username/stream', name: 'stream', component: StreamView},
		{path: '/users/:username/search/:user2', name: 'searchProfile', component: SearchedProfile, props: true},
		{path: '/users/:username/profile', name: 'myProfile', component: MyProfile},
		{path: '/users/:uploader/photos/:photoID', name: 'viewImage', component: ImageView, props: true}
	]
})

// Navigation guard to check for route changes
router.beforeEach((to, from, next) => {
	// Check if the user is navigating to the login page
	if ((to.path === '/') || (to.path === '/login')) {
	  // Log out the user by removing the bearer token from sessionStorage
	  sessionStorage.removeItem('username')
	  sessionStorage.removeItem('bearerToken')
	}
	// Continue with the navigation
	next()
  })

export default router
