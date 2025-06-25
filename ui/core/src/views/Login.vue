<template>
  <div>
    <!--<sign-message id="sign-message"></sign-message>-->
    <vue-headful :title="$t('Proxeus - Sign in', 'Proxeus - Log in')" />
    <h1 class="text-center">{{ $t("Sign in title", "Log in") }}</h1>
    <div
      class="login-form container-fluid px-4 pt-2 mt-3 bg-light"
      :class="{ 'login-form-sm': !metamaskLoginAvailable }"
    >
      <div class="row">
        <div
          class="col-12 mt-3 text-center d-flex align-items-center"
          v-if="app.wallet && metamaskLoginAvailable"
        >
          <div class="mid align-self-center w-100" v-show="signing === false">
            <h2 class="mb-3 font-weight-bold">{{ $t("Account signature") }}</h2>
            <p class="light-text">
              {{
                $t(
                  "Account to sign in explanation",
                  "Use your MetaMask account to log in."
                )
              }}
            </p>
            <p class="text-danger" v-if="walletErrorMessage">
              {{ walletErrorMessage }}
            </p>
            <button
              class="btn btn-primary px-3"
              @click="metamaskLogin"
              v-if="metamaskLoginAvailable"
            >
              {{ $t("Sign in with signature", "Log in with signature") }}
            </button>
          </div>
          <div class="mid align-self-center w-100" v-show="signing">
            <spinner
              background="transparent"
              style="position: relative"
              :margin="60"
            ></spinner>
            <p class="text-primary">
              {{
                $t(
                  "Account signature explanation",
                  "Please sign with MetaMask. If MetaMask doesn't pop up automatically, click on the MetaMask icon in the top right corner of your browser."
                )
              }}
            </p>
          </div>
        </div>
      </div>
    </div>
    <div class="modal fade" ref="tcModal" tabindex="-1" role="dialog">
      <div class="modal-dialog modal-lg" role="document">
        <div class="modal-content">
          <iframe
            frameborder="0"
            width="100%"
            height="100%"
            :src="$t('Terms & Conditions link', '')"
          ></iframe>
          <div class="modal-footer">
            <button type="button" @click="acceptTC" class="btn btn-primary">
              {{ $t("Accept") }}
            </button>
            <button
              type="button"
              class="btn btn-secondary"
              data-dismiss="modal"
            >
              {{ $t("Cancel") }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Spinner from "../components/Spinner";
import mafdc from "@/mixinApp";

export default {
  mixins: [mafdc],
  name: "login",
  components: {
    Spinner,
  },
  data() {
    return {
      account: undefined,
      pwlogin: true,
      email: "",
      password: "",
      hasError: false,
      walletErrorMessage: "",
      loginErrorMessage: "",
      metamaskLoginAvailable: false,
      signing: false,
    };
  },
  created() {
    if (typeof window.ethereum !== "undefined") {
      this.metamaskLoginAvailable = true;
    }
  },
  mounted() {
    this.$refs.inputEmail && this.$refs.inputEmail.focus();
  },
  computed: {
    wallet() {
      return this.app.wallet;
    },
  },
  methods: {
    login() {
      this.account = this.email;
      this.pwlogin = true;
      if (this.checkTermsAndConditions()) {
        axios
          .post("/api/login", { email: this.email, password: this.password })
          .then(
            (res) => {
              this.hasError = false;
              window.location = res.data.location || "/admin/workflow";
            },
            (err) => {
              this.loginErrorMessage =
                "You have entered an invalid username or password";
              this.hasError = true;
              this.$nextTick(() => {
                this.$refs.inputEmail.focus();
              });
              this.app.handleError(err);
            }
          );
      }
    },
    checkTermsAndConditions() {
      if (
        this.$t("Terms & Conditions link", "") === "Terms & Conditions link"
      ) {
        return true;
      }
      const rememberAccept = localStorage.getItem("acc_" + this.account);
      if (rememberAccept && rememberAccept === "yes") {
        return true;
      }
      $(this.$refs.tcModal).modal("show");
      $(this.$refs.tcModal).on("hide.bs.modal", (e) => {
        this.declineTC();
      });
    },
    declineTC() {
      this.signing = false;
    },
    acceptTC() {
      localStorage.setItem("acc_" + this.account, "yes");
      this.app.acknowledgeFirstLogin();
      $(this.$refs.tcModal).modal("hide");
      if (this.pwlogin) {
        this.login();
      } else {
        this.metamaskLogin();
      }
    },
    async metamaskLogin() {
      if (window.ethereum) {
        try {
          await window.ethereum.enable();
          await this.wallet.wallet.setupDefaultAccount();
        } catch (e) {
          this.walletErrorMessage = "Please grant access to MetaMask.";
          return;
        }
      }
      axios.get("/api/challenge").then(
        (response) => {
          this.signing = true;
          const challenge = response.data;
          const account = this.wallet.getCurrentAddress();
          this.account = account;
          if (this.account === undefined) {
            this.walletErrorMessage = "Please log in to MetaMask.";
            this.signing = false;
            return;
          }
          this.pwlogin = false;
          if (this.checkTermsAndConditions()) {
            console.log(this.account);
            this.wallet
              .signMessage(challenge, this.account)
              .then((signature) => {
                axios.post("/api/login", { signature }).then(
                  (res) => {
                    if (res.status >= 200 && res.status <= 299) {
                      window.location = res.data.location || "/admin/workflow";
                    } else {
                      this.walletErrorMessage = "Could not verify signature.";
                    }
                  },
                  (err) => {
                    this.signing = false;
                    this.walletErrorMessage = "Could not verify signature.";
                    this.app.handleError(err); // TODO: check if this is working
                  }
                );
              })
              .catch(() => {
                this.signing = false;
                this.walletErrorMessage = "Could not Sign Message.";
              });
          }
        },
        (err) => {
          this.app.handleError(err);
          this.$notify({
            group: "app",
            title: "Error",
            text: "Login failed. Please try again or if the error persists contact the platform operator.",
            type: "warning",
          });
        }
      );
    },
  },
};
</script>

<style lang="scss" scoped>
@use "@/assets/styles/variables.scss";

@media print {
  * {
    display: block !important;
    visibility: visible !important;
    opacity: 1 !important;
    height: 100% !important;
    min-height: 100% !important;
    overflow: visible !important;
  }
  .modal-body {
    max-height: 100%;
    overflow: visible;
  }

  .modal-backdrop {
    display: none;
  }

  .modal {
    position: relative;
    display: block;
    opacity: 1;
    overflow: visible;
  }

  .modal-open {
    overflow: visible;
    overflow: visible;
  }

  .modal-dialog {
    transform: none;
    overflow: visible;
  }
}

.modal-body {
  h1 {
    font-size: $h1-font-size;
  }

  h2 {
    font-size: $h2-font-size;
  }

  h3 {
    font-size: $h3-font-size;
  }

  ol {
    padding-left: 0;
    margin-bottom: 0;
  }
}

h2 {
  padding-top: 1rem;
  font-size: 1.5rem !important;
}

h3 {
  padding-top: 0.5rem;
  font-size: 1.2rem !important;
}

.c10:before,
.c7:before {
  font-size: $h1-font-size;
}

.c11:before {
  font-size: $h2-font-size;
}

.modal-header {
  background: $light;
}

.modal-footer {
  background: $light;
  border-bottom-right-radius: $border-radius;
  border-bottom-left-radius: $border-radius;
}

.modal-body {
  overflow-y: auto;
  max-height: calc(100vh - 280px);
}

.login-form {
  overflow: auto;
  margin: 0 auto;
  margin-top: 50px;
  height: 100%;
  max-width: 600px;
  padding-top: 40px;
  padding-bottom: 40px;
  border-radius: $border-radius;
}

.login-form-sm {
  max-width: 350px;
}

.form-signin {
  width: 100%;
  max-width: 330px;
  padding: 2rem;
  margin: 0 auto;
  z-index: 1000;

  .checkbox {
    font-weight: 400;
  }

  .form-control {
    position: relative;
    box-sizing: border-box;
    height: auto;
    padding: 10px;
    font-size: 16px;
  }
}

.tc-list {
  margin-bottom: 0;
  margin-top: 1.5rem;
  li {
    padding-left: 0.5rem;
    padding-bottom: 0.75rem;
  }
}

.modal-body {
  ol,
  li {
    list-style-type: none;
  }

  .lst-kix_list_33-5 > li {
    counter-increment: lst-ctn-kix_list_33-5;
  }

  .lst-kix_list_21-8 > li {
    counter-increment: lst-ctn-kix_list_21-8;
  }

  .lst-kix_list_5-0 > li {
    counter-increment: lst-ctn-kix_list_5-0;
  }

  ol.lst-kix_list_2-3.start {
    counter-reset: lst-ctn-kix_list_2-3 0;
  }

  ol.lst-kix_list_37-2.start {
    counter-reset: lst-ctn-kix_list_37-2 0;
  }

  .lst-kix_list_45-2 > li {
    counter-increment: lst-ctn-kix_list_45-2;
  }

  .lst-kix_list_42-1 > li:before {
    content: "" counter(lst-ctn-kix_list_42-1, decimal) " ";
  }

  .lst-kix_list_40-1 > li {
    counter-increment: lst-ctn-kix_list_40-1;
  }

  ol.lst-kix_list_5-3.start {
    counter-reset: lst-ctn-kix_list_5-3 0;
  }

  .lst-kix_list_38-6 > li {
    counter-increment: lst-ctn-kix_list_38-6;
  }

  .lst-kix_list_4-3 > li {
    counter-increment: lst-ctn-kix_list_4-3;
  }

  .lst-kix_list_42-7 > li:before {
    content: "" counter(lst-ctn-kix_list_42-7, decimal) " ";
  }

  .lst-kix_list_51-1 > li {
    counter-increment: lst-ctn-kix_list_51-1;
  }

  ol.lst-kix_list_17-1.start {
    counter-reset: lst-ctn-kix_list_17-1 0;
  }

  .lst-kix_list_42-5 > li:before {
    content: "" counter(lst-ctn-kix_list_42-5, decimal) " ";
  }

  .lst-kix_list_42-3 > li:before {
    content: "" counter(lst-ctn-kix_list_42-3, decimal) " ";
  }

  .lst-kix_list_32-8 > li {
    counter-increment: lst-ctn-kix_list_32-8;
  }

  .lst-kix_list_24-8 > li:before {
    content: "" counter(lst-ctn-kix_list_24-8, lower-roman) ". ";
  }

  .lst-kix_list_24-2 > li:before {
    content: "(" counter(lst-ctn-kix_list_24-2, lower-latin) ") ";
  }

  ol.lst-kix_list_54-3.start {
    counter-reset: lst-ctn-kix_list_54-3 0;
  }

  .lst-kix_list_24-4 > li:before {
    content: "" counter(lst-ctn-kix_list_24-4, decimal) ") ";
  }

  ol.lst-kix_list_34-2.start {
    counter-reset: lst-ctn-kix_list_34-2 0;
  }

  .lst-kix_list_9-4 > li {
    counter-increment: lst-ctn-kix_list_9-4;
  }

  .lst-kix_list_24-6 > li:before {
    content: "" counter(lst-ctn-kix_list_24-6, lower-roman) ") ";
  }

  ol.lst-kix_list_14-1.start {
    counter-reset: lst-ctn-kix_list_14-1 0;
  }

  .lst-kix_list_23-6 > li:before {
    content: "" counter(lst-ctn-kix_list_23-6, decimal) ". ";
  }

  .lst-kix_list_23-2 > li:before {
    content: "" counter(lst-ctn-kix_list_23-2, lower-roman) ") ";
  }

  .lst-kix_list_23-0 > li:before {
    content: "" counter(lst-ctn-kix_list_23-0, decimal) ") ";
  }

  .lst-kix_list_23-8 > li:before {
    content: "" counter(lst-ctn-kix_list_23-8, lower-roman) ". ";
  }

  .lst-kix_list_3-6 > li {
    counter-increment: lst-ctn-kix_list_3-6;
  }

  ol.lst-kix_list_44-0.start {
    counter-reset: lst-ctn-kix_list_44-0 0;
  }

  .lst-kix_list_24-0 > li:before {
    content: "Article " counter(lst-ctn-kix_list_24-0, upper-roman) ". ";
  }

  .lst-kix_list_23-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_23-4, lower-latin) ") ";
  }

  ol.lst-kix_list_51-3.start {
    counter-reset: lst-ctn-kix_list_51-3 0;
  }

  .lst-kix_list_43-3 > li:before {
    content: "" counter(lst-ctn-kix_list_43-3, decimal) " ";
  }

  ol.lst-kix_list_31-2.start {
    counter-reset: lst-ctn-kix_list_31-2 0;
  }

  .lst-kix_list_22-2 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) ". ";
  }

  .lst-kix_list_22-6 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) "."
      counter(lst-ctn-kix_list_22-3, decimal) "."
      counter(lst-ctn-kix_list_22-4, decimal) "."
      counter(lst-ctn-kix_list_22-5, decimal) "."
      counter(lst-ctn-kix_list_22-6, decimal) ". ";
  }

  .lst-kix_list_43-7 > li:before {
    content: "" counter(lst-ctn-kix_list_43-7, decimal) " ";
  }

  .lst-kix_list_43-1 > li:before {
    content: "" counter(lst-ctn-kix_list_43-1, decimal) " ";
  }

  .lst-kix_list_22-0 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) ". ";
  }

  .lst-kix_list_22-8 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) "."
      counter(lst-ctn-kix_list_22-3, decimal) "."
      counter(lst-ctn-kix_list_22-4, decimal) "."
      counter(lst-ctn-kix_list_22-5, decimal) "."
      counter(lst-ctn-kix_list_22-6, decimal) "."
      counter(lst-ctn-kix_list_22-7, decimal) "."
      counter(lst-ctn-kix_list_22-8, decimal) ". ";
  }

  ol.lst-kix_list_22-5.start {
    counter-reset: lst-ctn-kix_list_22-5 0;
  }

  .lst-kix_list_5-7 > li {
    counter-increment: lst-ctn-kix_list_5-7;
  }

  ol.lst-kix_list_34-7.start {
    counter-reset: lst-ctn-kix_list_34-7 0;
  }

  .lst-kix_list_43-5 > li:before {
    content: "" counter(lst-ctn-kix_list_43-5, decimal) " ";
  }

  .lst-kix_list_22-4 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) "."
      counter(lst-ctn-kix_list_22-3, decimal) "."
      counter(lst-ctn-kix_list_22-4, decimal) ". ";
  }

  ol.lst-kix_list_32-8.start {
    counter-reset: lst-ctn-kix_list_32-8 0;
  }

  ol.lst-kix_list_25-5.start {
    counter-reset: lst-ctn-kix_list_25-5 0;
  }

  .lst-kix_list_41-7 > li:before {
    content: "" counter(lst-ctn-kix_list_41-7, decimal) " ";
  }

  ol.lst-kix_list_24-4.start {
    counter-reset: lst-ctn-kix_list_24-4 0;
  }

  .lst-kix_list_6-4 > li {
    counter-increment: lst-ctn-kix_list_6-4;
  }

  .lst-kix_list_41-1 > li:before {
    content: "" counter(lst-ctn-kix_list_41-1, decimal) " ";
  }

  .lst-kix_list_40-7 > li:before {
    content: "" counter(lst-ctn-kix_list_40-7, decimal) " ";
  }

  ol.lst-kix_list_27-4.start {
    counter-reset: lst-ctn-kix_list_27-4 0;
  }

  .lst-kix_list_40-5 > li:before {
    content: "" counter(lst-ctn-kix_list_40-5, decimal) " ";
  }

  ol.lst-kix_list_15-2.start {
    counter-reset: lst-ctn-kix_list_15-2 0;
  }

  ol.lst-kix_list_42-1.start {
    counter-reset: lst-ctn-kix_list_42-1 0;
  }

  .lst-kix_list_41-5 > li:before {
    content: "" counter(lst-ctn-kix_list_41-5, decimal) " ";
  }

  ol.lst-kix_list_37-7.start {
    counter-reset: lst-ctn-kix_list_37-7 0;
  }

  ol.lst-kix_list_2-8.start {
    counter-reset: lst-ctn-kix_list_2-8 0;
  }

  .lst-kix_list_41-3 > li:before {
    content: "" counter(lst-ctn-kix_list_41-3, decimal) " ";
  }

  ol.lst-kix_list_5-8.start {
    counter-reset: lst-ctn-kix_list_5-8 0;
  }

  .lst-kix_list_1-3 > li {
    counter-increment: lst-ctn-kix_list_1-3;
  }

  .lst-kix_list_40-3 > li:before {
    content: "" counter(lst-ctn-kix_list_40-3, decimal) " ";
  }

  ol.lst-kix_list_12-2.start {
    counter-reset: lst-ctn-kix_list_12-2 0;
  }

  .lst-kix_list_40-1 > li:before {
    content: "" counter(lst-ctn-kix_list_40-1, decimal) " ";
  }

  .lst-kix_list_42-2 > li {
    counter-increment: lst-ctn-kix_list_42-2;
  }

  ol.lst-kix_list_38-3.start {
    counter-reset: lst-ctn-kix_list_38-3 0;
  }

  .lst-kix_list_24-8 > li {
    counter-increment: lst-ctn-kix_list_24-8;
  }

  ol.lst-kix_list_3-4.start {
    counter-reset: lst-ctn-kix_list_3-4 0;
  }

  .lst-kix_list_21-8 > li:before {
    content: "" counter(lst-ctn-kix_list_21-8, decimal) ". ";
  }

  ol.lst-kix_list_18-2.start {
    counter-reset: lst-ctn-kix_list_18-2 0;
  }

  .lst-kix_list_26-8 > li:before {
    content: "" counter(lst-ctn-kix_list_26-8, lower-roman) ". ";
  }

  .lst-kix_list_21-0 > li:before {
    content: "" counter(lst-ctn-kix_list_21-0, decimal) ". ";
  }

  .lst-kix_list_13-1 > li {
    counter-increment: lst-ctn-kix_list_13-1;
  }

  .lst-kix_list_26-4 > li:before {
    content: "" counter(lst-ctn-kix_list_26-4, decimal) ") ";
  }

  ol.lst-kix_list_36-1.start {
    counter-reset: lst-ctn-kix_list_36-1 0;
  }

  .lst-kix_list_42-4 > li {
    counter-increment: lst-ctn-kix_list_42-4;
  }

  .lst-kix_list_21-4 > li:before {
    content: "" counter(lst-ctn-kix_list_21-4, decimal) ". ";
  }

  .lst-kix_list_26-0 > li:before {
    content: "Article " counter(lst-ctn-kix_list_26-0, upper-roman) ". ";
  }

  .lst-kix_list_31-4 > li {
    counter-increment: lst-ctn-kix_list_31-4;
  }

  .lst-kix_list_31-2 > li {
    counter-increment: lst-ctn-kix_list_31-2;
  }

  ol.lst-kix_list_38-8.start {
    counter-reset: lst-ctn-kix_list_38-8 0;
  }

  ol.lst-kix_list_19-5.start {
    counter-reset: lst-ctn-kix_list_19-5 0;
  }

  .lst-kix_list_35-8 > li {
    counter-increment: lst-ctn-kix_list_35-8;
  }

  ol.lst-kix_list_39-6.start {
    counter-reset: lst-ctn-kix_list_39-6 0;
  }

  ol.lst-kix_list_45-6.start {
    counter-reset: lst-ctn-kix_list_45-6 0;
  }

  ol.lst-kix_list_26-3.start {
    counter-reset: lst-ctn-kix_list_26-3 0;
  }

  .lst-kix_list_45-5 > li:before {
    content: "" counter(lst-ctn-kix_list_45-5, decimal) " ";
  }

  .lst-kix_list_25-0 > li:before {
    content: "" counter(lst-ctn-kix_list_25-0, decimal) ") ";
  }

  ol.lst-kix_list_10-8.start {
    counter-reset: lst-ctn-kix_list_10-8 0;
  }

  .lst-kix_list_39-0 > li:before {
    content: "" counter(lst-ctn-kix_list_39-0, decimal) " ";
  }

  ol.lst-kix_list_40-7.start {
    counter-reset: lst-ctn-kix_list_40-7 0;
  }

  .lst-kix_list_37-0 > li {
    counter-increment: lst-ctn-kix_list_37-0;
  }

  ol.lst-kix_list_21-4.start {
    counter-reset: lst-ctn-kix_list_21-4 0;
  }

  .lst-kix_list_44-1 > li:before {
    content: "" counter(lst-ctn-kix_list_44-1, decimal) " ";
  }

  .lst-kix_list_45-1 > li:before {
    content: "" counter(lst-ctn-kix_list_45-1, decimal) " ";
  }

  ol.lst-kix_list_20-6.start {
    counter-reset: lst-ctn-kix_list_20-6 0;
  }

  .lst-kix_list_13-8 > li {
    counter-increment: lst-ctn-kix_list_13-8;
  }

  .lst-kix_list_2-2 > li {
    counter-increment: lst-ctn-kix_list_2-2;
  }

  .lst-kix_list_44-5 > li:before {
    content: "" counter(lst-ctn-kix_list_44-5, decimal) " ";
  }

  ol.lst-kix_list_4-7.start {
    counter-reset: lst-ctn-kix_list_4-7 0;
  }

  .lst-kix_list_26-2 > li {
    counter-increment: lst-ctn-kix_list_26-2;
  }

  .lst-kix_list_40-8 > li {
    counter-increment: lst-ctn-kix_list_40-8;
  }

  .lst-kix_list_27-4 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) "."
      counter(lst-ctn-kix_list_27-3, decimal) "."
      counter(lst-ctn-kix_list_27-4, decimal) ". ";
  }

  .lst-kix_list_20-2 > li {
    counter-increment: lst-ctn-kix_list_20-2;
  }

  .lst-kix_list_6-6 > li {
    counter-increment: lst-ctn-kix_list_6-6;
  }

  ol.lst-kix_list_15-7.start {
    counter-reset: lst-ctn-kix_list_15-7 0;
  }

  .lst-kix_list_13-6 > li {
    counter-increment: lst-ctn-kix_list_13-6;
  }

  ol.lst-kix_list_14-6.start {
    counter-reset: lst-ctn-kix_list_14-6 0;
  }

  .lst-kix_list_39-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_39-4, lower-roman) ") ";
  }

  .lst-kix_list_39-8 > li:before {
    content: "" counter(lst-ctn-kix_list_39-8, decimal) " ";
  }

  ol.lst-kix_list_26-8.start {
    counter-reset: lst-ctn-kix_list_26-8 0;
  }

  .lst-kix_list_19-6 > li {
    counter-increment: lst-ctn-kix_list_19-6;
  }

  .lst-kix_list_27-0 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) ". ";
  }

  .lst-kix_list_20-4 > li {
    counter-increment: lst-ctn-kix_list_20-4;
  }

  ol.lst-kix_list_44-5.start {
    counter-reset: lst-ctn-kix_list_44-5 0;
  }

  ol.lst-kix_list_33-6.start {
    counter-reset: lst-ctn-kix_list_33-6 0;
  }

  ol.lst-kix_list_20-1.start {
    counter-reset: lst-ctn-kix_list_20-1 0;
  }

  .lst-kix_list_25-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_25-4, lower-latin) ") ";
  }

  .lst-kix_list_19-4 > li {
    counter-increment: lst-ctn-kix_list_19-4;
  }

  .lst-kix_list_35-1 > li {
    counter-increment: lst-ctn-kix_list_35-1;
  }

  .lst-kix_list_46-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_24-1 > li {
    counter-increment: lst-ctn-kix_list_24-1;
  }

  .lst-kix_list_25-8 > li:before {
    content: "" counter(lst-ctn-kix_list_25-8, lower-roman) ". ";
  }

  ol.lst-kix_list_50-0.start {
    counter-reset: lst-ctn-kix_list_50-0 0;
  }

  ol.lst-kix_list_51-8.start {
    counter-reset: lst-ctn-kix_list_51-8 0;
  }

  .lst-kix_list_51-8 > li {
    counter-increment: lst-ctn-kix_list_51-8;
  }

  ol.lst-kix_list_26-6.start {
    counter-reset: lst-ctn-kix_list_26-6 0;
  }

  .lst-kix_list_44-5 > li {
    counter-increment: lst-ctn-kix_list_44-5;
  }

  .lst-kix_list_37-2 > li {
    counter-increment: lst-ctn-kix_list_37-2;
  }

  .lst-kix_list_15-2 > li {
    counter-increment: lst-ctn-kix_list_15-2;
  }

  .lst-kix_list_22-5 > li {
    counter-increment: lst-ctn-kix_list_22-5;
  }

  .lst-kix_list_28-8 > li:before {
    content: "" counter(lst-ctn-kix_list_28-8, lower-roman) ". ";
  }

  .lst-kix_list_24-6 > li {
    counter-increment: lst-ctn-kix_list_24-6;
  }

  .lst-kix_list_17-3 > li {
    counter-increment: lst-ctn-kix_list_17-3;
  }

  .lst-kix_list_28-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_28-4, lower-latin) ") ";
  }

  .lst-kix_list_39-3 > li {
    counter-increment: lst-ctn-kix_list_39-3;
  }

  .lst-kix_list_28-3 > li {
    counter-increment: lst-ctn-kix_list_28-3;
  }

  ol.lst-kix_list_25-0.start {
    counter-reset: lst-ctn-kix_list_25-0 0;
  }

  .lst-kix_list_35-6 > li {
    counter-increment: lst-ctn-kix_list_35-6;
  }

  ol.lst-kix_list_13-0.start {
    counter-reset: lst-ctn-kix_list_13-0 0;
  }

  ol.lst-kix_list_32-3.start {
    counter-reset: lst-ctn-kix_list_32-3 0;
  }

  ol.lst-kix_list_50-2.start {
    counter-reset: lst-ctn-kix_list_50-2 0;
  }

  .lst-kix_list_11-5 > li {
    counter-increment: lst-ctn-kix_list_11-5;
  }

  .lst-kix_list_28-0 > li:before {
    content: "" counter(lst-ctn-kix_list_28-0, decimal) ") ";
  }

  .lst-kix_list_27-8 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) "."
      counter(lst-ctn-kix_list_27-3, decimal) "."
      counter(lst-ctn-kix_list_27-4, decimal) "."
      counter(lst-ctn-kix_list_27-5, decimal) "."
      counter(lst-ctn-kix_list_27-6, decimal) "."
      counter(lst-ctn-kix_list_27-7, decimal) "."
      counter(lst-ctn-kix_list_27-8, decimal) ". ";
  }

  ol.lst-kix_list_20-3.start {
    counter-reset: lst-ctn-kix_list_20-3 0;
  }

  .lst-kix_list_4-1 > li {
    counter-increment: lst-ctn-kix_list_4-1;
  }

  .lst-kix_list_19-1 > li:before {
    content: "" counter(lst-ctn-kix_list_19-1, decimal) ". ";
  }

  .lst-kix_list_27-8 > li {
    counter-increment: lst-ctn-kix_list_27-8;
  }

  .lst-kix_list_19-3 > li:before {
    content: "" counter(lst-ctn-kix_list_19-3, decimal) ". ";
  }

  ol.lst-kix_list_38-0.start {
    counter-reset: lst-ctn-kix_list_38-0 0;
  }

  .lst-kix_list_15-0 > li {
    counter-increment: lst-ctn-kix_list_15-0;
  }

  ol.lst-kix_list_6-6.start {
    counter-reset: lst-ctn-kix_list_6-6 0;
  }

  .lst-kix_list_47-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_39-5 > li {
    counter-increment: lst-ctn-kix_list_39-5;
  }

  .lst-kix_list_32-6 > li {
    counter-increment: lst-ctn-kix_list_32-6;
  }

  ol.lst-kix_list_29-6.start {
    counter-reset: lst-ctn-kix_list_29-6 0;
  }

  .lst-kix_list_11-0 > li {
    counter-increment: lst-ctn-kix_list_11-0;
  }

  ol.lst-kix_list_1-5.start {
    counter-reset: lst-ctn-kix_list_1-5 0;
  }

  ol.lst-kix_list_9-6.start {
    counter-reset: lst-ctn-kix_list_9-6 0;
  }

  ol.lst-kix_list_16-3.start {
    counter-reset: lst-ctn-kix_list_16-3 0;
  }

  .lst-kix_list_22-7 > li {
    counter-increment: lst-ctn-kix_list_22-7;
  }

  .lst-kix_list_47-3 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_4-5.start {
    counter-reset: lst-ctn-kix_list_4-5 0;
  }

  .lst-kix_list_37-0 > li:before {
    content: "" counter(lst-ctn-kix_list_37-0, decimal) " ";
  }

  .lst-kix_list_26-7 > li {
    counter-increment: lst-ctn-kix_list_26-7;
  }

  .lst-kix_list_41-0 > li {
    counter-increment: lst-ctn-kix_list_41-0;
  }

  .lst-kix_list_33-7 > li {
    counter-increment: lst-ctn-kix_list_33-7;
  }

  .lst-kix_list_5-2 > li {
    counter-increment: lst-ctn-kix_list_5-2;
  }

  .lst-kix_list_37-7 > li {
    counter-increment: lst-ctn-kix_list_37-7;
  }

  .lst-kix_list_37-8 > li:before {
    content: "" counter(lst-ctn-kix_list_37-8, upper-latin) ". ";
  }

  .lst-kix_list_37-6 > li:before {
    content: "" counter(lst-ctn-kix_list_37-6, upper-roman) ". ";
  }

  .lst-kix_list_28-5 > li {
    counter-increment: lst-ctn-kix_list_28-5;
  }

  ol.lst-kix_list_35-0.start {
    counter-reset: lst-ctn-kix_list_35-0 0;
  }

  .lst-kix_list_22-0 > li {
    counter-increment: lst-ctn-kix_list_22-0;
  }

  .lst-kix_list_46-1 > li:before {
    content: "o  ";
  }

  ol.lst-kix_list_33-4.start {
    counter-reset: lst-ctn-kix_list_33-4 0;
  }

  .lst-kix_list_21-6 > li {
    counter-increment: lst-ctn-kix_list_21-6;
  }

  .lst-kix_list_10-3 > li {
    counter-increment: lst-ctn-kix_list_10-3;
  }

  ol.lst-kix_list_1-0.start {
    counter-reset: lst-ctn-kix_list_1-0 0;
  }

  ol.lst-kix_list_13-3.start {
    counter-reset: lst-ctn-kix_list_13-3 0;
  }

  .lst-kix_list_51-6 > li {
    counter-increment: lst-ctn-kix_list_51-6;
  }

  .lst-kix_list_26-0 > li {
    counter-increment: lst-ctn-kix_list_26-0;
  }

  ol.lst-kix_list_4-0.start {
    counter-reset: lst-ctn-kix_list_4-0 0;
  }

  .lst-kix_list_27-4 > li {
    counter-increment: lst-ctn-kix_list_27-4;
  }

  .lst-kix_list_38-2 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) " ";
  }

  .lst-kix_list_38-4 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) "."
      counter(lst-ctn-kix_list_38-3, decimal) "."
      counter(lst-ctn-kix_list_38-4, decimal) " ";
  }

  ol.lst-kix_list_43-2.start {
    counter-reset: lst-ctn-kix_list_43-2 0;
  }

  .lst-kix_list_16-1 > li {
    counter-increment: lst-ctn-kix_list_16-1;
  }

  .lst-kix_list_27-1 > li {
    counter-increment: lst-ctn-kix_list_27-1;
  }

  .lst-kix_list_38-4 > li {
    counter-increment: lst-ctn-kix_list_38-4;
  }

  ol.lst-kix_list_43-4.start {
    counter-reset: lst-ctn-kix_list_43-4 0;
  }

  .lst-kix_list_17-5 > li {
    counter-increment: lst-ctn-kix_list_17-5;
  }

  ol.lst-kix_list_9-4.start {
    counter-reset: lst-ctn-kix_list_9-4 0;
  }

  ol.lst-kix_list_41-8.start {
    counter-reset: lst-ctn-kix_list_41-8 0;
  }

  .lst-kix_list_21-3 > li {
    counter-increment: lst-ctn-kix_list_21-3;
  }

  ol.lst-kix_list_30-1.start {
    counter-reset: lst-ctn-kix_list_30-1 0;
  }

  .lst-kix_list_36-4 > li:before {
    content: "" counter(lst-ctn-kix_list_36-4, lower-latin) ". ";
  }

  ol.lst-kix_list_13-5.start {
    counter-reset: lst-ctn-kix_list_13-5 0;
  }

  ol.lst-kix_list_13-8.start {
    counter-reset: lst-ctn-kix_list_13-8 0;
  }

  .lst-kix_list_11-7 > li {
    counter-increment: lst-ctn-kix_list_11-7;
  }

  .lst-kix_list_5-5 > li {
    counter-increment: lst-ctn-kix_list_5-5;
  }

  .lst-kix_list_36-2 > li:before {
    content: "" counter(lst-ctn-kix_list_36-2, lower-roman) ". ";
  }

  ol.lst-kix_list_43-7.start {
    counter-reset: lst-ctn-kix_list_43-7 0;
  }

  .lst-kix_list_44-7 > li {
    counter-increment: lst-ctn-kix_list_44-7;
  }

  .lst-kix_list_16-8 > li {
    counter-increment: lst-ctn-kix_list_16-8;
  }

  .lst-kix_list_38-8 > li {
    counter-increment: lst-ctn-kix_list_38-8;
  }

  ol.lst-kix_list_18-4.start {
    counter-reset: lst-ctn-kix_list_18-4 0;
  }

  ol.lst-kix_list_29-1.start {
    counter-reset: lst-ctn-kix_list_29-1 0;
  }

  ol.lst-kix_list_1-2.start {
    counter-reset: lst-ctn-kix_list_1-2 0;
  }

  .lst-kix_list_20-4 > li:before {
    content: "" counter(lst-ctn-kix_list_20-4, decimal) ". ";
  }

  ol.lst-kix_list_22-7.start {
    counter-reset: lst-ctn-kix_list_22-7 0;
  }

  ol.lst-kix_list_6-1.start {
    counter-reset: lst-ctn-kix_list_6-1 0;
  }

  .lst-kix_list_20-2 > li:before {
    content: "" counter(lst-ctn-kix_list_20-2, decimal) ". ";
  }

  .lst-kix_list_54-1 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) " ";
  }

  ol.lst-kix_list_16-8.start {
    counter-reset: lst-ctn-kix_list_16-8 0;
  }

  ol.lst-kix_list_33-1.start {
    counter-reset: lst-ctn-kix_list_33-1 0;
  }

  .lst-kix_list_4-8 > li {
    counter-increment: lst-ctn-kix_list_4-8;
  }

  .lst-kix_list_55-5 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_9-1.start {
    counter-reset: lst-ctn-kix_list_9-1 0;
  }

  ol.lst-kix_list_24-2.start {
    counter-reset: lst-ctn-kix_list_24-2 0;
  }

  .lst-kix_list_54-7 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) "."
      counter(lst-ctn-kix_list_54-3, decimal) "."
      counter(lst-ctn-kix_list_54-4, decimal) "."
      counter(lst-ctn-kix_list_54-5, decimal) "."
      counter(lst-ctn-kix_list_54-6, decimal) "."
      counter(lst-ctn-kix_list_54-7, decimal) " ";
  }

  .lst-kix_list_19-8 > li {
    counter-increment: lst-ctn-kix_list_19-8;
  }

  .lst-kix_list_55-3 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_35-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_35-6, decimal) ") ";
  }

  .lst-kix_list_36-5 > li {
    counter-increment: lst-ctn-kix_list_36-5;
  }

  .lst-kix_list_3-2 > li:before {
    content: "" counter(lst-ctn-kix_list_3-2, decimal) ". ";
  }

  .lst-kix_list_8-1 > li:before {
    content: "o  ";
  }

  .lst-kix_list_30-7 > li:before {
    content: "" counter(lst-ctn-kix_list_30-7, decimal) " ";
  }

  .lst-kix_list_43-3 > li {
    counter-increment: lst-ctn-kix_list_43-3;
  }

  .lst-kix_list_26-2 > li:before {
    content: "(" counter(lst-ctn-kix_list_26-2, lower-latin) ") ";
  }

  .lst-kix_list_21-6 > li:before {
    content: "" counter(lst-ctn-kix_list_21-6, decimal) ". ";
  }

  .lst-kix_list_41-3 > li {
    counter-increment: lst-ctn-kix_list_41-3;
  }

  ol.lst-kix_list_54-6.start {
    counter-reset: lst-ctn-kix_list_54-6 0;
  }

  ol.lst-kix_list_35-3.start {
    counter-reset: lst-ctn-kix_list_35-3 0;
  }

  ol.lst-kix_list_4-2.start {
    counter-reset: lst-ctn-kix_list_4-2 0;
  }

  ol.lst-kix_list_27-2.start {
    counter-reset: lst-ctn-kix_list_27-2 0;
  }

  ol.lst-kix_list_18-7.start {
    counter-reset: lst-ctn-kix_list_18-7 0;
  }

  .lst-kix_list_25-5 > li {
    counter-increment: lst-ctn-kix_list_25-5;
  }

  ol.lst-kix_list_11-6.start {
    counter-reset: lst-ctn-kix_list_11-6 0;
  }

  ol.lst-kix_list_6-4.start {
    counter-reset: lst-ctn-kix_list_6-4 0;
  }

  .lst-kix_list_45-3 > li:before {
    content: "" counter(lst-ctn-kix_list_45-3, decimal) " ";
  }

  .lst-kix_list_17-1 > li:before {
    content: "" counter(lst-ctn-kix_list_17-1, decimal) ". ";
  }

  .lst-kix_list_32-3 > li {
    counter-increment: lst-ctn-kix_list_32-3;
  }

  ol.lst-kix_list_27-1.start {
    counter-reset: lst-ctn-kix_list_27-1 0;
  }

  .lst-kix_list_16-5 > li:before {
    content: "" counter(lst-ctn-kix_list_16-5, decimal) ". ";
  }

  ol.lst-kix_list_22-2.start {
    counter-reset: lst-ctn-kix_list_22-2 0;
  }

  .lst-kix_list_30-7 > li {
    counter-increment: lst-ctn-kix_list_30-7;
  }

  ol.lst-kix_list_41-5.start {
    counter-reset: lst-ctn-kix_list_41-5 0;
  }

  .lst-kix_list_44-7 > li:before {
    content: "" counter(lst-ctn-kix_list_44-7, decimal) " ";
  }

  ol.lst-kix_list_29-3.start {
    counter-reset: lst-ctn-kix_list_29-3 0;
  }

  .lst-kix_list_16-4 > li {
    counter-increment: lst-ctn-kix_list_16-4;
  }

  ol.lst-kix_list_35-8.start {
    counter-reset: lst-ctn-kix_list_35-8 0;
  }

  .lst-kix_list_38-1 > li {
    counter-increment: lst-ctn-kix_list_38-1;
  }

  ol.lst-kix_list_16-5.start {
    counter-reset: lst-ctn-kix_list_16-5 0;
  }

  .lst-kix_list_41-7 > li {
    counter-increment: lst-ctn-kix_list_41-7;
  }

  ol.lst-kix_list_35-5.start {
    counter-reset: lst-ctn-kix_list_35-5 0;
  }

  .lst-kix_list_2-6 > li:before {
    content: "" counter(lst-ctn-kix_list_2-6, decimal) ". ";
  }

  ol.lst-kix_list_54-8.start {
    counter-reset: lst-ctn-kix_list_54-8 0;
  }

  .lst-kix_list_14-5 > li {
    counter-increment: lst-ctn-kix_list_14-5;
  }

  .lst-kix_list_7-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_27-6 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) "."
      counter(lst-ctn-kix_list_27-3, decimal) "."
      counter(lst-ctn-kix_list_27-4, decimal) "."
      counter(lst-ctn-kix_list_27-5, decimal) "."
      counter(lst-ctn-kix_list_27-6, decimal) ". ";
  }

  .lst-kix_list_23-2 > li {
    counter-increment: lst-ctn-kix_list_23-2;
  }

  ol.lst-kix_list_30-4.start {
    counter-reset: lst-ctn-kix_list_30-4 0;
  }

  .lst-kix_list_30-0 > li {
    counter-increment: lst-ctn-kix_list_30-0;
  }

  ol.lst-kix_list_11-1.start {
    counter-reset: lst-ctn-kix_list_11-1 0;
  }

  .lst-kix_list_18-5 > li:before {
    content: "" counter(lst-ctn-kix_list_18-5, decimal) ". ";
  }

  .lst-kix_list_13-6 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) "."
      counter(lst-ctn-kix_list_13-3, decimal) "."
      counter(lst-ctn-kix_list_13-4, decimal) "."
      counter(lst-ctn-kix_list_13-5, decimal) "."
      counter(lst-ctn-kix_list_13-6, decimal) ". ";
  }

  .lst-kix_list_10-6 > li {
    counter-increment: lst-ctn-kix_list_10-6;
  }

  .lst-kix_list_1-7 > li {
    counter-increment: lst-ctn-kix_list_1-7;
  }

  ol.lst-kix_list_41-0.start {
    counter-reset: lst-ctn-kix_list_41-0 0;
  }

  ol.lst-kix_list_54-5.start {
    counter-reset: lst-ctn-kix_list_54-5 0;
  }

  .lst-kix_list_39-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_39-6, decimal) ") ";
  }

  .lst-kix_list_54-6 > li {
    counter-increment: lst-ctn-kix_list_54-6;
  }

  .lst-kix_list_29-0 > li {
    counter-increment: lst-ctn-kix_list_29-0;
  }

  ol.lst-kix_list_24-7.start {
    counter-reset: lst-ctn-kix_list_24-7 0;
  }

  .lst-kix_list_31-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_31-3, lower-latin) ") ";
  }

  .lst-kix_list_43-6 > li {
    counter-increment: lst-ctn-kix_list_43-6;
  }

  .lst-kix_list_10-2 > li:before {
    content: "(" counter(lst-ctn-kix_list_10-2, lower-roman) ") ";
  }

  ol.lst-kix_list_1-7.start {
    counter-reset: lst-ctn-kix_list_1-7 0;
  }

  .lst-kix_list_4-6 > li:before {
    content: "" counter(lst-ctn-kix_list_4-6, decimal) ". ";
  }

  ol.lst-kix_list_29-8.start {
    counter-reset: lst-ctn-kix_list_29-8 0;
  }

  .lst-kix_list_25-6 > li:before {
    content: "" counter(lst-ctn-kix_list_25-6, decimal) ". ";
  }

  .lst-kix_list_46-7 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_34-2 > li {
    counter-increment: lst-ctn-kix_list_34-2;
  }

  .lst-kix_list_12-2 > li {
    counter-increment: lst-ctn-kix_list_12-2;
  }

  .lst-kix_list_9-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_9-5, upper-latin) ") ";
  }

  .lst-kix_list_29-6 > li:before {
    content: "" counter(lst-ctn-kix_list_29-6, lower-roman) ") ";
  }

  ol.lst-kix_list_22-0.start {
    counter-reset: lst-ctn-kix_list_22-0 0;
  }

  ol.lst-kix_list_41-3.start {
    counter-reset: lst-ctn-kix_list_41-3 0;
  }

  .lst-kix_list_12-2 > li:before {
    content: "(" counter(lst-ctn-kix_list_12-2, lower-roman) ") ";
  }

  .lst-kix_list_11-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_11-6, decimal) ") ";
  }

  .lst-kix_list_32-7 > li:before {
    content: "" counter(lst-ctn-kix_list_32-7, decimal) " ";
  }

  ol.lst-kix_list_30-6.start {
    counter-reset: lst-ctn-kix_list_30-6 0;
  }

  .lst-kix_list_1-2 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) " ";
  }

  ol.lst-kix_list_11-3.start {
    counter-reset: lst-ctn-kix_list_11-3 0;
  }

  .lst-kix_list_49-3 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_1-0 > li {
    counter-increment: lst-ctn-kix_list_1-0;
  }

  ol.lst-kix_list_41-2.start {
    counter-reset: lst-ctn-kix_list_41-2 0;
  }

  .lst-kix_list_18-7 > li {
    counter-increment: lst-ctn-kix_list_18-7;
  }

  .lst-kix_list_48-7 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_50-4 > li {
    counter-increment: lst-ctn-kix_list_50-4;
  }

  .lst-kix_list_29-7 > li {
    counter-increment: lst-ctn-kix_list_29-7;
  }

  ol.lst-kix_list_11-4.start {
    counter-reset: lst-ctn-kix_list_11-4 0;
  }

  .lst-kix_list_28-2 > li:before {
    content: "" counter(lst-ctn-kix_list_28-2, lower-roman) ") ";
  }

  .lst-kix_list_51-4 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) "."
      counter(lst-ctn-kix_list_51-3, decimal) "."
      counter(lst-ctn-kix_list_51-4, decimal) " ";
  }

  .lst-kix_list_14-1 > li:before {
    content: "" counter(lst-ctn-kix_list_14-1, lower-latin) ") ";
  }

  .lst-kix_list_25-8 > li {
    counter-increment: lst-ctn-kix_list_25-8;
  }

  .lst-kix_list_14-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_14-4, lower-latin) ") ";
  }

  .lst-kix_list_51-1 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) " ";
  }

  ol.lst-kix_list_45-4.start {
    counter-reset: lst-ctn-kix_list_45-4 0;
  }

  .lst-kix_list_37-5 > li {
    counter-increment: lst-ctn-kix_list_37-5;
  }

  .lst-kix_list_14-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_14-5, lower-roman) ") ";
  }

  .lst-kix_list_14-7 > li:before {
    content: "" counter(lst-ctn-kix_list_14-7, lower-latin) ". ";
  }

  .lst-kix_ek190fjkpfvz-7 > li:before {
    content: "\0025cb  ";
  }

  .lst-kix_ek190fjkpfvz-6 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_34-6 > li {
    counter-increment: lst-ctn-kix_list_34-6;
  }

  .lst-kix_list_51-8 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) "."
      counter(lst-ctn-kix_list_51-3, decimal) "."
      counter(lst-ctn-kix_list_51-4, decimal) "."
      counter(lst-ctn-kix_list_51-5, decimal) "."
      counter(lst-ctn-kix_list_51-6, decimal) "."
      counter(lst-ctn-kix_list_51-7, decimal) "."
      counter(lst-ctn-kix_list_51-8, decimal) " ";
  }

  .lst-kix_list_51-7 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) "."
      counter(lst-ctn-kix_list_51-3, decimal) "."
      counter(lst-ctn-kix_list_51-4, decimal) "."
      counter(lst-ctn-kix_list_51-5, decimal) "."
      counter(lst-ctn-kix_list_51-6, decimal) "."
      counter(lst-ctn-kix_list_51-7, decimal) " ";
  }

  .lst-kix_ek190fjkpfvz-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_32-1.start {
    counter-reset: lst-ctn-kix_list_32-1 0;
  }

  .lst-kix_list_20-7 > li {
    counter-increment: lst-ctn-kix_list_20-7;
  }

  ol.lst-kix_list_35-6.start {
    counter-reset: lst-ctn-kix_list_35-6 0;
  }

  .lst-kix_ek190fjkpfvz-4 > li:before {
    content: "\0025cb  ";
  }

  ol.lst-kix_list_28-3.start {
    counter-reset: lst-ctn-kix_list_28-3 0;
  }

  .lst-kix_ek190fjkpfvz-3 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_44-1 > li {
    counter-increment: lst-ctn-kix_list_44-1;
  }

  .lst-kix_list_32-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_32-3, lower-latin) ") ";
  }

  .lst-kix_list_50-0 > li {
    counter-increment: lst-ctn-kix_list_50-0;
  }

  .lst-kix_list_14-8 > li:before {
    content: "" counter(lst-ctn-kix_list_14-8, lower-roman) ". ";
  }

  .lst-kix_list_39-7 > li {
    counter-increment: lst-ctn-kix_list_39-7;
  }

  .lst-kix_list_32-0 > li:before {
    content: "" counter(lst-ctn-kix_list_32-0, decimal) " ";
  }

  .lst-kix_list_3-2 > li {
    counter-increment: lst-ctn-kix_list_3-2;
  }

  ol.lst-kix_list_15-5.start {
    counter-reset: lst-ctn-kix_list_15-5 0;
  }

  ol.lst-kix_list_22-3.start {
    counter-reset: lst-ctn-kix_list_22-3 0;
  }

  .lst-kix_list_5-4 > li {
    counter-increment: lst-ctn-kix_list_5-4;
  }

  ol.lst-kix_list_24-6.start {
    counter-reset: lst-ctn-kix_list_24-6 0;
  }

  .lst-kix_list_5-1 > li:before {
    content: "" counter(lst-ctn-kix_list_5-1, decimal) ". ";
  }

  .lst-kix_list_5-7 > li:before {
    content: "" counter(lst-ctn-kix_list_5-7, decimal) ". ";
  }

  .lst-kix_list_5-8 > li:before {
    content: "" counter(lst-ctn-kix_list_5-8, decimal) ". ";
  }

  .lst-kix_list_5-4 > li:before {
    content: "" counter(lst-ctn-kix_list_5-4, decimal) ". ";
  }

  .lst-kix_list_5-5 > li:before {
    content: "" counter(lst-ctn-kix_list_5-5, decimal) ". ";
  }

  .lst-kix_list_50-1 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) " ";
  }

  .lst-kix_list_50-2 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) " ";
  }

  .lst-kix_list_6-1 > li:before {
    content: "" counter(lst-ctn-kix_list_6-1, decimal) ". ";
  }

  .lst-kix_list_6-3 > li:before {
    content: "" counter(lst-ctn-kix_list_6-3, decimal) ". ";
  }

  ol.lst-kix_list_32-6.start {
    counter-reset: lst-ctn-kix_list_32-6 0;
  }

  .lst-kix_list_6-8 > li {
    counter-increment: lst-ctn-kix_list_6-8;
  }

  .lst-kix_list_6-0 > li:before {
    content: "" counter(lst-ctn-kix_list_6-0, decimal) ". ";
  }

  .lst-kix_list_6-4 > li:before {
    content: "" counter(lst-ctn-kix_list_6-4, decimal) ". ";
  }

  ol.lst-kix_list_14-8.start {
    counter-reset: lst-ctn-kix_list_14-8 0;
  }

  ol.lst-kix_list_15-0.start {
    counter-reset: lst-ctn-kix_list_15-0 0;
  }

  ol.lst-kix_list_44-7.start {
    counter-reset: lst-ctn-kix_list_44-7 0;
  }

  .lst-kix_list_2-5 > li {
    counter-increment: lst-ctn-kix_list_2-5;
  }

  .lst-kix_list_50-8 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) "."
      counter(lst-ctn-kix_list_50-3, decimal) "."
      counter(lst-ctn-kix_list_50-4, decimal) "."
      counter(lst-ctn-kix_list_50-5, decimal) "."
      counter(lst-ctn-kix_list_50-6, decimal) "."
      counter(lst-ctn-kix_list_50-7, decimal) "."
      counter(lst-ctn-kix_list_50-8, decimal) " ";
  }

  .lst-kix_list_50-5 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) "."
      counter(lst-ctn-kix_list_50-3, decimal) "."
      counter(lst-ctn-kix_list_50-4, decimal) "."
      counter(lst-ctn-kix_list_50-5, decimal) " ";
  }

  .lst-kix_list_6-7 > li:before {
    content: "" counter(lst-ctn-kix_list_6-7, decimal) ". ";
  }

  .lst-kix_list_6-6 > li:before {
    content: "" counter(lst-ctn-kix_list_6-6, decimal) ". ";
  }

  ol.lst-kix_list_10-6.start {
    counter-reset: lst-ctn-kix_list_10-6 0;
  }

  ol.lst-kix_list_27-6.start {
    counter-reset: lst-ctn-kix_list_27-6 0;
  }

  .lst-kix_list_7-6 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_18-5 > li {
    counter-increment: lst-ctn-kix_list_18-5;
  }

  ol.lst-kix_list_19-7.start {
    counter-reset: lst-ctn-kix_list_19-7 0;
  }

  ol.lst-kix_list_6-2.start {
    counter-reset: lst-ctn-kix_list_6-2 0;
  }

  .lst-kix_list_15-5 > li {
    counter-increment: lst-ctn-kix_list_15-5;
  }

  ol.lst-kix_list_24-1.start {
    counter-reset: lst-ctn-kix_list_24-1 0;
  }

  ol.lst-kix_list_36-3.start {
    counter-reset: lst-ctn-kix_list_36-3 0;
  }

  .lst-kix_list_7-2 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_27-2 > li {
    counter-increment: lst-ctn-kix_list_27-2;
  }

  .lst-kix_list_31-0 > li {
    counter-increment: lst-ctn-kix_list_31-0;
  }

  ol.lst-kix_list_22-8.start {
    counter-reset: lst-ctn-kix_list_22-8 0;
  }

  .lst-kix_list_34-8 > li:before {
    content: "" counter(lst-ctn-kix_list_34-8, decimal) " ";
  }

  .lst-kix_list_31-0 > li:before {
    content: "" counter(lst-ctn-kix_list_31-0, decimal) " ";
  }

  .lst-kix_list_12-6 > li {
    counter-increment: lst-ctn-kix_list_12-6;
  }

  ol.lst-kix_list_23-0.start {
    counter-reset: lst-ctn-kix_list_23-0 0;
  }

  .lst-kix_list_52-7 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_9-8 > li {
    counter-increment: lst-ctn-kix_list_9-8;
  }

  .lst-kix_list_13-4 > li {
    counter-increment: lst-ctn-kix_list_13-4;
  }

  .lst-kix_list_52-3 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_31-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_31-4, lower-roman) ") ";
  }

  .lst-kix_list_15-7 > li:before {
    content: "" counter(lst-ctn-kix_list_15-7, lower-latin) ". ";
  }

  .lst-kix_list_17-7 > li {
    counter-increment: lst-ctn-kix_list_17-7;
  }

  .lst-kix_list_4-5 > li:before {
    content: "" counter(lst-ctn-kix_list_4-5, decimal) ". ";
  }

  .lst-kix_list_51-4 > li {
    counter-increment: lst-ctn-kix_list_51-4;
  }

  .lst-kix_list_15-1 > li:before {
    content: "Section " counter(lst-ctn-kix_list_15-0, upper-roman) "."
      counter(lst-ctn-kix_list_15-1, decimal) " ";
  }

  ol.lst-kix_list_1-4.start {
    counter-reset: lst-ctn-kix_list_1-4 0;
  }

  .lst-kix_list_15-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_15-3, lower-roman) ") ";
  }

  .lst-kix_list_22-1 > li {
    counter-increment: lst-ctn-kix_list_22-1;
  }

  .lst-kix_list_42-7 > li {
    counter-increment: lst-ctn-kix_list_42-7;
  }

  ol.lst-kix_list_4-4.start {
    counter-reset: lst-ctn-kix_list_4-4 0;
  }

  .lst-kix_list_33-2 > li {
    counter-increment: lst-ctn-kix_list_33-2;
  }

  .lst-kix_list_40-5 > li {
    counter-increment: lst-ctn-kix_list_40-5;
  }

  .lst-kix_list_30-2 > li {
    counter-increment: lst-ctn-kix_list_30-2;
  }

  ol.lst-kix_list_9-2.start {
    counter-reset: lst-ctn-kix_list_9-2 0;
  }

  .lst-kix_list_53-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_11-2 > li {
    counter-increment: lst-ctn-kix_list_11-2;
  }

  .lst-kix_list_33-8 > li {
    counter-increment: lst-ctn-kix_list_33-8;
  }

  .lst-kix_list_12-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_12-3, upper-latin) ") ";
  }

  ol.lst-kix_list_31-4.start {
    counter-reset: lst-ctn-kix_list_31-4 0;
  }

  .lst-kix_list_32-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_32-6, decimal) ") ";
  }

  ol.lst-kix_list_54-0.start {
    counter-reset: lst-ctn-kix_list_54-0 6;
  }

  .lst-kix_list_33-2 > li:before {
    content: "" counter(lst-ctn-kix_list_33-0, decimal) "."
      counter(lst-ctn-kix_list_33-1, decimal) "."
      counter(lst-ctn-kix_list_33-2, decimal) " ";
  }

  .lst-kix_list_16-3 > li {
    counter-increment: lst-ctn-kix_list_16-3;
  }

  .lst-kix_list_13-3 > li {
    counter-increment: lst-ctn-kix_list_13-3;
  }

  ol.lst-kix_list_40-5.start {
    counter-reset: lst-ctn-kix_list_40-5 0;
  }

  .lst-kix_list_10-4 > li {
    counter-increment: lst-ctn-kix_list_10-4;
  }

  .lst-kix_list_14-1 > li {
    counter-increment: lst-ctn-kix_list_14-1;
  }

  .lst-kix_list_34-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_34-4, lower-roman) ") ";
  }

  .lst-kix_list_35-3 > li {
    counter-increment: lst-ctn-kix_list_35-3;
  }

  .lst-kix_list_28-0 > li {
    counter-increment: lst-ctn-kix_list_28-0;
  }

  .lst-kix_list_13-5 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) "."
      counter(lst-ctn-kix_list_13-3, decimal) "."
      counter(lst-ctn-kix_list_13-4, decimal) "."
      counter(lst-ctn-kix_list_13-5, decimal) ". ";
  }

  ol.lst-kix_list_28-8.start {
    counter-reset: lst-ctn-kix_list_28-8 0;
  }

  .lst-kix_list_36-1 > li {
    counter-increment: lst-ctn-kix_list_36-1;
  }

  .lst-kix_list_36-7 > li {
    counter-increment: lst-ctn-kix_list_36-7;
  }

  .lst-kix_list_29-4 > li {
    counter-increment: lst-ctn-kix_list_29-4;
  }

  .lst-kix_list_54-4 > li {
    counter-increment: lst-ctn-kix_list_54-4;
  }

  .lst-kix_list_33-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_33-6, decimal) ") ";
  }

  .lst-kix_list_32-4 > li {
    counter-increment: lst-ctn-kix_list_32-4;
  }

  .lst-kix_list_33-8 > li:before {
    content: "" counter(lst-ctn-kix_list_33-8, decimal) " ";
  }

  ol.lst-kix_list_29-0.start {
    counter-reset: lst-ctn-kix_list_29-0 0;
  }

  .lst-kix_list_34-2 > li:before {
    content: "" counter(lst-ctn-kix_list_34-0, decimal) "."
      counter(lst-ctn-kix_list_34-1, decimal) "."
      counter(lst-ctn-kix_list_34-2, decimal) " ";
  }

  ol.lst-kix_list_54-1.start {
    counter-reset: lst-ctn-kix_list_54-1 0;
  }

  .lst-kix_list_34-5 > li {
    counter-increment: lst-ctn-kix_list_34-5;
  }

  .lst-kix_list_30-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_30-5, upper-latin) ") ";
  }

  .lst-kix_list_35-0 > li:before {
    content: "" counter(lst-ctn-kix_list_35-0, decimal) " ";
  }

  .lst-kix_list_35-1 > li:before {
    content: "" counter(lst-ctn-kix_list_35-1, decimal) ". ";
  }

  .lst-kix_list_35-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_35-4, lower-roman) ") ";
  }

  ol.lst-kix_list_40-6.start {
    counter-reset: lst-ctn-kix_list_40-6 0;
  }

  .lst-kix_list_28-8 > li {
    counter-increment: lst-ctn-kix_list_28-8;
  }

  .lst-kix_list_1-1 > li {
    counter-increment: lst-ctn-kix_list_1-1;
  }

  .lst-kix_list_30-1 > li:before {
    content: "" counter(lst-ctn-kix_list_30-1, decimal) ". ";
  }

  .lst-kix_list_3-0 > li:before {
    content: "" counter(lst-ctn-kix_list_3-0, decimal) ". ";
  }

  .lst-kix_list_30-2 > li:before {
    content: "" counter(lst-ctn-kix_list_30-0, decimal) "."
      counter(lst-ctn-kix_list_30-1, decimal) "."
      counter(lst-ctn-kix_list_30-2, decimal) " ";
  }

  .lst-kix_list_4-0 > li {
    counter-increment: lst-ctn-kix_list_4-0;
  }

  .lst-kix_list_31-6 > li {
    counter-increment: lst-ctn-kix_list_31-6;
  }

  ol.lst-kix_list_37-5.start {
    counter-reset: lst-ctn-kix_list_37-5 0;
  }

  .lst-kix_list_3-3 > li:before {
    content: "" counter(lst-ctn-kix_list_3-3, decimal) ". ";
  }

  .lst-kix_list_38-3 > li {
    counter-increment: lst-ctn-kix_list_38-3;
  }

  ol.lst-kix_list_10-7.start {
    counter-reset: lst-ctn-kix_list_10-7 0;
  }

  .lst-kix_list_53-2 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_53-3 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_17-1 > li {
    counter-increment: lst-ctn-kix_list_17-1;
  }

  .lst-kix_list_11-1 > li:before {
    content: "" counter(lst-ctn-kix_list_11-0, decimal) "."
      counter(lst-ctn-kix_list_11-1, decimal) ". ";
  }

  ol.lst-kix_list_31-0.start {
    counter-reset: lst-ctn-kix_list_31-0 0;
  }

  .lst-kix_list_11-0 > li:before {
    content: "" counter(lst-ctn-kix_list_11-0, decimal) " ";
  }

  ol.lst-kix_list_9-3.start {
    counter-reset: lst-ctn-kix_list_9-3 0;
  }

  .lst-kix_list_8-8 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_35-4 > li {
    counter-increment: lst-ctn-kix_list_35-4;
  }

  .lst-kix_list_43-5 > li {
    counter-increment: lst-ctn-kix_list_43-5;
  }

  .lst-kix_list_4-8 > li:before {
    content: "" counter(lst-ctn-kix_list_4-8, decimal) ". ";
  }

  .lst-kix_list_21-5 > li {
    counter-increment: lst-ctn-kix_list_21-5;
  }

  .lst-kix_list_14-2 > li {
    counter-increment: lst-ctn-kix_list_14-2;
  }

  .lst-kix_list_16-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_40-6 > li {
    counter-increment: lst-ctn-kix_list_40-6;
  }

  ol.lst-kix_list_4-8.start {
    counter-reset: lst-ctn-kix_list_4-8 0;
  }

  .lst-kix_list_16-3 > li:before {
    content: "" counter(lst-ctn-kix_list_16-3, decimal) ". ";
  }

  ol.lst-kix_list_37-4.start {
    counter-reset: lst-ctn-kix_list_37-4 0;
  }

  .lst-kix_list_41-1 > li {
    counter-increment: lst-ctn-kix_list_41-1;
  }

  .lst-kix_list_44-0 > li {
    counter-increment: lst-ctn-kix_list_44-0;
  }

  .lst-kix_list_17-8 > li:before {
    content: "" counter(lst-ctn-kix_list_17-8, decimal) ". ";
  }

  .lst-kix_list_7-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_27-7.start {
    counter-reset: lst-ctn-kix_list_27-7 0;
  }

  ol.lst-kix_list_19-6.start {
    counter-reset: lst-ctn-kix_list_19-6 0;
  }

  ol.lst-kix_list_9-7.start {
    counter-reset: lst-ctn-kix_list_9-7 0;
  }

  .lst-kix_list_2-4 > li:before {
    content: "" counter(lst-ctn-kix_list_2-4, decimal) ". ";
  }

  .lst-kix_list_7-3 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_48-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_9-7 > li {
    counter-increment: lst-ctn-kix_list_9-7;
  }

  .lst-kix_list_13-8 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) "."
      counter(lst-ctn-kix_list_13-3, decimal) "."
      counter(lst-ctn-kix_list_13-4, decimal) "."
      counter(lst-ctn-kix_list_13-5, decimal) "."
      counter(lst-ctn-kix_list_13-6, decimal) "."
      counter(lst-ctn-kix_list_13-7, decimal) "."
      counter(lst-ctn-kix_list_13-8, decimal) ". ";
  }

  .lst-kix_list_31-1 > li:before {
    content: "" counter(lst-ctn-kix_list_31-1, decimal) ". ";
  }

  .lst-kix_list_18-7 > li:before {
    content: "" counter(lst-ctn-kix_list_18-7, decimal) ". ";
  }

  .lst-kix_list_48-1 > li:before {
    content: "o  ";
  }

  .lst-kix_list_52-6 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_15-4 > li:before {
    content: "" counter(lst-ctn-kix_list_15-4, decimal) ") ";
  }

  ol.lst-kix_list_19-1.start {
    counter-reset: lst-ctn-kix_list_19-1 0;
  }

  .lst-kix_list_10-4 > li:before {
    content: "" counter(lst-ctn-kix_list_10-4, decimal) " ";
  }

  .lst-kix_list_10-8 > li:before {
    content: "" counter(lst-ctn-kix_list_10-8, decimal) " ";
  }

  .lst-kix_list_15-0 > li:before {
    content: "Article " counter(lst-ctn-kix_list_15-0, upper-roman) ". ";
  }

  ol.lst-kix_list_14-3.start {
    counter-reset: lst-ctn-kix_list_14-3 0;
  }

  ol.lst-kix_list_44-2.start {
    counter-reset: lst-ctn-kix_list_44-2 0;
  }

  ol.lst-kix_list_32-5.start {
    counter-reset: lst-ctn-kix_list_32-5 0;
  }

  .lst-kix_list_12-8 > li {
    counter-increment: lst-ctn-kix_list_12-8;
  }

  ol.lst-kix_list_14-4.start {
    counter-reset: lst-ctn-kix_list_14-4 0;
  }

  .lst-kix_list_9-7 > li:before {
    content: "" counter(lst-ctn-kix_list_9-7, decimal) " ";
  }

  .lst-kix_list_2-4 > li {
    counter-increment: lst-ctn-kix_list_2-4;
  }

  .lst-kix_list_29-4 > li:before {
    content: "" counter(lst-ctn-kix_list_29-4, decimal) ") ";
  }

  .lst-kix_list_53-6 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_32-2.start {
    counter-reset: lst-ctn-kix_list_32-2 0;
  }

  .lst-kix_list_11-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_11-4, lower-roman) ") ";
  }

  ol.lst-kix_list_37-0.start {
    counter-reset: lst-ctn-kix_list_37-0 0;
  }

  .lst-kix_list_29-0 > li:before {
    content: "Article " counter(lst-ctn-kix_list_29-0, upper-roman) ". ";
  }

  ol.lst-kix_list_19-2.start {
    counter-reset: lst-ctn-kix_list_19-2 0;
  }

  .lst-kix_list_54-2 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) " ";
  }

  .lst-kix_list_12-0 > li:before {
    content: "" counter(lst-ctn-kix_list_12-0, decimal) " ";
  }

  .lst-kix_list_1-4 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) "."
      counter(lst-ctn-kix_list_1-3, decimal) "."
      counter(lst-ctn-kix_list_1-4, decimal) " ";
  }

  ol.lst-kix_list_44-3.start {
    counter-reset: lst-ctn-kix_list_44-3 0;
  }

  .lst-kix_list_1-6 > li {
    counter-increment: lst-ctn-kix_list_1-6;
  }

  .lst-kix_list_34-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_34-5, upper-latin) ") ";
  }

  ol.lst-kix_list_45-8.start {
    counter-reset: lst-ctn-kix_list_45-8 0;
  }

  .lst-kix_list_33-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_33-5, upper-latin) ") ";
  }

  ol.lst-kix_list_19-3.start {
    counter-reset: lst-ctn-kix_list_19-3 0;
  }

  .lst-kix_list_2-0 > li:before {
    content: "" counter(lst-ctn-kix_list_2-0, decimal) ". ";
  }

  ol.lst-kix_list_9-8.start {
    counter-reset: lst-ctn-kix_list_9-8 0;
  }

  .lst-kix_list_1-8 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) "."
      counter(lst-ctn-kix_list_1-3, decimal) "."
      counter(lst-ctn-kix_list_1-4, decimal) "."
      counter(lst-ctn-kix_list_1-5, decimal) "."
      counter(lst-ctn-kix_list_1-6, decimal) "."
      counter(lst-ctn-kix_list_1-7, decimal) "."
      counter(lst-ctn-kix_list_1-8, decimal) " ";
  }

  .lst-kix_list_49-1 > li:before {
    content: "o  ";
  }

  .lst-kix_list_34-1 > li:before {
    content: "" counter(lst-ctn-kix_list_34-1, decimal) ". ";
  }

  .lst-kix_list_19-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_44-6.start {
    counter-reset: lst-ctn-kix_list_44-6 0;
  }

  .lst-kix_list_29-6 > li {
    counter-increment: lst-ctn-kix_list_29-6;
  }

  .lst-kix_list_47-2 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_31-3.start {
    counter-reset: lst-ctn-kix_list_31-3 0;
  }

  ol.lst-kix_list_19-8.start {
    counter-reset: lst-ctn-kix_list_19-8 0;
  }

  ol.lst-kix_list_36-4.start {
    counter-reset: lst-ctn-kix_list_36-4 0;
  }

  ol.lst-kix_list_14-7.start {
    counter-reset: lst-ctn-kix_list_14-7 0;
  }

  .lst-kix_list_19-7 > li:before {
    content: "" counter(lst-ctn-kix_list_19-7, decimal) ". ";
  }

  .lst-kix_list_9-2 > li {
    counter-increment: lst-ctn-kix_list_9-2;
  }

  ol.lst-kix_list_23-8.start {
    counter-reset: lst-ctn-kix_list_23-8 0;
  }

  .lst-kix_list_24-5 > li {
    counter-increment: lst-ctn-kix_list_24-5;
  }

  .lst-kix_list_46-3 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_38-1 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) " ";
  }

  ol.lst-kix_list_15-6.start {
    counter-reset: lst-ctn-kix_list_15-6 0;
  }

  .lst-kix_list_37-2 > li:before {
    content: "" counter(lst-ctn-kix_list_37-0, decimal) "."
      counter(lst-ctn-kix_list_37-1, decimal) "."
      counter(lst-ctn-kix_list_37-2, decimal) " ";
  }

  .lst-kix_list_25-6 > li {
    counter-increment: lst-ctn-kix_list_25-6;
  }

  .lst-kix_list_37-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_37-5, upper-latin) ") ";
  }

  .lst-kix_list_46-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_31-8.start {
    counter-reset: lst-ctn-kix_list_31-8 0;
  }

  .lst-kix_list_18-1 > li:before {
    content: "" counter(lst-ctn-kix_list_18-1, decimal) ". ";
  }

  .lst-kix_list_38-8 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) "."
      counter(lst-ctn-kix_list_38-3, decimal) "."
      counter(lst-ctn-kix_list_38-4, decimal) "."
      counter(lst-ctn-kix_list_38-5, decimal) "."
      counter(lst-ctn-kix_list_38-6, decimal) "."
      counter(lst-ctn-kix_list_38-7, decimal) "."
      counter(lst-ctn-kix_list_38-8, decimal) " ";
  }

  .lst-kix_list_40-3 > li {
    counter-increment: lst-ctn-kix_list_40-3;
  }

  .lst-kix_list_44-3 > li {
    counter-increment: lst-ctn-kix_list_44-3;
  }

  .lst-kix_list_45-4 > li {
    counter-increment: lst-ctn-kix_list_45-4;
  }

  .lst-kix_list_23-4 > li {
    counter-increment: lst-ctn-kix_list_23-4;
  }

  ol.lst-kix_list_23-1.start {
    counter-reset: lst-ctn-kix_list_23-1 0;
  }

  ol.lst-kix_list_45-5.start {
    counter-reset: lst-ctn-kix_list_45-5 0;
  }

  ol.lst-kix_list_32-0.start {
    counter-reset: lst-ctn-kix_list_32-0 0;
  }

  .lst-kix_list_23-1 > li {
    counter-increment: lst-ctn-kix_list_23-1;
  }

  .lst-kix_list_2-7 > li:before {
    content: "" counter(lst-ctn-kix_list_2-7, decimal) ". ";
  }

  .lst-kix_list_2-7 > li {
    counter-increment: lst-ctn-kix_list_2-7;
  }

  .lst-kix_list_24-2 > li {
    counter-increment: lst-ctn-kix_list_24-2;
  }

  .lst-kix_list_27-5 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) "."
      counter(lst-ctn-kix_list_27-3, decimal) "."
      counter(lst-ctn-kix_list_27-4, decimal) "."
      counter(lst-ctn-kix_list_27-5, decimal) ". ";
  }

  ol.lst-kix_list_10-3.start {
    counter-reset: lst-ctn-kix_list_10-3 0;
  }

  .lst-kix_list_22-3 > li {
    counter-increment: lst-ctn-kix_list_22-3;
  }

  .lst-kix_list_39-7 > li:before {
    content: "" counter(lst-ctn-kix_list_39-7, decimal) " ";
  }

  .lst-kix_list_55-7 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_10-1 > li:before {
    content: "(" counter(lst-ctn-kix_list_10-1, lower-latin) ") ";
  }

  .lst-kix_list_18-4 > li:before {
    content: "" counter(lst-ctn-kix_list_18-4, decimal) ". ";
  }

  ol.lst-kix_list_15-1.start {
    counter-reset: lst-ctn-kix_list_15-1 0;
  }

  ol.lst-kix_list_15-4.start {
    counter-reset: lst-ctn-kix_list_15-4 0;
  }

  ol.lst-kix_list_45-3.start {
    counter-reset: lst-ctn-kix_list_45-3 0;
  }

  .lst-kix_list_36-6 > li:before {
    content: "" counter(lst-ctn-kix_list_36-6, decimal) ". ";
  }

  .lst-kix_list_50-6 > li {
    counter-increment: lst-ctn-kix_list_50-6;
  }

  .lst-kix_list_36-0 > li:before {
    content: "" counter(lst-ctn-kix_list_36-0, upper-latin) " ";
  }

  ol.lst-kix_list_40-2.start {
    counter-reset: lst-ctn-kix_list_40-2 0;
  }

  ol.lst-kix_list_44-8.start {
    counter-reset: lst-ctn-kix_list_44-8 0;
  }

  .lst-kix_list_34-0 > li {
    counter-increment: lst-ctn-kix_list_34-0;
  }

  ol.lst-kix_list_45-0.start {
    counter-reset: lst-ctn-kix_list_45-0 0;
  }

  .lst-kix_list_26-4 > li {
    counter-increment: lst-ctn-kix_list_26-4;
  }

  ol.lst-kix_list_5-7.start {
    counter-reset: lst-ctn-kix_list_5-7 0;
  }

  ol.lst-kix_list_28-7.start {
    counter-reset: lst-ctn-kix_list_28-7 0;
  }

  .lst-kix_list_20-8 > li:before {
    content: "" counter(lst-ctn-kix_list_20-8, decimal) ". ";
  }

  .lst-kix_list_39-1 > li {
    counter-increment: lst-ctn-kix_list_39-1;
  }

  .lst-kix_list_46-6 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_29-7 > li:before {
    content: "" counter(lst-ctn-kix_list_29-7, lower-latin) ". ";
  }

  .lst-kix_list_9-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_9-4, lower-roman) ") ";
  }

  .lst-kix_list_30-8 > li {
    counter-increment: lst-ctn-kix_list_30-8;
  }

  ol.lst-kix_list_36-2.start {
    counter-reset: lst-ctn-kix_list_36-2 0;
  }

  .lst-kix_list_1-1 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) " ";
  }

  .lst-kix_list_11-7 > li:before {
    content: "" counter(lst-ctn-kix_list_11-7, decimal) " ";
  }

  .lst-kix_list_49-4 > li:before {
    content: "o  ";
  }

  .lst-kix_list_55-1 > li:before {
    content: "o  ";
  }

  ol.lst-kix_list_10-5.start {
    counter-reset: lst-ctn-kix_list_10-5 0;
  }

  .lst-kix_list_54-5 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) "."
      counter(lst-ctn-kix_list_54-3, decimal) "."
      counter(lst-ctn-kix_list_54-4, decimal) "."
      counter(lst-ctn-kix_list_54-5, decimal) " ";
  }

  .lst-kix_list_14-7 > li {
    counter-increment: lst-ctn-kix_list_14-7;
  }

  ol.lst-kix_list_31-5.start {
    counter-reset: lst-ctn-kix_list_31-5 0;
  }

  .lst-kix_list_48-8 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_40-4.start {
    counter-reset: lst-ctn-kix_list_40-4 0;
  }

  .lst-kix_list_28-1 > li:before {
    content: "" counter(lst-ctn-kix_list_28-1, lower-latin) ") ";
  }

  .lst-kix_list_30-8 > li:before {
    content: "" counter(lst-ctn-kix_list_30-8, decimal) " ";
  }

  .lst-kix_list_35-7 > li:before {
    content: "" counter(lst-ctn-kix_list_35-7, decimal) " ";
  }

  .lst-kix_list_30-5 > li {
    counter-increment: lst-ctn-kix_list_30-5;
  }

  .lst-kix_list_26-6 > li:before {
    content: "" counter(lst-ctn-kix_list_26-6, lower-roman) ") ";
  }

  .lst-kix_list_8-2 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_36-7.start {
    counter-reset: lst-ctn-kix_list_36-7 0;
  }

  ol.lst-kix_list_31-6.start {
    counter-reset: lst-ctn-kix_list_31-6 0;
  }

  .lst-kix_list_12-0 > li {
    counter-increment: lst-ctn-kix_list_12-0;
  }

  ol.lst-kix_list_45-1.start {
    counter-reset: lst-ctn-kix_list_45-1 0;
  }

  ol.lst-kix_list_40-0.start {
    counter-reset: lst-ctn-kix_list_40-0 0;
  }

  .lst-kix_list_8-5 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_28-0.start {
    counter-reset: lst-ctn-kix_list_28-0 0;
  }

  .lst-kix_list_26-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_26-3, lower-roman) ") ";
  }

  .lst-kix_list_3-6 > li:before {
    content: "" counter(lst-ctn-kix_list_3-6, decimal) ". ";
  }

  .lst-kix_list_21-7 > li:before {
    content: "" counter(lst-ctn-kix_list_21-7, decimal) ". ";
  }

  ol.lst-kix_list_5-0.start {
    counter-reset: lst-ctn-kix_list_5-0 0;
  }

  .lst-kix_list_39-4 > li {
    counter-increment: lst-ctn-kix_list_39-4;
  }

  ol.lst-kix_list_31-7.start {
    counter-reset: lst-ctn-kix_list_31-7 0;
  }

  .lst-kix_list_16-6 > li:before {
    content: "" counter(lst-ctn-kix_list_16-6, decimal) ". ";
  }

  ol.lst-kix_list_10-1.start {
    counter-reset: lst-ctn-kix_list_10-1 0;
  }

  .lst-kix_list_18-2 > li {
    counter-increment: lst-ctn-kix_list_18-2;
  }

  ol.lst-kix_list_5-6.start {
    counter-reset: lst-ctn-kix_list_5-6 0;
  }

  .lst-kix_list_22-6 > li {
    counter-increment: lst-ctn-kix_list_22-6;
  }

  .lst-kix_list_25-0 > li {
    counter-increment: lst-ctn-kix_list_25-0;
  }

  .lst-kix_list_45-4 > li:before {
    content: "" counter(lst-ctn-kix_list_45-4, decimal) " ";
  }

  .lst-kix_list_45-7 > li:before {
    content: "" counter(lst-ctn-kix_list_45-7, decimal) " ";
  }

  ol.lst-kix_list_28-6.start {
    counter-reset: lst-ctn-kix_list_28-6 0;
  }

  .lst-kix_list_19-3 > li {
    counter-increment: lst-ctn-kix_list_19-3;
  }

  ol.lst-kix_list_28-5.start {
    counter-reset: lst-ctn-kix_list_28-5 0;
  }

  .lst-kix_list_44-8 > li:before {
    content: "" counter(lst-ctn-kix_list_44-8, decimal) " ";
  }

  .lst-kix_list_23-7 > li {
    counter-increment: lst-ctn-kix_list_23-7;
  }

  .lst-kix_list_44-6 > li {
    counter-increment: lst-ctn-kix_list_44-6;
  }

  ol.lst-kix_list_10-2.start {
    counter-reset: lst-ctn-kix_list_10-2 0;
  }

  .lst-kix_list_40-0 > li {
    counter-increment: lst-ctn-kix_list_40-0;
  }

  .lst-kix_list_17-2 > li:before {
    content: "" counter(lst-ctn-kix_list_17-2, decimal) ". ";
  }

  ol.lst-kix_list_5-5.start {
    counter-reset: lst-ctn-kix_list_5-5 0;
  }

  ol.lst-kix_list_40-1.start {
    counter-reset: lst-ctn-kix_list_40-1 0;
  }

  .lst-kix_list_45-7 > li {
    counter-increment: lst-ctn-kix_list_45-7;
  }

  .lst-kix_list_17-5 > li:before {
    content: "" counter(lst-ctn-kix_list_17-5, decimal) ". ";
  }

  .lst-kix_list_6-2 > li {
    counter-increment: lst-ctn-kix_list_6-2;
  }

  ol.lst-kix_list_36-6.start {
    counter-reset: lst-ctn-kix_list_36-6 0;
  }

  .lst-kix_list_27-2 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) ". ";
  }

  .lst-kix_list_22-3 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) "."
      counter(lst-ctn-kix_list_22-3, decimal) ". ";
  }

  ol.lst-kix_list_28-4.start {
    counter-reset: lst-ctn-kix_list_28-4 0;
  }

  .lst-kix_list_43-8 > li {
    counter-increment: lst-ctn-kix_list_43-8;
  }

  ol.lst-kix_list_23-3.start {
    counter-reset: lst-ctn-kix_list_23-3 0;
  }

  ol.lst-kix_list_5-4.start {
    counter-reset: lst-ctn-kix_list_5-4 0;
  }

  .lst-kix_list_52-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_5-1.start {
    counter-reset: lst-ctn-kix_list_5-1 0;
  }

  .lst-kix_list_25-3 > li {
    counter-increment: lst-ctn-kix_list_25-3;
  }

  .lst-kix_list_16-6 > li {
    counter-increment: lst-ctn-kix_list_16-6;
  }

  ol.lst-kix_list_28-1.start {
    counter-reset: lst-ctn-kix_list_28-1 0;
  }

  .lst-kix_list_31-7 > li:before {
    content: "" counter(lst-ctn-kix_list_31-7, decimal) " ";
  }

  ol.lst-kix_list_10-0.start {
    counter-reset: lst-ctn-kix_list_10-0 0;
  }

  .lst-kix_list_3-8 > li {
    counter-increment: lst-ctn-kix_list_3-8;
  }

  .lst-kix_list_4-6 > li {
    counter-increment: lst-ctn-kix_list_4-6;
  }

  .lst-kix_list_51-7 > li {
    counter-increment: lst-ctn-kix_list_51-7;
  }

  .lst-kix_list_4-2 > li:before {
    content: "" counter(lst-ctn-kix_list_4-2, decimal) ". ";
  }

  ol.lst-kix_list_23-6.start {
    counter-reset: lst-ctn-kix_list_23-6 0;
  }

  .lst-kix_list_17-4 > li {
    counter-increment: lst-ctn-kix_list_17-4;
  }

  .lst-kix_list_36-3 > li:before {
    content: "" counter(lst-ctn-kix_list_36-3, decimal) ". ";
  }

  .lst-kix_list_26-1 > li {
    counter-increment: lst-ctn-kix_list_26-1;
  }

  .lst-kix_list_9-1 > li:before {
    content: "" counter(lst-ctn-kix_list_9-1, decimal) ". ";
  }

  .lst-kix_list_15-8 > li {
    counter-increment: lst-ctn-kix_list_15-8;
  }

  .lst-kix_list_40-8 > li:before {
    content: "" counter(lst-ctn-kix_list_40-8, decimal) " ";
  }

  ol.lst-kix_list_36-8.start {
    counter-reset: lst-ctn-kix_list_36-8 0;
  }

  .lst-kix_list_54-1 > li {
    counter-increment: lst-ctn-kix_list_54-1;
  }

  .lst-kix_list_37-8 > li {
    counter-increment: lst-ctn-kix_list_37-8;
  }

  .lst-kix_list_31-3 > li {
    counter-increment: lst-ctn-kix_list_31-3;
  }

  .lst-kix_list_41-4 > li:before {
    content: "" counter(lst-ctn-kix_list_41-4, decimal) " ";
  }

  ol.lst-kix_list_23-5.start {
    counter-reset: lst-ctn-kix_list_23-5 0;
  }

  .lst-kix_list_10-1 > li {
    counter-increment: lst-ctn-kix_list_10-1;
  }

  .lst-kix_list_49-7 > li:before {
    content: "o  ";
  }

  ol.lst-kix_list_28-2.start {
    counter-reset: lst-ctn-kix_list_28-2 0;
  }

  .lst-kix_list_12-6 > li:before {
    content: "" counter(lst-ctn-kix_list_12-6, decimal) ". ";
  }

  ol.lst-kix_list_23-4.start {
    counter-reset: lst-ctn-kix_list_23-4 0;
  }

  .lst-kix_list_55-4 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_32-1 > li {
    counter-increment: lst-ctn-kix_list_32-1;
  }

  .lst-kix_list_54-8 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) "."
      counter(lst-ctn-kix_list_54-3, decimal) "."
      counter(lst-ctn-kix_list_54-4, decimal) "."
      counter(lst-ctn-kix_list_54-5, decimal) "."
      counter(lst-ctn-kix_list_54-6, decimal) "."
      counter(lst-ctn-kix_list_54-7, decimal) "."
      counter(lst-ctn-kix_list_54-8, decimal) " ";
  }

  .lst-kix_list_13-2 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) ". ";
  }

  ol.lst-kix_list_5-2.start {
    counter-reset: lst-ctn-kix_list_5-2 0;
  }

  .lst-kix_list_2-1 > li {
    counter-increment: lst-ctn-kix_list_2-1;
  }

  ol.lst-kix_list_33-5.start {
    counter-reset: lst-ctn-kix_list_33-5 0;
  }

  ol.lst-kix_list_9-0.start {
    counter-reset: lst-ctn-kix_list_9-0 0;
  }

  ol.lst-kix_list_30-0.start {
    counter-reset: lst-ctn-kix_list_30-0 0;
  }

  ol.lst-kix_list_13-4.start {
    counter-reset: lst-ctn-kix_list_13-4 0;
  }

  .lst-kix_list_30-6 > li {
    counter-increment: lst-ctn-kix_list_30-6;
  }

  ol.lst-kix_list_40-3.start {
    counter-reset: lst-ctn-kix_list_40-3 0;
  }

  .lst-kix_list_42-8 > li:before {
    content: "" counter(lst-ctn-kix_list_42-8, decimal) " ";
  }

  ol.lst-kix_list_50-1.start {
    counter-reset: lst-ctn-kix_list_50-1 0;
  }

  .lst-kix_list_13-0 > li {
    counter-increment: lst-ctn-kix_list_13-0;
  }

  ol.lst-kix_list_43-8.start {
    counter-reset: lst-ctn-kix_list_43-8 0;
  }

  ol.lst-kix_list_20-2.start {
    counter-reset: lst-ctn-kix_list_20-2 0;
  }

  .lst-kix_list_42-0 > li:before {
    content: "" counter(lst-ctn-kix_list_42-0, decimal) " ";
  }

  ol.lst-kix_list_23-2.start {
    counter-reset: lst-ctn-kix_list_23-2 0;
  }

  .lst-kix_list_54-0 > li {
    counter-increment: lst-ctn-kix_list_54-0;
  }

  .lst-kix_list_42-2 > li:before {
    content: "" counter(lst-ctn-kix_list_42-2, decimal) " ";
  }

  .lst-kix_list_43-0 > li {
    counter-increment: lst-ctn-kix_list_43-0;
  }

  ol.lst-kix_list_50-6.start {
    counter-reset: lst-ctn-kix_list_50-6 0;
  }

  .lst-kix_list_35-7 > li {
    counter-increment: lst-ctn-kix_list_35-7;
  }

  .lst-kix_list_29-8 > li {
    counter-increment: lst-ctn-kix_list_29-8;
  }

  .lst-kix_list_42-6 > li:before {
    content: "" counter(lst-ctn-kix_list_42-6, decimal) " ";
  }

  .lst-kix_list_24-7 > li {
    counter-increment: lst-ctn-kix_list_24-7;
  }

  ol.lst-kix_list_30-5.start {
    counter-reset: lst-ctn-kix_list_30-5 0;
  }

  .lst-kix_list_42-4 > li:before {
    content: "" counter(lst-ctn-kix_list_42-4, decimal) " ";
  }

  .lst-kix_list_18-8 > li {
    counter-increment: lst-ctn-kix_list_18-8;
  }

  ol.lst-kix_list_26-7.start {
    counter-reset: lst-ctn-kix_list_26-7 0;
  }

  ol.lst-kix_list_10-4.start {
    counter-reset: lst-ctn-kix_list_10-4 0;
  }

  ol.lst-kix_list_36-5.start {
    counter-reset: lst-ctn-kix_list_36-5 0;
  }

  .lst-kix_list_24-7 > li:before {
    content: "" counter(lst-ctn-kix_list_24-7, lower-latin) ". ";
  }

  .lst-kix_list_1-4 > li {
    counter-increment: lst-ctn-kix_list_1-4;
  }

  ol.lst-kix_list_1-6.start {
    counter-reset: lst-ctn-kix_list_1-6 0;
  }

  ol.lst-kix_list_9-5.start {
    counter-reset: lst-ctn-kix_list_9-5 0;
  }

  ol.lst-kix_list_40-8.start {
    counter-reset: lst-ctn-kix_list_40-8 0;
  }

  ol.lst-kix_list_20-7.start {
    counter-reset: lst-ctn-kix_list_20-7 0;
  }

  .lst-kix_list_24-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_24-3, lower-roman) ") ";
  }

  .lst-kix_list_24-5 > li:before {
    content: "" counter(lst-ctn-kix_list_24-5, lower-latin) ") ";
  }

  ol.lst-kix_list_43-3.start {
    counter-reset: lst-ctn-kix_list_43-3 0;
  }

  ol.lst-kix_list_16-4.start {
    counter-reset: lst-ctn-kix_list_16-4 0;
  }

  .lst-kix_list_6-5 > li {
    counter-increment: lst-ctn-kix_list_6-5;
  }

  .lst-kix_list_23-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_23-3, decimal) ") ";
  }

  .lst-kix_list_23-7 > li:before {
    content: "" counter(lst-ctn-kix_list_23-7, lower-latin) ". ";
  }

  ol.lst-kix_list_23-7.start {
    counter-reset: lst-ctn-kix_list_23-7 0;
  }

  .lst-kix_list_23-1 > li:before {
    content: "" counter(lst-ctn-kix_list_23-1, lower-latin) ") ";
  }

  .lst-kix_list_24-1 > li:before {
    content: "Section " counter(lst-ctn-kix_list_24-0, upper-roman) "."
      counter(lst-ctn-kix_list_24-1, decimal) " ";
  }

  .lst-kix_list_2-8 > li {
    counter-increment: lst-ctn-kix_list_2-8;
  }

  ol.lst-kix_list_26-2.start {
    counter-reset: lst-ctn-kix_list_26-2 0;
  }

  .lst-kix_list_23-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_23-5, lower-roman) ") ";
  }

  ol.lst-kix_list_36-0.start {
    counter-reset: lst-ctn-kix_list_36-0 0;
  }

  ol.lst-kix_list_4-6.start {
    counter-reset: lst-ctn-kix_list_4-6 0;
  }

  ol.lst-kix_list_39-5.start {
    counter-reset: lst-ctn-kix_list_39-5 0;
  }

  ol.lst-kix_list_3-0.start {
    counter-reset: lst-ctn-kix_list_3-0 0;
  }

  ol.lst-kix_list_29-2.start {
    counter-reset: lst-ctn-kix_list_29-2 0;
  }

  .lst-kix_list_25-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_25-5, lower-roman) ") ";
  }

  .lst-kix_list_25-7 > li:before {
    content: "" counter(lst-ctn-kix_list_25-7, lower-latin) ". ";
  }

  .lst-kix_list_3-5 > li {
    counter-increment: lst-ctn-kix_list_3-5;
  }

  ol.lst-kix_list_1-1.start {
    counter-reset: lst-ctn-kix_list_1-1 0;
  }

  ol.lst-kix_list_18-3.start {
    counter-reset: lst-ctn-kix_list_18-3 0;
  }

  ol.lst-kix_list_33-0.start {
    counter-reset: lst-ctn-kix_list_33-0 0;
  }

  ol.lst-kix_list_45-2.start {
    counter-reset: lst-ctn-kix_list_45-2 0;
  }

  .lst-kix_list_9-3 > li {
    counter-increment: lst-ctn-kix_list_9-3;
  }

  ol.lst-kix_list_15-3.start {
    counter-reset: lst-ctn-kix_list_15-3 0;
  }

  ol.lst-kix_list_6-0.start {
    counter-reset: lst-ctn-kix_list_6-0 0;
  }

  .lst-kix_list_4-2 > li {
    counter-increment: lst-ctn-kix_list_4-2;
  }

  .lst-kix_list_5-1 > li {
    counter-increment: lst-ctn-kix_list_5-1;
  }

  ol.lst-kix_list_21-3.start {
    counter-reset: lst-ctn-kix_list_21-3 0;
  }

  .lst-kix_list_38-5 > li {
    counter-increment: lst-ctn-kix_list_38-5;
  }

  ol.lst-kix_list_25-6.start {
    counter-reset: lst-ctn-kix_list_25-6 0;
  }

  ol.lst-kix_list_51-7.start {
    counter-reset: lst-ctn-kix_list_51-7 0;
  }

  ol.lst-kix_list_32-4.start {
    counter-reset: lst-ctn-kix_list_32-4 0;
  }

  .lst-kix_list_26-5 > li:before {
    content: "" counter(lst-ctn-kix_list_26-5, lower-latin) ") ";
  }

  .lst-kix_list_10-0 > li {
    counter-increment: lst-ctn-kix_list_10-0;
  }

  .lst-kix_list_26-1 > li:before {
    content: "Section " counter(lst-ctn-kix_list_26-0, upper-roman) "."
      counter(lst-ctn-kix_list_26-1, decimal) " ";
  }

  .lst-kix_list_21-1 > li:before {
    content: "" counter(lst-ctn-kix_list_21-1, decimal) ". ";
  }

  ol.lst-kix_list_15-8.start {
    counter-reset: lst-ctn-kix_list_15-8 0;
  }

  .lst-kix_list_34-3 > li {
    counter-increment: lst-ctn-kix_list_34-3;
  }

  .lst-kix_list_45-3 > li {
    counter-increment: lst-ctn-kix_list_45-3;
  }

  .lst-kix_list_10-2 > li {
    counter-increment: lst-ctn-kix_list_10-2;
  }

  ol.lst-kix_list_45-7.start {
    counter-reset: lst-ctn-kix_list_45-7 0;
  }

  .lst-kix_list_21-5 > li:before {
    content: "" counter(lst-ctn-kix_list_21-5, decimal) ". ";
  }

  .lst-kix_list_45-2 > li:before {
    content: "" counter(lst-ctn-kix_list_45-2, decimal) " ";
  }

  ol.lst-kix_list_27-8.start {
    counter-reset: lst-ctn-kix_list_27-8 0;
  }

  .lst-kix_list_21-0 > li {
    counter-increment: lst-ctn-kix_list_21-0;
  }

  .lst-kix_list_25-1 > li:before {
    content: "" counter(lst-ctn-kix_list_25-1, lower-latin) ") ";
  }

  ol.lst-kix_list_51-2.start {
    counter-reset: lst-ctn-kix_list_51-2 0;
  }

  ol.lst-kix_list_31-1.start {
    counter-reset: lst-ctn-kix_list_31-1 0;
  }

  .lst-kix_list_45-6 > li:before {
    content: "" counter(lst-ctn-kix_list_45-6, decimal) " ";
  }

  .lst-kix_list_44-6 > li:before {
    content: "" counter(lst-ctn-kix_list_44-6, decimal) " ";
  }

  .lst-kix_list_44-2 > li:before {
    content: "" counter(lst-ctn-kix_list_44-2, decimal) " ";
  }

  .lst-kix_list_39-1 > li:before {
    content: "" counter(lst-ctn-kix_list_39-0, decimal) "."
      counter(lst-ctn-kix_list_39-1, decimal) ". ";
  }

  .lst-kix_list_16-7 > li {
    counter-increment: lst-ctn-kix_list_16-7;
  }

  .lst-kix_list_27-7 > li {
    counter-increment: lst-ctn-kix_list_27-7;
  }

  .lst-kix_list_45-1 > li {
    counter-increment: lst-ctn-kix_list_45-1;
  }

  .lst-kix_list_34-1 > li {
    counter-increment: lst-ctn-kix_list_34-1;
  }

  .lst-kix_list_38-7 > li {
    counter-increment: lst-ctn-kix_list_38-7;
  }

  .lst-kix_list_16-5 > li {
    counter-increment: lst-ctn-kix_list_16-5;
  }

  .lst-kix_list_3-7 > li {
    counter-increment: lst-ctn-kix_list_3-7;
  }

  .lst-kix_list_50-5 > li {
    counter-increment: lst-ctn-kix_list_50-5;
  }

  .lst-kix_list_22-5 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) "."
      counter(lst-ctn-kix_list_22-3, decimal) "."
      counter(lst-ctn-kix_list_22-4, decimal) "."
      counter(lst-ctn-kix_list_22-5, decimal) ". ";
  }

  .lst-kix_list_43-2 > li:before {
    content: "" counter(lst-ctn-kix_list_43-2, decimal) " ";
  }

  .lst-kix_list_43-6 > li:before {
    content: "" counter(lst-ctn-kix_list_43-6, decimal) " ";
  }

  .lst-kix_list_21-2 > li {
    counter-increment: lst-ctn-kix_list_21-2;
  }

  .lst-kix_list_45-8 > li {
    counter-increment: lst-ctn-kix_list_45-8;
  }

  .lst-kix_list_22-1 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) ". ";
  }

  .lst-kix_list_23-3 > li {
    counter-increment: lst-ctn-kix_list_23-3;
  }

  .lst-kix_list_9-5 > li {
    counter-increment: lst-ctn-kix_list_9-5;
  }

  .lst-kix_list_5-8 > li {
    counter-increment: lst-ctn-kix_list_5-8;
  }

  .lst-kix_list_27-0 > li {
    counter-increment: lst-ctn-kix_list_27-0;
  }

  .lst-kix_list_41-6 > li {
    counter-increment: lst-ctn-kix_list_41-6;
  }

  ol.lst-kix_list_19-4.start {
    counter-reset: lst-ctn-kix_list_19-4 0;
  }

  .lst-kix_list_38-0 > li {
    counter-increment: lst-ctn-kix_list_38-0;
  }

  .lst-kix_list_32-0 > li {
    counter-increment: lst-ctn-kix_list_32-0;
  }

  ol.lst-kix_list_2-2.start {
    counter-reset: lst-ctn-kix_list_2-2 0;
  }

  .lst-kix_list_50-3 > li {
    counter-increment: lst-ctn-kix_list_50-3;
  }

  ol.lst-kix_list_38-4.start {
    counter-reset: lst-ctn-kix_list_38-4 0;
  }

  ol.lst-kix_list_25-1.start {
    counter-reset: lst-ctn-kix_list_25-1 0;
  }

  ol.lst-kix_list_21-8.start {
    counter-reset: lst-ctn-kix_list_21-8 0;
  }

  .lst-kix_list_34-8 > li {
    counter-increment: lst-ctn-kix_list_34-8;
  }

  .lst-kix_list_27-5 > li {
    counter-increment: lst-ctn-kix_list_27-5;
  }

  .lst-kix_list_40-6 > li:before {
    content: "" counter(lst-ctn-kix_list_40-6, decimal) " ";
  }

  ol.lst-kix_list_44-4.start {
    counter-reset: lst-ctn-kix_list_44-4 0;
  }

  .lst-kix_list_41-6 > li:before {
    content: "" counter(lst-ctn-kix_list_41-6, decimal) " ";
  }

  .lst-kix_list_20-5 > li:before {
    content: "" counter(lst-ctn-kix_list_20-5, decimal) ". ";
  }

  .lst-kix_list_41-2 > li:before {
    content: "" counter(lst-ctn-kix_list_41-2, decimal) " ";
  }

  .lst-kix_list_20-1 > li:before {
    content: "" counter(lst-ctn-kix_list_20-1, decimal) ". ";
  }

  .lst-kix_list_12-3 > li {
    counter-increment: lst-ctn-kix_list_12-3;
  }

  ol.lst-kix_list_3-5.start {
    counter-reset: lst-ctn-kix_list_3-5 0;
  }

  .lst-kix_list_21-7 > li {
    counter-increment: lst-ctn-kix_list_21-7;
  }

  .lst-kix_list_14-4 > li {
    counter-increment: lst-ctn-kix_list_14-4;
  }

  .lst-kix_list_10-7 > li {
    counter-increment: lst-ctn-kix_list_10-7;
  }

  .lst-kix_list_54-7 > li {
    counter-increment: lst-ctn-kix_list_54-7;
  }

  .lst-kix_list_25-4 > li {
    counter-increment: lst-ctn-kix_list_25-4;
  }

  .lst-kix_list_18-1 > li {
    counter-increment: lst-ctn-kix_list_18-1;
  }

  ol.lst-kix_list_14-5.start {
    counter-reset: lst-ctn-kix_list_14-5 0;
  }

  .lst-kix_list_32-7 > li {
    counter-increment: lst-ctn-kix_list_32-7;
  }

  .lst-kix_list_40-2 > li:before {
    content: "" counter(lst-ctn-kix_list_40-2, decimal) " ";
  }

  ol.lst-kix_list_37-1.start {
    counter-reset: lst-ctn-kix_list_37-1 0;
  }

  .lst-kix_list_36-4 > li {
    counter-increment: lst-ctn-kix_list_36-4;
  }

  .lst-kix_list_29-1 > li {
    counter-increment: lst-ctn-kix_list_29-1;
  }

  .lst-kix_list_43-7 > li {
    counter-increment: lst-ctn-kix_list_43-7;
  }

  ol.lst-kix_list_32-7.start {
    counter-reset: lst-ctn-kix_list_32-7 0;
  }

  .lst-kix_list_43-2 > li {
    counter-increment: lst-ctn-kix_list_43-2;
  }

  ol.lst-kix_list_41-1.start {
    counter-reset: lst-ctn-kix_list_41-1 0;
  }

  .lst-kix_list_19-4 > li:before {
    content: "" counter(lst-ctn-kix_list_19-4, decimal) ". ";
  }

  ol.lst-kix_list_37-8.start {
    counter-reset: lst-ctn-kix_list_37-8 0;
  }

  .lst-kix_list_31-5 > li {
    counter-increment: lst-ctn-kix_list_31-5;
  }

  ol.lst-kix_list_24-5.start {
    counter-reset: lst-ctn-kix_list_24-5 0;
  }

  .lst-kix_list_47-6 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_47-8 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_27-5.start {
    counter-reset: lst-ctn-kix_list_27-5 0;
  }

  ol.lst-kix_list_22-4.start {
    counter-reset: lst-ctn-kix_list_22-4 0;
  }

  .lst-kix_list_42-1 > li {
    counter-increment: lst-ctn-kix_list_42-1;
  }

  ol.lst-kix_list_11-2.start {
    counter-reset: lst-ctn-kix_list_11-2 0;
  }

  .lst-kix_list_19-6 > li:before {
    content: "" counter(lst-ctn-kix_list_19-6, decimal) ". ";
  }

  .lst-kix_list_47-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_17-2 > li {
    counter-increment: lst-ctn-kix_list_17-2;
  }

  ol.lst-kix_list_41-6.start {
    counter-reset: lst-ctn-kix_list_41-6 0;
  }

  .lst-kix_list_42-5 > li {
    counter-increment: lst-ctn-kix_list_42-5;
  }

  .lst-kix_list_32-2 > li {
    counter-increment: lst-ctn-kix_list_32-2;
  }

  .lst-kix_list_36-2 > li {
    counter-increment: lst-ctn-kix_list_36-2;
  }

  .lst-kix_list_20-5 > li {
    counter-increment: lst-ctn-kix_list_20-5;
  }

  ol.lst-kix_list_37-3.start {
    counter-reset: lst-ctn-kix_list_37-3 0;
  }

  .lst-kix_list_37-1 > li:before {
    content: "" counter(lst-ctn-kix_list_37-0, decimal) "."
      counter(lst-ctn-kix_list_37-1, decimal) ". ";
  }

  .lst-kix_list_46-2 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_44-1.start {
    counter-reset: lst-ctn-kix_list_44-1 0;
  }

  .lst-kix_list_54-2 > li {
    counter-increment: lst-ctn-kix_list_54-2;
  }

  ol.lst-kix_list_24-0.start {
    counter-reset: lst-ctn-kix_list_24-0 0;
  }

  .lst-kix_list_18-3 > li {
    counter-increment: lst-ctn-kix_list_18-3;
  }

  .lst-kix_list_37-3 > li:before {
    content: "" counter(lst-ctn-kix_list_37-0, decimal) "."
      counter(lst-ctn-kix_list_37-1, decimal) "."
      counter(lst-ctn-kix_list_37-2, decimal) "."
      counter(lst-ctn-kix_list_37-3, decimal) " ";
  }

  ol.lst-kix_list_35-7.start {
    counter-reset: lst-ctn-kix_list_35-7 0;
  }

  .lst-kix_list_37-3 > li {
    counter-increment: lst-ctn-kix_list_37-3;
  }

  .lst-kix_list_18-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_38-7 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) "."
      counter(lst-ctn-kix_list_38-3, decimal) "."
      counter(lst-ctn-kix_list_38-4, decimal) "."
      counter(lst-ctn-kix_list_38-5, decimal) "."
      counter(lst-ctn-kix_list_38-6, decimal) "."
      counter(lst-ctn-kix_list_38-7, decimal) " ";
  }

  .lst-kix_list_3-0 > li {
    counter-increment: lst-ctn-kix_list_3-0;
  }

  .lst-kix_list_18-2 > li:before {
    content: "" counter(lst-ctn-kix_list_18-2, decimal) ". ";
  }

  ol.lst-kix_list_27-0.start {
    counter-reset: lst-ctn-kix_list_27-0 0;
  }

  ol.lst-kix_list_54-4.start {
    counter-reset: lst-ctn-kix_list_54-4 0;
  }

  ol.lst-kix_list_11-7.start {
    counter-reset: lst-ctn-kix_list_11-7 0;
  }

  ol.lst-kix_list_14-2.start {
    counter-reset: lst-ctn-kix_list_14-2 0;
  }

  .lst-kix_list_41-4 > li {
    counter-increment: lst-ctn-kix_list_41-4;
  }

  .lst-kix_list_30-4 > li {
    counter-increment: lst-ctn-kix_list_30-4;
  }

  .lst-kix_list_38-5 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) "."
      counter(lst-ctn-kix_list_38-3, decimal) "."
      counter(lst-ctn-kix_list_38-4, decimal) "."
      counter(lst-ctn-kix_list_38-5, decimal) " ";
  }

  .lst-kix_list_25-2 > li {
    counter-increment: lst-ctn-kix_list_25-2;
  }

  .lst-kix_list_27-1 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) ". ";
  }

  .lst-kix_list_48-2 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_27-3 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) "."
      counter(lst-ctn-kix_list_27-3, decimal) ". ";
  }

  .lst-kix_list_48-4 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_18-8 > li:before {
    content: "" counter(lst-ctn-kix_list_18-8, decimal) ". ";
  }

  ol.lst-kix_list_27-3.start {
    counter-reset: lst-ctn-kix_list_27-3 0;
  }

  .lst-kix_list_17-6 > li {
    counter-increment: lst-ctn-kix_list_17-6;
  }

  ol.lst-kix_list_4-3.start {
    counter-reset: lst-ctn-kix_list_4-3 0;
  }

  .lst-kix_list_10-7 > li:before {
    content: "" counter(lst-ctn-kix_list_10-7, decimal) " ";
  }

  .lst-kix_list_20-1 > li {
    counter-increment: lst-ctn-kix_list_20-1;
  }

  .lst-kix_list_10-5 > li:before {
    content: "" counter(lst-ctn-kix_list_10-5, decimal) " ";
  }

  ol.lst-kix_list_18-6.start {
    counter-reset: lst-ctn-kix_list_18-6 0;
  }

  .lst-kix_list_29-3 > li {
    counter-increment: lst-ctn-kix_list_29-3;
  }

  ol.lst-kix_list_54-2.start {
    counter-reset: lst-ctn-kix_list_54-2 0;
  }

  .lst-kix_list_9-2 > li:before {
    content: "" counter(lst-ctn-kix_list_9-0, decimal) "."
      counter(lst-ctn-kix_list_9-1, decimal) "."
      counter(lst-ctn-kix_list_9-2, decimal) " ";
  }

  .lst-kix_list_46-4 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_37-6.start {
    counter-reset: lst-ctn-kix_list_37-6 0;
  }

  ol.lst-kix_list_14-0.start {
    counter-reset: lst-ctn-kix_list_14-0 0;
  }

  .lst-kix_list_12-5 > li {
    counter-increment: lst-ctn-kix_list_12-5;
  }

  .lst-kix_list_31-1 > li {
    counter-increment: lst-ctn-kix_list_31-1;
  }

  .lst-kix_list_9-0 > li:before {
    content: "" counter(lst-ctn-kix_list_9-0, decimal) " ";
  }

  ol.lst-kix_list_24-3.start {
    counter-reset: lst-ctn-kix_list_24-3 0;
  }

  .lst-kix_list_23-5 > li {
    counter-increment: lst-ctn-kix_list_23-5;
  }

  .lst-kix_list_11-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_11-3, lower-latin) ") ";
  }

  .lst-kix_list_6-3 > li {
    counter-increment: lst-ctn-kix_list_6-3;
  }

  ol.lst-kix_list_1-3.start {
    counter-reset: lst-ctn-kix_list_1-3 0;
  }

  .lst-kix_list_29-1 > li:before {
    content: "Section " counter(lst-ctn-kix_list_29-0, upper-roman) "."
      counter(lst-ctn-kix_list_29-1, decimal) " ";
  }

  .lst-kix_list_29-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_29-3, lower-roman) ") ";
  }

  .lst-kix_list_49-8 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_9-8 > li:before {
    content: "" counter(lst-ctn-kix_list_9-8, decimal) " ";
  }

  .lst-kix_list_28-6 > li {
    counter-increment: lst-ctn-kix_list_28-6;
  }

  .lst-kix_list_28-7 > li:before {
    content: "" counter(lst-ctn-kix_list_28-7, lower-latin) ". ";
  }

  .lst-kix_list_1-7 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) "."
      counter(lst-ctn-kix_list_1-3, decimal) "."
      counter(lst-ctn-kix_list_1-4, decimal) "."
      counter(lst-ctn-kix_list_1-5, decimal) "."
      counter(lst-ctn-kix_list_1-6, decimal) "."
      counter(lst-ctn-kix_list_1-7, decimal) " ";
  }

  .lst-kix_list_50-1 > li {
    counter-increment: lst-ctn-kix_list_50-1;
  }

  .lst-kix_list_49-6 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_1-5 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) "."
      counter(lst-ctn-kix_list_1-3, decimal) "."
      counter(lst-ctn-kix_list_1-4, decimal) "."
      counter(lst-ctn-kix_list_1-5, decimal) " ";
  }

  .lst-kix_list_28-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_28-5, lower-roman) ") ";
  }

  .lst-kix_list_5-6 > li {
    counter-increment: lst-ctn-kix_list_5-6;
  }

  ol.lst-kix_list_22-6.start {
    counter-reset: lst-ctn-kix_list_22-6 0;
  }

  .lst-kix_list_2-1 > li:before {
    content: "" counter(lst-ctn-kix_list_2-1, decimal) ". ";
  }

  .lst-kix_list_49-0 > li:before {
    content: "\0025a0  ";
  }

  .lst-kix_list_2-3 > li:before {
    content: "" counter(lst-ctn-kix_list_2-3, decimal) ". ";
  }

  .lst-kix_list_11-8 > li {
    counter-increment: lst-ctn-kix_list_11-8;
  }

  .lst-kix_list_35-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_35-3, lower-latin) ") ";
  }

  ol.lst-kix_list_24-8.start {
    counter-reset: lst-ctn-kix_list_24-8 0;
  }

  .lst-kix_list_30-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_30-4, lower-roman) ") ";
  }

  .lst-kix_list_20-8 > li {
    counter-increment: lst-ctn-kix_list_20-8;
  }

  .lst-kix_list_9-1 > li {
    counter-increment: lst-ctn-kix_list_9-1;
  }

  .lst-kix_list_26-7 > li:before {
    content: "" counter(lst-ctn-kix_list_26-7, lower-latin) ". ";
  }

  ol.lst-kix_list_1-8.start {
    counter-reset: lst-ctn-kix_list_1-8 0;
  }

  .lst-kix_list_6-0 > li {
    counter-increment: lst-ctn-kix_list_6-0;
  }

  .lst-kix_list_3-5 > li:before {
    content: "" counter(lst-ctn-kix_list_3-5, decimal) ". ";
  }

  .lst-kix_list_40-2 > li {
    counter-increment: lst-ctn-kix_list_40-2;
  }

  ol.lst-kix_list_30-8.start {
    counter-reset: lst-ctn-kix_list_30-8 0;
  }

  ol.lst-kix_list_11-5.start {
    counter-reset: lst-ctn-kix_list_11-5 0;
  }

  .lst-kix_list_8-6 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_11-1 > li {
    counter-increment: lst-ctn-kix_list_11-1;
  }

  ol.lst-kix_list_16-6.start {
    counter-reset: lst-ctn-kix_list_16-6 0;
  }

  .lst-kix_list_44-4 > li {
    counter-increment: lst-ctn-kix_list_44-4;
  }

  ol.lst-kix_list_41-4.start {
    counter-reset: lst-ctn-kix_list_41-4 0;
  }

  ol.lst-kix_list_22-1.start {
    counter-reset: lst-ctn-kix_list_22-1 0;
  }

  .lst-kix_list_53-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_33-4 > li {
    counter-increment: lst-ctn-kix_list_33-4;
  }

  .lst-kix_list_21-3 > li:before {
    content: "" counter(lst-ctn-kix_list_21-3, decimal) ". ";
  }

  ol.lst-kix_list_30-2.start {
    counter-reset: lst-ctn-kix_list_30-2 0;
  }

  .lst-kix_list_4-4 > li {
    counter-increment: lst-ctn-kix_list_4-4;
  }

  ol.lst-kix_list_29-4.start {
    counter-reset: lst-ctn-kix_list_29-4 0;
  }

  .lst-kix_list_45-8 > li:before {
    content: "" counter(lst-ctn-kix_list_45-8, decimal) " ";
  }

  .lst-kix_list_31-8 > li {
    counter-increment: lst-ctn-kix_list_31-8;
  }

  .lst-kix_list_25-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_25-3, decimal) ") ";
  }

  ol.lst-kix_list_4-1.start {
    counter-reset: lst-ctn-kix_list_4-1 0;
  }

  .lst-kix_list_39-2 > li {
    counter-increment: lst-ctn-kix_list_39-2;
  }

  .lst-kix_list_16-2 > li:before {
    content: "" counter(lst-ctn-kix_list_16-2, decimal) ". ";
  }

  .lst-kix_list_26-6 > li {
    counter-increment: lst-ctn-kix_list_26-6;
  }

  .lst-kix_list_15-3 > li {
    counter-increment: lst-ctn-kix_list_15-3;
  }

  .lst-kix_list_44-4 > li:before {
    content: "" counter(lst-ctn-kix_list_44-4, decimal) " ";
  }

  ol.lst-kix_list_30-3.start {
    counter-reset: lst-ctn-kix_list_30-3 0;
  }

  .lst-kix_list_37-6 > li {
    counter-increment: lst-ctn-kix_list_37-6;
  }

  ol.lst-kix_list_11-0.start {
    counter-reset: lst-ctn-kix_list_11-0 0;
  }

  ol.lst-kix_list_18-8.start {
    counter-reset: lst-ctn-kix_list_18-8 0;
  }

  .lst-kix_list_3-3 > li {
    counter-increment: lst-ctn-kix_list_3-3;
  }

  .lst-kix_list_39-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_39-3, lower-latin) ") ";
  }

  ol.lst-kix_list_6-3.start {
    counter-reset: lst-ctn-kix_list_6-3 0;
  }

  .lst-kix_list_45-0 > li:before {
    content: "" counter(lst-ctn-kix_list_45-0, decimal) " ";
  }

  .lst-kix_list_17-6 > li:before {
    content: "" counter(lst-ctn-kix_list_17-6, decimal) ". ";
  }

  ol.lst-kix_list_16-2.start {
    counter-reset: lst-ctn-kix_list_16-2 0;
  }

  .lst-kix_list_42-8 > li {
    counter-increment: lst-ctn-kix_list_42-8;
  }

  .lst-kix_list_43-0 > li:before {
    content: "" counter(lst-ctn-kix_list_43-0, decimal) " ";
  }

  .lst-kix_list_43-8 > li:before {
    content: "" counter(lst-ctn-kix_list_43-8, decimal) " ";
  }

  .lst-kix_list_19-5 > li {
    counter-increment: lst-ctn-kix_list_19-5;
  }

  .lst-kix_list_28-2 > li {
    counter-increment: lst-ctn-kix_list_28-2;
  }

  .lst-kix_list_22-7 > li:before {
    content: "" counter(lst-ctn-kix_list_22-0, decimal) "."
      counter(lst-ctn-kix_list_22-1, decimal) "."
      counter(lst-ctn-kix_list_22-2, decimal) "."
      counter(lst-ctn-kix_list_22-3, decimal) "."
      counter(lst-ctn-kix_list_22-4, decimal) "."
      counter(lst-ctn-kix_list_22-5, decimal) "."
      counter(lst-ctn-kix_list_22-6, decimal) "."
      counter(lst-ctn-kix_list_22-7, decimal) ". ";
  }

  ol.lst-kix_list_35-2.start {
    counter-reset: lst-ctn-kix_list_35-2 0;
  }

  .lst-kix_list_34-7 > li:before {
    content: "" counter(lst-ctn-kix_list_34-7, decimal) " ";
  }

  .lst-kix_list_55-8 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_6-7 > li {
    counter-increment: lst-ctn-kix_list_6-7;
  }

  .lst-kix_list_26-3 > li {
    counter-increment: lst-ctn-kix_list_26-3;
  }

  .lst-kix_list_15-6 > li:before {
    content: "" counter(lst-ctn-kix_list_15-6, lower-roman) ") ";
  }

  .lst-kix_list_11-4 > li {
    counter-increment: lst-ctn-kix_list_11-4;
  }

  .lst-kix_list_22-4 > li {
    counter-increment: lst-ctn-kix_list_22-4;
  }

  .lst-kix_list_52-4 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_36-7 > li:before {
    content: "" counter(lst-ctn-kix_list_36-7, lower-latin) ". ";
  }

  ol.lst-kix_list_6-8.start {
    counter-reset: lst-ctn-kix_list_6-8 0;
  }

  .lst-kix_list_13-7 > li {
    counter-increment: lst-ctn-kix_list_13-7;
  }

  .lst-kix_list_20-7 > li:before {
    content: "" counter(lst-ctn-kix_list_20-7, decimal) ". ";
  }

  ol.lst-kix_list_6-5.start {
    counter-reset: lst-ctn-kix_list_6-5 0;
  }

  .lst-kix_list_41-8 > li:before {
    content: "" counter(lst-ctn-kix_list_41-8, decimal) " ";
  }

  ol.lst-kix_list_29-5.start {
    counter-reset: lst-ctn-kix_list_29-5 0;
  }

  .lst-kix_list_51-2 > li {
    counter-increment: lst-ctn-kix_list_51-2;
  }

  .lst-kix_list_41-0 > li:before {
    content: "" counter(lst-ctn-kix_list_41-0, decimal) " ";
  }

  ol.lst-kix_list_6-7.start {
    counter-reset: lst-ctn-kix_list_6-7 0;
  }

  .lst-kix_list_24-0 > li {
    counter-increment: lst-ctn-kix_list_24-0;
  }

  .lst-kix_list_54-4 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) "."
      counter(lst-ctn-kix_list_54-3, decimal) "."
      counter(lst-ctn-kix_list_54-4, decimal) " ";
  }

  .lst-kix_list_33-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_33-3, lower-latin) ") ";
  }

  .lst-kix_list_53-8 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_55-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_16-1.start {
    counter-reset: lst-ctn-kix_list_16-1 0;
  }

  ol.lst-kix_list_35-4.start {
    counter-reset: lst-ctn-kix_list_35-4 0;
  }

  ol.lst-kix_list_54-7.start {
    counter-reset: lst-ctn-kix_list_54-7 0;
  }

  ol.lst-kix_list_29-7.start {
    counter-reset: lst-ctn-kix_list_29-7 0;
  }

  .lst-kix_list_35-0 > li {
    counter-increment: lst-ctn-kix_list_35-0;
  }

  .lst-kix_list_40-4 > li:before {
    content: "" counter(lst-ctn-kix_list_40-4, decimal) " ";
  }

  ol.lst-kix_list_30-7.start {
    counter-reset: lst-ctn-kix_list_30-7 0;
  }

  .lst-kix_list_51-5 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) "."
      counter(lst-ctn-kix_list_51-3, decimal) "."
      counter(lst-ctn-kix_list_51-4, decimal) "."
      counter(lst-ctn-kix_list_51-5, decimal) " ";
  }

  .lst-kix_list_41-2 > li {
    counter-increment: lst-ctn-kix_list_41-2;
  }

  .lst-kix_list_14-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_14-3, decimal) ") ";
  }

  .lst-kix_list_14-0 > li:before {
    content: "" counter(lst-ctn-kix_list_14-0, decimal) ") ";
  }

  ol.lst-kix_list_18-5.start {
    counter-reset: lst-ctn-kix_list_18-5 0;
  }

  .lst-kix_list_6-1 > li {
    counter-increment: lst-ctn-kix_list_6-1;
  }

  .lst-kix_list_51-2 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) " ";
  }

  .lst-kix_list_51-3 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) "."
      counter(lst-ctn-kix_list_51-3, decimal) " ";
  }

  .lst-kix_list_14-6 > li:before {
    content: "" counter(lst-ctn-kix_list_14-6, decimal) ". ";
  }

  .lst-kix_ek190fjkpfvz-8 > li:before {
    content: "\0025a0  ";
  }

  .lst-kix_list_9-0 > li {
    counter-increment: lst-ctn-kix_list_9-0;
  }

  ol.lst-kix_list_25-3.start {
    counter-reset: lst-ctn-kix_list_25-3 0;
  }

  .lst-kix_list_51-6 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) "."
      counter(lst-ctn-kix_list_51-1, decimal) "."
      counter(lst-ctn-kix_list_51-2, decimal) "."
      counter(lst-ctn-kix_list_51-3, decimal) "."
      counter(lst-ctn-kix_list_51-4, decimal) "."
      counter(lst-ctn-kix_list_51-5, decimal) "."
      counter(lst-ctn-kix_list_51-6, decimal) " ";
  }

  ol.lst-kix_list_11-8.start {
    counter-reset: lst-ctn-kix_list_11-8 0;
  }

  .lst-kix_list_14-2 > li:before {
    content: "" counter(lst-ctn-kix_list_14-2, lower-roman) ") ";
  }

  .lst-kix_ek190fjkpfvz-1 > li:before {
    content: "\0025cb  ";
  }

  .lst-kix_ek190fjkpfvz-2 > li:before {
    content: "\0025a0  ";
  }

  ol.lst-kix_list_12-0.start {
    counter-reset: lst-ctn-kix_list_12-0 0;
  }

  .lst-kix_ek190fjkpfvz-5 > li:before {
    content: "\0025a0  ";
  }

  ol.lst-kix_list_41-7.start {
    counter-reset: lst-ctn-kix_list_41-7 0;
  }

  .lst-kix_list_32-2 > li:before {
    content: "" counter(lst-ctn-kix_list_32-0, decimal) "."
      counter(lst-ctn-kix_list_32-1, decimal) "."
      counter(lst-ctn-kix_list_32-2, decimal) " ";
  }

  .lst-kix_list_32-1 > li:before {
    content: "" counter(lst-ctn-kix_list_32-1, decimal) ". ";
  }

  ol.lst-kix_list_21-6.start {
    counter-reset: lst-ctn-kix_list_21-6 0;
  }

  ol.lst-kix_list_3-7.start {
    counter-reset: lst-ctn-kix_list_3-7 0;
  }

  .lst-kix_list_28-7 > li {
    counter-increment: lst-ctn-kix_list_28-7;
  }

  .lst-kix_list_31-7 > li {
    counter-increment: lst-ctn-kix_list_31-7;
  }

  .lst-kix_list_5-0 > li:before {
    content: "" counter(lst-ctn-kix_list_5-0, decimal) ". ";
  }

  .lst-kix_list_14-8 > li {
    counter-increment: lst-ctn-kix_list_14-8;
  }

  .lst-kix_list_5-3 > li:before {
    content: "" counter(lst-ctn-kix_list_5-3, decimal) ". ";
  }

  .lst-kix_list_36-8 > li {
    counter-increment: lst-ctn-kix_list_36-8;
  }

  .lst-kix_list_5-2 > li:before {
    content: "" counter(lst-ctn-kix_list_5-2, decimal) ". ";
  }

  .lst-kix_list_5-6 > li:before {
    content: "" counter(lst-ctn-kix_list_5-6, decimal) ". ";
  }

  ol.lst-kix_list_12-5.start {
    counter-reset: lst-ctn-kix_list_12-5 0;
  }

  .lst-kix_list_50-0 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) " ";
  }

  .lst-kix_list_50-3 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) "."
      counter(lst-ctn-kix_list_50-3, decimal) " ";
  }

  .lst-kix_list_50-4 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) "."
      counter(lst-ctn-kix_list_50-3, decimal) "."
      counter(lst-ctn-kix_list_50-4, decimal) " ";
  }

  .lst-kix_list_6-2 > li:before {
    content: "" counter(lst-ctn-kix_list_6-2, decimal) ". ";
  }

  ol.lst-kix_list_35-1.start {
    counter-reset: lst-ctn-kix_list_35-1 0;
  }

  ol.lst-kix_list_3-2.start {
    counter-reset: lst-ctn-kix_list_3-2 0;
  }

  .lst-kix_list_50-7 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) "."
      counter(lst-ctn-kix_list_50-3, decimal) "."
      counter(lst-ctn-kix_list_50-4, decimal) "."
      counter(lst-ctn-kix_list_50-5, decimal) "."
      counter(lst-ctn-kix_list_50-6, decimal) "."
      counter(lst-ctn-kix_list_50-7, decimal) " ";
  }

  .lst-kix_list_6-8 > li:before {
    content: "" counter(lst-ctn-kix_list_6-8, decimal) ". ";
  }

  .lst-kix_list_51-0 > li:before {
    content: "" counter(lst-ctn-kix_list_51-0, decimal) " ";
  }

  .lst-kix_list_50-6 > li:before {
    content: "" counter(lst-ctn-kix_list_50-0, decimal) "."
      counter(lst-ctn-kix_list_50-1, decimal) "."
      counter(lst-ctn-kix_list_50-2, decimal) "."
      counter(lst-ctn-kix_list_50-3, decimal) "."
      counter(lst-ctn-kix_list_50-4, decimal) "."
      counter(lst-ctn-kix_list_50-5, decimal) "."
      counter(lst-ctn-kix_list_50-6, decimal) " ";
  }

  .lst-kix_list_6-5 > li:before {
    content: "" counter(lst-ctn-kix_list_6-5, decimal) ". ";
  }

  ol.lst-kix_list_42-4.start {
    counter-reset: lst-ctn-kix_list_42-4 0;
  }

  .lst-kix_list_7-4 > li:before {
    content: "o  ";
  }

  .lst-kix_list_22-2 > li {
    counter-increment: lst-ctn-kix_list_22-2;
  }

  .lst-kix_list_44-8 > li {
    counter-increment: lst-ctn-kix_list_44-8;
  }

  .lst-kix_list_52-1 > li:before {
    content: "o  ";
  }

  .lst-kix_list_24-3 > li {
    counter-increment: lst-ctn-kix_list_24-3;
  }

  .lst-kix_list_13-7 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) "."
      counter(lst-ctn-kix_list_13-3, decimal) "."
      counter(lst-ctn-kix_list_13-4, decimal) "."
      counter(lst-ctn-kix_list_13-5, decimal) "."
      counter(lst-ctn-kix_list_13-6, decimal) "."
      counter(lst-ctn-kix_list_13-7, decimal) ". ";
  }

  ol.lst-kix_list_34-4.start {
    counter-reset: lst-ctn-kix_list_34-4 0;
  }

  .lst-kix_list_7-8 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_15-6 > li {
    counter-increment: lst-ctn-kix_list_15-6;
  }

  .lst-kix_list_4-7 > li {
    counter-increment: lst-ctn-kix_list_4-7;
  }

  .lst-kix_list_51-5 > li {
    counter-increment: lst-ctn-kix_list_51-5;
  }

  ol.lst-kix_list_2-5.start {
    counter-reset: lst-ctn-kix_list_2-5 0;
  }

  .lst-kix_list_15-5 > li:before {
    content: "" counter(lst-ctn-kix_list_15-5, lower-latin) ") ";
  }

  .lst-kix_list_31-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_31-6, decimal) ") ";
  }

  .lst-kix_list_31-8 > li:before {
    content: "" counter(lst-ctn-kix_list_31-8, decimal) " ";
  }

  .lst-kix_list_45-6 > li {
    counter-increment: lst-ctn-kix_list_45-6;
  }

  ol.lst-kix_list_26-0.start {
    counter-reset: lst-ctn-kix_list_26-0 0;
  }

  .lst-kix_list_52-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_27-3 > li {
    counter-increment: lst-ctn-kix_list_27-3;
  }

  .lst-kix_list_54-3 > li {
    counter-increment: lst-ctn-kix_list_54-3;
  }

  .lst-kix_list_4-1 > li:before {
    content: "" counter(lst-ctn-kix_list_4-1, decimal) ". ";
  }

  .lst-kix_list_31-2 > li:before {
    content: "" counter(lst-ctn-kix_list_31-0, decimal) "."
      counter(lst-ctn-kix_list_31-1, decimal) "."
      counter(lst-ctn-kix_list_31-2, decimal) " ";
  }

  .lst-kix_list_36-0 > li {
    counter-increment: lst-ctn-kix_list_36-0;
  }

  ol.lst-kix_list_33-3.start {
    counter-reset: lst-ctn-kix_list_33-3 0;
  }

  .lst-kix_list_4-3 > li:before {
    content: "" counter(lst-ctn-kix_list_4-3, decimal) ". ";
  }

  .lst-kix_list_1-8 > li {
    counter-increment: lst-ctn-kix_list_1-8;
  }

  ol.lst-kix_list_38-1.start {
    counter-reset: lst-ctn-kix_list_38-1 0;
  }

  .lst-kix_list_10-5 > li {
    counter-increment: lst-ctn-kix_list_10-5;
  }

  .lst-kix_list_24-4 > li {
    counter-increment: lst-ctn-kix_list_24-4;
  }

  .lst-kix_list_33-1 > li {
    counter-increment: lst-ctn-kix_list_33-1;
  }

  .lst-kix_list_16-2 > li {
    counter-increment: lst-ctn-kix_list_16-2;
  }

  ol.lst-kix_list_39-3.start {
    counter-reset: lst-ctn-kix_list_39-3 0;
  }

  ol.lst-kix_list_16-7.start {
    counter-reset: lst-ctn-kix_list_16-7 0;
  }

  .lst-kix_list_32-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_32-4, lower-roman) ") ";
  }

  .lst-kix_list_26-5 > li {
    counter-increment: lst-ctn-kix_list_26-5;
  }

  .lst-kix_list_19-2 > li {
    counter-increment: lst-ctn-kix_list_19-2;
  }

  .lst-kix_list_33-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_33-4, lower-roman) ") ";
  }

  .lst-kix_list_30-3 > li {
    counter-increment: lst-ctn-kix_list_30-3;
  }

  .lst-kix_list_38-2 > li {
    counter-increment: lst-ctn-kix_list_38-2;
  }

  ol.lst-kix_list_43-6.start {
    counter-reset: lst-ctn-kix_list_43-6 0;
  }

  .lst-kix_list_45-5 > li {
    counter-increment: lst-ctn-kix_list_45-5;
  }

  .lst-kix_list_12-1 > li:before {
    content: "(" counter(lst-ctn-kix_list_12-1, lower-latin) ") ";
  }

  .lst-kix_list_33-0 > li:before {
    content: "" counter(lst-ctn-kix_list_33-0, decimal) " ";
  }

  .lst-kix_list_53-7 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_32-8 > li:before {
    content: "" counter(lst-ctn-kix_list_32-8, decimal) " ";
  }

  .lst-kix_list_23-6 > li {
    counter-increment: lst-ctn-kix_list_23-6;
  }

  ol.lst-kix_list_13-6.start {
    counter-reset: lst-ctn-kix_list_13-6 0;
  }

  ol.lst-kix_list_25-8.start {
    counter-reset: lst-ctn-kix_list_25-8 0;
  }

  .lst-kix_list_39-0 > li {
    counter-increment: lst-ctn-kix_list_39-0;
  }

  .lst-kix_list_34-0 > li:before {
    content: "" counter(lst-ctn-kix_list_34-0, decimal) " ";
  }

  .lst-kix_list_21-4 > li {
    counter-increment: lst-ctn-kix_list_21-4;
  }

  .lst-kix_list_39-6 > li {
    counter-increment: lst-ctn-kix_list_39-6;
  }

  .lst-kix_list_13-3 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) "."
      counter(lst-ctn-kix_list_13-3, decimal) ". ";
  }

  .lst-kix_list_42-6 > li {
    counter-increment: lst-ctn-kix_list_42-6;
  }

  ol.lst-kix_list_43-5.start {
    counter-reset: lst-ctn-kix_list_43-5 0;
  }

  .lst-kix_list_34-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_34-6, decimal) ") ";
  }

  .lst-kix_list_12-5 > li:before {
    content: "" counter(lst-ctn-kix_list_12-5, decimal) " ";
  }

  ol.lst-kix_list_13-7.start {
    counter-reset: lst-ctn-kix_list_13-7 0;
  }

  .lst-kix_list_43-4 > li {
    counter-increment: lst-ctn-kix_list_43-4;
  }

  .lst-kix_list_18-4 > li {
    counter-increment: lst-ctn-kix_list_18-4;
  }

  .lst-kix_list_42-0 > li {
    counter-increment: lst-ctn-kix_list_42-0;
  }

  .lst-kix_list_12-7 > li:before {
    content: "" counter(lst-ctn-kix_list_12-7, decimal) " ";
  }

  .lst-kix_list_50-7 > li {
    counter-increment: lst-ctn-kix_list_50-7;
  }

  ol.lst-kix_list_51-0.start {
    counter-reset: lst-ctn-kix_list_51-0 0;
  }

  ol.lst-kix_list_21-1.start {
    counter-reset: lst-ctn-kix_list_21-1 0;
  }

  ol.lst-kix_list_50-8.start {
    counter-reset: lst-ctn-kix_list_50-8 0;
  }

  .lst-kix_list_25-1 > li {
    counter-increment: lst-ctn-kix_list_25-1;
  }

  .lst-kix_list_13-1 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) ". ";
  }

  .lst-kix_list_32-5 > li {
    counter-increment: lst-ctn-kix_list_32-5;
  }

  .lst-kix_list_22-8 > li {
    counter-increment: lst-ctn-kix_list_22-8;
  }

  .lst-kix_list_35-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_35-5, upper-latin) ") ";
  }

  ol.lst-kix_list_33-2.start {
    counter-reset: lst-ctn-kix_list_33-2 0;
  }

  ol.lst-kix_list_2-6.start {
    counter-reset: lst-ctn-kix_list_2-6 0;
  }

  ol.lst-kix_list_20-5.start {
    counter-reset: lst-ctn-kix_list_20-5 0;
  }

  ol.lst-kix_list_13-1.start {
    counter-reset: lst-ctn-kix_list_13-1 0;
  }

  .lst-kix_list_3-4 > li:before {
    content: "" counter(lst-ctn-kix_list_3-4, decimal) ". ";
  }

  .lst-kix_list_8-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_30-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_30-6, decimal) ") ";
  }

  .lst-kix_list_8-7 > li:before {
    content: "o  ";
  }

  .lst-kix_list_3-8 > li:before {
    content: "" counter(lst-ctn-kix_list_3-8, decimal) ". ";
  }

  .lst-kix_list_8-3 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_43-0.start {
    counter-reset: lst-ctn-kix_list_43-0 0;
  }

  .lst-kix_list_3-7 > li:before {
    content: "" counter(lst-ctn-kix_list_3-7, decimal) ". ";
  }

  .lst-kix_list_8-4 > li:before {
    content: "o  ";
  }

  .lst-kix_list_19-1 > li {
    counter-increment: lst-ctn-kix_list_19-1;
  }

  ol.lst-kix_list_26-4.start {
    counter-reset: lst-ctn-kix_list_26-4 0;
  }

  .lst-kix_list_35-8 > li:before {
    content: "" counter(lst-ctn-kix_list_35-8, decimal) " ";
  }

  ol.lst-kix_list_50-3.start {
    counter-reset: lst-ctn-kix_list_50-3 0;
  }

  .lst-kix_list_37-4 > li {
    counter-increment: lst-ctn-kix_list_37-4;
  }

  .lst-kix_list_16-8 > li:before {
    content: "" counter(lst-ctn-kix_list_16-8, decimal) ". ";
  }

  .lst-kix_list_16-7 > li:before {
    content: "" counter(lst-ctn-kix_list_16-7, decimal) ". ";
  }

  .lst-kix_list_17-8 > li {
    counter-increment: lst-ctn-kix_list_17-8;
  }

  .lst-kix_list_50-8 > li {
    counter-increment: lst-ctn-kix_list_50-8;
  }

  .lst-kix_list_4-7 > li:before {
    content: "" counter(lst-ctn-kix_list_4-7, decimal) ". ";
  }

  .lst-kix_list_17-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_39-8 > li {
    counter-increment: lst-ctn-kix_list_39-8;
  }

  .lst-kix_list_16-4 > li:before {
    content: "" counter(lst-ctn-kix_list_16-4, decimal) ". ";
  }

  ol.lst-kix_list_3-3.start {
    counter-reset: lst-ctn-kix_list_3-3 0;
  }

  .lst-kix_list_11-3 > li {
    counter-increment: lst-ctn-kix_list_11-3;
  }

  ol.lst-kix_list_50-4.start {
    counter-reset: lst-ctn-kix_list_50-4 0;
  }

  ol.lst-kix_list_18-1.start {
    counter-reset: lst-ctn-kix_list_18-1 0;
  }

  .lst-kix_list_30-1 > li {
    counter-increment: lst-ctn-kix_list_30-1;
  }

  ol.lst-kix_list_39-7.start {
    counter-reset: lst-ctn-kix_list_39-7 0;
  }

  .lst-kix_list_17-7 > li:before {
    content: "" counter(lst-ctn-kix_list_17-7, decimal) ". ";
  }

  ol.lst-kix_list_38-2.start {
    counter-reset: lst-ctn-kix_list_38-2 0;
  }

  .lst-kix_list_33-0 > li {
    counter-increment: lst-ctn-kix_list_33-0;
  }

  .lst-kix_list_17-3 > li:before {
    content: "" counter(lst-ctn-kix_list_17-3, decimal) ". ";
  }

  .lst-kix_list_17-4 > li:before {
    content: "" counter(lst-ctn-kix_list_17-4, decimal) ". ";
  }

  .lst-kix_list_2-8 > li:before {
    content: "" counter(lst-ctn-kix_list_2-8, decimal) ". ";
  }

  .lst-kix_list_10-0 > li:before {
    content: "" counter(lst-ctn-kix_list_10-0, decimal) ". ";
  }

  ol.lst-kix_list_21-7.start {
    counter-reset: lst-ctn-kix_list_21-7 0;
  }

  ol.lst-kix_list_43-1.start {
    counter-reset: lst-ctn-kix_list_43-1 0;
  }

  .lst-kix_list_18-3 > li:before {
    content: "" counter(lst-ctn-kix_list_18-3, decimal) ". ";
  }

  .lst-kix_list_18-6 > li {
    counter-increment: lst-ctn-kix_list_18-6;
  }

  ol.lst-kix_list_3-8.start {
    counter-reset: lst-ctn-kix_list_3-8 0;
  }

  ol.lst-kix_list_39-8.start {
    counter-reset: lst-ctn-kix_list_39-8 0;
  }

  .lst-kix_list_7-7 > li:before {
    content: "o  ";
  }

  .lst-kix_list_36-5 > li:before {
    content: "" counter(lst-ctn-kix_list_36-5, lower-roman) ". ";
  }

  .lst-kix_list_31-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_31-5, upper-latin) ") ";
  }

  ol.lst-kix_list_20-4.start {
    counter-reset: lst-ctn-kix_list_20-4 0;
  }

  .lst-kix_list_4-0 > li:before {
    content: "" counter(lst-ctn-kix_list_4-0, decimal) ". ";
  }

  .lst-kix_list_36-1 > li:before {
    content: "" counter(lst-ctn-kix_list_36-1, lower-latin) ". ";
  }

  .lst-kix_list_52-2 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_25-2.start {
    counter-reset: lst-ctn-kix_list_25-2 0;
  }

  .lst-kix_list_15-8 > li:before {
    content: "" counter(lst-ctn-kix_list_15-8, lower-roman) ". ";
  }

  ol.lst-kix_list_38-7.start {
    counter-reset: lst-ctn-kix_list_38-7 0;
  }

  .lst-kix_list_15-7 > li {
    counter-increment: lst-ctn-kix_list_15-7;
  }

  .lst-kix_list_4-4 > li:before {
    content: "" counter(lst-ctn-kix_list_4-4, decimal) ". ";
  }

  .lst-kix_list_9-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_9-3, lower-latin) ") ";
  }

  ol.lst-kix_list_13-2.start {
    counter-reset: lst-ctn-kix_list_13-2 0;
  }

  ol.lst-kix_list_33-7.start {
    counter-reset: lst-ctn-kix_list_33-7 0;
  }

  .lst-kix_list_29-8 > li:before {
    content: "" counter(lst-ctn-kix_list_29-8, lower-roman) ". ";
  }

  ol.lst-kix_list_3-6.start {
    counter-reset: lst-ctn-kix_list_3-6 0;
  }

  ol.lst-kix_list_51-5.start {
    counter-reset: lst-ctn-kix_list_51-5 0;
  }

  .lst-kix_list_32-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_32-5, upper-latin) ") ";
  }

  .lst-kix_list_12-4 > li:before {
    content: "" counter(lst-ctn-kix_list_12-4, decimal) " ";
  }

  .lst-kix_list_5-3 > li {
    counter-increment: lst-ctn-kix_list_5-3;
  }

  .lst-kix_list_33-1 > li:before {
    content: "" counter(lst-ctn-kix_list_33-1, decimal) ". ";
  }

  .lst-kix_list_1-0 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) " ";
  }

  ol.lst-kix_list_38-5.start {
    counter-reset: lst-ctn-kix_list_38-5 0;
  }

  .lst-kix_list_11-8 > li:before {
    content: "" counter(lst-ctn-kix_list_11-8, decimal) " ";
  }

  ol.lst-kix_list_2-0.start {
    counter-reset: lst-ctn-kix_list_2-0 0;
  }

  .lst-kix_list_49-5 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_13-0 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) ". ";
  }

  .lst-kix_list_13-4 > li:before {
    content: "" counter(lst-ctn-kix_list_13-0, decimal) "."
      counter(lst-ctn-kix_list_13-1, decimal) "."
      counter(lst-ctn-kix_list_13-2, decimal) "."
      counter(lst-ctn-kix_list_13-3, decimal) "."
      counter(lst-ctn-kix_list_13-4, decimal) ". ";
  }

  ol.lst-kix_list_26-5.start {
    counter-reset: lst-ctn-kix_list_26-5 0;
  }

  .lst-kix_list_55-6 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_38-6.start {
    counter-reset: lst-ctn-kix_list_38-6 0;
  }

  .lst-kix_list_54-6 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) "."
      counter(lst-ctn-kix_list_54-3, decimal) "."
      counter(lst-ctn-kix_list_54-4, decimal) "."
      counter(lst-ctn-kix_list_54-5, decimal) "."
      counter(lst-ctn-kix_list_54-6, decimal) " ";
  }

  ol.lst-kix_list_2-1.start {
    counter-reset: lst-ctn-kix_list_2-1 0;
  }

  ol.lst-kix_list_51-6.start {
    counter-reset: lst-ctn-kix_list_51-6 0;
  }

  .lst-kix_list_4-5 > li {
    counter-increment: lst-ctn-kix_list_4-5;
  }

  ol.lst-kix_list_33-8.start {
    counter-reset: lst-ctn-kix_list_33-8 0;
  }

  .lst-kix_list_12-8 > li:before {
    content: "" counter(lst-ctn-kix_list_12-8, decimal) " ";
  }

  .lst-kix_list_55-2 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_17-7.start {
    counter-reset: lst-ctn-kix_list_17-7 0;
  }

  ol.lst-kix_list_12-6.start {
    counter-reset: lst-ctn-kix_list_12-6 0;
  }

  .lst-kix_list_23-8 > li {
    counter-increment: lst-ctn-kix_list_23-8;
  }

  .lst-kix_list_35-5 > li {
    counter-increment: lst-ctn-kix_list_35-5;
  }

  ol.lst-kix_list_26-1.start {
    counter-reset: lst-ctn-kix_list_26-1 0;
  }

  .lst-kix_list_19-2 > li:before {
    content: "" counter(lst-ctn-kix_list_19-2, decimal) ". ";
  }

  ol.lst-kix_list_3-1.start {
    counter-reset: lst-ctn-kix_list_3-1 0;
  }

  ol.lst-kix_list_21-0.start {
    counter-reset: lst-ctn-kix_list_21-0 0;
  }

  .lst-kix_list_47-7 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_42-5.start {
    counter-reset: lst-ctn-kix_list_42-5 0;
  }

  .lst-kix_list_45-0 > li {
    counter-increment: lst-ctn-kix_list_45-0;
  }

  ol.lst-kix_list_51-4.start {
    counter-reset: lst-ctn-kix_list_51-4 0;
  }

  .lst-kix_list_36-6 > li {
    counter-increment: lst-ctn-kix_list_36-6;
  }

  .lst-kix_list_2-3 > li {
    counter-increment: lst-ctn-kix_list_2-3;
  }

  .lst-kix_list_47-4 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_1-2 > li {
    counter-increment: lst-ctn-kix_list_1-2;
  }

  .lst-kix_list_19-8 > li:before {
    content: "" counter(lst-ctn-kix_list_19-8, decimal) ". ";
  }

  ol.lst-kix_list_20-8.start {
    counter-reset: lst-ctn-kix_list_20-8 0;
  }

  .lst-kix_list_19-5 > li:before {
    content: "" counter(lst-ctn-kix_list_19-5, decimal) ". ";
  }

  .lst-kix_list_47-1 > li:before {
    content: "o  ";
  }

  ol.lst-kix_list_34-8.start {
    counter-reset: lst-ctn-kix_list_34-8 0;
  }

  .lst-kix_list_50-2 > li {
    counter-increment: lst-ctn-kix_list_50-2;
  }

  .lst-kix_list_37-7 > li:before {
    content: "" counter(lst-ctn-kix_list_37-7, lower-latin) ". ";
  }

  ol.lst-kix_list_17-2.start {
    counter-reset: lst-ctn-kix_list_17-2 0;
  }

  .lst-kix_list_13-2 > li {
    counter-increment: lst-ctn-kix_list_13-2;
  }

  ol.lst-kix_list_21-5.start {
    counter-reset: lst-ctn-kix_list_21-5 0;
  }

  .lst-kix_list_38-0 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) " ";
  }

  .lst-kix_list_19-7 > li {
    counter-increment: lst-ctn-kix_list_19-7;
  }

  .lst-kix_list_14-3 > li {
    counter-increment: lst-ctn-kix_list_14-3;
  }

  .lst-kix_list_37-4 > li:before {
    content: "(" counter(lst-ctn-kix_list_37-4, lower-roman) ") ";
  }

  .lst-kix_list_12-1 > li {
    counter-increment: lst-ctn-kix_list_12-1;
  }

  .lst-kix_list_51-3 > li {
    counter-increment: lst-ctn-kix_list_51-3;
  }

  .lst-kix_list_33-3 > li {
    counter-increment: lst-ctn-kix_list_33-3;
  }

  ol.lst-kix_list_25-4.start {
    counter-reset: lst-ctn-kix_list_25-4 0;
  }

  .lst-kix_list_38-3 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) "."
      counter(lst-ctn-kix_list_38-3, decimal) " ";
  }

  .lst-kix_list_38-6 > li:before {
    content: "" counter(lst-ctn-kix_list_38-0, decimal) "."
      counter(lst-ctn-kix_list_38-1, decimal) "."
      counter(lst-ctn-kix_list_38-2, decimal) "."
      counter(lst-ctn-kix_list_38-3, decimal) "."
      counter(lst-ctn-kix_list_38-4, decimal) "."
      counter(lst-ctn-kix_list_38-5, decimal) "."
      counter(lst-ctn-kix_list_38-6, decimal) " ";
  }

  ol.lst-kix_list_34-3.start {
    counter-reset: lst-ctn-kix_list_34-3 0;
  }

  ol.lst-kix_list_2-4.start {
    counter-reset: lst-ctn-kix_list_2-4 0;
  }

  .lst-kix_list_34-4 > li {
    counter-increment: lst-ctn-kix_list_34-4;
  }

  .lst-kix_list_41-8 > li {
    counter-increment: lst-ctn-kix_list_41-8;
  }

  .lst-kix_list_48-0 > li:before {
    content: "\0025cf  ";
  }

  ol.lst-kix_list_50-5.start {
    counter-reset: lst-ctn-kix_list_50-5 0;
  }

  .lst-kix_list_2-5 > li:before {
    content: "" counter(lst-ctn-kix_list_2-5, decimal) ". ";
  }

  .lst-kix_list_48-6 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_18-6 > li:before {
    content: "" counter(lst-ctn-kix_list_18-6, decimal) ". ";
  }

  .lst-kix_list_39-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_39-5, upper-latin) ") ";
  }

  .lst-kix_list_14-6 > li {
    counter-increment: lst-ctn-kix_list_14-6;
  }

  .lst-kix_list_54-5 > li {
    counter-increment: lst-ctn-kix_list_54-5;
  }

  ol.lst-kix_list_42-3.start {
    counter-reset: lst-ctn-kix_list_42-3 0;
  }

  ol.lst-kix_list_39-2.start {
    counter-reset: lst-ctn-kix_list_39-2 0;
  }

  .lst-kix_list_10-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_10-3, upper-latin) ") ";
  }

  .lst-kix_list_15-4 > li {
    counter-increment: lst-ctn-kix_list_15-4;
  }

  .lst-kix_list_2-6 > li {
    counter-increment: lst-ctn-kix_list_2-6;
  }

  .lst-kix_list_36-8 > li:before {
    content: "" counter(lst-ctn-kix_list_36-8, lower-roman) ". ";
  }

  .lst-kix_list_28-1 > li {
    counter-increment: lst-ctn-kix_list_28-1;
  }

  .lst-kix_list_46-8 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_3-4 > li {
    counter-increment: lst-ctn-kix_list_3-4;
  }

  .lst-kix_list_29-5 > li:before {
    content: "" counter(lst-ctn-kix_list_29-5, lower-latin) ") ";
  }

  .lst-kix_list_20-0 > li:before {
    content: "\0025cf  ";
  }

  .lst-kix_list_9-6 > li:before {
    content: "(" counter(lst-ctn-kix_list_9-6, decimal) ") ";
  }

  ol.lst-kix_list_34-5.start {
    counter-reset: lst-ctn-kix_list_34-5 0;
  }

  ol.lst-kix_list_51-1.start {
    counter-reset: lst-ctn-kix_list_51-1 0;
  }

  .lst-kix_list_20-6 > li:before {
    content: "" counter(lst-ctn-kix_list_20-6, decimal) ". ";
  }

  .lst-kix_list_23-0 > li {
    counter-increment: lst-ctn-kix_list_23-0;
  }

  ol.lst-kix_list_12-1.start {
    counter-reset: lst-ctn-kix_list_12-1 0;
  }

  .lst-kix_list_11-5 > li:before {
    content: "(" counter(lst-ctn-kix_list_11-5, upper-latin) ") ";
  }

  ol.lst-kix_list_50-7.start {
    counter-reset: lst-ctn-kix_list_50-7 0;
  }

  .lst-kix_list_54-3 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) "."
      counter(lst-ctn-kix_list_54-1, decimal) "."
      counter(lst-ctn-kix_list_54-2, decimal) "."
      counter(lst-ctn-kix_list_54-3, decimal) " ";
  }

  ol.lst-kix_list_21-2.start {
    counter-reset: lst-ctn-kix_list_21-2 0;
  }

  .lst-kix_list_20-6 > li {
    counter-increment: lst-ctn-kix_list_20-6;
  }

  ol.lst-kix_list_25-7.start {
    counter-reset: lst-ctn-kix_list_25-7 0;
  }

  .lst-kix_list_1-3 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) "."
      counter(lst-ctn-kix_list_1-3, decimal) " ";
  }

  ol.lst-kix_list_34-6.start {
    counter-reset: lst-ctn-kix_list_34-6 0;
  }

  .lst-kix_list_28-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_28-3, decimal) ") ";
  }

  ol.lst-kix_list_42-0.start {
    counter-reset: lst-ctn-kix_list_42-0 0;
  }

  ol.lst-kix_list_2-7.start {
    counter-reset: lst-ctn-kix_list_2-7 0;
  }

  .lst-kix_list_27-7 > li:before {
    content: "" counter(lst-ctn-kix_list_27-0, decimal) "."
      counter(lst-ctn-kix_list_27-1, decimal) "."
      counter(lst-ctn-kix_list_27-2, decimal) "."
      counter(lst-ctn-kix_list_27-3, decimal) "."
      counter(lst-ctn-kix_list_27-4, decimal) "."
      counter(lst-ctn-kix_list_27-5, decimal) "."
      counter(lst-ctn-kix_list_27-6, decimal) "."
      counter(lst-ctn-kix_list_27-7, decimal) ". ";
  }

  .lst-kix_list_25-7 > li {
    counter-increment: lst-ctn-kix_list_25-7;
  }

  ol.lst-kix_list_39-4.start {
    counter-reset: lst-ctn-kix_list_39-4 0;
  }

  .lst-kix_list_49-2 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_35-2 > li:before {
    content: "" counter(lst-ctn-kix_list_35-0, decimal) "."
      counter(lst-ctn-kix_list_35-1, decimal) "."
      counter(lst-ctn-kix_list_35-2, decimal) " ";
  }

  .lst-kix_list_3-1 > li {
    counter-increment: lst-ctn-kix_list_3-1;
  }

  .lst-kix_list_30-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_30-3, lower-latin) ") ";
  }

  ol.lst-kix_list_34-0.start {
    counter-reset: lst-ctn-kix_list_34-0 0;
  }

  .lst-kix_list_26-8 > li {
    counter-increment: lst-ctn-kix_list_26-8;
  }

  ol.lst-kix_list_39-1.start {
    counter-reset: lst-ctn-kix_list_39-1 0;
  }

  .lst-kix_list_3-1 > li:before {
    content: "" counter(lst-ctn-kix_list_3-1, decimal) ". ";
  }

  .lst-kix_list_14-0 > li {
    counter-increment: lst-ctn-kix_list_14-0;
  }

  .lst-kix_list_44-0 > li:before {
    content: "" counter(lst-ctn-kix_list_44-0, decimal) " ";
  }

  ol.lst-kix_list_17-4.start {
    counter-reset: lst-ctn-kix_list_17-4 0;
  }

  .lst-kix_list_33-6 > li {
    counter-increment: lst-ctn-kix_list_33-6;
  }

  ol.lst-kix_list_12-3.start {
    counter-reset: lst-ctn-kix_list_12-3 0;
  }

  .lst-kix_list_44-2 > li {
    counter-increment: lst-ctn-kix_list_44-2;
  }

  .lst-kix_list_21-2 > li:before {
    content: "" counter(lst-ctn-kix_list_21-2, decimal) ". ";
  }

  .lst-kix_list_2-0 > li {
    counter-increment: lst-ctn-kix_list_2-0;
  }

  .lst-kix_list_15-1 > li {
    counter-increment: lst-ctn-kix_list_15-1;
  }

  .lst-kix_list_36-3 > li {
    counter-increment: lst-ctn-kix_list_36-3;
  }

  .lst-kix_list_11-2 > li:before {
    content: "" counter(lst-ctn-kix_list_11-0, decimal) "."
      counter(lst-ctn-kix_list_11-1, decimal) "."
      counter(lst-ctn-kix_list_11-2, decimal) " ";
  }

  ol.lst-kix_list_42-2.start {
    counter-reset: lst-ctn-kix_list_42-2 0;
  }

  .lst-kix_list_53-1 > li:before {
    content: "o  ";
  }

  .lst-kix_list_40-4 > li {
    counter-increment: lst-ctn-kix_list_40-4;
  }

  ol.lst-kix_list_12-4.start {
    counter-reset: lst-ctn-kix_list_12-4 0;
  }

  .lst-kix_list_35-2 > li {
    counter-increment: lst-ctn-kix_list_35-2;
  }

  .lst-kix_list_25-2 > li:before {
    content: "" counter(lst-ctn-kix_list_25-2, lower-roman) ") ";
  }

  .lst-kix_list_29-5 > li {
    counter-increment: lst-ctn-kix_list_29-5;
  }

  .lst-kix_list_16-1 > li:before {
    content: "" counter(lst-ctn-kix_list_16-1, decimal) ". ";
  }

  ol.lst-kix_list_39-0.start {
    counter-reset: lst-ctn-kix_list_39-0 0;
  }

  .lst-kix_list_54-8 > li {
    counter-increment: lst-ctn-kix_list_54-8;
  }

  .lst-kix_list_12-4 > li {
    counter-increment: lst-ctn-kix_list_12-4;
  }

  .lst-kix_list_44-3 > li:before {
    content: "" counter(lst-ctn-kix_list_44-3, decimal) " ";
  }

  .lst-kix_list_39-2 > li:before {
    content: "" counter(lst-ctn-kix_list_39-0, decimal) "."
      counter(lst-ctn-kix_list_39-1, decimal) "."
      counter(lst-ctn-kix_list_39-2, decimal) " ";
  }

  .lst-kix_list_12-7 > li {
    counter-increment: lst-ctn-kix_list_12-7;
  }

  .lst-kix_list_34-7 > li {
    counter-increment: lst-ctn-kix_list_34-7;
  }

  .lst-kix_list_30-0 > li:before {
    content: "" counter(lst-ctn-kix_list_30-0, decimal) " ";
  }

  ol.lst-kix_list_17-3.start {
    counter-reset: lst-ctn-kix_list_17-3 0;
  }

  .lst-kix_list_51-0 > li {
    counter-increment: lst-ctn-kix_list_51-0;
  }

  .lst-kix_list_43-4 > li:before {
    content: "" counter(lst-ctn-kix_list_43-4, decimal) " ";
  }

  .lst-kix_list_7-1 > li:before {
    content: "o  ";
  }

  .lst-kix_list_13-5 > li {
    counter-increment: lst-ctn-kix_list_13-5;
  }

  .lst-kix_list_48-3 > li:before {
    content: "\0025aa  ";
  }

  .lst-kix_list_9-6 > li {
    counter-increment: lst-ctn-kix_list_9-6;
  }

  .lst-kix_list_29-2 > li {
    counter-increment: lst-ctn-kix_list_29-2;
  }

  ol.lst-kix_list_42-6.start {
    counter-reset: lst-ctn-kix_list_42-6 0;
  }

  .lst-kix_list_20-3 > li {
    counter-increment: lst-ctn-kix_list_20-3;
  }

  .lst-kix_list_11-6 > li {
    counter-increment: lst-ctn-kix_list_11-6;
  }

  .lst-kix_list_52-8 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_17-8.start {
    counter-reset: lst-ctn-kix_list_17-8 0;
  }

  ol.lst-kix_list_34-1.start {
    counter-reset: lst-ctn-kix_list_34-1 0;
  }

  .lst-kix_list_1-5 > li {
    counter-increment: lst-ctn-kix_list_1-5;
  }

  .lst-kix_list_28-4 > li {
    counter-increment: lst-ctn-kix_list_28-4;
  }

  .lst-kix_list_37-1 > li {
    counter-increment: lst-ctn-kix_list_37-1;
  }

  ol.lst-kix_list_17-5.start {
    counter-reset: lst-ctn-kix_list_17-5 0;
  }

  .lst-kix_list_15-2 > li:before {
    content: "(" counter(lst-ctn-kix_list_15-2, lower-latin) ") ";
  }

  .lst-kix_list_10-8 > li {
    counter-increment: lst-ctn-kix_list_10-8;
  }

  .lst-kix_list_10-6 > li:before {
    content: "" counter(lst-ctn-kix_list_10-6, decimal) " ";
  }

  .lst-kix_list_40-7 > li {
    counter-increment: lst-ctn-kix_list_40-7;
  }

  ol.lst-kix_list_12-7.start {
    counter-reset: lst-ctn-kix_list_12-7 0;
  }

  .lst-kix_list_41-5 > li {
    counter-increment: lst-ctn-kix_list_41-5;
  }

  .lst-kix_list_53-4 > li:before {
    content: "\0025aa  ";
  }

  ol.lst-kix_list_42-8.start {
    counter-reset: lst-ctn-kix_list_42-8 0;
  }

  ol.lst-kix_list_12-8.start {
    counter-reset: lst-ctn-kix_list_12-8 0;
  }

  .lst-kix_list_20-3 > li:before {
    content: "" counter(lst-ctn-kix_list_20-3, decimal) ". ";
  }

  .lst-kix_list_54-0 > li:before {
    content: "" counter(lst-ctn-kix_list_54-0, decimal) " ";
  }

  .lst-kix_list_27-6 > li {
    counter-increment: lst-ctn-kix_list_27-6;
  }

  .lst-kix_list_29-2 > li:before {
    content: "(" counter(lst-ctn-kix_list_29-2, lower-latin) ") ";
  }

  .lst-kix_list_21-1 > li {
    counter-increment: lst-ctn-kix_list_21-1;
  }

  ol.lst-kix_list_17-6.start {
    counter-reset: lst-ctn-kix_list_17-6 0;
  }

  .lst-kix_list_28-6 > li:before {
    content: "" counter(lst-ctn-kix_list_28-6, decimal) ". ";
  }

  .lst-kix_list_1-6 > li:before {
    content: "" counter(lst-ctn-kix_list_1-0, decimal) "."
      counter(lst-ctn-kix_list_1-1, decimal) "."
      counter(lst-ctn-kix_list_1-2, decimal) "."
      counter(lst-ctn-kix_list_1-3, decimal) "."
      counter(lst-ctn-kix_list_1-4, decimal) "."
      counter(lst-ctn-kix_list_1-5, decimal) "."
      counter(lst-ctn-kix_list_1-6, decimal) " ";
  }

  .lst-kix_list_42-3 > li {
    counter-increment: lst-ctn-kix_list_42-3;
  }

  .lst-kix_list_43-1 > li {
    counter-increment: lst-ctn-kix_list_43-1;
  }

  .lst-kix_list_33-7 > li:before {
    content: "" counter(lst-ctn-kix_list_33-7, decimal) " ";
  }

  ol.lst-kix_list_42-7.start {
    counter-reset: lst-ctn-kix_list_42-7 0;
  }

  .lst-kix_list_40-0 > li:before {
    content: "" counter(lst-ctn-kix_list_40-0, decimal) " ";
  }

  .lst-kix_list_34-3 > li:before {
    content: "(" counter(lst-ctn-kix_list_34-3, lower-latin) ") ";
  }

  .lst-kix_list_2-2 > li:before {
    content: "" counter(lst-ctn-kix_list_2-2, decimal) ". ";
  }
}
</style>

<style scoped>
.longer {
  width: 50%;
  height: 100%;
  max-width: none;
}
.breiter {
  height: 80%;
}
</style>
