var cglbl = {
  getLocale: function () {
    var lang = document.getElementsByTagName('HTML')[0].getAttribute('lang')
    if (!lang) {
      lang = navigator.language || navigator.userLanguage
    }
    return lang
  },
  showMsg: function (kind, msg, selector) {
    var deleteBtn = '<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>'
    var absoluteStyles = false ? 'style="position:absolute;width:100%;"' : ''
    var msgEl = $('<div class="alert alert-' + kind + ' alert-dismissible" role="alert" ' + absoluteStyles + '>' + deleteBtn + msg + '</div>')
    var msgContainer
    if (selector && typeof selector === 'string' && selector.length) {
      msgContainer = $(selector)
    } else {
      msgContainer = $('.msg-container.msg-container-active')
    }

    if (!msgContainer.is(':visible')) {
      msgContainer.slideDown('fast')
    }
    var childs = msgContainer.children('div.alert')
    var top = 0
    if (childs.length > 0) {
      var lastChild = childs.last()
      var topPx = lastChild.css('top')
      topPx = topPx.replace(/px/g, '')
      top = parseInt(topPx)
      top = top + 2
    }
    msgEl.css('top', top + 'px')
    msgEl.on('close.bs.alert', function () {
      $(this).remove()
      if (!msgContainer.children('div.alert').length) {
        msgContainer.slideUp('fast')
      }
    })
    setTimeout(function () {
      if (msgEl.length) {
        msgEl.fadeOut({
          complete: function () {
            msgEl.trigger('close.bs.alert')
          }
        })
      }
    }, 8000)
    msgContainer.append(msgEl)
  },
  redirect: function (url) {
    if (!url) {
      url = window.location.href
    } else if (/^\//.test(url)) {
      url = window.location.origin + url
    }
    // given for completeness, essentially an alias to window.location.href
    //         window.location.href = url;
    window.location.replace(url)
  },
  getScrollBarWidth: function () {
    if (!this.scrollBarWith) {
      var inner = document.createElement('p')
      inner.style.width = '100%'
      inner.style.height = '200px'

      var outer = document.createElement('div')
      outer.style.position = 'absolute'
      outer.style.top = '0px'
      outer.style.left = '0px'
      outer.style.visibility = 'hidden'
      outer.style.width = '200px'
      outer.style.height = '150px'
      outer.style.overflow = 'hidden'
      outer.appendChild(inner)

      document.body.appendChild(outer)
      var w1 = inner.offsetWidth
      outer.style.overflow = 'scroll'
      var w2 = inner.offsetWidth
      if (w1 == w2) w2 = outer.clientWidth

      document.body.removeChild(outer)

      this.scrollBarWith = (w1 - w2)
    }
    return this.scrollBarWith
  },
  sizeOf: function (obj) {
    var size = 0; var key
    for (key in obj) {
      if (obj.hasOwnProperty(key)) size++
    }
    return size
  },

  hex: {
    encode: function (s) {
      var h = ''
      for (var i = 0; i < s.length; i++) {
        h += s.charCodeAt(i).toString(16)
      }
      return h
    },
    decode: function (h) {
      var s = ''
      for (var i = 0; i < h.length; i += 2) {
        s += String.fromCharCode(parseInt(h.substr(i, 2), 16))
      }
      return s
    }
  },
  randomId: function () {
    return Math.random().toString(36).substring(2, 18) + Math.random().toString(36).substring(2, 18)
  }
}

