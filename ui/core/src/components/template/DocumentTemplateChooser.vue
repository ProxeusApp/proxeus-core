<template>
<form enctype="multipart/form-data" novalidate>
  <div v-if="savedFile && lang">
    <b-modal v-model="modalShow" class="b-modal"
             :title="$t('Confirm')"
             :ok-title="$t('Yes')"
             :cancel-title="$t('No')"
             :header-bg-variant="headerBgVariant"
             @hide="onDialogHide"
             @ok="onDialogOk">
      <div class="d-block">{{$t('This action can\'t be undone.')}}</div>
      <div class="d-block">{{$t('Are you sure, you want to delete this file?')}}</div>
      <div class="d-block bold bg-primary"><span class="lang-code">{{lang.Code}}</span>
        <span class="file-name">{{anyFile.name}}</span>
      </div>
    </b-modal>
  </div>
  <div class="file_drop_box m-0 w-100 position-relative" v-bind="$attrs"
       :class="{active: isActive, unsaved:unsavedFile}">
    <h5 class="text-center font-weight-bold bg-txt" v-if="readOnly === false && lang && lang.Code">
      {{ lang.Code }}
    </h5>
    <input type="file" :disabled="uploadPending"
           @change="filesChange($event.target.files)"
           ref="fileInput"
           accept=".odt"
           class="input-file">
    <drop-file-design :detail="$t('choose.or.drag.file.odt', 'Choose an ODT file or drag it here.')"/>
    <table class="file-upload-btns btns-table" v-if="anyFile" :class="{active: isActive, unsaved:unsavedFile}">
      <tr>
        <td class="min" v-show="isWrite()">
          <button @click="removeFile" type="button" class="align-left ml-auto delete-btn fub-btn sb bg-primary"
                  :class="{unsaved:unsavedFile}">
            <i class="mdi mdi-close"></i>
          </button>
        </td>
        <td class="max impcnt">
          <div style="padding-right:6px;" @click="setActive" class="fl-middle-btn fub-btn bg-primary"
               :class="{active: isActive, unsaved:unsavedFile}">
            <span class="material-icons align-middle file-name-icon"
                  style="line-height: unset;">insert_drive_file</span> <span
            style="vertical-align: middle;">{{ anyFileName }}</span>
          </div>
        </td>
        <td class="min" v-show="isWrite() && lang">
          <button :class="{disabled:!unsavedFile}" @click="uploadFile" type="button"
                  class="align-left ml-auto fub-btn sb bg-primary"
                  style="border-color: #40e1d1;border-radius: 0;color: #40e1d1;">
            <i v-if="!uploadPending" class="mdi mdi-upload"></i>
            <i v-else class="mdi mdi-loading mdi-spin"
               style="animation: mdi-spin 0.6s ease-in 0s infinite normal none running;"></i>
          </button>
        </td>
        <td class="min" v-show="lang">
          <a :class="{disabled:!savedFile}" :href="savedFile?downloadUrl:'javascript:void(0)'"
             class="ml-auto raw-download fub-btn sb bg-primary">
            <i class="mdi mdi-download"></i>
          </a>
        </td>
      </tr>
    </table>
  </div>
  <hr v-if="multiDivider">
</form>
</template>

<script>
import bModal from 'bootstrap-vue/es/components/modal/modal'
import bBtn from 'bootstrap-vue/es/components/button/button'
import bModalDirective from 'bootstrap-vue/es/directives/modal/modal'
import DropFileDesign from './DropFileDesign'
import mafdc from '@/mixinApp'

const uploadUrl = '/api/admin/template/upload/:id/:lang'
const deleteUrlWhenUnsaved = '/api/admin/template/ide/delete/:id/:lang'
const deleteUrl = '/api/admin/template/delete/:id/:lang'
const downloadUrlWhenActive = '/api/admin/template/ide/download/:id?raw'
const downloadUrlOtherwise = '/api/admin/template/download/:id/:lang?raw'

