export default {
  computed: {
    app: {
      get () {
        // simple proxy
        return this.$root.$children[0]
      },
      set () {
        // read only!
      }
    }
  }
}
