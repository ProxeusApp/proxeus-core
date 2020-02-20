<template>
<textarea @input="onInput" class="form-control" style="background: #ffffff00;" v-model="translation"></textarea>
</template>

<script>
export default {
  name: 'i18n-trans-cell',
  props: ['langKey', 'langCode', 'translations'],
  created () {

  },
  methods: {
    onInput (e) {
      let t = e.target
      if (t.scrollHeight > t.clientHeight) t.style.height = (t.scrollHeight + 20) + 'px'
    }
  },
  computed: {
    translation: {
      get () {
        if (!this.translations) {
          return ''
        }
        if (this.translations[this.langKey]) {
          return this.translations[this.langKey][this.langCode] || ''
        }
        return ''
      },
      set (value) {
        this.$emit('input', { lang: this.langCode, k: this.langKey, value: value })
        if (this.$i18n.locale() === this.langCode) {
          this.$i18n.add(this.langCode, {
            [this.langKey]: value
          })
        }
      }
    }
  }
}
</script>

<style>
  td.tdmin textarea.tt-style {
    width: 300px;
    overflow: hidden;
  }

  td.tdmax textarea.tt-style {
    width: 100%;
    min-height: 100%;
  }

  textarea.tt-style {
    border: 1px solid #edeff1;
    width: 100%;
    border-radius: 0;
    padding: 0px 4px;
    height: 100%;
    min-height: 100%;
    border-left: none;
    border-bottom: none;
  }
</style>
