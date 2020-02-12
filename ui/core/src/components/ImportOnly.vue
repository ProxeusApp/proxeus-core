<template>
<div>
  <div class="tabcontent">
    <h5>
      <strong>{{$t('Import a Proxeus database file')}}</strong>
    </h5>
    <file-drop-box @dropped="dropImport" text="Drop or click to select a Proxeus exported database file to import data."></file-drop-box>
    <div class="mt-2">
      <table class="imptbl" style="width: 100%;">
        <tbody>
        <tr>
          <td style="text-align: left">
            <checkbox :disabled="importFile ? false : true"
                      :label="$t('Skip existing entries')"
                      v-model="skipExistingEntries"/>
            <span
              class="text-muted">{{$t('Import skip explanation','Choose whether existing entries should be skipped or overwritten.')}}</span>
          </td>
          <td>
            <div v-if="importFile">{{importFile.name}}</div>
          </td>
          <td style="text-align: right">
            <button @click="importClick" :disabled="importFile?false:true" type="button" class="btn btn-primary">
              <span>{{$t('Import')}}</span>
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
  <div v-if="lastImportResults" class="tabcontent" style="margin-top:15px;position:relative;">
    <span class="imexclose" @click="app.loadLastImportResults(updateImportResults,'delete')">&#10799;</span>
    <h5>
      <strong>{{$t('Last Import')}}</strong>
    </h5>
    <imex-results :imexResult="lastImportResults"/>
  </div>
</div>
</template>

<script>
import Checkbox from './Checkbox'
import ImexResults from './ImexResults'
import FileDropBox from './template/FileDropBox'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'import-only',
  components: { FileDropBox, ImexResults, Checkbox },
  data () {
    return {
      skipExistingEntries: false,
      importFile: null,
      lastImportResults: null
    }
  },
  created () {
  },
  beforeDestroy () {
  },
  mounted () {
    this.app.loadLastImportResults(this.updateImportResults)
  },
  methods: {
    updateImportResults (results) {
      this.lastImportResults = results
    },
    mapSize (obj) {
      let size = 0
      if (obj) {
        for (let key in obj) {
          if (obj.hasOwnProperty(key)) {
            size++
          }
        }
      }
      return size
    },
    dropImport (file) {
      this.importFile = file
    },
    importClick () {
      this.app.importData(this.importFile, this.skipExistingEntries, this.updateImportResults)
    }
  }
}
</script>

<style>
  .imexclose {
    position: absolute;
    top: 2px;
    right: 3px;
    font-size: 28px;
    color: grey;
    cursor: pointer;
    line-height: 1;
  }
</style>
