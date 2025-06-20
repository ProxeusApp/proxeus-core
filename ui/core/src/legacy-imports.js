/*
 *
 * Legacy imports
 *
 */

import interact from "interactjs";

import flatpickr from "flatpickr";
window.flatpickr = flatpickr;

import $ from "jquery";
window.jQuery = window.$ = $;

window.Handlebars = require("handlebars");

import FTG from "./libs/legacy/global.js";
window.FTG = FTG;

import Translator from "./libs/number-translator.js";
window.Translator = Translator;

require("./libs/legacy/formbuilder/jquery.html-svg-connect.js");
require("./libs/legacy/formbuilder/contextMenu/contextmenu.js");
require("./libs/legacy/formbuilder/splitpane/split-pane.js");

require("./libs/legacy/formbuilder/typeahead.js");
window.interact = interact;

window.ace = require("brace");
require("brace/ext/language_tools");
require("brace/mode/handlebars");
require("brace/mode/javascript");
require("brace/mode/json");
require("brace/mode/yaml");
require("brace/snippets/javascript");
require("brace/snippets/handlebars");
require("brace/snippets/yaml");
require("brace/snippets/json");
require("brace/theme/tomorrow_night");

export default {};
