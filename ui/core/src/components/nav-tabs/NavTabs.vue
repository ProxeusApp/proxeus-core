<template>
    <div class="mod-nav-tabs">
        <ul class="nav">
            <li v-for="tab in tabs"
                class="nav-item"
                :key="tab.title">
                <button @click="selectTab(tab)"
                        class="nav-link"
                        :class="{'active': tab.isActive}">{{ tab.title }}
                </button>
            </li>
        </ul>
        <div class="tabs-content">
            <slot></slot>
        </div>
    </div>
</template>

<script>
export default {
  name: 'tabs',
  data () {
    return {
      tabs: []
    }
  },
  created () {
    this.tabs = this.$children
  },
  methods: {
    selectTab (selectedTab) {
      this.tabs.forEach(tab => {
        tab.isActive = (tab.title === selectedTab.title)

        if (tab.isActive && tab.onSelect) {
          tab.onSelect(tab)
        }
      })
    }
  }
}
</script>
