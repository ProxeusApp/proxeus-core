<template>
<div>
  <vue-headful :title="$t('Users title','Proxeus - Users')"/>
  <top-nav :title="$t('Users')">
  </top-nav>
  <div class="main-container">
    <b-modal v-model="modalShow" class="b-modal" :title="$t('Invite')" :ok-title="$t('Invite')"
             :cancel-title="$t('Cancel')"
             :header-bg-variant="'light'" @hide="onDialogHide">
      <div ref="inviteCont" class="form-group" style="text-align:left;">
        <animated-input :max="100" :label="'Email'" name="email" v-model="invite.email"/>
        <small class="text-muted">{{$t('Enter the email of the user you would like to invite.')}}</small>
        <div class="form-group">
          <label>Role</label>
          <simple-select :unselect="false" name="role" v-model="invite.role" :idProp="'role'" :labelProp="'name'"
                         :options="roles"/>
          <small
            class="text-muted">{{$t('Select the role the user is going to get, in case the invitation is going to be accepted.')}}
          </small>
        </div>
      </div>
      <template slot="modal-footer">
      <button @click="onDialogHide" class="btn btn-secondary">
        Cancel
      </button>
      <button @click="onDialogOk" class="btn btn-primary">
        Send invitation
      </button>
      </template>
    </b-modal>
    <list-group class="user-list" :prependFunc="prependFunc" icon="person" nodeType="user" path="user">
      <button slot="addBtn" @click="modalShow=true" type="button"
              class="btn btn-primary">
        Invite user
      </button>
      <!--                <button slot="addBtn" type="button" @click="toggleNewItemFormVisible" class="btn btn-primary btn-round plus-btn mshadow-dark">-->
      <!--                    <i class="material-icons">add</i>-->
      <!--                </button>-->
      <div slot="newItemForm"
           class="pb-3 new-item-form list-group-item-action bg-light container-fluid" v-show="newItemFormVisible"
           style="position:relative;">
        <i @click="newItemFormVisible=false" style="position: absolute;right: 10px;top:10px;cursor:pointer;"
           class="material-icons">clear
        </i>
        <h2 class="pt-3 pb-2">Create new User</h2>
        <div class="row">
          <div class="col-sm-6">
            <div class="form-group">
              <label for="newNameInput">Name</label>
              <input type="text" ref="newElementName" v-model.trim="newElement.name" name="newNameInput"
                     id="newNameInput" class="form-control" placeholder="Username" required>
            </div>
          </div>
          <div class="col-sm-6">
            <div class="form-group">
              <label for="userRoleSelect">User Role</label>
              <select v-model.trim="newElement.role" class="form-control" id="userRoleSelect"
                      aria-describedby="roleHelp" placeholder="role" required>
                <option v-for="role in roles" :value="role.role">{{ role.name }}</option>
              </select>
              <small id="roleHelp" class="light-text">The users role</small>
            </div>
          </div>
        </div>
        <div class="row">
          <div class="col-sm-12">
            <div class="form-group">
              <label for="etherAddressInput">Ethereum Account Address</label>
              <input type="text" v-model.trim="newElement.etherPK"
                     name="newElementEtherAddress" placeholder="0x245F..."
                     id="etherAddressInput" class="form-control" required>
              <small class="light-text">The users ethereum account address</small>
            </div>
          </div>
        </div>
        <div class="row">
          <div class="col-sm-12">
            <button type="button" @click="createUser" class="btn btn-primary"
                    :disabled="newElement.name === ''">Create user
            </button>
          </div>
        </div>
      </div>
      <template scope="element">
      <td class="tdmax impcnt" style="width:10%">
        <div>{{humanizeRole(element.role)}}</div>
        <small class="light-text">{{$t('Role')}}</small>
      </td>
      </template>
    </list-group>
  </div>
  <user-settings-modal :userSrc="userSrc" v-if="userSrc"></user-settings-modal>
</div>
</template>

