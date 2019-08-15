import Vue from 'vue'
import Router from 'vue-router'

import UserAdmin from '../views/UserAdmin.vue'
import DocumentViewer from '../views/DocumentViewer.vue'
import DocumentVerification from '../views/DocumentVerification.vue'
import SignatureRequests from '../views/SignatureRequests.vue'
import ManagementList from '../views/ManagementList.vue'
import DocumentFlow from '../views/DocumentFlow.vue'
import DocumentPayment from '../views/DocumentPayment.vue'
import DocumentCreate from '../views/DocumentCreate.vue'
import Documents from '../views/Documents.vue'
import UserImportExport from '../views/UserImportExport.vue'
import NotFound from '../views/user/Errors/NotFound.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/document/:id',
      name: 'DocumentFlow',
      component: DocumentFlow
    }, {
      path: '/user',
      name: 'UserAdmin',
      component: UserAdmin,
      redirect: { name: 'Documents' },
      children: [
        {
          path: '/',
          redirect: { name: 'Documents' }
        },
        {
          path: 'data',
          name: 'UserImportExport',
          component: UserImportExport
        },
        {
          path: 'document',
          name: 'Documents',
          component: Documents
        },
        {
          path: '/document/:documentId/payment',
          name: 'DocumentPayment',
          component: DocumentPayment,
          props: true
        },
        {
          path: 'document-verify',
          name: 'DocumentVerification',
          component: DocumentVerification
        }, {
          path: 'signature-requests',
          name: 'SignatureRequests',
          component: SignatureRequests
        }, {
          path: 'management-list',
          name: 'ManagementList',
          component: ManagementList
        }, {
          path: 'document/create',
          name: 'document-create',
          component: DocumentCreate
        }, {
          path: 'document/:id',
          name: 'DocumentViewer',
          component: DocumentViewer,
          props: true
        }, {
          path: '*',
          name: 'NotFound',
          component: NotFound
        }
      ]
    },
    {
      path: '*',
      name: 'NotFound',
      component: NotFound
    }]
})
