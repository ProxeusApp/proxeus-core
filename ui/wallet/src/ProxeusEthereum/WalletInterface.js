import serviceConfig from './config/service-config'
import { PROXEUS_FS_ABI, XES_TOKEN_ABI } from './config/ABI'

import ProxeusWallet from './ProxeusWallet'
import ProxeusFS from './services/ProxeusFS'
import MetamaskWallet from './MetamaskWallet'
import { keccak256 } from 'js-sha3'
import MetamaskUtil from './MetamaskUtil'

import Web3 from 'web3'

import getTransactionReceiptMined from './helpers/getTransactionReceiptMined'

class WalletInterface {
  // TODO improve checking that current network matches what is expected
  // TODO: network param only for compatibility reasons with blockchain/dapp
  constructor (network = 'goerli', proxeusFSAddress, forceProxeusWallet = false) {
    this.useProxeusWallet = forceProxeusWallet || typeof window.ethereum === 'undefined'

    // make sure we are using the web3 we want and not the one provided by metamask
    this.web3 = new Web3(Web3.givenProvider || 'ws://localhost:8545')

    this.web3.eth.getTransactionReceiptMined = getTransactionReceiptMined
    this.serviceConfig = serviceConfig[network]
    if (proxeusFSAddress) {
      this.serviceConfig.PROXEUS_FS_ADDRESS = proxeusFSAddress
    }
    this.metamaskUtil = new MetamaskUtil()

    if (this.useProxeusWallet) {
      // connect to the network using what was given in the constructor
      this.web3.setProvider(
        new this.web3.providers.HttpProvider(
          'https://' + network + '.infura.io/'))
    } else {
      if (window.ethereum) {
        this.web3.setProvider(window.ethereum)
      }
    }

    // contracts must be set after provider is defined,
    // otherwise they will have an invalid provider
    // add the XES smart contract to the config
    this.xesTokenContract = new this.web3.eth.Contract(
      XES_TOKEN_ABI,
      this.serviceConfig.XES_TOKEN_ADDRESS,
      { gas: this.serviceConfig.DEFAULT_GAS_REGULAR }
    )
    this.setProxeusFsContract(this.serviceConfig.PROXEUS_FS_ADDRESS)

    if (this.useProxeusWallet) {
      this.wallet = new ProxeusWallet(this.web3, this.xesTokenContract)
    } else {
      this.wallet = new MetamaskWallet(this.web3, this.xesTokenContract)

      // set the default from address to use on the proxeusFS smart contract
      this.wallet.setupDefaultAccount().then(() => {
        if (this.wallet.getCurrentAddress() !== null) { this.proxeusFSContract.options.from = this.wallet.getCurrentAddress() }
      })
    }
  }

  signMessage (message) {
    return this.wallet.signMessage(message)
  }

  getCurrentAddress () {
    return this.wallet.getCurrentAddress()
  }

  hasAccount () {
    return this.wallet.getCurrentAddress() !== undefined
  }

  importKeystore (keystore, password) {
    // optional function on wallet
    if (!this.wallet.importKeystore) {
      return null
    }

    return this.wallet.importKeystore(keystore, password)
  }

  exportKeystore (password) {
    // optional function on wallet
    if (!this.wallet.exportKeystore) {
      return null
    }

    return this.wallet.exportKeystore(password)
  }

  importPrivateKey (privateKey) {
    // optional function on wallet
    if (!this.wallet.importPrivateKey) {
      return null
    }

    return this.wallet.importPrivateKey(privateKey)
  }

  exportPrivateKey () {
    // optional function on wallet
    if (!this.wallet.exportPrivateKey) {
      return null
    }

    return this.wallet.exportPrivateKey()
  }

  storePGPPublicKey (pgpPublicKey) {
    pgpPublicKey = btoa(pgpPublicKey)
    // TODO: Store keys in constants
    localStorage.setItem('pgpPk', pgpPublicKey)
  }

