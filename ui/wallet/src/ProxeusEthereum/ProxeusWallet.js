class ProxeusWallet {
  constructor (web3, xesTokenContract) {
    this.web3 = web3
    this.xesTokenContract = xesTokenContract
  }

  loadKeystore (keystore, password) {
    // will be needed in the future for storing and loading into memory
    this.globalPassword = password

    this.web3.eth.accounts.wallet.decrypt(keystore, password)
  }

  async signMessage (message) {
    return this.web3.eth.accounts.sign(message,
      this.web3.eth.accounts.wallet[this.web3.eth.defaultAccount].privateKey).signature
  }

  async transferXES (to, amount) {
    return this.xesTokenContract.methods.transfer(to, amount)
      .send({ from: this.web3.eth.defaultAccount })
  }

  getCurrentAddress () {
    if ((this.web3.eth.defaultAccount === undefined ||
      this.web3.eth.defaultAccount === '0x0000000000000000000000000000000000000000') &&
      this.web3.eth.accounts.wallet[0]) {
      this.web3.eth.defaultAccount = this.web3.eth.accounts.wallet[0].address
    }
    return this.web3.eth.defaultAccount
  }

  importPrivateKey (privateKey) {
    this.web3.eth.accounts.wallet.add(privateKey)
    // set as default account
    this.web3.eth.defaultAccount = this.web3.eth.accounts.wallet[0].address
  }

  exportPrivateKey () {
    return this.web3.eth.accounts.wallet[0].privateKey
  }

  importKeystore (keystore, password) {
    try {
      this.web3.eth.accounts.wallet.decrypt(keystore, password)

      // set as default account if anything came from the keystore
      if (this.web3.eth.accounts.wallet[0] !== null) {
        this.web3.eth.defaultAccount = this.web3.eth.accounts.wallet[0].address
      }
      return true
    } catch (e) {
      return false
    }
  }

  exportKeystore (password) {
    return this.web3.eth.accounts.wallet.encrypt(password)
  }
}

export { ProxeusWallet as default }
