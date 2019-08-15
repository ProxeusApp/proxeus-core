export default {
  components: {
  },
  actionable: {
    created: function () {
      this.hello()
    },
    methods: {
      hello: function () {
        console.log('hello from mixin!')
        this.settings.label.value = 'wtf! it works'
      }
    }
  }
}
