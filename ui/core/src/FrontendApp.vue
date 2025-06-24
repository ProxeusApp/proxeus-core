<template>
  <div id="frontend-app" :class="{ 'frontend-headless': showHeader === false }">
    <header class="mb-4" v-if="showHeader === true">
      <frontend-navbar></frontend-navbar>
    </header>
    <main class="frontend-main" role="main">
      <notifications
        group="app"
        classes="alert"
        position="top center"
        :duration="4000"
        style="z-index: 500000000"
      />
      <frontend-view></frontend-view>
    </main>
    <div class="container-fluid" v-if="showFooter === true">
      <div class="row">
        <div class="footer p-3 col-sm-12 d-flex justify-content-between">
          <div></div>
          <div>
            <a
              target="_blank"
              href="https://doc.proxeus.org/#/handbook"
              class="footer-link"
              >Handbook
            </a>
            <a
              target="_blank"
              :href="$t('Privacy Policy url', '')"
              class="footer-link ml-3"
              >Privacy Policy</a
            >
            <a
              target="_blank"
              v-bind:href="$t('Terms & Conditions link')"
              class="footer-link ml-3"
              >Terms &amp; Conditions</a
            >
            <a
              target="_blank"
              href="https://github.com/ProxeusApp/proxeus-core/blob/master/LICENSE"
              class="footer-link ml-3"
              >License</a
            >
            <a
              target="_blank"
              href="https://github.com/ProxeusApp/proxeus-core"
              class="footer-link ml-3"
              >Source Code</a
            >
          </div>
          <div>
            <a
              target="_blank"
              href="https://proxeus.org/"
              class="footer-link ml-3 small"
            >
              {{ $t("Powered by") }}
              <img src="/static/proxeus-logo-white.svg" alt="" />
            </a>
          </div>
        </div>
      </div>
    </div>
    <blocker
      :text1="$t('Common blocker text 1', 'JUST A MOMENT')"
      :text2="$t('Common blocker text 2', 'PROCESSING')"
      :text3="$t('Common blocker text 3', 'PROCESSING')"
      :setup="commonUIBlocker"
    />
    <blocker
      :text1="$t('Reconnecting blocker text 1', 'SORRY')"
      :text2="$t('Reconnecting blocker text 2', 'JUST A MOMENT')"
      :text3="$t('Reconnecting blocker text 3', 'CONNECTING')"
      :showAnimation="true"
      :setup="setupUIBlocker"
    />
  </div>
</template>

<script>
import FrontendNavbar from "@/components/frontend/FrontendNavbar";
import _ from "lodash";
import baseApp from "./baseApp";
import Blocker from "./components/Blocker";

export default {
  mixins: [baseApp],
  name: "Frontend",
  components: {
    Blocker,
    FrontendNavbar,
  },
  data() {
    return {
      blockUI: null,
      unblockUI: null,
      connectionLostBlockUI: null,
      connectionLostUnblockUI: null,
    };
  },
  created() {
    this.$root.$on("service-on", this.onServiceOn);
    this.$root.$on("service-off", this.onServiceOff);
  },
  beforeDestroy() {
    this.$root.$off("service-on", this.onServiceOn);
    this.$root.$off("service-off", this.onServiceOff);
  },
  methods: {
    commonUIBlocker(blockClb, unblockClb) {
      this.blockUI = blockClb;
      this.unblockUI = unblockClb;
    },
    setupUIBlocker(blockClb, unblockClb) {
      this.connectionLostBlockUI = blockClb;
      this.connectionLostUnblockUI = unblockClb;
    },
    onServiceOn() {
      if (this.connectionLostUnblockUI) {
        this.connectionLostUnblockUI();
      }
    },
    onServiceOff() {
      if (this.connectionLostBlockUI) {
        this.connectionLostBlockUI();
      }
    },
  },
  computed: {
    showHeader() {
      return (
        _.get(this, "$route.matched[0].props.default.showHeader", true) === true
      );
    },
    showFooter() {
      return (
        _.get(this, "$route.matched[0].props.default.showFooter", true) === true
      );
    },
  },
};
</script>

<style lang="scss">
@use "assets/styles/variables.scss";
@use "assets/styles/fonts.scss";
@use "assets/styles/buttons.scss";

@use "~@mdi/font/scss/materialdesignicons.scss";
@use "assets/styles/global.scss";

$mdi-font-path: "~@mdi/font/fonts";

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
  div {
    min-width: 200px;
  }
}

.footer-link {
  color: white;
  &:hover {
    color: $gray-300;
  }
  img {
    height: 20px;
    margin-left: 5px;
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
