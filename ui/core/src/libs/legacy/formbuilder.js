import FT_FormBuilderCompiler from './formbuilder-compiler'

var FT_FormBuilder = function (jqEl, options) {
  this.init = function (jqEl, options) {
    var formBuilder = this
    this.compiler = new FT_FormBuilderCompiler({
      requestComp: function (id, callback) {
        var localComp
        if (formBuilder.componentsTab.components[id]) {
          localComp = formBuilder.componentsTab.components[id]
        }
        if (localComp) {
          if (callback === undefined) {
            return localComp
          } else {
            callback(localComp)
            return
          }
        }
        return formBuilder.options.component.requestComp(id, callback)
      },
      injectFormCompiler: true,
      formCompiler: {
        getComponentsHolder: function () {
          return formBuilder.workspace
        }
      },
      i18n: options.i18n
    })
    options = $.extend({
      file: {
        requestTypes: function (callback) {}
      },
      component: {
        requestComp: function (id, callback) {},
        searchComp: function (text, callback) {},
        storeComp: function (comp, callback) {},
        deleteComp: function (id, callback) {}
      },
      componentsTab: {
        components: {}
      },
      vars: []
    }, options)
    this.fbid = this.randomId()
    this.el = jqEl
    this.el.addClass('hcbuild-main')
    if (options && options.readOnly === true) {
      this.el.addClass('read-only')
    }
    this.el.append(formBuilderHtmlConstruct)
    this.vars = options.vars
    this.options = options
    this.workspace = new FT_Workspace(this, this.el.find('#htmlFormWorkspace'))
    this.componentsTab = new FT_ComponentsTab(this, this.el.find('#htmlComponents'))
    this.builderTab = new FT_BuilderTab(this, this.el.find('#htmlComponentBuilder'))
    this.settingsTab = new FT_SettingsTab(this, this.el.find('#htmlComponentSettings'))
    this.workspace.fb.options.scrollableContainer = this.workspace.el.find('.wsbody')

    $('.hcbuilder a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
      var previous = $(e.relatedTarget).attr('href')
      if (previous === '#htmlComponents') {
        formBuilder.componentsTab.tabHidden()
      } else if (previous === '#htmlComponentBuilder') {
        formBuilder.builderTab.tabHidden()
      } else if (previous === '#htmlComponentSettings') {
        formBuilder.settingsTab.tabHidden()
      }
      var target = $(e.target).attr('href') // activated tab
      if (target === '#htmlComponents') {
        formBuilder.componentsTab.tabShown()
      } else if (target === '#htmlComponentBuilder') {
        formBuilder.builderTab.tabShown()
      } else if (target === '#htmlComponentSettings') {
        formBuilder.settingsTab.tabShown()
      }
    })

    this.workspace.loadData()
    this.componentsTab.loadData()
    this.setFileTypes = function (fTypes) {
      formBuilder.fileTypes = fTypes
    }
    if (!this.fileTypes) {
      this.options.file.requestTypes(this.setFileTypes)
    }
    var $scrollingDiv = $('.hcbuild-main .hcbuilder')
    var scrollContainerEl = this.options.scrollableContainer
    this.$sticky = this.el.find('.fb-sticky-right')
    this.activeTabObj = this.componentsTab
    formBuilder.lastScrollPos = formBuilder.options.scrollableContainer.scrollTop()
    formBuilder.options.scrollableContainer.scroll(function () {
      formBuilder.currentScrollPos = formBuilder.options.scrollableContainer.scrollTop()
      var scrollOffsetY = (-1 * (formBuilder.lastScrollPos - formBuilder.currentScrollPos))
      formBuilder.lastScrollPos = formBuilder.currentScrollPos
      formBuilder.adjustHeight(formBuilder.activeTabObj, scrollOffsetY)
    })
    this.setupCopyPasteEvent()

    if (this.options.userAllowedToEditComponents === false) {
      $('.nav-item .htmlComponentBuilder').hide()
    }
  }
  this.el = null
  this.getData = function () {
    return this.workspace.storeComponents()
  }
  this.randomId = function () {
    return this.compiler.randomId()
  }
  this.setupCopyPasteEvent = function () {
    var fb = this
    fb['vKey'] = 86
    fb['cKey'] = 67
    var isField = function (t) {
      return /input|textarea|select/i.test(t.tagName)
    }
    $(document).keydown(function (e) {
      try {
        if (!isField(e.target)) {
          var ctrlDown = e.ctrlKey || e.metaKey // Mac support
          if (ctrlDown && e.keyCode == fb.cKey) {
            return !fb.workspace.onWSCopyComponent()
          } else if (ctrlDown && e.keyCode == fb.vKey) {
            return !fb.workspace.onWSPasteComponent()
          }
        }
      } catch (err) {
        console.log(err)
      }
      return true
    })
    /** .keyup(function(e) {
            if (e.keyCode == fb.ctrlKey || e.keyCode == fb.cmdKey){}
        });**/
  }
  this.adjustHeight = function (t, scrollOffset) {}
  this.calcContainerHeight = function (t, force) {
    if (t.el.is(':visible')) {
      if (force || !t.containerHeightWithoutInnerBody) {
        if (!t.innerBody || !t.innerBody.length) {
          t.innerBody = t.el.find('.fb-inner-body')
          t.innerBody.css({
            'height': '0',
            'min-height': '0'
          })
          t.containerHeightWithoutInnerBody = t.el.children('.panel').outerHeight(true)
          t.innerBody.css({
            'height': '',
            'min-height': ''
          })
        }
        this.adjustHeight(t, 0)
      }
    }
  }
  // provide tmplId for template cache/performance
  this.compileTemplate = function (templStr, settings, tmplId) {
    return this.compiler.compileTemplate(templStr, settings, tmplId)
  }
  this.getPossibleHeight = function (pos) {
    if (!pos) {
      pos = this.el.position()
    }
    if (!this.cachedContainerHeight) {
      this.cachedContainerHeight = this.options.scrollableContainer.outerHeight(false)
    }
    if (pos.top > 0) {
      return this.cachedContainerHeight - pos.top
    }
    return this.cachedContainerHeight
  }
  this.translateSettingsAndCompileTemplate = function (tmplId, templStr, jsonSettings, callback) {
    var _ = this
    this.translateSettings(jsonSettings, tmplId, function (translatedSettings, tplId) {
      var cachedTmpl = _.compiler.getTemplateCached(tplId)
      if (!cachedTmpl) {
        cachedTmpl = _.compiler.cacheTemplate(tplId, templStr)
      }
      var compiledComp = _.compiler.cTemplate(templStr, translatedSettings, tplId)
      if (compiledComp) {
        compiledComp = _.compiler.addClassOnRootElement(compiledComp, cachedTmpl.isGrpParent ? 'fb-component fbc-grp-parent' : 'fb-component')
      }
      callback(compiledComp)
    })
  }
  this.translateSettings = function (jsonSettings, tmplId, callback) {
    var _ = this
    if (typeof jsonSettings === 'string') {
      jsonSettings = $.parseJSON(jsonSettings)
    } else {
      jsonSettings = _.deepCopySettings(jsonSettings)
    }
    if (_.options.i18n && _.options.i18n.onDisplay && _.options.i18n.onTranslate && _.options.i18n.isCovered) {
      _.compiler.translateSettings(jsonSettings, tmplId, {
        onDisplay: _.options.i18n.onDisplay,
        onTranslate: _.options.i18n.onTranslate,
        isCovered: _.options.i18n.isCovered,
        done: callback
      })
    } else {
      callback(jsonSettings, tmplId)
    }
  }
  this.deepCopySettings = function (settingsOfComp, settingsToMergeInto) {
    return this.compiler.deepCopySettings(settingsOfComp, settingsToMergeInto)
  }
  this.isEnum = function (o) {
    return this.compiler.isEnum(o)
  }
  this.isHiddenField = function (key) {
    return this.compiler.isHiddenField(key)
  }
  this.htmlRenderer = function (htmlData, jsonSettings, callback) {
    if (htmlData) {
      callback(htmlData)
    }
  }
  this.randomIntBetween = function (min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min)
  }
  this.getCompParentOf = function ($t) {
    var p = $t.nextParentWithClass('fbc-dz')
    if (p.hasClass('fbc-grp')) {
      return p.nextParentWithClass('fbc-grp-parent')
    }
    return p
  }
  this.workspace = null
  this.componentsTab = null
  this.builderTab = null
  this.settingsTab = null
  this.init(jqEl, options)
}

/**
 * Workspace core model:
 *
 * components = { "id": settings:{...}}
 *
 * id = elementId in the workspace
 * settings = specific settings of the element in the workspace
 */
