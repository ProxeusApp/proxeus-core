import './typeahead'

/*! https://mths.be/startswith v0.2.0 by @mathias */
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
        if (string.charCodeAt(start + index) !=
                    searchString.charCodeAt(index)) {
          return false
        }
      }
      return true
    }
    if (defineProperty) {
      defineProperty(String.prototype, 'startsWith', {
        value: startsWith,
        configurable: true,
        writable: true
      })
    } else {
      String.prototype.startsWith = startsWith
    }
  }())
}

function VisFlowchart (wfData, pOptions) {
  this.randomId = function () {
    this.rndIndex++
    if (this.rndIndex % 1000 === 0) {
      this.rndHash = this.rnd()
    }
    return this.rndIndex + '_' + this.rndHash
  }
  this.rnd = function () {
    return Math.random().toString(36).substring(2, 18)
  }
  this.rndIndex = 0
  this.rndHash = this.rnd()
  this.container = pOptions.html.$container
  this.s = {
    startId: 'start111',
    conNode: 'conNode',
    conEdge: 'conEdge'
  }
  this.options = {
    readOnly: false,
    nodes: {
      start: {
        connections: {
          from: [
            {
              node: {
                shape: 'dot',
                size: 10,
                color: {
                  background: '#8dd5ff',
                  border: '#7b7b7b',
                  highlight: {
                    background: '#8dd5ff',
                    border: '#7b7b7b'
                  }
                },
                borderWidth: 1,
                borderWidthSelected: 1
              },
              edge: {
                arrowStrikethrough: false,
                arrows: 'to',
                width: 2,
                hoverWidth: 4,
                selectionWidth: 4,
                color: {
                  color: '#45bbff',
                  highlight: '#45bbff',
                  hover: '#45bbff'
                }
              }
            }
          ],
          to: 0
        },
        shape: 'icon',
        icon: {
          face: 'FontAwesome',
          code: '\uf152',
          size: 80,
          color: '#45bbff'
        }
      },
      con: {}
    },
    events: {
      hoverIn: null/* function(){} */,
      hoverOut: null/* function(){} */,
      insert: null/* function(){} */,
      remove: null/* function(){} */,
      click: null/* function(){} */,
      dblclick: null/* function(){} */
    },
    node: {
      connections: {
        from: [
          {
            node: {
              shape: 'dot',
              size: 10,
              color: {
                background: '#45bbff',
                border: '#7b7b7b',
                highlight: { background: '#8dd5ff', border: '#7b7b7b' },
                hover: { background: '#8dd5ff', border: '#7b7b7b' }
              },
              borderWidth: 1,
              borderWidthSelected: 1
            },
            edge: {
              arrowStrikethrough: false,
              arrows: 'to',
              width: 1,
              hoverWidth: 2,
              selectionWidth: 2,
              color: { color: '#45bbff', highlight: '#45bbff', hover: '#45bbff' }
            }
          }],
        to: Infinity,
        space: 1.1
      },
      shadow: {
        enabled: true,
        color: 'rgba(0,0,0,0.5)',
        size: 10,
        x: 5,
        y: 5
      }
    }
  }
  $.extend(true, this.options, pOptions)
  var nNodeOptions = {}
  var nn
  this.readOnly = this.options.readOnly
  for (var nk in this.options.nodes) {
    if (this.options.nodes.hasOwnProperty(nk)) {
      if (nk === 'con') {
        nn = this.options.nodes[nk]
      } else {
        nn = {}
        $.extend(true, nn, this.options.node, this.options.nodes[nk],
          { connections: null, events: null })
        delete nn.connections
        delete nn.events
      }
      nNodeOptions[nk] = nn
    }
  }

  this.networkOptions = {
    groups: nNodeOptions,
    interaction: {
      dragNodes: true,
      tooltipDelay: 200,
      hover: true,
      navigationButtons: true,
      selectable: true,
      selectConnectedEdges: true,
      keyboard: {
        enabled: true,
        bindToWindow: false
      }
    },
    manipulation: {
      enabled: false
    },
    edges: {
      font: {
        size: 10
      },
      smooth: {
        enabled: true,
        forceDirection: 'none'
      }
    },
    nodes: {
      font: {
        size: 18
      }
    },
    physics: {
      // enabled:false,
      stabilization: {
        enabled: true,
        iterations: 50000,
        updateInterval: 100,
        onlyDynamicEdges: false,
        fit: false
      },
      barnesHut: {
        gravitationalConstant: -850,
        centralGravity: 0.00001,
        springLength: 200,
        springConstant: 0.0885,
        damping: 0.6,
        avoidOverlap: 0
      },
      repulsion: {
        centralGravity: 0.002,
        springLength: 200,
        springConstant: 0.08,
        nodeDistance: 100,
        damping: 0.8
      },
      timestep: 0.5,
      adaptiveTimestep: true,
      maxVelocity: 50,
      minVelocity: 0.4,
      solver: 'barnesHut'
    }
  }
  var nData
  if ($.isArray(wfData)) {
    nData = this.workflowDataToVis(wfData)
  } else {
    nData = this.setFlowData(wfData)
  }
  this._initUI()
  this._initNetwork(nData, this.networkOptions)
  this._initEvents()
  this.stabilize()
  this.fit()
  window.network2 = this.network
}

