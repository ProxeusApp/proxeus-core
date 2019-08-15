export default {
  data () {
    return {
      inputTimeout: undefined
    }
  },
  methods: {
    delayedInputEvent (cb, timeout = 100) {
      clearTimeout(this.inputTimeout)
      this.inputTimeout = setTimeout(() => {
        cb()
      }, timeout)
    }
  }
}
