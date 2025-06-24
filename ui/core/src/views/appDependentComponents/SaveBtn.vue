<template>
<button v-if="item && app.amIWriteGrantedFor(item)" type="button" @click="click"
        class="btn btn-primary">
  Save
</button>
<button v-else type="button" :title="$t('read only')" disabled="disabled" class="btn btn-primary btn-round plus-btn">
  <i class="material-icons">remove_red_eye</i>
  <div style="position:relative;">
    <span :title="$t('read only')" class="read-only-entity">{{$t('read only')}}</span>
  </div>
</button>
</template>

<script>
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'save-btn',
  props: {
    click: {
      type: Function,
      default: function () {}
    },
    item: {
      type: Object,
      default: null
    }
  }
}
</script>

<style lang="scss">
  @use "../../assets/styles/variables";
  @use "~bootstrap/scss/mixins";
  @use "../../assets/styles/sidebar.scss";

  .brand-name {
    letter-spacing: 2px;
  }

  .read-only-entity {
    top: -65px;
    right: 32px;
    position: absolute;
    white-space: nowrap;
    color: $white;
    padding: 8px;
    background: $primary;
    border-radius: 10px;
    font-size: 16px;
  }
</style>
