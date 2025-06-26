<template>
<div class="modal fade bd-example-modal-lg" :id="mid" tabindex="-1" role="dialog"
     aria-labelledby="myLargeModalLabel"
     aria-hidden="true">
  <div class="modal-dialog modal-lg">
    <div class="modal-content" v-if="src">
      <spinner v-show="loadingPdf" background="transparent"></spinner>
      <a v-if="download && filename" :href="src" class="btn btn-primary"><span class="align-middle"><small
        class="text-truncate file-name text-white" v-if="filename">{{ filename }}</small></span>
        <i class="mdi mdi-download ml-2"></i>
      </a>
      <VuePdfEmbed
        v-for="i in numPages"
        :key="i"
        :source="pdfSrc"
        @loading-failed="errorHandler"
        :page="i"
      ></VuePdfEmbed>
    </div>
  </div>
</div>
</template>

<script>
import VuePdfEmbed from 'vue-pdf-embed'
import Spinner from '@/components/Spinner'

export default {
  name: 'pdf-modal',
  props: ['src', 'mid', 'filename', 'download'],
  components: {
    VuePdfEmbed,
    Spinner
  },
  data () {
    return {
      numPages: undefined,
      pdfSrc: undefined,
      loadingPdf: false
    }
  },
  computed: {},
  methods: {
    load () {
      this.loadingPdf = true
      this.pdfSrc = undefined
      this.numPages = undefined
      this.pdfSrc = pdf.createLoadingTask(this.src)

      this.pdfSrc.promise.then(pdf => {
        this.numPages = pdf.numPages
        this.loadingPdf = false
        this.pdfSrc._worker.destroy()
      }, (err) => {
        this.loadingPdf = false
        console.log(err)
      })
    },
    errorHandler (e) {
      console.log(e)
    }
  }
}
</script>

<style lang="scss" scoped>
  .modal-content {
    border: none;
    background: transparent;

    > div {
      margin-bottom: 5px;
    }
  }
</style>
