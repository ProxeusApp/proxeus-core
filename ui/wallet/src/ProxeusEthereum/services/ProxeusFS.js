export default class {
  constructor (web3, proxeusFSContract) {
    this.web3 = web3
    this.contract = proxeusFSContract
  }

  async XESAmountPerFile ({ providers }) {
    return this.contract.methods.getXESPrice().call()
  }

  /**
   ** Free calls
   **/

  /*
   * Returns expiry information of a file
   */
  fileExpiry (hash) {
  }

  fileGetPerm (hash, address, write) {
    return this.contract.methods.fileGetPerm(hash, address, write).call()
  }

  fileHasSP (hash, spAddress) {
    return this.contract.methods.fileHasSP(hash, spAddress).call()
  }

  /*
   * Returns file metadata
   */
  fileInfo (hash) {
    return this.contract.methods.fileInfo(hash).call()
  }

  /*
   * Returns a list of files
   */
  async fileList () {
    return this.contract.methods.fileList().call()
  }

  /*
   * Returns a list of signers for a specific  file
   */
  getFileSigners (hash) {
    return this.contract.methods.getFileSigners(hash).call()
  }

  /*
   * File hash verification if registered on the Contract
   */
  async fileVerify (hash) {
    return this.contract.methods.verifyFile(hash).call()
  }

  /*
   * Returns metadata about a Service Provider
   */
  spInfo (hash) {
    return this.contract.methods.spInfo(hash).call()
  }

  /*
   * Returns a list of storage providers
   */
  spList () {
    return this.contract.methods.spList().call()
  }

  /*
   * Paid calls
   */

  fileRemove (hash) {
    return this.contract.methods.fileRemove(hash).send()
  }

  signFile ({ hash }) {
    return this.contract.methods.signFile(hash).send()
  }

  fileSetPerm (hash, address) {
    return this.contract.methods.fileSetPerm(hash, address).send()
  }

  createFileDefinedSigners ({ hash, definedSigners, expiry, replaces, providers }) {
    return this.contract.methods.createFileDefinedSigners(
      hash,
      definedSigners,
      expiry / 1000 || 0,
      this.web3.utils.fromAscii('0x0'),
      providers
    ).send()
  }

  createFileUndefinedSigners ({ hash, data, mandatorySigners, expiry, replaces, providers, xes }) {
    return this.contract.methods.registerFile(
      hash,
      data
    ).send()
  }

  createFileUndefinedSignersEstimateGas (
    { from, hash, data, mandatorySigners, expiry, replaces, providers }, cb, xes) {
    const opt = from ? { from } : {}
    return this.contract.methods.registerFile(
      hash
    ).estimateGas(opt, cb)
  }

  createFileThumbnail (hash, pParent, pPublic) {
    return this.contract.methods.createFileThumbnail(
      hash,
      pParent,
      pPublic
    ).send()
  }

  fileRequestSign (hash, address) {
    return this.contract.methods.fileRequestSign(hash, address).send()
  }

  fileRequestAccess () {
  }

  getRequestSignEvents (address, fromBlock) {
    return this.contract.getPastEvents('RequestSign', {
      filter: { to: address },
      fromBlock: fromBlock || 0
    })
  }

  getNotifySignEvents (hash, address, fromBlock) {
    return this.contract.getPastEvents('NotifySign', {
      filter: { hash: hash, who: address },
      fromBlock: fromBlock || 0
    })
  }
}
