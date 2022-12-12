<template>
<div class="verification-file-entry" v-bind="$attrs">
  <div v-if="singleFile" class="text-center p-3 border border-2 border-light">
    <div class="pb-1">
      <spinner v-show="loading" background="transparent" color="#eee" :margin="10" cls="position-relative"></spinner>
      <i v-show="valid" class="text-success material-icons md-60 mdi mdi-check-circle"></i>
      <i v-show="notFound" class="text-danger material-icons md-60 mdi mdi-close-circle"></i>
      <i v-show="invalid || errorValidating" class="text-danger material-icons md-60 mdi mdi-alert-circle"></i>
    </div>
    <div v-if="!loading">
      <div class="break-word mx-auto mt-2 mb-0">
        <h4 v-show="valid"
            class="text-success">{{ $t('Verified file is valid', 'The file {filename} is valid.', {filename: file.name}) }}</h4>
        <h4 v-show="invalid"
            class="text-danger">{{ $t('Verified file revoked', 'The file {filename} has been revoked.', {filename: file.name}) }}</h4>
        <h4 v-show="notFound"
            class="text-danger">{{ $t('Verified file is invalid', 'The file {filename} is invalid.', {filename: file.name}) }}</h4>
      </div>
      <div class="text-center mt-2">
        <div v-show="valid" class="mb-2">
          <span class="text-hint">{{ $t('Issued by') }} </span>
          <a class="break-word font-weight-bold" :title="creator" target="_blank"
             :href="'https://' + network + 'etherscan.io/address/' + creator" v-if="contract">
            {{ creator | addressOrHash}}
            <i class="material-icons md-14 mdi mdi-launch"></i>
          </a>
        </div>
        <div v-if="showDetails === true || !valid">
          <hr v-show="!errorValidating" class="w-75"/>
          <div class="container text-md-left text-center">
            <div v-show="timestamp" class="row mb-2">
              <div
                class="offset-md-2 col-md-3 pr-0 text-hint">{{ $t('Verified file timestamp', 'Registration date:') }}
              </div>
              <div class="col-md-5">{{ timestamp }}</div>
            </div>
            <div v-show="valid && hash" class="row mb-2">
              <div class="offset-md-2 col-md-3 pr-0 text-hint">{{ $t('Verified file hash', 'File Hash:') }}</div>
              <div class="col-md-5 break-word">{{ hash | addressOrHash }}
              </div>
            </div>
            <div class="row mb-2" v-show="contract">
              <div
                class="offset-md-2 col-md-3 pr-0 text-hint">{{ $t('Verified file contract addr', 'Contract Address:') }}
              </div>
              <div class="col-md-5">
                <a class="break-word" :title="contract" target="_blank"
                   :href="'https://' + network + 'etherscan.io/address/' + contract" v-if="contract">
                  {{ contract | addressOrHash}}
                  <i class="material-icons md-14 mdi mdi-launch"></i>
                </a>
              </div>
            </div>
            <div class="row" v-show="tx">
              <div class="offset-md-2 col-md-3 pr-0 text-hint">{{ $t('Verified file transaction', 'Transaction:') }}
              </div>
              <div class="col-md-5">
                <a class="break-word" :title="tx" target="_blank" :href="'https://' + network + 'etherscan.io/tx/' + tx"
                   v-if="tx">
                  {{ tx | addressOrHash}}
                  <i class="material-icons md-14 mdi mdi-launch"></i>
                </a>
              </div>
            </div>
            <div v-show="signatures.length > 0" class="row mb-2">
              <div class="offset-md-2 col-md-3 pr-0 text-hint">{{ $t('Verified file signers', 'Signers:') }}</div>
              <div class="col-md-5">
                <div class="signer" v-for="signature in signatures">
                  {{ signature.address | addressOrHash}} {{$t('on')}}
                  <a class="break-word text-hint" :title="signature.txHash" target="_blank"
                     :href="'https://' + network + 'etherscan.io/tx/' + signature.txHash" v-if="signature.txHash">
                    {{ signature.time }}
                    <i class="material-icons md-14 mdi mdi-launch"></i>
                  </a>
                </div>
              </div>
            </div>
          </div>

          <div class="text-hint w-100 px-md-5 mt-2" :class="{'pt-2' : valid, 'pt-0' : !valid}">
            <p v-show="valid">
              {{ $t('Verified file is valid 2', 'This file has been recognized as genuine and valid.') }}
            </p>
            <p
              v-show="notFound">{{ $t('File not registered explanation', 'No entry was found for this file. This could mean it was never registered or it was manipulated. If you believe that this finding is an error, please contact the issuer.') }}</p>
            <p
              v-show="invalid">{{ $t('Verified file invalid explanation', 'This file has been recognised as authentic but has since been declared invalid. It may have expired or been recalled by the issuer.') }}</p>
            <p
              v-show="errorValidating">{{ $t('Error while verifying file', 'File hash could not been verified due to technical problems. Please try again.') }}</p>
          </div>

          <hr class="w-75" v-show="valid"/>
        </div>
      </div>
      <button v-if="valid" type="button" @click="showDetails = !showDetails"
              class="more-info-variant btn btn-light-round mt-3 mb-2 text-hint">
        <i class="material-icons md-16 mdi"
           :class="{'mdi-arrow-down-drop-circle-outline': !showDetails, 'mdi-arrow-up-drop-circle-outline' : showDetails}"></i>
        {{showDetails === false ? $t('Show Details') : $t('Hide Details')}}
      </button>
    </div>

  </div>
  <div v-else class="border border-2 border-light mb-2 container">
    <div class="row list-item">
      <div class="bg-light col-lg-1 px-1 status-holder text-center">
        <div class="check py-2">
          <spinner v-show="loading" background="transparent" color="#eee" :margin="16"
                   cls="position-relative"></spinner>
          <i v-show="notFound" class="text-danger material-icons md-36 mdi mdi-close-circle"></i>
          <i v-show="invalid || errorValidating" class="text-danger material-icons md-36 mdi mdi-alert-circle"></i>
          <i v-show="loading === false && status" class="text-success material-icons md-36 mdi mdi-check-circle"></i>
        </div>
      </div>
      <div class="py-2 pl-3 pr-2 col-lg-8">
        <div>
          <p class="filename mb-0">{{ file.name }}</p>
          <div v-show="!loading">
            <small v-show="valid"
                   class="text-success">{{ $t('Verified file is valid', 'The file {filename} is valid.', {filename: file.name}) }}
            </small>
            <small v-show="invalid"
                   class="text-danger">{{ $t('Verified file revoked', 'The file {filename} has been revoked.', {filename: file.name}) }}
            </small>
            <small v-show="notFound"
                   class="text-danger">{{ $t('Verified file is invalid', 'The file {filename} is invalid.', {filename: file.name}) }}
            </small>
          </div>
        </div>
        <div v-show="valid">
          <small class="text-hint">{{ $t('Issued by') }}</small>
          <a class="break-word font-weight-bold small" :title="creator" target="_blank"
             :href="'https://' + network + 'etherscan.io/address/' + creator" v-if="contract">
            {{ creator }}
            <i class="material-icons md-14 mdi mdi-launch"></i>
          </a>
        </div>

        <small class="text-hint" v-show="notFound">
          {{ $t('File not registered explanation', 'No entry was found for this file. This could mean it was never registered or it was manipulated. If you believe that this finding is an error, please contact the issuer.') }}
        </small>
        <small class="text-hint" v-show="invalid">
          {{ $t('Verified file invalid explanation', 'This file has been recognised as authentic but has since been declared invalid. It may have expired or been recalled by the issuer.') }}
        </small>
        <small class="text-hint" v-show="errorValidating">
          {{ $t('Error while verifying file', 'File could not been verified due to technical problems. Please try again.') }}
        </small>

        <small v-show="showDetails" class="text-hint text-left">

          <hr class="my-2"/>
          <div class="container mt-2">
            <div v-show="timestamp" class="row mb-2 ">
              <div class="col-lg-4 pr-0 pl-0">{{ $t('Verified file timestamp', 'Registration date:') }}</div>
              <div class="col-lg-8 pl-0">{{ timestamp }}</div>
            </div>
            <div v-show="valid && hash" class="row mb-2 ">
              <div class="col-lg-4 pr-0 pl-0 ">{{ $t('Verified file hash', 'File Hash:') }}</div>
              <div class="col-lg-8 pl-0 break-word">{{ hash | addressOrHash }}
              </div>
            </div>
            <div class="row mb-2" v-show="contract">
              <div class="col-lg-4 pl-0 pr-0 text-hint">{{ $t('Verified file contract addr', 'Contract Address:') }}
              </div>
              <div class="col-lg-8 pl-0">
                <a class="break-word text-hint" :title="contract" target="_blank"
                   :href="'https://' + network + 'etherscan.io/address/' + contract" v-if="contract">
                  {{ contract | addressOrHash}}
                  <i class="material-icons md-14 mdi mdi-launch"></i>
                </a>
              </div>
            </div>
            <div class="row" v-show="tx">
              <div class="col-lg-4 pr-0 pl-0 text-hint ">{{ $t('Verified file transaction', 'Transaction:') }}</div>
              <div class="col-lg-8 pl-0">
                <a class="break-word text-hint" :title="tx" target="_blank"
                   :href="'https://' + network + 'etherscan.io/tx/' + tx" v-if="tx">
                  {{ tx | addressOrHash}}
                  <i class="material-icons md-14 mdi mdi-launch"></i>
                </a>
              </div>
            </div>
          </div>
        </small>
      </div>
      <div class="col-lg-3" v-if="valid">
        <div class="py-4 pr-1 text-center text-lg-right">
          <button type="button" @click="showDetails = !showDetails"
                  class="more-info-variant btn btn-light-round btn-sm my-0 text-hint">
            <i class="material-icons md-16 mdi"
               :class="{'mdi-arrow-down-drop-circle-outline': !showDetails, 'mdi-arrow-up-drop-circle-outline' : showDetails}"></i>
            {{showDetails === false ? $t('Show Details') : $t('Hide Details')}}
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import Spinner from '../Spinner'

