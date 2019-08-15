<template>
<div class="e404">
  <h1 class="title-404 text-center">{{$t('Oops')}}</h1>
  <h2 class="subtitle-404 text-center light-text"
      v-if="payload && payload.what">{{ payload.what }} {{$t('not Found')}}</h2>
  <h2 class="subtitle-404 text-center light-text mt-3" v-else>{{$t('404 - Not Found')}}</h2>
  <p
    class="message-404 light-text text-center px-5">{{$t('We can\'t seem to find what you\'re looking for. Please contact our Support Team.')}}</p>
</div>
</template>

<script>
export default {
  name: 'not-found',
  props: ['payload']
}
</script>
