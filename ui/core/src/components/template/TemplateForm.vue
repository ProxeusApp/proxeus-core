<template>
<div class="card" v-if="compiledForm">
  <a class="card-header" data-toggle="collapse"
     :href="'#f-'+form.id"
     role="button" aria-expanded="false" :aria-controls="'#f-'+form.id">{{ form.name }}
  </a>
  <form class="card-body collapse form-compiled" data-parent="#formsContainer" :id="'f-'+form.id">{{compiledForm()}}</form>
</div>
</template>

<script>
import formCompilerAdapter from '../../libs/formcompiler-adapter.js'

export default {
  mixins: [formCompilerAdapter],
  name: 'template-form',
  props: ['form', 'comps'],
  watch: {
    form: function (newForm) {
      if (newForm.data && newForm.data.formSrc) {
        this.compileForm()
      }
    }
  },
  mounted () {
    const self = this
    $(document).ready(function () {
      if (!self.form) return
      $('#f-' + self.form.id).on('shown.bs.collapse', function () {
        self.$scrollTo('#f-' + self.form.id, 500, {
          container: '.layout-fixed-scroll-view',
          offset: -60
        })
      })
    })
  },
  created () {
    if (this.form && this.form.data && this.form.data.formSrc) {
      this.compileForm()
    }
  },
  data () {
    return {
      compiledForm: null
    }
  },
  methods: {
    compileForm () {
      const self = this
      this.adapterCompile(this.form.data.formSrc, this.comps, (compiledForm) => {
        if (!window.___copyMe) {
          window.___copyMe = function (e) {
            $(e.target).closest('.var-box').copyHtmlToClipboard()
          }
        }
        self.compiledForm = compiledForm
        self.$nextTick(() => {
          const $form = $('#f-' + self.form.id + ' > form')
          $form.on('dynamicFormScriptExecuted', function () {
            if (self.form.data && self.form.data.data) {
              $form.fillForm(self.form.data.data)
            }
            const changeOptions = {
              fileUrl: '/api/admin/form/test/file/' + self.form.id,
              url: '/api/admin/form/test/data/' + self.form.id,
              success: function (data, textStatus, xhr, myRe) {
                if (xhr.status >= 200 && xhr.status <= 299) {
                  self.$emit('updatedFormField', self.form.id)
                }
              },
              error: function (xhr, a, b, myRe) {

              }
            }
            $form.assignSubmitOnChange(changeOptions)
            $form.on('formFieldsAdded', function (event, parent) {
              if (parent && parent.length) {
                parent.assignSubmitOnChange(changeOptions)
              }
            })
            $form.find('.field-parent').each(function () {
              var _t = $(this)
              var _i = _t.find('input,select,textarea').first()
              var varName = _i.attr('name')
              var ___inputPrefix = 'input.'
              _t.append(
                '<div class="var-box"><span class="var-path" style="z-index:-1;position: absolute;background: none;top: -5px;right: 0;"><span class="e" style="">{{</span><span class="v">input</span><span class="d">.</span>' +
                varName +
                '<span class="e" style="">}}</span></span><span class="var-path" draggable="true" onclick="window.___copyMe(event)" ondragstart="libreHub.dragStart(event, libreHub.getVarJSON(\'' +
                ___inputPrefix + varName +
                '\')' +
                ');"><span class="e">{{</span><span class="v">input</span><span class="d">.</span>' +
                varName +
                '<span class="e">}}</span></span></div>'
              )
            })
          })
        })
      })
    }
  }
}
</script>

<style lang="scss">

  .var-box {
    /*display: none;*/
    position: absolute;
    right: 21px;
    top: 12px;
    > span {
      border: 1px solid;
    }
    > i {
      /*position: absolute;*/
      right: -6px;
      font-size: 13px;
      color: #00829d;
      background: rgba(255, 255, 255, 0.72);
      border-radius: 4px;
      padding: 2px;
      cursor: pointer;
    }

    .var-path {
      font-family: monospace, serif;
      color: #00778f;
      background: rgba(255, 255, 255, 0.72);
      padding: 5px;
      border-radius: 4px;
      white-space: nowrap;

      > span.e {
        color: #000000;
      }

      > span.v {
        color: #062a85;
      }

      > span.d {
        color: #7a00ff;
        font-weight: bold;
      }
    }
  }

  .form-group:hover {
    .var-box {
      display: block;
    }
  }

</style>
