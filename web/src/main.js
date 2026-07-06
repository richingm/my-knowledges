import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import '@blocknote/core/style.css'
import '@blocknote/core/fonts/inter.css'

const app = createApp(App)
app.use(router)
app.mount('#app')
