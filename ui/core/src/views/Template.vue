<template>
  <div class="h-100 TplIde">
    <vue-headful
      :title="$t('Template title prefix', 'Proxeus - ')+(template && template.name || $t('Template title', 'Template'))"/>
    <top-nav :title="template && template.name ? template.name : $t('Template')" :sm="true"
             :returnToRoute="{name:'Templates'}">
      <span v-if="hasRenderedDocument">
        <div slot="buttons" v-for="(type, i) in downloadButtons" :key="i">
          <a :href="`/api/admin/template/ide/download/${id}?format=${type}`"
             class="filedownload-btn px-2" target="_blank">
            <table style="text-align: center;">
              <tr>
                <td><span class="d-block text-uppercase">{{ type }}</span></td>
              </tr>
              <tr>
                <td>
                  <i class="mdi mdi-download text-info"></i>
                </td>
              </tr>
            </table>
          </a>
        </div>
      </span>
      <a target="_blank"
         href="https://docs.google.com/document/d/1-vJsTrU3w8dEcDr3-nV5owtxqHWSjzEf2uk6m9-cMIs/preview"
         class="btn btn-link btn-sm"><span
        class="material-icons">help</span>
      </a>
      <button slot="buttons" type="button" class="btn btn-link" @click="previewVisible = !previewVisible">
        <i class="material-icons">view_column</i>
      </button>
      <button slot="buttons" v-if="template && app.userIsUserOrHigher()" type="button" class="btn btn-link"
              @click="openPermissionDialog">
        <span>{{$t('Share')}}</span>
      </button>
      <button slot="buttons"
              v-if="app.userIsUserOrHigher() && template"
              @click="app.exportData('&id='+template.id, null, '/api/template/export','Template_'+template.id)"
              type="button" class="btn btn-link">
        <span>{{$t('Export')}}</span></button>
      <button type="button" v-if="template && app.amIWriteGrantedFor(template)" slot="buttons"
              class="btn btn-link" @click="infoToggled = !infoToggled">
        <span>{{$t('Edit infos')}}</span></button>
      <save-btn slot="buttons" class="ml-2" :item="template" :click="save"/>
    </top-nav>
    <div class="container-fluid layout-fixed-scroll-container">
      <div class="toggle-row row py-2 position-absolute" v-show="infoToggled" v-if="template">
        <name-and-detail-input v-model="template"/>
      </div>
      <div class="toggle-row-backdrop" @click="infoToggled = !infoToggled" v-show="infoToggled"></div>
      <div class="row h-100 align-items-start">
        <div :class="{'col-6 layout-fixed-scroll-view': previewVisible, 'col-12': previewVisible === false}">
          <div class="row mt-3" v-if="hasNoLangs()">
            <div class="col-sm-12 mt-1">
              <document-template-chooser :readOnly="true"
                                         :is-active="readOnly.active"
                                         @on-file-upload-pending="fileUploadPending"
                                         @on-file-upload-success="fileUploadSuccess"
                                         @on-file-upload-fail="fileUploadFail"
                                         @unsavedFile-remove="fileUnsavedRemoved"
                                         @dropped="fileDropped"
                                         @setActive="setActive"
                                         @setInactive="setInactive"
                                         @setServerFileActive="setServerFileActive">
              </document-template-chooser>
            </div>
          </div>
          <div class="row mt-3" v-if="meta !== null && meta && template">
            <div class="col" v-for="lang in meta.activeLangs" :key="lang.Code">
              <document-template-chooser :readOnly="template && !app.amIWriteGrantedFor(template)"
                                         :lang="lang"
                                         :is-active="lang.active"
                                         :multiDivider="true"
                                         @on-file-upload-pending="fileUploadPending"
                                         @on-file-upload-success="fileUploadSuccess"
                                         @on-file-upload-fail="fileUploadFail"
                                         @unsavedFile-remove="fileUnsavedRemoved"
                                         @dropped="fileDropped"
                                         @setActive="setActive"
                                         @setInactive="setInactive"
                                         :saveFunc="getSaveFunc"
                                         @setServerFileActive="setServerFileActive"
                                         :initFile="getFileForLang(lang.Code)"/>
            </div>
          </div>
          <div class="row">
            <div class="col-sm-12" id="formsContainer">
              <search-box v-on:search="search"></search-box>
              <span v-for="form in forms" :form="form" :key="form.id">
                <template-form v-if="form && components"
                               :comps="components"
                               @updatedFormField="handleFormUpdate" />
              </span>
              <div style="display: none">
                <div id="log"></div>
                <div id="inputContainer">
                  <form id="form1">
                    <input type="text" id="msg" placeholder="json" style="width:100%;"/>
                  </form>
                  <form id="form2">
                    <input type="text" id="plainText" placeholder="plain text" style="width:100%;">
                  </form>
                  <input type="text" id="plainTextDirect" placeholder="plain text direct"
                         onkeyup="plainTextDirect(this)"
                         style="width:100%;">
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="col-6 px-0 layout-fixed-scroll-view right-0"
             :class="{overflowHidden:loadingPdf}" v-show="previewVisible">
          <spinner v-show="loadingPdf" cls="position-sticky"/>
          <div class="viewer-container h-100">
            <div id="viewer" class="pdfViewer px-2 pt-2 h-100">
              <object @load="pdfLoaded()" type="application/pdf" :data="src" class="pdf" v-if="src"></object>
            </div>
          </div>
        </div>
        <div id="libreAssistanceContainer">
        </div>
      </div>

    </div>
    <permission-dialog v-if="template" :save="save" :publicLink="app.makeURL('/p/template/'+template.id)"
                       v-model="template"
                       :setup="setupDialog"/>
  </div>