var FT_Workspace = function (fb, jqEl, comps) {
  this.init = function (fb, jqEl, comps) {
    this.fb = fb
    this.el = jqEl
    if (this.el.hasClass('ws-holder')) {
      this.body = this.el
    } else {
      this.body = this.el.find('.ws-holder')
    }
    this.form = this.body.find('form')

    this.testMain = this.el.find('.hcbuild-workspace-test-main')
    this.testBody = this.testMain.find('.workspace-test-body')
    this.testMain.find('.submit-test').click(this.formSubmitTestEvent)
    this.testMain.find('.clear-form-data').click(this.formClearDataEvent)
    this.switches = {
      test: this.el.find('.ws-mode.mode-test'),
      action: this.el.find('.ws-mode.mode-action'),
      connection: this.el.find('.ws-mode.ws-connections input'),
      workspace: this.el.find('.ws-mode.mode-workspace')
    }
    let self = this
    this.switches.test.click(function (e) {
      e.preventDefault()
      self.switches.test.removeClass('active').addClass('active')
      self.switches.action.removeClass('active')
      self.switches.workspace.removeClass('active')
      fb.workspace.switchMode(true)
    })
    this.switches.action.click(function (e) {
      e.preventDefault()
      if (self.switches.action.hasClass('active')) {
        return
      }
      if (self.switches.test.hasClass('active')) {
        fb.workspace.switchMode(false)
      }
      fb.workspace.actionManager.lastStateEnabled = true
      fb.workspace.actionManager.enable(true)
      self.switches.action.removeClass('active').addClass('active')
      self.switches.test.removeClass('active')
      self.switches.workspace.removeClass('active')
    })
    this.switches.workspace.click(function (e) {
      e.preventDefault()
      if (self.switches.test.hasClass('active')) {
        fb.workspace.switchMode(false)
        fb.workspace.actionManager.lastStateEnabled = false
        fb.workspace.actionManager.enable(false)
      }
      if (self.switches.action.hasClass('active')) {
        fb.workspace.actionManager.lastStateEnabled = false
        fb.workspace.actionManager.enable(false)
      }

      self.switches.workspace.removeClass('active').addClass('active')
      self.switches.test.removeClass('active')
      self.switches.action.removeClass('active')
    })
    if (comps) {
      this.components = comps
    }
    var workspaceTab = this
    $.contextMenu({
      selector: '.hcbuild-main .hcbuild-workspace-body .fb-component',
      build: function ($trigger, e) {
        var items
        if ($trigger && $trigger.length) {
          workspaceTab.selectComponent(null, $trigger)
          if ($trigger.hasClass('fbc-grp-parent')) {
            items = {
              'paste': {
                name: 'Paste',
                icon: 'paste'
              },
              'delete': {
                name: 'Delete',
                icon: 'delete'
              }
            }
          } else {
            items = {
              'copy': {
                name: 'Copy',
                icon: 'copy'
              },
              'paste': {
                name: 'Paste',
                icon: 'paste'
              },
              'delete': {
                name: 'Delete',
                icon: 'delete'
              }
            }
            if (workspaceTab.fb.options.userAllowedToEditComponents) {
              items.copyToComp = {
                name: 'Copy to Components',
                icon: 'copy'
              }
            }
          }
          if (!localStorage.getItem('ws-clipboard')) {
            delete items['paste']
          }
        }
        return {
          callback: function (key, options) {
            if (key === 'copy') {
              workspaceTab.onWSCopyComponent()
            } else if (key === 'copyToComp') {
              workspaceTab.onCopyComponent(this, true)
            } else if (key === 'paste') {
              workspaceTab.onWSPasteComponent()
              workspaceTab.connectionManager.update(this)
            } else if (key === 'delete') {
              workspaceTab.onDeleteComponent(this)
              workspaceTab.connectionManager.update()
            }
          },
          items: items
        }
      }
    })
    this.connectionManager.setup(fb)
  }
  this.loadData = function (pFormSrc) {
    if (!pFormSrc) {
      pFormSrc = this.fb.options.data
    } else {
      this.fb.options.data = pFormSrc
    }

    var comps
    if (pFormSrc) {
      comps = this.fb.compiler.getComps(pFormSrc)
    }

    if (comps) {
      this.components = comps
      this.render('', this.components)
    } else {
      this.renderNew()
    }
  }
  this.recompileWorkspace = function (htmlStr, jsonSettings, compByIdCallback, doneCallback) {
    var _ = this.fb
    if (jsonSettings) {
      _.workspace.components = _.compiler.getComps(jsonSettings)
      if (htmlStr) {
        var $html, _compId, c, cm
        $html = $(htmlStr)
        if ($html.length) { // compatibility check
          _.compiler.compsMainLoop(_.workspace.components, function (compId, comp, compMain) {
            compMain = _.compiler.getCompMainObject(comp)
            if (!compMain['_compId']) {
              _compId = $html.find('#' + compId).attr('data-dfsId')
              if (_compId) {
                compMain['_compId'] = _compId
                c = _.workspace.components[compId]
                if (c) {
                  cm = _.compiler.getCompMainObject(c)
                  cm['_compId'] = _compId
                }
              } else {
                // delete because we can't find the corresponding comp reference
                delete _.workspace.components[compId]
              }
            }
            return true
          })
        }
      }
      _.translateSettings(_.workspace.components, null, function (translatedSettings) {
        doneCallback('<form>' + _.compiler.cForm({
          form: _.compiler.getConvertToNewestFormStructure(translatedSettings),
          componentsOnly: true
        }) + '</form>')
      })
    }
  }
  this.updateMetaData = function (compSettings, metaData, updateOrderOfAll) {
    var c = this.fb.compiler;
    var compMain
    if (compSettings && metaData) {
      compMain = c.getCompMainObject(compSettings)
      if (metaData.compId) {
        compMain['_compId'] = metaData.compId
      }
      if (!isNaN(metaData.order)) {
        compMain['_order'] = metaData.order
      }
    }
    if (updateOrderOfAll) {
      var ws = this;
      var $fbComp
      c.compsLoop(this.components, function (compId, comp) {
        compMain = c.getCompMainObject(comp)
        $fbComp = ws.getActiveBody().find('#' + compId)
        ws.updateMetaData(ws.components[compId], {
          order: $fbComp.index()
        }, null)
        return true
      })
    }
  }
  this.render = function (htmlStr, jsonSettings) {
    var workspace = this
    try {
      if (typeof jsonSettings === 'string') {
        jsonSettings = $.parseJSON(jsonSettings)
      } else {
        jsonSettings = this.fb.deepCopySettings(jsonSettings)
      }
    } catch (parseError) {
      jsonSettings = null
    }
    if (jsonSettings && this.fb.compiler.sizeOf(jsonSettings) > 0) {
      // workspace.reInitDragAndDropEvent()
      this.recompileWorkspace(htmlStr, jsonSettings, function (compId) {
        return workspace.fb.componentsTab.getComponentById(compId)
      }, function (compiled) {
        // workspace.compiledWorkspaceData = compiled;
        workspace.fb.htmlRenderer(compiled, jsonSettings, function (renderedWorkspaceData) {
          var $renderedWorkspace = $(renderedWorkspaceData)
          workspace.beforeInserting($renderedWorkspace)
          workspace.body.html($renderedWorkspace)
          workspace.form = workspace.body.children('form:first')
          workspace.form.addClass('fbc-dz fws-main')
          workspace.fb.settingsTab.clearBody()
          workspace.checkFormChilds()
          workspace.form.find('.fb-component').each(function () {
            this.addEventListener('click', function (e) {
              return workspace.selectComponent(e, $(this))
            }, false)
          })

          setTimeout(function () {
            workspace.setStuffWeDontNeedToStore(workspace.form)
          }, 600)
          workspace.reInitDragAndDropEvent()
          workspace.afterInserting(workspace.body)
        })
      })
    } else {
      this.renderNew()
    }
  }
  this.renderNew = function () {
    if (!this.form || this.form.length === 0) {
      var form = this.body.find('form')
      if (!form || form.length === 0) {
        this.body.append('<form class="fbc-dz fws-main"></form>')
        form = this.body.find('form')
      }
      this.form = form
    } else {
      this.form.empty()
    }
    this.form.addClass('fbc-dz fws-main')
    this.checkFormChilds()
    // this.compiledWorkspaceData = $(this.form.toString());
    this.fb.settingsTab.clearBody()
    this.reInitDragAndDropEvent()
    this.actionManager.init(this.body)
  }
  this.checkFormChilds = function () {
    this.form.css('height', '')
    if (this.form.parent().height() > this.form.height()) {
      this.form.css('height', '100%')
    }
  }
  this.insertHtmlStyleHelpers = function ($fbComp) {
    if ($fbComp.children('.fb-ws-only.fbdbg').length === 0) {
      $fbComp.prepend('<i class="fb-ws-only fbdbg"></i>')
    }
    if ($fbComp.children('.fb-ws-only.fbdfg').length === 0) {
      $fbComp.append('<i class="fb-ws-only fbdfg"><i class="left"></i><i class="right"></i></i>')
    }
  }
  this.reInitDragAndDropEvent = function () {
    var workspace = this
    setTimeout(function () {
      workspace.fb.calcContainerHeight(workspace.fb.componentsTab, true)
    }, 600)

    this.connectionManager.update()
    // this.oldDragAndDrop();return;

    if (!this.dragAndDropManager) {
      this.dragAndDropManager = {
        ws: this,
        dragActive: false,
        lastConUpdaterTo: [],
        onCompDropped: function ($comp, $newParent, $oldParent, newComp, index) {
          if (newComp) {
            var compId = this.ws.fb.randomId()
            $comp.attr('id', compId)
            var dfsId = $comp.attr('data-dfsId')
            this.ws.components[compId] = this.ws.fb.deepCopySettings(this.ws.fb.componentsTab.getSettingsById(dfsId))
          }
          $comp.data('lastFbcGrpIndex', this.dataIndex)
          this._setupGrpCompConnection($comp, $newParent, $oldParent, index)
          this.ws.addComponent($comp, newComp, index)
          this.ws.checkFormChilds()
        },
        onCompRemoved: function ($comp, $oldParent) {
          var _ = this
          this.ws.el.find('i[compref="' + $comp.attr('id') + '"]').each(function () {
            var $pointer = $(this)
            var $endpoint = $pointer.data('lastEndpoint')
            if ($endpoint && $endpoint.length) {
              _.ws.actionManager.clearDestinations($pointer, $endpoint)
            }
          })
          this._removeGrpCompImport($comp, $oldParent)
          this.ws.removeComponent($comp, $oldParent)
          this.ws.checkFormChilds()
        },
        _updateSettings: function (compId) {
          var settignsCompId = this.ws.fb.settingsTab.currentCompId()
          try {
            if (settignsCompId === compId) {
              this.ws.fb.settingsTab.updateSettings(settignsCompId, this.ws.components[compId])
            }
          } catch (e) {
            console.log(e)
          }
        },
        _setupGrpCompConnection: function ($comp, $newParent, $oldParent, index) {
          var grpParentCompId, childCompId, dataIndex, childCompMain, parentMain, childImports
          if ($oldParent) {
            this._removeGrpCompImport($comp, $oldParent)
          }
          // add import if it is a group parent
          if ($newParent && $newParent.length && $newParent.hasClass('fbc-grp-parent')) {
            grpParentCompId = $newParent.attr('id')
            dataIndex = $comp.nextParentWithClass('fbc-grp').attr('data-index')
            if (dataIndex !== undefined && dataIndex !== null && dataIndex !== '') {
              childCompId = $comp.attr('id')
              dataIndex = dataIndex + ''
              if (grpParentCompId && childCompId) {
                childCompMain = this.ws.fb.compiler.getCompMainObject(this.ws.components[childCompId])
                childCompMain['_grouped'] = true
                parentMain = this.ws.fb.compiler.getCompMainObject(this.ws.components[grpParentCompId])
                if (parentMain) {
                  if (!parentMain['_import']) {
                    parentMain['_import'] = {}
                  }
                  if (!parentMain['_import'][dataIndex]) {
                    parentMain['_import'][dataIndex] = []
                  }
                  childImports = parentMain['_import'][dataIndex]
                  if (isNaN(index)) {
                    index = $comp.index()
                  }
                  childImports.splice(index, 0, childCompId)
                }
                this._updateSettings(grpParentCompId)
                this._updateSettings(childCompId)
              }
            }
          }
        },
        _removeGrpCompImport: function ($comp, $oldParent) {
          var grpParentCompId, childCompId, dataIndex, childCompMain, parentMain, childImports
          // remove import if it was a group parent
          if ($oldParent && $oldParent.length && $oldParent.attr('id') && $oldParent.hasClass('fbc-grp-parent')) {
            grpParentCompId = $oldParent.attr('id')
            childCompId = $comp.attr('id')
            dataIndex = $comp.data('lastFbcGrpIndex')
            if (dataIndex === undefined || dataIndex === null || dataIndex === '') {
              dataIndex = $comp.nextParentWithClass('fbc-grp').attr('data-index')
            }
            childCompMain = this.ws.fb.compiler.getCompMainObject(this.ws.components[childCompId])
            if (childCompMain) {
              delete childCompMain['_grouped']
            }
            parentMain = this.ws.fb.compiler.getCompMainObject(this.ws.components[grpParentCompId])
            if (parentMain && parentMain['_import'] && parentMain['_import'][dataIndex + '']) {
              childImports = parentMain['_import'][dataIndex + '']
              if (childImports.length) {
                for (var i = 0; i < childImports.length; ++i) {
                  if (childImports[i] === childCompId) {
                    childImports.splice(i, 1)
                    break
                  }
                }
                if (childImports.length === 0) {
                  delete parentMain['_import'][dataIndex + '']
                }
              }
            }
            this._updateSettings(grpParentCompId)
            this._updateSettings(childCompId)
          }
        },
        attachDragEvent: function ($el, isRefComp) {
          if ($el.hasClass('fb-component')) {
            this._attachDragEvent($el, isRefComp)
          } else {
            var $eles, i
            $eles = $el.find('.fb-component')
            for (i = 0; i < $eles.length; ++i) {
              this._attachDragEvent($($eles[i]), isRefComp)
            }
          }
        },
        _attachDragEvent: function ($fbComp, isRefComp) {
          if (isRefComp) {
            $fbComp.data('isRefComp', true)
            this._makeDraggable($fbComp)
            return
          }
          var $eles, i
          $eles = $fbComp.find('.fb-drop-here')
          var $fbDropHere;
          var hasDrops = false
          for (i = 0; i < $eles.length; ++i) {
            $fbDropHere = $($eles[i])
            var $compChildsHolder = $fbDropHere.parent()
            $compChildsHolder.addClass('fbc-dz').addClass('fbc-grp').attr('data-index', $fbDropHere.attr('data-index'))
            $fbDropHere.remove()
            hasDrops = true
          }
          $eles = $fbComp.find('.fb-component')
          for (i = 0; i < $eles.length; ++i) {
            this._attachDragEvent($($eles[i]), isRefComp)
          }
          this._makeDraggable($fbComp)
        },
        _createDragHandle: function ($el, sides) {
          if ($el.children('.fb-ws-only.fbdh').length === 0) {
            var html = ''
            for (var i = 0; i < sides.length; ++i) {
              html += '<i data-side="' + sides[i] + '" class="fb-ws-only fbdh dh' + sides[i] + '"> <i></i><i></i> <i></i><i></i> <i></i><i></i> </i>'
            }
            $el.append(html)
          }
        },
        _makeDraggable: function ($el) {
          if ($el.data('_alreadyDraggable')) {
            return
          }
          $el.data('_alreadyDraggable', true)
          this.ws.insertHtmlStyleHelpers($el)
          this._createDragHandle($el, ['left', 'top', 'right', 'bottom'])

          var _ = this
          // prevent from any runtime comp action during build time
          $el.find('input,button,select,textarea,a').each(function () {
            var $t = $(this)
            $t.unbind()
            this.setAttribute('disabled', 'disabled')
            for (var i = 0; i < this.attributes.length; i++) {
              var attrib = this.attributes[i]
              if (attrib.specified === true) {
                if (/^on\w+/.test(attrib.name)) {
                  this.removeAttribute(attrib.name)
                } else if (/href/i.test(attrib.name)) {
                  this.setAttribute(attrib.name, 'JavaScript:void(0);')
                }
              }
            }
            $t.on('focus input paste change click submit select', function (e) {
              e.preventDefault()
              e.stopPropagation()
              $(this).blur()
              return false
            })
          })
          var scr = _.ws.fb.options.scrollableContainer
          var autoScroll
          if (scr && scr.length) {
            autoScroll = {
              container: scr[0]
            }
          }
          interact($el[0]).styleCursor(false).allowFrom('.fbdh').draggable({
            inertia: false,
            autoScroll: autoScroll,
            onstart: function (e) {
              if (e.interaction.downEvent.button === 2) {
                return false
              }

              _.$dz = _.ws.body.find('.fws-main')
              if (_.$dz && _.$dz.length) {
                _.dragActive = true
              } else {
                _.dragActive = false
                return false
              }
              _.$grpWhenStarted = null
              _.$parentWhenStarted = null
              _.dataIndex = null
              var $containers = _.ws.body.find('.fbc-dz')
              $containers.each(function () {
                // this.removeEventListener("touchmove", _._dragenter, false);
                this.removeEventListener('mouseover', _._dragenter, true)
                this.removeEventListener('mouseout', _._dragleave, true)
                // this.addEventListener("touchmove", _._dragenter, false);
                this.addEventListener('mouseover', _._dragenter, true)
                this.addEventListener('mouseout', _._dragleave, true)
              })
              _.$dragComp = $(e.target)
              _.$dz.addClass('drag-started')
              _._setDim(_.$dragComp, e.interaction.startOffset)
              _.tl.x = 0, _.tl.y = 0
              _.lastX = e.clientX
              _.lastY = e.clientY
              _.mousePos.y = e.clientY
              _.mousePos.x = e.clientX
              var cloneComp = function ($comp, id) {
                var ids = []
                var tid = $comp.attr('id')
                $comp.find('[id]').each(function () {
                  var _id = $(this).attr('id')
                  if (_id) {
                    ids.push(_id)
                  }
                })
                var cpyComp = $comp.toString()
                for (var i = 0; i < ids.length; i++) {
                  cpyComp = cpyComp.replace(new RegExp(ids[i], 'g'), 'cpy_' + ids[i])
                }
                if (tid) {
                  if (id) {
                    cpyComp = cpyComp.replace(new RegExp(tid, 'g'), id)
                  } else {
                    cpyComp = cpyComp.replace(new RegExp(tid, 'g'), 'cpy_' + tid)
                  }
                }
                var $dragCompDropCopy = $(cpyComp)
                $dragCompDropCopy.attr('id', id)
                return $dragCompDropCopy
              }
              if (_.$dragComp.data('isRefComp')) {
                _.isRefComp = true
                _.$dragComp = cloneComp(_.$dragComp, _.ws.fb.randomId())
                _.$dragComp.addClass('fbc-new')
                // _.$dragComp.css({visibility:'hidden'});
                _._getOutOfWorkspaceHolder().append(_.$dragComp)
                _.outOfMovingZone = true
                _.$dropzone = null
              } else {
                _._getOutOfWorkspaceHolder()
                _.outOfMovingZone = false
                _.isRefComp = false
                _.$dropzone = _.$dragComp.nextParentWithClass('fbc-dz')
                _.addDragActiveClass(_.$dropzone)
                _.$grpWhenStarted = _.$dragComp.nextParentWithClass('fbc-grp')
                if (_.$grpWhenStarted && _.$grpWhenStarted.length) {
                  _.dataIndex = _.$grpWhenStarted.attr('data-index')
                }
                _.$parentWhenStarted = _.ws.fb.getCompParentOf(_.$dragComp)
              }
              _.$dragMirror = _._createDragMirror(_.$dragComp)
              _.$dragComp.addClass('dragging')
              _._updateOrgPos()
              _.ws.actionManager.hide()
              _.ws.connectionManager.hide()
              _._isInsideForm()
              return true
            },
            onmove: function (e, a, b, c, d) {
              if (e.interaction.downEvent.button === 2) {
                return false
              }
              _.mousePos.y = e.clientY
              _.mousePos.x = e.clientX
              _.tl.x += _.lastX - e.clientX
              _.tl.y += _.lastY - e.clientY
              _.lastX = e.clientX
              _.lastY = e.clientY

              if (!_._isInsideForm()) {
                _.updateLater()
                return _._checkForLeave(e)
              }
              var $t = $(_.elementFromPoint(e))
              if ($t.hasClass('fb-component')) {
                if (!_.$dragMirror[0].contains($t[0]) && _.tryCloseToComponent($t)) {
                  return true
                }
              } else if ($t.hasClass('fbc-dz')) {
                if (!_.$dragMirror[0].contains($t[0]) && _.tryDropzoneInsert(e, $t)) {
                  return true
                }
              }
              if (_._isInsideDropzone() && _.tryDropzoneInsert(e)) {
                return true
              }
              if (_.tryCloseToComponent(_.$dragMirror.prev()) || _.tryCloseToComponent(_.$dragMirror.next())) {
                return true
              }
              _.updateLater()
              return _._checkForLeave(e)
            },
            onend: function (event) {
              _.dragActive = false
              setTimeout(function () {
                _._endDragging()
                if (_.$dz[0].contains(_.$dragComp[0])) {
                  _.outOfMovingZone = false
                }
                if (_.isRefComp) {
                  if (_.outOfMovingZone) {
                    _.$dragComp.remove()
                  } else {
                    _.$dragComp.removeClass('fbc-new')
                    _.attachDragEvent(_.$dragComp)
                    _.onCompDropped(_.$dragComp, _.ws.fb.getCompParentOf(_.$dragComp), _.$parentWhenStarted, _.isRefComp, _.$dragMirror.index())
                  }
                } else {
                  if (_.outOfMovingZone) {
                    _.onCompRemoved(_.$dragComp, _.$parentWhenStarted)
                    _.$dragComp.remove()
                  } else {
                    _.onCompDropped(_.$dragComp, _.ws.fb.getCompParentOf(_.$dragComp), _.$parentWhenStarted, _.isRefComp, _.$dragMirror.index())
                  }
                }
                _.ws.el.find('.fbc-copy-holder').children().each(function () {
                  var $t = $(this)
                  _.onCompRemoved($t)
                  $t.remove()
                })
                _.ws.body.find('.drag-active').each(function () {
                  _.remDragActiveClass($(this))
                })
              }, 0)

              setTimeout(function () {
                _._endDragging()
                _._removeDragMirror(_.$dragMirror)
                _.ws.actionManager.mayShow(_.ws.body, true)
                _.ws.connectionManager.show(true)
              }, 500)
            }
          })
        },
        _endDragging: function () {
          this.$dragComp.removeClass('dragging')
          this.$dragComp.css('transform', '')
          this.$dz.removeClass('drag-started')
          this.$dropzone = null
        },
        _updatePos: function (target) {
          if (!target) {
            target = this.$dragComp
          }
          this.tl.y = this.mousePos.y - this.dragCompPosition.top - this.dim.handleY
          this.tl.x = this.mousePos.x - this.dragCompPosition.left - this.dim.handleX
          target[0].style.webkitTransform = target[0].style.transform = 'translate(' + this.tl.x + 'px, ' + this.tl.y + 'px)'
        },
        dim: {
          h: 0,
          w: 0,
          h2: 0,
          w2: 0,
          mt: 0,
          ml: 0,
          handleX: 0,
          handleY: 0,
          handlePercentX: 0,
          handlePercentY: 0
        },
        mousePos: {
          x: 0,
          y: 0
        },
        _setDim: function ($t, startPos) {
          this.dim.mt = parseFloat($t.css('marginTop'))
          this.dim.ml = parseFloat($t.css('marginLeft'))
          this.dim.h = $t.outerHeight(true)
          this.dim.w = $t.outerWidth(true)
          this.dim.h2 = this.dim.h / 2
          this.dim.w2 = this.dim.w / 2
          if (startPos) {
            this.dim.handlePercentX = startPos.left + this.dim.ml
            this.dim.handlePercentX = 100 / this.dim.w * this.dim.handlePercentX
            this.dim.handlePercentY = startPos.top + this.dim.mt
            this.dim.handlePercentY = 100 / this.dim.h * this.dim.handlePercentY
          }
          this.dim.handleX = this.dim.handlePercentX * this.dim.w / 100
          this.dim.handleY = this.dim.handlePercentY * this.dim.h / 100
        },
        tl: {
          x: 0,
          y: 0
        },
        dragCompPosition: {
          scrollOffsetY: 0,
          top: 0,
          left: 0
        },
        _updateOrgPos: function (target) {
          if (!target) {
            target = this.$dragComp
          }
          this.bcr = target[0].getBoundingClientRect()
          this.dragCompPosition.top = (this.bcr.top + (-1 * this.tl.y)) - this.dim.mt
          this.dragCompPosition.left = (this.bcr.left + (-1 * this.tl.x)) - this.dim.ml
        },
        currentScrollPos: 0,
        lastScrollPos: 0,
        _getOutOfWorkspaceHolder: function () {
          this.$copyCompHolder = this.ws.el.find('.fbc-copy-holder')
          if (!this.$copyCompHolder.length) {
            this.$copyCompHolder = $('<form class="fbc-copy-holder"></form>')
            this.ws.el.append(this.$copyCompHolder)
          }
          return this.$copyCompHolder
        },
        _createDragMirror: function ($t) {
          var $mirror = $('<div class="drag-mirror"><div class="drag-mirror-bg"></div></div>')
          $t.replaceWith($mirror)
          $mirror.append($t)
          return $mirror
        },
        _removeDragMirror: function ($t) {
          var $fbc = $t.find('.fb-component:first')
          $t.replaceWith($fbc)
        },
        init: function () {
          var _ = this
          var $scr = _.ws.fb.options.scrollableContainer
          if ($scr && $scr.length && !$scr.data('init')) {
            $scr.scroll(function () {
              if (_.dragActive) {
                _.updateLater()
              }
            })
            $scr.data('init', true)
          }
          _.tryCloseToComponent = function ($t) {
            if ($t && $t.length) {
              if ($t[0] === _.$dragMirror.parent()[0]) {
                // target drag is the child of this dropzone
                return false
              }
              var c = $t[0].getBoundingClientRect()
              if (_.mousePos.x > c.left && _.mousePos.x < c.right) {
                var updateDragActiveClass = function ($t) {
                  $t = $t.nextParentWithClass('fbc-dz')
                  if (!_.$dropzone) {
                    _.addDragActiveClass($t)
                    _.$dropzone = $t
                  } else if ($t[0] !== _.$dropzone[0]) {
                    _.remDragActiveClass(_.$dropzone)
                    _.addDragActiveClass($t)
                    _.$dropzone = $t
                  }
                }
                if (_.mousePos.y > c.bottom - ((c.bottom - c.top) * 0.2) && _.mousePos.y < c.bottom) {
                  $t.after(_.$dragMirror)
                  updateDragActiveClass($t)
                  _.outOfMovingZone = false
                  _.updateCoords()
                  return true
                } else if (_.mousePos.y > c.top && _.mousePos.y < c.top + ((c.bottom - c.top) * 0.2)) {
                  $t.before(_.$dragMirror)
                  updateDragActiveClass($t)
                  _.outOfMovingZone = false
                  _.updateCoords()
                  return true
                }
              }
            }
            return false
          }
          _._dragenter = function (e) {
            if (!_.dragActive || _.enteredTarget || !_._isInsideForm()) {
              return true
            }
            var $t = $(_.elementFromPoint(e))
            if ($t.hasClass('fb-component')) {
              if (!_.$dragMirror[0].contains($t[0]) && _.tryCloseToComponent($t)) {
                _.enteredTarget = e.target
                _.changeDZfromComp($t)
                return true
              }
              $t = $t.nextParentWithClass('fbc-dz')
            }
            if (!$t.hasClass('fbc-dz')) {
              return true
            }
            if ($t && $t.length) {
              if (!_.$dragMirror[0].contains($t[0]) && _.tryDropzoneInsert(e, $t)) {
                _.enteredTarget = e.target
                return true
              }
            }
            _.updateCoords()
            return true
          }
          _.elementFromPoint = function (e) {
            return document.elementFromPoint(e.clientX, e.clientY)
          }
          _.updateCoords = function () {
            this._setDim(this.$dragComp)
            this._updateOrgPos()
            this.updateLater()
          }
          _.updateLater = function () {
            var _ = this
            _._updatePos()
            setTimeout(function () {
              // rerender after UI thread of the browser is done
              _._updatePos()
            }, 0)
          }
          _.tryDropzoneInsert = function (e, $dz) {
            var r, offset
            if (!$dz) {
              if (_.$dropzone && _.$dropzone.length) {
                $dz = _.$dropzone
              } else {
                $dz = _.$dz
              }
            } else {
              if (!$dz.hasClass('fbc-dz')) {
                return false
              }
              if (_.$dropzone && _.$dropzone.length) {
                if ($dz[0] === _.$dropzone[0]) {
                  return false
                }
              }
            }
            if (!$dz || $dz.length === 0) {
              return false
            }
            var isSameParent = $dz[0] === _.$dragMirror.parent()[0]
            if (isSameParent) {
              if ($dz[0] === _.$dz[0]) {
                // do not change position on root
                return false
              }
              offset = 0.1
            } else {
              if ($dz[0] === _.$dz[0]) {
                offset = 0.1
              } else {
                offset = 0.5
              }
            }

            if (!r) {
              r = $dz[0].getBoundingClientRect()
            }
            if (_.mousePos.y < r.top + ((r.bottom - r.top) * offset) && _.mousePos.y > r.top) {
              _.outOfMovingZone = false
              if (isSameParent) {
                if (!_.$dragMirror[0].previousSibling) {
                  // already at this position
                  return false
                }
                _.$dragMirror.detach()
              } else {
                _.remDragActiveClass(_.$dropzone)
                _.addDragActiveClass($dz)
              }
              $dz.prepend(_.$dragMirror)
              _.$dropzone = $dz
              _.updateCoords()
              return true
            } else if ((_.mousePos.y > r.bottom - ((r.bottom - r.top) * offset) && _.mousePos.y < r.bottom) || $dz.children().length === 0 || $dz[0] === _.$dz[0]) {
              if (isSameParent) {
                if (!_.$dragMirror[0].nextSibling) {
                  // already at this position
                  return false
                }
                _.$dragMirror.detach()
              } else {
                _.remDragActiveClass(_.$dropzone)
                _.addDragActiveClass($dz)
              }
              _.outOfMovingZone = false
              $dz.append(_.$dragMirror)
              _.$dropzone = $dz
              _.updateCoords()
              return true
            }
            return false
          }

          _.changeDZfromComp = function ($t) {
            $t = $t.nextParentWithClass('fbc-dz')
            if ($t && $t.length) {
              if (this.$dropzone && this.$dropzone.length) {
                if (this.$dropzone[0] !== $t[0]) {
                  this.remDragActiveClass(this.$dropzone)
                }
              }
              this.$dropzone = $t
              this.addDragActiveClass($t)
            }
          }
          _.addDragActiveClass = function ($t) {
            $t.addClass('drag-active')
          }
          _.remDragActiveClass = function ($t) {
            if ($t) {
              $t.removeClass('drag-active')
            }
          }
          _._dragleave = function (e) {
            if (!_.dragActive || !_.enteredTarget) {
              return true
            }
            if (_.enteredTarget === e.target) {
              _.enteredTarget = null
            }
            return true
            if (e.target.contains(_.$dragMirror[0])) {
              _._leave()
              return true
            }
            _._checkForLeave(e)
            return true
          }
        },
        _checkForLeave: function (e) {
          if (this.outOfMovingZone) {
            return false
          }
          this.dzRect = this.$dz[0].getBoundingClientRect()
          if (this._isInsideForm()) {
            return false
          }
          this._leave()
          return true
        },
        _leave: function () {
          if (this.outOfMovingZone) {
            return
          }
          this.outOfMovingZone = true
          if (this.$dropzone) {
            this.remDragActiveClass(this.$dropzone)
          }
          this.$dropzone = null
          // this.remDragActiveClass(this.$dz)
          this.$dragMirror.detach()
          this.$copyCompHolder.append(this.$dragMirror)
          this.ws.checkFormChilds()
          this.updateCoords()
        },
        _isInsideDropzone: function () {
          if (this.$dropzone && this.$dropzone.length) {
            return this._isInsideOf(this.$dropzone[0].getBoundingClientRect())
          }
          return false
        },
        _isInsideForm: function () {
          if (!this.dzRect) {
            this.dzRect = this.$dz[0].getBoundingClientRect()
          }
          return this._isInsideOf(this.dzRect)
        },
        _isInsideOf: function (r) {
          return this.mousePos.y > r.top && this.mousePos.y < r.bottom && this.mousePos.x > r.left && this.mousePos.x < r.right
        },
        dragCompPos: {},
        $dragComp: null,
        $element: null
      }
      this.dragAndDropManager.init()
    }
  }
  this.updateDfsId = function (oldId, newId) {
    var visibleFbComp = this.body.find(".fb-component[data-dfsId='" + oldId + "']")
    visibleFbComp.attr('data-dfsId', newId)
    this.updateMetaData(this.components[visibleFbComp.attr('id')], {
      compId: newId
    })
  }
  this.componentChanged = function (dfsId, component) {
    var _ = this
    _.body.find('.fb-component[data-dfsId="' + dfsId + '"]').each(function (e) {
      _.updateComponentWithSettingsMerge($(this), dfsId, component)
    })
  }
  this.reAssignComponent = function (workspaceCompEl, dfsId) {
    workspaceCompEl.attr('data-dfsId', dfsId)
    this.updateComponentWithSettingsMerge(workspaceCompEl, dfsId, this.fb.componentsTab.getComponentById(dfsId))
  }
  this.updateComponentWithSettingsMerge = function (workspaceCompEl, dfsId, component) {
    var compId = workspaceCompEl.attr('id')
    this.fb.settingsTab.currentCompId({
      release: true
    }) // to ensure the settings on the ui are created fresh
    var newSettings = this.fb.deepCopySettings(component.settings, this.components[compId])
    this.body.find('.fb-component.selected').removeClass('selected')
    this.updateComponent(compId, newSettings)
  }
  this.updateComponent = function (compId, componentSettings, doneCallback) {
    var workspace = this
    var replaceEl = workspace.form.find('#' + compId)
    var dfsId = replaceEl.attr('data-dfsId')
    try {
      var template = workspace.fb.componentsTab.getTemplateById(dfsId)
      if (compId && template && componentSettings) {
        workspace.components[compId] = componentSettings
        workspace.fb.translateSettingsAndCompileTemplate(dfsId, template, componentSettings, function (compiledComp) {
          workspace.fb.htmlRenderer(compiledComp, componentSettings, function (renderedData) {
            workspace.updateMetaData(workspace.components[compId], {
              compId: dfsId,
              order: replaceEl.index()
            }, true)
            workspace.wsSettingsChanged()
            var newEl = $(renderedData)
            newEl.attr('id', compId)
            newEl.attr('data-dfsId', dfsId)
            newEl.addClass('fb-component')
            // newEl.addClass("row");
            newEl.addClass('selected')
            workspace.beforeInserting(newEl, workspace.form)
            var formBuilder = workspace.fb

            newEl[0].addEventListener('click', function (e) {
              return formBuilder.workspace.selectComponent(e, $(this))
            }, false)
            newEl.find('.fb-component').each(function () {
              this.addEventListener('click', function (e) {
                return formBuilder.workspace.selectComponent(e, $(this))
              }, false)
            })

            replaceEl.replaceWith(newEl)
            var p = newEl.nextParentWithClass('fbc-grp')
            if (p && p.length) {
              newEl.data('lastFbcGrpIndex', p.attr('data-index'))
            }
            workspace.afterInserting(newEl, workspace.form)
            if ($.isFunction(doneCallback)) {
              try {
                doneCallback()
              } catch (er) {
                console.log(er)
              }
            }
          })
        })
      } else {
        this.fb.settingsTab.uiNotification({
          status: 'warning',
          message: 'Cannot render without the corresponding component!'
        })
        replaceEl.addClass('no-comp-ref')
      }
    } catch (ee) {
      this.fb.settingsTab.uiNotification({
        status: 'warning',
        message: 'Cannot render without the corresponding component!'
      })
      console.log(ee)
      /* maybe an old element in the workspace, reference to the components lost */
      /* mark it so the user knows the settings couldn't be updated */
      replaceEl.addClass('no-comp-ref')
    }
  }
  this.addComponent = function (jqEl, isNew, index) {
    var compId = jqEl.attr('id');
    var dfsId = jqEl.attr('data-dfsId')
    var p = jqEl.nextParentWithClass('fbc-grp')
    if (p && p.length) {
      jqEl.data('lastFbcGrpIndex', p.attr('data-index'))
    }
    if (isNaN(index)) {
      index = jqEl.index()
    }
    if (!isNew && compId && compId.length > 5) {
      this.updateMetaData(this.components[compId], {
        compId: dfsId,
        order: index
      }, true)
      return
    }
    if (!compId || compId.length < 5) {
      compId = this.fb.randomId()
      jqEl.attr('id', compId)
    }

    var formBuilder = this.fb
    if (!this.components[compId]) {
      this.components[compId] = formBuilder.deepCopySettings(formBuilder.componentsTab.getSettingsById(dfsId))
    }
    jqEl.off('click')
    jqEl[0].addEventListener('click', function (e) {
      return formBuilder.workspace.selectComponent(e, $(this))
    }, false)
    jqEl.find('.fb-component').each(function () {
      $(this).off('click')
      this.addEventListener('click', function (e) {
        return formBuilder.workspace.selectComponent(e, $(this))
      }, false)
    })

    this.updateMetaData(this.components[compId], {
      compId: dfsId,
      order: index
    }, true)
    this.beforeInserting(jqEl, this.form)
    this.afterInserting(jqEl, this.form)
  }
  this.removeComponent = function (jqEl, $parent) {
    var compId = jqEl.attr('id')
    if (compId) {
      var parent
      if ($parent && $parent.length && $parent.attr('id') && $parent.hasClass('fbc-grp-parent')) {
        parent = this.components[$parent.attr('id')]
      }
      this.removeCompAndRefs(compId, this.components[compId], parent)
      this.actionManager.remove(compId)
      delete this.components[compId]
      this.updateMetaData(null, null, true)
    }
    this.fb.settingsTab.componentRemoved(compId)
  }
  this.removeImport = function (importCompId, comp) {
    var compMain = this.fb.compiler.getCompMainObject(comp)
    if (compMain && compMain['_import']) {
      var imports = compMain['_import']
      for (var k in imports) {
        if (imports.hasOwnProperty(k)) {
          var importsArray = imports[k]
          if (importsArray && importsArray.length) {
            for (var i = 0; i < importsArray.length; ++i) {
              if (importsArray[i] === importCompId) {
                importsArray.splice(i, 1)
                break
              }
            }
            if (importsArray.length === 0) {
              delete imports[k]
            }
          }
        }
      }
    }
  }
  this.removeCompAndRefs = function (thisCompId, compToRemove, parent) {
    if (compToRemove) {
      if (parent) {
        this.removeImport(thisCompId, parent)
      }
      var compMain = this.fb.compiler.getCompMainObject(compToRemove)
      if (compMain && compMain['_import']) {
        var imports = compMain['_import']
        for (var k in imports) {
          if (imports.hasOwnProperty(k)) {
            var importsArray = imports[k]
            if (importsArray && importsArray.length) {
              for (var i = 0; i < importsArray.length; ++i) {
                var childComp = this.components[importsArray[i]]
                if (childComp) {
                  var childCompMain = this.fb.compiler.getCompMainObject(childComp)
                  if (childCompMain && childCompMain['_import']) {
                    this.removeCompAndRefs(importsArray[i], this.components[importsArray[i]])
                  }
                  this.actionManager.remove(importsArray[i])
                  delete this.components[importsArray[i]]
                }
              }
            }
          }
        }
      }
    }
  }
  this.selectComponent = function (event, jqEl, withoutFocus) {
    if (!event || jqEl[0] === event.target || jqEl[0] === $(event.target).nextParentWithClass('fb-component')[0]) {
      var compId = jqEl.attr('id')
      this.selected = compId
      this.highlightComponent(compId)
      this.fb.settingsTab.editSettings(compId, this.components[compId], withoutFocus)
    }
    return true
  }
  this.highlightComponent = function (compId) {
    if (compId) {
      this.form.find('.fb-component').removeClass('selected')
      this.form.find('#' + compId).addClass('selected')
      if (this.isInTestMode) {
        this.testBody.find('.fb-component').removeClass('selected')
        this.testBody.find('#' + compId).addClass('selected')
      }
    }
  }
  this.unHighlightComponent = function (compId) {
    if (compId) {
      this.form.find('#' + compId).removeClass('selected')
      if (this.isInTestMode) {
        this.testBody.find('#' + compId).removeClass('selected')
      }
      this.selected = null
    }
  }
  this.onCopyComponent = function (jqEl, copyToComp) {
    if (copyToComp) {
      try {
        var templateCopy = this.fb.componentsTab.getTemplateById(jqEl.attr('data-dfsId'))
        var oldId = jqEl.attr('id')
        var newId = this.fb.randomId()
        var currentSettingsCopy = this.fb.deepCopySettings(this.components[oldId])
        this.fb.componentsTab.addComponent(newId, {
          template: templateCopy,
          settings: currentSettingsCopy
        })
      } catch (tryToFindTheTemplate) {
        this.fb.settingsTab.uiNotification({
          status: 'warning',
          message: 'Cannot render without the corresponding component!'
        })
        jqEl.addClass('no-comp-ref')
      }
    }
  }
  this.onWSCopyComponent = function () {
    if (this.selected) {
      var comp = this.components[this.selected]
      var mainComp = fb.compiler.getCompMainObject(comp)
      if (!mainComp['_import']) {
        localStorage.setItem('ws-clipboard', JSON.stringify(comp))
        return true
      }
    }
    return false
  }
  this.onWSPasteComponent = function () {
    var pastedCompStr = localStorage.getItem('ws-clipboard')
    var success = false
    if (this.selected && pastedCompStr) {
      try {
        var pastedCompJson = JSON.parse(pastedCompStr)
        var compMain = this.fb.compiler.getCompMainObject(pastedCompJson)
        var dfsId = compMain['_compId']
        var template = this.fb.componentsTab.getTemplateById(dfsId)
        if (dfsId && template) {
          var component = {
            template: template,
            settings: pastedCompJson
          }
          var _ = this
          success = true
          this.fb.componentsTab.componentTojqElement(dfsId, component, function (newJqEl) {
            var newCompId = _.fb.randomId()
            _.components[newCompId] = pastedCompJson
            newJqEl.attr('id', newCompId)
            _.beforeInserting(newJqEl, _.form)
            var jqEl
            if (_.selected) {
              jqEl = _.form.find('#' + _.selected)
            }
            if (jqEl && jqEl.length) {
              jqEl.after(newJqEl)
              var p = newJqEl.nextParentWithClass('fbc-grp')
              if (p && p.length) {
                newJqEl.data('lastFbcGrpIndex', p.attr('data-index'))
              }
              _.dragAndDropManager._setupGrpCompConnection(newJqEl, _.fb.getCompParentOf(newJqEl))
            } else {
              _.form.prepend(newJqEl)
            }

            newJqEl[0].addEventListener('click', function (e) {
              return _.selectComponent(e, $(this))
            }, false)
            newJqEl.find('.fb-component').each(function () {
              this.addEventListener('click', function (e) {
                return _.selectComponent(e, $(this))
              }, false)
            })
            _.updateMetaData(pastedCompJson, {
              compId: dfsId
            }, true)
            _.afterInserting(newJqEl, _.form)
            _.actionManager.mayShow(_.body)
          })
        }
      } catch (tryToFindTheTemplate) {
        if (this.selected) {
          this.form.find('#' + this.selected).addClass('no-comp-ref')
        }
      }
    }
    return success
  }
  this.onDeleteComponent = function ($target) {
    var $oldParent = this.fb.getCompParentOf($target)
    this.removeComponent($target.detach(), $oldParent)
    $target.remove()
  }
  this.storeComponents = function () {
    if (this.sizeOf(this.components) === 0) {
      return null
    }
    var formJsonSrc = JSON.stringify(this.components)
    return {
      formSrc: this.fb.compiler.getConvertToNewestFormStructure(JSON.parse(formJsonSrc))
    }
  }
  this.beforeStoring = function (target) {
    target.find('[fb_name]').each(function () {
      var t = $(this)
      t.off()
      t.renameAttr('fb_name', 'name')
    })
  }
  this.beforeInserting = function (newEl, target) {
    if (!newEl || !newEl.length) {
      newEl = target
    }
    newEl.find('[name]').renameAttr('name', 'fb_name')
  }
  this.afterInserting = function (newEl) {
    this.dragAndDropManager.attachDragEvent(newEl)
    this.setStuffWeDontNeedToStore(newEl)
  }
  this.setStuffWeDontNeedToStore = function (target) {
    if (!target || !target.length) {
      return
    }
    try {
      this.actionManager.init(target)
    } catch (e) {
      console.error(e)
    }
  }
  this.removeStuffWeDontNeedToStore = function (target) {
    target.find('.fb-component.selected').removeClass('selected')
    target.find('.fb-component .fb-ws-only').remove()
    target.removeClass('fbc-dz')
    $('.hcbuild-main').find('.fbc-copy-holder').remove()
    target.find('.fbc-dz').removeClass('fbc-dz')
    target.find('.gu-transit').removeClass('gu-transit')
    target.find('.fb-component.context-menu-active').removeClass('context-menu-active')
    // sync json with html
    var compSettingsSize = this.sizeOf(this.components)
    var htmlCompSize = target.find('.fb-component').length
    if (compSettingsSize !== htmlCompSize) {
      var htmlComps = {}
      target.find('.fb-component').each(function () {
        htmlComps[this.id] = true
      })
      var key
      for (key in this.components) {
        if (this.components.hasOwnProperty(key)) {
          if (!htmlComps[key]) {
            delete this.components[key]
          }
        }
      }
    }
    // remove fbonly data
    this.fb.compiler.deepLoopOverJson(this.components, {
      'string': function (value, keyOrIndex, obj) {
        if (/^_fbonly_.*$/.test(keyOrIndex)) {
          delete obj[keyOrIndex]
          return true
        }
        return true
      }
    })
  }
  this.connectionManager = {
    ws: this,
    fb: null,
    $connectionLayer: null,
    enabled: true,
    comConEnabled: false,
    settingsConEnabled: true,
    hidden: false,
    setup: function (fb) {
      this.fb = fb
      this.$connectionLayer = this.fb.el.find('.connectionLayer')
      if (this.enabled) {
        this.show()
      } else {
        this.hide()
      }
    },
    hide: function () {
      if (!this.hidden) {
        this.hidden = true
        this.ws.fb.el.removeClass('cm-active')
      }
    },
    show: function (update) {
      if (this.enabled) {
        if (update) {
          this.update()
        }
        this.ws.fb.el.addClass('cm-active')
        this.hidden = false
        this.cleaned = false
      }
    },
    enable: function (enable) {
      if (enable) {
        this.enabled = true
        if (this.comConEnabled) {
          this.addLocalCss()
        }
        this.update()
      } else {
        this.enabled = false
        if (this.comConEnabled) {
          this.clearLocalCss()
        }
        this.$connectionLayer.find('svg.connections').empty()
      }
    },
    update: function (tfbcomp, reconstructed) {
      if (!this.enabled) {
        return
      }
      this.hidden = false
      var targetBody = this.ws.getActiveBody()
      var _ = this
      _.$connectionLayer.find('svg.connections').empty()
      _.$connectionLayer.find('svg.am-connections').empty()
      if (_.comConEnabled && _.fb.componentsTab.tabActive) {
        var connections = []
        targetBody.find('.fb-component').each(function () {
          var t = $(this)
          var id = t.attr('data-dfsid')
          var c = '.hcbuild-main .fb-workspace .fb-component[data-dfsid="' + id + '"] {' +
            'border-right: 2px solid ' + _.getRandomColor(id) + ';' +
            '}'
          _.addLocalCss(c)
          c = '.hcbuild-main .htmlComponents .fb-component[data-dfsid="' + id + '"]{' +
            'border-left: 2px solid ' + _.getRandomColor(id) + ';' +
            '}'
          _.addLocalCss(c)
          var startId
          t = t.findVisibleInputOrRefElement()
          startId = t.attr('id')
          if (!startId) {
            startId = _.randomId()
            t.attr('id', startId)
          }
          var newPath = {
            offset: _.getRandomLineOffsetFor(id),
            orientation: 'vertical',
            stroke: _.getRandomColor(id),
            strokeWidth: 2,
            start: '#' + startId,
            end: ".hcbuild-comp-body .fb-component[data-dfsid='" + id + "']"
          }
          if (tfbcomp && tfbcomp[0] == t[0]) {
            newPath['animate'] = true
          }
          connections.push(newPath)
        })
        if (connections.length > 0) {
          _.$connectionLayer.find('svg.connections').HTMLSVGconnect({
            paths: connections
          })
        }
      } else if (_.settingsConEnabled && _.fb.settingsTab.tabActive) {
        if (_.comConEnabled) {
          _.clearLocalCss()
        }
        _.fb.settingsTab.connectionsWereActive = true
        return
      }
      _.fb.settingsTab.connectionsWereActive = false
    },
    cleaned: false,
    empty: function () {
      if (!this.cleaned) {
        this.nameConnections = []
        this.nameConExists = {}
        this.cleaned = true
        this.$connectionLayer.find('svg.am-connections').empty()
        this.$connectionLayer.find('svg.connections').empty()
      }
    },
    showAction: function ($arrayItem, uipp, destCompId) {
      if ($arrayItem.length) {
        this.cleaned = false
        this.$connectionLayer.find('svg.am-connections').empty()
        var connections = []
        var endId = $arrayItem.attrID()
        if (uipp) {
          var $uipp = this.ws.el.find(".wsbody [uipp='" + uipp + "']")
          if ($uipp.length) {
            // connection from src input to src settings
            var newPath = {
              animate: '0.4s',
              offset: 24,
              orientation: 'vertical',
              stroke: '#1e00ff',
              strokeWidth: 1,
              start: '#' + $uipp.attrID(),
              end: '.hcbuilder-settings-body #' + endId
            }
            connections.push(newPath)
          }
        }
        if (destCompId) {
          // connection from dest comp to src settings
          newPath = {
            animate: '0.4s',
            offset: 14,
            orientation: 'vertical',
            stroke: '#af36d2',
            strokeWidth: 1,
            start: '#' + destCompId,
            end: '.hcbuilder-settings-body #' + endId
          }
          connections.push(newPath)
        }
        if (connections.length) {
          this.$connectionLayer.find('svg.am-connections').HTMLSVGconnect({
            paths: connections
          })
        }
      }
    },
    hideAction: function ($arrayItem, uipp, destCompId) {

    },
    nameConnections: [],
    nameConExists: {},
    showName: function (connection) {
      if (!this.nameConExists[connection.start + connection.end]) {
        this.nameConnections.push(connection)
        this.nameConExists[connection.start + connection.end] = true
      }
      this.$connectionLayer.find('svg.connections').empty()
      if (this.nameConnections.length > 0) {
        this.cleaned = false
        this.$connectionLayer.find('svg.connections').HTMLSVGconnect({
          paths: this.nameConnections
        })
      }
    },
    getDraglineSvg: function () {
      return this.$connectionLayer.find('svg.dragline')
    },
    localCssMap: {},
    addLocalCss: function (css) {
      if (!this.localCssEl) {
        var c = ''
        c = '.hcbuild-main .fb-workspace .fb-component {' +
          'border-right: 2px solid #c3c3c3;' +
          '}'
        this.localCssMap[c] = 1
        c = '.hcbuild-main .htmlComponents  .fb-component {' +
          'border-left: 2px solid #c3c3c3;' +
          '}'
        this.localCssMap[c] = 1
        c = '.hcbuild-main .am-endpoint-main>i {' +
          'border-right: none !important;' +
          '}'
        this.localCssMap[c] = 1
        this.localCssEl = $('<style type="text/css"></style>')
        this.fb.el.append(this.localCssEl)
      }
      if (css) {
        this.localCssMap[css] = 1
      }
      var allCss = '';
      var key
      for (key in this.localCssMap) {
        if (this.localCssMap.hasOwnProperty(key)) {
          allCss += key
        }
      }
      this.localCssEl.html(allCss)
    },
    clearLocalCss: function (css) {
      if (this.localCssEl && this.localCssEl.length) {
        this.localCssEl.empty()
      }
    },
    colorSet: [
      '#00ffff',
      '#0000ff',
      '#a52a2a',
      '#00ffff',
      '#00008b',
      '#008b8b',
      '#bdb76b',
      '#8b008b',
      '#556b2f',
      '#ff8c00',
      '#8b0000',
      '#e9967a',
      '#9400d3',
      '#ff00ff',
      '#ffd700',
      '#008000',
      '#4b0082',
      '#add8e6',
      '#90ee90',
      '#ffb6c1',
      '#00ff00',
      '#ff00ff',
      '#800000',
      '#000080',
      '#808000',
      '#ffa500',
      '#ffc0cb',
      '#800080'
    ],
    _randomColor: {},
    _colorAlreadySet: {},
    getRandomColor: function (id) {
      if (!this._randomColor[id]) {
        var tryCount = 6;
        var i = 0;
        var color
        for (; i < tryCount; ++i) {
          color = this.colorSet[this.fb.randomIntBetween(0, this.colorSet.length - 1)]
          if (!this._colorAlreadySet[color]) {
            this._colorAlreadySet[color] = 1
            return this._randomColor[id] = color
          }
        }
        return this._randomColor[id] = color
      }
      return this._randomColor[id]
    },
    getRandomLineOffsetForMap: {},
    getRandomLineOffsetForMapNr: {},
    getRandomLineOffsetFor: function (id) {
      if (!this.getRandomLineOffsetForMap[id]) {
        var newInt = this.fb.randomIntBetween(0, 20)
        for (var i = 0; i < 4; ++i) {
          if (!this.getRandomLineOffsetForMapNr[newInt]) {
            this.getRandomLineOffsetForMap[id] = newInt
            this.getRandomLineOffsetForMapNr[newInt] = 1
            return this.getRandomLineOffsetForMap[id]
          }
          newInt = this.fb.randomIntBetween(0, 20)
        }
        this.getRandomLineOffsetForMap[id] = newInt
      }
      return this.getRandomLineOffsetForMap[id]
    }
  } // end of connectionManager
  /**
     "action":{
             "source":[{"_destComp":"", "_index":"", "regex":"regex"}],
             "destination":{
                  "transition":{
                      "all": [
                          "fade"
                      ],
                      "selected": 0
                  }
             }
         }
     **/
  var _ = this.actionManager = {
    ws: this,
    $pointLayer: null,
    $draglineSvg: null,
    lastStateEnabled: false,
    drawedConnections: {},
    transition: {
      'all': [
        'none',
        'slide',
        'fade'

      ],
      'selected': 0
    },
    firstTime: true,
    enable: function (enable) {
      if (this.firstTime) {
        this._setup()
        this.firstTime = false
      }
      if (enable === undefined) {
        enable = this.lastStateEnabled
      }
      var $switch = this.ws.el.find('.ws-mode.ws-action')
      $switch.removeClass('disabled').find('input').prop('disabled', false).prop('checked', enable)
      if (enable) {
        this.mayShow(this.ws.form)
      } else {
        this.disable()
      }
    },
    disable: function () {
      if (this.firstTime) {
        this._setup()
        this.firstTime = false
      }
      var $switch = this.ws.el.find('.ws-mode.ws-action')
      if (this.ws.isInTestMode) {
        $switch.addClass('disabled').find('input').prop('disabled', true)
      }
      $switch.find('input').prop('checked', false)
      this.hide()
    },
    hide: function () {
      if (this.firstTime) {
        this._setup()
        this.firstTime = false
      }
      this.ws.fb.el.removeClass('am-active')
    },
    mayShow: function ($target, forceInit) {
      if (!this.lastStateEnabled && forceInit) {
        this.init($target)
      }
      if (this.lastStateEnabled) {
        this.init($target)
        this.show($target)
      }
    },
    show: function ($target) {
      var _ = this
      setTimeout(function () {
        _.ws.fb.el.addClass('am-active')
      }, 0)
    },
    init: function ($target) {
      if (this.firstTime) {
        this._setup()
        this.firstTime = false
      }
      if ($target) {
        this._drawPointerAndLines($target)
        this._checkDrawedConnections()
      }
    },
    _setup: function () {
      this.$pointLayer = this.ws.fb.el.find('.fb-workspace .am-point-layer')
      this.$draglineSvg = this.ws.connectionManager.getDraglineSvg()
      this._initDropzone()
      this.$pointLayer.off()
      this.$pointLayer.on({
        mouseenter: function () {
          $(this).data('line').attr('class', 'hover')
        },
        mouseleave: function () {
          var t = $(this)
          if (!t.data('dragging')) {
            t.data('line').removeAttr('class')
          }
        }
      }, 'i.am-pointer.connected')

      this.$pointLayer.on({
        mouseenter: function () {
          var lines = $(this).data('lines')
          for (var l = 0; l < lines.length; ++l) {
            $(lines[l]).attr('class', 'hover')
          }
        },
        mouseleave: function () {
          var lines = $(this).data('lines')
          for (var l = 0; l < lines.length; ++l) {
            if (!$(lines[l]).data('dragging')) {
              $(lines[l]).removeAttr('class')
            }
          }
        }
      }, 'o.am-overlay.connected')
    },
    /*
         @return: true=connection accepted, false=rejected
         */
    onActionConnected: function (data) {
      var srcComp = this.ws.components[data.srcCompId]
      var srcCompFieldObj = this.ws.fb.compiler.getCompFieldObjectByName(srcComp, data.name)
      if (!srcCompFieldObj.action) {
        srcCompFieldObj.action = {
          source: [data.source]
        }
      } else {
        if (!$.isArray(srcCompFieldObj.action['source'])) {
          srcCompFieldObj.action['source'] = []
        }
        var sourceAlreadyExists = this._getSourceItemFromComp(srcCompFieldObj, data.source['_index'], data.source['_destCompId'])
        if (sourceAlreadyExists) {
          return false
        } else {
          srcCompFieldObj.action['source'].push(data.source)
        }
      }
      var destComp = this.ws.components[data.source._destCompId]
      var destCompFieldObj = this.ws.fb.compiler.getCompMainObject(destComp)
      if (!destCompFieldObj.action) {
        destCompFieldObj.action = {
          'destination': {
            'transition': this.transition
          }
        }
      } else {
        if (!destCompFieldObj.action['destination'] || !destCompFieldObj.action['destination']['transition']) {
          destCompFieldObj.action['destination'] = {
            'transition': this.transition
          }
        }
      }
      var settignsCompId = this.ws.fb.settingsTab.currentCompId()
      try {
        if (settignsCompId === data.srcCompId) {
          this.ws.fb.settingsTab.updateSettings(settignsCompId, srcComp)
        } else if (settignsCompId === data.source._destCompId) {
          this.ws.fb.settingsTab.updateSettings(settignsCompId, destComp)
        }
      } catch (e) {
        console.log(e)
      }
      return true
    },

    /*
        @return: true=has more destinations, false=was the last
         */
    onActionReleased: function (data) {
      data.destCounter = 0
      var destCompId = data.source._destCompId
      this.ws.fb.compiler.compsLoop(this.ws.components, function (cmpId, cmp, subComp) {
        if (subComp && subComp.action && subComp.action.source) {
          var sources = subComp.action.source
          if (sources.length) {
            for (var ss = 0; ss < sources.length; ++ss) {
              if (sources[ss]['_destCompId'] === destCompId) {
                data.destCounter++
                if (data.destCounter > 1) {
                  return false
                }
                return true
              }
            }
          }
        }
        return true
      })

      if (data.destCounter === 1) {
        var destComp = this.ws.components[data.source._destCompId]
        var destCompFieldObj = this.ws.fb.compiler.getCompMainObject(destComp)
        if (destCompFieldObj['action']) {
          this.ws.fb.compiler.cleanActionDst(destCompFieldObj)
        }
      }
      var srcComp = this.ws.components[data.srcCompId]
      var srcCompFieldObj = this.ws.fb.compiler.getCompFieldObjectByName(srcComp, data.name)
      if (srcCompFieldObj['action']) {
        if (srcCompFieldObj['action']['source']) {
          var sources = srcCompFieldObj['action']['source']
          for (var d = 0; d < sources.length; ++d) {
            if (sources[d]['_index'] === data.source['_index'] && sources[d]['_destCompId'] === data.source['_destCompId']) {
              sources.splice(d, 1)
            }
          }
          this.ws.fb.compiler.cleanActionSrc(srcCompFieldObj)
        }
      }
      var settignsCompId = this.ws.fb.settingsTab.currentCompId()
      if (settignsCompId === data.srcCompId) {
        this.ws.fb.settingsTab.updateSettings(settignsCompId, srcComp)
      } else if (settignsCompId === destCompId) {
        this.ws.fb.settingsTab.updateSettings(settignsCompId, destComp)
      }
      return data.destCounter > 1
    },
    _drawPointerAndLines: function ($target) {
      var _ = this;
      var key
      if ($target.hasClass('fb-component')) {
        this.update($target.nextAll('.fb-component').addBack())
      } else {
        if (this.ws.components) {
          this._markForRemove()
          for (key in _.ws.components) {
            if (_.ws.components.hasOwnProperty(key)) {
              _._drawPointerAndLinesOf(key, _.ws.components[key])
            }
          }
          this._removeMarkedOnes()
        } else {
          this.$pointLayer.empty()
          this.$draglineSvg.empty()
        }
      }
    },
    update: function ($fbComps) {
      if ($fbComps && $fbComps.length) {
        var $target, i2, $childComps, $t
        for (var i = 0; i < $fbComps.length; ++i) {
          $target = $($fbComps[i])
          $childComps = $target.find('.fb-component')
          for (i2 = 0; i2 < $childComps.length; ++i2) {
            $t = $($childComps[i2])
            this._update($t, $t.attr('id'))
          }
          this._update($target, $target.attr('id'))
        }
        this._checkDrawedConnections()
      }
    },
    _update: function ($target, key) {
      this._markForRemove(key)
      var comp = this.ws.components[key]
      this._drawPointerAndLinesOf(key, comp)
      this._removeMarkedOnes()
    },
    remove: function (compId) {
      if (compId) {
        this.ws.fb.compiler.compLoop(this.ws.components[compId], function (mainComp, subComp) {
          if (subComp && subComp.action) {
            delete subComp.action
          }
        })
        this.mayShow(this.ws.body)
      }
    },
    _checkDrawedConnections: function () {
      var compId = this.ws.fb.settingsTab.currentCompId()
      if (compId && this.drawedConnections[compId]) {
        this.ws.fb.settingsTab.updateSettings(compId, this.ws.components[compId])
      }
      this.drawedConnections = {}
    },

    _markForRemove: function (compId) {
      if (compId) {
        var _ = this
        this.$pointLayer.children('i[compRef="' + compId + '"]').each(function () {
          var t = $(this)
          t.data('line').addClass('am-remove')
          t.data('srcOverlay').addClass('am-remove')
          t.addClass('am-remove')
        })
      } else {
        this.$pointLayer.children().each(function () {
          var t = $(this)
          var tagName = this.tagName.toUpperCase()
          if (tagName === 'I') {
            t.data('line').addClass('am-remove')
          }
          t.addClass('am-remove')
        })
      }
    },
    _removeMarkedOnes: function () {
      this.$pointLayer.children('.am-remove').each(function () {
        var t = $(this)
        var tagName = this.tagName.toUpperCase()
        if (tagName === 'I') {
          var $line = t.data('line')
          try {
            var lines = t.data('srcOverlay').data('lines')
            for (var l = 0; l < lines.length; ++l) {
              if (lines[l] === $line[0]) {
                lines.splice(l, 1)
              }
            }
          } catch (e) {}
          $line.remove()
        }
        t.remove()
      })
    },
    _drawPointerAndLinesOf: function (compId, comp) {
      if ($.isArray(comp)) {
        for (var i = 0; i < comp.length; ++i) {
          this._eachCompObj(compId, comp[i], comp[0])
        }
      } else {
        this._eachCompObj(compId, comp, comp)
      }
    },
    _eachCompObj: function (compId, comp, compMain) {
      var fbComp, fields
      fbComp = this.ws.form.find('#' + compId)
      this._createEndpointIfNecessary(fbComp)
      if (comp && comp.name) {
        fields = fbComp.find("[fb_name='" + comp.name + "'],[name='" + comp.name + "']")
        for (var b = 0; b < fields.length; ++b) {
          this._eachField(this.ws.fb, $(fields[b]), fbComp, comp, compMain, b, compId)
        }
      }
    },
    _createEndpointIfNecessary: function (fbComp) {
      var endpoint = fbComp.children('.am-endpoint-main').find('.am-endpoint')
      if (!endpoint.length) {
        endpoint = $('<div class="fb-ws-only am-endpoint-main"><i></i><div class="am-endpoint"></div></div>')
        fbComp.prepend(endpoint)
        endpoint = endpoint.find('.am-endpoint')
      }
      return endpoint
    },
    _eachField: function (_, t, fbComp, comp, compMain, _index, compId) {
      var o = {
        $pointer: null,
        pointSelector: null,
        $srcOverlay: null,
        pointerId: null,
        t: t,
        wasConnectedTo: null,
        sourceItem: null,
        comp: comp,
        _index: _index,
        compId: compId
      }
      var sourceItems;
      var $endpoint;
      var psel
      psel = "i[compRef='" + compId + "'][_index='" + _index + "'][_name='" + comp.name + "']"
      sourceItems = this._getSourceItemsFromComp(comp, _index)
      var hasSrcConnections = false
      // make src connections
      for (var si = 0; si < sourceItems.length; ++si) {
        o.sourceItem = sourceItems[si]
        if (o.sourceItem && o.sourceItem['_destCompId']) {
          var destComp = _.compiler.getCompMainObject(this.ws.components[o.sourceItem['_destCompId']])
          if (destComp && destComp.action && destComp.action.destination) {
            $endpoint = this._createEndpointIfNecessary(this.ws.form.find('#' + o.sourceItem['_destCompId']))
            if ($endpoint && $endpoint.length) {
              o.pointSelector = psel + "[connectedTo='" + o.sourceItem['_destCompId'] + "']"
              this._getOrCreatePointer(o)
              if (o.sourceItem['_fbonly_uipp'] !== o.pointerId) {
                o.sourceItem['_fbonly_uipp'] = o.pointerId
                this.drawedConnections[o.sourceItem['_destCompId']] = 1
              }
              this._drawStartPoint(o.$pointer, o.$srcOverlay)
              this._drawConnected(o.$pointer, $endpoint, true)
              hasSrcConnections = true
            }
          } else {
            this._removeSourceItemFromComp(comp, o.sourceItem)
            this.drawedConnections[compId] = 1
          }
        }
      }
      // make dest visible or cleanup
      if (compMain.action && compMain.action.destination) {
        var hasSourcesConnected = false
        _.compiler.compsLoop(this.ws.components, function (cmpId, cmp, subComp) {
          if (subComp && subComp.action && subComp.action.source) {
            var sources = subComp.action.source
            if (sources.length) {
              for (var ss = 0; ss < sources.length; ++ss) {
                if (sources[ss]['_destCompId'] === compId) {
                  hasSourcesConnected = true
                  return false
                }
              }
            }
          }
          return true
        })
        if (hasSourcesConnected) {
          fbComp.find('.am-endpoint-main').addClass('connected-to')
        } else {
          this.drawedConnections[compId] = 1
          this.ws.fb.compiler.cleanActionDst(compMain)
        }
      } else {
        var m = this
        _.compiler.compsLoop(this.ws.components, function (cmpId, cmp, subComp) {
          if (subComp && subComp.action && subComp.action.source) {
            var sources = subComp.action.source
            if (sources.length) {
              for (var ss = 0; ss < sources.length; ++ss) {
                if (sources[ss]['_destCompId'] === compId) {
                  sources.splice(ss, 1)
                  m.drawedConnections[cmpId] = 1
                  var $connectedPointer = m.$pointLayer.find("i[compRef='" + cmpId + "'][connectedTo='" + compId + "']")
                  if ($connectedPointer.length) {
                    if (!m._removeIfStartPointerExists($connectedPointer)) {
                      m._drawStartPoint($connectedPointer)
                    }
                  }
                  break
                }
              }
              m.ws.fb.compiler.cleanActionSrc(subComp)
            }
          }
          return true
        })
      }

      o.pointSelector = psel + ':not(.connected)'
      if (sourceItems.length === 0) {
        this._getOrCreatePointer(o)
        this._drawStartPoint(o.$pointer, o.$srcOverlay)
      } else {
        t.findVisibleInputOrRefElement()
        if (!t.data('isBool')) {
          // setup not connected pointer for not radios or checkboxes
          this._getOrCreatePointer(o)
          this._drawStartPoint(o.$pointer, o.$srcOverlay)
        }
      }
      if (o.$srcOverlay) {
        if (hasSrcConnections) {
          o.$srcOverlay.addClass('connected')
        } else {
          o.$srcOverlay.removeClass('connected')
        }
      }
    },
    clearDestinations: function ($pointer, $endpoint) {
      var destFbComp = $endpoint.nextParentWithClass('fb-component')
      var ap = {
        name: $pointer.attr('_name'),
        srcCompId: $pointer.attr('compRef'),
        source: {
          '_fbonly_uipp': $pointer.data('srcOverlay').attr('id'),
          '_destCompId': destFbComp.attr('id'),
          '_index': parseInt($pointer.attr('_index'))
        }
      }
      var hasMoreDestinations = this.onActionReleased(ap)
      $pointer.removeClass('connected')
      $pointer.removeAttr('connectedTo')
      if (!hasMoreDestinations) {
        $endpoint.parent().removeClass('connected-to')
      }
    },
    _dragenter: function (e) {
      if (!_.dragActive) {
        return
      }
      _.$dragPointer[0].hidden = true
      let target = document.elementFromPoint(e.clientX, e.clientY)
      _.$dragPointer[0].hidden = false
      let $t = $(target)
      if ($t.hasClass('am-endpoint')) {
        _.$targetEndpoint = $t
        _.intrDz.ondragenter({
          target: _.$targetEndpoint[0],
          relatedTarget: _.$dragPointer[0]
        })
      }
    },
    _dragleave: function (e) {
      if (!_.dragActive) {
        return
      }
      if (_.$targetEndpoint) {
        _.intrDz.ondragleave({
          target: _.$targetEndpoint[0],
          relatedTarget: _.$dragPointer[0]
        })
        _.$targetEndpoint = null
      }
    },
    _dropcheck: function () {
      if (_.$targetEndpoint && _.$targetEndpoint.length && _._isInside(_.$targetEndpoint[0].getBoundingClientRect())) {
        _.intrDz.ondrop({
          target: _.$targetEndpoint[0],
          relatedTarget: _.$dragPointer[0]
        })
      }
    },
    _isInside(r) {
      return r && this.lastY > r.top && this.lastY < r.bottom && this.lastX > r.left && this.lastX < r.right
    },
    _enableDropzoneEntpoints(pointer) {
      var _ = this
      var $endpoints = this.ws.body.find('.fws-main:first')
      $endpoints = $endpoints.find('.am-endpoint')
      $endpoints.each(function () {
        // this.removeEventListener("touchmove", _._dragenter, false);
        this.removeEventListener('mouseover', _._dragenter, false)
        this.removeEventListener('mouseout', _._dragleave, false)

        // this.addEventListener("touchmove", _._dragenter, false);
        this.addEventListener('mouseover', _._dragenter, false)
        this.addEventListener('mouseout', _._dragleave, false)
        _.intrDz.ondropactivate({
          target: this,
          relatedTarget: pointer
        })
      })
    },
    _disableDropzoneEntpoints(pointer) {
      var _ = this
      var $endpoints = this.ws.body.find('.fws-main:first')
      $endpoints = $endpoints.find('.am-endpoint')
      $endpoints.each(function () {
        // this.removeEventListener("touchmove", _._dragenter, false);
        this.removeEventListener('mouseover', _._dragenter, false)
        this.removeEventListener('mouseout', _._dragleave, false)
        _.intrDz.ondropdeactivate({
          target: this,
          relatedTarget: pointer
        })
      })
    },
    _initDropzone: function () {
      var _ = this
      _.currentScrollX = 0
      _.currentScrollY = 0
      _.lastScrollX = 0
      _.lastScrollY = 0
      var $scr = _.ws.fb.options.scrollableContainer
      if ($scr && $scr.length && !$scr.data('am-init')) {
        _.lastScrollX = $scr.scrollLeft()
        _.lastScrollY = $scr.scrollTop()
        $scr.scroll(function () {
          if (_.dragActive) {
            var scroll = $(this)
            var y = scroll.scrollTop()
            var x = scroll.scrollLeft()
            _.currentScrollX += x - _.lastScrollX
            _.currentScrollY += y - _.lastScrollY
            _.lastScrollX = x
            _.lastScrollY = y
          }
        })
        $scr.data('am-init', true)
      }

      _.intrDz = {
        ondropactivate: function (e) {
          var $pointer = $(e.relatedTarget)
          var $endpoint = $(e.target)
          if (_._isConnectionPossible($pointer, $endpoint)) {
            $endpoint.addClass('hcb-changed3')
            $endpoint.parent().addClass('drag-active')
            $endpoint.nextParentWithClass('fbc-dz').addClass('dz-drag-active')
          }
        },
        ondragenter: function (e) {
          var $pointer = $(e.relatedTarget)
          var $endpoint = $(e.target)
          if (_._isConnectionPossible($pointer, $endpoint)) {
            $endpoint.addClass('maybe-connected')
            $endpoint.parent().addClass('maybe-connected-to')
          }
        },
        ondragleave: function (e) {
          var $pointer = $(e.relatedTarget)
          var $endpoint = $(e.target)
          $endpoint.removeClass('maybe-connected')
          $endpoint.parent().removeClass('maybe-connected-to')
          if ($pointer.data('connected')) {
            $pointer.data('connected', false)
            _.clearDestinations($pointer, $endpoint)
          } else {
            if (_._isConnectionPossible($pointer, $endpoint)) {
              $endpoint.addClass('hcb-changed3')
              $endpoint.parent().addClass('drag-active')
              $endpoint.nextParentWithClass('fbc-dz').addClass('dz-drag-active')
            }
          }
        },
        ondrop: function (e) {
          var $pointer = $(e.relatedTarget)
          var $endpoint = $(e.target)
          if (_._isConnectionPossible($pointer, $endpoint)) {
            var destFbComp = $endpoint.nextParentWithClass('fb-component')
            var $src = $pointer.data('valueSrc')
            var conditionVal = $src.val()
            var connectionAccepted = _.onActionConnected({
              name: $pointer.attr('_name'),
              srcCompId: $pointer.attr('compRef'),
              source: {
                '_fbonly_uipp': $pointer.data('srcOverlay').attr('id'),
                '_destCompId': destFbComp.attr('id'),
                '_index': parseInt($pointer.attr('_index')),
                'comment': '',
                'regex': conditionVal
              }
            })
            if (connectionAccepted) {
              $endpoint.removeClass('hcb-changed3')
              _._drawConnected($pointer, $endpoint)
              try {
                if (!$pointer.data('valueSrc').data('isBool')) {
                  // setup not connected pointer for not radios or checkboxes
                  var o = {
                    $pointer: null,
                    pointSelector: null,
                    $srcOverlay: null,
                    pointerId: null,
                    t: null,
                    wasConnectedTo: null,
                    sourceItem: null,
                    comp: null,
                    _index: null,
                    compId: null
                  }
                  o.$srcOverlay = $pointer.data('srcOverlay')
                  o.pointerId = o.$srcOverlay.attr('id')
                  o.t = $pointer.data('valueSrc')
                  o._index = $pointer.attr('_index')
                  o.compId = $pointer.attr('compRef')
                  o.comp = _.ws.fb.compiler.getCompFieldObjectByName(_.ws.components[o.compId], $pointer.attr('_name'))
                  o.pointSelector = "i[compRef='" + o.compId + "'][_index='" + o._index + "'][_name='" + o.comp.name + "']:not(.connected)"
                  _._getOrCreatePointer(o)
                  _._drawStartPoint(o.$pointer, o.$srcOverlay)
                  _.ws.connectionManager.update()
                }
              } catch (she) {
                console.log(she)
              }
            } else {
              $endpoint.removeClass('maybe-connected')
              $endpoint.parent().removeClass('maybe-connected-to')
            }
          }
        },
        ondropdeactivate: function (e) {
          var $endpoint = $(e.target)
          $endpoint.removeClass('hcb-changed3')
          $endpoint.parent().removeClass('drag-active')
          $endpoint.nextParentWithClass('fbc-dz').removeClass('dz-drag-active')
        }
      }
    },
    _getOrCreatePointer: function (o) {
      o.$pointer = this.$pointLayer.children(o.pointSelector)
      o.wasConnectedTo = null
      if (o.$pointer.length) {
        o.$pointer.data('line').removeClass('am-remove')
        o.$pointer.removeClass('am-remove')
        o.$srcOverlay = o.$pointer.data('srcOverlay')
        o.$srcOverlay.removeClass('am-remove')
        o.pointerId = o.$srcOverlay.attr('id')
        o.wasConnectedTo = o.$pointer.attr('connectedTo')
        this._updatePointerRefs(o.$pointer, o.pointerId, o.t)
      } else {
        o.pointerId = this.ws.fb.randomId()
        o.$pointer = $('<i id="' + o.pointerId + '" class="am-pointer"></i>')
        this.$pointLayer.append(o.$pointer)

        if (!o.$srcOverlay) {
          o.pointerId = this.ws.fb.randomId()
          o.$srcOverlay = $('<o id="' + o.pointerId + '" class="am-overlay"></o>')
          // o.$srcOverlay = $('<o oFor="'+o.pointerId+'" class="am-overlay"></o>');
          o.$srcOverlay.data('lines', [])
          this.$pointLayer.append(o.$srcOverlay)
        }
        o.$pointer.data('srcOverlay', o.$srcOverlay)

        o.$pointer.attr('compRef', o.compId)
        o.$pointer.attr('_index', o._index)
        o.$pointer.attr('_name', o.comp.name)

        o.pointerId = o.$srcOverlay.attr('id')
        this._updatePointerRefs(o.$pointer, o.pointerId, o.t)

        var $line = $(document.createElementNS('http://www.w3.org/2000/svg', 'line'))
        $line.attr('stroke-linecap', 'round')
        o.$pointer.data('line', $line)
        this.$draglineSvg.append($line)
        o.$srcOverlay.data('lines').push($line[0])

        var _ = this
        var scr = _.ws.fb.options.scrollableContainer
        var autoScroll
        if (scr && scr.length) {
          autoScroll = {
            container: scr[0]
          }
        }
        interact(o.$pointer[0]).draggable({
          // enable inertial throwing
          inertia: true,
          autoScroll: autoScroll,
          onstart: function (e) {
            _.dragActive = true
            _.currentScrollX = 0
            _.currentScrollY = 0
            _.$dragPointer = $(e.target)
            _.ws.body.find('.fws-main').addClass('drag-started')
            _._enableDropzoneEntpoints(e.target)
            _._drawStartConnection(_.$dragPointer)
            _.lastX = e.clientX
            _.lastY = e.clientY
            _.x = 0
            _.y = 0
            _.$dragPointer.data('srcOverlay').addClass('hcb-grow')
            _.$dragPointer.data('dragging', true)
            _.$dragPointer.data('line').attr('class', 'hover')
            _.$dragPointer.data('line').data('dragging', true)
            _.$targetEndpoint = _.$dragPointer.data('lastEndpoint')
            if (_.$targetEndpoint && _.$targetEndpoint.length) {
              _.wasConnected = true
              _.teRect = _.$targetEndpoint[0].getBoundingClientRect()
            }
          },
          onmove: function (e) {
            _.x += (e.clientX - _.lastX + (_.currentScrollX))
            _.y += (e.clientY - _.lastY + (_.currentScrollY))
            _.lastX = e.clientX
            _.lastY = e.clientY
            _.currentScrollX = 0
            _.currentScrollY = 0

            e.target.style.webkitTransform = e.target.style.transform = 'translate(' + _.x + 'px, ' + _.y + 'px)'

            _._checkInitialConnection(e)

            var $pointer = $(e.target)
            var pos = $pointer.positionOfUnderlying(_.$pointLayer)
            $pointer.data('line').attr('x2', pos.x).attr('y2', pos.y)
            var $endpoint = $pointer.data('lastEndpoint')
            if (!$pointer.data('connected') && $endpoint) {
              $endpoint.addClass('hcb-changed3')
              $endpoint.parent().addClass('drag-active')
              $endpoint.nextParentWithClass('fbc-dz').addClass('dz-drag-active')
              $pointer.data('lastEndpoint', null)
            }
          },
          onend: function (e) {
            _._disableDropzoneEntpoints(e.target)
            _._checkInitialConnection(e)
            _._dropcheck(e)
            var $pointer = $(e.target)
            $pointer.data('srcOverlay').removeClass('hcb-grow')
            $pointer.data('dragging', false)
            try {
              $pointer.data('line').data('dragging', false)
              $pointer.data('line').removeAttr('class')
            } catch (e) {}
            if (!$pointer.data('connected')) {
              // reset
              if (!_._removeIfStartPointerExists($pointer)) {
                e.target.style.webkitTransform = e.target.style.transform = 'translate(0, 0)'
                var pos = $pointer.data('src').positionOfUnderlying(_.$pointLayer)
                $pointer.css({
                  top: pos.y + 'px',
                  left: pos.x + 'px'
                })
                $pointer.data('line').attr('visibility', 'hidden')
                $pointer.data('lastEndpoint', null)
              }
            }
            _.ws.body.find('.fws-main').removeClass('drag-started')
          }
        })
      }
    },
    _checkInitialConnection: function (e) {
      if (this.wasConnected && !this._isInside(this.teRect)) {
        this._dragleave(e)
        if (!this.$targetEndpoint) {
          this.wasConnected = false
          this.teRect = null
        }
      }
    },
    _removeIfStartPointerExists: function ($pointer) {
      var $otherStartPointer = this.$pointLayer.find("i[compRef='" + $pointer.attr('compRef') + "'][_index='" + $pointer.attr('_index') + "'][_name='" + $pointer.attr('_name') + "']:not(.connected)")
      if ($otherStartPointer.length > 1) {
        $pointer.removeClass('connected')
        $pointer.removeAttr('connectedTo')
        if (this._otherPointersConnected($pointer)) {
          $pointer.data('srcOverlay').addClass('connected')
        } else {
          $pointer.data('srcOverlay').removeClass('connected')
        }
        var $line = $pointer.data('line')
        try {
          var lines = $pointer.data('srcOverlay').data('lines')
          for (var l = 0; l < lines.length; ++l) {
            if (lines[l] === $line[0]) {
              lines.splice(l, 1)
            }
          }
        } catch (e) {}
        $line.remove()
        $pointer.remove()
        return true
      }
      $pointer.data('srcOverlay').removeClass('connected')
      return false
    },
    _otherPointersConnected: function ($pointer) {
      var $otherPointersConnected = this.$pointLayer.find("i.connected[compRef='" + $pointer.attr('compRef') + "'][_index='" + $pointer.attr('_index') + "'][_name='" + $pointer.attr('_name') + "']")
      return $otherPointersConnected.length > 0
    },
    _updatePointerRefs: function ($pointer, pId, $t) {
      var $visibleEl = $t.findVisibleInputOrRefElement()
      $visibleEl.attr('uipp', pId)
      $pointer.data('src', $visibleEl)
      $pointer.data('valueSrc', $t)
      try { // for test mode
        var mirrorT = this.ws.testBody.find('#' + $t.attr('id'))
        if (mirrorT.length) {
          mirrorT.findVisibleInputOrRefElement().attr('uipp', pId)
        }
      } catch (eee) {
        console.log(eee)
      }
    },
    _getSourceItemsFromComp: function (comp, index) {
      var sourceItem = []
      if (comp.action && comp.action.source && comp.action.source.length) {
        for (var d = 0; d < comp.action.source.length; ++d) {
          if (comp.action.source[d]['_index'] === index) {
            sourceItem.push(comp.action.source[d])
          }
        }
      }
      return sourceItem
    },
    _getSourceItemFromComp: function (comp, index, compId) {
      if (comp.action && comp.action.source && comp.action.source.length) {
        for (var d = 0; d < comp.action.source.length; ++d) {
          if (comp.action.source[d]['_index'] === index && comp.action.source[d]['_destCompId'] === compId) {
            return comp.action.source[d]
          }
        }
      }
      return null
    },
    _removeSourceItemFromComp: function (comp, srcItem) {
      if (comp.action && comp.action.source && comp.action.source.length) {
        for (var d = 0; d < comp.action.source.length; ++d) {
          if (comp.action.source[d]['_index'] === srcItem['_index'] && comp.action.source[d]['_destCompId'] === srcItem['_destCompId']) {
            comp.action.source.splice(d, 1)
            break
          }
        }
        this.ws.fb.compiler.cleanActionSrc(comp)
      }
    },
    _drawStartPoint: function ($pointer, $srcOverlay) {
      if (!$srcOverlay) {
        $srcOverlay = $pointer.data('srcOverlay')
      }
      $pointer.removeClass('connected')
      var $visibleEl = $pointer.data('src')
      var startPos = $visibleEl.positionOfUnderlying(this.$pointLayer)
      var h2, w2, h, w
      var margin = 6
      h = $visibleEl.height() + margin
      w = $visibleEl.width() + margin
      h2 = h / 2
      w2 = w / 2

      $srcOverlay.css({
        top: (startPos.y - h2) + 'px',
        left: (startPos.x - w2) + 'px',
        height: h + 'px',
        width: w + 'px'
      })
      $pointer.css({
        top: startPos.y + 'px',
        left: startPos.x + 'px'
      })

      $pointer.data('line').attr('visibility', 'hidden')
      $pointer.data('connected', false)
    },
    _drawConnected: function ($pointer, $endpoint, start) {
      $endpoint.removeClass('maybe-connected')
      $endpoint.parent().removeClass('maybe-connected-to')
      $endpoint.addClass('connected')
      $endpoint.parent().addClass('connected-to')
      if (start) {
        this._drawStartConnection($pointer, start)
      } else {
        $pointer.data('srcOverlay').addClass('connected')
      }
      var _ = this
      // ensure the render thread is updating the layout first
      setTimeout(function () {
        var pos = $endpoint.positionOfUnderlying(_.$pointLayer)
        $pointer.css({
          top: parseInt(pos.y - $pointer[0].offsetHeight / 2) + 'px',
          left: parseInt(pos.x - $pointer[0].offsetWidth / 2) + 'px'
        })
        pos = $pointer.positionOfUnderlying(_.$pointLayer)
        $pointer.data('line').attr('x2', parseInt(pos.x)).attr('y2', parseInt(pos.y))
      }, 0)
      $pointer[0].style.webkitTransform = $pointer[0].style.transform = 'translate(0, 0)'
      $pointer.data('lastEndpoint', $endpoint)
      $pointer.attr('connectedTo', $endpoint.parent().parent().attr('id'))
      $pointer.data('connected', true)
      $pointer.addClass('connected')
    },
    _drawStartConnection: function ($pointer, start) {
      $pointer.data('srcOverlay').addClass('connected')
      var pos = $pointer.data('src').positionOfUnderlying(this.$pointLayer)
      var updateLine = function () {
        $pointer.data('line').attr('x1', pos.x + $pointer[0].offsetWidth / 2).attr('y1', pos.y + $pointer[0].offsetHeight / 2).attr('visibility', 'visible')
      }
      if ($pointer.width() === 0) {
        // the zero second pause allows the browser to process the rerender event
        setTimeout(updateLine, 0)
      } else {
        updateLine()
        $pointer.data('line').attr('x2', pos.x + 1).attr('y2', pos.y + 1).attr('visibility', 'visible')
      }
    },
    _isConnectionPossible: function ($pointer, $endpoint) {
      // not allowed to connect with itself
      var $endpointComp = $endpoint.parent().parent()
      var destCompId = $endpointComp.attr('id')
      if (destCompId === $pointer.attr('compRef')) {
        return false
      }
      // not allowed to connect with his parent
      var $p = this.ws.fb.getCompParentOf($pointer.data('src'))
      if ($endpointComp[0] === $p[0]) {
        return false
      }
      // not allowed to connect again with this endpoint
      var _index, srcCompId
      _index = parseInt($pointer.attr('_index'))
      srcCompId = $pointer.attr('compRef')
      var srcComp = this.ws.components[srcCompId]
      var srcCompFieldObj = this.ws.fb.compiler.getCompFieldObjectByName(srcComp, $pointer.attr('_name'))
      var sourceAlreadyExists = this._getSourceItemFromComp(srcCompFieldObj, _index, destCompId)
      if (sourceAlreadyExists) {
        return false
      }
      return true
    }
  } // end of actionManager
  this.sizeOf = function (obj) {
    var size = 0;
    var key
    for (key in obj) {
      if (obj.hasOwnProperty(key)) size++
    }
    return size
  }
  this.getActiveBody = function () {
    if (this.isInTestMode) {
      return this.testBody
    } else {
      return this.body
    }
  }
  this.switchMode = function (testMode) {
    if (testMode) {
      this.fb.el.addClass('fb-test-mode-active')
      this.isInTestMode = true
      this.actionManager.disable()
      this.wsSettingsChanged()
    } else {
      this.fb.el.removeClass('fb-test-mode-active')
      this.isInTestMode = false
      this.testMain.hide()
      this.cleanupTestMode()
      this.body.html(this.form)
      this.formDetached = false
      this.body.parent().show()
      this.connectionManager.update()
      this.actionManager.enable()
    }
  }
  this.cleanupTestMode = function () {
    if (this.testBody && this.testBody.length) {
      this.testBody.find('[name]').each(function () {
        var t = $(this)
        t.off()
        t.remove()
      })
      this.testBody.children('form').remove()
      this.testBody.empty()
    }
  }
  this.wsSettingsChanged = function () {
    if (this.isInTestMode) {
      var _ = this
      this.testBody.removeClass('test-succeeded').removeClass('test-failed')
      var settingsCopy = _.fb.deepCopySettings(_.components)
      _.fb.translateSettings(settingsCopy, null, function (translatedSettings) {
        if (!_.formDetached) {
          _.body.parent().hide()
          _.form = _.body.find('form:first')
          _.form.detach()
          _.formDetached = true
        }
        _.cleanupTestMode()
        try {
          var htmlForm = _.fb.compiler.cForm({
            form: _.fb.compiler.getConvertToNewestFormStructure(translatedSettings)
          })
          var $form = $(htmlForm)
          _.testBody.html($form)
          $form.on('action-animation-started', function () {
            _.connectionManager.hide()
          })
          $form.on('action-animation-ended', function () {
            _.connectionManager.show(true)
          })
          $form.find('.fb-component').each(function () {
            var $t = $(this)
            _.insertHtmlStyleHelpers($t)
            $t.click(function (e) {
              return _.selectComponent(e, $(this), true)
            })
          })
          _.highlightComponent(_.fb.settingsTab.currentCompId())

          if (_.fb.options.test && $.isFunction(_.fb.options.test.onActive)) {
            _.fb.options.test.onActive($form, _.fb.compiler.getConvertToNewestFormStructure(_.components))
          }
          _.testMain.show()
        } catch (dontCare) {
          console.log(dontCare)
        }
        _.connectionManager.update()
      })
    }
  }
  this.formSubmitTestEvent = function () {
    try {
      if (fb.options.test && fb.options.test.onSubmit && $.isFunction(fb.options.test.onSubmit)) {
        var theForm = fb.workspace.testBody.children('form')
        fb.options.test.onSubmit(theForm, fb.workspace.components, function (success) {
          if (fb.workspace.lastSubmitTimeoutId) {
            clearTimeout(fb.workspace.lastSubmitTimeoutId)
          }
          fb.workspace.testBody.removeClass('test-succeeded').removeClass('test-failed')
          var removeClassLater = function () {
            fb.workspace.lastSubmitTimeoutId = setTimeout(function () {
              fb.workspace.testBody.removeClass('test-succeeded').removeClass('test-failed')
            }, 2000)
          }
          setTimeout(function () {
            if (success === true) {
              fb.workspace.testBody.addClass('test-succeeded')
              removeClassLater()
            } else if (success === false) {
              fb.workspace.testBody.addClass('test-failed')
              removeClassLater()
            }
          }, 0)
        })
      }
    } catch (e) {
      console.log(e)
    }
  }
  this.formClearDataEvent = function () {
    try {
      fb.options.test.onReset()
    } catch (e) {
      console.log(e)
    }
    fb.workspace.wsSettingsChanged()
  }
  this.cleanElement = function (targetEl) {
    targetEl.removeClass('gu-transit')
    targetEl.removeClass('selected')
    targetEl.find('.fb-ws-only').remove()
    targetEl.removeClass('fbc-dz')
    $('.hcbuild-main').find('.fbc-copy-holder').remove()
    targetEl.find('.fbc-dz').removeClass('fbc-dz')
    targetEl.removeClass('context-menu-active')
  }
  this.clearWorkspace = function () {
    if (this.isInTestMode) {
      this.switchMode(false)
      this.switches.test.prop('checked', false)
    }
    delete this.components
    // try{this.compiledWorkspaceData.remove();this.compiledWorkspaceData = null;}catch(dontCare){this.compiledWorkspaceData = null;}
    this.body.empty()
    this.form = null
  }
  this.fb = null
  this.el = null
  this.components = {} /** {"compId":{..settings..}**/
  this.init(fb, jqEl, comps)
}

