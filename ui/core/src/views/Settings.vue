<template>
  <div>
    <vue-headful :title="$t('Settings title', 'Proxeus - Settings')"/>
    <top-nav :title="$t('Settings')">
      <button slot="buttons" v-if="app.userIsUserOrHigher()"
              @click="app.exportData('', null, '/api/settings/export', 'Settings')" type="button"
              class="btn btn-link ml-2">
        <span>{{$t('Export')}}</span>
      </button>
    </top-nav>
    <div ref="scrollable" class="main-container layout-fixed-scroll-container" style="overflow: auto;">
      <settings-inner/>
    </div>
  </div>
</template>

<script>
import TopNav from '@/components/layout/TopNav'
// import LanguageDropDown from '@/views/appDependentComponents/LanguageDropDown'
import SettingsInner from './appDependentComponents/SettingsInner'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'settings',
  components: {
    SettingsInner,
    TopNav
    // LanguageDropDown
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
