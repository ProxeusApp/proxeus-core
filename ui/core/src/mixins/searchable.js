export default {
  data () {
    return {
      searchTimeout: undefined
    }
  },
  methods: {
    timeoutSearch (term, cb, timeout = 100) {
      clearTimeout(this.searchTimeout)

      this.searchTimeout = setTimeout(() => {
        cb(term)
      }, timeout)
    }
  }
}
