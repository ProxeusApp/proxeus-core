<template>
<responsive-sidebar class="sidebar-app" :toggled="toggled" v-if="app.userIsCreatorOrHigher()">
  <table style="height: 100%;width: 100%;">
    <tr>
      <td valign="top" style="vertical-align:top">
        <router-link :to="{name:'Workflows'}" class="navbar-brand mb-0 w-100 px-3 mx-0">
          <svg width="140" viewBox="0 0 187 55" xmlns="http://www.w3.org/2000/svg"><title color="#fff">logo</title>
            <g fill="none" fill-rule="evenodd">
              <path d="M70 18.187h4.98c3.453 0 5.825 2.155 5.825 5.292 0 3.091-2.372 5.292-5.825 5.292h-2.255v6.042H70V18.187zm4.862 8.055c1.997 0 3.171-1.1 3.171-2.763 0-1.662-1.174-2.763-3.17-2.763h-2.138v5.526h2.137zm15.503 1.499h-1.574v7.072h-2.724V18.187h5.402c3.1 0 5.309 1.99 5.309 4.777 0 2.154-1.316 3.864-3.43 4.52l5.309 7.329h-3.265l-5.027-7.072zm.846-2.365c1.879 0 2.795-.913 2.795-2.318 0-1.429-.916-2.365-2.795-2.365h-2.42v4.683h2.42zm10.382 1.124c0-4.8 3.805-8.5 8.527-8.5 4.744 0 8.526 3.7 8.526 8.5s-3.782 8.5-8.526 8.5c-4.722 0-8.527-3.7-8.527-8.5zm14.282 0c0-3.325-2.467-5.924-5.755-5.924-3.289 0-5.755 2.599-5.755 5.924s2.466 5.924 5.755 5.924c3.288 0 5.755-2.599 5.755-5.924zm11.72-.117l-5.73 8.43h3.17l4.416-6.463 4.416 6.463h3.124l-5.708-8.43 5.38-8.196h-3.148l-4.087 6.159-4.087-6.159h-3.148l5.403 8.196zm14.517-8.196v16.626h9.889v-2.67h-7.164v-4.449h5.755v-2.622h-5.755v-4.215H152v-2.67h-9.889zm15.503 10.163c0 4.098 2.795 6.65 6.718 6.65s6.694-2.552 6.694-6.65V18.187h-2.724V28.42c0 2.576-1.598 4.004-3.97 4.004-2.373 0-3.993-1.428-3.993-4.004V18.187h-2.725V28.35zm18.462 3.817c1.058 1.709 3.101 2.833 5.614 2.833 3.218 0 5.309-1.92 5.309-4.73 0-2.318-1.55-3.77-3.829-4.753l-1.292-.586c-1.667-.702-2.42-1.475-2.42-2.552 0-1.148 1.081-1.803 2.209-1.803 1.198 0 1.95.445 2.819 1.475l1.879-1.944C185.144 18.656 183.618 18 181.69 18c-2.536 0-4.956 1.733-4.956 4.402 0 2.154 1.245 3.747 3.758 4.847l1.316.562c1.644.726 2.42 1.452 2.42 2.693 0 1.194-.94 1.92-2.561 1.92-1.456 0-2.748-.936-3.523-1.99l-2.068 1.733z"
                    fill="#ebedf3"/>
              <path d="M24.93 0c.394.072.824.109 1.218.253.43.145.86.29 1.253.507 6.09 3.546 12.214 7.056 18.268 10.638 2.042 1.194 3.08 3.076 3.116 5.428.036 7.056.036 14.111 0 21.131 0 2.352-1.074 4.234-3.08 5.428-6.09 3.582-12.179 7.092-18.268 10.674-.931.543-1.934.724-3.009.76 0-.108-.035-.217-.035-.326V0h.537zM0 16.174c.143-.506.25-1.013.43-1.483.25-.76.788-1.592 1.253-2.027v29.563c-.107-.145-.214-.253-.286-.362-.68-.904-1.146-1.918-1.29-3.04 0-.072-.035-.108-.071-.18C0 31.155 0 23.665 0 16.175zm22.029 38.247c-.574-.326-1.11-.615-1.648-.94-.036-.037-.072-.11-.108-.182V1.846c0-.325.072-.506.394-.65.43-.218.86-.471 1.362-.76V54.42zm-4.084-2.351c-.537-.327-1.074-.616-1.576-.906-.071-.036-.107-.18-.143-.289-.036-.072 0-.145 0-.217V4.198c0-.326.072-.507.358-.652.466-.217.896-.507 1.361-.796v49.32zM13.898 5.138v44.543c-.072-.036-.144-.036-.18-.072-.43-.254-.895-.507-1.325-.797-.107-.072-.143-.217-.215-.361-.035-.073 0-.145 0-.217V6.622c0-.362.108-.543.43-.724.394-.217.824-.507 1.29-.76zM9.779 47.329a30.57 30.57 0 0 1-1.577-.94c-.071-.037-.107-.146-.143-.218-.036-.072 0-.145 0-.217V8.937c0-.289.072-.47.323-.578.465-.254.895-.507 1.397-.796v39.766zM5.731 9.878v35.063c-.573-.326-1.11-.652-1.648-.941a.793.793 0 0 1-.107-.18V11.326c0-.29.107-.47.358-.615.466-.254.896-.543 1.397-.833z"
                    fill-rule="nonzero" fill="#ebedf3"/>
            </g>
          </svg>
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
                class="material-icons mdi mdi-code-block-tags"></span><span
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
