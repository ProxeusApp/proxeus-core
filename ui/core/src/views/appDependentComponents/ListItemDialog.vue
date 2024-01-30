<template>
    <b-modal :modal-class="{noBtns:noBtns}" :cancel-disabled="noBtns" :ok-disabled="noBtns" v-if="elements && elements.length" v-model="modalShow" class="b-modal"
             :title="title || $t('Confirm')"
             :ok-title="$t('Yes')"
             :cancel-title="$t('No')"
             :header-bg-variant="headerBgVariant"
             @hide="onDialogHide"
             @ok="onDialogOk">
        <slot/>
        <table class="nicetbl tblspacing">
            <tbody>
            <!--        <list-item :defaultName="$t('Unnamed')" :timestamps="true" :iconFa="iconFa" :icon="icon" :element="item"/>-->

            <list-item :defaultName="$t('Unnamed')" :to="element.getLink?element.getLink():()=>{}" :index="index"
                       v-for="(element, index) in elements"
                       :key="element.id" :timestamps="timestamp" :error="element.error" :iconFa="iconFa || element.iconFa"
                       :icon="icon || element.icon" :element="element" :_blank="_blank">
            </list-item>
            </tbody>
        </table>
    </b-modal>
</template>
<script>
import { BModal, VBModal } from 'bootstrap-vue'
// import bModalDirective from 'bootstrap-vue/es/directives/modal/modal'
// import ListGroup from '@/components/ListGroup'
import ListItem from '@/components/ListItem.vue'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'list-item-dialog',
  components: {
    ListItem,
    // ListGroup,
    BModal
  },
  directives: {
    'b-modal': VBModal
  },
  props: {
    setup: { type: Function },
    sureFunc: {
      type: Function,
      default: () => {
      }
    },
    iconFa: {
      type: String,
      default: ''
    },
    icon: {
      type: String,
      default: ''
    },
    title: {
      type: String,
      default: ''
    },
    noBtns: {
      type: Boolean,
      default: false
    },
    _blank: {
      type: Boolean,
      default: false
    },
    timestamp: {
      type: Boolean,
      default: true
    }
  },
  created () {
    if (this.setup) {
      this.setup(this.showNow)
    }
  },
  data () {
    return {
      modalShow: false,
      headerBgVariant: 'light',
      elements: null
    }
  },
  methods: {
    onDialogHide () {
      this.modalShow = false
    },
    onDialogOk () {
      this.modalShow = false
      if (this.elements) {
        if (this.elements.length === 1) {
          this.sureFunc(this.elements[0])
        } else {
          this.sureFunc(this.elements)
        }
        this.elements = null
      }
    },
    showNow (e, item) {
      if (e) {
        e.stopPropagation()
      }
      if (item && item.id) {
        // single
        this.elements = [item]
      } else {
        this.elements = item
      }
      if (this.elements) {
        this.modalShow = true
      }
    }
  }
}
</script>

<style lang="scss">
    .modal.noBtns .modal-footer {
        display: none;
    }
</style>
