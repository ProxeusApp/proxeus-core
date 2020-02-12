<template>
<span ref="trigger" class="ltrigger"><i v-show="loaderVisible" class="mdi mdi-loading mdi-spin"></i></span>
</template>

<script>
export default {
  name: 'trigger',
  props: {
    init: {
      type: Function,
      default: null
    },
    options: {
      type: Object,
      default: () => {
        return {
          // circumstances under which the observer's callback is invoked
          root: null, // defaults to the browser viewport if not specified or if null
          threshold: '1' // the degree of intersection between the target element and its root (0 - 1)
          // threshold of 1.0 means that when 100% of the target is visible within
          // the element specified by the root option, the callback is invoked
        }
      }
      // Whether you're using the viewport or some other element as the root,the API works the same way,
      // executing a callback function you provide whenever the visibility of the target element changes
      // so that it crosses desired amounts of intersection with the root
    }
  },
  data () {
    return {
      observer: null,
      loaderVisible: true
    }
  },
  mounted () {
    this.observer = new IntersectionObserver(entries => {
      this.handleIntersect(entries[0])
    }, this.options)
    if (this.init) {
      this.init(this.start, this.stop, this.hide)
    }
  },

  beforeDestroy () {
    this.stop()
  },

  methods: {
    handleIntersect (entry) {
      if (entry.isIntersecting) {
        this.show()
        this.$emit('trigger', this.start, this.stop, this.hide)
      }
    },
    stop () {
      this.hide()
      if (this.observer) {
        this.observer.disconnect()
      }
    },
    start () {
      this.show()
      if (this.observer && this.$refs.trigger) {
        this.observer.observe(this.$refs.trigger)
      }
    },
    show () {
      this.loaderVisible = true
    },
    hide () {
      this.loaderVisible = false
    }
  }
}
</script>

<style>
  .ltrigger {
    display: block;
    text-align: center;
    width: 100%;
    height: 10px;
    position: relative;
    margin: 0 auto;
  }

  .ltrigger .mdi-loading.mdi-spin {
    font-size: 40px;
    animation: mdi-spin 0.6s ease-in 0s infinite normal none running;
    top: -59px;
    z-index: 10;
    position: absolute;
    pointer-events: none;
  }
</style>
