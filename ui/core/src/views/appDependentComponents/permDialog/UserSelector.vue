<template>
<div ref="main">
  <vue-tags-input v-model="vueTagsInputTarget" :tags="toBeGranted" :disabled="disabled()"
                  @before-adding-tag="beforeAddingTag"
                  :autocomplete-items="addresses" :max-tags="maxItems" :add-only-from-autocomplete="true"
                  :add-on-blur="false" :add-from-paste="false" @tags-changed="updateVueTagsInput"
                  :placeholder="$t('Enter names, email or blockchain addresses...')">
    <div slot="autocomplete-item"
         slot-scope="props"
         class="my-item"
         @click="props.performAdd(props.item)">
      <user-item :item="props.item"/>
    </div>
    <div slot="tag-left" slot-scope="props" class="my-tag-left"
         @click="props.performOpenEdit(props.index)">
      <div></div>
      <user-item :showContent="false" :item="props.tag"/>
    </div>
    <div slot="tag-center" slot-scope="props" class="my-tag-center"
         @click="props.performOpenEdit(props.index)">
      <user-item style="min-width: 278px;" :showPhoto="false" :item="props.tag"/>
    </div>
  </vue-tags-input>
  <button @click="addGranted" style="margin-top:4px;width:100%;"
          v-show="canAdd()" type="button"
          class="btn btn-primary">{{addBtnText}}
  </button>
</div>
</template>

<script>
import VueTagsInput from '@johmun/vue-tags-input'
import UserItem from './UserItem'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'user-selector',
  components: {
    UserItem,
    'vue-tags-input': VueTagsInput
  },
  props: {
    excludes: {
      type: Function,
      default: function () { return null }
    },
    disabled: {
      type: Function,
      default: function () { return false }
    },
    maxItems: {
      type: Number,
      default: 8
    },
    value: {
      type: Array,
      default: null
    },
    addBtnText: {
      type: String,
      default: 'Add'
    },
    dependencyFulfilled: {
      type: Function,
      default: function () { return true }
    },
    uri: { type: String, default: '' }
  },
  data () {
    return {
      toBeGranted: [],
      addresses: [],
      vueTagsInputTarget: ''
    }
  },
  watch: {
    'vueTagsInputTarget': 'vueTagsInputSearch'
  },
  mounted () {
    console.log(this.$refs.main)
    if (this.$refs.main) {
      let inputs = this.$refs.main.getElementsByClassName('ti-new-tag-input')
      if (inputs && inputs.length === 1) {
        inputs[0].addEventListener('keydown', this.keyboardHandler)
      }
    }
  },
  beforeDestroy () {
    if (this.$refs.main) {
      let inputs = this.$refs.main.getElementsByClassName('ti-new-tag-input')
      if (inputs && inputs.length === 1) {
        inputs[0].removeEventListener('keydown', this.keyboardHandler)
      }
    }
  },
  methods: {
    canAdd () {
      return this.toBeGranted && this.toBeGranted.length && this.dependencyFulfilled()
    },
    beforeAddingTag (tag, clb) {
      // fixes vue-tag-input bug
      // bug: when entering chars in the input field and then typing  backspace, the autocompletion dropdown disappears, typing enter at that time adds a new tag.
      if (this.autocompleteDropDownExists()) {
        tag.addTag()
      }
    },
    keyboardHandler (e) {
      if (e.which === 13 || e.keyCode === 13) {
        if (!this.autocompleteDropDownExists()) {
          e.preventDefault()
          e.stopPropagation()
          if (this.canAdd()) {
            this.addGranted()
          }
          return true
        }
      }
    },
    autocompleteDropDownExists () {
      if (this.$refs.main) {
        let a = this.$refs.main.getElementsByClassName('ti-autocomplete')
        return a && a.length === 1
      }
      return false
    },
    updateVueTagsInput (newTags) {
      console.log('updateVueTagsInput')
      this.addresses = []
      this.toBeGranted = newTags
      this.vueTagsInputTarget = ''
    },
    vueTagsInputSearch () {
      this.getAddresses(this.vueTagsInputTarget)
    },
    getAddresses (query) {
      axios.post(this.uri + '?c=' + query, { limit: 6, exclude: this.excludes() }).then(res => {
        if (res.data) {
          for (let i = 0; i < res.data.length; i++) {
            // create text prop for tagInput component, we don't need it but it causes an error
            // in the vue-tags-input if we do not provide it
            if (!res.data[i].name) {
              res.data[i].name = '-'
            }
            res.data[i].text = res.data[i].name
          }
        }
        this.addresses = res.data
        // fixes vue-tag-input bug that does not display any items after adding one
        let inputs = this.$refs.main.getElementsByClassName('ti-new-tag-input')
        if (inputs && inputs.length === 1) {
          inputs[0].blur()
          inputs[0].focus()
        }
      }, (err) => {
        this.addresses = []
        this.app.handleError(err)
      })
    },
    addGranted () {
      let newGranted = []
      for (let i = 0; i < this.toBeGranted.length; i++) {
        if (this.toBeGranted[i] && this.toBeGranted[i].id && !this.exists(this.toBeGranted[i])) {
          newGranted.push(this.toBeGranted[i])
        }
      }
      if (newGranted.length) {
        this.$emit('added', newGranted)
        if (this.value && this.value.length) {
          this.$emit('input', newGranted.concat(this.value))
        } else {
          this.$emit('input', newGranted)
        }
      }
      this.toBeGranted = []
    },
    exists (item) {
      if (this.value && this.value.length) {
        for (let i = 0; i < this.value.length; i++) {
          if (item.id === this.value[i].id) {
            return true
          }
        }
      }
      return false
    }
  }
}
</script>

