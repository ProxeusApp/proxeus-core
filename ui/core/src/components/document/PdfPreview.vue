<template>
<div class="preview-wrapper d-flex flex-column">
  <div class="actions justify-content-center" :class="{show:showActions}"
       v-if="canDownload">
    <button class="btn btn-primary mx-1" @click="showPdf">
      <i class="mdi mdi-fullscreen text-white"></i>
    </button>
    <a :href="src" class="btn btn-primary mx-1" title="Certified PDF original">
      <i class="mdi mdi-file-pdf text-white"></i>
    </a>
    <a :href="getSrcForTmplRef('docx')" class="btn btn-primary mx-1" title="Microsoft Word un-certified copy">
      <i class="mdi mdi-file-word text-white"></i>
    </a>
    <button v-if="this.showSig!=false" class="btn btn-primary mx-1" @click="signature_request(doc)" title="Request a Signature">
      <i class="mdi mdi-pen text-white"></i>
    </button>
  </div>
  <div class="pdf-preview mb-3 mr-2" :class="{toggled:showActions}">
    <div class="tmpl-filename bg-primary text-white w-100 px-2" v-if="name">
      <small :title="name">{{ name }}</small>
      <span style="display: inline-block;float: right;margin-right: -15px;">
        <table>
          <tr><td>
            <a :href="dynamicDownloadSrc('pdf')" class="float-right" title="PDF un-certified copy">
              <i class="mdi mdi-file-pdf text-white"></i>
            </a>
          </td><td>
            <a :href="dynamicDownloadSrc('docx')" class="float-right px-1" title="Microsoft Word un-certified copy">
              <i class="mdi mdi-file-word text-white"></i>
            </a>
          </td></tr>
        </table>
          </span>
    </div>
    <div class="language-chooser p-0" v-if="shouldDisplayLanguageChooser">
      <div class="input-group input-group-sm mb-">
        <div class="input-group-prepend">
          <label class="input-group-text bg-light p-1 maxwidth" for="inputGroupSelect01">
            <small class="ellipsis">{{$t('Language')}}</small>
          </label>
        </div>
        <select class="custom-select custom-select-sm maxwidth" style="padding: .2rem .4rem;"
                v-model="selectedLanguage"
                id="inputGroupSelect01" @change="error=false">
          <option v-for="lang in getAvailableLangs(languages)" :key="lang" :value="lang">{{ lang }}</option>
        </select>
      </div>
    </div>
    <div class="error w-100 d-flex flex-column align-items-center justify-content-center" v-if="error === true">
      <i class="material-icons">error</i>
      <p class="mt-1">{{$t('Could not load PDF')}}</p>
    </div>
    <spinner v-show="loaded !== true && error !== true" :margin="0" background="transparent" color="#333"
             cls="position-relative no-padding-top mt-0"/>
    <button v-show="loaded && error === false" class="btn btn-link p-0 border-0" @click.prevent="mainAction">
      <pdf :src="getSrc" @loaded="pdfLoaded" @error="pdfError"/>
      <span class="filename d-inline-block t-ellipsis" v-if="filename"><small :title="filename">{{ filename }}</small></span>
    </button>
    <pdf-modal class="pdfwkaround" :src="getSrc" :mid="'modal' + _uid" ref="pdfMod" :filename="filename"
               :download="canDownload"/>
    <signing-modal :src="getSrc" :mid="'smodal' + _uid" ref="signMod"/>
  </div>
</div>
</template>

