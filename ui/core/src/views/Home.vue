<template>
<div class="frontend">
  <vue-headful :title="$t('Home title', 'Proxeus')"/>
  <div class="container">
    <div class="row">
      <div class="col text-center">
        <h1 class="mb-3">{{$t('Proxeus Platform Title', 'Proxeus Platform')}}</h1>
        <p class="mb-3">{{$t('Proxeus Platform Description', 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry standard dummy text ever since the 1500s, when an Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry standard dummy text ever since the 1500s')}}</p>
        <a href="/register" class="btn btn-lg btn-primary">
          {{$t('Sign up', 'Login or sign up')}}
        </a>
      </div>
    </div>
    <hr class="mt-5"/>
    <div class="row mt-5">
      <div class="col text-center">
        <h2>{{$t('Document verification title', 'Document verification')}}</h2>
        <p class="mt-1 text-center">{{$t('Verify a document Description', 'Verify a document Lorem Ipsum is simply dummy text of the printing a nting and typesetting industry')}}</p>
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