<script>
import bModal from 'bootstrap-vue/es/components/modal/modal'
import bModalDirective from 'bootstrap-vue/es/directives/modal/modal'
import TopNav from '@/components/layout/TopNav'
import ListGroup from '../components/ListGroup'
import ListItem from '../components/ListItem'
import mafdc from '@/mixinApp'
import AnimatedInput from '../components/AnimatedInput'
import SimpleSelect from '../components/SimpleSelect'

export default {
  mixins: [mafdc],
  name: 'users',
  components: {
    SimpleSelect,
    AnimatedInput,
    ListGroup,
    ListItem,
    TopNav,
    'b-modal': bModal
  },
  directives: {
    'b-modal': bModalDirective
  },
  data () {
    return {
      prependItem: null,
      userSrc: null,
      newItemFormVisible: false,
      modalShow: false,
      newElement: {
        name: '',
        etherPK: '',
        role: 5
      },
      invite: { role: 1 }
    }
  },
  computed: {
    roles () {
      return this.app.roles
    }
  },
  methods: {
    cleanErr () {
      $(this.$refs.inviteCont).cleanFieldErrors()
    },
    onDialogHide () {
      this.modalShow = false
      this.invite = { role: 1 }
    },
    onDialogOk () {
      axios.post('/api/admin/invite', this.invite).then(res => {
        this.cleanErr()
        this.modalShow = false
        this.invite = { role: 1 }
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Invitation sent'),
          type: 'success'
        })
      }, (err) => {
        this.cleanErr()
        this.app.handleError(err)
        if (err.response && err.response.status === 422) {
          $(this.$refs.inviteCont).showFieldErrors({ errors: err.response.data })
        } else {
          if (err.response && err.response.data && typeof err.response.data === 'string') {
            this.$notify({
              group: 'app',
              title: this.$t('Warning'),
              text: err.response.data,
              type: 'warning'
            })
          } else {
            this.$notify({
              group: 'app',
              title: this.$t('Warning'),
              text: this.$t('There was an unexpected error. Please try again or if the error persists contact the platform operator.'),
              type: 'warning'
            })
          }
        }
      })
    },
    prependFunc (f) {
      this.prependItem = f
    },
    humanizeRole (role) {
      const unknownRole = this.$t('Unknown Role')
      if (this.app.roles) {
        const roleConfig = this.app.roles.find(r => {
          return r.role === role
        })
        return roleConfig ? roleConfig.name : unknownRole
      }
      return unknownRole
    },
    resetNewElement () {
      this.newElement = {
        name: '',
        walletAddress: '',
        role: 5
      }
    },
    createUser () {
      this.newElement.role = parseInt(this.newElement.role)
      this.newElement.id = ''
      axios.post('/api/admin/user/update', this.newElement).then(response => {
        if (response.data.id) {
          this.resetNewElement()
          if (this.prependItem) {
            this.prependItem(response.data)
          }
        }
        this.resetNewElement()
        this.newItemFormVisible = false
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: 'Created new user',
          type: 'success'
        })
      }, (e) => {
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: 'Could not create new user. Please try again or if the error persists contact the platform operator.\n',
          type: 'error'
        })
        this.app.handleError(e)
      })
    },
    toggleNewItemFormVisible () {
      this.newItemFormVisible = !this.newItemFormVisible
      this.$nextTick(() => this.$refs.newElementName.focus())
    }
  }
}
</script>

<style lang="scss">
  .user-list .mlist-group-icon {
    border: 2px solid #e8e8e8;
  }

  .modal-body {
    .hcbuild-main .tab-pane .panel .panel-body {
      height: calc(100vh - 250px - 120px) !important;
    }

    .hcbuild-main .hcbuild-workspace-body {
      height: calc(100vh - 145px - 120px) !important;
    }

    .hcbuild-workspace-test-main .panel-body {
      height: calc(100vh - 190px - 120px) !important;
    }
  }

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

  .flex-col-truncate {
    overflow: hidden;
    min-width: 0;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

</style>
