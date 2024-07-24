<template>
    <div class="app-workflow">
        <vue-headful
                :title="$t('Workflow title prefix', 'Proxeus - ')+(workflow && workflow.name || $t('Workflow title', 'Workflow'))"/>
        <top-nav :title="workflow && workflow.name ? workflow.name : 'Workflow'" :returnToRoute="{name:'Workflows'}">
            <button slot="buttons" v-if="workflow && app.userIsUserOrHigher()" class="btn btn-link" @click="openPermissionDialog">
                <span>{{$t('Share')}}</span>
            </button>
            <button slot="buttons" v-if="workflow && app.amIWriteGrantedFor(workflow)" class="btn btn-link" @click="publishModalShow = true">
                <span>{{$t('Publish')}}</span>
            </button>
            <button slot="buttons" v-if="app.userIsUserOrHigher()" @click="app.exportData('&id='+id, null, '/api/workflow/export', 'Workflow_'+id)"
                    type="button" class="btn btn-link">
                <span>{{$t('Export')}}</span>
            </button>
            <button slot="buttons" v-if="workflow && app.amIWriteGrantedFor(workflow)" class="btn btn-link" @click="infoToggled = !infoToggled">
                <span>{{$t('Edit infos')}}</span>
            </button>
            <a slot="buttons" target="_blank" :href="'/document/'+id" class="btn btn-primary ml-1">
              Run
            </a>
            <button slot="buttons" class="btn btn-primary ml-2" @click="save" :disabled="app.amIWriteGrantedFor(workflow) === false">Save</button>
        </top-nav>
        <div class="container-fluid p-0">
            <div class="modal priceErrorModal" ref="priceErrorModal" tabindex="-1" role="dialog">
              <div class="modal-dialog" role="document">
                <div class="modal-content">
                  <div class="modal-header">
                    <h5 class="modal-title">{{$t('Could not save workflow')}}</h5>
                    <button @click="togglePriceErrorModal(false)" type="button" class="close" data-dismiss="modal" aria-label="Close">
                      <span aria-hidden="true">&times;</span>
                    </button>
                  </div>
                  <div class="modal-body">
                    <p>{{$t('Before setting a workflow price make sure to set your ethereum address in your account settings on the top right.')}}</p>
                  </div>
                  <div class="modal-footer">
                    <button @click="togglePriceErrorModal(false)" type="button" class="btn btn-primary" data-dismiss="modal">{{$t('Close')}}</button>
                  </div>
                </div>
              </div>
            </div>

            <form id="inputFormId" novalidate="novalidate" class="h-100">
                <div class="toggle-row row py-2 position-absolute" v-show="infoToggled" v-if="workflow">
                    <name-and-detail-input v-model="workflow" :display-price="true"/>
                </div>
                <div class="toggle-row-backdrop" @click="infoToggled = !infoToggled" v-show="infoToggled"></div>
                <div class="workflow-wrapper">
                    <div class="dynamic-top h-100">
                        <div class="flow-chart-main">
                            <textarea class="flow-chart-data" name="data.flow" style="display: none;"></textarea>
                            <textarea class="progress-flow" name="data.progressFlow" style="display: none;"></textarea>
                            <div class="wfchart-main">
                                <div class="fcf-main" style="width: 100%;">
                                    <div class="flow-chart-finder">
                                        <div>
                                            <div class="w-100 flow-chart-search-box">
                                                <div class="field-parent">
                                                    <form autocomplete="off">
                                                        <div class="input-group w-100 search-box">
                                                            <div class="input-group-prepend">
                                                              <span class="input-group-text" id="basic-addon1" @click="focusInputElementAdd">
                                                                <i class="material-icons">add</i>
                                                              </span>
                                                            </div>
                                                              <input type="search" autocomplete="off" ref="workflowElementAdd"
                                                                     class="form-control text-field flow-chart-finder-input"
                                                                     :placeholder="$t('Search for workflow elements…')"
                                                                     style="box-sizing: border-box;">
                                                            <div class="input-group-append fcf-layout-switch p-1 bg-light border">
                                                                <button type="button"
                                                                        class="btn layout-btn border grid-layout-btn btn-default">
                                                                    <span class="mdi mdi-view-module"></span>
                                                                </button>
                                                                <button type="button"
                                                                        class="btn layout-btn border row-layout-btn btn-primary">
                                                                    <span class="mdi mdi-format-list-bulleted"></span>
                                                                </button>
                                                            </div>
                                                        </div>
                                                    </form>
                                                </div>
                                            </div>
                                            <div style="display: none;" class="flow-chart-finder-result-main">
                                                <div class="flow-chart-finder-cnt">
                                                    <div class="flow-chart-finder-result">
                                                        <div class="fcfr-body"></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="wfchart-wrapper bg-graphpaper border-0"></div>
                                <div class="wfchart-snavi">
                                    <div class="wfc-snn">
                                        <!--<div data-id="collection" class="flow-chart-node">-->
                                        <!--<div class="flow-chart-node-inner">-->
                                        <!--<div cla                            <div class="flow-chart-finder-result">
                                                    <div class="fcfr-body"></div>
                                                </div>ss="fci-node">-->
                                        <!--<i class="node-icon mdi mdi-circle-outline fcn-collection"-->
                                        <!--aria-hidden="true"></i>-->
                                        <!--<div class="flow-chart-finder-simple mt-2">collection</div>-->
                                        <!--</div>-->
                                        <!--</div>-->
                                        <!--</div>-->
                                        <div data-id="condition" class="flow-chart-node flow-chart-node-condition">
                                            <div class="flow-chart-node-inner">
                                                <div class="fci-node">
                                                    <i class="node-icon fcn-condition mdi mdi-rhombus-outline"
                                                       aria-hidden="true"></i>
                                                    <div class="flow-chart-finder-simple">{{$t('condition')}}</div>
                                                </div>
                                            </div>
                                        </div>
                                        <div data-id="placeholder" class="flow-chart-node flow-chart-node-placeholder">
                                            <div class="flow-chart-node-inner">
                                                <div class="fci-node">
                                                    <i class="node-icon fcn-placeholder mdi mdi-hexagon"
                                                       aria-hidden="true"></i>
                                                    <div class="flow-chart-finder-simple">{{$t('placeholder')}}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <ul class='custom-menu'>
                            <li data-action="restrictAccess" for-nodes="template">Restrict Access</li>
                            <li data-action="contextDefineUser" for-nodes="contextSwitch">Define new user</li>
                            <li data-action="connectorFunc" for-nodes="connector">Select remote function</li>
                        </ul>
                    </div>
                </div>
            </form>
            <div id="condition_dialog" class="modal fade" data-backdrop="true"
                 tabindex="-1" role="dialog" aria-labelledby="mySmallModalLabel" aria-hidden="true">
                <div class="modal-dialog modal-lg" role="document">
                    <div class="modal-content">
                        <div class="modal-header bg-light">
                            <h5 class="modal-title">{{$t('Condition')}}</h5>
                            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>
                        <div class="modal-body">
                            <form novalidate="novalidate">
                                <input type="hidden" name="id" value="">
                                <div class="row">
                                    <div class="col-md-6">
                                        <div class="field-parent">
                                            <label for="cond_name" class="ft-label">{{$t('Name')}}</label>
                                            <input type="text"
                                                   id="cond_name"
                                                   :placeholder="$t('Name')"
                                                   class="form-control before-info-text"
                                                   name="condition.name">
                                        </div>
                                    </div>

                                    <div class="col-md-6">
                                        <div class="field-parent">
                                            <label for="cond_descr" class="ft-label">{{$t('Description')}}</label>
                                            <textarea :placeholder="$t('Description')"
                                                      id="cond_descr"
                                                      class="form-control before-info-text"
                                                      name="condition.detail"></textarea>
                                        </div>
                                    </div>
                                </div>

                                <div class="row" style="margin-bottom:20px;">
                                    <div class="col-md-6">
                                        <span class="ft-label">{{$t('Javascript')}}</span>
                                        <textarea id="condition_jsCode"
                                                  name="condition.jsCode"></textarea>
                                    </div>
                                    <div class="col-md-6">
                                        <span class="ft-label">{{$t('Cases')}}</span>
                                        <textarea id="condition_cases"
                                                  name="condition.cases"></textarea>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div class="modal-footer">
                            <div class="row">
                                <div class="col-md-12">
                                    <label for="condition_cases_parser" style="margin-right:10px;">
                                        {{$t('auto change cases by parsing the script')}}
                                        <input id="condition_cases_parser" type="checkbox" checked>
                                    </label>
                                    <button type="button" class="btn btn-primary"
                                            @click="wfm.conditionDialogHandler.save(this)">{{$t('Save')}}
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div id="placeholder_dialog" class="modal fade" data-backdrop="true"
                 tabindex="-1" role="dialog" aria-labelledby="mySmallModalLabel" aria-hidden="true">
                <div class="modal-dialog modal-lg" role="document">
                    <div class="modal-content">
                        <div class="modal-header bg-light">
                            <h5 class="modal-title">{{$t('Placeholder')}}</h5>
                            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>
                        <div class="modal-body">
                            <form novalidate="novalidate">
                                <input type="hidden" name="id" value="">
                                <div class="row">
                                    <div class="col-md-12">
                                        <div class="field-parent">
                                            <label for="placeholder_description" class="ft-label">{{$t('Description')}}</label>
                                            <textarea :placeholder="$t('Description')"
                                                      id="placeholder_description"
                                                      class="form-control before-info-text"
                                                      name="placeholder.detail"></textarea>
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div class="modal-footer">
                            <div class="row">
                                <div class="col-md-12">
                                    <button type="button" class="btn btn-primary"
                                            @click="wfm.placeholderDialogHandler.save(this)">{{$t('Save')}}
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <b-modal v-model="publishModalShow" class="b-modal"
                 :title="$t('Confirm')"
                 :ok-title="$t('Yes')"
                 :cancel-title="$t('No')"
                 header-bg-variant="light"
                 @hide="publishModalShow = false"
                 @ok="publish">
          <div class="d-block">{{$t('If you publish this workflow anyone can start this workflow.')}}</div>
          <div class="d-block mt-3">{{$t('Are you sure you want to publish this workflow?')}}</div>
          <div class="d-block bold bg-primary">
          </div>
        </b-modal>

        <permission-dialog v-if="workflow" :save="save" :publicLink="app.makeURL('/p/workflow/'+id)"
                           :publicLink2="app.makeURL('/document/'+id)" v-model="workflow" :setup="setupDialog"/>
        <list-item-dialog :_blank="true" :title="$t('Couldn\'t publish')" :setup="setupPublishResponseDialog" :timestamp="false" :noBtns="true" >
            <div class="d-block">{{$t('Some entities couldn\'t be changed.')}}</div>
            <div class="d-block">{{$t('Please review the following entities carefully in order to publish.')}}</div>
        </list-item-dialog>
    </div>