VisFlowchart.prototype.workflowDataToVis = function (wfData) {
  var dataResult = { nodes: [], edges: [] }
  var startCon = this.getConnectionSettings('start')
  var newNode = {
    id: this.s.startId,
    label: 'start',
    physics: false,
    group: 'start',
    connections: startCon
  }
  startCon.from = []
  this.mergeConnectionFromSettings(newNode, { value: newNode.id }, 0)
  this.mergeNodeEvents(newNode)
  this.createConNode(dataResult, newNode, 0)
  var incrementToCount = function (id, toNodesMap) {
    if (typeof toNodesMap[id] === 'number') {
      toNodesMap[id]++
    } else {
      toNodesMap[id] = 1
    }
  }

  var startNode = newNode
  dataResult.nodes.push(newNode)
  if (wfData) {
    var w
    var cnd
    if (wfData.length > 0) {
      var toNodesMap = {}
      var nodeMap = {}
      dataResult.edges.push(
        this.createEdge('start', { from: newNode.id, to: wfData[0].id }, 0))
      incrementToCount(wfData[0].id, toNodesMap)
      try {
        startNode.x = wfData[0].p.start.x
        startNode.y = wfData[0].p.start.y
      } catch (dntcr) {
      }
      newNode.connections.from[0].hide(this)

      for (var i = 0; i < wfData.length; ++i) {
        w = wfData[i]
        if (w.type === 'step') {
          w.type = 'form'
        }
        newNode = this._createNode(w)
        nodeMap[newNode.id] = newNode
        dataResult.nodes.push(newNode)
        newNode._data = w.data
        if (w.type === 'condition') {
          newNode._cases = []
          for (var i2 = 0; i2 < w.cases.length; i2++) {
            w.cases[i2].name = w.cases[i2].label
            newNode._cases.push(
              { name: w.cases[i2].label, value: w.cases[i2].value })
            this.mergeConnectionFromSettings(newNode, w.cases[i2], i2)
            this.createConNode(dataResult, newNode, i2)
          }
        } else {
          this.mergeConnectionFromSettings(newNode, { value: newNode.id }, 0)
          newNode.uriName = w.uriName
          this.createConNode(dataResult, newNode, 0)
        }
        this.mergeNodeEvents(newNode)
        if (w.connectedTo && w.connectedTo.length) {
          for (var c = 0; c < w.connectedTo.length; ++c) {
            cnd = w.connectedTo[c]
            incrementToCount(cnd.targetId, toNodesMap)
            if (w.type === 'condition') {
              dataResult.edges.push(
                this.createEdge(w.type, { from: w.id, to: cnd.targetId }, c, cnd))
              for (var b = 0; b < newNode.connections.from.length; ++b) {
                if (newNode.connections.from[b].value == cnd.value) {
                  newNode.connections.from[b].hide(this)
                }
              }
            } else {
              dataResult.edges.push(
                this.createEdge(w.type, { from: w.id, to: cnd.targetId }, 0))
              newNode.connections.from[0].hide(this)
            }
          }
        }
        this._fireEvent(newNode, 'insert')
      }
      for (var toNodeId in toNodesMap) {
        if (toNodeId && toNodesMap.hasOwnProperty(toNodeId) && nodeMap[toNodeId] && nodeMap[toNodeId].connections) {
          nodeMap[toNodeId].connections.toCount = toNodesMap[toNodeId]
        }
      }
    }
  }
  return dataResult
}
VisFlowchart.prototype.fit = function (a) {
  if (this.network) {
    this.network.fit(a)
  }
}
VisFlowchart.prototype.stabilize = function (a) {
  if (this.network) {
    if (a === undefined) {
      a = 5000000
    }
    this.network.stabilize(a)
  }
}
VisFlowchart.prototype._initUI = function () {
  var jqVisNav = this.container.find('.vis-navigation')
  var jqFullscreenBtn = $(
    '<div class="vis-button vis-fullscreen"><i class="fa fa-window-maximize" aria-hidden="true"></i></div>')
  var _ = this
  if (!_.ui) {
    _.ui = {}
  }
  jqFullscreenBtn.click(function () {
    var o = _.options.html
    var fullscreenEl = o.$container
    var mainEl = fullscreenEl.parent()
    var fullscreenWrapperEl = null
    if (o && o.$wrapper) {
      fullscreenWrapperEl = o.$wrapper
    }
    if (_.ui.fullscreen) {
      mainEl.css({ height: '', width: '' })
      if (fullscreenWrapperEl && fullscreenWrapperEl.length) {
        fullscreenWrapperEl.removeAttr('style')
      }
      fullscreenEl.removeAttr('style')
      _.ui.fullscreen = false
      $(this).find('i').css({ color: 'black' })
    } else {
      _.ui.fullscreen = true
      if (fullscreenWrapperEl && fullscreenWrapperEl.length) {
        fullscreenWrapperEl.css({
          position: 'fixed',
          'z-index': 990000,
          width: '100%',
          height: '100%',
          left: 0,
          top: 0,
          bottom: 0,
          right: 0
        })
        fullscreenEl.css({
          width: '100%',
          height: '100%'
        })
      } else {
        fullscreenEl.css({
          position: 'fixed',
          'z-index': 990000,
          width: '100%',
          height: '100%',
          left: 0,
          top: 0,
          bottom: 0,
          right: 0
        })
      }
      $(this).find('i').css({ color: 'green' })
      mainEl.css({ height: '100%', width: '100%' })
    }
  })
  jqVisNav.append(jqFullscreenBtn)
}
VisFlowchart.prototype.nodes = null
VisFlowchart.prototype.edges = null
VisFlowchart.prototype.network = null
VisFlowchart.prototype.seed = 2

VisFlowchart.prototype.destroy = function () {
  if (this.network !== null) {
    try {
      this.container.unbind()
    } catch (dontCare1) {
    }
    try {
      this.container.find('.vis-button').unbind()
    } catch (dontCare2) {
    }
    try {
      this.container.find('canvas').unbind()
    } catch (dontCare3) {
    }
    this.network.destroy()
    this.network = null
  }
}

VisFlowchart.prototype._initNetwork = function (dataP, optionsP) {
  this.network = new window.vis.Network(this.options.html.$container[0], dataP,
    optionsP)
}
VisFlowchart.prototype.getData = function () {
  var data = { edges: [], nodes: [] }
  var id, n, i
  for (i = 0; i < this.network.body.nodeIndices.length; ++i) {
    id = this.network.body.nodeIndices[i]
    if (!id.startsWith(this.s.conNode)) {
      n = this.getNodeById(id)
      data.nodes.push(n)
    }
  }
  for (i = 0; i < this.network.body.edgeIndices.length; ++i) {
    id = this.network.body.edgeIndices[i]
    if (!id.startsWith(this.s.conEdge)) {
      n = this.getEdgeById(id)
      data.edges.push(n)
    }
  }
  return data
}

