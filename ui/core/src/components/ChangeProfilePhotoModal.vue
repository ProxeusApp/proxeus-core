<template>
  <b-modal
    class="profilephoto-modal b-modal"
    v-model="show"
    :title="$t('Select profile photo')"
    :header-bg-variant="'light'"
    @hide="onDialogHide"
  >
    <change-profile-photo :setup="getCropFunc" @onFileDrop="onFileDrop" />
    <template slot="modal-footer">
      <button v-show="file" @click="onDialogOk" class="btn btn-primary">
        {{ $t("Set as profile photo") }}
      </button>
      <button @click="onDialogHide" class="btn btn-secondary">
        {{ $t("Cancel") }}
      </button>
    </template>
  </b-modal>
</template>

<script>
import { BModal, bModalDirective } from "bootstrap-vue";
import mafdc from "@/mixinApp";
import ChangeProfilePhoto from "./ChangeProfilePhoto";

export default {
  mixins: [mafdc],
  name: "change-profile-photo-modal",
  components: {
    ChangeProfilePhoto,
    "b-modal": bModal,
  },
  directives: {
    "b-modal": bModalDirective,
  },
  props: {
    setup: {
      type: Function,
      default: () => null,
    },
  },
  data() {
    return {
      show: false,
      cropFunc: null,
      file: null,
    };
  },
  created() {
    if (this.setup) {
      this.setup(this.openDialog);
    }
  },
  methods: {
    onFileDrop(f) {
      this.file = f;
    },
    getCropFunc(f) {
      this.cropFunc = f;
    },
    cropped(file, croppedBlob) {
      axios
        .post("/api/my/profile/photo", croppedBlob, {
          headers: {
            "File-Name": encodeURI(file.name),
            "Content-Type": croppedBlob.type,
          },
        })
        .then(
          (response) => {
            this.show = false;
            this.$root.$emit("profile-photo-update");
            this.$emit("onDialogOk");
          },
          (err) => {
            this.app.handleError(err);
            this.$notify({
              group: "app",
              title: this.$t("Error"),
              text: this.$t(
                "There was an error uploading the image. Please try again or if the error persists contact the platform operator."
              ),
              type: "error",
            });
          }
        );
    },
    openDialog() {
      this.show = true;
    },
    onDialogHide() {
      this.show = false;
      this.$emit("onDialogHide");
    },
    onDialogOk() {
      if (this.cropFunc) {
        this.cropFunc(this.cropped);
      }
    },
  },
};
</script>
