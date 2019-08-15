<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Reset password title', 'Proxeus - Reset password')"/>
  <h1 class="text-center">{{$t('Reset password')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light">
    <div class="row">
      <div class="col-12">
        <form v-show="!done" class="text-center" @submit.prevent="request">
          <div class="form-group mt-3 field-parent">
            <label for="inputEmail" class="sr-only">{{$t('Email address')}}</label>
            <input @input="cleanErr" type="text" id="inputEmail" ref="inputEmail" v-model.trim="email" name="email"
                   class="form-control"
                   :placeholder="$t('Email address')" required
                   autofocus>
          </div>
          <span class="text-muted"
                style="display: inline-block;">{{$t('Password reset request explanation', 'Request a password reset by providing your email and clicking the button below.')}}</span>
          <button class="btn btn-primary px-3 mt-3" type="submit">{{$t('Request')}}</button>
        </form>
        <div v-show="done">
          <div class="my-3">{{$t('Email sent.')}}</div>

          <a href="/" class="btn btn-primary" style="float: left;">{{$t('Home')}}</a>
          <a href="/login" class="btn btn-primary" style="float: right;">{{$t('Sign in')}}</a>
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
  name: 'PasswordResetRequest',
  data () {
    return {
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
  },
  mounted () {
    this.$refs.inputEmail && this.$refs.inputEmail.focus()
  },
  methods: {
    cleanErr () {
      $(this.$refs.inputEmail).cleanFieldErrors()
    },
    request () {
      axios.post('/api/reset/password', { email: this.email }).then(res => {
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
            text: this.$t('Something wrong happened, please try again.'),
            type: 'warning'
          })
        }
        this.$nextTick(() => {
          this.$refs.inputEmail.focus()
        })
      })
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
