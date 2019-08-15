<template>
<div>
  <div class="verification animscroll mt-2 col-lg-10 col-lg-offset-1">
    <div class="animscroll">
      <spinner v-show="!wallet || initialized === false" background="transparent" color="#eee"></spinner>
      <div class="light-text mx-auto mt-5"
           v-show="!wallet || initialized === false">{{ $t('Connecting with the blockchain ...') }}
      </div>
      <div class="dropbox bg-white pl-0" v-if="initialized"
           :class="{hasDocuments}">
        <file-drop-box @dropped="drop"></file-drop-box>
        <p class="w-100 text-center" v-if="!hasDocuments">
            <span class="text-hint py-2 px-4 d-inline-block">
             {{ $t('Verify file here', 'You can verify the authenticity of your document here. The hash of your document will be compared to the hashes, which were registered on the blockchain upon document creation.') }}
            </span>
        </p>
      </div>
      <div class="stats col-sm-12 p-0" v-show="hasDocuments">
        <verification-file-entry :singleFile="files.length === 1" v-for="file in files" :key="file.name" :file="file"
                                 :wallet="wallet"
                                 :thumbnail="thumbnail"></verification-file-entry>
      </div>
    </div>
  </div>
  <div class="powered-by row">
    <div class="offset-4 col-4 mt-4 text-center">
      <p class="light-text smaller m-0 p-x-0 mb-1">{{ $t('Powered by') }}</p>
      <a href="https://proxeus.com" target="_blank" class="d-block">
        <proxeus-logo :width="90"></proxeus-logo>
      </a>
    </div>
  </div>
</div>
</template>

<script>
import FileDropBox from '../template/FileDropBox'
import VerificationFileEntry from './VerificationFileEntry'
import Spinner from '../Spinner'
import ProxeusLogo from '../ProxeusLogo'
import mafdc from '@/mixinApp'

export default {
  mixins: [mafdc],
  name: 'verification',
  props: {
    thumbnail: {
      type: Boolean,
      default: false
    }
  },
  components: {
    FileDropBox,
    Spinner,
    VerificationFileEntry,
    ProxeusLogo
  },
  data () {
    return {
      files: [],
      hasDocuments: false,
      initialized: true
    }
  },
  created () {
  },
  computed: {
  },
  methods: {
    drop (file) {
      if (this.initialized === false) {
        throw new Error('Wallet not initialized yet')
      }
      if (this.files.find(f => file.name === f.name)) {
        return
      }
      this.hasDocuments = true
      this.files.unshift(file)
    },
    wallet () {
      return this.app.wallet
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "../../assets/styles/variables";

  /deep/ .file_drop_box {
    border: 2px dashed #dddddd;
    min-width: 155px;
    vertical-align: middle;
    border-radius: $border-radius;

    &.active {
      background: #eeeeee;
      border: 2px dashed #aaaaaa;
    }
  }

  .list-enter-active, .list-leave-active {
    transition: all 1s;
  }

  .list-move {
    transition: transform 1s;
  }

  .list-enter, .list-leave-to {
    opacity: 0;
    transform: translateY(30px);
  }

  .verification {
    overflow-y: visible;
    overflow-x: hidden;
    margin: 0 auto;
  }

  .dropbox {
    width: 100%;
    &.hasDocuments {
      margin-bottom: 2rem;
    }
  }

  .dropbox.col-sm-12 {
    padding-right: 0 !important;
  }

  @media (max-width: 768px) {
    .dropbox.col-md-6 {
      padding-right: 0 !important;
    }
  }

  .powered-by {
    display: none;
  }

  .frontend-headless {
    .powered-by {
      display: block;
    }
  }

</style>