VisFlowchart.prototype.createConFromNode = function (node, i, from) {
  var n = this.getMergedNode(node, i, from)
  n.physics = false
  n.id = this.s.conNode + '_' + node.id + (from.value ? '_' + from.value : '') +
        '_' + i
  n.group = 'con'
  n._isCon = true
  n.hidden = this.readOnly
  n._from_ = { parentId: node.id, index: i }
  if (typeof node.x === 'number') {
    n.x = node.x
  }
  if (typeof node.y === 'number') {
    n.y = node.y
  }
  var e = this.createEdge(node.group, { from: node.id, to: n.id }, i, from)
  if (e.arrows) {
    delete e.arrows
  }
  e.id = this.s.conEdge + '_' + e.id
  e._isCon = true
  e.hidden = this.readOnly
  if (from.notAvailable) {
    n.hidden = true
    e.hidden = true
    e.physics = false
  }
  from.ne = {
    n: n,
    e: e,
    hide: function (fc, f) {
      f.notAvailable = true
      if (fc.network) {
        fc.network.body.data.nodes.update([{ id: this.n.id, hidden: true }])
        fc.network.body.data.edges.update(
          [{ id: this.e.id, hidden: true, physics: false }])
      } else {
        this.n.hidden = true
        this.e.hidden = true
        this.e.physics = false
      }
    },
    show: function (fc, f) {
      f.notAvailable = false
      if (fc.network) {
        fc.network.body.data.nodes.update([{ id: this.n.id, hidden: false }])
        fc.network.body.data.edges.update(
          [{ id: this.e.id, hidden: false, physics: true }])
      } else {
        this.n.hidden = false
        this.e.hidden = false
        this.e.physics = true
      }
    }
  }
  return from.ne
}

VisFlowchart.prototype.getMergedNode = function (node, i, from) {
  var defaultSettings
  if (typeof i === 'number') {
    try {
      defaultSettings = this.options.nodes[node.group].connections.from[i].node
    } catch (dontCare) {
    }
  }
  if (!defaultSettings) {
    try {
      defaultSettings = this.options.nodes[node.group].connections.from[0].node
    } catch (dontCare) {
    }
  }
  var newSettings = undefined
  if (from && from.node) {
    newSettings = from.node
  }
  var newNode = {}
  $.extend(true, newNode, this.options.node.connections.from[0].node,
    defaultSettings, newSettings)
  return newNode
}

