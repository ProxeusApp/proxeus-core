import Vue from 'vue'
import Router from 'vue-router'

import Home from '../views/Home.vue'
import Validation from '../views/Validation.vue'
// import Login from '../views/Login.vue'
import PasswordResetRequest from '../views/PasswordResetRequest.vue'
import PasswordReset from '../views/PasswordReset.vue'
import RegisterRequest from '../views/RegisterRequest.vue'
import Register from '../views/Register.vue'
import EmailChangeRequest from '../views/EmailChangeRequest.vue'
import EmailChange from '../views/EmailChange.vue'
import AdminLogin from '../views/AdminLogin.vue'
import RegistrationSuccess from '../views/frontend/RegistrationSuccess.vue'
import NotFound from '../views/frontend/Errors/NotFound.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    }, {
      path: '/validation',
      name: 'Validation',
      component: Validation,
      props: { showHeader: false, showFooter: false }
    }, {
      path: '/change/email',
      name: 'EmailChangeRequest',
      component: EmailChangeRequest
    }, {
      path: '/change/email/:token',
      name: 'EmailChange',
      component: EmailChange
    }, {
      path: '/register',
      name: 'RegisterRequest',
      component: RegisterRequest
    }, {
      path: '/register/:token',
      name: 'Register',
      component: Register
    }, {
      path: '/reset/password',
      name: 'PasswordResetRequest',
      component: PasswordResetRequest
    }, {
      path: '/reset/password/:token',
      name: 'PasswordReset',
      component: PasswordReset
    }, {
      path: '/login',
      name: 'AdminLogin',
      component: AdminLogin
    },
    {
      path: '/guest',
      name: 'RegistrationSuccess',
      component: RegistrationSuccess
    }, {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }]
})