/**
 * ComponentsTab core model:
 *
 * components = {"someDfsId":{template:"mustache js template...", settings:{...}}}
 *
 * template: handlebars template string you have defined in builderTab
 * settings: settings you have defined in the builderTab
 */
var FT_ComponentsTab = function (fb, jqEl) {
  this.init = function (fb, el) {
    this.tabActive = true
    this.fb = fb
    this.el = el
    this.searchField = this.el.find('.fb-comp-search-field')
    if (this.el.hasClass('hcbuild-comp-body')) {
      this.body = this.el
    } else {
      this.body = this.el.find('.hcbuild-comp-body')
    }
    if (this.fb.options.componentsTab.components) {
      this.components = this.fb.options.componentsTab.components
      this.render()
    }
    this.defaultDescriptor = {
      // name: "",
      // validate: {required: true}
    }
    var componentsTab = this
    if (componentsTab.fb.options.userAllowedToEditComponents) {
      $.contextMenu({
        selector: '.hcbuild-main .htmlComponents .fb-component',
        callback: function (key, options) {
          if (key === 'edit') {
            componentsTab.onEditComponent(this)
          } else if (key === 'copy') {
            componentsTab.onCopyComponent(this)
          } else if (key === 'paste') {
            componentsTab.onPasteComponent(this)
          } else if (key === 'delete') {
            componentsTab.onDeleteComponent(this)
          }
        },
        items: {
          'edit': {
            name: 'Edit',
            icon: 'edit'
          },
          // "cut": {name: "Cut", icon: "cut"},
          'copy': {
            name: 'Copy',
            icon: 'copy'
          },
          'paste': {
            name: 'Paste',
            icon: 'paste'
          },
          'delete': {
            name: 'Delete',
            icon: 'delete'
          }
        }
      })
    }
    this.saveBtnEl = this.el.find('.save-btn')
    if (componentsTab.fb.options.userAllowedToEditComponents !== true) {
      this.saveBtnEl.hide()
    }
    this.saveBtnEl.click(function () {
      componentsTab.storeComponents()
    })
    this.body.dblclick(function (event) {
      if (componentsTab.fb.options.userAllowedToEditComponents) {
        var fbComponent = $(event.target)
        var maxIndex = 80
        var index = 1
        while (!fbComponent.hasClass('fb-component')) {
          fbComponent = fbComponent.parent()
            ++index
          if (maxIndex < index) {
            break
          }
        }
        if (fbComponent.hasClass('fb-component')) {
          componentsTab.onEditComponent(fbComponent)
        }
      }
    })
    this.searchField.keydown(function (event) {
      if (event.keyCode === 13) {
        event.preventDefault()
        return false
      }
    })
    this.searchField.keyup(function (event) {
      if (event.keyCode === 13) {
        event.preventDefault()
        return false
      }
      componentsTab.searchComponentsAndRender($(this).val())
    })
  }
  this.searchComponentsAndRender = function (text) {
    var _ = this
    _.fb.options.component.searchComp(text, function (data) {
      if (data) {
        try {
          var searchedComponents = {}
          var compsCount = 0
          for (var key in data) {
            if (data.hasOwnProperty(key)) {
              if (typeof data[key].settings === 'string') {
                data[key].settings = JSON.parse(data[key].settings)
              }
              data[key].settings = _.fb.deepCopySettings(data[key].settings, _.getDefaultSettingsDescriptor())
              searchedComponents[key] = {
                  template: data[key].template,
                  settings: data[key].settings
                }
                ++compsCount
            }
          }
          if (compsCount > 0) {
            _.components = searchedComponents
            _.render()
            return
          }
        } catch (errorOnSearch) {
          console.log(errorOnSearch)
        }
        _.components = {}
        _.render()
      } else {
        _.components = {}
        _.render()
      }
    })
  }
  this.getDefaultSettingsDescriptor = function () {
    // return this.fb.deepCopySettings(this.defaultDescriptor);
    return null
  }
  this.loadData = function () {
    var componentsTab = this
    this.fb.options.component.requestComp(null, function (data) {
      if (data) {
        var count = 0
        if (!componentsTab.components) {
          componentsTab.components = {}
        }
        for (var key in data) {
          if (data.hasOwnProperty(key)) {
            if (data[key].settings && data[key].template) {
              if (typeof data[key].settings === 'string') {
                data[key].settings = JSON.parse(data[key].settings)
              }
              data[key].settings = componentsTab.fb.deepCopySettings(data[key].settings, componentsTab.getDefaultSettingsDescriptor())
              // if(data[key].settings.action && data[key].settings.action.startsWith("alert")){
              //   data[key].settings.action = "";
              // }
              componentsTab.components[key] = {
                  template: data[key].template,
                  settings: data[key].settings
                }
                ++count
            }
          }
        }

        if (count > 0) {
          componentsTab.render()
        } else {
          componentsTab.loadDefaultData()
        }
      } else {
        componentsTab.loadDefaultData()
      }
    })
  }
  this.loadDefaultData = function () {
    this.components = fbDefaultComponents
    for (var key in this.components) {
      if (this.components.hasOwnProperty(key)) {
        this.components[key].settings = this.fb.deepCopySettings(this.components[key].settings, this.getDefaultSettingsDescriptor())
        this.storeQueue[key] = true
      }
    }
    this.componentsChanged()
    this.render()
  }
  this.fb = null
  this.el = null
  this.templates = null
  this.copiedComponentDfsId = null
  this.components = {} // {"someDfsId":{template:"", settings:{}}};
  // same as components but this one is not being emptied when searching for components
  // it is needed to prevent from compile issues when re-rending components on the workspace
  this.backupComponents = {}
  this.storeQueue = {}
  this.getComponentById = function (dfsId) {
    return this.components[dfsId]
  }
  this.updateDfsId = function (oldId, newId) {
    this.components[newId] = this.components[oldId]
    delete this.components[oldId]
    if (this.copiedComponentDfsId === oldId) {
      this.copiedComponentDfsId = newId
    }
    this.body.find(".fb-component[data-dfsId='" + oldId + "']").attr('data-dfsId', newId)
    this.fb.workspace.updateDfsId(oldId, newId)
    this.fb.builderTab.updateDfsId(oldId, newId)
  }
  this.storeComponents = function () {
    var componentsTab = this
    for (var key in this.storeQueue) {
      if (this.storeQueue.hasOwnProperty(key)) {
        if (this.components[key]) {
          var updateComp = this.components[key]
          var newComp = {
            id: key,
            template: updateComp.template,
            settings: updateComp.settings
          }
          this.fb.options.component.storeComp(newComp, function (data) {
            if (data && data.oldId && data.newId) {
              componentsTab.updateDfsId(data.oldId, data.newId)
            }
          })
        }
      }
    }
    this.storeQueue = {}
    this.componentsSaved()
  }
  this.getSettingsById = function (dfsId) {
    if (this.components[dfsId]) {
      return this.components[dfsId].settings
    }
    return this.backupComponents[dfsId].settings
  }
  this.getTemplateById = function (dfsId) {
    if (this.components[dfsId]) {
      return this.components[dfsId].template
    }
    return this.backupComponents[dfsId].template
  }
  this.addComponent = function (dfsId, component, jqEl) {
    var _ = this
    var isNew = false
    if (_.components[dfsId]) {
      isNew = false
    } else {
      isNew = true
    }
    _.components[dfsId] = component
    _.backupComponents[dfsId] = component
    _.componentTojqElement(dfsId, component, function (newJqEl) {
      if (isNew) {
        if (jqEl && jqEl.length) {
          jqEl.after(newJqEl)
        } else {
          _.body.append(newJqEl)
        }
      } else {
        _.body.find(".fb-component[data-dfsId='" + dfsId + "']").replaceWith(newJqEl)
      }
      _.fb.workspace.dragAndDropManager.attachDragEvent(newJqEl, true)
      _.storeQueue[dfsId] = true
      _.componentsChanged()
    })
    if (!isNew) {
      _.fb.workspace.componentChanged(dfsId, component)
    }
  }
  this.componentTojqElement = function (dfsId, component, callback) {
    var _ = this.fb
    _.translateSettingsAndCompileTemplate(dfsId, component.template, component.settings, function (compiledTmpl) {
      _.htmlRenderer(compiledTmpl, component.settings, function (renderedHtml) {
        var newJqEl = $(renderedHtml).attr('data-dfsId', dfsId)
        newJqEl.addClass('fb-component')
        // newJqEl.addClass("row");
        callback(newJqEl)
      })
    })
  }
  this.onEditComponent = function (jqEl) {
    var dfsId = jqEl.attr('data-dfsId')
    this.fb.builderTab.editComponent(dfsId, this.components[dfsId])
  }
  this.onCopyComponent = function (jqEl) {
    var dfsId = jqEl.attr('data-dfsId')
    this.copiedComponentDfsId = dfsId
  }
  this.onPasteComponent = function (jqEl) {
    if (this.copiedComponentDfsId) {
      this.addComponent(this.fb.randomId(), this.fb.deepCopySettings(this.components[this.copiedComponentDfsId]), jqEl)
    }
  }
  this.onDeleteComponent = function (jqEl) {
    var dfsId = jqEl.attr('data-dfsId')
    // delete remote
    this.fb.options.component.deleteComp(dfsId, function () {
      // deleted
    })
    // delete local
    this.body.find(".fb-component[data-dfsId='" + dfsId + "']").remove()
    delete this.components[dfsId]
  }
  this.render = function () {
    var _ = this
    var size = _.sizeOf(_.components)
    var deliveredCollector = []
    if (size === 0) {
      _.body.empty()
    } else {
      if (!_.fb.workspace.dragAndDropManager) {
        _.fb.workspace.reInitDragAndDropEvent()
      }
      for (var key in _.components) {
        if (_.components.hasOwnProperty(key)) {
          _.backupComponents[key] = _.components[key]
          _.componentTojqElement(key, _.components[key], function (jqElement) {
            deliveredCollector.push(jqElement)
            if (size == deliveredCollector.length) {
              deliveredCollector.sort(function (a, b) { // sort by id ignoring letters in it:
                var aId = $(a).attr('data-dfsid') + ''
                var bId = $(b).attr('data-dfsid') + ''
                return (parseInt(bId.substring(2)) < parseInt(aId.substring(2))) ? 1 : -1
              })
              _.body.empty()
              // append
              var $refComp
              for (var i = 0; i < deliveredCollector.length; ++i) {
                $refComp = $(deliveredCollector[i])
                _.body.append($refComp)
                try {
                  _.fb.workspace.dragAndDropManager.attachDragEvent($refComp, true)
                } catch (eee) {
                  console.log(eee)
                }
              }
              setTimeout(function () {
                _.fb.workspace.connectionManager.update()
              }, 2000)
            }
          })
        }
      }
    }
  }
  this.sizeOf = function (obj) {
    var size = 0;
    var key
    for (key in obj) {
      if (obj.hasOwnProperty(key)) size++
    }
    return size
  }
  this.showTab = function () {
    this.fb.el.find('.hcbuilder a.' + this.el.attr('id')).tab('show')
  }
  this.tabShown = function () {
    this.tabActive = true
    this.fb.activeTabObj = this
    fb.calcContainerHeight(this, true)
    if (this.fb.options.readOnly !== true) {
      this.searchField.focus()
    }
    this.fb.workspace.connectionManager.update()
  }
  this.tabHidden = function () {
    this.tabActive = false
    this.fb.workspace.connectionManager.update()
  }
  this.componentsChanged = function () {
    this.saveBtnEl.addClass('btn-warning')
    this.saveBtnEl.removeClass('btn-secondary')
    this.saveBtnEl.addClass('hcb-changed')
  }
  this.componentsSaved = function () {
    this.saveBtnEl.addClass('btn-secondary')
    this.saveBtnEl.removeClass('btn-warning')
    this.saveBtnEl.removeClass('hcb-changed')
  }
  this.init(fb, jqEl)
}