  loadPGPPublicKey () {
    return atob(localStorage.getItem('pgpPk'))
  }

  exportWalletToBlob (password = '') {
    if (password === '') {
      try {
        const authObj = localStorage.getItem('mnidmao')
        password = JSON.parse(atob(authObj)).password
      } catch (e) {
        return false
      }
    }

    const encryptedKeystore = this.exportKeystore(password)
    if (encryptedKeystore.length === 0) {
      return false
    }

    return new Blob([
      btoa(JSON.stringify({
        keystore: encryptedKeystore,
        pgpKeys: this.web3.eth.accounts.wallet.PGPKeys
      }))], { type: 'text/plain' })
  }

  importWalletFromBlob (blob, password) {
    const reader = new FileReader()
    // first set the reader event listener and wrap it in a promise
    const promise = new Promise((resolve, reject) => {
      reader.addEventListener('loadend', () => {
        try {
          const decoder = new TextDecoder()
          const parsed = JSON.parse(atob(decoder.decode(reader.result)))

          const importedKeystore = this.importKeystore(parsed.keystore, password)

          if (importedKeystore) {
            this.web3.eth.accounts.wallet.PGPKeys = parsed.pgpKeys

            this.proxeusFSContract.options.from = this.getCurrentAddress()

            // save the pgp keys, encrypted keystore, and password on local storage
            localStorage.setItem('mnidmao', btoa(JSON.stringify({ password })))
            localStorage.setItem('mnidmpgp',
              btoa(JSON.stringify(parsed.pgpKeys[this.getCurrentAddress()])))
            localStorage.setItem('mnidmks',
              btoa(JSON.stringify(parsed.keystore)))

            resolve(true)
          }

          reject(
            new Error('wallet.importWalletFromBlob.error.incorrectPassword'))
        } catch (error) {
          reject(new Error('wallet.importWalletFromBlob.error.invalidFile'))
        }
      })
    })

    // call the reader on the blob which will trigger the above event listener
    reader.readAsArrayBuffer(blob)

    return promise
  }

  /*
   * Saves into local storage the list of private keys existing in the
   * wallet as encrypted keystores.
   * The last element of the list is the array of PGP keys of the wallet.
   *
   * @param password - the password that protects the keystore
   */
  saveWallet (password) {
    let encryptedKeystore = this.exportKeystore(password)
    if (encryptedKeystore.length === 0) {
      return false
    }

    let pgp = this.exportPGPKeyPair(this.getCurrentAddress())

    let authObj = {
      password
    }

    authObj = btoa(JSON.stringify(authObj))
    encryptedKeystore = btoa(JSON.stringify(encryptedKeystore))
    pgp = btoa(JSON.stringify(pgp))

    localStorage.setItem('mnidmao', authObj)
    localStorage.setItem('mnidmpgp', pgp)
    localStorage.setItem('mnidmks', encryptedKeystore)

    return true
  }

  /*
   * Reads a list of keystores from local storage and loads it into the wallet.
   * The list should have as the last element an array of PGP keys.
   *
   * @param password - the password that protects the keystore
   */
  loadWallet (password = '') {
    const authObj = localStorage.getItem('mnidmao')
    let pgp = localStorage.getItem('mnidmpgp')
    let encryptedKeystore = this.getKeystoreFromLocalStorage()

    if (password === '') {
      try {
        password = JSON.parse(atob(authObj)).password
      } catch (e) {
        return false
      }
    } else {
      let authObj = {
        password
      }

      authObj = btoa(JSON.stringify(authObj))
      localStorage.setItem('mnidmao', authObj)
    }

    encryptedKeystore = JSON.parse(atob(encryptedKeystore))
    if (this.importKeystore(encryptedKeystore, password) !== true) {
      return false
    }

    pgp = JSON.parse(atob(pgp))
    this.importPGPKeyPair(this.getCurrentAddress(), pgp.publicKey,
      pgp.privateKey)

    this.proxeusFSContract.options.from = this.getCurrentAddress()

    return true
  }