// storePositions is causing labels to reset there positions sometimes
// this methods is suppoest to fix that but so far it doesn't
VisFlowchart.prototype._fixEdgeLabelPositions = function (edgesIds) {
  if (this.lastFixTime === undefined) {
    this.lastFixTime = 0
  }
  if (new Date() - this.lastFixTime > 100) {
    this.lastFixTime = new Date()
    var update = []
    if (edgesIds === undefined || edgesIds.length === 0) {
      edgesIds = this.network.body.edgeIndices
    }
    for (var i = 0; i < edgesIds.length; i++) {
      if (this.getEdgeById(edgesIds[i]).label) {
        update.push({ id: edgesIds[i], allowedToMoveX: false, allowedToMoveY: false })
      }
    }
    this.network.body.data.edges.update(update)
  }
}
VisFlowchart.prototype.getFlowData = function () {
  this.network.storePositions()
  // because labels are of the conditions are messed up sometimes after calling storePositions
  // this._fixEdgeLabelPositions();
  this.network.stopSimulation()
  var output = { start: {}, nodes: {} }
  var targetNode = null
  var layoutNode = null
  var
    ioNode = null
  var startP = null
  var startNode = this.network.body.nodes[this.s.startId]
  var id, i, c, e, conn
  if (startNode && startNode.edges && startNode.edges.length) {
    layoutNode = startNode
    startP = { x: layoutNode.x, y: layoutNode.y }
    output.start.p = startP
    if (startNode.edges && startNode.edges.length > 0) {
      for (i = 0; i < startNode.edges.length; ++i) {
        if (!startNode.edges[i].toId.startsWith(this.s.conNode)) {
          output.start.node = startNode.edges[i].toId
        }
      }
    }
  } else {
    // throw "start not connected!";
  }

  for (i = 0; i < this.network.body.nodeIndices.length; ++i) {
    id = this.network.body.nodeIndices[i]
    if (id === this.s.startId || id.startsWith(this.s.conNode)) {
      continue
    }
    targetNode = this.getNodeById(id)
    layoutNode = this.network.body.nodes[targetNode.id]
    output.nodes[id] =
            ioNode = {
              id: targetNode.id,
              name: targetNode.label,
              detail: targetNode.detail ? targetNode.detail : '',
              type: targetNode.group,
              p: { x: layoutNode.x, y: layoutNode.y },
              conns: []
            }
    if (targetNode._cases) {
      ioNode.cases = targetNode._cases
    }
    if (targetNode._data) {
      ioNode.data = targetNode._data
    }
    for (c = 0; c < layoutNode.edges.length; ++c) {
      e = layoutNode.edges[c].options
      if (!e.id.startsWith(this.s.conEdge) && e.to !== targetNode.id) {
        e = this.getEdgeById(e.id)
        conn = { id: e.to/*, "data":{"label": e.label, "value": e.cvalue} */ }
        if (ioNode.id !== e.cvalue) {
          conn.value = e.cvalue
        }
        ioNode.conns.push(conn)
      }
    }
  }
  return output
}
VisFlowchart.prototype.reload = function (wfData) {
  this.destroy()
  this._initNetwork(this.setFlowData(wfData), this.networkOptions)
  this._initUI()
  this._initEvents()
}
VisFlowchart.prototype.setFlowData = function (wfData) {
  var dataResult = { nodes: [], edges: [] }
  var startCon = this.getConnectionSettings('start')
  var newNode = {
    id: this.s.startId,
    label: 'start',
    physics: false,
    group: 'start',
    connections: startCon
  }
  startCon.from = []
  this.mergeConnectionFromSettings(newNode, { value: newNode.id }, 0)
  this.mergeNodeEvents(newNode)
  this.createConNode(dataResult, newNode, 0)
  var incrementToCount = function (id, toNodesMap) {
    if (typeof toNodesMap[id] === 'number') {
      toNodesMap[id]++
    } else {
      toNodesMap[id] = 1
    }
  }

  var convCase = function (cid, cs) {
    var c = { value: cid }
    if (cs.name) {
      c.name = cs.name
      if (cs.detail) {
        c.title = cs.detail
      }
    }
    return c
  }

  var startNode = newNode
  dataResult.nodes.push(newNode)
  if (wfData) {
    var toNodesMap = {}
    var nodeMap = {}
    if (wfData.start) {
      if (wfData.start.p) {
        startNode.x = wfData.start.p.x
        startNode.y = wfData.start.p.y
      }
      if (wfData.start.node && wfData.nodes &&
                wfData.nodes[wfData.start.node]) {
        dataResult.edges.push(
          this.createEdge('start', { from: newNode.id, to: wfData.start.node },
            0))
        incrementToCount(wfData.start.node, toNodesMap)
        newNode.connections.from[0].hide(this)
      }
    }
    var ioNode; var i = 0
    if (wfData.nodes) {
      for (var id in wfData.nodes) {
        if (wfData.nodes.hasOwnProperty(id)) {
          ioNode = wfData.nodes[id]
          if (ioNode.type === 'step') {
            ioNode.type = 'form'
          }
          newNode = this._createNode(ioNode, id)
          nodeMap[newNode.id] = newNode
          dataResult.nodes.push(newNode)
          i = 0
          if (ioNode.cases) {
            var cond
            for (; i < ioNode.cases.length; i++) {
              cond = ioNode.cases[i]
              cond = convCase(cond.value, cond)
              this.mergeConnectionFromSettings(newNode, cond, i)
              this.createConNode(dataResult, newNode, i)
            }
            newNode._cases = ioNode.cases
          }
          newNode._data = ioNode.data
          try {
            if (i == 0 && ioNode.type &&
                            this.options.nodes[ioNode.type].connections.from.length > 0) {
              this.mergeConnectionFromSettings(newNode, { value: newNode.id }, 0)
              this.createConNode(dataResult, newNode, 0)
            }
          } catch (dntc) {
            console.error(dntc)
          }
          this.mergeNodeEvents(newNode)
          if (ioNode.conns && ioNode.conns.length) {
            var cnd, conn
            for (var c = 0; c < ioNode.conns.length; ++c) {
              cnd = undefined
              conn = ioNode.conns[c]
              if (!conn) {
                continue
              }
              if (ioNode.cases) {
                for (var cv = 0; cv < ioNode.cases.length; cv++) {
                  if (ioNode.cases[cv] &&
                                        ioNode.cases[cv].value === conn.value) {
                    cnd = convCase(conn.value, ioNode.cases[cv])
                    break
                  }
                }
              }
              incrementToCount(conn.id, toNodesMap)
              dataResult.edges.push(
                this.createEdge(ioNode.type, { from: id, to: conn.id }, c, cnd))
              if (!newNode.connections.fromInfinity) {
                for (var b = 0; b < newNode.connections.from.length; ++b) {
                  if ((!cnd && newNode.connections.from.length == 1) ||
                                        (cnd && newNode.connections.from[b].value === cnd.value)) {
                    newNode.connections.from[b].hide(this)
                  }
                }
              }
            }
          }
          this._fireEvent(newNode, 'insert')
        }
      }
      for (var toNodeId in toNodesMap) {
        if (toNodeId && toNodesMap.hasOwnProperty(toNodeId) && nodeMap[toNodeId] && nodeMap[toNodeId].connections) {
          nodeMap[toNodeId].connections.toCount = toNodesMap[toNodeId]
        }
      }
    }
  }
  return dataResult
}

VisFlowchart.prototype._storeSeed = function (wfData) {
  try {
    if (!wfData[0].p) {
      wfData[0].p = {}
    }
    wfData[0].p.seed = this.network.getSeed()
  } catch (dontCare) {
  }
}

