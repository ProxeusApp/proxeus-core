<template>
  <div>
    <div class="d-flex mb-3 align-items-center" v-if="!configOnly">
      <button class="btn btn-primary" :class="{active:createNew}" @click="tabClick(true)"
              @focus="$event.target.blur();">New
      </button>
      <button class="btn btn-primary ml-2" :class="{active:!createNew}" @click="tabClick(false)"
              @focus="$event.target.blur();">Import
      </button>
    </div>
    <div v-if="createNew" class="tabcontent" ref="inputs">
      <div v-if="newFormReady()" class="init-settings form-group" ref="fields">
        <nav-tabs class="mt-0">
          <tab v-if="app.me === null" :selected="app.me === null" title="Initial user settings (will be ignored if email exists already)">
            <animated-input name="user.email" :max="100" :label="$t('New Email')" v-model="user.email">
              <span
                class="text-muted">{{$t('Initial user email explanation','Set the email of the initial user.')}}</span>
            </animated-input>
            <div class="form-group">
              <simple-select :unselect="false" style="margin-top: 15px;" name="user.role" v-model="user.role"
                             :idProp="'role'"
                             :labelProp="'name'" :options="app.roles"/>
              <span
                class="text-muted">{{$t('Initial user role explanation','Set the role for the initial user, root is recommended.')}}</span>
            </div>
            <animated-input type="password" name="user.password" :max="100" :label="$t('Password')"
                            v-model="user.password">
              <span class="text-muted"
                    style="white-space: normal;">{{$t('Initial user password explanation','Set the password of the initial user.')}}</span>
            </animated-input>
          </tab>
          <tab title="System settings" :selected="app.me !== null">
            <animated-input name="settings.dataDir" :max="100" :label="$t('Data dir','Data directory')"
                            v-model="settings.dataDir">
              <div class="alert-danger">{{$t('Warning data dir', 'Warning: by changing the data directory you will loose all the data like the users, workflows, etc...Do not do it unless you are aware of the implications.')}}
              </div>
              <span class="text-muted">{{$t('Data dir explanation','Set the database directory path. All data will be stored here.')}}</span>
            </animated-input>
            <animated-input name="settings.sessionExpiry" :max="100" :label="$t('Session expiry')"
                            v-model="settings.sessionExpiry">
              <span class="text-muted">{{$t('Session expiry explanation','Set the session expiry like 1h as one hour, 1m as one minute or 1s as one second.')}}</span>
            </animated-input>
            <animated-input name="settings.cacheExpiry" :max="100" :label="$t('Cache expiry')"
                            v-model="settings.cacheExpiry">
              <span class="text-muted">{{$t('Cache expiry explanation','Set the common cache expiry which will be used for email tokens or similar like 1h as one hour, 1m as one minute or 1s as one second.')}}</span>
              <simple-select :unselect="false" style="margin-top: 15px;" name="settings.defaultRole"
                             v-model="settings.defaultRole" :idProp="'role'" :labelProp="'name'" :options="defaultUserRoles"/>
              <span class="text-muted">{{$t('Default role explanation','Select the default role that is going to be used for new registrations.')}}</span>
            </animated-input>
            <animated-input name="settings.documentServiceUrl" :max="100" :label="$t('Document Service URL')"
                            v-model="settings.documentServiceUrl">
              <span class="text-muted">{{$t('Document Service URL explanation','Set the Document Service URL which will be used to render documents.')}}</span>
            </animated-input>
            <animated-input name="settings.platformDomain" :max="100" :label="$t('Platform Domain')"
                            v-model="settings.platformDomain">
              <span class="text-muted">{{$t('Platform Domain explanation','Set the Domain this Platform instance is identifying as (used for example for sending links to this instance)')}}</span>

            </animated-input>
            <div class="d-none">
              <animated-input name="settings.defaultWorkflowIds" :label="$t('Default Workflows')"
                              v-model="settings.defaultWorkflowIds">
                <span class="text-muted">{{$t('Default workflow ids explanation','Comma separated ids of workflows you want your new users to inherit (if any)')}}</span>
              </animated-input>
            </div>
          </tab>
          <tab title="Blockchain settings">
            <animated-input name="settings.blockchainNet" :max="100" :label="$t('Blockchain net')"
                            v-model="settings.blockchainNet">
              <span class="text-muted">{{$t('Blockchain net explanation','Set the ethereum blockchain net like mainnet/goerli/sepolia/polygon-mumbai/polygon.')}}</span>
            </animated-input>
            <animated-input name="settings.infuraApiKey" :max="100" :label="$t('Infura API Key')"
                            v-model="settings.infuraApiKey">
              <span class="text-muted">{{$t('Infura API Key explanation','API Key to access Infura node.')}}</span>
            </animated-input>
            <animated-input name="settings.blockchainContractAddress" :max="100"
                            :label="$t('Blockchain contract address')" v-model="settings.blockchainContractAddress">
              <span class="text-muted">{{$t('Blockchain contract address explanation','Set the ethereum contract address which will be used to register files and verify them.')}}</span>

            </animated-input>
            <div class="form-group">
              <simple-select :unselect="false" style="margin-top: 15px;" name="settings.airdropEnabled"
                             v-model="settings.airdropEnabled" :idProp="'value'" :labelProp="'label'"
                             :options="this.airdropoptions"/>
              <span class="text-muted">{{$t('Airdrop Enable Explanation','Enables/Disables the XES & Ether airdrop feature for new users on testnet (sepolia, goerli, mumbai). The Amount and Wallet to be used is configured in the platform configuration.')}}</span>
            </div>
            <animated-input :disabled="settings.airdropEnabled!='true'" name="settings.airdropAmountXES" :max="100"
                            :label="$t('Airdrop Amount XES')" v-model="settings.airdropAmountXES">
              <span class="text-muted">{{$t('Airdrop Amount XES Explanation','Set the amount of XES to be airdropped to newly registered users.')}}</span>
            </animated-input>
            <animated-input :disabled="settings.airdropEnabled!='true'" name="settings.airdropAmountEther" :max="100"
                            :label="$t('Airdrop Amount Ether')" v-model="settings.airdropAmountEther">
              <span class="text-muted">{{$t('Airdrop Amount Ether Explanation','Set the amount of Ether to be airdropped to newly registered users.')}}</span>
            </animated-input>
          </tab>
          <tab title="E-Mail settings">
            <animated-input name="settings.emailFrom" :max="100" :label="$t('Email from')"
                            v-model="settings.emailFrom">
              <span class="text-muted">{{$t('Email from explanation','Set the email that is being used to send out emails.')}}</span>
            </animated-input>
            <animated-input name="settings.sparkpostApiKey" :max="100" :label="$t('Sparkpost API Key')"
                            v-model="settings.sparkpostApiKey">
              <span class="text-muted">{{$t('Sparkpost API Key explanation','Set the Sparkpost API key which will be used to send out emails.')}}</span>
            </animated-input>
          </tab>
        </nav-tabs>
      </div>
      <div v-else>
        {{$t('loading...')}}
      </div>
    </div>
    <div v-else class="tabcontent">
      <import-only></import-only>
    </div>
    <button v-if="configOnly && !importResultsAvailable" type="button" @click="powerUp"
            class="btn btn-primary">
      Save
    </button>
    <button v-else-if="createNew && !importResultsAvailable" type="button" @click="powerUp"
            class="btn btn-primary">
      Power up
    </button>
    <button v-if="importResultsAvailable && configured" type="button" @click="redirectToHome"
            class="btn btn-primary">
      Redirect to home
    </button>
  </div>
