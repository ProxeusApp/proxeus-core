<template>
  <div>
    <div v-if="!listOnly" class="d-flex flex-row w-100 align-items-center">
      <slot name="addBtn">
        <button type="button" @click="toggleNewItemFormVisible" class="btn btn-primary btn-round">
          <i class="material-icons">add</i>
        </button>
      </slot>
      <div class="ml-3">
        <!-- TODO: Uncomment if multi row actions implemented -->
        <!--<div class="custom-control custom-checkbox">-->
        <!--<input type="checkbox" v-model="selectedAll" class="custom-control-input" id="selectAll">-->
        <!--<label class="ml-2 custom-control-label" for="selectAll">Select all</label>-->
        <!--</div>-->
      </div>
      <button type="button" v-show="itemsSelected && false" @click="remove" class="btn btn-danger ml-3">Remove
      </button>
    </div>
    <div class="p-0 w-100 ml-0 mt-3">
      <search-box v-on:search="search" ref="searchBox"/>
    </div>
    <div class="loading" v-show="loading === true">
      <spinner background="white"/>
    </div>
    <div class="no-elements mt-3" v-show="loading === false && !elements && newItemFormVisible === false">
      <span class="text-muted">No elements found</span>
    </div>
    <div class="list-group" v-show="loading === false">
      <slot name="newItemForm">
        <div v-show="newItemFormVisible" v-if="!listOnly"
             class="p-0 new-item-form list-group-item-action flex-column align-items-start bg-light">
          <h2 class="pt-2 pl-2 text-center">Create new {{ path }}</h2>
          <div class="d-flex flex-row w-100 justify-content-between align-items-center">
            <div class="py-2 px-3 w-auto">
              <label for="newNameInput">Name</label>
            </div>
            <div class="py-2 px-3 w-50">
              <input type="text" ref="newItemName" v-model.trim="newElement.name" name="newNameInput"
                     id="newNameInput" class="form-control" required>
            </div>
            <div class="py-2 px-3 w-auto">
              <label for="newDetailInput">Detail</label>
            </div>
            <div class="py-2 px-3 w-50">
              <input type="text" v-model.trim="newElement.detail" name="newDetailInput" id="newDetailInput"
                     class="form-control">
            </div>
            <div class="py-2 px-3 w-auto">
              <button type="button" @click="create" class="btn btn-primary btn-round"
                      :disabled="newElement.name === ''">
                <span class="material-icons">check</span>
              </button>
            </div>
          </div>
        </div>
      </slot>
      <slot name="list">
        <div v-for="(element, index) in elements" :key="element.id"
             class="p-0 flex-column list-group-item" :data-index="index">
          <div class="d-flex flex-row w-100 justify-content-between">
            <!-- TODO: Uncomment if multi row actions implemented -->
            <!--<div class="py-2 px-3 w-auto align-self-center" v-if="enableCheckboxes">-->
            <!--<label class="custom-control custom-checkbox">-->
            <!--<input type="checkbox" v-model="element.selected" name="multiSelect"-->
            <!--class="custom-control-input">-->
            <!--<span class="custom-control-label" for="multiSelect"></span>-->
            <!--</label>-->
            <!--</div>-->
            <router-link class="p-0 list-group-item-action flex-column align-items-start"
                         :to="getLink(element)">
              <div class="d-flex flex-row w-100 justify-content-between">
                <div class="py-2 px-3 w-auto align-self-center">
                  <i class="list-group-icon material-icons" :class="[iconFa ? iconFa : '']">{{ icon || ''}}</i>
                </div>
                <div class="py-2 px-3 d-flex flex-column flex-col-truncate align-self-center"
                     :class="{'w-50':timestamps,'w-100':timestamps===false}">
                  <h5 class="mb-0 pb-0 d-inline-block w-100 flex-col-truncate font-weight-bold">{{ element.name }}</h5>
                  <small class="text-muted mt-1 d-inline-block w-100 flex-col-truncate" v-if="element.detail">
                    {{ element.detail }}
                  </small>
                </div>
                <div class="py-2 px-3 w-25 d-flex flex-column lg-small align-self-center" v-if="timestamps">
                  <h5 class="mb-0 pb-0">{{ element.created | moment('DD.MM.YY - HH:mm') }}</h5>
                  <small class="text-muted mt-1 d-inline-block">Created</small>
                </div>
                <div class="py-2 px-3 w-25 lg-small align-self-center" v-if="timestamps">
                  <h5 class="mb-0 pb-0">{{ element.updated | moment('DD.MM.YY - HH:mm') }}</h5>
                  <small class="text-muted mt-1 d-inline-block">Updated</small>
                </div>
              </div>
            </router-link>
          </div>
        </div>
      </slot>
    </div>
  </div>
