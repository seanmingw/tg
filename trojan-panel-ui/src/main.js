import Vue from 'vue'
import 'normalize.css/normalize.css'
import ElementUI from 'element-ui'
import '@/styles/element-variables.scss'
import 'element-ui/lib/theme-chalk/index.css'
import '@/styles/index.scss'
import App from './App'
import store from '@/store'
import router from '@/router'
import Cookies from 'js-cookie'
import '@/icons'
import '@/permission'
import i18n from '@/lang'

import VCA from '@vue/composition-api'
import JsonEditorVue from 'json-editor-vue'

Vue.use(VCA)
Vue.use(JsonEditorVue)

Vue.use(ElementUI, {
  size: Cookies.get('size') || 'small', // set element-ui default size
  i18n: (key, value) => i18n.t(key, value)
})

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  store,
  i18n,
  render: (h) => h(App)
})
