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

/*
 *
 * My components
 *
 */
import DocumentApp from './DocumentApp.vue'
import router from './router/user'
import store from './store'

/*
 *
 * Global
 *
 */
import 'bootstrap'
import './global'
import './legacy-user-imports'

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

// @ts-ignore
Vue.config.productionTip = false

/*
 *
 * We use a router wrapper to be able to trigger error pages
 * without redirecting
 *
 */
// @ts-ignore
window.eventBus = new Vue()

Vue.use(ErrorPage, {
  resolver: (component) => {
    return require('./views/user/Errors/' + component).default
  }
})

const vm = new Vue({ // eslint-disable-line no-unused-vars
  el: '#document-app',
  router,
  store,
  components: { DocumentApp },
  template: '<DocumentApp/>'
})
