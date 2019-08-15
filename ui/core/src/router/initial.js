import Vue from 'vue'
import Router from 'vue-router'

import Settings from '../views/appDependentComponents/SettingsInner.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '*',
      name: 'Settings',
      component: Settings,
      props: { configOnly: false }
    }]
})
