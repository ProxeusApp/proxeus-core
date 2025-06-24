<template>
<div class="search-box" :class="[cls ? cls: '']">
  <input type="text" class="form-control" :placeholder="placeholder || 'Search for...'" ref="searchbox" v-bind="$attrs"
         v-model="searchTerm">
</div>
</template>
<script>
export default {
  name: 'search-box',
  props: ['cls', 'placeholder'],
  watch: {
    searchTerm: function (term) {
      this.$emit('search', term)
    }
  },
  mounted () {
    this.$refs.searchbox && this.$refs.searchbox.focus()
  },
  data () {
    return {
      searchTerm: ''
    }
  },
  methods: {
    focus () {
      this.$refs.searchbox.focus()
    }
  }
}
</script>
<style lang="scss" scoped>
  @use "../assets/styles/variables";

  .search-box .form-control {
    height: 60px;
  }

  .search-box .form-control:focus {
    border-color: $gray-300 !important;
  }

</style>
