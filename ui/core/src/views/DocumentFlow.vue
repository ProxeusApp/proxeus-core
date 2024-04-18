<template>
<div class="document-view">
  <vue-headful :title="$t('Flow title prefix','Proxeus - ')+(name || $t('Document title','Document'))"/>
  <top-nav :returnToRoute="{name:'Documents'}" bg="#ffffff" :title="name" :sidebarToggler="false" class="border-bottom-0 bg-white">
    <language-drop-down slot="buttons" style="margin-top: 3px;"/>
  </top-nav>
  <div class="container-fluid document-inner-view bg-light" v-if="status && name">
    <div class="row h-100">
      <div class="col-md-3 mt-3">
        <ul class="nav mod_stepnav flex-column nav-pills">
          <li v-for="(step, index) in status.steps" :key="index" class="nav-item">
            <span class="btn-link nav-link text-truncate"
                  :class="{active:index === status.steps.length - 1 && isConfirmationStep === false}">
              {{ step.name }}
            </span>
          </li>
          <li class="nav-item">
            <span class="btn-link nav-link active text-truncate" v-if="isConfirmationStep">{{$t('Confirmation')}}</span>
          </li>
        </ul>
      </div>
      <div class="col-md-6 mt-3">
        <div class="document-form card bg-transparent w-100 border-0" style="position:relative;">
          <div v-show="loading" class="v-fade"
               style="position:absolute;top:0;left:0;bottom:0;right:0;margin-top: 50px;">
            <div>
              <spinner background="transparent"/>
            </div>
            <div class="processing-hint">
              {{$t('Document processing hint', 'Processing your request. This might take a while.')}}

            </div>
          </div>
          <div><!--to prevent from flickering-->
            <div class="card-header border-bottom-0" style="visibility: hidden;">
              <h2 class="text-white my-0 py-0">A</h2>
            </div>
            <div class="card-body card-form-body bg-white border-0">
            </div>
          </div>
          <div v-show="loading === false && isConfirmationStep === false" class="v-fade"
               style="position:absolute;top:0;left:0;bottom:0;right:0;">
            <div class="card-header border-bottom-0">
              <h2 class="text-white my-0 py-0">{{ getCurrentName()}}</h2>
            </div>
            <div id="docForm" class="card-body card-form-body bg-white border-0">
              <div class="form form-compiled" v-append.sync="formSrc" @appended="appended"
                   v-show="isConfirmationStep === false"></div>
            </div>
          </div>

          <div v-show="loading === false && isConfirmationStep" class="v-fade"
               style="position:absolute;top:0;left:0;bottom:0;right:0;">
            <div v-if="!isConnectedAccount" class="alert alert-warning mb-3" role="alert">
              {{$t('Document register help', 'Make sure that you are logged into the correct account in MetaMask. The Document-Registration must be made from the ethereum-account that you have connected in your account settings. Switch Metamask account and make sure that you have set your ethereum address in your account account settings on the top right.')}}
            </div>
            <div class="card-header border-bottom-0">
              <h2 class="text-white my-0 py-0">{{$t('Confirmation')}}</h2>
            </div>
            <div class="card-body card-form-body bg-white border-0">
              <div class="last-step">
                <form v-on:submit.prevent="confirm">
                  <div class="mb-2" v-show="submitting === false && pricePerFile">
                    <div class="easy-read">{{ pricePerFileFormatted }} per File.</div>
                    <span class="light-text">(Ethereum transaction cost excluded)</span>
                  </div>
                  <div class="mb-2" v-show="submitting === false">
                    <animated-input label="Additional Text data" :max="32" v-model="data"/>
                    <span class="light-text">Information entered here will be publicly visible to anyone.</span>
                  </div>
                  <button type="submit"
                          class="btn btn-primary"
                          :disabled="submitting || !isConnectedAccount" v-show="submitting === false">Submit
                  </button>
                  <spinner background="transparent" style="position:relative;" :margin="60" v-if="submitting"></spinner>
                  <div class="text-center">
                    <transition name="fade" mode="out-in">
                      <p class="text-primary px-5"
                         v-if="approvingXES">Please confirm the release of XES from your wallet in MetaMask.</p>
                      <p class="text-primary px-5"
                         v-if="confirmingDocs">Please confirm the registration of your document on the blockchain in MetaMask.</p>
                    </transition>
                    <p class="light-text"
                       v-show="submitting">If MetaMask doesn't pop up automatically, click on the MetaMask icon in the top right corner of your browser.</p>
                  </div>
                </form>
              </div>
            </div>
          </div>

          <div class="card-footer bg-primary d-flex flex-row">
            <button type="submit" class="btn btn-primary border mr-auto" @click="previous"
                    :disabled="loading"
                    v-if="status.hasPrev && isConfirmationStep === false">
              <i class="material-icons">chevron_left</i>
            </button>
            <button type="submit" class="btn btn-primary border mr-auto" @click="cancelConfirmationStep"
                    :disabled="loading"
                    v-else-if="isConfirmationStep === true && submitting === false">
              <i class="material-icons">chevron_left</i>
            </button>

            <button type="submit" class="btn btn-primary border ml-auto" @click="next"
                    :disabled="loading"
                    v-if="status.hasNext">
              <i class="material-icons">chevron_right</i>
            </button>
            <button type="submit" class="btn btn-primary border ml-auto" @click="confirmationStep"
                    :disabled="loading"
                    v-else-if="status.hasNext === false && isConfirmationStep === false">
              <i class="material-icons">chevron_right</i>
            </button>
          </div>
        </div>
      </div>
      <div class="col-md-3 mt-3 document-scroll-view" :class="{'has-pdf':documentPreviews != null}">
        <template v-if="documentPreviews">
          <div class="d-flex flex-row flex-wrap">
            <template v-for="doc in documentPreviews">
              <pdf-preview v-if="anyLangAvailable(doc)"
                           :key="doc.id"
                           :name="doc.name"
                           :languages="doc.langs"
                           :doc="doc" :wfId="id"
                           :locale="locale"
                           :langSelectorVisible="isConfirmationStep === false"/>
            </template>
          </div>
        </template>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import TopNav from '@/components/layout/TopNav'