$.fn.fillForm = function (data, keyPrefix) {
  var isRootCall = false
  if (!data) return
  if (!keyPrefix) {
    isRootCall = true
    keyPrefix = ''
  }
  var isFileObject = function (fobj) {
    return fobj && fobj.path && fobj.name
  }
  for (var key2 in data) {
    if (data.hasOwnProperty(key2)) {
      if (data[key2] !== null && typeof data[key2] === 'object' && !$.isArray(data[key2])) {
        // console.log("is Object: "+key2);
        if (isFileObject(data[key2])) {
          var fullKey = keyPrefix + key2
          var element = this.find("[name='" + fullKey + "']")
          if (element.length) {
            element.checkAndFillInputValue(fullKey, data[key2], element)
          }
        } else {
          this.fillForm(data[key2], keyPrefix + key2 + '.')
        }
      } else {
        var fullKey = keyPrefix + key2
        var element = this.find("[name='" + fullKey + "']")
        if ($.isArray(data[key2]) && element.length < data[key2].length) {
          try {
            var hitTimes = data[key2].length - element.length
            var inputFieldPushBtn = $(element.first().attr('addBtnRef'))
            for (var h = 0; h < hitTimes; h++) {
              inputFieldPushBtn.click()
            }
            element = this.find("[name='" + fullKey + "']")
          } catch (someErrorWhenPushingMissingFields) {
            console.log(someErrorWhenPushingMissingFields)
          }
        }
        if (!element.length && key2 === 'id') {
          var selectName = keyPrefix
          if (selectName.length > 0 && selectName.charAt(selectName.length - 1) == '.') {
            selectName = selectName.substring(0, selectName.length - 1)
          }
          element = this.find("[name='" + selectName + "']")
        }
        if (element.length) {
          var nt
          for (var e = 0; e < element.length; ++e) {
            nt = $(element[e])
            if ($.isArray(data[key2])) {
              for (var omg = 0; omg < data[key2].length; ++omg) {
                if (nt.checkAndFillInputValue(fullKey, data[key2][omg], element)) {
                  data[key2][omg] = '_![DELETE-MARK]!_'
                  var newArray = []
                  for (var om = 0; om < data[key2].length; ++om) {
                    if (data[key2][om] !== '_![DELETE-MARK]!_') {
                      newArray.push(data[key2][om])
                    }
                  }
                  data[key2] = newArray
                  break
                }
              }
            } else {
              nt.checkAndFillInputValue(fullKey, data[key2], element)
            }
          }
        }
      }
    }
  }
  if (isRootCall) {
    this.trigger('fillForm-done', [data])
  }
}
$.fn.doCompAction = function (init) {
  if ($.isFunction(this.data('fm_action'))) {
    this.data('fm_action').apply(this, arguments)
  }
}
$.fn.checkAndFillInputValue = function (fullKey, dataVal, element) {
  var _ = this
  if (_.attr('type') === 'radio') {
    _.removeAttr('checked')
    if (_.attr('name') === fullKey && _.attr('value') === (dataVal === null ? '' : dataVal)) {
      _.prop('checked', true)
      _.parents('.field-parent').attr('data-status', 'success')
      _.attr('data-status', 'success')
      _.doCompAction(true)
      _.trigger('change', [{ init: true }])
      return true
    }
  } else if (_.attr('type') === 'checkbox') {
    var cbLength = element.length
    if (cbLength == 1) { // single
      _.removeAttr('checked')
      _.prop('checked', !!dataVal)
      _.parents('.field-parent').attr('data-status', 'success')
      _.attr('data-status', 'success')
      _.doCompAction(true)
      _.trigger('change', [{ init: true }])
      return true
    } else { // multi
      _.removeAttr('checked')
      if (_.attr('name') === fullKey && _.attr('value') === (dataVal === null ? '' : dataVal)) {
        _.prop('checked', true)
        _.parents('.field-parent').attr('data-status', 'success')
        _.attr('data-status', 'success')
        _.doCompAction(true)
        _.trigger('change', [{ init: true }])
        return true
      }
    }
  } else {
    if (_.attr('type') === 'file') {
      try {
        if (dataVal) {
          _.doCompAction(true)
          var path
          var name = ''
          if (typeof dataVal === 'object') {
            if (dataVal.name) {
              name = dataVal.name
            }
          }
          var id
          if (_.attr('data-label')) {
            id = '#' + _.attr('data-label').replace('#', '')
            $(id).html(name)
          }
          var up = _.data('uriPrefix')
          if (/.*\/$/.test(up) === false) {
            up += '/'
          }
          path = up + _.attr('name')
          if (_.attr('data-img')) {
            var id = '#' + _.attr('data-img').replace('#', '')
            $(id).attr('src', path)
          }
          if (_.attr('imgRef')) {
            $(_.attr('imgRef')).attr('src', path)
          }
          return true
        }
      } catch (imgSrcEx) {
        console.log(imgSrcEx)
      }
    } else {
      if (_[0].value !== undefined) {
        _.val(dataVal)
        _.parents('.field-parent').attr('data-status', 'success')
        _.attr('data-status', 'success')
        _.doCompAction(true)
        if (_.is('select')) {
          _.trigger('change', [{ init: true }])
        }
        return true
      } else {
        if (_.attr('data-fill') === 'false') {

        } else {
          _.parents('.field-parent').attr('data-status', 'success')
          _.attr('data-status', 'success')
          _.doCompAction(true)
          _.html(dataVal)
          return true
        }
      }
    }
  }
  return false
}

