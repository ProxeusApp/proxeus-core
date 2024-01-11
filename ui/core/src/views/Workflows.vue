<template>
<div style="height:100%;">
  <vue-headful :title="$t('Workflows title','Proxeus - Workflows')"/>
  <top-nav :title="$t('Workflows')"></top-nav>
  <div class="main-container">
    <list-group :deleteElementFunc="provideDeleteFunc" :prependFunc="prependFunc" :iconFa="iconFa" nodeType="workflow"
                path="workflow" :display-price="true" @error="handleListError">
      <template scope="element">
        <button :title="$t('copy this workflow')" @click="areYouSureDialog($event, element)" type="button"
                class="btn btn-primary btn-sm">
          Duplicate
        </button>
        <a target="_blank" :title="$t('run this workflow')" @click="$event.stopPropagation();" :href="'/document/'+element.id"
           class="btn btn-primary btn-sm ml-2">
          Run
        </a>
        <button v-if="element && app.amIWriteGrantedFor(element)" :title="$t('delete this workflow')" @click="areYouSureDeleteDialog($event, element)" type="button"
                class="btn btn-primary btn-sm ml-2">
            Delete
        </button>
      </template>
    </list-group>
    <list-item-dialog :setup="setupCopyDialog" :sureFunc="copyDialogSureAction" :iconFa="iconFa">
      <div class="d-block">{{$t('Are you sure, you want to copy this item?')}}</div>
    </list-item-dialog>
    <list-item-dialog :setup="setupDeleteDialog" :sureFunc="deleteDialogSureAction" :iconFa="iconFa">
      <div class="d-block">{{$t('This action can\'t be undone.')}}</div>
      <div class="d-block">{{$t('Are you sure, you want to delete?')}}</div>
    </list-item-dialog>
  </div>
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
</div>
</template>

<script>
import ListGroup from '@/components/ListGroup'
import TopNav from '@/components/layout/TopNav'
import axios from 'axios'
import mafdc from '@/mixinApp'
import ListItemDialog from './appDependentComponents/ListItemDialog'

export default {
  mixins: [mafdc],
  name: 'Workflows',
  components: {
    ListItemDialog,
    ListGroup,
    TopNav
  },
  data () {
    return {
      iconFa: 'mdi mdi-source-branch',
      deleteTarget: null,
      prependItem: null,
      copyItem: null,
      deleteFunc: () => {},
      showCopyDialog: () => {},
      deleteDialogShowFunc: () => {}
    }
  },
  methods: {
    handleListError (priceError) {
      if (priceError.response.data === 'can not set price without eth addr') {
        this.togglePriceErrorModal(true)
      }
    },
    togglePriceErrorModal (show) {
      if (show === true) {
        $(this.$refs.priceErrorModal).show()
        return
      }
      $(this.$refs.priceErrorModal).hide()
    },
    setupDeleteDialog (showFunc) {
      this.deleteDialogShowFunc = showFunc
    },
    provideDeleteFunc (fun) {
      this.deleteFunc = fun
    },
    deleteDialogSureAction () {
      axios.get('/api/admin/workflow/' + this.deleteTarget.id + '/delete').then(() => {
        this.deleteFunc(this.deleteTarget.id)
        this.deleteTarget = null
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Workflow deleted'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not delete Workflow. Please try again or if the error persists contact the platform operator.\n'),
          type: 'error'
        })
      })
    },
    areYouSureDeleteDialog (e, item) {
      this.deleteTarget = item
      this.deleteDialogShowFunc(e, item)
    },

    setupCopyDialog (showFunc) {
      this.showCopyDialog = showFunc
    },
    copyDialogSureAction () {
      this.copy(this.copyItem)
      this.copyItem = null
    },
    areYouSureDialog (e, item) {
      e.stopPropagation()
      this.copyItem = item
      this.showCopyDialog(e, item)
    },
    prependFunc (f) {
      this.prependItem = f
    },
    copy (item) {
      // get the full item first as we are loading only meta data into the list
      axios.get('/api/admin/workflow/' + item.id).then(response => {
        if (response.data) {
          const item = response.data
          // empty id to ensure a new entry is created
          item.id = ''
          item.name = item.name + this.$t(' Copy')
          axios.post('/api/admin/workflow/update', item).then(response => {
            this.$notify({
              group: 'app',
              title: this.$t('Success'),
              text: this.$t('Workflow copied'),
              type: 'success'
            })
            if (response.data && response.data.id) {
              this.prependItem(response.data)
            }
          }, (err) => {
            this.app.handleError(err)
            this.$notify({
              group: 'app',
              title: this.$t('Error'),
              text: this.$t('Could not copy Workflow. Please try again or if the error persists contact the platform operator.\n'),
              type: 'error'
            })
          })
        }
      }, (error) => {
        this.app.handleError(error)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not copy Workflow. Please try again or if the error persists contact the platform operator.\n'),
          type: 'error'
        })
      })
    }
  }
}
</script>

<style>
  .modal.priceErrorModal .modal-content{
    border: 3px solid;
  }
</style>
