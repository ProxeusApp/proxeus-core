<template>
<responsive-sidebar class="sidebar-app mshadow-right-light" v-if="app.userIsCreatorOrHigher()">
  <table style="height: 100%;width: 100%;">
    <tr>
      <td valign="top" style="vertical-align:top">
        <router-link :to="{name:'Workflows'}" class="navbar-brand p-0 m-0">
          <img :src="$t('Sidebar Logo','/static/proxeus_blue.jpg')" alt="" class="d-inline-block align-top">
        </router-link>
        <nav class="collapse show sidebar-sticky">
          <nav class="nav main-nav flex-column main-nav-backend">
            <li class="nav-item" v-if="app.userIsCreatorOrHigher()">
              <router-link :to="{name:'Workflows'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu Workflows','Workflows')"><span
                class="material-icons mdi mdi-source-branch"></span><span
                class="nav-link-title">{{$t('Menu Workflows','Workflows')}}</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="app.userIsCreatorOrHigher()">
              <router-link :to="{name:'Templates'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu Templates','Templates')"><span
                class="material-icons mdi mdi-file-xml"></span><span
                class="nav-link-title">{{$t('Menu Templates','Templates')}}</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="app.userIsCreatorOrHigher()">
              <router-link :to="{name:'Forms'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu Forms','Forms')"><span
                class="material-icons">view_quilt</span><span class="nav-link-title">{{$t('Menu Forms','Forms')}}</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="superadmin">
              <router-link :to="{name:'Users'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu Users','Users')">
                <span class="material-icons">people</span> <span
                class="nav-link-title">{{ $t('backend.menu.user', 'Users') }}</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="superadmin">
              <router-link :to="{name:'I18n'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu Internationalization','Internationalization')">
                <span class="material-icons">language</span><span
                class="nav-link-title">{{$t('Menu Internationalization','Internationalization')}}</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="superadmin">
              <router-link :to="{name:'AdminImportExport'}" class="nav-link" data-toggle="tooltip"
                           data-placement="right"
                           data-boundary="window" :title="$t('Menu Data','Data')">
                <span class="material-icons mdi mdi-database"></span><span
                class="nav-link-title">{{$t('Menu Data','Data')}}</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="app.userIsRoot()">
              <router-link :to="{name:'Settings'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                           data-boundary="window" :title="$t('Menu System Settings','System Settings')">
                <span class="material-icons">settings_applications</span><span
                class="nav-link-title">{{$t('Menu System Settings','System Settings')}}</span>
              </router-link>
            </li>
          </nav>
        </nav>
      </td>
    </tr>
    <tr>
      <td valign="bottom" style="vertical-align:bottom">
        <ul class="nav secondary-nav flex-column">
          <li class="nav-item">
            <router-link :to="{name:'Support'}" class="nav-link" data-toggle="tooltip" data-placement="right"
                         data-boundary="window" :title="$t('Menu Help / Support','Help / Support')"><span
              class="mdi mdi-help-circle-outline material-icons"></span><span
              class="nav-link-title">{{$t('Menu Help / Support','Help / Support')}}</span>
            </router-link>
          </li>
          <li class="nav-item" v-if="app.userIsUserOrHigher()">
            <a href="/user/document" class="nav-link" data-toggle="tooltip" data-placement="right"
               data-boundary="window" :title="$t('Menu Frontend','User View')"><span
              class="material-icons">event_note</span><span
              class="nav-link-title">{{$t('Menu Frontend','User View')}}</span>
            </a>
          </li>
        </ul>
      </td>
    </tr>
  </table>
</responsive-sidebar>
<!--<nav class="col sidebar sidebar-app px-0 pt-0" :class="{toggled:toggled}">-->

<!--</nav>-->
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
  computed: {
    superadmin () {
      return this.app.userIsSuperAdmin()
    }
  },
  methods: {
    logout () {
      axios.post('/api/logout', null).then(response => {
        window.location.replace('/')
      }, (err) => {
        this.app.handleError(err)
      })
    }
  }
}
</script>

<style lang="scss">
  @import "../../assets/styles/variables";
  @import "~bootstrap/scss/mixins";
  @import "../../assets/styles/sidebar.scss";

  .brand-name {
    letter-spacing: 2px;
  }
</style>
