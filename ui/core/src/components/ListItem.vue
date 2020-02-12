<template>
  <tr class="mlist-item" :class="{'animate-new-entry':element.isNew}" :data-index="index"
      @click="goToLink()">
    <td class="tdmin">
      <img v-if="element.photo" class="mlist-group-icon" :src="element.photo"/>
      <i v-else class="mlist-group-icon material-icons" :class="[iconFa ? iconFa : '']">{{ icon || ''}}</i>
    </td>
    <td>
      <div style="display:block;" class="fregular">{{ element.name||defaultName||'-' }}
        <slot name="nameExtra"/>
      </div>
      <small style="display:block;" class="light-text fregular" v-if="element.email">{{element.email}}</small>
      <small style="display:block;" class="light-text" v-if="element.etherPK">{{ element.etherPK }}</small>
      <small style="display:block;" class="light-text" v-else-if="element.detail">{{ element.detail }}</small>
      <small style="display:block;" class="light-text" v-else>-</small>

      <small style="display:block;" class="error" v-if="error">{{error}}</small>
    </td>
    <td v-if="$parent.lastExs || $parent.lastImps">
      <div v-if="$parent.lastImps && $parent.lastImps.hasOwnProperty(element.id)" class="easy-read"><span class="badge"
                                                                                                          :class="$parent.lastImps[element.id] === ''?'badge-info':'badge-danger'">{{$t('Imported')}}</span>
      </div>
      <div v-if="$parent.lastExs && $parent.lastExs.hasOwnProperty(element.id)" class="easy-read"><span class="badge"
                                                                                                        :class="$parent.lastExs[element.id] === ''?'badge-info':'badge-danger'">{{$t('Exported')}}</span>
      </div>
    </td>
    <td v-if="timestamps">
      <div class="easy-read stime">{{ element.created | moment('DD.MM.YY - HH:mm') }}</div>
      <small class="light-text">{{$t('Created')}}</small>
    </td>
    <td v-if="timestamps">
      <div class="easy-read stime">{{ element.updated | moment('DD.MM.YY - HH:mm') }}</div>
      <small class="light-text">{{$t('Updated')}}</small>
    </td>
    <td class="pl-2" style="text-align: right;">
      <div v-if="price" class="easy-read stime">{{ this.price }} XES</div>
    </td>
    <td class="tdmin align-items-center">
      <div class="d-flex align-items-center">
        <slot/>
      </div>
    </td>
  </tr>
</template>

<script>

export default {
  name: 'list-item',
  created () {
    window.$parent = this.$parent
  },
  props: {
    index: Number,
    element: Object,
    to: null,
    error: String,
    iconFa: String,
    icon: String,

    defaultName: {
      type: String,
      default: ''
    },
    timestamps: {
      type: Boolean,
      default: true
    },
    _blank: {
      type: Boolean,
      default: false
    },
    price: Number
  },
  methods: {
    goToLink () {
      let p = typeof this.to === 'string' ? { path: this.to } : this.to
      if (this._blank) {
        let routeData = this.$router.resolve(p)
        window.open(routeData.href, '_blank')
        return
      }
      this.$router.push(p)
    }
  }
}
</script>

<style lang="scss">
  .mlist-item .error {
    color: #dc3545;
    background-color: transparent;
  }
</style>
