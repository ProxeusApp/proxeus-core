<template>
<div class="registration-success p-3">
  <h1>Successfully Registered</h1>
  <p
    class="success-text my-3">Your account has been registered successfully. Please wait until we have confirmed the activation of your account. If you have received the account activation confirmation and you are still unable to log in, please make sure you are using the correct MetaMask account using the Ethereum Address you listed in the Access Agreement.</p>
  <a target="_blank"
     href="https://docs.google.com/document/d/1nUqB31uWRVokDOeKgzhv49rTbkLQAS8nmv8bCumDMEM/preview"
     class="btn btn-primary">Handbook
  </a>
  <a target="_blank" href="https://proxeus.com/en/privacy-policy/" class="btn btn-secondary ml-2">Privacy Policy</a>
</div>
</template>

<script>
export default {
  name: 'registration-success'
}
</script>

<style scoped>
  .registration-success {
    max-width: 700px;
    margin: 0 auto;
    text-align: center;
  }

  .success-text {
    font-size: 1.1rem;
  }
</style>