export default {
  mixins: [mafdc],
  name: 'document-template-chooser',
  components: {
    DropFileDesign,
    'b-modal': bModal,
    'b-btn': bBtn
  },
  directives: {
    'b-modal': bModalDirective
  },
  props: {
    multiDivider: {
      type: Boolean,
      default: false
    },
    lang: Object,
    initFile: Object,
    readOnly: Boolean,
    renderFile: {
      type: Boolean,
      default: true
    },
    renderLang: {
      type: Boolean,
      default: true
    },
    isActive: {
      type: Boolean
    },
    hasDocuments: {
      type: Boolean
    },
    saveFunc: {
      type: Function
    }
  },
  data () {
    return {
      modalShow: false,
      headerBgVariant: 'light',
      uploadPending: false,
      unsavedFile: undefined,
      savedFile: undefined
    }
  },
  created () {
    if (this.saveFunc) {
      this.saveFunc(this.uploadFile)
    }
  },
  beforeDestroy () {
  },
  mounted () {
    this.reset()
  },
  methods: {
    isWrite () {
      return !this.readOnly
    },
    onDialogHide () {
      this.modalShow = false
    },
    onDialogOk () {
      this.modalShow = false
      this.removeConfirmFromDialog()
    },
    uploadFile () {
      if (!this.unsavedFile || this.uploadPending) {
        return
      }
      this.uploadPending = true
      this.$emit('on-file-upload-pending', this.lang.Code, this.unsavedFile)
      axios.post(uploadUrl.replace(':id', this.$parent.id).replace(':lang', this.lang.Code), this.unsavedFile, {
        headers: {
          'File-Name': encodeURI(this.unsavedFile.name),
          'Content-Type': this.unsavedFile.type
        }
      }).then(response => {
        this.$emit('on-file-upload-success', this.lang.Code, this.unsavedFile)
        this.savedFile = this.unsavedFile
        this.unsavedFile = null
        this.uploadPending = false
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('The template was saved successfully'),
          type: 'success'
        })
      }, (err) => {
        this.uploadPending = false
        this.$emit('on-file-upload-fail', this.lang.Code, this.unsavedFile)
        this.notifyError(this.$t('Could not save template. Please try again or if the error persists contact the platform operator.'))
        this.app.handleError(err)
      })
    },
    removeFile () {
      if (this.unsavedFile) {
        this.unsavedFile = null
        axios.get(deleteUrlWhenUnsaved.replace(':id', this.$parent.id).replace(':lang', this.getLang().Code))
          .then(response => {
            this.$emit('unsavedFile-remove', this.getLang(), this.unsavedFile, !this.savedFile)
          }, (err) => {
            this.app.handleError(err)
          })
      } else {
        this.modalShow = true
      }
    },
    removeConfirmFromDialog () {
      axios.get(deleteUrl.replace(':id', this.$parent.id).replace(':lang', this.lang.Code)).then(response => {
        this.savedFile = null
        this.$emit('setInactive', this.lang, this.unsavedFile)
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Successfully deleted the template.'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Couldn\'t delete template. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    },
    filesChange (fileList) {
      if (fileList.length > 0 && fileList[0]) {
        if (fileList[0].name && /.*\.odt/i.test(fileList[0].name)) {
          this.unsavedFile = fileList[0]
          this.$emit('dropped', this.getLang(), this.unsavedFile)
        }
      }
    },
    reset () {
      this.savedFile = this.initFile
    },
    setActive () {
      if (this.unsavedFile) {
        this.$emit('setActive', this.getLang(), this.unsavedFile)
      } else {
        this.$emit('setActive', this.getLang(), this.savedFile)
        this.$emit('setServerFileActive', this.getLang())
      }
    },
    getLang () {
      if (this.lang) {
        return this.lang
      }
      return { 'Code': 'none' }
    }
  },
  computed: {
    downloadUrl: {
      get () {
        try {
          if (this.isActive) {
            return downloadUrlWhenActive.replace(':id', this.$parent.id)
          }
          return downloadUrlOtherwise.replace(':id', this.$parent.id).replace(':lang', this.getLang().Code)
        } catch (e) {
        }
        return downloadUrlWhenActive.replace(':id', this.$parent.id)
      }
    },
    anyFileName: {
      get () {
        var n = this.anyFile.name
        const max = 14
        if (n.length > max) {
          n = n.substring(0, max) + '..' + n.substring(n.length - 4)
        }
        return n
      }
    },
    anyFile: {
      get () {
        if (this.unsavedFile) {
          return this.unsavedFile
        }
        return this.savedFile
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~bootstrap/scss/functions";
  @import "../../assets/styles/variables.scss";

  .icon-unpersisted {
    position: absolute;
    margin: 0.5rem;
  }

  .input-file {
    opacity: 0; /* invisible but it's there! */
    width: 100%;
    height: 100%;
    position: absolute;
    cursor: pointer;
  }

  .dnd-and-select-file-icon {
    pointer-events: none;
    position: relative;
    width: 100%;
    margin-top: 110px;
  }

  .dnd-and-select-file-icon .outer {
    display: table;
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
  }

  .dnd-and-select-file-icon .middle {
    display: table-cell;
    vertical-align: middle;
  }

  .dnd-and-select-file-icon .inner {
    position: relative;
    margin-left: auto;
    margin-right: auto;
    width: 60%;
  }

  .dnd-and-select-file-icon .arrow-cont {
    border: 6px solid #dddddd;
    border-top: none;
    width: 100%;
    height: 40px;
  }

  .dnd-and-select-file-icon .mdi {
    color: #dddddd;
    width: 100%;
    font-size: 5em;
    text-align: center;
    font-weight: 100;
    top: -50px;
    position: absolute;
  }

  .btns-table {
    width: 100%;
    border-collapse: separate;
    border-spacing: 2px;
    padding: 2px;
    background: transparent;
  }

  .file-upload-btns {
    position: relative;
    border-radius: 0;

    .align-middle {
      line-height: 1;
    }

  }

  .file-upload-btns .fub-btn i {
    font-size: 22px;
    vertical-align: middle;
    text-align: center;
    line-height: unset;
  }

  .file-upload-btns .fub-btn {
    height: 40px;
    vertical-align: middle;
    border-radius: 0;
    outline: none;
    border: 1px solid #06255f;
    color: #ffffff;
    display: inline-block;
    font-weight: 400;
    text-align: center;
    white-space: nowrap;
    vertical-align: middle;
    -webkit-user-select: none;
    -moz-user-select: none;
    user-select: none;
    font-size: 1rem;
    transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
    text-decoration: none;
    background-color: transparent;
    -webkit-text-decoration-skip: objects;
  }

  .file-upload-btns .fub-btn {
    cursor: pointer;
  }

  .file-upload-btns .fub-btn:hover {
    background-color: #06255f !important;
  }

  .delete-btn.unsaved {
    border-color: #40e1d1;
    color: #40e1d1;
    border-radius: 0;
  }

  .file_drop_box {
    border: 2px dashed #dddddd;
    min-width: 155px;

    &.active {
      border: 2px dashed #40e1d1;
    }

    &.unsaved {
      background: rgba(67, 255, 239, 0.2);
      border: 2px dashed #40e1d1;
    }

    .heading-small {
      font-size: 1.2rem;
    }

  }

  .sb {
    width: 40px;
  }

  .fl-middle-btn {
    width: 100%;
    color: white;
    vertical-align: middle;

    &.active {
      text-decoration: underline !important;
    }

    &.unsaved {
      color: #40e1d1;
    }

    &.unsaved {
      color: #40e1d1;
    }

  }

  .lang-code {
    color: #40e1d1;
    text-transform: uppercase;
    font-weight: bold;
    padding-left: 5px;
    padding-right: 5px;
  }

  .file-name {
    color: white;
    padding-right: 5px;
  }

  .file-upload-btns .fub-btn.disabled {
    cursor: not-allowed;
    background-color: #e6e6e6 !important;
    color: #bbbbbb !important;
    border-color: #bbbbbb !important;
  }

  td.min {
    width: 1%;
    white-space: nowrap;
  }

  td.max {
    max-width: 1px;
    white-space: nowrap;
  }

  .file_drop_box:hover {
    background: rgba(67, 255, 239, 0.4);
  }

  .file_drop_box p {
    padding: 10px 0;
    text-align: center;
  }
</style>
