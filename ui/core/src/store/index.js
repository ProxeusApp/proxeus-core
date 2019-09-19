import Vue from 'vue'
import Vuex from 'vuex'
import forms from './forms'
import signatures from './signatures'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  modules: {
    forms,
    signatures
  },
  strict: debug
})
