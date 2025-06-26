<template>
  <div style="height: 100%">
    <notifications
      group="app"
      classes="alert"
      position="top center"
      :duration="6500"
      :width="400"
      style="z-index: 500000000000"
    />
    <router-view></router-view>
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
import baseApp from "./baseApp";
import Blocker from "./components/Blocker";

export default {
  components: { Blocker },
  mixins: [baseApp],
  name: "App",
  data() {
    return {
      moreInApp: "yes...",
      blockUI: null,
      unblockUI: null,
      connectionLostBlockUI: null,
      connectionLostUnblockUI: null,
    };
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
  created() {
    window.$root = this.$root;
    this.$root.$on("service-on", this.onServiceOn);
    this.$root.$on("service-off", this.onServiceOff);

    if (!this.app.checkUserHasSession()) {
      this.app.redirectToLogin(window.location.pathname);
    }
  },
  beforeDestroy() {
    this.$root.$off("service-on", this.onServiceOn);
    this.$root.$off("service-off", this.onServiceOff);
  },
};
</script>

<style lang="scss">
@import "@/assets/styles/app.scss";
</style>
