import Vue from 'vue'
import App from './App.vue'
import router from './router/router'
import store from './store/store'
import i18n from './lang'
import './plugins/element.js'
import './directive/premissionBtn'
import './assets/css/public.css'
import './element-variables.scss'
import VueParticles from 'vue-particles'
import VueClipboard from 'vue-clipboard2'

import { messages } from './assets/js/common'
import gvgcount from './assets/js/gvgcount'

import { formatTimeToStr } from './assets/js/date'


// 引入字体文件
import '@/assets/icon/iconfont.css'
import '@/assets/icon/iconfont.js'


Vue.prototype.$url = process.env.VUE_APP_LOGOUT_URL


Vue.config.devtools = true;
Vue.use(VueParticles)
Vue.use(VueClipboard)
//全局挂载提示框 success/warning/info/error
Vue.prototype.$message = messages
Vue.prototype.$gvgcount = gvgcount

Vue.prototype.$formatTimeToStr = formatTimeToStr
Vue.config.productionTip = false
new Vue({
    i18n,
    router,
    store,
    render: h => h(App)
}).$mount('#app')