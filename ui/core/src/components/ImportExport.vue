<template>
<div class="main-container">
  <div>
    <nav-tabs class="mt-0">
      <tab title="Import" :selected="true">
        <import-only/>
      </tab>
      <tab title="Export">
        <div class="tabcontent">
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
            </tbody>
          </table>
          <button @click="exportClick" :disabled="getParams()?false:true"
                  type="button" class="btn btn-primary mt-3"><span>{{$t('Export')}}</span></button>
        </div>
        <div v-if="lastExportResults" class="tabcontent" style="margin-top:15px;position: relative;">
          <span class="imexclose" @click="app.loadLastExportResults(updateExportResults,'delete')">&#10799;</span>
          <h5>
            <strong>{{$t('Last Export')}}</strong>
          </h5>
          <imex-results :imexResult="lastExportResults"/>
        </div>
      </tab>
    </nav-tabs>
  </div>
</div>
</template>

<script>
import mafdc from '@/mixinApp'
import ImexResults from '../components/ImexResults'
import Checkbox from '../components/Checkbox'
import ImportOnly from './ImportOnly'
import NavTabs from '@/components/nav-tabs/NavTabs'
import Tab from '@/components/nav-tabs/Tab'

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
    NavTabs,
    Tab
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
        const lg = $(this.$refs.scrollable)
        lg.css('height', ($(document.body).height() - 136) + 'px')
      }
    },
    updateExportResults (results) {
      this.lastExportResults = results
    },
    mapSize (obj) {
      let size = 0
      if (obj) {
        for (const key in obj) {
          if (Object.prototype.hasOwnProperty.call(obj, key)) {
            size++
          }
        }
      }
      return size
    },
    exportClick () {
      const params = this.getParams()
      const entities = params.split(',')
      entities.sort()
      let dbName = ''
      for (let i = 0; i < entities.length; i++) {
        dbName += dbName === '' ? entities[i] : '_' + entities[i]
      }
      this.app.exportData(params, this.updateExportResults, null, dbName)
    },
    getParams () {
      let params = ''
      for (const key in this.ex) {
        const exHasProp = Object.prototype.hasOwnProperty.call(this.ex, key)
        if (exHasProp && this.ex[key]) {
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
  @import "@/assets/styles/variables.scss";

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
    background: $gray-100;
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