  lockAccount () {
    localStorage.removeItem('mnidmao')
  }

  logout () {
    localStorage.removeItem('mnidmao')
    localStorage.removeItem('mnidmks')
    localStorage.removeItem('mnidmpgp')
    this.web3.eth.accounts.wallet.PGPKeys = {}
    // can't be undefined because Eth has a watched on the property and won't let it go undefined
    this.web3.eth.defaultAccount = '0x0000000000000000000000000000000000000000'
    this.web3.eth.accounts.wallet.clear()
  }

  getKeystoreFromLocalStorage () {
    return localStorage.getItem('mnidmks')
  }

  getPasswordFromLocalStorage () {
    try {
      const authObj = localStorage.getItem('mnidmao')
      return JSON.parse(atob(authObj)).password
    } catch (e) {
      return false
    }
  }

  /*
   * Imports a PGP key pair into the wallet.
   *
   * @param address - the ethereum wallet address to which the PGP key pair relates to
   * @param publicKey
   * @param privateKey
   */
  importPGPKeyPair (address, publicKey, privateKey) {
    if (!this.web3.eth.accounts.wallet.PGPKeys) this.web3.eth.accounts.wallet.PGPKeys = {}
    this.web3.eth.accounts.wallet.PGPKeys[address] = { publicKey, privateKey }
  }

  /*
   * Returns the PGP key pair assigned to the requested address.
   *
   * @param address - the ethereum wallet address to which the PGP key pair relates to
   * @return Object with two fields: "privateKey" and "publicKey"
   */
  exportPGPKeyPair (address) {
    return this.web3.eth.accounts.wallet.PGPKeys[address]
  }

  createNewAccount () {
    return this.web3.eth.accounts.wallet.create(1)
  }

  /*
   * Everything below this comment needs to be moved outside of the wallet libray
   */

  hashFile (arrBuffer) {
    return keccak256(arrBuffer)
  }

  getDocumentRegistrationTx (hash, proxeusFSContract) {
    const contract = (proxeusFSContract === undefined) ? this.proxeusFSContract : proxeusFSContract.contract
    // this one is based on events and not promises, so can't use async
    return new Promise((resolve, reject) => {
      contract.getPastEvents('UpdatedEvent', {
        filter: { hash: hash },
        fromBlock: 0
      }, (error, result) => {
        if (error) {
          reject(error)
        }
        if (result === undefined || result.length === 0) {
          reject(new Error('No event found'))
        } else {
          resolve(result[0].transactionHash)
        }
      })
    })
  }

  setProxeusFsContract (address) {
    // add the document registry smart contract to the config
    this.proxeusFSContract = new this.web3.eth.Contract(
      PROXEUS_FS_ABI,
      address,
      { gas: this.serviceConfig.DEFAULT_GAS_REGULAR }
    )

    // Attach proxeus FS Service
    this.proxeusFS = new ProxeusFS(this.web3, this.proxeusFSContract)
  }

  /**
   * Returns and array of ProxeusFS contract instances which refer to past contract addresses
   *
   * @return array of past contract instances
   */
  getAllProxeusFsServices () {
    const proxeusFSPastContracts = []
    for (const address of this.serviceConfig.PROXEUS_FS_PAST_ADDRESSES) {
      const proxeusFSContract = new this.web3.eth.Contract(
        PROXEUS_FS_ABI,
        address,
        { gas: this.serviceConfig.DEFAULT_GAS_REGULAR }
      )
      proxeusFSPastContracts.push(new ProxeusFS(this.web3, proxeusFSContract))
    }
    return proxeusFSPastContracts
  }

  /**
   * Returns the network provided by the client's browser wallet
   *
   * @return string
   */
  async getClientProvidedNetwork () {
    const netId = await this.web3.eth.getChainId()
    switch (netId) {
      case 5:
        return 'goerli'
      case 11155111:
        return 'sepolia'
      case 137:
        return 'polygon-mainnet'
      case 80001:
        return 'polygon-mumbai'
      default:
        return 'mainnet'
    }
  }

