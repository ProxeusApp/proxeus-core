<template>
    <div class="mod-tab" v-show="this.isActive">
        <slot></slot>
    </div>
</template>

<script>
export default {
  name: 'tab',
  props: {
    title: { required: true },
    selected: { default: false },
    onSelect: { type: Function }
  },
  computed: {
    href () {
      return '#' + this.name.toLowerCase().replace(/ /g, '-')
    }
  },
  updated () {
  },
  data () {
    return {
      isActive: false
    }
  },
  mounted () {
    this.isActive = this.selected
  }
}
</script>
