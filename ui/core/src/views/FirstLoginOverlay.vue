<template>
  <div class="modal fade" ref="firstLoginModal" tabindex="-1" role="dialog">
    <div class="modal-dialog longer" role="document">
      <div class="modal-content breiter">
        <div class="modal-header">
          <h5 class="modal-title">{{$t('Welcome!')}}</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body px-3">
          <iframe frameborder="0" width="100%" height="100%" :src="previewUrl"></iframe>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-primary" data-dismiss="modal">{{$t('Dont show this again')}}</button>
          <button type="button" @click="showInTheFuture" class="btn btn-secondary" data-dismiss="modal">{{$t('Close')}}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'FirstLoginOverlay',
  props: {
    keyz: { type: String, required: true }, // A key to assign to the overlay. the overlay could be shown somewhere else,
    previewUrl: { type: String, required: true }
  },
  mounted () {
    const delayedMessage = localStorage.getItem('showFirstLoginMessageOn-' + this.$props.keyz)

    if (delayedMessage && new Date(delayedMessage).getTime() < new Date().getTime()) {
      $(this.$refs.firstLoginModal).modal('show')
      localStorage.removeItem('showFirstLoginMessageOn-' + this.$props.keyz)
    }
  },
  methods: {
    showInTheFuture: function () {
      const milliseconds = 1 * 60 * 60 * 1000 // hours * minutes * seconds * milliseconds
      const now = new Date()
      const showAt = new Date(now.getTime() + milliseconds)

      localStorage.setItem('showFirstLoginMessageOn-' + this.$props.keyz, showAt)
      console.log('Will show the message again later...', showAt)

      return true
    }
  }
}
</script>

<style scoped>
  .longer {
    width: 50%;
    height: 100%;
    max-width: none;
  }
  .breiter {
    height: 80%;
  }
</style>
