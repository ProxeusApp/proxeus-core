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
          <iframe v-if="app.userIsCreatorOrHigher()" frameborder="0" width="100%" height="100%" src="https://docs.google.com/document/d/e/2PACX-1vSp1Mk4mlyvxRf36qUXFWGFeqrCZEhTniMGQRCPfiGvgdS9pB1wDjED89LZzqqWbD72R9s388oAMJ-4/pub?embedded=true"></iframe>
          <iframe v-else frameborder="0" width="100%" height="100%" src="https://docs.google.com/document/d/e/2PACX-1vQfPfXcYSnQXvZC3FNvnAZ_D49xnmOQ8V2xv6ftNTbPgQXz6Jze-AUPtBU-XiQKmwlNbdypiVsyxOkm/pub?embedded=true"></iframe>
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
  mounted () {
    const firstLogin = sessionStorage.getItem('firstLogin')
    const delayedMessage = localStorage.getItem('showFirstLoginMessageOn')

    if (firstLogin || (delayedMessage && new Date(delayedMessage).getTime() < new Date().getTime())) {
      $(this.$refs.firstLoginModal).modal('show')
      sessionStorage.removeItem('firstLogin')
      localStorage.removeItem('showFirstLoginMessageOn')
    }
  },
  methods: {
    showInTheFuture: function () {
      const milliseconds = 1 * 60 * 60 * 1000 // hours * minutes * seconds * milliseconds
      const now = new Date()
      const showAt = new Date(now.getTime() + milliseconds)

      localStorage.setItem('showFirstLoginMessageOn', showAt)
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