VisFlowchart.prototype._initEvents = function () {
  var _ = this
  this.network.on('beforeDrawing', function (ev) {
    // _vfThis.beforeDrawing.call(this, ev);
    var tn, t, r, f, flen, n, id, pos
    for (var i = 0; i < this.body.nodeIndices.length; ++i) {
      id = this.body.nodeIndices[i]
      if (id.startsWith(_.s.conNode)) {
        continue
      }
      tn = this.body.data.nodes.get(id)
      if (tn && tn.connections && tn.connections.from) {
        pos = this.body.nodes[id]
        if (pos) {
          flen = tn.connections.from.length
          for (var c = 0; c < flen; ++c) {
            f = tn.connections.from[c]
            if (f.ne && f.ne.n.id !== _.draggingCon) {
              n = this.body.nodes[f.ne.n.id]
              if (n) {
                t = 2 * Math.PI * (1 + c) / flen
                if (!tn.connections.radius) {
                  r = pos.shape.width / tn.connections.space
                } else {
                  r = tn.connections.radius
                }
                n.x = Math.round(pos.x + r * Math.cos(t))
                n.y = Math.round(pos.y + r * Math.sin(t))
              }
            }
          }
        }
      }
    }
  })
  this.conPossible = {
    fc: this,
    valid: {
      enabled: true,
      color: 'rgba(0, 206, 6, 0.79)',
      size: 12,
      x: 3,
      y: 3
    },
    invalid: {
      enabled: true,
      color: 'rgba(255, 0, 0, 0.62)',
      size: 12,
      x: 3,
      y: 3
    },
    lastNodeId: '',
    updated: {},
    resetted: false,
    Yes: function (n, u) {
      try {
        this.update(n, this.valid, u)
      } catch (e) {
        console.error(e)
      }
    },
    No: function (n, u) {
      try {
        this.update(n, this.invalid, u)
      } catch (e) {
        console.error(e)
      }
    },
    Reset: function () {
      if (!this.resetted) {
        for (var key in this.updated) {
          if (!this.updated.hasOwnProperty(key)) continue
          this.clean(key, this.updated[key])
        }
        this.resetted = true
        this.lastNodeId = null
      }
    },
    clean: function (key, up) {
      delete this.updated[key]
      up.n.shadow = up.s
      up.u.update([up.n.id])
    },
    update: function (n, s, u) {
      var up
      if (this.lastNodeId && this.lastNodeId != n.id &&
                (up = this.updated[this.lastNodeId])) {
        this.clean(this.lastNodeId, up)
        this.lastNodeId = null
      }
      if (!this.updated[n.id]) {
        this.lastNodeId = n.id
        up = { n: n, s: n.shadow, u: u }
        n.shadow = s
        this.updated[n.id] = up
        this.resetted = false
        up.u.update([up.n.id])
      }
    }
  }
  this.checkStart = function (ev) {
    if (ev) {
      for (var i = 0; i < ev.nodes.length; ++i) {
        if (ev.nodes[i] === this.s.startId) {
          return false
        }
      }
    }
    return true
  }
  this.network.on('select', function (ev) {
    if (ev.nodes.length == 0 && ev.edges.length == 0) {
      _.currentSelected = null
      _.hideDeleteOnUI(ev)
    } else {
      if (_.checkStart(ev)) {
        if (_._canDelete(ev)) {
          _.currentSelected = ev
          _.showDeleteOnUI(ev)
        }
      } else {
        _.currentSelected = null
        _.hideDeleteOnUI(ev)
      }
    }
  })
  this.network.on('click', function (ev) {
    if (ev.nodes.length == 1) {
      _._fireEvent(_.getNodeById(ev.nodes[0]), 'click')
    }
  })
  this.draggingCon = null
  this.network.on('dragStart', function (ev) {
    if (ev.nodes.length == 1) {
      var nId = ev.nodes[0]
      var draggingNode = this.body.data.nodes.get(nId)
      if (draggingNode._isCon) {
        _.draggingCon = nId
      }
    }
  })
  this.network.on('dragging', function (ev) {
    if (ev.nodes.length == 1) {
      var nId = ev.nodes[0]
      var draggingNode = this.body.nodes[nId]
      if (draggingNode.options.group === 'con') {
        _.getNodeToConnect(nId, draggingNode)
      }
    }
  })
  this.getNodeToConnect = function (nId, draggingNode) {
    var s = draggingNode.shape
    var pointerObj = {
      left: s.left,
      top: s.top,
      right: s.left + s.width,
      bottom: s.top + s.height
    }
    // get the overlapping node but NOT the temporary node;
    var node = undefined
    for (var i2 = 0; i2 < this.network.body.nodeIndices.length; i2++) {
      var nodeId = this.network.body.nodeIndices[i2]
      if (nodeId !== nId) {
        node = this.network.body.nodes[nodeId]
        if (node.options.group !== 'con' &&
                    node.isOverlappingWith(pointerObj)) {
          node = this.network.body.nodes[nodeId]
          break
        }
        node = undefined
      }
    }
    if (node !== undefined &&
            node.options.id !== draggingNode.options._from_.parentId) {
      var nodeTo = node.options
      if (node.options.id !== nId) {
        var nodeFrom = this.network.body.data.nodes.get(nId)
        var syncedNode = this.network.body.data.nodes.get(node.options.id)
        if (syncedNode && syncedNode.connections &&
                    syncedNode.connections.to > 0 &&
                    syncedNode.connections.toCount < syncedNode.connections.to) {
          if (nodeFrom.connections && nodeFrom.connections.fromInfinity &&
                        this.network.body.nodes[nodeTo.id].edges.some(
                          function (element, index, array) {
                            return element.fromId === nodeFrom.id
                          })) {
            _.conPossible.No(nodeTo, this.network.body.data.nodes)
          } else {
            _.conPossible.Yes(nodeTo, this.network.body.data.nodes)
            return nodeTo
          }
        } else {
          _.conPossible.No(nodeTo, this.network.body.data.nodes)
        }
      } else {
        _.conPossible.No(nodeTo, this.network.body.data.nodes)
      }
    } else {
      _.conPossible.Reset()
    }
    return null
  }
  this.network.on('dragEnd', function (ev) {
    if (ev.nodes.length == 1) {
      var nId = ev.nodes[0]
      var draggingNode = this.body.data.nodes.get(nId)
      if (draggingNode._isCon) {
        _.network.unselectAll()
        var nodeTo = _.getNodeToConnect(nId, this.body.nodes[nId])
        if (nodeTo) {
          var nodeFrom = this.body.data.nodes.get(draggingNode._from_.parentId)
          var targetFrom = nodeFrom.connections.from[draggingNode._from_.index]
          var newEdge = _.createEdge(nodeFrom.group, {
            from: draggingNode._from_.parentId,
            to: nodeTo.id
          }, draggingNode._from_.index, targetFrom)
          this.body.data.edges.getDataSet().add(newEdge)
          this.body.data.nodes.get(nodeTo.id).connections.toCount++
          if (!nodeFrom.connections.fromInfinity) {
            targetFrom.hide(_)
          }
        }
        _.conPossible.Reset()
        _.draggingCon = null
      }
    }
  })
  this.network.on('hoverNode', function (ev, b, c, d) {
    if (ev.node) {
      var node = _.getNodeById(ev.node)
      if (node) {
        _._fireEvent(node, 'hoverIn')
      }
    }
  })
  this.network.on('doubleClick', function (ev, b, c, d) {
    if (ev.nodes && ev.nodes.length == 1) {
      var node = _.getNodeById(ev.nodes[0])
      if (node) {
        _._fireEvent(node, 'dblclick')
      }
    }
  })
  this.network.on('blurNode', function (ev, b, c, d) {
    if (ev.node) {
      var node = _.getNodeById(ev.node)
      if (node) {
        _._fireEvent(node, 'hoverOut')
      }
    }
  })

  function showModalDialog (title, text, inputValue, buttonFunc, suggestFunc) {
    $('#single-input-dialog #modal-title').text(title)
    $('#single-input-dialog #modal-text').text(text)
    $('#single-input-dialog #modal-button').unbind().click(buttonFunc)
    $('#single-input-dialog').modal({ show: true, backdrop: true })

    const input = $('#single-input-dialog #modal-input')
    input.typeahead('destroy')
    input.val(inputValue)

    if (suggestFunc) {
      input.text('')
      input.on('click', function () {
        input.typeahead('lookup').focus()
      })
      input.typeahead(
        {
          minLength: 0,
          source: suggestFunc,
          items: 20
        })
    }
  }

  function modalDialogGetVal () {
    return $('#single-input-dialog #modal-input').val()
  }

  this.network.on('oncontext', function (ev, b, c, d) {
    if (ev.nodes && ev.nodes.length == 1) {
      const node = _.getNodeById(ev.nodes[0])
      ev.event.preventDefault()
      $('.custom-menu li').unbind().click(function () {
        if (!node._data) {
          node._data = {}
        }
        switch ($(this).attr('data-action')) {
          case 'restrictAccess':
            showModalDialog(
              'Restrict access',
              'Define read access for documents produced by the template',
              node._data.readAccess,
              function () {
                node._data.readAccess = modalDialogGetVal()
              })
            break
          case 'contextDefineUser':
            showModalDialog(
              'Define new user',
              'Set public address of the user to which you want to change the execution',
              node._data.expectUser,
              function () {
                node._data.expectUser = modalDialogGetVal()
              })
            break
          case 'connectorFunc':
            const varConnectorsList = function (text, process) {
              $.get('/api/admin/connector/' + node.id + '/functions', function (data) {
                const l = data.reduce((res, x) => {
                  if (x.includes(text)) {
                    res.push(x)
                  }
                  return res
                }, [])
                process(l)
              })
            }
            showModalDialog(
              'Select connector function',
              'Choose from available remote functions',
              node._data.connectorFunc,
              function () {
                node._data.connectorFunc = modalDialogGetVal()
              },
              varConnectorsList
            )
            break
        }
        $('.custom-menu').hide(100)
      })
      let count = 0
      $('.custom-menu li').each(function () {
        if (!$(this).attr('for-nodes').includes(node.group)) {
          $(this).hide()
        } else {
          count++
          $(this).show()
        }
      })
      if (count > 0) {
        $('.custom-menu').finish().toggle(100)
        $('.custom-menu').css({
          top: ev.event.pageY + 'px',
          left: ev.event.pageX + 'px'
        })
      }
    }
  })
  return this
}
VisFlowchart.prototype._fireEvent = function (node, t) {
  if (!node._isCon) {
    if (this.options.events[t]) {
      this.options.events[t].apply(this, [node])
    }
    if (node.events && node.events[t]) {
      node.events[t].apply(this, [node])
    }
  }
}

