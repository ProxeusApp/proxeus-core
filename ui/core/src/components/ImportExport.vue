<template>
<div class="main-container">
  <div>
    <table>
      <tr>
        <td>
          <button class="btn btn-default tabbtn" :class="{active:tabImport}" @click="tabImport=true"
                  @focus="$event.target.blur();" type="button">{{$t('Import')}}
          </button>
        </td>
        <td>
          <button class="btn btn-default tabbtn" :class="{active:!tabImport}" @click="tabImport=false"
                  @focus="$event.target.blur();" type="button">{{$t('Export')}}
          </button>
        </td>
      </tr>
    </table>
    <div ref="scrollable" style="height: 100%;overflow: auto;" class="tabcontent">
      <div v-show="tabImport" ref="inputs">
        <import-only/>
      </div>
      <div v-show="!tabImport">
        <div class="tabcontent">
          <h5>
            <strong>{{$t('Export')}}</strong>
          </h5>
          <table style="width: 100%;">
            <tbody>
            <tr v-if="chooser || app.userIsAdminOrHigher()">
              <td colspan="2">
                <table class="dsel" style="width: 100%;">
                  <tbody>
                  <tr>
                    <td>
                      <div class="row">
                        <div class="col-md-3">
                          <checkbox
                            :title="$t('Export UserData explanation','Choose to include the data that is being resulted by workflow executions')"
                            :label="'UserData'" v-model="ex.UserData"/>
                        </div>
                        <div class="col-md-3">
                          <checkbox
                            :title="$t('Export Doc Template explanation','Choose to include the document template entries')"
                            :label="'Template'" v-model="ex.Template"/>
                        </div>
                        <div class="col-md-3">
                          <checkbox :title="$t('Export Form explanation','Choose to include the forms')" :label="'Form'"
                                    v-model="ex.Form"/>
                        </div>
                        <div class="col-md-3">
                          <checkbox :title="$t('Export Workflow explanation','Choose to include the workflows')"
                                    :label="'Workflow'" v-model="ex.Workflow"/>
                        </div>
                        <div class="col-md-3" v-if="app.userIsSuperAdmin()">
                          <checkbox :title="$t('Export User explanation','Choose to include the users')" :label="'User'"
                                    v-model="ex.User"/>
                        </div>
                        <div class="col-md-3" v-if="app.userIsSuperAdmin()">
                          <checkbox :title="$t('Export I18n explanation','Choose to include internationalization')"
                                    :label="'I18n'" v-model="ex.I18n"/>
                        </div>
                        <div class="col-md-3" v-if="app.userIsRoot()">
                          <checkbox :title="$t('Export Settings explanation','Choose to include the system settings')"
                                    :label="'Settings'" v-model="ex.Settings"/>
                        </div>
                      </div>
                    </td>
                  </tr>
                  </tbody>
                </table>
              </td>
            </tr>
            <tr>
              <td class="tdmax" style="text-align: left;">
                <div class="text-muted">{{$t('Export explanation','Export data by clicking on the button.')}}</div>
              </td>
              <td class="tdmin" style="text-align: right;">
                <button style="white-space: nowrap;" @click="exportClick" :disabled="getParams()?false:true"
                        type="button" class="btn btn-primary"><span>&#8659; {{$t('Export')}}</span></button>
              </td>
            </tr>
            </tbody>
          </table>

        </div>
        <div v-if="lastExportResults" class="tabcontent" style="margin-top:15px;position: relative;">
          <span class="imexclose" @click="app.loadLastExportResults(updateExportResults,'delete')">&#10799;</span>
          <h5>
            <strong>{{$t('Last Export')}}</strong>
          </h5>
          <imex-results :imexResult="lastExportResults"/>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import mafdc from '@/mixinApp'
import FileDropBox from '../components/template/FileDropBox'
import ImexResults from '../components/ImexResults'
import Checkbox from '../components/Checkbox'
import ImportOnly from './ImportOnly'

export default {
  mixins: [mafdc],
  name: 'import-export',
  props: {
    chooser: {
      type: Boolean,
      default: false
    },
    exportTypes: {
      type: Array,
      default: () => {
        return []
      }
    }
  },
  components: {
    ImportOnly,
    Checkbox,
    ImexResults,
    FileDropBox
  },
  data () {
    return {
      ex: {
        User: false,
        UserData: false,
        Form: false,
        I18n: false,
        Workflow: false,
        Template: false
      },
      tabImport: true,
      lastExportResults: null
    }
  },
  created () {
    if (this.exportTypes) {
      for (let i = 0; i < this.exportTypes.length; i++) {
        this.ex[this.exportTypes[i]] = true
      }
    }
    window.addEventListener('resize', this.setListGroupHeight)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.setListGroupHeight)
  },
  mounted () {
    this.setListGroupHeight()
    this.app.loadLastExportResults(this.updateExportResults)
  },
  methods: {
    setListGroupHeight () {
      if (this.$refs.scrollable) {
        let lg = $(this.$refs.scrollable)
        lg.css('height', ($(document.body).height() - 136) + 'px')
      }
    },
    updateExportResults (results) {
      this.lastExportResults = results
    },
    mapSize (obj) {
      let size = 0
      if (obj) {
        for (let key in obj) {
          if (obj.hasOwnProperty(key)) {
            size++
          }
        }
      }
      return size
    },
    exportClick () {
      let params = this.getParams()
      let entities = params.split(',')
      entities.sort()
      let dbName = ''
      for (let i = 0; i < entities.length; i++) {
        dbName += dbName === '' ? entities[i] : '_' + entities[i]
      }
      this.app.exportData(params, this.updateExportResults, null, dbName)
    },
    getParams () {
      let params = ''
      for (let key in this.ex) {
        if (this.ex.hasOwnProperty(key) && this.ex[key]) {
          if (params !== '') {
            params += ',' + key
          } else {
            params += key
          }
        }
      }
      return params
    }
  }
}
</script>

<style lang="scss">
  .imexclose {
    position: absolute;
    top: 2px;
    right: 3px;
    font-size: 28px;
    color: grey;
    cursor: pointer;
    line-height: 1;
  }

  .dsel .fancy-checkbox {
    margin-right: 25px;
  }

  .imptbl {
    background: #0900ff08;
  }

  .imptbl td {
    padding: 0 14px;
    padding-bottom: 8px;
  }

  .tabbtn {
    margin-bottom: -2px;
    border: 1px solid #cecece;
    border-bottom: none !important;
  }

  .tabbtn.active {
    border: 1px solid #062a85;
    background: white;
  }

  .spanel > .spanel-title {
    position: absolute;
    background: white;
    top: -16px;
    font-size: 18px;
  }

  .spanel {
    padding: 15px;
    border: 1px solid #dedede;
    position: relative;
    margin-bottom: 15px;
  }
</style>
