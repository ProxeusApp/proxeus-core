<template>
<div :title="title" class="field-parent" style="margin-top: 15px;">
  <div class="fancy-checkbox" :class="{disabled:disabled}">
    <input v-model="changing" :disabled="disabled" type="checkbox" :id="id">
    <label :for="id"><span><i></i></span>{{label}}</label>
  </div>
</div>
</template>

<script>
export default {
  name: 'checkbox',
  props: {
    title: {
      type: String,
      default: ''
    },
    label: {
      type: String,
      default: 'Undefined label'
    },
    value: {
      type: null,
      default: null
    },
    disabled: {
      type: Boolean,
      default: false
    },
    input: {
      type: Function
    }
  },
  data () {
    return {
      id: null,
      me: {},
      profile: { Name: '', Detail: '', hello: '' }
    }
  },
  created () {
    this.id = this.uuid()
  },
  methods: {
    uuid () {
      const S4 = function () {
        return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1)
      }
      return (S4() + S4() + '-' + S4() + '-' + S4() + '-' + S4() + '-' + S4() + S4() + S4())
    }
  },
  computed: {
    changing: {
      get () {
        return this.value
      },
      set (a) {
        this.$emit('input', a)
        if (this.input) {
          this.input()
        }
      }
    }
  }
}
</script>

<style>

</style>
