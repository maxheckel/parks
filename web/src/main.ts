import './assets/main.css'

import {createApp, createSSRApp} from 'vue'
import App from './App.vue'
import router from './router'

const app = createSSRApp(App)

app.use(router)

app.mount('#app')
