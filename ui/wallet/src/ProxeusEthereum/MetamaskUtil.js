class MetamaskUtil {
  constructor () {
    this.abi = require('human-standard-token-abi')
    this.ethUtil = require('ethereumjs-util')
    this.vreme = new (require('vreme'))()

    this.MIN_GAS_PRICE_GWEI_BN = new this.ethUtil.BN(1)
    this.GWEI_FACTOR = new this.ethUtil.BN(1e9)
    this.MIN_GAS_PRICE_BN = this.MIN_GAS_PRICE_GWEI_BN.mul(this.GWEI_FACTOR)

    this.valueTable = {
      wei: '1000000000000000000',
      kwei: '1000000000000000',
      mwei: '1000000000000',
      gwei: '1000000000',
      szabo: '1000000',
      finney: '1000',
      ether: 1,
      kether: 0.001,
      mether: 0.000001,
      gether: 0.000000001,
      tether: 0.000000000001
    }
    this.multiple = new this.ethUtil.BN('10000', 10)

    this.bnTable = {}
    for (const currency in this.valueTable) {
      this.bnTable[currency] = new this.ethUtil.BN(this.valueTable[currency],
        10)
    }
  }

  /*
   * formatData :: ( date: <Unix Timestamp> ) -> String
   */
  formatDate (date) {
    return this.vreme.format(new Date(date), 'March 16 2014 14:30')
  }

  valuesFor (obj) {
    if (!obj) return []
    return Object.keys(obj)
      .map((key) => { return obj[key] })
  }

  addressSummary (
    address, firstSegLength = 10, lastSegLength = 4, includeHex = true) {
    if (!address) return ''
    let checked = this.checksumAddress(address)
    if (!includeHex) {
      checked = this.ethUtil.stripHexPrefix(checked)
    }
    return checked ? checked.slice(0, firstSegLength) + '...' +
      checked.slice(checked.length - lastSegLength) : '...'
  }

  miniAddressSummary (address) {
    if (!address) return ''
    var checked = this.checksumAddress(address)
    return checked ? checked.slice(0, 4) + '...' + checked.slice(-4) : '...'
  }

  isValidAddress (address) {
    var prefixed = this.ethUtil.addHexPrefix(address)
    if (address === '0x0000000000000000000000000000000000000000') return false
    return (this.isAllOneCase(prefixed) &&
      this.ethUtil.isValidAddress(prefixed)) ||
      this.ethUtil.isValidChecksumAddress(prefixed)
  }

  isInvalidChecksumAddress (address) {
    var prefixed = this.ethUtil.addHexPrefix(address)
    if (address === '0x0000000000000000000000000000000000000000') return false
    return !this.isAllOneCase(prefixed) &&
      !this.ethUtil.isValidChecksumAddress(prefixed) &&
      this.ethUtil.isValidAddress(prefixed)
  }

  isAllOneCase (address) {
    if (!address) return true
    var lower = address.toLowerCase()
    var upper = address.toUpperCase()
    return address === lower || address === upper
  }

  // Takes wei Hex, returns wei BN, even if input is null
  numericBalance (balance) {
    if (!balance) return new this.ethUtil.BN(0, 16)
    var stripped = this.ethUtil.stripHexPrefix(balance)
    return new this.ethUtil.BN(stripped, 16)
  }

  // Takes  hex, returns [beforeDecimal, afterDecimal]
  parseBalance (balance) {
    var beforeDecimal, afterDecimal
    const wei = this.numericBalance(balance)
    var weiString = wei.toString()
    const trailingZeros = /0+$/

    beforeDecimal = weiString.length > 18 ? weiString.slice(
      0, weiString.length - 18) : '0'
    afterDecimal = ('000000000000000000' + wei).slice(-18)
      .replace(trailingZeros, '')
    if (afterDecimal === '') { afterDecimal = '0' }
    return [beforeDecimal, afterDecimal]
  }

  // Takes wei hex, returns an object with three properties.
  // Its "formatted" property is what we generally use to render values.
  formatBalance (balance, decimalsToKeep, needsParse = true) {
    var parsed = needsParse ? this.parseBalance(balance) : balance.split('.')
    var beforeDecimal = parsed[0]
    var afterDecimal = parsed[1]
    var formatted = 'None'
    if (decimalsToKeep === undefined) {
      if (beforeDecimal === '0') {
        if (afterDecimal !== '0') {
          var sigFigs = afterDecimal.match(/^0*(.{2})/) // default: grabs 2 most significant digits
          if (sigFigs) { afterDecimal = sigFigs[0] }
          formatted = '0.' + afterDecimal
        }
      } else {
        formatted = beforeDecimal + '.' + afterDecimal.slice(0, 3)
      }
    } else {
      afterDecimal += Array(decimalsToKeep).join('0')
      formatted = beforeDecimal + '.' + afterDecimal.slice(0, decimalsToKeep)
    }
    return formatted
  }

  generateBalanceObject (formattedBalance, decimalsToKeep = 1) {
    var balance = formattedBalance.split(' ')[0]
    var label = formattedBalance.split(' ')[1]
    var beforeDecimal = balance.split('.')[0]
    var afterDecimal = balance.split('.')[1]
    var shortBalance = this.shortenBalance(balance, decimalsToKeep)

    if (beforeDecimal === '0' && afterDecimal.substr(0, 5) === '00000') {
      // eslint-disable-next-line eqeqeq
      if (afterDecimal == 0) {
        balance = '0'
      } else {
        balance = '<1.0e-5'
      }
    } else if (beforeDecimal !== '0') {
      balance = `${beforeDecimal}.${afterDecimal.slice(0, decimalsToKeep)}`
    }

    return { balance, label, shortBalance }
  }

  shortenBalance (balance, decimalsToKeep = 1) {
    var truncatedValue
    var convertedBalance = parseFloat(balance)
    if (convertedBalance > 1000000) {
      truncatedValue = (balance / 1000000).toFixed(decimalsToKeep)
      return `${truncatedValue}m`
    } else if (convertedBalance > 1000) {
      truncatedValue = (balance / 1000).toFixed(decimalsToKeep)
      return `${truncatedValue}k`
    } else if (convertedBalance === 0) {
      return '0'
    } else if (convertedBalance < 0.001) {
      return '<0.001'
    } else if (convertedBalance < 1) {
      var stringBalance = convertedBalance.toString()
      if (stringBalance.split('.')[1].length > 3) {
        return convertedBalance.toFixed(3)
      } else {
        return stringBalance
      }
    } else {
      return convertedBalance.toFixed(decimalsToKeep)
    }
  }

  dataSize (data) {
    var size = data ? this.ethUtil.stripHexPrefix(data).length : 0
    return size + ' bytes'
  }

  // Takes a BN and an ethereum currency name,
  // returns a BN in wei
  normalizeToWei (amount, currency) {
    try {
      return amount.mul(this.bnTable.wei).div(this.bnTable[currency])
    } catch (e) {}
    return amount
  }

  normalizeEthStringToWei (str) {
    const parts = str.split('.')
    let eth = new this.ethUtil.BN(parts[0], 10).mul(this.bnTable.wei)
    if (parts[1]) {
      var decimal = parts[1]
      while (decimal.length < 18) {
        decimal += '0'
      }
      if (decimal.length > 18) {
        decimal = decimal.slice(0, 18)
      }
      const decimalBN = new this.ethUtil.BN(decimal, 10)
      eth = eth.add(decimalBN)
    }
    return eth
  }

  normalizeNumberToWei (n, currency) {
    var enlarged = n * 10000
    var amount = new this.ethUtil.BN(String(enlarged), 10)
    return this.normalizeToWei(amount, currency).div(this.multiple)
  }

  readableDate (ms) {
    var date = new Date(ms)
    var month = date.getMonth()
    var day = date.getDate()
    var year = date.getFullYear()
    var hours = date.getHours()
    var minutes = '0' + date.getMinutes()
    var seconds = '0' + date.getSeconds()

    var dateStr = `${month}/${day}/${year}`
    var time = `${hours}:${minutes.substr(-2)}:${seconds.substr(-2)}`
    return `${dateStr} ${time}`
  }

  isHex (str) {
    return Boolean(str.match(/^(0x)?[0-9a-fA-F]+$/))
  }

  bnMultiplyByFraction (targetBN, numerator, denominator) {
    const numBN = new this.ethUtil.BN(numerator)
    const denomBN = new this.ethUtil.BN(denominator)
    return targetBN.mul(numBN).div(denomBN)
  }

  getContractAtAddress (tokenAddress) {
    return global.eth.contract(this.abi).at(tokenAddress)
  }

  exportAsFile (filename, data) {
    // source: https://stackoverflow.com/a/33542499 by Ludovic Feltz
    const blob = new Blob([data], { type: 'text/csv' })
    if (window.navigator.msSaveOrOpenBlob) {
      window.navigator.msSaveBlob(blob, filename)
    } else {
      const elem = window.document.createElement('a')
      elem.target = '_blank'
      elem.href = window.URL.createObjectURL(blob)
      elem.download = filename
      document.body.appendChild(elem)
      elem.click()
      document.body.removeChild(elem)
    }
  }

  allNull (obj) {
    return Object.entries(obj).every(([key, value]) => value === null)
  }

  getTokenAddressFromTokenObject (token) {
    return Object.values(token)[0].address.toLowerCase()
  }

  /**
   * Safely checksumms a potentially-null address
   *
   * @param {String} [address] - address to checksum
   * @returns {String} - checksummed address
   */
  checksumAddress (address) {
    return address ? this.ethUtil.toChecksumAddress(address) : ''
  }
}

export { MetamaskUtil as default }
