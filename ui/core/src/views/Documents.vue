<template>
<div class="documents" style="height:100%;">
  <vue-headful :title="$t('Documents title', 'Proxeus - Documents')"/>
  <top-nav :title="$t('Documents')" bg="#ffffff" class="border-bottom-0"/>
  <div class="main-container">
    <list-group :deleteElementFunc="provideDeleteFunc" :defaultName="$t('Unnamed')" :linkResolver="getDocRoute"
                :iconFa="'mdi mdi-view-carousel'" :fixedPath="'/api/user/document'" nodeType="UserData">
      <router-link slot="addBtn" :to="{name:'document-create'}" class="btn btn-primary">
        Create document
      </router-link>
      <template scope="element">
          <div v-if="element.finished === false" class="easy-read"><span
            class="badge badge-info">{{$t('inline draft badge','Draft')}}</span></div>
          <div v-else></div>
          <small class="light-text"></small>
          <button v-if="element && app.amIWriteGrantedFor(element)" :title="$t('delete this item')" @click="areYouSureDialog($event, element)" type="button"
                  class="btn btn-primary btn-round ml-3" style="z-index: 1;padding: 6px;display: inline-block;">
            <i class="material-icons">
              delete_forever
            </i>
          </button>
      </template>
    </list-group>
    <list-item-dialog :setup="setupDeleteDialog" :sureFunc="deleteDialogAction" :iconFa="'mdi mdi-view-carousel'">
      <div class="d-block mb-2">{{$t('Are you sure, you want to delete?')}}</div>
      <div class="d-block mb-2">{{$t('workflow delete xes alert', 'If you have paid XES for this workflow, you will not be able to start the workflow again, unless you make another payment.')}}</div>
      <div class="d-block mb-2">{{$t('This action can\'t be undone.')}}</div>
    </list-item-dialog>
  </div>
  <first-login-overlay keyz="documents"
                       preview-url="https://docs.google.com/document/d/1PoJMmdBt8bu1tfqbOBJ3g87V-z1yjBrhlNZpOW2wQWk/preview">
  </first-login-overlay>
</div>
</template>

<script>
import ListGroup from '@/components/ListGroup'
import TopNav from '@/components/layout/TopNav'
import ListItem from '../components/ListItem'
import mafdc from '@/mixinApp'
import ListItemDialog from './appDependentComponents/ListItemDialog'
import FirstLoginOverlay from '@/views/FirstLoginOverlay'

export default {
  mixins: [mafdc],
  name: 'documents',
  components: {
    ListItemDialog,
    ListItem,
    ListGroup,
    TopNav,
    FirstLoginOverlay
  },
  data () {
    return {
      deleteTarget: null,
      deleteFunc: null,
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
    deleteDialogAction () {
      axios.get('/api/document/' + this.deleteTarget.id + '/delete').then(() => {
        this.deleteFunc(this.deleteTarget.id)
        this.deleteTarget = null
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Item deleted'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not delete item. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    },
    areYouSureDialog (e, item) {
      this.deleteTarget = item
      this.deleteDialogShowFunc(e, item)
    },
    getDocRoute (doc) {
      return {
        name: doc.finished ? 'DocumentViewer' : 'DocumentFlow',
        params: {
          id: doc.finished ? doc.id : doc.workflowID
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .list-group-item {
    border-left: 0;
    border-right: 0;
    line-height: 1;

    &[data-index="0"] {
      border-top: 0;
    }

    .lg-small h5 {
      font-weight: 400;
      font-size: .9rem;
    }
  }

  .flex-col-truncate {
    overflow: hidden;
    min-width: 0;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
