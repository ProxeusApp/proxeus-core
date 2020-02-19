<template>
<div class="document-view">
  <vue-headful :title="$t('Flow title prefix','Proxeus - ')+(name || $t('Document title','Document'))"/>
  <top-nav :returnToRoute="{name:'Documents'}" :sidebarToggler="false" class="border-bottom-0">
    <td slot="td-start" class="tdmin">
      <router-link :to="{name:'Documents'}" class="navbar-brand d-flex mr-auto">
        <svg width="140" viewBox="0 0 187 55" xmlns="http://www.w3.org/2000/svg"><title color="#fff">
          Proxeus</title>
          <g fill="none" fill-rule="evenodd">
            <path
                    d="M70 18.187h4.98c3.453 0 5.825 2.155 5.825 5.292 0 3.091-2.372 5.292-5.825 5.292h-2.255v6.042H70V18.187zm4.862 8.055c1.997 0 3.171-1.1 3.171-2.763 0-1.662-1.174-2.763-3.17-2.763h-2.138v5.526h2.137zm15.503 1.499h-1.574v7.072h-2.724V18.187h5.402c3.1 0 5.309 1.99 5.309 4.777 0 2.154-1.316 3.864-3.43 4.52l5.309 7.329h-3.265l-5.027-7.072zm.846-2.365c1.879 0 2.795-.913 2.795-2.318 0-1.429-.916-2.365-2.795-2.365h-2.42v4.683h2.42zm10.382 1.124c0-4.8 3.805-8.5 8.527-8.5 4.744 0 8.526 3.7 8.526 8.5s-3.782 8.5-8.526 8.5c-4.722 0-8.527-3.7-8.527-8.5zm14.282 0c0-3.325-2.467-5.924-5.755-5.924-3.289 0-5.755 2.599-5.755 5.924s2.466 5.924 5.755 5.924c3.288 0 5.755-2.599 5.755-5.924zm11.72-.117l-5.73 8.43h3.17l4.416-6.463 4.416 6.463h3.124l-5.708-8.43 5.38-8.196h-3.148l-4.087 6.159-4.087-6.159h-3.148l5.403 8.196zm14.517-8.196v16.626h9.889v-2.67h-7.164v-4.449h5.755v-2.622h-5.755v-4.215H152v-2.67h-9.889zm15.503 10.163c0 4.098 2.795 6.65 6.718 6.65s6.694-2.552 6.694-6.65V18.187h-2.724V28.42c0 2.576-1.598 4.004-3.97 4.004-2.373 0-3.993-1.428-3.993-4.004V18.187h-2.725V28.35zm18.462 3.817c1.058 1.709 3.101 2.833 5.614 2.833 3.218 0 5.309-1.92 5.309-4.73 0-2.318-1.55-3.77-3.829-4.753l-1.292-.586c-1.667-.702-2.42-1.475-2.42-2.552 0-1.148 1.081-1.803 2.209-1.803 1.198 0 1.95.445 2.819 1.475l1.879-1.944C185.144 18.656 183.618 18 181.69 18c-2.536 0-4.956 1.733-4.956 4.402 0 2.154 1.245 3.747 3.758 4.847l1.316.562c1.644.726 2.42 1.452 2.42 2.693 0 1.194-.94 1.92-2.561 1.92-1.456 0-2.748-.936-3.523-1.99l-2.068 1.733z"
                    fill="#062a85"></path>
            <path
                    d="M24.93 0c.394.072.824.109 1.218.253.43.145.86.29 1.253.507 6.09 3.546 12.214 7.056 18.268 10.638 2.042 1.194 3.08 3.076 3.116 5.428.036 7.056.036 14.111 0 21.131 0 2.352-1.074 4.234-3.08 5.428-6.09 3.582-12.179 7.092-18.268 10.674-.931.543-1.934.724-3.009.76 0-.108-.035-.217-.035-.326V0h.537zM0 16.174c.143-.506.25-1.013.43-1.483.25-.76.788-1.592 1.253-2.027v29.563c-.107-.145-.214-.253-.286-.362-.68-.904-1.146-1.918-1.29-3.04 0-.072-.035-.108-.071-.18C0 31.155 0 23.665 0 16.175zm22.029 38.247c-.574-.326-1.11-.615-1.648-.94-.036-.037-.072-.11-.108-.182V1.846c0-.325.072-.506.394-.65.43-.218.86-.471 1.362-.76V54.42zm-4.084-2.351c-.537-.327-1.074-.616-1.576-.906-.071-.036-.107-.18-.143-.289-.036-.072 0-.145 0-.217V4.198c0-.326.072-.507.358-.652.466-.217.896-.507 1.361-.796v49.32zM13.898 5.138v44.543c-.072-.036-.144-.036-.18-.072-.43-.254-.895-.507-1.325-.797-.107-.072-.143-.217-.215-.361-.035-.073 0-.145 0-.217V6.622c0-.362.108-.543.43-.724.394-.217.824-.507 1.29-.76zM9.779 47.329a30.57 30.57 0 0 1-1.577-.94c-.071-.037-.107-.146-.143-.218-.036-.072 0-.145 0-.217V8.937c0-.289.072-.47.323-.578.465-.254.895-.507 1.397-.796v39.766zM5.731 9.878v35.063c-.573-.326-1.11-.652-1.648-.941a.793.793 0 0 1-.107-.18V11.326c0-.29.107-.47.358-.615.466-.254.896-.543 1.397-.833z"
                    fill-rule="nonzero" fill="#062a85"></path>
          </g>
        </svg>
      </router-link>
    </td>
    <td slot="buttons" class="tdmin">
      <language-drop-down style="margin-top: 3px;"/>
    </td>
  </top-nav>

