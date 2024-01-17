// declare function require(name:string);

/*
 *
 * Vue
 *
 */
import Vue from 'vue'
import Web3 from 'web3' // eslint-disable-line no-unused-vars

/*
 *
 * Vendor native Vue Libs
 *
 */
import VueAppend from 'vue-append'
import VTooltip from 'v-tooltip'
import ErrorPage from 'vue-error-page'
import VueClipboard from 'vue-clipboard2'

/*
 *
 * My components
 *
 */
import App from './App.vue'
import router from './router/app'
import store from './store'

/*
 *
 * Global
 *
 */

import './global'
import './legacy-imports'
import VueScrollTo from 'vue-scrollto'
import { events } from 'vue-notification/src/events' // eslint-disable-line no-unused-vars

/*
 *
 * Vue component global registration. Use global components only if really needed
 * Normally it's better to import them in local components.
 *
 */
Vue.use(require('vue-moment'))
Vue.use(VueScrollTo)
Vue.use(VTooltip)
Vue.use(VueAppend)

Vue.use(VueClipboard)

// @ts-ignore
Vue.config.productionTip = false

/*
 *
 * We use a router wrapper to be able to trigger error pages
 * without redirecting
 *
 */
window.eventBus = new Vue()

// Vue.prototype.$notify = (params) => {
//   duration: 10000
//   // if (typeof params === 'string') {
//   //   params = { title: '', text: params }
//   // }
//   //
//   // if (typeof params === 'object') {
//   //   events.$emit('add', params)
//   // }
// }

Vue.use(ErrorPage, {
  resolver: (component) => {
    return require('./views/app/Errors/' + component).default
  }
})

const vm = new Vue({ // eslint-disable-line no-unused-vars
  el: '#app',
  router,
  store,
  components: {
    App
  },
  template: '<App/>'
})
