<template>
  <div class="modal fade smodal" :id="mid" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel"
       aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ $t('Request Signature', 'Request Signature') }}</h5>
          <button type="button" class="close ml-2" data-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body">
          <user-selector :excludes="getMapOfExistingUsersInPermItem"
                         @added="grantAdded"
                         :dependencyFulfilled="isGrantSelectorNotEmpty"
                         v-model="granted"
                         :uri="'/api/admin/user/list'"/>
        </div>
        <div class="modal-footer p-2">
          <button class="btn btn-primary ml-2" @click="cancel">{{ $t('Cancel', 'Cancel') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import UserSelector from '../../views/appDependentComponents/permDialog/UserSelector.vue'

export default {
  name: 'sig-modal',
  props: ['src', 'mid', 'docid'],
  components: {
    UserSelector
  },
  data () {
    return {
      signatory: null,
      item: null,
      granted: null,
      owner: null,
      grantReadWriteSelect: 1
    }
  },
  methods: {
    send (addr) {
      this.$notify({
        group: 'app',
        title: this.$t('Signature request'),
        text: this.$t('Requesting signature...'),
        type: 'info'
      })
      var bodyFormData = new FormData()
      bodyFormData.set('signatory', addr)
      axios({
        method: 'post',
        url: this.src.replace('/file', '/signingRequests').replace('?format=pdf', '/add'),
        data: bodyFormData,
        config: { headers: { 'Content-Type': 'multipart/form-data' } }
      }).then(response => {
        this.$notify({
          group: 'app',
          title: this.$t('Signature request'),
          text: this.$t('Signature request successfully sent.'),
          type: 'success'
        })
        $('#signingModal').modal('hide')
        this.signatory = ''
        this.$router.go()
      }, (err) => {
        console.log(err)
        this.$notify({
          group: 'app',
          title: this.$t('Failed to request signature'),
          text: this.$t('The signatory has no either no Ethereum address defined or was already requested.'),
          type: 'error',
          duration: 5000
        })
      }).catch(e => {
        console.log(e)
      })
    },
    cancel () {
      this.signatory = ''
      $('.smodal').modal('hide')
    },
    grantAdded (usrList) {
      this.send(usrList[0].etherPK)
    },
    isGrantSelectorNotEmpty () {
      return this.grantReadWriteSelect !== null && this.grantReadWriteSelect !== undefined
    },
    onGrantSelect (item, id) {
      this.signatory = item
    },
    deleteGrant (id) {
      this.granted = this.granted.filter(item => item.id !== id)
    },
    OnSelectedShowIconOnly (strEl, change) {
      const e = $(strEl)
      e.find('.my-explanation').remove()
      change($('<div>').append(e).html())
    },
    getMapOfExistingUsersInPermItem () {
      return null
    },
    refreshGrantList () {
      this.loading = true
    }
  }
}
</script>

<style scoped>
  .modal-dialog{
    top: 55px
  }
  .modal-body p{
    margin: 25px
  }
  .modal-body p textarea{
    width: 100%
  }

</style>
