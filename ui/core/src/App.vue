<template>
<div style="height:100%;">
  <notifications group="app" classes="alert" position="top center" :duration="6500" :width="400" style="z-index: 500000000000;"/>
  <router-view></router-view>
  <blocker :text1="$t('Common blocker text 1','JUST A MOMENT')" :text2="$t('Common blocker text 2','PROCESSING')"
           :text3="$t('Common blocker text 3','PROCESSING')" :setup="commonUIBlocker"/>
  <blocker :text1="$t('Reconnecting blocker text 1','SORRY')" :text2="$t('Reconnecting blocker text 2','JUST A MOMENT')"
           :text3="$t('Reconnecting blocker text 3','CONNECTING')" :showAnimation="true" :setup="setupUIBlocker"/>
</div>
</template>

<script>
import baseApp from './baseApp'
import Blocker from './components/Blocker'

export default {
  components: { Blocker },
  mixins: [baseApp],
  name: 'App',
  data () {
    return {
      moreInApp: 'yes...',
      blockUI: null,
      unblockUI: null,
      connectionLostBlockUI: null,
      connectionLostUnblockUI: null
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
    }
  },
  created () {
    window.$root = this.$root
    this.$root.$on('service-on', this.onServiceOn)
    this.$root.$on('service-off', this.onServiceOff)
  },
  beforeDestroy () {
    this.$root.$off('service-on', this.onServiceOn)
    this.$root.$off('service-off', this.onServiceOff)
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

  .navbar h1 {
    margin-bottom: 0;
  }

  .app-main {
    height: auto;

    @media (max-width: 767px) {
      max-width: 100% !important;
    }
  }

  ::-moz-selection {
    background: $info;
    color: $primary;
  }

  ::selection {
    background: $info;
    color: $primary;
  }

  .btn.btn-sm.topnav-back {
    border: 0;
    border-right: 1px solid $gray-300;
    border-radius: 0;
    margin-left: -24px;
    padding-top: 0;
    padding-bottom: 0;
    height: 58px;
    vertical-align: middle;
  }

  .navbar.navbar-expand-lg {
    z-index: 100;
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

  .toggle-row {
    z-index: 2223273;
    width: 100%;
    background: white;
    border: 1px solid #gray-200;
    box-shadow: 0 45px 100px rgba(0, 0, 0, .35);
    position: relative;
  }

  .toggle-row-backdrop {
    background: rgba(0, 0, 0, .2);
    position: fixed;
    z-index: 98;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
  }

</style>
