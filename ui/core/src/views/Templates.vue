<template>
<div>
  <vue-headful :title="$t('Templates title', 'Proxeus - Templates')"/>
  <top-nav :title="$t('Templates')"></top-nav>
  <div class="main-container">
    <list-group :deleteElementFunc="provideDeleteFunc" :iconFa="iconFa" nodeType="template" path="template">
      <template scope="element">
        <td class="tdmin">
          <button v-if="element && app.amIWriteGrantedFor(element)" :title="$t('delete this template')" @click="areYouSureDeleteDialog($event, element)" type="button"
                  class="btn btn-primary btn-round" style="z-index: 1;padding: 6px;display: inline-block;">
            <i class="material-icons">
              delete_forever
            </i>
          </button>
        </td>
      </template>
    </list-group>
    <list-item-dialog :setup="setupDeleteDialog" :sureFunc="deleteDialogSureAction" :iconFa="iconFa">
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
  name: 'Templates',
  components: {
    ListItemDialog,
    ListGroup,
    TopNav
  },
  data () {
    return {
      iconFa: 'mdi mdi-file-xml',
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
      axios.get('/api/admin/template/' + this.deleteTarget.id + '/delete').then(() => {
        this.deleteFunc(this.deleteTarget.id)
        this.deleteTarget = null
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Template deleted'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not delete Template. Please try again or if the error persists contact the platform operator.\n'),
          type: 'error'
        })
      })
    },
    areYouSureDeleteDialog (e, item) {
      this.deleteTarget = item
      this.deleteDialogShowFunc(e, item)
    }
  }
}
</script>
