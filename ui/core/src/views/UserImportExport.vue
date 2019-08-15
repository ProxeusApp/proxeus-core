<template>
<div style="height:100%;">
  <vue-headful :title="$t('Data title','Proxeus - Data')"/>
  <top-nav :title="$t('Data')" bg="#ffffff" class="border-bottom-0"/>
  <import-export :chooser="false" :exportTypes="['UserData']"/>
</div>
</template>

<script>
import TopNav from '../components/layout/TopNav'
import importExport from '@/components/ImportExport'

export default {
  name: 'UserImportExport',
  components: {
    TopNav,
    importExport
  }
}
</script>

<style lang="scss">
</style>
