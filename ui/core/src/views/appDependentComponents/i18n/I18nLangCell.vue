<template>
<checkbox style="margin-top:5px;" label="" v-model="languageActive"/>
</template>

<script>
import mafdc from '@/mixinApp'
import Checkbox from '@/components/Checkbox.vue'

export default {
  components: { Checkbox },
  mixins: [mafdc],
  name: 'i18n-lang-cell',
  props: ['lang'],
  computed: {
    languageActive: {
      get () {
        return this.lang.Enabled
      },
      set (value) {
        axios.post('/api/admin/i18n/lang', { [this.lang.Code]: value }).then(res => {
          this.app.loadMeta()
          this.$notify({
            group: 'app',
            title: this.$t('Success'),
            text: value ? this.$t('Activated language') : this.$t('Deactivated language'),
            type: 'success'
          })
        }, (err) => {
          this.app.handleError(err)
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not save language. Please try again or if the error persists contact the platform operator.\n'),
            type: 'error'
          })
        })
      }
    }
  }
}
</script>
