<template>
  <div class="documents" style="height:100%;">
    <vue-headful :title="$t('Documents title', 'Proxeus - Document Payment')"/>
    <top-nav :title="$t('Document Payment', 'Document Payment')" bg="#ffffff" class="border-bottom-0" :returnToRoute="{name:'Documents'}"/>
    <div class="main-container">
      <div class="p-4 bg-light" style="position: relative">
        <h1 class="mb-3">{{$t('Document', 'Document')}} </h1>
        <p class="text-danger" v-if="walletErrorMessage">{{ walletErrorMessage }}</p>
        <h2 class="mb-1">{{$t('Name', 'Name')}} </h2>
        <p>{{workflow.name}}</p>
        <h2 class="mb-1">{{$t('Description', 'Description')}}</h2>
        <p>{{workflow.detail}}</p>
        <h2 class="mb-1">{{$t('Price', 'Price')}}</h2>
        <p>{{workflow.price}} XES</p>
        <h2 class="mb-1">{{$t('Owner', 'Owner')}}</h2>
        <p>{{workflow.ownerEthAddress}}</p>
        <br/>
        <div v-show="submitting" class="v-fade" style="position:absolute;top:0;left:0;bottom:0;right:0;margin-top: 50px;">
          <div>
            <spinner background="transparent"/>
            <div style="margin: 155px auto 100px;position: relative; text-align: center;">Please wait for your transaction to be confirmed.</div>
          </div>
        </div>
        <button type="button" class="btn btn-primary mr-4 mb-1" @click="payDocument()" :disabled="submitting || !isConnectedAccount"
                :title="$t('Buy Document', 'Buy Document')">
          <span>{{$t('Buy Document', 'Buy Document')}}</span>
        </button>
        <div v-if="!isConnectedAccount" class="alert alert-warning mb-3" role="alert">
          {{$t('Document buy help 2', 'Make sure that you are logged into the correct account in MetaMask. The payment must be made from the ethereum-account that you have connected in your account settings. Switch Metamask account and make sure that you have set your ethereum address in your account account settings on the top right.')}}
        </div>
        <div class="alert alert-light mb-2" role="alert">
          {{$t('Document buy help 1', 'By clicking "Buy Document" you confirm that you wish to purchase this document creation workflow for one use. You can only purchase one instance of this workflow at a time.')}}
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import TopNav from '@/components/layout/TopNav'
import Spinner from '@/components/Spinner'
import mafdc from '@/mixinApp'
import web3 from 'web3'

export default {
  mixins: [mafdc],
  name: 'DocumentPayment',
  props: ['documentId'],
  components: {
    TopNav,
    Spinner
  },
  data () {
    return {
      accountEthAddress: '',
      walletErrorMessage: '',
      workflow: {},
      submitting: false,
      isConnectedAccount: false,
      me: ''
    }
  },
  created () {
  },
  async mounted () {
    this.getDocument()
    this.setAccountEthAddress()
    this.checkPaymentAndRedirectIfExists()
    this.checkConnectedAccount()
  },
  methods: {
    checkPaymentAndRedirectIfExists () {
      console.log('Check payment ' + this.documentId)
      axios.get('/api/admin/workflow/' + this.documentId + '/payment').then(response => {
        if (response.data) {
          this.$router.push({ name: 'DocumentFlow', params: { id: response.data.workflowID } })
        }
      }, (error) => {
        if (error.response.status !== 404) {
          console.log(error)
        } else {
          console.log('no payment found')
        }
      })
    },
    async checkConnectedAccount () {
      console.log('Check connected account ' + this.documentId)
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
          this.isConnectedAccount = this.me.toLowerCase() === this.accountEthAddress.toLowerCase()
        }
      } catch (e) {
        console.log(e)
      }
    },
    async setAccountEthAddress () {
      let response = await axios.get('/api/me')
      if (response.data.etherPK) {
        this.accountEthAddress = response.data.etherPK
      }
    },
    getDocument () {
      console.log('Getting document ' + this.documentId)
      axios.get('/api/admin/workflow/' + this.documentId).then(response => {
        if (!response.data) {
          this.workflow = {}
        } else {
          this.workflow = response.data
        }
      }, (error) => {
        this.app.handleError(error)
        if (error.response && error.response.status === 404) {
          this.$_error('NotFound', { what: 'Workflow' })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not load Workflow'),
            type: 'error'
          })
          this.$router.push({ to: 'Workflows' })
        }
      })
    },
    wallet () {
      return this.app.wallet
    },
    async checkPaymentReceived (txHash) {
      let response = await axios.get('/api/admin/workflow/' + this.documentId + '/payment', {
        params: {
          txHash: txHash
        }
      })
      if (response.data) {
        return true
      }
      return false
    },
    sleep (seconds) {
      return new Promise(resolve => setTimeout(resolve, seconds * 1000))
    },
    redirectWorkflow () {
      this.$router.push({ name: 'DocumentFlow', params: { id: this.documentId } })
    },
    async payDocument () {
      this.submitting = true

      if (this.me === null) {
        this.submitting = false
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Please login to your wallet and refresh'),
          type: 'error'
        })
        this.walletErrorMessage = this.$t('Please login to your wallet and refresh')
      }

      this.nonce = await this.app.wallet.proxeusFS.web3.eth.getTransactionCount(this.me)
      this.nonce++
      let self = this

      const xesAmountWei = web3.utils.toWei(this.workflow.price.toString(), 'ether')
      this.app.wallet.transferXES(this.workflow.ownerEthAddress, xesAmountWei)
        .then(async (result) => {
          this.submitting = false
          let txHash = result.transactionHash
          let paymentReceived = false
          let tryCount = 0
          do {
            paymentReceived = await self.checkPaymentReceived(txHash)
            tryCount++
            if (tryCount > 10) { // assume in ~30 seconds blockchain nodes are in sync
              this.$notify({
                group: 'app',
                title: this.$t('Error'),
                text: this.$t('Payment failed. Please try again.'),
                type: 'error'
              })
              return
            }
            if (paymentReceived !== true) {
              await self.sleep(3000)
            }
          } while (paymentReceived !== true)
          axios.post('/api/admin/workflow/' + this.documentId + '/payment/' + txHash).then(response => {
          }).catch((error) => {
            console.log(error)
            this.$notify({
              group: 'app',
              title: this.$t('Error'),
              text: this.$t('Payment failed. Please try again.'),
              type: 'error'
            })
          })
          self.redirectWorkflow()
        }).catch((error) => {
          this.submitting = false
          console.log(error)
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Payment failed. Please try again.'),
            type: 'error'
          })
        })
    }
  }
}
</script>

<style scoped>
  .alert {
    font-size: 0.85em;
  }
</style>