/**
 * SettingsTab
 *
 * This Class is responsible for the modification of the component settings
 */
var FT_SettingsTab = function (fb, jqEl) {
  this.init = function (fb, jqEl) {
    this.fb = fb
    this.el = jqEl
    this.autoSave = !!fb.options.autoSaveSettings
    this.el.find('.save-btn').click(function () {
      fb.settingsTab.storeSettings()
    })
    this.autoSaveBtn = this.el.find('.auto-save-btn')
    this.autoSaveBtn.prop('checked', this.autoSave)
    this.autoSaveBtn.click(function () {
      fb.settingsTab.autoSave = $(this).is(':checked')
    })

    if (this.el.hasClass('hcbuilder-settings-body')) {
      this.body = this.el
    } else {
      this.body = this.el.find('.hcbuilder-settings-body')
    }
  }
  this.fb = null
  this.el = null
  this.templates = {
    settings: {
      label: '<label class="fbs-lbl {{pathClass}}" for="{{id}}">{{varToLabel label}}</label>',
      hiddenField: '{{#ifEq type "checkbox"}}' +
        '<input style="display: none;" class="hidden {{pathClass}}" id="{{id}}" type="{{type}}" name="{{path}}" {{#unless typeText}}{{#val}}checked="checked"{{/val}}{{/unless}}  {{#if typeText}}value="{{val}}"{{/if}}>' +
        '{{else}}' +
        '<input style="display: none;" class="hidden {{pathClass}}" id="{{id}}" type="{{type}}" name="{{path}}" {{#unless typeText}}{{#val}}checked="checked"{{/val}}{{/unless}}  {{#if typeText}}value="{{val}}"{{/if}}>' +
        '{{/ifEq}}',
      field: '   <table class="hcbuilder-settings-tbl">' +
        '       <tbody><tr>' +
        '           {{#if createI18nEl}}' +
        '           <td class="switch-td">' +
        '               <label class="i18n-toggle">' +
        '                   <input class="i18n-switch" data-targetid="{{id}}" type="checkbox" {{#if i18nChecked}}checked="checked"{{/if}}>' +
        '                   <div class="slider"><span class="i18n-vertical">i18n</span></div>' +
        '               </label>' +
        '           </td>' +
        '           {{/if}}' +
        '           <td>' +
        '               {{#ifEq type "checkbox"}}' +
        '               <label class="switch"><input class="{{pathClass}} {{#if typeText}}form-control{{/if}} {{#if i18nChecked}}i18n-active{{/if}}" id="{{id}}" type="{{type}}" name="{{path}}" {{#unless typeText}}{{#val}}checked="checked"{{/val}}{{/unless}}  {{#if typeText}}value="{{escapeForAttr val}}"{{/if}}> <div class="slider"></div> </label>' +
        '               {{else}}' +
        '               <input class="{{pathClass}} {{#if typeText}}form-control{{/if}} {{#if i18nChecked}}i18n-active{{/if}}" id="{{id}}" type="{{type}}" name="{{path}}" {{#unless typeText}}{{#val}}checked="checked"{{/val}}{{/unless}}  {{#if typeText}}value="{{escapeForAttr val}}"{{/if}}>' +
        '               {{/ifEq}}' +
        '           </td>' +
        '       </tr></tbody>' +
        '   </table>',
      enumField: '   <table class="hcbuilder-settings-tbl enum-field">' +
        '       <tbody><tr>' +
        '           <td>' +
        '               {{#each val.all}}' +
        '               <div class="fancy-el fancy-radio float-left">' +
        '                   <input type="radio" {{#ifEq @index @root.val.selected}}checked="checked"{{/ifEq}} id="{{@root.path}}{{@index}}" name="{{@root.path}}.selected" value="{{@index}}" />' +
        '                   <label for="{{@root.path}}{{@index}}"><span><i></i></span><p>{{this}}</p></label>' +
        '               </div>' +
        '               {{/each}}' +
        '               {{#each val.all}}' +
        '                   <input type="hidden" name="{{@root.path}}.all[{{@index}}]" value="{{this}}"/>' +
        '               {{/each}}' +
        '           </td>' +
        '       </tr></tbody>' +
        '   </table>',
      enumFieldSelect: '   <table class="hcbuilder-settings-tbl enum-field">' +
        '       <tbody><tr>' +
        '           {{#ifEq label "file"}}<td><table><tr><td>exact</td><td>' +
        '               <label class="switch" style="margin-right: 10px;margin-left: 5px;"><input class="file-exact {{pathClass}}exact " type="checkbox" data-path="{{path}}" data-key="{{key}}" name="{{path}}.exact" {{#if val.exact}}checked="checked"{{/if}}> <div class="slider"></div> </label>' +
        '           </td></tr></table></td>{{/ifEq}}' +
        '           <td>' +
        '               <select class="{{pathClass}}" id="{{id}}" name="{{@root.path}}.selected">' +
        '                   {{#each val.all}}' +
        '                   <option {{#ifEq @index @root.val.selected}}selected="selected"{{/ifEq}} value="{{@index}}">{{this}}</option>' +
        '                   {{/each}}' +
        '               </select>' +
        '               {{#each val.all}}' +
        '                   <input type="hidden" name="{{@root.path}}.all[{{@index}}]" value="{{this}}"/>' +
        '               {{/each}}' +
        '           </td>' +
        '       </tr></tbody>' +
        '   </table>'
    }
  }
  this.templates.startIndention =
    '<div{{#if hidden}} style="display:none;" {{/if}} class="{{pathClass}} relative-parent row {{#if path}}root-row {{/if}}indention{{#if isValidate}} validation-parent{{/if}}">' +
    '   <div class="col-md-{{#if path}}3{{else}}1{{/if}}">' +
    '       {{#if path}}<label class="fbs-lbl {{pathClass}}" name="{{path}}" >{{label}}</label>{{/if}}' +
    '       {{#if isArray}}' +
    '       <button type="button" name="{{path}}" class="array-add btn btn-secondary" >' +
    '           <span class="fa fa-plus" aria-hidden="true"></span>' +
    '       </button>' +
    '       {{/if}}' +
    '   </div>' +
    '   <div class="col-md-{{#if path}}9{{else}}11{{/if}}">' +
    '       <div class="array-childs">'
  this.templates.endIndention =
    '       </div>' +
    '   </div>' +
    '</div>'

  this.templates.enumRowSingle =
    '<div class="{{#if select}}fb-field-group row{{else}}relative-parent fb-field-group row {{#if path}}root-row {{/if}}{{/if}}">' +
    '   {{#if label}}' +
    '   <div class="col-md-{{#ifEq label "file"}}2{{else}}3{{/ifEq}}">' +
    this.templates.settings.label +
    '   </div>' +
    '   {{/if}}' +
    '   <div class="col-md-{{#ifEq label "file"}}6{{else}}{{#if select}}5{{else}}{{#if label}}9{{else}}12{{/if}}{{/if}}{{/ifEq}}">{{#if select}}' +
    this.templates.settings.enumFieldSelect + '{{else}}' + this.templates.settings.enumField +
    '   {{/if}}</div>' +
    '   {{#unless select}}{{#if label}}' +
    '   {{/if}}{{/unless}}' +
    '</div>'

  this.templates.inputRowSingleDate =
    '<div class="fb-field-group row vdf">' +
    '   <div class="col-md-12">' +
    '       <div class="row">' +
    '           <div class="col-md-3">' +
    this.templates.settings.label +
    '           </div>' +
    '           <div class="col-md-{{#if label}}9{{else}}12{{/if}}">' +
    '               <table class="hcbuilder-settings-tbl">' +
    '                   <tbody name="{{path}}"><tr>' +
    '                       <td class="td-min vdf-ext-main vdf-before vdf-inactive">' +
    '                           <table class="vdf-input"><tbody><tr>' +
    '                               <td class="td-min">' +
    '                                   <input class="form-control vdf-date vdf-var" placeholder="now or variable name" type="text" data-dval="undefined" name="{{path}}.before.date" value="{{val.before.date}}"/>' +
    '                                   <span class="{{#if val.before.date}}vdf-true {{/if}}glyphicon glyphicon-calendar" aria-hidden="true"></span>' +
    '                               </td>' +
    '                               <td class="td-min">' +
    '                                   <input class="form-control vdf-pattern" placeholder="dd.MM.yyyy HH:mm:ss" type="text" data-dval="undefined" name="{{path}}.before.pattern" value="{{val.before.pattern}}"/>' +
    '                                   <span class="{{#if val.before.pattern}}vdf-true {{/if}}glyphicon glyphicon-equalizer" aria-hidden="true"></span>' +
    '                               </td>' +
    '                               <td class="td-min">' +
    '                                   <input class="form-control vdf-val" placeholder="01.01.2017 23:59:59" type="text" data-dval="undefined" name="{{path}}.before.val" value="{{val.before.val}}"/>' +
    '                                   <span class="{{#if val.before.val}}vdf-true {{/if}}glyphicon glyphicon-dashboard" aria-hidden="true"></span>' +
    '                               </td>' +
    '                           </tr></tbody></table>' +
    '                           <span class="vdf-ext{{#if val.before}} vdf-true{{/if}} glyphicon glyphicon-step-backward" aria-hidden="true"></span>' +
    '                       </td>' +
    '                       <td>' +
    '                           <input class="form-control" id="{{id}}" placeholder="dd.MM.yyyy HH:mm:ss" type="{{type}}" name="{{path}}.pattern" value="{{val.pattern}}" />' +
    '                       </td>' +
    '                       <td class="td-min vdf-ext-main vdf-after vdf-inactive">' +
    '                           <table class="vdf-input"><tbody><tr>' +
    '                               <td class="td-min">' +
    '                                   <input class="form-control vdf-date vdf-var" placeholder="now or variable name" type="text" data-dval="undefined" name="{{path}}.after.date" value="{{val.after.date}}"/>' +
    '                                   <span class="{{#if val.after.date}}vdf-true {{/if}}glyphicon glyphicon-calendar" aria-hidden="true"></span>' +
    '                               </td>' +
    '                               <td class="td-min">' +
    '                                   <input class="form-control vdf-pattern" placeholder="dd.MM.yyyy HH:mm:ss" type="text" data-dval="undefined" name="{{path}}.after.pattern" value="{{val.after.pattern}}"/>' +
    '                                   <span class="{{#if val.after.pattern}}vdf-true {{/if}}glyphicon glyphicon-equalizer" aria-hidden="true"></span>' +
    '                               </td>' +
    '                               <td class="td-min">' +
    '                                   <input class="form-control vdf-val" placeholder="01.01.2017 23:59:59" type="text" data-dval="undefined" name="{{path}}.after.val" value="{{val.after.val}}"/>' +
    '                                   <span class="{{#if val.after.val}}vdf-true {{/if}}glyphicon glyphicon-dashboard" aria-hidden="true"></span>' +
    '                               </td>' +
    '                           </tr></tbody></table>' +
    '                           <span class="vdf-ext{{#if val.after}} vdf-true{{/if}} glyphicon glyphicon-step-forward" aria-hidden="true"></span>' +
    '                       </td>' +
    '                   </tr></tbody>' +
    '               </table>' +
    '           </div>' +
    '       </div>' +
    '   </div>' +
    '</div>'

  this.templates.inputRowSingle =
    '<div class="relative-parent fb-field-group row {{#if path}}root-row {{/if}}">' +
    '   {{#if label}}' +
    '   <div class="col-md-3">' +
    this.templates.settings.label +
    '   </div>' +
    '   {{/if}}' +
    '   <div class="col-md-{{#if label}}9{{else}}12{{/if}}">' +
    this.templates.settings.field +
    '   </div>' +
    '   {{#if label}}' +
    '   {{/if}}' +
    '</div>'

  this.templates.inputRowMultiple =
    '<div class="fb-field-group row">' +
    '   {{#if label}}' +
    '   <div class="col-md-3">' +
    this.templates.settings.label +
    '   </div>' +
    '   {{/if}}' +
    '   <div class="col-md-{{#if label}}5{{else}}6{{/if}}">' +
    this.templates.settings.field +
    '   </div>' +
    '</div>'
  this.templates.array = {
    item: {
      start: '<div class="relative-parent array-item">',
      end: '<button type="button" class="btn btn-light btn-sm array-item-move array-item-move-up">' +
        '       <small><span class="fa fa-arrow-up text-secondary" aria-hidden="true"></span></small>' +
        '   </button>' +
        '<button type="button" class="btn btn-light btn-sm array-item-move array-item-move-down">' +
        '       <small><span class="fa fa-arrow-down text-secondary" aria-hidden="true"></span></small>' +
        '   </button>' +
        '<button type="button" class="btn btn-danger array-item-del">' +
        '       <span class="fa fa-minus" aria-hidden="true">3</span>' +
        '   </button>' +
        '</div>'
    }
  }
  this.templates.validate = {
    dropDown: '<div class="dropdown validate-dropdown">' +
      '  <button class="btn dropdown-toggle validate-dropdown-toggle" type="button" id="{{id}}" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">' +
      '   {{label}}' +
      '  </button>' +
      '  <div class="dropdown-menu fb-validation-menu" aria-labelledby="{{id}}">{{lis}}</div>' +
      '</div>',
    dropDownLi: '<li class="{{#if disabled}}disabled{{/if}}">' +
      '   <a href="javascript:void(0)" data-key="{{key}}" data-exact="{{val.exact}}" data-kind="{{val.kind}}" name="{{path}}.{{key}}">{{key}}</a>' +
      '</li>',
    deleteBtn: '<td class="td-min"><button type="button" name="{{path}}" class="btn btn-danger validation-del">' +
      '   <span class="fa fa-minus" aria-hidden="true"></span>' +
      '</button></td>'
  }
  this.validate = {
    dropdown: {
      defaults: {
        'required': true,
        'email': true,
        'number': true,
        'matches': 'regex',
        'max': 10,
        'min': 1,
        'url': true,
        'phoneNr': true
      },
      replaceLabel: function (json, fb, baseEle) {
        var _ = fb.settingsTab
        var validateLabel = baseEle.find("label[name$='validate']")
        if (validateLabel && validateLabel.length) {
          validateLabel.each(function () {
            var _thisVlabel = $(this)
            var validate = null
            var vlPath = _thisVlabel.attr('name')
            var isFile = false
            if ($.isArray(json)) {
              var arrIndex = vlPath.match(/^\[(\d+)\]/)
              if (arrIndex && arrIndex.length == 2) {
                validate = json[arrIndex[1]]
                isFile = validate['_file']
                validate = validate['validate']
              }
            } else {
              isFile = json['_file']
              validate = json['validate']
            }
            if (validate) {
              var validationDropDownEl = $(_.fb.compileTemplate(_.templates.validate.dropDown, {
                id: _.fb.randomId(),
                label: _thisVlabel.toString(),
                lis: _.validate.dropdown.createLis(fb, vlPath, validate, isFile)
              }, 'templates.validate.dropDown'))
              // _.createValidationSelectEvent(fb, validationDropDownEl);
              _thisVlabel.replaceWith(validationDropDownEl)
            }
          })
          _.validate.deleteButton.addAllDeleteButton(fb, baseEle)
        }
      },
      createLis: function (fb, path, validate, isFile) {
        if (!validate) {
          validate = {}
        }
        var html = ''
        for (var key in fb.settingsTab.validate.dropdown.defaults) {
          if (fb.settingsTab.validate.dropdown.defaults.hasOwnProperty(key)) {
            html += fb.compileTemplate(fb.settingsTab.templates.validate.dropDownLi, {
              disabled: (key in validate),
              path: path,
              key: key
            }, 'templates.validate.dropDownLi')
          }
        }
        if (isFile) {
          key = 'file'
          html += fb.compileTemplate(fb.settingsTab.templates.validate.dropDownLi, {
            disabled: (key in validate),
            path: path,
            key: key,
            val: validate.file ? validate.file : {}
          }, 'templates.validate.dropDownLi')
        }
        return html
      },
      createSelectEvent: function (fb, baseEl, json) {
        var _ = fb.settingsTab
        baseEl.find('.fb-validation-menu li a').click(function () {
          var _this = $(this)
          var _li = _this.parent()
          if (_li.hasClass('disabled')) {
            return false
          }
          _li.addClass('disabled')
          var path, type, htmlStr, val, key
          path = _this.attr('name')
          key = _this.attr('data-key')
          val = _.validate.dropdown.defaults[key]
          if (typeof val === 'boolean') {
            type = 'checkbox'
            htmlStr = _.transformer.json.createInputRow({
              fb: fb,
              val: val
            }, key, path, type, true)
          } else if (key === 'file' || fb.isEnum(val)) {
            if (!val) {
              val = {}
            }
            val.exact = _this.attr('data-exact')
            val.kind = _this.attr('data-kind')
            htmlStr = _.transformer.json.enumToHtmlInput({
              fb: fb,
              val: val
            }, key, path)
          } else /* if(typeof val === 'string') */ {
            type = 'text'
            htmlStr = _.transformer.json.createInputRow({
              fb: fb,
              val: val
            }, key, path, type, true)
          }
          var newEl = $(htmlStr)
          _.validate.file.createExactSwitchEvent(fb, newEl)

          var newDelBtnEl = _.validate.deleteButton.addDeleteButton(fb, path, newEl.find('.hcbuilder-settings-tbl>tbody>tr>td:last-child'))
          _.validate.deleteButton.createDeleteEvent(fb, newDelBtnEl)
          _.validate.date.createDateEvents(fb, newEl)
          var rootRowEl = _this.nextParentWithClass('root-row')
          rootRowEl.find('.array-childs').append(newEl)
          // Setup events for the newly added validation input (PC-526)
          _.transformer.json.createEvents(fb, fb.settingsTab.body)
          _.settingsChanged(_this, 'click')
        })
      }
    },
    file: {
      selectedTitle: function (fb, _t) {
        var selected = _t.find(':selected').text()
        if (!selected) {
          selected = _t.children().first().text()
        }
        if (_t.data('exact')) {
          try {
            _t.attr('title', fb.fileTypes['Exact'][selected]['MIME']['Value'])
          } catch (e) {
            console.log(e)
          }
        } else {
          try {
            var arr = fb.fileTypes['Vague'][selected]
            var t = ''
            for (var i = 0; i < arr.length; i++) {
              if (t !== '') {
                t += ', ' + arr[i].Extension
              } else {
                t += arr[i].Extension
              }
            }
            _t.attr('title', t)
          } catch (e) {
            console.log(e)
          }
        }
      },
      createExactSwitchEvent: function (fb, el) {
        var _ = fb.settingsTab
        var fileExactClickEvent
        fileExactClickEvent = function () {
          var _this = $(this)
          var exact = _this.is(':checked')
          var path = _this.attr('data-path')
          var key = 'file'
          var val = {
            exact: exact
          }
          var hStr = _.transformer.json.enumToHtmlInput({
            fb: fb,
            val: val
          }, key, path)
          var newEl = $(hStr)
          var fileExactSwitch = newEl.find('input.file-exact')
          fileExactSwitch.click(fileExactClickEvent)
          var $sel = newEl.find('select')
          $sel.data('exact', fileExactSwitch.is(':checked'))
          $sel.change(function () {
            var _t = $(this)
            _.validate.file.selectedTitle(fb, _t)
            fb.settingsTab.settingsChanged(_t, 'change')
          })
          _.validate.file.selectedTitle(fb, $sel)
          var newDelBtnEl = _.validate.deleteButton.addDeleteButton(fb, path, newEl.find('.hcbuilder-settings-tbl>tbody>tr>td:last-child'))
          _.validate.deleteButton.createDeleteEvent(fb, newDelBtnEl)
          _this.closest('.fb-field-group.row').replaceWith(newEl)
          fb.settingsTab.settingsChanged(_this, 'click')
        }
        var fileExactSwitch = el.find('input.file-exact')
        if (fileExactSwitch.length) {
          fileExactSwitch.each(function () {
            var _this = $(this)
            _this.unbind('click')
            _this.click(fileExactClickEvent)
            var $sel = _this.closest('.enum-field').find('select')
            $sel.data('exact', _this.is(':checked'))
            $sel.unbind('change')
            $sel.change(function () {
              var _t = $(this)
              _.validate.file.selectedTitle(fb, _t)
              fb.settingsTab.settingsChanged(_t, 'change')
            })
            _.validate.file.selectedTitle(fb, $sel)
          })
        }
      }
    },
    date: {
      createDateEvents: function (fb, baseEle) {
        baseEle.find('.vdf-ext-main .vdf-ext').click(function (e) {
          var t = $(this)
          var p = t.parent()
          var vdf = p.nextParentWithClass('vdf')
          var vdfMains = vdf.find('.vdf-ext-main.vdf-active')
          vdfMains.each(function () {
            var b = $(this)
            var tt = b.find('.vdf-ext')
            if (tt[0] != t) {
              b.removeClass('td-auto').addClass('td-min')
              b.find('.td-auto').removeClass('td-auto').addClass('td-min')
              b.removeClass('vdf-active').addClass('vdf-inactive')
            }
          })
          p.removeClass('vdf-inactive').addClass('vdf-active')
        })
        baseEle.find('.vdf-input span.glyphicon').bind('click', function () {
          var t = $(this)
          var p = t.nextParentWithClass('vdf-ext-main')
          p.find('.td-auto').removeClass('td-auto').addClass('td-min')
          t.parent().removeClass('td-min').addClass('td-auto')
          p.removeClass('td-min').addClass('td-auto')
        })
      }
    },
    deleteButton: {
      addAllDeleteButton: function (fb, baseEle) {
        var skipValidateFor = ['datePattern']
        baseEle.find('[name*="validate."]').each(function () {
          var t = $(this)
          var m = /^((\[\d+\]\.)?(validate\.[A-Za-z]+)\.?)/.exec(t.attr('name'))
          if (t.attr('type') !== 'hidden' && m && m.length > 3 && !skipValidateFor.find(s => m[3].endsWith(s))) {
            fb.settingsTab.validate.deleteButton.addDeleteButton(fb, m[0], t.nextParentWithClass('hcbuilder-settings-tbl').parent().find('.hcbuilder-settings-tbl>tbody>tr>td:last-child'))
          }
        })
      },
      addDeleteButton: function (fb, path, el) {
        var _ = fb.settingsTab
        if (el.find('.validation-del').length === 0) {
          var htmlStr = fb.compileTemplate(_.templates.validate.deleteBtn, {
            path: path
          }, 'templates.validate.deleteBtn')
          var newDelEl = $(htmlStr)
          // _.validate.deleteButton.createDeleteEvent(fb, newDelEl);
          el.after(newDelEl)
          return newDelEl
        }
        return null
      },
      createDeleteEvent: function (fb, baseEle) {
        if (!baseEle) {
          return
        }
        var allDelBtns = baseEle.find('button.validation-del')
        allDelBtns.unbind('click')
        allDelBtns.click(function () {
          var _this = $(this)
          var parent = _this.nextParentWithClass('indention')
          parent.find('.fb-validation-menu li a[name="' + _this.attr('name') + '"]').parent().removeClass('disabled')
          _this.nextParentWithClass('fb-field-group').remove()
          fb.settingsTab.settingsChanged(_this, 'click')
        })
      }
    },
    createAllEvents: function (fb, baseEle, json) {
      fb.settingsTab.validate.deleteButton.createDeleteEvent(fb, baseEle)
      fb.settingsTab.validate.dropdown.createSelectEvent(fb, baseEle, json)
      fb.settingsTab.validate.date.createDateEvents(fb, baseEle)
    }
  }
  this.transformer = {
    json: {
      i18n: {
        autocompleteDefaults: fb.options.i18n ? {
          updater: function (item) {
            var i18nData = fb.options.i18n.onSelect(item.id)
            this.$element.val(item.id)
            if (!this.$element.hasClass('i18n-active') && fb.options.i18n.isCovered(i18nData)) {
              this.$element.addClass('i18n-active')
            }
            fb.settingsTab.settingsChanged(this.$element)
            return fb.options.i18n.onDisplay(item.id)
          },
          showHintOnFocus: true,
          autoSelect: false,
          source: fb.options.i18n.onSearch,
          fitToElement: true,
          items: 6
        } : null,
        initEvents: function (fb, baseEle) {
          if (fb.options.i18n) {
            baseEle.find('input.i18n-switch').each(function () {
              var i18nSwitch = $(this)
              fb.settingsTab.transformer.json.i18n.createClickEvent(fb, i18nSwitch)
              if (i18nSwitch.is(':checked')) {
                fb.settingsTab.transformer.json.i18n.initAutocomplete(fb, i18nSwitch.nextParentWithClass('hcbuilder-settings-tbl').find('#' + i18nSwitch.attr('data-targetid')))
              }
            })
          }
        },
        createClickEvent: function (fb, i18nSwitch) {
          i18nSwitch.click(function () {
            var _this = $(this)
            if (_this.is(':checked')) {
              fb.settingsTab.transformer.json.i18n.initAutocomplete(fb, _this.nextParentWithClass('hcbuilder-settings-tbl').find('#' + _this.attr('data-targetid')), true)
            } else {
              fb.settingsTab.transformer.json.i18n.destroyAutocomplete(fb, _this.nextParentWithClass('hcbuilder-settings-tbl').find('#' + _this.attr('data-targetid')))
            }
            fb.settingsTab.settingsChanged(_this, 'click')
          })
        },
        initAutocomplete: function (fb, inputEl, focus) {
          inputEl.attr('autocomplete', 'off')
          inputEl.attr('spellcheck', 'false')
          inputEl.typeahead(fb.settingsTab.transformer.json.i18n.autocompleteDefaults)
          if (focus) {
            if (inputEl[0].selectionStart || inputEl[0].selectionStart == '0') {
              var elemLen = inputEl[0].value.length
              // Firefox/Chrome
              inputEl[0].selectionStart = elemLen
              inputEl[0].selectionEnd = elemLen
            }
            inputEl.focus()
          }
        },
        destroyAutocomplete: function (fb, inputEl) {
          inputEl.typeahead('destroy')
          inputEl.unbind()
          inputEl.removeClass('i18n-active')
          fb.settingsTab.transformer.json.bindInputEvents(inputEl, fb.settingsTab)
        }
      },
      toHtmlInput: function (fb, json) {
        var buffer = {
          html: '',
          objPath: ''
        }
        if ($.isArray(json)) {
          this.arrayToHtmlInput({
            fb: fb,
            key: '',
            val: json,
            buffer: buffer,
            rootWrap: true
          }, '')
        } else {
          this.objectToHtmlInput({
            fb: fb,
            key: '',
            val: json,
            buffer: buffer
          }, '')
        }
        var $html = $(buffer.html)
        fb.settingsTab.validate.dropdown.replaceLabel(json, fb, $html)
        setTimeout(function () {
          $html.find('.actionsource .array-item').each(function () {
            $(this).find('.actionsourcecomment, .actionsourceregex').on('click', function () {
              var $arrayItem = $(this).nextParentWithClass('array-item')
              var $uipp = $arrayItem.find('.actionsource_fbonly_uipp')
              var $destCompId = $arrayItem.find('.actionsource_destCompId')
              var uipp, destCompId
              if ($destCompId.length) {
                destCompId = $destCompId.val()
              }
              if ($uipp.length) {
                uipp = $uipp.val()
              }
              fb.workspace.connectionManager.showAction($arrayItem, uipp, destCompId)
            })
          })

          $html.find('input.var-name').each(function () {
            $(this).attr('pattern', '^[a-zA-Z]?[a-zA-Z0-9]+$')

            $(this).on('click', function () {
              $(this).typeahead('lookup').focus()
              var $varName = $(this)
              var name = $varName.val()
              if (name) {
                fb.workspace.getActiveBody().find('.fb-component.selected').find("[fb_name='" + name + "'],[name='" + name + "']").each(function () {
                  var t = $(this)
                  var n = t.attrFBName()
                  if (n) {
                    var p = $varName.parent().nextParentWithClass('var-name')
                    if (p.length) {
                      var startId
                      var endId = p.attrID()
                      t = t.findVisibleInputOrRefElement()
                      startId = t.attrID()
                      fb.workspace.connectionManager.showName({
                        animate: '0.3s',
                        offset: 10,
                        orientation: 'vertical',
                        stroke: '#007800',
                        strokeWidth: 1,
                        start: '#' + startId,
                        end: '.hcbuilder-settings-body #' + endId
                      })
                    }
                  }
                })
              }
            })
          })
        }, 100)

        fb.settingsTab.body.html($html)
        this.createEvents(fb, fb.settingsTab.body, json)
      },
      createEvents: function (fb, baseEle, json) {
        var _ = fb
        var changables = baseEle.find('input[name], select[name]')
        var nameSelector = "input[name='name'], input[name*='.name']"
        var nameEl = baseEle.find(nameSelector)
        var _fbOptions = _.options
        nameEl.each(function () {
          var _thisNameEl = $(this)
          var nameRootRowEl = _thisNameEl.nextParentWithClass('root-row')
          if (nameRootRowEl.hasClass('root-row')) {
            var rootRowParent = nameRootRowEl.parent()
            rootRowParent.prepend(nameRootRowEl.detach())
            _thisNameEl = nameRootRowEl.find(nameSelector)
            nameRootRowEl.find("label[for='" + _thisNameEl.attr('id') + "']").addClass('var-name')
            nameRootRowEl.addClass('var-name')
          }
          _thisNameEl.addClass('var-name')
          _thisNameEl.attr('autocomplete', 'off')
          _thisNameEl.attr('spellcheck', 'false')
          _thisNameEl.typeahead({
            minLength: 0,
            autoSelect: true,
            source: _fbOptions.varNameFunction ? _fbOptions.varNameFunction : _fbOptions.vars,
            items: 20
          })
        })

        var connectorSelector = 'input.connector'
        var conEl = baseEle.find(connectorSelector)
        conEl.each(function () {
          var _thisNameEl = $(this)
          _thisNameEl.attr('autocomplete', 'off')
          _thisNameEl.attr('spellcheck', 'false')
          _thisNameEl.typeahead({
            minLength: 0,
            autoSelect: true,
            source: _fbOptions.varConnectorsFunction,
            items: 20
          })
          $(_thisNameEl).on('click', function () {
            $(_thisNameEl).typeahead('lookup').focus()
          })
        })

        var settingsTab = _.settingsTab
        baseEle.find('.array-add').click(function () {
          var arrayAddEl = $(this)
          var parent = arrayAddEl.nextParentWithClass('indention')
          var arrayChilds = parent.find('.array-childs:first')
          var lastArrayItem = arrayChilds.children('.array-item:last')
          var newArrayItem = $(lastArrayItem.toString())

          var inputs = newArrayItem.find('input[name]')
          inputs.each(function () {
            var _in = $(this)
            var newId = _.randomId()
            var oldId = _in.attr('id')
            var i18nSwitch = newArrayItem.find('[data-targetid="' + oldId + '"]')
            if (i18nSwitch.length) {
              i18nSwitch.attr('data-targetid', newId)
              _in.removeClass('i18n-active')
            }
            newArrayItem.find('label[for="' + oldId + '"]').attr('for', newId)
            _in.attr('id', newId)
            if (!(/^(\[\d+\]\.)?(validate\..*)$/.test(_in.attr('name'))) && _in.attr('type') === 'text') {
              _in.val('')
            }
            _in.attr('name', settingsTab.transformer.json.incrementArrayPattern(arrayAddEl.attr('name'), _in.attr('name')))
          })
          newArrayItem.find('button[name],label[name],a[name]').each(function () {
            $(this).attr('name', settingsTab.transformer.json.incrementArrayPattern(arrayAddEl.attr('name'), $(this).attr('name')))
          })
          settingsTab.transformer.json.createEvents(_, newArrayItem, json)
          arrayChilds.append(newArrayItem)
          settingsTab.settingsChanged(arrayAddEl, 'click')
        })

        function moveArrayItem(target, direction) {
          let $arrayMoveUpEl = target
          let _arrayItem = $arrayMoveUpEl.parent()

          if (direction === 'up') {
            _arrayItem.prev().before(_arrayItem)
          } else if (direction === 'down') {
            _arrayItem.next().after(_arrayItem)
          } else {
            return
          }

          let _arrayIndention = _arrayItem.nextParentWithClass('indention')
          let arrayAddBtn = _arrayIndention.find('.array-add')
          let _arrayParentPath = arrayAddBtn.attr('name')

          let i = 0
          _arrayIndention.find('.array-item').not(_arrayIndention.find('.array-item .array-item')).each(function () {
            $(this).find('input[name],button[name],label[name]').each(function () {
              $(this).attr('name', settingsTab.transformer.json.moveArrayPattern(_arrayParentPath, $(this).attr('name'), i))
            })
            i++
          })

          settingsTab.settingsChanged($arrayMoveUpEl, 'click')
        }

        baseEle.find('.array-item-move-up').click(function () {
          moveArrayItem($(this), 'up')
        })

        baseEle.find('.array-item-move-down').click(function () {
          moveArrayItem($(this), 'down')
        })
        settingsTab.transformer.json.bindInputEvents(changables, settingsTab)
        baseEle.find('button').each(function () {
          var cgl = $(this)
          cgl.bind('keyup', function (event) {
            return settingsTab.transformer.json.keyupEvent($(this), event, settingsTab)
          })
          cgl.bind('keydown', function (event) {
            return settingsTab.transformer.json.keydownEvent($(this), event, settingsTab)
          })
        })
        baseEle.find('.array-item-del').click(function () {
          var arrayDelBtn = $(this)
          var _arrayItem = arrayDelBtn.parent()
          if (!arrayDelBtn.hasClass('simple-item')) {
            if (_arrayItem.is(':last-child')) {} else {
              var _arrayIndention = _arrayItem.nextParentWithClass('indention')
              var arrayAddBtn = _arrayIndention.find('.array-add')
              var _arrayParentPath = arrayAddBtn.attr('name')
              while (true) {
                _arrayItem = _arrayItem.next('.array-item')
                if (_arrayItem.length == 0) {
                  break
                }
                _arrayItem.find('input[name],button[name],label[name]').each(function () {
                  $(this).attr('name', settingsTab.transformer.json.decrementArrayPattern(_arrayParentPath, $(this).attr('name')))
                })
              }
            }
          }
          arrayDelBtn.parent().remove()
          settingsTab.settingsChanged(arrayDelBtn, 'click')
        })
        settingsTab.transformer.json.i18n.initEvents(_, baseEle)
        settingsTab.validate.createAllEvents(_, baseEle, json)
        settingsTab.validate.file.createExactSwitchEvent(_, baseEle)
      },
      bindInputEvents: function (changables, settingsTab) {
        changables.each(function () {
          var cgl = $(this)
          cgl.bind('keyup', function (event) {
            return settingsTab.transformer.json.keyupEvent($(this), event, settingsTab)
          })
          cgl.bind('keydown', function (event) {
            return settingsTab.transformer.json.keydownEvent($(this), event, settingsTab)
          })
          cgl.bind('change', function (event) {
            if ($(this).attr('type') !== 'text') {
              settingsTab.settingsChanged($(this), 'change')
            } else {
              if (settingsTab.autoSave) {
                settingsTab.settingsChanged($(this), 'keyup')
              }
            }
          })
        })
      },
      objectToHtmlInput: function (o, key, path, isNested) {
        var iHideIt = false
        if (key) {
          if (o.fb.isHiddenField(key)) {
            o.isParentHidden = true
            iHideIt = true
          } else {
            o.buffer.html += this.startIndention(o, key, path, false)
          }
        }
        // var multipleInputs = o.fb.settingsTab.sizeOf(o.val) == 2

        var multipleInputs = !!isNested

        var json = o.val
        for (var k in json) {
          if (json.hasOwnProperty(k)) {
            if (json[k] !== null) {
              var subpath = (!path ? k : (path + '.' + k))
              o.val = json[k]
              if (typeof json[k] === 'boolean') {
                this.booleanToHtmlInput(o, k, subpath, multipleInputs)
              } else if ($.isArray(json[k])) {
                this.arrayToHtmlInput(o, k, subpath)
              } else if (/^(\[\d+\]\.)?(validate\.file.*)$/.test(subpath) || o.fb.isEnum(json[k])) {
                this.enumToHtmlInput(o, k, subpath)
              } else if (typeof json[k] === 'object') {
                if ((/^(\[\d+\]\.)?(validate\.date)$/.test(subpath))) {
                  this.stringToHtmlInput(o, k, subpath, multipleInputs)
                } else if (o.fb.options.i18n.isCovered(json[k])) {
                  this.stringToHtmlInput(o, k, subpath, multipleInputs)
                } else {
                  this.objectToHtmlInput(o, k, subpath, true)
                }
              } else {
                this.stringToHtmlInput(o, k, subpath, multipleInputs)
              }
            }
          }
        }
        if (iHideIt) {
          o.isParentHidden = false
        } else {
          if (key) {
            o.buffer.html += this.endIndention(o)
          }
        }
        return o.buffer
      },
      // fb, path, index, json, buffer, rootWrap
      arrayToHtmlInput: function (o, key, path) {
        var iHideIt = false
        if (o.fb.isHiddenField(key)) {
          o.isParentHidden = true
          iHideIt = true
        } else {
          if (key || o.rootWrap) {
            o.buffer.html += this.startIndention(o, key, path, true)
          }
        }
        var json = o.val
        for (var i = 0; i < json.length; ++i) {
          if (!o.isParentHidden) {
            o.buffer.html += o.fb.settingsTab.templates.array.item.start
          }
          var subpath = (!path ? (!key ? '[' + i + ']' : key) : (path + '[' + i + ']'))
          o.val = json[i]
          if (typeof o.val === 'boolean') {
            this.booleanToHtmlInput(o, '', subpath)
          } else if ($.isArray(o.val)) {
            this.arrayToHtmlInput(o, '', subpath)
          } else if (o.fb.isEnum(o.val)) {
            this.enumToHtmlInput(o, '', subpath)
          } else if (typeof o.val === 'object') {
            this.objectToHtmlInput(o, '', subpath)
          } else {
            this.stringToHtmlInput(o, '', subpath)
          }
          if (!o.isParentHidden) {
            o.buffer.html += o.fb.settingsTab.templates.array.item.end
          }
        }

        if (iHideIt) {
          o.isParentHidden = false
        } else {
          if (key || o.rootWrap) {
            o.buffer.html += this.endIndention(o)
            o.rootWrap = false
          }
        }
        return o.buffer
      },
      fileTypesToEnum: function (fb, exact, kind) {
        var en = {
          all: [],
          selected: 0
        }
        var obj
        if (exact === true) {
          obj = fb.fileTypes['Exact']
        } else {
          obj = fb.fileTypes['Vague']
        }
        if (obj) {
          for (var key in obj) {
            if (obj.hasOwnProperty(key)) {
              en.all.push(key)
            }
          }
        }
        if (kind) {
          for (var i = 0; i < en.all.length; i++) {
            if (en.all[i] === kind) {
              en.selected = i
              break
            }
          }
        }
        return en
      },
      enumToHtmlInput: function (o, key, path) {
        if (o.fb.isHiddenField(key)) {
          return ''
        }
        var isFile = (/^(\[\d+\]\.)?(validate\.file.*)$/.test(path))
        if (isFile) {
          o.val = $.extend({}, o.val, this.fileTypesToEnum(o.fb, o.val ? o.val.exact : false, o.val ? o.val.kind : null))
        }
        var html = o.fb.compileTemplate(o.fb.settingsTab.templates.enumRowSingle, {
          id: o.fb.randomId(),
          label: '' + key,
          val: o.val,
          path: path,
          exact: o.val ? o.val.exact : false,
          pathClass: this.toPathClass(path),
          select: (/^(\[\d+\]\.)?(validate\..*)$/.test(path)) || isFile
        }, 'templates.enumRowSingle')
        if (o && o.buffer && o.buffer.html) {
          o.buffer.html += html
        } else {
          return html
        }
      },
      stringToHtmlInput: function (o, key, path, multipleInputs) {
        o.buffer.html += this.createInputRow(o, key, path, 'text', multipleInputs)
      },
      booleanToHtmlInput: function (o, key, path, multipleInputs) {
        o.buffer.html += this.createInputRow(o, key, path, 'checkbox', multipleInputs)
      },
      createInputRow: function (o, key, path, type, multipleInputs) {
        if (o.isParentHidden || o.fb.isHiddenField(key)) {
          return o.fb.compileTemplate(o.fb.settingsTab.templates.settings.hiddenField, {
            id: o.fb.randomId(),
            label: '' + key,
            orgVal: o.val,
            val: o.val,
            path: path,
            pathClass: this.toPathClass(path),
            type: type,
            typeText: type === 'text'
          }, 'templates.settings.hiddenField')
        }
        if (/^(\[\d+\]\.)?(validate\.date)$/.test(path)) {
          return o.fb.compileTemplate(o.fb.settingsTab.templates.inputRowSingleDate, {
            id: o.fb.randomId(),
            label: '' + key,
            orgVal: o.val,
            val: o.val,
            path: path,
            pathClass: this.toPathClass(path),
            type: type
          }, 'templates.inputRowSingleDate')
        } else {
          var tmplData = {
            id: o.fb.randomId(),
            label: '' + key,
            orgVal: o.val,
            val: o.val,
            path: path,
            pathClass: this.toPathClass(path),
            type: type,
            typeText: type === 'text',
            createI18nEl: false,
            i18nChecked: false
          }
          if ((/^(\[\d+\]\.)?(action\..*)$/.test(path))) {
            return o.fb.compileTemplate(o.fb.settingsTab.templates.inputRowSingle, tmplData, 'templates.inputRowSingle')
          } else if (!(/^(\[\d+\]\.)?(name|validate\..*)$/.test(path))) {
            if (o.fb.options.i18n) {
              if (type === 'text') {
                if ($.isFunction(o.fb.options.i18n.isCovered)) {
                  tmplData.i18nChecked = o.fb.options.i18n.isCovered(o.val)
                }
                if (o.fb.options.enableI18n && o.fb.options.enableI18n === true) {
                  tmplData.createI18nEl = true
                }
              }
              if ($.isFunction(o.fb.options.i18n.onDisplay)) {
                tmplData.val = o.fb.options.i18n.onDisplay(o.val)
              }
            }
          } else if (/^(\[\d+\]\.)?(name)$/.test(path)) {
            multipleInputs = false
          } else {
            multipleInputs = true
          }
          if (multipleInputs) {
            return o.fb.compileTemplate(o.fb.settingsTab.templates.inputRowMultiple, tmplData, 'templates.inputRowMultiple')
          }
          return o.fb.compileTemplate(o.fb.settingsTab.templates.inputRowSingle, tmplData, 'templates.inputRowSingle')
        }
      },
      toPathClass: function (path) {
        if (path) {
          path = path.replace(/\.|\[\d+\]/g, '')
          if (path === 'label') {
            return ''
          }
          return path
        }
        return ''
      },
      incrementArrayPattern: function (parent, target) {
        var arrayChild = target.substring(parent.length)
        var matches = arrayChild.match(/^\[(\d+)\]/)
        var pathSuffix = arrayChild.substring(matches[0].length)
        return parent + '[' + (parseInt(matches[1]) + 1) + ']' + pathSuffix
      },
      decrementArrayPattern: function (parent, target) {
        var arrayChild = target.substring(parent.length)
        var matches = arrayChild.match(/^\[(\d+)\]/)
        var pathSuffix = arrayChild.substring(matches[0].length)
        return parent + '[' + (parseInt(matches[1]) - 1) + ']' + pathSuffix
      },
      moveArrayPattern: function (parent, target, index) {
        var arrayChild = target.substring(parent.length)
        var matches = arrayChild.match(/^\[(\d+)\]/)
        var pathSuffix = arrayChild.substring(matches[0].length)
        return parent + '[' + (index) + ']' + pathSuffix
      },
      startIndention: function (o, key, path, isArray) {
        var tmplData = {
          path: path,
          hidden: o.isParentHidden,
          pathClass: this.toPathClass(path),
          label: '',
          isArray: isArray,
          isValidate: /^(\[\d+\]\.)?(validate)$/.test(path)
        }
        if (key) {
          key = '' + key
          tmplData.label = key.charAt(0).toUpperCase() + key.slice(1)
        }
        return o.fb.compileTemplate(o.fb.settingsTab.templates.startIndention, tmplData, 'templates.startIndention')
      },
      endIndention: function (o) {
        return o.fb.settingsTab.templates.endIndention
      },
      keyupEvent: function (t, event, settingsTab) {
        var code = event.keyCode || event.which
        if (code) {
          if (settingsTab.autoSave) {
            if (event.ctrlKey) {
              if (event.which === 86) { // Check for the Ctrl key being pressed, and if the key = [V] (86)
                settingsTab.settingsChanged(t, 'keyup')
                return true
              }
            } else {
              settingsTab.settingsChanged(t, 'keyup')
              return true
            }
          }
          if (event.ctrlKey && event.which === 86) { // Check for the Ctrl key being pressed, and if the key = [V] (86)
            settingsTab.settingsChanged(t, 'keyup')
            return true
          }
          if (!event.ctrlKey && (code > 46 || code < 9 ||
              (code !== 9 &&
                code !== 16 &&
                code !== 17 &&
                code !== 18 &&
                code !== 33 &&
                code !== 34 &&
                code !== 35 &&
                code !== 36 &&
                code !== 37 &&
                code !== 38 &&
                code !== 39 &&
                code !== 40 &&
                code !== 45 &&
                code !== 46))) {
            settingsTab.settingsChanged(t, 'keyup')
          }
        }
        return true
      },
      keydownEvent: function (t, event, settingsTab) {
        var code = event.keyCode || event.which
        if (code) {
          if (event.ctrlKey && event.which === 83) { // Check for the Ctrl key being pressed, and if the key = [S] (83)
            event.preventDefault()
            settingsTab.storeSettings()
            return false
          }
        }
        return true
      }
    },
    html: {
      /**
       * data-dval ["undefined" = don't assign, "null" assign empty, "some val" assign with this if empty]
       * @param fb
       * @param jqEl
       * @param o
       * @returns {*}
       */
      toJsonSettings: function (fb, jqEl, o) {
        if (!o) {
          o = $.isArray(fb.settingsTab.currentSettings) ? [] : {}
        }
        // var values = this.serializeArray();
        var allInputs = jqEl.find('input[name], select[name]')
        var targetInput = null
        for (var i = 0; i < allInputs.length; ++i) {
          targetInput = $(allInputs[i])
          var defaultValueIfEmpty = targetInput.attr('data-dval')
          var val
          if (targetInput.length) {
            if (targetInput.hasClass('file-exact') && targetInput.attr('type') === 'checkbox') {
              this.setValueByJsString(fb, o, targetInput.attr('name'), targetInput.is(':checked'))
            } else if (targetInput.attr('type') === 'radio' || targetInput.attr('type') === 'checkbox') {
              var cbOrRbCollection = []
              var tName = targetInput.attr('name')
              for (; i < allInputs.length; ++i) {
                var cbOrRb = $(allInputs[i])
                if (tName === cbOrRb.attr('name')) {
                  cbOrRbCollection.push(cbOrRb)
                } else {
                  --i
                  break
                }
              }
              if (cbOrRbCollection.length > 1) {
                if (targetInput.attr('type') === 'radio') {
                  for (var rb = 0; rb < cbOrRbCollection.length; ++rb) {
                    if (cbOrRbCollection[rb].is(':checked')) {
                      val = cbOrRbCollection[rb].val()
                      if (defaultValueIfEmpty) {
                        if (defaultValueIfEmpty === 'undefined') {
                          if (!val) {
                            break
                          }
                        } else if (defaultValueIfEmpty === 'null') {
                          val = ''
                        } else {
                          if (!val) {
                            val = defaultValueIfEmpty
                          }
                        }
                      }
                      this.setValueByJsString(fb, o, tName, val === 'on' ? true : val)
                      break
                    }
                  }
                }
              } else {
                val = cbOrRbCollection[0].val()
                if (defaultValueIfEmpty) {
                  if (defaultValueIfEmpty === 'undefined') {
                    if (!val) {
                      continue
                    }
                  } else if (defaultValueIfEmpty === 'null') {
                    val = ''
                  } else {
                    if (!val) {
                      val = defaultValueIfEmpty
                    }
                  }
                }
                this.setValueByJsString(fb, o, tName, val === 'on' ? cbOrRbCollection[0].is(':checked') : val)
              }
            } else {
              if (targetInput[0].value !== undefined) {
                val = targetInput.val()
              } else {
                val = targetInput.html()
              }
              if (targetInput.parents('.hcbuilder-settings-tbl').find('input.i18n-switch').is(':checked')) {
                val = fb.options.i18n.onSelect(val)
              }
              if (defaultValueIfEmpty) {
                if (defaultValueIfEmpty === 'undefined') {
                  if (!val) {
                    continue
                  }
                } else if (defaultValueIfEmpty === 'null') {
                  val = ''
                } else {
                  if (!val) {
                    val = defaultValueIfEmpty
                  }
                }
              }
              this.setValueByJsString(fb, o, targetInput.attr('name'), val)
            }
          }
        }
        // Ensure that validate is always at least an empty object so you can still add
        // validation rules after removing all of them
        if (o instanceof Array) {
          o.length && o.forEach((element) => {
            if (!element.validate) {
              element.validate = {}
            }
            return element
          })
        } else if (typeof o === 'object' && !o.validate) {
          o.validate = {}
        }
        fb.compiler.deepLoopOverJson(o, {
          'object': function (value, keyOrIndex, obj) {
            if (keyOrIndex === 'file' && value && value.all) {
              obj.file = {
                exact: value.exact,
                kind: value.all[value.selected]
              }
            }
            return true
          }
        })
        var oldCm = fb.compiler.getCompMainObject(fb.settingsTab.currentSettings)
        if (oldCm && oldCm['_import']) {
          var cm = fb.compiler.getCompMainObject(o)
          if (cm) {
            cm['_import'] = oldCm['_import']
          }
        }
        return o
      },
      setValueByJsString: function (fb, rootObj, name, value) {
        var o = rootObj
        var strArray = name.split(/(\[\d+\])|(\.)/g)
        var key = ''
        var arrayIndex = 0
        var arrayPointer = $.isArray(o)
        for (var i = 0; i < strArray.length; ++i) {
          if (strArray[i] && strArray[i] !== '') {
            if (/^[_a-zA-Z0-9]+.*$/.test(strArray[i])) { // key
              key = strArray[i]
              arrayPointer = false
            } else if (/^\[\d+\]$/.test(strArray[i])) { // array
              if (key !== '') {
                if (!$.isArray(o[key])) {
                  o[key] = []
                }
                o = o[key]
                arrayPointer = true
              }
              try {
                arrayIndex = parseInt(strArray[i].match(/(\d+)/g)[0])
              } catch (e) {}
            } else if (strArray[i] === '.') { // object
              if (arrayPointer) {
                if (typeof o[arrayIndex] !== 'object') {
                  o[arrayIndex] = {}
                }
                o = o[arrayIndex]
              } else {
                if (typeof o[key] !== 'object') {
                  o[key] = {}
                }
                o = o[key]
              }
              arrayPointer = false
            }
          }
        }
        if ((/^\d+\.?\d*$/.test(value)) && key !== 'name') { // is number
          value = parseFloat(value)
        }
        if (arrayPointer) {
          o[arrayIndex] = value
        } else {
          o[key] = value
        }
      }
    }
  }
  this.sizeOf = function (obj) {
    var size = 0;
    var key
    for (key in obj) {
      if (obj.hasOwnProperty(key)) size++
    }
    return size
  }
  this.uiNotification = function (obj) {
    var selector = '.inline-notification .inline-notification-holder'
    var notifEle = this.body.find(selector)
    if (!notifEle || notifEle.length === 0) {
      var template =
        '<div class="inline-notification">' +
        '   <div class="row">' +
        '       <div class="col-md-2"></div>' +
        '       <div class="col-md-8 inline-notification-holder">' +
        '           <div class="alert alert-{{status}} alert-dismissible" role="alert" style="width: 100%;z-index: 2;">' +
        '               <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span></button>' +
        '               {{message}}' +
        '           </div>' +
        '       </div>' +
        '       <div class="col-md-2"></div>' +
        '   </div>' +
        '</div>'
      var htmlStr = this.fb.compileTemplate(template, obj)
      this.body.prepend(htmlStr)
    } else {
      var alertTemplate =
        '           <div class="alert alert-{{status}} alert-dismissible" role="alert" style="width: 100%;z-index: 2;">' +
        '               <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span></button>' +
        '               {{message}}' +
        '           </div>'
      notifEle.html(this.fb.compileTemplate(alertTemplate, obj))
    }
  }
  this.compId = null
  this.currentSettings = null
  this.tabActive = false
  this.storeSettings = function (doneCallback) {
    var newSettings = this.transformer.html.toJsonSettings(this.fb, this.body)
    var main = this.fb.compiler.getCompMainObject(newSettings)
    delete main['_order']
    this.currentSettings = newSettings
    var _ = this
    this.fb.workspace.updateComponent(this.compId, this.currentSettings, function () {
      // done updating comp
      _.settingsSaved()
      if ($.isFunction(doneCallback)) {
        try {
          doneCallback()
        } catch (er) {
          console.log(er)
        }
      }
    })
  }
  this.settingsChanged = function (target, etype) {
    if (this.autoSave) {
      var _ = this.fb
      this.storeSettings(function () {
        if ((etype == 'click' || etype == 'change') || (target && target.length && target.hasClass('var-name'))) {
          _.workspace.connectionManager.update()
        }
      })
    } else {
      var saveBtn = this.el.find('.save-btn')
      saveBtn.addClass('btn-warning')
      saveBtn.removeClass('btn-secondary')
      saveBtn.addClass('hcb-changed')
    }
  }
  this.settingsSaved = function () {
    var saveBtn = this.el.find('.save-btn')
    saveBtn.addClass('btn-secondary')
    saveBtn.removeClass('btn-warning')
    saveBtn.removeClass('hcb-changed')
  }
  this.editSettings = function (compId, settingsJson, withoutFocus) {
    if (this.compId === compId) {
      this.showTab(false, withoutFocus)
    } else {
      this.renderSettings(compId, settingsJson, withoutFocus)
    }
  }
  this.updateSettings = function (compId, settingsJson) {
    if (this.tabActive) {
      this.renderSettings(compId, settingsJson, true)
    }
  }
  this.renderSettings = function (compId, settingsJson, withoutFocus) {
    this.fb.workspace.connectionManager.nameConnections = []
    this.fb.workspace.connectionManager.nameConExists = {}
    this.compId = compId
    this.currentSettings = settingsJson
    this.transformer.json.toHtmlInput(this.fb, settingsJson)
    this.settingsSaved()
    this.showTab(true, withoutFocus)
  }
  this.currentCompId = function (options) {
    if (options && options.release) {
      this.compId = null
    }
    return this.compId
  }
  this.componentRemoved = function (compId) {
    if (this.compId === compId) {
      this.clearBody()
    }
  }
  this.clearBody = function () {
    this.body.empty()
    this.compId = null
    this.settingsSaved()
    if (this.tabActive) {
      this.fb.componentsTab.showTab()
    }
  }
  this.showTab = function (reconstructed, withoutFocus) {
    if (!this.tabActive) {
      this.fb.el.find('.hcbuilder a.' + this.el.attr('id')).tab('show')
    }
    if (!withoutFocus) {
      this.focusInput()
    }
    this.fb.workspace.connectionManager.update(null, reconstructed)
  }
  this.tabShown = function () {
    this.tabActive = true
    this.fb.calcContainerHeight(this, true)
    if (!this.scrollEventAttached) {
      this.scrollEventAttached = true
      var fb = this.fb
      // this.innerBody.scroll(function(){
      //     fb.workspace.connectionManager.empty();
      // });
    }
    this.fb.activeTabObj = this
    this.fb.workspace.highlightComponent(this.compId)
    this.focusInput()
    this.fb.workspace.connectionManager.update()
  }
  this.tabHidden = function () {
    this.tabActive = false
    this.fb.workspace.unHighlightComponent(this.compId)
    this.fb.workspace.connectionManager.update()
  }
  this.focusInput = function () {}
  this.init(fb, jqEl)
}

