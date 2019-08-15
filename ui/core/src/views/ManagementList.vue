<template>
<div class="management-list">
  <vue-headful title="Proxeus - Management List"/>
  <top-nav title="Management List" bg="#f7f8f9" class="border-bottom-0"/>
  <div class="col-12">
    <div class="table-responsive">
      <table class="table">
        <thead>
        <tr>
          <th>ID</th>
          <th>Owner</th>
          <th>Consignment ID</th>
          <th>Timestamp</th>
          <th>Signatory</th>
          <th>Action</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="item in myList">
          <td>{{item.id}}</td>
          <td>{{item.owner}}</td>
          <td>{{item.consignmentID}}</td>
          <td>{{item.timestamp}}</td>
          <td>{{item.signatory}}</td>
          <td>
            <button class="btn btn-primary">My Action</button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</div>
</template>

<script>
import SearchBox from '@/components/SearchBox'
import ListGroup from '@/components/ListGroup'
import TopNav from '@/components/layout/TopNav'
import Table from 'bootstrap-vue/es/components/table/table'

import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'management-list',
  components: {
    Table,
    SearchBox,
    ListGroup,
    TopNav
  },
  data () {
    return {
      theList: null
    }
  },
  computed: {
    myList: {
      get () {
        if (!this.omfg) {
          console.log('setinterval')
          var _this = this
          this.omfg = setInterval(function () {
            console.log('interval..')
            axios.get('/api/management-list', null).then(response => {
              _this.myList = response.data
            }, (err) => {
              this.app.handleError(err)
            })
          }, 1000)
        }
        console.log('my items...')

        return this.theList
      },
      set (l) {
        console.log('set...')
        console.log(l)
        this.theList = l
      }
    },
    items () {
      if (!this.omfg) {
        var _this = this
        this.omfg = setInterval(function () {
          console.log('interval..')
          var a = Math.floor(Math.random() * 9) + 1
          var b = []
          for (var i = 0; i < a; i++) {
            b.push('omfg_' + i)
          }
          _this.items = b
        }, 1000)
      }

      console.log('my items...')
      this.items = ['a', 'bb', 'c']
      return this.items
    }
  }
}
</script>

<style lang="scss" scoped>
  .management-list {
    overflow: hidden;
    height: auto;
  }
</style>
