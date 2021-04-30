import FT_FormBuilderCompiler from './legacy/formbuilder-compiler'
import Vue from 'vue'

export default {
  data () {
    return {
      ft: null,
      components: null
    }
  },
  created () {
    try {
      this.ft = new FT_FormBuilderCompiler({
        i18n: {
          onTranslate: function (keyArray, callback, notAsync) {
            if (keyArray instanceof Array) {
              const translations = []
              keyArray.forEach(key => {
                translations.push(Vue.i18n.translate(key))
              })
              callback(translations)
            }
          }
        },
        requestComp: this.requestComp
      })
    } catch (e) {
      console.log(e)
    }
  },
  methods: {
    requestComp (id, cb) {
      return this.components[id]
    },
    adapterCompile (form, components, callback) {
      try {
        console.log('adapterCompile')
        this.components = components
        this.ft.compileForm({
          done (frm) {
            callback(frm)
          },
          form: form
        })
      } catch (e) {
        console.log(e)
      }
    }
  }
}