VisFlowchart.prototype.calcPointsCirc = function (pos, i, max, from) {
  this.calcConPoints(pos, i, max, this.network.body.nodes[from.ne.n.id])
}

VisFlowchart.prototype.calcConPoints = function (pos, i, max, n) {
  var t = 2 * Math.PI * i / max
  var r = pos.shape.width / this.options.conSpace
  n.x = Math.round(pos.x + r * Math.cos(t))
  n.y = Math.round(pos.y + r * Math.sin(t))
}

VisFlowchart.prototype.isInCircle = function (pos, circlesObj) {
  var t
  for (var key in circlesObj) {
    if (!circlesObj.hasOwnProperty(key)) continue
    t = circlesObj[key]
    if (pos.x < t.right && pos.x > t.left && pos.y > t.top &&
            pos.y < t.bottom) {
      return t
    }
  }
  return null
}

VisFlowchart.prototype.circle = function (id, ctx, ref) {
  ctx.circle(ref.x, ref.y, ref.radius)
  ref.id = id
  ref.top = ref.y - ref.radius
  ref.left = ref.x - ref.radius
  ref.right = ref.x + ref.radius
  ref.bottom = ref.y + ref.radius
}

VisFlowchart.prototype.getNodeById = function (id) {
  return this.network.body.data.nodes.get(id)
}
VisFlowchart.prototype.getEdgeById = function (id) {
  return this.network.body.data.edges.get(id)
}
VisFlowchart.prototype._canDelete = function (ev) {
  var n = null
  var i
  for (i = 0; i < ev.nodes.length; ++i) {
    n = this.getNodeById(ev.nodes[i])
    if (n._isCon) {
      return false
    }
  }
  if (ev.nodes.length == 1) {
    return true
  }
  if (ev.edges.length == 1) {
    n = this.getEdgeById(ev.edges[0])
    if (n._isCon) {
      return false
    }
  }
  return true
}
VisFlowchart.prototype._deleteElements = function (ev) {
  var conNodes = []
  var i, e
  if (ev.edges) {
    for (i = 0; i < ev.edges.length; ++i) {
      e = this.getEdgeById(ev.edges[i])
      if (e._isCon) {
        conNodes.push(e.to)
      } else {
        var nodeTo = this.getNodeById(e.to)
        if (nodeTo) {
          this.network.body.data.nodes.get(nodeTo.id).connections.toCount--
        }
        var nodeFrom = this.getNodeById(e.from)
        if (nodeFrom) {
          console.log(nodeFrom)
          try {
            for (var ca = 0; ca < nodeFrom.connections.from.length; ++ca) {
              if (nodeFrom.connections.from[ca].value === e.cvalue) {
                if (!nodeFrom.connections.fromInfinity) {
                  nodeFrom.connections.from[ca].show(this)
                }
              }
            }
          } catch (ee) {
            console.error(ee)
          }
        }
      }
    }
  }
  if (ev.nodes) {
    for (i = 0; i < ev.nodes.length; ++i) {
      this._fireEvent(this.getNodeById(ev.nodes[i]), 'remove')
    }
  }
  if (ev.conNodes && ev.conNodes.length > 0) {
    this.network.body.data.nodes.getDataSet().remove(ev.conNodes)
  }

  if (ev.conEdges && ev.conEdges.length > 0) {
    this.network.body.data.edges.getDataSet().remove(ev.conEdges)
  }

  if (conNodes.length > 0) {
    this.network.body.data.nodes.getDataSet().remove(conNodes)
  }
  this.network.body.data.edges.getDataSet().remove(ev.edges)
  this.network.body.data.nodes.getDataSet().remove(ev.nodes)
  this.network.startSimulation()
}