/**
 * BuilderTab
 *
 * This Class is responsible for the construction of components
 */
var FT_BuilderTab = function (fb, jqEl) {
  this.init = function (fb, jqEl) {
    this.fb = fb
    this.el = jqEl
    var _ = this
    try {
      var splitLayoutMain = this.el.find('.split-layout-main')
      if (!splitLayoutMain.attr('id')) {
        splitLayoutMain.attr('id', this.fb.randomId())
      }
      this.layout = splitLayoutMain.layout({
        closable: false, // pane can open & close
        resizable: true, // when open, pane can be resized
        slidable: true, // when closed, pane can 'slide' open over other panes - closes on mouse-out
        livePaneResizing: true,
        stateManagement__enabled: true, // automatic cookie load & save enabled by default
        showDebugMessages: false, // log and/or display messages from debugging & testing code
        minSize: 0.1,
        size: 0.5,
        maxSize: 0.9,
        west: {},
        center: {}
      })
    } catch (exceptionWhenInitSplitPane) {
      console.log(exceptionWhenInitSplitPane)
      this.layout = {
        resizeAll: function () {}
      }
    }

    this.templateEditor = ace.edit('hcbuild-editor')
    this.templateEditor.session.setMode('ace/mode/handlebars')
    this.templateEditor.setTheme('ace/theme/tomorrow_night')
    // this.templateEditor.setAutoScrollEditorIntoView(true);
    // enable autocompletion and snippets
    this.aceEditorOptions = {
      maxLines: Infinity,
      autoScrollEditorIntoView: true,
      enableBasicAutocompletion: false,
      enableSnippets: true,
      wrapBehavioursEnabled: true,
      // autoScrollEditorIntoView: true,
      enableLiveAutocompletion: true
    }
    this.templateEditor.setOptions(this.aceEditorOptions)
    this.templateEditor.setShowPrintMargin(false)

    this.templateEditor.commands.addCommand({
      name: 'saveFile',
      bindKey: {
        win: 'Ctrl-S',
        mac: 'Command-S',
        sender: 'editor|cli'
      },
      exec: function (env, args, request) {
        try {
          fb.builderTab.storeComponent()
        } catch (ex) {
          console.log(ex)
        }
      }
    })
    this.templateEditor.on('input', function (a, b, c) {
      try {
        if (b.replacingContentItsNotAChange) {
          b.replacingContentItsNotAChange = false
          return
        } else {
          fb.builderTab.editorsValueChanged(b)
        }
        var template = b.getValue()
        fb.builderTab.render(template, fb.builderTab.settingsEditor.getValue())
      } catch (ex) {
        console.log(ex)
      }
    })
    this.templateEditor.renderer.on('afterRender', function () {
      _.aceEditorResize($(_.templateEditor.container))
    })
    // this.templateEditor.renderer.setScrollMargin(10, 10, 10, 10);
    // this.templateEditor.getSession().setUseWrapMode(false);
    this.templateEditor.$blockScrolling = Infinity

    this.settingsEditor = ace.edit('hcbuild-json-settings')
    this.settingsEditor.session.setMode('ace/mode/json')
    this.settingsEditor.setTheme('ace/theme/tomorrow_night')
    // this.settingsEditor.setAutoScrollEditorIntoView(true);
    // enable autocompletion and snippets
    this.settingsEditor.setOptions(this.aceEditorOptions)
    this.settingsEditor.setShowPrintMargin(false)
    // this.settingsEditor.renderer.setScrollMargin(10, 10, 10, 10);
    this.settingsEditor.$blockScrolling = Infinity
    this.settingsEditor.commands.addCommand({
      name: 'saveFile',
      bindKey: {
        win: 'Ctrl-S',
        mac: 'Command-S',
        sender: 'editor|cli'
      },
      exec: function (env, args, request) {
        try {
          fb.builderTab.storeComponent()
        } catch (ex) {
          console.log(ex)
        }
      }
    })
    this.settingsEditor.on('input', function (a, b, c) {
      try {
        if (b.replacingContentItsNotAChange) {
          b.replacingContentItsNotAChange = false
          return
        } else {
          fb.builderTab.editorsValueChanged(b)
        }
        fb.builderTab.render(fb.builderTab.templateEditor.getValue(), b.getValue())
      } catch (ex) {
        console.log(ex)
      }
    })
    this.settingsEditor.renderer.on('afterRender', function () {
      _.aceEditorResize($(_.settingsEditor.container))
    })
    this.el.find('.new-btn').click(function () {
      fb.builderTab.newComponent(true)
    })
    this.el.find('.new-copy-btn').click(function () {
      fb.builderTab.newComponent(false)
    })
    this.el.find('.build-btn').click(function () {
      fb.builderTab.storeComponent()
    })
    var fullscreenTarget = this.el.parent().find('.fullscreen-target')
    var fullscreenEvent = function (fullscreenBtn, changeState) {
      if (fullscreenTarget.hasClass('fullscreen')) {
        if (changeState) {
          var sticky = fullscreenTarget.parents('.fb-sticky-right')
          if (sticky.length) {
            if (fullscreenTarget._lastTrans) {
              sticky.css('transform', fullscreenTarget._lastTrans)
            }
          }
          fullscreenTarget.removeClass('fullscreen')
          fullscreenTarget.css({
            'height': ''
          })
          fullscreenTarget.find('.panel-body').css({
            'height': ''
          })
          fullscreenBtn.find('span.glyphicon').removeClass('glyphicon-resize-small')
          fullscreenBtn.find('span.glyphicon').addClass('glyphicon-fullscreen')
        } else {
          fullscreenTarget.height($(window).height())
          fullscreenTarget.find('.panel-body').height($(window).height() - 55)
        }
      } else {
        if (changeState) {
          var sticky = fullscreenTarget.parents('.fb-sticky-right')
          if (sticky.length) {
            fullscreenTarget._lastTrans = sticky.css('transform')
            sticky.css('transform', '')
          }
          fullscreenTarget.addClass('fullscreen')
          fullscreenTarget.height($(window).height())
          fullscreenTarget.find('.panel-body').height($(window).height() - 55)
          fullscreenBtn.find('span.glyphicon').removeClass('glyphicon-fullscreen')
          fullscreenBtn.find('span.glyphicon').addClass('glyphicon-resize-small')
        }
      }
      setTimeout(function () {
        _.layout.resizeAll()
        _.aceEditorResizeAll()
      }, 450)
      setTimeout(function () {
        _.layout.resizeAll()
        _.aceEditorResizeAll()
      }, 850)
      setTimeout(function () {
        _.layout.resizeAll()
        _.aceEditorResizeAll()
      }, 1250)
    }
    var fullscreenBtn = this.el.find('.fullscreen-btn')
    this.el.find('.fullscreen-btn').click(function () {
      fullscreenEvent($(this), true)
    })
    $(window).resize(function () {
      fullscreenEvent(fullscreenBtn, false)
    })
    this.el.addClass('new-entry')
  }
  this.getBasePath = function () {
    return '/static/js'
  }
  this.fb = null
  this.el = null
  this.compId = null
  this.dfsId = null
  this.component = null
  this.templateEditor = null
  this.settingsEditor = null
  this.possibleToStore = false
  this.updateDfsId = function (oldId, newId) {
    if (this.dfsId === oldId) {
      this.dfsId = newId
    }
  }
  this.render = function (template, settings, notChanged) {
    if (template || settings) {
      var html = fb.compileTemplate(template, settings)
      if (html) {
        try {
          var $html = $(html)
          $html.find('script').each(function () {
            eval($(this).text())
          })
          var $v = this.el.find('.hcbuild-viewer')
          $v.empty()
          $v.append($html)
          if (notChanged) {
            return
          }
          this.possibleToStore = true
          this.enableSaveBtn(true)
          return
        } catch (e) {
          console.log(e)
        }
      }
    }
    this.possibleToStore = false
    this.enableSaveBtn(false)
    this.el.find('.hcbuild-viewer').empty()
  }
  this.storeComponent = function () {
    if (!this.possibleToStore) {
      return
    }
    if (!this.dfsId) {
      this.dfsId = this.fb.randomId()
    }
    var settings = $.parseJSON(this.settingsEditor.getValue())
    var template = this.templateEditor.getValue()
    this.fb.compiler.cacheTemplate(this.dfsId, template)
    this.fb.componentsTab.addComponent(this.dfsId, {
      template: template,
      settings: settings
    })
    this.editMode()
    this.editorsValueSaved()
  }
  // {"someDfsId":{template:"", settings:{}}}
  this.editComponent = function (dfsId, component) {
    this.editMode()
    this.editorsValueSaved()
    this.dfsId = dfsId
    this.component = component
    this.templateEditor.replacingContentItsNotAChange = true
    this.templateEditor.setValue(component.template, -1)
    this.templateEditor.session.getUndoManager().reset()
    this.templateEditor.session.getUndoManager().markClean()
    this.templateEditor.selection.moveCursorToScreen(0, 0)
    this.settingsEditor.replacingContentItsNotAChange = true
    this.settingsEditor.setValue(JSON.stringify(this.cleanSettingsBeforeSetToEditor(component.settings), null, '\t'), -1)
    this.settingsEditor.session.getUndoManager().reset()
    this.settingsEditor.session.getUndoManager().markClean()
    this.settingsEditor.selection.moveCursorToScreen(0, 0)
    this.showTab()
    this.enableSaveBtn(false)
    this.possibleToStore = false
    this.render(component.template, component.settings, true)
  }
  this.cleanSettingsBeforeSetToEditor = function (settings) {
    var settingsCopy = this.fb.deepCopySettings(settings)
    if ($.isArray(settingsCopy)) {
      for (var i = 0; i < settingsCopy.length; ++i) {
        delete settingsCopy[i].id
      }
    } else {
      delete settingsCopy.id
    }
    return settingsCopy
  }
  this.newComponent = function (clean) {
    this.dfsId = fb.randomId()
    if (clean) {
      this.templateEditor.replacingContentItsNotAChange = true
      this.templateEditor.setValue('')
      this.templateEditor.selection.moveCursorToScreen(0, 0)
      this.settingsEditor.replacingContentItsNotAChange = true
      this.settingsEditor.setValue('')
      this.settingsEditor.selection.moveCursorToScreen(0, 0)
      this.newEntryMode()
      this.render()
    } else {
      this.newEntryPlusMode()
    }
  }
  this.editMode = function () {
    this.el.removeClass('new-entry')
    this.el.removeClass('new-entry-plus')
  }
  this.newEntryMode = function () {
    this.el.removeClass('new-entry-plus')
    this.el.addClass('new-entry')
  }
  this.newEntryPlusMode = function () {
    this.el.removeClass('new-entry')
    this.el.addClass('new-entry-plus')
  }
  this.editorsValueChanged = function (editor) {
    var saveBtn = this.el.find('.build-btn')
    saveBtn.addClass('btn-warning')
    saveBtn.removeClass('btn-secondary')
    saveBtn.addClass('hcb-changed')
    try {
      _.aceEditorResize($(editor.container))
    } catch (eee) {}
  }
  this.enableSaveBtn = function (enable) {
    var saveBtn = this.el.find('.build-btn')
    if (enable) {
      saveBtn.removeAttr('disabled')
    } else {
      saveBtn.prop('disabled', true)
    }
  }
  this.editorsValueSaved = function () {
    var saveBtn = this.el.find('.build-btn')
    saveBtn.addClass('btn-secondary')
    saveBtn.removeClass('btn-warning')
    saveBtn.removeClass('hcb-changed')
  }
  this.showTab = function () {
    this.fb.el.find('.hcbuilder a.' + this.el.attr('id')).tab('show')
  }
  this.tabShown = function () {
    var _ = this
    _.layout.resizeAll()
    _.aceEditorResizeAll()
    this.fb.activeTabObj = null
  }
  this.aceEditorResizeAll = function () {
    this.templateEditor.resize()
    this.settingsEditor.resize()
  }
  this.aceEditorResize = function (aceEditorEl) {
    var ace_content = aceEditorEl.find('.ace_content')
    var ace_gutter = aceEditorEl.find('.ace_gutter')
    var stretchEditorWidth = ace_gutter.width()
    var stretchEditorHeight = ace_gutter.height()
    stretchEditorWidth = stretchEditorWidth + ace_content.width()
    aceEditorEl.width(stretchEditorWidth)
    aceEditorEl.height(stretchEditorHeight)
  }
  this.tabHidden = function () {}
  this.init(fb, jqEl)
}

