<template>
  <div class="documents" style="height:100%;">
    <vue-headful :title="$t('Documents title', 'Proxeus - Signature Requests')"/>
    <top-nav :title="$t('Signature Requests', 'Signature Requests')" bg="#ffffff" class="border-bottom-0"/>
    <div class="main-container">
      <div class="col-sm-12 mt-3">
        <div class="row">
          <div class="col-sm-12 col-md-12">
            <div v-if="appBlockchainNet !== clientProvidedNet" class="alert alert-warning mb-3" role="alert">
              {{$t('Document buy help 3', 'Make sure you are logged in to the "{network}" network in metamask.', { network: appBlockchainNet })}}
            </div>
            <div v-if="!isConnectedAccount" class="alert alert-warning mb-3" role="alert">
              {{$t('Document sign help', 'Make sure that you are logged into the correct account in MetaMask. The signature must be made from the ethereum-account that you have connected in your account settings. Switch Metamask account and make sure that you have set your ethereum address in your account account settings on the top right.')}}
            </div>
            <p class="text-danger" v-if="walletErrorMessage">{{ walletErrorMessage }}</p>
            <h3 class="mb-4">{{$t('This page lists the documents for which your signature has been requested. You can either sign or reject to sign the documents.',
              'This page lists the documents for which your signature has been requested. You can either sign or reject to sign the documents.')}}</h3>
            <div class="alert alert-light" role="alert" v-if="documents.length === 0">{{$t('You currently have no Signature Requests.', 'You currently have no Signature Requests.')}}</div>
            <div v-else v-for="(doc, index) in documents" class="card mb-1">
              <div class="card-body p-0">
                <signature-request :doc="doc" :index="index" :isConnectedAccount="isConnectedAccount" :appBlockchainNet="appBlockchainNet" :clientProvidedNet="clientProvidedNet">
                </signature-request>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopNav from '@/components/layout/TopNav.vue'
import mafdc from '@/mixinApp'
import SignatureRequest from '@/components/SignatureRequest.vue'

export default {
  mixins: [mafdc],
  name: 'signature-requests',
  components: {
    TopNav,
    SignatureRequest
  },
  data () {
    return {
      documents: [],
      nonce: 0,
      walletErrorMessage: '',
      accountEthAddress: '',
      isConnectedAccount: false,
      appBlockchainNet: '',
      clientProvidedNet: '',
      me: ''
    }
  },
  created () {
    this.getSigningRequests()
  },
  async mounted () {
    this.setAccountEthAddress()
    this.checkConnectedAccount()
    this.checkBlockchainNetwork()
  },
  methods: {
    async checkBlockchainNetwork () {
      const r = await axios.get('/api/config')
      if (!r.data) {
        return
      }
      if (!r.data.blockchainNet) {
        return
      }
      this.appBlockchainNet = r.data.blockchainNet

      try {
        this.clientProvidedNet = await this.app.wallet.getClientProvidedNetwork()
      } catch (e) {
        console.log(e)
      }
    },
    async setAccountEthAddress () {
      const response = await axios.get('/api/me')
      if (response.data.etherPK) {
        this.accountEthAddress = response.data.etherPK
      }
    },
    async checkConnectedAccount () {
      let response2 = null
      try {
        response2 = await axios.get('/api/me')
        if (response2.data.etherPK) {
          this.accountEthAddress = response2.data.etherPK

          await this.app.wallet.getClientProvidedNetwork()
          this.me = window.ethereum.selectedAddress || this.app.wallet.getCurrentAddress()

          if (this.me === '' || this.me === null ||
            this.accountEthAddress === '' || this.accountEthAddress === null) {
            this.isConnectedAccount = false
            return
          }
          // check lowercase because window.ethereum.selectedAddress returns lowercase hash
          this.isConnectedAccount = this.me.toLowerCase() === this.accountEthAddress.toLowerCase()
        }
      } catch (e) {
        console.log(e)
      }
    },
    getSigningRequests () {
      axios.get('/api/user/document/signingRequests').then(async response => {
        this.documents = response.data
      }, (err) => {
        this.app.handleError(err)
        if (err.response && err.response.status === 404) {
        } else if (err.response.status !== 404) {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not load documents. Please try again or if the error persists contact the platform operator.\n'),
            type: 'error'
          })
          this.$router.push({ name: 'Documents' })
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  .file-previews button {
    height: 45px;
  }
</style>
