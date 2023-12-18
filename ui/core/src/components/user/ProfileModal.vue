<template>
<b-modal class="profile-modal b-modal" v-model="show"
         :title="$t('Account title', 'Account')"
         :header-bg-variant="'light'" @hide="onDialogHide">
  <animated-input :max="80" :label="$t('Name')" v-model="me.name"/>
  <animated-input :max="100" :label="$t('More about me')" v-model="me.detail"/>
  <animated-input :max="100" :disabled="true" :label="$t('Email')" v-model="me.email"
                  :action="{Name:$t('inline change badge','change'), Func:changeEmail}"/>
  <div v-show="showNewEmailField" class="">
    <form v-show="!done" @submit.prevent="changeEmailRequest">
      <animated-input ref="inputEmail" name="email" :max="100" :label="$t('New Email')" v-model="email"/>
      <div>
        <span class="text-muted"
              style="white-space: normal;">{{$t('New email explanation','Provide a valid email and proceed by clicking the button below to receive a link.')}}</span>
      </div>
      <table style="width: 100%;">
        <tr>
          <td style="text-align:right">
            <button class="btn btn-primary px-3 mt-1" type="submit">{{$t('Validate new email', 'Validate')}}</button>
          </td>
        </tr>
      </table>
    </form>
    <div v-show="done">
      <div class="my-3">{{$t('Email sent','Email sent.')}}</div>
    </div>
  </div>
  <animated-input :max="100" :label="$t('Password')" v-model="pass" type="password"/>
  <span class="text-muted"
        style="white-space: normal;">{{$t('Password explanation','Leave it empty, if you don\'t want to change it.')}}</span>
  <animated-input v-if="app.wallet" ref="inputBcAddr" name="etherPK" :max="100" :disabled="true" :label="$t('Ethereum Address')"
                  v-model="me.etherPK" :action="{Name:$t('inline change badge','change'), Func:updateEthereumAddress}"/>
  <span v-show="walletErrorMessage" class="error">{{walletErrorMessage}}</span>
  <br>
  <div class="sub-title">{{$t('Privacy settings')}}</div>
  <checkbox :label="$t('Want to be found')" v-model="me.wantToBeFound"/>
  <span class="text-muted"
        style="white-space: normal;">{{$t('Want to be found explanation','Uncheck this property if you want to be found only by your blockchain address.')}}</span>
    <br><br>
    <div class="sub-title">{{$t('Delete account')}}</div>
    <div v-if="showConfirm">
      <button @click="showConfirm = false" class="btn btn-secondary mt-1 ml-1">{{$t('Cancel')}}</button>
      <button class="btn btn-danger mt-1 float-right mr-1" @click="deleteAccount">{{$t('Confirm account deletion')}}</button>
    </div>
    <button v-if="!showConfirm" class="btn btn-danger mt-1 ml-1" @click="showConfirm = true">{{$t('Delete account')}}</button>
    <div class="text-muted mt-1 ml-1"
          style="white-space: normal;">{{$t('Account deletion can not be undone.')}}</div>
    <span v-show="deleteErrorMessage" class="error">{{deleteErrorMessage}}</span>
    <api-key :user="me" />
    <template slot="modal-footer">
  <button @click="onDialogHide" class="btn btn-secondary">
    <i class="material-icons">cancel</i>
  </button>
  <button @click="onDialogOk" class="btn btn-primary">
    <i class="material-icons">save</i>
  </button>
  </template>
</b-modal>
</template>

<script>
// import bModal from 'bootstrap-vue/es/components/modal/modal'
// import bModalDirective from 'bootstrap-vue/es/directives/modal/modal'
import AnimatedInput from '../AnimatedInput.vue'
import Checkbox from '../Checkbox.vue'
import mafdc from '@/mixinApp'
import ApiKey from './ApiKey.vue'