</template>

<script>
'use strict'
import vis from 'vis'
import TopNav from '@/components/layout/TopNav'
import VisFlowchart from '../libs/vis-flowchart.js'
import formChangeAlert from '../mixins/form-change-alert'
import bModal from 'bootstrap-vue/es/components/modal/modal'
import bModalDirective from 'bootstrap-vue/es/directives/modal/modal'
import PermissionDialog from './appDependentComponents/permDialog/PermissionDialog'
import NameAndDetailInput from '../components/NameAndDetailInput'
import mafdc from '@/mixinApp'
import ListItemDialog from './appDependentComponents/ListItemDialog'

window.vis = vis
export default {
  name: 'workflow',
  mixins: [mafdc, formChangeAlert],
  components: {
    ListItemDialog,
    NameAndDetailInput,
    PermissionDialog,
    TopNav,
    // eslint-disable-next-line vue/no-unused-components
    'b-modal': bModal
  },
  directives: {
    'b-modal': bModalDirective
  },
  computed: {
    id () {
      if (this.workflow && this.workflow.id) {
        return this.workflow.id
      }
      try {
        return this.$route.params.id
      } catch (e) {
      }
      return ''
    },
    // eslint-disable-next-line vue/no-async-in-computed-properties
    async docUrl () {
      const url = window.location.href
      return url.replace('admin/workflow', 'document')
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
  data () {
    return {
      infoToggled: false,
      workflow: null,
      wfm: null,
      network: null,
      lastJsonValue: null,
      openPermissionDialog: () => {},
      showPublishResponseDialog: () => {},
      publishModalShow: false
    }
  },
  mounted () {
    document.addEventListener('keydown', this.keyboardHandler)
    axios.get('/api/admin/workflow/' + this.id).then(response => {
      if (!response.data) {
        this.workflow = {}
      } else {
        this.workflow = response.data
      }
      this.snapshot(this.workflow, this.skipFromSnapshot)
      $(document).ready(() => {
        this.init(this.app.amIWriteGrantedFor(this.workflow))
      })
    }, (error) => {
      this.app.handleError(error)
      if (error.response && error.response.status === 404) {
        this.$_error('NotFound', { what: 'Workflow' })
      } else {
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not load Workflow. Please try again or if the error persists contact the platform operator.\n'),
          type: 'error'
        })
        this.$router.push({ to: 'Workflows' })
      }
    })
  },
  beforeDestroy () {
    document.removeEventListener('keydown', this.keyboardHandler)
    if (this.wfm) {
      this.wfm.close()
    }
  },
  methods: {
    focusInputElementAdd () {
      this.$refs.workflowElementAdd.focus()
      this.$refs.workflowElementAdd.click()
    },
    setupPublishResponseDialog (showFunc) {
      this.showPublishResponseDialog = showFunc
    },
    hasUnsavedChanges () {
      this.syncWorkflowWithFlowchart()
      return !this.compare(this.workflow, this.skipFromSnapshot)
    },
    skipFromSnapshot (value, keyOrIndex, obj, fullpath) {
      if (/^.*p\.(y|x)$/m.test(fullpath)) {
        return true// skip
      }
      return false
    },
    keyboardHandler (e) {
      if (e.ctrlKey && (e.which === 48 || e.keyCode === 48)) {
        // fit on ctrl+0
        this.wfm.fit()
        return false
      }
      if (e.which === 46) {
        this.wfm.Delete()
        e.preventDefault()
        return false
      } else if (e.ctrlKey && e.which === 83) { // Check for the Ctrl key being pressed, and if the key = [S] (83)
        this.save()
        e.preventDefault()
        return false
      }
    },
    setupDialog (openPermissionDialog) {
      this.openPermissionDialog = openPermissionDialog
    },
    syncWorkflowWithFlowchart () {
      this.getWFlow()
      if (this.wfm) {
        this.workflow.data.flow = this.wfm.getData()
      }
    },
    getWFlow () {
      if (!this.workflow) {
        this.workflow = { data: {} }
      }
      if (!this.workflow.data) {
        this.workflow.data = {}
      }
      return this.workflow.data.flow
    },
    publish () {
      this.workflow.published = true

      this.infoToggled = false
      this.syncWorkflowWithFlowchart()

      const workflowClone = Object.assign({}, this.workflow)
      if (!workflowClone.price) {
        workflowClone.price = 0
      }
      axios.post('/api/admin/workflow/update?publish&id=' + this.id, workflowClone).then(response => {
        this.snapshot(this.workflow, this.skipFromSnapshot)
        if (response.status === 207) {
          // show collected errors
          if (response.data) {
            const elements = []
            for (const key in response.data) {
              if (Object.prototype.hasOwnProperty.call(response.data, key)) {
                const item = response.data[key].Item
                if (item) {
                  item.getLink = function () {
                    if (this.id && this.type && (/workflow|form|template/.test(this.type))) {
                      return '/admin/' + this.type + '/' + this.id
                    }
                  }
                  if (item.type === 'workflow') {
                    item.iconFa = 'mdi mdi-source-branch'
                  }
                  if (item.type === 'form') {
                    item.icon = 'view_quilt'
                  }
                  if (item.type === 'template') {
                    item.iconFa = 'mdi mdi-code-block-tags'
                  }
                  item.error = response.data[key].Error
                  elements.push(item)
                }
              }
            }
            if (elements.length) {
              this.showPublishResponseDialog(null, elements)
            }
          }
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Success'),
            text: this.$t('Workflow published'),
            type: 'success'
          })
        }
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not publish workflow. Please try again or if the error persists contact the platform operator.\n'),
          type: 'error'
        })
      })
    },
    togglePriceErrorModal (show) {
      if (show === true) {
        $(this.$refs.priceErrorModal).show()
        return
      }
      $(this.$refs.priceErrorModal).hide()
    },
    save () {
      if (!this.app.amIWriteGrantedFor(this.workflow)) {
        return
      }
      this.infoToggled = false
      this.syncWorkflowWithFlowchart()
      const workflowClone = Object.assign({}, this.workflow)
      if (!workflowClone.price) {
        workflowClone.price = 0
      }
      if (workflowClone.price % 1 !== 0) {
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Decimal numbers are not allowed for the XES Price'),
          type: 'error'
        })
        return
      }

      axios.post('/api/admin/workflow/update?id=' + this.id, workflowClone).then(response => {
        this.snapshot(this.workflow, this.skipFromSnapshot)
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Workflow saved'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        if (err.response.data === 'can not set price without eth addr') {
          this.togglePriceErrorModal(true)
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not save workflow. Please try again or if the error persists contact the platform operator.\n'),
            type: 'error'
          })
        }
      })
    },
    init (writeRights) {
      const self = this
      const myWFM = {
        urlPrefix: '/api/admin',
        frontendUrlPrefix: '/admin',
        initUsedNodesMap: { workflow: { [self.id]: true } },
        usedNodesMap: null,
        $main: null,
        $chartMain: null,
        $finderResultMain: null,
        $finderResult: null,
        localLayoutSettings: null,
        nodeData: null,
        leftMenuNodes: {
          collection: {
            detail: '',
            type: 'collection',
            name: 'collection'
          },
          condition: {
            detail: '',
            type: 'condition',
            name: self.$t('condition')
          },
          placeholder: {
            detail: '',
            type: 'placeholder',
            name: self.$t('placeholder')
          }
        },
        lastSearchTime: 2,
        lastNodeDeleteTime: 1,
        resetUsedNodesMap: function () {
          this.usedNodesMap = $.extend({}, this.initUsedNodesMap)
          return this.usedNodesMap
        },
        reload: function (flowData) {
          this.drawWf(flowData)
        },
        getData: function () {
          if (this.network) {
            return this.network.getFlowData()
          }
        },
        Delete: function () {
          this.network.Delete()
        },
        close: function () {
          window.removeEventListener('resize', this.heightFixer)
        },
        initConditionDialogHandler: function () {
          this.conditionDialogHandler = {
            $dialog: $('#condition_dialog'),
            $accessModal: $('#access_dialog'),
            $cases: $('#condition_cases'),
            $autoParseCases: null,
            init: true,
            casesEditor: null,
            jsEditor: null,
            currentNode: null,
            callback: null,
            saveEvent: function (env, args, request) {
              this.save(this)
            },
            onInputEvent: function (a, b, c) {
              try {
                var jsonStr = b.getValue()
                this.lastJsonValue = $.parseJSON(jsonStr)
              } catch (parseException) {
              }
            },
            open: function (node, callback, isNewNode) {
              this.isNewNode = isNewNode
              this.currentNode = node
              this.callback = callback
              if (this.init) {
                this.$autoParseCases = this.$dialog.find('#condition_cases_parser')
                // ace.require("ace/ext/language_tools");
                this.jsEditor = ace.edit('condition_jsCode')
                this.jsEditor.session.setMode('ace/mode/javascript')
                this.jsEditor.setTheme('ace/theme/chrome')
                // jsEditor.setAutoScrollEditorIntoView(true);
                // enable autocompletion and snippets
                this.jsEditor.setOptions({
                  maxLines: 400,
                  //                    autoScrollEditorIntoView: true,
                  enableBasicAutocompletion: false,
                  enableSnippets: true,
                  wrapBehavioursEnabled: true,
                  // autoScrollEditorIntoView: true,
                  enableLiveAutocompletion: true
                })
                this.jsEditor.setShowPrintMargin(false)
                this.jsEditor.commands.addCommand({
                  name: 'saveFile',
                  bindKey: {
                    win: 'Ctrl-S',
                    mac: 'Command-S',
                    sender: 'editor|cli'
                  },
                  exec: this.saveEvent
                })
                this.jsEditor.on('input', this.onInputEvent)
                this.jsEditor.renderer.on('afterRender', function () {
                })
                // jsEditor.renderer.setScrollMargin(10, 10, 10, 10);
                // jsEditor.getSession().setUseWrapMode(false);
                this.jsEditor.$blockScrolling = Infinity
                this.jsEditor.getSession().setValue('', -1)

                this.casesEditor = ace.edit('condition_cases')
                this.casesEditor.session.setMode('ace/mode/javascript')
                this.casesEditor.setTheme('ace/theme/chrome')
                // casesEditor.setAutoScrollEditorIntoView(true);
                // enable autocompletion and snippets
                this.casesEditor.setOptions({
                  maxLines: 400,
                  //                    autoScrollEditorIntoView: true,
                  enableBasicAutocompletion: false,
                  enableSnippets: true,
                  wrapBehavioursEnabled: true,
                  // autoScrollEditorIntoView: true,
                  enableLiveAutocompletion: true
                })
                this.casesEditor.setShowPrintMargin(false)
                this.casesEditor.commands.addCommand({
                  name: 'saveFile',
                  bindKey: {
                    win: 'Ctrl-S',
                    mac: 'Command-S',
                    sender: 'editor|cli'
                  },
                  exec: this.saveEvent
                })
                this.casesEditor.on('input', this.onInputEvent)
                this.casesEditor.renderer.on('afterRender', function () {
                })
                // casesEditor.renderer.setScrollMargin(10, 10, 10, 10);
                // casesEditor.getSession().setUseWrapMode(false);
                this.casesEditor.$blockScrolling = Infinity
                this.casesEditor.getSession().setValue('', -1)

                this.init = false
              }
              var js = `
function condition(){
  if( input["someVar"] == "someValue" ){
    /*must match the value in the cases*/
    return "someValue";
  }else{
    /*must match the value in the cases*/
    return "something else";
  }
}
                                        `
              if (node._data && node._data.js) {
                js = node._data.js
              }
              this.jsEditor.getSession().setValue(js, -1)
              var cases = `
[
  {
    "name": "Some value explanation",
    "value": "someValue"
  },
  {
    "name": "Something else explanation",
    "value": "something else"
  }
]
                                        `
              if (node._cases) {
                cases = JSON.stringify(node._cases, null, '\t')
              }
              this.casesEditor.getSession().setValue(cases, -1)
              if (!node.label) {
                node.label = node.name
              }
              this.$dialog.find('[name=\'condition.name\']').val(node.label)
              this.$dialog.find('[name=\'condition.detail\']').val(node.detail)
              this.$dialog.modal({ show: true, backdrop: true })
            },
            save: function (el) {
              try {
                var js = this.jsEditor.getValue()
                if (this.$autoParseCases.is(':checked')) {
                  this.currentNode._cases = parseFunctionCases(js, this.casesEditor.getValue())
                } else {
                  try {
                    this.currentNode._cases = JSON.parse(this.casesEditor.getValue())
                  } catch (eee) {
                    this.currentNode._cases = []
                  }
                }

                if (this.currentNode._cases) {
                  this.casesEditor.getSession().setValue(JSON.stringify(this.currentNode._cases, null, '\t'), -1)
                } else {
                  this.casesEditor.getSession().setValue('', -1)
                }
                if (!this.currentNode._data) {
                  this.currentNode._data = {}
                }
                this.currentNode._data.js = js
                this.currentNode.label = this.$dialog.find('[name=\'condition.name\']').val()
                this.currentNode.name = this.currentNode.label
                this.currentNode.detail = this.$dialog.find('[name=\'condition.detail\']').val()
                if ($.isFunction(this.callback)) {
                  this.callback(this.currentNode)
                }
                this.$dialog.modal('hide')
              } catch (e) {
                console.log(e)
              }
            }
          }
        },
        initPlaceholderDialogHandler: function () {
          this.placeholderDialogHandler = {
            $dialog: $('#placeholder_dialog'),
            lineSeparator: '\x02\n\x02',
            currentNode: null,
            callback: null,
            open: function (node, callback, isNewNode) {
              this.isNewNode = isNewNode
              this.currentNode = node
              this.callback = callback

              const formatted = node.label?.replaceAll(this.lineSeparator, ' ')

              this.$dialog.find('[name=\'placeholder.detail\']').val(formatted)
              this.$dialog.modal({ show: true, backdrop: true })
            },
            save: function () {
              try {
                const detail = this.$dialog.find('[name=\'placeholder.detail\']').val()

                const symbolsPerLine = 20

                let formatted = ''
                let pointer = 0
                let split = false

                for (const symbol of detail) {
                  if (symbol === '\n') {
                    pointer = 0
                  }

                  if (pointer >= symbolsPerLine) {
                    split = true
                  } else {
                    split = false
                  }

                  if (split && symbol === ' ') {
                    formatted = `${formatted}${this.lineSeparator}`
                    pointer = 0
                  } else {
                    formatted = `${formatted}${symbol}`
                  }

                  pointer++
                }

                this.currentNode.label = formatted
                this.currentNode.name = formatted

                if ($.isFunction(this.callback)) {
                  this.callback(this.currentNode)
                }
                this.$dialog.modal('hide')
              } catch (e) {
                console.log(e)
              }
            }
          }
        },
        readOnly: false,
        main: function () {
          this.readOnly = !writeRights
          this.$main = $('.flow-chart-main')
          var _ = this
          this.$finderResultMain = this.$main.find('.flow-chart-finder-result-main')
          this.$finderResult = this.$finderResultMain.find('.flow-chart-finder-cnt')
          this.$finderResultBody = this.$finderResult.find('.fcfr-body')

          this.$finderResultCloseBtn = this.$finderResultMain.find('.fcfr-close')
          this.$searchInput = this.$main.find('.flow-chart-finder-input')
          this.$searchInput.on('click keyup', function () {
            _.searchNow($(this).val())
          })
          this.$finderResultCloseBtn.click(function () {
            _.hideFinderResult()
          })

          this.$fcfMain = this.$main.find('.wfchart-main .fcf-main')
          this.$legendBtn = this.$main.find('.wfchart-main .btn-legend')
          this.$legendBtn.click(function () {
            _.showLegend()
            return false
          })
          var gridBtn = this.$main.find('.flow-chart-finder .fcf-layout-switch .grid-layout-btn')
          var rowBtn = this.$main.find('.flow-chart-finder .fcf-layout-switch .row-layout-btn')
          gridBtn.click(function () {
            rowBtn.removeClass('btn-primary').addClass('btn-default')
            gridBtn.removeClass('btn-default').addClass('btn-primary')
            _.$finderResult.removeClass('fcf-row').addClass('fcf-grid')
            _.layoutSwitch('grid')
          })
          rowBtn.click(function () {
            gridBtn.removeClass('btn-primary').addClass('btn-default')
            rowBtn.removeClass('btn-default').addClass('btn-primary')
            _.$finderResult.removeClass('fcf-grid').addClass('fcf-row')
            _.layoutSwitch('row')
          })
          try {
            _.localLayoutSettings = $.parseJSON(Cookies.get(window.location.pathname))
          } catch (e) {
            _.localLayoutSettings = { layout: 'row' }
          }
          if (_.localLayoutSettings.layout === 'row') {
            rowBtn.click()
          } else {
            gridBtn.click()
          }
          this.tmpl.finderCompiled = Handlebars.compile(this.tmpl.finder)
          this.setupDragAndDrop()
          this.$chartMain = this.$main.find('.wfchart-wrapper')
          this.heightFixer = function () {
            if (_.$chartMain) {
              if (this.minusHeight === undefined) {
                this.minusHeight = 110
              }
              if (this.minusFinderHeight === undefined) {
                this.minusFinderHeight = _.$fcfMain.height()
              }
              let minusHeight = this.minusHeight
              if (_.readOnly) {
                _.$fcfMain.hide()
                minusHeight = minusHeight - this.minusFinderHeight
              } else {
                _.$fcfMain.show()
              }
              _.$chartMain.css('height', ($(document.body).height() - minusHeight - 2) + 'px')
            }
          }
          this.heightFixer()
          window.addEventListener('resize', this.heightFixer)
          this.attachMouseOffsetWorkaroundListener()
          this.$main.on('click mousedown', function (e) {
            _.hideLegend()
            _.handleContextMenu(e)
          })
          this.$chartMain.on('click', function () {
            _.hideFinderResult()
            return true
          })
          this.$main.find('.wfchart-snavi').find('.fci-node:not(.disabled)').each(function () {
            if (_.readOnly) {
              return
            }
            _.attachDragAndDrop($(this))
          })
          this.initConditionDialogHandler()
          this.initPlaceholderDialogHandler()
        },
        searchNow: function (v) {
          var _ = this
          if (!v) {
            v = ''
          }
          if (_.lastSearchText === v && _.deletedNodesAfterLastSearch() && _.lastSearchUnderThreeMinutes() &&
              _.finderResultNodesCount() > 0) {
            _.showFinderResult()
          } else {
            _.$finderResultBody.empty()
            _.nodeData = {}
            _.lastSearchTime = new Date().getTime()
            for (var k in _.nodeIconMap) {
              if (Object.prototype.hasOwnProperty.call(_.nodeIconMap, k)) {
                _.searchNodeSpecific(k, v)
              }
            }
          }
          _.lastSearchText = v
        },
        deletedNodesAfterLastSearch: function () {
          var b = this.lastNodeDeleteTime < this.lastSearchTime
          return b
        },
        lastSearchUnderThreeMinutes: function () {
          var b = new Date().getTime() < (this.lastSearchTime + (1000 * 60 * 3))
          return b
        },
        blacklisted: {},
        searchNodeSpecific: function (nodeKind, search) {
          if (nodeKind === 'user' || nodeKind === 'condition' || nodeKind === 'collection') {
            return
          }
          var _ = this
          if (_.blacklisted[nodeKind]) {
            return
          }
          $.ajax({
            myCustomObj: { kind: nodeKind },
            type: 'GET',
            data: {
              c: search,
              i: 0,
              l: 20,
              e: JSON.stringify(myWFM._makeGetParamOfExcludes(nodeKind)),
              q: JSON.stringify({ metaOnly: true })
            },
            url: '/api/admin/' + nodeKind + '/list',
            error: function (res) {
              if (res && res.status && res.status === 401) {
                window.redirectToLogin()
              }
            },
            success: function (data) {
              if (data && data.length) {
                for (var i = 0; i < data.length; ++i) {
                  if (data[i] && data[i].name) {
                    myWFM.createFinderItem({
                      name: data[i].name,
                      detail: data[i].detail ? data[i].detail : data[i].description,
                      kind: this.myCustomObj.kind,
                      data: data[i]
                    })
                  }
                }
                this.myCustomObj = null
                if (!myWFM.hideFinderResultIfEmpty()) {
                  myWFM.showFinderResult()
                }
              }
            }
          })
        },
        _makeGetParamOfExcludes: function (t) {
          if (this.usedNodesMap && t && this.usedNodesMap[t]) {
            return this.usedNodesMap[t]
          }
          return {}
        },
        layoutSwitch: function (layout) {
          if (this.localLayoutSettings.layout !== layout) {
            this.localLayoutSettings.layout = layout
            this.storeSettings()
          }
        },
        showLegend: function () {
          if (!this.legendVisible) {
            this.$legendBtn.blur()
            var chartMain = this.$main.find('.wfchart-main')
            if (!this.$legend) {
              var rows = ''
              for (var k in this.nodeIconMap) {
                if (Object.prototype.hasOwnProperty.call(this.nodeIconMap, k)) {
                  rows += this._createFinderItem({ name: k, kind: k })
                }
              }
              this.$legend = $('<div class="fcl fcl-legend flow-chart-finder-cnt fcf-row" style="display:none;">' +
                  rows + '</div>')
              chartMain.append(this.$legend)
            }
            this.$legend.fadeIn('fast')
            this.legendVisible = true
          }
        },
        hideLegend: function () {
          if (this.legendVisible && this.$legend) {
            this.$legendBtn.blur()
            this.legendVisible = false
            this.$legend.fadeOut('fast')
          }
        },
        handleContextMenu: function (e) {
          // If the clicked element is not the menu
          if (!$(e.target).parents('.custom-menu').length > 0) {
            // Hide it
            $('.custom-menu').hide(100)
          }
        },
        storeSettings: function () {
          try {
            Cookies.set(window.location.pathname, JSON.stringify(this.localLayoutSettings))
          } catch (someErrorWhenStoringCookies) {
            console.log(someErrorWhenStoringCookies)
          }
        },
        hideFinderResultIfEmpty: function () {
          if (this.finderResultVisible && this.finderResultNodesCount() === 0) {
            this.hideFinderResult()
            return true
          }
          return false
        },
        hideFinderResult: function () {
          if (this.finderResultVisible) {
            this.$finderResultMain.fadeOut('fast')
            this.finderResultVisible = false
          }
        },
        finderResultVisible: false,
        showFinderResult: function () {
          if (!this.finderResultVisible) {
            if (this.beforeShowFinder) {
              this.beforeShowFinder()
            }
            if (this.finderResultNodesCount() > 0) {
              this.$finderResultMain.fadeIn('fast')
              this.finderResultVisible = true
              return true
            }
          }
          return false
        },
        finderResultNodesCount: function () {
          return this.$finderResultBody.children().length
        },
        onNodeInsert: function (data, position) {
          if (data) {
            data.x = position.x
            data.y = position.y
            this.network.addNode(data)
          }
        },
        attachMouseOffsetWorkaroundListener: function () {
          var _ = this
          this.workaroundListener = function (e) {
            if (!_.dragActive) {
              if (_.$dragComp && _.$dragComp.length && _.$dragComp.data('dropped')) {
                if (_.cloned) {
                  _.onNodeInsert(_.nodeData[_.$newTarget.attr('data-id')], {
                    x: e.offsetX,
                    y: e.offsetY
                  })
                  delete _.nodeData[_.$newTarget.attr('data-id')]
                  _.$newTarget.remove()
                } else {
                  if (_.$newTarget.attr('data-id') === 'condition') {
                    _.conditionDialogHandler.open($.extend({}, _.leftMenuNodes[_.$newTarget.attr('data-id')]),
                      function (n) {
                        n.cases = n._cases
                        n.data = n._data
                        _.onNodeInsert(n, { x: e.offsetX, y: e.offsetY })
                      }, true)
                  } else if (_.$newTarget.attr('data-id') === 'placeholder') {
                    _.placeholderDialogHandler.open($.extend({}, _.leftMenuNodes[_.$newTarget.attr('data-id')]),
                      function (n) {
                        n.data = n._data
                        _.onNodeInsert(n, { x: e.offsetX, y: e.offsetY })
                      }, true)
                  } else {
                    _.onNodeInsert(_.leftMenuNodes[_.$newTarget.attr('data-id')], {
                      x: e.offsetX,
                      y: e.offsetY
                    })
                  }
                }
                _.$dragComp.data('dropped', false)
                setTimeout(function () {
                  _.hideFinderResultIfEmpty()
                }, 800)
              }
            }
          }
          this.$chartMain.on('mouseover touchmove touchend', this.workaroundListener)
        },
        droppedOverDZ: false,
        setupDragAndDrop: function () {
          var _ = this
          _.dim = {}
          _.tl = {}
          _.mousePos = {}
          _.dragCompPosition = {}
          interact('.flow-chart-main .wfchart-wrapper').dropzone({
            accept: '.fci-node',
            overlap: 'pointer',
            ondropactivate: function (event) {
              _.$dropzone = $(event.target)
              _.$dropzone.addClass('drag-active')
            },
            ondragenter: function (event) {
              //                        _.$dropzone = $(event.target);
            },
            ondragleave: function (event) {

              // if(!this.outOfMovingZone){
              //                        _._checkForLeave();
              // }
            },
            ondrop: function (e, a, b, c, d) {
              _.droppedOverDZ = true
            },
            ondropdeactivate: function (event) {
              var $dropzone = $(event.target)
              $dropzone.removeClass('hcb-changed3')
              $dropzone.removeClass('drag-active')
            }
          })
        },
        attachDragAndDrop: function ($el) {
          var _ = this
          interact($el[0]).styleCursor(false).draggable({
            // enable inertial throwing
            inertia: false,
            onstart: function (event) {
              if (event?.interaction?.downEvent?.button === 2) {
                return false
              }
              _.dragActive = true
              _.$dragComp = $(event.target)
              _.$dragComp.data('dropped', false)
              _.$newTarget = _.$dragComp.parents('.flow-chart-node')
              var snn = _.$newTarget.parents('.wfc-snn')
              if (snn && snn.length) {
                _.cloned = false
                _.$clone = null
              } else {
                _.cloned = true
                var curPos = _.$newTarget.position()
                _.$clone = _.$newTarget.clone()
                _.$clone.find('.fci-node').css('visibility', 'hidden')
                _.$newTarget.replaceWith(_.$clone)
                _.$newTarget.css({
                  position: 'absolute',
                  left: (curPos.left) + 'px',
                  top: (curPos.top) + 'px'
                })
                _.$finderResultMain.append(_.$newTarget)
              }

              _.dim.mt = parseFloat(_.$dragComp.css('marginTop'))
              _.dim.ml = parseFloat(_.$dragComp.css('marginLeft'))
              _.dim.h = _.$dragComp.outerHeight(true)
              _.dim.w = _.$dragComp.outerWidth(true)
              _.dim.h2 = _.dim.h / 2
              _.dim.w2 = _.dim.w / 2
              if (event.interaction.startOffset) {
                _.dim.handlePercentX = event.interaction.startOffset.left + _.dim.ml
                _.dim.handlePercentX = 100 / _.dim.w * _.dim.handlePercentX
                _.dim.handlePercentY = event.interaction.startOffset.top + _.dim.mt
                _.dim.handlePercentY = 100 / _.dim.h * _.dim.handlePercentY
              }
              _.dim.handleX = _.dim.handlePercentX * _.dim.w / 100
              _.dim.handleY = _.dim.handlePercentY * _.dim.h / 100

              _.tl.x = 0; _.tl.y = 0
              _.$dragComp.addClass('dragging')

              _.bcr = _.$dragComp[0].getBoundingClientRect()
              _.dragCompPosition.top = (_.bcr.top) - _.dim.mt
              _.dragCompPosition.left = (_.bcr.left) - _.dim.ml
            },
            onmove: function (event, a, b, c, d) {
              if (event?.interaction?.downEvent?.button === 2) {
                return false
              }
              _.mousePos.y = event.pageY
              _.mousePos.x = event.pageX
              _.tl.x += event.dx
              _.tl.y += event.dy
              const translatePos = 'translate(' + _.tl.x + 'px, ' + _.tl.y + 'px)'
              _.$dragComp[0].style.webkitTransform = translatePos
              _.$dragComp[0].setAttribute('style', 'transform: ' + translatePos)
            },
            onend: function (e, a, b, c, d) {
              _.wasInside = false
              _.dragActive = false
              var revert = function () {
                _.$dragComp.removeClass('dragging')
                if (_.localLayoutSettings.layout === 'row') {
                  _.$dragComp[0].style.webkitTransform = _.$dragComp[0].style.transform = 'translate(1px, -3px)'
                } else {
                  _.$dragComp[0].style.webkitTransform = _.$dragComp[0].style.transform = 'translate(2px, 2px)'
                }
                // revert after animation
                setTimeout(function () {
                  if (_.cloned) {
                    _.$newTarget.detach()
                  }

                  _.$newTarget.css({ position: '', left: '', top: '' })
                  _.$dragComp.css('transform', '')
                  if (_.cloned) {
                    _.$clone.replaceWith(_.$newTarget)
                  } else {
                    _.$dragComp.fadeIn('fast')
                  }
                }, 10)
              }
              if (_.droppedOverDZ) {
                var rect = _.$finderResult[0].getBoundingClientRect()
                window.droppedOverDZ = false
                if (rect.bottom < e.clientY) {
                  _.$dragComp.data('dropped', true)
                  var drect = _.$dragComp[0].getBoundingClientRect()
                  var isMobileBrowser = window.isMobileBrowser()
                  if (isMobileBrowser) {
                    _.$searchInput.blur()
                  }
                  // remove from list and add to chart
                  _.$dragComp.fadeOut('fast', function () {
                    if (_.cloned) {
                      _.$newTarget.hide()
                    } else {
                      revert()
                    }
                  })
                  if (_.cloned) {
                    var thisClone = _.$clone
                    thisClone.hide('fast', function () {
                      thisClone.remove()
                    })
                  }
                  if (isMobileBrowser) {
                    var ee = new Event('touchend')
                    ee.offsetX = drect.left
                    ee.offsetY = drect.top
                    ee.pageX = ee.clientX = ee.offsetX
                    ee.pageY = ee.clientY = ee.offsetY
                    _.$chartMain[0].dispatchEvent(ee)
                  }

                  // _.workaroundListener(e)
                  return false
                } else {
                  revert()
                }
              } else {
                revert()
              }
            }
          })
        },
        tmpl: {
          finder:
              '<div data-id="{{id}}" class="flow-chart-node">' +
              '    <div class="flow-chart-node-inner">' +
              '    <div class="fci-node">' +
              '       <i class="{{iconClass}}" aria-hidden="true"></i>' +
              '       <table>' +
              '           <tbody>' +
              '               <tr><td><div class="flow-chart-finder-simple">' +
              '                   {{name}}' +
              '               </div></td></tr>' +
              '               <tr><td><div class="flow-chart-finder-detail">' +
              '                   {{detail}}' +
              '               </div></td></tr>' +
              '           </tbody>' +
              '       </table>' +
              '    </div>' +
              '    </div>' +
              '</div>'
        },
        nodeIconMap: {
          ibmsender: 'fcn-ibmsender node-icon mdi mdi-send',
          mailsender: 'fcn-externalnode node-icon mdi mdi-send',
          priceretriever: 'fcn-externalnode node-icon mdi mdi-send',
          externalNode: 'fcn-externalnode node-icon mdi mdi-link-variant',
          condition: 'fcn-condition node-icon mdi mdi-circle-outline',
          placeholder: 'fcn-placeholder node-icon mdi mdi-hexagon',
          user: 'fcn-usr node-icon mdi mdi-account',
          form: 'fcn-form node-icon mdi mdi-view-quilt',
          workflow: 'fcn-wflow node-icon mdi mdi-source-branch',
          template: 'fcn-tmpl node-icon mdi mdi-code-block-tags'
        },
        // d={name:"name", detail:"detail", kind:"condition|user|form|workflow|template"};
        _createFinderItem: function (d) {
          d.iconClass = this.nodeIconMap[d.kind]
          return this.tmpl.finderCompiled(d)
        },
        createFinderItem: function (d) {
          if (d && d.data && d.data.id) {
            d.id = d.kind + '_' + d.data.id
            var $elExists = this.$finderResultBody.find('div[data-id="' + d.id + '"]')
            if ($elExists.length > 0) {
              return
            }
            if (d.kind === 'condition' && d.data.data) {
              if (d.data.data.cases) {
                if (typeof d.data.data.cases === 'string') {
                  d.data.cases = JSON.parse(d.data.data.cases)
                } else {
                  d.data.cases = d.data.data.cases
                }
              }
            }
            d.data.type = d.kind

            if (!d.detail) {
              d.detail = '-'
            }
            if (!this.nodeData) {
              this.nodeData = {}
            }
            this.nodeData[d.id] = d.data
            var $el = $(this._createFinderItem(d))
            this.attachDragAndDrop($el.find('.fci-node'))
            this.$finderResultBody.append($el)
            return $el
          }
        },
        onDblClick: function (n) {
          console.log(this)
          var url = window.location.origin + this.wfm.frontendUrlPrefix + '/' + n.group + '/' + n.id
          window.open(url, '_blank')
        },
        onDblClickWithName: function (n) {
          var url = window.location.origin + this.wfm.frontendUrlPrefix + '/' + n.group + '/' + n.label + '/' + n.id
          window.open(url, '_blank')
        },
        conditionDblClick: function (node) {
          console.log(this)
          var _ = this
          this.wfm.conditionDialogHandler.open(node, function (n) {
            _.updateNode(n)
          })
        },
        placeholderDblClick: function (node) {
          var _ = this
          this.wfm.placeholderDialogHandler.open(node, function (n) {
            _.updateNode(n)
          })
        },
        destroy: function () {
          if (this.network) {
            this.network.destroy()
            this.network = null
          }
        },
        fit: function () {
          this.network.fit()
        },
        drawWf: function (wfData) {
          const _ = this
          this.destroy()
          this.resetUsedNodesMap()
          this.network = new VisFlowchart(wfData, {
            html: {
              $wrapper: _.$main,
              $container: _.$chartMain
            },
            showIDOnTooltip: true,
            readOnly: _.readOnly,
            nodes: {
              start: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#45bbff',
                          highlight: { background: '#8dd5ff' },
                          hover: { background: '#8dd5ff' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: { color: '#45bbff', highlight: '#45bbff', hover: '#45bbff' }
                      }
                    }
                  ]
                },
                icon: {
                  face: 'Material Icons',
                  code: 'radio_button_unchecked',
                  size: 80,
                  color: '#45bbff'
                }
              },
              collection: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#23d6d6',
                          highlight: { background: '#23d6d6' },
                          hover: { background: '#23d6d6' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: { color: '#23d6d6', highlight: '#23d6d6', hover: '#23d6d6' }
                      }
                    }
                  ],
                  fromInfinity: true,
                  to: Infinity
                },
                font: {
                  color: '#343434',
                  size: 15,
                  mod: 'bold'
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '\uF765',
                  color: '#1fa6a6'
                }
              },
              form: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#0eaa64',
                          highlight: { background: '#99d29a' },
                          hover: { background: '#99d29a' }
                        },
                        borderWidthSelected: 1
                      },
                      edge: {
                        color: { color: '#0eaa64', highlight: '#3c763d', hover: '#3c763d' }
                      }
                    }],
                  to: Infinity,
                  space: 1.1
                },
                borderWidth: 10,
                borderWidthSelected: 1,
                font: {
                  color: '#343434',
                  size: 15,
                  mod: 'bold',
                  bold: {
                    color: '#343434',
                    size: 14, // px
                    face: 'arial',
                    vadjust: 0,
                    mod: 'bold'
                  }
                },
                icon: {
                  face: 'Material Icons',
                  code: 'view_quilt',
                  color: '#0eaa64'
                },
                events: {
                  hoverIn: function () {
                  },
                  hoverOut: function () {
                  },
                  remove: function () {
                  },
                  click: function () {
                  },
                  dblclick: _.onDblClick
                }
              },
              ibmsender: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#8688ff',
                          highlight: { background: '#5f5ff0' },
                          hover: { background: '#5f5ff0' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: {
                          color: '#8688ff',
                          highlight: '#a8a5ff',
                          hover: '#a8a5ff'
                        }
                      }
                    }],
                  to: Infinity,
                  space: 1.1
                },
                font: {
                  color: '#343434',
                  size: 15,
                  mod: 'bold',
                  bold: {
                    color: '#343434',
                    size: 14, // px
                    face: 'arial',
                    vadjust: 0,
                    mod: 'bold'
                  }
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '\uf48a',
                  color: '#5353c0'
                },
                events: {
                  hoverIn: function () {
                  },
                  hoverOut: function () {
                  },
                  remove: function () {
                  },
                  click: function () {
                  },
                  dblclick: function () {
                  }
                }
              },
              mailsender: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#8688ff',
                          highlight: { background: '#5f5ff0' },
                          hover: { background: '#5f5ff0' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: {
                          color: '#8688ff',
                          highlight: '#a8a5ff',
                          hover: '#a8a5ff'
                        }
                      }
                    }],
                  to: Infinity,
                  space: 1.1
                },
                font: {
                  color: '#343434',
                  size: 15,
                  mod: 'bold',
                  bold: {
                    color: '#343434',
                    size: 14, // px
                    face: 'arial',
                    vadjust: 0,
                    mod: 'bold'
                  }
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '󰒊',
                  color: '#5353c0'
                },
                events: {
                  hoverIn: function () {
                  },
                  hoverOut: function () {
                  },
                  remove: function () {
                  },
                  click: function () {
                  },
                  dblclick: function () {
                  }
                }
              },
              priceretriever: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#8688ff',
                          highlight: { background: '#5f5ff0' },
                          hover: { background: '#5f5ff0' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: {
                          color: '#8688ff',
                          highlight: '#a8a5ff',
                          hover: '#a8a5ff'
                        }
                      }
                    }],
                  to: Infinity,
                  space: 1.1
                },
                font: {
                  color: '#343434',
                  size: 15,
                  mod: 'bold',
                  bold: {
                    color: '#343434',
                    size: 14, // px
                    face: 'arial',
                    vadjust: 0,
                    mod: 'bold'
                  }
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '󰒊',
                  color: '#5150c0'
                },
                events: {
                  hoverIn: function () {
                  },
                  hoverOut: function () {
                  },
                  remove: function () {
                  },
                  click: function () {
                  },
                  dblclick: function () {
                  }
                }
              },
              externalNode: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#8688ff',
                          highlight: { background: '#5f5ff0' },
                          hover: { background: '#5f5ff0' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: {
                          color: '#8688ff',
                          highlight: '#a8a5ff',
                          hover: '#a8a5ff'
                        }
                      }
                    }],
                  to: Infinity,
                  space: 1.1
                },
                font: {
                  color: '#343434',
                  size: 15,
                  mod: 'bold',
                  bold: {
                    color: '#343434',
                    size: 14, // px
                    face: 'arial',
                    vadjust: 0,
                    mod: 'bold'
                  }
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '󰌹',
                  color: '#5150c0'
                },
                events: {
                  dblclick: _.onDblClickWithName
                }
              },
              condition: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#f4a646',
                          highlight: { background: '#f9cd95' },
                          hover: { background: '#f9cd95' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        font: { align: 'middle' },
                        dashes: true,
                        arrows: 'to',
                        color: { color: '#f4a646', highlight: '#e88000', hover: '#e88000' }
                      }
                    }],
                  to: Infinity,
                  space: 0.6
                },
                font: {
                  color: '#343434',
                  size: 15
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '󰜌',
                  color: '#f0a30a'
                },
                events: {
                  dblclick: _.conditionDblClick
                }
              },
              placeholder: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#f5e203',
                          highlight: { background: '#f5e203' },
                          hover: { background: '#f5e203' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        font: { align: 'middle' },
                        arrows: 'to',
                        color: { color: '#f5e203', highlight: '#f5e203', hover: '#f5e203' }
                      }
                    }],
                  to: Infinity,
                  space: 0.6
                },
                font: {
                  color: '#343434',
                  size: 15
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '\u2B22',
                  color: '#f5e203'
                },
                events: {
                  dblclick: _.placeholderDblClick
                }
              },
              workflow: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#e40070',
                          highlight: { background: '#f3a3ca' },
                          hover: { background: '#f3a3ca' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: { color: '#e40070', highlight: '#d60069', hover: '#d60069' }
                      }
                    }],
                  to: Infinity
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '󰘬',
                  color: '#e40070'
                },
                events: {
                  dblclick: _.onDblClick
                }
              },
              user: {
                connections: {
                  from: [
                    {
                      node: {
                        color: {
                          background: '#7200ff',
                          highlight: { background: '#c393ff' },
                          hover: { background: '#c393ff' }
                        },
                        borderWidthSelected: 3
                      },
                      edge: {
                        color: { color: '#7200ff', highlight: '#5b00cc', hover: '#5b00cc' }
                      }
                    }],
                  to: Infinity
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '\uf004',
                  color: '#7200ff'
                },
                events: {
                  dblclick: function () {
                  }
                }
              },
              template: {
                connections: {
                  from: [
                    {
                      node: {
                        shape: 'dot',
                        size: 8,
                        color: {
                          background: '#ff30ec',
                          border: '#7b7b7b',
                          highlight: { background: '#ffabf7', border: '#7b7b7b' },
                          hover: { background: '#ffabf7', border: '#7b7b7b' }
                        },
                        borderWidth: 1,
                        borderWidthSelected: 5
                      },
                      edge: {
                        color: { color: '#ff30ec', highlight: '#e22cd1', hover: '#e22cd1' }
                      }
                    }],
                  to: Infinity
                },
                font: {
                  color: '#343434',
                  size: 15
                },
                icon: {
                  face: 'Material Design Icons',
                  code: '󱲆',
                  color: '#ff30ec'
                },
                events: {
                  dblclick: _.onDblClick
                }
              }
            },
            events: {
              insert: function (n) {
                var g = n.group
                if (g === 'step') {
                  g = 'form'
                }
                if (!myWFM.usedNodesMap[g]) {
                  myWFM.usedNodesMap[g] = {}
                }
                myWFM.usedNodesMap[g][n.id] = true
              },
              remove: function (n) {
                var g = n.group
                if (g === 'step') {
                  g = 'form'
                }
                if (!myWFM.usedNodesMap[g]) {
                  myWFM.usedNodesMap[g] = {}
                }
                delete myWFM.usedNodesMap[g][n.id]
                myWFM.lastNodeDeleteTime = new Date().getTime()
              }
            },
            node: {
              connections: {
                from: [
                  {
                    node: {
                      shape: 'dot',
                      size: 10,
                      color: {
                        border: '#7b7b7b',
                        highlight: { border: '#7b7b7b' },
                        hover: { border: '#7b7b7b' }
                      },
                      borderWidth: 1,
                      borderWidthSelected: 3
                    },
                    edge: {
                      arrowStrikethrough: false,
                      arrows: 'to',
                      width: 1,
                      hoverWidth: 2
                    }
                  }],
                space: 1.1,
                to: Infinity
              },
              shape: 'icon',
              icon: {
                size: 50
              },
              shadow: {
                enabled: true,
                color: 'rgba(0,0,0,0.1)',
                size: 10,
                x: 5,
                y: 5
              }
            }
          })
          this.network.wfm = this
        }
      }

      self.wfm = myWFM;

      (function () {
        'use strict'
        myWFM.beforeShowFinder = function () {
          if (myWFM.$finderResultMain.find('i.fcn-usr, i.fcn-tmpl, i.fcn-wflow').length === 0) {
            for (var i = 0; i < myWFM.dummyNodes.length; ++i) {
              myWFM.createFinderItem(myWFM.dummyNodes[i])
            }
          }
        }
        myWFM.dummyAddNodeIndex = 0
        myWFM.dummyNodes = []
        myWFM.addNode = function (t, n, d) {
          myWFM.dummyAddNodeIndex++
          var dummy = { i: myWFM.dummyAddNodeIndex, name: n, detail: d, kind: t }
          dummy.data = { id: dummy.kind + dummy.i, name: dummy.name, description: dummy.detail }
          myWFM.dummyNodes.push(dummy)
          myWFM.createFinderItem(dummy)
        }
        myWFM.main()
      })()

      var myRegexMatcher = new RegExp('return\s*((.*));', 'gi')

      function parseFunctionCases (src, casesStr) {
        if (src) {
          var caseValues = src.match(myRegexMatcher)
          if (caseValues) {
            for (var i = 0; i < caseValues.length; ++i) {
              caseValues[i] = caseValues[i].trim()
              if (caseValues[i].startsWith('return')) {
                caseValues[i] = caseValues[i].replace('return', '')
              }
              caseValues[i] = caseValues[i].trim()
              if (caseValues[i].endsWith(';')) {
                caseValues[i] = caseValues[i].replace(';', '')
              }
              caseValues[i] = caseValues[i].trim()
            }
            var val
            var lab
            for (var c = 0; c < caseValues.length; ++c) {
              val = JSON.parse(caseValues[c])
              if (typeof val === 'string') {
                lab = val
              } else {
                lab = caseValues[c]
              }
              caseValues[c] = { name: lab, value: val }
            }
            var targetJson
            if (casesStr || casesStr.trim()) {
              targetJson = caseValues
            } else {
              targetJson = caseValues
            }
            return targetJson
          }
        }
        return null
      }

      myWFM.reload(self.getWFlow())
    }
  }
}

