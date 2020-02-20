<template>
  <div>
    <div class="new-api-key" style="position:relative;">
        <div v-show="createNew">
            <form ref="form" v-on:submit.prevent="createNewApiKey">
                <animated-input :max="80" :label="$t('Name')" v-model="newName"/>
                <span class="text-muted" style="white-space: normal;">{{$t('Create a API key')}}</span>
                <button type="button" style="position: absolute;right: 0;top: -3px;" @click="createNewApiKey" :disabled="!newName" class="btn btn-primary">
                    Save
                </button>
            </form>
        </div>
        <div v-if="result" class="created-api-key" style="
        position: relative;
    background: #acf9ee5e;
    padding: 30px;">
            <span @click="clearResult" style="    position: absolute;
    right: 5px;
    top: 3px;
    padding: 2px 5px;
    cursor: pointer;
    background: #00000017;
    font-weight: bold;
    border: 1px solid #d6d6d6;
    border-radius: 44%;
    line-height: 1;">X</span>
            <span style="margin-right: 10px;">{{result.Name}}</span>
            <span style="font-weight: bold;">{{result.Key}}</span>
            <p>
                <span class="text-muted" style="white-space: normal;">{{$t('Please copy this API key now as it will be not fully readable afterwards.')}}</span>
            </p>
        </div>
    </div>
      <div style="position:relative;">
          <button v-if="app.me && (app.me.id === user.id)" :disabled="createNew" style="padding:4px;position:absolute;right:0;top:-5px;" :title="$t('Create a API key')" type="button" @click="add" class="btn btn-sm btn-primary">
              Add key
          </button>
          <div class="sub-title">{{$t('API keys')}}</div>
          <table class="table hidden-api-keys" style="width: 100%;margin: 4px;">
              <tr v-for="item in user.apiKeys" >
                  <td style="text-align: left;">
                      <span>{{item.Name}}</span>
                  </td>
                  <td style="text-align: left;">
                      <span>{{item.Key}}</span>
                  </td>
                  <td v-if="app.me && (app.me.role >= 100 || app.me.id === user.id)" style="text-align: right;padding-right: 10px;">
                      <button style="float:right;padding:2px;" :title="$t('Delete API key')" type="button" @click="deleteApiKey(item.Key);" class="btn btn-danger btn-sm">
                          Remove
                      </button>
                  </td>
              </tr>
          </table>
      </div>
  </div>
</template>

<script>
import AnimatedInput from '../AnimatedInput'
import Checkbox from '../Checkbox'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'api-key',
  components: {
    Checkbox,
    AnimatedInput
  },
  props: {
    user: {
      type: Object,
      default: {}
    }
  },
  data () {
    return {
      newName: '',
      result: null,
      createNew: false
    }
  },
  created () {
  },
  mounted () {
  },
  computed: {
  },
  methods: {
    add () {
      this.createNew = true
      this.clearResult()
      this.$nextTick(() => {
        if (this.$refs.form) {
          $(this.$refs.form).find('input').first().focus()
        }
      })
    },
    reloadKeys () {
      let include = {}
      include[this.user.id] = true
      axios.post('/api/admin/user/list', { include: include, limit: 1 }).then(response => {
        if (response.data && response.data.length === 1 && response.data[0].id === this.user.id) {
          let usr = response.data[0]
          this.user.apiKeys = usr.apiKeys
        }
      }, (err) => {
        this.app.handleError(err)
      })
    },
    clearResult () {
      this.result = null
    },
    deleteApiKey (hiddenApiKey) {
      axios.delete('/api/user/create/api/key/' + this.user.id + '?hiddenApiKey=' + hiddenApiKey).then(response => {
        this.user.apiKeys = this.user.apiKeys.filter(item => item.Key !== hiddenApiKey)
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('The API-Key was successfully deleted.'),
          type: 'success'
        })
        this.reloadKeys()
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not delete API-Key. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    },
    createNewApiKey () {
      if (this.newName === '') {
        return
      }
      axios.get('/api/user/create/api/key/' + this.user.id + '?name=' + this.newName).then(response => {
        this.result = { 'Name': this.newName, 'Key': response.data }
        this.newName = ''
        this.createNew = false
        this.reloadKeys()
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('The API-Key was successfully saved.'),
          type: 'success'
        })
      }, (err) => {
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not save API-Key. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    }
  }
}
</script>

<style>

</style>
