<template>
<div class="wallet-account">
  <ul class="nav nav-tabs">
    <li class="nav-item">
      <a class="nav-link active" data-toggle="tab" href="#main" role="tab">Active</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" data-toggle="tab" href="#transactions" role="tab">Transactions</a>
    </li>
    <li class="nav-item">
      <a class="nav-link" data-toggle="tab" href="#settings" role="tab">Settings</a>
    </li>
  </ul>
  <div class="tab-content p-3" id="myTabContent">
    <div class="tab-pane fade show active text-center" id="main" role="tabpanel" aria-labelledby="home-tab">
      <h1 class="text-center">Your account</h1>
      <small class="account-hash">
        <div class="form-group light-text">
          <strong class="hash-title">Hash:</strong>
          <span class="hash">0x87d786366c650fa09e77deb5d42e293edb9b6478</span>
        </div>
        <div class="balance mb-3">
          <span class="xes xes-nr">{{ balance }}</span>
          <!--<span class="xes">XES</span>-->
          <!--<span class="eth">/ 1.422 ETH</span>-->
        </div>
        <hr>
        <button type="button" class="btn btn-primary m-1">Deposit</button>
        <button type="button" class="btn btn-primary m-1">Send</button>
      </small>
    </div>
    <div class="tab-pane fade" id="transactions" role="tabpanel" aria-labelledby="profile-tab">
      <transaction-history></transaction-history>
    </div>
    <div class="tab-pane fade" id="settings" role="tabpanel" aria-labelledby="contact-tab">...</div>
  </div>
</div>
</template>

<script>
import TransactionHistory from './TransactionHistory'

export default {
  name: 'account',
  props: ['balance'],
  components: {
    TransactionHistory
  }
}
</script>

<style lang="scss" scoped>
  @use "../../assets/styles/variables";

  .account-hash {
    word-wrap: break-word;
  }

  .nav-tabs {
    background: theme-color("info");

    .nav-link {
      color: white;
      &.active {
        color: theme-color("primary");
      }
    }
  }

  .xes {
    font-size: 2rem;
    font-weight: 100;

    &.xes-nr {
      font-weight: 400;
    }
  }

  .eth {
    font-size: 1.5rem;
    font-weight: 100;
  }

  .hash {
    max-width: 100%;
    display: inline-block;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
  }

  .hash-title {
    overflow: hidden;
    display: inline-block;
  }
</style>
