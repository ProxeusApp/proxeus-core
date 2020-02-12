function transform (wfData) {
  if (!String.prototype.startsWith) {
    (function () {
      'use strict'
      var defineProperty = (function () {
        try {
          var object = {}
          var $defineProperty = Object.defineProperty
          var result = $defineProperty(object, object, object) &&
            $defineProperty
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
        //
        var pos = position ? Number(position) : 0
        if (pos != pos) { // better
          pos = 0
        }
        var start = Math.min(Math.max(pos, 0), stringLength)
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
          'value': startsWith,
          'configurable': true,
          'writable': true
        })
      } else {
        String.prototype.startsWith = startsWith
      }
    }())
  }

  var NetworkStub = function () {
    this.body = {}
    this.body.nodes = {}
    this.body.edges = {}
    this.body.data = {}
    this.body.data.nodes = {}
    this.body.data.nodes._data = {}
    this.body.data.edges = {}
    this.body.data.edges._data = {}
  }

  NetworkStub.prototype.fit = function () {}
  NetworkStub.prototype.storePositions = function () {}

  // import VisFlowchart from './vis-flowchart'
  // Pass in the objects to merge as arguments.
  // For a deep extend, set the first argument to `true`.
  var extend = function () {
    // Variables
    var extended = {}
    var deep = false
    var i = 0
    var length = arguments.length

    // Check if a deep merge
    if (Object.prototype.toString.call(arguments[0]) === '[object Boolean]') {
      deep = arguments[0]
      i++
    }

    // Merge the object into the extended object
    var merge = function (obj) {
      for (var prop in obj) {
        if (Object.prototype.hasOwnProperty.call(obj, prop)) {
          // If deep merge and property is an object, merge properties
          if (deep &&
            Object.prototype.toString.call(obj[prop]) === '[object Object]') {
            extended[prop] = extend(true, extended[prop], obj[prop])
          } else {
            extended[prop] = obj[prop]
          }
        }
      }
    }

    // Loop through each object and conduct a merge
    for (; i < length; i++) {
      var obj = arguments[i]
      merge(obj)
    }

    return extended
  }

  function VisFlowchart (wfData, pOptions) {
    this.randomId = function () {
      return Math.random().toString(36).substring(2, 18) +
        Math.random().toString(36).substring(2, 18)
    }
    this.s = {
      startId: 'start111',
      conNode: 'conNode',
      conEdge: 'conEdge'
    }
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
        'hoverIn': null/* function(){} */,
        'hoverOut': null/* function(){} */,
        'insert': null/* function(){} */,
        'remove': null/* function(){} */,
        'click': null/* function(){} */,
        'dblclick': null/* function(){} */
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
                color: {
                  color: '#45bbff',
                  highlight: '#45bbff',
                  hover: '#45bbff'
                }
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

    this.options = extend(true, this.options, pOptions)

    var nNodeOptions = {}
    var nn
    for (var nk in this.options.nodes) {
      if (this.options.nodes.hasOwnProperty(nk)) {
        if (nk === 'con') {
          nn = this.options.nodes[nk]
        } else {
          nn = {}
          nn = extend(true, nn, this.options.node, this.options.nodes[nk],
            { connections: null, events: null })
          delete nn['connections']
          delete nn['events']
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
        keyboard: true
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
          'forceDirection': 'none'
        }
      },
      nodes: {
        font: {
          size: 18
        }
      },
      // layout:{randomSeed:0},
      // layout: {randomSeed:seed}, // just to make sure the layout is the same when the locale is changed
      // // locale: document.getElementById('locale').value,
      physics: {
        // enabled:false,
        stabilization: {
          enabled: true,
          // iterations: 10,
          // updateInterval: 50,
          onlyDynamicEdges: false,
          fit: true
        },
        barnesHut: {
          'gravitationalConstant': -850,
          'centralGravity': 0.00001,
          'springLength': 200,
          'springConstant': 0.0885,
          'damping': 0.6,
          'avoidOverlap': 0
        },
        repulsion: {
          centralGravity: 0.002,
          springLength: 200,
          springConstant: 0.08,
          nodeDistance: 100,
          damping: 0.8
        },
        timestep: 1,
        adaptiveTimestep: false,
        maxVelocity: 50,
        minVelocity: 0.4,
        solver: 'barnesHut'
      }
    }
  }

  VisFlowchart.prototype.workflowDataToVis = function (wfData) {
    var dataResult = { nodes: [], edges: [] }
    var startCon = this.getConnectionSettings('start')
    var newNode = {
      'id': this.s.startId,
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
            newNode['uriName'] = w.uriName
            this.createConNode(dataResult, newNode, 0)
          }
          this.mergeNodeEvents(newNode)
          if (w.connectedTo && w.connectedTo.length) {
            for (var c = 0; c < w.connectedTo.length; ++c) {
              cnd = w.connectedTo[c]
              incrementToCount(cnd.targetId, toNodesMap)
              if (w.type === 'condition') {
                dataResult.edges.push(
                  this.createEdge(w.type, { from: w.id, to: cnd.targetId }, c,
                    cnd))
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
  VisFlowchart.prototype.network = new NetworkStub()
  VisFlowchart.prototype.seed = 2

  VisFlowchart.prototype.createConFromNode = function (node, i, from) {
    var n = this.getMergedNode(node, i, from)
    n.physics = false
    n.id = this.s.conNode + '_' + node.id +
      (from.value ? '_' + from.value : '') +
      '_' + i
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
      e: e
    }
    return from['ne']
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
    newNode = extend(true, newNode,
      this.options.node.connections.from[0].node,
      defaultSettings, newSettings)
    return newNode
  }

  VisFlowchart.prototype.getFlowData = function (wfData) {
    var output = { start: {}, nodes: {} }
    var targetNode = null
    var layoutNode = null
    var
      ioNode = null
    var startP = null
    var startNode = this.network.body.nodes[this.s.startId]
    var id, i, c, e, conn
    var startNodeToId = this.getStartNodeId()
    if (startNode && startNodeToId) {
      layoutNode = startNode
      startP = { 'x': layoutNode.x, 'y': layoutNode.y }
      output.start.p = startP
      output.start.node = startNodeToId
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
          p: { 'x': layoutNode.x, 'y': layoutNode.y },
          conns: []
        }
      if (targetNode._cases) {
        ioNode.cases = targetNode._cases
      }
      if (targetNode._data) {
        ioNode.data = targetNode._data
      }

      var edg = this.getEdgesByNode(layoutNode)

      for (c = 0; c < edg.length; ++c) {
        e = edg[c]
        if (!e.id.startsWith(this.s.conEdge) && e.to !== targetNode.id) {
          e = this.getEdgeById(e.id)
          conn = { 'id': e.to/*, "data":{"label": e.label, "value": e.cvalue} */ }
          if (ioNode.id !== e.cvalue) {
            conn['value'] = e.cvalue
          }
          ioNode.conns.push(conn)
        }
      }
    }
    return output
  }

  VisFlowchart.prototype.getEdgesByNode = function (node) {
    var edges = []
    var keys = Object.keys(this.network.body.edges)

    keys.forEach(key => {
      if (this.network.body.edges[key].from === node.id ||
        this.network.body.edges[key].to === node.id) {
        edges.push(this.network.body.edges[key])
      }
    })
    return edges
  }

  VisFlowchart.prototype.getStartNodeId = function () {
    var keys = Object.keys(this.network.body.edges)
    var nodeId = ''
    keys.forEach(key => {
      if (this.network.body.edges[key] &&
        this.network.body.edges[key].from === this.s.startId) {
        nodeId = this.network.body.edges[key].to
      }
    })
    return nodeId
  }

  VisFlowchart.prototype.getNodeById = function (id) {
    return this.network.body.data.nodes._data[id]
  }
  VisFlowchart.prototype.getEdgeById = function (id) {
    return this.network.body.data.edges._data[id]
  }

  VisFlowchart.prototype._createNode = function (w, id) {
    var con = this.getConnectionSettings(w.type)
    con.from = []
    con.toCount = 0
    if (w.description) {
      w.detail = w.description
    }
    var n = {
      'id': w.id,
      'label': w.name,
      'detail': w.detail,
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

  VisFlowchart.prototype.createTooltip = function () {}

  VisFlowchart.prototype.mergeConnectionFromSettings = function (
    newNode, settings, i) {
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
    newNode.connections.from[i] = extend(true, newNode.connections.from[i],
      this.options.node.connections.from[0], nodeSpecificFromStyles, settings,
      {
        hide: function (fc) {
          if (this.ne && fc) {
          }
        },
        show: function (fc) {
          if (this.ne && fc) {
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
    return extend(true, {}, this.options.node.connections,
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
    newEdge = extend(true, newEdge,
      this.options.node.connections.from[0].edge,
      defaultSettings, newSettings)
    return newEdge
  }

  var vf = new VisFlowchart(wfData, {
    showIDOnTooltip: true,
    nodes: {
      start: {
        connections: {
          from: [
            {
              node: {
                color: {
                  background: '#45bbff',
                  highlight: { background: '#8dd5ff' },
                  hover: { background: '#8dd5ff' }
                }
              },
              edge: {
                color: {
                  color: '#45bbff',
                  highlight: '#45bbff',
                  hover: '#45bbff'
                }
              }
            }
          ]
        },
        icon: {
          face: 'Material Icons',
          code: 'radio_button_unchecked',
          size: 80,
          color: '#45bbff'
        }
      },
      collection: {
        connections: {
          from: [
            {
              node: {
                color: {
                  background: '#23d6d6',
                  highlight: { background: '#23d6d6' },
                  hover: { background: '#23d6d6' }
                }
              },
              edge: {
                color: {
                  color: '#23d6d6',
                  highlight: '#23d6d6',
                  hover: '#23d6d6'
                }
              }
            }
          ],
          fromInfinity: true,
          to: Infinity
        },
        font: {
          color: '#343434',
          size: 15,
          mod: 'bold'
        },
        icon: {
          face: 'Material Design Icons',
          code: '\uF765',
          color: '#1fa6a6'
        }
      },
      form: {
        connections: {
          from: [
            {
              node: {
                color: {
                  background: '#60aa61',
                  highlight: { background: '#99d29a' },
                  hover: { background: '#99d29a' }
                }
              },
              edge: {
                color: {
                  color: '#60aa61',
                  highlight: '#3c763d',
                  hover: '#3c763d'
                }
              }
            }],
          to: Infinity,
          space: 1.1
        },
        font: {
          color: '#343434',
          size: 15,
          mod: 'bold',
          bold: {
            color: '#343434',
            size: 14, // px
            face: 'arial',
            vadjust: 0,
            mod: 'bold'
          }
        },
        icon: {
          face: 'Material Icons',
          code: 'view_quilt',
          color: '#0eaa64'
        },
        events: {
        }
      },
      condition: {
        connections: {
          from: [
            {
              node: {
                color: {
                  background: '#f4a646',
                  highlight: { background: '#f9cd95' },
                  hover: { background: '#f9cd95' }
                }
              },
              edge: {
                font: { align: 'middle' },
                dashes: true,
                arrows: 'to',
                color: {
                  color: '#f4a646',
                  highlight: '#e88000',
                  hover: '#e88000'
                }
              }
            }],
          to: Infinity,
          space: 0.6
        },
        font: {
          color: '#343434',
          size: 15
        },
        icon: {
          face: 'Material Design Icons',
          code: '\uf70B',
          color: '#f0a30a'
        },
        events: {
        }
      },
      workflow: {
        connections: {
          from: [
            {
              node: {
                color: {
                  background: '#e40070',
                  highlight: { background: '#f3a3ca' },
                  hover: { background: '#f3a3ca' }
                }
              },
              edge: {
                color: {
                  color: '#e40070',
                  highlight: '#d60069',
                  hover: '#d60069'
                }
              }
            }],
          to: Infinity
        },
        icon: {
          face: 'Material Design Icons',
          code: '\uf62C',
          color: '#e40070'
        },
        events: {
        }
      },
      user: {
        connections: {
          from: [
            {
              node: {
                color: {
                  background: '#7200ff',
                  highlight: { background: '#c393ff' },
                  hover: { background: '#c393ff' }
                }
              },
              edge: {
                color: {
                  color: '#7200ff',
                  highlight: '#5b00cc',
                  hover: '#5b00cc'
                }
              }
            }],
          to: Infinity
        },
        icon: {
          face: 'Material Design Icons',
          code: '\uf004',
          color: '#7200ff'
        },
        events: {
        }
      },
      template: {
        connections: {
          from: [
            {
              node: {
                shape: 'dot',
                size: 8,
                color: {
                  background: '#ff30ec',
                  border: '#7b7b7b',
                  highlight: { background: '#ffabf7', border: '#7b7b7b' },
                  hover: { background: '#ffabf7', border: '#7b7b7b' }
                },
                borderWidth: 1,
                borderWidthSelected: 1
              },
              edge: {
                color: {
                  color: '#ff30ec',
                  highlight: '#e22cd1',
                  hover: '#e22cd1'
                }
              }
            }],
          to: Infinity
        },
        font: {
          color: '#343434',
          size: 15
        },
        icon: {
          face: 'Material Design Icons',
          code: '\uf22E',
          color: '#ff30ec'
        },
        events: {
        }
      }
    },
    events: {
    },
    node: {
      connections: {
        from: [
          {
            node: {
              shape: 'dot',
              size: 10,
              color: {
                border: '#7b7b7b',
                highlight: { border: '#7b7b7b' },
                hover: { border: '#7b7b7b' }
              },
              borderWidth: 1,
              borderWidthSelected: 1
            },
            edge: {
              arrowStrikethrough: false,
              arrows: 'to',
              width: 1,
              hoverWidth: 2,
              selectionWidth: 2
            }
          }],
        space: 1.1,
        to: Infinity
      },
      shape: 'icon',
      icon: {
        size: 50
      },
      shadow: {
        enabled: true,
        color: 'rgba(0,0,0,0.1)',
        size: 10,
        x: 5,
        y: 5
      }
    }
  })

  wfData = vf.workflowDataToVis(wfData)

  wfData.nodes.forEach(node => {
    vf.network.body.nodes[node.id] = node
  })

  wfData.edges.forEach(edge => {
    vf.network.body.edges[edge.id] = edge
  })

  vf.network.body.data.nodes.length = wfData.nodes.length
  wfData.nodes.forEach(node => {
    vf.network.body.data.nodes._data[node.id] = node
  })

  vf.network.body.data.edges.length = wfData.edges.length
  wfData.edges.forEach(edge => {
    vf.network.body.data.edges._data[edge.id] = edge
  })

  vf.network.body.nodeIndices = []
  wfData.nodes.forEach(node => {
    vf.network.body.nodeIndices.push(node.id)
  })

  return vf.getFlowData(wfData)
}
