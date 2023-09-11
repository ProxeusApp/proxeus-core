<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Sign in title', 'Proxeus - Log in')"/>
  <h1 class="text-center">{{$t('Sign in','Log in')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light" :class="{'login-form-sm': !metamaskLoginAvailable}">
    <div class="row">
      <div class="col-6 mt-3 text-center d-flex align-items-center border-right" v-if="app.wallet && metamaskLoginAvailable">
        <div class="mid align-self-center w-100">
          <h2 class="mb-3 font-weight-bold">{{$t('Wallet signature')}}</h2>
          <p class="light-text">{{$t('Use your MetaMask Wallet to log in.')}}</p>
          <p class="text-danger" v-if="walletErrorMessage">{{ walletErrorMessage }}</p>
          <button class="btn btn-primary px-3" @click="loginWithSignature"
                  v-if="metamaskLoginAvailable">{{$t('Log in with Metamask')}}
          </button>
        </div>
      </div>
      <div :class="{'col-6': metamaskLoginAvailable, 'col-12': metamaskLoginAvailable === false}">
        <form class="text-center" @submit.prevent="login">
          <div class="form-group mt-3">
            <label for="inputEmail" class="sr-only">{{$t('Email address')}}</label>
            <input type="text" id="inputEmail" ref="inputEmail" v-model="email" name="email" class="form-control"
                   :placeholder="$t('Email address')" required autofocus>
          </div>
          <div class="form-group">
            <label for="inputPassword" class="sr-only">{{$t('Password')}}</label>
            <input type="password" id="inputPassword" v-model="password" name="password" class="form-control mt-2"
                   :placeholder="$t('Password')" required>
          </div>
          <div class="text-danger mb-3" v-show="hasError">{{ loginErrorMessage }}</div>
          <div class="form-group">
            <router-link :to="{path:'/reset/password'}" class="" data-toggle="tooltip"
                         data-boundary="window" :title="$t('Forgot password')">
              {{$t('Forgot password')}}
            </router-link>
            |
            <router-link :to="{path:'/register'}" class="" data-toggle="tooltip"
                         data-boundary="window" :title="$t('Sign up','Log in')">
              {{$t('Sign up','Log in')}}
            </router-link>
          </div>
          <button class="btn btn-primary px-3" type="submit">{{$t('Sign in','Log in')}}</button>
        </form>
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
  name: 'AdminLogin',
  data () {
    return {
      account: undefined,
      pwlogin: false,
      email: '',
      password: '',
      hasError: false,
      walletErrorMessage: '',
      loginErrorMessage: '',
      metamaskLoginAvailable: false,
      challenge: null
    }
  },
  created () {
    if (typeof window.ethereum !== 'undefined') {
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
    login () {
      this.account = this.email
      this.pwlogin = true
      if (this.checkTermsAndConditions()) {
        axios.post('/api/login', { email: this.email, password: this.password }).then(res => {
          this.app.initUserHasSession()
          this.hasError = false
          window.location = res.data.location || '/admin/workflow'
        }, (err) => {
          this.app.deleteUserHasSession()
          this.app.handleError(err)
          this.loginErrorMessage = this.$t('You have entered an invalid username or password')
          this.hasError = true
          this.$nextTick(() => {
            this.$refs.inputEmail.focus()
          })
        })
      }
    },
    checkTermsAndConditions () {
      if (this.$t('Terms & Conditions link', '') === 'Terms & Conditions link') {
        return true
      }
      const rememberAccept = localStorage.getItem('acc_' + this.account)
      if (rememberAccept && rememberAccept === 'yes') {
        return true
      }
      $(this.$refs.tcModal).modal('show')
    },
    acceptTC () {
      localStorage.setItem('acc_' + this.account, 'yes')
      this.app.acknowledgeFirstLogin()
      $(this.$refs.tcModal).modal('hide')

      if (this.pwlogin) {
        this.login()
      } else {
        this.metamaskLogin()
      }
    },
    loginWithSignature () {
      if (!this.challenge) {
        if (typeof window.ethereum !== 'undefined') {
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
      this.pwlogin = false
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

  .modal-body {
    h1 {
      line-height: 48px;
    }
  }

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
