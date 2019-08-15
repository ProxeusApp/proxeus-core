// declare function require(name:string);

/*
 *
 * Vue
 *
 */
import Vue from 'vue'
import Web3 from 'web3'

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
import InitialApp from './InitialApp.vue'
import router from './router/initial'

/*
 *
 * Global
 *
 */

import './global'

import './legacy-imports'

import VueScrollTo from 'vue-scrollto'

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

Vue.use(ErrorPage, {
  resolver: (component) => {
    return require('./views/app/Errors/' + component).default
  }
})

new Vue({
  el: '#app',
  router,
  components: { InitialApp },
  template: '<InitialApp/>'
})
