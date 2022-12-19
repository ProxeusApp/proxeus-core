import WalletInterface from './libs/WalletAdapter'

export default {
  data () {
    return {
      me: null,
      meta: null,
      translations: null,
      blockchainNet: 'goerli',
      blockchainProxeusFSAddress: '',
      roles: [],
      uiBlocked: false,
      intervalID: null,
      lastExportResults: null,
      lastImportResults: null
    }
  },
  methods: {
    makeURL (uri) {
      let origin = location.origin
      if (/.*\/$/.test(origin)) {
        origin = origin.substring(0, origin.length - 1)
      }
      if (/^\/.*/.test(uri) === false) {
        uri = '/' + uri
      }
      return origin + uri
    },
    isLangAvailable (lang) {
      if (this.meta && this.meta.activeLangs) {
        for (let i = 0; i < this.meta.activeLangs.length; i++) {
          if (this.meta.activeLangs[i] &&
            this.meta.activeLangs[i].Code === lang) {
            return true
          }
        }
      }
      return false
    },
    userIsRoot () {
      if (this.me) {
        return this.me.role === 100
      }
      return false
    },
    userIsSuperAdmin () {
      if (this.me) {
        return this.me.role > 10
      }
      return false
    },
    userIsAdminOrHigher () {
      if (this.me) {
        return this.me.role >= 10
      }
      return false
    },
    userIsCreatorOrHigher () {
      if (this.me) {
        return this.me.role >= 7
      }
      return false
    },
    userIsUserOrHigher () {
      if (this.me) {
        return this.me.role >= 5
      }
      return false
    },
    amIWriteGrantedFor (item) {
      try {
        // check authority
        if (this.me.role >= 100) {
          return true
        }
      } catch (e) {
      }
      try {
        // check public
        if (item.publicByID[0] === 2) {
          return true
        }
      } catch (e) {
      }

      try {
        // check owner
        if (this.me.id === item.owner) {
          return true
        }
        if (item.grant && item.grant[this.me.id] &&
          item.grant[this.me.id][0] === 2) {
          return true
        }
      } catch (e) {
      }
      try {
        // check group
        if (this.me.role !== 0) {
          if (item.group <= this.me.role && item.groupAndOthers.rights[0] === 2) {
            // same or higher group has write rights
            return true
          }
        }
      } catch (e) {
      }

      // check others
      try {
        if (this.me.role >= 1 && item.groupAndOthers.rights[1] === 2) {
          // others have write rights
          return true
        }
      } catch (e) {
      }
      return false
    },
    handleError (o) {
      if (this.isConToServerLostError(o)) {
        // couldn't reach the server
        if (!this.intervalID) {
          this.$root.$emit('service-off')
          this.intervalID = setInterval(this.pingService, 2000)
        }
      }
    },
    isConToServerLostError (o) {
      let txt = ''
      try {
        txt = o.data || o.request.response || o.request.responseText
      } catch (e) {}
      return (o.request && !o.response) ||
        (window.webpackHotUpdate && /^.*\(ECONNREFUSED\)\.$/m.test(txt))
    },
    pingService () {
      console.log('pingService')
      axios.get('/api/config').then((r) => {
        if (r.data) {
          this.setConfig(r.data)
        }
        if (this.intervalID) {
          clearInterval(this.intervalID)
          this.intervalID = null
          this.$root.$emit('service-on')
        }
      }, (err) => {
        this.handleError(err)
      })
    },
    getSelectedLang () {
      return this.$cookie.get('lang') || this.fallbackLang()
    },
    setSelectedLang (lang) {
      if (lang) {
        this.$cookie.set('lang', lang, { expires: '1Y' })
        this.reloadI18n()
      } else {
        this.$cookie.delete('lang')
        this.$i18n.set(this.fallbackLang())
      }
    },
    loadMe (clb) {
      axios.get('/api/me').then((response) => {
        this.me = response.data
        this.$root.$emit('me', this.me)
        if (clb) {
          try {
            clb(this.me)
          } catch (e) {
            console.log(e)
          }
        }
      }, (err) => {
        this.handleError(err)
      })
    },
    loadLastExportResults (clb, delParams) {
      let url = '/api/export/results'
      if (delParams) {
        url += '?' + delParams
      }
      axios.get(url).then((response) => {
        this.lastExportResults = response.data
        this.$root.$emit('lastExportResults', this.lastExportResults)
        if (clb) {
          try {
            clb(this.lastExportResults)
          } catch (e) {
            console.log(e)
          }
        }
      }, (err) => {
        this.handleError(err)
      })
    },
    loadLastImportResults (clb, delParams) {
      let url = '/api/import/results'
      if (delParams) {
        url += '?' + delParams
      }
      axios.get(url).then((response) => {
        this.lastImportResults = response.data
        this.$root.$emit('lastImportResults', this.lastImportResults)
        if (clb) {
          try {
            clb(this.lastImportResults)
          } catch (e) {
            console.log(e)
          }
        }
      }, (err) => {
        this.handleError(err)
      })
    },
    exportData (params, cb, url, name) {
      if (!url) {
        url = '/api/export?include=' + params
      } else {
        url = url + '?include=' + params
      }
      this.blockUI()
      axios({
        url: url,
        method: 'GET',
        responseType: 'blob',
        timeout: 0
      }, {
        timeout: 0
      }).then((response) => {
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        if (!name) {
          name = 'Proxeus'
        } else {
          name = 'Proxeus-' + name
        }
        const now = new Date()
        name += '_' + now.getDate() + '-' + (now.getMonth() + 1) + '-' +
          now.getFullYear() + '_' + now.getHours() + '-' + now.getMinutes()
        link.setAttribute('download', name + '.db')
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        this.unblockUI()
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('The data was successfully exported.'),
          type: 'success'
        })
        this.loadLastExportResults(cb)
      }, (err) => {
        this.handleError(err)
        this.unblockUI()
        let text = this.$t('There was an error exporting the data. Please try again or if the error persists contact the platform operator.')
        try {
          if (err.response && typeof err.response.data === 'string') {
            text = err.response.data
          }
        } catch (e) {}
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: text,
          type: 'error'
        })
      })
    },
    importData (file, skipExisting, cb) {
      if (!file) {
        return
      }
      this.blockUI()
      if (!skipExisting) {
        skipExisting = false
      }
      axios.post('/api/import?skipExisting=' + skipExisting, file, {
        headers: {
          'File-Name': encodeURI(file.name),
          'Content-Type': file.type
        },
        timeout: 0
      }).then(response => {
        this.unblockUI()
        this.$notify({
          group: 'app',
          title: this.$t('Success'),
          text: this.$t('Import was successful.'),
          type: 'success'
        })
        this.loadLastImportResults(cb)
      }, (err) => {
        this.unblockUI()
        let text = this.$t('There was an error while importing the data. Please try again or if the error persists contact the platform operator.')
        try {
          text = err.response.data
        } catch (e) {}
        this.$notify({
          group: 'app',
          title: this.$t('Error'),
          text: text,
          type: 'error'
        })
        this.handleError(err)
      })
    },
    loadMeta (clb) {
      axios.get('/api/i18n/meta').then((response) => {
        this.meta = response.data
        if (this.meta && this.meta.langFallback) {
          this.$i18n.fallback(this.meta.langFallback)
          this.$i18n.add(this.meta.langFallback, this.meta.fallbackTranslations)
        }
        if (clb) {
          try {
            clb(this.meta)
          } catch (e) {
            console.log(e)
          }
        }
        this.reloadI18n()
      }, (err) => {
        this.handleError(err)
      })
    },
    updateLangLabel () {
      if (this.meta && this.meta.activeLangs) {
        for (let i = 0; i < this.meta.activeLangs.length; i++) {
          this.meta.activeLangs[i].label = this.$t(
            this.meta.activeLangs[i].Code)
        }
        return this.meta.activeLangs
      }
    },
    fallbackLang () {
      if (this.meta) {
        return this.meta.langFallback
      }
      return null
    },
    getSelectedLangIndex () {
      if (this.meta && this.meta.activeLangs) {
        for (let i = 0; i < this.meta.activeLangs.length; i++) {
          if (this.meta.activeLangs[i].Code === this.getSelectedLang()) {
            return i
          }
        }
      }
      return null
    },
    reloadI18n () {
      axios.get('/api/i18n/all').then(
        (response) => {
          this.$i18n.add(this.getSelectedLang(), response.data)
          this.$i18n.set(this.getSelectedLang())
          this.updateLangLabel()
          this.$root.$emit('translations-updated')
        }, (err) => {
          this.handleError(err)
        })
    },
    loadConfig () {
      axios.get('/api/config').then(r => {
        if (r.data) {
          this.setConfig(r.data)
        }
      }, (err) => {
        this.handleError(err)
      })
    },
    setConfig (d) {
      if (d.blockchainNet) {
        this.blockchainNet = d.blockchainNet
      }
      if (d.blockchainProxeusFSAddress) {
        this.blockchainProxeusFSAddress = d.blockchainProxeusFSAddress
      }
      if (d.roles) {
        this.roles = d.roles
      }
      if (this.blockchainNet && this.blockchainProxeusFSAddress) {
        this.wallet = new WalletInterface(this.blockchainNet, this.blockchainProxeusFSAddress)
      }
    },
    acknowledgeFirstLogin () {
      // Show the following overlays starting now (they can be delayed)
      localStorage.setItem('showFirstLoginMessageOn-documents', new Date())
      if (this.userIsCreatorOrHigher()) {
        localStorage.setItem('showFirstLoginMessageOn-admin', new Date())
      }
    }
  },
  computed: {
    app: {
      get () {
        return this.$root.$children[0]
      },
      set (a) {
      }
    }
  },
  created () {
    const tmpLangToPreventFromWarnings = 'en'
    this.$i18n.fallback(tmpLangToPreventFromWarnings)
    this.$i18n.set(tmpLangToPreventFromWarnings)
    this.loadMeta()
    this.loadConfig()
    this.loadMe()

    // when accounts loaded register accountsChanged handler and reload page if user changes the account
    window.ethereum.on('accountsChanged', function () {
      window.location.reload()
    })
    window.ethereum.on('chainChanged', function (e) {
      console.log('Chain changed')
      window.location.reload()
    })
  }
}