export default {
  name: 'verification-file-entry',
  props: ['file', 'wallet', 'singleFile'],
  components: {
    Spinner
  },
  data () {
    return {
      loading: false,
      status: false,
      isFileInvalidated: undefined,
      tx: undefined,
      hash: undefined,
      network: process.env.VUE_APP_BLOCKCHAIN_NET === 'goerli' ? 'goerli.' : '',
      timestamp: undefined,
      creator: undefined,
      signatures: [],
      contract: undefined,
      errorValidating: false,
      showDetails: false
    }
  },
  mounted () {
    this.verify()
  },
  computed: {
    valid () {
      return this.errorValidating === false && this.loading === false && this.loading === false && this.status === true
    },
    invalid () {
      return this.errorValidating === false && this.loading === false && this.status === false &&
        this.isFileInvalidated === true
    },
    notFound () {
      return this.errorValidating === false && this.loading === false && this.status === false &&
        this.isFileInvalidated !== true
    }
  },
  filters: {
    addressOrHash: function (value) {
      if (typeof value !== 'string' || value.length === 0) {
        return value
      }
      let shortValue = value.substring(0, 11)
      shortValue += '...' + value.substring(value.length - 11, value.length)
      return shortValue
    }
  },
  methods: {
    async verify () {
      this.loading = true
      try {
        // use hash if it already comes with the file object
        if (this.file.hash) {
          this.hash = this.file.hash
        } else {
          // calculate hash from file object
          this.hash = await this.hashFile(this.file)
          this.hash = '0x' + this.hash
        }

        this.isFileInvalidated = false

        let result
        try {
          result = await this.wallet().verifyHash(this.hash)
        } catch (e) {
          result = false
        }
        if (result) {
          const transaction = await this.wallet().web3.eth.getTransaction(result)
          this.creator = transaction.from
          this.contract = transaction.to
          const block = await this.wallet().web3.eth.getBlock(transaction.blockNumber)
          // *1000 is conversion to seconds
          this.timestamp = (new Date(block.timestamp * 1000)).toUTCString()

          const signersArr = await this.wallet().proxeusFS.contract.methods.getFileSigners(this.hash).call()
          Promise.all(signersArr.map(async signerAddr => {
            const registrationTxBlock = await this.wallet().getRegistrationTxBlock(signerAddr)
            return {
              address: signerAddr,
              txHash: registrationTxBlock.txHash,
              block: block,
              time: (new Date(registrationTxBlock.block.timestamp * 1000)).toUTCString()
            }
          })).then((sn) => {
            this.signatures = sn
          })

          this.loading = false
          this.tx = result
          this.status = true
          // emit event to alert parent component of status
          this.$emit('updateFileState', true)
        } else {
          this.status = false
          this.tx = null
          this.loading = false
        }
      } catch (e) {
        this.status = false
        this.loading = false
        this.tx = null
        this.errorValidating = true
      }
    },
    async hashFile (file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()

        reader.onload = (e) => {
          const hash = this.wallet().hashFile(reader.result)
          resolve(hash)
        }

        reader.onerror = (e) => {
          reject(e)
        }

        reader.readAsArrayBuffer(file)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "../../assets/styles/variables";

  .break-word {
    word-wrap: break-word;
  }

  .verification-file-entry {
    border-radius: $border-radius;

    ::v-deep .spinner {
      min-height: auto !important;
    }

    ::v-deep .sk-circle {
      width: 36px;
      height: 36px;
    }

    .list-item {
      min-height: 96px;
    }

    .filename {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      word-wrap: break-word;
    }

    .status-holder {
      .material-icons {
        position: relative;
        /*top: 14px;*/
      }
    }

    .btn {
      i.material-icons {
        position: relative;
        top: -2px;
        left: 1px;
      }
    }
  }

  .material-icons.md-60 {
    font-size: 60px;
    line-height: 40px;
    height: 53px;
  }
</style>
