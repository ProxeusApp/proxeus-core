/*
 *
 * Legacy imports
 *
 */
import flatpickr from 'flatpickr'

import FTG from './libs/legacy/global.js'

import Translator from './libs/number-translator.js'

window.Handlebars = require('handlebars')

window.flatpickr = flatpickr

window.$ = window.jQuery = require('jquery')
window.FTG = FTG
window.Translator = Translator

export default {}
