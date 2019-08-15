<template>
<div class="wallet-login text-center px-3 pt-3">
  <h2>Login to Proxeus Wallet</h2>
  <div v-show="authenticating === false">
    <div class="form-group">
      <input type="password" class="form-control" v-model="password" name="password" placeholder="Password">
    </div>
    <button type="button" @click="authenticate" class="btn btn-primary">Login</button>
    <button type="button" class="btn btn-link btn-restore mt-2" @click="restoreFromSeed">Restore from seed</button>
  </div>
  <div class="spinner-wrapper">
    <spinner background="transparent" color="#eee" v-show="authenticating" cls="position-relative"
             :margin="0"></spinner>
  </div>
</div>
</template>

<script>
import Spinner from '../Spinner'

export default {
  name: 'login',
  props: {
    wallet: {
      type: Object,
      required: true
    }
  },
  components: {
    Spinner
  },
  created () {
    this.authenticating = true
    if (this.restoreAuth()) {
      this.authenticating = false
      this.$emit('authenticated')
    } else {
      this.authenticating = false
    }
  },
  data () {
    return {
      password: '',
      encryptedKeystore: undefined,
      authenticating: false
    }
  },
  methods: {
    authenticate () {
      this.authenticating = true
      this.wallet.setupKeystore(this.password).then((encryptedKeystore) => {
        this.storeAuth({ password: this.password, encryptedKeystore })
        // import a private key, also async
        this.wallet.importPrivateKey('792ab5748955ea50d1641055615400da75daed092de6f95ae90bc876e4010717').then(() => {
          this.authenticating = false
          this.$emit('authenticated')
        }).catch((e) => {
          this.authenticating = false
        })
      }).catch((e) => {
        this.authenticating = false
      })
    },
    storeAuth (authObj) {
      let authStr = btoa(JSON.stringify(authObj))
      console.log(authStr)
      this.$cookie.set('test', 'hallo')
      this.$cookie.set('mnidm', authStr, {
        secure: window.location.hostname !== 'localhost',
        expires: '1h',
        domain: window.location.hostname,
        samesite: 'Strict'
      })
    },
    restoreAuth () {
      let authObj = this.$cookie.get('mnidm')
      if (authObj) {
        authObj = JSON.parse(atob(authObj))
        this.password = authObj.password || ''
        this.encryptedKeystore = authObj.encryptedKeystore || undefined
      }
      return !!(this.password.length && this.encryptedKeystore)
    },
    restoreFromSeed () {

    }
  }
}
</script>

<style scoped>
  .btn-restore {
    text-decoration: underline;
  }
</style>
