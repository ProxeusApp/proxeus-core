<template>
  <b-modal
    class="perm-modal b-modal"
    v-model="show"
    :title="$t('Share with others')"
    :ok-title="$t('Done')"
    :header-bg-variant="'light'"
    @hide="onDialogHide"
    @ok="onDialogOk"
  >
    <div class="small alert alert-warning" role="alert">
      {{
        $t(
          "share warning",
          "This feature is for advanced users. Make sure you understand how the sharing mechanisms work before you change anything.  Consult the handbook for more information."
        )
      }}
    </div>
    <div
      class="perm-change-wrapper"
      style="margin-bottom: 15px; position: relative"
    >
      <div class="sub-title">Owner</div>
      <div v-if="owner" class="perm-change-parent">
        <user-selector
          v-show="changeOwnerClicked"
          :excludes="getMapOfExistingUsersInPermItem"
          :disabled="isOwnerDisabled"
          @added="ownerChanged"
          :maxItems="1"
          :dependencyFulfilled="isGrantSelectorNotEmpty"
          v-model="granted"
          :addBtnText="$t('Change')"
          :uri="'/api/admin/user/list'"
        />
        <i v-if="loading" class="mdi mdi-loading mdi-spin"></i>
        <table class="perm-granted-tbl" v-else-if="owner" style="width: 100%">
          <tr>
            <td class="tdmin">
              <img
                v-if="owner.photo"
                class="mlist-group-icon"
                :src="owner.photo"
              />
              <i v-else class="mlist-group-icon material-icons">person</i>
            </td>
            <td class="tdmax impcnt" style="padding-left: 10px">
              <span>{{ owner.name }}</span>
              <small class="light-text easy-read" v-if="owner.email">{{
                owner.email
              }}</small>
              <small class="light-text easy-read" v-if="owner.etherPK">{{
                owner.etherPK
              }}</small>
              <small class="light-text easy-read" v-else-if="owner.detail">{{
                owner.detail
              }}</small>
              <small class="light-text easy-read" v-else>-</small>
            </td>
            <td class="tdmin" v-if="ownerEnabled">
              <small
                v-show="!changeOwnerClicked"
                class="mt-1 mr-1 perm-link-btn"
                @click="changeOwnerClicked = true"
                >{{ $t("inline change badge", "change") }}
              </small>
            </td>
            <td class="tdmin">
              <small class="light-text mt-1 mr-1">owner</small>
            </td>
          </tr>
        </table>
      </div>
      <div v-else class="py-1 light-text" style="font-size: 0.85em">
        {{ $t("This entity doesn't have an owner.") }}
      </div>
    </div>
    <div v-show="publicLink" style="position: relative; margin-bottom: 15px">
      <div class="perm-change-wrapper">
        <div class="sub-title d-flex align-items-center">
          <div>{{ $t("Share by link") }}</div>
        </div>
        <div class="perm-subtitle mb-1">
          <small>{{ $t("Share with everyone who has the link.") }}</small>
        </div>
        <div ref="publicInput" class="perm-change-parent">
          <table class="nicetbl">
            <tr v-show="publicLink2" class="share-lnk-row">
              <td>
                <button
                  type="button"
                  class="btn btn-default"
                  :class="{ focus: exeLink }"
                  @click="toggleExeLink"
                >
                  <i class="material-icons"> play_arrow </i>
                  {{ $t("execute link") }}
                </button>
              </td>
              <td>
                <button
                  type="button"
                  class="btn btn-default"
                  :class="{ focus: buildLink }"
                  @click="toggleBuildLink"
                >
                  <i class="material-icons"> build </i>
                  {{ $t("build link") }}
                </button>
              </td>
            </tr>
            <tr>
              <td colspan="2">
                <input
                  ref="publicLinkInput"
                  class="share-link form-control"
                  v-model="publicLinkReadOnly"
                  readonly
                  :disabled="publicReadWriteSelect === 0"
                  onclick="this.setSelectionRange(0, this.value.length)"
                  onfocus="this.setSelectionRange(0, this.value.length)"
                  type="text"
                  :placeholder="$t('Link...')"
                />
              </td>
              <td>
                <read-write-selector
                  :unselect="false"
                  class="perm-grant-selector"
                  :disableWithVisibleLayer="false"
                  :disabled="!publicEnabled"
                  :provideNone="true"
                  v-model="publicReadWriteSelect"
                  :onSelectedChange="OnSelectedShowIconOnly"
                />
              </td>
            </tr>
          </table>
          <div class="text-center">
            <button type="button" class="btn btn-default" @click="copyLink">
              Copy to clipboard
            </button>
          </div>
          <div
            v-show="!publicEnabled"
            class="perm-disabled"
            data-toggle="tooltip"
            data-placement="top"
            title="You don't have permissions to share by link."
          ></div>
        </div>
      </div>
      <div
        v-show="!publicEnabled"
        class="py-1 light-text"
        style="font-size: 0.85em"
      >
        {{ $t("You don't have permissions to share by link.") }}
      </div>
    </div>
    <div v-if="!advancedClicked" class="mt-1 mr-1" style="text-align: right">
      <span class="perm-link-btn" @click="advancedClicked = true">{{
        $t("advanced")
      }}</span>
    </div>
    <div v-else>
      <div
        v-show="groupEnabled"
        class="perm-change-wrapper"
        style="margin-bottom: 15px; position: relative"
      >
        <div class="sub-title d-flex align-items-center">
          <div>{{ $t("Group and others") }}</div>
          <div
            class="perm-link-btn ml-auto"
            @click="
              advancedClicked = false;
              changeGroupClicked = false;
            "
          >
            {{ $t("hide") }}
          </div>
        </div>
        <div v-if="changeGroupClicked" class="perm-subtitle">
          {{
            $t(
              "Group and others explanation",
              "Define write and read rights for a group or any other people. Other people are authorized users with a different group."
            )
          }}
        </div>
        <div v-if="changeGroupClicked" class="perm-change-parent">
          <table
            style="width: 100%; text-align: center; vertical-align: middle"
          >
            <tr style="text-align: left">
              <td>
                <table>
                  <tr>
                    <td class="tdmin" style="padding-right: 10px">
                      {{ $t("Group") }}
                    </td>
                    <td>
                      <simple-select
                        v-model="groupRole"
                        style="width: 120px"
                        class=""
                        :idProp="'role'"
                        :labelProp="'name'"
                        :options="app.roles"
                      />
                    </td>
                  </tr>
                </table>
              </td>
              <td>Others</td>
            </tr>
            <tr>
              <td style="width: 50%">
                <read-write-selector
                  :disabled="groupRole === null || groupRole === undefined"
                  v-model="groupRights"
                  :selected="groupRightsFunc"
                />
              </td>
              <td style="width: 50%">
                <read-write-selector v-model="othersRights" />
              </td>
            </tr>
          </table>
        </div>
        <div v-else class="mt-1 mr-1" style="text-align: center">
          <span class="perm-link-btn" @click="changeGroupClicked = true">{{
            $t("define")
          }}</span>
        </div>
      </div>
      <div style="position: relative">
        <div class="perm-change-wrapper" style="position: relative">
          <div class="sub-title">{{ $t("Grant people") }}</div>
          <div class="perm-subtitle mb-1">
            <small>
              {{
                $t(
                  "Grant people explanation",
                  "Grant access to a specific person by entering a name, email or blockchain address and selecting read or write rights."
                )
              }}
            </small>
          </div>
          <div ref="grantInput" class="perm-change-parent">
            <table class="nicetbl">
              <tr>
                <td>
                  <user-selector
                    :excludes="getMapOfExistingUsersInPermItem"
                    :disabled="isGrantDisabled"
                    @added="grantAdded"
                    :dependencyFulfilled="isGrantSelectorNotEmpty"
                    v-model="granted"
                    :uri="'/api/admin/user/list'"
                  />
                </td>
                <td class="tdmin">
                  <read-write-selector
                    :unselect="false"
                    class="perm-grant-selector"
                    :disableWithVisibleLayer="false"
                    :disabled="!grantEnabled"
                    :provideNone="false"
                    v-model="grantReadWriteSelect"
                    :onSelectedChange="OnSelectedShowIconOnly"
                  />
                </td>
              </tr>
            </table>
            <div
              v-show="!grantEnabled"
              class="perm-disabled"
              data-toggle="tooltip"
              data-placement="top"
              title="You don't have permissions to grant other people access."
            ></div>
          </div>
        </div>
        <div
          v-show="!grantEnabled"
          class="py-1 light-text"
          style="font-size: 0.85em"
        >
          {{
            $t("You don't have write permissions to grant other people access.")
          }}
        </div>
        <div
          class="perm-granted-list perm-change-parent"
          v-if="loading || (granted && granted.length > 0)"
        >
          <i v-if="loading" class="mdi mdi-loading mdi-spin"></i>
          <table
            class="nicetbl perm-granted-tbl"
            v-else-if="granted"
            style="width: 100%"
          >
            <tr
              v-for="(element, index) in granted"
              :key="element.id"
              :data-index="index"
            >
              <td
                class="tdmin"
                v-if="grantEnabled"
                @click="deleteGrant(element.id)"
              >
                <i class="material-icons perm-delete">clear</i>
              </td>
              <td class="tdmin">
                <img
                  v-if="element.photo"
                  class="mlist-group-icon"
                  :src="element.photo"
                />
                <i v-else class="mlist-group-icon material-icons">person</i>
              </td>
              <td class="tdmax impcnt">
                <span class="px-2 flex-col-truncate font-weight-bold">{{
                  element.name
                }}</span>
                <small
                  class="px-2 light-text mt-1 flex-col-truncate"
                  v-if="element.etherPK"
                  >{{ element.etherPK }}</small
                >
                <small
                  class="px-2 light-text mt-1 flex-col-truncate"
                  v-else-if="element.detail"
                  >{{ element.detail }}
                </small>
                <small class="px-2 light-text mt-1 flex-col-truncate" v-else>{{
                  $t("-")
                }}</small>
              </td>
              <td class="tdmin" v-if="hasGrantVal(element.id, 1)">
                <h5 class="mb-0 pb-0 pt-1 mr-1" style="text-align: center">
                  <i class="material-icons" style="text-align: center"
                    >visibility</i
                  >
                </h5>
                <small class="light-text mt-1 mr-1 d-inline-block">{{
                  $t("Can read")
                }}</small>
              </td>
              <td class="tdmin" v-else-if="hasGrantVal(element.id, 2)">
                <h5 class="mb-0 pb-0 pt-1 mr-1" style="text-align: center">
                  <i class="material-icons" style="text-align: center"
                    >create</i
                  >
                </h5>
                <small class="light-text mt-1 mr-1 d-inline-block">{{
                  $t("Can write")
                }}</small>
              </td>
            </tr>
          </table>
        </div>
      </div>
    </div>
    <i slot="modal-ok" class="material-icons"> save </i>
  </b-modal>
