<template>
<div class="verification">
  <vue-headful :title="$t('Verification title', 'Proxeus - Verification')"/>
  <top-nav :title="$t('Document verification')" bg="#f7f8f9" class="border-bottom-0"/>
  <div class="col-12">
    <verification :thumbnail="true"></verification>
  </div>
</div>
</template>

<script>
import SearchBox from '@/components/SearchBox'
import ListGroup from '@/components/ListGroup'
import TopNav from '@/components/layout/TopNav'
import Verification from '@/components/document/Verification'

export default {
  name: 'document-verification',
  components: {
    SearchBox,
    ListGroup,
    TopNav,
    Verification
  }
}
</script>

<style lang="scss" scoped>
  .verification {
    overflow: hidden;
    height: auto;
  }
</style>
