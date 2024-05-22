<template>
<div style="position: relative;">
  <file-drop-box v-if="!file" :text="$t('Choose a profile photo or drag it here')" @dropped="drop"
                 :multiple="false"></file-drop-box>
  <div v-else>
    <div class="crper_cont">
      <vue-cropper
        ref='cropper'
        :guides="true"
        :aspectRatio="1/1"
        :view-mode="2"
        drag-mode="crop"
        :auto-crop-area="0.5"
        :min-container-width="250"
        :min-container-height="180"
        :minCropBoxWidth="80"
        :background="true"
        :rotatable="true"
        :src="imgSrc"
        alt="Source Image"
        :img-style="{ 'width': '400px', 'height': '300px' }">
      </vue-cropper>
    </div>
    <table style="width:100%">
      <tr>
        <td style="text-align: left;"><span class="crper_btn"
                                            @click="file=null">{{$t('choose another profile photo','choose another')}}</span>
        </td>
        <td class="tdmin"><span class="crper_btn" @click="rotateLeft"><i class="material-icons">undo</i>{{$t('rotate left','left')}}</span>
        </td>
        <td class="tdmin"><span class="crper_btn" @click="rotateRight"><i class="material-icons">redo</i>{{$t('rotate right','right')}}</span>
        </td>
      </tr>
    </table>
  </div>
</div>
</template>

<script>
import VueCropper from 'vue-cropperjs'
import FileDropBox from './template/FileDropBox.vue'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'change-profile-photo',
  components: {
    FileDropBox,
    VueCropper
  },
  props: {
    setup: {
      type: Function,
      default: () => null
    }
  },
  data () {
    return {
      imgSrc: '',
      file: null
    }
  },
  created () {
    if (this.setup) {
      this.setup(this.cropCanvasToBlob)
    }
  },
  methods: {
    cropCanvasToBlob (cb) {
      this.$refs.cropper.getCroppedCanvas().toBlob((blob) => {
        if (cb) {
          cb(this.file, blob)
        }
      })
    },
    drop (file) {
      if (!file.type.includes('image/')) {
        alert('Please select an image file')
        return
      }
      if (typeof FileReader === 'function') {
        const reader = new FileReader()
        reader.onload = (event) => {
          this.imgSrc = event.target.result
          // rebuild cropperjs with the updated source
          this.$refs.cropper.replace(event.target.result)
        }
        this.file = file
        this.$emit('onFileDrop', this.file)
        reader.readAsDataURL(file)
      } else {
        alert('Sorry, FileReader API not supported')
      }
    },
    rotateLeft () {
      // guess what this does :)
      this.$refs.cropper.rotate(-90)
    },
    rotateRight () {
      // guess what this does :)
      this.$refs.cropper.rotate(90)
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .crper_cont {
    width: 100%;
    height: 300px;
    border: 1px solid #ececec;
    display: inline-block;
  }

  .crper_btn {
    vertical-align: middle;
    padding: 4px;
    display: inline-block;
  }

  .crper_btn i {
    vertical-align: middle;
    padding-right: 4px;
  }

  .crper_btn:hover {
    background: #00000036;
    cursor: pointer;
  }
</style>