var formBuilderHtmlConstruct =
  '<div id="htmlFormWorkspace" class="col-sm-6 fb-workspace">' +
  '   ' +
  `<div class="row"><div class="col-sm-12"><ul class="nav nav-fb">
            <li class="nav-item">
            <a class="nav-link active ws-mode mode-workspace" href="#">Workspace</a>
            </li>
            <li class="nav-item">
            <a class="nav-link ws-mode mode-action" href="#">Action</a>
            </li>
            <li class="nav-item">
            <a class="nav-link ws-mode mode-test" href="#">Test</a>
            </li>
            </ul></div></div>` +
  '<div class="col-md-4 d-none">' +
  '           <label class="ws-mode ws-connections" style="float:right;">' +
  '               <input type="checkbox" checked="checked">' +
  '               <div class="slider"><span class="ws-mode-text">connection</span></div>' +
  '           </label>' +
  '       </div>' +
  '   ' +
  '   <div class="hcbuild-workspace-body wsbody container-fluid">' +
  '      <div class="connectionLayer">' +
  '       <svg class="dragline" height="0" width="0">' +
  '       </svg>' +
  '       <svg class="fbcm connections" height="0" width="0"></svg>' +
  '       <svg class="fbcm am-connections" height="0" width="0"></svg>' +
  '   </div><div class="am-point-layer am-point-main"></div><div class="ws-holder"></div></div>' +
  '   <div class="hcbuild-workspace-test-main wsbody" style="display:none;">' +
  '       <div class="panel panel-default">' +
  '           <div class="panel-body workspace-test-body"></div>' +
  '           <div class="panel-footer">' +
  '               <div class="workspace-test-btns row">' +
  '                   <div class="col-md-12">' +
  '                       <button type="button" class="submit-test btn btn-primary" style="float:right;margin-left: 5px;"><span class="fa fa-play text-white" aria-hidden="true"></span></button>' +
  '                       <button type="button" class="clear-form-data btn btn-primary" style="float:right;"><span class="fa fa-redo-alt text-white" aria-hidden="true"></span></button>' +
  '                   </div>' +
  '               </div>' +
  '           </div>' +
  '       </div>' +
  '   </div>' +
  '</div>' +
  '<div class="col-sm-6 hcbuilder">' +
  '   <div class="fb-sticky-right-nn"><div class="hcbuilder-main"><ul class="nav nav-fb">' +
  '       <li class="nav-item"><a class="htmlComponents nav-link active" data-toggle="tab" href="#htmlComponents">Components</a></li>' +
  '       <li class="nav-item"><a class="htmlComponentSettings nav-link" data-toggle="tab" href="#htmlComponentSettings">Properties</a></li>' +
  '       <li class="nav-item"><a class="htmlComponentBuilder nav-link" data-toggle="tab" href="#htmlComponentBuilder">Builder</a></li>' +
  '   </ul>' +
  '   <div class="tab-content">' +
  '       <div id="htmlComponents" class="htmlComponents tab-pane in active">' +
  '           <div class="panel panel-default">' +
  '               <div class="panel-heading">' +
  '                   <div class="input-group input-group-lg fb-comp-search-group">' +
  '                       <span class="input-group-addon" id="sizing-addon1"><span class="fa fa-search" aria-hidden="true"></span></span>' +
  '                       <input type="text" class="fb-comp-search-field form-control" placeholder="Search: textarea, checkbox, radio ..." aria-describedby="sizing-addon1">' +
  '                   </div>' +
  '               </div>' +
  '               <div class="panel-body fb-inner-body hcbuilder-components-body"><form><div class="hcbuild-comp-body"></div></form></div>' +
  '               <div class="panel-footer">' +
  '                   <button type="button" class="save-btn btn btn-secondary"><span class="fa fa-save" aria-hidden="true"></span></button>' +
  '               </div>' +
  '           </div>' +
  '       </div>' +
  '       <div id="htmlComponentSettings" class="htmlComponentSettings tab-pane">' +
  '           <div class="panel panel-default">' +
  '               <div class="panel-body fb-inner-body hcbuilder-settings-body"></div>' +
  '               <div class="panel-footer d-flex flex-row">' +
  '                   <button type="button" class="save-btn btn btn-secondary"><span class="fa fa-save" aria-hidden="true"></span></button>' +
  '                   <div class="switch-main"><span>auto save</span><label class="switch"> <input class="auto-save-btn" type="checkbox"> <div class="slider"></div> </label></div>' +
  '                   <div class="clearfix"></div>' +
  '               </div>' +
  '           </div>' +
  '       </div>' +
  '       <div id="htmlComponentBuilder" class="htmlComponentBuilder tab-pane fullscreen-target">' +
  '           <div class="panel panel-default" >' +
  '               <div class="panel-body fb-inner-body split-layout-main" >' +
  '                                       <div class="ui-layout-west ">' +
  '                                           <div class="hcbuild-container ace-workaround">' +
  '                                               <div id="hcbuild-editor" class="hcbuild-editor" ></div>' +
  '                                           </div>' +
  '                                       </div>' +
  '                                       <div class="ui-layout-center">' +
  '                                           <div class="hcbuild-container ace-workaround">' +
  '                                               <div id="hcbuild-json-settings" class="hcbuild-json-settings" ></div>' +
  '                                           </div>' +
  '                                       </div>' +
  '                               <div class="ui-layout-south hcbuild-container">' +
  '                                   <div class="hcbuild-container"><div class="hcbuild-viewer" ></div></div>' +
  '                               </div>' +
  '               </div>' +
  '               <div class="panel-footer">' +
  '                   <button type="button" class="new-btn btn btn-primary" ><span class="fa fa-plus" aria-hidden="true"></span></button>' +
  '                   <button type="button" class="build-btn btn btn-secondary"><span class="fa fa-save" aria-hidden="true"></span></button>' +
  '                   <button type="button" class="fullscreen-btn btn btn-secondary"><span class="fa fa-expand" aria-hidden="true"></span></button>' +
  '               </div>' +
  '           </div>' +
  '       </div>' +
  '   </div></div></div>' +
  '</div>'

