<template>
<div class="col-sm-12">
  <div class="row">
    <div class="col-md-4">
      <animated-input :input="input" label="Name" v-model="value.name" :max="maxNameLength"/>
    </div>
    <div class="col-md-4">
      <animated-input :input="input" label="Detail" v-model="value.detail" :max="maxDetailLength"/>
    </div>
    <div v-if="displayPrice === true" class="col-md-4">
      <animated-input :input="input" label="Price (XES)" v-model.number="value.price"/>
    </div>
  </div>
</div>
</template>

<script>
import AnimatedInput from './AnimatedInput'

export default {
  components: { AnimatedInput },
  name: 'name-and-detail-input',
  props: {
    value: {
      type: Object
    },
    input: {
      type: Function,
      default: function () {
      }
    },
    displayPrice: {
      type: Boolean,
      default: false
    },
    maxNameLength: {
      type: Number
    },
    maxDetailLength: {
      type: Number
    }
  }
}
</script>

<style>

</style>