</script>

<style lang="scss">
    @import "../assets/styles/variables";
    @import "~@mdi/font/scss/_variables.scss";
    @import "~@mdi/font/scss/_functions.scss";
    @import "../assets/styles/vis.min.scss";
    @import "../assets/styles/vis-styles.scss";

    .fcl-legend i.node-icon.mdi {
        font-size: 30px;
    }

    .app-workflow {
      .search-box {
        .input-group-text {
          border-radius: 0;
        }
        input.form-control {
          border-width: 1px !important;
          border-radius: 0;
          border-left: 0 !important;
        }
      }
      .topnav {
        margin-bottom: 0!important;
      }
    }

    .app-workflow .toggle-row {
        margin: 0;
    }

    .workflow-wrapper {
        overflow: hidden;
    }

    .vis-network > canvas {
        height: 100%;
    }

    .wfchart-wrapper {
        background-repeat: repeat;
    }

    div.vis-tooltip {
        background: rgba(0, 0, 0, 0.55);
        color: white;
        font-family: arial;
        text-shadow: 1px 1px #000000;
    }

    div.vis-tooltip .wf-node-tooltip .wfnt-name {
        color: inherit;
    }

    div.vis-tooltip .wf-node-tooltip .wfnt-desc {
        color: inherit;
    }

    div.vis-tooltip .wf-node-tooltip .wfnt-label {
        color: inherit;
    }

    div.vis-tooltip .wf-node-tooltip .wfnt-cases {
        color: inherit;
    }

    div.vis-tooltip .wf-node-tooltip .wfnt-case {
        color: inherit;
    }

    $trans: transparent;
    $block: #ffffff;
    $line: #eeeeee;
    $gridSize: 60px;
    $subdivisions: 3;
    $gridSizeOverSubdivisions: 20px;
    $lineAlpha: .6;
    $sublineAlpha: .4;

    .bg-graphpaper {
        background-color: $block;
        background-image: linear-gradient(rgba($line, $sublineAlpha) 1px, $trans 1px), /*sub horiz*/
        linear-gradient($line 1px, $trans 1px), /*main horiz*/
        linear-gradient(90deg, rgba($line, $sublineAlpha) 1px, $trans 1px), /*sub vert*/
        linear-gradient(90deg, rgba($line, $lineAlpha) 1px, $trans 1px), /*main vert*/
        linear-gradient($trans 3px, $block 3px, $block $gridSize - 2, $trans $gridSize - 2), /*nub horiz*/
        linear-gradient(90deg, rgba($line, $lineAlpha) 3px, $trans 3px, $trans $gridSize - 2, rgba($line, $lineAlpha) $gridSize - 2) /*nub vert*/
    ;
        background-size: $gridSizeOverSubdivisions $gridSizeOverSubdivisions,
        $gridSize $gridSize,
        $gridSizeOverSubdivisions $gridSizeOverSubdivisions,
        $gridSize $gridSize,
        $gridSize $gridSize,
        $gridSize $gridSize;
    }

    .vis-zoomExtends:before {
        font-family: "Material Design Icons";
        content: mdi('image-filter-center-focus');
    }

    .vis-zoomIn:before {
        font-family: "Material Design Icons";
        content: mdi('magnify-plus-outline');
    }

    .vis-zoomOut:before {
        font-family: "Material Design Icons";
        content: mdi('magnify-minus-outline');
    }

    .vis-down:before {
        font-family: "Material Design Icons";
        content: mdi('arrow-down');
    }

    .vis-left:before {
        font-family: "Material Design Icons";
        content: mdi('arrow-left');
    }

    .vis-up:before {
        font-family: "Material Design Icons";
        content: mdi('arrow-up');
    }

    .vis-right:before {
        font-family: "Material Design Icons";
        content: mdi('arrow-right');
    }

    .vis-fullscreen {
        right: 15px;
        bottom: 90px;

        i:before {
            font-family: "Material Design Icons";
            content: mdi('fullscreen');
        }
    }

    .flow-chart-finder-input {
      border-top: none !important;
      //border: 1px solid $gray-200 !important;
      /*text-indent: 8px;*/
    }
    .flow-chart-finder-input::placeholder {
      /*color: #00A655;*/
      /*text-indent: 8px;*/
    }

    .custom-menu {
        display: none;
        z-index: 1000;
        position: absolute;
        overflow: hidden;
        border: 1px solid #cccccc;
        white-space: nowrap;
        font-family: sans-serif;
        background: #ffffff;
        color: #333333;
        border-radius: 5px;
        padding: 0;
    }

    /* Each of the items in the list */
    .custom-menu li {
        padding: 8px 12px;
        cursor: pointer;
        list-style-type: none;
        transition: all .3s ease;
        user-select: none;
    }

    .custom-menu li:hover {
        background-color: #ddeeff;
    }
    .modal.priceErrorModal .modal-content{
      border: 3px solid;
    }
    .flow-chart-finder .input-group-text {
      width: 64px;
      background-color: #00d499;
      border: 1px solid $gray-300;
      border-left: none;
      cursor: pointer;
      .material-icons {
        width: 73px;
        color: white;
      }
    }
</style>
