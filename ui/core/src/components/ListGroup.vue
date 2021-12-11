<template>
  <div>
    <div v-if="!listOnly" class="d-flex flex-row w-100 align-items-center mb-3">
      <slot name="addBtn">
        <button :title="$t('create a new one')" type="button" @click="toggleNewItemFormVisible" class="btn btn-primary">
          Create new
        </button>
        <button v-if="elements && elements.length>0" @click="exportData" type="button" class="btn btn-link ml-auto">Export</button>
      </slot>
      <button type="button" v-show="itemsSelected && false" @click="remove" class="btn btn-secondary ml-3">Remove</button>
    </div>
    <div class="p-0 ml-0" style="position:relative;">
      <search-box v-on:search="search" ref="searchBox"/>
    </div>
    <div ref="listGroup" class="mbottominset">
      <slot name="newItemForm">
        <div v-show="newItemFormVisible" v-if="!listOnly"
             class="p-0 new-item-form list-group-item-action flex-column align-items-start bg-light"
             style="position: relative;">
          <i @click="newItemFormVisible=false" style="position: absolute;right: 10px;top:10px;cursor:pointer;"
             class="material-icons">
            clear
          </i>
          <h2 class="pt-2 pl-2 text-center">Create new {{ path }}</h2>
          <div class="d-flex flex-row w-100 justify-content-between align-items-center">
            <div class="py-2 px-2 w-auto">
              <label for="newNameInput">Name</label>
            </div>
            <div class="py-2 px-2 w-50">
              <input type="text" ref="newItemName" v-model.trim="newElement.name" name="newNameInput"
                     id="newNameInput" class="form-control" required>
            </div>
            <div class="py-2 px-2 w-auto">
              <label for="newDetailInput">Detail</label>
            </div>
            <div class="py-2 px-2 w-50">
              <input type="text" v-model.trim="newElement.detail" name="newDetailInput" id="newDetailInput"
                     class="form-control">
            </div>
            <div v-if="displayPrice" class="py-2 px-2 w-auto">
              <label for="newPriceInput">Price&nbsp;(XES)</label>
            </div>
            <div v-if="displayPrice" class="py-2 px-2 w-30">
              <input type="number" v-model.number="newElement.price" name="newPriceInput" id="newPriceInput"
                     class="form-control">
            </div>
            <div class="py-2 px-2 w-auto">
              <button type="button" @click="create" class="btn btn-primary btn-round"
                      :disabled="newElement.name === '' || displayPrice === true && newElement.price === ''">
                <span class="material-icons">check</span>
              </button>
            </div>
          </div>
        </div>
      </slot>
      <slot name="list">
        <div v-if="!elements || elements.length===0" class="no-elements mt-3">
          <div v-if="loading" style="height: 60px;width: 100%;"></div>
          <span v-else class="light-text">{{ $t('No elements found') }}</span>
        </div>
        <table v-else class="table mt-3">
          <tbody>
          <list-item :defaultName="defaultName" :to="getLink(element)" :index="index" v-for="(element, index) in elements"
                     :key="element.id" :timestamps="timestamps" :iconFa="iconFa" :icon="icon" :element="element" :price="element.price">
            <slot v-bind="element"/>
          </list-item>
          </tbody>
        </table>
      </slot>
      <trigger :init="triggerInit" @trigger="bottomTrigger"/>
    </div>
  </div>
</template>

