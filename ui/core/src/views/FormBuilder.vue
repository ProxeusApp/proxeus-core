<template>
    <div class="">
        <vue-headful :title="$t('FormBuilder title prefix','Proxeus - ')+(form && form.name || $t('Formbuilder title', 'Formbuilder'))" v-if="!userSrc"/>
        <top-nav :title="form.name || $t('Formbuilder title', 'Formbuilder')" :sm="true" v-if="modalMode === false && form" class="mb-0"
                 :returnToRoute="{name:'Forms'}">
            <button slot="buttons" v-if="form && app.userIsUserOrHigher()" class="btn btn-link" @click="openPermissionDialog">
                {{ $t('Share') }}
            </button>
            <button slot="buttons" v-if="form && app.userIsUserOrHigher()" @click="app.exportData('&id='+form.id, null, '/api/form/export','Form_'+form.id)" type="button" class="btn btn-link">{{$t('Export')}}</button>
            <button slot="buttons" v-if="app.amIWriteGrantedFor(form)" class="btn btn-link" @click="infoToggled = !infoToggled">
                <span>{{$t('Edit infos')}}</span>
            </button>
          <button slot="buttons" class="btn btn-primary ml-2" @click="save" :disabled="app.amIWriteGrantedFor(form) === false">Save</button>
        </top-nav>
        <div class="container-fluid">
            <div class="toggle-row row py-3 position-absolute" v-if="form" v-show="infoToggled">
                <name-and-detail-input v-model="form"/>
            </div>
            <form id="inputFormId" novalidate="novalidate" class="app-formbuilder row h-100">
                <div class="dynamic-top col-sm-12 p-layout-fixed-scroll-container">
                    <div class="toggle-row-backdrop" @click="infoToggled = !infoToggled" v-show="infoToggled"></div>
                    <div id="formBuilderMainTag" class="row">
                    </div>
                </div>
            </form>
        </div>
        <permission-dialog v-if="form" :save="save" :publicLink="app.makeURL('/p/form/'+form.id)" v-model="form" :setup="setupDialog"/>
    </div>
</template>

