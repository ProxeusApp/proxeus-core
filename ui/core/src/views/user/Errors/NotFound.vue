<template>
<div class="e404">
  <vue-headful :title="$t('404 Not Found title', 'Proxeus - 404 Not Found')"/>
  <h1 class="title-404 text-center">{{$t('404')}}</h1>
  <h2 class="subtitle-404 text-center light-text"
      v-if="payload && payload.what">{{ payload.what }} {{$t('not Found')}}</h2>
  <h2 class="subtitle-404 text-center light-text" v-else>{{$t('Not Found')}}</h2>
</div>
</template>

<script>
export default {
  name: 'not-found',
  props: ['payload']
}
</script>

<style scoped>
  .subtitle-404 {
    font-size: 5vw;
    font-weight: 400;
    line-height: 1;
  }

  .title-404 {
    font-size: 33vw;
    font-weight: 900;
    line-height: 1;
    margin-bottom: 0;
  }
</style>
