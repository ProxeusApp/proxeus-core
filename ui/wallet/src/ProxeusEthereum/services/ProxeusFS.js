export default class {
  constructor (web3, proxeusFSContract) {
    this.web3 = web3
    this.contract = proxeusFSContract
  }

  /*
   * Returns a list of signers for a specific  file
   */
  async getFileSigners (hash) {
    return this.contract.methods.getFileSigners(hash).call()
  }

  /*
   * File hash verification if registered on the Contract
   */
  async fileVerify (hash) {
    return this.contract.methods.verifyFile(hash).call()
  }

  /**
   * Sign file
   */
  async signFile ({
    hash
  }) {
    return this.contract.methods.signFile(hash).send()
  }

  /**
   * Create a file with undefined signers
   */
  async createFileUndefinedSigners ({
    hash,
    data
  }) {
    return this.contract.methods.registerFile(
      hash,
      data
    ).send()
  }

  /**
   * Estimate gas for freate a file with undefined signers
   */
  async createFileUndefinedSignersEstimateGas ({
    from,
    hash
  }, cb, xes) {
    const opt = from ? {
      from
    } : {}
    return this.contract.methods.registerFile(
      hash
    ).estimateGas(opt, cb)
  }
}
