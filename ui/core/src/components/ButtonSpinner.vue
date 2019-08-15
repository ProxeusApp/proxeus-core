<template>
<div class="spinner spinner-btn ml-1">
  <div class="double-bounce1"></div>
  <div class="double-bounce2"></div>
</div>
</template>

<script>
export default {
  name: 'button-spinner'
}
</script>

<style scoped>
  .spinner {
    vertical-align: middle;
    width: 1rem;
    height: 1rem;
    display: inline-block;
    position: relative;
    margin: 0px auto;
  }

  .spinner.btn-spinner {
    height: 1rem;
  }

  .double-bounce1, .double-bounce2 {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    background-color: #ffffff;
    opacity: 0.6;
    position: absolute;
    top: 0;
    left: 0;

    -webkit-animation: sk-bounce 2.0s infinite ease-in-out;
    animation: sk-bounce 2.0s infinite ease-in-out;
  }

  .double-bounce2 {
    -webkit-animation-delay: -1.0s;
    animation-delay: -1.0s;
  }

  @-webkit-keyframes sk-bounce {
    0%, 100% {
      -webkit-transform: scale(0.0)
    }
    50% {
      -webkit-transform: scale(1.0)
    }
  }

  @keyframes sk-bounce {
    0%, 100% {
      transform: scale(0.0);
      -webkit-transform: scale(0.0);
    }
    50% {
      transform: scale(1.0);
      -webkit-transform: scale(1.0);
    }
  }
</style>
