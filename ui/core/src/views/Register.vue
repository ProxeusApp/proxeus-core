<template>
<div>
  <!--<sign-message id="sign-message"></sign-message>-->
  <vue-headful :title="$t('Set password title', 'Proxeus - Set password')"/>
  <h1 class="text-center">{{$t('Set password')}}</h1>
  <div class="login-form container-fluid px-4 pt-2 mt-3 bg-light">
    <div class="row">
      <div class="col-12">
        <form v-show="!done" class="text-center" @submit.prevent="reset">
          <div class="form-group mt-3 field-parent">
            <label for="inputPassword" class="sr-only">{{$t('Password')}}</label>
            <input @input="cleanErr" type="password" id="inputPassword" ref="inputPassword" v-model="password"
                   name="password" class="form-control"
                   :placeholder="$t('Password')" required autofocus>
          </div>
          <div v-show="pwResetRequest" class="form-group">
            <a href="/register" class="" id="signupcontent">{{$t('Sign up')}}</a>
          </div>
          <span class="text-muted"
                style="display: inline-block;">{{$t('Set password explanation', 'Set a password by clicking the button below.')}}</span>
          <div style="width:100%;">
            <button class="btn btn-primary px-3 mt-3" type="submit">{{$t('Ok')}}</button>
          </div>
        </form>
        <div v-show="done">
          <div class="my-3">{{$t('Registration completed. You can sign in now.')}}</div>

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
  name: 'Register',
  data () {
    return {
      account: undefined,
      email: '',
      password: '',
      loadingChallenge: false,
      challenge: null,
      done: false,
      pwResetRequest: false
    }
  },
  created () {
  },
  mounted () {
    this.$refs.inputPassword && this.$refs.inputPassword.focus()
  },
  computed: {
    token () {
      try {
        return this.$route.params.token
      } catch (e) {}
      return ''
    }
  },
  methods: {
    cleanErr () {
      $(this.$refs.inputPassword).cleanFieldErrors()
    },
    reset () {
      axios.post('/api/register/' + this.token, { password: this.password }).then(res => {
        this.cleanErr()
        this.done = true
      }, (err) => {
        this.cleanErr()
        console.log(err)
        this.app.handleError(err)
        if (err.response && err.response.status === 422) {
          $(this.$refs.inputPassword).showFieldErrors({ errors: err.response.data })
        } else if (err.response && err.response.status === 400) {
          $(this.$refs.inputPassword).showFieldErrors({ errors: { password: [{ msg: this.$t('Please sign up first.') }] } })
          this.pwResetRequest = true
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Warning'),
            text: this.$t('There was an unexpected error. Please try again or if the error persists contact the platform operator.'),
            type: 'warning'
          })
        }
        this.$nextTick(() => {
          this.$refs.inputPassword.focus()
        })
      })
    }
  }
}
</script>

<style lang="scss">
  @use "../assets/styles/variables.scss";

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