<script>
import TopNav from '@/components/layout/TopNav'
// eslint-disable-next-line no-unused-vars,camelcase
import _ from 'lodash'
// eslint-disable-next-line camelcase
import FT_FormBuilder from '../libs/legacy/formbuilder'
import NameAndDetailInput from '../components/NameAndDetailInput'
import PermissionDialog from './appDependentComponents/permDialog/PermissionDialog'
import SaveBtn from './appDependentComponents/SaveBtn'
import mafdc from '@/mixinApp'
import formChangeAlert from '../mixins/form-change-alert'
export default {
  mixins: [mafdc, formChangeAlert],
  name: 'form-builder',
  props: ['userSrc', 'modal'],
  components: {
    SaveBtn,
    PermissionDialog,
    NameAndDetailInput,
    TopNav
  },
  beforeDestroy () {
    $.contextMenu('destroy')
    document.removeEventListener('keydown', this.keyboardHandler)
  },
  data () {
    return {
      formSrc: null,
      myFormBuilder: null,
      form: null,
      infoToggled: false,
      components: [],
      templateEditorHasChanges: false,
      openPermissionDialog: function () {}
    }
  },
  computed: {
    id () {
      return this.$route.params.id
    },
    modalMode () {
      return this.modal || false
    },
    superadmin () {
      return this.app.userIsSuperAdmin()
    },
    userIsAdminOrHigher () {
      return this.app.userIsAdminOrHigher()
    }
  },
  watch: {
    infoToggled (newValue, oldValue) {
      if (newValue === true) {
        this.$nextTick(() => {
          this.$refs.inputName && this.$refs.inputName.focus()
        })
      }
    }
  },
  created () {
    // We use this to clean all instanciated context menus as it is
    // a likely source of bugs
    $.contextMenu('destroy')
    document.addEventListener('keydown', this.keyboardHandler)
    if (this.userSrc) {
      this.formSrc = this.userSrc
      this.loadComponents((store) => {
        this.initFb(store)
      })
    } else {
      axios.get('/api/admin/form/' + this.id).then(response => {
        // JSON responses are automatically parsed.
        this.form = response.data
        this.snapshot(this.form, this.skipFromSnapshot)
        if (response.data && response.data.id) {
          if (response.data.data) {
            this.formSrc = response.data.data.formSrc
          }
          this.loadComponents((store) => {
            this.initFb(store)
          })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not load form. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
          this.$router.push({ name: 'Forms' })
        }
      }, (error) => {
        this.app.handleError(error)
        if (error.response && error.response.status === 404) {
          this.$_error('NotFound', { what: 'Form' })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not load form. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
          this.$router.push({ name: 'Forms' })
        }
      })
    }
  },
  methods: {
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
      if (this.myFormBuilder) {
        let frmData = this.myFormBuilder.getData()
        if (frmData) {
          if (!this.form) {
            this.form = {}
          }
          this.form.data = frmData
        }
      }
      return !this.compare(this.form, this.skipFromSnapshot)
    },
    skipFromSnapshot (value, keyOrIndex, obj, fullpath) {
      if (/^_fbonly_uipp$/m.test(keyOrIndex)) {
        return true// skip
      }
      return false
    },
    loadComponents (done) {
      axios.get('/api/admin/form/component?l=1000').then(res => {
        this.components = {}
        var has = false
        Object.entries(res.data).sort((a, b) => {
          return a[1].sortIndex - b[1].sortIndex
        }).forEach((entry) => {
          has = true
          this.components[entry[0]] = entry[1]
        })
        done(!has)
      }, (err) => {
        this.app.handleError(err)
      })
    },
    save () {
      if (!this.form) {
        this.form = {}
      }
      if (!this.app.amIWriteGrantedFor(this.form)) {
        return
      }
      let formToSave
      this.form.data = this.myFormBuilder.getData()
      formToSave = this.form
      let self = this
      $.ajax({
        url: '/api/admin/form/update?id=' + this.$route.params.id,
        type: 'POST',
        data: JSON.stringify(formToSave),
        contentType: 'application/json',
        error: function (error) {
          if (error && error.status && error.status === 401) {
            redirectToLogin()
          } else {
            error = error.responseJSON
            if (error && error.msg) {
              self.$notify({
                group: 'app',
                title: self.$t('Error'),
                text: error.msg,
                type: 'error'
              })
            }
          }
        },
        success: function (res) {
          self.snapshot(self.form, self.skipFromSnapshot)
          self.hasFormChanges = false
          self.$notify({
            group: 'app',
            title: self.$t('Success'),
            text: self.$t('Saved form'),
            type: 'success'
          })
        }
      })
    },
    initFb (store) {
      let _this = this
      jQuery(function () {
        let myHtmlEntry = function (id, langCode, text) {
          return '<div class="fb-i18n-entry"><p class="fb-i18n-key">' + id + '</p><span class="fb-i18n-lang">' +
              langCode + '</span><p class="fb-i18n-text">' + text + '</p></div>'
        }
        let formSrc = _this.formSrc

        let formBuilderOptions = {
          userAllowedToEditComponents: _this.superadmin,
          enableI18n: _this.superadmin,
          autoSaveSettings: true,
          varNameFunction: function (text, process) {
            $.get('/api/admin/template/vars?c=' + text, function (data) {
              process(data)
            })
          },
          varConnectorsFunction: function (text, process) {
            $.get('/api/admin/connector/list?c=' + text, function (data) {
              process(data.map(x => x.name))
            })
          },
          varUserFunction: function (text, process) {
            $.get('/api/admin/form/users?c=' + text, function (data) {
              process(data)
            })
          },
          file: {
            requestTypes: function (callback) {
              if (callback) {
                $.get('/api/admin/form/file/types', function (data) {
                  callback(data)
                })
              }
            }
          },
          test: {
            onActive: function ($form, formSrc) {
              $.get('/api/admin/form/test/data/' + _this.id, function (data) {
                $form.fillForm(data)
              })
              setFormSrcFirstIfNecessary(formSrc)
              let changeOptions = {
                beforeSend: function () {
                  setFormSrcFirstIfNecessary(formSrc)
                },
                fileUrl: '/api/admin/form/test/file/' + _this.id,
                url: '/api/admin/form/test/data/' + _this.id,
                data: { formSrc: JSON.stringify(formSrc) }
              }
              $form.assignSubmitOnChange(changeOptions)
              $form.on('formFieldsAdded', function (event, parent) {
                if (parent && parent.length) {
                  parent.assignSubmitOnChange(changeOptions)
                }
              })
            },
            onSubmit: function (formEle, formSrc, cb) {
              let d = formEle.serializeFormToObject()
              $.ajax({
                type: 'POST',
                url: '/api/admin/form/test/data/' + _this.id + '?s=true',
                data: JSON.stringify(d),
                contentType: 'application/json; charset=utf-8',
                success: function (data, textStatus, jqXHR) {
                  if (jqXHR.status === 200) {
                    formEle.cleanFieldErrors()
                    if ($.isFunction(cb)) {
                      cb(true)
                    }
                  }
                },
                error: function (res, textStatus, errorThrown) {
                  if (res && res.status) {
                    if (res.status === 401) {
                      redirectToLogin()
                    }
                    if (res.status === 412 || res.status === 422) {
                      let data = res.responseJSON
                      formEle.showFieldErrors(data)
                    } else {
                      if ($.isFunction(cb)) {
                        cb(false)
                      }
                    }
                  }
                }
              })
            },
            onReset: function () {
              $.get('/api/admin/form/test/data/' + _this.id + '?reset=true', function (data) {
                console.log('reset success')
              })
            }
          },
          i18n: {
            onTranslate: function (keyArray, callback, notAsync) {
              $.ajax({
                type: 'POST',
                url: '/api/admin/i18n/translate',
                data: JSON.stringify(keyArray),
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                async: !notAsync,
                success: function (data) {
                  callback(data)
                },
                error: function (res) {
                  if (res && res.status && res.status === 401) {
                    redirectToLogin()
                  } else {
                    callback(keyArray)
                  }
                }
              })
            },
            onSearch: function (text, callback) {
              $.get('/api/admin/i18n/search?c=' + text).done(function (data) {
                let dataList = []
                if (data) {
                  let id, langCode, dataIDObj, langText
                  for (id in data) {
                    dataIDObj = data[id]
                    if (data.hasOwnProperty(id) && dataIDObj) {
                      for (langCode in dataIDObj) {
                        langText = dataIDObj[langCode]
                        if (data.hasOwnProperty(id) && langText) {
                          dataList.push({
                            id: id,
                            name: myHtmlEntry(id, langCode, langText)
                          })
                        }
                      }
                    }
                  }
                } else {
                  if (/^[a-zA-Z0-9\.\_\-]+$/.test(text)) {
                    dataList.push({
                      id: text,
                      name: myHtmlEntry(text, '??',
                        '<span class="fb-i18n-key-new">{{i18n("i18n.code_not_found")}}</span>')
                    })
                  }
                }
                try {
                  callback(dataList)
                } catch (doneE) {
                }
              }).fail(function () {
              })
            },
            // TODO refactor as it is not needed outside the formbuilder
            onSelect: function (data) {
              return { i18n: data }
            },
            onDisplay: function (data) {
              if (typeof data === 'object' && data.hasOwnProperty('i18n') && data['i18n']) {
                return data['i18n']
              }
              return data
            },
            isCovered: function (data) {
              return typeof data === 'object' && data.hasOwnProperty('i18n') && data['i18n']
            }
            // TODO -------
          },
          data: formSrc,
          component: {
            requestComp: function (id, callback) {
              if (id) {
                if (callback) {
                  callback(_this.components[id])
                } else {
                  return _this.components[id] || {}
                }
              } else {
                if (callback) {
                  callback(_this.components)
                } else {
                  return _this.components
                }
              }
            },
            searchComp: function (text, callback) {
              let getOpts = {
                type: 'GET',
                thisCallback: callback,
                url: '/api/admin/form/component?l=1000' + (text ? '&c=' + text : ''),
                success: function (data) {
                  this.thisCallback(data)
                },
                error: function () {
                  if (this.thisCallback) {
                    this.thisCallback()
                  }
                }
              }
              $.ajax(getOpts)
            },
            storeComp: function (comp, callback) {
              $.ajax({
                dfsId: comp.id,
                url: '/api/admin/form/component?id=' + comp.id,
                method: 'POST',
                data: JSON.stringify(comp),
                contentType: 'application/json; charset=utf-8',
                success: function (data) {
                  if (data && data.id) {
                    if (this.dfsId !== data.id) {
                      callback({ oldId: this.dfsId, newId: data.id })
                    }
                  }
                },
                error: function (res) {
                  if (res && res.status && res.status === 401) {
                    redirectToLogin()
                  }
                }
              })
            },
            deleteComp: function (id, callback) {
              $.ajax({
                dfsId: id,
                url: '/api/admin/form/component/' + id,
                method: 'DELETE',
                success: function (data) {
                  callback()
                },
                error: function (res) {
                  if (res && res.status && res.status === 401) {
                    redirectToLogin()
                  }
                }
              })
            }
          }
        }
        if (!_this.app.amIWriteGrantedFor(_this.form)) {
          formBuilderOptions.readOnly = true
        }
        _this.myFormBuilder = {
          $main: $('#formBuilderMainTag'),
          fb: null,
          init () {
            this.fb = new FT_FormBuilder(this.$main, formBuilderOptions)
          },
          getData () {
            return this.fb.getData()
          }
        }
        _this.myFormBuilder.init()
        window.__fb = _this.myFormBuilder
        if (store) {
          _this.myFormBuilder.fb.componentsTab.storeComponents()
        }
      })
      let lastFromSrc

      function setFormSrcFirstIfNecessary (formSrc) {
        let currentFormSrc = JSON.stringify(formSrc)
        if (currentFormSrc !== lastFromSrc) {
          lastFromSrc = currentFormSrc
          $.ajax({
            type: 'POST',
            url: '/api/admin/form/test/setFormSrc/' + _this.id,
            data: currentFormSrc,
            contentType: 'application/json',
            async: false
          })
        }
      }
    }
  }
}
</script>