  async XESAmountPerFile ({ providers }) {
    const tokensRaw = await this.proxeusFS.XESAmountPerFile({ providers })
    return this.metamaskUtil.formatBalance(this.web3.utils.toHex(tokensRaw))
  }

  async verifyHash (hash) {
    const result = await this.proxeusFS.fileVerify(hash)

    if (result && result[0] === true) {
      return this.getDocumentRegistrationTx(hash)
    } else {
      for (const proxeusFSPastContract of this.getAllProxeusFsServices()) {
        const result = await proxeusFSPastContract.fileVerify(hash)
        if (result && result[0] === true) {
          return this.getDocumentRegistrationTx(hash, proxeusFSPastContract)
        }
      }
    }
    throw new Error('Could not verify hash.')
  }

  async fileVerify (hash) {
    const result = await this.proxeusFS.fileVerify(hash)

    if (result && result[0] === true) {
      return result
    } else {
      for (const proxeusFSPastContract of this.getAllProxeusFsServices()) {
        const result = await proxeusFSPastContract.fileVerify(hash)
        if (result && result[0] === true) {
          return result
        }
      }
    }
    throw new Error('Could not verify hash.')
  }

  formatBalance (decimalsToKeep, tokensRaw) {
    if (decimalsToKeep !== undefined) {
      return this.metamaskUtil.formatBalance(this.web3.utils.toHex(tokensRaw),
        decimalsToKeep)
    } else {
      return this.web3.utils.fromWei(tokensRaw, 'ether')
    }
  }

  transferETH (to, amount) {
    return this.web3.eth.sendTransaction({
      to: to,
      value: amount,
      gas: this.serviceConfig.DEFAULT_GAS_REGULAR
    })
  }

  async getETHBalance (decimalsToKeep) {
    return this.formatBalance(decimalsToKeep, await this.web3.eth.getBalance(this.getCurrentAddress()))
  }

  // optional callback parameter
  transferXES (to, amount, callback) {
    return this.wallet.transferXES(to, amount, callback)
  }

  async getXESBalance (decimalsToKeep, address) {
    if (address === undefined) {
      address = this.getCurrentAddress()
    }
    const tokensRaw = await this.xesTokenContract.methods.balanceOf(address)
      .call()

    return this.formatBalance(decimalsToKeep, tokensRaw)
  }

  approveXES (spender, value) {
    return this.xesTokenContract.methods.approve(spender,
      this.web3.utils.toWei(value.toString()))
      .send({ from: this.getCurrentAddress() })
  }

  async getAllowance ({ spender, decimalsToKeep }) {
    const owner = this.getCurrentAddress()
    const tokensRaw = await this.xesTokenContract.methods.allowance(owner,
      spender).call()

    if (decimalsToKeep === undefined) {
      return this.web3.utils.fromWei(tokensRaw)
    } else {
      return this.metamaskUtil.formatBalance(this.web3.utils.toHex(tokensRaw),
        decimalsToKeep)
    }
  }

  txMined (tx) {
    return this.web3.eth.getTransactionReceiptMined(tx)
  }

  async getFileSignedEvent (signerAddress) {
    const arrEvents = await this.proxeusFSContract.getPastEvents(
      'FileSignedEvent', {
        filter: { signer: signerAddress },
        fromBlock: 0
      })
    return arrEvents[0] || null
  }

  async getRegistrationTxBlock (signerAddress) {
    const event = await this.getFileSignedEvent(signerAddress)
    if (event === null) {
      return null
    }
    const block = await this.getBlock(event.blockHash)
    return {
      txHash: event.transactionHash,
      block: block
    }
  }

  async getBlock (blockHash) {
    return await this.web3.eth.getBlock(blockHash)
  }
}

export default WalletInterface
