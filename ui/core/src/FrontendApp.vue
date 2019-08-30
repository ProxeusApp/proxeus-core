<template>
<div id="frontend-app" :class="{'frontend-headless' : showHeader === false}">
  <header class="mb-4" v-if="showHeader === true">
    <frontend-navbar></frontend-navbar>
  </header>
  <main class="frontend-main" role="main">
    <notifications group="app" classes="alert" position="top center" :duration="4000" style="z-index: 500000000;"/>
    <frontend-view></frontend-view>
  </main>
  <div class="container-fluid" v-if="showFooter === true">
    <div class="row">
      <div class="footer p-3 col-sm-12 d-flex justify-content-center">
        <a target="_blank"
           href="https://docs.google.com/document/d/1C3B1oNY6lOv8Q_AvbKhwlySrS6qTiRl3raPLV6OXr7w/preview"
           class="footer-link">Handbook
        </a>
        <a target="_blank" href="https://proxeus.com/en/privacy-policy/" class="footer-link ml-3">Privacy Policy</a>
      </div>
    </div>
  </div>
  <blocker :text1="$t('Common blocker text 1','JUST A MOMENT')" :text2="$t('Common blocker text 2','PROCESSING')"
           :text3="$t('Common blocker text 3','PROCESSING')" :setup="commonUIBlocker"/>
  <blocker :text1="$t('Reconnecting blocker text 1','SORRY')" :text2="$t('Reconnecting blocker text 2','JUST A MOMENT')"
           :text3="$t('Reconnecting blocker text 3','CONNECTING')" :showAnimation="true" :setup="setupUIBlocker"/>
</div>
</template>

<script>
import FrontendNavbar from '@/components/frontend/FrontendNavbar'
import _ from 'lodash'
import baseApp from './baseApp'
import Blocker from './components/Blocker'

export default {
  mixins: [baseApp],
  name: 'Frontend',
  components: {
    Blocker,
    FrontendNavbar
  },
  data () {
    return {
      blockUI: null,
      unblockUI: null,
      connectionLostBlockUI: null,
      connectionLostUnblockUI: null
    }
  },
  created () {
    this.$root.$on('service-on', this.onServiceOn)
    this.$root.$on('service-off', this.onServiceOff)
  },
  beforeDestroy () {
    this.$root.$off('service-on', this.onServiceOn)
    this.$root.$off('service-off', this.onServiceOff)
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
  computed: {
    showHeader () {
      return _.get(this, '$route.matched[0].props.default.showHeader', true) === true
    },
    showFooter () {
      return _.get(this, '$route.matched[0].props.default.showFooter', true) === true
    }
  }
}
</script>

<style lang="scss">
  @import "assets/styles/variables.scss";
  @import "~bootstrap/scss/bootstrap";
  @import "assets/styles/fonts.scss";
  @import "assets/styles/buttons.scss";

  $mdi-font-path: "~@mdi/font/fonts";
  @import "~@mdi/font/scss/materialdesignicons.scss";
  @import "assets/styles/global.scss";

  html {
    position: relative;
    min-height: 100%;
  }

  body {
    padding-bottom: 100px;
  }

  h1 {
    font-size: 48px;
  }

  ::-moz-selection {
    background: $info;
    color: $primary;
  }

  ::selection {
    background: $info;
    color: $primary;
  }

  .footer {
    bottom: 0;
    position: fixed;
    background: $primary;
  }

  .footer-link {
    color: white;
    &:hover {
      color: $gray-300;
    }
  }

  .text-hint {
    color: $text-hint;
  }

  .border-2 {
    border-width: 2px !important;
  }

  .btn-lighter {
    border: 2px solid $gray-200;
    background: none;
    color: $primary;
    border-radius: 25px;
    &:hover {
      background: $gray-100;
      border-color: $gray-100;
    }
  }

  .btn-light-round {
    @extend .btn-light;
    border-radius: 25px;
  }

  .smaller {
    font-size: 0.7rem;
  }

</style>
