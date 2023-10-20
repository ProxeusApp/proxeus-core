<template>
<div class="doc-viewer" v-if="document">
  <vue-headful
    :title="$t('Document viewer title prefix','Proxeus - ')+(document && document.name || $t('Document title','Document'))"/>
  <top-nav :title="document.name || $t('Unnamed document')" bg="#f7f8f9" class="border-bottom-0"
           :returnToRoute="{name:'Documents'}">
    <span slot="buttons" class="btn btn-secondary" v-show="hasChanges">Saving ...</span>
    <!--<button slot="buttons" class="btn btn-primary ml-3" @click="infoToggled = !infoToggled">Edit infos</button>-->
    <button slot="buttons" v-if="document && app.userIsUserOrHigher()" style="height: 40px;"
            @click="app.exportData('&id='+document.id, null, '/api/userdata/export','UserData_'+document.id)"
            type="button" class="btn btn-primary ml-2">
      <span>{{$t('Export')}}</span></button>
  </top-nav>
  <div class="col-sm-12 mt-3">
    <div class="row">
      <name-and-detail-input :input="changed" v-model="document"/>
    </div>
    <hr/>
    <div class="row" v-if="documentPreviews">
      <!--<div class="col-sm-12 col-md-6">-->
      <!--<div class="card">-->
      <!--<div class="card-header bg-white"><h2 class="mb-0">Uploaded files</h2></div>-->
      <!--<div class="card-body p-0"><div class="file-previews bg-light">-->
      <!--<pdf-preview v-for="doc in documentPreviews" :key="doc.id" :src="doc.src"-->
      <!--v-if="documentPreviews" :filename="doc.filename"></pdf-preview>-->
      <!--</div></div>-->
      <!--</div>-->
      <!--</div>-->
      <div class="col-sm-12 col-md-12">
        <h2>Documents</h2>
        <div class="card">
          <!--<div class="card-header bg-white"><h2 class="mb-0">Documents</h2></div>-->
          <div class="card-body p-0 docsandsign">
            <div class="file-previews bg-light d-flex flex-row flex-wrap">
              <table>
              <tr v-for="(doc, index) in documentPreviews">
                <td>
                <pdf-preview :key="index"
                             :src="docsPath(index)"
                             v-if="documentPreviews" :item="doc" :filename="doc.name"/>
                </td>
                <td>
                  <signature-request-list :id="id" :doc="doc" :index="index"></signature-request-list>
                </td>
              </tr>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import PdfPreview from '@/components/document/PdfPreview.vue'
// import SearchBox from '@/components/SearchBox.vue'
// import ListGroup from '@/components/ListGroup.vue'
import TopNav from '@/components/layout/TopNav.vue'
import NameAndDetailInput from '@/components/NameAndDetailInput.vue'
import mafdc from '@/mixinApp'
import SignatureRequestList from '@/components/SignatureRequestList.vue'

export default {
  mixins: [mafdc],
  name: 'document-viewer',
  props: {
    id: {
      required: true
    }
  },
  components: {
    NameAndDetailInput,
    PdfPreview,
    // SearchBox,
    // ListGroup,
    TopNav,
    SignatureRequestList
  },
  data () {
    return {
      document: undefined,
      hasChanges: false,
      infoToggled: false,
      isLastRequest: true
    }
  },
  created () {
    axios.get('/api/user/document/' + this.id).then(response => {
      // JSON responses are automatically parsed.
      this.document = response.data
    }, (err) => {
      this.app.handleError(err)
      if (err.response && err.response.status === 404) {
        this.$_error('NotFound', { what: 'Document' })
      } else {
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not load document. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
        this.$router.push({ name: 'Documents' })
      }
    })
  },
  computed: {
    documentPreviews () {
      if (this.document && this.document.data && this.document.data.docs) {
        return this.document.data.docs
      }
      return false
    }
  },
  methods: {
    docsPath (index) {
      return '/api/user/document/file/' + this.document.id + '/docs[' + index + ']'
    },
    search (term) {
      axios.get('/api/user/document?c=' + term).then(response => {
        // JSON responses are automatically parsed.
        this.documents = response.data
      }, (err) => {
        this.app.handleError(err)
      })
    },
    changed () {
      this.hasChanges = true
      if (this.isLastRequest) {
        this.isLastRequest = false
        setTimeout(() => {
          axios.post('/api/document/' + this.id + '/name', { name: this.document.name, detail: this.document.detail })
            .then(response => {
              // Request succeeded
            }, (err) => {
              this.app.handleError(err)
              this.$notify({
                group: 'app',
                title: this.$t('Error'),
                text: this.$t('Could not change name. Please try again or if the error persists contact the platform operator.'),
                type: 'error'
              })
            })
          this.isLastRequest = true
          this.hasChanges = false
        }, 1500)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .list-group-item {
    border-left: 0;
    border-right: 0;

    .lg-small h5 {
      font-weight: 400;
      font-size: .9rem;
    }
  }

  .list-group-icon {
    color: #747782;
    font-size: 2rem;
  }

  ::v-deep .pdf-preview {
    max-width: 150px !important;
    margin-left: 50px !important;
    margin-right: 20px !important;
    display: inline-block !important;
  }

  .file-previews {
    padding-top: 1rem;
    padding-bottom: 1rem;
  .preview-wrapper {
    width: 300px;
  }
  td {
    vertical-align:top;
  }
  }

  .docsandsign {
  td div {
    padding-top: 0px;  }
  }

</style>