VisFlowchart.prototype.addNode = function (data) {
  var newNode = this._createNode(data)
  if (data.x !== undefined && data.y !== undefined) {
    var coords = this.network.DOMtoCanvas({ x: data.x, y: data.y })
    newNode.x = coords.x
    newNode.y = coords.y
  }
  if (data.cases) {
    newNode._cases = data.cases
    for (var i2 = 0; i2 < data.cases.length; i2++) {
      this.mergeConnectionFromSettings(newNode, data.cases[i2], i2)
    }
  } else {
    this.mergeConnectionFromSettings(newNode, { value: newNode.id }, 0)
  }
  newNode._data = data.data
  if (!newNode._data) {
    newNode._data = {}
  }
  this.mergeNodeEvents(newNode)
  newNode.uriName = data.uriName
  for (var i = 0; i < newNode.connections.from.length; i++) {
    var ne = this.createConFromNode(newNode, i, newNode.connections.from[i])
    this.network.body.data.nodes.getDataSet().add(ne.n)
    this.network.body.data.edges.getDataSet().add(ne.e)
  }
  this.network.body.data.nodes.getDataSet().add(newNode)
  this._fireEvent(newNode, 'insert')
}

VisFlowchart.prototype.updateNode = function (data) {
  try {
    var tnode = this.network.body.nodes[data.id]

    if (data._cases) {
      var cf; var c; var leave; var toRemove = []; var toLeave = []; var updateEdges = []
      var visRemove = { conNodes: [], conEdges: [], edges: [] }
      var findRelatedEdgeId = function (fc, node, caseValue) {
        var e
        for (var i = 0; i < node.edges.length; i++) {
          if (node.edges[i].fromId == node.id) {
            e = fc.getEdgeById(node.edges[i].id)
            if (e && !e._isCon && e.cvalue === caseValue) {
              return e.id
            }
          }
        }
        return null
      }
      for (var f = 0; f < data.connections.from.length; f++) {
        cf = data.connections.from[f]
        if (cf) {
          leave = null
          for (var i = 0; i < data._cases.length; i++) {
            c = data._cases[i]
            if (c && c.value === cf.value) {
              leave = cf
              cf.name = c.name
              toLeave.push(i)
            }
          }
          if (leave) {
            updateEdges.push({ id: cf.ne.e.id, label: cf.name })
          } else {
            var e = findRelatedEdgeId(this, tnode, cf.value)
            if (e) {
              visRemove.edges.push(e)
            }
            visRemove.conEdges.push(cf.ne.e.id)
            visRemove.conNodes.push(cf.ne.e.to)
            toRemove.push(f)
          }
        }
      }

      for (var r = toRemove.length - 1; r >= 0; r--) {
        data.connections.from.splice(toRemove[r], 1)
      }
      var shouldAppendCase = function (cases, i) {
        for (var ci = 0; ci < cases.length; ci++) {
          if (i == cases[ci]) {
            return false
          }
        }
        return true
      }
      var fromIndex = 0
      for (var i = 0; i < data._cases.length; i++) {
        if (shouldAppendCase(toLeave, i)) {
          fromIndex = data.connections.from.length
          this.mergeConnectionFromSettings(data, data._cases[i], fromIndex)
          var ne = this.createConFromNode(data, fromIndex,
            data.connections.from[fromIndex])
          this.network.body.data.nodes.getDataSet().add(ne.n)
          this.network.body.data.edges.getDataSet().add(ne.e)
          ne.n.x = tnode.x
          ne.n.y = tnode.y
        }
      }
      this._deleteElements(visRemove)
    }

    this.network.body.data.nodes.update([
      {
        id: data.id,
        _cases: data._cases,
        _data: data._data,
        label: data.label,
        title: this.createTooltip(data, data, this.options.showIDOnTooltip)
      }])
    if (updateEdges.length) {
      this.network.body.data.edges.update(updateEdges)
    }
  } catch (e) {
    console.log(e)
  }
}

VisFlowchart.prototype._createNode = function (w, id) {
  var con = this.getConnectionSettings(w.type)
  con.from = []
  con.toCount = 0
  if (w.description) {
    w.detail = w.description
  }
  var n = {
    id: w.id,
    label: w.name,
    detail: w.detail,
    group: w.type,
    connections: con
  }
  if (!n.id) {
    n.id = id
    if (!n.id) {
      n.id = this.randomId()
    }
  }
  n.title = this.createTooltip(w, n, this.options.showIDOnTooltip)
  if (w.p) {
    n.x = w.p.x
    n.y = w.p.y
  }
  return n
}

VisFlowchart.prototype.mergeConnectionFromSettings = function (newNode, settings, i) {
  newNode.connections.from.push({})
  var nodeSpecificFromStyles
  try {
    nodeSpecificFromStyles = this.options.nodes[newNode.group].connections.from[i]
  } catch (dontCare1) {
  }
  if (!nodeSpecificFromStyles) {
    try {
      nodeSpecificFromStyles = this.options.nodes[newNode.group].connections.from[0]
    } catch (dontCare1) {
    }
  }
  $.extend(true, newNode.connections.from[i],
    this.options.node.connections.from[0], nodeSpecificFromStyles, settings, {
      hide: function (fc) {
        if (this.ne && fc) {
          this.ne.hide(fc, this)
        }
      },
      show: function (fc) {
        if (this.ne && fc) {
          this.ne.show(fc, this)
        }
      }
    })
}
VisFlowchart.prototype.getConnectionSettings = function (group) {
  var nodeSpecificFromStyles
  try {
    nodeSpecificFromStyles = this.options.nodes[group].connections
  } catch (dontCare1) {
  }
  return $.extend(true, {}, this.options.node.connections,
    nodeSpecificFromStyles, { from: null })
}
VisFlowchart.prototype.mergeNodeEvents = function (newNode) {
  try {
    newNode.events = this.options.nodes[newNode.group].events
  } catch (dontCare1) {
  }
}

