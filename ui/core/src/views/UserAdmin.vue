<template>
<div class="container-fluid h-100">
  <div class="row h-100 app-main-holder">
    <sidebar-user :toggled="sidebarToggled" v-if="showSidebar"></sidebar-user>
    <main class="app-main col px-0" role="main">
      <app-view></app-view>
    </main>
  </div>
</div>
</template>

<script>
import SidebarUser from '@/views/appDependentComponents/SidebarUser'
import _ from 'lodash'

export default {
  name: 'admin',
  components: {
    SidebarUser
  },
  data () {
    return {
      sidebarToggled: false
    }
  },
  watch: {
    '$route': function (route) {
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
