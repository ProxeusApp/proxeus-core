<template>
<div class="frontend">
  <vue-headful :title="$t('Home title', 'Proxeus')"/>
  <div class="container">
    <div class="row">
      <div class="col text-center">
        <h1 class="mb-3">{{$t('Proxeus Platform Title', 'Proxeus Platform')}}</h1>
        <p class="mb-3">{{$t('Proxeus Platform Description', 'Proxeus is all about document-centered processes that interact with the Blockchain. Our app lets you drag & drop a workflow in the wink of an eye. Proxeus workflows may comprise forms, document templates and conditions. Proxeus enables you to create tamper-proof and easily verifiable documents by registering them on the Ethereum Blockchain.')}}</p>
        <a href="/register" class="btn btn-lg btn-primary" id="signupcontent">{{$t('Sign up', 'Sign up')}}</a>
      </div>
    </div>
    <hr class="mt-5"/>
    <div class="row mt-5">
      <div class="col text-center">
        <h2>{{$t('Document verification title', 'Document verification')}}</h2>
        <p class="mt-1 text-center">{{$t('Verify a document Description', 'Did you receive a document that has been certified by Proxeus? Verify it here!')}}</p>
        <button v-if="!displayVerification" class="btn btn-lg btn-primary mx-auto" @click="displayVerification = true">
          {{$t('Verify a document')}}
        </button>
      </div>
    </div>
    </div>
  <div class="row mt-2">
    <div class="col">
      <verification v-if="displayVerification"></verification>
    </div>
  </div>
</div>
</template>

<script>
import FileDropBox from '@/components/template/FileDropBox'
import Verification from '@/components/document/Verification'
import FrontendNavbar from '@/components/frontend/FrontendNavbar'

export default {
  name: 'home',
  components: {
    FileDropBox,
    FrontendNavbar,
    Verification
  },
  data () {
    return {
      displayVerification: false
    }
  }
}
</script>

<style lang="scss" scoped>
  .frontend {
    margin-bottom: 100px;
  }
</style>