<script>
import SearchBox from './SearchBox'
import Trigger from './Trigger'
import ListItem from './ListItem'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'list-group',
  components: {
    ListItem,
    Trigger,
    SearchBox
  },
  props: {
    deleteElementFunc: { type: Function, default: null },
    defaultName: { type: String, default: '' },
    linkPrefix: {},
    nodeType: String,
    nodeTypeSuffix: String,
    path: String,
    icon: String,
    iconFa: String,
    listOnly: Boolean,
    payment: Boolean,
    select: Boolean,
    fixedPath: {
      required: false,
      type: String
    },
    linkResolver: {
      type: Function,
      default: null
    },
    prependFunc: {
      type: Function,
      default: null
    },
    timestamps: {
      type: Boolean,
      default: true
    },
    enableCheckboxes: {
      type: Boolean,
      default: true
    },
    displayPrice: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      elements: [],
      searchTerm: '',
      selectedAll: false,
      listElements: null,
      loading: false,
      newItemFormVisible: false,
      listIndex: 0,
      reachedTheEnd: false,
      triggerApi: {},
      lastImps: null,
      lastExs: null,
      delImpsUrl: '',
      delExsUrl: '',
      newElement: {
        name: '',
        detail: '',
        price: 0
      }
    }
  },
  watch: {
    selectedAll (newValue) {
      const self = this
      this.elements && this.elements.map((value, key) => {
        // Use Vue $set so Vue can react to the new object keys
        self.$set(self.elements[key], 'selected', newValue)
      })
    }
  },
  mounted () {
    this.setListGroupHeight()
    this.$refs.searchBox.focus()
  },
  created () {
    if (this.prependFunc) {
      this.prependFunc(this.prependItem)
    }
    if (this.deleteElementFunc) {
      this.deleteElementFunc(this.deleteElement)
    }
    // for better support on mobile devices because 100vh doesn't work correctly
    window.addEventListener('resize', this.setListGroupHeight)
    this.loading = true
    this.$root.$on('lastExportResults', this.lastExportResults)
    this.$root.$on('lastImportResults', this.lastImportResults)
    if (this.app.lastExportResults) {
      this.lastExportResults(this.app.lastExportResults)
    } else {
      this.app.loadLastExportResults()
    }
    if (this.app.lastImportResults) {
      this.lastImportResults(this.app.lastImportResults)
    } else {
      this.app.loadLastImportResults()
    }
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.setListGroupHeight)
    this.$root.$off('lastExportResults', this.lastExportResults)
    this.$root.$off('lastImportResults', this.lastImportResults)
  },
  computed: {
    itemsSelected () {
      let itemsSelected = false
      this.elements && this.elements.map((value, key) => {
        if (value.selected) {
          itemsSelected = true
        }
      })
      return itemsSelected
    }
  },
  methods: {
    delImpsAction () {
      this.app.loadLastImportResults(null, this.delImpsUrl)
    },
    delExsAction () {
      this.app.loadLastExportResults(null, this.delExsUrl)
    },
    lastExportResults (res) {
      this.lastExs = null
      this.delExsUrl = ''
      if (res && res.results) {
        for (const key in res.results) {
          const resHasOwnProp = Object.prototype.hasOwnProperty.call(res.results, key)
          if (resHasOwnProp && key.toLowerCase() === this.nodeType.toLowerCase()) {
            if (res.results[key]) {
              this.lastExs = res.results[key]
              this.delExsUrl = 'delete=' + key
            }
          }
        }
      }
    },
    lastImportResults (res) {
      this.lastImps = null
      this.delImpsUrl = ''
      if (res && res.results) {
        for (const key in res.results) {
          const resHasOwnProp = Object.prototype.hasOwnProperty.call(res.results, key)
          if (resHasOwnProp && key.toLowerCase() === this.nodeType.toLowerCase()) {
            if (res.results[key]) {
              this.lastImps = res.results[key]
              this.delImpsUrl = 'delete=' + key
            }
          }
        }
      }
    },
    prependItem (item) {
      if (!this.elements) {
        this.elements = []
      }
      if (item) {
        item.isNew = true
        setTimeout(() => { item.isNew = false }, 5000)
        this.elements.unshift(item)
        $(this.$refs.listGroup).animate({ scrollTop: 0 }, 500)
      }
    },
    deleteElement (id) {
      this.elements = this.elements.filter(v => {
        return v.id !== id
      })
    },
    triggerInit (startFunc, stopFunc, hideFunc) {
      this.triggerApi.start = startFunc
      this.triggerApi.stop = stopFunc
      this.triggerApi.hide = hideFunc
      this.triggerApi.start()
    },
    exportData () {
      let params = this.nodeType
      if (this.searchTerm) {
        params += '&contains=' + this.searchTerm
      }
      const url = '/api/' + this.nodeType.toLowerCase() + '/export'
      this.app.exportData(params, null, url, this.nodeType)
    },
    createParams () {
      let params = '?'
      if (this.searchTerm) {
        params += 'c=' + this.searchTerm + '&'
      }
      params += 'i=' + this.listIndex
      if (this.path === 'user') {
        params += '&m=f'
      }
      return params
    },
    getReqLink () {
      let p = this.fixedPath
      if (!p) {
        p = '/api/admin/' + this.path + '/list'
      }
      return p + this.createParams()
    },
    getLink (element) {
      if (this.linkResolver) {
        return this.linkResolver(element)
      }
      const prefix = this.linkPrefix !== undefined ? this.linkPrefix : '/admin'
      const nodeTypeSuffix = this.nodeTypeSuffix ? '/' + this.nodeTypeSuffix : ''
      const link = prefix + '/' + this.nodeType + '/' + element.id + nodeTypeSuffix

      // no payment if i am owner or price is 0
      if (this.payment === false ||
        element.owner === this.app.me.id ||
        element.price <= 0 ||
        element.paid === true) {
        return link
      }
      return link + '/payment'
    },
    toggleNewItemFormVisible () {
      this.newItemFormVisible = !this.newItemFormVisible
      this.$nextTick(() => this.$refs.newItemName.focus())
    },
    create: function () {
      this.newElement.id = ''
      const newElementClone = Object.assign({}, this.newElement)

      if (newElementClone.price && newElementClone.price % 1 !== 0) {
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Decimal numbers are not allowed for the XES Price field.'),
          type: 'error'
        })
        return
      }

      axios.post('/api/admin/' + this.path + '/update', newElementClone).then(response => {
        if (response.data.id) {
          this.newElement = {}
          this.prependItem(response.data)
        }
        this.resetNewElement()
        this.newItemFormVisible = false
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: 'The new item has been created successfully.',
          type: 'success'
        })
      }, (err) => {
        this.$emit('error', err)
        this.app.handleError(err)
        this.$notify({
          group: 'app',
          title: 'Error',
          text: 'The element could not be created. Please try again or if the error persists contact the platform operator.',
          type: 'error'
        })
      })
    },
    remove () {

    },
    resetNewElement () {
      this.newElement = {
        name: '',
        detail: ''
      }
    },
    bottomTrigger (startFunc, stopFunc, hideFunc) {
      if (!this.reachedTheEnd) {
        this.loading = true
        const searchTermBeforeReq = this.searchTerm
        axios.get(this.getReqLink()).then(response => {
          if (searchTermBeforeReq === this.searchTerm && response.data && response.data.length > 0) {
            hideFunc()
            this.elements = [...this.elements, ...response.data]
            this.listIndex += 1
          } else {
            this.reachedTheEnd = true
            stopFunc()
          }
          this.loading = false
        }, (err) => {
          // on err
          stopFunc()
          this.loading = false
          this.app.handleError(err)
        })
      }
    },
    search (term) {
      this.searchTerm = term
      this.listIndex = 0
      this.reachedTheEnd = false
      if (this.triggerApi.stop) {
        this.triggerApi.stop()
      }
      this.elements = []
      this.$emit('searched', term)
      this.loading = true
      axios.get(this.getReqLink()).then(response => {
        if (response.data && response.data.length > 0) {
          if (this.triggerApi.hide) {
            this.triggerApi.hide()
          }
          this.elements = response.data
          this.listIndex += 1
          if (this.triggerApi.start) {
            this.triggerApi.start()
          }
        } else {
          this.reachedTheEnd = true
          if (this.triggerApi.stop) {
            this.triggerApi.stop()
          }
        }
        this.loading = false
      }, error => {
        this.loading = false
        if (this.triggerApi.stop) {
          this.triggerApi.stop()
        }
        this.app.handleError(error)
      })
    },
    setListGroupHeight () {
      if (this.$refs.listGroup) {
        const lg = $(this.$refs.listGroup)
        lg.css('height', ($(document.body).height() - 166) + 'px')
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .new-item-form {
    border-bottom: 1px solid #dadada;
  }

  .flex-col-truncate {
    overflow: hidden;
    min-width: 0;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .justify-content-between {
    overflow-x: hidden;
  }

  .animate-new-entry {
    background: linear-gradient(-45deg, #ffffff, #40e1d1, #ffffff, #40e1d1);
    background-size: 400% 400%;
    -webkit-animation: NewRowEntryGradient 5s ease infinite;
    -moz-animation: NewRowEntryGradient 5s ease infinite;
    animation: NewRowEntryGradient 5s ease infinite;
  }

  @-webkit-keyframes NewRowEntryGradient {
    0% {
      background-position: 0% 50%
    }
    50% {
      background-position: 100% 50%
    }
    100% {
      background-position: 0% 50%
    }
  }

  @-moz-keyframes NewRowEntryGradient {
    0% {
      background-position: 0% 50%
    }
    50% {
      background-position: 100% 50%
    }
    100% {
      background-position: 0% 50%
    }
  }

  @keyframes NewRowEntryGradient {
    0% {
      background-position: 0% 50%
    }
    50% {
      background-position: 100% 50%
    }
    100% {
      background-position: 0% 50%
    }
  }
</style>