var FTGlobal = function () {
  this.translations = {}
  this.javaPattern = {
    toJsTimePattern: function (javaPat) {
      var m = javaPat.match(/(H|h).*(m|s|S)/g)
      if (m && m.length) {
        return m[0]
      }
      return null
    },
    toJsDatePattern: function (javaPat) {
      var m = javaPat.match(/(M|d|y).*(M|d|y)/g)
      if (m && m.length) {
        return m[0].replace('MM', 'mm').replace('yyyy', 'yy')
      }
      return null
    },
    toJsPattern: function (javaPat) {
      return javaPat.replace('MM', 'mm').replace('yyyy', 'yy')
    }
  }
}

FTGlobal.prototype.translate = function (key, value) {
  if (key) {
    if (value) {
      this.translations[key] = value
    } else {
      value = this.translations[key]
      if (value) {
        return value
      } else {
        return key
      }
    }
  }
  return ''
}

FTGlobal.prototype.createDatepicker = function (elem, options) {
  try {
    var normalize = function (pattern) {
      return pattern.replace('dd', 'd')
        .replace('DD', 'd')
        .replace('mm', 'm')
        .replace('MM', 'm')
        .replace('YYYY', 'Y')
        .replace('yyyy', 'Y')
        .replace('hh', 'H')
        .replace('HH', 'H')
        .replace('ii', 'i')
        .replace('II', 'i')
    }
    var opt = {}
    if (options && options.datePattern) {
      opt.dateFormat = normalize(options.datePattern)
      opt.enableTime = opt.dateFormat.indexOf('H') >= 0 || opt.dateFormat.indexOf('i') >= 0
    }
    try {
      opt.defaultDate = elem.val()
      var pickr = flatpickr(elem[0], opt)
      if (elem.val()) {
        if (pickr) {
          pickr.setDate(elem.val())
        }
      }
    } catch (e) { console.log(e) }
  } catch (dateExc) {
    console.log(dateExc)
  }
}

FTGlobal.prototype.fileTypeEmpty = function (file) {
  return file.type === '' || file.type === undefined || file.type === null
}

FTGlobal.prototype.randomId = function () {
  return Math.floor((1 + Math.random()) * 0x10000)
    .toString(16)
    .substring(1)
}

var FTG = new FTGlobal()

if (!String.prototype.startsWith) {
  (function () {
    'use strict' // needed to support `apply`/`call` with `undefined`/`null`
    var defineProperty = (function () {
      // IE 8 only supports `Object.defineProperty` on DOM elements
      try {
        var object = {}
        var $defineProperty = Object.defineProperty
        var result = $defineProperty(object, object, object) && $defineProperty
      } catch (error) {}
      return result
    }())
    var toString = {}.toString
    var startsWith = function (search) {
      if (this == null) {
        throw TypeError()
      }
      var string = String(this)
      if (search && toString.call(search) == '[object RegExp]') {
        throw TypeError()
      }
      var stringLength = string.length
      var searchString = String(search)
      var searchLength = searchString.length
      var position = arguments.length > 1 ? arguments[1] : undefined
      // `ToInteger`
      var pos = position ? Number(position) : 0
      if (pos != pos) { // better `isNaN`
        pos = 0
      }
      var start = Math.min(Math.max(pos, 0), stringLength)
      // Avoid the `indexOf` call if no match is possible
      if (searchLength + start > stringLength) {
        return false
      }
      var index = -1
      while (++index < searchLength) {
        if (string.charCodeAt(start + index) != searchString.charCodeAt(index)) {
          return false
        }
      }
      return true
    }
    if (defineProperty) {
      defineProperty(String.prototype, 'startsWith', {
        'value': startsWith,
        'configurable': true,
        'writable': true
      })
    } else {
      String.prototype.startsWith = startsWith
    }
  }())
}

if (typeof String.prototype.endsWith !== 'function') {
  String.prototype.endsWith = function (suffix) {
    return this.indexOf(suffix, this.length - suffix.length) !== -1
  }
}

