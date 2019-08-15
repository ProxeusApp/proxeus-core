/**
 * Created by Serafina on 26.10.2016.
 */

var Translator = function () {
  this.RADIX = 10
  this.TOKEN_LENGTH = 3
  this.MAX_NUMBERS = 9
  this.SINGLE_INDEX = 0
  this.TEN_INDEX = 1
  this.HUNDRED_INDEX = 2

  // start function
  this.toWords = function (locale, val, doNotUseFallback) {
    locale = locale.toLowerCase()
    doNotUseFallback = doNotUseFallback || false
    var numbers = this.tokenize(val, this.TOKEN_LENGTH)

    // parse input:
    if (isNaN(val)) {
      return FTG.translate('ct.not.a.number').msgFormat(val)
    } else if (numbers.length * this.TOKEN_LENGTH > this.MAX_NUMBERS) {
      return FTG.translate('ct.number.is.to.long').msgFormat(this.MAX_NUMBERS)
    }

    var dictionary = this.dictionary[locale]

    if (!dictionary && !doNotUseFallback) {
      dictionary = this.dictionary.en
    } else if (!dictionary && doNotUseFallback) {
      throw {
        name: 'Error',
        message: FTG.translate('ct.lang.not.supported')
      }
    }

    // Deal with exceptions - zero
    if (numbers[this.SINGLE_INDEX] === 0 && numbers.length === 1) {
      return dictionary.exceptions[numbers[this.SINGLE_INDEX]]
    }

    var words = []
    if (this[locale]) {
      for (var idx = 0, max = numbers.length; idx < max; idx++) {
        words.unshift(
          this[locale].translate(this.tokenize(numbers[idx], 1), idx, max,
            dictionary))
      }
    }
    return words.join('')
  }

  // returns array of numbers
  this.tokenize = function (val, tokenLength) {
    if (val === 0) {
      return [0]
    }
    var tokens = []
    var base = Math.pow(this.RADIX, tokenLength)
    while (val) {
      tokens.push(val % base)
      val = parseInt(val / base, this.RADIX)
    }
    return tokens
  }

  /* >>>>>>>>> AVAILABLE TRANSLATING FUNCTIONS: */
  var _ = this

  this.en = function () {}
  this.de = function () {}

  // English:
  this.en.translate = function (numbers, index, max, dictionary) {
    var hundred = ''
    var ten = ''
    var single = ''
    var radix = ' ' + _.getRadix(numbers, index, dictionary)

    if (numbers[_.HUNDRED_INDEX]) {
      hundred = numbers[_.TEN_INDEX] || numbers[this.SINGLE_INDEX]
        ? _.getOnes(numbers[_.HUNDRED_INDEX], dictionary) + ' ' +
        dictionary.hundred + ' ' + dictionary.delimiters[1] + ' '
        : _.getOnes(numbers[_.HUNDRED_INDEX], dictionary) + ' ' +
        dictionary.hundred
    }

    if (numbers[_.TEN_INDEX]) {
      ten = _.getTeens(numbers[_.SINGLE_INDEX], dictionary)
    }

    if (numbers[_.TEN_INDEX] >= 2) {
      ten = numbers[_.SINGLE_INDEX]
        ? _.getTens(numbers[_.TEN_INDEX], dictionary) +
        dictionary.delimiters[0] +
        _.getOnes(numbers[_.SINGLE_INDEX], dictionary)
        : _.getTens(numbers[_.TEN_INDEX], dictionary)
    }

    if (!numbers[_.TEN_INDEX]) {
      single = _.getOnes(numbers[_.SINGLE_INDEX], dictionary)
    }

    if (index + 1 < max && (numbers[_.HUNDRED_INDEX] || numbers[_.TEN_INDEX] ||
        numbers[_.SINGLE_INDEX])) {
      hundred = ' ' + hundred
    }

    if (index === 0 && index + 1 < max && !numbers[_.HUNDRED_INDEX] &&
      (numbers[_.TEN_INDEX] || numbers[_.SINGLE_INDEX])) {
      hundred = ' ' + dictionary.delimiters[1] + ' '
    }

    return hundred + ten + single + radix
  }

  // German:
  this.de.translate = function (numbers, index, max, dictionary) {
    var token = ''
    var hundred = ''
    var ten = ''
    var single = ''
    var radix = _.getRadix(numbers, index, dictionary)

    if (numbers[_.HUNDRED_INDEX]) {
      hundred = numbers[_.TEN_INDEX] || numbers[this.SINGLE_INDEX]
        ? _.getOnes(numbers[_.HUNDRED_INDEX], dictionary) + '' +
        dictionary.hundred + '' + dictionary.delimiters[1] + ''
        : _.getOnes(numbers[_.HUNDRED_INDEX], dictionary) + '' +
        dictionary.hundred
    }

    if (numbers[_.TEN_INDEX]) {
      ten = _.getTeens(numbers[_.SINGLE_INDEX], dictionary)
    }

    if (numbers[_.TEN_INDEX] >= 2) {
      ten = numbers[_.SINGLE_INDEX]
        ? _.getOnes(numbers[_.SINGLE_INDEX], dictionary) +
        dictionary.delimiters[0] + _.getTens(numbers[_.TEN_INDEX], dictionary)
        : _.getTens(numbers[_.TEN_INDEX], dictionary)
    }

    if (!numbers[_.TEN_INDEX]) {
      single = _.getOnes(numbers[_.SINGLE_INDEX], dictionary)
    }

    if (index + 1 < max && (numbers[_.HUNDRED_INDEX] || numbers[_.TEN_INDEX] ||
        numbers[_.SINGLE_INDEX])) {
      hundred = '' + hundred
    }

    // handle exceptions
    if (radix != 'tausend' && radix != '') {
      token = hundred + ten + single + radix
      if (token.substr(token.indexOf(radix) - 3, token.indexOf(radix)) ==
        'ein') {
        single = 'eine'
      } else {
        radix += 'en' // add plural ending
      }
      radix = ' ' + radix + ' ' // add
    } else if (radix == '' && single == 'ein') {
      single += 's'
    }
    return hundred + ten + single + radix
  }
  /* TRANSLATING FUNCTIONS END <<<<<<<<< */

  this.getOnes = function (number, dictionary) {
    return dictionary.ones[number]
  }
  this.getTens = function (number, dictionary) {
    return dictionary.tens[number]
  }
  this.getTeens = function (number, dictionary) {
    return dictionary.teens[number]
  }
  this.getRadix = function (numbers, index, dictionary) {
    var radix = ''
    if (index > 0 && (numbers[this.HUNDRED_INDEX] || numbers[this.TEN_INDEX] ||
        numbers[this.SINGLE_INDEX])) {
      radix = dictionary.radix[index]
    }

    return radix
  }

  // all available dictionaries:
  this.dictionary = {
    en: {
      zero: 'zero',
      ones: [
        '',
        'one',
        'two',
        'three',
        'four',
        'five',
        'six',
        'seven',
        'eight',
        'nine'],
      teens: [
        'ten',
        'eleven',
        'twelve',
        'thirteen',
        'fourteen',
        'fifteen',
        'sixteen',
        'seventeen',
        'eighteen',
        'nineteen'],
      tens: [
        '',
        '',
        'twenty',
        'thirty',
        'forty',
        'fifty',
        'sixty',
        'seventy',
        'eighty',
        'ninety'],
      hundred: 'hundred',
      radix: ['', 'thousand', 'million'],
      delimiters: ['-', 'and']
    },
    de: {
      zero: 'null',
      ones: [
        '',
        'ein',
        'zwei',
        'drei',
        'vier',
        'fünf',
        'sechs',
        'sieben',
        'acht',
        'neun'],
      teens: [
        'zehn',
        'elf',
        'zwölf',
        'dreizehn',
        'vierzehn',
        'fünfzehn',
        'sechszehn',
        'siebzehn',
        'achtzehn',
        'neunzehn'],
      tens: [
        '',
        '',
        'zwanzig',
        'dreissig',
        'vierzig',
        'fünfzig',
        'sechszig',
        'siebzig',
        'achtzig',
        'neunzig'],
      hundred: 'hundert',
      radix: ['', 'tausend', 'Million'],
      delimiters: ['und', '']
    }
  }
}

export default Translator
