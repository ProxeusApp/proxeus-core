/*
 *
 * Vue
 *
 */
import Vue from 'vue'
import Web3 from 'web3'
import ErrorPage from 'vue-error-page'

import FTG from './libs/legacy/global.js'

/*
 *
 * My components
 *
 */
import FrontendApp from './FrontendApp.vue'
import store from './store'

/*
 *
 * Routing
 *
 */
import router from './router/frontend'

/*
 *
 * i18n
 *
 */
import './global'

// @ts-ignore
window.$ = window.jQuery = require('jquery')

/*
 *
 * We use a router wrapper to be able to trigger error pages
 * without redirecting
 *
 */
window.eventBus = new Vue()

Vue.use(ErrorPage, {
  resolver: (component) => {
    return require('./views/frontend/Errors/' + component).default
  },
  tagName: 'frontend-view'
})

new Vue({
  el: '#frontend-app',
  router,
  store,
  components: { FrontendApp },
  template: '<FrontendApp/>'
})
