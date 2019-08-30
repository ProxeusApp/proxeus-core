<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Change Email title', 'Proxeus - Change Email')"/>
  <h1 class="text-center">{{$t('Change Email')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light">
    <div class="row">
      <div class="col-12">
        <div v-show="badRequest" class="form-group">
          <span class="error">{{$t('Please validate your email first.')}}</span>
          <a href="/change/email" class="">{{$t('Request a link for validation')}}</a>
        </div>
        <div v-show="done">
          <div class="my-3">{{$t('Email changed. You can use the new one now.')}}</div>

          <a href="/" class="btn btn-primary" style="float: left;">{{$t('Home')}}</a>
          <a href="/login" class="btn btn-primary" style="float: right;">{{$t('Sign in','Log in')}}</a>
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
  name: 'EmailChange',
  data () {
    return {
      account: undefined,
      email: '',
      password: '',
      metamaskLoginAvailable: false,
      loadingChallenge: false,
      challenge: null,
      done: false,
      badRequest: false
    }
  },
  created () {
  },
  mounted () {
    axios.post('/api/change/email/' + this.token, { password: this.password }).then(res => {
      this.done = true
      if (this.app.me) {
        this.app.loadMe()
      }
    }, (err) => {
      this.app.handleError(err)
      if (err.response && err.response.status === 422) {
      } else if (err.response && err.response.status === 400) {
        this.badRequest = true
      } else {
        this.$notify({
          group: 'app',
          title: this.$t('Warning'),
          text: this.$t('An unexpected error occurred. Please try again or if the error persists contact the platform operator.'),
          type: 'warning'
        })
      }
    })
  },
  computed: {
    token () {
      try {
        return this.$route.params.token
      } catch (e) {
      }
      return ''
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
