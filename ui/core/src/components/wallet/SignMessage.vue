<template>
<div class="modal fade" data-backdrop="false" v-bind="$attrs" tabindex="-1" role="dialog" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header align-items-center bg-light px-3">
        <h5 class="modal-title font-weight-bold">Proxeus Wallet</h5>
        <div class="balance ml-auto">28 XES</div>
        <button type="button" class="btn btn-link text-primary ml-1 pr-0">
          <i class="material-icons mdi mdi-settings"></i>
        </button>
      </div>
      <div class="modal-body px-3">
        <h2 class="text-center">Sign message</h2>
        <h3 class="message">Message:</h3>
        <p>Please sign this message to login to the PROXEUS App.</p>
        <h3 class="fees">Transaction fees:</h3>
        <p>0 XES</p>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button>
        <button type="button" class="btn btn-primary" data-dismiss="modal">Sign message</button>
      </div>
    </div>
  </div>
</div>
</template>

<script>

export default {
  name: 'sign-message'
}
</script>

<style scoped>
  .modal-dialog {
    box-shadow: 0 40px 100px rgba(0, 0, 0, .35);
    top: 50px;
  }

</style>
