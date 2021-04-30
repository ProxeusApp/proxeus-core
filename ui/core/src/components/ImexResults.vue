<template>
<div>
  <div v-if="imexResult && mapSize(imexResult.results)===0" class="text-muted"
       style="text-align: center;">{{$t('Processed none')}}
  </div>
  <table v-else-if="res.length>0" style="width: 100%;">
    <tbody>
    <tr>
      <td>
        <table class="results">
          <tbody>
          <tr style="background: #00000008;border-bottom: 1px solid #c5c5c5;">
            <th class="tdmax impcnt" style="min-width: 150px;">
              <strong
                v-if="imexResult">{{imexResult.filename}} {{ imexResult.timestamp | moment('DD.MM.YY - HH:mm') }}
              </strong>
            </th>
            <th class="tdmin" v-for="(item, key) in resultsArr">{{item.name}}</th>
          </tr>
          <tr>
            <td class="tdmin" style="text-align: left;">
              <strong>{{$t('Processed')}}</strong>
            </td>
            <td v-for="(item, key) in resultsArr">{{item.count}}</td>
          </tr>
          <tr style="color:#bc0a00;background: rgba(255, 0, 0, 0.03);">
            <td class="tdmin" style="text-align: left;">
              <strong>{{$t('Errors')}}</strong>
            </td>
            <td v-for="(item, key) in resultsArr" :style="item.errCount===0?'color:green;':''">{{item.errCount}}</td>
          </tr>
          </tbody>
        </table>
      </td>
    </tr>
    <tr v-if="errors.length>0">
      <td>
        <div class="res-subtitle" style="color:#bc0a00">{{$t('Errors in detail')}}</div>
        <table class="results-err">
          <tbody>
          <tr v-for="(rowItem, index) in errors">
            <td>
              <strong>{{rowItem.name}}</strong>
            </td>
            <td style="color:#bc0a00;background: rgba(255, 0, 0, 0.03);">
              <table v-if="rowItem.errs && rowItem.errs.length>0">
                <tbody>
                <tr v-for="(item, index2) in rowItem.errs">
                  <td style="color:grey;">{{item.id}}</td>
                  <td>{{item.err}}</td>
                </tr>
                </tbody>
              </table>
            </td>
          </tr>
          </tbody>
        </table>
      </td>
    </tr>
    </tbody>
  </table>
</div>
</template>

<script>
import moment from 'moment'

export default {
  name: 'imex-results',
  components: {
    moment
  },
  props: {
    imexResult: {
      type: Object,
      default: null
    }
  },
  data () {
    return {
      resultsArr: [],
      errors: []
    }
  },
  computed: {
    res () {
      this.resultsArr = []
      this.errors = []
      const results = this.imexResult.results
      if (results) {
        for (const key in results) {
          if (results.hasOwnProperty(key)) {
            const errs = []
            this.resultsArr.push(
              { name: key, errCount: this.errCountOf(results[key]), count: this.mapSize(results[key]) })
            for (const id in results[key]) {
              if (results[key].hasOwnProperty(id) && results[key][id]) {
                errs.push({ id: id, err: results[key][id] })
              }
            }
            if (errs.length > 0) {
              this.errors.push({ name: key, errs: errs })
            }
          }
        }
      }
      return this.resultsArr
    }
  },
  methods: {
    mapSize (obj) {
      let size = 0
      if (obj) {
        for (const key in obj) {
          if (obj.hasOwnProperty(key)) {
            size++
          }
        }
      }
      return size
    },
    errCountOf (item) {
      let size = 0
      for (const id in item) {
        if (item.hasOwnProperty(id) && item[id]) {
          size++
        }
      }
      return size
    }
  }
}
</script>

<style lang="scss">
  .res-subtitle {
    padding: 10px;
    text-align: center;
    margin-top: 30px;
    font-weight: bold;
    font-size: 18px;
    background: #00000003;
  }

  .results, .results-err {
    width: 100%;
  }

  .results > tbody > tr > th:first-child, .results > tbody > tr > td:first-child {
    padding: 10px;
    text-align: left;
  }

  .results > tbody > tr > th, .results > tbody > tr > td {
    text-align: center;
  }

  .results-err th, .results-err td {
    padding: 10px;
  }

  table.results-err > tbody > tr:nth-child(odd) {
    background: #00000008;
  }
</style>