</template>

<script>
import SearchBox from './SearchBox'
import Spinner from '@/components/Spinner'
import moment from 'moment'

export default {
  name: 'list-group',
  components: {
    SearchBox,
    Spinner,
    moment
  },
  props: {
    linkPrefix: {},
    nodeType: String,
    nodeTypeSuffix: String,
    path: String,
    icon: String,
    iconFa: String,
    listOnly: Boolean,
    select: Boolean,
    fixedPath: {
      required: false,
      type: String
    },
    timestamps: {
      type: Boolean,
      default: true
    },
    listQuery: {
      type: String,
      default: '/list?c='
    },
    enableCheckboxes: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      deleteElementFunc: { type: Function, default: null },

      elements: [],
      searchTerm: '',
      selectedAll: false,
      listElements: null,
      loading: false,
      lastRequest: null,
      loadingBuffer: false,
      newItemFormVisible: false,
      newElement: {
        name: '',
        detail: ''
      }
    }
  },
  watch: {
    selectedAll (newValue) {
      let self = this
      this.elements && this.elements.map((value, key) => {
        // Use Vue $set so Vue can react to the new object keys
        self.$set(self.elements[key], 'selected', newValue)
      })
    }
  },
  mounted () {
    this.$refs.searchBox.focus()
  },
  created () {
    this.loading = true
    const p = this.fixedPath || '/api/admin/' + this.path + '/list'
    axios.get(p).then(response => {
      this.elements = response.data
      this.loading = false
    }).catch(e => {
      this.loading = false
      this.$notify({
        group: 'app',
        title: 'Error',
        text: 'Could not load elements.',
        type: 'error'
      })
    })
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
    getLink (element) {
      const prefix = this.linkPrefix !== undefined ? this.linkPrefix : '/admin'
      const nodeTypeSuffix = this.nodeTypeSuffix ? '/' + this.nodeTypeSuffix : ''
      return prefix + '/' + this.nodeType + '/' + element.id + nodeTypeSuffix
    },
    toggleNewItemFormVisible () {
      this.newItemFormVisible = !this.newItemFormVisible
      this.$nextTick(() => this.$refs.newItemName.focus())
    },
    create () {
      axios.post('/api/admin/' + this.path + '/update', this.newElement).then(response => {
        if (response.data.id) {
          this.newElement = response.data
          if (!this.elements) {
            this.elements = []
          }
          this.elements.unshift(this.newElement)
        }
        this.resetNewElement()
        this.newItemFormVisible = false
        this.$notify({
          group: 'app',
          title: 'Success',
          text: 'Created new item',
          type: 'success'
        })
      }).catch(error => {
        console.log(error.response)
        this.$notify({
          group: 'app',
          title: 'Could not create element.',
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
    search (term) {
      this.$emit('searched', term)
      if (this.lastRequest !== null) {
        this.lastRequest.cancel()
        this.loading = true
      }
      this.searchTerm = term
      this.loadingBuffer = true
      setTimeout(() => {
        if (this.loadingBuffer === true && this.lastRequest !== null) {
          // this.elements = null
          // this.loadingBuffer = false
          this.loading = true
        }
      }, 100)
      var CancelToken = axios.CancelToken
      this.lastRequest = CancelToken.source()
      let source = this.lastRequest
      const p = this.fixedPath ? this.fixedPath + '?c=' + this.searchTerm : '/api/admin/' + this.path + this.listQuery +
        this.searchTerm
      axios.get(p, {
        cancelToken: source.token
      }).then(response => {
        this.lastRequest = null
        this.elements = response.data || []
        this.loading = false
        this.loadingBuffer = false
      }, error => {
        if (axios.isCancel(error)) {
        } else {
          this.lastRequest = null
          this.lastRequest = null
          this.loading = false
          this.loadingBuffer = false
        }
      })
      this.isLastRequest = true
      // }, 100)
    }
  }
}
</script>

<style lang="scss" scoped>

  .loading /deep/ .spinner {
    position: relative;
    background: transparent;

    .sk-circle {
      top: 20px !important;
      margin-top: 30px !important;
    }
  }

  .loading /deep/ .sk-circle .sk-child:before {
    background-color: #aaaaaa !important;
  }

  .new-item-form {
    border-bottom: 1px solid #dadada;
  }

  .list-group-item {
    border-left: 0;
    border-right: 0;
    line-height: 1;

    &[data-index="0"] {
      border-top: 0;
    }

    .lg-small h5 {
      font-weight: 400;
      font-size: .9rem;
    }
  }

  .list-group-icon {
    color: #747782;
    font-size: 2rem;
  }

  .flex-col-truncate {
    overflow: hidden;
    min-width: 0;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
