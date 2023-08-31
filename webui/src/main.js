import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import LinkToUserProfile from './components/LinkToUserProfile.vue'
import NoLinkToUserProfile  from './components/NoLinkToUserProfile.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("LinkToUserProfile", LinkToUserProfile);
app.component("NoLinkToUserProfile", NoLinkToUserProfile);
app.use(router)
app.mount('#app')