$.fn.nextParentWithClass = function (attrClass) {
  var maxIndex = 80
  var index = 1
  var parent = this
  while (!parent.hasClass(attrClass)) {
    parent = parent.parent()
    ++index
    if (maxIndex < index) {
      break
    }
  }
  return parent
}
$.fn.cleanFieldErrors = function (data) {
  var _this = this
  if (!data) {
    data = {}
  }
  if (!data.parentClass) {
    data.parentClass = 'field-parent'
  }
  var tagName = this[0].tagName.toUpperCase()
  if (tagName === 'INPUT' || tagName === 'TEXTAREA' || tagName === 'SELECT') {
    _this = this.nextParentWithClass(data.parentClass)
  }
  _this.find('span.error').remove()
  _this.find('.error').removeClass('error')
  return this
}
$.fn.placeError = function (obj, parentClass) {
  if (this && this.length > 0) {
    var stringToId = function (str) {
      var hash = 0
      var char
      if (!str || str.length == 0) return hash
      for (var ind = 0; ind < str.length; ind++) {
        char = str.charCodeAt(ind)
        hash = ((hash << 5) - hash) + char
        hash = hash & hash // Convert to 32bit integer
      }
      return hash
    }
    var fieldParent = this.nextParentWithClass(parentClass)
    if (!fieldParent.length) {
      fieldParent = this.parent()
    }
    if (fieldParent.find('.errors').length === 1) {
      fieldParent = fieldParent.find('.errors')
    }
    var spanErrorId = 'err' + stringToId(obj.msg)
    if (fieldParent.find('#' + spanErrorId + '').length == 0) {
      var span = $('<span id="' + spanErrorId + '" class="error">' + obj.msg + '</span>')
      // span.html(JSON.stringify(o));
      fieldParent.append(span)
    }
  }
}
$.fn.isElVisible = function () {
  if (this.length === 0) {
    return true
  }
  var style = window.getComputedStyle(this[0])
  return (style.display !== 'none')
}
$.fn.showFieldErrors = function (data) {
  this.cleanFieldErrors(data)

  if (data) {
    if (!data.parentClass) {
      data.parentClass = 'field-parent'
    }
    if (data.errors) {
      var isInputField = function (jsEl) {
        var tagName = jsEl.tagName.toUpperCase()
        return tagName === 'INPUT' || tagName === 'TEXTAREA' || tagName === 'SELECT'
      }
      var o = null
      if (this.length == 1 && isInputField(this[0])) {
        if ($.isArray(data.errors)) {
          for (var i = 0; i < data.errors.length; i++) {
            o = data.errors[i]
            if (o && o.msg) {
              var err = o
              var nindex = this.attr('nindex')
              if (!isNaN(err.i) && !isNaN(nindex)) {
                if (err.i == parseInt(nindex)) {
                  this.placeError(err, data.parentClass)
                }
              } else {
                this.placeError(err, data.parentClass)
              }
            }
          }
        } else {
          for (var c in data.errors) {
            if (data.errors.hasOwnProperty(c)) {
              o = data.errors[c]
              if (c === this.attr('name')) {
                for (var a in o) {
                  var err = o[a]
                  var nindex = this.attr('nindex')
                  if (!isNaN(err.i) && !isNaN(nindex)) {
                    if (err.i == parseInt(nindex)) {
                      this.placeError(err, data.parentClass)
                    }
                  } else {
                    this.placeError(err, data.parentClass)
                  }
                }
              }
            }
          }
        }
      } else {
        var newThis = this
        if (isInputField(this[0])) {
          newThis = $(this[0]).parents('form')
        }
        o = null
        for (var c in data.errors) {
          if (data.errors.hasOwnProperty(c)) {
            o = data.errors[c]
            for (var a in o) {
              var err = o[a]
              if (err.i) {
                try {
                  var tel = $(newThis.find('[name="' + c + '"]')[err.i])

                  tel.placeError(err, data.parentClass)
                } catch (abc) {
                  console.log(abc)
                }
              } else {
                newThis.find('[name="' + c + '"]').each(function () {
                  var $this = $(this)

                  $this.placeError(err, data.parentClass)
                })
              }
            }
          }
        }
      }
    }
  }
  return this
}

$.fn.enableOrDisableContent = function (enable, colorStr) {
  var contentDisablerEl = this.find('.content-disabler')
  if (enable) {
    contentDisablerEl.remove()
  } else {
    if (!contentDisablerEl.length) {
      this.css('position', 'relative')
      if (!colorStr) {
        colorStr = 'rgba(236, 236, 236, 0.54)'
      }
      this.append('<div class="content-disabler" style="position:absolute;width:100%;height:100%;top:0;left:0;background:' + colorStr + ';"></div>')
    }
  }
  return parent
}
$.fn.fileElement = function (element) {
  if (!element.v) {
    this.find('.view-btn').hide()
  }
  if (!element.d) {
    this.find('.download-btn').hide()
  }
  if (!element.r) {
    this.find('.tf-delete').hide()
  }
}
jQuery.fn.center = function () {
  this.css('margin-top', Math.max(0, (($(window).height() - $(this).outerHeight()) / 2) +
            $(window).scrollTop()) + 'px')
  this.css('margin-left', Math.max(0, (($(window).width() - $(this).outerWidth()) / 2) +
            $(window).scrollLeft()) + 'px')
  return this
}

