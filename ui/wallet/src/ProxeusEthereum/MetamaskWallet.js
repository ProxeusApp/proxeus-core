class MetamaskWallet {
  constructor (web3, xesTokenContract) {
    this.web3 = web3
    this.xesTokenContract = xesTokenContract
  }

  async setupDefaultAccount () {
    this.web3.eth.defaultAccount = (await this.web3.eth.getAccounts())[0]
  }

  async signMessage (message) {
    let address = this.getCurrentAddress()
    return this.web3.eth.personal.sign(message, address)
  }

  async transferXES (to, amount) {
    return this.xesTokenContract.methods.transfer(to, amount)
      .send({ from: this.getCurrentAddress() })
  }

  getCurrentAddress () {
    this.web3.eth.getAccounts().then((address) => {
      if (this.web3.eth.defaultAccount !== address[0]) {
        console.log('user changed the account')
      }
    })
    return this.web3.eth.defaultAccount
  }
}

export { MetamaskWallet as default }
