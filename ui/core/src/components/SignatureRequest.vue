<template>
    <div class="file-previews bg-light d-flex flex-row flex-wrap py-4">
      <pdf-preview :key="index" :src="docsPath(doc, index)" :showSig="false"
               v-if="documentPreviews" :item="doc" :filename="doc.name"/>
    <div class="ml-3">
      <table class="signatures">
        <tr>
          <td>
            <h5 class="pb-0">{{ $t('Document hash', 'Document hash') }}:</h5>
          </td>
          <td>{{doc.hash}}</td>
        </tr>
        <tr>
          <td><h5 class="pb-0">{{ $t('From', 'From') }}:</h5></td>
          <td>{{requestorNameOrAddress}}</td>
        </tr>
        <tr v-if="doc.requestedAt">
          <td><h5 class="pb-0">{{ $t('Requested at', 'Requested at') }}:</h5></td>
          <td>{{doc.requestedAt}}</td>
        </tr>
        <tr>
          <td><h5 class="pb-0">{{ $t('Status', 'Status') }}:</h5></td>
          <td>{{signatureStatus}}</td>
        </tr>
        <tr v-if="doc.rejected">
          <td><h5 class="pb-0">{{ $t('Rejected at', 'Rejected at') }}:</h5></td>
          <td>{{doc.rejectedAt}}</td>
        </tr>
        <tr v-if="doc.revoked">
          <td><h5 class="pb-0">{{ $t('Revoked at', 'Revoked at') }}:</h5></td>
          <td>{{doc.revokedAt}}</td>
        </tr>
        <tr>
          <td><h5>{{$t('Signers','Signers')}}:</h5></td>
          <td>
            <div v-if="signers.length <= 0">{{$t('None','None')}}</div>
            <div v-else>
              <ul>
                <li v-for="signer in signers">
                  {{ signerName(signer) }}
                </li>
              </ul>
            </div>
          </td>
        </tr>
      </table>
        <div v-if="isSignatureStatusPending">
          <hr>
          <button type="button" class="btn btn-primary mr-4" @click="signFile()" :disabled="processing || !isConnectedAccount || appBlockchainNet !== clientProvidedNet"
                  :title="$t('Sign file', 'Sign file')">
            <i class="material-icons mr-1">check_circle_outline</i>
            <span>{{$t('Sign', 'Sign')}}</span>
          </button>
          <button type="button" class="btn btn-primary" @click="rejectFile()" :disabled="processing || !isConnectedAccount || appBlockchainNet !== clientProvidedNet"
                  :title="$t('Reject', 'Reject')">
            <i class="material-icons mr-1">remove_circle_outline</i>
            <span>{{$t('Reject', 'Reject')}}</span>
          </button>
        </div>
      </div>
    </div>
</template>

<script>
import mafdc from '@/mixinApp'
import PdfPreview from '@/components/document/PdfPreview'

export default {
  mixins: [mafdc],
  name: 'signature-request',
  components: {
    PdfPreview
  },
  props: {
    index: null,
    doc: null,
    isConnectedAccount: {
      type: Boolean,
      default: false
    },
    appBlockchainNet: {
      type: String,
      default: ''
    },
    clientProvidedNet: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      signers: [],
      me: '',
      processing: false
    }
  },
  async mounted () {
    await this.app.wallet.getClientProvidedNetwork() // wait until wallet is ready
    this.getSigners()
    this.setAccount()
  },
  computed: {
    requestorNameOrAddress () {
      return this.doc.requestorName ? this.doc.requestorName + ' (' + this.doc.requestorAddr + ')' : this.doc.requestorAddr
    },
    isSignedByMe () {
      return this.signers.includes(this.me)
    },
    signatureStatus () {
      if (this.isSignedByMe === true) {
        return 'Signed'
      }
      if (this.doc.revoked === true) {
        return 'Revoked'
      }
      if (this.doc.rejected === true) {
        return 'Rejected'
      }
      return 'Pending'
    },
    isSignatureStatusPending () {
      return this.signatureStatus === 'Pending'
    }
  },
  methods: {
    setAccount () {
      this.me = this.app.wallet.getCurrentAddress()
    },
    signerName (signerAddress) {
      const myAccountTrans = this.$t('My Account', 'My Account')
      return signerAddress === this.me ? myAccountTrans + ' (' + signerAddress + ')' : signerAddress
    },
    async getSigners () {
      this.signers = await this.app.wallet.proxeusFS.getFileSigners(this.doc.hash)
    },
    rejectFile () {
      axios.post('/api/user/document/signingRequests/' + this.doc.id + '/' + this.doc.docID + '/reject')
        .then(response => {
          this.doc.rejected = true
          this.deductSignerCount()
        }, (err) => {
          if (err.response && err.response.status === 404) {
            this.$notify({
              group: 'app',
              title: this.$t('Error'),
              text: this.$t('Unable to reject file. Make sure you have an eth address set in your account settings.'),
              type: 'error'
            })
            return
          }

          this.doc.rejected = false
          this.app.handleError(err)
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Unable to reject file. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
        })
    },
    deductSignerCount () {
      this.$store.dispatch('UPDATE_SIGNERS_COUNT', { sigCount: this.$store.getters.signatureRequestCount - 1 })
    },
    async signFile () {
      this.processing = true
      this.setAccount()

      if (this.me === null) {
        this.processing = false
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Please login to your wallet.'),
          type: 'error'
        })
      }

      this.nonce = await this.app.wallet.proxeusFS.web3.eth.getTransactionCount(this.me)
      this.nonce++
      const self = this
      this.app.wallet.proxeusFS.signFile({
        from: this.me,
        hash: this.doc.hash,
        nonce: this.nonce
      }).then((result) => {
        this.processing = false
        self.getSigners()
        this.deductSignerCount()
      }).catch((error) => {
        console.log(error)
        this.processing = false
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not sign document. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    },
    docsPath (doc, index) {
      return '/api/user/document/file/' + doc.id + '/' + this.doc.docID
    },
    documentPreviews () {
      return true
    },
    async hashFile (file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()

        reader.onload = (e) => {
          const hash = this.app.wallet.hashFile(reader.result)
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
    margin-bottom: -1em;
    td {
      padding-bottom: 0.8em;
      padding-right: 3em;
    }
  }
</style>