$.fn.updateInputFile = function () {
  this.each(function () {
    var _ = $(this)
    var up = _.data('uriPrefix')
    if (up) {
      if (/.*\/$/.test(up) === false) {
        up += '/'
      }
      var fileUrl = up + _.attr('name')
      var imgId = _.attr('data-img')
      if (imgId) {
        $('#' + imgId.replace('#', '')).attr('src', fileUrl + '?v=' + new Date().getSeconds())
      }
    }
  })
}
$.fn.assignSubmitOnChange = function (options) {
  if (!this || !this[0]) {
    return
  }
  var changeEvent = function (e, arg2, arg3) {
    var changeOptions = {
      ajaxRequest: true
    }
    if (arg2) {
      changeOptions = $.extend(changeOptions, arg2)
    }
    var _this = $(this)
    var changedData = {}
    _this.doCompAction()
    if (_this.attr('type') === 'file') {
      var files = this.files
      var file
      if (files && files.length) {
        file = files[0]
      }
      if (file) {
        if (changeOptions.ajaxRequest && (options.fileUrl || options.url)) {
          if (!file || _this.data('_lastData') === file.size) {
            return true
          }
          _this.data('_lastData', file.size)
          var up = _this.data('uriPrefix')
          if (/.*\/$/.test(up) === false) {
            up += '/'
          }
          var fileUrl = up + _this.attr('name')
          var myReq = $.ajax({
            url: fileUrl,
            type: 'POST',
            responsibleEl: _this,
            data: file,
            processData: false,
            fileFieldName: _this.attr('name'),
            beforeSend: function (request) {
              this.myReq = myReq
              if ($.isFunction(options.beforeSend)) {
                try {
                  var newFile = options.beforeSend.apply(this, [changedData, null, _this, this.myReq])
                  if (newFile && file && newFile.size !== file.size) {
                    file = newFile
                    if (!file || _this.data('_lastData') === file.size) {
                      this.myReq.abort()
                      return true
                    }
                    _this.data('_lastData', file.size)
                    this.data = file
                  }
                } catch (dontCareEx) {
                  console.log(dontCareEx)
                }
              }
              var status = 'pending'
              this.responsibleEl.parents('.field-parent').attr('data-status', status)
              this.responsibleEl.attr('data-status', status)
              this.responsibleEl.trigger('onComStatusChange', [status, file.name])
              var encodedFilename = encodeURI(file.name)
              request.setRequestHeader('Content-disposition', 'attachment; filename=' + encodedFilename)
              var type = file.type
              if (FTG.fileTypeEmpty(file)) {
                try {
                  type = FTG.MimeType.lookup(encodedFilename)
                } catch (eee) {}
              }
              request.setRequestHeader('Content-Type', type)
              request.setRequestHeader('File-Name', encodedFilename)
              if (this.fileFieldName) {
                request.setRequestHeader('File-FieldName', this.fileFieldName)
              }
              request.setRequestHeader('File-Size', file.size)
            },
            success: function (data, textStatus, xhr) {
              if (xhr.status == 205) {
                FTG.jsRedirect()
              }
              var status = 'success'
              if (data && data.path && /\//.test(data.path)) {
                var imgId = this.responsibleEl.attr('data-img')
                if (imgId) {
                  $('#' + imgId.replace('#', '')).attr('src', data.path + '?v=' + new Date().getSeconds())
                }
              } else {
                this.responsibleEl.updateInputFile()
              }
              this.responsibleEl.parents('.field-parent').attr('data-status', status)
              this.responsibleEl.attr('data-status', status)
              this.responsibleEl.trigger('onComStatusChange', [status, (data ? data['name'] : '')])
              this.responsibleEl.cleanFieldErrors()
              try {
                if ($.isFunction(options.success)) {
                  options.success.apply(this, [data, textStatus, xhr, this.myReq])
                }
              } catch (eee) { console.log(eee) }
            },
            error: function (data, a2, a3, a4) {
              var status = 'error'
              this.responsibleEl.parents('.field-parent').attr('data-status', status)
              this.responsibleEl.attr('data-status', status)
              this.responsibleEl.trigger('onComStatusChange', [status])
              console.log(data.responseJSON)
              if (data.status == 412 || data.status === 422) {
                this.responsibleEl.showFieldErrors(data.responseJSON)
              }
              if (data.status == 502) {
                this.responsibleEl.showFieldErrors({ errors: [{ field: this.responsibleEl.attr('name'), message: FTG.translate('file.limit.exceeded') }] })
              }
              try {
                if ($.isFunction(options.error)) {
                  options.error.apply(this, [data, a2, a3, this.myReq])
                }
              } catch (eee) { console.log(eee) }
            }
          })
        }
      }
    } else {
      var nindex = _this.attr('nindex')
      var form = _this.parents('form')
      var setChangedValue = function (ev) {
        var _t = $(this)
        var thisName = _t.attr('name')
        var val = _t.val()
        if (_t.hasClass('vcode') && _t.attr('data-ref')) {
          val = _t.parent().find('#' + _t.attr('data-ref')).val() + '[' + val + ']'
        }
        var hasNIndexAttr, nindex
        if (_t.attr('type') === 'checkbox') {
          if (_t.is(':checked')) {
            if (val) {
              if ($.isArray(changedData[thisName])) {
                hasNIndexAttr = !isNaN(_t.attr('nindex'))
                if (hasNIndexAttr) {
                  nindex = parseInt(_t.attr('nindex'))
                  changedData[thisName][nindex] = val
                } else {
                  changedData[thisName].push(val)
                }
              } else {
                changedData[thisName] = val
              }
            } else {
              if ($.isArray(changedData[thisName])) {
                hasNIndexAttr = !isNaN(_t.attr('nindex'))
                if (hasNIndexAttr) {
                  nindex = parseInt(_t.attr('nindex'))
                  changedData[thisName][nindex] = true
                } else {
                  changedData[thisName].push(true)
                }
              } else {
                changedData[thisName] = true
              }
            }
          }
        } else {
          if ($.isArray(changedData[thisName])) {
            hasNIndexAttr = !isNaN(_t.attr('nindex'))
            if (hasNIndexAttr) {
              nindex = parseInt(_t.attr('nindex'))
              changedData[thisName][nindex] = val
            } else {
              if (val) {
                changedData[thisName].push(val)
              }
            }
          } else {
            changedData[thisName] = val
          }
        }
      }
      var allArrayItemFields = form.find("[name='" + _this.attr('name') + "']")
      var oneOfTheFieldsIsACheckbox = false
      allArrayItemFields.each(function () {
        if ($(this).attr('type') === 'checkbox') {
          oneOfTheFieldsIsACheckbox = true
          return false
        }
      })
      if (!isNaN(nindex) || (allArrayItemFields.length > 1 && oneOfTheFieldsIsACheckbox)) {
        changedData[_this.attr('name')] = []
        allArrayItemFields.each(setChangedValue)
      } else {
        if (allArrayItemFields.length > 1 && oneOfTheFieldsIsACheckbox) {
          changedData[_this.attr('name')] = []
          allArrayItemFields.each(setChangedValue)
        } else {
          _this.each(setChangedValue)
        }
      }
      var dataBeforeSend = JSON.stringify(changedData)
      if (changeOptions.ajaxRequest && options.url) {
        if (_this.attr('type') !== 'checkbox' && _this.attr('type') !== 'radio') {
          if (_this.data('_lastData') === dataBeforeSend) {
            return true
          }
          _this.data('_lastData', dataBeforeSend)
        }
        var myReq = $.ajax({
          responsibleEl: _this,
          type: 'POST',
          url: options.url,
          data: dataBeforeSend,
          contentType: 'application/json',
          beforeSend: function (request) {
            this.myReq = myReq
            if ($.isFunction(options.beforeSend)) {
              try {
                var newDat = options.beforeSend.apply(this, [changedData, dataBeforeSend, _this, this.myReq])
                if (newDat) {
                  changedData = newDat
                  dataBeforeSend = JSON.stringify(changedData)
                  if (_this.attr('type') !== 'checkbox' && _this.attr('type') !== 'radio') {
                    if (_this.data('_lastData') === dataBeforeSend) {
                      this.myReq.abort()
                      return true
                    }
                    _this.data('_lastData', dataBeforeSend)
                  }
                  this.data = dataBeforeSend
                }
              } catch (dontCareEx) {
                console.log(dontCareEx)
              }
            }
            var status = 'pending'
            this.responsibleEl.parents('.field-parent').attr('data-status', status)
            this.responsibleEl.attr('data-status', status)
            this.responsibleEl.trigger('onComStatusChange', [status])
          },
          success: function (data, textStatus, xhr) {
            var status = 'success'
            if (data && data.status) {
              status = data.status
            }
            this.responsibleEl.parents('.field-parent').attr('data-status', status)
            this.responsibleEl.attr('data-status', status)
            this.responsibleEl.trigger('onComStatusChange', [status])
            this.responsibleEl.cleanFieldErrors()
            try {
              if ($.isFunction(options.success)) {
                options.success.apply(this, [data, textStatus, xhr, this.myReq])
              }
            } catch (eee) {}
          },
          error: function (data, a2, a3, a4, a5) {
            var status = 'error'
            this.responsibleEl.parents('.field-parent').attr('data-status', status)
            this.responsibleEl.attr('data-status', status)
            this.responsibleEl.trigger('onComStatusChange', [status])
            if (data.status == 412 || data.status === 422) {
              this.responsibleEl.showFieldErrors(data.responseJSON)
            }
            try {
              if ($.isFunction(options.error)) {
                options.error.apply(this, [data, a2, a3, this.myReq])
              }
            } catch (eee) {}
          }
        })
      }
    }
  }
  var tagName = this[0].tagName.toLowerCase()

  if (tagName === 'input') {
    var type = this.attr('type')
    if (type === 'file') {
      this.data('uriPrefix', options.fileUrl ? options.fileUrl : options.url)
      this.updateInputFile()
    }
    if (/radio|checkbox|file/g.test(type)) {
      this.bind('change input', changeEvent)
    } else if (type === 'text') {
      this.bind('input paste change', changeEvent)
      if (this.hasClass('simple-date-field')) {
        this.bind('input paste change', changeEvent)
      }
    } else {
      this.bind('input paste change', changeEvent)
    }
  } else if (tagName === 'textarea') {
    this.bind('input change', changeEvent)
  } else if (tagName === 'select') {
    this.bind('change input', changeEvent)
    this.each(function () {
      var s = $(this)
      if (s.attr('data-status') !== 'success' && s.val()) {
        s.change()
      }
    })
  } else {
    this.find("input[type='file']").data('uriPrefix', options.fileUrl ? options.fileUrl : options.url)
    this.find("input[type='file']").updateInputFile()
    var selects = this.find('select')
    selects.bind('paste change input', changeEvent)
    selects.each(function () {
      var s = $(this)
      if (s.attr('data-status') !== 'success' && s.val()) {
        s.change()
      }
    })
    this.find('input.simple-date-field').bind('input paste change', changeEvent)
    this.find('textarea').bind('input paste change', changeEvent)
    this.find("input[type='text']").bind('input paste change', changeEvent)
    this.find("input[type='radio'], input[type='checkbox'], input[type='file']").bind('change', changeEvent)
  }
}
$.fn.serializeFormToObject = function (excludeFileFormFields) {
  var data = {}
  var allElementsHavingNameAttr = this.find('textarea[name], input[name], select[name]')
  var targetElement = null
  var key2 = null
  // var value = null;
  for (var i = 0; i < allElementsHavingNameAttr.length; ++i) {
    targetElement = $(allElementsHavingNameAttr[i])
    // if (excludeFileFormFields === true) {
    if (!targetElement.is(':visible') || targetElement.attr('type') === 'file') {
      continue
    }
    // }
    key2 = targetElement.attr('name')
    // value = data[key2];
    if (key2) {
      if (targetElement.length) {
        if (targetElement.attr('type') === 'radio') {
          data[key2] = this.find("input[name='" + key2 + "']:checked").val()
        } else if (targetElement.attr('type') === 'checkbox') {
          var checkBox = this.find("input[name='" + key2 + "']")
          if (checkBox.length === 1) { // single
            data[key2] = checkBox.is(':checked')
          } else { // multi
            let vals = []
            this.find("input[name='" + key2 + "']:checked").each(function (index) {
              vals.push($(this).val())
            })
            data[key2] = vals
          }
        } else {
          if (data[key2] === undefined && targetElement.attr('nindex') !== undefined) {
            data[key2] = []
          }
          let val
          if (targetElement[0].value !== undefined) {
            val = targetElement.val()
          } else {
            val = targetElement.html()
          }
          if (Array.isArray(data[key2])) {
            data[key2].push(val)
          } else {
            data[key2] = val
          }
        }
      }
      if (data[key2] === undefined || data[key2] === null) {
        data[key2] = ''
      }
    }
  }
  return data
}
$.fn.isValAsFilenameValid = function () {
  var fname = this.val()
  if (fname && fname.length > 0 && fname.length < 101) {
    var rg1 = /^[^\\/:\*\?"<>\|]+$/ // forbidden characters \ / : * ? " < > |
    var rg2 = /^\./ // cannot start with dot (.)
    var rg3 = /^(nul|prn|con|lpt[0-9]|com[0-9])(\.|$)/i // forbidden file names
    return rg1.test(fname) && !rg2.test(fname) && !rg3.test(fname)
  }
  return false
}
$.fn.makeSameHeight = function (options) {
  var _this = this
  if (!options) {
    options = { delay: 10 }
  }
  if (!options.delay) {
    options.delay = 10
  }
  var makeSameHeightFunc = function () {
    var maxHeight = 0
    var currentHeight = 0
    var maxHeightEl = null
    _this.each(function () {
      var t = $(this)
      t.css('height', '')
      currentHeight = t.height()
      if (maxHeight < currentHeight) {
        maxHeight = currentHeight
        maxHeightEl = t
      }
    })
    _this.each(function () {
      var t = $(this)
      currentHeight = t.height()
      if (currentHeight < maxHeight) {
        t.height(maxHeight)
      }
    })
  }
  $(window).resize(makeSameHeightFunc)
  setTimeout(makeSameHeightFunc, options.delay)
}

String.prototype.msgFormat = function () {
  var args = arguments
  return this.replace(/\{(\d+)\}/g, function () {
    return args[arguments[1]]
  })
}
jQuery.fn.copyHtmlToClipboard = function () {
  if (this.length === 0) {
    return this
  }
  var node = this[0]
  try {
    var range, selection
    selection = window.getSelection()
    range = document.createRange()
    range.selectNodeContents(node)
    selection.removeAllRanges()
    selection.addRange(range)
    document.execCommand('copy')
    $('#copyToclipboardInfo_').remove()
    var cInfo = $('<span id="copyToclipboardInfo_" draggable="true" style="border-color:transparent;border-radius: 4px;background: #4ef2d0ba;padding:10px;position:absolute;top:-40px;left:50px;z-index:100000000;white-space: nowrap;">copied to clipboard</span>')
    $(node).append(cInfo)
    setTimeout(function () {
      cInfo.remove()
    }, 4000)
  } catch (e) {
    console.log(e)
  }
}
/*!
 * JavaScript Cookie v2.1.2
 * https://github.com/js-cookie/js-cookie
 *
 * Copyright 2006, 2015 Klaus Hartl & Fagner Brack
 * Released under the MIT license
 */
!(function (e) { if (typeof define === 'function' && define.amd)define(e); else if (typeof exports === 'object')module.exports = e(); else { var n = window.Cookies; var t = window.Cookies = e(); t.noConflict = function () { return window.Cookies = n, t } } }(function () { function e () { for (var e = 0, n = {}; e < arguments.length; e++) { var t = arguments[e]; for (var o in t)n[o] = t[o] } return n } function n (t) { function o (n, r, i) { var c; if (typeof document !== 'undefined') { if (arguments.length > 1) { if (i = e({ path: '/' }, o.defaults, i), typeof i.expires === 'number') { var a = new Date(); a.setMilliseconds(a.getMilliseconds() + 864e5 * i.expires), i.expires = a } try { c = JSON.stringify(r), /^[\{\[]/.test(c) && (r = c) } catch (s) {} return r = t.write ? t.write(r, n) : encodeURIComponent(String(r)).replace(/%(23|24|26|2B|3A|3C|3E|3D|2F|3F|40|5B|5D|5E|60|7B|7D|7C)/g, decodeURIComponent), n = encodeURIComponent(String(n)), n = n.replace(/%(23|24|26|2B|5E|60|7C)/g, decodeURIComponent), n = n.replace(/[\(\)]/g, escape), document.cookie = [n, '=', r, i.expires ? '; expires=' + i.expires.toUTCString() : '', i.path ? '; path=' + i.path : '', i.domain ? '; domain=' + i.domain : '', i.secure ? '; secure' : ''].join('') }n || (c = {}); for (var p = document.cookie ? document.cookie.split('; ') : [], u = /(%[0-9A-Z]{2})+/g, d = 0; d < p.length; d++) { var f = p[d].split('='); var l = f.slice(1).join('='); l.charAt(0) === '"' && (l = l.slice(1, -1)); try { var m = f[0].replace(u, decodeURIComponent); if (l = t.read ? t.read(l, m) : t(l, m) || l.replace(u, decodeURIComponent), this.json) try { l = JSON.parse(l) } catch (s) {} if (n === m) { c = l; break }n || (c[m] = l) } catch (s) {} } return c } } return o.set = o, o.get = function (e) { return o(e) }, o.getJSON = function () { return o.apply({ json: !0 }, [].slice.call(arguments)) }, o.defaults = {}, o.remove = function (n, t) { o(n, '', e(t, { expires: -1 })) }, o.withConverter = n, o } return n(function () {}) }))

export default FTG