VisFlowchart.prototype.createConNode = function (dataResult, newNode, i) {
  var ne = this.createConFromNode(newNode, i, newNode.connections.from[i])
  dataResult.nodes.push(ne.n)
  dataResult.edges.push(ne.e)
  ne.n.x = newNode.x
  ne.n.y = newNode.y
}

VisFlowchart.prototype.createEdge = function (grp, con, i, more) {
  var newEdge = this.createMergedEdge(grp, i, more)
  newEdge.from = con.from
  newEdge.to = con.to
  newEdge.id = this.randomId()
  if (more) {
    if (more.label) {
      newEdge.label = more.label
    } else {
      newEdge.label = more.name
    }
    newEdge.cvalue = more.value
    if (more.title) {
      newEdge.title = more.title
    }
  } else {
    newEdge.cvalue = con.from
  }
  return newEdge
}

VisFlowchart.prototype.createMergedEdge = function (grp, i, more) {
  var defaultSettings
  if (typeof i === 'number') {
    try {
      defaultSettings = this.options.nodes[grp].connections.from[i].edge
    } catch (settingsWithGivenIndexDoesNotExist) {
    }
  }
  if (!defaultSettings) {
    try {
      defaultSettings = this.options.nodes[grp].connections.from[0].edge
    } catch (settingsWithTypeDoesNotExist) {
    }
  }
  var newSettings = undefined
  if (more && more.edge) {
    newSettings = more.edge
  }
  var newEdge = {}
  $.extend(true, newEdge, this.options.node.connections.from[0].edge, defaultSettings, newSettings)
  return newEdge
}

VisFlowchart.prototype.createTooltip = function (data, n, showID) {
  var id, name, detail, cases
  id = data.id ? data.id : n.id
  if (data.name) {
    name = data.name
  }
  if (!name && n.label) {
    name = n.label
  }
  if (!name) {
    name = ''
  }
  if (data.detail) {
    detail = data.detail
  }
  if (!detail && n.detail) {
    detail = n.detail
  }
  if (!detail) {
    detail = ''
  }
  if (data.cases) {
    cases = data.cases
  }
  if (!cases && n._cases) {
    cases = n._cases
  }
  return '<div class="wf-node-tooltip">' +
        (showID ? ('<p>' + (id) + '</p>') : '') + '<strong class="wfnt-name">' +
        (name) + '</strong><br/><span class="wfnt-desc">' + (detail) + '</span>' +
        this.createTooltipConditionCases(cases) + '</div>'
}

VisFlowchart.prototype.createTooltipConditionCases = function (cases) {
  if (cases && cases.length) {
    var html = '<ul class="wfnt-cases">'; var name
    for (var i = 0; i < cases.length; ++i) {
      name = cases[i].name
      if (!name) {
        name = cases[i].label
      }
      html += '<li class="wfnt-case" ><span class="wfnt-label">' + name +
                '</span></li>'
    }
    html += '</ul>'
    return html
  }
  return ''
}

VisFlowchart.prototype.createEdgeTooltip = function (label) {
  return '<div class="wf-node-tooltip"><span class="wfnt-desc">' + label +
        '</span></div>'
}

VisFlowchart.prototype.absoluteCenter = function (_this, parent) {
  _this.css('left', parent.width() / 2 - _this.width() / 2)
}

VisFlowchart.prototype.showDeleteOnUI = function (ev) {
  if (this.readOnly) {
    return
  }
  var jqVisNetworkEl = this.container.find('.vis-network')
  var jqWfUiEl = jqVisNetworkEl.find('.wfvis-ui-div .wfvis-ui-delete')
  if (!jqWfUiEl.length) {
    var html = '<div class="wfvis-ui-div" style="position: absolute;left: 48%;top: 0;z-index:100;"><div class="wfvis-ui-delete" style="display:none;"><table><tbody><tr><td><i class="material-icons" aria-hidden="true">delete_outline</i></td><td>delete selected</td></tr></tbody></table></div></div>'
    jqWfUiEl = $(html)
    var jqWfUIMainEl = jqWfUiEl
    jqVisNetworkEl.prepend(jqWfUiEl)
    jqWfUiEl = jqWfUiEl.find('.wfvis-ui-delete')
    var _vfThis = this
    $(jqVisNetworkEl).resize(function () {
      _vfThis.absoluteCenter(jqWfUIMainEl, $(this))
    })
    _vfThis.absoluteCenter(jqWfUIMainEl, jqVisNetworkEl)
    jqWfUiEl.click(function () {
      console.log('delete....')
      $(this).hide()
      _vfThis.Delete()
    })
  }
  jqWfUiEl.show()
  this.absoluteCenter(jqVisNetworkEl.find('.wfvis-ui-div'), jqVisNetworkEl)
}

VisFlowchart.prototype.Delete = function () {
  if (this.readOnly) {
    return
  }
  if (this.currentSelected) {
    if (this.checkStart(this.currentSelected) &&
            this._canDelete(this.currentSelected)) {
      this.hideDeleteOnUI()
      this._deleteElements(this.currentSelected)
    }
    this.currentSelected = null
  }
}

VisFlowchart.prototype.hideDeleteOnUI = function (ev) {
  var jqVisNetworkEl = this.container.find('.vis-network')
  var jqWfUiEl = jqVisNetworkEl.find('.wfvis-ui-div .wfvis-ui-delete')
  jqWfUiEl.hide()
}

export default VisFlowchart
