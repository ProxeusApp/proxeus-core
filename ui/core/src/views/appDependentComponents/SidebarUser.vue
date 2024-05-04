<template>
<responsive-sidebar class="sidebar-light">
  <table style="height: 100%;width: 100%;">
    <tr>
      <td valign="top" style="vertical-align:top">
        <a class="navbar-brand m-0 w-100 p-0 m-0" href="#">
          <img :src="$t('Sidebar User Logo','/static/proxeus-logo.svg')" alt="" class="d-inline-block align-top">
        </a>
        <nav class="collapse show sidebar-sticky">
          <nav class="nav main-nav flex-column main-nav-frontend">
            <li class="nav-item">
              <router-link :to="{name:'Documents'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu Documents','Documents')"><span
                class="material-icons mdi mdi-view-carousel"></span><span
                class="nav-link-title">{{$t('Menu Documents','Documents')}}</span>
              </router-link>
            </li>
            <li class="nav-item">
              <router-link :to="{name:'DocumentVerification'}" class="nav-link" data-toggle="tooltip"
                           data-placement="right"
                           data-boundary="window" :title="$t('Menu Verification','Verification')"><span
                class="material-icons mdi mdi-shield-check"></span><span
                class="nav-link-title">{{$t('Menu Verification','Verification')}}</span>
              </router-link>
            </li>
            <li class="nav-item">
              <router-link :to="{name:'SignatureRequests'}" class="nav-link" data-toggle="tooltip"
                           data-placement="right"
                           data-boundary="window" :title="$t('Menu Signature Requests','Signature Requests')"><span
                class="material-icons mdi mdi-pen"></span><span
                class="nav-link-title">{{$t('Menu Signature Requests','Signature Requests')}} ({{ signatureRequestCount }})</span>
              </router-link>
            </li>
          </nav>
        </nav>
      </td>
    </tr>
    <tr>
      <td valign="bottom" style="vertical-align:bottom">
        <ul class="nav secondary-nav">
          <li class="nav-item" v-if="userCanAccessBackend">
            <router-link :to="{name:'UserImportExport'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                         data-boundary="window" :title="$t('Menu Data', 'Data')"><span
              class="material-icons mdi mdi-database"></span><span
              class="nav-link-title">{{$t('Menu Data','Data')}}</span>
            </router-link>
          </li>
          <li class="nav-item" v-if="userCanAccessBackend">
            <a href="/admin/workflow" class="nav-link" data-toggle="tooltip" data-placement="right"
               data-boundary="window" :title="$t('Menu Backend', 'Admin Panel')"><span
              class="material-icons">event_note</span><span
              class="nav-link-title">{{$t('Menu Backend', 'Admin Panel')}}</span>
            </a>
          </li>

          <!--<li class="nav-item dropdown">-->
          <!--<a class="nav-link dropdown-toggle" href="#" id="navbarDropdownWallet" role="button" data-toggle="dropdown"-->
          <!--aria-haspopup="true" aria-expanded="false">-->
          <!--<span class="material-icons mdi mdi-wallet"></span>-->
          <!--<span class="nav-link-title">Wallet</span>-->
          <!--</a>-->
          <!--</li>-->
        </ul>
      </td>
    </tr>
  </table>
</responsive-sidebar>
</template>

<script>
import 'bootstrap'
import ResponsiveSidebar from './ResponsiveSidebar'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  components: { ResponsiveSidebar },
  name: 'sidebar',
  props: ['user', 'toggled'],
  data () {
    return {
      documents: []
    }
  },
  watch: {
    toggled (newValue) {
      this.toggleTooltips(newValue)
    }
  },
  computed: {
    me () {
      return this.app.me
    },
    userCanAccessBackend () {
      return this.app.userIsCreatorOrHigher()
    },
    signatureRequestCount () {
      return this.$store.getters.signatureRequestCount
    }
  },
  created () {
    this.getSigningRequests()
  },
  methods: {
    logout () {
      axios.post('/api/logout', null).then(response => {
        window.location.replace('/')
      }, (err) => {
        this.app.handleError(err)
      })
    },
    async getSigners (docHash) {
      const signers = await this.app.wallet.proxeusFS.getFileSigners(docHash)
      return signers
    },
    getSigningRequests () {
      axios.get('/api/user/document/signingRequests').then(async response => {
        if (response.data) {
          let sigCount = 0
          this.documents = response.data
          for (let i = 0, len = this.documents.length; i < len; i++) {
            const response2 = await this.getSigners(this.documents[i].hash)
            console.log(response2)
            if (!response2.includes(this.me.etherPK)) {
              if (!this.documents[i].rejected) {
                sigCount++
              }
            }
          }
          console.log(sigCount)
          this.$store.dispatch('UPDATE_SIGNERS_COUNT', { sigCount: sigCount })
        }
      }, (err) => {
        console.log(err)
      })
    }
  }
}
</script>