</template>

<script>
// import VueTagsInput from '@johmun/vue-tags-input'
import { BModal, VBModal } from "bootstrap-vue";
import SimpleSelect from "@/components/SimpleSelect";
import ReadWriteSelector from "./ReadWriteSelector";
// import UserItem from './UserItem'
import mafdc from "@/mixinApp";
import formChangeAlert from "../../../mixins/form-change-alert";
import UserSelector from "./UserSelector";

export default {
  mixins: [mafdc, formChangeAlert],
  name: "permission-dialog",
  components: {
    UserSelector,
    // UserItem,
    ReadWriteSelector,
    SimpleSelect,
    // 'vue-tags-input': VueTagsInput,
    "b-modal": BModal,
  },
  directives: {
    "b-modal": VBModal,
  },
  props: {
    setup: {
      type: Function,
      default: () => null,
    },
    publicLink: {
      type: String,
      default: "",
    },
    publicLink2: {
      type: String,
      default: "",
    },
    value: {
      type: Object,
      default: null,
    },
    save: {
      type: Function,
      default: function () {},
    },
  },
  data() {
    return {
      exeLinkStr: "",
      buildLinkStr: "",
      exeLink: false,
      buildLink: true,
      loaderVisible: true,
      groupRole: null,
      groupRights: null,
      othersRights: null,
      show: false,
      item: null,
      granted: null,
      loading: true,
      owner: null,
      publicReadWriteSelect: 0,
      grantReadWriteSelect: 1,
      ownerEnabled: false,
      groupEnabled: false,
      grantEnabled: false,
      publicEnabled: false,
      changeOwnerClicked: false,
      changeGroupClicked: false,
      advancedClicked: false,
      publicLinkReadOnly: "",
    };
  },
  watch: {
    groupRole: "updateGroupRightsSelector",
  },
  created() {
    if (this.setup) {
      this.setup(this.openDialog);
    }
    this.publicLinkReadOnly = this.publicLink;
    if (this.publicLink2) {
      this.exeLinkStr = this.publicLink2;
      this.buildLinkStr = this.publicLink;
    }
    this.reset();
  },
  methods: {
    toggleExeLink() {
      this.exeLink = true;
      this.buildLink = false;
      this.publicLinkReadOnly = this.exeLinkStr;
    },
    toggleBuildLink() {
      this.exeLink = false;
      this.buildLink = true;
      this.publicLinkReadOnly = this.buildLinkStr;
    },
    copyLink() {
      if (this.$refs.publicLinkInput) {
        if (this.$refs.publicLinkInput !== document.activeElement) {
          this.$refs.publicLinkInput.focus();
        }
        this.$refs.publicLinkInput.setSelectionRange(
          0,
          this.$refs.publicLinkInput.value.length
        );
        document.execCommand("copy");
        this.$notify({
          group: "app",
          title: this.$t("Success"),
          text: this.$t("copied to clipboard"),
          type: "success",
        });
      }
    },
    hasUnsavedChanges() {
      return false;
    },
    updateGroupRightsSelector() {
      if (this.groupRole === null || this.groupRole === undefined) {
        this.groupRights = null;
      }
    },
    groupRightsFunc() {
      return this.groupRights;
    },
    reset() {
      if (this.value) {
        this.item = {
          owner: this.value.owner,
          groupAndOthers: this.value.groupAndOthers,
          grant: this.value.grant,
          publicByID: this.value.publicByID,
        };
      } else {
        this.item = {
          owner: null,
          groupAndOthers: { group: null, rights: null },
          grant: null,
          publicByID: null,
        };
      }
      this.groupRole = null;
      this.snapshot(this.item);
      if (this.item.groupAndOthers) {
        this.groupRole = this.item.groupAndOthers.group;
        if (this.item.groupAndOthers.rights) {
          if (this.item.groupAndOthers.rights.length > 0) {
            this.groupRights = this.item.groupAndOthers.rights[0];
          }
          if (this.item.groupAndOthers.rights.length > 1) {
            this.othersRights = this.item.groupAndOthers.rights[1];
          }
        }
      }
      if (this.item.publicByID) {
        if (this.item.publicByID.length > 0) {
          this.publicReadWriteSelect = this.item.publicByID[0];
        }
      }
      this.changeOwnerClicked = false;
      this.changeGroupClicked = false;
      this.advancedClicked = false;
      this.ownerEnabled = false;
      this.groupEnabled = false;
      this.grantEnabled = false;
      this.publicEnabled = false;
      this.updateMe();
    },
    updateMe() {
      this.app.loadMe(this.updatedMe);
    },
    updatedMe(me) {
      if (me && this.item) {
        if (me.id === this.item.owner) {
          this.owner = me;
          this.ownerEnabled = true;
          this.groupEnabled = true;
          this.grantEnabled = true;
          this.publicEnabled = true;
        } else if (
          this.item.groupAndOthers &&
          ((this.item.groupAndOthers.group <= me.role &&
            this.item.groupAndOthers.rights &&
            this.item.groupAndOthers.rights.length > 0 &&
            this.item.groupAndOthers.rights[0] === 2) ||
            (this.item.groupAndOthers.rights &&
              this.item.groupAndOthers.rights.length > 1 &&
              this.item.groupAndOthers.rights[1] === 2))
        ) {
          this.publicEnabled = true;
          this.grantEnabled = true;
        } else if (
          this.item.grant &&
          this.item.grant[me.id] &&
          this.item.grant[me.id][0] === 2
        ) {
          this.publicEnabled = true;
          this.grantEnabled = true;
        }
        if (!this.publicEnabled || !this.grantEnabled) {
          if (this.app.amIWriteGrantedFor(this.item)) {
            this.publicEnabled = true;
            this.grantEnabled = true;
          }
        }
      }
    },
    hasGrantVal(id, val) {
      id = id + "";
      return (
        this.item.grant && this.item.grant[id] && this.item.grant[id][0] === val
      );
    },
    openDialog() {
      this.reset();
      this.show = true;
      this.refreshGrantList();
      if (this.publicLink2) {
        if (this.exeLink) {
          this.toggleExeLink();
        } else {
          this.toggleBuildLink();
        }
      }
    },
    onDialogHide() {
      this.show = false;
      this.$emit("onDialogHide");
    },
    onDialogOk() {
      this.show = false;
      if (!this.item.groupAndOthers) {
        this.item.groupAndOthers = {};
      }
      if (
        !this.item.groupAndOthers.rights ||
        this.item.groupAndOthers.rights.length !== 2
      ) {
        this.item.groupAndOthers.rights = [0, 0];
      }
      this.item.groupAndOthers.rights[0] = this.groupRights;
      this.item.groupAndOthers.rights[1] = this.othersRights;
      this.item.groupAndOthers.group = this.groupRole;
      if (
        !isNaN(this.publicReadWriteSelect) &&
        this.publicReadWriteSelect > 0
      ) {
        this.item.publicByID = [this.publicReadWriteSelect];
      } else {
        this.item.publicByID = null;
      }
      if (!this.compare(this.item)) {
        const merger = this.value;
        merger.owner = this.item.owner;
        merger.groupAndOthers = this.item.groupAndOthers;
        merger.grant = this.item.grant;
        merger.publicByID = this.item.publicByID;
        this.$emit("input", merger);
        this.save();
      }
      this.$emit("onDialogOk");
    },
    grantAdded(usrList) {
      if (!this.item.grant) {
        this.item.grant = {};
      }
      for (let i = 0; i < usrList.length; i++) {
        this.item.grant[usrList[i].id] = [this.grantReadWriteSelect];
      }
    },
    ownerChanged(usrList) {
      if (usrList && usrList.length === 1 && usrList[0] && usrList[0].id) {
        this.owner = usrList[0];
        this.item.owner = usrList[0].id;
      }
    },
    isGrantDisabled() {
      return !this.grantEnabled;
    },
    isOwnerDisabled() {
      return !this.ownerEnabled;
    },
    isGrantSelectorNotEmpty() {
      return (
        this.grantReadWriteSelect !== null &&
        this.grantReadWriteSelect !== undefined
      );
    },
    onGrantSelect(item, id) {},
    deleteGrant(id) {
      if (id && this.item && this.item.grant) {
        delete this.item.grant[id];
      }
      this.granted = this.granted.filter((item) => item.id !== id);
    },
    OnSelectedShowIconOnly(strEl, change) {
      const e = $(strEl);
      e.find(".my-explanation").remove();
      change($("<div>").append(e).html());
    },
    getMapOfExistingUsersInPermItem() {
      let size = 0;
      const include = {};
      if (this.item.owner) {
        include[this.item.owner] = true;
        size++;
      }
      if (this.item.grant) {
        for (const key in this.item.grant) {
          if (Object.prototype.hasOwnProperty.call(this.item.grant, key)) {
            include[key] = true;
            size++;
          }
        }
      }
      if (size === 0) {
        return null;
      }
      return include;
    },
    refreshGrantList() {
      this.loading = true;
      if (this.item) {
        const include = this.getMapOfExistingUsersInPermItem();
        if (include) {
          axios.post("/api/admin/user/list", { include: include }).then(
            (response) => {
              const granted = [];
              for (let i = 0; i < response.data.length; i++) {
                if (response.data[i].id === this.item.owner) {
                  this.owner = response.data[i];
                } else if (
                  this.item.grant &&
                  response.data[i].id &&
                  this.item.grant[response.data[i].id]
                ) {
                  granted.push(response.data[i]);
                }
              }
              this.granted = granted;
              this.loading = false;
            },
            (err) => {
              this.app.handleError(err);
            }
          );
        } else {
          this.loading = false;
        }
      } else {
        this.loading = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
@use "@/assets/styles/variables" as *;

input.share-link {
  width: 100%;
  height: 40px;
  padding: 3px 5px;
  border: 1px solid #e0e0e0;
  color: #062a85;
  border-top-right-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
}

input[disabled].share-link {
  background-color: rgb(245, 245, 245);
  color: #d6d6d6;
}

.perm-link-btn {
  cursor: pointer;
  text-align: center;
}

.perm-link-btn:hover {
  text-decoration: underline;
}

.share-lnk-row .btn:focus,
.btn.focus {
  outline: 0;
  box-shadow: 0 0 0 0.05rem #062a85;
  background: transparent;
  border: 1px solid $primary;
}

.share-lnk-row .btn {
  padding-left: 5px;
  padding-right: 5px;
}

.perm-modal .modal-footer .btn-secondary {
  display: none;
}

.perm-input {
  width: 100%;
}

.perm-change-wrapper {
  padding: 2px;
}

.perm-change-parent {
  background: $gray-100;
  border: 1px solid #f5f5f5;
  padding: $spacer * 0.5;
}

.perm-input input {
  width: 100%;
  height: 40px;
}

.perm-selp {
  position: relative;
}

.perm-selp i {
  position: absolute;
  left: 12px;
  top: 8px;
  pointer-events: none;
}

.perm-sel {
  padding-left: 42px;
  padding-top: 2px;
  padding-bottom: 2px;
  height: 40px;
  padding-right: 28px;
}

.tdmin > small,
.tdmin > h5 {
  white-space: nowrap;
}

.ss-sel-main .ss-list li i {
  vertical-align: middle;
}

.ss-sel-main {
  height: 40px;
}

.ss-sel-main.perm-grant-selector .ss-select {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
}

.ss-select {
  height: 100%;
}

.perm-disabled {
  position: absolute;
  background: #ffffff7a;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  cursor: not-allowed;
  border: 1px dashed #e2e2e2;
}

.perm-delete {
  text-align: center;
  cursor: pointer;
  color: #9a9a9a;
}

.perm-delete:hover {
  color: #40e1d1;
}

.perm-granted-list {
  max-height: 180px !important;
  overflow: auto;
}
</style>
