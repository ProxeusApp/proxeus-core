<template>
<div class="field-parent mb-3" v-bind="$attrs">
  <div class="form-group mb-0">
    <label>{{label}}</label>
    <input :maxlength="max" :autofocus="autofocus" :disabled="disabled" @input="change" :name="name" ref="field"
           class="form-control" :class="{'has-content':hasContent}" :type="type" placeholder="">
    <span class="focus-border"></span>
    <a v-if="action.Func && action.Name" class="faction" href="JavaScript:void(0);"
       @click="action.Func">{{action.Name}}
    </a>
    <a v-if="action.Link && action.Name" class="faction" :target="action.Target" :href="action.Link">{{action.Name}}</a>
  </div>
  <div class="helptext text-muted" v-if="helptext"><small>{{ helptext }}</small></div>
  <div class="errors"></div>
  <slot></slot>
</div>
</template>

<script>
export default {
  name: 'animated-input',
  props: {
    name: {
      type: String,
      default: null
    },
    type: {
      type: String,
      default: 'text'
    },
    label: {
      type: String,
      default: 'Undefined label'
    },
    value: {
      type: null,
      default: null
    },
    action: {
      type: Object,
      default: () => { return { Name: '', Func: null } }
    },
    disabled: {
      type: Boolean,
      default: false
    },
    autofocus: {
      type: Boolean,
      default: false
    },
    max: {
      type: Number,
      default: 1000
    },
    input: {
      type: Function
    },
    helptext: {}
  },
  data () {
    return {
      initDone: false,
      hasContent: false
    }
  },
  watch: {
    'value': 'valueChanged'
  },
  mounted () {
    if (this.$refs.field && this.value) {
      this.$refs.field.value = this.value
      this.hasContent = true
    }
    this.initDone = true
  },
  methods: {
    valueChanged () {
      if (this.$refs.field) {
        if (this.$refs.field.value !== this.value) {
          this.$refs.field.value = this.value
        }
      }
      if (this.value) {
        this.hasContent = true
      } else {
        this.hasContent = false
      }
    },
    change (e) {
      if (!this.initDone) {
        return
      }
      if (e.target) {
        if (e.target.value) {
          this.hasContent = true
        } else {
          this.hasContent = false
        }
        this.$emit('input', e.target.value)
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
