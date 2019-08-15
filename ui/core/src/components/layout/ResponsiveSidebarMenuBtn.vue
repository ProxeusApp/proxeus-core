<template>
<i @click="$root.$emit('menutoggler-click');" class="material-icons sidebar-menutoggler">
  menu
</i>
</template>
<script>
export default {
  name: 'responsive-sidebar-menu-btn'
}
</script>
