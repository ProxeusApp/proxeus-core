<template>
<div class="file_drop_box m-0 w-100 position-relative" v-bind="$attrs">
  <input :multiple="multiple" type="file"
         @change="filesChange($event.target.files)"
         ref="fileInput"
         class="input-file">
  <div class="text-center w-100 drop-icon py-2">
    <i class="mdi mdi-arrow-collapse-down md-48"></i>
  </div>
  <div class="text-center pt-2 pm-1 pb-2">
    <span v-if="text" class="h5 text-muted p-0 pb-3"> {{ text }} </span>
  </div>
</div>
</template>

<script>

export default {
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
  @use "../../assets/styles/variables.scss";

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
