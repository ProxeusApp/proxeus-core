<template>
<span v-if="showSignInBtn()||hasContent()" ref="myroot" class="ml-3 trp-main"
      style="display: inline-block;position:relative;">
        <div v-if="hasContent()">
             <div v-if="path() === '/'" class="return-workflows">
               <a href="/user/document" class="btn btn-primary text-white">{{$t('Return to Workflows', 'Return to Workflows')}}</a>
             </div>
            <img v-show="topRightFallback === false" @click="clickTopRightPhoto($event)" @error="errorForTopRight"
                 class="trp-photo"
                 :src="rawSrc+profilePhotoSrc"/>
        <i v-show="topRightFallback" @click="clickTopRightPhoto($event)" class="trp-photo material-icons" style="
    font-size: 35px;text-align: center;vertical-align: middle;background: rgb(249, 249, 249);
    border: 2px solid rgb(236, 236, 236);">person</i>
        </div>
        <table style="margin-right: 15px;" v-else-if="showSignInBtn()"><tbody><tr>
            <td class="tdmin impcnt"><a href="/register" class="btn btn-primary text-white">{{$t('Sign up')}}</a></td>
            <td class="tdmin impcnt"><a style="margin-left: 14px;" href="/login"
                                        class="btn btn-primary text-white">{{$t('Sign in','Log in')}}</a></td>
        </tr></tbody></table>
        <span v-show="showOverlay" class="trp-overlay">
            <div class="trp-overlay-body">
                <table v-if="hasContent()">
                    <tr>
                        <td class="tdmin">
                            <div @click="openProfileImgDialog" class="trp-pchange">
                                <img v-show="topRightInnerFallback === false" style="width: 100px;"
                                     class="trp-photo-inner" @error="errorForTopRightInner"
                                     :src="rawSrc+profilePhotoSrc"/>
                                <i v-show="topRightInnerFallback" class="trp-photo-inner material-icons">person</i>
                                <span class="trp-pchange-btn">{{$t('Change profile photo', 'Change')}}</span>
                            </div>
                        </td>
                        <td class="impcnt" style="max-width: 280px;">
                            <div class="px-2 pb-1 font-weight-bold">{{ name() }}</div>
                            <div style="font-size: 10px;" class="px-2 pb-1 light-text mt-1"
                                 v-if="me().etherPK">{{ me().etherPK }}</div>
                            <div style="font-size: 10px;" class="px-2 pb-1 light-text mt-1"
                                 v-else-if="me().detail">{{ me().detail }}</div>
                            <div style="font-size: 10px;" class="px-2 pb-1 light-text mt-1" v-else>-</div>
                        </td>
                    </tr>
                </table>
            </div>
            <div class="trp-overlay-footer">
                <table style="width: 100%;">
                    <tbody>
                    <tr>
                        <td style="text-align:left"><button @click="openProfileModal"
                                                            class="btn btn-primary">{{$t('User account', 'Account')}}</button></td>
                        <td style="text-align:right"><button @click="logout"
                                                             class="btn btn-secondary">{{$t('Sign out')}}</button></td>
                    </tr>
                    </tbody>
                </table>
                <div class="clearfix"></div>
            </div>
        </span>
        <change-profile-photo-modal :setup="setupProfileImgDialog"/>
        <profile-modal :setup="setupProfileModal"/>
    </span>
</template>

<script>
import ProfileModal from './ProfileModal'
import mafdc from '@/mixinApp'
import ChangeProfilePhotoModal from '../ChangeProfilePhotoModal'

export default {
  mixins: [mafdc],
  components: {
    ChangeProfilePhotoModal,
    ProfileModal
  },
  name: 'top-right-profile',
  props: {
    options: {
      type: Array,
      default: () => []
    }
  },
  data () {
    return {
      rawSrc: '/api/my/profile/photo',
      profilePhotoSrc: '',
      topRightFallback: false,
      topRightInnerFallback: false,
      showOverlay: false,
      openProfileImgDialog: false,
      openProfileModal: false,
      hasContentFlag: false
    }
  },
  created () {
    this.$root.$on('profile-photo-update', this.updateProfilePhoto)
  },
  beforeDestroy () {
    this.$root.$off('profile-photo-update', this.updateProfilePhoto)
  },
  methods: {
    hasContent () { // to prevent from the settings dialog being gone when editing
      if (!this.hasContentFlag && this.name()) {
        this.hasContentFlag = true
      }
      return this.hasContentFlag
    },
    name () {
      return this.me().email || this.me().name
    },
    me () {
      if (this.app.me) {
        return this.app.me
      }
      return {}
    },
    errorForTopRight () {
      this.topRightFallback = true
    },
    errorForTopRightInner () {
      this.topRightInnerFallback = true
    },
    showSignInBtn () {
      return this.path() !== '/login'
    },
    path () {
      return window.location.pathname
    },
    setupProfileImgDialog (opid) {
      this.openProfileImgDialog = opid
    },
    setupProfileModal (o) {
      this.openProfileModal = o
    },
    updateProfilePhoto (ev) {
      this.topRightFallback = false
      this.topRightInnerFallback = false
      this.profilePhotoSrc = '?v=' + Date.now()
    },
    clickTopRightPhoto (ev) {
      if (!this.showOverlay) {
        this.showOverlay = true
        document.addEventListener('click', this.hide, true)
      } else {
        this.justHide(ev)
      }
    },
    hide (ev) {
      if (this.showOverlay) {
        if (ev && this.canHide(ev.target)) {
          this.justHide(ev)
          return true
        }
      }
      return false
    },
    justHide (ev) {
      this.showOverlay = false
      document.removeEventListener('click', this.hide)
    },
    canHide (n) {
      if (n) {
        while (true) {
          if (n === this.$refs.myroot) {
            return false
          }
          if (n === undefined || n === document.body) {
            return true
          }
          n = n.parentNode
        }
      }
      return false
    },
    logout () {
      axios.post('/api/logout', null).then(response => {
        window.location.replace('/')
      }, (err) => {
        this.app.handleError(err)
      })
    }
  }
}
</script>
<style>
  i.trp-photo-inner.material-icons {
    font-size: 100px;
  }

  .trp-overlay-body {
    padding: 25px;
    background: white;
  }

  .trp-overlay-footer {
    padding: 10px 14px;
    position: relative;
    bottom: 0;
    background: #ececec;
    border: 1px solid #d2d2d2;
  }

  .trp-overlay {
    position: absolute;
    top: 60px;
    right: 0;
    -webkit-box-shadow: 0px 0px 15px 2px rgba(0, 0, 0, 0.14);
    -moz-box-shadow: 0px 0px 15px 2px rgba(0, 0, 0, 0.14);
    box-shadow: 0px 0px 15px 2px rgba(0, 0, 0, 0.14);
    border: 1px solid rgb(226, 226, 226);
    z-index: 10;
  }

  .trp-photo {
    width: 45px;
    height: 45px;
    border-radius: 100%;
    border: 2px solid #f3f3f3;
    cursor: pointer;
    background: white;
  }

  .trp-photo-inner {
    width: 100px;
  }

  .trp-pchange {
    background: #e2e2e2;
    position: relative;
    border-radius: 48%;
    overflow: hidden;
    border: 2px solid #f5f5f5;
    cursor: pointer;
  }

  span.trp-pchange-btn {
    background: #000000ba;
    position: absolute;
    bottom: 0px;
    left: 0;
    width: 100%;
    text-align: center;
    color: white;
    font-size: 12px;
    padding: 4px;
    padding-bottom: 7px;
  }

  .return-workflows {
    float: left;
    margin-right: 30px;
  }
</style>
