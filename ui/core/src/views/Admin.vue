<template>
<div class="container-fluid" style="height:100%;">
  <vue-headful :title="$t('Admin title', 'Proxeus - Admin')"/>
  <div class="row" style="height:100%;">
    <sidebar :toggled="sidebarToggled" v-if="showSidebar"></sidebar>
    <main class="app-main col px-0" role="main">
      <app-view></app-view>
    </main>
    <first-login-overlay></first-login-overlay>
  </div>
</div>
</template>

<script>
import Sidebar from '@/views/appDependentComponents/Sidebar'
import FirstLoginOverlay from '@/views/FirstLoginOverlay'
import _ from 'lodash'

export default {
  name: 'admin',
  components: {
    Sidebar,
    FirstLoginOverlay
  },
  data () {
    return {
      sidebarToggled: false
    }
  },
  watch: {
    '$route': function () {
      this.toggleSidebar()
    }
  },
  computed: {
    showSidebar () {
      return _.get(this, '$route.matched[0].props.default.showSidebar', true) === true
    }
  },
  methods: {
    toggleSidebar () {
      this.sidebarToggled = _.get(this, '$route.matched[1].props.default.sidebarToggled', false) === true
    }
  },
  mounted () {
    this.toggleSidebar()
  }
}
</script>
<style>
  .app-main {
    overflow-x: auto;
  }
</style>
