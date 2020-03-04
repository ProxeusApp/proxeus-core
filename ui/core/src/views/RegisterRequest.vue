<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Sign up title', 'Proxeus - Sign up')"/>
  <h1 class="text-center">{{$t('Sign up')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light">
    <div class="row">
      <div class="col-6 text-center d-flex align-items-center border-right" v-if="app.wallet && metamaskLoginAvailable">
        <div class="mid align-self-center w-100">
          <h2 class="mb-3 font-weight-bold">{{$t('Wallet signature')}}</h2>
          <p class="light-text">{{$t('Use your MetaMask Wallet to sign up.')}}</p>
          <p class="text-danger" v-if="walletErrorMessage">{{ walletErrorMessage }}</p>
          <button class="btn btn-primary px-3" @click="loginWithSignature"
                  v-if="metamaskLoginAvailable">{{$t('Sign up with Metamask')}}
          </button>
        </div>
      </div>
      <div :class="{'col-6': app.wallet && metamaskLoginAvailable, 'col-12': !(app.wallet && metamaskLoginAvailable)}">
        <form v-show="!done" class="text-center" @submit.prevent="request">
          <div class="form-group mt-3 field-parent">
            <label for="inputEmail" class="sr-only">{{$t('Email address')}}</label>
            <input @input="cleanErr" type="text" id="inputEmail" ref="inputEmail" v-model.trim="email" name="email"
                   class="form-control"
                   :placeholder="$t('Email address')" required
                   autofocus>
          </div>
          <span class="text-muted"
                style="display: inline-block;">{{$t('Sign up explanation', 'Sign up by providing your email and clicking the button below.')}}</span>
          <button class="btn btn-primary px-3 mt-3" type="submit" id="signupcontent">{{$t('Sign up')}}</button>
        </form>
        <div v-show="done">
          <div
            class="my-3">{{$t('Sign up email sent explanation', 'Email sent. Please check your emails and visit the provided link to proceed.')}}
          </div>

          <a href="/" class="btn btn-primary" style="float: left;">{{$t('Home')}}</a>
          <a href="/login" class="btn btn-primary" style="float: right;">{{$t('Sign in')}}</a>
        </div>
      </div>
    </div>
  </div>
  <div class="modal fade" ref="tcModal" tabindex="-1" role="dialog">
    <div class="modal-dialog longer" role="document">
      <div class="modal-content breiter">
        <iframe frameborder="0" width="100%" height="100%" :src="$t('Terms & Conditions link', '')"></iframe>
        <div class="modal-footer">
          <button type="button" @click="acceptTC" class="btn btn-primary">{{$t('Accept')}}</button>
          <button type="button" class="btn btn-secondary" data-dismiss="modal">{{$t('Cancel')}}</button>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'RegisterRequest',
  data () {
    return {
      walletErrorMessage: '',

      account: undefined,
      email: '',
      password: '',
      metamaskLoginAvailable: false,
      loadingChallenge: false,
      challenge: null,
      done: false
    }
  },
  created () {
    if (window.web3 !== undefined) {
      this.metamaskLoginAvailable = true
    }
  },
  mounted () {
    this.$refs.inputEmail && this.$refs.inputEmail.focus()
    if (this.app.wallet) {
      this.account = this.app.wallet.getCurrentAddress()
    }
  },
  methods: {
    cleanErr () {
      $(this.$refs.inputEmail).cleanFieldErrors()
    },
    request () {
      axios.post('/api/register', { email: this.email }).then(res => {
        this.cleanErr()
        this.done = true
      }, (err) => {
        this.cleanErr()
        console.log(err)
        this.app.handleError(err)
        if (err.response && err.response.status === 422) {
          $(this.$refs.inputEmail).showFieldErrors({ errors: err.response.data })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Warning'),
            text: this.$t('There was an unexpected error. Please try again or if the error persists contact the platform operator.'),
            type: 'warning'
          })
        }
        this.$nextTick(() => {
          this.$refs.inputEmail.focus()
        })
      })
    },
    acceptTC () {
      localStorage.setItem('acc_' + this.account, 'yes')
      this.app.acknowledgeFirstLogin()
      $(this.$refs.tcModal).modal('hide')
      this.metamaskLogin()
    },
    checkTermsAndConditions () {
      if (this.$t('Terms & Conditions link', '') == 'Terms & Conditions link') {
        return true
      }
      let rememberAccept = localStorage.getItem('acc_' + this.account)
      if (rememberAccept && rememberAccept === 'yes') {
        return true
      }
      $(this.$refs.tcModal).modal('show')
    },
    loginWithSignature () {
      if (!this.challenge) {
        if (window.web3 !== undefined) {
          this.metamaskLoginAvailable = true
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
          await this.app.wallet.wallet.setupDefaultAccount()
        } catch (e) {
          console.log(e)
          this.walletErrorMessage = this.$t('Please grant access to MetaMask.')
          return
        }
      } else {
        this.walletErrorMessage = this.$t('Please grant access to MetaMask.')
        return
      }
      this.account = this.app.wallet.getCurrentAddress()
      if (this.account === undefined) {
        this.walletErrorMessage = this.$t('Please sign in to MetaMask.')
        return
      }
      if (this.checkTermsAndConditions()) {
        this.app.wallet.signMessage(this.challenge, this.account).then((signature) => {
          axios.post('/api/login', { signature }).then((res) => {
            this.challenge = ''
            if (res.status >= 200 && res.status <= 299) {
              window.location = res.data.location || '/admin/workflow'
            } else {
              this.walletErrorMessage = this.$t('Could not verify signature.')
            }
          }, (err) => {
            this.challenge = ''
            this.app.handleError(err)
            this.walletErrorMessage = this.$t('Could not verify signature.')
          })
        }).catch(() => {
          this.walletErrorMessage = this.$t('Could not Sign Message.')
        })
      }
    }
  }
}
</script>

<style lang="scss">
  @import "../assets/styles/variables.scss";

  .login-form {
    overflow: auto;
    margin: 0 auto;
    margin-top: 50px;
    height: 100%;
    max-width: 600px;
    padding-top: 40px;
    padding-bottom: 40px;
    border-radius: $border-radius;
  }

  .login-form-sm {
    max-width: 350px;
  }

  .form-signin {
    width: 100%;
    max-width: 330px;
    padding: 2rem;
    margin: 0 auto;
    z-index: 1000;

    .checkbox {
      font-weight: 400;
    }

    .form-control {
      position: relative;
      box-sizing: border-box;
      height: auto;
      padding: 10px;
      font-size: 16px;
    }
  }
</style>

<style scoped>
  .longer {
    width: 50%;
    height: 100%;
    max-width: none;
  }
  .breiter {
    height: 80%;
  }
</style>
