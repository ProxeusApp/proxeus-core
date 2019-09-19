import Vue from 'vue'

const state = {
  signatureRequestCount: 0
}

const mutations = {
  SET_SIGNATURE_REQUEST_COUNT (state, { signatureRequestCount }) {
    state.signatureRequestCount = signatureRequestCount
  }
}

const getters = {
  signatureRequestCount: (state) => {
    return state.signatureRequestCount
  }
}

const actions = {
  UPDATE_SIGNERS_COUNT ({ commit }, { sigCount }) {
    commit('SET_SIGNATURE_REQUEST_COUNT', { signatureRequestCount: sigCount })
  }
}

export default {
  state, mutations, actions, getters
}
