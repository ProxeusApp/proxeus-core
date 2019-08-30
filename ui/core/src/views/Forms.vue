<template>
<div>
  <vue-headful :title="$t('Forms title', 'Proxeus - Forms')"/>
  <top-nav :title="$t('Forms')"></top-nav>
  <div class="main-container">
    <list-group :deleteElementFunc="provideDeleteFunc" :prependFunc="prependFunc" :icon="icon" nodeType="form" path="form">
      <template scope="element">
      <td class="tdmin">
        <button :title="$t('copy this form')" @click="areYouSureDialog($event, element)" type="button"
                class="btn btn-primary btn-round mshadow-light" style="z-index: 1;padding: 6px;display: inline-block;">
          <i class="material-icons">file_copy</i>
        </button>
      </td>
      <td class="tdmin">
        <button v-if="element && app.amIWriteGrantedFor(element)" :title="$t('delete this form')" @click="areYouSureDeleteDialog($event, element)" type="button"
                class="btn btn-primary btn-round mshadow-light" style="z-index: 1;padding: 6px;display: inline-block;">
          <i class="material-icons">
            delete_forever
          </i>
        </button>
      </td>
      </template>
    </list-group>
    <list-item-dialog :setup="setupCopyDialog" :sureFunc="copyDialogSureAction" :icon="icon">
      <div class="d-block fregular">{{$t('Are you sure, you want to copy this item?')}}</div>
    </list-item-dialog>
    <list-item-dialog :setup="setupDeleteDialog" :sureFunc="deleteDialogSureAction" :icon="icon">
      <div class="d-block fregular">{{$t('This action can\'t be undone.')}}</div>
      <div class="d-block">{{$t('Are you sure, you want to delete?')}}</div>
    </list-item-dialog>
  </div>
</div>
</template>

<script>
import ListGroup from '@/components/ListGroup'
import TopNav from '@/components/layout/TopNav'
import ListItemDialog from './appDependentComponents/ListItemDialog'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'Forms',
  components: {
    ListItemDialog,
    ListGroup,
    TopNav
  },
  data () {
    return {
      icon: 'view_quilt',
      prependItem: null,
      copyItem: null,
      showCopyDialog: () => {},
      deleteTarget: null,
      deleteFunc: () => {},
      deleteDialogShowFunc: () => {}
    }
  },
  methods: {
    setupDeleteDialog (showFunc) {
      this.deleteDialogShowFunc = showFunc
    },
    provideDeleteFunc (fun) {
      this.deleteFunc = fun
    },
    deleteDialogSureAction () {
      axios.get('/api/admin/form/' + this.deleteTarget.id + '/delete').then(() => {
        this.deleteFunc(this.deleteTarget.id)
        this.deleteTarget = null
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Form successfully deleted'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not delete Form. Please try again or if the error persists contact the platform operator.'),
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
      axios.get('/api/admin/form/' + item.id).then(response => {
        if (response.data) {
          let item = response.data
          // empty id to ensure a new entry is created
          item.id = ''
          item.name = item.name + this.$t(' Copy')
          axios.post('/api/admin/form/update', item).then(response => {
            this.$notify({
              group: 'app',
              title: this.$t('Success'),
              text: this.$t('Form successfully copied.'),
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
              text: this.$t('Could not copy Form. Please try again or if the error persists contact the platform operator.'),
              type: 'error'
            })
          })
        }
      }, (error) => {
        this.app.handleError(error)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not copy Form. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    }
  }
}
</script>
