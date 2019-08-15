<template>
<div class="doc-create">
  <vue-headful :title="$t('New Document title', 'Proxeus - New Document')"/>
  <top-nav :title="$t('New Document')" bg="#f7f8f9" :returnToRoute="{name:'Documents'}">
  </top-nav>
  <div class="main-container">
    <list-group icon="description" :select="false" :timestamps="false"
                linkPrefix="" :enableCheckboxes="false" nodeType="document" fixedPath="/api/user/workflow/list"
                path="workflow" :listOnly="true" :payment="true">
    </list-group>
  </div>
</div>
</template>

<script>
import ListGroup from '@/components/ListGroup'
import TopNav from '@/components/layout/TopNav'

export default {
  name: 'document-create',
  components: {
    ListGroup,
    TopNav
  }
}
</script>

<style scoped>

</style>
