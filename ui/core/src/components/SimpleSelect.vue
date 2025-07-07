<template>
  <div :tabindex="disabled ? false : '0'" ref="myroot" class="ss-sel-main">
    <div @click="clickDropDown($event)" class="ss-select" :name="name">
      <table style="width: 100%; height: 100%">
        <tr>
          <td>
            <span class="ss-selected">
              <span
                ref="selected"
                class="sss"
                style="vertical-align: middle"
              ></span>
            </span>
          </td>
          <td class="min">
            <span class="ss-arrow" :class="{ selected: show }">
              <i class="mdi mdi-menu-down"></i>
            </span>
          </td>
        </tr>
      </table>
    </div>
    <div class="ss-list" v-show="shouldShowDropDown()">
      <div
        ref="unselect"
        v-if="unselect"
        @click="unselectClick"
        class="ss-unselect"
        v-html="unselectedHTML"
      ></div>
      <ul ref="items">
        <li
          tabindex="-1"
          v-for="(option, id) in options"
          :key="get(option, idProp) || id"
          class="cursor-pointer outline-none"
          @click.prevent="select($event, option, id)"
        >
          <slot :option="option">
            <span v-if="typeof option === 'string'">{{ option }}</span>
            <span v-else-if="get(option, labelProp)">{{
              get(option, labelProp)
            }}</span>
            <span v-else
              >your option value: {{ JSON.stringify(option) }}, please provide
              labelProp</span
            >
          </slot>
        </li>
      </ul>
    </div>
    <div class="errors"></div>
    <div
      v-show="disableWithVisibleLayer && disabled"
      class="ss-sel-disabled"
    ></div>
  </div>
</template>

