<template>
  <nav class="navbar navbar-expand-lg py-0 topnav d-flex flex-row flex-nowrap"
       v-bind="$attrs"
       :style="{background:bg ? bg : ''}">
    <router-link class="topnav-back btn btn-sm btn-light mr-3" :to="returnToRoute"
                 v-if="returnToRoute && isShareLinkAndUserOrHigher">
      <span class="material-icons mdi md-36 mdi-chevron-left"></span>
    </router-link>

    <h1 class="navbar-text">
      {{ title }}
    </h1>
    <div class="topnav-buttons ml-auto">
      <div class="ml-auto d-flex align-items-center">
        <slot name="buttons"></slot>
        <top-right-profile/>
      </div>
    </div>
  </nav>
</template>

<script>
import TopRightProfile from '../user/TopRightProfile'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  components: {
    TopRightProfile
  },
  name: 'top-nav',
  props: [
    'title', 'sm', 'bg', 'returnToRoute'
  ],
  computed: {
    isShareLinkAndUserOrHigher () {
      return /^\/p\//.test(window.location.pathname) ? this.app.userIsUserOrHigher() : true
    }
  }
}
</script>
