<template>
    <div class="mod-tab" v-if="this.isActive">
        <slot></slot>
    </div>
</template>

<script>
export default {
  name: 'tab',
  props: {
    title: { required: true },
    selected: { default: false }
  },
  computed: {
    href () {
      return '#' + this.name.toLowerCase().replace(/ /g, '-')
    }
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
