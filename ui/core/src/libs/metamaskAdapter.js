export function metamaskGetAccount (web3) {
  return new Promise((resolve, reject) => {
    // In Metamask's current version of web3, getAccounts is not a Promise :(
    web3.eth.getAccounts((err, accounts) => {
      if (err) {
        reject(err)
      }
      if (!(accounts && accounts.length && accounts[0])) {
        reject(new Error('Please login to Metamask'))
      }
      const account = accounts[0]
      resolve({ account })
    })
  })
}

export function metamaskSign ({ challenge, account }, web3) {
  return new Promise((resolve, reject) => {
    const method = 'personal_sign'
    const params = [challenge, account]

    // In Metamask's current version of web3, sendAsync is not a Promise :(
    web3.currentProvider.sendAsync({ method, params, account }, (err, result) => {
      if (err) {
        reject(err)
      }
      const signature = result.result
      if (signature) {
        resolve({ signature })
      }
      reject(new Error('Could not sign.'))
    })
  })
}
