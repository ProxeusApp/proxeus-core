<template>
<div v-if="app.wallet" class="modal fade" v-bind="$attrs" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
     aria-hidden="true">
  <div class="modal-dialog" :class="{'modal-lg': authenticated, 'modal-sm': authenticated === false}">
    <div class="modal-content">
      <div class="modal-header align-items-center bg-light px-3">
        <h5 class="modal-title font-weight-bold">
          <img src="../../assets/images/logo-blue.png" width="15" alt="Proxeus Logo" class="mr-1">
          Proxeus Wallet
        </h5>
        <div class="balance ml-auto" v-if="balance && authenticated">{{ balance }}</div>
        <button type="button" class="btn btn-link text-primary ml-1 pr-0">
          <i class="material-icons mdi mdi-settings"></i>
        </button>
      </div>
      <div class="modal-body px-0 pt-0">
        <keep-alive>
          <component :is="currentViewComponent"
                     v-if="currentViewComponent"
                     :wallet="wallet"
                     :balance="balance"
                     @authenticated="handleAuthenticated"></component>
        </keep-alive>
      </div>
      <!--<div class="modal-footer">-->
      <!--<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>-->
      <!--</div>-->
    </div>
  </div>
</div>
</template>

<script>
/*
 * Dynamic View Components
 */
import Account from './Account'
import Login from './Login'
import SignMessage from './SignMessage'
import TransactionHistory from './TransactionHistory'

/*
 * Wallet Library
 */
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'wallet',
  components: {
    Account,
    Login,
    SignMessage,
    TransactionHistory
  },
  created () {
  },
  mounted () {
  },
  data () {
    return {
      authenticated: false,
      balance: undefined
    }
  },
  computed: {
    currentViewComponent () {
      return this.authenticated === true ? Account : Login
    },
    wallet () {
      return this.app.wallet
    }
  },
  methods: {
    handleAuthenticated () {
      this.refreshBalance()
      this.authenticated = true
    },
    refreshBalance () {
      return new Promise((resolve, reject) => {
        wallet.getXESBalance('0x74Efe378f700436B7987435A8466Ab638e7B7C92', 5).then(res => {
          if (!res) {
            reject(new Error('Could not get balance'))
          }
          this.balance = res
          resolve(res)
        })
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  .modal {
    .modal-dialog {
      position: absolute;
      right: 50px;
      top: 45px;
      min-width: 300px;
      transition: all 350ms !important;

      &.modal-sm {
      }

      &.modal-lg {
        min-width: 550px;
      }
      .modal-content {
        box-shadow: 0 40px 100px rgba(0, 0, 0, .35);
      }
    }
  }

</style>
