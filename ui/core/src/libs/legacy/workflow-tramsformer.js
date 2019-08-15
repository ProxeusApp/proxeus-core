
function VisFlowchart (wfData, pOptions) {
  this.randomId = function () {
    return Math.random().toString(36).substring(2, 18) + Math.random().toString(36).substring(2, 18)
  }
  // keep
  this.s = {
    startId: 'start111',
    conNode: 'conNode',
    conEdge: 'conEdge'
  }
  // keep
  this.options = {
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
                color: { color: '#45bbff', highlight: '#45bbff', hover: '#45bbff' }
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
    node: {
      connections: {
        from: [{
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
  var nNodeOptions = {}
  var nn
  for (var nk in this.options.nodes) {
    if (this.options.nodes.hasOwnProperty(nk)) {
      if (nk === 'con') {
        nn = this.options.nodes[nk]
      } else {
        nn = {}
        $.extend(true, nn, this.options.node, this.options.nodes[nk], { connections: null, events: null })
        delete nn['connections']
        delete nn['events']
      }
      nNodeOptions[nk] = nn
    }
  }
}

VisFlowchart.prototype.workflowDataToVis = function (wfData) {
  var dataResult = { nodes: [], edges: [] }
  var startCon = this.getConnectionSettings('start')
  var newNode = { 'id': this.s.startId, label: 'start', physics: false, group: 'start', connections: startCon }
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
      dataResult.edges.push(this.createEdge('start', { from: newNode.id, to: wfData[0].id }, 0))
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
            newNode._cases.push({ name: w.cases[i2].label, value: w.cases[i2].value })
            this.mergeConnectionFromSettings(newNode, w.cases[i2], i2)
            this.createConNode(dataResult, newNode, i2)
          }
        } else {
          this.mergeConnectionFromSettings(newNode, { value: newNode.id }, 0)
          newNode['uriName'] = w.uriName
          this.createConNode(dataResult, newNode, 0)
        }
        this.mergeNodeEvents(newNode)
        if (w.connectedTo && w.connectedTo.length) {
          for (var c = 0; c < w.connectedTo.length; ++c) {
            cnd = w.connectedTo[c]
            incrementToCount(cnd.targetId, toNodesMap)
            if (w.type === 'condition') {
              dataResult.edges.push(this.createEdge(w.type, { from: w.id, to: cnd.targetId }, c, cnd))
              for (var b = 0; b < newNode.connections.from.length; ++b) {
                if (newNode.connections.from[b].value == cnd.value) {
                  newNode.connections.from[b].hide(this)
                }
              }
            } else {
              dataResult.edges.push(this.createEdge(w.type, { from: w.id, to: cnd.targetId }, 0))
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
VisFlowchart.prototype.nodes = null
VisFlowchart.prototype.edges = null
VisFlowchart.prototype.network = null
VisFlowchart.prototype.seed = 2

// keep
VisFlowchart.prototype.createConFromNode = function (node, i, from) {
  var n = this.getMergedNode(node, i, from)
  n.physics = false
  n.id = this.s.conNode + '_' + node.id + (from.value ? '_' + from.value : '') + '_' + i
  n.group = 'con'
  n._isCon = true
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
  if (from.notAvailable) {
    n.hidden = true
    e.hidden = true
    e.physics = false
  }
  from['ne'] = {
    n: n,
    e: e,
    hide: function (fc, f) {
      f.notAvailable = true
      if (fc.network) {
        fc.network.body.data.nodes.update([{ id: this.n.id, hidden: true }])
        fc.network.body.data.edges.update([{ id: this.e.id, hidden: true, physics: false }])
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
        fc.network.body.data.edges.update([{ id: this.e.id, hidden: false, physics: true }])
      } else {
        this.n.hidden = false
        this.e.hidden = false
        this.e.physics = true
      }
    }
  }
  return from['ne']
}

// keep
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
  $.extend(true, newNode, this.options.node.connections.from[0].node, defaultSettings, newSettings)
  return newNode
}

VisFlowchart.prototype.getNodeById = function (id) {
  return this.network.body.data.nodes.get(id)
}
VisFlowchart.prototype.getEdgeById = function (id) {
  return this.network.body.data.edges.get(id)
}

// keep
VisFlowchart.prototype._createNode = function (w, id) {
  var con = this.getConnectionSettings(w.type)
  con.from = []
  con.toCount = 0
  if (w.description) {
    w.detail = w.description
  }
  var n = { 'id': w.id, 'label': w.name, 'detail': w.detail, group: w.type, connections: con }
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

// Keep
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
  return '<div class="wf-node-tooltip">' + (showID ? ('<p>' + (id) + '</p>') : '') + '<strong class="wfnt-name">' + (name) + '</strong><br/><span class="wfnt-desc">' + (detail) + '</span>' + this.createTooltipConditionCases(cases) + '</div>'
}

// keep
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
  $.extend(true, newNode.connections.from[i], this.options.node.connections.from[0], nodeSpecificFromStyles, settings, {
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
// keep
VisFlowchart.prototype.getConnectionSettings = function (group) {
  var nodeSpecificFromStyles
  try {
    nodeSpecificFromStyles = this.options.nodes[group].connections
  } catch (dontCare1) {
  }
  return $.extend(true, {}, this.options.node.connections, nodeSpecificFromStyles, { from: null })
}
// keep
VisFlowchart.prototype.mergeNodeEvents = function (newNode) {
  try {
    newNode.events = this.options.nodes[newNode.group].events
  } catch (dontCare1) {
  }
}

// keep
VisFlowchart.prototype.createConNode = function (dataResult, newNode, i) {
  var ne = this.createConFromNode(newNode, i, newNode.connections.from[i])
  dataResult.nodes.push(ne.n)
  dataResult.edges.push(ne.e)
  ne.n.x = newNode.x
  ne.n.y = newNode.y
}

// keep
VisFlowchart.prototype.createEdge = function (grp, con, i, more) {
  var newEdge = this.createMergedEdge(grp, i, more)
  newEdge.from = con.from
  newEdge.to = con.to
  newEdge.id = con.from + '-' + con.to + '_' + i + '_' + this.randomId()
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

// keep
VisFlowchart.prototype.createMergedEdge = function (grp, i, more) {
  var defaultSettings
  if (typeof i === 'number') {
    try {
      defaultSettings = this.options.nodes[grp].connections.from[i].edge
    } catch (dontCare) {
    }
  }
  if (!defaultSettings) {
    try {
      defaultSettings = this.options.nodes[grp].connections.from[0].edge
    } catch (dontCare) {
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

export default VisFlowchart
