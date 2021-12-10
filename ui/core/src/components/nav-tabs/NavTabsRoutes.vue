<template>
    <div class="mod-nav-tabs">
        <ul class="nav">
            <li v-for="tab in tabs"
                class="nav-item"
                :key="tab.title">
                <router-link :to="tab.to"
                        class="nav-link">{{ tab.title }}
                </router-link>
            </li>
        </ul>
    </div>
</template>

<script>
export default {
  name: 'nav-tabs-routes',
  props: {
    tabs: {
      type: Array,
      required: true
    }
  }
}
</script>
