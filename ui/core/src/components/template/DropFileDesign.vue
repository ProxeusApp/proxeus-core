<template>
<div>
  <div class="dnd-and-select-file-icon">
    <div class="outer">
      <div class="middle">
        <div class="inner">
          <div class="ar-parent">
            <div class="ar-line"></div>
            <div class="ar-pos ar-down-arrow"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="text-center pt-4 pm-1 pb-2">
    <span v-if="text" class="h5 text-hint p-0 pb-3"> {{ text }} </span>
    <span v-else class="h5 text-hint p-0 mt-3 mb-0"> {{ detail }} </span>
  </div>
</div>
</template>

<script>
export default {
  name: 'drop-file-design',
  props: {
    text: {
      type: String,
      default: ''
    },
    detail: {
      type: String,
      default: function () {
        return this.$t('choose.or.drag.file', 'Choose a file or drag it here.')
      }
    }
  }
}
</script>

<style lang="scss">
  @import "~bootstrap/scss/functions";
  @import "../../assets/styles/variables.scss";

  .dnd-and-select-file-icon {
    pointer-events: none;
    position: relative;
    width: 100%;
    margin-top: 110px;
  }

  .dnd-and-select-file-icon .outer {
    display: table;
    position: absolute;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
  }

  .dnd-and-select-file-icon .middle {
    display: table-cell;
    vertical-align: middle;
  }

  .ar-parent {
    height: 300px;
    width: 100px;
  }

  .ar-parent > * {
    box-sizing: border-box;
    transition: top 0.2s;
    -webkit-transition: top 0.2s;
  }

  .ar-pos {
    position: relative;
    margin: 0 auto;
  }

  .ar-down-arrow {
    height: 60%;
    width: 60%;
  }

  .ar-down-arrow:before, .ar-down-arrow:after {
    content: '';
    position: absolute;
    top: 0;
    width: 100%;
    border: 1px solid #dddddd;
    transition: all 0.3s;
    -webkit-transition: all 0.3s;
  }

  .ar-down-arrow:before {
    transform: rotateZ(50deg);
    -webkit-transform: rotateZ(50deg);
    left: -32%;
  }

  .ar-down-arrow:after {
    transform: rotateZ(-50deg);
    -webkit-transform: rotateZ(-50deg);
    left: 32%;
  }

  input.input-file:hover + div .ar-line:before {
    top: 8px;
  }

  input.input-file:hover + div .ar-down-arrow:before, input.input-file:hover + div .ar-down-arrow:after {
    top: 8px;
  }

  .ar-parent {
    position: absolute;
    max-height: 80px;
    width: 70px;
    top: -56px;
    display: inline-block;
    margin: 0 auto;
    left: 0;
    right: 0;
  }

  .ar-line:before {
    overflow: visible;
    content: "";
    padding-top: 100%;
  }

  .ar-line:before {
    transition: top 0.3s;
    -webkit-transition: top 0.3s;
    top: 0;
    position: relative;
    height: 100%;
    width: 1px;
    border-left: 2px solid #dddddd;
    display: inline-block;
    margin-left: 48.6%;
    margin-bottom: -24%;
  }

  .dnd-and-select-file-icon .inner {
    overflow: visible;
    position: relative;
    margin-left: auto;
    margin-right: auto;
    width: 60%;
  }

</style>
