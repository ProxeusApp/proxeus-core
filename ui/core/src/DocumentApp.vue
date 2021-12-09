<template>
<div id="document-app">
  <notifications group="app" classes="alert" position="top center" :duration="4000" style="z-index: 500000000;"/>
  <router-view></router-view>
  <blocker :text1="$t('Common blocker text 1','JUST A MOMENT')" :text2="$t('Common blocker text 2','PROCESSING')"
           :text3="$t('Common blocker text 3','PROCESSING')" :setup="commonUIBlocker"/>
  <blocker :text1="$t('Reconnecting blocker text 1','SORRY')" :text2="$t('Reconnecting blocker text 2','JUST A MOMENT')"
           :text3="$t('Reconnecting blocker text 3','CONNECTING')" :showAnimation="true" :setup="setupUIBlocker"/>
</div>
</template>

<script>
// import SidebarUser from '@/views/appDependentComponents/SidebarUser'
import baseApp from './baseApp'
import Blocker from './components/Blocker'

export default {
  mixins: [baseApp],
  name: 'document-app',
  components: {
    Blocker
    // SidebarUser
  },
  created () {
    this.$root.$on('service-on', this.onServiceOn)
    this.$root.$on('service-off', this.onServiceOff)
  },
  beforeDestroy () {
    this.$root.$off('service-on', this.onServiceOn)
    this.$root.$off('service-off', this.onServiceOff)
  },
  data () {
    return {
      sidebarToggled: false,
      blockUI: null,
      unblockUI: null,
      connectionLostBlockUI: null,
      connectionLostUnblockUI: null
    }
  },
  watch: {
    $route: function (route) {
      this.toggleSidebar()
    }
  },
  computed: {
    showSidebar () {
      return _.get(this, '$route.matched[0].props.default.showSidebar', true) === true
    }
  },
  methods: {
    commonUIBlocker (blockClb, unblockClb) {
      this.blockUI = blockClb
      this.unblockUI = unblockClb
    },
    setupUIBlocker (blockClb, unblockClb) {
      this.connectionLostBlockUI = blockClb
      this.connectionLostUnblockUI = unblockClb
    },
    onServiceOn () {
      if (this.connectionLostUnblockUI) {
        this.connectionLostUnblockUI()
      }
    },
    onServiceOff () {
      if (this.connectionLostBlockUI) {
        this.connectionLostBlockUI()
      }
    },
    toggleSidebar () {
      this.sidebarToggled = _.get(this, '$route.matched[1].props.default.sidebarToggled', false) === true
    }
  },
  mounted () {
    this.toggleSidebar()
  }
}
</script>

<style lang="scss">
  @import "~bootstrap/scss/functions";
  @import "assets/styles/variables.scss";
  @import "~bootstrap/scss/bootstrap";
  @import "assets/styles/fonts.scss";
  @import "assets/styles/buttons.scss";

  $mdi-font-path: "~@mdi/font/fonts";
  @import "~@mdi/font/scss/materialdesignicons.scss";

  @import "assets/styles/modals.scss";
  @import "assets/styles/fancy-radio-checkbox.scss";

  @import "assets/styles/forms.scss";
  @import "assets/styles/alerts.scss";
  @import "assets/styles/global.scss";

  @import "assets/styles/flatpickr.scss";
  @import "assets/styles/nav-tabs.scss";

  ::-moz-selection {
    background: $info;
    color: $primary;
  }

  ::selection {
    background: $info;
    color: $primary;
  }

  .navbar h1 {
    margin-bottom: 0;
  }

  .app-main {
    @media (max-width: 767px) {
      max-width: 100% !important;
    }
  }

  .app-main {
    //min-width: 900px;
    > div > .container-fluid, > div > .main-container {
      height: calc(100vh - 62px);
      overflow-y: scroll;
    }
    .topnav {
      box-shadow: 0 0 7px rgba(0,0,0,.1);
    }
  }

  #document-app .sidebar {
     position: relative;
     z-index: 999999999 !important;
   }

  .btn.btn-sm.topnav-back {
    border: 0;
    border-right: 1px solid $gray-300;
    border-radius: 0;
    margin-left: -24px;
    padding-top: 0;
    padding-bottom: 0;
    background: transparent;
    height: 60px;
    vertical-align: middle;
    display: flex;
    align-items: center;
  }

  .navbar h1.navbar-text {
    display: inline-block;
    padding-top: 0.75rem;
    padding-bottom: 0.75rem;
    margin-top: 0;
    margin-bottom: 0;
    overflow: hidden;
    min-width: 0;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
