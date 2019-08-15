<template>
<div class="topnav mshadow-light">
  <table style="width: 100%;height: 58px;min-height: 58px;">
    <tbody style="width: 100%;">
    <tr style="width: 100%;">
      <td class="tdmin" v-if="isShareLinkAndUserOrHigher()">
        <router-link class="topnav-back btn btn-light mr-3" :to="returnToRoute" v-if="returnToRoute">
          <span class="material-icons pt-10 mdi md-36 mdi-chevron-left"></span>
        </router-link>
      </td>
      <td class="tdmin" v-if="sidebarToggler && isShareLinkAndUserOrHigher()">
        <responsive-sidebar-menu-btn/>
      </td>
      <slot name="td-start"/>
      <td class="tdmax impcnt">
                    <span class="title" style="width: 100%;padding: 4px 0;">
                    {{ title }}
                    </span>
      </td>
      <slot name="td"/>
      <td class="tdmin">
        <slot name="buttons"/>
      </td>
      <td class="tdmin">
        <top-right-profile/>
      </td>
    </tr>
    </tbody>
  </table>
</div>
</template>

<script>
import TopRightProfile from '../user/TopRightProfile'
import ResponsiveSidebarMenuBtn from './ResponsiveSidebarMenuBtn'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  components: {
    ResponsiveSidebarMenuBtn,
    TopRightProfile
  },
  name: 'top-nav',
  props: {
    title: {
      type: String,
      default: ''
    },
    returnToRoute: {
      type: Object,
      default: null
    },
    sidebarToggler: {
      type: Boolean,
      default: true
    }
  },
  methods: {
    isShareLinkAndUserOrHigher () {
      return /^\/p\//.test(location.pathname) ? this.app.userIsUserOrHigher() : true
    }
  }
}
</script>

<style scoped>

  .topnav-buttons {
    white-space: nowrap;
    text-overflow: ellipsis;
  }
</style>
