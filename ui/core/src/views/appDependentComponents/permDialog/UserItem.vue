<template>
<table class="nicetbl">
  <tr>
    <td v-show="showPhoto" class="tdmin">
      <img v-if="item.photo" class="mlist-group-icon" :src="item.photo"/>
      <i v-else class="mlist-group-icon material-icons">person</i>
    </td>
    <td v-show="showContent" class="tdmax impcnt">
      <span class="font-weight-bold">{{ item.name }}</span>
      <small class="light-text" v-if="item.email">{{ item.email }}</small>
      <small class="light-text" v-if="item.etherPK">{{ item.etherPK }}</small>
      <small class="light-text" v-else-if="item.detail">{{ item.detail }}
      </small>
      <small class="light-text" v-else>{{$t('-')}}</small>
    </td>
  </tr>
</table>
</template>

<script>
export default {
  name: 'user-item',
  props: {
    item: null,
    showPhoto: { type: Boolean, default: true },
    showContent: { type: Boolean, default: true }
  }
}
</script>