export default {
  mixins: [mafdc],
  name: 'profile-modal',
  components: {
    ApiKey,
    Checkbox,
    AnimatedInput,
    'b-modal': bModal
  },
  directives: {
    'b-modal': bModalDirective
  },
  props: {
    setup: {
      type: Function,
      default: () => null
    }
  },
  data () {
    return {
      show: false,
      done: false,
      me: {},
      pass: '',
      email: '',
      showNewEmailField: false,
      walletErrorMessage: '',
      deleteErrorMessage: '',
      challenge: '',
      showConfirm: false
    }
  },
  created () {
    if (this.setup) {
      this.setup(this.openDialog)
    }
  },
  mounted () {
  },
  computed: {
    wallet () {
      return this.app.wallet
    }
  },
  methods: {
    async deleteAccount () {
      try {
        const result = await axios.post('/api/user/delete')
        if (result.status === 200) {
          window.location.href = '/'
          return
        }
      } catch (e) {
        this.deleteErrorMessage = this.$t('There was en error deleting your account. Please try later or contact support.')
      }
    },
    updateEthereumAddress () {
      this.walletErrorMessage = ''
      if (!this.challenge) {
        if (typeof window.ethereum !== 'undefined') {
          axios.get('/api/challenge').then((response) => {
            this.challenge = response.data
            this.metamaskLogin()
          }, (err) => {
            this.app.handleError(err)
          })
        }
      } else {
        this.metamaskLogin()
      }
    },
    async metamaskLogin () {
      if (window.ethereum) {
        try {
          await window.ethereum.enable()
          await this.wallet.wallet.setupDefaultAccount()
        } catch (e) {
          this.walletErrorMessage = this.$t('Please grant access to MetaMask.')
          return
        }
      } else {
        this.walletErrorMessage = this.$t('Please grant access to MetaMask.')
        return
      }
      this.account = this.wallet.getCurrentAddress()
      if (this.account === undefined) {
        this.walletErrorMessage = this.$t('Please sign in to MetaMask.', 'Please log in to MetaMask')
        return
      }
      this.wallet.signMessage(this.challenge, this.account).then((signature) => {
        $(this.$refs.inputBcAddr.$refs.field).cleanFieldErrors()
        axios.post('/api/change/bcaddress', { signature }).then((res) => {
          this.challenge = null
          if (res.status >= 200 && res.status <= 299) {
            this.me.etherPK = this.account
          } else {
            this.walletErrorMessage = this.$t('Could not verify signature.')
          }
        }, (err) => {
          this.challenge = null
          this.app.handleError(err)
          if (err.response && err.response.status === 422) {
            $(this.$refs.inputBcAddr.$refs.field).showFieldErrors({ errors: err.response.data })
          } else {
            this.walletErrorMessage = this.$t('Could not verify signature.')
          }
        })
      }).catch(() => {
        this.walletErrorMessage = this.$t('Could not Sign Message.')
      })
    },
    cleanErr () {
      $(this.$refs.inputEmail.$refs.field).cleanFieldErrors()
    },
    changeEmailRequest () {
      axios.post('/api/change/email', { email: this.email }).then(res => {
        this.cleanErr()
        this.done = true
      }, (err) => {
        this.cleanErr()
        this.app.handleError(err)
        if (err.response && err.response.status === 422) {
          $(this.$refs.inputEmail.$refs.field).showFieldErrors({ errors: err.response.data })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Warning'),
            text: this.$t('There was an error changing the email address. Please try again or if the error persists contact the platform operator.'),
            type: 'warning'
          })
        }
      })
    },
    changeEmail () {
      this.showNewEmailField = true
      this.$nextTick(() => {
        $(this.$refs.inputEmail.$refs.field).focus()
      })
    },
    updateMe () {
      this.app.loadMe(this.updatedMe)
    },
    updatedMe (me) {
      this.me = me
    },
    openDialog () {
      this.show = true
      if (this.app.me) {
        this.me = this.app.me
      } else {
        this.updateMe()
      }
    },
    onDialogHide () {
      this.show = false
      this.$emit('onDialogHide')
    },
    onDialogOk () {
      this.save()
    },
    save () {
      if (this.pass) {
        this.me.password = this.pass
      }
      axios.post('/api/me', this.me).then(response => {
        this.pass = ''
        delete this.me.password
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Successfully saved changes.'),
          type: 'success'
        })
        this.show = false
        this.updateMe()
        this.$emit('onDialogOk')
      }, (err) => {
        this.pass = ''
        delete this.me.password
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('There was an error saving the changes. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    }
  }
}
</script>

<style>

</style>
