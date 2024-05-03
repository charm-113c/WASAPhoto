import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import PhotoContainer from './components/PhotoContainer.vue';
import PopUpMsg from './components/PopUpMsg.vue';
import Sidebar from './components/SidebarView.vue';

import './assets/main.css'
import './assets/stream.css'
import './assets/photoContainer.css'
import './assets/singleImage.css'
import './assets/commentSection.css'
import './assets/searchedProfile.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("PhotoContainer", PhotoContainer);
app.component("PopUpMsg", PopUpMsg);
app.component("Sidebar", Sidebar);
app.use(router)
app.mount('#app')
