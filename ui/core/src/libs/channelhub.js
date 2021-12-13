import ReconnectingWebSocket from './ReconnectingWebsocket'

var ChannelHub = function (url, protocols) {
  this.channels = {}
  this.sysChanls = {}
  if (window && window.WebSocket) {
    var _ = this
    if (!(/^wss?:\/\/.*/.test(url))) {
      if (!(/^\/.*/.test(url))) {
        url = '/' + url
      }
      var loc = window.location
      var p = 'ws:'
      if (loc.protocol) {
        p = loc.protocol.replace('http', 'ws')
      }
      url = p + '//' + loc.host + url
    }
    _.ok = false
    this.ws = new ReconnectingWebSocket(url, protocols, {
      debug: false,
      reconnectInterval: 3000,
      maxReconnectInterval: 8000
    })
    this.ws.onopen = function (evt) {
      console.log('onopen')
      console.log(evt)
      _.ok = true
      for (var key in _.channels) {
        if (_.channels.hasOwnProperty(key)) {
          _.channels[key].onopen()
        }
      }
    }
    this.ws.onclose = function (evt) {
      console.log('onclose')
      console.log(evt)
      _.ok = false
    }
    this.ws.onerror = function (evt) {
      console.log('onerror')
      console.log(evt)
    }
    this.ws.onconnecting = function (evt, b) {
      console.log('onconnecting')
      console.log(evt)
      console.log(b)
    }
    this.ws.onmessage = function (evt) {
      var cMsg = JSON.parse(evt.data)
      if (cMsg.m === 'sys') {
        if (cMsg.d && cMsg.d.length) {
          for (var i = 0; i < cMsg.d.length; i++) {
            if (cMsg.d[i]) {
              _.sysChanls[cMsg.d[i]] = true
            }
          }
        }
      } else {
        if (cMsg.cid) {
          var chanl = _.channels[cMsg.cid]
          if (chanl) {
            if (chanl.onMessage) {
              chanl.onMessage(cMsg)
            }
          }
        }
      }
    }
  } else {
    throw 'no WebSocket support!'
  }
  this.Channel = function (map) {
    if (map) {
      var chanl
      for (var key in map) {
        if (map.hasOwnProperty(key)) {
          chanl = this._getchanl(key)
          chanl.public = map[key]
        }
      }
    }
  }
  this.channel = this.Channel
  this.Send = function (obj) {
    var to = function (cid) {
      if (cid && typeof cid === 'string') {
        this.cid = cid
        var chanl = this.ch._getchanl(cid)
        var msg = this
        this.promiseExe = function (ok, err) {
          msg._ok = ok
          msg._err = err
          chanl.send(msg)
        }
        this.promise = new Promise(this.promiseExe)

        this.then = function (ok, err) {
          return this.promise.then(ok, err)
        }
        delete this.ch
        return this
      } else {
        throw 'you must provide a channel id!'
      }
    }
    return { m: 'pub', d: obj, ch: this, To: to, to: to }
  }
  this.send = this.Send
  this._send = function (obj) {
    if (this.ok) {
      var newmsg = { m: obj.m }
      if (obj.d !== undefined) {
        newmsg.d = obj.d
      }
      if (obj.cid !== undefined) {
        newmsg.cid = obj.cid
      }
      if (obj.rid !== undefined) {
        newmsg.rid = obj.rid
      }
      if (obj.u !== undefined) {
        newmsg.u = obj.u
      }
      this.ws.send(JSON.stringify(newmsg))
      return true
    }
    return false
  }
  this._getchanl = function (key) {
    var chanl = this.channels[key]
    if (!chanl) {
      chanl = this._newChannel(key)
      this.channels[key] = chanl
      chanl.subscribe()
    }
    return chanl
  }
  this._newChannel = function (id) {
    return {
      ch: this,
      id: id,
      public: null,
      msgs: {},
      msgIndex: 0,
      onMessage: function (cMsg) {
        if (cMsg) {
          if (cMsg.rid) {
            var msg = this.msgs[cMsg.rid]
            if (msg) {
              this.okResp(msg, cMsg)
              delete this.msgs[cMsg.rid]
            }
          }
          if (cMsg.m === 'sub') {
            this.ok = cMsg.d === 'ok'
            var sendLater = this.sendLater
            this.sendLater = []
            for (var i = 0; i < sendLater.length; i++) {
              this.send(sendLater[i])
            }
          } else {
            if (this.public && this.public.onMessage) {
              this.public.onMessage(cMsg.d, cMsg)
            }
          }
        }
      },
      subscribe: function () {
        this.send({ m: 'sub', cid: this.id })
      },
      sendLater: [],
      subLater: [],
      onopen: function () {
        var subs = this.subLater
        this.subLater = []
        for (var i = 0; i < subs.length; i++) {
          this.send(subs[i])
        }
      },
      send: function (msg) {
        if (msg && msg.cid) {
          if (msg.m === 'sub') {
            if (!this.ch._send(msg)) {
              this.subLater.push(msg)
            }
          } else {
            this.msgIndex++
            msg.rid = this.msgIndex
            if (this.ok === undefined) {
              if (this.ch.sysChanls[msg.cid]) {
                this.errResp(msg)
              } else {
                this.sendLater.push(msg)
              }
            } else if (this.ok) {
              if (this.ch._send(msg)) {
                if (msg.m == 'pub') {
                  this.okResp(msg, { d: 'ok' })
                } else {
                  this.msgs[msg.rid] = msg
                }
              } else {
                this.errResp(msg)
              }
            } else {
              this.errResp(msg)
            }
          }
        }
      },
      okResp: function (msg, cMsg) {
        if (msg._ok) {
          msg._ok(cMsg.d, cMsg)
        }
      },
      errResp: function (msg) {
        if (msg._err) {
          msg._err('not possible!')
        }
      }
    }
  }
}

export default ChannelHub
