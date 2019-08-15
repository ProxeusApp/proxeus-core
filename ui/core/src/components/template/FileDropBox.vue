<template>
<div class="file_drop_box m-0 w-100 position-relative" v-bind="$attrs">
  <input :multiple="multiple" type="file"
         @change="filesChange($event.target.files)"
         ref="fileInput"
         class="input-file">
  <drop-file-design :text="text"/>
</div>
</template>

<script>
import DropFileDesign from './DropFileDesign'

export default {
  components: { DropFileDesign },
  name: 'file-drop-box',
  props: {
    multiple: {
      type: Boolean,
      default: true
    },
    text: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      uploadFieldName: 'docs',
      unsavedFile: undefined,
      savedFile: undefined
    }
  },
  methods: {
    filesChange (fileList) {
      if (fileList && fileList.length) {
        for (var i = 0; i < fileList.length; i++) {
          if (fileList[i] && fileList[i].name) {
            this.$emit('dropped', fileList[i])
          }
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~bootstrap/scss/functions";
  @import "../../assets/styles/variables.scss";

  .icon-unpersisted {
    position: absolute;
    margin: 0.5rem;
  }

  .input-file {
    opacity: 0; /* invisible but it's there! */
    width: 100%;
    height: 100%;
    position: absolute;
    cursor: pointer;
  }

  .file_drop_box {
    border: 2px dashed #dddddd;
    min-width: 155px;

    &.active {
      border: 2px dashed #40e1d1;
    }
    &.unsaved {
      background: rgba(67, 255, 239, 0.2);
      border: 2px dashed #40e1d1;
    }

    .heading-small {
      font-size: 1.2rem;
    }
  }

  .file-upload-btns .fub-btn.disabled {
    cursor: not-allowed;
    background-color: #e6e6e6 !important;
    color: #bbbbbb !important;
    border-color: #bbbbbb !important;
  }

  .file_drop_box:hover {
    background: rgba(67, 255, 239, 0.4);
  }

  .file_drop_box p {
    padding: 10px 0;
    text-align: center;
  }
</style>
