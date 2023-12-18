<template>
<div class="view-support">
  <vue-headful :title="$t('Help / Support title', 'Proxeus - Help / Support')"/>
  <top-nav :title="$t('Help / Support')">
  </top-nav>
  <div class="container-fluid pt-3">
    <div class="row mb-3">
      <div class="col-sm-12">
        <section class="pl-2">
          <p class="mt-1">{{$t('Support description paragraph', 'More detailed information on how to use the Proxeus platform can be found in the Handbook.')}}</p>
          <a target="_blank"
            :href="$t('Support Button Link', 'https://doc.proxeus.org/#/handbook')"
             class="btn btn-primary-dark btn-lg mt-1">{{$t('Handbook')}}
          </a>
        </section>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import TopNav from '../components/layout/TopNav.vue'

export default {
  name: 'support',
  components: {
    TopNav
  }
}
</script>

<style scoped>

</style>
