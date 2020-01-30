import Vue from 'vue'
import Router from 'vue-router'

import Admin from '../views/Admin.vue'
import Dashboard from '../views/Dashboard.vue'
import I18n from '../views/I18n.vue'
import Settings from '../views/Settings.vue'
import FormBuilder from '../views/FormBuilder.vue'
import Forms from '../views/Forms.vue'
import User from '../views/User.vue'
import Users from '../views/Users.vue'
import Support from '../views/Support.vue'
import Templates from '../views/Templates.vue'
import Template from '../views/Template.vue'
import Workflow from '../views/Workflow.vue'
import Workflows from '../views/Workflows.vue'
import External from '../views/External.vue'
import NotFound from '../views/app/Errors/NotFound.vue'
import AdminImportExport from '../views/AdminImportExport'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/admin',
      alias: '/p',
      name: 'Admin',
      redirect: '/admin/workflow',
      component: Admin,
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: Dashboard
        }, {
          path: 'form/:id',
          name: 'FormBuilder',
          component: FormBuilder,
          props: { sidebarToggled: true }
        }, {
          path: 'form',
          name: 'Forms',
          component: Forms
        }, {
          path: 'i18n',
          name: 'I18n',
          component: I18n,
          props: { sidebarToggled: true }
        }, {
          path: 'settings',
          name: 'Settings',
          component: Settings,
          props: { sidebarToggled: true }
        }, {
          path: 'data',
          name: 'AdminImportExport',
          component: AdminImportExport,
          props: { sidebarToggled: true }
        }, {
          path: 'user/:id',
          name: 'User',
          component: User,
          props: true
        }, {
          path: 'user',
          name: 'Users',
          component: Users
        }, {
          path: 'template/:id',
          name: 'Template',
          component: Template,
          props: { sidebarToggled: true }
        }, {
          path: 'support',
          name: 'Support',
          component: Support,
          props: { sidebarToggled: false }
        }, {
          path: 'template',
          name: 'Templates',
          component: Templates,
          props: { sidebarToggled: false }
        }, {
          path: 'workflow',
          name: 'Workflows',
          component: Workflows
        }, {
          path: 'workflow/:id',
          name: 'Workflow',
          component: Workflow,
          props: { sidebarToggled: true }
        }, {
          path: 'externalNode/:name/:id',
          name: 'External configuration',
          component: External,
          props: { sidebarToggled: true }
        }
      ]
    }, {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }
  ]
})