</template>

<script>
import DocumentTemplateChooser from '@/components/template/DocumentTemplateChooser'
import TemplateForm from '@/components/template/TemplateForm'
import _ from 'lodash'
import Spinner from '@/components/Spinner'
import SearchBox from '@/components/SearchBox'
// import bModal from 'bootstrap-vue/es/components/modal/modal'
// import bBtn from 'bootstrap-vue/es/components/button/button'
import bModalDirective from 'bootstrap-vue/es/directives/modal/modal'
import libreHub from '../libs/libreTmplAssistance.min.js'

import Searchable from '../mixins/searchable'
import FormChangeAlert from '../mixins/form-change-alert'
import DelayedInputEvent from '../mixins/delayedInputEvent'

import axios from 'axios'
import TopNav from '@/components/layout/TopNav'
import NameAndDetailInput from '../components/NameAndDetailInput'
import mafdc from '@/mixinApp'
import PermissionDialog from './appDependentComponents/permDialog/PermissionDialog'
import SaveBtn from './appDependentComponents/SaveBtn'

const CancelToken = axios.CancelToken

window.libreHub = libreHub

export default {
  name: 'Template',
  mixins: [mafdc, Searchable, FormChangeAlert, DelayedInputEvent],
  components: {
    SaveBtn,
    PermissionDialog,
    NameAndDetailInput,
    TopNav,
    DocumentTemplateChooser,
    Spinner,
    SearchBox,
    TemplateForm
    // 'b-modal': bModal,
    // 'b-btn': bBtn
  },
  directives: {
    'b-modal': bModalDirective
  },
  // watch: {
  //   infoToggled (newValue, oldValue) {
  //     if (newValue === true) {
  //       this.$nextTick(() => {
  //         this.$refs.inputName && this.$refs.inputName.focus()
  //       })
  //     }
  //   }
  // },
  data () {
    return {
      previewVisible: true,
      infoToggled: false,
      searchTerm: '',
      retries: 2,
      hasChanges: false,
      forms: [],
      components: null,
      template: null,
      openPermissionDialog: function () {
      },
      hasRenderedDocument: false,
      downloadButtons: ['docx', 'doc', 'odt', 'pdf'],
      extDownloadOptions: [
        {
          os: 'linux_x86_64',
          label: 'Linux x86_64',
          icon: 'fa-linux'
        }, {
          os: 'mac_x86_64',
          label: 'Mac x86_64',
          icon: 'fa-apple'
        }, {
          os: 'win_x86_64',
          label: 'Windows x86_64',
          icon: 'fa-windows'
        }, {
          os: 'win_x86',
          label: 'Windows x86',
          icon: 'fa-windows'
        }],
      meta: null,

      interval: null,
      activeFileSize: null,
      lastModified: null,
      loadingPdf: false,
      lastFormID: '',
      src: null,
      lastRequest: undefined,
      dataTimeout: undefined,
      numPages: undefined,
      hasUnpersistedBoxes: false,
      unpersistedBoxes: {},
      lastUploadedFile: undefined,
      lastActiveLang: undefined,
      boxesSaveFuncs: [],
      readOnly: {
        active: false
      },
      saveAfterUpload: false
    }
  },
  computed: {
    id () {
      return this.$route.params.id
    }
  },
  created () {
    this.loadComponents()
    this.loadTemplate(
      (resp) => {
        this.template = resp.data
        $('.panel-body').collapse()
      },
      (err) => {
        if (err.response && err.response.status === 404) {
          this.$_error('NotFound', { what: 'Template' })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t(
              'Could not load template. Please try again or if the error persists contact the platform operator.\n'),
            type: 'error'
          })
          this.$router.push({ name: 'Templates' })
        }
      })
    this.loadLanguageMeta()
    this.loadForms()
    document.addEventListener('keydown', this.keyboardHandler)
    libreHub.updateFileName = function () {}
    // $(document).ready(function () {
    //   libreHub.main({
    //     $container: $('#libreAssistanceContainer'),
    //     fileName: '',
    //     searchVar: function (resultObj, clbk) {
    //       $.get('/api/admin/form/vars?limit=20&wip=1&id=' + resultObj.text, function (data, a, b, c) {
    //         if ($.isArray(data)) {
    //           clbk(data, resultObj)
    //         }
    //       })
    //     }
    //   })
    // })
  },
  beforeDestroy () {
    document.removeEventListener('keydown', this.keyboardHandler)
  },
  methods: {
    hasNoLangs () {
      if (this.app.meta && this.app.meta.activeLangs && this.app.meta.activeLangs.length > 0) {
        return false
      }
      return true
    },
    getSaveFunc (saveFunc) {
      this.boxesSaveFuncs.push(saveFunc)
    },
    keyboardHandler (e) {
      if (e.ctrlKey && e.which === 83) { // Check for the Ctrl key being pressed, and if the key = [S] (83)
        this.save()
        e.preventDefault()
        return false
      }
    },
    setupDialog (openPermissionDialog) {
      this.openPermissionDialog = openPermissionDialog
    },
    hasUnsavedChanges () {
      return this.hasUnpersistedBoxes
    },
    save () {
      if (!this.app.amIWriteGrantedFor(this.template)) {
        return
      }
      if (this.hasUnpersistedBoxes) {
        this.saveAfterUpload = true
        for (let i = 0; i < this.boxesSaveFuncs.length; i++) {
          this.boxesSaveFuncs[i]()
        }
      } else {
        this.saveAfterUpload = false
        this.loadTemplate(
          (resp) => {
            if (resp.data) {
              this.template.data = resp.data.data
              this.hasChanges = true
              axios.post(`/api/admin/template/update?id=${this.template.id}`,
                this.template).then(response => {
                this.hasChanges = false
                this.$notify({
                  group: 'app',
                  title: this.$t('Success'),
                  text: this.$t('Template saved'),
                  type: 'success'
                })
                this.infoToggled = false
              }, (err) => {
                this.app.handleError(err)
                this.hasChanges = false
                console.log(err)
                this.$notify({
                  group: 'app',
                  title: this.$t('Error'),
                  text: this.$t(
                    'Could not update template. Please try again or if the error persists contact the platform operator.\n'),
                  type: 'error'
                })
              })
            }
          }, (err) => {
            console.log(err)
          })
      }
    },
    loadForms (searchTerm) {
      const listUrl = '/api/admin/template/ide/form' + (searchTerm ? '?c=' + searchTerm : '')
      axios.get(listUrl).then((response) => {
        this.forms = response.data
      }, (err) => {
        this.app.handleError(err)
      })
    },
    loadComponents () {
      axios.get('/api/admin/form/component?l=1000').then(res => {
        this.components = res.data
      }, (err) => {
        this.app.handleError(err)
      })
    },
    loadLanguageMeta () {
      axios.get('/api/admin/i18n/meta').then((response) => {
        this.meta = response.data
      }, (err) => {
        this.app.handleError(err)
        this.notifyError(this.$t('Could not load languages.'))
      })
    },
    loadTemplate (success, error) {
      axios.get('/api/admin/template/' + this.id).then((response) => {
        success.apply(this, [response])
      }, (err) => {
        this.app.handleError(err)
        error.apply(this, [err])
      })
    },
    search (term) {
      this.searchTerm = term
      this.timeoutSearch(term, () => {
        this.loadForms(term)
      }, 100)
    },
    getFileForLang (lang) {
      return (this.template.data && this.template.data[lang] && this.template.data[lang].name)
        ? this.template.data[lang]
        : null
    },
    downloadIdeFile (noError) {
      this.loadingPdf = true
      this.lastRequest = CancelToken.source()
      axios.get('/api/admin/template/ide/download/' + this.id, {
        responseType: 'blob',
        cancelToken: this.lastRequest.token
      }).then(res => {
        this.hasRenderedDocument = true
        this.src = URL.createObjectURL(res.data) + '#toolbar=0&navpanes=0'
        this.loadingPdf = false
      }, (err) => {
        this.app.handleError(err)
        this.loadingPdf = false
        if (axios.isCancel(err)) {
        } else {
          this.src = ''
          if (noError === undefined) {
            this.notifyError(this.$t('Could not render PDF. Please try again.'))
          }
        }
      })
    },
    setServerFileActive (lang) {
      this.retries = 2
      this.loadingPdf = true
      if (this.interval) {
        clearInterval(this.interval)
      }
      if (this.lastActiveLang !== lang.Code) {
        axios.get(`/api/admin/template/ide/active/${this.id}/${lang.Code}`).then((res) => {
          this.lastActiveLang = lang.Code
          this.downloadIdeFile()
          this.toggleActiveBox(lang)
        }, (err) => {
          this.app.handleError(err)
          this.loadingPdf = false
          this.notifyError(this.$t('Could not connect to the server. Please try again.'))
        })
      } else {
        this.downloadIdeFile()
      }
    },
    setActive (lang, file) {
      this.retries = 2
      this.toggleActiveBox(lang)
      this.refresh(file, lang)
    },
    setInactive (lang, file) {
      if (lang.Code === this.lastActiveLang) {
        this.retries = 2
        this.$set(lang, 'active', false)
        this.hasRenderedDocument = false
        this.lastModified = new Date()
        this.activeFileSize = 0
        clearInterval(this.interval)
        this.downloadIdeFile(true)
      }
    },
    toggleActiveBox (lang) {
      this.meta.activeLangs.forEach(lang => {
        this.$set(lang, 'active', false)
      })

      if (_.isObject(lang)) {
        this.readOnly.active = false
        this.$set(lang, 'active', true)
      } else {
        this.readOnly.active = true
      }
    },
    refresh (file, lang) {
      this.lastModified = new Date()
      this.activeFileSize = 0

      clearInterval(this.interval)
      if (file && file.name) {
        libreHub.updateFileName(file.name)
      }
      this.interval = setInterval(() => {
        if (file && this.lastModified && file.lastModifiedDate &&
          file.lastModifiedDate.getTime() !== this.lastModified.getTime() &&
          this.activeFileSize !== file.size
        ) {
          this.activeFileSize = file.size
          this.lastModified = file.lastModifiedDate
          this.postFile(file, lang)
        }
      }, 250)
    },
    postFile (file, lang) {
      this.loadingPdf = true
      var curr = lang.Code + '' + file.name + '' + file.lastModifiedDate.getTime()
      if (this.lastUploadedFile !== curr) {
        this.lastUploadedFile = curr
        // active call not needed when calling upload
        axios.post('/api/admin/template/ide/upload/' + this.id + '/' + lang.Code, file, {
          headers: {
            'File-Name': encodeURI(file.name),
            'Content-Type': file.type
          }
        }).then((/* res */) => {
          this.lastActiveLang = lang.Code
          this.numPages = null
          this.downloadIdeFile()
        }, (err) => {
          this.app.handleError(err)
          this.loadingPdf = false
          this.notifyError(this.$t('Couldn\'t upload file.'))
        })
      } else {
        if (this.lastActiveLang !== lang.Code) {
          axios.get('/api/admin/template/ide/active/' + this.id + '/' + lang.Code).then((/* res */) => {
            this.lastActiveLang = lang.Code
            this.downloadIdeFile()
          }, (err) => {
            this.app.handleError(err)
            this.loadingPdf = false
            this.notifyError(this.$t('Could not connect to the server. Please try again.'))
          })
        } else {
          this.downloadIdeFile()
        }
      }
    },
    fileUnsavedRemoved (lang, file, empty) {
      if (lang.Code === this.lastActiveLang) {
        this.downloadIdeFile(true)
        if (empty) {
          this.$set(lang, 'active', false)
          if (this.interval) {
            clearInterval(this.interval)
          }
        }
      }
      delete this.unpersistedBoxes[lang.Code]
      for (var prop in this.unpersistedBoxes) {
        if (!Object.prototype.hasOwnProperty.call(this.unpersistedBoxes, prop)) continue
        this.hasUnpersistedBoxes = true
        return
      }
      this.hasUnpersistedBoxes = false
    },
    fileDropped (lang, file) {
      this.toggleActiveBox(lang)
      this.refresh(file, lang)

      if (_.isObject(lang)) {
        this.hasUnpersistedBoxes = true
        this.unpersistedBoxes[lang.Code] = file
      }
    },
    fileUploadPending (lang, file) {
      this.loadingPdf = false
    },
    fileUploadSuccess (lang, file) {
      delete this.unpersistedBoxes[lang]
      for (var prop in this.unpersistedBoxes) {
        if (Object.prototype.hasOwnProperty.call(this.unpersistedBoxes, prop)) {
          this.hasUnpersistedBoxes = true
          return
        }
      }
      this.hasUnpersistedBoxes = false
      this.unpersistedBoxes = {}
      this.hasFormChanges = false
      if (this.saveAfterUpload) {
        this.save()
      }
    },
    fileUploadFail (file, lang) {
      this.loadingPdf = false
    },
    handleFormUpdate (formID) {
      this.lastFormID = formID
      if (this.hasRenderedDocument === false) {
        return
      }
      if (this.loadingPdf === true && this.lastRequest) {
        this.lastRequest.cancel()
      }

      this.loadingPdf = true

      clearTimeout(this.dataTimeout)

      this.dataTimeout = setTimeout(() => {
        this.downloadIdeFile()
      }, 3000)
    },
    pdfLoaded () {
      this.loadingPdf = false
      this.hasRenderedDocument = true
    },
    handlePdfError (/* err */) {
      this.hasRenderedDocument = false
      this.loadingPdf = false
      if (this.retries > 0) {
        this.retries = this.retries - 1
        this.downloadIdeFile()
      }
    },
    hideModal () {
      this.$refs.extModal.hide()
    },
    notifyError (error) {
      this.$notify({
        group: 'app',
        title: this.$t('Error'),
        text: error || this.$t(
          'Could not connect to the server. Please try again or if the error persists contact the platform operator.\n'),
        type: 'error',
        duration: 5000
      })
    }
  }
}
</script>