<script>
import pdf from 'vue-pdf'
import PdfModal from '@/components/document/PdfModal.vue'
import SigningModal from '@/components/document/SigningModal.vue'
import Spinner from '@/components/Spinner.vue'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'pdf-preview',
  props: ['item', 'src', 'filename', 'name', 'languages', 'doc', 'wfId', 'locale', 'langSelectorVisible', 'showSig'],
  components: {
    pdf,
    PdfModal,
    SigningModal,
    Spinner
  },
  data () {
    return {
      loaded: false,
      error: false,
      lang: null,
      showActions: false
    }
  },
  computed: {
    getSrc () {
      return this.src ? this.src + '?format=pdf' : this.dynamicLangPreviewSrc
    },
    canDownload () {
      return !!this.src
    },
    shouldDisplayLanguageChooser () {
      return !this.src && this.langSelectorVisible
    },
    templateLanguage () {
      return this.selectedLanguage
    },
    selectedLanguage: {
      get () {
        if (this.lang) {
          return this.lang
        }
        if (this.languages) {
          let l = null
          this.languages.forEach(lang => {
            if (this.locale === lang) {
              l = lang
            }
          })
          return l || this.languages[0]
        }
        return this.locale
      },
      set (newLanguage) {
        this.lang = newLanguage
        this.loaded = false
        this.$emit('changedTemplateLanguage', newLanguage)
      }
    },
    dynamicLangPreviewSrc () {
      return this.dynamicDownloadSrc('pdf')
    }
  },
  methods: {
    langAvailable (lang) {
      return this.app.isLangAvailable(lang)
    },
    getAvailableLangs (languages) {
      return languages.filter((l) => { this.langAvailable(l) })
    },
    getSrcForTmplRef (format) {
      if (this.item && this.item.ref) {
        return this.src.replace(/[^\/]+$/, this.item.ref) + '?format=' + format
      }
      return this.src + '?format=' + format
    },
    pdfLoaded (doc) {
      this.loaded = true
    },
    pdfError (doc) {
      this.loaded = true
      this.error = true
    },
    mainAction () {
      if (this.canDownload) {
        this.toggleActions()
      } else {
        this.showPdf()
      }
    },
    toggleActions () {
      this.showActions = !this.showActions
    },
    showPdf () {
      const $modal = $('#modal' + this._uid)
      $modal.unbind('shown.bs.modal')
      // this is needed as a workaround because sometimes the pdf canvas height and width stays 0
      $modal.on('shown.bs.modal', function (e) {
        $('.pdfwkaround canvas').css('min-height', '100px')
      })
      $modal.modal('show')
      this.$nextTick(() => {
        this.$refs.pdfMod.load()
      })
    },
    dynamicDownloadSrc (format) {
      return `/api/document/${this.wfId}/preview/${this.doc.id}/${this.selectedLanguage}/` + format
    },
    signature_request () {
      const $smodal = $('#smodal' + this._uid)
      $smodal.modal('show')
      this.$nextTick(() => {
        this.$refs.pdfMod.load()
      })
    }

  }
}
</script>
<style lang="scss" scoped>
  @import "../../assets/styles/variables";

  .maxwidth {
    max-width: 120px;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .ellipsis {
    width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .actions {
    opacity: 0;
    height: 0;
    transition: transform 200ms;
    transform: translateY(30px);
    &.show {
      display: flex;
      opacity: 1;
      height: auto;
      transform: translateY(0px);
    }
  }

  .error {
    background: #666666;
    height: 100px;
    color: white;
    border-top-left-radius: 0;
    border-top-right-radius: 0;
  }

  .btn-file-download {
    border-radius: 0;
  }

  .pdf-preview {
    vertical-align: top;
    position: relative;
    max-width: 120px;
    min-width: 120px;
    min-height: 100px;
    margin: 0 auto;
    transition: box-shadow 150ms;
    border-radius: 2px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24),
    0 10px 0 -5px #fafafa, /* The second layer */
    0 10px 1px -4px rgba(0, 0, 0, 0.2), /* The second layer shadow */
    0 20px 0 -10px #fafafa, /* The third layer */
    0 20px 1px -9px rgba(0, 0, 0, 0.2); /* The third layer shadow */

    &.toggled, &:hover {
      box-shadow: 1px 1px 30px rgba(0, 0, 0, .2);
    }

    .tmpl-filename {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      background: #f7f7f7;
      color: theme-color("primary");
      border-top-left-radius: .2rem;
      border-top-right-radius: .2rem;
    }

    .filename {
      max-width: 150px;
      padding: 5px;
      white-space: normal;
      font-size: 13px;
    }

    .language-chooser {
      border-radius: 0;
      position: relative;
      width: 100%;

      .input-group-text {
        border-left: 0;
        border-bottom-width: 1px;
        border-top-width: 1px;
      }

      .custom-select {
        border-right: 0;
      }

      * {
        border-radius: 0;
      }
    }
  }
</style>