// import PdfModal from '@/components/document/PdfModal'
import PdfPreview from '@/components/document/PdfPreview'
// import ButtonSpinner from '@/components/ButtonSpinner'
import Spinner from '@/components/Spinner'
import LanguageDropDown from '@/views/appDependentComponents/LanguageDropDown'

import formCompilerAdapter from '../libs/formcompiler-adapter'
// import TopRightProfile from '../components/user/TopRightProfile'
import mafdc from '@/mixinApp'
import AnimatedInput from '../components/AnimatedInput'

export default {
  mixins: [mafdc, formCompilerAdapter],
  name: 'document-flow',
  components: {
    AnimatedInput,
    // TopRightProfile,
    // ButtonSpinner,
    Spinner,
    // PdfModal,
    PdfPreview,
    TopNav,
    LanguageDropDown
  },
  data () {
    return {
      name: '',
      components: null,
      data: '',
      status: null,
      loading: false,
      formSrc: null,
      documentPreviews: null,
      loadingPreview: false,
      isConfirmationStep: false,
      isLastStep: true,
      submitting: false,
      approvingXES: false,
      confirmingDocs: false,
      pricePerFile: undefined,
      nonce: undefined,
      accountEthAddress: '',
      isConnectedAccount: false,
      me: null
    }
  },
  computed: {
    pricePerFileFormatted () {
      return this.pricePerFile ? this.pricePerFile + ' XES' : ''
    },
    id () {
      return this.$route.params.id
    },
    locale () {
      return this.app.getSelectedLang()
    },
    changeOptions () {
      return {
        url: '/api/document/' + this.id + '/data',
        fileUrl: '/api/document/' + this.id + '/file'
      }
    },
    blockchainNet () {
      return this.app.blockchainNet
    },
    wallet () {
      return this.app.wallet
    }
  },
  created () {
    this.loadPricePerFile()
    this.initialScreen()
  },
  async mounted () {
    this.setAccountEthAddress()
    this.checkConnectedAccount()
  },
  beforeRouteLeave (to, from, next) {
    if (this.submitting) {
      const answer = window.confirm(this.$t('Blockchain progress alert warning',
        'Its highly recommended to stay on this page until the metamask transaction has been confirmed, else your payment might not be successfully processed. ' +
        '\n\nClick on "Cancel" to stay on this page or "OK" to leave.'))
      if (answer) {
        next()
      } else {
        next(false)
      }
    } else {
      next()
    }
  },
  methods: {
    async setAccountEthAddress () {
      const response = await axios.get('/api/me')
      if (response.data.etherPK) {
        this.accountEthAddress = response.data.etherPK
      }
    },
    async checkConnectedAccount () {
      let response2 = null
      try {
        response2 = await axios.get('/api/me')
        if (response2.data.etherPK) {
          this.accountEthAddress = response2.data.etherPK

          await this.app.wallet.getClientProvidedNetwork()
          this.me = window.ethereum.selectedAddress || this.app.wallet.getCurrentAddress()

          if (this.me === '' || this.me === null ||
            this.accountEthAddress === '' || this.accountEthAddress === null) {
            this.isConnectedAccount = false
            return
          }
          // check lowercase because window.ethereum.selectedAddress returns lowercase hash
          this.isConnectedAccount = this.me.toLowerCase() === this.accountEthAddress.toLowerCase()
        }
      } catch (e) {
        console.log(e)
      }
    },
    xesTransfer (cb) {
      return new Promise((resolve, reject) => {
        this.wallet.approveXES(this.wallet.serviceConfig.PROXEUS_FS_ADDRESS, 1)
          .then((result) => {
            resolve(result)
          })
          .catch(e => {
            reject(e)
          })
      })
    },
    getCurrentName () {
      try {
        return this.status.targetName
      } catch (e) {
        return ''
      }
    },
    anyLangAvailable (doc) {
      if (doc && doc.langs) {
        for (let i = 0; i < doc.langs.length; i++) {
          if (this.app.isLangAvailable(doc.langs[i])) {
            return true
          }
        }
      }
      return false
    },
    async loadPricePerFile () {
      try {
        this.pricePerFile = 0
      } catch (e) {
        console.log(e)
      }
    },
    initialScreen () {
      this.loading = true
      axios.get('/api/document/' + this.id).then((response) => {
        this.loading = false
        this.handleDocumentResponse(response)
      }, (error) => {
        this.app.handleError(error)
        this.loading = false
        try {
          let printedError = false
          if (error.response &&
            error.response.status > 400 &&
            error.response.data &&
            typeof error.response.data === 'string') {
            printedError = true
            this.$notify({
              group: 'app',
              title: this.$t('Error'),
              text: error.response.data + '. Please try again or if the error persists contact the platform operator.',
              type: 'error'
            })
          }
          if (error.response && error.response.status === 404) {
            this.$_error('NotFound', { what: 'Document' })
            this.$router.history.updateRoute(this.$router.match('*'))
          } else if (error.response && error.response.status === 401) {
            redirectToLogin()
          } else if (error.response && error.response.status === 422) {
            this.handleDocumentResponse(error.response)
          } else {
            if (!printedError && error.response && typeof error.response.data === 'string') {
              this.$notify({
                group: 'app',
                title: this.$t('Error'),
                text: error.response.data + '. Please try again or if the error persists contact the platform operator.',
                type: 'error'
              })
              return
            }
            if (!printedError) {
              this.$notify({
                group: 'app',
                title: this.$t('Error'),
                text: this.$t('An unexpected error occurred. Please try again or if the error persists contact the platform operator.'),
                type: 'error'
              })
            }
          }
        } catch (e) {
          console.log(e)
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('An unexpected error occurred. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
        }
      })
    },
    previous () {
      this.loading = true
      axios.get('/api/document/' + this.id + '/prev').then((response) => {
        this.loading = false
        this.handleDocumentResponse(response)
      }, (err) => {
        this.app.handleError(err)
      })
    },
    next (lastStep = false, callback) {
      this.loading = true
      axios.post('/api/document/' + this.id + '/next' + (lastStep === true ? '?confirm' : ''), this.getFormData())
        .then((response) => {
          this.loading = false
          if (response.status >= 200 && response.status <= 299) {
            if (lastStep === true) {
              this.handleDocumentResponse(response)
              callback()
              return
            }
            this.handleDocumentResponse(response)
          }
        }, (err) => {
          this.app.handleError(err)
          this.loading = false
          if (err.response && typeof err.response.data === 'string') {
            this.$notify({
              group: 'app',
              title: this.$t('Error'),
              text: err.response.data + '. Please try again or if the error persists contact the platform operator.',
              type: 'error'
            })
            return
          }
          err.response && err.response.data && $('.form-compiled > form').showFieldErrors(err.response.data)

          this.$nextTick(() => {
            document.querySelector('.error') && this.$scrollTo('.error', 500, {
              container: '#docForm',
              offset: -150
            })
          })
        })
    },
    confirmationStep () {
      this.next(true, () => {
        this.isConfirmationStep = true
      })
    },
    cancelConfirmationStep () {
      this.isConfirmationStep = false

      // Handle case where there is no formdata (only a template in the workflow) -> automatically go to confirm
      if (!this.status.data) {
        if (this.status.hasNext === false) {
          this.confirmationStep()
        } else if (this.status.hasNext === true) {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not complete. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
        }
      }
    },
    submitDoc () {
      const templateLanguageConfig = {}
      this.$children.forEach(pdfPreviewComp => {
        if (pdfPreviewComp.templateLanguage) {
          templateLanguageConfig[pdfPreviewComp.doc.id] = pdfPreviewComp.templateLanguage
        }
      })
      axios.post('/api/document/' + this.id + '/next?final', templateLanguageConfig).then((response) => {
        if (response.data.id) {
          this.$router.push({ name: 'DocumentViewer', params: { id: response.data.id } })
        } else {
          this.submitting = false
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not complete. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
        }
      }, (err) => {
        this.app.handleError(err)
        this.submitting = false
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Could not complete. Please try again or if the error persists contact the platform operator.'),
          type: 'error'
        })
      })
    },
    verifyHash (hash) {
      /* sha3 */
      return this.wallet.verifyHash(hash)
    },
    async confirmDoc (hash) {
      return new Promise((resolve, reject) => {
        const account = this.wallet.getCurrentAddress()

        if (account === null) {
          this.submitting = false
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Please login to your wallet.'),
            type: 'error'
          })
          reject(new Error(this.$t('Please login to your wallet.')))
        }

        this.nonce = this.wallet.proxeusFS.web3.eth.getTransactionCount(account)
        this.nonce++
        this.wallet.proxeusFS.createFileUndefinedSigners({
          from: account,
          hash,
          data: (this.data === '' ? '0x00' : this.wallet.proxeusFS.web3.utils.fromAscii(this.data)),
          mandatorySigners: 0,
          expiry: 0,
          providers: [],
          nonce: this.nonce,
          xes: 0
        }).then((result) => {
          resolve(result.transactionHash)
        }).catch((error) => {
          console.warn(error.stack)
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not register document. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
          this.submitting = false
          reject(new Error(this.$t('Could not create transaction. Please try again or if the error persists contact the platform operator.')))
        })
      })
    },
    async confirm () {
      const numFiles = this.status.docs ? this.status.docs.length : 0

      if (numFiles === 0) {
        this.submitDoc()

        return
      }

      let clientProvidedNet
      try {
        clientProvidedNet = await this.app.wallet.getClientProvidedNetwork()
      } catch (e) {
        console.log(e)
        return
      }
      if (this.blockchainNet !== clientProvidedNet) {
        this.$notify({
          group: 'app',
          title: 'Error',
          text: 'Could not register document. Please switch to ' + this.blockchainNet + ' Network in MetaMask and try again.',
          type: 'error',
          duration: 5000
        })
        return
      }
      this.submitting = true

      try {
        if (numFiles > 0) {
          this.confirmingDocs = true
          await this.asyncForEach(this.status.docs, async (doc) => {
            await this.confirmDoc(doc.hash)
          })
          this.confirmingDocs = false
        }

        this.submitting = false
        this.submitDoc()
      } catch (e) {
        console.log(e)
        this.submitting = false
        this.confirmingDocs = false
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: this.$t('Please make sure you are logged in to Metamask and refresh.'),
          type: 'error'
        })
      }
    },
    async asyncForEach (array, callback) {
      for (let index = 0; index < array.length; index++) {
        await callback(array[index], index, array)
      }
    },
    handleDocumentResponse (response) {
      $('.form-compiled').hide()

      if (!response.data) {
        return
      }
      if (response.data.name) {
        this.name = response.data.name
      }
      this.status = response.data.status
      if (this.status && this.status.docs) {
        this.documentPreviews = []
        this.status.docs.forEach(doc => {
          doc.loaded = false
          this.documentPreviews.unshift(doc)
        })
      }
      if (this.isConfirmationStep === false && this.status.hasNext === false) {
        this.next(true, () => {
          this.isConfirmationStep = true
        })
      } else if (this.status.data) {
        this.destroyDatepicker()
        this.compile(this.status.data, (form) => {
          this.formSrc = form
        })
      } else {
        if (this.isConfirmationStep === false && this.status.hasNext === false) {
          this.confirmationStep()
        } else if (this.isConfirmationStep === false && this.status.hasNext === true) {
          this.$notify({
            group: 'app',
            title: this.$t('Error'),
            text: this.$t('Could not complete. Please try again or if the error persists contact the platform operator.'),
            type: 'error'
          })
        }
      }
    },
    destroyDatepicker () {
      $('.simple-date-field').each(function () {
        if (this._flatpickr) {
          this._flatpickr.destroy()
        }
      })
    },
    appended () {
      this.initForm()
    },
    initForm () {
      const self = this
      const $form = $('.form-compiled > form')

      if ($form.length) {
        $('.card-body').scrollTop(0)

        if (self.status.userData) {
          $form.fillForm(self.status.userData)
        }

        $form.assignSubmitOnChange(self.changeOptions)
        $form.on('formFieldsAdded', function (event, parent) {
          if (parent && parent.length) {
            parent.assignSubmitOnChange(self.changeOptions)
          }
        })
        $('.form-compiled').fadeIn()
      }
    },
    getFormData () {
      const excludeFileFormFields = true
      return $('.form-compiled > form').serializeFormToObject(excludeFileFormFields)
    },
    compile (form, cb) {
      if (this.components) {
        this.adapterCompile(form, this.components, cb)
      } else {
        axios.get('/api/form/component?l=1000').then((res) => {
          this.components = res.data
          this.adapterCompile(form, this.components, cb)
        }, (err) => {
          this.app.handleError(err)
        })
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "../assets/styles/variables";
  @import "~bootstrap/scss/mixins";

  hr {
    border-top: 1px solid rgb(6, 42, 133) !important;
  }

  .form-heading {
    border-bottom: 1px solid #40e1d1;
  }

  .document-inner-view {
    overflow: auto;
    height: calc(100vh - 65px);
  }

  @media (max-width: 768px) {
    .form-heading {
      color: $primary;
    }
  }

  .navbar-nav .ss-sel-main {
    top: 4px;
  }

  .navbar-brand svg {
    height: 50px;
  }

  .navbar {
    .dropdown-menu {
      right: 0;
      left: auto;
    }
  }

  ::v-deep .navbar .dropdown .nav-link {
    color: white !important;
  }

  .card {
    border-radius: 0;

    .card-footer {
      border-radius: 0;
      box-shadow: 0 -20px 20px -10px rgba(0, 0, 0, .1);
      z-index: 9;
      min-height: 64px;
    }
  }

  .document-scroll-view {
    //height: calc(100vh - 68px) !important;
    overflow-y: visible;
    position: relative;
  }

  .card-form-body {
    height: calc(100vh - 220px);
    overflow-y: auto;
    overflow-x: hidden;

    &::-webkit-scrollbar {
      -webkit-appearance: none;
      width: 9px;
      height: 9px;
    }

    &::-webkit-scrollbar-thumb {
      border-radius: 7px;
      border: 1px solid white; /* Angleichen mit Hintergrundfarbe-nicht transparent! */
      background-color: rgba(0, 0, 0, .5);
    }

  }

  .document-form {
    position: relative;
    .card-header {
      background-color:$primary;
    }
    .processing-hint {
      position: absolute;
      margin: 180px auto;
      width: 100%;
      text-align: center;
      color: $primary;
    }
  }

  .paper-stack {
    border: none;
    box-shadow: 0 1px 5px rgba(0, 0, 0, 0.2); /* The top layer shadow */
  }

  .mod_stepnav {
    .nav-link {
      cursor: default;
      color: $gray-500;
      border-left: 4px solid transparent;

      &.active {
        font-weight: 500;
        color: $primary;
        border-left: 4px solid $info;
        background: transparent;
      }
    }
  }

  .fade-enter-active {
    transition: transform 0.5s, opacity 0.5s;
  }

  .fade-leave-active {
    transition: opacity 0.3s;
  }

  .fade-enter {
    opacity: 0;
    transform: translateY(-100%);
  }

  .fade-enter-to, .fade-leave {
    opacity: 1;
    transform: translateY(0);
  }

  .fade-leave-to {
    opacity: 0;
  }

  .v-fade {
    display: block !important; /* override v-show display: none */
    transition: opacity 0.5s;
  }

  .v-fade[style*="display: none;"] {
    opacity: 0;
    pointer-events: none; /* disable user interaction */
    user-select: none; /* disable user selection */
  }
</style>