<style lang="scss">
  @import "../assets/styles/librehub.scss";
</style>

<style lang="scss" scoped>

  .pdf {
    width: 100%;
    height: 100%;
  }

  div#formsContainer .search-box {
    margin-bottom: 6px;
  }

  div#formsContainer {
    padding-bottom: 20px;
  }

  .TplIde {

    .card {
      border-radius: 0;
      border: none;
    }

    .filedownload-btn {
      display: inline-block;
      text-decoration: none !important;
      border: 1px solid transparent;
      border-radius: 0;

      .mdi-download {
        font-size: 30px;
        line-height: 0.1 !important;
        max-height: 10px !important;
        bottom: -5px;
        position: relative;
      }

      span {
        line-height: 1;
        font-size: 12px;
      }

    }

    .filedownload-btn:focus {
      outline: 0 !important;
      border: none !important;
    }

    .filedownload-btn:hover {
      background: rgba(0, 0, 0, .1);
    }

    .navbar {
      -webkit-box-shadow: 1px 1px 20px rgba(0, 0, 0, .2);
      -moz-box-shadow: 1px 1px 20px rgba(0, 0, 0, .2);
      box-shadow: 1px 1px 20px rgba(0, 0, 0, .2);
    }

    .navbar-text {
      overflow: hidden;
      min-width: 0;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .topnav-buttons {
      white-space: nowrap;
      text-overflow: ellipsis;

      > * {
        display: inline-block;
      }

    }

    .float-left {
      float: none !important;
    }

    .popover {
      background: rgba(0, 0, 0, .5);
    }

    .overflowHidden {
      overflow: hidden;
    }

  }

  .layout-fixed-scroll-container {
    height: calc(100% - 60px);
    position: relative;
  }

  .layout-fixed-scroll-view {
    position: absolute;
    overflow-y: auto;
    overflow-x: hidden;
    top: 0;
    bottom: 0;
  }

  .right-0 {
    right: 0;
  }

  div.viewer-container {
    width: 100%;
    background: #333333;
    height: 100%;
    position: relative;
  }

  .form-panel {
    border: 1px solid #eeeeee;
    margin-bottom: 1rem;

    h4 {
      font-size: 1rem;
    }

  }

  .pdfViewer {
    width: 100%;
    height: 100%;
    padding-bottom: 5px;
  }

</style>
