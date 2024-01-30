<template>
<div class="container">
  <vue-headful :title="$t('Verification title','Proxeus - Verification')"/>
  <verification></verification>
</div>
</template>

<script>
import Verification from '@/components/document/Verification.vue'

export default {
  name: 'mercuria-validation',
  components: {
    Verification
  }
}
</script>

<style scoped>

</style>
