<template>
<tr class="mshadow-light">
  <td v-bind:class="{'tdmin':cols.length>0}" class="easy-read">{{ lk }}</td>
  <td v-for="(lang,key) in langList" :key="lang.Code"
      v-bind:class="{active:cols.indexOf(key) > -1,'tdmin tdmoremin':cols.indexOf(key) === -1}"
      v-on:click="activateLangCol(key, $event)">
    <div class="iwd" style="position:relative;">
      <i18n-trans-cell v-model="translation" :translations="translations" :langKey="lk" :langCode="lang.Code"
                       @input="handleTranslationChange"></i18n-trans-cell>
      <span class="bg-txt i18n-bg-txt">{{lang.Code}}</span>
    </div>
  </td>
  <td class="tdmin tcntr">
    <button :disabled="hasChanges !== true" type="button" class="btn btn-secondary apdng" @click="saveTranslation"><span
      class="material-icons">check</span></button>
  </td>
</tr>
</template>

<script>
import I18nTransCell from '@/views/appDependentComponents/i18n/I18nTransCell'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'translation',
  props: {
    activateLangCol: Function,
    update: Function,
    translations: Object,
    cols: Array,
    lk: String,
    langList: Array
  },
  components: {
    I18nTransCell
  },
  watch: {
    'translation': 'updateProxy'
  },
  computed: {},
  methods: {
    updateProxy () {
      if (this.update) {
        this.update(this.translation)
      }
    },
    handleTranslationChange (value) {
      this.hasChanges = true
    },
    saveTranslation () {
      let t = {}
      if (this.translations[this.lk]) {
        t[this.lk] = this.translations[this.lk]
        if (this.lk.length <= 4) {
          console.log('updateLangLabel....4')
          this.app.updateLangLabel()
          this.$root.$emit('translations-updated')
        }
        axios.post('/api/admin/i18n/update', t).then(response => {
          this.$notify({
            group: 'app',
            title: this.$t('Success'),
            text: this.$t('Saved translation'),
            type: 'success'
          })
          this.hasChanges = false
        }, (err) => {
          this.app.handleError(err)
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not save translation. Please try again or if the error persists contact the platform operator.\n'),
            type: 'error'
          })
        })
      } else {
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Translation key already exists.'),
          type: 'error'
        })
      }
    }
  },
  data () {
    return {
      hasChanges: false,
      translation: null
    }
  }
}
</script>
<style>
  .i18n-tbl td.tdmax {
    background-color: #f9f9f9;
  }

  .i18n-tbl .iwd {
    height: 100%;
    width: 100%;
    min-height: 100%;
    background: #00000006;
    overflow: hidden;
  }

  .i18n-tbl span.i18n-bg-txt {
    bottom: 6px;
    right: 12px;
    line-height: 1;
    width: auto;
    height: auto;
    font-size: 14px;
    color: #40e1d1;
  }

  .tcntr {
    text-align: center;
  }

  .apdng {
    padding: 1px 5px;
    margin: 2px 8px;
  }
</style>
