<template>
<simple-select class="langdrpdown" :unselect="false" v-if="hasLangs()" :selected="app.getSelectedLang" idProp="Code"
               :onSelect="onSelect" :options="app.meta.activeLangs"/>
</template>
<script>
import SimpleSelect from '@/components/SimpleSelect'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  components: { SimpleSelect },
  name: 'language-drop-down',
  methods: {
    hasLangs () {
      if (this.app.meta &&
        this.app.meta.activeLangs &&
        this.app.meta.activeLangs.length > 0 &&
        this.app.meta.activeLangs[0].label) {
        return true
      }
      return false
    },
    onSelect (item, index) {
      this.app.setSelectedLang(item.Code)
    },
    updateSimpleSelect () {
      try {
        this.$forceUpdate()
        this.$children[0].$forceUpdate()
      } catch (e) {
      }
    },
    translationsUpdated () {
      this.updateSimpleSelect()
    },
    metaUpdated () {
      this.updateSimpleSelect()
    }
  },
  created () {
    this.$root.$on('translations-updated', this.translationsUpdated)
    this.$root.$on('meta-updated', this.metaUpdated)
  },
  beforeDestroy () {
    this.$root.$off('meta-updated', this.metaUpdated)
    this.$root.$off('translations-updated', this.translationsUpdated)
  },
  computed: {
    selected: {
      get () {
        return this.app.getSelectedLangIndex()
      },
      set (a) {
        // do nothing.. read only
      }
    }
  }
}
</script>
<style>
  .langdrpdown .sss > span, .langdrpdown .ss-list li {
    white-space: nowrap;
    vertical-align: middle;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .langdrpdown .sss > span {
    max-width: 200px;
    display: inline-block;
    padding: 2px 0;
  }

  .langdrpdown .ss-list li {
    max-width: 226px;
  }
</style>
