<template>
  <div class="container-fluid">
    <vue-headful :title="$t('Admin title', 'Proxeus - Admin')"/>
    <div class="row flex-nowrap" style="height:100%;min-width:900px;">
      <sidebar :toggled="sidebarToggled" v-if="showSidebar"></sidebar>
      <main class="app-main col px-0" role="main">
        <app-view></app-view>
      </main>
      <first-login-overlay
        keyz="admin"
        :preview-url="$t('First login admin', 'https://doc.proxeus.org/#/handbook')">
      </first-login-overlay>
    </div>
  </div>
</template>

<script>
import Sidebar from '@/views/appDependentComponents/Sidebar'
import FirstLoginOverlay from '@/views/FirstLoginOverlay'
import _ from 'lodash'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'admin',
  components: {
    Sidebar,
    FirstLoginOverlay
  },
  computed: {
    sidebarToggled () {
      return this.$route.meta.sidebarToggled === true
    },
    showSidebar () {
      return _.get(this, '$route.matched[0].props.default.showSidebar', true) === true
    }
  }
}
</script>