<script>
export default {
  name: "simple-select",
  props: {
    name: {
      type: String,
      default: null,
    },
    options: {
      type: Array,
      default: () => [],
    },
    selected: {
      type: null,
      default: null,
    },
    value: {
      type: null,
      default: null,
    },
    disabled: {
      type: Boolean,
      default: () => false,
    },
    // possible to unselect after a selection?
    unselect: {
      type: Boolean,
      default: true,
    },
    disableWithVisibleLayer: {
      type: Boolean,
      default: () => true,
    },
    unselectedValue: {
      type: null,
      default: null,
    },
    // define the appearance of the default status of the select
    unselectedHTML: {
      type: String,
      default: "&#8226;&#8226;&#8226;",
    },
    // to access the id inside your option object
    idProp: {
      type: String,
      default: "id",
    },
    // whether v-model is set by idProp or by the object in options
    selectedValueIsIDProp: {
      type: Boolean,
      default: true,
    },
    // to access the label inside your option object
    labelProp: {
      type: String,
      default: "label",
    },
    onSelect: {
      type: Function,
      default: () => null,
    },
    onSelectedChange: {
      type: Function,
      default: () => null,
    },
  },
  data() {
    return {
      show: false,
      lastSelectedLi: null,
      lastTargetLi: null,
      selectedIndex: -1,
      selectedKey: null,
      lastOptionsLength: 0,
      renderNow: true,
      focused: false,
    };
  },
  updated() {
    if (this.renderNow) {
      this.renderSelected();
    }
  },
  beforeDestroy() {
    if (this.$refs.myroot) {
      this.$refs.myroot.removeEventListener("focus", this.focusHandler);
      this.$refs.myroot.removeEventListener("blur", this.blurHandler);
    }
  },
  created() {
    if (this.options && this.options.length > 0) {
      this.lastOptionsLength = this.options.length;
    }
  },
  mounted() {
    let selectedVal = "";
    if (typeof this.selected === "function") {
      selectedVal = this.selected();
    } else if (this.selected !== null) {
      selectedVal = this.selected;
    } else {
      selectedVal = this.value;
    }
    if (this.$refs.myroot) {
      this.$refs.myroot.addEventListener("focus", this.focusHandler);
      this.$refs.myroot.addEventListener("blur", this.blurHandler);
    }
    if (this.options && this.options.length > 0) {
      this.lastOptionsLength = this.options.length;
      // use id prop
      for (let i = 0; i < this.options.length; i++) {
        var maybeSelected = this.get(this.options[i], this.idProp);
        if (maybeSelected === this.value || maybeSelected === selectedVal) {
          this.selectedIndex = i;
          this.selectedKey = this.value || selectedVal;
          this.lastSelectedLi = this.lastTargetLi =
            this.$refs.items.children[i];
          this.setActive(this.lastSelectedLi);
          this.renderSelected();
          return;
        }
      }
      if (
        (selectedVal === null ||
          selectedVal === undefined ||
          selectedVal === "") &&
        !isNaN(this.value)
      ) {
        this.selectedIndex = this.value;
      } else if (!isNaN(selectedVal)) {
        this.selectedIndex = selectedVal;
      }
      // use array index
      if (
        this.selectedIndex !== null &&
        this.selectedIndex !== undefined &&
        !isNaN(this.selectedIndex) &&
        this.selectedIndex >= 0 &&
        this.selectedIndex < this.options.length
      ) {
        this.selectedKey = this.selectedIndex;
        this.renderSelected();
        return;
      }
    }
    this.insertSelectedHTMLStr(this.unselectedHTML);
  },
  methods: {
    focusHandler(e) {
      if (this.disabled) {
        return;
      }
      this.focused = true;
      if (this.$refs.myroot) {
        this.$refs.myroot.addEventListener("keydown", this.keyboardHandler);
      }
    },
    unselectClick() {
      if (this.show) {
        this.justHide();
      }
      this.selectedIndex = -1;
      this.selectedKey = null;
      this.updateVModel(null);
    },
    blurHandler(e) {
      this.hide(e);
      if (this.$refs.myroot) {
        this.$refs.myroot.removeEventListener("keydown", this.keyboardHandler);
      }
    },
    keyboardHandler(e) {
      if (this.disabled) {
        return;
      }
      if (e.which === 38 || e.keyCode === 38) {
        // up
        const li = this.getTargetLI();
        if (li) {
          if (li.previousSibling) {
            this.changeTargetLI(li.previousSibling);
            return false;
          } else if (this.unselect) {
            this.changeTargetLI(this.$refs.unselect);
            return false;
          }
        }
      } else if (e.which === 40 || e.keyCode === 40) {
        // down
        if (!this.show) {
          this.clickDropDown(e, true);
          return false;
        }
        let li = this.getTargetLI();
        if (li) {
          if (
            li === this.$refs.unselect &&
            this.$refs.items &&
            this.$refs.items.children &&
            this.$refs.items.children.length
          ) {
            li = this.$refs.items.children[0];
            this.changeTargetLI(li);
            return false;
          } else if (li.nextSibling) {
            this.changeTargetLI(li.nextSibling);
            return false;
          }
        }
      } else if (
        e.which === 9 || // tab
        e.keyCode === 9 ||
        e.which === 32 || // space
        e.keyCode === 32
      ) {
        if (
          e.which === 32 || // space
          e.keyCode === 32
        ) {
          if (!this.show) {
            this.clickDropDown(e, true);
            return false;
          }
        }
        // close
        this.justHide(e);
      } else if (e.which === 13 || e.keyCode === 13) {
        // enter
        // select
        const target = this.getTargetLI();
        if (target === this.$refs.unselect) {
          this.unselectClick();
          this.removeTargetClass(target);
        } else {
          const dummyEv = { target: target };
          const index = this.getIndexOfLI(dummyEv.target);
          if (index !== -1) {
            this.select(dummyEv, this.options[index], index);
          }
        }
      }
    },
    getIndexOfLI(li) {
      for (let i = 0; i < this.$refs.items.children.length; i++) {
        if (this.$refs.items.children[i] === li) {
          return i;
        }
      }
      return -1;
    },
    getTargetLI(noDefault) {
      if (this.lastTargetLi) {
        return this.lastTargetLi;
      }
      if (this.lastSelectedLi) {
        return this.lastSelectedLi;
      }
      if (
        !noDefault &&
        this.$refs.items &&
        this.$refs.items.children &&
        this.$refs.items.children.length > 0 &&
        this.$refs.items.children.length === this.options.length
      ) {
        this.selectedIndex = 0;
        this.selectedKey = this.get(
          this.options[this.selectedIndex],
          this.idProp
        );
        return this.$refs.items.children[0];
      }
      return null;
    },
    shouldShowDropDown() {
      return !this.disabled && this.show && this.options && this.options.length;
    },
    renderSelected() {
      let selectedValue = "";
      // update when selected function changed
      if (
        this.renderNow &&
        !this.show &&
        this.selected &&
        typeof this.selected === "function"
      ) {
        selectedValue = this.selected();
        if (selectedValue !== this.selectedKey) {
          if (this.options) {
            for (let i = 0; i < this.options.length; i++) {
              if (selectedValue === this.get(this.options[i], this.idProp)) {
                this.selectedIndex = i;
                this.selectedKey = selectedValue;
                this.lastSelectedLi = this.lastTargetLi =
                  this.$refs.items.children[i];
                this.selectItem(this.$refs.items.children[i]);
                return;
              }
            }
          }
        }
      }
      // updated
      if (
        this.selectedIndex > -1 &&
        this.$refs.items &&
        this.$refs.items.children &&
        this.$refs.items.children.length === this.options.length &&
        this.selectedIndex < this.options.length
      ) {
        if (this.lastOptionsLength !== this.options.length) {
          // check if the key can be found after the options changed
          let foundItAgainObj;
          for (let i = 0; i < this.$refs.items.children.length; i++) {
            this.removeActiveClass(this.$refs.items.children[i]);
            this.removeTargetClass(this.$refs.items.children[i]);
            if (this.selectedKey && this.options) {
              if (this.selectedKey === this.get(this.options[i], this.idProp)) {
                foundItAgainObj = this.options[i];
                this.selectedIndex = i;
                this.lastSelectedLi = this.lastTargetLi =
                  this.$refs.items.children[i];
                break;
              }
            }
          }
          if (!foundItAgainObj) {
            // key does not exist in the new options.. unselect
            this.resetAll();
          }
          this.lastOptionsLength = this.options.length;
          return;
        }
        // just update the selection with the li item
        this.selectItem(this.$refs.items.children[this.selectedIndex]);
        return;
      }
      this.lastOptionsLength = this.options.length;
      this.insertSelectedHTMLStr(this.unselectedHTML);
    },
    resetAll() {
      this.lastSelectedLi = null;
      this.lastTargetLi = null;
      this.selectedIndex = -1;
      this.selectedKey = null;
      this.$emit("input", this.unselectedValue);
      this.lastOptionsLength = this.options.length;
      this.insertSelectedHTMLStr(this.unselectedHTML);
    },
    get(obj, k) {
      try {
        return obj[k];
      } catch (e) {}
      return null;
    },
    hide(ev) {
      if (this.show) {
        if (ev && this.canHide(ev.target)) {
          this.justHide(ev);
          return true;
        }
      }
      return false;
    },
    justHide() {
      this.show = false;
      if (
        this.$refs.items &&
        this.$refs.items.children &&
        this.$refs.items.children.length
      ) {
        for (let i = 0; i < this.$refs.items.children.length; i++) {
          this.removeTargetClass(this.$refs.items.children[i]);
        }
      }
      document.removeEventListener("click", this.hide);
      this.renderNow = true;
    },
    removeTargetClass(t) {
      if (t) {
        t.className = t.className.replace(/\blitarget\b/g, "");
      }
    },
    canHide(n) {
      if (n) {
        //while (true) {
        if (n === this.$refs.myroot) {
          return false;
        }
        if (n === undefined || n === document.body) {
          return true;
        }
        n = n.parentNode;
        //}
      }
      return false;
    },
    clickDropDown(ev, noKeydownHandler) {
      if (this.show) {
        this.justHide(ev);
      } else {
        // to set the target class
        this.changeTargetLI(this.getTargetLI(true));
        if (this.$refs.myroot) {
          this.$refs.myroot.focus();
        }
        if (!noKeydownHandler && this.$refs.myroot) {
          this.$refs.myroot.addEventListener("keydown", this.keyboardHandler);
        }
        this.renderNow = false;
        this.show = true;
        document.addEventListener("click", this.hide, true);
      }
    },
    select(ev, item, id) {
      if (this.disabled) {
        this.renderNow = true;
        return;
      }
      if (ev && ev.target) {
        const t = this.getDPLi(ev.target);
        if (t && t.outerHTML) {
          this.selectedIndex = id;
          this.selectedKey = this.get(item, this.idProp);
          this.selectItem(t);
          if (this.onSelect) {
            this.onSelect(item, id);
          }
          this.updateVModel(item);
          this.justHide(ev);
        }
      }
      this.renderNow = true;
    },
    updateVModel(item) {
      if (this.selectedValueIsIDProp) {
        this.$emit("input", this.selectedKey);
      } else {
        this.$emit("input", item);
      }
    },
    changeTargetLI(li) {
      const t = this.getTargetLI();
      if (t !== li) {
        if (t) {
          t.className = t.className.replace(/\blitarget\b/g, "");
        }
        this.lastTargetLi = li;
      }
      try {
        const name = "litarget";
        const arr = li.className.split(" ");
        if (arr.indexOf(name) === -1) {
          li.className += " " + name;
        }
      } catch (e) {}
    },
    removeActiveClass(t) {
      if (t) {
        t.className = t.className.replace(/\bliactive\b/g, "");
      }
    },
    setActive(li) {
      try {
        const name = "liactive";
        const arr = li.className.split(" ");
        if (arr.indexOf(name) === -1) {
          li.className += " " + name;
        }
      } catch (e) {}
    },
    selectItem(li) {
      if (this.lastSelectedLi !== li) {
        if (this.lastSelectedLi) {
          this.removeActiveClass(this.lastSelectedLi);
          this.removeTargetClass(this.lastSelectedLi);
        }
        this.lastSelectedLi = li;
        this.lastTargetLi = li;
        this.setActive(li);
      }
      if (!li || !li.outerHTML) {
        return;
      }
      const n = li.outerHTML.replace(/<li[^>]+>/m, "").replace(/\<\/li\>$/, "");
      let newSelectedItem;
      const before = function (custom) {
        if (custom) {
          newSelectedItem = custom;
          return;
        }
        newSelectedItem = n;
      };
      if (this.onSelectedChange) {
        this.onSelectedChange(n, before);
      }
      if (!newSelectedItem) {
        newSelectedItem = n;
      }
      this.insertSelectedHTMLStr(newSelectedItem);
    },
    insertSelectedHTMLStr(htmlStr) {
      this.$refs.selected.innerHTML = htmlStr;
    },
    getDPLi(t) {
      if (t) {
        if (this.$refs.items === t.parentNode) {
          return t;
        }
        return this.getDPLi(t.parentNode);
      }
      return null;
    },
  },
};
</script>

