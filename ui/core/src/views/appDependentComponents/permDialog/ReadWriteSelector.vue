<template>
<simple-select :unselect="unselect" class="read-write-selector" v-model="valueProxy" :selected="selected"
               :disableWithVisibleLayer="disableWithVisibleLayer" :disabled="disabled"
               :onSelect="onSelect" :onSelectedChange="onSelectedChange" :options="options">
  <template slot-scope="{option}">
  <div style="padding: 2px;">
    <i v-if="option.id==2" class="material-icons">create</i>
    <i v-else-if="option.id==1" class="material-icons">visibility</i>
    <i v-else-if="option.id==0" class="material-icons">block</i>
    <span class="my-explanation" style="margin-left: 4px;">{{option.label}}</span>
  </div>
  </template>
</simple-select>
</template>

<script>
import SimpleSelect from '@/components/SimpleSelect'

export default {
  name: 'read-write-selector',
  components: {
    SimpleSelect
  },
  props: {
    value: {
      type: null,
      default: null
    },
    selected: {
      type: null,
      default: null
    },
    provideNone: {
      type: Boolean,
      default: true
    },
    unselect: {
      type: Boolean,
      default: true
    },
    disabled: {
      type: Boolean,
      default: false
    },
    disableWithVisibleLayer: {
      type: Boolean,
      default: () => true
    },
    onSelect: {
      type: Function,
      default: () => null
    },
    onSelectedChange: {
      type: Function,
      default: () => null
    }
  },
  data () {
    let $t = this.$t
    return {
      options: [{ id: 0, label: $t('None') }, { id: 1, label: $t('Can read') }, { id: 2, label: $t('Can write') }],
      valueProxy: null
    }
  },
  watch: {
    'valueProxy': 'updateSelected'
  },
  methods: {
    updateSelected () {
      this.$emit('input', this.valueProxy)
    }
  },
  created () {
    this.valueProxy = this.value
    if (this.provideNone === false) {
      this.options.splice(0, 1)
    }
  }
}
</script>

<style>
  .read-write-selector {
    min-width: 80px;
  }
</style>
