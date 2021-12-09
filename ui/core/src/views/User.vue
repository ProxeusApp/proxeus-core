<template>
  <div v-if="user">
    <vue-headful :title="$t('User title prefix','Proxeus - ')+(user.name || $t('User title', 'User'))"/>
    <top-nav :title="user.name || $t('User title', 'User')" :returnToRoute="{name:'Users'}">
      <button slot="buttons" v-if="user && app.userIsUserOrHigher()" style="height: 40px;"
              @click="app.exportData('&id='+user.id, null, '/api/user/export', 'User_'+user.id)" type="button"
              class="btn btn-primary ml-2">
        <span>Export</span></button>
    </top-nav>
    <div class="container-fluid">
      <div class="row mb-3">
        <div class="col-md-5">
          <div class="profile-img-main" style="width: 100%;" :style="[topRightInnerFallback?{'opacity':0.2}:'']">
            <div class="trp-pchange mt-2" style="display: inline-block;border: 2px solid whitesmoke;width: 100%;">
              <img v-show="topRightInnerFallback === false" @error="errorForTopRightInner" :src="rawSrc+profilePhotoSrc"
                   style="width: 100%;"/>
              <i v-show="topRightInnerFallback" class="material-icons" style="font-size: 25vw;color: grey;">person</i>
            </div>
          </div>
        </div>
        <div class="col-md-7">
          <div class="form-group mt-3">
            <animated-input :max="80" :label="$t('Name')" v-model="user.name"/>
            <animated-input :max="100" :label="$t('More about me')" v-model="user.detail"/>
            <animated-input :max="100" :disabled="true" :label="$t('Email')" v-model="user.email"/>
            <animated-input :max="100" :disabled="true" :label="$t('Ethereum Address')" v-model="user.etherPK"/>
            <br>
            <div class="sub-title">{{$t('Privacy settings')}}</div>
            <checkbox :label="$t('Want to be found')" :disabled="true" v-model="user.wantToBeFound"/>
            <span class="text-muted"
                  style="white-space: normal;">{{$t('Want to be found explanation','Uncheck this property if you want to be found only by your blockchain address.')}}</span>
          </div>
          <div class="form-group">
            <label>Role</label>
            <simple-select v-model="user.role" :idProp="'role'" :labelProp="'name'" :options="roles"/>
            <small id="roleHelp" class="form-text light-text">{{$t('Select the users Role')}}</small>
          </div>
          <div class="form-group">
            <api-key :user="user"/>
          </div>
          <!-- Legacy markup (real row/col features not applicable here) -->
          <div class="col-sm-12">
            <hr>
          </div>
          <form id="userForm" class="form-compiled" v-if="userForm">{{userForm()}}</form>
          <button type="button" class="btn btn-primary" :class="{saving:saving}"
                  @click="saveUserForm" v-if="user">
            Save
          </button>
        </div>
      </div>
      <!--      <user-settings-modal :userSrc="user.data.userSrc" v-if="user" @saved="updateUserForm"></user-settings-modal>-->
    </div>
  </div>
</template>

<script>
import TopNav from '@/components/layout/TopNav'
// import FT_FormBuilderCompiler from '../libs/legacy/formbuilder-compiler'

import SimpleSelect from '../components/SimpleSelect'
import AnimatedInput from '../components/AnimatedInput'
import Checkbox from '../components/Checkbox'
import formChangeAlert from '../mixins/form-change-alert'
import mafdc from '@/mixinApp'
import ApiKey from '../components/user/ApiKey'

export default {
  mixins: [mafdc, formChangeAlert],
  name: 'user',
  props: ['id'],
  components: {
    ApiKey,
    Checkbox,
    AnimatedInput,
    SimpleSelect,
    TopNav
  },
  data () {
    return {
      saving: false,
      user: null,
      userForm: null,
      rawSrc: '/api/profile/photo',
      profilePhotoSrc: '',
      topRightFallback: false,
      topRightInnerFallback: false
    }
  },
  computed: {
    roles () {
      return this.app.roles
    }
  },
  mounted () {
    this.updateUserForm()
    // this.loadComponents(() => this.updateUserForm())
  },
  watch: {
    userForm () {
      this.fillUserForm()
    }
  },
  created () {
    this.$root.$on('profile-photo-update', this.updateProfilePhoto)
  },
  beforeDestroy () {
    this.$root.$off('profile-photo-update', this.updateProfilePhoto)
  },
  methods: {
    errorForTopRight () {
      this.topRightFallback = true
    },
    errorForTopRightInner () {
      this.topRightInnerFallback = true
    },
    updateProfilePhoto (ev) {
      this.topRightFallback = false
      this.topRightInnerFallback = false
      this.profilePhotoSrc = '?id=' + this.user.id + '&v=' + Date.now()
    },
    loadComponents (done) {
      axios.get('/api/admin/form/component?l=1000').then(res => {
        this.components = res.data
        done()
      }, (err) => {
        this.app.handleError(err)
      })
    },
    saveUserForm () {
      this.saving = true
      // let formData = $('#userForm').serializeFormToObject()
      const data = {
        name: this.user.name,
        detail: this.user.detail,
        role: parseInt(this.user.role)
        // data: {
        //   user: formData
        // }
      }
      axios.post('/api/admin/user/update?id=' + this.user.id, data).then(res => {
        this.user = res.data
        this.snapshot(this.user)
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Saved user form'),
          type: 'success'
        })
        if (this.user.id === this.app.me.id) {
          this.app.loadMe()
        }
      }, (err) => {
        this.app.handleError(err)
        this.saving = false
      })
    },
    hasUnsavedChanges () {
      return !this.compare(this.user)
    },
    updateUserForm () {
      const self = this
      axios.get('/api/admin/user/' + self.$route.params.id).then(response => {
        self.user = response.data
        self.snapshot(self.user)
        self.profilePhotoSrc = '?id=' + self.user.id
      }, (error) => {
        this.app.handleError(error)
        if (error.response && error.response.status === 404) {
          self.$_error('NotFound', { what: 'User' })
        } else {
          self.$notify({
            group: 'app',
            title: self.$t('Error'),
            text: self.$t('Could not load user'),
            type: 'error'
          })
          self.$router.push({ name: 'Users' })
        }
      })
    },
    fillUserForm () {

    }
  }
}
</script>

<style scoped>
  .profile-img-main {
    text-align: center;
    cursor: default;
  }

  .profile-img-main > * {
    cursor: default;
  }

  @-webkit-keyframes sk-bounce {
    0%, 100% {
      -webkit-transform: scale(0.0)
    }
    50% {
      -webkit-transform: scale(1.0)
    }
  }

  @keyframes sk-bounce {
    0%, 100% {
      transform: scale(0.0);
      -webkit-transform: scale(0.0);
    }
    50% {
      transform: scale(1.0);
      -webkit-transform: scale(1.0);
    }
  }
</style>