<!--  <li class="nav-item">-->
<!--    <language-drop-down/>-->
<!--  </li>-->
  <div class="container-fluid document-inner-view bg-light" v-if="status && name">
    <div class="row h-100">
      <div class="col-md-3 mt-3">
        <h5 class="mb-4" :title="name">
          <i class="mdi mdi-file-document material-icons md-24 mr-2"></i>
          {{ name }}
        </h5>
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
              <div class="form form-compiled" v-append="formSrc" @appended="appended"
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
        <div class="d-flex flex-row flex-wrap">
          <pdf-preview v-if="documentPreviews && anyLangAvailable(doc)" v-for="doc in documentPreviews" :key="doc.id"
                       :name="doc.name" :languages="doc.langs" :doc="doc" :wfId="id" :locale="locale"
                       :langSelectorVisible="isConfirmationStep === false"/>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import TopNav from '@/components/layout/TopNav'
import PdfModal from '@/components/document/PdfModal'
import PdfPreview from '@/components/document/PdfPreview'
import ButtonSpinner from '@/components/ButtonSpinner'
import Spinner from '@/components/Spinner'
import LanguageDropDown from '@/views/appDependentComponents/LanguageDropDown'

import formCompilerAdapter from '../libs/formcompiler-adapter'
import TopRightProfile from '../components/user/TopRightProfile'
import mafdc from '@/mixinApp'
import AnimatedInput from '../components/AnimatedInput'

export default {
  mixins: [mafdc, formCompilerAdapter],
  name: 'document-flow',
  components: {
    AnimatedInput,
    TopRightProfile,
    ButtonSpinner,
    Spinner,
    PdfModal,
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
      let response = await axios.get('/api/me')
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
        $('.form-compiled').hide()
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
        $('.form-compiled').hide()
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
            $('.form-compiled').hide()
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
      let templateLanguageConfig = {}
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
    confirmDoc (hash) {
      return new Promise(async (resolve, reject) => {
        let account = this.wallet.getCurrentAddress()

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

        this.nonce = await this.wallet.proxeusFS.web3.eth.getTransactionCount(account)
        this.nonce++
        this.wallet.proxeusFS.createFileUndefinedSigners({
          from: account,
          hash,
          data: (this.data === '' ? '0x00' : web3.fromAscii(this.data)),
          mandatorySigners: 0,
          expiry: 0,
          providers: [],
          nonce: this.nonce,
          xes: 0
        }).then((result) => {
          resolve(result.transactionHash)
        }).catch((error) => {
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
      let self = this
      const $form = $('.form-compiled > form')
      $form.on('dynamicFormScriptExecuted', function () {
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
      })
    },
    getFormData () {
      let excludeFileFormFields = true
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

  /deep/ .navbar .dropdown .nav-link {
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
