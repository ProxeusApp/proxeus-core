<template>
  <nav
    ref="sidebar"
    class="sidebar px-0 pt-0"
    :class="{ toggled: toggled, menutoggled: menutoggled }"
  >
    <div ref="bgLayerContainer" style="display: none">
      <div
        ref="bgLayer"
        v-show="menutoggled"
        @click="menuToggle"
        class="sidebar-bg-layer"
      >
        <i class="material-icons sidebar-close-btn"> close </i>
      </div>
    </div>
    <slot></slot>
  </nav>
</template>

<script>
import "bootstrap";

import mafdc from "@/mixinApp";

export default {
  mixins: [mafdc],
  name: "responsive-sidebar",
  props: ["user", "toggled"],
  watch: {
    toggled(newValue) {
      this.toggleTooltips(newValue);
    },
  },
  data() {
    return {
      menutoggled: false,
      bgLayerOutside: false,
    };
  },
  mounted() {
    if (window) {
      window.addEventListener("resize", this.resizeListener, true);
    }
    this.$root.$on("menutoggler-click", this.menuToggle);
    this.$refs.sidebar.addEventListener("click", this.anyClick, false);

    this.resizeListener();
    this.toggleTooltips(this.toggled);
  },
  beforeDestroy() {
    if (window) {
      window.removeEventListener("resize", this.resizeListener);
    }
    this.$root.$off("menutoggler-click", this.menuToggle);
    this.$refs.sidebar.removeEventListener("click", this.anyClick);
    $(".tooltip.show").each(function () {
      $(this).remove();
    });
  },
  methods: {
    resizeListener() {
      if (window.innerWidth <= 660) {
        if (!this.bgLayerOutside) {
          this.open();
        }
      } else {
        if (this.bgLayerOutside) {
          this.close();
        }
      }
    },
    anyClick(ev) {
      if (this.menutoggled && ev.defaultPrevented) {
        this.menuToggle();
      }
    },
    menuToggle() {
      if (this.menutoggled) {
        this.menutoggled = false;
        if (this.bgLayerOutside) {
          this.close();
        }
      } else {
        this.menutoggled = true;
        if (!this.bgLayerOutside) {
          this.open();
        }
      }
    },
    toggleTooltips(toggled) {
      if (toggled) {
        $('[data-toggle="tooltip"]').tooltip();
      } else {
        $('[data-toggle="tooltip"]').tooltip("dispose");
      }
    },
    open() {
      if (this.$refs.bgLayer) {
        this.bgLayerOutside = true;
        document.body.appendChild(this.$refs.bgLayer);
      }
    },
    close() {
      if (this.$refs.bgLayerContainer) {
        this.bgLayerOutside = false;
        this.$refs.bgLayerContainer.appendChild(this.$refs.bgLayer);
      }
    },
    logout() {
      axios.post("/api/logout", null).then(
        (response) => {
          window.location.replace("/");
        },
        (err) => {
          this.app.handleError(err);
        }
      );
    },
  },
};
</script>
<style lang="scss">
@use "@/assets/styles/variables" as *;
@use "~bootstrap/scss/mixins";
@use "../../assets/styles/sidebar.scss";

.brand-name {
  letter-spacing: 2px;
}
</style>