<style>

  .vue-tags-input {
    background-color: transparent;
  }

  .vue-tags-input .ti-new-tag-input {
    background: transparent;
    /*color: #b7c4c9;*/
  }

  .vue-tags-input .ti-input {
    padding: 4px;
    transition: border-bottom 200ms ease;
    min-height: 40px;
    border: 2px solid #dee2e6;
  }

  .vue-tags-input .ti-new-tag-input-wrapper {
    font-size: 16px;
    padding: 1px 5px;
  }

  .vue-tags-input .ti-new-tag-input-wrapper input {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-left: 2px;
    width: 100%;
  }

  /* some stylings for the autocomplete layer */
  .vue-tags-input .ti-autocomplete {
    background: white;
    border: 1px solid #dee2e6;
    border-top: none;
    padding: 3px;
  }

  .vue-tags-input .ti-autocomplete li {
    border-bottom: 1px solid #f3f3f3;
  }

  .vue-tags-input .ti-item > div {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    padding-right: 5px;
  }

  .vue-tags-input .ti-item.ti-selected-item {
    background: rgba(0, 0, 0, 0.09);
    color: #062a85;
  }

  .vue-tags-input ::-webkit-input-placeholder {
    color: #a4b1b6;
  }

  .vue-tags-input ::-moz-placeholder {
    color: #a4b1b6;
  }

  .vue-tags-input :-ms-input-placeholder {
    color: #a4b1b6;
  }

  .vue-tags-input :-moz-placeholder {
    color: #a4b1b6;
  }

  /* default styles for all the tags */
  .vue-tags-input .ti-tag.ti-valid {
    position: relative;
    background-color: #f1f1f1;
    color: #062a85;
    padding: 1px 5px;
  }

  .vue-tags-input .ti-tag.ti-deletion-mark {
    position: relative;
    background-color: #f1f1f1 !important;
    color: #062a85;
    padding: 1px 5px;
  }

  /* if a tag or the user input is a duplicate, it should be crossed out */
  .vue-tags-input .ti-duplicate span,
  .vue-tags-input .ti-new-tag-input.ti-duplicate {
    text-decoration: line-through;
  }

  /* if the user presses backspace, the complete tag should be crossed out, to mark it for deletion */
  .vue-tags-input .ti-tag:after {
    transition: transform .2s;
    position: absolute;
    content: '';
    height: 2px;
    width: 108%;
    left: -4%;
    top: calc(50% - 1px);
    background-color: #cd0e1b;
    transform: scaleX(0);
  }

  .vue-tags-input .ti-deletion-mark:after {
    transform: scaleX(1);
  }

  .vue-tags-input .ti-tag-center > span {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 330px;
    display: block;
    padding: 4px 2px;
    font-size: 16px;
  }
</style>
