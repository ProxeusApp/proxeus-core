export default {
  data () {
    return {
      lastSnapshot: ''
    }
  },
  beforeRouteLeave (to, from, next) {
    if (this.hasUnsavedChangesMethodImplemented() && this.hasUnsavedChanges()) {
      const answer = window.confirm(
        'Do you really want to leave? You have unsaved changes!')
      if (answer) {
        next()
      } else {
        next(false)
      }
    } else {
      next()
    }
  },
  methods: {
    snapshot (obj, skipper) {
      this.lastSnapshot = this.serializeToStr(obj, skipper)
    },
    compare (newObj, skipper) {
      var comparedWith = this.serializeToStr(newObj, skipper)
      return this.lastSnapshot === comparedWith
    },
    serializeToStr (obj, skipper) {
      var sorted = []
      var func = function (value, keyOrIndex, obj, fullpath) {
        if (this.skipper && this.skipper(value, keyOrIndex, obj, fullpath)) {
          return true
        }
        if (!isNaN(value)) { // prevent from taking 0 or negative numbers into account
          value = '' + value
        }
        if (value) {
          sorted.push(fullpath + ':' + (value + '').trim())
        }
        return true
      }
      this.jsonCrawler(obj, {
        skipper: skipper,
        string: func,
        number: func,
        boolean: func
      })
      sorted.sort()
      var newSnapshot = ''
      for (var i = 0; i < sorted.length; i++) {
        newSnapshot += sorted[i]
      }
      return newSnapshot
    },
    jsonCrawler (obj, options, path) {
      var makeChildPath = function (p, c) {
        if (p === undefined || p === '') {
          return c + ''
        }
        return p + '.' + c
      }
      if (typeof obj === 'object') {
        for (var i in obj) {
          if (obj.hasOwnProperty(i)) {
            if (typeof obj[i] === 'object') {
              try {
                if (options[(typeof obj[i]) + ''](obj[i], i, obj,
                  makeChildPath(path, i))) {
                  if (this.jsonCrawler(obj[i], options,
                    makeChildPath(path, i))) {
                    continue
                  } else {
                    return false
                  }
                }
              } catch (functionForTypeMissing) {
                if (this.jsonCrawler(obj[i], options, makeChildPath(path, i))) {
                  continue
                } else {
                  return false
                }
              }
            } else {
              // check for number, boolean or string
              if (obj[i] !== null || obj[i] !== undefined) {
                try {
                  if (options[(typeof obj[i]) + ''](obj[i], i, obj,
                    makeChildPath(path, i))) {
                    continue
                  } else {
                    return false
                  }
                } catch (functionForTypeMissing) {}
              }
            }
          }
        }
      }
      return true
    },
    beforeUnload (e) {
      if (this.hasUnsavedChangesMethodImplemented() &&
        this.hasUnsavedChanges()) {
        e.preventDefault()
        // custom messages are not supported anymore in beforeunload
        // https://stackoverflow.com/questions/38879742/is-it-possible-to-display-a-custom-message-in-the-beforeunload-popup
        // this return value is going to be used as a flag on the new browsers
        return e.returnValue = 'Do you really want to leave? You have unsaved changes!'
      }
    },
    hasUnsavedChangesMethodImplemented () {
      if (!this.hasUnsavedChanges) {
        return false
      }
      return true
    }
  },
  created () {
    window.addEventListener('beforeunload', this.beforeUnload)
  },
  beforeDestroy () {
    window.removeEventListener('beforeunload', this.beforeUnload)
  }
}
