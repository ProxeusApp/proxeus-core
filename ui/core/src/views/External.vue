<template>
  <iframe v-bind:src="iframeUrl" style="width: 100%; height: 100%; border: 0;" />
</template>

<script>
'use strict'
import mafdc from '@/mixinApp'

export default {
  name: 'External',
  mixins: [mafdc],
  components: {},
  data () {
    return {
      iframeUrl: '/api/admin/external/' + this.$route.params.name + '/' + this.$route.params.id
    }
  }
}
</script>
