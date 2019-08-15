<template>
<div>
  <vue-headful :title="$t('Settings title', 'Proxeus - Settings')"/>
  <top-nav :title="$t('Settings')">
    <td slot="td" class="tdmin" v-if="app.userIsUserOrHigher()">
      <button style="height: 40px;" @click="app.exportData('', null, '/api/settings/export', 'Settings')" type="button"
              class="btn btn-primary ml-2">
        <i style="font-style: normal;font-size: 18px;">&#8659;</i>
        <span>{{$t('Export')}}</span></button>
    </td>
  </top-nav>
  <div class="container-fluid" style="padding:20px 20px;">
    <settings-inner/>
  </div>
</div>
</template>

<script>
import TopNav from '@/components/layout/TopNav'
import LanguageDropDown from '@/views/appDependentComponents/LanguageDropDown'
import SettingsInner from './appDependentComponents/SettingsInner'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'settings',
  components: {
    SettingsInner,
    TopNav,
    LanguageDropDown
  },
  created () {
  },
  mounted () {
  },
  beforeDestroy () {
  },
  updated () {
  },
  methods: {},
  data () {
    return {}
  }
}
</script>
