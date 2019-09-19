import Vue from 'vue'

const state = {
  components: [
    {
      id: 'C1',
      settings: {},
      code: '<template>\n' +
      '    <div class="form-group">\n' +
      '        <label for="exampleInputEmail1" >{{ settings.label }}</label>\n' +
      '        <input type="email" class="form-control" :placeholder="settings.placeholder" v-model="userData.value" :id="_uid">\n' +
      '      <small id="emailHelp" class="form-text text-muted">{{ settings.help }}</small>\n' +
      '    </div>\n' +
      '</template>\n' +
      '<script>\n' +
      '    module.exports = {\n' +
      '        mounted () {\n' +
      '          console.log(this.$parent.eventBus);\n' +
      '          console.log(\'tis was component\');\n' +
      '        },\n' +
      '        methods: {\n' +
      '          actionClicked(event) {\n' +
      '            this.$emit(\'actionFromInput\', event)\n' +
      '          }\n' +
      '        },\n' +
      '        data() {\n' +
      '            return {\n' +
      '                compName: \'Single-line\',\n' +
      '                icon: \'textbox\',\n' +
      '                userData: {\n' +
      '                    value: \'\'\n' +
      '                },\n' +
      '                settings: {\n' +
      '        \t\t\t\t\thelp: \'help texttttt\',\n' +
      '        \t\t\t\t\tlabel: \'Simple Field\',\n' +
      '        \t\t\t\t\tname: \'\',\n' +
      '        \t\t\t\t\tplaceholder: \'Placeholder\',\n' +
      '        \t\t\t\t\tvalidate: [\n' +
      '        \t\t\t\t\t]\n' +
      '               }\n' +
      '            }\n' +
      '        }\n' +
      '    }\n' +
      '</script>'
    }, {
      id: 'C1c',
      settings: {},
      code: '<template>\n' +
      '    <div class="form-group">\n' +
      '        <label :for="_uid" >{{ settings.label }}</label>\n' +
      '      <textarea class="form-control" v-model="userData.value" :id="_uid" :placeholder="settings.placeholder"></textarea>\n' +
      '      <small class="form-text text-muted">{{ settings.help }}</small>\n' +
      '    </div>\n' +
      '</template>\n' +
      '<script>\n' +
      '    module.exports = {\n' +
      '        mounted () {\n' +
      '          console.log(this.$parent.eventBus);\n' +
      '          console.log(\'tis was component\');\n' +
      '        },\n' +
      '        methods: {\n' +
      '          actionClicked(event) {\n' +
      '            this.$emit(\'actionFromInput\', event)\n' +
      '          }\n' +
      '        },\n' +
      '        data() {\n' +
      '            return {\n' +
      '                compName: \'Textarea\',\n' +
      '                icon: \'cursor-text\',\n' +
      '                userData: {\n' +
      '                    value: \'\'\n' +
      '                },\n' +
      '                settings: {\n' +
      '        \t\t\t\t\thelp: \'help text\',\n' +
      '        \t\t\t\t\tlabel: \'Simple Textarea\',\n' +
      '        \t\t\t\t\tname: \'\',\n' +
      '        \t\t\t\t\tplaceholder: \'Placeholder\',\n' +
      '        \t\t\t\t\tvalidate: [\n' +
      '        \t\t\t\t\t]\n' +
      '               }\n' +
      '            }\n' +
      '        }\n' +
      '    }\n' +
      '</script>'
    }],
  formComponents: [
    {
      id: 'FC1',
      cid: 'C1',
      settings: {}
    }, {
      id: 'FC2',
      cid: 'C1c',
      settings: {}
    }
  ]
}

const mutations = {
  SET_COMPONENTS (state, { components }) {
    state.components = components
  },
  SET_FORMCOMPONENT_SETTING (state, { component, setting }) {
    let fc = state.formComponents.find(
      formComponent => formComponent.id === component.id
    )
    // fc.settings[setting.k] = setting.v
    component.settings[setting.k] = setting.v
    // Vue.set(fc.settings, setting.k, setting.v)
    state.formComponents.splice(state.formComponents.indexOf(fc), 1, component)
  },
  SET_FORMCOMPONENT_SETTINGS (state, { component, setting }) {
    let fc = state.formComponents.find(
      formComponent => formComponent.cid === component.cid)
    fc.settings = settings
  },
  SET_COMPONENT_CODE (state, { component, code }) {
    let comp = getters.byId(component.id)
    // console.log(comp)
    // component.code = code
    if (comp) {
      state.components = state.components.map((c) => {
        if (c.id === component.id) {
          c.code = code
        }
        return c
      })
    }
  },
  UPDATE_COMPONENT (state, { component }) {
    const comp = getters.byId(state, component.id)
    if (comp) {
      state.components.splice(state.components.indexOf(comp), 1, component)
    }
  },
  ADD_BASECOMPONENT (state, { comp }) {
    state.components.push(comp)
  },
  ADD_COMPONENT (state, { component }) {
    state.components.push(component)
  },
  SET_FORMCOMPONENTS (state, forms) {
    state.formComponents = forms
  },
  ADD_FORMCOMPONENT (state, { comp, form }) {
    state.formComponents.push(comp)
  }
}

const getters = {
  byId: (state) => (id) => {
    console.log('test')
    return state.components.find(component => component.id === id)
  },
  tt: (state) => (id) => {
    return state.formComponents.find(formComponent => formComponent.cid === id)
  }
}

const actions = {
  UPDATE_COMPONENT ({ commit }, { comp }) {
    commit('UPDATE_COMPONENT', { component: comp })
  },
  UPDATE_COMPONENT_CODE ({ commit }, { comp, code }) {
    commit('SET_COMPONENT_CODE', { component: comp, code: code })
  }
}

export default {
  state, mutations, actions, getters
}
