/*! jQuery UI - v1.9.2 - 2012-11-23
 * http://jqueryui.com
 * Includes: jquery.ui.core.js, jquery.ui.widget.js, jquery.ui.mouse.js, jquery.ui.position.js, jquery.ui.accordion.js, jquery.ui.autocomplete.js, jquery.ui.button.js, jquery.ui.datepicker.js, jquery.ui.dialog.js, jquery.ui.draggable.js, jquery.ui.droppable.js, jquery.ui.effect.js, jquery.ui.effect-blind.js, jquery.ui.effect-bounce.js, jquery.ui.effect-clip.js, jquery.ui.effect-drop.js, jquery.ui.effect-explode.js, jquery.ui.effect-fade.js, jquery.ui.effect-fold.js, jquery.ui.effect-highlight.js, jquery.ui.effect-pulsate.js, jquery.ui.effect-scale.js, jquery.ui.effect-shake.js, jquery.ui.effect-slide.js, jquery.ui.effect-transfer.js, jquery.ui.menu.js, jquery.ui.progressbar.js, jquery.ui.resizable.js, jquery.ui.selectable.js, jquery.ui.slider.js, jquery.ui.sortable.js, jquery.ui.spinner.js, jquery.ui.tabs.js, jquery.ui.tooltip.js
 * Copyright (c) 2012 jQuery Foundation and other contributors Licensed MIT */

!(function (t, e) {
  function i (e, i) {
    var n; var o; var a; var r = e.nodeName.toLowerCase()
    return r === 'area' ? (n = e.parentNode, o = n.name, !(!e.href || !o || n.nodeName.toLowerCase() !== 'map') && !!(a = t('img[usemap=#' + o + ']')[0]) && s(a)) : (/input|select|textarea|button|object/.test(r) ? !e.disabled : r === 'a' ? e.href || i : i) && s(e)
  }

  function s (e) {
    return t.expr.filters.visible(e) && !t(e).parents().addBack().filter(function () {
      return t.css(this, 'visibility') === 'hidden'
    }).length
  }

  var n = 0; var o = /^ui-id-\d+$/
  t.ui = t.ui || {}, t.ui.version || (t.extend(t.ui, {
    version: '1.9.2',
    keyCode: {
      BACKSPACE: 8,
      COMMA: 188,
      DELETE: 46,
      DOWN: 40,
      END: 35,
      ENTER: 13,
      ESCAPE: 27,
      HOME: 36,
      LEFT: 37,
      NUMPAD_ADD: 107,
      NUMPAD_DECIMAL: 110,
      NUMPAD_DIVIDE: 111,
      NUMPAD_ENTER: 108,
      NUMPAD_MULTIPLY: 106,
      NUMPAD_SUBTRACT: 109,
      PAGE_DOWN: 34,
      PAGE_UP: 33,
      PERIOD: 190,
      RIGHT: 39,
      SPACE: 32,
      TAB: 9,
      UP: 38
    }
  }), t.fn.extend({
    _focus: t.fn.focus,
    focus: function (e, i) {
      return typeof e === 'number' ? this.each(function () {
        var s = this
        setTimeout(function () {
          t(s).focus(), i && i.call(s)
        }, e)
      }) : this._focus.apply(this, arguments)
    },
    scrollParent: function () {
      var e
      return e = t.ui.ie && /(static|relative)/.test(this.css('position')) || /absolute/.test(this.css('position')) ? this.parents().filter(function () {
        return /(relative|absolute|fixed)/.test(t.css(this, 'position')) && /(auto|scroll)/.test(t.css(this, 'overflow') + t.css(this, 'overflow-y') + t.css(this, 'overflow-x'))
      }).eq(0) : this.parents().filter(function () {
        return /(auto|scroll)/.test(t.css(this, 'overflow') + t.css(this, 'overflow-y') + t.css(this, 'overflow-x'))
      }).eq(0), /fixed/.test(this.css('position')) || !e.length ? t(document) : e
    },
    zIndex: function (i) {
      if (i !== e) return this.css('zIndex', i)
      if (this.length) {
        for (var s, n, o = t(this[0]); o.length && o[0] !== document;) {
          if (((s = o.css('position')) === 'absolute' || s === 'relative' || s === 'fixed') && (n = parseInt(o.css('zIndex'), 10), !isNaN(n) && n !== 0)) return n
          o = o.parent()
        }
      }
      return 0
    },
    uniqueId: function () {
      return this.each(function () {
        this.id || (this.id = 'ui-id-' + ++n)
      })
    },
    removeUniqueId: function () {
      return this.each(function () {
        o.test(this.id) && t(this).removeAttr('id')
      })
    }
  }), t.extend(t.expr[':'], {
    data: t.expr.createPseudo ? t.expr.createPseudo(function (e) {
      return function (i) {
        return !!t.data(i, e)
      }
    }) : function (e, i, s) {
      return !!t.data(e, s[3])
    },
    focusable: function (e) {
      return i(e, !isNaN(t.attr(e, 'tabindex')))
    },
    tabbable: function (e) {
      var s = t.attr(e, 'tabindex'); var n = isNaN(s)
      return (n || s >= 0) && i(e, !n)
    }
  }), t(function () {
    var e = document.body; var i = e.appendChild(i = document.createElement('div'))
    i.offsetHeight, t.extend(i.style, {
      minHeight: '100px',
      height: 'auto',
      padding: 0,
      borderWidth: 0
    }), t.support.minHeight = i.offsetHeight === 100, t.support.selectstart = 'onselectstart' in i, e.removeChild(i).style.display = 'none'
  }), t('<a>').outerWidth(1).jquery || t.each(['Width', 'Height'], function (i, s) {
    function n (e, i, s, n) {
      return t.each(o, function () {
        i -= parseFloat(t.css(e, 'padding' + this)) || 0, s && (i -= parseFloat(t.css(e, 'border' + this + 'Width')) || 0), n && (i -= parseFloat(t.css(e, 'margin' + this)) || 0)
      }), i
    }

    var o = s === 'Width' ? ['Left', 'Right'] : ['Top', 'Bottom']; var a = s.toLowerCase(); var r = {
      innerWidth: t.fn.innerWidth,
      innerHeight: t.fn.innerHeight,
      outerWidth: t.fn.outerWidth,
      outerHeight: t.fn.outerHeight
    }
    t.fn['inner' + s] = function (i) {
      return i === e ? r['inner' + s].call(this) : this.each(function () {
        t(this).css(a, n(this, i) + 'px')
      })
    }, t.fn['outer' + s] = function (e, i) {
      return typeof e !== 'number' ? r['outer' + s].call(this, e) : this.each(function () {
        t(this).css(a, n(this, e, !0, i) + 'px')
      })
    }
  }), t('<a>').data('a-b', 'a').removeData('a-b').data('a-b') && (t.fn.removeData = (function (e) {
    return function (i) {
      return arguments.length ? e.call(this, t.camelCase(i)) : e.call(this)
    }
  }(t.fn.removeData))), (function () {
    var e = /msie ([\w.]+)/.exec(navigator.userAgent.toLowerCase()) || []
    t.ui.ie = !!e.length, t.ui.ie6 = parseFloat(e[1], 10) === 6
  }()), t.fn.extend({
    disableSelection: function () {
      return this.bind((t.support.selectstart ? 'selectstart' : 'mousedown') + '.ui-disableSelection', function (t) {
        t.preventDefault()
      })
    },
    enableSelection: function () {
      return this.unbind('.ui-disableSelection')
    }
  }), t.extend(t.ui, {
    plugin: {
      add: function (e, i, s) {
        var n; var o = t.ui[e].prototype
        for (n in s) o.plugins[n] = o.plugins[n] || [], o.plugins[n].push([i, s[n]])
      },
      call: function (t, e, i) {
        var s; var n = t.plugins[e]
        if (n && t.element[0].parentNode && t.element[0].parentNode.nodeType !== 11) for (s = 0; s < n.length; s++) t.options[n[s][0]] && n[s][1].apply(t.element, i)
      }
    },
    contains: t.contains,
    hasScroll: function (e, i) {
      if (t(e).css('overflow') === 'hidden') return !1
      var s = i && i === 'left' ? 'scrollLeft' : 'scrollTop'; var n = !1
      return e[s] > 0 || (e[s] = 1, n = e[s] > 0, e[s] = 0, n)
    },
    isOverAxis: function (t, e, i) {
      return t > e && t < e + i
    },
    isOver: function (e, i, s, n, o, a) {
      return t.ui.isOverAxis(e, s, o) && t.ui.isOverAxis(i, n, a)
    }
  }))
}(jQuery)), (function (t, e) {
  var i = 0; var s = Array.prototype.slice; var n = t.cleanData
  t.cleanData = function (e) {
    for (var i, s = 0; (i = e[s]) != null; s++) {
      try {
        t(i).triggerHandler('remove')
      } catch (t) {
      }
    }
    n(e)
  }, t.widget = function (e, i, s) {
    var n; var o; var a; var r; var l = e.split('.')[0]
    e = e.split('.')[1], n = l + '-' + e, s || (s = i, i = t.Widget), t.expr[':'][n.toLowerCase()] = function (e) {
      return !!t.data(e, n)
    }, t[l] = t[l] || {}, o = t[l][e], a = t[l][e] = function (t, e) {
      return this._createWidget ? void (arguments.length && this._createWidget(t, e)) : new a(t, e)
    }, t.extend(a, o, {
      version: s.version,
      _proto: t.extend({}, s),
      _childConstructors: []
    }), (r = new i()).options = t.widget.extend({}, r.options), t.each(s, function (e, n) {
      t.isFunction(n) && (s[e] = (function () {
        var t = function () {
          return i.prototype[e].apply(this, arguments)
        }; var s = function (t) {
          return i.prototype[e].apply(this, t)
        }
        return function () {
          var e; var i = this._super; var o = this._superApply
          return this._super = t, this._superApply = s, e = n.apply(this, arguments), this._super = i, this._superApply = o, e
        }
      }()))
    }), a.prototype = t.widget.extend(r, { widgetEventPrefix: o ? r.widgetEventPrefix : e }, s, {
      constructor: a,
      namespace: l,
      widgetName: e,
      widgetBaseClass: n,
      widgetFullName: n
    }), o ? (t.each(o._childConstructors, function (e, i) {
      var s = i.prototype
      t.widget(s.namespace + '.' + s.widgetName, a, i._proto)
    }), delete o._childConstructors) : i._childConstructors.push(a), t.widget.bridge(e, a)
  }, t.widget.extend = function (i) {
    for (var n, o, a = s.call(arguments, 1), r = 0, l = a.length; r < l; r++) for (n in a[r]) o = a[r][n], a[r].hasOwnProperty(n) && o !== e && (t.isPlainObject(o) ? i[n] = t.isPlainObject(i[n]) ? t.widget.extend({}, i[n], o) : t.widget.extend({}, o) : i[n] = o)
    return i
  }, t.widget.bridge = function (i, n) {
    var o = n.prototype.widgetFullName || i
    t.fn[i] = function (a) {
      var r = typeof a === 'string'; var l = s.call(arguments, 1); var h = this
      return a = !r && l.length ? t.widget.extend.apply(null, [a].concat(l)) : a, r ? this.each(function () {
        var s; var n = t.data(this, o)
        return n ? t.isFunction(n[a]) && a.charAt(0) !== '_' ? (s = n[a].apply(n, l), s !== n && s !== e ? (h = s && s.jquery ? h.pushStack(s.get()) : s, !1) : void 0) : t.error("no such method '" + a + "' for " + i + ' widget instance') : t.error('cannot call methods on ' + i + " prior to initialization; attempted to call method '" + a + "'")
      }) : this.each(function () {
        var e = t.data(this, o)
        e ? e.option(a || {})._init() : t.data(this, o, new n(a, this))
      }), h
    }
  }, t.Widget = function () {
  }, t.Widget._childConstructors = [], t.Widget.prototype = {
    widgetName: 'widget',
    widgetEventPrefix: '',
    defaultElement: '<div>',
    options: { disabled: !1, create: null },
    _createWidget: function (e, s) {
      s = t(s || this.defaultElement || this)[0], this.element = t(s), this.uuid = i++, this.eventNamespace = '.' + this.widgetName + this.uuid, this.options = t.widget.extend({}, this.options, this._getCreateOptions(), e), this.bindings = t(), this.hoverable = t(), this.focusable = t(), s !== this && (t.data(s, this.widgetName, this), t.data(s, this.widgetFullName, this), this._on(!0, this.element, {
        remove: function (t) {
          t.target === s && this.destroy()
        }
      }), this.document = t(s.style ? s.ownerDocument : s.document || s), this.window = t(this.document[0].defaultView || this.document[0].parentWindow)), this._create(), this._trigger('create', null, this._getCreateEventData()), this._init()
    },
    _getCreateOptions: t.noop,
    _getCreateEventData: t.noop,
    _create: t.noop,
    _init: t.noop,
    destroy: function () {
      this._destroy(), this.element.unbind(this.eventNamespace).removeData(this.widgetName).removeData(this.widgetFullName).removeData(t.camelCase(this.widgetFullName)), this.widget().unbind(this.eventNamespace).removeAttr('aria-disabled').removeClass(this.widgetFullName + '-disabled ui-state-disabled'), this.bindings.unbind(this.eventNamespace), this.hoverable.removeClass('ui-state-hover'), this.focusable.removeClass('ui-state-focus')
    },
    _destroy: t.noop,
    widget: function () {
      return this.element
    },
    option: function (i, s) {
      var n; var o; var a; var r = i
      if (arguments.length === 0) return t.widget.extend({}, this.options)
      if (typeof i === 'string') {
        if (r = {}, n = i.split('.'), i = n.shift(), n.length) {
          for (o = r[i] = t.widget.extend({}, this.options[i]), a = 0; a < n.length - 1; a++) o[n[a]] = o[n[a]] || {}, o = o[n[a]]
          if (i = n.pop(), s === e) return o[i] === e ? null : o[i]
          o[i] = s
        } else {
          if (s === e) return this.options[i] === e ? null : this.options[i]
          r[i] = s
        }
      }
      return this._setOptions(r), this
    },
    _setOptions: function (t) {
      var e
      for (e in t) this._setOption(e, t[e])
      return this
    },
    _setOption: function (t, e) {
      return this.options[t] = e, t === 'disabled' && (this.widget().toggleClass(this.widgetFullName + '-disabled ui-state-disabled', !!e).attr('aria-disabled', e), this.hoverable.removeClass('ui-state-hover'), this.focusable.removeClass('ui-state-focus')), this
    },
    enable: function () {
      return this._setOption('disabled', !1)
    },
    disable: function () {
      return this._setOption('disabled', !0)
    },
    _on: function (e, i, s) {
      var n; var o = this
      typeof e !== 'boolean' && (s = i, i = e, e = !1), s ? (i = n = t(i), this.bindings = this.bindings.add(i)) : (s = i, i = this.element, n = this.widget()), t.each(s, function (s, a) {
        function r () {
          if (e || !0 !== o.options.disabled && !t(this).hasClass('ui-state-disabled')) return (typeof a === 'string' ? o[a] : a).apply(o, arguments)
        }

        typeof a !== 'string' && (r.guid = a.guid = a.guid || r.guid || t.guid++)
        var l = s.match(/^(\w+)\s*(.*)$/); var h = l[1] + o.eventNamespace; var c = l[2]
        c ? n.delegate(c, h, r) : i.bind(h, r)
      })
    },
    _off: function (t, e) {
      e = (e || '').split(' ').join(this.eventNamespace + ' ') + this.eventNamespace, t.unbind(e).undelegate(e)
    },
    _delay: function (t, e) {
      var i = this
      return setTimeout(function () {
        return (typeof t === 'string' ? i[t] : t).apply(i, arguments)
      }, e || 0)
    },
    _hoverable: function (e) {
      this.hoverable = this.hoverable.add(e), this._on(e, {
        mouseenter: function (e) {
          t(e.currentTarget).addClass('ui-state-hover')
        },
        mouseleave: function (e) {
          t(e.currentTarget).removeClass('ui-state-hover')
        }
      })
    },
    _focusable: function (e) {
      this.focusable = this.focusable.add(e), this._on(e, {
        focusin: function (e) {
          t(e.currentTarget).addClass('ui-state-focus')
        },
        focusout: function (e) {
          t(e.currentTarget).removeClass('ui-state-focus')
        }
      })
    },
    _trigger: function (e, i, s) {
      var n; var o; var a = this.options[e]
      if (s = s || {}, i = t.Event(i), i.type = (e === this.widgetEventPrefix ? e : this.widgetEventPrefix + e).toLowerCase(), i.target = this.element[0], o = i.originalEvent) for (n in o) n in i || (i[n] = o[n])
      return this.element.trigger(i, s), !(t.isFunction(a) && !1 === a.apply(this.element[0], [i].concat(s)) || i.isDefaultPrevented())
    }
  }, t.each({ show: 'fadeIn', hide: 'fadeOut' }, function (e, i) {
    t.Widget.prototype['_' + e] = function (s, n, o) {
      typeof n === 'string' && (n = { effect: n })
      var a; var r = n ? !0 === n || typeof n === 'number' ? i : n.effect || i : e
      typeof (n = n || {}) === 'number' && (n = { duration: n }), a = !t.isEmptyObject(n), n.complete = o, n.delay && s.delay(n.delay), a && t.effects && (t.effects.effect[r] || !1 !== t.uiBackCompat && t.effects[r]) ? s[e](n) : r !== e && s[r] ? s[r](n.duration, n.easing, o) : s.queue(function (i) {
        t(this)[e](), o && o.call(s[0]), i()
      })
    }
  }), !1 !== t.uiBackCompat && (t.Widget.prototype._getCreateOptions = function () {
    return t.metadata && t.metadata.get(this.element[0])[this.widgetName]
  })
}(jQuery)), (function (t, e) {
  var i = !1
  t(document).mouseup(function (t) {
    i = !1
  }), t.widget('ui.mouse', {
    version: '1.9.2',
    options: { cancel: 'input,textarea,button,select,option', distance: 1, delay: 0 },
    _mouseInit: function () {
      var e = this
      this.element.bind('mousedown.' + this.widgetName, function (t) {
        return e._mouseDown(t)
      }).bind('click.' + this.widgetName, function (i) {
        if (!0 === t.data(i.target, e.widgetName + '.preventClickEvent')) return t.removeData(i.target, e.widgetName + '.preventClickEvent'), i.stopImmediatePropagation(), !1
      }), this.started = !1
    },
    _mouseDestroy: function () {
      this.element.unbind('.' + this.widgetName), this._mouseMoveDelegate && t(document).unbind('mousemove.' + this.widgetName, this._mouseMoveDelegate).unbind('mouseup.' + this.widgetName, this._mouseUpDelegate)
    },
    _mouseDown: function (e) {
      if (!i) {
        this._mouseStarted && this._mouseUp(e), this._mouseDownEvent = e
        var s = this; var n = e.which === 1
        var o = !(typeof this.options.cancel !== 'string' || !e.target.nodeName) && t(e.target).closest(this.options.cancel).length
        return !(n && !o && this._mouseCapture(e) && (this.mouseDelayMet = !this.options.delay, this.mouseDelayMet || (this._mouseDelayTimer = setTimeout(function () {
          s.mouseDelayMet = !0
        }, this.options.delay)), this._mouseDistanceMet(e) && this._mouseDelayMet(e) && (this._mouseStarted = !1 !== this._mouseStart(e), !this._mouseStarted) ? (e.preventDefault(), 0) : (!0 === t.data(e.target, this.widgetName + '.preventClickEvent') && t.removeData(e.target, this.widgetName + '.preventClickEvent'), this._mouseMoveDelegate = function (t) {
          return s._mouseMove(t)
        }, this._mouseUpDelegate = function (t) {
          return s._mouseUp(t)
        }, t(document).bind('mousemove.' + this.widgetName, this._mouseMoveDelegate).bind('mouseup.' + this.widgetName, this._mouseUpDelegate), e.preventDefault(), i = !0, 0)))
      }
    },
    _mouseMove: function (e) {
      return !t.ui.ie || document.documentMode >= 9 || e.button ? this._mouseStarted ? (this._mouseDrag(e), e.preventDefault()) : (this._mouseDistanceMet(e) && this._mouseDelayMet(e) && (this._mouseStarted = !1 !== this._mouseStart(this._mouseDownEvent, e), this._mouseStarted ? this._mouseDrag(e) : this._mouseUp(e)), !this._mouseStarted) : this._mouseUp(e)
    },
    _mouseUp: function (e) {
      return t(document).unbind('mousemove.' + this.widgetName, this._mouseMoveDelegate).unbind('mouseup.' + this.widgetName, this._mouseUpDelegate), this._mouseStarted && (this._mouseStarted = !1, e.target === this._mouseDownEvent.target && t.data(e.target, this.widgetName + '.preventClickEvent', !0), this._mouseStop(e)), !1
    },
    _mouseDistanceMet: function (t) {
      return Math.max(Math.abs(this._mouseDownEvent.pageX - t.pageX), Math.abs(this._mouseDownEvent.pageY - t.pageY)) >= this.options.distance
    },
    _mouseDelayMet: function (t) {
      return this.mouseDelayMet
    },
    _mouseStart: function (t) {
    },
    _mouseDrag: function (t) {
    },
    _mouseStop: function (t) {
    },
    _mouseCapture: function (t) {
      return !0
    }
  })
}(jQuery)), (function (t, e) {
  function i (t, e, i) {
    return [parseInt(t[0], 10) * (d.test(t[0]) ? e / 100 : 1), parseInt(t[1], 10) * (d.test(t[1]) ? i / 100 : 1)]
  }

  function s (e, i) {
    return parseInt(t.css(e, i), 10) || 0
  }

  t.ui = t.ui || {}
  var n; var o = Math.max; var a = Math.abs; var r = Math.round; var l = /left|center|right/; var h = /top|center|bottom/
  var c = /[\+\-]\d+%?/; var u = /^\w+/; var d = /%$/; var p = t.fn.position
  t.position = {
    scrollbarWidth: function () {
      if (n !== e) return n
      var i; var s
      var o = t("<div style='display:block;width:50px;height:50px;overflow:hidden;'><div style='height:100px;width:auto;'></div></div>")
      var a = o.children()[0]
      return t('body').append(o), i = a.offsetWidth, o.css('overflow', 'scroll'), s = a.offsetWidth, i === s && (s = o[0].clientWidth), o.remove(), n = i - s
    },
    getScrollInfo: function (e) {
      var i = e.isWindow ? '' : e.element.css('overflow-x'); var s = e.isWindow ? '' : e.element.css('overflow-y')
      var n = i === 'scroll' || i === 'auto' && e.width < e.element[0].scrollWidth
      var o = s === 'scroll' || s === 'auto' && e.height < e.element[0].scrollHeight
      return { width: n ? t.position.scrollbarWidth() : 0, height: o ? t.position.scrollbarWidth() : 0 }
    },
    getWithinInfo: function (e) {
      var i = t(e || window); var s = t.isWindow(i[0])
      return {
        element: i,
        isWindow: s,
        offset: i.offset() || { left: 0, top: 0 },
        scrollLeft: i.scrollLeft(),
        scrollTop: i.scrollTop(),
        width: s ? i.width() : i.outerWidth(),
        height: s ? i.height() : i.outerHeight()
      }
    }
  }, t.fn.position = function (e) {
    if (!e || !e.of) return p.apply(this, arguments)
    e = t.extend({}, e)
    var n; var d; var f; var g; var m; var v = t(e.of); var b = t.position.getWithinInfo(e.within); var _ = t.position.getScrollInfo(b)
    var y = v[0]; var w = (e.collision || 'flip').split(' '); var C = {}
    return y.nodeType === 9 ? (d = v.width(), f = v.height(), g = {
      top: 0,
      left: 0
    }) : t.isWindow(y) ? (d = v.width(), f = v.height(), g = {
      top: v.scrollTop(),
      left: v.scrollLeft()
    }) : y.preventDefault ? (e.at = 'left top', d = f = 0, g = {
      top: y.pageY,
      left: y.pageX
    }) : (d = v.outerWidth(), f = v.outerHeight(), g = v.offset()), m = t.extend({}, g), t.each(['my', 'at'], function () {
      var t; var i; var s = (e[this] || '').split(' ')
      s.length === 1 && (s = l.test(s[0]) ? s.concat(['center']) : h.test(s[0]) ? ['center'].concat(s) : ['center', 'center']), s[0] = l.test(s[0]) ? s[0] : 'center', s[1] = h.test(s[1]) ? s[1] : 'center', t = c.exec(s[0]), i = c.exec(s[1]), C[this] = [t ? t[0] : 0, i ? i[0] : 0], e[this] = [u.exec(s[0])[0], u.exec(s[1])[0]]
    }), w.length === 1 && (w[1] = w[0]), e.at[0] === 'right' ? m.left += d : e.at[0] === 'center' && (m.left += d / 2), e.at[1] === 'bottom' ? m.top += f : e.at[1] === 'center' && (m.top += f / 2), n = i(C.at, d, f), m.left += n[0], m.top += n[1], this.each(function () {
      var l; var h; var c = t(this); var u = c.outerWidth(); var p = c.outerHeight(); var y = s(this, 'marginLeft')
      var x = s(this, 'marginTop'); var z = u + y + s(this, 'marginRight') + _.width
      var P = p + x + s(this, 'marginBottom') + _.height; var k = t.extend({}, m)
      var S = i(C.my, c.outerWidth(), c.outerHeight())
      e.my[0] === 'right' ? k.left -= u : e.my[0] === 'center' && (k.left -= u / 2), e.my[1] === 'bottom' ? k.top -= p : e.my[1] === 'center' && (k.top -= p / 2), k.left += S[0], k.top += S[1], t.support.offsetFractions || (k.left = r(k.left), k.top = r(k.top)), l = {
        marginLeft: y,
        marginTop: x
      }, t.each(['left', 'top'], function (i, s) {
        t.ui.position[w[i]] && t.ui.position[w[i]][s](k, {
          targetWidth: d,
          targetHeight: f,
          elemWidth: u,
          elemHeight: p,
          collisionPosition: l,
          collisionWidth: z,
          collisionHeight: P,
          offset: [n[0] + S[0], n[1] + S[1]],
          my: e.my,
          at: e.at,
          within: b,
          elem: c
        })
      }), t.fn.bgiframe && c.bgiframe(), e.using && (h = function (t) {
        var i = g.left - k.left; var s = i + d - u; var n = g.top - k.top; var r = n + f - p; var l = {
          target: { element: v, left: g.left, top: g.top, width: d, height: f },
          element: { element: c, left: k.left, top: k.top, width: u, height: p },
          horizontal: s < 0 ? 'left' : i > 0 ? 'right' : 'center',
          vertical: r < 0 ? 'top' : n > 0 ? 'bottom' : 'middle'
        }
        d < u && a(i + s) < d && (l.horizontal = 'center'), f < p && a(n + r) < f && (l.vertical = 'middle'), o(a(i), a(s)) > o(a(n), a(r)) ? l.important = 'horizontal' : l.important = 'vertical', e.using.call(this, t, l)
      }), c.offset(t.extend(k, { using: h }))
    })
  }, t.ui.position = {
    fit: {
      left: function (t, e) {
        var i; var s = e.within; var n = s.isWindow ? s.scrollLeft : s.offset.left; var a = s.width
        var r = t.left - e.collisionPosition.marginLeft; var l = n - r; var h = r + e.collisionWidth - a - n
        e.collisionWidth > a ? l > 0 && h <= 0 ? (i = t.left + l + e.collisionWidth - a - n, t.left += l - i) : t.left = h > 0 && l <= 0 ? n : l > h ? n + a - e.collisionWidth : n : l > 0 ? t.left += l : h > 0 ? t.left -= h : t.left = o(t.left - r, t.left)
      },
      top: function (t, e) {
        var i; var s = e.within; var n = s.isWindow ? s.scrollTop : s.offset.top; var a = e.within.height
        var r = t.top - e.collisionPosition.marginTop; var l = n - r; var h = r + e.collisionHeight - a - n
        e.collisionHeight > a ? l > 0 && h <= 0 ? (i = t.top + l + e.collisionHeight - a - n, t.top += l - i) : t.top = h > 0 && l <= 0 ? n : l > h ? n + a - e.collisionHeight : n : l > 0 ? t.top += l : h > 0 ? t.top -= h : t.top = o(t.top - r, t.top)
      }
    },
    flip: {
      left: function (t, e) {
        var i; var s; var n = e.within; var o = n.offset.left + n.scrollLeft; var r = n.width
        var l = n.isWindow ? n.scrollLeft : n.offset.left; var h = t.left - e.collisionPosition.marginLeft
        var c = h - l; var u = h + e.collisionWidth - r - l
        var d = e.my[0] === 'left' ? -e.elemWidth : e.my[0] === 'right' ? e.elemWidth : 0
        var p = e.at[0] === 'left' ? e.targetWidth : e.at[0] === 'right' ? -e.targetWidth : 0
        var f = -2 * e.offset[0]
        c < 0 ? ((i = t.left + d + p + f + e.collisionWidth - r - o) < 0 || i < a(c)) && (t.left += d + p + f) : u > 0 && ((s = t.left - e.collisionPosition.marginLeft + d + p + f - l) > 0 || a(s) < u) && (t.left += d + p + f)
      },
      top: function (t, e) {
        var i; var s; var n = e.within; var o = n.offset.top + n.scrollTop; var r = n.height
        var l = n.isWindow ? n.scrollTop : n.offset.top; var h = t.top - e.collisionPosition.marginTop; var c = h - l
        var u = h + e.collisionHeight - r - l
        var d = e.my[1] === 'top' ? -e.elemHeight : e.my[1] === 'bottom' ? e.elemHeight : 0
        var p = e.at[1] === 'top' ? e.targetHeight : e.at[1] === 'bottom' ? -e.targetHeight : 0
        var f = -2 * e.offset[1]
        c < 0 ? (s = t.top + d + p + f + e.collisionHeight - r - o, t.top + d + p + f > c && (s < 0 || s < a(c)) && (t.top += d + p + f)) : u > 0 && (i = t.top - e.collisionPosition.marginTop + d + p + f - l, t.top + d + p + f > u && (i > 0 || a(i) < u) && (t.top += d + p + f))
      }
    },
    flipfit: {
      left: function () {
        t.ui.position.flip.left.apply(this, arguments), t.ui.position.fit.left.apply(this, arguments)
      },
      top: function () {
        t.ui.position.flip.top.apply(this, arguments), t.ui.position.fit.top.apply(this, arguments)
      }
    }
  }, (function () {
    var e; var i; var s; var n; var o; var a = document.getElementsByTagName('body')[0]; var r = document.createElement('div')
    e = document.createElement(a ? 'div' : 'body'), s = {
      visibility: 'hidden',
      width: 0,
      height: 0,
      border: 0,
      margin: 0,
      background: 'none'
    }, a && t.extend(s, { position: 'absolute', left: '-1000px', top: '-1000px' })
    for (o in s) e.style[o] = s[o]
    e.appendChild(r), (i = a || document.documentElement).insertBefore(e, i.firstChild), r.style.cssText = 'position: absolute; left: 10.7432222px;', n = t(r).offset().left, t.support.offsetFractions = n > 10 && n < 11, e.innerHTML = '', i.removeChild(e)
  }()), !1 !== t.uiBackCompat && (function (t) {
    var i = t.fn.position
    t.fn.position = function (s) {
      if (!s || !s.offset) return i.call(this, s)
      var n = s.offset.split(' '); var o = s.at.split(' ')
      return n.length === 1 && (n[1] = n[0]), /^\d/.test(n[0]) && (n[0] = '+' + n[0]), /^\d/.test(n[1]) && (n[1] = '+' + n[1]), o.length === 1 && (/left|center|right/.test(o[0]) ? o[1] = 'center' : (o[1] = o[0], o[0] = 'center')), i.call(this, t.extend(s, {
        at: o[0] + n[0] + ' ' + o[1] + n[1],
        offset: e
      }))
    }
  }(jQuery))
}(jQuery)), (function (t, e) {
  var i; var s; var n; var o; var a = 'ui-button ui-widget ui-state-default ui-corner-all'
  var r = 'ui-button-icons-only ui-button-icon-only ui-button-text-icons ui-button-text-icon-primary ui-button-text-icon-secondary ui-button-text-only'
  var l = function () {
    var e = t(this).find(':ui-button')
    setTimeout(function () {
      e.button('refresh')
    }, 1)
  }; var h = function (e) {
    var i = e.name; var s = e.form; var n = t([])
    return i && (n = s ? t(s).find("[name='" + i + "']") : t("[name='" + i + "']", e.ownerDocument).filter(function () {
      return !this.form
    })), n
  }
}(jQuery)), (function (t, e) {
  var i = 'ui-dialog ui-widget ui-widget-content ui-corner-all '
  var s = { buttons: !0, height: !0, maxHeight: !0, maxWidth: !0, minHeight: !0, minWidth: !0, width: !0 }
  var n = { maxHeight: !0, maxWidth: !0, minHeight: !0, minWidth: !0 }
  t.widget('ui.dialog', {
    version: '1.9.2',
    options: {
      autoOpen: !0,
      buttons: {},
      closeOnEscape: !0,
      closeText: 'close',
      dialogClass: '',
      draggable: !0,
      hide: null,
      height: 'auto',
      maxHeight: !1,
      maxWidth: !1,
      minHeight: 150,
      minWidth: 150,
      modal: !1,
      position: {
        my: 'center',
        at: 'center',
        of: window,
        collision: 'fit',
        using: function (e) {
          var i = t(this).css(e).offset().top
          i < 0 && t(this).css('top', e.top - i)
        }
      },
      resizable: !0,
      show: null,
      stack: !0,
      title: '',
      width: 300,
      zIndex: 1e3
    },
    _create: function () {
      this.originalTitle = this.element.attr('title'), typeof this.originalTitle !== 'string' && (this.originalTitle = ''), this.oldPosition = {
        parent: this.element.parent(),
        index: this.element.parent().children().index(this.element)
      }, this.options.title = this.options.title || this.originalTitle
      var e; var s; var n; var o; var a; var r = this; var l = this.options; var h = l.title || '&#160;'
      e = (this.uiDialog = t('<div>')).addClass(i + l.dialogClass).css({
        display: 'none',
        outline: 0,
        zIndex: l.zIndex
      }).attr('tabIndex', -1).keydown(function (e) {
        l.closeOnEscape && !e.isDefaultPrevented() && e.keyCode && e.keyCode === t.ui.keyCode.ESCAPE && (r.close(e), e.preventDefault())
      }).mousedown(function (t) {
        r.moveToTop(!1, t)
      }).appendTo('body'), this.element.show().removeAttr('title').addClass('ui-dialog-content ui-widget-content').appendTo(e), s = (this.uiDialogTitlebar = t('<div>')).addClass('ui-dialog-titlebar  ui-widget-header  ui-corner-all  ui-helper-clearfix').bind('mousedown', function () {
        e.focus()
      }).prependTo(e), n = t("<a href='#'></a>").addClass('ui-dialog-titlebar-close  ui-corner-all').attr('role', 'button').click(function (t) {
        t.preventDefault(), r.close(t)
      }).appendTo(s), (this.uiDialogTitlebarCloseText = t('<span>')).addClass('ui-icon ui-icon-closethick').text(l.closeText).appendTo(n), o = t('<span>').uniqueId().addClass('ui-dialog-title').html(h).prependTo(s), a = (this.uiDialogButtonPane = t('<div>')).addClass('ui-dialog-buttonpane ui-widget-content ui-helper-clearfix'), (this.uiButtonSet = t('<div>')).addClass('ui-dialog-buttonset').appendTo(a), e.attr({
        role: 'dialog',
        'aria-labelledby': o.attr('id')
      }), s.find('*').add(s).disableSelection(), this._hoverable(n), this._focusable(n), l.draggable && t.fn.draggable && this._makeDraggable(), l.resizable && t.fn.resizable && this._makeResizable(), this._createButtons(l.buttons), this._isOpen = !1, t.fn.bgiframe && e.bgiframe(), this._on(e, {
        keydown: function (i) {
          if (l.modal && i.keyCode === t.ui.keyCode.TAB) {
            var s = t(':tabbable', e); var n = s.filter(':first'); var o = s.filter(':last')
            return i.target !== o[0] || i.shiftKey ? i.target === n[0] && i.shiftKey ? (o.focus(1), !1) : void 0 : (n.focus(1), !1)
          }
        }
      })
    },
    _init: function () {
      this.options.autoOpen && this.open()
    },
    _destroy: function () {
      var t; var e = this.oldPosition
      this.overlay && this.overlay.destroy(), this.uiDialog.hide(), this.element.removeClass('ui-dialog-content ui-widget-content').hide().appendTo('body'), this.uiDialog.remove(), this.originalTitle && this.element.attr('title', this.originalTitle), (t = e.parent.children().eq(e.index)).length && t[0] !== this.element[0] ? t.before(this.element) : e.parent.append(this.element)
    },
    widget: function () {
      return this.uiDialog
    },
    close: function (e) {
      var i; var s; var n = this
      if (this._isOpen && !1 !== this._trigger('beforeClose', e)) {
        return this._isOpen = !1, this.overlay && this.overlay.destroy(), this.options.hide ? this._hide(this.uiDialog, this.options.hide, function () {
          n._trigger('close', e)
        }) : (this.uiDialog.hide(), this._trigger('close', e)), t.ui.dialog.overlay.resize(), this.options.modal && (i = 0, t('.ui-dialog').each(function () {
          this !== n.uiDialog[0] && (s = t(this).css('z-index'), isNaN(s) || (i = Math.max(i, s)))
        }), t.ui.dialog.maxZ = i), this
      }
    },
    isOpen: function () {
      return this._isOpen
    },
    moveToTop: function (e, i) {
      var s; var n = this.options
      return n.modal && !e || !n.stack && !n.modal ? this._trigger('focus', i) : (n.zIndex > t.ui.dialog.maxZ && (t.ui.dialog.maxZ = n.zIndex), this.overlay && (t.ui.dialog.maxZ += 1, t.ui.dialog.overlay.maxZ = t.ui.dialog.maxZ, this.overlay.$el.css('z-index', t.ui.dialog.overlay.maxZ)), s = {
        scrollTop: this.element.scrollTop(),
        scrollLeft: this.element.scrollLeft()
      }, t.ui.dialog.maxZ += 1, this.uiDialog.css('z-index', t.ui.dialog.maxZ), this.element.attr(s), this._trigger('focus', i), this)
    },
    open: function () {
      if (!this._isOpen) {
        var e; var i = this.options; var s = this.uiDialog
        return this._size(), this._position(i.position), s.show(i.show), this.overlay = i.modal ? new t.ui.dialog.overlay(this) : null, this.moveToTop(!0), (e = this.element.find(':tabbable')).length || (e = this.uiDialogButtonPane.find(':tabbable')).length || (e = s), e.eq(0).focus(), this._isOpen = !0, this._trigger('open'), this
      }
    },
    _createButtons: function (e) {
      var i = this; var s = !1
      this.uiDialogButtonPane.remove(), this.uiButtonSet.empty(), typeof e === 'object' && e !== null && t.each(e, function () {
        return !(s = !0)
      }), s ? (t.each(e, function (e, s) {
        var n, o
        s = t.isFunction(s) ? {
          click: s,
          text: e
        } : s, s = t.extend({ type: 'button' }, s), o = s.click, s.click = function () {
          o.apply(i.element[0], arguments)
        }, n = t('<button></button>', s).appendTo(i.uiButtonSet), t.fn.button && n.button()
      }), this.uiDialog.addClass('ui-dialog-buttons'), this.uiDialogButtonPane.appendTo(this.uiDialog)) : this.uiDialog.removeClass('ui-dialog-buttons')
    },
    _makeDraggable: function () {
      function e (t) {
        return { position: t.position, offset: t.offset }
      }

      var i = this; var s = this.options
      this.uiDialog.draggable({
        cancel: '.ui-dialog-content, .ui-dialog-titlebar-close',
        handle: '.ui-dialog-titlebar',
        containment: 'document',
        start: function (s, n) {
          t(this).addClass('ui-dialog-dragging'), i._trigger('dragStart', s, e(n))
        },
        drag: function (t, s) {
          i._trigger('drag', t, e(s))
        },
        stop: function (n, o) {
          s.position = [o.position.left - i.document.scrollLeft(), o.position.top - i.document.scrollTop()], t(this).removeClass('ui-dialog-dragging'), i._trigger('dragStop', n, e(o)), t.ui.dialog.overlay.resize()
        }
      })
    },
    _makeResizable: function (e) {
      function i (t) {
        return {
          originalPosition: t.originalPosition,
          originalSize: t.originalSize,
          position: t.position,
          size: t.size
        }
      }

      e = void 0 === e ? this.options.resizable : e
      var s = this; var n = this.options; var o = this.uiDialog.css('position')
      var a = typeof e === 'string' ? e : 'n,e,s,w,se,sw,ne,nw'
      this.uiDialog.resizable({
        cancel: '.ui-dialog-content',
        containment: 'document',
        alsoResize: this.element,
        maxWidth: n.maxWidth,
        maxHeight: n.maxHeight,
        minWidth: n.minWidth,
        minHeight: this._minHeight(),
        handles: a,
        start: function (e, n) {
          t(this).addClass('ui-dialog-resizing'), s._trigger('resizeStart', e, i(n))
        },
        resize: function (t, e) {
          s._trigger('resize', t, i(e))
        },
        stop: function (e, o) {
          t(this).removeClass('ui-dialog-resizing'), n.height = t(this).height(), n.width = t(this).width(), s._trigger('resizeStop', e, i(o)), t.ui.dialog.overlay.resize()
        }
      }).css('position', o).find('.ui-resizable-se').addClass('ui-icon ui-icon-grip-diagonal-se')
    },
    _minHeight: function () {
      var t = this.options
      return t.height === 'auto' ? t.minHeight : Math.min(t.minHeight, t.height)
    },
    _position: function (e) {
      var i; var s = []; var n = [0, 0]
      e ? ((typeof e === 'string' || typeof e === 'object' && '0' in e) && ((s = e.split ? e.split(' ') : [e[0], e[1]]).length === 1 && (s[1] = s[0]), t.each(['left', 'top'], function (t, e) {
        +s[t] === s[t] && (n[t] = s[t], s[t] = e)
      }), e = {
        my: s[0] + (n[0] < 0 ? n[0] : '+' + n[0]) + ' ' + s[1] + (n[1] < 0 ? n[1] : '+' + n[1]),
        at: s.join(' ')
      }), e = t.extend({}, t.ui.dialog.prototype.options.position, e)) : e = t.ui.dialog.prototype.options.position, (i = this.uiDialog.is(':visible')) || this.uiDialog.show(), this.uiDialog.position(e), i || this.uiDialog.hide()
    },
    _setOptions: function (e) {
      var i = this; var o = {}; var a = !1
      t.each(e, function (t, e) {
        i._setOption(t, e), t in s && (a = !0), t in n && (o[t] = e)
      }), a && this._size(), this.uiDialog.is(':data(resizable)') && this.uiDialog.resizable('option', o)
    },
    _setOption: function (e, s) {
      var n; var o; var a = this.uiDialog
      switch (e) {
        case 'buttons':
          this._createButtons(s)
          break
        case 'closeText':
          this.uiDialogTitlebarCloseText.text('' + s)
          break
        case 'dialogClass':
          a.removeClass(this.options.dialogClass).addClass(i + s)
          break
        case 'disabled':
          s ? a.addClass('ui-dialog-disabled') : a.removeClass('ui-dialog-disabled')
          break
        case 'draggable':
          (n = a.is(':data(draggable)')) && !s && a.draggable('destroy'), !n && s && this._makeDraggable()
          break
        case 'position':
          this._position(s)
          break
        case 'resizable':
          (o = a.is(':data(resizable)')) && !s && a.resizable('destroy'), o && typeof s === 'string' && a.resizable('option', 'handles', s), o || !1 === s || this._makeResizable(s)
          break
        case 'title':
          t('.ui-dialog-title', this.uiDialogTitlebar).html('' + (s || '&#160;'))
      }
      this._super(e, s)
    },
    _size: function () {
      var e; var i; var s; var n = this.options; var o = this.uiDialog.is(':visible')
      this.element.show().css({
        width: 'auto',
        minHeight: 0,
        height: 0
      }), n.minWidth > n.width && (n.width = n.minWidth), e = this.uiDialog.css({
        height: 'auto',
        width: n.width
      }).outerHeight(), i = Math.max(0, n.minHeight - e), n.height === 'auto' ? t.support.minHeight ? this.element.css({
        minHeight: i,
        height: 'auto'
      }) : (this.uiDialog.show(), s = this.element.css('height', 'auto').height(), o || this.uiDialog.hide(), this.element.height(Math.max(s, i))) : this.element.height(Math.max(n.height - e, 0)), this.uiDialog.is(':data(resizable)') && this.uiDialog.resizable('option', 'minHeight', this._minHeight())
    }
  }), t.extend(t.ui.dialog, {
    uuid: 0,
    maxZ: 0,
    getTitleId: function (t) {
      var e = t.attr('id')
      return e || (this.uuid += 1, e = this.uuid), 'ui-dialog-title-' + e
    },
    overlay: function (e) {
      this.$el = t.ui.dialog.overlay.create(e)
    }
  }), t.extend(t.ui.dialog.overlay, {
    instances: [],
    oldInstances: [],
    maxZ: 0,
    events: t.map('focus,mousedown,mouseup,keydown,keypress,click'.split(','), function (t) {
      return t + '.dialog-overlay'
    }).join(' '),
    create: function (e) {
      this.instances.length === 0 && (setTimeout(function () {
        t.ui.dialog.overlay.instances.length && t(document).bind(t.ui.dialog.overlay.events, function (e) {
          if (t(e.target).zIndex() < t.ui.dialog.overlay.maxZ) return !1
        })
      }, 1), t(window).bind('resize.dialog-overlay', t.ui.dialog.overlay.resize))
      var i = this.oldInstances.pop() || t('<div>').addClass('ui-widget-overlay')
      return t(document).bind('keydown.dialog-overlay', function (s) {
        var n = t.ui.dialog.overlay.instances
        n.length !== 0 && n[n.length - 1] === i && e.options.closeOnEscape && !s.isDefaultPrevented() && s.keyCode && s.keyCode === t.ui.keyCode.ESCAPE && (e.close(s), s.preventDefault())
      }), i.appendTo(document.body).css({
        width: this.width(),
        height: this.height()
      }), t.fn.bgiframe && i.bgiframe(), this.instances.push(i), i
    },
    destroy: function (e) {
      var i = t.inArray(e, this.instances); var s = 0
      i !== -1 && this.oldInstances.push(this.instances.splice(i, 1)[0]), this.instances.length === 0 && t([document, window]).unbind('.dialog-overlay'), e.height(0).width(0).remove(), t.each(this.instances, function () {
        s = Math.max(s, this.css('z-index'))
      }), this.maxZ = s
    },
    height: function () {
      var e, i
      return t.ui.ie ? (e = Math.max(document.documentElement.scrollHeight, document.body.scrollHeight), i = Math.max(document.documentElement.offsetHeight, document.body.offsetHeight), e < i ? t(window).height() + 'px' : e + 'px') : t(document).height() + 'px'
    },
    width: function () {
      var e, i
      return t.ui.ie ? (e = Math.max(document.documentElement.scrollWidth, document.body.scrollWidth), i = Math.max(document.documentElement.offsetWidth, document.body.offsetWidth), e < i ? t(window).width() + 'px' : e + 'px') : t(document).width() + 'px'
    },
    resize: function () {
      var e = t([])
      t.each(t.ui.dialog.overlay.instances, function () {
        e = e.add(this)
      }), e.css({ width: 0, height: 0 }).css({
        width: t.ui.dialog.overlay.width(),
        height: t.ui.dialog.overlay.height()
      })
    }
  }), t.extend(t.ui.dialog.overlay.prototype, {
    destroy: function () {
      t.ui.dialog.overlay.destroy(this.$el)
    }
  })
}(jQuery)), (function (t, e) {
  t.widget('ui.draggable', t.ui.mouse, {
    version: '1.9.2',
    widgetEventPrefix: 'drag',
    options: {
      addClasses: !0,
      appendTo: 'parent',
      axis: !1,
      connectToSortable: !1,
      containment: !1,
      cursor: 'auto',
      cursorAt: !1,
      grid: !1,
      handle: !1,
      helper: 'original',
      iframeFix: !1,
      opacity: !1,
      refreshPositions: !1,
      revert: !1,
      revertDuration: 500,
      scope: 'default',
      scroll: !0,
      scrollSensitivity: 20,
      scrollSpeed: 20,
      snap: !1,
      snapMode: 'both',
      snapTolerance: 20,
      stack: !1,
      zIndex: !1
    },
    _create: function () {
      this.options.helper != 'original' || /^(?:r|a|f)/.test(this.element.css('position')) || (this.element[0].style.position = 'relative'), this.options.addClasses && this.element.addClass('ui-draggable'), this.options.disabled && this.element.addClass('ui-draggable-disabled'), this._mouseInit()
    },
    _destroy: function () {
      this.element.removeClass('ui-draggable ui-draggable-dragging ui-draggable-disabled'), this._mouseDestroy()
    },
    _mouseCapture: function (e) {
      var i = this.options
      return !(this.helper || i.disabled || t(e.target).is('.ui-resizable-handle') || (this.handle = this._getHandle(e), !this.handle || (t(!0 === i.iframeFix ? 'iframe' : i.iframeFix).each(function () {
        t('<div class="ui-draggable-iframeFix" style="background: #fff;"></div>').css({
          width: this.offsetWidth + 'px',
          height: this.offsetHeight + 'px',
          position: 'absolute',
          opacity: '0.001',
          zIndex: 1e3
        }).css(t(this).offset()).appendTo('body')
      }), 0)))
    },
    _mouseStart: function (e) {
      var i = this.options
      return this.helper = this._createHelper(e), this.helper.addClass('ui-draggable-dragging'), this._cacheHelperProportions(), t.ui.ddmanager && (t.ui.ddmanager.current = this), this._cacheMargins(), this.cssPosition = this.helper.css('position'), this.scrollParent = this.helper.scrollParent(), this.offset = this.positionAbs = this.element.offset(), this.offset = {
        top: this.offset.top - this.margins.top,
        left: this.offset.left - this.margins.left
      }, t.extend(this.offset, {
        click: { left: e.pageX - this.offset.left, top: e.pageY - this.offset.top },
        parent: this._getParentOffset(),
        relative: this._getRelativeOffset()
      }), this.originalPosition = this.position = this._generatePosition(e), this.originalPageX = e.pageX, this.originalPageY = e.pageY, i.cursorAt && this._adjustOffsetFromHelper(i.cursorAt), i.containment && this._setContainment(), !1 === this._trigger('start', e) ? (this._clear(), !1) : (this._cacheHelperProportions(), t.ui.ddmanager && !i.dropBehaviour && t.ui.ddmanager.prepareOffsets(this, e), this._mouseDrag(e, !0), t.ui.ddmanager && t.ui.ddmanager.dragStart(this, e), !0)
    },
    _mouseDrag: function (e, i) {
      if (this.position = this._generatePosition(e), this.positionAbs = this._convertPositionTo('absolute'), !i) {
        var s = this._uiHash()
        if (!1 === this._trigger('drag', e, s)) return this._mouseUp({}), !1
        this.position = s.position
      }
      return this.options.axis && this.options.axis == 'y' || (this.helper[0].style.left = this.position.left + 'px'), this.options.axis && this.options.axis == 'x' || (this.helper[0].style.top = this.position.top + 'px'), t.ui.ddmanager && t.ui.ddmanager.drag(this, e), !1
    },
    _mouseStop: function (e) {
      var i = !1
      t.ui.ddmanager && !this.options.dropBehaviour && (i = t.ui.ddmanager.drop(this, e)), this.dropped && (i = this.dropped, this.dropped = !1)
      for (var s = this.element[0], n = !1; s && (s = s.parentNode);) s == document && (n = !0)
      if (!n && this.options.helper === 'original') return !1
      if (this.options.revert == 'invalid' && !i || this.options.revert == 'valid' && i || !0 === this.options.revert || t.isFunction(this.options.revert) && this.options.revert.call(this.element, i)) {
        var o = this
        t(this.helper).animate(this.originalPosition, parseInt(this.options.revertDuration, 10), function () {
          !1 !== o._trigger('stop', e) && o._clear()
        })
      } else !1 !== this._trigger('stop', e) && this._clear()
      return !1
    },
    _mouseUp: function (e) {
      return t('div.ui-draggable-iframeFix').each(function () {
        this.parentNode.removeChild(this)
      }), t.ui.ddmanager && t.ui.ddmanager.dragStop(this, e), t.ui.mouse.prototype._mouseUp.call(this, e)
    },
    cancel: function () {
      return this.helper.is('.ui-draggable-dragging') ? this._mouseUp({}) : this._clear(), this
    },
    _getHandle: function (e) {
      var i = !this.options.handle || !t(this.options.handle, this.element).length
      return t(this.options.handle, this.element).find('*').addBack().each(function () {
        this == e.target && (i = !0)
      }), i
    },
    _createHelper: function (e) {
      var i = this.options
      var s = t.isFunction(i.helper) ? t(i.helper.apply(this.element[0], [e])) : i.helper == 'clone' ? this.element.clone().removeAttr('id') : this.element
      return s.parents('body').length || s.appendTo(i.appendTo == 'parent' ? this.element[0].parentNode : i.appendTo), s[0] == this.element[0] || /(fixed|absolute)/.test(s.css('position')) || s.css('position', 'absolute'), s
    },
    _adjustOffsetFromHelper: function (e) {
      typeof e === 'string' && (e = e.split(' ')), t.isArray(e) && (e = {
        left: +e[0],
        top: +e[1] || 0
      }), 'left' in e && (this.offset.click.left = e.left + this.margins.left), 'right' in e && (this.offset.click.left = this.helperProportions.width - e.right + this.margins.left), 'top' in e && (this.offset.click.top = e.top + this.margins.top), 'bottom' in e && (this.offset.click.top = this.helperProportions.height - e.bottom + this.margins.top)
    },
    _getParentOffset: function () {
      this.offsetParent = this.helper.offsetParent()
      var e = this.offsetParent.offset()
      return this.cssPosition == 'absolute' && this.scrollParent[0] != document && t.contains(this.scrollParent[0], this.offsetParent[0]) && (e.left += this.scrollParent.scrollLeft(), e.top += this.scrollParent.scrollTop()), (this.offsetParent[0] == document.body || this.offsetParent[0].tagName && this.offsetParent[0].tagName.toLowerCase() == 'html' && t.ui.ie) && (e = {
        top: 0,
        left: 0
      }), {
        top: e.top + (parseInt(this.offsetParent.css('borderTopWidth'), 10) || 0),
        left: e.left + (parseInt(this.offsetParent.css('borderLeftWidth'), 10) || 0)
      }
    },
    _getRelativeOffset: function () {
      if (this.cssPosition == 'relative') {
        var t = this.element.position()
        return {
          top: t.top - (parseInt(this.helper.css('top'), 10) || 0) + this.scrollParent.scrollTop(),
          left: t.left - (parseInt(this.helper.css('left'), 10) || 0) + this.scrollParent.scrollLeft()
        }
      }
      return { top: 0, left: 0 }
    },
    _cacheMargins: function () {
      this.margins = {
        left: parseInt(this.element.css('marginLeft'), 10) || 0,
        top: parseInt(this.element.css('marginTop'), 10) || 0,
        right: parseInt(this.element.css('marginRight'), 10) || 0,
        bottom: parseInt(this.element.css('marginBottom'), 10) || 0
      }
    },
    _cacheHelperProportions: function () {
      this.helperProportions = { width: this.helper.outerWidth(), height: this.helper.outerHeight() }
    },
    _setContainment: function () {
      var e = this.options
      if (e.containment == 'parent' && (e.containment = this.helper[0].parentNode), e.containment != 'document' && e.containment != 'window' || (this.containment = [e.containment == 'document' ? 0 : t(window).scrollLeft() - this.offset.relative.left - this.offset.parent.left, e.containment == 'document' ? 0 : t(window).scrollTop() - this.offset.relative.top - this.offset.parent.top, (e.containment == 'document' ? 0 : t(window).scrollLeft()) + t(e.containment == 'document' ? document : window).width() - this.helperProportions.width - this.margins.left, (e.containment == 'document' ? 0 : t(window).scrollTop()) + (t(e.containment == 'document' ? document : window).height() || document.body.parentNode.scrollHeight) - this.helperProportions.height - this.margins.top]), /^(document|window|parent)$/.test(e.containment) || e.containment.constructor == Array) e.containment.constructor == Array && (this.containment = e.containment); else {
        var i = t(e.containment); var s = i[0]
        if (!s) return
        var n = (i.offset(), t(s).css('overflow') != 'hidden')
        this.containment = [(parseInt(t(s).css('borderLeftWidth'), 10) || 0) + (parseInt(t(s).css('paddingLeft'), 10) || 0), (parseInt(t(s).css('borderTopWidth'), 10) || 0) + (parseInt(t(s).css('paddingTop'), 10) || 0), (n ? Math.max(s.scrollWidth, s.offsetWidth) : s.offsetWidth) - (parseInt(t(s).css('borderLeftWidth'), 10) || 0) - (parseInt(t(s).css('paddingRight'), 10) || 0) - this.helperProportions.width - this.margins.left - this.margins.right, (n ? Math.max(s.scrollHeight, s.offsetHeight) : s.offsetHeight) - (parseInt(t(s).css('borderTopWidth'), 10) || 0) - (parseInt(t(s).css('paddingBottom'), 10) || 0) - this.helperProportions.height - this.margins.top - this.margins.bottom], this.relative_container = i
      }
    },
    _convertPositionTo: function (e, i) {
      i || (i = this.position)
      var s = e == 'absolute' ? 1 : -1
      var n = (this.options, this.cssPosition != 'absolute' || this.scrollParent[0] != document && t.contains(this.scrollParent[0], this.offsetParent[0]) ? this.scrollParent : this.offsetParent)
      var o = /(html|body)/i.test(n[0].tagName)
      return {
        top: i.top + this.offset.relative.top * s + this.offset.parent.top * s - (this.cssPosition == 'fixed' ? -this.scrollParent.scrollTop() : o ? 0 : n.scrollTop()) * s,
        left: i.left + this.offset.relative.left * s + this.offset.parent.left * s - (this.cssPosition == 'fixed' ? -this.scrollParent.scrollLeft() : o ? 0 : n.scrollLeft()) * s
      }
    },
    _generatePosition: function (e) {
      var i = this.options
      var s = this.cssPosition != 'absolute' || this.scrollParent[0] != document && t.contains(this.scrollParent[0], this.offsetParent[0]) ? this.scrollParent : this.offsetParent
      var n = /(html|body)/i.test(s[0].tagName); var o = e.pageX; var a = e.pageY
      if (this.originalPosition) {
        var r
        if (this.containment) {
          if (this.relative_container) {
            var l = this.relative_container.offset()
            r = [this.containment[0] + l.left, this.containment[1] + l.top, this.containment[2] + l.left, this.containment[3] + l.top]
          } else r = this.containment
          e.pageX - this.offset.click.left < r[0] && (o = r[0] + this.offset.click.left), e.pageY - this.offset.click.top < r[1] && (a = r[1] + this.offset.click.top), e.pageX - this.offset.click.left > r[2] && (o = r[2] + this.offset.click.left), e.pageY - this.offset.click.top > r[3] && (a = r[3] + this.offset.click.top)
        }
        if (i.grid) {
          var h = i.grid[1] ? this.originalPageY + Math.round((a - this.originalPageY) / i.grid[1]) * i.grid[1] : this.originalPageY
          a = r && (h - this.offset.click.top < r[1] || h - this.offset.click.top > r[3]) ? h - this.offset.click.top < r[1] ? h + i.grid[1] : h - i.grid[1] : h
          var c = i.grid[0] ? this.originalPageX + Math.round((o - this.originalPageX) / i.grid[0]) * i.grid[0] : this.originalPageX
          o = r && (c - this.offset.click.left < r[0] || c - this.offset.click.left > r[2]) ? c - this.offset.click.left < r[0] ? c + i.grid[0] : c - i.grid[0] : c
        }
      }
      return {
        top: a - this.offset.click.top - this.offset.relative.top - this.offset.parent.top + (this.cssPosition == 'fixed' ? -this.scrollParent.scrollTop() : n ? 0 : s.scrollTop()),
        left: o - this.offset.click.left - this.offset.relative.left - this.offset.parent.left + (this.cssPosition == 'fixed' ? -this.scrollParent.scrollLeft() : n ? 0 : s.scrollLeft())
      }
    },
    _clear: function () {
      this.helper.removeClass('ui-draggable-dragging'), this.helper[0] == this.element[0] || this.cancelHelperRemoval || this.helper.remove(), this.helper = null, this.cancelHelperRemoval = !1
    },
    _trigger: function (e, i, s) {
      return s = s || this._uiHash(), t.ui.plugin.call(this, e, [i, s]), e == 'drag' && (this.positionAbs = this._convertPositionTo('absolute')), t.Widget.prototype._trigger.call(this, e, i, s)
    },
    plugins: {},
    _uiHash: function (t) {
      return {
        helper: this.helper,
        position: this.position,
        originalPosition: this.originalPosition,
        offset: this.positionAbs
      }
    }
  }), t.ui.plugin.add('draggable', 'connectToSortable', {
    start: function (e, i) {
      var s = t(this).data('draggable'); var n = s.options; var o = t.extend({}, i, { item: s.element })
      s.sortables = [], t(n.connectToSortable).each(function () {
        var i = t.data(this, 'sortable')
        i && !i.options.disabled && (s.sortables.push({
          instance: i,
          shouldRevert: i.options.revert
        }), i.refreshPositions(), i._trigger('activate', e, o))
      })
    },
    stop: function (e, i) {
      var s = t(this).data('draggable'); var n = t.extend({}, i, { item: s.element })
      t.each(s.sortables, function () {
        this.instance.isOver ? (this.instance.isOver = 0, s.cancelHelperRemoval = !0, this.instance.cancelHelperRemoval = !1, this.shouldRevert && (this.instance.options.revert = !0), this.instance._mouseStop(e), this.instance.options.helper = this.instance.options._helper, s.options.helper == 'original' && this.instance.currentItem.css({
          top: 'auto',
          left: 'auto'
        })) : (this.instance.cancelHelperRemoval = !1, this.instance._trigger('deactivate', e, n))
      })
    },
    drag: function (e, i) {
      var s = t(this).data('draggable'); var n = this
      t.each(s.sortables, function (o) {
        var a = !1; var r = this
        this.instance.positionAbs = s.positionAbs, this.instance.helperProportions = s.helperProportions, this.instance.offset.click = s.offset.click, this.instance._intersectsWith(this.instance.containerCache) && (a = !0, t.each(s.sortables, function () {
          return this.instance.positionAbs = s.positionAbs, this.instance.helperProportions = s.helperProportions, this.instance.offset.click = s.offset.click, this != r && this.instance._intersectsWith(this.instance.containerCache) && t.ui.contains(r.instance.element[0], this.instance.element[0]) && (a = !1), a
        })), a ? (this.instance.isOver || (this.instance.isOver = 1, this.instance.currentItem = t(n).clone().removeAttr('id').appendTo(this.instance.element).data('sortable-item', !0), this.instance.options._helper = this.instance.options.helper, this.instance.options.helper = function () {
          return i.helper[0]
        }, e.target = this.instance.currentItem[0], this.instance._mouseCapture(e, !0), this.instance._mouseStart(e, !0, !0), this.instance.offset.click.top = s.offset.click.top, this.instance.offset.click.left = s.offset.click.left, this.instance.offset.parent.left -= s.offset.parent.left - this.instance.offset.parent.left, this.instance.offset.parent.top -= s.offset.parent.top - this.instance.offset.parent.top, s._trigger('toSortable', e), s.dropped = this.instance.element, s.currentItem = s.element, this.instance.fromOutside = s), this.instance.currentItem && this.instance._mouseDrag(e)) : this.instance.isOver && (this.instance.isOver = 0, this.instance.cancelHelperRemoval = !0, this.instance.options.revert = !1, this.instance._trigger('out', e, this.instance._uiHash(this.instance)), this.instance._mouseStop(e, !0), this.instance.options.helper = this.instance.options._helper, this.instance.currentItem.remove(), this.instance.placeholder && this.instance.placeholder.remove(), s._trigger('fromSortable', e), s.dropped = !1)
      })
    }
  }), t.ui.plugin.add('draggable', 'cursor', {
    start: function (e, i) {
      var s = t('body'); var n = t(this).data('draggable').options
      s.css('cursor') && (n._cursor = s.css('cursor')), s.css('cursor', n.cursor)
    },
    stop: function (e, i) {
      var s = t(this).data('draggable').options
      s._cursor && t('body').css('cursor', s._cursor)
    }
  }), t.ui.plugin.add('draggable', 'opacity', {
    start: function (e, i) {
      var s = t(i.helper); var n = t(this).data('draggable').options
      s.css('opacity') && (n._opacity = s.css('opacity')), s.css('opacity', n.opacity)
    },
    stop: function (e, i) {
      var s = t(this).data('draggable').options
      s._opacity && t(i.helper).css('opacity', s._opacity)
    }
  }), t.ui.plugin.add('draggable', 'scroll', {
    start: function (e, i) {
      var s = t(this).data('draggable')
      s.scrollParent[0] != document && s.scrollParent[0].tagName != 'HTML' && (s.overflowOffset = s.scrollParent.offset())
    },
    drag: function (e, i) {
      var s = t(this).data('draggable'); var n = s.options; var o = !1
      s.scrollParent[0] != document && s.scrollParent[0].tagName != 'HTML' ? (n.axis && n.axis == 'x' || (s.overflowOffset.top + s.scrollParent[0].offsetHeight - e.pageY < n.scrollSensitivity ? s.scrollParent[0].scrollTop = o = s.scrollParent[0].scrollTop + n.scrollSpeed : e.pageY - s.overflowOffset.top < n.scrollSensitivity && (s.scrollParent[0].scrollTop = o = s.scrollParent[0].scrollTop - n.scrollSpeed)), n.axis && n.axis == 'y' || (s.overflowOffset.left + s.scrollParent[0].offsetWidth - e.pageX < n.scrollSensitivity ? s.scrollParent[0].scrollLeft = o = s.scrollParent[0].scrollLeft + n.scrollSpeed : e.pageX - s.overflowOffset.left < n.scrollSensitivity && (s.scrollParent[0].scrollLeft = o = s.scrollParent[0].scrollLeft - n.scrollSpeed))) : (n.axis && n.axis == 'x' || (e.pageY - t(document).scrollTop() < n.scrollSensitivity ? o = t(document).scrollTop(t(document).scrollTop() - n.scrollSpeed) : t(window).height() - (e.pageY - t(document).scrollTop()) < n.scrollSensitivity && (o = t(document).scrollTop(t(document).scrollTop() + n.scrollSpeed))), n.axis && n.axis == 'y' || (e.pageX - t(document).scrollLeft() < n.scrollSensitivity ? o = t(document).scrollLeft(t(document).scrollLeft() - n.scrollSpeed) : t(window).width() - (e.pageX - t(document).scrollLeft()) < n.scrollSensitivity && (o = t(document).scrollLeft(t(document).scrollLeft() + n.scrollSpeed)))), !1 !== o && t.ui.ddmanager && !n.dropBehaviour && t.ui.ddmanager.prepareOffsets(s, e)
    }
  }), t.ui.plugin.add('draggable', 'snap', {
    start: function (e, i) {
      var s = t(this).data('draggable'); var n = s.options
      s.snapElements = [], t(n.snap.constructor != String ? n.snap.items || ':data(draggable)' : n.snap).each(function () {
        var e = t(this); var i = e.offset()
        this != s.element[0] && s.snapElements.push({
          item: this,
          width: e.outerWidth(),
          height: e.outerHeight(),
          top: i.top,
          left: i.left
        })
      })
    },
    drag: function (e, i) {
      for (var s = t(this).data('draggable'), n = s.options, o = n.snapTolerance, a = i.offset.left, r = a + s.helperProportions.width, l = i.offset.top, h = l + s.helperProportions.height, c = s.snapElements.length - 1; c >= 0; c--) {
        var u = s.snapElements[c].left; var d = u + s.snapElements[c].width; var p = s.snapElements[c].top
        var f = p + s.snapElements[c].height
        if (u - o < a && a < d + o && p - o < l && l < f + o || u - o < a && a < d + o && p - o < h && h < f + o || u - o < r && r < d + o && p - o < l && l < f + o || u - o < r && r < d + o && p - o < h && h < f + o) {
          if (n.snapMode != 'inner') {
            var g = Math.abs(p - h) <= o; var m = Math.abs(f - l) <= o; var v = Math.abs(u - r) <= o
            var b = Math.abs(d - a) <= o
            g && (i.position.top = s._convertPositionTo('relative', {
              top: p - s.helperProportions.height,
              left: 0
            }).top - s.margins.top), m && (i.position.top = s._convertPositionTo('relative', {
              top: f,
              left: 0
            }).top - s.margins.top), v && (i.position.left = s._convertPositionTo('relative', {
              top: 0,
              left: u - s.helperProportions.width
            }).left - s.margins.left), b && (i.position.left = s._convertPositionTo('relative', {
              top: 0,
              left: d
            }).left - s.margins.left)
          }
          var _ = g || m || v || b
          if (n.snapMode != 'outer') {
            var g = Math.abs(p - l) <= o; var m = Math.abs(f - h) <= o; var v = Math.abs(u - a) <= o
            var b = Math.abs(d - r) <= o
            g && (i.position.top = s._convertPositionTo('relative', {
              top: p,
              left: 0
            }).top - s.margins.top), m && (i.position.top = s._convertPositionTo('relative', {
              top: f - s.helperProportions.height,
              left: 0
            }).top - s.margins.top), v && (i.position.left = s._convertPositionTo('relative', {
              top: 0,
              left: u
            }).left - s.margins.left), b && (i.position.left = s._convertPositionTo('relative', {
              top: 0,
              left: d - s.helperProportions.width
            }).left - s.margins.left)
          }
          !s.snapElements[c].snapping && (g || m || v || b || _) && s.options.snap.snap && s.options.snap.snap.call(s.element, e, t.extend(s._uiHash(), { snapItem: s.snapElements[c].item })), s.snapElements[c].snapping = g || m || v || b || _
        } else s.snapElements[c].snapping && s.options.snap.release && s.options.snap.release.call(s.element, e, t.extend(s._uiHash(), { snapItem: s.snapElements[c].item })), s.snapElements[c].snapping = !1
      }
    }
  }), t.ui.plugin.add('draggable', 'stack', {
    start: function (e, i) {
      var s = t(this).data('draggable').options; var n = t.makeArray(t(s.stack)).sort(function (e, i) {
        return (parseInt(t(e).css('zIndex'), 10) || 0) - (parseInt(t(i).css('zIndex'), 10) || 0)
      })
      if (n.length) {
        var o = parseInt(n[0].style.zIndex) || 0
        t(n).each(function (t) {
          this.style.zIndex = o + t
        }), this[0].style.zIndex = o + n.length
      }
    }
  }), t.ui.plugin.add('draggable', 'zIndex', {
    start: function (e, i) {
      var s = t(i.helper); var n = t(this).data('draggable').options
      s.css('zIndex') && (n._zIndex = s.css('zIndex')), s.css('zIndex', n.zIndex)
    },
    stop: function (e, i) {
      var s = t(this).data('draggable').options
      s._zIndex && t(i.helper).css('zIndex', s._zIndex)
    }
  })
}(jQuery)), (function (t, e) {
  t.widget('ui.droppable', {
    version: '1.9.2',
    widgetEventPrefix: 'drop',
    options: {
      accept: '*',
      activeClass: !1,
      addClasses: !0,
      greedy: !1,
      hoverClass: !1,
      scope: 'default',
      tolerance: 'intersect'
    },
    _create: function () {
      var e = this.options; var i = e.accept
      this.isover = 0, this.isout = 1, this.accept = t.isFunction(i) ? i : function (t) {
        return t.is(i)
      }, this.proportions = {
        width: this.element[0].offsetWidth,
        height: this.element[0].offsetHeight
      }, t.ui.ddmanager.droppables[e.scope] = t.ui.ddmanager.droppables[e.scope] || [], t.ui.ddmanager.droppables[e.scope].push(this), e.addClasses && this.element.addClass('ui-droppable')
    },
    _destroy: function () {
      for (var e = t.ui.ddmanager.droppables[this.options.scope], i = 0; i < e.length; i++) e[i] == this && e.splice(i, 1)
      this.element.removeClass('ui-droppable ui-droppable-disabled')
    },
    _setOption: function (e, i) {
      e == 'accept' && (this.accept = t.isFunction(i) ? i : function (t) {
        return t.is(i)
      }), t.Widget.prototype._setOption.apply(this, arguments)
    },
    _activate: function (e) {
      var i = t.ui.ddmanager.current
      this.options.activeClass && this.element.addClass(this.options.activeClass), i && this._trigger('activate', e, this.ui(i))
    },
    _deactivate: function (e) {
      var i = t.ui.ddmanager.current
      this.options.activeClass && this.element.removeClass(this.options.activeClass), i && this._trigger('deactivate', e, this.ui(i))
    },
    _over: function (e) {
      var i = t.ui.ddmanager.current
      i && (i.currentItem || i.element)[0] != this.element[0] && this.accept.call(this.element[0], i.currentItem || i.element) && (this.options.hoverClass && this.element.addClass(this.options.hoverClass), this._trigger('over', e, this.ui(i)))
    },
    _out: function (e) {
      var i = t.ui.ddmanager.current
      i && (i.currentItem || i.element)[0] != this.element[0] && this.accept.call(this.element[0], i.currentItem || i.element) && (this.options.hoverClass && this.element.removeClass(this.options.hoverClass), this._trigger('out', e, this.ui(i)))
    },
    _drop: function (e, i) {
      var s = i || t.ui.ddmanager.current
      if (!s || (s.currentItem || s.element)[0] == this.element[0]) return !1
      var n = !1
      return this.element.find(':data(droppable)').not('.ui-draggable-dragging').each(function () {
        var e = t.data(this, 'droppable')
        if (e.options.greedy && !e.options.disabled && e.options.scope == s.options.scope && e.accept.call(e.element[0], s.currentItem || s.element) && t.ui.intersect(s, t.extend(e, { offset: e.element.offset() }), e.options.tolerance)) return n = !0, !1
      }), !n && !!this.accept.call(this.element[0], s.currentItem || s.element) && (this.options.activeClass && this.element.removeClass(this.options.activeClass), this.options.hoverClass && this.element.removeClass(this.options.hoverClass), this._trigger('drop', e, this.ui(s)), this.element)
    },
    ui: function (t) {
      return {
        draggable: t.currentItem || t.element,
        helper: t.helper,
        position: t.position,
        offset: t.positionAbs
      }
    }
  }), t.ui.intersect = function (e, i, s) {
    if (!i.offset) return !1
    var n = (e.positionAbs || e.position.absolute).left; var o = n + e.helperProportions.width
    var a = (e.positionAbs || e.position.absolute).top; var r = a + e.helperProportions.height; var l = i.offset.left
    var h = l + i.proportions.width; var c = i.offset.top; var u = c + i.proportions.height
    switch (s) {
      case 'fit':
        return l <= n && o <= h && c <= a && r <= u
      case 'intersect':
        return l < n + e.helperProportions.width / 2 && o - e.helperProportions.width / 2 < h && c < a + e.helperProportions.height / 2 && r - e.helperProportions.height / 2 < u
      case 'pointer':
        var d = (e.positionAbs || e.position.absolute).left + (e.clickOffset || e.offset.click).left
        var p = (e.positionAbs || e.position.absolute).top + (e.clickOffset || e.offset.click).top
        return t.ui.isOver(p, d, c, l, i.proportions.height, i.proportions.width)
      case 'touch':
        return (a >= c && a <= u || r >= c && r <= u || a < c && r > u) && (n >= l && n <= h || o >= l && o <= h || n < l && o > h)
      default:
        return !1
    }
  }, t.ui.ddmanager = {
    current: null,
    droppables: { default: [] },
    prepareOffsets: function (e, i) {
      var s = t.ui.ddmanager.droppables[e.options.scope] || []; var n = i ? i.type : null
      var o = (e.currentItem || e.element).find(':data(droppable)').addBack()
      t:for (var a = 0; a < s.length; a++) {
        if (!(s[a].options.disabled || e && !s[a].accept.call(s[a].element[0], e.currentItem || e.element))) {
          for (var r = 0; r < o.length; r++) {
            if (o[r] == s[a].element[0]) {
              s[a].proportions.height = 0
              continue t
            }
          }
          s[a].visible = s[a].element.css('display') != 'none', s[a].visible && (n == 'mousedown' && s[a]._activate.call(s[a], i), s[a].offset = s[a].element.offset(), s[a].proportions = {
            width: s[a].element[0].offsetWidth,
            height: s[a].element[0].offsetHeight
          })
        }
      }
    },
    drop: function (e, i) {
      var s = !1
      return t.each(t.ui.ddmanager.droppables[e.options.scope] || [], function () {
        this.options && (!this.options.disabled && this.visible && t.ui.intersect(e, this, this.options.tolerance) && (s = this._drop.call(this, i) || s), !this.options.disabled && this.visible && this.accept.call(this.element[0], e.currentItem || e.element) && (this.isout = 1, this.isover = 0, this._deactivate.call(this, i)))
      }), s
    },
    dragStart: function (e, i) {
      e.element.parentsUntil('body').bind('scroll.droppable', function () {
        e.options.refreshPositions || t.ui.ddmanager.prepareOffsets(e, i)
      })
    },
    drag: function (e, i) {
      e.options.refreshPositions && t.ui.ddmanager.prepareOffsets(e, i), t.each(t.ui.ddmanager.droppables[e.options.scope] || [], function () {
        if (!this.options.disabled && !this.greedyChild && this.visible) {
          var s = t.ui.intersect(e, this, this.options.tolerance)
          var n = s || this.isover != 1 ? s && this.isover == 0 ? 'isover' : null : 'isout'
          if (n) {
            var o
            if (this.options.greedy) {
              var a = this.options.scope
              var r = this.element.parents(':data(droppable)').filter(function () {
                return t.data(this, 'droppable').options.scope === a
              })
              r.length && (o = t.data(r[0], 'droppable'), o.greedyChild = n == 'isover' ? 1 : 0)
            }
            o && n == 'isover' && (o.isover = 0, o.isout = 1, o._out.call(o, i)), this[n] = 1, this[n == 'isout' ? 'isover' : 'isout'] = 0, this[n == 'isover' ? '_over' : '_out'].call(this, i), o && n == 'isout' && (o.isout = 0, o.isover = 1, o._over.call(o, i))
          }
        }
      })
    },
    dragStop: function (e, i) {
      e.element.parentsUntil('body').unbind('scroll.droppable'), e.options.refreshPositions || t.ui.ddmanager.prepareOffsets(e, i)
    }
  }
}(jQuery)), jQuery.effects || (function (t, e) {
  var i = !1 !== t.uiBackCompat; var s = 'ui-effects-'
  t.effects = { effect: {} }, (function (e, i) {
    function s (t, e, i) {
      var s = d[e.type] || {}
      return t == null ? i || !e.def ? null : e.def : (t = s.floor ? ~~t : parseFloat(t), isNaN(t) ? e.def : s.mod ? (t + s.mod) % s.mod : t < 0 ? 0 : s.max < t ? s.max : t)
    }

    function n (t) {
      var i = c(); var s = i._rgba = []
      return t = t.toLowerCase(), g(h, function (e, n) {
        var o; var a = n.re.exec(t); var r = a && n.parse(a); var l = n.space || 'rgba'
        if (r) return o = i[l](r), i[u[l].cache] = o[u[l].cache], s = i._rgba = o._rgba, !1
      }), s.length ? (s.join() === '0,0,0,0' && e.extend(s, a.transparent), i) : a[t]
    }

    function o (t, e, i) {
      return i = (i + 1) % 1, 6 * i < 1 ? t + (e - t) * i * 6 : 2 * i < 1 ? e : 3 * i < 2 ? t + (e - t) * (2 / 3 - i) * 6 : t
    }

    var a
    var r = 'backgroundColor borderBottomColor borderLeftColor borderRightColor borderTopColor color columnRuleColor outlineColor textDecorationColor textEmphasisColor'.split(' ')
    var l = /^([\-+])=\s*(\d+\.?\d*)/; var h = [{
      re: /rgba?\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*(?:,\s*(\d+(?:\.\d+)?)\s*)?\)/,
      parse: function (t) {
        return [t[1], t[2], t[3], t[4]]
      }
    }, {
      re: /rgba?\(\s*(\d+(?:\.\d+)?)\%\s*,\s*(\d+(?:\.\d+)?)\%\s*,\s*(\d+(?:\.\d+)?)\%\s*(?:,\s*(\d+(?:\.\d+)?)\s*)?\)/,
      parse: function (t) {
        return [2.55 * t[1], 2.55 * t[2], 2.55 * t[3], t[4]]
      }
    }, {
      re: /#([a-f0-9]{2})([a-f0-9]{2})([a-f0-9]{2})/,
      parse: function (t) {
        return [parseInt(t[1], 16), parseInt(t[2], 16), parseInt(t[3], 16)]
      }
    }, {
      re: /#([a-f0-9])([a-f0-9])([a-f0-9])/,
      parse: function (t) {
        return [parseInt(t[1] + t[1], 16), parseInt(t[2] + t[2], 16), parseInt(t[3] + t[3], 16)]
      }
    }, {
      re: /hsla?\(\s*(\d+(?:\.\d+)?)\s*,\s*(\d+(?:\.\d+)?)\%\s*,\s*(\d+(?:\.\d+)?)\%\s*(?:,\s*(\d+(?:\.\d+)?)\s*)?\)/,
      space: 'hsla',
      parse: function (t) {
        return [t[1], t[2] / 100, t[3] / 100, t[4]]
      }
    }]; var c = e.Color = function (t, i, s, n) {
      return new e.Color.fn.parse(t, i, s, n)
    }; var u = {
      rgba: {
        props: {
          red: { idx: 0, type: 'byte' },
          green: { idx: 1, type: 'byte' },
          blue: { idx: 2, type: 'byte' }
        }
      },
      hsla: {
        props: {
          hue: { idx: 0, type: 'degrees' },
          saturation: { idx: 1, type: 'percent' },
          lightness: { idx: 2, type: 'percent' }
        }
      }
    }; var d = { byte: { floor: !0, max: 255 }, percent: { max: 1 }, degrees: { mod: 360, floor: !0 } }; var p = c.support = {}
    var f = e('<p>')[0]; var g = e.each
    f.style.cssText = 'background-color:rgba(1,1,1,.5)', p.rgba = f.style.backgroundColor.indexOf('rgba') > -1, g(u, function (t, e) {
      e.cache = '_' + t, e.props.alpha = { idx: 3, type: 'percent', def: 1 }
    }), c.fn = e.extend(c.prototype, {
      parse: function (o, r, l, h) {
        if (o === i) return this._rgba = [null, null, null, null], this;
        (o.jquery || o.nodeType) && (o = e(o).css(r), r = i)
        var d = this; var p = e.type(o); var f = this._rgba = []
        return r !== i && (o = [o, r, l, h], p = 'array'), p === 'string' ? this.parse(n(o) || a._default) : p === 'array' ? (g(u.rgba.props, function (t, e) {
          f[e.idx] = s(o[e.idx], e)
        }), this) : p === 'object' ? (o instanceof c ? g(u, function (t, e) {
          o[e.cache] && (d[e.cache] = o[e.cache].slice())
        }) : g(u, function (e, i) {
          var n = i.cache
          g(i.props, function (t, e) {
            if (!d[n] && i.to) {
              if (t === 'alpha' || o[t] == null) return
              d[n] = i.to(d._rgba)
            }
            d[n][e.idx] = s(o[t], e, !0)
          }), d[n] && t.inArray(null, d[n].slice(0, 3)) < 0 && (d[n][3] = 1, i.from && (d._rgba = i.from(d[n])))
        }), this) : void 0
      },
      is: function (t) {
        var e = c(t); var i = !0; var s = this
        return g(u, function (t, n) {
          var o; var a = e[n.cache]
          return a && (o = s[n.cache] || n.to && n.to(s._rgba) || [], g(n.props, function (t, e) {
            if (a[e.idx] != null) return i = a[e.idx] === o[e.idx]
          })), i
        }), i
      },
      _space: function () {
        var t = []; var e = this
        return g(u, function (i, s) {
          e[s.cache] && t.push(i)
        }), t.pop()
      },
      transition: function (t, e) {
        var i = c(t); var n = i._space(); var o = u[n]; var a = this.alpha() === 0 ? c('transparent') : this
        var r = a[o.cache] || o.to(a._rgba); var l = r.slice()
        return i = i[o.cache], g(o.props, function (t, n) {
          var o = n.idx; var a = r[o]; var h = i[o]; var c = d[n.type] || {}
          h !== null && (a === null ? l[o] = h : (c.mod && (h - a > c.mod / 2 ? a += c.mod : a - h > c.mod / 2 && (a -= c.mod)), l[o] = s((h - a) * e + a, n)))
        }), this[n](l)
      },
      blend: function (t) {
        if (this._rgba[3] === 1) return this
        var i = this._rgba.slice(); var s = i.pop(); var n = c(t)._rgba
        return c(e.map(i, function (t, e) {
          return (1 - s) * n[e] + s * t
        }))
      },
      toRgbaString: function () {
        var t = 'rgba('; var i = e.map(this._rgba, function (t, e) {
          return t == null ? e > 2 ? 1 : 0 : t
        })
        return i[3] === 1 && (i.pop(), t = 'rgb('), t + i.join() + ')'
      },
      toHslaString: function () {
        var t = 'hsla('; var i = e.map(this.hsla(), function (t, e) {
          return t == null && (t = e > 2 ? 1 : 0), e && e < 3 && (t = Math.round(100 * t) + '%'), t
        })
        return i[3] === 1 && (i.pop(), t = 'hsl('), t + i.join() + ')'
      },
      toHexString: function (t) {
        var i = this._rgba.slice(); var s = i.pop()
        return t && i.push(~~(255 * s)), '#' + e.map(i, function (t) {
          return t = (t || 0).toString(16), t.length === 1 ? '0' + t : t
        }).join('')
      },
      toString: function () {
        return this._rgba[3] === 0 ? 'transparent' : this.toRgbaString()
      }
    }), c.fn.parse.prototype = c.fn, u.hsla.to = function (t) {
      if (t[0] == null || t[1] == null || t[2] == null) return [null, null, null, t[3]]
      var e; var i; var s = t[0] / 255; var n = t[1] / 255; var o = t[2] / 255; var a = t[3]; var r = Math.max(s, n, o)
      var l = Math.min(s, n, o); var h = r - l; var c = r + l; var u = 0.5 * c
      return e = l === r ? 0 : s === r ? 60 * (n - o) / h + 360 : n === r ? 60 * (o - s) / h + 120 : 60 * (s - n) / h + 240, i = u === 0 || u === 1 ? u : u <= 0.5 ? h / c : h / (2 - c), [Math.round(e) % 360, i, u, a == null ? 1 : a]
    }, u.hsla.from = function (t) {
      if (t[0] == null || t[1] == null || t[2] == null) return [null, null, null, t[3]]
      var e = t[0] / 360; var i = t[1]; var s = t[2]; var n = t[3]; var a = s <= 0.5 ? s * (1 + i) : s + i - s * i; var r = 2 * s - a
      return [Math.round(255 * o(r, a, e + 1 / 3)), Math.round(255 * o(r, a, e)), Math.round(255 * o(r, a, e - 1 / 3)), n]
    }, g(u, function (t, n) {
      var o = n.props; var a = n.cache; var r = n.to; var h = n.from
      c.fn[t] = function (t) {
        if (r && !this[a] && (this[a] = r(this._rgba)), t === i) return this[a].slice()
        var n; var l = e.type(t); var u = l === 'array' || l === 'object' ? t : arguments; var d = this[a].slice()
        return g(o, function (t, e) {
          var i = u[l === 'object' ? t : e.idx]
          i == null && (i = d[e.idx]), d[e.idx] = s(i, e)
        }), h ? (n = c(h(d)), n[a] = d, n) : c(d)
      }, g(o, function (i, s) {
        c.fn[i] || (c.fn[i] = function (n) {
          var o; var a = e.type(n); var r = i === 'alpha' ? this._hsla ? 'hsla' : 'rgba' : t; var h = this[r]()
          var c = h[s.idx]
          return a === 'undefined' ? c : (a === 'function' && (n = n.call(this, c), a = e.type(n)), n == null && s.empty ? this : (a === 'string' && (o = l.exec(n)) && (n = c + parseFloat(o[2]) * (o[1] === '+' ? 1 : -1)), h[s.idx] = n, this[r](h)))
        })
      })
    }), g(r, function (t, i) {
      e.cssHooks[i] = {
        set: function (t, s) {
          var o; var a; var r = ''
          if (e.type(s) !== 'string' || (o = n(s))) {
            if (s = c(o || s), !p.rgba && s._rgba[3] !== 1) {
              for (a = i === 'backgroundColor' ? t.parentNode : t; (r === '' || r === 'transparent') && a && a.style;) {
                try {
                  r = e.css(a, 'backgroundColor'), a = a.parentNode
                } catch (t) {
                }
              }
              s = s.blend(r && r !== 'transparent' ? r : '_default')
            }
            s = s.toRgbaString()
          }
          try {
            t.style[i] = s
          } catch (t) {
          }
        }
      }, e.fx.step[i] = function (t) {
        t.colorInit || (t.start = c(t.elem, i), t.end = c(t.end), t.colorInit = !0), e.cssHooks[i].set(t.elem, t.start.transition(t.end, t.pos))
      }
    }), e.cssHooks.borderColor = {
      expand: function (t) {
        var e = {}
        return g(['Top', 'Right', 'Bottom', 'Left'], function (i, s) {
          e['border' + s + 'Color'] = t
        }), e
      }
    }, a = e.Color.names = {
      aqua: '#00ffff',
      black: '#000000',
      blue: '#0000ff',
      fuchsia: '#ff00ff',
      gray: '#808080',
      green: '#008000',
      lime: '#00ff00',
      maroon: '#800000',
      navy: '#000080',
      olive: '#808000',
      purple: '#800080',
      red: '#ff0000',
      silver: '#c0c0c0',
      teal: '#008080',
      white: '#ffffff',
      yellow: '#ffff00',
      transparent: [null, null, null, 0],
      _default: '#ffffff'
    }
  }(jQuery)), (function () {
    function i () {
      var e; var i
      var s = this.ownerDocument.defaultView ? this.ownerDocument.defaultView.getComputedStyle(this, null) : this.currentStyle
      var n = {}
      if (s && s.length && s[0] && s[s[0]]) for (i = s.length; i--;) e = s[i], typeof s[e] === 'string' && (n[t.camelCase(e)] = s[e]); else for (e in s) typeof s[e] === 'string' && (n[e] = s[e])
      return n
    }

    function s (e, i) {
      var s; var n; var a = {}
      for (s in i) n = i[s], e[s] !== n && (o[s] || !t.fx.step[s] && isNaN(parseFloat(n)) || (a[s] = n))
      return a
    }

    var n = ['add', 'remove', 'toggle']; var o = {
      border: 1,
      borderBottom: 1,
      borderColor: 1,
      borderLeft: 1,
      borderRight: 1,
      borderTop: 1,
      borderWidth: 1,
      margin: 1,
      padding: 1
    }
    t.each(['borderLeftStyle', 'borderRightStyle', 'borderBottomStyle', 'borderTopStyle'], function (e, i) {
      t.fx.step[i] = function (t) {
        (t.end !== 'none' && !t.setAttr || t.pos === 1 && !t.setAttr) && (jQuery.style(t.elem, i, t.end), t.setAttr = !0)
      }
    }), t.effects.animateClass = function (e, o, a, r) {
      var l = t.speed(o, a, r)
      return this.queue(function () {
        var o; var a = t(this); var r = a.attr('class') || ''; var h = l.children ? a.find('*').addBack() : a
        h = h.map(function () {
          return { el: t(this), start: i.call(this) }
        }), (o = function () {
          t.each(n, function (t, i) {
            e[i] && a[i + 'Class'](e[i])
          })
        })(), h = h.map(function () {
          return this.end = i.call(this.el[0]), this.diff = s(this.start, this.end), this
        }), a.attr('class', r), h = h.map(function () {
          var e = this; var i = t.Deferred(); var s = jQuery.extend({}, l, {
            queue: !1,
            complete: function () {
              i.resolve(e)
            }
          })
          return this.el.animate(this.diff, s), i.promise()
        }), t.when.apply(t, h.get()).done(function () {
          o(), t.each(arguments, function () {
            var e = this.el
            t.each(this.diff, function (t) {
              e.css(t, '')
            })
          }), l.complete.call(a[0])
        })
      })
    }, t.fn.extend({
      _addClass: t.fn.addClass,
      addClass: function (e, i, s, n) {
        return i ? t.effects.animateClass.call(this, { add: e }, i, s, n) : this._addClass(e)
      },
      _removeClass: t.fn.removeClass,
      removeClass: function (e, i, s, n) {
        return i ? t.effects.animateClass.call(this, { remove: e }, i, s, n) : this._removeClass(e)
      },
      _toggleClass: t.fn.toggleClass,
      toggleClass: function (i, s, n, o, a) {
        return typeof s === 'boolean' || s === e ? n ? t.effects.animateClass.call(this, s ? { add: i } : { remove: i }, n, o, a) : this._toggleClass(i, s) : t.effects.animateClass.call(this, { toggle: i }, s, n, o)
      },
      switchClass: function (e, i, s, n, o) {
        return t.effects.animateClass.call(this, { add: i, remove: e }, s, n, o)
      }
    })
  }()), (function () {
    function n (e, i, s, n) {
      return t.isPlainObject(e) && (i = e, e = e.effect), e = { effect: e }, i == null && (i = {}), t.isFunction(i) && (n = i, s = null, i = {}), (typeof i === 'number' || t.fx.speeds[i]) && (n = s, s = i, i = {}), t.isFunction(s) && (n = s, s = null), i && t.extend(e, i), s = s || i.duration, e.duration = t.fx.off ? 0 : typeof s === 'number' ? s : s in t.fx.speeds ? t.fx.speeds[s] : t.fx.speeds._default, e.complete = n || i.complete, e
    }

    function o (e) {
      return !(e && typeof e !== 'number' && !t.fx.speeds[e] && (typeof e !== 'string' || t.effects.effect[e] || i && t.effects[e]))
    }

    t.extend(t.effects, {
      version: '1.9.2',
      save: function (t, e) {
        for (var i = 0; i < e.length; i++) e[i] !== null && t.data(s + e[i], t[0].style[e[i]])
      },
      restore: function (t, i) {
        var n, o
        for (o = 0; o < i.length; o++) i[o] !== null && ((n = t.data(s + i[o])) === e && (n = ''), t.css(i[o], n))
      },
      setMode: function (t, e) {
        return e === 'toggle' && (e = t.is(':hidden') ? 'show' : 'hide'), e
      },
      getBaseline: function (t, e) {
        var i, s
        switch (t[0]) {
          case 'top':
            i = 0
            break
          case 'middle':
            i = 0.5
            break
          case 'bottom':
            i = 1
            break
          default:
            i = t[0] / e.height
        }
        switch (t[1]) {
          case 'left':
            s = 0
            break
          case 'center':
            s = 0.5
            break
          case 'right':
            s = 1
            break
          default:
            s = t[1] / e.width
        }
        return { x: s, y: i }
      },
      createWrapper: function (e) {
        if (e.parent().is('.ui-effects-wrapper')) return e.parent()
        var i = { width: e.outerWidth(!0), height: e.outerHeight(!0), float: e.css('float') }
        var s = t('<div></div>').addClass('ui-effects-wrapper').css({
          fontSize: '100%',
          background: 'transparent',
          border: 'none',
          margin: 0,
          padding: 0
        }); var n = { width: e.width(), height: e.height() }; var o = document.activeElement
        try {
          o.id
        } catch (t) {
          o = document.body
        }
        return e.wrap(s), (e[0] === o || t.contains(e[0], o)) && t(o).focus(), s = e.parent(), e.css('position') === 'static' ? (s.css({ position: 'relative' }), e.css({ position: 'relative' })) : (t.extend(i, {
          position: e.css('position'),
          zIndex: e.css('z-index')
        }), t.each(['top', 'left', 'bottom', 'right'], function (t, s) {
          i[s] = e.css(s), isNaN(parseInt(i[s], 10)) && (i[s] = 'auto')
        }), e.css({
          position: 'relative',
          top: 0,
          left: 0,
          right: 'auto',
          bottom: 'auto'
        })), e.css(n), s.css(i).show()
      },
      removeWrapper: function (e) {
        var i = document.activeElement
        return e.parent().is('.ui-effects-wrapper') && (e.parent().replaceWith(e), (e[0] === i || t.contains(e[0], i)) && t(i).focus()), e
      },
      setTransition: function (e, i, s, n) {
        return n = n || {}, t.each(i, function (t, i) {
          var o = e.cssUnit(i)
          o[0] > 0 && (n[i] = o[0] * s + o[1])
        }), n
      }
    }), t.fn.extend({
      effect: function () {
        function e (e) {
          function i () {
            t.isFunction(o) && o.call(n[0]), t.isFunction(e) && e()
          }

          var n = t(this); var o = s.complete; var a = s.mode;
          (n.is(':hidden') ? a === 'hide' : a === 'show') ? i() : r.call(n[0], s, i)
        }

        var s = n.apply(this, arguments); var o = s.mode; var a = s.queue; var r = t.effects.effect[s.effect]
        var l = !r && i && t.effects[s.effect]
        return t.fx.off || !r && !l ? o ? this[o](s.duration, s.complete) : this.each(function () {
          s.complete && s.complete.call(this)
        }) : r ? !1 === a ? this.each(e) : this.queue(a || 'fx', e) : l.call(this, {
          options: s,
          duration: s.duration,
          callback: s.complete,
          mode: s.mode
        })
      },
      _show: t.fn.show,
      show: function (t) {
        if (o(t)) return this._show.apply(this, arguments)
        var e = n.apply(this, arguments)
        return e.mode = 'show', this.effect.call(this, e)
      },
      _hide: t.fn.hide,
      hide: function (t) {
        if (o(t)) return this._hide.apply(this, arguments)
        var e = n.apply(this, arguments)
        return e.mode = 'hide', this.effect.call(this, e)
      },
      __toggle: t.fn.toggle,
      toggle: function (e) {
        if (o(e) || typeof e === 'boolean' || t.isFunction(e)) return this.__toggle.apply(this, arguments)
        var i = n.apply(this, arguments)
        return i.mode = 'toggle', this.effect.call(this, i)
      },
      cssUnit: function (e) {
        var i = this.css(e); var s = []
        return t.each(['em', 'px', '%', 'pt'], function (t, e) {
          i.indexOf(e) > 0 && (s = [parseFloat(i), e])
        }), s
      }
    })
  }()), (function () {
    var e = {}
    t.each(['Quad', 'Cubic', 'Quart', 'Quint', 'Expo'], function (t, i) {
      e[i] = function (e) {
        return Math.pow(e, t + 2)
      }
    }), t.extend(e, {
      Sine: function (t) {
        return 1 - Math.cos(t * Math.PI / 2)
      },
      Circ: function (t) {
        return 1 - Math.sqrt(1 - t * t)
      },
      Elastic: function (t) {
        return t === 0 || t === 1 ? t : -Math.pow(2, 8 * (t - 1)) * Math.sin((80 * (t - 1) - 7.5) * Math.PI / 15)
      },
      Back: function (t) {
        return t * t * (3 * t - 2)
      },
      Bounce: function (t) {
        for (var e, i = 4; t < ((e = Math.pow(2, --i)) - 1) / 11;) ;
        return 1 / Math.pow(4, 3 - i) - 7.5625 * Math.pow((3 * e - 2) / 22 - t, 2)
      }
    }), t.each(e, function (e, i) {
      t.easing['easeIn' + e] = i, t.easing['easeOut' + e] = function (t) {
        return 1 - i(1 - t)
      }, t.easing['easeInOut' + e] = function (t) {
        return t < 0.5 ? i(2 * t) / 2 : 1 - i(-2 * t + 2) / 2
      }
    })
  }())
}(jQuery)), (function (t, e) {
  var i = /up|down|vertical/; var s = /up|left|vertical|horizontal/
  t.effects.effect.blind = function (e, n) {
    var o; var a; var r; var l = t(this); var h = ['position', 'top', 'bottom', 'left', 'right', 'height', 'width']
    var c = t.effects.setMode(l, e.mode || 'hide'); var u = e.direction || 'up'; var d = i.test(u)
    var p = d ? 'height' : 'width'; var f = d ? 'top' : 'left'; var g = s.test(u); var m = {}; var v = c === 'show'
    l.parent().is('.ui-effects-wrapper') ? t.effects.save(l.parent(), h) : t.effects.save(l, h), l.show(), a = (o = t.effects.createWrapper(l).css({ overflow: 'hidden' }))[p](), r = parseFloat(o.css(f)) || 0, m[p] = v ? a : 0, g || (l.css(d ? 'bottom' : 'right', 0).css(d ? 'top' : 'left', 'auto').css({ position: 'absolute' }), m[f] = v ? r : a + r), v && (o.css(p, 0), g || o.css(f, r + a)), o.animate(m, {
      duration: e.duration,
      easing: e.easing,
      queue: !1,
      complete: function () {
        c === 'hide' && l.hide(), t.effects.restore(l, h), t.effects.removeWrapper(l), n()
      }
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.bounce = function (e, i) {
    var s; var n; var o; var a = t(this); var r = ['position', 'top', 'bottom', 'left', 'right', 'height', 'width']
    var l = t.effects.setMode(a, e.mode || 'effect'); var h = l === 'hide'; var c = l === 'show'; var u = e.direction || 'up'
    var d = e.distance; var p = e.times || 5; var f = 2 * p + (c || h ? 1 : 0); var g = e.duration / f; var m = e.easing
    var v = u === 'up' || u === 'down' ? 'top' : 'left'; var b = u === 'up' || u === 'left'; var _ = a.queue()
    var y = _.length
    for ((c || h) && r.push('opacity'), t.effects.save(a, r), a.show(), t.effects.createWrapper(a), d || (d = a[v === 'top' ? 'outerHeight' : 'outerWidth']() / 3), c && (o = { opacity: 1 }, o[v] = 0, a.css('opacity', 0).css(v, b ? 2 * -d : 2 * d).animate(o, g, m)), h && (d /= Math.pow(2, p - 1)), (o = {})[v] = 0, s = 0; s < p; s++) n = {}, n[v] = (b ? '-=' : '+=') + d, a.animate(n, g, m).animate(o, g, m), d = h ? 2 * d : d / 2
    h && (n = { opacity: 0 }, n[v] = (b ? '-=' : '+=') + d, a.animate(n, g, m)), a.queue(function () {
      h && a.hide(), t.effects.restore(a, r), t.effects.removeWrapper(a), i()
    }), y > 1 && _.splice.apply(_, [1, 0].concat(_.splice(y, f + 1))), a.dequeue()
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.clip = function (e, i) {
    var s; var n; var o; var a = t(this); var r = ['position', 'top', 'bottom', 'left', 'right', 'height', 'width']
    var l = t.effects.setMode(a, e.mode || 'hide') === 'show'; var h = (e.direction || 'vertical') === 'vertical'
    var c = h ? 'height' : 'width'; var u = h ? 'top' : 'left'; var d = {}
    t.effects.save(a, r), a.show(), s = t.effects.createWrapper(a).css({ overflow: 'hidden' }), o = (n = a[0].tagName === 'IMG' ? s : a)[c](), l && (n.css(c, 0), n.css(u, o / 2)), d[c] = l ? o : 0, d[u] = l ? 0 : o / 2, n.animate(d, {
      queue: !1,
      duration: e.duration,
      easing: e.easing,
      complete: function () {
        l || a.hide(), t.effects.restore(a, r), t.effects.removeWrapper(a), i()
      }
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.drop = function (e, i) {
    var s; var n = t(this); var o = ['position', 'top', 'bottom', 'left', 'right', 'opacity', 'height', 'width']
    var a = t.effects.setMode(n, e.mode || 'hide'); var r = a === 'show'; var l = e.direction || 'left'
    var h = l === 'up' || l === 'down' ? 'top' : 'left'; var c = l === 'up' || l === 'left' ? 'pos' : 'neg'
    var u = { opacity: r ? 1 : 0 }
    t.effects.save(n, o), n.show(), t.effects.createWrapper(n), s = e.distance || n[h === 'top' ? 'outerHeight' : 'outerWidth'](!0) / 2, r && n.css('opacity', 0).css(h, c === 'pos' ? -s : s), u[h] = (r ? c === 'pos' ? '+=' : '-=' : c === 'pos' ? '-=' : '+=') + s, n.animate(u, {
      queue: !1,
      duration: e.duration,
      easing: e.easing,
      complete: function () {
        a === 'hide' && n.hide(), t.effects.restore(n, o), t.effects.removeWrapper(n), i()
      }
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.explode = function (e, i) {
    function s () {
      d.css({ visibility: 'visible' }), t(v).remove(), p || d.hide(), i()
    }

    var n; var o; var a; var r; var l; var h; var c = e.pieces ? Math.round(Math.sqrt(e.pieces)) : 3; var u = c; var d = t(this)
    var p = t.effects.setMode(d, e.mode || 'hide') === 'show'; var f = d.show().css('visibility', 'hidden').offset()
    var g = Math.ceil(d.outerWidth() / u); var m = Math.ceil(d.outerHeight() / c); var v = []
    for (n = 0; n < c; n++) {
      for (r = f.top + n * m, h = n - (c - 1) / 2, o = 0; o < u; o++) {
        a = f.left + o * g, l = o - (u - 1) / 2, d.clone().appendTo('body').wrap('<div></div>').css({
          position: 'absolute',
          visibility: 'visible',
          left: -o * g,
          top: -n * m
        }).parent().addClass('ui-effects-explode').css({
          position: 'absolute',
          overflow: 'hidden',
          width: g,
          height: m,
          left: a + (p ? l * g : 0),
          top: r + (p ? h * m : 0),
          opacity: p ? 0 : 1
        }).animate({
          left: a + (p ? 0 : l * g),
          top: r + (p ? 0 : h * m),
          opacity: p ? 1 : 0
        }, e.duration || 500, e.easing, function () {
          v.push(this), v.length === c * u && s()
        })
      }
    }
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.fade = function (e, i) {
    var s = t(this); var n = t.effects.setMode(s, e.mode || 'toggle')
    s.animate({ opacity: n }, { queue: !1, duration: e.duration, easing: e.easing, complete: i })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.fold = function (e, i) {
    var s; var n; var o = t(this); var a = ['position', 'top', 'bottom', 'left', 'right', 'height', 'width']
    var r = t.effects.setMode(o, e.mode || 'hide'); var l = r === 'show'; var h = r === 'hide'; var c = e.size || 15
    var u = /([0-9]+)%/.exec(c); var d = !!e.horizFirst; var p = l !== d; var f = p ? ['width', 'height'] : ['height', 'width']
    var g = e.duration / 2; var m = {}; var v = {}
    t.effects.save(o, a), o.show(), s = t.effects.createWrapper(o).css({ overflow: 'hidden' }), n = p ? [s.width(), s.height()] : [s.height(), s.width()], u && (c = parseInt(u[1], 10) / 100 * n[h ? 0 : 1]), l && s.css(d ? {
      height: 0,
      width: c
    } : {
      height: c,
      width: 0
    }), m[f[0]] = l ? n[0] : c, v[f[1]] = l ? n[1] : 0, s.animate(m, g, e.easing).animate(v, g, e.easing, function () {
      h && o.hide(), t.effects.restore(o, a), t.effects.removeWrapper(o), i()
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.highlight = function (e, i) {
    var s = t(this); var n = ['backgroundImage', 'backgroundColor', 'opacity']
    var o = t.effects.setMode(s, e.mode || 'show'); var a = { backgroundColor: s.css('backgroundColor') }
    o === 'hide' && (a.opacity = 0), t.effects.save(s, n), s.show().css({
      backgroundImage: 'none',
      backgroundColor: e.color || '#ffff99'
    }).animate(a, {
      queue: !1,
      duration: e.duration,
      easing: e.easing,
      complete: function () {
        o === 'hide' && s.hide(), t.effects.restore(s, n), i()
      }
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.pulsate = function (e, i) {
    var s; var n = t(this); var o = t.effects.setMode(n, e.mode || 'show'); var a = o === 'show'; var r = o === 'hide'
    var l = a || o === 'hide'; var h = 2 * (e.times || 5) + (l ? 1 : 0); var c = e.duration / h; var u = 0; var d = n.queue()
    var p = d.length
    for (!a && n.is(':visible') || (n.css('opacity', 0).show(), u = 1), s = 1; s < h; s++) n.animate({ opacity: u }, c, e.easing), u = 1 - u
    n.animate({ opacity: u }, c, e.easing), n.queue(function () {
      r && n.hide(), i()
    }), p > 1 && d.splice.apply(d, [1, 0].concat(d.splice(p, h + 1))), n.dequeue()
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.puff = function (e, i) {
    var s = t(this); var n = t.effects.setMode(s, e.mode || 'hide'); var o = n === 'hide'
    var a = parseInt(e.percent, 10) || 150; var r = a / 100
    var l = { height: s.height(), width: s.width(), outerHeight: s.outerHeight(), outerWidth: s.outerWidth() }
    t.extend(e, {
      effect: 'scale',
      queue: !1,
      fade: !0,
      mode: n,
      complete: i,
      percent: o ? a : 100,
      from: o ? l : {
        height: l.height * r,
        width: l.width * r,
        outerHeight: l.outerHeight * r,
        outerWidth: l.outerWidth * r
      }
    }), s.effect(e)
  }, t.effects.effect.scale = function (e, i) {
    var s = t(this); var n = t.extend(!0, {}, e); var o = t.effects.setMode(s, e.mode || 'effect')
    var a = parseInt(e.percent, 10) || (parseInt(e.percent, 10) === 0 ? 0 : o === 'hide' ? 0 : 100)
    var r = e.direction || 'both'; var l = e.origin
    var h = { height: s.height(), width: s.width(), outerHeight: s.outerHeight(), outerWidth: s.outerWidth() }
    var c = { y: r !== 'horizontal' ? a / 100 : 1, x: r !== 'vertical' ? a / 100 : 1 }
    n.effect = 'size', n.queue = !1, n.complete = i, o !== 'effect' && (n.origin = l || ['middle', 'center'], n.restore = !0), n.from = e.from || (o === 'show' ? {
      height: 0,
      width: 0,
      outerHeight: 0,
      outerWidth: 0
    } : h), n.to = {
      height: h.height * c.y,
      width: h.width * c.x,
      outerHeight: h.outerHeight * c.y,
      outerWidth: h.outerWidth * c.x
    }, n.fade && (o === 'show' && (n.from.opacity = 0, n.to.opacity = 1), o === 'hide' && (n.from.opacity = 1, n.to.opacity = 0)), s.effect(n)
  }, t.effects.effect.size = function (e, i) {
    var s; var n; var o; var a = t(this)
    var r = ['position', 'top', 'bottom', 'left', 'right', 'width', 'height', 'overflow', 'opacity']
    var l = ['position', 'top', 'bottom', 'left', 'right', 'overflow', 'opacity']
    var h = ['width', 'height', 'overflow']; var c = ['fontSize']
    var u = ['borderTopWidth', 'borderBottomWidth', 'paddingTop', 'paddingBottom']
    var d = ['borderLeftWidth', 'borderRightWidth', 'paddingLeft', 'paddingRight']
    var p = t.effects.setMode(a, e.mode || 'effect'); var f = e.restore || p !== 'effect'; var g = e.scale || 'both'
    var m = e.origin || ['middle', 'center']; var v = a.css('position'); var b = f ? r : l
    var _ = { height: 0, width: 0, outerHeight: 0, outerWidth: 0 }
    p === 'show' && a.show(), s = {
      height: a.height(),
      width: a.width(),
      outerHeight: a.outerHeight(),
      outerWidth: a.outerWidth()
    }, e.mode === 'toggle' && p === 'show' ? (a.from = e.to || _, a.to = e.from || s) : (a.from = e.from || (p === 'show' ? _ : s), a.to = e.to || (p === 'hide' ? _ : s)), o = {
      from: {
        y: a.from.height / s.height,
        x: a.from.width / s.width
      },
      to: { y: a.to.height / s.height, x: a.to.width / s.width }
    }, g !== 'box' && g !== 'both' || (o.from.y !== o.to.y && (b = b.concat(u), a.from = t.effects.setTransition(a, u, o.from.y, a.from), a.to = t.effects.setTransition(a, u, o.to.y, a.to)), o.from.x !== o.to.x && (b = b.concat(d), a.from = t.effects.setTransition(a, d, o.from.x, a.from), a.to = t.effects.setTransition(a, d, o.to.x, a.to))), g !== 'content' && g !== 'both' || o.from.y !== o.to.y && (b = b.concat(c).concat(h), a.from = t.effects.setTransition(a, c, o.from.y, a.from), a.to = t.effects.setTransition(a, c, o.to.y, a.to)), t.effects.save(a, b), a.show(), t.effects.createWrapper(a), a.css('overflow', 'hidden').css(a.from), m && (n = t.effects.getBaseline(m, s), a.from.top = (s.outerHeight - a.outerHeight()) * n.y, a.from.left = (s.outerWidth - a.outerWidth()) * n.x, a.to.top = (s.outerHeight - a.to.outerHeight) * n.y, a.to.left = (s.outerWidth - a.to.outerWidth) * n.x), a.css(a.from), g !== 'content' && g !== 'both' || (u = u.concat(['marginTop', 'marginBottom']).concat(c), d = d.concat(['marginLeft', 'marginRight']), h = r.concat(u).concat(d), a.find('*[width]').each(function () {
      var i = t(this)
      var s = { height: i.height(), width: i.width(), outerHeight: i.outerHeight(), outerWidth: i.outerWidth() }
      f && t.effects.save(i, h), i.from = {
        height: s.height * o.from.y,
        width: s.width * o.from.x,
        outerHeight: s.outerHeight * o.from.y,
        outerWidth: s.outerWidth * o.from.x
      }, i.to = {
        height: s.height * o.to.y,
        width: s.width * o.to.x,
        outerHeight: s.height * o.to.y,
        outerWidth: s.width * o.to.x
      }, o.from.y !== o.to.y && (i.from = t.effects.setTransition(i, u, o.from.y, i.from), i.to = t.effects.setTransition(i, u, o.to.y, i.to)), o.from.x !== o.to.x && (i.from = t.effects.setTransition(i, d, o.from.x, i.from), i.to = t.effects.setTransition(i, d, o.to.x, i.to)), i.css(i.from), i.animate(i.to, e.duration, e.easing, function () {
        f && t.effects.restore(i, h)
      })
    })), a.animate(a.to, {
      queue: !1,
      duration: e.duration,
      easing: e.easing,
      complete: function () {
        a.to.opacity === 0 && a.css('opacity', a.from.opacity), p === 'hide' && a.hide(), t.effects.restore(a, b), f || (v === 'static' ? a.css({
          position: 'relative',
          top: a.to.top,
          left: a.to.left
        }) : t.each(['top', 'left'], function (t, e) {
          a.css(e, function (e, i) {
            var s = parseInt(i, 10); var n = t ? a.to.left : a.to.top
            return i === 'auto' ? n + 'px' : s + n + 'px'
          })
        })), t.effects.removeWrapper(a), i()
      }
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.shake = function (e, i) {
    var s; var n = t(this); var o = ['position', 'top', 'bottom', 'left', 'right', 'height', 'width']
    var a = t.effects.setMode(n, e.mode || 'effect'); var r = e.direction || 'left'; var l = e.distance || 20
    var h = e.times || 3; var c = 2 * h + 1; var u = Math.round(e.duration / c)
    var d = r === 'up' || r === 'down' ? 'top' : 'left'; var p = r === 'up' || r === 'left'; var f = {}; var g = {}; var m = {}
    var v = n.queue(); var b = v.length
    for (t.effects.save(n, o), n.show(), t.effects.createWrapper(n), f[d] = (p ? '-=' : '+=') + l, g[d] = (p ? '+=' : '-=') + 2 * l, m[d] = (p ? '-=' : '+=') + 2 * l, n.animate(f, u, e.easing), s = 1; s < h; s++) n.animate(g, u, e.easing).animate(m, u, e.easing)
    n.animate(g, u, e.easing).animate(f, u / 2, e.easing).queue(function () {
      a === 'hide' && n.hide(), t.effects.restore(n, o), t.effects.removeWrapper(n), i()
    }), b > 1 && v.splice.apply(v, [1, 0].concat(v.splice(b, c + 1))), n.dequeue()
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.slide = function (e, i) {
    var s; var n = t(this); var o = ['position', 'top', 'bottom', 'left', 'right', 'width', 'height']
    var a = t.effects.setMode(n, e.mode || 'show'); var r = a === 'show'; var l = e.direction || 'left'
    var h = l === 'up' || l === 'down' ? 'top' : 'left'; var c = l === 'up' || l === 'left'; var u = {}
    t.effects.save(n, o), n.show(), s = e.distance || n[h === 'top' ? 'outerHeight' : 'outerWidth'](!0), t.effects.createWrapper(n).css({ overflow: 'hidden' }), r && n.css(h, c ? isNaN(s) ? '-' + s : -s : s), u[h] = (r ? c ? '+=' : '-=' : c ? '-=' : '+=') + s, n.animate(u, {
      queue: !1,
      duration: e.duration,
      easing: e.easing,
      complete: function () {
        a === 'hide' && n.hide(), t.effects.restore(n, o), t.effects.removeWrapper(n), i()
      }
    })
  }
}(jQuery)), (function (t, e) {
  t.effects.effect.transfer = function (e, i) {
    var s = t(this); var n = t(e.to); var o = n.css('position') === 'fixed'; var a = t('body'); var r = o ? a.scrollTop() : 0
    var l = o ? a.scrollLeft() : 0; var h = n.offset()
    var c = { top: h.top - r, left: h.left - l, height: n.innerHeight(), width: n.innerWidth() }; var u = s.offset()
    var d = t('<div class="ui-effects-transfer"></div>').appendTo(document.body).addClass(e.className).css({
      top: u.top - r,
      left: u.left - l,
      height: s.innerHeight(),
      width: s.innerWidth(),
      position: o ? 'fixed' : 'absolute'
    }).animate(c, e.duration, e.easing, function () {
      d.remove(), i()
    })
  }
}(jQuery)), (function (t, e) {
  var i = !1
  t.widget('ui.menu', {
    version: '1.9.2',
    defaultElement: '<ul>',
    delay: 300,
    options: {
      icons: { submenu: 'ui-icon-carat-1-e' },
      menus: 'ul',
      position: { my: 'left top', at: 'right top' },
      role: 'menu',
      blur: null,
      focus: null,
      select: null
    },
    _create: function () {
      this.activeMenu = this.element, this.element.uniqueId().addClass('ui-menu ui-widget ui-widget-content ui-corner-all').toggleClass('ui-menu-icons', !!this.element.find('.ui-icon').length).attr({
        role: this.options.role,
        tabIndex: 0
      }).bind('click' + this.eventNamespace, t.proxy(function (t) {
        this.options.disabled && t.preventDefault()
      }, this)), this.options.disabled && this.element.addClass('ui-state-disabled').attr('aria-disabled', 'true'), this._on({
        'mousedown .ui-menu-item > a': function (t) {
          t.preventDefault()
        },
        'click .ui-state-disabled > a': function (t) {
          t.preventDefault()
        },
        'click .ui-menu-item:has(a)': function (e) {
          var s = t(e.target).closest('.ui-menu-item')
          !i && s.not('.ui-state-disabled').length && (i = !0, this.select(e), s.has('.ui-menu').length ? this.expand(e) : this.element.is(':focus') || (this.element.trigger('focus', [!0]), this.active && this.active.parents('.ui-menu').length === 1 && clearTimeout(this.timer)))
        },
        'mouseenter .ui-menu-item': function (e) {
          var i = t(e.currentTarget)
          i.siblings().children('.ui-state-active').removeClass('ui-state-active'), this.focus(e, i)
        },
        mouseleave: 'collapseAll',
        'mouseleave .ui-menu': 'collapseAll',
        focus: function (t, e) {
          var i = this.active || this.element.children('.ui-menu-item').eq(0)
          e || this.focus(t, i)
        },
        blur: function (e) {
          this._delay(function () {
            t.contains(this.element[0], this.document[0].activeElement) || this.collapseAll(e)
          })
        },
        keydown: '_keydown'
      }), this.refresh(), this._on(this.document, {
        click: function (e) {
          t(e.target).closest('.ui-menu').length || this.collapseAll(e), i = !1
        }
      })
    },
    _destroy: function () {
      this.element.removeAttr('aria-activedescendant').find('.ui-menu').addBack().removeClass('ui-menu ui-widget ui-widget-content ui-corner-all ui-menu-icons').removeAttr('role').removeAttr('tabIndex').removeAttr('aria-labelledby').removeAttr('aria-expanded').removeAttr('aria-hidden').removeAttr('aria-disabled').removeUniqueId().show(), this.element.find('.ui-menu-item').removeClass('ui-menu-item').removeAttr('role').removeAttr('aria-disabled').children('a').removeUniqueId().removeClass('ui-corner-all ui-state-hover').removeAttr('tabIndex').removeAttr('role').removeAttr('aria-haspopup').children().each(function () {
        var e = t(this)
        e.data('ui-menu-submenu-carat') && e.remove()
      }), this.element.find('.ui-menu-divider').removeClass('ui-menu-divider ui-widget-content')
    },
    _keydown: function (e) {
      function i (t) {
        return t.replace(/[\-\[\]{}()*+?.,\\\^$|#\s]/g, '\\$&')
      }

      var s; var n; var o; var a; var r; var l = !0
      switch (e.keyCode) {
        case t.ui.keyCode.PAGE_UP:
          this.previousPage(e)
          break
        case t.ui.keyCode.PAGE_DOWN:
          this.nextPage(e)
          break
        case t.ui.keyCode.HOME:
          this._move('first', 'first', e)
          break
        case t.ui.keyCode.END:
          this._move('last', 'last', e)
          break
        case t.ui.keyCode.UP:
          this.previous(e)
          break
        case t.ui.keyCode.DOWN:
          this.next(e)
          break
        case t.ui.keyCode.LEFT:
          this.collapse(e)
          break
        case t.ui.keyCode.RIGHT:
          this.active && !this.active.is('.ui-state-disabled') && this.expand(e)
          break
        case t.ui.keyCode.ENTER:
        case t.ui.keyCode.SPACE:
          this._activate(e)
          break
        case t.ui.keyCode.ESCAPE:
          this.collapse(e)
          break
        default:
          l = !1, n = this.previousFilter || '', o = String.fromCharCode(e.keyCode), a = !1, clearTimeout(this.filterTimer), o === n ? a = !0 : o = n + o, r = new RegExp('^' + i(o), 'i'), s = this.activeMenu.children('.ui-menu-item').filter(function () {
            return r.test(t(this).children('a').text())
          }), (s = a && s.index(this.active.next()) !== -1 ? this.active.nextAll('.ui-menu-item') : s).length || (o = String.fromCharCode(e.keyCode), r = new RegExp('^' + i(o), 'i'), s = this.activeMenu.children('.ui-menu-item').filter(function () {
            return r.test(t(this).children('a').text())
          })), s.length ? (this.focus(e, s), s.length > 1 ? (this.previousFilter = o, this.filterTimer = this._delay(function () {
            delete this.previousFilter
          }, 1e3)) : delete this.previousFilter) : delete this.previousFilter
      }
      l && e.preventDefault()
    },
    _activate: function (t) {
      this.active.is('.ui-state-disabled') || (this.active.children("a[aria-haspopup='true']").length ? this.expand(t) : this.select(t))
    },
    refresh: function () {
      var e; var i = this.options.icons.submenu; var s = this.element.find(this.options.menus)
      s.filter(':not(.ui-menu)').addClass('ui-menu ui-widget ui-widget-content ui-corner-all').hide().attr({
        role: this.options.role,
        'aria-hidden': 'true',
        'aria-expanded': 'false'
      }).each(function () {
        var e = t(this); var s = e.prev('a')
        var n = t('<span>').addClass('ui-menu-icon ui-icon ' + i).data('ui-menu-submenu-carat', !0)
        s.attr('aria-haspopup', 'true').prepend(n), e.attr('aria-labelledby', s.attr('id'))
      }), (e = s.add(this.element)).children(':not(.ui-menu-item):has(a)').addClass('ui-menu-item').attr('role', 'presentation').children('a').uniqueId().addClass('ui-corner-all').attr({
        tabIndex: -1,
        role: this._itemRole()
      }), e.children(':not(.ui-menu-item)').each(function () {
        var e = t(this);
        /[^\-—–\s]/.test(e.text()) || e.addClass('ui-widget-content ui-menu-divider')
      }), e.children('.ui-state-disabled').attr('aria-disabled', 'true'), this.active && !t.contains(this.element[0], this.active[0]) && this.blur()
    },
    _itemRole: function () {
      return { menu: 'menuitem', listbox: 'option' }[this.options.role]
    },
    focus: function (t, e) {
      var i, s
      this.blur(t, t && t.type === 'focus'), this._scrollIntoView(e), this.active = e.first(), s = this.active.children('a').addClass('ui-state-focus'), this.options.role && this.element.attr('aria-activedescendant', s.attr('id')), this.active.parent().closest('.ui-menu-item').children('a:first').addClass('ui-state-active'), t && t.type === 'keydown' ? this._close() : this.timer = this._delay(function () {
        this._close()
      }, this.delay), (i = e.children('.ui-menu')).length && /^mouse/.test(t.type) && this._startOpening(i), this.activeMenu = e.parent(), this._trigger('focus', t, { item: e })
    },
    _scrollIntoView: function (e) {
      var i, s, n, o, a, r
      this._hasScroll() && (i = parseFloat(t.css(this.activeMenu[0], 'borderTopWidth')) || 0, s = parseFloat(t.css(this.activeMenu[0], 'paddingTop')) || 0, n = e.offset().top - this.activeMenu.offset().top - i - s, o = this.activeMenu.scrollTop(), a = this.activeMenu.height(), r = e.height(), n < 0 ? this.activeMenu.scrollTop(o + n) : n + r > a && this.activeMenu.scrollTop(o + n - a + r))
    },
    blur: function (t, e) {
      e || clearTimeout(this.timer), this.active && (this.active.children('a').removeClass('ui-state-focus'), this.active = null, this._trigger('blur', t, { item: this.active }))
    },
    _startOpening: function (t) {
      clearTimeout(this.timer), t.attr('aria-hidden') === 'true' && (this.timer = this._delay(function () {
        this._close(), this._open(t)
      }, this.delay))
    },
    _open: function (e) {
      var i = t.extend({ of: this.active }, this.options.position)
      clearTimeout(this.timer), this.element.find('.ui-menu').not(e.parents('.ui-menu')).hide().attr('aria-hidden', 'true'), e.show().removeAttr('aria-hidden').attr('aria-expanded', 'true').position(i)
    },
    collapseAll: function (e, i) {
      clearTimeout(this.timer), this.timer = this._delay(function () {
        var s = i ? this.element : t(e && e.target).closest(this.element.find('.ui-menu'))
        s.length || (s = this.element), this._close(s), this.blur(e), this.activeMenu = s
      }, this.delay)
    },
    _close: function (t) {
      t || (t = this.active ? this.active.parent() : this.element), t.find('.ui-menu').hide().attr('aria-hidden', 'true').attr('aria-expanded', 'false').end().find('a.ui-state-active').removeClass('ui-state-active')
    },
    collapse: function (t) {
      var e = this.active && this.active.parent().closest('.ui-menu-item', this.element)
      e && e.length && (this._close(), this.focus(t, e))
    },
    expand: function (t) {
      var e = this.active && this.active.children('.ui-menu ').children('.ui-menu-item').first()
      e && e.length && (this._open(e.parent()), this._delay(function () {
        this.focus(t, e)
      }))
    },
    next: function (t) {
      this._move('next', 'first', t)
    },
    previous: function (t) {
      this._move('prev', 'last', t)
    },
    isFirstItem: function () {
      return this.active && !this.active.prevAll('.ui-menu-item').length
    },
    isLastItem: function () {
      return this.active && !this.active.nextAll('.ui-menu-item').length
    },
    _move: function (t, e, i) {
      var s
      this.active && (s = t === 'first' || t === 'last' ? this.active[t === 'first' ? 'prevAll' : 'nextAll']('.ui-menu-item').eq(-1) : this.active[t + 'All']('.ui-menu-item').eq(0)), s && s.length && this.active || (s = this.activeMenu.children('.ui-menu-item')[e]()), this.focus(i, s)
    },
    nextPage: function (e) {
      var i, s, n
      return this.active ? void (this.isLastItem() || (this._hasScroll() ? (s = this.active.offset().top, n = this.element.height(), this.active.nextAll('.ui-menu-item').each(function () {
        return (i = t(this)).offset().top - s - n < 0
      }), this.focus(e, i)) : this.focus(e, this.activeMenu.children('.ui-menu-item')[this.active ? 'last' : 'first']()))) : void this.next(e)
    },
    previousPage: function (e) {
      var i, s, n
      return this.active ? void (this.isFirstItem() || (this._hasScroll() ? (s = this.active.offset().top, n = this.element.height(), this.active.prevAll('.ui-menu-item').each(function () {
        return (i = t(this)).offset().top - s + n > 0
      }), this.focus(e, i)) : this.focus(e, this.activeMenu.children('.ui-menu-item').first()))) : void this.next(e)
    },
    _hasScroll: function () {
      return this.element.outerHeight() < this.element.prop('scrollHeight')
    },
    select: function (e) {
      this.active = this.active || t(e.target).closest('.ui-menu-item')
      var i = { item: this.active }
      this.active.has('.ui-menu').length || this.collapseAll(e, !0), this._trigger('select', e, i)
    }
  })
}(jQuery)), (function (t, e) {
  t.widget('ui.progressbar', {
    version: '1.9.2',
    options: { value: 0, max: 100 },
    min: 0,
    _create: function () {
      this.element.addClass('ui-progressbar ui-widget ui-widget-content ui-corner-all').attr({
        role: 'progressbar',
        'aria-valuemin': this.min,
        'aria-valuemax': this.options.max,
        'aria-valuenow': this._value()
      }), this.valueDiv = t("<div class='ui-progressbar-value ui-widget-header ui-corner-left'></div>").appendTo(this.element), this.oldValue = this._value(), this._refreshValue()
    },
    _destroy: function () {
      this.element.removeClass('ui-progressbar ui-widget ui-widget-content ui-corner-all').removeAttr('role').removeAttr('aria-valuemin').removeAttr('aria-valuemax').removeAttr('aria-valuenow'), this.valueDiv.remove()
    },
    value: function (t) {
      return void 0 === t ? this._value() : (this._setOption('value', t), this)
    },
    _setOption: function (t, e) {
      t === 'value' && (this.options.value = e, this._refreshValue(), this._value() === this.options.max && this._trigger('complete')), this._super(t, e)
    },
    _value: function () {
      var t = this.options.value
      return typeof t !== 'number' && (t = 0), Math.min(this.options.max, Math.max(this.min, t))
    },
    _percentage: function () {
      return 100 * this._value() / this.options.max
    },
    _refreshValue: function () {
      var t = this.value(); var e = this._percentage()
      this.oldValue !== t && (this.oldValue = t, this._trigger('change')), this.valueDiv.toggle(t > this.min).toggleClass('ui-corner-right', t === this.options.max).width(e.toFixed(0) + '%'), this.element.attr('aria-valuenow', t)
    }
  })
}(jQuery)), (function (t, e) {
  t.widget('ui.resizable', t.ui.mouse, {
    version: '1.9.2',
    widgetEventPrefix: 'resize',
    options: {
      alsoResize: !1,
      animate: !1,
      animateDuration: 'slow',
      animateEasing: 'swing',
      aspectRatio: !1,
      autoHide: !1,
      containment: !1,
      ghost: !1,
      grid: !1,
      handles: 'e,s,se',
      helper: !1,
      maxHeight: null,
      maxWidth: null,
      minHeight: 10,
      minWidth: 10,
      zIndex: 1e3
    },
    _create: function () {
      var e = this; var i = this.options
      if (this.element.addClass('ui-resizable'), t.extend(this, {
        _aspectRatio: !!i.aspectRatio,
        aspectRatio: i.aspectRatio,
        originalElement: this.element,
        _proportionallyResizeElements: [],
        _helper: i.helper || i.ghost || i.animate ? i.helper || 'ui-resizable-helper' : null
      }), this.element[0].nodeName.match(/canvas|textarea|input|select|button|img/i) && (this.element.wrap(t('<div class="ui-wrapper" style="overflow: hidden;"></div>').css({
        position: this.element.css('position'),
        width: this.element.outerWidth(),
        height: this.element.outerHeight(),
        top: this.element.css('top'),
        left: this.element.css('left')
      })), this.element = this.element.parent().data('resizable', this.element.data('resizable')), this.elementIsWrapper = !0, this.element.css({
        marginLeft: this.originalElement.css('marginLeft'),
        marginTop: this.originalElement.css('marginTop'),
        marginRight: this.originalElement.css('marginRight'),
        marginBottom: this.originalElement.css('marginBottom')
      }), this.originalElement.css({
        marginLeft: 0,
        marginTop: 0,
        marginRight: 0,
        marginBottom: 0
      }), this.originalResizeStyle = this.originalElement.css('resize'), this.originalElement.css('resize', 'none'), this._proportionallyResizeElements.push(this.originalElement.css({
        position: 'static',
        zoom: 1,
        display: 'block'
      })), this.originalElement.css({ margin: this.originalElement.css('margin') }), this._proportionallyResize()), this.handles = i.handles || (t('.ui-resizable-handle', this.element).length ? {
        n: '.ui-resizable-n',
        e: '.ui-resizable-e',
        s: '.ui-resizable-s',
        w: '.ui-resizable-w',
        se: '.ui-resizable-se',
        sw: '.ui-resizable-sw',
        ne: '.ui-resizable-ne',
        nw: '.ui-resizable-nw'
      } : 'e,s,se'), this.handles.constructor == String) {
        this.handles == 'all' && (this.handles = 'n,e,s,w,se,sw,ne,nw')
        var s = this.handles.split(',')
        this.handles = {}
        for (var n = 0; n < s.length; n++) {
          var o = t.trim(s[n])
          var a = t('<div class="ui-resizable-handle ' + ('ui-resizable-' + o) + '"></div>')
          a.css({ zIndex: i.zIndex }), o == 'se' && a.addClass('ui-icon ui-icon-gripsmall-diagonal-se'), this.handles[o] = '.ui-resizable-' + o, this.element.append(a)
        }
      }
      this._renderAxis = function (e) {
        e = e || this.element
        for (var i in this.handles) {
          if (this.handles[i].constructor == String && (this.handles[i] = t(this.handles[i], this.element).show()), this.elementIsWrapper && this.originalElement[0].nodeName.match(/textarea|input|select|button/i)) {
            var s = t(this.handles[i], this.element); var n = 0
            n = /sw|ne|nw|se|n|s/.test(i) ? s.outerHeight() : s.outerWidth()
            var o = ['padding', /ne|nw|n/.test(i) ? 'Top' : /se|sw|s/.test(i) ? 'Bottom' : /^e$/.test(i) ? 'Right' : 'Left'].join('')
            e.css(o, n), this._proportionallyResize()
          }
          t(this.handles[i]).length
        }
      }, this._renderAxis(this.element), this._handles = t('.ui-resizable-handle', this.element).disableSelection(), this._handles.mouseover(function () {
        if (!e.resizing) {
          if (this.className) var t = this.className.match(/ui-resizable-(se|sw|ne|nw|n|e|s|w)/i)
          e.axis = t && t[1] ? t[1] : 'se'
        }
      }), i.autoHide && (this._handles.hide(), t(this.element).addClass('ui-resizable-autohide').mouseenter(function () {
        i.disabled || (t(this).removeClass('ui-resizable-autohide'), e._handles.show())
      }).mouseleave(function () {
        i.disabled || e.resizing || (t(this).addClass('ui-resizable-autohide'), e._handles.hide())
      })), this._mouseInit()
    },
    _destroy: function () {
      this._mouseDestroy()
      var e = function (e) {
        t(e).removeClass('ui-resizable ui-resizable-disabled ui-resizable-resizing').removeData('resizable').removeData('ui-resizable').unbind('.resizable').find('.ui-resizable-handle').remove()
      }
      if (this.elementIsWrapper) {
        e(this.element)
        var i = this.element
        this.originalElement.css({
          position: i.css('position'),
          width: i.outerWidth(),
          height: i.outerHeight(),
          top: i.css('top'),
          left: i.css('left')
        }).insertAfter(i), i.remove()
      }
      return this.originalElement.css('resize', this.originalResizeStyle), e(this.originalElement), this
    },
    _mouseCapture: function (e) {
      var i = !1
      for (var s in this.handles) t(this.handles[s])[0] == e.target && (i = !0)
      return !this.options.disabled && i
    },
    _mouseStart: function (e) {
      var s = this.options; var n = this.element.position(); var o = this.element
      this.resizing = !0, this.documentScroll = {
        top: t(document).scrollTop(),
        left: t(document).scrollLeft()
      }, (o.is('.ui-draggable') || /absolute/.test(o.css('position'))) && o.css({
        position: 'absolute',
        top: n.top,
        left: n.left
      }), this._renderProxy()
      var a = i(this.helper.css('left')); var r = i(this.helper.css('top'))
      s.containment && (a += t(s.containment).scrollLeft() || 0, r += t(s.containment).scrollTop() || 0), this.offset = this.helper.offset(), this.position = {
        left: a,
        top: r
      }, this.size = this._helper ? { width: o.outerWidth(), height: o.outerHeight() } : {
        width: o.width(),
        height: o.height()
      }, this.originalSize = this._helper ? { width: o.outerWidth(), height: o.outerHeight() } : {
        width: o.width(),
        height: o.height()
      }, this.originalPosition = { left: a, top: r }, this.sizeDiff = {
        width: o.outerWidth() - o.width(),
        height: o.outerHeight() - o.height()
      }, this.originalMousePosition = {
        left: e.pageX,
        top: e.pageY
      }, this.aspectRatio = typeof s.aspectRatio === 'number' ? s.aspectRatio : this.originalSize.width / this.originalSize.height || 1
      var l = t('.ui-resizable-' + this.axis).css('cursor')
      return t('body').css('cursor', l == 'auto' ? this.axis + '-resize' : l), o.addClass('ui-resizable-resizing'), this._propagate('start', e), !0
    },
    _mouseDrag: function (t) {
      var e = this.helper; var i = (this.options, this.originalMousePosition); var s = this.axis
      var n = t.pageX - i.left || 0; var o = t.pageY - i.top || 0; var a = this._change[s]
      if (!a) return !1
      var r = a.apply(this, [t, n, o])
      return this._updateVirtualBoundaries(t.shiftKey), (this._aspectRatio || t.shiftKey) && (r = this._updateRatio(r, t)), r = this._respectSize(r, t), this._propagate('resize', t), e.css({
        top: this.position.top + 'px',
        left: this.position.left + 'px',
        width: this.size.width + 'px',
        height: this.size.height + 'px'
      }), !this._helper && this._proportionallyResizeElements.length && this._proportionallyResize(), this._updateCache(r), this._trigger('resize', t, this.ui()), !1
    },
    _mouseStop: function (e) {
      this.resizing = !1
      var i = this.options; var s = this
      if (this._helper) {
        var n = this._proportionallyResizeElements; var o = n.length && /textarea/i.test(n[0].nodeName)
        var a = o && t.ui.hasScroll(n[0], 'left') ? 0 : s.sizeDiff.height; var r = o ? 0 : s.sizeDiff.width
        var l = { width: s.helper.width() - r, height: s.helper.height() - a }
        var h = parseInt(s.element.css('left'), 10) + (s.position.left - s.originalPosition.left) || null
        var c = parseInt(s.element.css('top'), 10) + (s.position.top - s.originalPosition.top) || null
        i.animate || this.element.css(t.extend(l, {
          top: c,
          left: h
        })), s.helper.height(s.size.height), s.helper.width(s.size.width), this._helper && !i.animate && this._proportionallyResize()
      }
      return t('body').css('cursor', 'auto'), this.element.removeClass('ui-resizable-resizing'), this._propagate('stop', e), this._helper && this.helper.remove(), !1
    },
    _updateVirtualBoundaries: function (t) {
      var e; var i; var n; var o; var a; var r = this.options
      a = {
        minWidth: s(r.minWidth) ? r.minWidth : 0,
        maxWidth: s(r.maxWidth) ? r.maxWidth : 1 / 0,
        minHeight: s(r.minHeight) ? r.minHeight : 0,
        maxHeight: s(r.maxHeight) ? r.maxHeight : 1 / 0
      }, (this._aspectRatio || t) && (e = a.minHeight * this.aspectRatio, n = a.minWidth / this.aspectRatio, i = a.maxHeight * this.aspectRatio, o = a.maxWidth / this.aspectRatio, e > a.minWidth && (a.minWidth = e), n > a.minHeight && (a.minHeight = n), i < a.maxWidth && (a.maxWidth = i), o < a.maxHeight && (a.maxHeight = o)), this._vBoundaries = a
    },
    _updateCache: function (t) {
      this.options, this.offset = this.helper.offset(), s(t.left) && (this.position.left = t.left), s(t.top) && (this.position.top = t.top), s(t.height) && (this.size.height = t.height), s(t.width) && (this.size.width = t.width)
    },
    _updateRatio: function (t, e) {
      var i = (this.options, this.position); var n = this.size; var o = this.axis
      return s(t.height) ? t.width = t.height * this.aspectRatio : s(t.width) && (t.height = t.width / this.aspectRatio), o == 'sw' && (t.left = i.left + (n.width - t.width), t.top = null), o == 'nw' && (t.top = i.top + (n.height - t.height), t.left = i.left + (n.width - t.width)), t
    },
    _respectSize: function (t, e) {
      var i = (this.helper, this._vBoundaries); var n = (this._aspectRatio || e.shiftKey, this.axis)
      var o = s(t.width) && i.maxWidth && i.maxWidth < t.width
      var a = s(t.height) && i.maxHeight && i.maxHeight < t.height
      var r = s(t.width) && i.minWidth && i.minWidth > t.width
      var l = s(t.height) && i.minHeight && i.minHeight > t.height
      r && (t.width = i.minWidth), l && (t.height = i.minHeight), o && (t.width = i.maxWidth), a && (t.height = i.maxHeight)
      var h = this.originalPosition.left + this.originalSize.width; var c = this.position.top + this.size.height
      var u = /sw|nw|w/.test(n); var d = /nw|ne|n/.test(n)
      r && u && (t.left = h - i.minWidth), o && u && (t.left = h - i.maxWidth), l && d && (t.top = c - i.minHeight), a && d && (t.top = c - i.maxHeight)
      var p = !t.width && !t.height
      return p && !t.left && t.top ? t.top = null : p && !t.top && t.left && (t.left = null), t
    },
    _proportionallyResize: function () {
      if (this.options, this._proportionallyResizeElements.length) {
        for (var e = this.helper || this.element, i = 0; i < this._proportionallyResizeElements.length; i++) {
          var s = this._proportionallyResizeElements[i]
          if (!this.borderDif) {
            var n = [s.css('borderTopWidth'), s.css('borderRightWidth'), s.css('borderBottomWidth'), s.css('borderLeftWidth')]
            var o = [s.css('paddingTop'), s.css('paddingRight'), s.css('paddingBottom'), s.css('paddingLeft')]
            this.borderDif = t.map(n, function (t, e) {
              return (parseInt(t, 10) || 0) + (parseInt(o[e], 10) || 0)
            })
          }
          s.css({
            height: e.height() - this.borderDif[0] - this.borderDif[2] || 0,
            width: e.width() - this.borderDif[1] - this.borderDif[3] || 0
          })
        }
      }
    },
    _renderProxy: function () {
      var e = this.element; var i = this.options
      if (this.elementOffset = e.offset(), this._helper) {
        this.helper = this.helper || t('<div style="overflow:hidden;"></div>')
        var s = t.ui.ie6 ? 1 : 0; var n = t.ui.ie6 ? 2 : -1
        this.helper.addClass(this._helper).css({
          width: this.element.outerWidth() + n,
          height: this.element.outerHeight() + n,
          position: 'absolute',
          left: this.elementOffset.left - s + 'px',
          top: this.elementOffset.top - s + 'px',
          zIndex: ++i.zIndex
        }), this.helper.appendTo('body').disableSelection()
      } else this.helper = this.element
    },
    _change: {
      e: function (t, e, i) {
        return { width: this.originalSize.width + e }
      },
      w: function (t, e, i) {
        var s = (this.options, this.originalSize)
        return { left: this.originalPosition.left + e, width: s.width - e }
      },
      n: function (t, e, i) {
        var s = (this.options, this.originalSize)
        return { top: this.originalPosition.top + i, height: s.height - i }
      },
      s: function (t, e, i) {
        return { height: this.originalSize.height + i }
      },
      se: function (e, i, s) {
        return t.extend(this._change.s.apply(this, arguments), this._change.e.apply(this, [e, i, s]))
      },
      sw: function (e, i, s) {
        return t.extend(this._change.s.apply(this, arguments), this._change.w.apply(this, [e, i, s]))
      },
      ne: function (e, i, s) {
        return t.extend(this._change.n.apply(this, arguments), this._change.e.apply(this, [e, i, s]))
      },
      nw: function (e, i, s) {
        return t.extend(this._change.n.apply(this, arguments), this._change.w.apply(this, [e, i, s]))
      }
    },
    _propagate: function (e, i) {
      t.ui.plugin.call(this, e, [i, this.ui()]), e != 'resize' && this._trigger(e, i, this.ui())
    },
    plugins: {},
    ui: function () {
      return {
        originalElement: this.originalElement,
        element: this.element,
        helper: this.helper,
        position: this.position,
        size: this.size,
        originalSize: this.originalSize,
        originalPosition: this.originalPosition
      }
    }
  }), t.ui.plugin.add('resizable', 'alsoResize', {
    start: function (e, i) {
      var s = t(this).data('resizable').options; var n = function (e) {
        t(e).each(function () {
          var e = t(this)
          e.data('resizable-alsoresize', {
            width: parseInt(e.width(), 10),
            height: parseInt(e.height(), 10),
            left: parseInt(e.css('left'), 10),
            top: parseInt(e.css('top'), 10)
          })
        })
      }
      typeof s.alsoResize !== 'object' || s.alsoResize.parentNode ? n(s.alsoResize) : s.alsoResize.length ? (s.alsoResize = s.alsoResize[0], n(s.alsoResize)) : t.each(s.alsoResize, function (t) {
        n(t)
      })
    },
    resize: function (e, i) {
      var s = t(this).data('resizable'); var n = s.options; var o = s.originalSize; var a = s.originalPosition; var r = {
        height: s.size.height - o.height || 0,
        width: s.size.width - o.width || 0,
        top: s.position.top - a.top || 0,
        left: s.position.left - a.left || 0
      }; var l = function (e, s) {
        t(e).each(function () {
          var e = t(this); var n = t(this).data('resizable-alsoresize'); var o = {}
          var a = s && s.length ? s : e.parents(i.originalElement[0]).length ? ['width', 'height'] : ['width', 'height', 'top', 'left']
          t.each(a, function (t, e) {
            var i = (n[e] || 0) + (r[e] || 0)
            i && i >= 0 && (o[e] = i || null)
          }), e.css(o)
        })
      }
      typeof n.alsoResize !== 'object' || n.alsoResize.nodeType ? l(n.alsoResize) : t.each(n.alsoResize, function (t, e) {
        l(t, e)
      })
    },
    stop: function (e, i) {
      t(this).removeData('resizable-alsoresize')
    }
  }), t.ui.plugin.add('resizable', 'animate', {
    stop: function (e, i) {
      var s = t(this).data('resizable'); var n = s.options; var o = s._proportionallyResizeElements
      var a = o.length && /textarea/i.test(o[0].nodeName)
      var r = a && t.ui.hasScroll(o[0], 'left') ? 0 : s.sizeDiff.height; var l = a ? 0 : s.sizeDiff.width
      var h = { width: s.size.width - l, height: s.size.height - r }
      var c = parseInt(s.element.css('left'), 10) + (s.position.left - s.originalPosition.left) || null
      var u = parseInt(s.element.css('top'), 10) + (s.position.top - s.originalPosition.top) || null
      s.element.animate(t.extend(h, u && c ? { top: u, left: c } : {}), {
        duration: n.animateDuration,
        easing: n.animateEasing,
        step: function () {
          var i = {
            width: parseInt(s.element.css('width'), 10),
            height: parseInt(s.element.css('height'), 10),
            top: parseInt(s.element.css('top'), 10),
            left: parseInt(s.element.css('left'), 10)
          }
          o && o.length && t(o[0]).css({
            width: i.width,
            height: i.height
          }), s._updateCache(i), s._propagate('resize', e)
        }
      })
    }
  }), t.ui.plugin.add('resizable', 'containment', {
    start: function (e, s) {
      var n = t(this).data('resizable'); var o = n.options; var a = n.element; var r = o.containment
      var l = r instanceof t ? r.get(0) : /parent/.test(r) ? a.parent().get(0) : r
      if (l) {
        if (n.containerElement = t(l), /document/.test(r) || r == document) {
          n.containerOffset = {
            left: 0,
            top: 0
          }, n.containerPosition = { left: 0, top: 0 }, n.parentData = {
            element: t(document),
            left: 0,
            top: 0,
            width: t(document).width(),
            height: t(document).height() || document.body.parentNode.scrollHeight
          }
        } else {
          var h = t(l); var c = []
          t(['Top', 'Right', 'Left', 'Bottom']).each(function (t, e) {
            c[t] = i(h.css('padding' + e))
          }), n.containerOffset = h.offset(), n.containerPosition = h.position(), n.containerSize = {
            height: h.innerHeight() - c[3],
            width: h.innerWidth() - c[1]
          }
          var u = n.containerOffset; var d = n.containerSize.height; var p = n.containerSize.width
          var f = t.ui.hasScroll(l, 'left') ? l.scrollWidth : p; var g = t.ui.hasScroll(l) ? l.scrollHeight : d
          n.parentData = { element: l, left: u.left, top: u.top, width: f, height: g }
        }
      }
    },
    resize: function (e, i) {
      var s = t(this).data('resizable'); var n = s.options; var o = (s.containerSize, s.containerOffset)
      var a = (s.size, s.position); var r = s._aspectRatio || e.shiftKey; var l = { top: 0, left: 0 }
      var h = s.containerElement
      h[0] != document && /static/.test(h.css('position')) && (l = o), a.left < (s._helper ? o.left : 0) && (s.size.width = s.size.width + (s._helper ? s.position.left - o.left : s.position.left - l.left), r && (s.size.height = s.size.width / s.aspectRatio), s.position.left = n.helper ? o.left : 0), a.top < (s._helper ? o.top : 0) && (s.size.height = s.size.height + (s._helper ? s.position.top - o.top : s.position.top), r && (s.size.width = s.size.height * s.aspectRatio), s.position.top = s._helper ? o.top : 0), s.offset.left = s.parentData.left + s.position.left, s.offset.top = s.parentData.top + s.position.top
      var c = Math.abs((s._helper, s.offset.left - l.left + s.sizeDiff.width))
      var u = Math.abs((s._helper ? s.offset.top - l.top : s.offset.top - o.top) + s.sizeDiff.height)
      var d = s.containerElement.get(0) == s.element.parent().get(0)
      var p = /relative|absolute/.test(s.containerElement.css('position'))
      d && p && (c -= s.parentData.left), c + s.size.width >= s.parentData.width && (s.size.width = s.parentData.width - c, r && (s.size.height = s.size.width / s.aspectRatio)), u + s.size.height >= s.parentData.height && (s.size.height = s.parentData.height - u, r && (s.size.width = s.size.height * s.aspectRatio))
    },
    stop: function (e, i) {
      var s = t(this).data('resizable'); var n = s.options; var o = (s.position, s.containerOffset)
      var a = s.containerPosition; var r = s.containerElement; var l = t(s.helper); var h = l.offset()
      var c = l.outerWidth() - s.sizeDiff.width; var u = l.outerHeight() - s.sizeDiff.height
      s._helper && !n.animate && /relative/.test(r.css('position')) && t(this).css({
        left: h.left - a.left - o.left,
        width: c,
        height: u
      }), s._helper && !n.animate && /static/.test(r.css('position')) && t(this).css({
        left: h.left - a.left - o.left,
        width: c,
        height: u
      })
    }
  }), t.ui.plugin.add('resizable', 'ghost', {
    start: function (e, i) {
      var s = t(this).data('resizable'); var n = s.options; var o = s.size
      s.ghost = s.originalElement.clone(), s.ghost.css({
        opacity: 0.25,
        display: 'block',
        position: 'relative',
        height: o.height,
        width: o.width,
        margin: 0,
        left: 0,
        top: 0
      }).addClass('ui-resizable-ghost').addClass(typeof n.ghost === 'string' ? n.ghost : ''), s.ghost.appendTo(s.helper)
    },
    resize: function (e, i) {
      var s = t(this).data('resizable')
      s.options, s.ghost && s.ghost.css({ position: 'relative', height: s.size.height, width: s.size.width })
    },
    stop: function (e, i) {
      var s = t(this).data('resizable')
      s.options, s.ghost && s.helper && s.helper.get(0).removeChild(s.ghost.get(0))
    }
  }), t.ui.plugin.add('resizable', 'grid', {
    resize: function (e, i) {
      var s = t(this).data('resizable'); var n = s.options; var o = s.size; var a = s.originalSize; var r = s.originalPosition
      var l = s.axis
      n._aspectRatio || e.shiftKey, n.grid = typeof n.grid === 'number' ? [n.grid, n.grid] : n.grid
      var h = Math.round((o.width - a.width) / (n.grid[0] || 1)) * (n.grid[0] || 1)
      var c = Math.round((o.height - a.height) / (n.grid[1] || 1)) * (n.grid[1] || 1);
      /^(se|s|e)$/.test(l) ? (s.size.width = a.width + h, s.size.height = a.height + c) : /^(ne)$/.test(l) ? (s.size.width = a.width + h, s.size.height = a.height + c, s.position.top = r.top - c) : /^(sw)$/.test(l) ? (s.size.width = a.width + h, s.size.height = a.height + c, s.position.left = r.left - h) : (s.size.width = a.width + h, s.size.height = a.height + c, s.position.top = r.top - c, s.position.left = r.left - h)
    }
  })
  var i = function (t) {
    return parseInt(t, 10) || 0
  }; var s = function (t) {
    return !isNaN(parseInt(t, 10))
  }
}(jQuery)), (function (t, e) {
  t.widget('ui.selectable', t.ui.mouse, {
    version: '1.9.2',
    options: { appendTo: 'body', autoRefresh: !0, distance: 0, filter: '*', tolerance: 'touch' },
    _create: function () {
      var e = this
      this.element.addClass('ui-selectable'), this.dragged = !1
      var i
      this.refresh = function () {
        (i = t(e.options.filter, e.element[0])).addClass('ui-selectee'), i.each(function () {
          var e = t(this); var i = e.offset()
          t.data(this, 'selectable-item', {
            element: this,
            $element: e,
            left: i.left,
            top: i.top,
            right: i.left + e.outerWidth(),
            bottom: i.top + e.outerHeight(),
            startselected: !1,
            selected: e.hasClass('ui-selected'),
            selecting: e.hasClass('ui-selecting'),
            unselecting: e.hasClass('ui-unselecting')
          })
        })
      }, this.refresh(), this.selectees = i.addClass('ui-selectee'), this._mouseInit(), this.helper = t("<div class='ui-selectable-helper'></div>")
    },
    _destroy: function () {
      this.selectees.removeClass('ui-selectee').removeData('selectable-item'), this.element.removeClass('ui-selectable ui-selectable-disabled'), this._mouseDestroy()
    },
    _mouseStart: function (e) {
      var i = this
      if (this.opos = [e.pageX, e.pageY], !this.options.disabled) {
        var s = this.options
        this.selectees = t(s.filter, this.element[0]), this._trigger('start', e), t(s.appendTo).append(this.helper), this.helper.css({
          left: e.clientX,
          top: e.clientY,
          width: 0,
          height: 0
        }), s.autoRefresh && this.refresh(), this.selectees.filter('.ui-selected').each(function () {
          var s = t.data(this, 'selectable-item')
          s.startselected = !0, e.metaKey || e.ctrlKey || (s.$element.removeClass('ui-selected'), s.selected = !1, s.$element.addClass('ui-unselecting'), s.unselecting = !0, i._trigger('unselecting', e, { unselecting: s.element }))
        }), t(e.target).parents().addBack().each(function () {
          var s = t.data(this, 'selectable-item')
          if (s) {
            var n = !e.metaKey && !e.ctrlKey || !s.$element.hasClass('ui-selected')
            return s.$element.removeClass(n ? 'ui-unselecting' : 'ui-selected').addClass(n ? 'ui-selecting' : 'ui-unselecting'), s.unselecting = !n, s.selecting = n, s.selected = n, n ? i._trigger('selecting', e, { selecting: s.element }) : i._trigger('unselecting', e, { unselecting: s.element }), !1
          }
        })
      }
    },
    _mouseDrag: function (e) {
      var i = this
      if (this.dragged = !0, !this.options.disabled) {
        var s = this.options; var n = this.opos[0]; var o = this.opos[1]; var a = e.pageX; var r = e.pageY
        if (n > a) {
          l = a
          a = n, n = l
        }
        if (o > r) {
          var l = r
          r = o, o = l
        }
        return this.helper.css({
          left: n,
          top: o,
          width: a - n,
          height: r - o
        }), this.selectees.each(function () {
          var l = t.data(this, 'selectable-item')
          if (l && l.element != i.element[0]) {
            var h = !1
            s.tolerance == 'touch' ? h = !(l.left > a || l.right < n || l.top > r || l.bottom < o) : s.tolerance == 'fit' && (h = l.left > n && l.right < a && l.top > o && l.bottom < r), h ? (l.selected && (l.$element.removeClass('ui-selected'), l.selected = !1), l.unselecting && (l.$element.removeClass('ui-unselecting'), l.unselecting = !1), l.selecting || (l.$element.addClass('ui-selecting'), l.selecting = !0, i._trigger('selecting', e, { selecting: l.element }))) : (l.selecting && ((e.metaKey || e.ctrlKey) && l.startselected ? (l.$element.removeClass('ui-selecting'), l.selecting = !1, l.$element.addClass('ui-selected'), l.selected = !0) : (l.$element.removeClass('ui-selecting'), l.selecting = !1, l.startselected && (l.$element.addClass('ui-unselecting'), l.unselecting = !0), i._trigger('unselecting', e, { unselecting: l.element }))), l.selected && (e.metaKey || e.ctrlKey || l.startselected || (l.$element.removeClass('ui-selected'), l.selected = !1, l.$element.addClass('ui-unselecting'), l.unselecting = !0, i._trigger('unselecting', e, { unselecting: l.element }))))
          }
        }), !1
      }
    },
    _mouseStop: function (e) {
      var i = this
      return this.dragged = !1, this.options, t('.ui-unselecting', this.element[0]).each(function () {
        var s = t.data(this, 'selectable-item')
        s.$element.removeClass('ui-unselecting'), s.unselecting = !1, s.startselected = !1, i._trigger('unselected', e, { unselected: s.element })
      }), t('.ui-selecting', this.element[0]).each(function () {
        var s = t.data(this, 'selectable-item')
        s.$element.removeClass('ui-selecting').addClass('ui-selected'), s.selecting = !1, s.selected = !0, s.startselected = !0, i._trigger('selected', e, { selected: s.element })
      }), this._trigger('stop', e), this.helper.remove(), !1
    }
  })
}(jQuery)), (function (t, e) {
  t.widget('ui.slider', t.ui.mouse, {
    version: '1.9.2',
    widgetEventPrefix: 'slide',
    options: {
      animate: !1,
      distance: 0,
      max: 100,
      min: 0,
      orientation: 'horizontal',
      range: !1,
      step: 1,
      value: 0,
      values: null
    },
    _create: function () {
      var e; var i; var s = this.options
      var n = this.element.find('.ui-slider-handle').addClass('ui-state-default ui-corner-all'); var o = []
      for (this._keySliding = !1, this._mouseSliding = !1, this._animateOff = !0, this._handleIndex = null, this._detectOrientation(), this._mouseInit(), this.element.addClass('ui-slider ui-slider-' + this.orientation + ' ui-widget ui-widget-content ui-corner-all' + (s.disabled ? ' ui-slider-disabled ui-disabled' : '')), this.range = t([]), s.range && (!0 === s.range && (s.values || (s.values = [this._valueMin(), this._valueMin()]), s.values.length && s.values.length !== 2 && (s.values = [s.values[0], s.values[0]])), this.range = t('<div></div>').appendTo(this.element).addClass('ui-slider-range ui-widget-header' + (s.range === 'min' || s.range === 'max' ? ' ui-slider-range-' + s.range : ''))), i = s.values && s.values.length || 1, e = n.length; e < i; e++) o.push("<a class='ui-slider-handle ui-state-default ui-corner-all' href='#'></a>")
      this.handles = n.add(t(o.join('')).appendTo(this.element)), this.handle = this.handles.eq(0), this.handles.add(this.range).filter('a').click(function (t) {
        t.preventDefault()
      }).mouseenter(function () {
        s.disabled || t(this).addClass('ui-state-hover')
      }).mouseleave(function () {
        t(this).removeClass('ui-state-hover')
      }).focus(function () {
        s.disabled ? t(this).blur() : (t('.ui-slider .ui-state-focus').removeClass('ui-state-focus'), t(this).addClass('ui-state-focus'))
      }).blur(function () {
        t(this).removeClass('ui-state-focus')
      }), this.handles.each(function (e) {
        t(this).data('ui-slider-handle-index', e)
      }), this._on(this.handles, {
        keydown: function (e) {
          var i; var s; var n; var o = t(e.target).data('ui-slider-handle-index')
          switch (e.keyCode) {
            case t.ui.keyCode.HOME:
            case t.ui.keyCode.END:
            case t.ui.keyCode.PAGE_UP:
            case t.ui.keyCode.PAGE_DOWN:
            case t.ui.keyCode.UP:
            case t.ui.keyCode.RIGHT:
            case t.ui.keyCode.DOWN:
            case t.ui.keyCode.LEFT:
              if (e.preventDefault(), !this._keySliding && (this._keySliding = !0, t(e.target).addClass('ui-state-active'), !1 === this._start(e, o))) return
          }
          switch (n = this.options.step, i = s = this.options.values && this.options.values.length ? this.values(o) : this.value(), e.keyCode) {
            case t.ui.keyCode.HOME:
              s = this._valueMin()
              break
            case t.ui.keyCode.END:
              s = this._valueMax()
              break
            case t.ui.keyCode.PAGE_UP:
              s = this._trimAlignValue(i + (this._valueMax() - this._valueMin()) / 5)
              break
            case t.ui.keyCode.PAGE_DOWN:
              s = this._trimAlignValue(i - (this._valueMax() - this._valueMin()) / 5)
              break
            case t.ui.keyCode.UP:
            case t.ui.keyCode.RIGHT:
              if (i === this._valueMax()) return
              s = this._trimAlignValue(i + n)
              break
            case t.ui.keyCode.DOWN:
            case t.ui.keyCode.LEFT:
              if (i === this._valueMin()) return
              s = this._trimAlignValue(i - n)
          }
          this._slide(e, o, s)
        },
        keyup: function (e) {
          var i = t(e.target).data('ui-slider-handle-index')
          this._keySliding && (this._keySliding = !1, this._stop(e, i), this._change(e, i), t(e.target).removeClass('ui-state-active'))
        }
      }), this._refreshValue(), this._animateOff = !1
    },
    _destroy: function () {
      this.handles.remove(), this.range.remove(), this.element.removeClass('ui-slider ui-slider-horizontal ui-slider-vertical ui-slider-disabled ui-widget ui-widget-content ui-corner-all'), this._mouseDestroy()
    },
    _mouseCapture: function (e) {
      var i; var s; var n; var o; var a; var r; var l; var h = this; var c = this.options
      return !c.disabled && (this.elementSize = {
        width: this.element.outerWidth(),
        height: this.element.outerHeight()
      }, this.elementOffset = this.element.offset(), i = {
        x: e.pageX,
        y: e.pageY
      }, s = this._normValueFromMouse(i), n = this._valueMax() - this._valueMin() + 1, this.handles.each(function (e) {
        var i = Math.abs(s - h.values(e))
        n > i && (n = i, o = t(this), a = e)
      }), !0 === c.range && this.values(1) === c.min && (a += 1, o = t(this.handles[a])), !1 !== this._start(e, a) && (this._mouseSliding = !0, this._handleIndex = a, o.addClass('ui-state-active').focus(), r = o.offset(), l = !t(e.target).parents().addBack().is('.ui-slider-handle'), this._clickOffset = l ? {
        left: 0,
        top: 0
      } : {
        left: e.pageX - r.left - o.width() / 2,
        top: e.pageY - r.top - o.height() / 2 - (parseInt(o.css('borderTopWidth'), 10) || 0) - (parseInt(o.css('borderBottomWidth'), 10) || 0) + (parseInt(o.css('marginTop'), 10) || 0)
      }, this.handles.hasClass('ui-state-hover') || this._slide(e, a, s), this._animateOff = !0, !0))
    },
    _mouseStart: function () {
      return !0
    },
    _mouseDrag: function (t) {
      var e = { x: t.pageX, y: t.pageY }; var i = this._normValueFromMouse(e)
      return this._slide(t, this._handleIndex, i), !1
    },
    _mouseStop: function (t) {
      return this.handles.removeClass('ui-state-active'), this._mouseSliding = !1, this._stop(t, this._handleIndex), this._change(t, this._handleIndex), this._handleIndex = null, this._clickOffset = null, this._animateOff = !1, !1
    },
    _detectOrientation: function () {
      this.orientation = this.options.orientation === 'vertical' ? 'vertical' : 'horizontal'
    },
    _normValueFromMouse: function (t) {
      var e, i, s, n, o
      return this.orientation === 'horizontal' ? (e = this.elementSize.width, i = t.x - this.elementOffset.left - (this._clickOffset ? this._clickOffset.left : 0)) : (e = this.elementSize.height, i = t.y - this.elementOffset.top - (this._clickOffset ? this._clickOffset.top : 0)), (s = i / e) > 1 && (s = 1), s < 0 && (s = 0), this.orientation === 'vertical' && (s = 1 - s), n = this._valueMax() - this._valueMin(), o = this._valueMin() + s * n, this._trimAlignValue(o)
    },
    _start: function (t, e) {
      var i = { handle: this.handles[e], value: this.value() }
      return this.options.values && this.options.values.length && (i.value = this.values(e), i.values = this.values()), this._trigger('start', t, i)
    },
    _slide: function (t, e, i) {
      var s, n, o
      this.options.values && this.options.values.length ? (s = this.values(e ? 0 : 1), this.options.values.length === 2 && !0 === this.options.range && (e === 0 && i > s || e === 1 && i < s) && (i = s), i !== this.values(e) && (n = this.values(), n[e] = i, o = this._trigger('slide', t, {
        handle: this.handles[e],
        value: i,
        values: n
      }), s = this.values(e ? 0 : 1), !1 !== o && this.values(e, i, !0))) : i !== this.value() && !1 !== (o = this._trigger('slide', t, {
        handle: this.handles[e],
        value: i
      })) && this.value(i)
    },
    _stop: function (t, e) {
      var i = { handle: this.handles[e], value: this.value() }
      this.options.values && this.options.values.length && (i.value = this.values(e), i.values = this.values()), this._trigger('stop', t, i)
    },
    _change: function (t, e) {
      if (!this._keySliding && !this._mouseSliding) {
        var i = { handle: this.handles[e], value: this.value() }
        this.options.values && this.options.values.length && (i.value = this.values(e), i.values = this.values()), this._trigger('change', t, i)
      }
    },
    value: function (t) {
      return arguments.length ? (this.options.value = this._trimAlignValue(t), this._refreshValue(), void this._change(null, 0)) : this._value()
    },
    values: function (e, i) {
      var s, n, o
      if (arguments.length > 1) return this.options.values[e] = this._trimAlignValue(i), this._refreshValue(), void this._change(null, e)
      if (!arguments.length) return this._values()
      if (!t.isArray(arguments[0])) return this.options.values && this.options.values.length ? this._values(e) : this.value()
      for (s = this.options.values, n = arguments[0], o = 0; o < s.length; o += 1) s[o] = this._trimAlignValue(n[o]), this._change(null, o)
      this._refreshValue()
    },
    _setOption: function (e, i) {
      var s; var n = 0
      switch (t.isArray(this.options.values) && (n = this.options.values.length), t.Widget.prototype._setOption.apply(this, arguments), e) {
        case 'disabled':
          i ? (this.handles.filter('.ui-state-focus').blur(), this.handles.removeClass('ui-state-hover'), this.handles.prop('disabled', !0), this.element.addClass('ui-disabled')) : (this.handles.prop('disabled', !1), this.element.removeClass('ui-disabled'))
          break
        case 'orientation':
          this._detectOrientation(), this.element.removeClass('ui-slider-horizontal ui-slider-vertical').addClass('ui-slider-' + this.orientation), this._refreshValue()
          break
        case 'value':
          this._animateOff = !0, this._refreshValue(), this._change(null, 0), this._animateOff = !1
          break
        case 'values':
          for (this._animateOff = !0, this._refreshValue(), s = 0; s < n; s += 1) this._change(null, s)
          this._animateOff = !1
          break
        case 'min':
        case 'max':
          this._animateOff = !0, this._refreshValue(), this._animateOff = !1
      }
    },
    _value: function () {
      var t = this.options.value
      return t = this._trimAlignValue(t)
    },
    _values: function (t) {
      var e, i, s
      if (arguments.length) return e = this.options.values[t], e = this._trimAlignValue(e)
      for (i = this.options.values.slice(), s = 0; s < i.length; s += 1) i[s] = this._trimAlignValue(i[s])
      return i
    },
    _trimAlignValue: function (t) {
      if (t <= this._valueMin()) return this._valueMin()
      if (t >= this._valueMax()) return this._valueMax()
      var e = this.options.step > 0 ? this.options.step : 1; var i = (t - this._valueMin()) % e; var s = t - i
      return 2 * Math.abs(i) >= e && (s += i > 0 ? e : -e), parseFloat(s.toFixed(5))
    },
    _valueMin: function () {
      return this.options.min
    },
    _valueMax: function () {
      return this.options.max
    },
    _refreshValue: function () {
      var e; var i; var s; var n; var o; var a = this.options.range; var r = this.options; var l = this; var h = !this._animateOff && r.animate
      var c = {}
      this.options.values && this.options.values.length ? this.handles.each(function (s) {
        i = (l.values(s) - l._valueMin()) / (l._valueMax() - l._valueMin()) * 100, c[l.orientation === 'horizontal' ? 'left' : 'bottom'] = i + '%', t(this).stop(1, 1)[h ? 'animate' : 'css'](c, r.animate), !0 === l.options.range && (l.orientation === 'horizontal' ? (s === 0 && l.range.stop(1, 1)[h ? 'animate' : 'css']({ left: i + '%' }, r.animate), s === 1 && l.range[h ? 'animate' : 'css']({ width: i - e + '%' }, {
          queue: !1,
          duration: r.animate
        })) : (s === 0 && l.range.stop(1, 1)[h ? 'animate' : 'css']({ bottom: i + '%' }, r.animate), s === 1 && l.range[h ? 'animate' : 'css']({ height: i - e + '%' }, {
          queue: !1,
          duration: r.animate
        }))), e = i
      }) : (s = this.value(), n = this._valueMin(), o = this._valueMax(), i = o !== n ? (s - n) / (o - n) * 100 : 0, c[this.orientation === 'horizontal' ? 'left' : 'bottom'] = i + '%', this.handle.stop(1, 1)[h ? 'animate' : 'css'](c, r.animate), a === 'min' && this.orientation === 'horizontal' && this.range.stop(1, 1)[h ? 'animate' : 'css']({ width: i + '%' }, r.animate), a === 'max' && this.orientation === 'horizontal' && this.range[h ? 'animate' : 'css']({ width: 100 - i + '%' }, {
        queue: !1,
        duration: r.animate
      }), a === 'min' && this.orientation === 'vertical' && this.range.stop(1, 1)[h ? 'animate' : 'css']({ height: i + '%' }, r.animate), a === 'max' && this.orientation === 'vertical' && this.range[h ? 'animate' : 'css']({ height: 100 - i + '%' }, {
        queue: !1,
        duration: r.animate
      }))
    }
  })
}(jQuery)), (function (t, e) {
  t.widget('ui.sortable', t.ui.mouse, {
    version: '1.9.2',
    widgetEventPrefix: 'sort',
    ready: !1,
    options: {
      appendTo: 'parent',
      axis: !1,
      connectWith: !1,
      containment: !1,
      cursor: 'auto',
      cursorAt: !1,
      dropOnEmpty: !0,
      forcePlaceholderSize: !1,
      forceHelperSize: !1,
      grid: !1,
      handle: !1,
      helper: 'original',
      items: '> *',
      opacity: !1,
      placeholder: !1,
      revert: !1,
      scroll: !0,
      scrollSensitivity: 20,
      scrollSpeed: 20,
      scope: 'default',
      tolerance: 'intersect',
      zIndex: 1e3
    },
    _create: function () {
      var t = this.options
      this.containerCache = {}, this.element.addClass('ui-sortable'), this.refresh(), this.floating = !!this.items.length && (t.axis === 'x' || /left|right/.test(this.items[0].item.css('float')) || /inline|table-cell/.test(this.items[0].item.css('display'))), this.offset = this.element.offset(), this._mouseInit(), this.ready = !0
    },
    _destroy: function () {
      this.element.removeClass('ui-sortable ui-sortable-disabled'), this._mouseDestroy()
      for (var t = this.items.length - 1; t >= 0; t--) this.items[t].item.removeData(this.widgetName + '-item')
      return this
    },
    _setOption: function (e, i) {
      e === 'disabled' ? (this.options[e] = i, this.widget().toggleClass('ui-sortable-disabled', !!i)) : t.Widget.prototype._setOption.apply(this, arguments)
    },
    _mouseCapture: function (e, i) {
      var s = this
      if (this.reverting) return !1
      if (this.options.disabled || this.options.type == 'static') return !1
      this._refreshItems(e)
      var n = null
      if (t(e.target).parents().each(function () {
        if (t.data(this, s.widgetName + '-item') == s) return n = t(this), !1
      }), t.data(e.target, s.widgetName + '-item') == s && (n = t(e.target)), !n) return !1
      if (this.options.handle && !i) {
        var o = !1
        if (t(this.options.handle, n).find('*').addBack().each(function () {
          this == e.target && (o = !0)
        }), !o) return !1
      }
      return this.currentItem = n, this._removeCurrentsFromItems(), !0
    },
    _mouseStart: function (e, i, s) {
      var n = this.options
      if (this.currentContainer = this, this.refreshPositions(), this.helper = this._createHelper(e), this._cacheHelperProportions(), this._cacheMargins(), this.scrollParent = this.helper.scrollParent(), this.offset = this.currentItem.offset(), this.offset = {
        top: this.offset.top - this.margins.top,
        left: this.offset.left - this.margins.left
      }, t.extend(this.offset, {
        click: { left: e.pageX - this.offset.left, top: e.pageY - this.offset.top },
        parent: this._getParentOffset(),
        relative: this._getRelativeOffset()
      }), this.helper.css('position', 'absolute'), this.cssPosition = this.helper.css('position'), this.originalPosition = this._generatePosition(e), this.originalPageX = e.pageX, this.originalPageY = e.pageY, n.cursorAt && this._adjustOffsetFromHelper(n.cursorAt), this.domPosition = {
        prev: this.currentItem.prev()[0],
        parent: this.currentItem.parent()[0]
      }, this.helper[0] != this.currentItem[0] && this.currentItem.hide(), this._createPlaceholder(), n.containment && this._setContainment(), n.cursor && (t('body').css('cursor') && (this._storedCursor = t('body').css('cursor')), t('body').css('cursor', n.cursor)), n.opacity && (this.helper.css('opacity') && (this._storedOpacity = this.helper.css('opacity')), this.helper.css('opacity', n.opacity)), n.zIndex && (this.helper.css('zIndex') && (this._storedZIndex = this.helper.css('zIndex')), this.helper.css('zIndex', n.zIndex)), this.scrollParent[0] != document && this.scrollParent[0].tagName != 'HTML' && (this.overflowOffset = this.scrollParent.offset()), this._trigger('start', e, this._uiHash()), this._preserveHelperProportions || this._cacheHelperProportions(), !s) for (var o = this.containers.length - 1; o >= 0; o--) this.containers[o]._trigger('activate', e, this._uiHash(this))
      return t.ui.ddmanager && (t.ui.ddmanager.current = this), t.ui.ddmanager && !n.dropBehaviour && t.ui.ddmanager.prepareOffsets(this, e), this.dragging = !0, this.helper.addClass('ui-sortable-helper'), this._mouseDrag(e), !0
    },
    _mouseDrag: function (e) {
      if (this.position = this._generatePosition(e), this.positionAbs = this._convertPositionTo('absolute'), this.lastPositionAbs || (this.lastPositionAbs = this.positionAbs), this.options.scroll) {
        var i = this.options; var s = !1
        this.scrollParent[0] != document && this.scrollParent[0].tagName != 'HTML' ? (this.overflowOffset.top + this.scrollParent[0].offsetHeight - e.pageY < i.scrollSensitivity ? this.scrollParent[0].scrollTop = s = this.scrollParent[0].scrollTop + i.scrollSpeed : e.pageY - this.overflowOffset.top < i.scrollSensitivity && (this.scrollParent[0].scrollTop = s = this.scrollParent[0].scrollTop - i.scrollSpeed), this.overflowOffset.left + this.scrollParent[0].offsetWidth - e.pageX < i.scrollSensitivity ? this.scrollParent[0].scrollLeft = s = this.scrollParent[0].scrollLeft + i.scrollSpeed : e.pageX - this.overflowOffset.left < i.scrollSensitivity && (this.scrollParent[0].scrollLeft = s = this.scrollParent[0].scrollLeft - i.scrollSpeed)) : (e.pageY - t(document).scrollTop() < i.scrollSensitivity ? s = t(document).scrollTop(t(document).scrollTop() - i.scrollSpeed) : t(window).height() - (e.pageY - t(document).scrollTop()) < i.scrollSensitivity && (s = t(document).scrollTop(t(document).scrollTop() + i.scrollSpeed)), e.pageX - t(document).scrollLeft() < i.scrollSensitivity ? s = t(document).scrollLeft(t(document).scrollLeft() - i.scrollSpeed) : t(window).width() - (e.pageX - t(document).scrollLeft()) < i.scrollSensitivity && (s = t(document).scrollLeft(t(document).scrollLeft() + i.scrollSpeed))), !1 !== s && t.ui.ddmanager && !i.dropBehaviour && t.ui.ddmanager.prepareOffsets(this, e)
      }
      this.positionAbs = this._convertPositionTo('absolute'), this.options.axis && this.options.axis == 'y' || (this.helper[0].style.left = this.position.left + 'px'), this.options.axis && this.options.axis == 'x' || (this.helper[0].style.top = this.position.top + 'px')
      for (var n = this.items.length - 1; n >= 0; n--) {
        var o = this.items[n]; var a = o.item[0]; var r = this._intersectsWithPointer(o)
        if (r && o.instance === this.currentContainer && !(a == this.currentItem[0] || this.placeholder[r == 1 ? 'next' : 'prev']()[0] == a || t.contains(this.placeholder[0], a) || this.options.type == 'semi-dynamic' && t.contains(this.element[0], a))) {
          if (this.direction = r == 1 ? 'down' : 'up', this.options.tolerance != 'pointer' && !this._intersectsWithSides(o)) break
          this._rearrange(e, o), this._trigger('change', e, this._uiHash())
          break
        }
      }
      return this._contactContainers(e), t.ui.ddmanager && t.ui.ddmanager.drag(this, e), this._trigger('sort', e, this._uiHash()), this.lastPositionAbs = this.positionAbs, !1
    },
    _mouseStop: function (e, i) {
      if (e) {
        if (t.ui.ddmanager && !this.options.dropBehaviour && t.ui.ddmanager.drop(this, e), this.options.revert) {
          var s = this; var n = this.placeholder.offset()
          this.reverting = !0, t(this.helper).animate({
            left: n.left - this.offset.parent.left - this.margins.left + (this.offsetParent[0] == document.body ? 0 : this.offsetParent[0].scrollLeft),
            top: n.top - this.offset.parent.top - this.margins.top + (this.offsetParent[0] == document.body ? 0 : this.offsetParent[0].scrollTop)
          }, parseInt(this.options.revert, 10) || 500, function () {
            s._clear(e)
          })
        } else this._clear(e, i)
        return !1
      }
    },
    cancel: function () {
      if (this.dragging) {
        this._mouseUp({ target: null }), this.options.helper == 'original' ? this.currentItem.css(this._storedCSS).removeClass('ui-sortable-helper') : this.currentItem.show()
        for (var e = this.containers.length - 1; e >= 0; e--) this.containers[e]._trigger('deactivate', null, this._uiHash(this)), this.containers[e].containerCache.over && (this.containers[e]._trigger('out', null, this._uiHash(this)), this.containers[e].containerCache.over = 0)
      }
      return this.placeholder && (this.placeholder[0].parentNode && this.placeholder[0].parentNode.removeChild(this.placeholder[0]), this.options.helper != 'original' && this.helper && this.helper[0].parentNode && this.helper.remove(), t.extend(this, {
        helper: null,
        dragging: !1,
        reverting: !1,
        _noFinalSort: null
      }), this.domPosition.prev ? t(this.domPosition.prev).after(this.currentItem) : t(this.domPosition.parent).prepend(this.currentItem)), this
    },
    serialize: function (e) {
      var i = this._getItemsAsjQuery(e && e.connected); var s = []
      return e = e || {}, t(i).each(function () {
        var i = (t(e.item || this).attr(e.attribute || 'id') || '').match(e.expression || /(.+)[-=_](.+)/)
        i && s.push((e.key || i[1] + '[]') + '=' + (e.key && e.expression ? i[1] : i[2]))
      }), !s.length && e.key && s.push(e.key + '='), s.join('&')
    },
    toArray: function (e) {
      var i = this._getItemsAsjQuery(e && e.connected); var s = []
      return e = e || {}, i.each(function () {
        s.push(t(e.item || this).attr(e.attribute || 'id') || '')
      }), s
    },
    _intersectsWith: function (t) {
      var e = this.positionAbs.left; var i = e + this.helperProportions.width; var s = this.positionAbs.top
      var n = s + this.helperProportions.height; var o = t.left; var a = o + t.width; var r = t.top; var l = r + t.height
      var h = this.offset.click.top; var c = this.offset.click.left
      var u = s + h > r && s + h < l && e + c > o && e + c < a
      return this.options.tolerance == 'pointer' || this.options.forcePointerForContainers || this.options.tolerance != 'pointer' && this.helperProportions[this.floating ? 'width' : 'height'] > t[this.floating ? 'width' : 'height'] ? u : o < e + this.helperProportions.width / 2 && i - this.helperProportions.width / 2 < a && r < s + this.helperProportions.height / 2 && n - this.helperProportions.height / 2 < l
    },
    _intersectsWithPointer: function (e) {
      var i = this.options.axis === 'x' || t.ui.isOverAxis(this.positionAbs.top + this.offset.click.top, e.top, e.height)
      var s = this.options.axis === 'y' || t.ui.isOverAxis(this.positionAbs.left + this.offset.click.left, e.left, e.width)
      var n = i && s; var o = this._getDragVerticalDirection(); var a = this._getDragHorizontalDirection()
      return !!n && (this.floating ? a && a == 'right' || o == 'down' ? 2 : 1 : o && (o == 'down' ? 2 : 1))
    },
    _intersectsWithSides: function (e) {
      var i = t.ui.isOverAxis(this.positionAbs.top + this.offset.click.top, e.top + e.height / 2, e.height)
      var s = t.ui.isOverAxis(this.positionAbs.left + this.offset.click.left, e.left + e.width / 2, e.width)
      var n = this._getDragVerticalDirection(); var o = this._getDragHorizontalDirection()
      return this.floating && o ? o == 'right' && s || o == 'left' && !s : n && (n == 'down' && i || n == 'up' && !i)
    },
    _getDragVerticalDirection: function () {
      var t = this.positionAbs.top - this.lastPositionAbs.top
      return t != 0 && (t > 0 ? 'down' : 'up')
    },
    _getDragHorizontalDirection: function () {
      var t = this.positionAbs.left - this.lastPositionAbs.left
      return t != 0 && (t > 0 ? 'right' : 'left')
    },
    refresh: function (t) {
      return this._refreshItems(t), this.refreshPositions(), this
    },
    _connectWith: function () {
      var t = this.options
      return t.connectWith.constructor == String ? [t.connectWith] : t.connectWith
    },
    _getItemsAsjQuery: function (e) {
      var i = []; var s = []; var n = this._connectWith()
      if (n && e) {
        for (l = n.length - 1; l >= 0; l--) {
          for (var o = t(n[l]), a = o.length - 1; a >= 0; a--) {
            var r = t.data(o[a], this.widgetName)
            r && r != this && !r.options.disabled && s.push([t.isFunction(r.options.items) ? r.options.items.call(r.element) : t(r.options.items, r.element).not('.ui-sortable-helper').not('.ui-sortable-placeholder'), r])
          }
        }
      }
      s.push([t.isFunction(this.options.items) ? this.options.items.call(this.element, null, {
        options: this.options,
        item: this.currentItem
      }) : t(this.options.items, this.element).not('.ui-sortable-helper').not('.ui-sortable-placeholder'), this])
      for (var l = s.length - 1; l >= 0; l--) {
        s[l][0].each(function () {
          i.push(this)
        })
      }
      return t(i)
    },
    _removeCurrentsFromItems: function () {
      var e = this.currentItem.find(':data(' + this.widgetName + '-item)')
      this.items = t.grep(this.items, function (t) {
        for (var i = 0; i < e.length; i++) if (e[i] == t.item[0]) return !1
        return !0
      })
    },
    _refreshItems: function (e) {
      this.items = [], this.containers = [this]
      var i = this.items
      var s = [[t.isFunction(this.options.items) ? this.options.items.call(this.element[0], e, { item: this.currentItem }) : t(this.options.items, this.element), this]]
      var n = this._connectWith()
      if (n && this.ready) {
        for (l = n.length - 1; l >= 0; l--) {
          for (var o = t(n[l]), a = o.length - 1; a >= 0; a--) {
            var r = t.data(o[a], this.widgetName)
            r && r != this && !r.options.disabled && (s.push([t.isFunction(r.options.items) ? r.options.items.call(r.element[0], e, { item: this.currentItem }) : t(r.options.items, r.element), r]), this.containers.push(r))
          }
        }
      }
      for (var l = s.length - 1; l >= 0; l--) {
        for (var h = s[l][1], c = s[l][0], a = 0, u = c.length; a < u; a++) {
          var d = t(c[a])
          d.data(this.widgetName + '-item', h), i.push({
            item: d,
            instance: h,
            width: 0,
            height: 0,
            left: 0,
            top: 0
          })
        }
      }
    },
    refreshPositions: function (e) {
      this.offsetParent && this.helper && (this.offset.parent = this._getParentOffset())
      for (n = this.items.length - 1; n >= 0; n--) {
        var i = this.items[n]
        if (i.instance == this.currentContainer || !this.currentContainer || i.item[0] == this.currentItem[0]) {
          var s = this.options.toleranceElement ? t(this.options.toleranceElement, i.item) : i.item
          e || (i.width = s.outerWidth(), i.height = s.outerHeight())
          o = s.offset()
          i.left = o.left, i.top = o.top
        }
      }
      if (this.options.custom && this.options.custom.refreshContainers) this.options.custom.refreshContainers.call(this); else {
        for (var n = this.containers.length - 1; n >= 0; n--) {
          var o = this.containers[n].element.offset()
          this.containers[n].containerCache.left = o.left, this.containers[n].containerCache.top = o.top, this.containers[n].containerCache.width = this.containers[n].element.outerWidth(), this.containers[n].containerCache.height = this.containers[n].element.outerHeight()
        }
      }
      return this
    },
    _createPlaceholder: function (e) {
      var i = (e = e || this).options
      if (!i.placeholder || i.placeholder.constructor == String) {
        var s = i.placeholder
        i.placeholder = {
          element: function () {
            var i = t(document.createElement(e.currentItem[0].nodeName)).addClass(s || e.currentItem[0].className + ' ui-sortable-placeholder').removeClass('ui-sortable-helper')[0]
            return s || (i.style.visibility = 'hidden'), i
          },
          update: function (t, n) {
            s && !i.forcePlaceholderSize || (n.height() || n.height(e.currentItem.innerHeight() - parseInt(e.currentItem.css('paddingTop') || 0, 10) - parseInt(e.currentItem.css('paddingBottom') || 0, 10)), n.width() || n.width(e.currentItem.innerWidth() - parseInt(e.currentItem.css('paddingLeft') || 0, 10) - parseInt(e.currentItem.css('paddingRight') || 0, 10)))
          }
        }
      }
      e.placeholder = t(i.placeholder.element.call(e.element, e.currentItem)), e.currentItem.after(e.placeholder), i.placeholder.update(e, e.placeholder)
    },
    _contactContainers: function (e) {
      for (var i = null, s = null, n = this.containers.length - 1; n >= 0; n--) {
        if (!t.contains(this.currentItem[0], this.containers[n].element[0])) {
          if (this._intersectsWith(this.containers[n].containerCache)) {
            if (i && t.contains(this.containers[n].element[0], i.element[0])) continue
            i = this.containers[n], s = n
          } else this.containers[n].containerCache.over && (this.containers[n]._trigger('out', e, this._uiHash(this)), this.containers[n].containerCache.over = 0)
        }
      }
      if (i) {
        if (this.containers.length === 1) this.containers[s]._trigger('over', e, this._uiHash(this)), this.containers[s].containerCache.over = 1; else {
          for (var o = 1e4, a = null, r = this.containers[s].floating ? 'left' : 'top', l = this.containers[s].floating ? 'width' : 'height', h = this.positionAbs[r] + this.offset.click[r], c = this.items.length - 1; c >= 0; c--) {
            if (t.contains(this.containers[s].element[0], this.items[c].item[0]) && this.items[c].item[0] != this.currentItem[0]) {
              var u = this.items[c].item.offset()[r]; var d = !1
              Math.abs(u - h) > Math.abs(u + this.items[c][l] - h) && (d = !0, u += this.items[c][l]), Math.abs(u - h) < o && (o = Math.abs(u - h), a = this.items[c], this.direction = d ? 'up' : 'down')
            }
          }
          if (!a && !this.options.dropOnEmpty) return
          this.currentContainer = this.containers[s], a ? this._rearrange(e, a, null, !0) : this._rearrange(e, null, this.containers[s].element, !0), this._trigger('change', e, this._uiHash()), this.containers[s]._trigger('change', e, this._uiHash(this)), this.options.placeholder.update(this.currentContainer, this.placeholder), this.containers[s]._trigger('over', e, this._uiHash(this)), this.containers[s].containerCache.over = 1
        }
      }
    },
    _createHelper: function (e) {
      var i = this.options
      var s = t.isFunction(i.helper) ? t(i.helper.apply(this.element[0], [e, this.currentItem])) : i.helper == 'clone' ? this.currentItem.clone() : this.currentItem
      return s.parents('body').length || t(i.appendTo != 'parent' ? i.appendTo : this.currentItem[0].parentNode)[0].appendChild(s[0]), s[0] == this.currentItem[0] && (this._storedCSS = {
        width: this.currentItem[0].style.width,
        height: this.currentItem[0].style.height,
        position: this.currentItem.css('position'),
        top: this.currentItem.css('top'),
        left: this.currentItem.css('left')
      }), (s[0].style.width == '' || i.forceHelperSize) && s.width(this.currentItem.width()), (s[0].style.height == '' || i.forceHelperSize) && s.height(this.currentItem.height()), s
    },
    _adjustOffsetFromHelper: function (e) {
      typeof e === 'string' && (e = e.split(' ')), t.isArray(e) && (e = {
        left: +e[0],
        top: +e[1] || 0
      }), 'left' in e && (this.offset.click.left = e.left + this.margins.left), 'right' in e && (this.offset.click.left = this.helperProportions.width - e.right + this.margins.left), 'top' in e && (this.offset.click.top = e.top + this.margins.top), 'bottom' in e && (this.offset.click.top = this.helperProportions.height - e.bottom + this.margins.top)
    },
    _getParentOffset: function () {
      this.offsetParent = this.helper.offsetParent()
      var e = this.offsetParent.offset()
      return this.cssPosition == 'absolute' && this.scrollParent[0] != document && t.contains(this.scrollParent[0], this.offsetParent[0]) && (e.left += this.scrollParent.scrollLeft(), e.top += this.scrollParent.scrollTop()), (this.offsetParent[0] == document.body || this.offsetParent[0].tagName && this.offsetParent[0].tagName.toLowerCase() == 'html' && t.ui.ie) && (e = {
        top: 0,
        left: 0
      }), {
        top: e.top + (parseInt(this.offsetParent.css('borderTopWidth'), 10) || 0),
        left: e.left + (parseInt(this.offsetParent.css('borderLeftWidth'), 10) || 0)
      }
    },
    _getRelativeOffset: function () {
      if (this.cssPosition == 'relative') {
        var t = this.currentItem.position()
        return {
          top: t.top - (parseInt(this.helper.css('top'), 10) || 0) + this.scrollParent.scrollTop(),
          left: t.left - (parseInt(this.helper.css('left'), 10) || 0) + this.scrollParent.scrollLeft()
        }
      }
      return { top: 0, left: 0 }
    },
    _cacheMargins: function () {
      this.margins = {
        left: parseInt(this.currentItem.css('marginLeft'), 10) || 0,
        top: parseInt(this.currentItem.css('marginTop'), 10) || 0
      }
    },
    _cacheHelperProportions: function () {
      this.helperProportions = { width: this.helper.outerWidth(), height: this.helper.outerHeight() }
    },
    _setContainment: function () {
      var e = this.options
      if (e.containment == 'parent' && (e.containment = this.helper[0].parentNode), e.containment != 'document' && e.containment != 'window' || (this.containment = [0 - this.offset.relative.left - this.offset.parent.left, 0 - this.offset.relative.top - this.offset.parent.top, t(e.containment == 'document' ? document : window).width() - this.helperProportions.width - this.margins.left, (t(e.containment == 'document' ? document : window).height() || document.body.parentNode.scrollHeight) - this.helperProportions.height - this.margins.top]), !/^(document|window|parent)$/.test(e.containment)) {
        var i = t(e.containment)[0]; var s = t(e.containment).offset(); var n = t(i).css('overflow') != 'hidden'
        this.containment = [s.left + (parseInt(t(i).css('borderLeftWidth'), 10) || 0) + (parseInt(t(i).css('paddingLeft'), 10) || 0) - this.margins.left, s.top + (parseInt(t(i).css('borderTopWidth'), 10) || 0) + (parseInt(t(i).css('paddingTop'), 10) || 0) - this.margins.top, s.left + (n ? Math.max(i.scrollWidth, i.offsetWidth) : i.offsetWidth) - (parseInt(t(i).css('borderLeftWidth'), 10) || 0) - (parseInt(t(i).css('paddingRight'), 10) || 0) - this.helperProportions.width - this.margins.left, s.top + (n ? Math.max(i.scrollHeight, i.offsetHeight) : i.offsetHeight) - (parseInt(t(i).css('borderTopWidth'), 10) || 0) - (parseInt(t(i).css('paddingBottom'), 10) || 0) - this.helperProportions.height - this.margins.top]
      }
    },
    _convertPositionTo: function (e, i) {
      i || (i = this.position)
      var s = e == 'absolute' ? 1 : -1
      var n = (this.options, this.cssPosition != 'absolute' || this.scrollParent[0] != document && t.contains(this.scrollParent[0], this.offsetParent[0]) ? this.scrollParent : this.offsetParent)
      var o = /(html|body)/i.test(n[0].tagName)
      return {
        top: i.top + this.offset.relative.top * s + this.offset.parent.top * s - (this.cssPosition == 'fixed' ? -this.scrollParent.scrollTop() : o ? 0 : n.scrollTop()) * s,
        left: i.left + this.offset.relative.left * s + this.offset.parent.left * s - (this.cssPosition == 'fixed' ? -this.scrollParent.scrollLeft() : o ? 0 : n.scrollLeft()) * s
      }
    },
    _generatePosition: function (e) {
      var i = this.options
      var s = this.cssPosition != 'absolute' || this.scrollParent[0] != document && t.contains(this.scrollParent[0], this.offsetParent[0]) ? this.scrollParent : this.offsetParent
      var n = /(html|body)/i.test(s[0].tagName)
      this.cssPosition != 'relative' || this.scrollParent[0] != document && this.scrollParent[0] != this.offsetParent[0] || (this.offset.relative = this._getRelativeOffset())
      var o = e.pageX; var a = e.pageY
      if (this.originalPosition && (this.containment && (e.pageX - this.offset.click.left < this.containment[0] && (o = this.containment[0] + this.offset.click.left), e.pageY - this.offset.click.top < this.containment[1] && (a = this.containment[1] + this.offset.click.top), e.pageX - this.offset.click.left > this.containment[2] && (o = this.containment[2] + this.offset.click.left), e.pageY - this.offset.click.top > this.containment[3] && (a = this.containment[3] + this.offset.click.top)), i.grid)) {
        var r = this.originalPageY + Math.round((a - this.originalPageY) / i.grid[1]) * i.grid[1]
        a = this.containment && (r - this.offset.click.top < this.containment[1] || r - this.offset.click.top > this.containment[3]) ? r - this.offset.click.top < this.containment[1] ? r + i.grid[1] : r - i.grid[1] : r
        var l = this.originalPageX + Math.round((o - this.originalPageX) / i.grid[0]) * i.grid[0]
        o = this.containment && (l - this.offset.click.left < this.containment[0] || l - this.offset.click.left > this.containment[2]) ? l - this.offset.click.left < this.containment[0] ? l + i.grid[0] : l - i.grid[0] : l
      }
      return {
        top: a - this.offset.click.top - this.offset.relative.top - this.offset.parent.top + (this.cssPosition == 'fixed' ? -this.scrollParent.scrollTop() : n ? 0 : s.scrollTop()),
        left: o - this.offset.click.left - this.offset.relative.left - this.offset.parent.left + (this.cssPosition == 'fixed' ? -this.scrollParent.scrollLeft() : n ? 0 : s.scrollLeft())
      }
    },
    _rearrange: function (t, e, i, s) {
      i ? i[0].appendChild(this.placeholder[0]) : e.item[0].parentNode.insertBefore(this.placeholder[0], this.direction == 'down' ? e.item[0] : e.item[0].nextSibling), this.counter = this.counter ? ++this.counter : 1
      var n = this.counter
      this._delay(function () {
        n == this.counter && this.refreshPositions(!s)
      })
    },
    _clear: function (e, i) {
      this.reverting = !1
      var s = []
      if (!this._noFinalSort && this.currentItem.parent().length && this.placeholder.before(this.currentItem), this._noFinalSort = null, this.helper[0] == this.currentItem[0]) {
        for (var n in this._storedCSS) this._storedCSS[n] != 'auto' && this._storedCSS[n] != 'static' || (this._storedCSS[n] = '')
        this.currentItem.css(this._storedCSS).removeClass('ui-sortable-helper')
      } else this.currentItem.show()
      this.fromOutside && !i && s.push(function (t) {
        this._trigger('receive', t, this._uiHash(this.fromOutside))
      }), !this.fromOutside && this.domPosition.prev == this.currentItem.prev().not('.ui-sortable-helper')[0] && this.domPosition.parent == this.currentItem.parent()[0] || i || s.push(function (t) {
        this._trigger('update', t, this._uiHash())
      }), this !== this.currentContainer && (i || (s.push(function (t) {
        this._trigger('remove', t, this._uiHash())
      }), s.push(function (t) {
        return function (e) {
          t._trigger('receive', e, this._uiHash(this))
        }
      }.call(this, this.currentContainer)), s.push(function (t) {
        return function (e) {
          t._trigger('update', e, this._uiHash(this))
        }
      }.call(this, this.currentContainer))))
      for (n = this.containers.length - 1; n >= 0; n--) {
        i || s.push(function (t) {
          return function (e) {
            t._trigger('deactivate', e, this._uiHash(this))
          }
        }.call(this, this.containers[n])), this.containers[n].containerCache.over && (s.push(function (t) {
          return function (e) {
            t._trigger('out', e, this._uiHash(this))
          }
        }.call(this, this.containers[n])), this.containers[n].containerCache.over = 0)
      }
      if (this._storedCursor && t('body').css('cursor', this._storedCursor), this._storedOpacity && this.helper.css('opacity', this._storedOpacity), this._storedZIndex && this.helper.css('zIndex', this._storedZIndex == 'auto' ? '' : this._storedZIndex), this.dragging = !1, this.cancelHelperRemoval) {
        if (!i) {
          this._trigger('beforeStop', e, this._uiHash())
          for (n = 0; n < s.length; n++) s[n].call(this, e)
          this._trigger('stop', e, this._uiHash())
        }
        return this.fromOutside = !1, !1
      }
      if (i || this._trigger('beforeStop', e, this._uiHash()), this.placeholder[0].parentNode.removeChild(this.placeholder[0]), this.helper[0] != this.currentItem[0] && this.helper.remove(), this.helper = null, !i) {
        for (n = 0; n < s.length; n++) s[n].call(this, e)
        this._trigger('stop', e, this._uiHash())
      }
      return this.fromOutside = !1, !0
    },
    _trigger: function () {
      !1 === t.Widget.prototype._trigger.apply(this, arguments) && this.cancel()
    },
    _uiHash: function (e) {
      var i = e || this
      return {
        helper: i.helper,
        placeholder: i.placeholder || t([]),
        position: i.position,
        originalPosition: i.originalPosition,
        offset: i.positionAbs,
        item: i.currentItem,
        sender: e ? e.element : null
      }
    }
  })
}(jQuery)), (function (t) {
  function e (t) {
    return function () {
      var e = this.element.val()
      t.apply(this, arguments), this._refresh(), e !== this.element.val() && this._trigger('change')
    }
  }
}(jQuery)), (function (t, e) {
  function i () {
    return ++n
  }

  function s (t) {
    return t.hash.length > 1 && t.href.replace(o, '') === location.href.replace(o, '').replace(/\s/g, '%20')
  }

  var n = 0; var o = /#.*$/
  t.widget('ui.tabs', {
    version: '1.9.2',
    delay: 300,
    options: {
      active: null,
      collapsible: !1,
      event: 'click',
      heightStyle: 'content',
      hide: null,
      show: null,
      activate: null,
      beforeActivate: null,
      beforeLoad: null,
      load: null
    },
    _create: function () {
      var e = this; var i = this.options; var s = i.active; var n = location.hash.substring(1)
      this.running = !1, this.element.addClass('ui-tabs ui-widget ui-widget-content ui-corner-all').toggleClass('ui-tabs-collapsible', i.collapsible).delegate('.ui-tabs-nav > li', 'mousedown' + this.eventNamespace, function (e) {
        t(this).is('.ui-state-disabled') && e.preventDefault()
      }).delegate('.ui-tabs-anchor', 'focus' + this.eventNamespace, function () {
        t(this).closest('li').is('.ui-state-disabled') && this.blur()
      }), this._processTabs(), s === null && (n && this.tabs.each(function (e, i) {
        if (t(i).attr('aria-controls') === n) return s = e, !1
      }), s === null && (s = this.tabs.index(this.tabs.filter('.ui-tabs-active'))), s !== null && s !== -1 || (s = !!this.tabs.length && 0)), !1 !== s && (s = this.tabs.index(this.tabs.eq(s))) === -1 && (s = !i.collapsible && 0), i.active = s, !i.collapsible && !1 === i.active && this.anchors.length && (i.active = 0), t.isArray(i.disabled) && (i.disabled = t.unique(i.disabled.concat(t.map(this.tabs.filter('.ui-state-disabled'), function (t) {
        return e.tabs.index(t)
      }))).sort()), !1 !== this.options.active && this.anchors.length ? this.active = this._findActive(this.options.active) : this.active = t(), this._refresh(), this.active.length && this.load(i.active)
    },
    _getCreateEventData: function () {
      return { tab: this.active, panel: this.active.length ? this._getPanelForTab(this.active) : t() }
    },
    _tabKeydown: function (e) {
      var i = t(this.document[0].activeElement).closest('li'); var s = this.tabs.index(i); var n = !0
      if (!this._handlePageNav(e)) {
        switch (e.keyCode) {
          case t.ui.keyCode.RIGHT:
          case t.ui.keyCode.DOWN:
            s++
            break
          case t.ui.keyCode.UP:
          case t.ui.keyCode.LEFT:
            n = !1, s--
            break
          case t.ui.keyCode.END:
            s = this.anchors.length - 1
            break
          case t.ui.keyCode.HOME:
            s = 0
            break
          case t.ui.keyCode.SPACE:
            return e.preventDefault(), clearTimeout(this.activating), void this._activate(s)
          case t.ui.keyCode.ENTER:
            return e.preventDefault(), clearTimeout(this.activating), void this._activate(s !== this.options.active && s)
          default:
            return
        }
        e.preventDefault(), clearTimeout(this.activating), s = this._focusNextTab(s, n), e.ctrlKey || (i.attr('aria-selected', 'false'), this.tabs.eq(s).attr('aria-selected', 'true'), this.activating = this._delay(function () {
          this.option('active', s)
        }, this.delay))
      }
    },
    _panelKeydown: function (e) {
      this._handlePageNav(e) || e.ctrlKey && e.keyCode === t.ui.keyCode.UP && (e.preventDefault(), this.active.focus())
    },
    _handlePageNav: function (e) {
      return e.altKey && e.keyCode === t.ui.keyCode.PAGE_UP ? (this._activate(this._focusNextTab(this.options.active - 1, !1)), !0) : e.altKey && e.keyCode === t.ui.keyCode.PAGE_DOWN ? (this._activate(this._focusNextTab(this.options.active + 1, !0)), !0) : void 0
    },
    _findNextTab: function (e, i) {
      for (var s = this.tabs.length - 1; t.inArray((e > s && (e = 0), e < 0 && (e = s), e), this.options.disabled) !== -1;) e = i ? e + 1 : e - 1
      return e
    },
    _focusNextTab: function (t, e) {
      return t = this._findNextTab(t, e), this.tabs.eq(t).focus(), t
    },
    _setOption: function (t, e) {
      return t === 'active' ? void this._activate(e) : t === 'disabled' ? void this._setupDisabled(e) : (this._super(t, e), t === 'collapsible' && (this.element.toggleClass('ui-tabs-collapsible', e), e || !1 !== this.options.active || this._activate(0)), t === 'event' && this._setupEvents(e), void (t === 'heightStyle' && this._setupHeightStyle(e)))
    },
    _tabId: function (t) {
      return t.attr('aria-controls') || 'ui-tabs-' + i()
    },
    _sanitizeSelector: function (t) {
      return t ? t.replace(/[!"$%&'()*+,.\/:;<=>?@\[\]\^`{|}~]/g, '\\$&') : ''
    },
    refresh: function () {
      var e = this.options; var i = this.tablist.children(':has(a[href])')
      e.disabled = t.map(i.filter('.ui-state-disabled'), function (t) {
        return i.index(t)
      }), this._processTabs(), !1 !== e.active && this.anchors.length ? this.active.length && !t.contains(this.tablist[0], this.active[0]) ? this.tabs.length === e.disabled.length ? (e.active = !1, this.active = t()) : this._activate(this._findNextTab(Math.max(0, e.active - 1), !1)) : e.active = this.tabs.index(this.active) : (e.active = !1, this.active = t()), this._refresh()
    },
    _refresh: function () {
      this._setupDisabled(this.options.disabled), this._setupEvents(this.options.event), this._setupHeightStyle(this.options.heightStyle), this.tabs.not(this.active).attr({
        'aria-selected': 'false',
        tabIndex: -1
      }), this.panels.not(this._getPanelForTab(this.active)).hide().attr({
        'aria-expanded': 'false',
        'aria-hidden': 'true'
      }), this.active.length ? (this.active.addClass('ui-tabs-active ui-state-active').attr({
        'aria-selected': 'true',
        tabIndex: 0
      }), this._getPanelForTab(this.active).show().attr({
        'aria-expanded': 'true',
        'aria-hidden': 'false'
      })) : this.tabs.eq(0).attr('tabIndex', 0)
    },
    _processTabs: function () {
      var e = this
      this.tablist = this._getList().addClass('ui-tabs-nav ui-helper-reset ui-helper-clearfix ui-widget-header ui-corner-all').attr('role', 'tablist'), this.tabs = this.tablist.find('> li:has(a[href])').addClass('ui-state-default ui-corner-top').attr({
        role: 'tab',
        tabIndex: -1
      }), this.anchors = this.tabs.map(function () {
        return t('a', this)[0]
      }).addClass('ui-tabs-anchor').attr({
        role: 'presentation',
        tabIndex: -1
      }), this.panels = t(), this.anchors.each(function (i, n) {
        var o; var a; var r; var l = t(n).uniqueId().attr('id'); var h = t(n).closest('li'); var c = h.attr('aria-controls')
        s(n) ? (o = n.hash, a = e.element.find(e._sanitizeSelector(o))) : (r = e._tabId(h), o = '#' + r, (a = e.element.find(o)).length || (a = e._createPanel(r)).insertAfter(e.panels[i - 1] || e.tablist), a.attr('aria-live', 'polite')), a.length && (e.panels = e.panels.add(a)), c && h.data('ui-tabs-aria-controls', c), h.attr({
          'aria-controls': o.substring(1),
          'aria-labelledby': l
        }), a.attr('aria-labelledby', l)
      }), this.panels.addClass('ui-tabs-panel ui-widget-content ui-corner-bottom').attr('role', 'tabpanel')
    },
    _getList: function () {
      return this.element.find('ol,ul').eq(0)
    },
    _createPanel: function (e) {
      return t('<div>').attr('id', e).addClass('ui-tabs-panel ui-widget-content ui-corner-bottom').data('ui-tabs-destroy', !0)
    },
    _setupDisabled: function (e) {
      t.isArray(e) && (e.length ? e.length === this.anchors.length && (e = !0) : e = !1)
      for (var i, s = 0; i = this.tabs[s]; s++) !0 === e || t.inArray(s, e) !== -1 ? t(i).addClass('ui-state-disabled').attr('aria-disabled', 'true') : t(i).removeClass('ui-state-disabled').removeAttr('aria-disabled')
      this.options.disabled = e
    },
    _setupEvents: function (e) {
      var i = {
        click: function (t) {
          t.preventDefault()
        }
      }
      e && t.each(e.split(' '), function (t, e) {
        i[e] = '_eventHandler'
      }), this._off(this.anchors.add(this.tabs).add(this.panels)), this._on(this.anchors, i), this._on(this.tabs, { keydown: '_tabKeydown' }), this._on(this.panels, { keydown: '_panelKeydown' }), this._focusable(this.tabs), this._hoverable(this.tabs)
    },
    _setupHeightStyle: function (e) {
      var i; var s; var n = this.element.parent()
      e === 'fill' ? (t.support.minHeight || (s = n.css('overflow'), n.css('overflow', 'hidden')), i = n.height(), this.element.siblings(':visible').each(function () {
        var e = t(this); var s = e.css('position')
        s !== 'absolute' && s !== 'fixed' && (i -= e.outerHeight(!0))
      }), s && n.css('overflow', s), this.element.children().not(this.panels).each(function () {
        i -= t(this).outerHeight(!0)
      }), this.panels.each(function () {
        t(this).height(Math.max(0, i - t(this).innerHeight() + t(this).height()))
      }).css('overflow', 'auto')) : e === 'auto' && (i = 0, this.panels.each(function () {
        i = Math.max(i, t(this).height('').height())
      }).height(i))
    },
    _eventHandler: function (e) {
      var i = this.options; var s = this.active; var n = t(e.currentTarget).closest('li'); var o = n[0] === s[0]
      var a = o && i.collapsible; var r = a ? t() : this._getPanelForTab(n)
      var l = s.length ? this._getPanelForTab(s) : t()
      var h = { oldTab: s, oldPanel: l, newTab: a ? t() : n, newPanel: r }
      e.preventDefault(), n.hasClass('ui-state-disabled') || n.hasClass('ui-tabs-loading') || this.running || o && !i.collapsible || !1 === this._trigger('beforeActivate', e, h) || (i.active = !a && this.tabs.index(n), this.active = o ? t() : n, this.xhr && this.xhr.abort(), l.length || r.length || t.error('jQuery UI Tabs: Mismatching fragment identifier.'), r.length && this.load(this.tabs.index(n), e), this._toggle(e, h))
    },
    _toggle: function (e, i) {
      function s () {
        o.running = !1, o._trigger('activate', e, i)
      }

      function n () {
        i.newTab.closest('li').addClass('ui-tabs-active ui-state-active'), a.length && o.options.show ? o._show(a, o.options.show, s) : (a.show(), s())
      }

      var o = this; var a = i.newPanel; var r = i.oldPanel
      this.running = !0, r.length && this.options.hide ? this._hide(r, this.options.hide, function () {
        i.oldTab.closest('li').removeClass('ui-tabs-active ui-state-active'), n()
      }) : (i.oldTab.closest('li').removeClass('ui-tabs-active ui-state-active'), r.hide(), n()), r.attr({
        'aria-expanded': 'false',
        'aria-hidden': 'true'
      }), i.oldTab.attr('aria-selected', 'false'), a.length && r.length ? i.oldTab.attr('tabIndex', -1) : a.length && this.tabs.filter(function () {
        return t(this).attr('tabIndex') === 0
      }).attr('tabIndex', -1), a.attr({
        'aria-expanded': 'true',
        'aria-hidden': 'false'
      }), i.newTab.attr({ 'aria-selected': 'true', tabIndex: 0 })
    },
    _activate: function (e) {
      var i; var s = this._findActive(e)
      s[0] !== this.active[0] && (s.length || (s = this.active), i = s.find('.ui-tabs-anchor')[0], this._eventHandler({
        target: i,
        currentTarget: i,
        preventDefault: t.noop
      }))
    },
    _findActive: function (e) {
      return !1 === e ? t() : this.tabs.eq(e)
    },
    _getIndex: function (t) {
      return typeof t === 'string' && (t = this.anchors.index(this.anchors.filter("[href$='" + t + "']"))), t
    },
    _destroy: function () {
      this.xhr && this.xhr.abort(), this.element.removeClass('ui-tabs ui-widget ui-widget-content ui-corner-all ui-tabs-collapsible'), this.tablist.removeClass('ui-tabs-nav ui-helper-reset ui-helper-clearfix ui-widget-header ui-corner-all').removeAttr('role'), this.anchors.removeClass('ui-tabs-anchor').removeAttr('role').removeAttr('tabIndex').removeData('href.tabs').removeData('load.tabs').removeUniqueId(), this.tabs.add(this.panels).each(function () {
        t.data(this, 'ui-tabs-destroy') ? t(this).remove() : t(this).removeClass('ui-state-default ui-state-active ui-state-disabled ui-corner-top ui-corner-bottom ui-widget-content ui-tabs-active ui-tabs-panel').removeAttr('tabIndex').removeAttr('aria-live').removeAttr('aria-busy').removeAttr('aria-selected').removeAttr('aria-labelledby').removeAttr('aria-hidden').removeAttr('aria-expanded').removeAttr('role')
      }), this.tabs.each(function () {
        var e = t(this); var i = e.data('ui-tabs-aria-controls')
        i ? e.attr('aria-controls', i) : e.removeAttr('aria-controls')
      }), this.panels.show(), this.options.heightStyle !== 'content' && this.panels.css('height', '')
    },
    enable: function (i) {
      var s = this.options.disabled
      !1 !== s && (i === e ? s = !1 : (i = this._getIndex(i), s = t.isArray(s) ? t.map(s, function (t) {
        return t !== i ? t : null
      }) : t.map(this.tabs, function (t, e) {
        return e !== i ? e : null
      })), this._setupDisabled(s))
    },
    disable: function (i) {
      var s = this.options.disabled
      if (!0 !== s) {
        if (i === e) s = !0; else {
          if (i = this._getIndex(i), t.inArray(i, s) !== -1) return
          s = t.isArray(s) ? t.merge([i], s).sort() : [i]
        }
        this._setupDisabled(s)
      }
    },
    load: function (e, i) {
      e = this._getIndex(e)
      var n = this; var o = this.tabs.eq(e); var a = o.find('.ui-tabs-anchor'); var r = this._getPanelForTab(o)
      var l = { tab: o, panel: r }
      s(a[0]) || (this.xhr = t.ajax(this._ajaxSettings(a, i, l)), this.xhr && this.xhr.statusText !== 'canceled' && (o.addClass('ui-tabs-loading'), r.attr('aria-busy', 'true'), this.xhr.success(function (t) {
        setTimeout(function () {
          r.html(t), n._trigger('load', i, l)
        }, 1)
      }).complete(function (t, e) {
        setTimeout(function () {
          e === 'abort' && n.panels.stop(!1, !0), o.removeClass('ui-tabs-loading'), r.removeAttr('aria-busy'), t === n.xhr && delete n.xhr
        }, 1)
      })))
    },
    _ajaxSettings: function (e, i, s) {
      var n = this
      return {
        url: e.attr('href'),
        beforeSend: function (e, o) {
          return n._trigger('beforeLoad', i, t.extend({ jqXHR: e, ajaxSettings: o }, s))
        }
      }
    },
    _getPanelForTab: function (e) {
      var i = t(e).attr('aria-controls')
      return this.element.find(this._sanitizeSelector('#' + i))
    }
  }), !1 !== t.uiBackCompat && (t.ui.tabs.prototype._ui = function (t, e) {
    return { tab: t, panel: e, index: this.anchors.index(t) }
  }, t.widget('ui.tabs', t.ui.tabs, {
    url: function (t, e) {
      this.anchors.eq(t).attr('href', e)
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { ajaxOptions: null, cache: !1 },
    _create: function () {
      this._super()
      var e = this
      this._on({
        tabsbeforeload: function (i, s) {
          return t.data(s.tab[0], 'cache.tabs') ? void i.preventDefault() : void s.jqXHR.success(function () {
            e.options.cache && t.data(s.tab[0], 'cache.tabs', !0)
          })
        }
      })
    },
    _ajaxSettings: function (e, i, s) {
      var n = this.options.ajaxOptions
      return t.extend({}, n, {
        error: function (t, e) {
          try {
            n.error(t, e, s.tab.closest('li').index(), s.tab[0])
          } catch (t) {
          }
        }
      }, this._superApply(arguments))
    },
    _setOption: function (t, e) {
      t === 'cache' && !1 === e && this.anchors.removeData('cache.tabs'), this._super(t, e)
    },
    _destroy: function () {
      this.anchors.removeData('cache.tabs'), this._super()
    },
    url: function (t) {
      this.anchors.eq(t).removeData('cache.tabs'), this._superApply(arguments)
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    abort: function () {
      this.xhr && this.xhr.abort()
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { spinner: '<em>Loading&#8230;</em>' },
    _create: function () {
      this._super(), this._on({
        tabsbeforeload: function (t, e) {
          if (t.target === this.element[0] && this.options.spinner) {
            var i = e.tab.find('span'); var s = i.html()
            i.html(this.options.spinner), e.jqXHR.complete(function () {
              i.html(s)
            })
          }
        }
      })
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { enable: null, disable: null },
    enable: function (e) {
      var i; var s = this.options;
      (e && !0 === s.disabled || t.isArray(s.disabled) && t.inArray(e, s.disabled) !== -1) && (i = !0), this._superApply(arguments), i && this._trigger('enable', null, this._ui(this.anchors[e], this.panels[e]))
    },
    disable: function (e) {
      var i; var s = this.options;
      (e && !1 === s.disabled || t.isArray(s.disabled) && t.inArray(e, s.disabled) === -1) && (i = !0), this._superApply(arguments), i && this._trigger('disable', null, this._ui(this.anchors[e], this.panels[e]))
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: {
      add: null,
      remove: null,
      tabTemplate: "<li><a href='#{href}'><span>#{label}</span></a></li>"
    },
    add: function (i, s, n) {
      n === e && (n = this.anchors.length)
      var o; var a; var r = this.options; var l = t(r.tabTemplate.replace(/#\{href\}/g, i).replace(/#\{label\}/g, s))
      var h = i.indexOf('#') ? this._tabId(l) : i.replace('#', '')
      return l.addClass('ui-state-default ui-corner-top').data('ui-tabs-destroy', !0), l.attr('aria-controls', h), o = n >= this.tabs.length, (a = this.element.find('#' + h)).length || (a = this._createPanel(h), o ? n > 0 ? a.insertAfter(this.panels.eq(-1)) : a.appendTo(this.element) : a.insertBefore(this.panels[n])), a.addClass('ui-tabs-panel ui-widget-content ui-corner-bottom').hide(), o ? l.appendTo(this.tablist) : l.insertBefore(this.tabs[n]), r.disabled = t.map(r.disabled, function (t) {
        return t >= n ? ++t : t
      }), this.refresh(), this.tabs.length === 1 && !1 === r.active && this.option('active', 0), this._trigger('add', null, this._ui(this.anchors[n], this.panels[n])), this
    },
    remove: function (e) {
      e = this._getIndex(e)
      var i = this.options; var s = this.tabs.eq(e).remove(); var n = this._getPanelForTab(s).remove()
      return s.hasClass('ui-tabs-active') && this.anchors.length > 2 && this._activate(e + (e + 1 < this.anchors.length ? 1 : -1)), i.disabled = t.map(t.grep(i.disabled, function (t) {
        return t !== e
      }), function (t) {
        return t >= e ? --t : t
      }), this.refresh(), this._trigger('remove', null, this._ui(s.find('a')[0], n[0])), this
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    length: function () {
      return this.anchors.length
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { idPrefix: 'ui-tabs-' },
    _tabId: function (e) {
      var s = e.is('li') ? e.find('a[href]') : e
      return s = s[0], t(s).closest('li').attr('aria-controls') || s.title && s.title.replace(/\s/g, '_').replace(/[^\w\u00c0-\uFFFF\-]/g, '') || this.options.idPrefix + i()
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { panelTemplate: '<div></div>' },
    _createPanel: function (e) {
      return t(this.options.panelTemplate).attr('id', e).addClass('ui-tabs-panel ui-widget-content ui-corner-bottom').data('ui-tabs-destroy', !0)
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    _create: function () {
      var t = this.options
      t.active === null && t.selected !== e && (t.active = t.selected !== -1 && t.selected), this._super(), t.selected = t.active, !1 === t.selected && (t.selected = -1)
    },
    _setOption: function (t, e) {
      if (t !== 'selected') return this._super(t, e)
      var i = this.options
      this._super('active', e !== -1 && e), i.selected = i.active, !1 === i.selected && (i.selected = -1)
    },
    _eventHandler: function () {
      this._superApply(arguments), this.options.selected = this.options.active, !1 === this.options.selected && (this.options.selected = -1)
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { show: null, select: null },
    _create: function () {
      this._super(), !1 !== this.options.active && this._trigger('show', null, this._ui(this.active.find('.ui-tabs-anchor')[0], this._getPanelForTab(this.active)[0]))
    },
    _trigger: function (t, e, i) {
      var s; var n; var o = this._superApply(arguments)
      return !!o && (t === 'beforeActivate' ? (s = i.newTab.length ? i.newTab : i.oldTab, n = i.newPanel.length ? i.newPanel : i.oldPanel, o = this._super('select', e, {
        tab: s.find('.ui-tabs-anchor')[0],
        panel: n[0],
        index: s.closest('li').index()
      })) : t === 'activate' && i.newTab.length && (o = this._super('show', e, {
        tab: i.newTab.find('.ui-tabs-anchor')[0],
        panel: i.newPanel[0],
        index: i.newTab.closest('li').index()
      })), o)
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    select: function (t) {
      if ((t = this._getIndex(t)) === -1) {
        if (!this.options.collapsible || this.options.selected === -1) return
        t = this.options.selected
      }
      this.anchors.eq(t).trigger(this.options.event + this.eventNamespace)
    }
  }), (function () {
    var e = 0
    t.widget('ui.tabs', t.ui.tabs, {
      options: { cookie: null },
      _create: function () {
        var t; var e = this.options
        e.active == null && e.cookie && ((t = parseInt(this._cookie(), 10)) === -1 && (t = !1), e.active = t), this._super()
      },
      _cookie: function (i) {
        var s = [this.cookie || (this.cookie = this.options.cookie.name || 'ui-tabs-' + ++e)]
        return arguments.length && (s.push(!1 === i ? -1 : i), s.push(this.options.cookie)), t.cookie.apply(null, s)
      },
      _refresh: function () {
        this._super(), this.options.cookie && this._cookie(this.options.active, this.options.cookie)
      },
      _eventHandler: function () {
        this._superApply(arguments), this.options.cookie && this._cookie(this.options.active, this.options.cookie)
      },
      _destroy: function () {
        this._super(), this.options.cookie && this._cookie(null, this.options.cookie)
      }
    })
  }()), t.widget('ui.tabs', t.ui.tabs, {
    _trigger: function (e, i, s) {
      var n = t.extend({}, s)
      return e === 'load' && (n.panel = n.panel[0], n.tab = n.tab.find('.ui-tabs-anchor')[0]), this._super(e, i, n)
    }
  }), t.widget('ui.tabs', t.ui.tabs, {
    options: { fx: null },
    _getFx: function () {
      var e; var i; var s = this.options.fx
      return s && (t.isArray(s) ? (e = s[0], i = s[1]) : e = i = s), s ? { show: i, hide: e } : null
    },
    _toggle: function (t, e) {
      function i () {
        n.running = !1, n._trigger('activate', t, e)
      }

      function s () {
        e.newTab.closest('li').addClass('ui-tabs-active ui-state-active'), o.length && r.show ? o.animate(r.show, r.show.duration, function () {
          i()
        }) : (o.show(), i())
      }

      var n = this; var o = e.newPanel; var a = e.oldPanel; var r = this._getFx()
      return r ? (n.running = !0, void (a.length && r.hide ? a.animate(r.hide, r.hide.duration, function () {
        e.oldTab.closest('li').removeClass('ui-tabs-active ui-state-active'), s()
      }) : (e.oldTab.closest('li').removeClass('ui-tabs-active ui-state-active'), a.hide(), s()))) : this._super(t, e)
    }
  }))
}(jQuery)), (function (t) {
  function e (e, i) {
    var s = (e.attr('aria-describedby') || '').split(/\s+/)
    s.push(i), e.data('ui-tooltip-id', i).attr('aria-describedby', t.trim(s.join(' ')))
  }

  function i (e) {
    var i = e.data('ui-tooltip-id'); var s = (e.attr('aria-describedby') || '').split(/\s+/); var n = t.inArray(i, s)
    n !== -1 && s.splice(n, 1), e.removeData('ui-tooltip-id'), (s = t.trim(s.join(' '))) ? e.attr('aria-describedby', s) : e.removeAttr('aria-describedby')
  }

  var s = 0
}(jQuery)), (function ($) {
  var min = Math.min; var max = Math.max; var round = Math.floor; var isStr = function (t) {
    return $.type(t) === 'string'
  }; var runPluginCallbacks = function (Instance, a_fn) {
    function g (t) {
      return t
    }

    if ($.isArray(a_fn)) {
      for (var i = 0, c = a_fn.length; i < c; i++) {
        var fn = a_fn[i]
        try {
          isStr(fn) && (fn = eval(fn)), $.isFunction(fn) && g(fn)(Instance)
        } catch (t) {
        }
      }
    }
  }
  $.layout = {
    version: '1.4.4',
    revision: 1.0404,
    browser: {},
    effects: {
      slide: {
        all: { duration: 'fast' },
        north: { direction: 'up' },
        south: { direction: 'down' },
        east: { direction: 'right' },
        west: { direction: 'left' }
      },
      drop: {
        all: { duration: 'slow' },
        north: { direction: 'up' },
        south: { direction: 'down' },
        east: { direction: 'right' },
        west: { direction: 'left' }
      },
      scale: { all: { duration: 'fast' } },
      blind: {},
      clip: {},
      explode: {},
      fade: {},
      fold: {},
      puff: {},
      size: { all: { easing: 'swing' } }
    },
    config: {
      optionRootKeys: 'effects,panes,north,south,west,east,center'.split(','),
      allPanes: 'north,south,west,east,center'.split(','),
      borderPanes: 'north,south,west,east'.split(','),
      oppositeEdge: { north: 'south', south: 'north', east: 'west', west: 'east' },
      offscreenCSS: { left: '-99999px', right: 'auto' },
      offscreenReset: 'offscreenReset',
      hidden: { visibility: 'hidden' },
      visible: { visibility: 'visible' },
      resizers: {
        cssReq: {
          position: 'absolute',
          padding: 0,
          margin: 0,
          fontSize: '1px',
          textAlign: 'left',
          overflow: 'hidden'
        },
        cssDemo: { background: '#DDD', border: 'none' }
      },
      togglers: {
        cssReq: {
          position: 'absolute',
          display: 'block',
          padding: 0,
          margin: 0,
          overflow: 'hidden',
          textAlign: 'center',
          fontSize: '1px',
          cursor: 'pointer',
          zIndex: 1
        },
        cssDemo: { background: '#AAA' }
      },
      content: {
        cssReq: { position: 'relative' },
        cssDemo: { overflow: 'auto', padding: '10px' },
        cssDemoPane: { overflow: 'hidden', padding: 0 }
      },
      panes: {
        cssReq: { position: 'absolute', margin: 0 },
        cssDemo: { padding: '10px', background: '#FFF', border: '1px solid #BBB', overflow: 'auto' }
      },
      north: {
        side: 'top',
        sizeType: 'Height',
        dir: 'horz',
        cssReq: { top: 0, bottom: 'auto', left: 0, right: 0, width: 'auto' }
      },
      south: {
        side: 'bottom',
        sizeType: 'Height',
        dir: 'horz',
        cssReq: { top: 'auto', bottom: 0, left: 0, right: 0, width: 'auto' }
      },
      east: {
        side: 'right',
        sizeType: 'Width',
        dir: 'vert',
        cssReq: { left: 'auto', right: 0, top: 'auto', bottom: 'auto', height: 'auto' }
      },
      west: {
        side: 'left',
        sizeType: 'Width',
        dir: 'vert',
        cssReq: { left: 0, right: 'auto', top: 'auto', bottom: 'auto', height: 'auto' }
      },
      center: {
        dir: 'center',
        cssReq: { left: 'auto', right: 'auto', top: 'auto', bottom: 'auto', height: 'auto', width: 'auto' }
      }
    },
    callbacks: {},
    getParentPaneElem: function (t) {
      var e = $(t); var i = e.data('layout') || e.data('parentLayout')
      if (i) {
        var s = i.container
        if (s.data('layoutPane')) return s
        var n = s.closest('.' + $.layout.defaults.panes.paneClass)
        if (n.data('layoutPane')) return n
      }
      return null
    },
    getParentPaneInstance: function (t) {
      var e = $.layout.getParentPaneElem(t)
      return e ? e.data('layoutPane') : null
    },
    getParentLayoutInstance: function (t) {
      var e = $.layout.getParentPaneElem(t)
      return e ? e.data('parentLayout') : null
    },
    getEventObject: function (t) {
      return typeof t === 'object' && t.stopPropagation ? t : null
    },
    parsePaneName: function (t) {
      var e = $.layout.getEventObject(t); var i = t
      return e && (e.stopPropagation(), i = $(this).data('layoutEdge')), i && !/^(west|east|north|south|center)$/.test(i) && ($.layout.msg('LAYOUT ERROR - Invalid pane-name: "' + i + '"'), i = 'error'), i
    },
    plugins: {
      draggable: !!$.fn.draggable,
      effects: {
        core: !!$.effects,
        slide: $.effects && ($.effects.slide || $.effects.effect && $.effects.effect.slide)
      }
    },
    onCreate: [],
    onLoad: [],
    onReady: [],
    onDestroy: [],
    onUnload: [],
    afterOpen: [],
    afterClose: [],
    scrollbarWidth: function () {
      return window.scrollbarWidth || $.layout.getScrollbarSize('width')
    },
    scrollbarHeight: function () {
      return window.scrollbarHeight || $.layout.getScrollbarSize('height')
    },
    getScrollbarSize: function (t) {
      var e = $('<div style="position: absolute; top: -10000px; left: -10000px; width: 100px; height: 100px; border: 0; overflow: scroll;"></div>').appendTo('body')
      var i = { width: e.outerWidth - e[0].clientWidth, height: 100 - e[0].clientHeight }
      return e.remove(), window.scrollbarWidth = i.width, window.scrollbarHeight = i.height, t.match(/^(width|height)$/) ? i[t] : i
    },
    disableTextSelection: function () {
      var t = $(document); var e = 'textSelectionDisabled'; var i = 'textSelectionInitialized'
      $.fn.disableSelection && (t.data(i) || t.on('mouseup', $.layout.enableTextSelection).data(i, !0), t.data(e) || t.disableSelection().data(e, !0))
    },
    enableTextSelection: function () {
      var t = $(document); var e = 'textSelectionDisabled'
      $.fn.enableSelection && t.data(e) && t.enableSelection().data(e, !1)
    },
    showInvisibly: function (t, e) {
      if (t && t.length && (e || t.css('display') === 'none')) {
        var i = t[0].style; var s = { display: i.display || '', visibility: i.visibility || '' }
        return t.css({ display: 'block', visibility: 'hidden' }), s
      }
      return {}
    },
    getElementDimensions: function (t, e) {
      var i; var s; var n; var o = { css: {}, inset: {} }; var a = o.css; var r = { bottom: 0 }; var l = $.layout.cssNum; var h = Math.round
      var c = t.offset()
      return o.offsetLeft = c.left, o.offsetTop = c.top, e || (e = {}), $.each('Left,Right,Top,Bottom'.split(','), function (l, h) {
        i = a['border' + h] = $.layout.borderWidth(t, h), s = a['padding' + h] = $.layout.cssNum(t, 'padding' + h), n = h.toLowerCase(), o.inset[n] = e[n] >= 0 ? e[n] : s, r[n] = o.inset[n] + i
      }), a.width = h(t.width()), a.height = h(t.height()), a.top = l(t, 'top', !0), a.bottom = l(t, 'bottom', !0), a.left = l(t, 'left', !0), a.right = l(t, 'right', !0), o.outerWidth = h(t.outerWidth()), o.outerHeight = h(t.outerHeight()), o.innerWidth = max(0, o.outerWidth - r.left - r.right), o.innerHeight = max(0, o.outerHeight - r.top - r.bottom), o.layoutWidth = h(t.innerWidth()), o.layoutHeight = h(t.innerHeight()), o
    },
    getElementStyles: function (t, e) {
      var i; var s; var n; var o; var a; var r; var l = {}; var h = t[0].style; var c = e.split(','); var u = 'Top,Bottom,Left,Right'.split(',')
      var d = 'Color,Style,Width'.split(',')
      for (o = 0; o < c.length; o++) if ((i = c[o]).match(/(border|padding|margin)$/)) for (a = 0; a < 4; a++) if (s = u[a], i === 'border') for (r = 0; r < 3; r++) n = d[r], l[i + s + n] = h[i + s + n]; else l[i + s] = h[i + s]; else l[i] = h[i]
      return l
    },
    cssWidth: function (t, e) {
      if (e <= 0) return 0
      var i = $.layout.browser; var s = i.boxModel ? i.boxSizing ? t.css('boxSizing') : 'content-box' : 'border-box'
      var n = $.layout.borderWidth; var o = $.layout.cssNum; var a = e
      return s !== 'border-box' && (a -= n(t, 'Left') + n(t, 'Right')), s === 'content-box' && (a -= o(t, 'paddingLeft') + o(t, 'paddingRight')), max(0, a)
    },
    cssHeight: function (t, e) {
      if (e <= 0) return 0
      var i = $.layout.browser; var s = i.boxModel ? i.boxSizing ? t.css('boxSizing') : 'content-box' : 'border-box'
      var n = $.layout.borderWidth; var o = $.layout.cssNum; var a = e
      return s !== 'border-box' && (a -= n(t, 'Top') + n(t, 'Bottom')), s === 'content-box' && (a -= o(t, 'paddingTop') + o(t, 'paddingBottom')), max(0, a)
    },
    cssNum: function (t, e, i) {
      t.jquery || (t = $(t))
      var s = $.layout.showInvisibly(t); var n = $.css(t[0], e, !0)
      var o = i && n == 'auto' ? n : Math.round(parseFloat(n) || 0)
      return t.css(s), o
    },
    borderWidth: function (t, e) {
      t.jquery && (t = t[0])
      var i = 'border' + e.substr(0, 1).toUpperCase() + e.substr(1)
      return $.css(t, i + 'Style', !0) === 'none' ? 0 : Math.round(parseFloat($.css(t, i + 'Width', !0)) || 0)
    },
    isMouseOverElem: function (t, e) {
      var i = $(e || this); var s = i.offset(); var n = s.top; var o = s.left; var a = o + i.outerWidth()
      var r = n + i.outerHeight(); var l = t.pageX; var h = t.pageY
      return $.layout.browser.msie && l < 0 && h < 0 || l >= o && l <= a && h >= n && h <= r
    },
    msg: function (t, e, i, s) {
      if ($.isPlainObject(t) && window.debugData) {
        typeof e === 'string' ? (s = i, i = e) : typeof i === 'object' && (s = i, i = null)
        var n = i || 'log( <object> )'; var o = $.extend({ sort: !1, returnHTML: !1, display: !1 }, s)
        !0 === e || o.display ? debugData(t, n, o) : window.console && console.log(debugData(t, n, o))
      } else if (e) alert(t); else if (window.console) console.log(t); else {
        var a = $('#layoutLogger')
        a.length || (a = (function () {
          var t = $.support.fixedPosition ? 'fixed' : 'absolute'
          var e = $('<div id="layoutLogger" style="position: ' + t + '; top: 5px; z-index: 999999; max-width: 25%; overflow: hidden; border: 1px solid #000; border-radius: 5px; background: #FBFBFB; box-shadow: 0 2px 10px rgba(0,0,0,0.3);"><div style="font-size: 13px; font-weight: bold; padding: 5px 10px; background: #F6F6F6; border-radius: 5px 5px 0 0; cursor: move;"><span style="float: right; padding-left: 7px; cursor: pointer;" title="Remove Console" onclick="$(this).closest(\'#layoutLogger\').remove()">X</span>Layout console.log</div><ul style="font-size: 13px; font-weight: none; list-style: none; margin: 0; padding: 0 0 2px;"></ul></div>').appendTo('body')
          return e.css('left', $(window).width() - e.outerWidth() - 5), $.ui.draggable && e.draggable({ handle: ':first-child' }), e
        }())), a.children('ul').append('<li style="padding: 4px 10px; margin: 0; border-top: 1px solid #CCC;">' + t.replace(/\</g, '&lt;').replace(/\>/g, '&gt;') + '</li>')
      }
    }
  }, (function () {
    var t = navigator.userAgent.toLowerCase()
    var e = /(chrome)[ \/]([\w.]+)/.exec(t) || /(webkit)[ \/]([\w.]+)/.exec(t) || /(opera)(?:.*version|)[ \/]([\w.]+)/.exec(t) || /(msie) ([\w.]+)/.exec(t) || t.indexOf('compatible') < 0 && /(mozilla)(?:.*? rv:([\w.]+)|)/.exec(t) || []
    var i = e[1] || ''; var s = e[2] || 0; var n = i === 'msie'; var o = document.compatMode; var a = $.support
    var r = void 0 !== a.boxSizing ? a.boxSizing : a.boxSizingReliable
    var l = !n || !o || o === 'CSS1Compat' || a.boxModel || !1; var h = $.layout.browser = {
      version: s,
      safari: i === 'webkit',
      webkit: i === 'chrome',
      msie: n,
      isIE6: n && s == 6,
      boxModel: l,
      boxSizing: !!(typeof r === 'function' ? r() : r)
    }
    i && (h[i] = !0), l || o || $(function () {
      h.boxModel = a.boxModel
    })
  }()), $.layout.defaults = {
    name: '',
    containerClass: 'ui-layout-container',
    inset: null,
    scrollToBookmarkOnLoad: !0,
    resizeWithWindow: !0,
    resizeWithWindowDelay: 200,
    resizeWithWindowMaxDelay: 0,
    maskPanesEarly: !1,
    onresizeall_start: null,
    onresizeall_end: null,
    onload_start: null,
    onload_end: null,
    onunload_start: null,
    onunload_end: null,
    initPanes: !0,
    showErrorMessages: !0,
    showDebugMessages: !1,
    zIndex: null,
    zIndexes: {
      pane_normal: 0,
      content_mask: 1,
      resizer_normal: 2,
      pane_sliding: 100,
      pane_animate: 1e3,
      resizer_drag: 1e4
    },
    errors: {
      pane: 'pane',
      selector: 'selector',
      addButtonError: 'Error Adding Button\nInvalid ',
      containerMissing: 'UI Layout Initialization Error\nThe specified layout-container does not exist.',
      centerPaneMissing: 'UI Layout Initialization Error\nThe center-pane element does not exist.\nThe center-pane is a required element.',
      noContainerHeight: "UI Layout Initialization Warning\nThe layout-container \"CONTAINER\" has no height.\nTherefore the layout is 0-height and hence 'invisible'!",
      callbackError: 'UI Layout Callback Error\nThe EVENT callback is not a valid function.'
    },
    panes: {
      applyDemoStyles: !1,
      closable: !0,
      resizable: !0,
      slidable: !0,
      initClosed: !1,
      initHidden: !1,
      contentSelector: '.ui-layout-content',
      contentIgnoreSelector: '.ui-layout-ignore',
      findNestedContent: !1,
      paneClass: 'ui-layout-pane',
      resizerClass: 'ui-layout-resizer',
      togglerClass: 'ui-layout-toggler',
      buttonClass: 'ui-layout-button',
      minSize: 0,
      maxSize: 0,
      spacing_open: 6,
      spacing_closed: 6,
      togglerLength_open: 50,
      togglerLength_closed: 50,
      togglerAlign_open: 'center',
      togglerAlign_closed: 'center',
      togglerContent_open: '',
      togglerContent_closed: '',
      resizerDblClickToggle: !0,
      autoResize: !0,
      autoReopen: !0,
      resizerDragOpacity: 1,
      maskContents: !1,
      maskObjects: !1,
      maskZindex: null,
      resizingGrid: !1,
      livePaneResizing: !1,
      liveContentResizing: !1,
      liveResizingTolerance: 1,
      sliderCursor: 'pointer',
      slideTrigger_open: 'click',
      slideTrigger_close: 'mouseleave',
      slideDelay_open: 300,
      slideDelay_close: 300,
      hideTogglerOnSlide: !1,
      preventQuickSlideClose: $.layout.browser.webkit,
      preventPrematureSlideClose: !1,
      tips: {
        Open: 'Open',
        Close: 'Close',
        Resize: 'Resize',
        Slide: 'Slide Open',
        Pin: 'Pin',
        Unpin: 'Un-Pin',
        noRoomToOpen: 'Not enough room to show this panel.',
        minSizeWarning: 'Panel has reached its minimum size',
        maxSizeWarning: 'Panel has reached its maximum size'
      },
      showOverflowOnHover: !1,
      enableCursorHotkey: !0,
      customHotkeyModifier: 'SHIFT',
      fxName: 'slide',
      fxSpeed: null,
      fxSettings: {},
      fxOpacityFix: !0,
      animatePaneSizing: !1,
      children: null,
      containerSelector: '',
      initChildren: !0,
      destroyChildren: !0,
      resizeChildren: !0,
      triggerEventsOnLoad: !1,
      triggerEventsDuringLiveResize: !0,
      onshow_start: null,
      onshow_end: null,
      onhide_start: null,
      onhide_end: null,
      onopen_start: null,
      onopen_end: null,
      onclose_start: null,
      onclose_end: null,
      onresize_start: null,
      onresize_end: null,
      onsizecontent_start: null,
      onsizecontent_end: null,
      onswap_start: null,
      onswap_end: null,
      ondrag_start: null,
      ondrag_end: null
    },
    north: { paneSelector: '.ui-layout-north', size: 'auto', resizerCursor: 'n-resize', customHotkey: '' },
    south: { paneSelector: '.ui-layout-south', size: 'auto', resizerCursor: 's-resize', customHotkey: '' },
    east: { paneSelector: '.ui-layout-east', size: 200, resizerCursor: 'e-resize', customHotkey: '' },
    west: { paneSelector: '.ui-layout-west', size: 200, resizerCursor: 'w-resize', customHotkey: '' },
    center: { paneSelector: '.ui-layout-center', minWidth: 0, minHeight: 0 }
  }, $.layout.optionsMap = {
    layout: 'name,instanceKey,stateManagement,effects,inset,zIndexes,errors,zIndex,scrollToBookmarkOnLoad,showErrorMessages,maskPanesEarly,outset,resizeWithWindow,resizeWithWindowDelay,resizeWithWindowMaxDelay,onresizeall,onresizeall_start,onresizeall_end,onload,onload_start,onload_end,onunload,onunload_start,onunload_end'.split(','),
    center: 'paneClass,contentSelector,contentIgnoreSelector,findNestedContent,applyDemoStyles,triggerEventsOnLoad,showOverflowOnHover,maskContents,maskObjects,liveContentResizing,containerSelector,children,initChildren,resizeChildren,destroyChildren,onresize,onresize_start,onresize_end,onsizecontent,onsizecontent_start,onsizecontent_end'.split(','),
    noDefault: 'paneSelector,resizerCursor,customHotkey'.split(',')
  }, $.layout.transformData = function (t, e) {
    var i; var s; var n; var o; var a; var r; var l; var h = e ? { panes: {}, center: {} } : {}
    if (typeof t !== 'object') return h
    for (s in t) for (i = h, a = t[s], n = s.split('__'), l = n.length - 1, r = 0; r <= l; r++) o = n[r], r === l ? $.isPlainObject(a) ? i[o] = $.layout.transformData(a) : i[o] = a : (i[o] || (i[o] = {}), i = i[o])
    return h
  }, $.layout.backwardCompatibility = {
    map: {
      applyDefaultStyles: 'applyDemoStyles',
      childOptions: 'children',
      initChildLayout: 'initChildren',
      destroyChildLayout: 'destroyChildren',
      resizeChildLayout: 'resizeChildren',
      resizeNestedLayout: 'resizeChildren',
      resizeWhileDragging: 'livePaneResizing',
      resizeContentWhileDragging: 'liveContentResizing',
      triggerEventsWhileDragging: 'triggerEventsDuringLiveResize',
      maskIframesOnResize: 'maskContents',
      useStateCookie: 'stateManagement.enabled',
      'cookie.autoLoad': 'stateManagement.autoLoad',
      'cookie.autoSave': 'stateManagement.autoSave',
      'cookie.keys': 'stateManagement.stateKeys',
      'cookie.name': 'stateManagement.cookie.name',
      'cookie.domain': 'stateManagement.cookie.domain',
      'cookie.path': 'stateManagement.cookie.path',
      'cookie.expires': 'stateManagement.cookie.expires',
      'cookie.secure': 'stateManagement.cookie.secure',
      noRoomToOpenTip: 'tips.noRoomToOpen',
      togglerTip_open: 'tips.Close',
      togglerTip_closed: 'tips.Open',
      resizerTip: 'tips.Resize',
      sliderTip: 'tips.Slide'
    },
    renameOptions: function (t) {
      function e (e, i) {
        for (var s, n = e.split('.'), o = n.length - 1, a = {
            branch: t,
            key: n[o]
          }, r = 0; r < o; r++) s = n[r], void 0 == a.branch[s] ? a.branch = i ? a.branch[s] = {} : {} : a.branch = a.branch[s]
        return a
      }

      var i; var s; var n; var o = $.layout.backwardCompatibility.map
      for (var a in o) i = e(a), void 0 !== (n = i.branch[i.key]) && (s = e(o[a], !0), s.branch[s.key] = n, delete i.branch[i.key])
    },
    renameAllOptions: function (t) {
      var e = $.layout.backwardCompatibility.renameOptions
      return e(t), t.defaults && (typeof t.panes !== 'object' && (t.panes = {}), $.extend(!0, t.panes, t.defaults), delete t.defaults), t.panes && e(t.panes), $.each($.layout.config.allPanes, function (i, s) {
        t[s] && e(t[s])
      }), t
    }
  }, $.fn.layout = function (opts) {
    function keyDown (t) {
      if (!t) return !0
      var e = t.keyCode
      if (e < 33) return !0
      var i; var s; var n; var o; var a = { 38: 'north', 40: 'south', 37: 'west', 39: 'east' }; var r = (t.altKey, t.shiftKey)
      var l = t.ctrlKey
      return l && e >= 37 && e <= 40 && options[a[e]].enableCursorHotkey ? o = a[e] : (l || r) && $.each(_c.borderPanes, function (t, a) {
        if (i = options[a], s = i.customHotkey, n = i.customHotkeyModifier, (r && n == 'SHIFT' || l && n == 'CTRL' || l && r) && s && e === (isNaN(s) || s <= 9 ? s.toUpperCase().charCodeAt(0) : s)) return o = a, !1
      }), !(o && $Ps[o] && options[o].closable && !state[o].isHidden && (toggle(o), t.stopPropagation(), t.returnValue = !1, 1))
    }

    function allowOverflow (t) {
      if (isInitialized()) {
        this && this.tagName && (t = this)
        var e
        if (isStr(t) ? e = $Ps[t] : $(t).data('layoutRole') ? e = $(t) : $(t).parents().each(function () {
          if ($(this).data('layoutRole')) return e = $(this), !1
        }), e && e.length) {
          var i = e.data('layoutEdge'); var s = state[i]
          if (s.cssSaved && resetOverflow(i), s.isSliding || s.isResizing || s.isClosed) return void (s.cssSaved = !1)
          var n = { zIndex: options.zIndexes.resizer_normal + 1 }; var o = {}; var a = e.css('overflow')
          var r = e.css('overflowX'); var l = e.css('overflowY')
          a != 'visible' && (o.overflow = a, n.overflow = 'visible'), r && !r.match(/(visible|auto)/) && (o.overflowX = r, n.overflowX = 'visible'), l && !l.match(/(visible|auto)/) && (o.overflowY = r, n.overflowY = 'visible'), s.cssSaved = o, e.css(n), $.each(_c.allPanes, function (t, e) {
            e != i && resetOverflow(e)
          })
        }
      }
    }

    function resetOverflow (t) {
      if (isInitialized()) {
        this && this.tagName && (t = this)
        var e
        if (isStr(t) ? e = $Ps[t] : $(t).data('layoutRole') ? e = $(t) : $(t).parents().each(function () {
          if ($(this).data('layoutRole')) return e = $(this), !1
        }), e && e.length) {
          var i = e.data('layoutEdge'); var s = state[i]; var n = s.cssSaved || {}
          s.isSliding || s.isResizing || e.css('zIndex', options.zIndexes.pane_normal), e.css(n), s.cssSaved = !1
        }
      }
    }

    var browser = $.layout.browser; var _c = $.layout.config; var cssW = $.layout.cssWidth; var cssH = $.layout.cssHeight
    var elDims = $.layout.getElementDimensions; var styles = $.layout.getElementStyles
    var evtObj = $.layout.getEventObject; var evtPane = $.layout.parsePaneName
    var options = $.extend(!0, {}, $.layout.defaults)
    var effects = options.effects = $.extend(!0, {}, $.layout.effects); var state = {
      id: 'layout' + $.now(),
      initialized: !1,
      paneResizing: !1,
      panesSliding: {},
      container: { innerWidth: 0, innerHeight: 0, outerWidth: 0, outerHeight: 0, layoutWidth: 0, layoutHeight: 0 },
      north: { childIdx: 0 },
      south: { childIdx: 0 },
      east: { childIdx: 0 },
      west: { childIdx: 0 },
      center: { childIdx: 0 }
    }; var children = { north: null, south: null, east: null, west: null, center: null }; var timer = {
      data: {},
      set: function (t, e, i) {
        timer.clear(t), timer.data[t] = setTimeout(e, i)
      },
      clear: function (t) {
        var e = timer.data
        e[t] && (clearTimeout(e[t]), delete e[t])
      }
    }; var _log = function (t, e, i) {
      var s = options
      return (s.showErrorMessages && !i || i && s.showDebugMessages) && $.layout.msg(s.name + ' / ' + t, !1 !== e), !1
    }; var _runCallbacks = function (evtName, pane, skipBoundEvents) {
      function g (t) {
        return t
      }

      var hasPane = pane && isStr(pane); var s = hasPane ? state[pane] : state; var o = hasPane ? options[pane] : options
      var lName = options.name; var lng = evtName + (evtName.match(/_/) ? '' : '_end')
      var shrt = lng.match(/_end$/) ? lng.substr(0, lng.length - 4) : ''; var fn = o[lng] || o[shrt]; var retVal = 'NC'
      var args = []; var $P = hasPane ? $Ps[pane] : 0
      if (hasPane && !$P) return retVal
      if (hasPane || $.type(pane) !== 'boolean' || (skipBoundEvents = pane, pane = ''), fn) {
        try {
          isStr(fn) && (fn.match(/,/) ? (args = fn.split(','), fn = eval(args[0])) : fn = eval(fn)), $.isFunction(fn) && (retVal = args.length ? g(fn)(args[1]) : hasPane ? g(fn)(pane, $Ps[pane], s, o, lName) : g(fn)(Instance, s, o, lName))
        } catch (t) {
          _log(options.errors.callbackError.replace(/EVENT/, $.trim((pane || '') + ' ' + lng)), !1), $.type(t) === 'string' && string.length && _log('Exception:  ' + t, !1)
        }
      }
      return skipBoundEvents || !1 === retVal || (hasPane ? (o = options[pane], s = state[pane], $P.triggerHandler('layoutpane' + lng, [pane, $P, s, o, lName]), shrt && $P.triggerHandler('layoutpane' + shrt, [pane, $P, s, o, lName])) : ($N.triggerHandler('layout' + lng, [Instance, s, o, lName]), shrt && $N.triggerHandler('layout' + shrt, [Instance, s, o, lName]))), hasPane && evtName === 'onresize_end' && resizeChildren(pane + '', !0), retVal
    }; var _fixIframe = function (t) {
      if (!browser.mozilla) {
        var e = $Ps[t]
        state[t].tagName === 'IFRAME' ? e.css(_c.hidden).css(_c.visible) : e.find('IFRAME').css(_c.hidden).css(_c.visible)
      }
    }; var cssSize = function (t, e) {
      return (_c[t].dir == 'horz' ? cssH : cssW)($Ps[t], e)
    }; var cssMinDims = function (t) {
      var e = $Ps[t]; var i = _c[t].dir; var s = { minWidth: 1001 - cssW(e, 1e3), minHeight: 1001 - cssH(e, 1e3) }
      return i === 'horz' && (s.minSize = s.minHeight), i === 'vert' && (s.minSize = s.minWidth), s
    }; var setOuterWidth = function (t, e, i) {
      var s; var n = t
      isStr(t) ? n = $Ps[t] : t.jquery || (n = $(t)), s = cssW(n, e), n.css({ width: s }), s > 0 ? i && n.data('autoHidden') && n.innerHeight() > 0 && (n.show().data('autoHidden', !1), browser.mozilla || n.css(_c.hidden).css(_c.visible)) : i && !n.data('autoHidden') && n.hide().data('autoHidden', !0)
    }; var setOuterHeight = function (t, e, i) {
      var s; var n = t
      isStr(t) ? n = $Ps[t] : t.jquery || (n = $(t)), s = cssH(n, e), n.css({
        height: s,
        visibility: 'visible'
      }), s > 0 && n.innerWidth() > 0 ? i && n.data('autoHidden') && (n.show().data('autoHidden', !1), browser.mozilla || n.css(_c.hidden).css(_c.visible)) : i && !n.data('autoHidden') && n.hide().data('autoHidden', !0)
    }; var _parseSize = function (t, e, i) {
      if (i || (i = _c[t].dir), isStr(e) && e.match(/%/) && (e = e === '100%' ? -1 : parseInt(e, 10) / 100), e === 0) return 0
      if (e >= 1) return parseInt(e, 10)
      var s = options; var n = 0
      if (i == 'horz' ? n = sC.innerHeight - ($Ps.north ? s.north.spacing_open : 0) - ($Ps.south ? s.south.spacing_open : 0) : i == 'vert' && (n = sC.innerWidth - ($Ps.west ? s.west.spacing_open : 0) - ($Ps.east ? s.east.spacing_open : 0)), e === -1) return n
      if (e > 0) return round(n * e)
      if (t == 'center') return 0
      var o = i === 'horz' ? 'height' : 'width'; var a = $Ps[t]; var r = o === 'height' && $Cs[t]
      var l = $.layout.showInvisibly(a); var h = a.css(o); var c = r ? r.css(o) : 0
      return a.css(o, 'auto'), r && r.css(o, 'auto'), e = o === 'height' ? a.outerHeight() : a.outerWidth(), a.css(o, h).css(l), r && r.css(o, c), e
    }; var getPaneSize = function (t, e) {
      var i = $Ps[t]; var s = options[t]; var n = state[t]; var o = e ? s.spacing_open : 0; var a = e ? s.spacing_closed : 0
      return !i || n.isHidden ? 0 : n.isClosed || n.isSliding && e ? a : _c[t].dir === 'horz' ? i.outerHeight() + o : i.outerWidth() + o
    }; var setSizeLimits = function (t, e) {
      if (isInitialized()) {
        var i = options[t]; var s = state[t]; var n = _c[t]; var o = n.dir
        var a = (n.sizeType.toLowerCase(), void 0 != e ? e : s.isSliding); var r = ($Ps[t], i.spacing_open)
        var l = _c.oppositeEdge[t]; var h = state[l]; var c = $Ps[l]
        var u = !c || !1 === h.isVisible || h.isSliding ? 0 : o == 'horz' ? c.outerHeight() : c.outerWidth()
        var d = (!c || h.isHidden ? 0 : options[l][!1 !== h.isClosed ? 'spacing_closed' : 'spacing_open']) || 0
        var p = o == 'horz' ? sC.innerHeight : sC.innerWidth; var f = cssMinDims('center')
        var g = o == 'horz' ? max(options.center.minHeight, f.minHeight) : max(options.center.minWidth, f.minWidth)
        var m = p - r - (a ? 0 : _parseSize('center', g, o) + u + d)
        var v = s.minSize = max(_parseSize(t, i.minSize), cssMinDims(t).minSize)
        var b = s.maxSize = min(i.maxSize ? _parseSize(t, i.maxSize) : 1e5, m); var _ = s.resizerPosition = {}
        var y = sC.inset.top; var w = sC.inset.left; var C = sC.innerWidth; var x = sC.innerHeight; var z = i.spacing_open
        switch (t) {
          case 'north':
            _.min = y + v, _.max = y + b
            break
          case 'west':
            _.min = w + v, _.max = w + b
            break
          case 'south':
            _.min = y + x - b - z, _.max = y + x - v - z
            break
          case 'east':
            _.min = w + C - b - z, _.max = w + C - v - z
        }
      }
    }; var calcNewCenterPaneDims = function () {
      var t = {
        top: getPaneSize('north', !0),
        bottom: getPaneSize('south', !0),
        left: getPaneSize('west', !0),
        right: getPaneSize('east', !0),
        width: 0,
        height: 0
      }
      return t.width = sC.innerWidth - t.left - t.right, t.height = sC.innerHeight - t.bottom - t.top, t.top += sC.inset.top, t.bottom += sC.inset.bottom, t.left += sC.inset.left, t.right += sC.inset.right, t
    }; var getHoverClasses = function (t, e) {
      var i = $(t); var s = i.data('layoutRole'); var n = i.data('layoutEdge'); var o = options[n][s + 'Class']; var a = '-' + n
      var r = '-open'; var l = '-closed'; var h = '-sliding'; var c = '-hover '; var u = i.hasClass(o + l) ? l : r
      var d = u === l ? r : l; var p = o + c + (o + a + c) + (o + u + c) + (o + a + u + c)
      return e && (p += o + d + c + (o + a + d + c)), s == 'resizer' && i.hasClass(o + h) && (p += o + h + c + (o + a + h + c)), $.trim(p)
    }; var addHover = function (t, e) {
      var i = $(e || this)
      t && i.data('layoutRole') === 'toggler' && t.stopPropagation(), i.addClass(getHoverClasses(i))
    }; var removeHover = function (t, e) {
      var i = $(e || this)
      i.removeClass(getHoverClasses(i, !0))
    }; var onResizerEnter = function (t) {
      var e = $(this).data('layoutEdge'); var i = state[e]
      $(document), i.isResizing || state.paneResizing || options.maskPanesEarly && showMasks(e, { resizing: !0 })
    }; var onResizerLeave = function (t, e) {
      var i = e || this; var s = $(i).data('layoutEdge'); var n = s + 'ResizerLeave'
      $(document), timer.clear(s + '_openSlider'), timer.clear(n), e ? options.maskPanesEarly && !state.paneResizing && hideMasks() : timer.set(n, function () {
        onResizerLeave(t, i)
      }, 200)
    }; var _create = function () {
      initOptions()
      var t = options; var e = state
      return e.creatingLayout = !0, runPluginCallbacks(Instance, $.layout.onCreate), !1 === _runCallbacks('onload_start') ? 'cancel' : (_initContainer(), initHotkeys(), $(window).bind('unload.' + sID, unload), runPluginCallbacks(Instance, $.layout.onLoad), t.initPanes && _initLayoutElements(), delete e.creatingLayout, state.initialized)
    }; var isInitialized = function () {
      return !(!state.initialized && !state.creatingLayout) || _initLayoutElements()
    }; var _initLayoutElements = function (t) {
      var e = options
      if (!$N.is(':visible')) {
        return !t && browser.webkit && $N[0].tagName === 'BODY' && setTimeout(function () {
          _initLayoutElements(!0)
        }, 50), !1
      }
      if (!getPane('center').length) return _log(e.errors.centerPaneMissing)
      if (state.creatingLayout = !0, $.extend(sC, elDims($N, e.inset)), initPanes(), e.scrollToBookmarkOnLoad) {
        var i = self.location
        i.hash && i.replace(i.hash)
      }
      return Instance.hasParentLayout ? e.resizeWithWindow = !1 : e.resizeWithWindow && $(window).bind('resize.' + sID, windowResize), delete state.creatingLayout, state.initialized = !0, runPluginCallbacks(Instance, $.layout.onReady), _runCallbacks('onload_end'), !0
    }; var createChildren = function (t, e) {
      var i = evtPane.call(this, t); var s = $Ps[i]
      if (s) {
        var n = $Cs[i]; var o = state[i]; var a = options[i]; var r = options.stateManagement || {}
        var l = e ? a.children = e : a.children
        if ($.isPlainObject(l)) l = [l]; else if (!l || !$.isArray(l)) return
        $.each(l, function (t, e) {
          $.isPlainObject(e) && (e.containerSelector ? s.find(e.containerSelector) : n || s).each(function () {
            var t = $(this); var s = t.data('layout')
            if (!s) {
              if (setInstanceKey({
                container: t,
                options: e
              }, o), r.includeChildren && state.stateData[i]) {
                var n = (state.stateData[i].children || {})[e.instanceKey]
                var a = e.stateManagement || (e.stateManagement = { autoLoad: !0 })
                !0 === a.autoLoad && n && (a.autoSave = !1, a.includeChildren = !0, a.autoLoad = $.extend(!0, {}, n))
              }
              (s = t.layout(e)) && refreshChildren(i, s)
            }
          })
        })
      }
    }; var setInstanceKey = function (t, e) {
      var i = t.container; var s = t.options; var n = s.stateManagement; var o = s.instanceKey || i.data('layoutInstanceKey')
      return o || (o = (n && n.cookie ? n.cookie.name : '') || s.name), o = o ? o.replace(/[^\w-]/gi, '_').replace(/_{2,}/g, '_') : 'layout' + ++e.childIdx, s.instanceKey = o, i.data('layoutInstanceKey', o), o
    }; var refreshChildren = function (t, e) {
      var i; var s = $Ps[t]; var n = children[t]; var o = state[t]
      $.isPlainObject(n) && ($.each(n, function (t, e) {
        e.destroyed && delete n[t]
      }), $.isEmptyObject(n) && (n = children[t] = null)), e || n || (e = s.data('layout')), e && (e.hasParentLayout = !0, i = e.options, setInstanceKey(e, o), n || (n = children[t] = {}), n[i.instanceKey] = e.container.data('layout')), Instance[t].children = children[t], e || createChildren(t)
    }; var windowResize = function () {
      var t = options; var e = Number(t.resizeWithWindowDelay)
      e < 10 && (e = 100), timer.clear('winResize'), timer.set('winResize', function () {
        timer.clear('winResize'), timer.clear('winResizeRepeater')
        var e = elDims($N, t.inset)
        e.innerWidth === sC.innerWidth && e.innerHeight === sC.innerHeight || resizeAll()
      }, e), timer.data.winResizeRepeater || setWindowResizeRepeater()
    }; var setWindowResizeRepeater = function () {
      var t = Number(options.resizeWithWindowMaxDelay)
      t > 0 && timer.set('winResizeRepeater', function () {
        setWindowResizeRepeater(), resizeAll()
      }, t)
    }; var unload = function () {
      _runCallbacks('onunload_start'), runPluginCallbacks(Instance, $.layout.onUnload), _runCallbacks('onunload_end')
    }; var _initContainer = function () {
      var t; var e; var i = $N[0]; var s = $('html'); var n = sC.tagName = i.tagName; var o = sC.id = i.id
      var a = sC.className = i.className; var r = options; var l = r.name; var h = 'position,margin,padding,border'
      var c = 'layoutCSS'; var u = {}; var d = 'hidden'; var p = $N.data('parentLayout'); var f = $N.data('layoutEdge')
      var g = p && f; var m = $.layout.cssNum
      sC.ref = (r.name ? r.name + ' layout / ' : '') + n + (o ? '#' + o : a ? '.[' + a + ']' : ''), sC.isBody = n === 'BODY', g || sC.isBody || (t = $N.closest('.' + $.layout.defaults.panes.paneClass), p = t.data('parentLayout'), f = t.data('layoutEdge'), g = p && f), $N.data({
        layout: Instance,
        layoutContainer: sID
      }).addClass(r.containerClass)
      var v = { destroy: '', initPanes: '', resizeAll: 'resizeAll', resize: 'resizeAll' }
      for (l in v) $N.bind('layout' + l.toLowerCase() + '.' + sID, Instance[v[l] || l])
      g && (Instance.hasParentLayout = !0, p.refreshChildren(f, Instance)), $N.data(c) || (sC.isBody ? ($N.data(c, $.extend(styles($N, h), {
        height: $N.css('height'),
        overflow: $N.css('overflow'),
        overflowX: $N.css('overflowX'),
        overflowY: $N.css('overflowY')
      })), s.data(c, $.extend(styles(s, 'padding'), {
        height: 'auto',
        overflow: s.css('overflow'),
        overflowX: s.css('overflowX'),
        overflowY: s.css('overflowY')
      }))) : $N.data(c, styles($N, h + ',top,bottom,left,right,width,height,overflow,overflowX,overflowY')))
      try {
        if (u = {
          overflow: d,
          overflowX: d,
          overflowY: d
        }, $N.css(u), r.inset && !$.isPlainObject(r.inset) && (e = parseInt(r.inset, 10) || 0, r.inset = {
          top: e,
          bottom: e,
          left: e,
          right: e
        }), sC.isBody) {
          r.outset ? $.isPlainObject(r.outset) || (e = parseInt(r.outset, 10) || 0, r.outset = {
            top: e,
            bottom: e,
            left: e,
            right: e
          }) : r.outset = {
            top: m(s, 'paddingTop'),
            bottom: m(s, 'paddingBottom'),
            left: m(s, 'paddingLeft'),
            right: m(s, 'paddingRight')
          }, s.css(u).css({
            height: '100%',
            border: 'none',
            padding: 0,
            margin: 0
          }), browser.isIE6 ? ($N.css({
            width: '100%',
            height: '100%',
            border: 'none',
            padding: 0,
            margin: 0,
            position: 'relative'
          }), r.inset || (r.inset = elDims($N).inset)) : ($N.css({
            width: 'auto',
            height: 'auto',
            margin: 0,
            position: 'absolute'
          }), $N.css(r.outset)), $.extend(sC, elDims($N, r.inset))
        } else {
          var b = $N.css('position')
          b && b.match(/(fixed|absolute|relative)/) || $N.css('position', 'relative'), $N.is(':visible') && ($.extend(sC, elDims($N, r.inset)), sC.innerHeight < 1 && _log(r.errors.noContainerHeight.replace(/CONTAINER/, sC.ref)))
        }
        m($N, 'minWidth') && $N.parent().css('overflowX', 'auto'), m($N, 'minHeight') && $N.parent().css('overflowY', 'auto')
      } catch (t) {
      }
    }; var initHotkeys = function (t) {
      t = t ? t.split(',') : _c.borderPanes, $.each(t, function (t, e) {
        var i = options[e]
        if (i.enableCursorHotkey || i.customHotkey) return $(document).bind('keydown.' + sID, keyDown), !1
      })
    }; var initOptions = function () {
      function t (t) {
        var e = options[t]; var i = options.panes
        e.fxSettings || (e.fxSettings = {}), i.fxSettings || (i.fxSettings = {}), $.each(['_open', '_close', '_size'], function (s, n) {
          var o = 'fxName' + n; var a = 'fxSpeed' + n; var r = 'fxSettings' + n
          var l = e[o] = e[o] || i[o] || e.fxName || i.fxName || 'none'
          var h = $.effects && ($.effects[l] || $.effects.effect && $.effects.effect[l])
          l !== 'none' && options.effects[l] && h || (l = e[o] = 'none')
          var c = options.effects[l] || {}; var u = c.all || null; var d = c[t] || null
          e[a] = e[a] || i[a] || e.fxSpeed || i.fxSpeed || null, e[r] = $.extend(!0, {}, u, d, i.fxSettings, e.fxSettings, i[r], e[r])
        }), delete e.fxName, delete e.fxSpeed, delete e.fxSettings
      }

      var e, i, s, n, o, a, r
      if (opts = $.layout.transformData(opts, !0), opts = $.layout.backwardCompatibility.renameAllOptions(opts), !$.isEmptyObject(opts.panes)) {
        for (e = $.layout.optionsMap.noDefault, o = 0, a = e.length; o < a; o++) s = e[o], delete opts.panes[s]
        for (e = $.layout.optionsMap.layout, o = 0, a = e.length; o < a; o++) s = e[o], delete opts.panes[s]
      }
      e = $.layout.optionsMap.layout
      var l = $.layout.config.optionRootKeys
      for (s in opts) n = opts[s], $.inArray(s, l) < 0 && $.inArray(s, e) < 0 && (opts.panes[s] || (opts.panes[s] = $.isPlainObject(n) ? $.extend(!0, {}, n) : n), delete opts[s])
      $.extend(!0, options, opts), $.each(_c.allPanes, function (n, o) {
        if (_c[o] = $.extend(!0, {}, _c.panes, _c[o]), i = options.panes, r = options[o], o === 'center') for (e = $.layout.optionsMap.center, n = 0, a = e.length; n < a; n++) s = e[n], opts.center[s] || !opts.panes[s] && r[s] || (r[s] = i[s]); else r = options[o] = $.extend(!0, {}, i, r), t(o), r.resizerClass || (r.resizerClass = 'ui-layout-resizer'), r.togglerClass || (r.togglerClass = 'ui-layout-toggler')
        r.paneClass || (r.paneClass = 'ui-layout-pane')
      })
      var h = opts.zIndex; var c = options.zIndexes
      h > 0 && (c.pane_normal = h, c.content_mask = max(h + 1, c.content_mask), c.resizer_normal = max(h + 2, c.resizer_normal)), delete options.panes
    }; var getPane = function (t) {
      var e = options[t].paneSelector
      if (e.substr(0, 1) === '#') return $N.find(e).eq(0)
      var i = $N.children(e).eq(0)
      return i.length ? i : $N.children('form:first').children(e).eq(0)
    }; var initPanes = function (t) {
      evtPane(t), $.each(_c.allPanes, function (t, e) {
        addPane(e, !0)
      }), initHandles(), $.each(_c.borderPanes, function (t, e) {
        $Ps[e] && state[e].isVisible && (setSizeLimits(e), makePaneFit(e))
      }), sizeMidPanes('center'), $.each(_c.allPanes, function (t, e) {
        afterInitPane(e)
      })
    }; var addPane = function (t, e) {
      if (e || isInitialized()) {
        var i; var s; var n; var o = options[t]; var a = state[t]; var r = _c[t]; var l = r.dir
        var h = (a.fx, o.spacing_open, t === 'center'); var c = {}; var u = $Ps[t]
        if (u ? removePane(t, !1, !0, !1) : $Cs[t] = !1, !(u = $Ps[t] = getPane(t)).length) return void ($Ps[t] = !1)
        if (!u.data('layoutCSS')) {
          u.data('layoutCSS', styles(u, 'position,top,left,bottom,right,width,height,overflow,zIndex,display,backgroundColor,padding,margin,border'))
        }
        Instance[t] = {
          name: t,
          pane: $Ps[t],
          content: $Cs[t],
          options: options[t],
          state: state[t],
          children: children[t]
        }, u.data({
          parentLayout: Instance,
          layoutPane: Instance[t],
          layoutEdge: t,
          layoutRole: 'pane'
        }).css(r.cssReq).css('zIndex', options.zIndexes.pane_normal).css(o.applyDemoStyles ? r.cssDemo : {}).addClass(o.paneClass + ' ' + o.paneClass + '-' + t).bind('mouseenter.' + sID, addHover).bind('mouseleave.' + sID, removeHover)
        var d; var p = {
          hide: '',
          show: '',
          toggle: '',
          close: '',
          open: '',
          slideOpen: '',
          slideClose: '',
          slideToggle: '',
          size: 'sizePane',
          sizePane: 'sizePane',
          sizeContent: '',
          sizeHandles: '',
          enableClosable: '',
          disableClosable: '',
          enableSlideable: '',
          disableSlideable: '',
          enableResizable: '',
          disableResizable: '',
          swapPanes: 'swapPanes',
          swap: 'swapPanes',
          move: 'swapPanes',
          removePane: 'removePane',
          remove: 'removePane',
          createChildren: '',
          resizeChildren: '',
          resizeAll: 'resizeAll',
          resizeLayout: 'resizeAll'
        }
        for (d in p) u.bind('layoutpane' + d.toLowerCase() + '.' + sID, Instance[p[d] || d])
        initContent(t, !1), h || (i = a.size = _parseSize(t, o.size), s = _parseSize(t, o.minSize) || 1, n = _parseSize(t, o.maxSize) || 1e5, i > 0 && (i = max(min(i, n), s)), a.autoResize = o.autoResize, a.isClosed = !1, a.isSliding = !1, a.isResizing = !1, a.isHidden = !1, a.pins || (a.pins = [])), a.tagName = u[0].tagName, a.edge = t, a.noRoom = !1, a.isVisible = !0, setPanePosition(t), l === 'horz' ? c.height = cssH(u, i) : l === 'vert' && (c.width = cssW(u, i)), u.css(c), l != 'horz' && sizeMidPanes(t, !0), state.initialized && (initHandles(t), initHotkeys(t)), o.initClosed && o.closable && !o.initHidden ? close(t, !0, !0) : o.initHidden || o.initClosed ? hide(t) : a.noRoom || u.css('display', 'block'), u.css('visibility', 'visible'), o.showOverflowOnHover && u.hover(allowOverflow, resetOverflow), state.initialized && afterInitPane(t)
      }
    }; var afterInitPane = function (t) {
      var e = $Ps[t]; var i = state[t]; var s = options[t]
      e && (e.data('layout') && refreshChildren(t, e.data('layout')), i.isVisible && (state.initialized ? resizeAll() : sizeContent(t), s.triggerEventsOnLoad ? _runCallbacks('onresize_end', t) : resizeChildren(t, !0)), s.initChildren && s.children && createChildren(t))
    }; var setPanePosition = function (t) {
      t = t ? t.split(',') : _c.borderPanes, $.each(t, function (t, e) {
        var i = $Ps[e]; var s = $Rs[e]; var n = (options[e], state[e]); var o = _c[e].side; var a = {}
        if (i) {
          switch (e) {
            case 'north':
              a.top = sC.inset.top, a.left = sC.inset.left, a.right = sC.inset.right
              break
            case 'south':
              a.bottom = sC.inset.bottom, a.left = sC.inset.left, a.right = sC.inset.right
              break
            case 'west':
              a.left = sC.inset.left
              break
            case 'east':
              a.right = sC.inset.right
          }
          i.css(a), s && n.isClosed ? s.css(o, sC.inset[o]) : s && !n.isHidden && s.css(o, sC.inset[o] + getPaneSize(e))
        }
      })
    }; var initHandles = function (t) {
      t = t ? t.split(',') : _c.borderPanes, $.each(t, function (t, e) {
        var i = $Ps[e]
        if ($Rs[e] = !1, $Ts[e] = !1, i) {
          var s = options[e]; var n = state[e]
          var o = (_c[e], s.paneSelector.substr(0, 1) === '#' ? s.paneSelector.substr(1) : '')
          var a = s.resizerClass; var r = s.togglerClass
          var l = (n.isVisible ? s.spacing_open : s.spacing_closed, '-' + e); var h = (n.isVisible, Instance[e])
          var c = h.resizer = $Rs[e] = $('<div></div>')
          var u = h.toggler = !!s.closable && ($Ts[e] = $('<div></div>'))
          !n.isVisible && s.slidable && c.attr('title', s.tips.Slide).css('cursor', s.sliderCursor), c.attr('id', o ? o + '-resizer' : '').data({
            parentLayout: Instance,
            layoutPane: Instance[e],
            layoutEdge: e,
            layoutRole: 'resizer'
          }).css(_c.resizers.cssReq).css('zIndex', options.zIndexes.resizer_normal).css(s.applyDemoStyles ? _c.resizers.cssDemo : {}).addClass(a + ' ' + a + l).hover(addHover, removeHover).hover(onResizerEnter, onResizerLeave).mousedown($.layout.disableTextSelection).mouseup($.layout.enableTextSelection).appendTo($N), $.fn.disableSelection && c.disableSelection(), s.resizerDblClickToggle && c.bind('dblclick.' + sID, toggle), u && (u.attr('id', o ? o + '-toggler' : '').data({
            parentLayout: Instance,
            layoutPane: Instance[e],
            layoutEdge: e,
            layoutRole: 'toggler'
          }).css(_c.togglers.cssReq).css(s.applyDemoStyles ? _c.togglers.cssDemo : {}).addClass(r + ' ' + r + l).hover(addHover, removeHover).bind('mouseenter', onResizerEnter).appendTo(c), s.togglerContent_open && $('<span>' + s.togglerContent_open + '</span>').data({
            layoutEdge: e,
            layoutRole: 'togglerContent'
          }).data('layoutRole', 'togglerContent').data('layoutEdge', e).addClass('content content-open').css('display', 'none').appendTo(u), s.togglerContent_closed && $('<span>' + s.togglerContent_closed + '</span>').data({
            layoutEdge: e,
            layoutRole: 'togglerContent'
          }).addClass('content content-closed').css('display', 'none').appendTo(u), enableClosable(e)), initResizable(e), n.isVisible ? setAsOpen(e) : (setAsClosed(e), bindStartSlidingEvents(e, !0))
        }
      }), sizeHandles()
    }; var initContent = function (t, e) {
      if (isInitialized()) {
        var i; var s = options[t]; var n = s.contentSelector; var o = Instance[t]; var a = $Ps[t]
        n && (i = o.content = $Cs[t] = s.findNestedContent ? a.find(n).eq(0) : a.children(n).eq(0)), i && i.length ? (i.data('layoutRole', 'content'), i.data('layoutCSS') || i.data('layoutCSS', styles(i, 'height')), i.css(_c.content.cssReq), s.applyDemoStyles && (i.css(_c.content.cssDemo), a.css(_c.content.cssDemoPane)), a.css('overflowX').match(/(scroll|auto)/) && a.css('overflow', 'hidden'), state[t].content = {}, !1 !== e && sizeContent(t)) : o.content = $Cs[t] = !1
      }
    }; var initResizable = function (t) {
      var e = $.layout.plugins.draggable
      t = t ? t.split(',') : _c.borderPanes, $.each(t, function (t, s) {
        var n = options[s]
        if (!e || !$Ps[s] || !n.resizable) return n.resizable = !1, !0
        var o; var a; var r = state[s]; var l = options.zIndexes; var h = _c[s]; var c = h.dir == 'horz' ? 'top' : 'left'
        var u = ($Ps[s], $Rs[s]); var d = n.resizerClass; var p = 0; var f = d + '-drag'; var g = d + '-' + s + '-drag'
        var m = d + '-dragging'; var v = d + '-' + s + '-dragging'; var b = d + '-dragging-limit'
        var _ = d + '-' + s + '-dragging-limit'; var y = !1
        r.isClosed || u.attr('title', n.tips.Resize).css('cursor', n.resizerCursor), u.draggable({
          containment: $N[0],
          axis: h.dir == 'horz' ? 'y' : 'x',
          delay: 0,
          distance: 1,
          grid: n.resizingGrid,
          helper: 'clone',
          opacity: n.resizerDragOpacity,
          addClasses: !1,
          zIndex: l.resizer_drag,
          start: function (t, e) {
            return n = options[s], r = state[s], a = n.livePaneResizing, !1 !== _runCallbacks('ondrag_start', s) && (r.isResizing = !0, state.paneResizing = s, timer.clear(s + '_closeSlider'), setSizeLimits(s), o = r.resizerPosition, p = e.position[c], u.addClass(f + ' ' + g), y = !1, void showMasks(s, { resizing: !0 }))
          },
          drag: function (t, e) {
            y || (e.helper.addClass(m + ' ' + v).css({
              right: 'auto',
              bottom: 'auto'
            }).children().css('visibility', 'hidden'), y = !0, r.isSliding && $Ps[s].css('zIndex', l.pane_sliding))
            var h = 0
            e.position[c] < o.min ? (e.position[c] = o.min, h = -1) : e.position[c] > o.max && (e.position[c] = o.max, h = 1), h ? (e.helper.addClass(b + ' ' + _), window.defaultStatus = h > 0 && s.match(/(north|west)/) || h < 0 && s.match(/(south|east)/) ? n.tips.maxSizeWarning : n.tips.minSizeWarning) : (e.helper.removeClass(b + ' ' + _), window.defaultStatus = ''), a && Math.abs(e.position[c] - p) >= n.liveResizingTolerance && (p = e.position[c], i(t, e, s))
          },
          stop: function (t, e) {
            $('body').enableSelection(), window.defaultStatus = '', u.removeClass(f + ' ' + g), r.isResizing = !1, state.paneResizing = !1, i(t, e, s, !0)
          }
        })
      })
      var i = function (t, e, i, s) {
        var n; var o = e.position; var a = _c[i]; var r = options[i]; var l = state[i]
        switch (i) {
          case 'north':
            n = o.top
            break
          case 'west':
            n = o.left
            break
          case 'south':
            n = sC.layoutHeight - o.top - r.spacing_open
            break
          case 'east':
            n = sC.layoutWidth - o.left - r.spacing_open
        }
        var h = n - sC.inset[a.side]
        if (s) !1 !== _runCallbacks('ondrag_end', i) && manualSizePane(i, h, !1, !0), hideMasks(!0), l.isSliding && showMasks(i, { resizing: !0 }); else {
          if (Math.abs(h - l.size) < r.liveResizingTolerance) return
          manualSizePane(i, h, !1, !0), sizeMasks()
        }
      }
    }; var sizeMask = function () {
      var t = $(this); var e = t.data('layoutMask'); var i = state[e]
      i.tagName == 'IFRAME' && i.isVisible && t.css({
        top: i.offsetTop,
        left: i.offsetLeft,
        width: i.outerWidth,
        height: i.outerHeight
      })
    }; var sizeMasks = function () {
      $Ms.each(sizeMask)
    }; var showMasks = function (t, e) {
      var i; var s; var n = _c[t]; var o = ['center']; var a = options.zIndexes
      var r = $.extend({ objectsOnly: !1, animation: !1, resizing: !0, sliding: state[t].isSliding }, e)
      r.resizing && o.push(t), r.sliding && o.push(_c.oppositeEdge[t]), n.dir === 'horz' && (o.push('west'), o.push('east')), $.each(o, function (t, e) {
        s = state[e], i = options[e], s.isVisible && (i.maskObjects || !r.objectsOnly && i.maskContents) && getMasks(e).each(function () {
          sizeMask.call(this), this.style.zIndex = s.isSliding ? a.pane_sliding + 1 : a.pane_normal + 1, this.style.display = 'block'
        })
      })
    }; var hideMasks = function (t) {
      if (t || !state.paneResizing) $Ms.hide(); else if (!t && !$.isEmptyObject(state.panesSliding)) for (var e, i, s = $Ms.length - 1; s >= 0; s--) i = $Ms.eq(s), e = i.data('layoutMask'), options[e].maskObjects || i.hide()
    }; var getMasks = function (t) {
      for (var e, i = $([]), s = 0, n = $Ms.length; s < n; s++) (e = $Ms.eq(s)).data('layoutMask') === t && (i = i.add(e))
      return i.length ? i : createMasks(t)
    }; var createMasks = function (t) {
      var e; var i; var s; var n; var o; var a = $Ps[t]; var r = state[t]; var l = options[t]; var h = options.zIndexes
      if (!l.maskContents && !l.maskObjects) return $([])
      for (o = 0; o < (l.maskObjects ? 2 : 1); o++) e = l.maskObjects && o == 0, i = document.createElement(e ? 'iframe' : 'div'), s = $(i).data('layoutMask', t), i.className = 'ui-layout-mask ui-layout-mask-' + t, n = i.style, n.background = '#FFF', n.position = 'absolute', n.display = 'block', e ? (i.src = 'about:blank', i.frameborder = 0, n.border = 0, n.opacity = 0, n.filter = "Alpha(Opacity='0')") : (n.opacity = 0.001, n.filter = "Alpha(Opacity='1')"), r.tagName == 'IFRAME' ? (n.zIndex = h.pane_normal + 1, $N.append(i)) : (s.addClass('ui-layout-mask-inside-pane'), n.zIndex = l.maskZindex || h.content_mask, n.top = 0, n.left = 0, n.width = '100%', n.height = '100%', a.append(i)), $Ms = $Ms.add(i)
      return $Ms
    }; var destroy = function (t, e) {
      $(window).unbind('.' + sID), $(document).unbind('.' + sID), typeof t === 'object' ? evtPane(t) : e = t, $N.clearQueue().removeData('layout').removeData('layoutContainer').removeClass(options.containerClass).unbind('.' + sID), $Ms.remove(), $.each(_c.allPanes, function (t, i) {
        removePane(i, !1, !0, e)
      })
      var i = 'layoutCSS'
      $N.data(i) && !$N.data('layoutRole') && $N.css($N.data(i)).removeData(i), sC.tagName === 'BODY' && ($N = $('html')).data(i) && $N.css($N.data(i)).removeData(i), runPluginCallbacks(Instance, $.layout.onDestroy), unload()
      for (var s in Instance) s.match(/^(container|options)$/) || delete Instance[s]
      return Instance.destroyed = !0, Instance
    }; var removePane = function (t, e, i, s) {
      if (isInitialized()) {
        var n = evtPane.call(this, t); var o = $Ps[n]; var a = $Cs[n]; var r = $Rs[n]; var l = $Ts[n]
        o && $.isEmptyObject(o.data()) && (o = !1), a && $.isEmptyObject(a.data()) && (a = !1), r && $.isEmptyObject(r.data()) && (r = !1), l && $.isEmptyObject(l.data()) && (l = !1), o && o.stop(!0, !0)
        var h = options[n]; var c = (state[n], 'layoutCSS'); var u = children[n]
        var d = $.isPlainObject(u) && !$.isEmptyObject(u); var p = void 0 !== s ? s : h.destroyChildren
        if (d && p && ($.each(u, function (t, e) {
          e.destroyed || e.destroy(!0), e.destroyed && delete u[t]
        }), $.isEmptyObject(u) && (u = children[n] = null, d = !1)), o && e && !d) o.remove(); else if (o && o[0]) {
          var f = h.paneClass; var g = f + '-' + n; var m = '-open'; var v = '-sliding'; var b = '-closed'
          var _ = [f, f + m, f + b, f + v, g, g + m, g + b, g + v]
          $.merge(_, getHoverClasses(o, !0)), o.removeClass(_.join(' ')).removeData('parentLayout').removeData('layoutPane').removeData('layoutRole').removeData('layoutEdge').removeData('autoHidden').unbind('.' + sID), d && a ? (a.width(a.width()), $.each(u, function (t, e) {
            e.resizeAll()
          })) : a && a.css(a.data(c)).removeData(c).removeData('layoutRole'), o.data('layout') || o.css(o.data(c)).removeData(c)
        }
        l && l.remove(), r && r.remove(), Instance[n] = $Ps[n] = $Cs[n] = $Rs[n] = $Ts[n] = !1, { removed: !0 }, i || resizeAll()
      }
    }; var _hidePane = function (t) {
      var e = $Ps[t]; var i = options[t]; var s = e[0].style
      i.useOffscreenClose ? (e.data(_c.offscreenReset) || e.data(_c.offscreenReset, {
        left: s.left,
        right: s.right
      }), e.css(_c.offscreenCSS)) : e.hide().removeData(_c.offscreenReset)
    }; var _showPane = function (t) {
      var e = $Ps[t]; var i = options[t]; var s = _c.offscreenCSS; var n = e.data(_c.offscreenReset); var o = e[0].style
      e.show().removeData(_c.offscreenReset), i.useOffscreenClose && n && (o.left == s.left && (o.left = n.left), o.right == s.right && (o.right = n.right))
    }; var hide = function (t, e) {
      if (isInitialized()) {
        var i = evtPane.call(this, t); var s = options[i]; var n = state[i]; var o = $Ps[i]; var a = $Rs[i]
        i !== 'center' && o && !n.isHidden && (state.initialized && !1 === _runCallbacks('onhide_start', i) || (n.isSliding = !1, delete state.panesSliding[i], a && a.hide(), !state.initialized || n.isClosed ? (n.isClosed = !0, n.isHidden = !0, n.isVisible = !1, state.initialized || _hidePane(i), sizeMidPanes(_c[i].dir === 'horz' ? '' : 'center'), (state.initialized || s.triggerEventsOnLoad) && _runCallbacks('onhide_end', i)) : (n.isHiding = !0, close(i, !1, e))))
      }
    }; var show = function (t, e, i, s) {
      if (isInitialized()) {
        var n = evtPane.call(this, t); var o = (options[n], state[n]); var a = $Ps[n]
        $Rs[n], n !== 'center' && a && o.isHidden && !1 !== _runCallbacks('onshow_start', n) && (o.isShowing = !0, o.isSliding = !1, delete state.panesSliding[n], !1 === e ? close(n, !0) : open(n, !1, i, s))
      }
    }; var toggle = function (t, e) {
      if (isInitialized()) {
        var i = evtObj(t); var s = evtPane.call(this, t); var n = state[s]
        i && i.stopImmediatePropagation(), n.isHidden ? show(s) : n.isClosed ? open(s, !!e) : close(s)
      }
    }; var _closePane = function (t, e) {
      var i = ($Ps[t], state[t])
      _hidePane(t), i.isClosed = !0, i.isVisible = !1, e && setAsClosed(t)
    }; var close = function (t, e, i, s) {
      function n () {
        d.isMoving = !1, bindStartSlidingEvents(o, !0)
        var t = _c.oppositeEdge[o]
        state[t].noRoom && (setSizeLimits(t), makePaneFit(t)), s || !state.initialized && !u.triggerEventsOnLoad || (r || _runCallbacks('onclose_end', o), r && _runCallbacks('onshow_end', o), l && _runCallbacks('onhide_end', o))
      }

      var o = evtPane.call(this, t)
      if (o !== 'center') {
        if (!state.initialized && $Ps[o]) return void _closePane(o, !0)
        if (isInitialized()) {
          var a; var r; var l; var h; var c = $Ps[o]; var u = ($Rs[o], $Ts[o], options[o]); var d = state[o]
          _c[o], $N.queue(function (t) {
            if (!c || !u.closable && !d.isShowing && !d.isHiding || !e && d.isClosed && !d.isShowing) return t()
            var s = !d.isShowing && !1 === _runCallbacks('onclose_start', o)
            return r = d.isShowing, l = d.isHiding, h = d.isSliding, delete d.isShowing, delete d.isHiding, s ? t() : (a = !i && !d.isClosed && u.fxName_close != 'none', d.isMoving = !0, d.isClosed = !0, d.isVisible = !1, l ? d.isHidden = !0 : r && (d.isHidden = !1), d.isSliding ? bindStopSlidingEvents(o, !1) : sizeMidPanes(_c[o].dir === 'horz' ? '' : 'center', !1), setAsClosed(o), void (a ? (lockPaneForFX(o, !0), c.hide(u.fxName_close, u.fxSettings_close, u.fxSpeed_close, function () {
              lockPaneForFX(o, !1), d.isClosed && n(), t()
            })) : (_hidePane(o), n(), t())))
          })
        }
      }
    }; var setAsClosed = function (t) {
      if ($Rs[t]) {
        var e = ($Ps[t], $Rs[t]); var i = $Ts[t]; var s = options[t]; var n = state[t]; var o = _c[t].side; var a = s.resizerClass
        var r = s.togglerClass; var l = '-' + t; var h = '-open'; var c = '-sliding'; var u = '-closed'
        e.css(o, sC.inset[o]).removeClass(a + h + ' ' + a + l + h).removeClass(a + c + ' ' + a + l + c).addClass(a + u + ' ' + a + l + u), n.isHidden && e.hide(), s.resizable && $.layout.plugins.draggable && e.draggable('disable').removeClass('ui-state-disabled').css('cursor', 'default').attr('title', ''), i && (i.removeClass(r + h + ' ' + r + l + h).addClass(r + u + ' ' + r + l + u).attr('title', s.tips.Open), i.children('.content-open').hide(), i.children('.content-closed').css('display', 'block')), syncPinBtns(t, !1), state.initialized && sizeHandles()
      }
    }; var open = function (t, e, i, s) {
      function n () {
        c.isMoving = !1, _fixIframe(r), c.isSliding || sizeMidPanes(_c[r].dir == 'vert' ? 'center' : '', !1), setAsOpen(r)
      }

      if (isInitialized()) {
        var o; var a; var r = evtPane.call(this, t); var l = $Ps[r]; var h = ($Rs[r], $Ts[r], options[r]); var c = state[r]
        _c[r], r !== 'center' && $N.queue(function (t) {
          if (!l || !h.resizable && !h.closable && !c.isShowing || c.isVisible && !c.isSliding) return t()
          if (c.isHidden && !c.isShowing) return t(), void show(r, !0)
          c.autoResize && c.size != h.size ? sizePane(r, h.size, !0, !0, !0) : setSizeLimits(r, e)
          var u = _runCallbacks('onopen_start', r)
          return u === 'abort' ? t() : (u !== 'NC' && setSizeLimits(r, e), c.minSize > c.maxSize ? (syncPinBtns(r, !1), !s && h.tips.noRoomToOpen && alert(h.tips.noRoomToOpen), t()) : (e ? bindStopSlidingEvents(r, !0) : c.isSliding ? bindStopSlidingEvents(r, !1) : h.slidable && bindStartSlidingEvents(r, !1), c.noRoom = !1, makePaneFit(r), a = c.isShowing, delete c.isShowing, o = !i && c.isClosed && h.fxName_open != 'none', c.isMoving = !0, c.isVisible = !0, c.isClosed = !1, a && (c.isHidden = !1), void (o ? (lockPaneForFX(r, !0), l.show(h.fxName_open, h.fxSettings_open, h.fxSpeed_open, function () {
            lockPaneForFX(r, !1), c.isVisible && n(), t()
          })) : (_showPane(r), n(), t()))))
        })
      }
    }; var setAsOpen = function (t, e) {
      var i = $Ps[t]; var s = $Rs[t]; var n = $Ts[t]; var o = options[t]; var a = state[t]; var r = _c[t].side; var l = o.resizerClass
      var h = o.togglerClass; var c = '-' + t; var u = '-open'; var d = '-closed'; var p = '-sliding'
      s.css(r, sC.inset[r] + getPaneSize(t)).removeClass(l + d + ' ' + l + c + d).addClass(l + u + ' ' + l + c + u), a.isSliding ? s.addClass(l + p + ' ' + l + c + p) : s.removeClass(l + p + ' ' + l + c + p), removeHover(0, s), o.resizable && $.layout.plugins.draggable ? s.draggable('enable').css('cursor', o.resizerCursor).attr('title', o.tips.Resize) : a.isSliding || s.css('cursor', 'default'), n && (n.removeClass(h + d + ' ' + h + c + d).addClass(h + u + ' ' + h + c + u).attr('title', o.tips.Close), removeHover(0, n), n.children('.content-closed').hide(), n.children('.content-open').css('display', 'block')), syncPinBtns(t, !a.isSliding), $.extend(a, elDims(i)), state.initialized && (sizeHandles(), sizeContent(t, !0)), !e && (state.initialized || o.triggerEventsOnLoad) && i.is(':visible') && (_runCallbacks('onopen_end', t), a.isShowing && _runCallbacks('onshow_end', t), state.initialized && _runCallbacks('onresize_end', t))
    }; var slideOpen = function (t) {
      function e () {
        n.isClosed ? n.isMoving || open(s, !0) : bindStopSlidingEvents(s, !0)
      }

      if (isInitialized()) {
        var i = evtObj(t); var s = evtPane.call(this, t); var n = state[s]; var o = options[s].slideDelay_open
        s !== 'center' && (i && i.stopImmediatePropagation(), n.isClosed && i && i.type === 'mouseenter' && o > 0 ? timer.set(s + '_openSlider', e, o) : e())
      }
    }; var slideClose = function (t) {
      function e () {
        o.isClosed ? bindStopSlidingEvents(s, !1) : o.isMoving || close(s)
      }

      if (isInitialized()) {
        var i = evtObj(t); var s = evtPane.call(this, t); var n = options[s]; var o = state[s]; var a = o.isMoving ? 1e3 : 300
        if (s !== 'center' && !o.isClosed && !o.isResizing) {
          if (n.slideTrigger_close === 'click') e(); else {
            if (n.preventQuickSlideClose && o.isMoving) return
            if (n.preventPrematureSlideClose && i && $.layout.isMouseOverElem(i, $Ps[s])) return
            i ? timer.set(s + '_closeSlider', e, max(n.slideDelay_close, a)) : e()
          }
        }
      }
    }; var slideToggle = function (t) {
      var e = evtPane.call(this, t)
      toggle(e, !0)
    }; var lockPaneForFX = function (t, e) {
      var i = $Ps[t]; var s = state[t]; var n = options[t]; var o = options.zIndexes
      e ? (showMasks(t, {
        animation: !0,
        objectsOnly: !0
      }), i.css({ zIndex: o.pane_animate }), t == 'south' ? i.css({ top: sC.inset.top + sC.innerHeight - i.outerHeight() }) : t == 'east' && i.css({ left: sC.inset.left + sC.innerWidth - i.outerWidth() })) : (hideMasks(), i.css({ zIndex: s.isSliding ? o.pane_sliding : o.pane_normal }), t == 'south' ? i.css({ top: 'auto' }) : t != 'east' || i.css('left').match(/\-99999/) || i.css({ left: 'auto' }), browser.msie && n.fxOpacityFix && n.fxName_open != 'slide' && i.css('filter') && i.css('opacity') == 1 && i[0].style.removeAttribute('filter'))
    }; var bindStartSlidingEvents = function (t, e) {
      var i = options[t]; var s = ($Ps[t], $Rs[t]); var n = i.slideTrigger_open.toLowerCase()
      !s || e && !i.slidable || (n.match(/mouseover/) ? n = i.slideTrigger_open = 'mouseenter' : n.match(/(click|dblclick|mouseenter)/) || (n = i.slideTrigger_open = 'click'), i.resizerDblClickToggle && n.match(/click/) && s[e ? 'unbind' : 'bind']('dblclick.' + sID, toggle), s[e ? 'bind' : 'unbind'](n + '.' + sID, slideOpen).css('cursor', e ? i.sliderCursor : 'default').attr('title', e ? i.tips.Slide : ''))
    }; var bindStopSlidingEvents = function (t, e) {
      function i (e) {
        timer.clear(t + '_closeSlider'), e.stopPropagation()
      }

      var s = options[t]; var n = state[t]; var o = (_c[t], options.zIndexes); var a = s.slideTrigger_close.toLowerCase()
      var r = e ? 'bind' : 'unbind'; var l = $Ps[t]; var h = $Rs[t]
      timer.clear(t + '_closeSlider'), e ? (n.isSliding = !0, state.panesSliding[t] = !0, bindStartSlidingEvents(t, !1)) : (n.isSliding = !1, delete state.panesSliding[t]), l.css('zIndex', e ? o.pane_sliding : o.pane_normal), h.css('zIndex', e ? o.pane_sliding + 2 : o.resizer_normal), a.match(/(click|mouseleave)/) || (a = s.slideTrigger_close = 'mouseleave'), h[r](a, slideClose), a === 'mouseleave' && (l[r]('mouseleave.' + sID, slideClose), h[r]('mouseenter.' + sID, i), l[r]('mouseenter.' + sID, i)), e ? a !== 'click' || s.resizable || (h.css('cursor', e ? s.sliderCursor : 'default'), h.attr('title', e ? s.tips.Close : '')) : timer.clear(t + '_closeSlider')
    }; var makePaneFit = function (t, e, i, s) {
      var n = options[t]; var o = state[t]; var a = _c[t]; var r = $Ps[t]; var l = $Rs[t]; var h = a.dir === 'vert'; var c = !1
      if ((t === 'center' || h && o.noVerticalRoom) && (c = o.maxHeight >= 0, c && o.noRoom ? (_showPane(t), l && l.show(), o.isVisible = !0, o.noRoom = !1, h && (o.noVerticalRoom = !1), _fixIframe(t)) : c || o.noRoom || (_hidePane(t), l && l.hide(), o.isVisible = !1, o.noRoom = !0)), t === 'center') ; else if (o.minSize <= o.maxSize) {
        if (c = !0, o.size > o.maxSize) sizePane(t, o.maxSize, i, !0, s); else if (o.size < o.minSize) sizePane(t, o.minSize, i, !0, s); else if (l && o.isVisible && r.is(':visible')) {
          var u = o.size + sC.inset[a.side]
          $.layout.cssNum(l, a.side) != u && l.css(a.side, u)
        }
        o.noRoom && (o.wasOpen && n.closable ? n.autoReopen ? open(t, !1, !0, !0) : o.noRoom = !1 : show(t, o.wasOpen, !0, !0))
      } else o.noRoom || (o.noRoom = !0, o.wasOpen = !o.isClosed && !o.isSliding, o.isClosed || (n.closable ? close(t, !0, !0) : hide(t, !0)))
    }; var manualSizePane = function (t, e, i, s, n) {
      if (isInitialized()) {
        var o = evtPane.call(this, t); var a = options[o]; var r = state[o]
        var l = n || a.livePaneResizing && !r.isResizing
        o !== 'center' && (r.autoResize = !1, sizePane(o, e, i, s, l))
      }
    }; var sizePane = function (t, e, i, s, n) {
      function o () {
        for (var t = f === 'width' ? u.outerWidth() : u.outerHeight(), s = [{
            pane: l,
            count: 1,
            target: e,
            actual: t,
            correct: e === t,
            attempt: e,
            cssSize: r
          }], o = s[0], h = {}, m = 'Inaccurate size after resizing the ' + l + '-pane.'; !(o.correct || (h = {
            pane: l,
            count: o.count + 1,
            target: e
          }, o.actual > e ? h.attempt = max(0, o.attempt - (o.actual - e)) : h.attempt = max(0, o.attempt + (e - o.actual)), h.cssSize = cssSize(l, h.attempt), u.css(f, h.cssSize), h.actual = f == 'width' ? u.outerWidth() : u.outerHeight(), h.correct = e === h.actual, s.length === 1 && (_log(m, !1, !0), _log(o, !1, !0)), _log(h, !1, !0), s.length > 3));) s.push(h), o = s[s.length - 1]
        c.size = e, $.extend(c, elDims(u)), c.isVisible && u.is(':visible') && (d && d.css(p, e + sC.inset[p]), sizeContent(l)), !i && !g && state.initialized && c.isVisible && _runCallbacks('onresize_end', l), i || (c.isSliding || sizeMidPanes(_c[l].dir == 'horz' ? '' : 'center', g, n), sizeHandles())
        var v = _c.oppositeEdge[l]
        e < a && state[v].noRoom && (setSizeLimits(v), makePaneFit(v, !1, i)), s.length > 1 && _log(m + '\nSee the Error Console for details.', !0, !0)
      }

      if (isInitialized()) {
        var a; var r; var l = evtPane.call(this, t); var h = options[l]; var c = state[l]; var u = $Ps[l]; var d = $Rs[l]
        var p = _c[l].side; var f = _c[l].sizeType.toLowerCase()
        var g = c.isResizing && !h.triggerEventsDuringLiveResize; var m = !0 !== s && h.animatePaneSizing
        l !== 'center' && $N.queue(function (t) {
          if (setSizeLimits(l), a = c.size, e = _parseSize(l, e), e = max(e, _parseSize(l, h.minSize)), (e = min(e, c.maxSize)) < c.minSize) return t(), void makePaneFit(l, !1, i)
          if (!n && e === a) return t()
          if (c.newSize = e, !i && state.initialized && c.isVisible && _runCallbacks('onresize_start', l), r = cssSize(l, e), m && u.is(':visible')) {
            var s = $.layout.effects.size[l] || $.layout.effects.size.all
            var d = h.fxSettings_size.easing || s.easing; var p = options.zIndexes; var g = {}
            g[f] = r + 'px', c.isMoving = !0, u.css({ zIndex: p.pane_animate }).show().animate(g, h.fxSpeed_size, d, function () {
              u.css({ zIndex: c.isSliding ? p.pane_sliding : p.pane_normal }), c.isMoving = !1, delete c.newSize, o(), t()
            })
          } else u.css(f, r), delete c.newSize, u.is(':visible') ? o() : c.size = e, t()
        })
      }
    }; var sizeMidPanes = function (t, e, i) {
      t = (t || 'east,west,center').split(','), $.each(t, function (t, s) {
        if ($Ps[s]) {
          var n = options[s]; var o = state[s]; var a = $Ps[s]; var r = ($Rs[s], !0); var l = {}
          var h = $.layout.showInvisibly(a); var c = calcNewCenterPaneDims()
          if ($.extend(o, elDims(a)), s === 'center') {
            if (!i && o.isVisible && c.width === o.outerWidth && c.height === o.outerHeight) return a.css(h), !0
            if ($.extend(o, cssMinDims(s), {
              maxWidth: c.width,
              maxHeight: c.height
            }), l = c, o.newWidth = l.width, o.newHeight = l.height, l.width = cssW(a, l.width), l.height = cssH(a, l.height), r = l.width >= 0 && l.height >= 0, !state.initialized && n.minWidth > c.width) {
              var u = n.minWidth - o.outerWidth; var d = options.east.minSize || 0
              var p = options.west.minSize || 0; var f = state.east.size; var g = state.west.size; var m = f; var v = g
              if (u > 0 && state.east.isVisible && f > d && (m = max(f - d, f - u), u -= f - m), u > 0 && state.west.isVisible && g > p && (v = max(g - p, g - u), u -= g - v), u === 0) return f && f != d && sizePane('east', m, !0, !0, i), g && g != p && sizePane('west', v, !0, !0, i), sizeMidPanes('center', e, i), void a.css(h)
            }
          } else {
            if (o.isVisible && !o.noVerticalRoom && $.extend(o, elDims(a), cssMinDims(s)), !i && !o.noVerticalRoom && c.height === o.outerHeight) return a.css(h), !0
            l.top = c.top, l.bottom = c.bottom, o.newSize = c.height, l.height = cssH(a, c.height), o.maxHeight = l.height, (r = o.maxHeight >= 0) || (o.noVerticalRoom = !0)
          }
          if (r ? (!e && state.initialized && _runCallbacks('onresize_start', s), a.css(l), s !== 'center' && sizeHandles(s), !o.noRoom || o.isClosed || o.isHidden || makePaneFit(s), o.isVisible && ($.extend(o, elDims(a)), state.initialized && sizeContent(s))) : !o.noRoom && o.isVisible && makePaneFit(s), a.css(h), delete o.newSize, delete o.newWidth, delete o.newHeight, !o.isVisible) return !0
          if (s === 'center') {
            var b = browser.isIE6 || !browser.boxModel
            $Ps.north && (b || state.north.tagName == 'IFRAME') && $Ps.north.css('width', cssW($Ps.north, sC.innerWidth)), $Ps.south && (b || state.south.tagName == 'IFRAME') && $Ps.south.css('width', cssW($Ps.south, sC.innerWidth))
          }
          !e && state.initialized && _runCallbacks('onresize_end', s)
        }
      })
    }; var resizeAll = function (t) {
      sC.innerWidth, sC.innerHeight
      if (evtPane(t), $N.is(':visible')) {
        if (!state.initialized) return void _initLayoutElements()
        if (!0 === t && $.isPlainObject(options.outset) && $N.css(options.outset), $.extend(sC, elDims($N, options.inset)), sC.outerHeight) {
          if (!0 === t && setPanePosition(), !1 === _runCallbacks('onresizeall_start')) return !1
          var e, i, s
          sC.innerHeight, sC.innerWidth, $.each(['south', 'north', 'east', 'west'], function (t, e) {
            $Ps[e] && (i = options[e], s = state[e], s.autoResize && s.size != i.size ? sizePane(e, i.size, !0, !0, !0) : (setSizeLimits(e), makePaneFit(e, !1, !0, !0)))
          }), sizeMidPanes('', !0, !0), sizeHandles(), $.each(_c.allPanes, function (t, i) {
            (e = $Ps[i]) && state[i].isVisible && _runCallbacks('onresize_end', i)
          }), _runCallbacks('onresizeall_end')
        }
      }
    }; var resizeChildren = function (t, e) {
      var i = evtPane.call(this, t)
      if (options[i].resizeChildren) {
        e || refreshChildren(i)
        var s = children[i]
        $.isPlainObject(s) && $.each(s, function (t, e) {
          e.destroyed || e.resizeAll()
        })
      }
    }; var sizeContent = function (t, e) {
      if (isInitialized()) {
        var i = evtPane.call(this, t)
        i = i ? i.split(',') : _c.allPanes, $.each(i, function (t, i) {
          function s (t) {
            return max(l.css.paddingBottom, parseInt(t.css('marginBottom'), 10) || 0)
          }

          function n () {
            var t = options[i].contentIgnoreSelector
            var e = a.nextAll().not('.ui-layout-mask').not(t || ':lt(0)'); var n = e.filter(':visible')
            var o = n.filter(':last');
            (h = {
              top: a[0].offsetTop,
              height: a.outerHeight(),
              numFooters: e.length,
              hiddenFooters: e.length - n.length,
              spaceBelow: 0
            }).spaceAbove = h.top, h.bottom = h.top + h.height, o.length ? h.spaceBelow = o[0].offsetTop + o.outerHeight() - h.bottom + s(o) : h.spaceBelow = s(a)
          }

          var o = $Ps[i]; var a = $Cs[i]; var r = options[i]; var l = state[i]; var h = l.content
          if (!o || !a || !o.is(':visible')) return !0
          if ((a.length || (initContent(i, !1), a)) && !1 !== _runCallbacks('onsizecontent_start', i)) {
            (!l.isMoving && !l.isResizing || r.liveContentResizing || e || void 0 == h.top) && (n(), h.hiddenFooters > 0 && o.css('overflow') === 'hidden' && (o.css('overflow', 'visible'), n(), o.css('overflow', 'hidden')))
            var c = l.innerHeight - (h.spaceAbove - l.css.paddingTop) - (h.spaceBelow - l.css.paddingBottom)
            a.is(':visible') && h.height == c || (setOuterHeight(a, c, !0), h.height = c), state.initialized && _runCallbacks('onsizecontent_end', i)
          }
        })
      }
    }; var sizeHandles = function (t) {
      var e = evtPane.call(this, t)
      e = e ? e.split(',') : _c.borderPanes, $.each(e, function (t, e) {
        var i; var s = options[e]; var n = state[e]; var o = $Ps[e]; var a = $Rs[e]; var r = $Ts[e]
        if (o && a) {
          var l; var h; var c; var u = _c[e].dir; var d = n.isClosed ? '_closed' : '_open'; var p = s['spacing' + d]
          var f = s['togglerAlign' + d]; var g = s['togglerLength' + d]
          if (p === 0) return void a.hide()
          if (n.noRoom || n.isHidden || a.show(), u === 'horz' ? (l = sC.innerWidth, n.resizerLength = l, h = $.layout.cssNum(o, 'left'), a.css({
            width: cssW(a, l),
            height: cssH(a, p),
            left: h > -9999 ? h : sC.inset.left
          })) : (l = o.outerHeight(), n.resizerLength = l, a.css({
            height: cssH(a, l),
            width: cssW(a, p),
            top: sC.inset.top + getPaneSize('north', !0)
          })), removeHover(s, a), r) {
            if (g === 0 || n.isSliding && s.hideTogglerOnSlide) return void r.hide()
            if (r.show(), !(g > 0) || g === '100%' || g > l) g = l, c = 0; else if (isStr(f)) {
              switch (f) {
                case 'top':
                case 'left':
                  c = 0
                  break
                case 'bottom':
                case 'right':
                  c = l - g
                  break
                case 'middle':
                case 'center':
                default:
                  c = round((l - g) / 2)
              }
            } else {
              var m = parseInt(f, 10)
              c = f >= 0 ? m : l - g + m
            }
            if (u === 'horz') {
              var v = cssW(r, g)
              r.css({
                width: v,
                height: cssH(r, p),
                left: c,
                top: 0
              }), r.children('.content').each(function () {
                (i = $(this)).css('marginLeft', round((v - i.outerWidth()) / 2))
              })
            } else {
              var b = cssH(r, g)
              r.css({
                height: b,
                width: cssW(r, p),
                top: c,
                left: 0
              }), r.children('.content').each(function () {
                (i = $(this)).css('marginTop', round((b - i.outerHeight()) / 2))
              })
            }
            removeHover(0, r)
          }
          state.initialized || !s.initHidden && !n.isHidden || (a.hide(), r && r.hide())
        }
      })
    }; var enableClosable = function (t) {
      if (isInitialized()) {
        var e = evtPane.call(this, t); var i = $Ts[e]; var s = options[e]
        i && (s.closable = !0, i.bind('click.' + sID, function (t) {
          t.stopPropagation(), toggle(e)
        }).css('visibility', 'visible').css('cursor', 'pointer').attr('title', state[e].isClosed ? s.tips.Open : s.tips.Close).show())
      }
    }; var disableClosable = function (t, e) {
      if (isInitialized()) {
        var i = evtPane.call(this, t); var s = $Ts[i]
        s && (options[i].closable = !1, state[i].isClosed && open(i, !1, !0), s.unbind('.' + sID).css('visibility', e ? 'hidden' : 'visible').css('cursor', 'default').attr('title', ''))
      }
    }; var enableSlidable = function (t) {
      if (isInitialized()) {
        var e = evtPane.call(this, t); var i = $Rs[e]
        i && i.data('draggable') && (options[e].slidable = !0, state[e].isClosed && bindStartSlidingEvents(e, !0))
      }
    }; var disableSlidable = function (t) {
      if (isInitialized()) {
        var e = evtPane.call(this, t); var i = $Rs[e]
        i && (options[e].slidable = !1, state[e].isSliding ? close(e, !1, !0) : (bindStartSlidingEvents(e, !1), i.css('cursor', 'default').attr('title', ''), removeHover(null, i[0])))
      }
    }; var enableResizable = function (t) {
      if (isInitialized()) {
        var e = evtPane.call(this, t); var i = $Rs[e]; var s = options[e]
        i && i.data('draggable') && (s.resizable = !0, i.draggable('enable'), state[e].isClosed || i.css('cursor', s.resizerCursor).attr('title', s.tips.Resize))
      }
    }; var disableResizable = function (t) {
      if (isInitialized()) {
        var e = evtPane.call(this, t); var i = $Rs[e]
        i && i.data('draggable') && (options[e].resizable = !1, i.draggable('disable').css('cursor', 'default').attr('title', ''), removeHover(null, i[0]))
      }
    }; var swapPanes = function (t, e) {
      function i (t) {
        var e = $Ps[t]; var i = $Cs[t]
        return !!e && {
          pane: t,
          P: !!e && e[0],
          C: !!i && i[0],
          state: $.extend(!0, {}, state[t]),
          options: $.extend(!0, {}, options[t])
        }
      }

      function s (t, e) {
        if (t) {
          var i; var s; var n = t.P; var o = t.C; var a = t.pane; var l = _c[e]; var h = $.extend(!0, {}, state[e]); var c = options[e]
          var u = { resizerCursor: c.resizerCursor }
          $.each('fxName,fxSpeed,fxSettings'.split(','), function (t, e) {
            u[e + '_open'] = c[e + '_open'], u[e + '_close'] = c[e + '_close'], u[e + '_size'] = c[e + '_size']
          }), $Ps[e] = $(n).data({
            layoutPane: Instance[e],
            layoutEdge: e
          }).css(_c.hidden).css(l.cssReq), $Cs[e] = !!o && $(o), options[e] = $.extend(!0, {}, t.options, u), state[e] = $.extend(!0, {}, t.state), i = new RegExp(c.paneClass + '-' + a, 'g'), n.className = n.className.replace(i, c.paneClass + '-' + e), initHandles(e), l.dir != _c[a].dir ? (s = r[e] || 0, setSizeLimits(e), s = max(s, state[e].minSize), manualSizePane(e, s, !0, !0)) : $Rs[e].css(l.side, sC.inset[l.side] + (state[e].isVisible ? getPaneSize(e) : 0)), t.state.isVisible && !h.isVisible ? setAsOpen(e, !0) : (setAsClosed(e), bindStartSlidingEvents(e, !0)), t = null
        }
      }

      if (isInitialized()) {
        var n = evtPane.call(this, t)
        if (state[n].edge = e, state[e].edge = n, !1 === _runCallbacks('onswap_start', n) || !1 === _runCallbacks('onswap_start', e)) return state[n].edge = n, void (state[e].edge = e)
        var o = i(n); var a = i(e); var r = {}
        r[n] = o ? o.state.size : 0, r[e] = a ? a.state.size : 0, $Ps[n] = !1, $Ps[e] = !1, state[n] = {}, state[e] = {}, $Ts[n] && $Ts[n].remove(), $Ts[e] && $Ts[e].remove(), $Rs[n] && $Rs[n].remove(), $Rs[e] && $Rs[e].remove(), $Rs[n] = $Rs[e] = $Ts[n] = $Ts[e] = !1, s(o, e), s(a, n), o = a = r = null, $Ps[n] && $Ps[n].css(_c.visible), $Ps[e] && $Ps[e].css(_c.visible), resizeAll(), _runCallbacks('onswap_end', n), _runCallbacks('onswap_end', e)
      }
    }; var syncPinBtns = function (t, e) {
      $.layout.plugins.buttons && $.each(state[t].pins, function (i, s) {
        $.layout.buttons.setPinState(Instance, $(s), t, e)
      })
    }; var $N = $(this).eq(0)
    if (!$N.length) return _log(options.errors.containerMissing)
    if ($N.data('layoutContainer') && $N.data('layout')) return $N.data('layout')
    var $Ps = {}; var $Cs = {}; var $Rs = {}; var $Ts = {}; var $Ms = $([]); var sC = state.container; var sID = state.id; var Instance = {
      options: options,
      state: state,
      container: $N,
      panes: $Ps,
      contents: $Cs,
      resizers: $Rs,
      togglers: $Ts,
      hide: hide,
      show: show,
      toggle: toggle,
      open: open,
      close: close,
      slideOpen: slideOpen,
      slideClose: slideClose,
      slideToggle: slideToggle,
      setSizeLimits: setSizeLimits,
      _sizePane: sizePane,
      sizePane: manualSizePane,
      sizeContent: sizeContent,
      swapPanes: swapPanes,
      showMasks: showMasks,
      hideMasks: hideMasks,
      initContent: initContent,
      addPane: addPane,
      removePane: removePane,
      createChildren: createChildren,
      refreshChildren: refreshChildren,
      enableClosable: enableClosable,
      disableClosable: disableClosable,
      enableSlidable: enableSlidable,
      disableSlidable: disableSlidable,
      enableResizable: enableResizable,
      disableResizable: disableResizable,
      allowOverflow: allowOverflow,
      resetOverflow: resetOverflow,
      destroy: destroy,
      initPanes: isInitialized,
      resizeAll: resizeAll,
      runCallbacks: _runCallbacks,
      hasParentLayout: !1,
      children: children,
      north: !1,
      south: !1,
      west: !1,
      east: !1,
      center: !1
    }
    return _create() === 'cancel' ? null : Instance
  }
}(jQuery)), (function (t) {
  t.layout && (t.ui || (t.ui = {}), t.ui.cookie = {
    acceptsCookies: !!navigator.cookieEnabled,
    read: function (e) {
      var i; var s; var n; var o = document.cookie; var a = o ? o.split(';') : []
      for (n = 0; i = a[n]; n++) if ((s = t.trim(i).split('='))[0] == e) return decodeURIComponent(s[1])
      return null
    },
    write: function (e, i, s) {
      var n = ''; var o = ''; var a = !1; var r = s || {}; var l = r.expires || null; var h = t.type(l)
      h === 'date' ? o = l : h === 'string' && l > 0 && (l = parseInt(l, 10), h = 'number'), h === 'number' && (o = new Date(), l > 0 ? o.setDate(o.getDate() + l) : (o.setFullYear(1970), a = !0)), o && (n += ';expires=' + o.toUTCString()), r.path && (n += ';path=' + r.path), r.domain && (n += ';domain=' + r.domain), r.secure && (n += ';secure'), document.cookie = e + '=' + (a ? '' : encodeURIComponent(i)) + n
    },
    clear: function (e) {
      t.ui.cookie.write(e, '', { expires: -1 })
    }
  }, t.cookie || (t.cookie = function (e, i, s) {
    var n = t.ui.cookie
    if (i === null) n.clear(e); else {
      if (void 0 === i) return n.read(e)
      n.write(e, i, s)
    }
  }), t.layout.plugins.stateManagement = !0, t.layout.defaults.stateManagement = {
    enabled: !1,
    autoSave: !0,
    autoLoad: !0,
    animateLoad: !0,
    includeChildren: !0,
    stateKeys: 'north.size,south.size,east.size,west.size,north.isClosed,south.isClosed,east.isClosed,west.isClosed,north.isHidden,south.isHidden,east.isHidden,west.isHidden',
    cookie: { name: '', domain: '', path: '', expires: '', secure: !1 }
  }, t.layout.optionsMap.layout.push('stateManagement'), t.layout.config.optionRootKeys.push('stateManagement'), t.layout.state = {
    saveCookie: function (e, i, s) {
      var n = e.options; var o = n.stateManagement; var a = t.extend(!0, {}, o.cookie, s || null)
      var r = e.state.stateData = e.readState(i || o.stateKeys)
      return t.ui.cookie.write(a.name || n.name || 'Layout', t.layout.state.encodeJSON(r), a), t.extend(!0, {}, r)
    },
    deleteCookie: function (e) {
      var i = e.options
      t.ui.cookie.clear(i.stateManagement.cookie.name || i.name || 'Layout')
    },
    readCookie: function (e) {
      var i = e.options; var s = t.ui.cookie.read(i.stateManagement.cookie.name || i.name || 'Layout')
      return s ? t.layout.state.decodeJSON(s) : {}
    },
    loadCookie: function (e) {
      var i = t.layout.state.readCookie(e)
      return i && !t.isEmptyObject(i) && (e.state.stateData = t.extend(!0, {}, i), e.loadState(i)), i
    },
    loadState: function (e, i, n) {
      if (t.isPlainObject(i) && !t.isEmptyObject(i)) {
        i = e.state.stateData = t.layout.transformData(i)
        var o = e.options.stateManagement
        if (n = t.extend({ animateLoad: !1, includeChildren: o.includeChildren }, n), e.state.initialized) {
          var a; var r; var l; var h; var c = !n.animateLoad
          if (t.each(t.layout.config.borderPanes, function (n, o) {
            p = i[o], t.isPlainObject(p) && (s = p.size, a = p.initClosed, r = p.initHidden, ar = p.autoResize, l = e.state[o], h = l.isVisible, ar && (l.autoResize = ar), h || e._sizePane(o, s, !1, !1, !1), !0 === r ? e.hide(o, c) : !0 === a ? e.close(o, !1, c) : !1 === a ? e.open(o, !1, c) : !1 === r && e.show(o, !1, c), h && e._sizePane(o, s, !1, !1, c))
          }), n.includeChildren) {
            var u, d
            t.each(e.children, function (e, s) {
              (u = i[e] ? i[e].children : 0) && s && t.each(s, function (t, e) {
                d = u[t], e && d && e.loadState(d)
              })
            })
          }
        } else {
          var p = t.extend(!0, {}, i)
          t.each(t.layout.config.allPanes, function (t, e) {
            p[e] && delete p[e].children
          }), t.extend(!0, e.options, p)
        }
      }
    },
    readState: function (e, i) {
      t.type(i) === 'string' && (i = { keys: i }), i || (i = {})
      var s; var n; var o; var a; var r; var l; var h; var c = e.options.stateManagement; var u = i.includeChildren
      var d = void 0 !== u ? u : c.includeChildren; var p = i.stateKeys || c.stateKeys
      var f = { isClosed: 'initClosed', isHidden: 'initHidden' }; var g = e.state; var m = t.layout.config.allPanes; var v = {}
      t.isArray(p) && (p = p.join(','))
      for (var b = 0, _ = (p = p.replace(/__/g, '.').split(',')).length; b < _; b++) s = p[b].split('.'), n = s[0], o = s[1], t.inArray(n, m) < 0 || void 0 != (a = g[n][o]) && (o == 'isClosed' && g[n].isSliding && (a = !0), (v[n] || (v[n] = {}))[f[o] ? f[o] : o] = a)
      return d && t.each(m, function (i, s) {
        l = e.children[s], r = g.stateData[s], t.isPlainObject(l) && !t.isEmptyObject(l) && ((h = v[s] || (v[s] = {})).children || (h.children = {}), t.each(l, function (e, i) {
          i.state.initialized ? h.children[e] = t.layout.state.readState(i) : r && r.children && r.children[e] && (h.children[e] = t.extend(!0, {}, r.children[e]))
        }))
      }), v
    },
    encodeJSON: function (e) {
      return ((window.JSON || {}).stringify || function (e) {
        var i; var s; var n; var o = []; var a = 0; var r = t.isArray(e)
        for (i in e) s = e[i], n = typeof s, n == 'string' ? s = '"' + s + '"' : n == 'object' && (s = parse(s)), o[a++] = (r ? '' : '"' + i + '":') + s
        return (r ? '[' : '{') + o.join(',') + (r ? ']' : '}')
      })(e)
    },
    decodeJSON: function (e) {
      try {
        return t.parseJSON ? t.parseJSON(e) : window.eval('(' + e + ')') || {}
      } catch (t) {
        return {}
      }
    },
    _create: function (e) {
      var i = t.layout.state; var s = e.options.stateManagement
      if (t.extend(e, {
        readCookie: function () {
          return i.readCookie(e)
        },
        deleteCookie: function () {
          i.deleteCookie(e)
        },
        saveCookie: function (t, s) {
          return i.saveCookie(e, t, s)
        },
        loadCookie: function () {
          return i.loadCookie(e)
        },
        loadState: function (t, s) {
          i.loadState(e, t, s)
        },
        readState: function (t) {
          return i.readState(e, t)
        },
        encodeJSON: i.encodeJSON,
        decodeJSON: i.decodeJSON
      }), e.state.stateData = {}, s.autoLoad) {
        if (t.isPlainObject(s.autoLoad)) t.isEmptyObject(s.autoLoad) || e.loadState(s.autoLoad); else if (s.enabled) {
          if (t.isFunction(s.autoLoad)) {
            var n = {}
            try {
              n = s.autoLoad(e, e.state, e.options, e.options.name || '')
            } catch (t) {
            }
            n && t.isPlainObject(n) && !t.isEmptyObject(n) && e.loadState(n)
          } else e.loadCookie()
        }
      }
    },
    _unload: function (e) {
      var i = e.options.stateManagement
      if (i.enabled && i.autoSave) {
        if (t.isFunction(i.autoSave)) {
          try {
            i.autoSave(e, e.state, e.options, e.options.name || '')
          } catch (t) {
          }
        } else e.saveCookie()
      }
    }
  }, t.layout.onCreate.push(t.layout.state._create), t.layout.onUnload.push(t.layout.state._unload))
}(jQuery)), (function (t) {
  t.layout && (t.layout.plugins.buttons = !0, t.layout.defaults.autoBindCustomButtons = !1, t.layout.optionsMap.layout.push('autoBindCustomButtons'), t.layout.buttons = {
    config: { borderPanes: 'north,south,west,east' },
    init: function (e) {
      var i; var s = e.options.name || ''
      t.each('toggle,open,close,pin,toggle-slide,open-slide'.split(','), function (n, o) {
        t.each(t.layout.buttons.config.borderPanes.split(','), function (n, a) {
          t('.ui-layout-button-' + o + '-' + a).each(function () {
            void 0 != (i = t(this).data('layoutName') || t(this).attr('layoutName')) && i !== s || e.bindButton(this, o, a)
          })
        })
      })
    },
    get: function (e, i, s, n) {
      var o = t(i); var a = e.options
      if (o.length && t.layout.buttons.config.borderPanes.indexOf(s) >= 0) {
        var r = a[s].buttonClass + '-' + n
        o.addClass(r + ' ' + r + '-' + s).data('layoutName', a.name)
      }
      return o
    },
    bind: function (e, i, s, n) {
      var o = t.layout.buttons
      switch (s.toLowerCase()) {
        case 'toggle':
          o.addToggle(e, i, n)
          break
        case 'open':
          o.addOpen(e, i, n)
          break
        case 'close':
          o.addClose(e, i, n)
          break
        case 'pin':
          o.addPin(e, i, n)
          break
        case 'toggle-slide':
          o.addToggle(e, i, n, !0)
          break
        case 'open-slide':
          o.addOpen(e, i, n, !0)
      }
      return e
    },
    addToggle: function (e, i, s, n) {
      return t.layout.buttons.get(e, i, s, 'toggle').click(function (t) {
        e.toggle(s, !!n), t.stopPropagation()
      }), e
    },
    addOpen: function (e, i, s, n) {
      return t.layout.buttons.get(e, i, s, 'open').attr('title', e.options[s].tips.Open).click(function (t) {
        e.open(s, !!n), t.stopPropagation()
      }), e
    },
    addClose: function (e, i, s) {
      return t.layout.buttons.get(e, i, s, 'close').attr('title', e.options[s].tips.Close).click(function (t) {
        e.close(s), t.stopPropagation()
      }), e
    },
    addPin: function (e, i, s) {
      var n = t.layout.buttons.get(e, i, s, 'pin')
      if (n.length) {
        var o = e.state[s]
        n.click(function (i) {
          t.layout.buttons.setPinState(e, t(this), s, o.isSliding || o.isClosed), o.isSliding || o.isClosed ? e.open(s) : e.close(s), i.stopPropagation()
        }), t.layout.buttons.setPinState(e, n, s, !o.isClosed && !o.isSliding), o.pins.push(i)
      }
      return e
    },
    setPinState: function (t, e, i, s) {
      var n = e.attr('pin')
      if (!n || s !== (n == 'down')) {
        var o = t.options[i]; var a = o.tips; var r = o.buttonClass + '-pin'; var l = r + '-' + i
        var h = r + '-up ' + l + '-up'; var c = r + '-down ' + l + '-down'
        e.attr('pin', s ? 'down' : 'up').attr('title', s ? a.Unpin : a.Pin).removeClass(s ? h : c).addClass(s ? c : h)
      }
    },
    syncPinBtns: function (e, i, s) {
      t.each(state[i].pins, function (n, o) {
        t.layout.buttons.setPinState(e, t(o), i, s)
      })
    },
    _load: function (e) {
      t.extend(e, {
        bindButton: function (i, s, n) {
          return t.layout.buttons.bind(e, i, s, n)
        },
        addToggleBtn: function (i, s, n) {
          return t.layout.buttons.addToggle(e, i, s, n)
        },
        addOpenBtn: function (i, s, n) {
          return t.layout.buttons.addOpen(e, i, s, n)
        },
        addCloseBtn: function (i, s) {
          return t.layout.buttons.addClose(e, i, s)
        },
        addPinBtn: function (i, s) {
          return t.layout.buttons.addPin(e, i, s)
        }
      })
      for (var i = 0; i < 4; i++) {
        var s = t.layout.buttons.config.borderPanes[i]
        e.state[s].pins = []
      }
      e.options.autoBindCustomButtons && t.layout.buttons.init(e)
    },
    _unload: function (t) {
    }
  }, t.layout.onLoad.push(t.layout.buttons._load))
}(jQuery)), (function (t) {
  t.layout.plugins.browserZoom = !0, t.layout.defaults.browserZoomCheckInterval = 1e3, t.layout.optionsMap.layout.push('browserZoomCheckInterval'), t.layout.browserZoom = {
    _init: function (e) {
      !1 !== t.layout.browserZoom.ratio() && t.layout.browserZoom._setTimer(e)
    },
    _setTimer: function (e) {
      if (!e.destroyed) {
        var i = e.options; var s = e.state; var n = e.hasParentLayout ? 5e3 : Math.max(i.browserZoomCheckInterval, 100)
        setTimeout(function () {
          if (!e.destroyed && i.resizeWithWindow) {
            var n = t.layout.browserZoom.ratio()
            n !== s.browserZoom && (s.browserZoom = n, e.resizeAll()), t.layout.browserZoom._setTimer(e)
          }
        }, n)
      }
    },
    ratio: function () {
      function e (t, e) {
        return (parseInt(t, 10) / parseInt(e, 10) * 100).toFixed()
      }

      var i; var s; var n; var o = window; var a = screen; var r = document; var l = r.documentElement || r.body; var h = t.layout.browser
      var c = h.version
      return !(!h.msie || c > 8) && (a.deviceXDPI && a.systemXDPI ? e(a.deviceXDPI, a.systemXDPI) : h.webkit && (i = r.body.getBoundingClientRect) ? e(i.left - i.right, r.body.offsetWidth) : h.webkit && (s = o.outerWidth) ? e(s, o.innerWidth) : !(!(s = a.width) || !(n = l.clientWidth)) && e(s, n))
    }
  }, t.layout.onReady.push(t.layout.browserZoom._init)
}(jQuery)), (function (t) {
  t.effects && (t.layout.defaults.panes.useOffscreenClose = !1, t.layout.plugins && (t.layout.plugins.effects.slideOffscreen = !0), t.layout.effects.slideOffscreen = t.extend(!0, {}, t.layout.effects.slide), t.effects.slideOffscreen = function (e) {
    return this.queue(function () {
      var i = t.effects; var s = e.options; var n = t(this); var o = n.data('layoutEdge'); var a = n.data('parentLayout').state
      var r = a[o].size; var l = this.style; var h = i.setMode(n, s.mode || 'show') == 'show'; var c = s.direction || 'left'
      var u = c == 'up' || c == 'down' ? 'top' : 'left'; var d = c == 'up' || c == 'left'
      var p = t.layout.config.offscreenCSS || {}; var f = t.layout.config.offscreenReset; var g = 'offscreenResetTop'
      var m = {}
      m[u] = (h ? d ? '+=' : '-=' : d ? '-=' : '+=') + r, h ? (n.data(g, {
        top: l.top,
        bottom: l.bottom
      }), d ? n.css(u, isNaN(r) ? '-' + r : -r) : c === 'right' ? n.css({
        left: a.container.layoutWidth,
        right: 'auto'
      }) : n.css({
        top: a.container.layoutHeight,
        bottom: 'auto'
      }), u === 'top' && n.css(n.data(f) || {})) : (n.data(g, {
        top: l.top,
        bottom: l.bottom
      }), n.data(f, { left: l.left, right: l.right })), n.show().animate(m, {
        queue: !1,
        duration: e.duration,
        easing: s.easing,
        complete: function () {
          n.data(g) && n.css(n.data(g)).removeData(g), h ? n.css(n.data(f) || {}).removeData(f) : n.css(p), e.callback && e.callback.apply(this, arguments), n.dequeue()
        }
      })
    })
  })
}(jQuery))