</template>

<script>
import AnimatedInput from '@/components/AnimatedInput'

import SimpleSelect from '@/components/SimpleSelect'
import ImportOnly from '../../components/ImportOnly'

import mafdc from '@/mixinApp'

import NavTabs from '@/components/nav-tabs/NavTabs'
import Tab from '@/components/nav-tabs/Tab'

export default {
  mixins: [mafdc],
  name: 'settings-inner',
  props: {
    configOnly: {
      type: Boolean,
      default: true
    }
  },
  components: {
    ImportOnly,
    SimpleSelect,
    AnimatedInput,
    NavTabs,
    Tab
  },
  created () {
    this.getInit(this.setInitialState)
  },
  mounted () {
    if (this.app.me === null) {
      this.$root.$on('lastImportResults', this.checkIfWeReceiveImportData)
    }
  },
  beforeDestroy () {
    if (this.app.me === null) {
      this.$root.$off('lastImportResults', this.checkIfWeReceiveImportData)
    }
  },
  updated () {
  },
  computed: {
    defaultUserRoles () {
      return this.app.roles.map(roleObj => { return { name: roleObj.name, role: roleObj.name.toLowerCase() } })
    }
  },
  methods: {
    getInit (cb) {
      axios.get('/api/init').then(res => {
        if (res.data) {
          this.settings = res.data.settings
          this.configured = res.data.configured
          if (cb) {
            cb()
          }
        }
      }, (err) => {
        if (err.response && err.response.status !== 401) {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not load default system settings', 'Could not load default values'),
            type: 'error'
          })
        }
      })
    },
    setInitialState () {
      this.initialized = this.configured
    },
    tabClick (createNew) {
      if (this.configured) {
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Already configured'),
          type: 'success'
        })
        return
      }
      this.createNew = createNew
    },
    checkIfWeReceiveImportData (results) {
      if (!this.app.me && !this.configured) {
        this.results = results
        this.getInit(this.shouldShowPowerUpBtn)
      }
    },
    shouldShowPowerUpBtn () {
      if (this.app.me === null && !this.initialized && this.configured && this.results) {
        this.importResultsAvailable = true
        this.initialized = true
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Initial screen get started', 'Click on the bottom right button to get started'),
          type: 'success'
        })
      }
      this.results = null
    },
    refresh () {
      setTimeout(this.redirectToHome, 4000)
    },
    redirectToHome () {
      window.location.href = '/'
    },
    newFormReady () {
      return this.app.roles && this.app.roles.length > 0 && this.settings !== null
    },
    cleanErr () {
      $(this.$refs.inputs).cleanFieldErrors()
    },
    powerUp () {
      axios.post('/api/init', { settings: this.settings, user: this.user }).then(res => {
        this.cleanErr()
        this.user = res.data
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Saved successfully'),
          type: 'success'
        })
        if (this.app.me === null) {
          this.refresh()
        } else {
          this.app.loadConfig()
        }
      }, (err) => {
        this.cleanErr()
        if (err.response && err.response.status === 422) {
          $(this.$refs.inputs).showFieldErrors({ errors: err.response.data })
        } else {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not save. Please try again or if the error persists contact the platform operator.\n'),
            type: 'error'
          })
        }
      })
    }
  },
  data () {
    return {
      user: { role: 100 },
      settings: null,
      createNew: true,
      configured: false,
      initialized: false,
      importResultsAvailable: false,
      results: null,
      airdropoptions: [
        {
          label: 'Airdrop enabled on Testnet (sepolia or goerli or mumbai)',
          value: 'true'
        },
        {
          label: 'Airdrop disabled',
          value: 'false'
        }
      ]
    }
  }
}
</script>

<style lang="scss">
  .tabbtn {
    margin-bottom: -2px;
    border: 1px solid #cecece;
    border-bottom: none !important;
  }

  .tabbtn.active {
    border: 1px solid #062a85;
    background: white;
  }

  .spanel > .spanel-title {
    position: absolute;
    background: white;
    top: -16px;
    font-size: 18px;
  }

  .spanel {
    padding: 15px;
    border: 1px solid #dedede;
    position: relative;
    margin-bottom: 15px;
  }
</style>