<style>

</style>

<style lang="less">

    .app-formbuilder {
        @import "../assets/styles/legacy/backend.less";
        @import "../libs/legacy/formbuilder/contextMenu/contextmenu.css";
        @import "../libs/legacy/formbuilder/splitpane/split-pane.css";
        @import "../assets/styles/legacy/formbuilder.less";

    }

    /* Hide Autoswitch for now */
    .switch-main {
        display: none !important;
    }

    #formBuilderMainTag {
        margin-top: 1rem;
    }

    .fb-component.selected {
        background: rgba(255, 255, 255, .9);
    }

    .hcbuilder-settings-body input:invalid {
        color: red !important;
        border: 1px solid red;
    }

    .app-formbuilder .hcbuild-main .hcbuilder-settings-body .validation-parent .fb-field-group .fbs-lbl {
        left: 0;
    }

    .app-formbuilder .hcbuild-main .split-layout-main {
        min-height: 0;
    }

    .app-formbuilder {

        .nav-fb {
            background: #062a85;
            /*border: 1px solid #cecece;*/
            border-bottom: 0;
            padding: 0px .75rem;
            border-top-left-radius: 4px;
            border-top-right-radius: 4px;

            .nav-link {
                color: white;
                background: #062a85;
                cursor: pointer;
                padding-top: 9px;
                padding-bottom: 4px;
                padding-left: 15px;
                padding-right: 15px;
                margin-right: .25rem;
                margin-left: 0px;
                font-size: .9rem;
                border: none;
                border-bottom: 5px solid transparent;
                &.active, &:hover {
                    background: #062a85;
                    color: white;
                    border: none;
                    border-radius: 0px;
                    border-bottom: 5px solid #40e1d1;
                }
                &:hover {
                    color: #eee;
                    border-color: #eee;
                    font-weight: normal;
                }
            }
        }

        .dropdown-menu.show {
            display: block !important;
        }

        .hcbuild-main .tab-pane .panel .panel-body {
            height: calc(100vh - 255px);
            overflow-y: auto;
            overflow-x: hidden;
        }

        .hcbuild-main .tab-pane.htmlComponentBuilder .panel .panel-body {
            height: calc(100vh - 185px);
            overflow-y: auto;
            overflow-x: hidden;
        }
        .hcbuild-main .tab-pane.htmlComponentSettings .panel .panel-body {
            height: calc(100vh - 185px);
        }

        .hcbuild-main .hcbuild-workspace-body {
            overflow-y: auto;
            overflow-x: hidden;
        }

        .hcbuild-main .hcbuild-workspace-body .ws-holder {
            height: calc(100vh - 130px);
        }

        .hcbuild-workspace-test-main .panel-body {
            height: calc(100vh - 190px);
            overflow-y: auto;
            overflow-x: hidden;
        }

        .container-fluid.hcbuild-workspace-body {
            padding-left: 0;
            padding-right: 0;
        }

        .help-block {
            color: #a8b0b8;
        }
        .navbar {
            display: flex;
        }
    }

</style>
