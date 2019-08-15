<template>
    <div class="file-previews bg-light d-flex flex-row flex-wrap py-4">
      <div class="alert alert-light" role="alert" v-if="requests.length === 0">{{$t('There are no Signature Requests for this File.', 'There are no Signature Requests for this File.')}}</div>
      <table>
            <template v-for="request in requests" class="mb-1">
              <tr><td>
                <table class="signatures">
              <tr><td rowspan="37"><span v-if="signatureStatus(request)=='Pending'"
                                        class="material-icons mdi mdi-clock-outline"></span>
                <span v-else-if="signatureStatus(request)=='Signed'"
                      class="material-icons mdi mdi-check"></span>
                <span v-else-if="signatureStatus(request)=='Rejected'"
                      class="material-icons mdi mdi-close"></span>
                <span v-if="signatureStatus(request)=='Revoked'"
                      class="material-icons mdi mdi-minus"></span></td></tr>
            <tr>
              <td>
                <h5 class="pb-0">{{ $t('Signatory', 'Signatory') }}:</h5>
              </td>
              <td>{{request.signatoryName}} ({{request.signatoryAddr}})</td>
              <td rowspan="2">
                <div v-if="isSignatureStatusPending(request)">
                <button type="button" class="btn btn-primary" @click="revokeFile(request.signatoryAddr)"
                        :title="$t('Revoke', 'Revoke')">
                  <i class="material-icons mr-1">remove_circle_outline</i>
                  <span>{{$t('Revoke Request', 'Revoke Request')}}</span>
                </button>
              </div></td>
            </tr>
                  <tr v-if="request.requestedAt">
                <td><h5 class="pb-0">{{ $t('Requested at', 'Requested at') }}:</h5></td>
                <td>{{request.requestedAt}}</td>
              </tr>
            <tr>
              <td><h5 class="pb-0">{{ $t('Status', 'Status') }}:</h5></td>
              <td>{{signatureStatus(request)}}</td>
            </tr>
              <tr v-if="request.rejected">
                <td><h5 class="pb-0">{{ $t('Rejected at', 'Rejected at') }}:</h5></td>
                <td>{{request.rejectedAt}}</td>
              </tr>
              <tr v-if="request.revoked">
                <td><h5 class="pb-0">{{ $t('Revoked at', 'Revoked at') }}:</h5></td>
                <td>{{request.revokedAt}}</td>
              </tr>
              <tr>&nbsp;</tr>
              </table></td></tr>
            </template>
      </table>
      </div>
</template>

<script>
import mafdc from '@/mixinApp'
import PdfPreview from '@/components/document/PdfPreview'
import moment from 'moment'

export default {
  mixins: [mafdc],
  name: 'signature-request-list',
  components: {
  },
  props: {
    id: null,
    index: null,
    doc: null,
    hash: null
  },
  data () {
    return {
      signers: [],
      me: '',
      requests: []
    }
  },
  async created () {
    await this.app.wallet.getClientProvidedNetwork() // wait until wallet is ready
    this.getRequests()
  },
  async mounted () {
    this.getSigners()
  },
  methods: {
    formatDate (date) {
      return moment(String(date)).format('MM.DD.YYYY hh:mm:ss')
    },
    isSignatureStatusPending (item) {
      return this.signatureStatus(item) === 'Pending'
    },
    isSigned (item) {
      return this.signers.includes(item.signatoryAddr)
    },
    signatureStatus (item) {
      if (item.rejected === true) {
        return 'Rejected'
      }
      if (item.revoked === true) {
        return 'Revoked'
      }
      if (this.isSigned(item) === true) {
        return 'Signed'
      }
      return 'Pending'
    },
    getRequests () {
      axios.get('/api/user/document/signingRequests/' + this.id + '/docs[' + this.index + ']').then(async response => {
        this.requests = response.data
        console.log(response)
      }, (err) => {
        this.app.handleError(err)
        if (err.response && err.response.status === 404) {
        } else if (err.response.status !== 404) {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not load documents'),
            type: 'error'
          })
        }
      })
    },
    async getSigners () {
      if (this.doc.Hash) {
        this.signers = await this.app.wallet.proxeusFS.getFileSigners(this.doc.Hash)
      }
      console.log(this.signers)
    },
    revokeFile (signatory) {
      var bodyFormData = new FormData()
      bodyFormData.set('signatory', signatory)
      axios({
        method: 'post',
        url: '/api/user/document/signingRequests/' + this.id + '/docs[' + this.index + ']/revoke',
        data: bodyFormData,
        config: { headers: { 'Content-Type': 'multipart/form-data' } }
      }).then(response => {
        this.$notify({
          group: 'app',
          title: 'Signature request',
          text: 'Revoked!',
          type: 'success'
        })
        this.$router.go()
      }, (err) => {
        this.$notify({
          group: 'app',
          title: 'Error',
          text: 'Failed to revoke signature',
          type: 'error'
        })
        return false
      }).catch(e => {
        console.log(e)
        return false
      })
    },
    docsPath (doc, index) {
      return '/api/user/document/file/' + doc.id + '/' + this.doc.docID
    },
    documentPreviews () {
      console
      return true
    },
    async hashFile (file) {
      return new Promise((resolve, reject) => {
        let reader = new FileReader()

        reader.onload = (e) => {
          let hash = this.app.wallet.hashFile(reader.result)
          resolve(hash)
        }
        reader.onerror = (e) => {
          reject(e)
        }
        reader.readAsArrayBuffer(file)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  .file-previews button {
    height: 45px;
  }
  ul, ul li {
    list-style: none;
    margin:0;
    padding:0;
  }
  h5 {
    font-weight: 900;
    margin-bottom: 0.1em;
  }
  h4 {
    font-size: 1.2em;
    font-weight: 900;
    margin-bottom: 1.1em;
  }
  table.signatures {
    margin-bottom: 2em;
    pading-bottom: 2em;

    td {
      vertical-align:top;
      padding-bottom: 0.0em;
      padding-right: 3em;
    }
  }
</style>