<style lang="scss">
.ss-list > ul > li {
  white-space: nowrap;
}

.ss-sel-main .ss-sel-disabled {
  position: absolute;
  background: #ffffff7a;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  cursor: not-allowed;
}

.ss-sel-main {
  position: relative;
}

.ss-sel-main:focus {
  outline: none;
}

.ss-sel-main .ss-list li {
  border-bottom: 1px solid #8080803b;
  cursor: pointer;
  padding: 6px;
  overflow: hidden;
}

.ss-sel-main .ss-list li.liactive {
  border-bottom: 1px solid #8080803b;
  cursor: default;
  color: rgba(128, 128, 128, 0.51) !important;
}

.ss-sel-main .ss-list li:hover,
.ss-sel-main .ss-unselect:hover,
.ss-sel-main .ss-list .litarget {
  background: #0000000d;
}

.ss-sel-main .ss-list {
  position: absolute;
  right: 0;
  background: white;
  z-index: 10000;
  min-width: 100%;
  border: 1px solid #dadada;
  -webkit-box-shadow: 6px 6px 19px -5px rgba(0, 0, 0, 0.5);
  -moz-box-shadow: 6px 6px 19px -5px rgba(0, 0, 0, 0.5);
  box-shadow: 6px 6px 19px -5px rgba(0, 0, 0, 0.5);
}

.ss-sel-main .ss-list ul {
  list-style-type: none;
  padding: 4px;
  margin: 0;
}

.ss-select {
  display: inline-block;
  position: relative;
  padding: 2px 5px;
  padding-left: 12px;
  cursor: default;
  border: 2px solid #dee2e6;
  min-width: 100%;
  background: white;
}

.ss-sel-main:focus .ss-select {
  border: 2px solid #062a85;
}

.ss-select .ss-arrow {
  vertical-align: middle;
  display: inline-block;
  position: relative;
  font-size: 22px;
  line-height: 1;
}

.ss-select .ss-arrow.selected {
  color: #a9a9a9;
}

.ss-select .sss * {
  vertical-align: middle;
}

.ss-sel-main td.min {
  width: 1%;
  white-space: nowrap;
}

.ss-sel-main .ss-unselect {
  text-align: center;
  vertical-align: middle;
  display: inline-block;
  width: 100%;
  border-bottom: 1px solid #8080805c;
  cursor: pointer;
  padding-top: 4px;
}
</style>