if (!String.prototype.startsWith) {
  (function () {
    'use strict' // needed to support `apply`/`call` with `undefined`/`null`
    var defineProperty = (function () {
      // IE 8 only supports `Object.defineProperty` on DOM elements
      try {
        var object = {}
        var $defineProperty = Object.defineProperty
        var result = $defineProperty(object, object, object) && $defineProperty
      } catch (error) {
        console.log(error)
      }
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
String.prototype.toUniqueHash = function () {
  var hash = 0;
  var i;
  var chr
  if (this.length === 0) return hash
  for (i = 0; i < this.length; i++) {
    chr = this.charCodeAt(i)
    hash = ((hash << 5) - hash) + chr
    hash |= 0 // Convert to 32bit integer
  }
  return (hash + '').replace('-', 'm')
}
$.fn.indexInParent = function () {
  var elem = this
  var i = 0
  while ((elem = elem.previousSibling) != null) ++i
  return i
}
$.fn.nextParentWithClass = function (attrClass) {
  var maxIndex = 800
  var index = 1
  var parent = this.parent()
  while (!parent.hasClass(attrClass)) {
    parent = parent.parent()
      ++index
    if (maxIndex < index) {
      return $([])
    }
  }
  if (parent.hasClass(attrClass)) {
    return parent
  }
  return $([])
}
$.fn.insertAt = function (index, element) {
  var lastIndex = this.children().size()
  if (index < 0) {
    index = Math.max(0, lastIndex + 1 + index)
  }
  this.append(element)
  if (index < lastIndex) {
    this.children().eq(index).before(this.children().last())
  }
  return this
}
$.fn.toString = function () {
  return this.prop('outerHTML')
}
$.fn.findVisibleInputOrRefElement = function () {
  if ((this.attr('type') === 'radio' || this.attr('type') === 'checkbox')) {
    this.data('isBool', true)
    var findLabel = this.parent()
    var tagName = findLabel[0].tagName.toUpperCase()
    if (tagName === 'LABEL' && findLabel.attr('for') === this.attr('id')) {
      return findLabel
    } else {
      findLabel = findLabel.find("label[for='" + this.attr('id') + "']")
    }
    if (!findLabel.length) {
      if (this.is(':visible')) {
        return this
      }
    }
    return findLabel
  } else {
    return this
  }
}
$.fn.positionOfUnderlying = function (ctxEl) {
  var start = this.offset();
  var mainLayer = ctxEl.offset()
  return {
    x: start.left + 0.5 * this.outerWidth() - mainLayer.left,
    y: start.top + 0.5 * this.outerHeight() - mainLayer.top
  }
}
$.fn.attrFBName = function () {
  var name = this.attr('fb_name')
  if (!name) {
    name = this.attr('name')
  }
  return name
}
$.fn.attrID = function () {
  var startId = this.attr('id')
  if (!startId) {
    startId = Math.random().toString(36).substring(2, 18) + Math.random().toString(36).substring(2, 18)
    this.attr('id', startId)
  }
  return startId
}
$.fn.renameAttr = function (oldName, newName) {
  if (oldName && newName) {
    this.each(function () {
      var _this = $(this)
      _this.attr(newName, _this.attr(oldName)).removeAttr(oldName)
    })
  }
  return this
}

var fbDefaultComponents = {
  'HC1': {
    'template': `
        <div class=" form-group">
            <div class="col-md-12">
              <h1>{{label}}</h1>
              <span class="help-block">{{help}}</span> 
            </div>
        </div>
        `,
    'settings': {
      'help': 'help text',
      'label': 'Headline'
    }
  },
  'HC2': {
    'template': `
      <div class=" form-group">
        <div class="col-md-12">
          <label class="control-label" for="{{id}}">{{label}}</label>
          <div class="field-parent">
            <input id="{{id}}" name="{{name}}" type="text" placeholder="{{escapeForAttr placeholder}}" class="form-control input-md" autocomplete="{{escapeForAttr autocomplete}}">
             <span class="help-block">{{help}}</span> 
          </div>
        </div>
        <div class="clearfix"></div>
      </div>
      `,
    'settings': {
      'autocomplete': 'on',
      'help': 'help text',
      'label': 'Simple field',
      'name': '',
      'placeholder': 'Placeholder',
      'validate': {
        'required': true
      }
    }
  },
  'HC3': {
    'template': `
     <div class="form-group">
    \t   <div class="col-md-12 field-parent">
    \t       <label class="control-label" >{{label}}</label>
    \t\t       {{#each values}}
    \t\t       <div class="fancy-radio">
                        <input type="radio" id="{{@root.id}}{{@index}}" name="{{@root.name}}" value="{{value}}" {{#ifEq @root.defaultValue value}}checked="checked"{{/ifEq}}/>
                        <label for="{{@root.id}}{{@index}}"><span><i></i></span>{{label}}</label>
                    </div>
                     {{#if help}}
                    <div class="row">
                        <div class="col-sm-12 col-md-12 col-lg-12">
        \t                <span style="margin-left:34px;" class="help-block">{{help}}</span>
        \t            </div>
                     </div>
    \t         {{/if}}
    \t\t       {{/each}}
    \t\t       <span class="help-block">{{help}}</span>
    \t   </div>
    \t   <div class="clearfix"></div>
       </div>
        `,
    'settings': {
      'defaultValue': '',
      'help': 'help text',
      'label': 'Radio buttons',
      'name': '',
      'validate': {
        'required': true
      },
      'values': [{
        'help': '',
        'label': '...',
        'value': 'val'
      }]
    }
  },
  'HC5': {
    'template': `
        <div class="form-group">
          <label class="col-sm-12 col-md-12 col-lg-12 control-label">{{label}}</label>
      \t  <div class="col-sm-12 col-md-12 col-lg-12 field-parent">
      \t       {{#each values}}
      \t\t     <div class="fancy-checkbox {{#enumEq @root.orientation "horizontal"}}float-left{{/enumEq}}">
      \t\t            <input type="checkbox" id="{{@root.id}}{{@index}}" name="{{@root.name}}" value="{{value}}" />
                          <label for="{{@root.id}}{{@index}}"><span><i></i></span>{{label}}</label>
                   </div>
                   {{#if help}}
                   {{#enumEq @root.orientation "horizontal"}}
          \t           <span style="margin-right:10px;float:left;margin-bottom:0;" class="help-block">{{help}}</span>
                   {{else}}
                      <div class="row">
                          <div class="col-sm-12 col-md-12 col-lg-12">
          \t                <span style="margin-left:34px;" class="help-block">{{help}}</span>
          \t            </div>
                       </div>
                   {{/enumEq}}
      \t         {{/if}}
      \t       {{/each}}
      \t       <div class="row">
      \t           <div class="col-sm-12 col-md-12 col-lg-12">
      \t               <span class="help-block">{{help}}</span>
      \t           </div>
      \t       </div>
      \t  </div>
      \t   <div class="clearfix"></div>
      </div>
        `,
    'settings': {
      'help': 'help text',
      'label': 'Checkbox',
      'name': '',
      'orientation': {
        'all': [
          'horizontal',
          'vertical'
        ],
        'selected': 1
      },
      'validate': {
        'required': true
      },
      'values': [{
        'help': '',
        'label': '...',
        'value': '1'
      }]
    }
  },
  'HC7': {
    'template': `
    <div class="form-group">
      <div class="col-md-12">
        <label class="control-label" for="{{id}}">{{label}}</label>
        <div class="field-parent">
          <textarea id="{{id}}" name="{{name}}" placeholder="{{escapeForAttr placeholder}}" class="form-control"></textarea>
           <span class="help-block">{{help}}</span> 
        </div>
      </div>
      <div class="clearfix"></div>
    </div>
        `,
    'settings': {
      'help': 'help text',
      'label': 'Multiline text',
      'name': '',
      'placeholder': 'Placeholder',
      'validate': {
        'required': true
      }
    }
  },
  'HC8': {
    'template': `
        <div class="form-group">
            <div class="col-md-12">
                <label class="control-label" for="{{id}}">{{label}}</label>
                <div class="field-parent">
                    <select id="{{id}}" name="{{name}}" class="form-control" aria-invalid="false">
                      {{#if placeholder}}<option value="" selected>{{placeholder}}</option>{{/if}}
                      {{#each this.values}}
                      <option value="{{value}}" {{#ifEq @root.defaultValue value}}selected{{/ifEq}}>{{label}}</option>
                      {{/each}}
                    </select>
                    <span class="help-block">{{help}}</span> 
                </div>
            </div>
            <div class="clearfix"></div>
        </div>
        `,
    'settings': {
      'defaultValue': '',
      'help': 'help text',
      'label': 'Drop-down list',
      'name': '',
      'placeholder': '',
      'validate': {
        'required': true
      },
      'values': [{
          'label': 'Select',
          'value': 'val1'
        },
        {
          'label': 'Value 2',
          'value': 'val2'
        }
      ]
    }
  },
  'HC9': {
    'template': `
      <div class="form-group">
        <div class="col-md-12">
          <label class="control-label" for="dfId{{id}}">{{label}}</label>
          <div class="field-parent">
            <div class="input-group date-field">
                <input placeholder="{{escapeForAttr Placeholder}}"vtype="text" class="text-field simple-date-field" name="{{name}}" id="dfId{{id}}" size="0" aria-invalid="false">
                <span class="input-group-addon">
                    <span class="fa fa-calendar"></span>
                </span>
                <script type="text/javascript">
                    $(function(){
                        var dateInputField = $("#dfId{{id}}");
                        FTG.createDatepicker(dateInputField, {
                            clickEventElement:dateInputField.parent().find(".input-group-addon"),
                            datePattern:'{{validate.datePattern}}'});
                    });
                </script>
            </div>
            <span class="help-block">{{help}}</span> 
          </div>
        </div>
        <div class="clearfix"></div>
      </div>
        `,
    'settings': {
      'Placeholder': 'Placeholder',
      'help': 'help text',
      'label': 'Date',
      'name': '',
      'validate': {
        'datePattern': 'dd.MM.yyyy HH:mm:ss',
        'required': true
      }
    }
  },
  'HC10': {
    'template': `
      <div class="array form-group">
      \t    <div class="col-md-12">
      \t\t    <div class="row">
      \t\t        <div class="col-md-12"><label class="control-label">{{this.0.label}}</label></div>
      \t        </div>
      \t    </div>
      \t    <div class="col-md-12" id="main{{this.0.id}}">
      \t        {{#eachCount this.0.initialRows}}
      \t\t\t  <div class="{{#enumEq @root.0.columnClass "auto"}}array-row{{else}}array-row row{{/enumEq}}">               
      \t\t\t    {{#each @root}}
      \t\t\t        <div {{#enumEq @root.0.columnClass "auto"}}class="array-item field-parent form-group" style="{{#unless @last}}padding-right:10px;{{/unless}}float:left;width:{{percentOfOneFrom @root.length}}%;"{{else}} class="array-item field-parent form-group {{enumName @root.0.columnClass}}"{{/enumEq}}>
      \t\t\t            {{#enumEq type "textarea"}}
      \t\t\t\t        <textarea id="{{id}}{{@countIndex}}" addBtnRef="#btn{{@root.0.id}}" name="{{name}}" nindex="{{@countIndex}}" placeholder="{{escapeForAttr placeholder}}" class="form-control"></textarea>
      \t\t\t            {{else}}
          \t\t\t            {{#enumEq type "date"}}
          \t\t\t                <div class="input-group date-field">
                                          <input  id="{{id}}{{@countIndex}}" data-datePattern="{{validate.datePattern}}" addBtnRef="#btn{{@root.0.id}}" name="{{name}}" nindex="{{@countIndex}}" placeholder="{{escapeForAttr placeholder}}" type="text" class="text-field simple-date-field" size="0" aria-invalid="false">
                                          <span class="input-group-addon">
                                              <span class="fa fa-calendar"></span>
                                          </span>
                                      </div>
                                  {{else}}
                                      {{#enumEq type "dropdown"}}
                                          <select id="select{{id}}{{@countIndex}}" addBtnRef="#btn{{@root.0.id}}" name="{{name}}" nindex="{{@countIndex}}" class="form-control" aria-invalid="false">
                                            {{#each dropdownValues}}
                                            <option value="{{value}}">{{label}}</option>
                                            {{/each}}
                                          </select>
                                      {{else}}
                  \t\t\t\t        <input id="{{id}}{{@countIndex}}" addBtnRef="#btn{{@root.0.id}}" name="{{name}}" nindex="{{@countIndex}}" type="text" placeholder="{{escapeForAttr placeholder}}" class="form-control input-md">
              \t\t\t            {{/enumEq}}
          \t\t\t            {{/enumEq}}
      \t\t\t            {{/enumEq}}
      \t\t\t        </div>
      \t\t\t     {{/each}}
      \t               <div class="clearfix"></div>           
      \t           </div>
      \t\t    {{/eachCount}}
      \t    </div>
      \t    <div class="col-md-12">
      \t\t    <div class="row">
      \t\t        <div class="col-md-12 form-group">
      \t\t\t        <button id="btn{{@root.0.id}}" type="button" class="btn btn-success pull-right" onclick="add{{this.0.id}}(this);"><span class="fa fa-plus" aria-hidden="true"></span></button>
      \t\t\t        <button id="btnDel{{@root.0.id}}" disabled type="button" class="btn btn-success pull-right" onclick="remove{{this.0.id}}(this);" style="margin:0 15px"><span class="fa fa-minus" aria-hidden="true"></span></button>
      \t\t        </div>
      \t        </div>
      \t    </div>
      \t    <div class="col-md-12">
      \t        <span class="help-block">{{this.0.help}}</span>
      \t    </div>
      \t\t    <div class="clearfix"></div>
      \t    <script type="text/javascript">
                      $(function(){
                          var mainEl = $("#main{{this.0.id}}");
      \t\t            initDate{{this.0.id}}(mainEl);
      \t\t            var btnDel = $("#btnDel{{@root.0.id}}");
                          if(mainEl.find(".array-row").length>{{this.0.initialRows}}){
      \t\t\t            btnDel.removeAttr("disabled");
      \t\t            }else{
      \t\t                btnDel.attr("disabled", "disabled");
      \t\t            }
                      });
      \t\t        function add{{this.0.id}}(ce){
      \t\t            var _this = $(ce);
      \t\t            var array = _this.parent().parent().parent().parent();
      \t\t            var arrayItem = array.find(".array-row:last");
      \t\t            var newArrayItem = $(arrayItem.prop("outerHTML"));
      \t\t            newArrayItem.find("input,textarea,select").each(function(){
      \t\t                var newIn = $(this);
      \t\t                newIn.attr("id",randomId{{this.0.id}}());
      \t\t            });
      \t\t            newArrayItem.find("[nindex]").each(function(e){
      \t\t                var _nItem = $(this);
      \t\t                try{
      \t\t                    _nItem.attr("nindex", parseInt(_nItem.attr("nindex"))+1);
      \t\t                }catch(incEx){console.log(incEx);}
      \t\t            });
      \t\t            newArrayItem.cleanFieldErrors();
      \t\t            arrayItem.after(newArrayItem);
      \t\t            func{{this.0.id}}2(arrayItem.parent());
      \t\t            console.log("added");
      \t\t            initDate{{this.0.id}}(newArrayItem);
      \t\t            _this.parents("form").trigger("formFieldsAdded",[newArrayItem]);
      \t\t            $("#btnDel{{@root.0.id}}").removeAttr("disabled");
      \t\t        }
      \t\t        function remove{{this.0.id}}(ce){
      \t\t            var _this = $(ce);
      \t\t            var array = _this.parent().parent().parent().parent();
      \t\t            var arrayItem = array.find(".array-row:last");
      \t\t            if(array.find(".array-row").length>{{this.0.initialRows}}){
      \t\t                arrayItem.remove();
      \t\t                if(array.find(".array-row").length=={{this.0.initialRows}}){
      \t\t                    _this.attr("disabled", "disabled");
      \t\t                }
      \t\t                var firstRow = array.find(".array-row:first-child");
      \t\t                firstRow.find("input,textarea,select").change();
      \t\t            }else{
      \t\t                _this.attr("disabled", "disabled");
      \t\t            }
      \t\t        }
      \t\t        function func{{this.0.id}}2(arrayParent){
      \t\t            {{#enumEq @root.0.columnClass "auto"}}
      \t\t\t            var arrayItems = arrayParent.children(".array-item");
      \t\t\t            var width = 100.0/arrayItems.length;
      \t\t\t            arrayItems.css({width:((width-0.2)+"%")});
      \t\t\t        {{/enumEq}}
      \t\t        }
      \t\t        function initDate{{this.0.id}}(jqMainEl){
      \t\t            var dateFields = jqMainEl.find(".simple-date-field");
      \t\t            dateFields.each(function(){
      \t\t                var dField = $(this);
      \t\t                dField.removeClass("hasDatepicker");
      \t\t                date{{this.0.id}}(dField, dField.attr("data-datePattern"));
      \t\t            });
      \t\t        }
      \t\t        function date{{this.0.id}}(jqEl, pattern){
                          var dateInputField = jqEl;
                          jqEl.removeClass("hasDatepicker");
                          FTG.createDatepicker(dateInputField, {
                              clickEventElement:dateInputField.parent().find(".input-group-addon"),
                              datePattern:pattern
                          });
      \t\t        }
      \t            function randomId{{this.0.id}}(){
                          return Math.random().toString(36).substring(2, 18) + Math.random().toString(36).substring(2, 18);
                      };
      \t    </script>
      </div>
        `,
    'settings': [{
        'columnClass': {
          'all': [
            'auto',
            'col-md-2',
            'col-md-4',
            'col-md-6',
            'col-md-8',
            'col-md-12'
          ],
          'selected': 0
        },
        'dropdownValues': [{
          'label': 'Label',
          'value': 'Value'
        }],
        'help': 'help text',
        'initialRows': 1,
        'label': 'Dynamic List',
        'name': '',
        'placeholder': 'Placeholder',
        'type': {
          'all': [
            'input',
            'textarea',
            'dropdown',
            'date'
          ],
          'selected': 0
        },
        'validate': {
          'required': true
        }
      },
      {
        'dropdownValues': [{
          'label': 'Label',
          'value': 'Value'
        }],
        'name': '',
        'placeholder': 'Placeholder',
        'type': {
          'all': [
            'input',
            'textarea',
            'dropdown',
            'date'
          ],
          'selected': 3
        },
        'validate': {
          'required': true
        }
      }
    ]
  },
  'HC11': {
    'template': `
      <div class="form-group">
      \t    {{#each this}}    <label class="col-md-12 control-label" for="{{id}}">{{label}}</label>
      \t    <div class="col-md-12 field-parent">
      \t\t        <input id="{{id}}" name="{{name}}" type="text" placeholder="{{escapeForAttr placeholder}}" class="form-control input-md" autocomplete="{{escapeForAttr autocomplete}}">
      \t\t       <p class="help-block">{{help}}</p>
      \t    </div>
      \t    {{/each}}
      </div>
        `,
    'settings': [{
      'autocomplete': 'on',
      'help': 'help text',
      'label': 'Grouped simple fields',
      'name': '',
      'placeholder': 'Placeholder',
      'validate': {
        'required': true
      }
    }]
  },
  'HC12': {
    'template': `
      <div>
          <div class="col-md-12">
              {{#eachCount rows}}
                  <div class="row">
                      {{#each @root.column}}
                          <div class="{{enumName @this @root._md}}">{{dropHere}}</div>
                      {{/each}}
                      <div class="clearfix"></div>
                  </div>
              {{/eachCount}}
              
          </div>
          <div class="clearfix"></div>
      </div>
        `,
    'settings': {
      '_md': {
        '100%': 'col-md-12',
        '16%': 'col-md-2',
        '25%': 'col-md-3',
        '33%': 'col-md-4',
        '42%': 'col-md-5',
        '50%': 'col-md-6',
        '58%': 'col-md-7',
        '67%': 'col-md-8',
        '75%': 'col-md-9',
        '8%': 'col-md-1',
        '84%': 'col-md-10',
        '92%': 'col-md-11',
        'auto': 'auto'
      },
      '_searchHelper': 'group,parent,grp,container,box',
      'access': {},
      'column': [{
        'all': [
          'auto',
          '8%',
          '16%',
          '25%',
          '33%',
          '42%',
          '50%',
          '58%',
          '67%',
          '75%',
          '84%',
          '92%',
          '100%'
        ],
        'selected': 12
      }],
      'rows': 1
    }
  }
}

export default FT_FormBuilder
