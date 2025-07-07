import he from "he";
import "./jquery.js";

var FT_FormBuilderCompiler = function (options) {
  this.__init__ = function (options) {
    this.registerTemplateCompilerFeatures();
    this.options = options;
    this.FORM_CONTEXT_VAR_NAME = "_F_M_F_";
    if (!this.options) {
      this.options = {};
    }
    if (!this.options.i18n) {
      this.options.i18n = {};
    }
    this.options.i18n = this.extend(true, this.options.i18n, {
      onSelect: function (data) {
        return { i18n: data };
      },
      onDisplay: function (data) {
        if (
          typeof data === "object" &&
          data.hasOwnProperty("i18n") &&
          data["i18n"]
        ) {
          return data["i18n"];
        }
        return data;
      },
      isCovered: function (data) {
        return (
          typeof data === "object" &&
          data.hasOwnProperty("i18n") &&
          data["i18n"]
        );
      },
    });
  };
  this.randomId = function () {
    return (
      Math.random().toString(36).substring(2, 18) +
      Math.random().toString(36).substring(2, 18)
    );
  };
  this.varToLabel = function (text) {
    if (text && text.length > 0) {
      text = text.charAt(0).toUpperCase() + text.slice(1);
      if (text.length > 1) {
        var ic = [];
        var _char = "";
        for (var i = 1; i < text.length; ++i) {
          _char = text.charAt(i);
          if (_char == _char.toUpperCase()) {
            ic.push(i);
          }
        }
        if (ic.length > 0) {
          var tmpt = "";
          tmpt += text.slice(0, ic[0]);
          for (var c = 0; c < ic.length - 1; ++c) {
            tmpt += " " + text.slice(ic[c], ic[c + 1]);
          }
          tmpt += " " + text.slice(ic[ic.length - 1]);
          text = tmpt;
        }
      }
    }
    return text;
  };
  this.escapeForAttr = function (unsafe) {
    if (typeof unsafe === "string") {
      return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
    }
    return unsafe;
  };
  this.registerTemplateCompilerFeatures = function () {
    var _c = this;
    Handlebars.registerHelper("escapeForAttr", this.escapeForAttr);
    Handlebars.registerHelper("varToLabel", this.varToLabel);
    Handlebars.registerHelper("randomId", this.randomId);
    Handlebars.registerHelper("enumEq", function (a, b, opts) {
      if (_c.isEnum(a)) {
        if (a.all[a.selected] === b) {
          return opts.fn(this);
        } else {
          return opts.inverse(this);
        }
      }
    });
    Handlebars.registerHelper("objectValOf", function (a, b, opts) {
      try {
        return a[b];
      } catch (e) {}
      return "";
    });
    Handlebars.registerHelper("enumName", function (a, b, opts) {
      try {
        if (opts) {
          return b[a.all[a.selected]];
        } else {
          return a.all[a.selected];
        }
      } catch (e) {}
      return "";
    });
    Handlebars.registerHelper("dropHere", function (opts) {
      var dropHerePH;
      try {
        var mainComp, compiled, dropHereIndex;
        if (!opts.data.root["_dropHere"]) {
          opts.data.root["_dropHere"] = {};
          opts.data.root["_dropHere"]["index"] = 0;
        } else {
          ++opts.data.root["_dropHere"]["index"];
        }
        dropHereIndex = opts.data.root["_dropHere"]["index"];
        dropHerePH =
          '<div class="fb-drop-here" data-index="' + dropHereIndex + '"></div>';
        mainComp = _c.getCompMainObject(opts.data.root);
        if (mainComp && mainComp["_import"]) {
          if (mainComp["_import"][dropHereIndex + ""]) {
            var childComps = mainComp["_import"][dropHereIndex + ""];
            if (mainComp.formCompilerId && childComps && childComps.length) {
              var htmlComps = dropHerePH;
              var formCompiler = _c.fcc[mainComp.formCompilerId];
              for (var i = 0; i < childComps.length; i++) {
                if (childComps[i]) {
                  compiled = formCompiler.compileFromGroup(childComps[i]);
                  if (compiled) {
                    htmlComps += compiled;
                  }
                }
              }
              return htmlComps;
            }
          }
        }
      } catch (e) {}
      return dropHerePH;
    });
    Handlebars.registerHelper("ifEq", function (a, b, opts) {
      if (a == b) {
        // Or === depending on your needs
        return opts.fn(this);
      } else {
        return opts.inverse(this);
      }
    });
    Handlebars.registerHelper(
      "math",
      function (lvalue, operator, rvalue, options) {
        lvalue = parseFloat(lvalue);
        rvalue = parseFloat(rvalue);
        return {
          "+": lvalue + rvalue,
          "-": lvalue - rvalue,
          "*": lvalue * rvalue,
          "/": lvalue / rvalue,
          "%": lvalue % rvalue,
        }[operator];
      }
    );
    Handlebars.registerHelper("percentOfOneFrom", function (number, options) {
      return 100.0 / parseFloat(number);
    });
    Handlebars.registerHelper("eachCouple", function (context, options) {
      var out = "";
      var data;
      if (options.data) {
        data = Handlebars.createFrame(options.data);
      }
      data.first = false;
      data.last = false;
      for (var i = 0, i2 = 1; i2 < context.length; i2 = i2 + 2, i = i2 - 1) {
        if (data) {
          if (i == 0) {
            data.first = true;
          }
          if (i == context.length - 1) {
            data.last = true;
          }
          data.index = i;
          data.index2 = i2;
          data.data1 = context[i];
          data.data2 = context[i2];
        }
        out += options.fn(context[i], { data: data });
      }
      return out;
    });
    Handlebars.registerHelper("eachCount", function (context, options) {
      var out = "";
      var data;
      if (options.data) {
        data = Handlebars.createFrame(options.data);
      }
      data.countFirst = false;
      data.countLast = false;
      for (var i = 0; i < context; i++) {
        if (data) {
          if (i == 0) {
            data.countFirst = true;
          }
          if (i == context - 1) {
            data.countLast = true;
          }
          data.countIndex = i;
        }
        out += options.fn(i, { data: data });
      }
      return out;
    });
    Handlebars.registerHelper("ifEven", function (conditional, options) {
      if (conditional % 2 == 0) {
        return options.fn(this);
      } else {
        return options.inverse(this);
      }
    });
  };
  this.compileForm = function (data) {
    if (!this.options || !this.options.requestComp) {
      return "";
    }
    var form = this.getConvertToNewestFormStructure(data.form);
    var i18n = this.extend(true, this.options.i18n, data.i18n);
    var _ = this;
    if (
      i18n &&
      this.isFunction(i18n.onDisplay) &&
      this.isFunction(i18n.onTranslate) &&
      this.isFunction(i18n.isCovered)
    ) {
      i18n.done = function (translatedFormSrc) {
        if (_.isFunction(data.done)) {
          data.done(_.cForm({ form: translatedFormSrc }));
        }
      };
      this.translateSettings(form, null, i18n);
    } else {
      if (this.isFunction(data.done)) {
        data.done(this.cForm({ form: form }));
        return;
      }
      return this.cForm({ form: form });
    }
  };
  this.getConvertToNewestFormStructure = function (form) {
    if (form && form.v && form.v === 2) {
      // newest version
    } else {
      form = { v: 2, components: form };
    }
    return form;
  };
  this.getComps = function (form) {
    if (form && form.v && form.v == 2) {
      // newest version
      return form.components;
    }
    return form;
  };
  this.FormCompilerClass = function (compiler, compWrapper) {
    this.c = compiler;
    this.compWrapper = compWrapper;
    this.compileFromGroup = function (cmpId) {
      if (cmpId) {
        if (this.compWrapper.components) {
          var cp = this.compWrapper.components[cmpId];
          if (cp) {
            var mp = this.c.getCompMainObject(cp);
            if (mp["_compId"]) {
              cp = this.c.deepCopySettings(cp);
              if (this.c.options.injectFormCompiler) {
                var o = this.c.options;
                // translate if coming from the formbuilder
                if (
                  o.i18n &&
                  o.i18n.onDisplay &&
                  o.i18n.onTranslate &&
                  o.i18n.isCovered
                ) {
                  this.c.translateSettings(cp, mp["_compId"], {
                    onDisplay: o.i18n.onDisplay,
                    onTranslate: o.i18n.onTranslate,
                    isCovered: o.i18n.isCovered,
                    notAsync: true,
                    done: function (settings, tmplId) {
                      cp = settings;
                    },
                  });
                }
              }
              return this.c.compileAndSetIdOnRootElement(
                mp["_compId"],
                cp,
                cmpId
              );
            }
          }
        }
      }
      return null;
    };
  };
  this.fcc = {};
  this.cForm = function (data) {
    var formId = this.randomId();
    var comps = data.form.components;
    var _ = this;
    var sortedComps = [];
    var currentComp;
    var c;
    var compMainObj;
    var compActions = {};
    var template;
    var id;
    var compiledComp;
    var allCompsCompiled = "";
    var validationExt = {};
    this.fcc[formId] = new this.FormCompilerClass(this, data.form);
    for (id in comps) {
      if (comps.hasOwnProperty(id)) {
        compMainObj = _.getCompMainObject(comps[id]);
        compMainObj["_comp_Id_"] = id;
        compMainObj["formCompilerId"] = formId;
        sortedComps.push(comps[id]);
      }
    }
    sortedComps.sort(function (a, b) {
      return (
        _.getCompMainObject(a)["_order"] - _.getCompMainObject(b)["_order"]
      );
    });
    _.actionCollector.collected = compActions;
    var htmlStr;
    if (data.componentsOnly) {
      htmlStr = "";
    } else {
      htmlStr = '<form id="' + formId + '">';
    }
    for (var i = 0; i < sortedComps.length; i++) {
      currentComp = sortedComps[i];
      compMainObj = this.getCompMainObject(currentComp);
      id = compMainObj["_comp_Id_"];
      if (!data.componentsOnly) {
        if (this.isArray(currentComp)) {
          if (currentComp.length > 0) {
            for (var a = 0; a < currentComp.length; a++) {
              _.actionCollector.collect(currentComp[a], id);
              this.clientValidationExt(validationExt, currentComp[a]);
            }
          }
        } else {
          _.actionCollector.collect(currentComp, id);
          this.clientValidationExt(validationExt, currentComp);
        }
      }
      if (compMainObj["_grouped"]) {
        continue;
      }
      if (compMainObj["_compId"]) {
        compiledComp = this.compileAndSetIdOnRootElement(
          compMainObj["_compId"],
          currentComp,
          id
        );
        if (data.onCompCompiled) {
          compiledComp = data.onCompCompiled(
            compiledComp,
            compMainObj,
            currentComp
          );
        }
        if (compiledComp) {
          allCompsCompiled += compiledComp;
        } else {
          console.error("error when compiling comp " + id);
          console.log(currentComp);
        }
      } else {
        console.error("error when compiling comp " + id);
        console.log(currentComp);
      }
    }
    if (!allCompsCompiled) {
      return "";
    }
    htmlStr += allCompsCompiled;
    if (!data.componentsOnly) {
      var fmf = this.FORM_CONTEXT_VAR_NAME;
      htmlStr +=
        '   <script type="text/javascript">\n' +
        '       if(!window["' +
        fmf +
        '"]){window["' +
        fmf +
        '"]={}}' +
        '       window["' +
        fmf +
        '"]["' +
        formId +
        '"]= ' +
        this.formRuntimeMain.toString() +
        ";\n" +
        '       window["' +
        fmf +
        '"]["' +
        formId +
        '"]=window["' +
        fmf +
        '"]["' +
        formId +
        '"]();\n' +
        '       window["' +
        fmf +
        '"]["' +
        formId +
        '"].main("' +
        fmf +
        '","' +
        formId +
        '",' +
        JSON.stringify(compActions, null, 2) +
        "," +
        JSON.stringify(validationExt, null, 2) +
        ");\n" +
        "   </script>" +
        "</form>";
    }

    delete this.fcc[formId];
    return htmlStr;
  };
  this.clientValidationExt = function (validationExt, comp) {
    if (comp.name && comp.validate) {
      if (comp.validate.phoneNr || comp.validate.required) {
        if (!validationExt[comp.name]) {
          validationExt[comp.name] = {};
        }
        if (comp.validate.phoneNr) {
          validationExt[comp.name].phoneNr = comp.validate.phoneNr;
        }
        if (comp.validate.required) {
          validationExt[comp.name].required = comp.validate.required;
        }
      }
    }
  };
  this.compileAndSetIdOnRootElement = function (compId, compSettings, htmlId) {
    var c, template, compiledComp, cachedTmpl;
    cachedTmpl = this.getTemplateCached(compId);
    if (!cachedTmpl) {
      c = this.options.requestComp(compId);
      if (!c || !c["template"]) {
        console.error("error component not found " + compId);
        if (this.options.injectFormCompiler && compId && htmlId) {
          return (
            '<div id="' +
            htmlId +
            '" data-dfsid="' +
            compId +
            '" class="fb-component"><span style="color:red;padding:20px;display:inline-block;">Error could not render this component! Looks like it does not exist anymore.</span></div>'
          );
        }
        return;
      }
      template = c["template"];
      cachedTmpl = this.cacheTemplate(compId, template);
    }
    compiledComp = this.cTemplate(template, compSettings, compId);
    if (compiledComp) {
      compiledComp = this.setIdAttrOnRootElement(compiledComp, htmlId);
      if (this.options.injectFormCompiler) {
        return this.addClassOnRootElement(
          this.addAttrOnRootElement(compiledComp, "data-dfsid", compId),
          cachedTmpl.isGrpParent
            ? "fb-component fbc-grp-parent"
            : "fb-component"
        );
      } else {
        return compiledComp;
      }
    } else {
      return "";
    }
  };
  this.compileTemplate = function (templStr, settingsStr, tmplId) {
    if (templStr && settingsStr) {
      var json = null;
      if (typeof settingsStr === "string") {
        json = JSON.parse(settingsStr);
      } else {
        json = this.deepCopySettings(settingsStr);
      }
      return this.cTemplate(templStr, json, tmplId, true);
    }
    return "";
  };
  this.preCompiledCache = {};
  this.cTemplate = function (templStr, json, tmplId, copied) {
    try {
      if (!copied) {
        json = this.deepCopySettings(json);
      }
      this.addFeaturesToTemplateJson(json);
      var html;
      if (tmplId) {
        var preCompTemplate = this.getTemplateCached(tmplId);
        if (preCompTemplate) {
          html = preCompTemplate.tmpl(json);
        } else {
          preCompTemplate = this.cacheTemplate(tmplId, templStr);
          html = preCompTemplate.tmpl(json);
        }
      } else {
        html = Handlebars.compile(templStr)(json);
      }
      if (/^templates\..*/.test(tmplId)) {
        return this.decodeHtml(html);
      } else {
        if ((html = this.isHtmlAndSingleElement(html)) && html) {
          return this.decodeHtml(html);
        }
      }
    } catch (e) {
      console.error(e);
    }
    if (tmplId) {
      console.error("error with " + tmplId);
    }
    return "";
  };
  this.releaseCacheTemplate = function (id) {
    if (this.preCompiledCache[id]) {
      delete this.preCompiledCache[id];
    }
  };
  this.getTemplateCached = function (id) {
    return this.preCompiledCache[id];
  };
  this.cacheTemplate = function (id, tmpl) {
    if (tmpl) {
      tmpl = tmpl.trim();
    }
    return (this.preCompiledCache[id] = {
      tmpl: Handlebars.compile(tmpl),
      isGrpParent: tmpl.indexOf("{{dropHere}}") !== -1,
    });
  };
  this.actionCollector = {
    p: this,
    collected: null,
    removePrivateFields: true,
    collectAll: function (comps, removePrivateFields) {
      if (removePrivateFields != undefined) {
        this.removePrivateFields = removePrivateFields;
      }
      this.collected = {};
      for (var id in comps) {
        if (comps.hasOwnProperty(id)) {
          if (this.p.isArray(comps[id])) {
            if (comps[id].length > 0) {
              for (var a = 0; a < comps[id].length; ++a) {
                this.collect(comps[id][a], id);
              }
            }
          } else {
            this.collect(comps[id], id);
          }
        }
      }
      return this.collected;
    },
    collect: function (comp, id) {
      if (comp.action) {
        this.put(id, comp.name, comp.action);
      }
    },
    put: function (id, name, action) {
      if (!this.collected[id]) {
        this.collected[id] = [];
      }
      if (this.removePrivateFields) {
        if (action.source && action.source.length) {
          for (var i = 0; i < action.source.length; ++i) {
            delete action.source[i]["comment"];
          }
        }
      }
      this.collected[id].push({ name: name, action: action });
    },
  };
  this.formRuntimeMain = function () {
    var $ = window.$;
    $.fn.fmf_fade = function ($form, init, show) {
      if (show) {
        if (!this.is(":visible")) {
          if (init) {
            this.show();
          } else {
            $form.trigger("action-animation-started", [
              { init: init, show: show, target: this },
            ]);
            this.css("opacity", 0);
            this.slideDown(250, "swing", function () {
              $(this).animate({ opacity: 1 }, 400, function () {
                $(this).css("opacity", 1);
                $form.trigger("action-animation-ended", [
                  { init: init, show: show, target: this },
                ]);
              });
            });
          }
        }
      } else {
        if (this.is(":visible")) {
          if (init) {
            this.hide();
          } else {
            $form.trigger("action-animation-started", [
              { init: init, show: show, target: this },
            ]);
            this.animate({ opacity: 0 }, 400, function () {
              $(this).slideUp(250, "swing", function () {
                $(this).css("opacity", 0);
                $form.trigger("action-animation-ended", [
                  { init: init, show: show, target: this },
                ]);
              });
            });
          }
        }
      }
    };

    $.fn.fmf_slide = function ($form, init, show) {
      if (show) {
        if (!this.is(":visible")) {
          if (init) {
            this.show();
          } else {
            $form.trigger("action-animation-started", [
              { init: init, show: show, target: this },
            ]);
            this.slideDown(400, "swing", function () {
              $form.trigger("action-animation-ended", [
                { init: init, show: show, target: this },
              ]);
            });
          }
        }
      } else {
        if (this.is(":visible")) {
          if (init) {
            this.hide();
          } else {
            $form.trigger("action-animation-started", [
              { init: init, show: show, target: this },
            ]);
            $(this).slideUp(400, "swing", function () {
              $form.trigger("action-animation-ended", [
                { init: init, show: show, target: this },
              ]);
            });
          }
        }
      }
    };
    $.fn.fmf_none = function ($form, init, show) {
      if (show) {
        if (!this.is(":visible")) {
          this.show();
        }
      } else {
        if (this.is(":visible")) {
          this.hide();
        }
      }
    };
    $.fn.fm_animate = function ($form, init, show, $caller) {
      var fmShow = this.data("fm_show");
      this.data("fm_hiding", !show);
      if (!fmShow) {
        fmShow = { show: {} };
        this.data("fm_show", fmShow);
      }
      if ($caller) {
        if (!fmShow.show[$caller.attr("name")]) {
          fmShow.show[$caller.attr("name")] = {};
        }
        fmShow.show[$caller.attr("name")]["active"] = show;
      }
      var fmo = $form.data("fmo");
      if (show) {
        if (fmo) {
          try {
            var srcs = fmo.actions[this.attr("id")];
            if (srcs) {
              var innersrcs;
              for (var i = 0; i < srcs.length; ++i) {
                innersrcs = srcs[i].action.source;
                if (innersrcs) {
                  for (var ss = 0; ss < innersrcs.length; ++ss) {
                    $form
                      .find(
                        "[name='" +
                          srcs[i].name +
                          "']:eq(" +
                          innersrcs[ss]._index +
                          ")"
                      )
                      .doCompAction(init);
                  }
                }
              }
            }
          } catch (err) {
            console.log(err);
          }
        }
      } else {
        for (var id in fmShow.show) {
          if (
            fmShow.show.hasOwnProperty(id) &&
            fmShow.show[id]["active"] &&
            !$form.find("#" + fmShow.show[id]["parent"]).data("fm_hiding")
          ) {
            this.data("fm_hiding", false);
            return;
          }
        }
        if (fmo) {
          try {
            var srcs = fmo.actions[this.attr("id")];
            if (srcs) {
              var innersrcs;
              for (var i = 0; i < srcs.length; ++i) {
                innersrcs = srcs[i].action.source;
                if (innersrcs) {
                  for (var ss = 0; ss < innersrcs.length; ++ss) {
                    $.fn.fm_animate.apply(
                      $form.find("#" + innersrcs[ss]._destCompId),
                      [$form, init, show]
                    );
                  }
                }
              }
            }
          } catch (err) {
            console.log(err);
          }
        }
      }
      $form.trigger("action-change", [
        { init: init, show: show, target: this, caller: $caller },
      ]);
      if (this.data("fm_transition")) {
        this["fmf_" + this.data("fm_transition")]($form, init, show);
      } else {
        this.fmf_fade($form, init, show);
      }
    };
    try {
      (function ($) {
        $.event.special.destroyed = {
          remove: function (o) {
            if (o.handler) {
              o.handler();
            }
          },
        };
      })(jQuery);
    } catch (specialEvent) {}

    return {
      main: function (formContextVar, formId, actions, validationExt) {
        var $targetForm = $("#" + formId);
        if (!$targetForm.length) {
          return;
        }
        var fmo = window[formContextVar][formId];
        fmo.form = $targetForm[0];
        fmo.$form = $targetForm;
        $targetForm.parent().bind("DOMNodeRemoved", function (e) {
          try {
            if (
              e &&
              e.target &&
              window[formContextVar][formId].form === e.target
            ) {
              delete window[formContextVar][formId];
            }
          } catch (abc) {}
        });
        $targetForm.bind("destroyed", function (e) {
          try {
            delete window[formContextVar][formId];
          } catch (abc) {}
        });
        $targetForm.data("fmo", fmo);
        fmo.actions = actions;
        if (validationExt) {
          var i1 = 0;
          var i2 = 0;
          for (var vName in validationExt) {
            if (validationExt.hasOwnProperty(vName)) {
              var flds = $targetForm.find("[name='" + vName + "']");
              i2 = 0;
              flds.each(function () {
                var id = $(this).attr("id");
                if (validationExt[vName].required) {
                  $(this).parent().append('<span class="frequired">*</span>');
                }
              });
              i1++;
            }
          }
        }
        var currentComp,
          currentCompItem,
          sources,
          $srcComp,
          destCompId,
          $destComp,
          $fields,
          $targetField;
        for (var id in actions) {
          if (actions.hasOwnProperty(id)) {
            currentComp = actions[id];
            $srcComp = $targetForm.find("#" + id);
            if ($srcComp.length) {
              for (var i = 0; i < currentComp.length; ++i) {
                currentCompItem = currentComp[i];
                if (currentCompItem.action.source) {
                  sources = currentCompItem.action.source;
                  $fields = $targetForm.find(
                    "[name='" + currentCompItem.name + "']"
                  );
                  if ($fields && $fields.length) {
                    for (var s = 0; s < sources.length; ++s) {
                      if (sources[s]) {
                        destCompId = sources[s]["_destCompId"];
                        if (destCompId) {
                          $destComp = $targetForm.find("#" + destCompId);
                          if ($destComp.length) {
                            $targetField = $($fields[sources[s]["_index"]]);
                            if ($targetField && $targetField.length) {
                              if (!$targetField.data("fm_actionData")) {
                                var fmShow = $destComp.data("fm_show");
                                if (!fmShow) {
                                  fmShow = { show: {} };
                                  $destComp.data("fm_show", fmShow);
                                }
                                if (!fmShow.show[currentCompItem.name]) {
                                  fmShow.show[currentCompItem.name] = {};
                                }
                                fmShow.show[currentCompItem.name]["parent"] =
                                  id;
                                fmShow.show[currentCompItem.name][
                                  "active"
                                ] = false;
                                $destComp.hide();
                                $targetField.data(
                                  "fm__index",
                                  sources[s]["_index"]
                                );
                                $targetField.data("fm_actionData", [
                                  sources[s],
                                ]);
                                $targetField.data("fm_actionDataAll", sources);
                                $targetField.data(
                                  "fm_action",
                                  this.actionEvent
                                );
                              } else {
                                var fmShow = $destComp.data("fm_show");
                                if (!fmShow) {
                                  fmShow = { show: {} };
                                  $destComp.data("fm_show", fmShow);
                                }
                                if (!fmShow.show[currentCompItem.name]) {
                                  fmShow.show[currentCompItem.name] = {};
                                }
                                fmShow.show[currentCompItem.name]["parent"] =
                                  id;
                                fmShow.show[currentCompItem.name][
                                  "active"
                                ] = false;
                                $destComp.hide();
                                $targetField
                                  .data("fm_actionData")
                                  .push(sources[s]);
                              }
                            }
                          }
                        }
                      }
                    }
                    if (
                      this.isRadio($($fields[0])) ||
                      this.isRadio($($fields[$fields.length - 1]))
                    ) {
                      var $cRorC;
                      var cRorCIndex = 0;
                      for (var rc = 0; rc < $fields.length; ++rc) {
                        $cRorC = $($fields[rc]);
                        if (this.isRadio($cRorC)) {
                          if (
                            !$cRorC.data("fm__index") &&
                            !$cRorC.data("fm_actionDataAll") &&
                            !$cRorC.data("fm_action")
                          ) {
                            $cRorC.data("fm__index", cRorCIndex);
                            $cRorC.data("fm_actionDataAll", sources);
                            $cRorC.data("fm_action", this.actionEvent);
                          }
                          ++cRorCIndex;
                        }
                      }
                    }
                  }
                }
                if (currentCompItem.action.destination) {
                  var d = currentCompItem.action.destination;
                  if (
                    d.transition &&
                    d.transition.all &&
                    d.transition.selected !== undefined
                  ) {
                    $srcComp.data(
                      "fm_transition",
                      d.transition.all[d.transition.selected]
                    );
                  }
                }
              }
            }
          }
        }
        $targetForm.trigger("dynamicFormScriptExecuted");
      },
      createRegex: function (r) {
        r = r + "";
        var regexStart = new RegExp("^\\/");
        var regexEnd = new RegExp("\\/(\\w*)$");
        var rStart = regexStart.exec(r);
        var rEnd = regexEnd.exec(r);
        if (rStart && rEnd && rStart.length > 0 && rEnd.length > 0) {
          r = r.substring(1, r.length - rEnd[0].length);
          if (rEnd.length > 1 && rEnd[1]) {
            return new RegExp(r, rEnd[1]);
          } else {
            return new RegExp(r);
          }
        } else {
          return new RegExp(r);
        }
      },
      actionEvent: function (init) {
        var actionDataAll = this.data("fm_actionDataAll");
        var actionData = this.data("fm_actionData");
        var $myForm = this.parents("form");
        var markedForHide = {};
        var fmo = $myForm.data("fmo");
        if ($myForm.length && fmo) {
          var $targetComp, regex, elId;
          if (this.attr("type") === "checkbox") {
            if (!this.is(":checked")) {
              if (actionData && actionData.length) {
                elId = actionData[0]["_destCompId"];
                $targetComp = $myForm.find("#" + elId);
                if ($targetComp.length) {
                  markedForHide[elId] = $targetComp;
                }
              }
            }
          } else if (actionDataAll) {
            for (var l = 0; l < actionDataAll.length; ++l) {
              elId = actionDataAll[l]["_destCompId"];
              $targetComp = $myForm.find("#" + elId);
              if ($targetComp.length) {
                if (this.attr("type") === "radio") {
                  markedForHide[elId] = $targetComp;
                } else {
                  if (actionData) {
                    for (var i = 0; i < actionData.length; ++i) {
                      if (actionData[i]["regex"]) {
                        regex = fmo.createRegex(actionData[i]["regex"]);
                        if (!regex.test(this.val())) {
                          markedForHide[elId] = $targetComp;
                        }
                      }
                    }
                  }
                }
              }
            }
          }
          if (actionData) {
            for (var c = 0; c < actionData.length; ++c) {
              $targetComp = $myForm.find("#" + actionData[c]["_destCompId"]);
              if ($targetComp.length) {
                if (
                  this.attr("type") === "radio" ||
                  this.attr("type") === "checkbox"
                ) {
                  if (this.data("fm__index") == actionData[c]["_index"]) {
                    if (this.is(":checked")) {
                      if (markedForHide[$targetComp.attr("id")]) {
                        delete markedForHide[$targetComp.attr("id")];
                      }
                      $targetComp.fm_animate($myForm, init, true, this);
                    }
                  }
                } else {
                  if (actionData[c]["regex"]) {
                    regex = fmo.createRegex(actionData[c]["regex"]);
                    if (regex.test(this.val())) {
                      if (markedForHide[$targetComp.attr("id")]) {
                        delete markedForHide[$targetComp.attr("id")];
                      }
                      $targetComp.fm_animate($myForm, init, true, this);
                    }
                  }
                }
              }
            }
          }
          for (elId in markedForHide) {
            if (markedForHide.hasOwnProperty(elId)) {
              markedForHide[elId].fm_animate($myForm, init, false, this);
            }
          }
        }
      },
      isRadio: function (t) {
        return t.attr("type") === "radio";
      },
    };
  };
  this.setIdAttrOnRootElement = function (htmlStr, id) {
    return this.addAttrOnRootElement(htmlStr, "id", id);
  };
  this.addAttrOnRootElement = function (htmlStr, attrName, attrValue) {
    var rootStartTagReg = new RegExp("<[^>]*");
    var m = rootStartTagReg.exec(htmlStr);
    if (m !== null && m.length) {
      var rootStartTag = m[0];
      rootStartTag = rootStartTag.replace(
        new RegExp(" " + attrName + "(=(\"|')(.*?)(\"|'))?", "g"),
        ""
      );
      rootStartTag =
        rootStartTag.replace(/\s*$/, "") +
        " " +
        attrName +
        '="' +
        attrValue +
        '"';
      htmlStr = htmlStr.replace(rootStartTagReg, rootStartTag);
    }
    return htmlStr;
  };
  this.addClassOnRootElement = function (htmlStr, classVal) {
    var rootStartTagReg = new RegExp("<[^>]*");
    var m = rootStartTagReg.exec(htmlStr);
    if (m !== null && m.length) {
      var rootStartTag = m[0];
      var classAttr = new RegExp("( class(=(\"|')(.*?)(\"|'))?)");
      m = classAttr.exec(rootStartTag);
      if (m !== null && m.length) {
        if (m.length > 3 && m[4] && m[4] !== '"' && m[4] !== "'") {
          classVal = m[4].trim() + " " + classVal.trim();
        }
        rootStartTag = rootStartTag.replace(
          /(\ class(\=(\"|\')(.*?)(\"|\'))?)/,
          ""
        );
      }
      rootStartTag =
        rootStartTag.replace(/\s*$/, "") + ' class="' + classVal + '"';
      htmlStr = htmlStr.replace(rootStartTagReg, rootStartTag);
    }
    return htmlStr;
  };
  this.translateSettings = function (json, tmplId, i18nObj) {
    var _ = this;
    if (json) {
      var keyArray = [];
      var refsToTheVal = [];
      var jsonSettingsCopy = _.deepCopySettings(json);
      _.deepLoopOverJson(jsonSettingsCopy, {
        object: function (value, keyOrIndex, obj) {
          if (i18nObj.isCovered(value)) {
            keyArray.push(i18nObj.onDisplay(value));
            refsToTheVal.push({ k: keyOrIndex, ref: obj });
          }
          return true;
        },
      });
      if (keyArray.length > 0) {
        i18nObj.onTranslate(
          keyArray,
          function (translatedArray) {
            for (var i = 0; i < translatedArray.length; ++i) {
              refsToTheVal[i].ref[refsToTheVal[i].k] = translatedArray[i];
            }
            i18nObj.done(jsonSettingsCopy, tmplId);
          },
          i18nObj.notAsync
        );
      } else {
        i18nObj.done(json, tmplId);
      }
    } else {
      i18nObj.done(json, tmplId);
    }
  };
  /**
   * search recursive for something in a json
   * @param obj {bla:["bla1", "bla2"]}
   * @param options {'string':function(value, keyOrIndex, obj){return false;//to break or true to keep searching},
   *                 'number':function(value, keyOrIndex, obj){..},
   *                 'boolean':function(value, keyOrIndex, obj){..}}
   * @returns {boolean} false to break or true to keep searching
   */
  this.deepLoopOverJson = function (obj, options) {
    var _ = this;
    if (typeof obj === "object") {
      for (var i in obj) {
        if (obj.hasOwnProperty(i)) {
          if (typeof obj[i] === "object") {
            try {
              if (options[typeof obj[i] + ""](obj[i], i, obj)) {
                if (_.deepLoopOverJson(obj[i], options)) {
                  continue;
                } else {
                  return false;
                }
              }
            } catch (functionForTypeMissing) {
              if (_.deepLoopOverJson(obj[i], options)) {
                continue;
              } else {
                return false;
              }
            }
          } else {
            // check for number, boolean or string
            if (obj[i] !== null || obj[i] !== undefined) {
              try {
                if (options[typeof obj[i] + ""](obj[i], i, obj)) {
                  continue;
                } else {
                  return false;
                }
              } catch (functionForTypeMissing) {}
            }
          }
        }
      }
    }
    return true;
  };
  this.addFeaturesToTemplateJson = function (json) {
    if (this.isArray(json)) {
      if (json.length > 0) {
        if (this.options.injectFormCompiler) {
          json["formCompilerId"] = this.getFormCompilerForInjection();
        }
        for (var i = 0; i < json.length; ++i) {
          json[i].id = this.randomId();
        }
      }
    } else {
      if (this.options.injectFormCompiler) {
        json["formCompilerId"] = this.getFormCompilerForInjection();
      }
      json.id = this.randomId();
    }
  };
  this.getFormCompilerForInjection = function () {
    if (!this.formCompilerIdForInjection) {
      if (
        this.options.formCompiler &&
        this.options.formCompiler.getComponentsHolder
      ) {
        this.formCompilerIdForInjection = this.randomId();
        this.fcc[this.formCompilerIdForInjection] = new this.FormCompilerClass(
          this,
          this.options.formCompiler.getComponentsHolder()
        );
      }
    }
    return this.formCompilerIdForInjection;
  };
  this.deepCopySettings = function (settingsOfComp, settingsToMergeInto) {
    if (settingsToMergeInto) {
      if (this.isArray(settingsToMergeInto) && !this.isArray(settingsOfComp)) {
        settingsToMergeInto = settingsToMergeInto[0];
      } else if (
        !this.isArray(settingsToMergeInto) &&
        this.isArray(settingsOfComp)
      ) {
        settingsToMergeInto = [settingsToMergeInto];
      }
      return this.extend(
        true,
        this.isArray(settingsOfComp) ? [] : {},
        settingsOfComp,
        settingsToMergeInto
      );
    }
    return this.extend(
      true,
      this.isArray(settingsOfComp) ? [] : {},
      settingsOfComp
    );
  };
  try {
    if (document && document.createElement) {
      this.textarea = document.createElement("textarea");
    }
  } catch (documentNotFound) {}
  this.isHtmlAndSingleElement = function (str) {
    var d = document.createElement("div");
    if (d && str) {
      str = str.trim();
      d.innerHTML = str;
      if (d.childNodes.length === 1 && d.childNodes[0].nodeType === 1) {
        return d.innerHTML;
      }
    }
    return null;
  };
  this.decodeHtml = function (html) {
    return he.decode(html, {
      isAttributeValue: false,
    });
  };
  this.getCompFieldObjectByName = function (comp, name) {
    name = name + "";
    if (this.isArray(comp)) {
      for (var i = 0; i < comp.length; ++i) {
        if (comp[i].name === name) {
          return comp[i];
        }
      }
    } else {
      if (comp.name === name) {
        return comp;
      }
    }
  };
  this.isGroupComp = function (comp, tmpl) {
    var compMain = this.getCompMainObject(comp);
    if (compMain["_import"]) {
      return true;
    }
    if (tmpl) {
      return tmpl.indexOf("{{dropHere}}") !== -1;
    }
    return false;
  };
  this.compsLoop = function (comps, callback) {
    var comp, k;
    for (k in comps) {
      if (comps.hasOwnProperty(k)) {
        comp = comps[k];
        if (this.isArray(comp)) {
          for (var i = 0; i < comp.length; ++i) {
            if (!callback(k, comp, comp[i])) {
              return;
            }
          }
        } else {
          if (!callback(k, comp, comp)) {
            return;
          }
        }
      }
    }
  };
  this.compsMainLoop = function (comps, callback) {
    var comp, k;
    for (k in comps) {
      if (comps.hasOwnProperty(k)) {
        comp = comps[k];
        if (this.isArray(comp)) {
          if (comp.length > 0) {
            if (!callback(k, comp, comp[0])) {
              return;
            }
          }
        } else {
          if (!callback(k, comp, comp)) {
            return;
          }
        }
      }
    }
  };
  this.cleanActionDst = function (comp) {
    if (comp.action) {
      delete comp.action.destination;
      if (this.sizeOf(comp.action) == 0) {
        delete comp.action;
      }
    }
  };
  this.cleanActionSrc = function (comp) {
    if (
      comp &&
      comp.action &&
      comp.action.source &&
      comp.action.source.length === 0
    ) {
      delete comp.action.source;
      if (this.sizeOf(comp.action) == 0) {
        delete comp.action;
      }
    }
  };
  this.sizeOf = function (obj) {
    var size = 0;
    var key;
    for (key in obj) {
      if (obj.hasOwnProperty(key)) size++;
    }
    return size;
  };
  this.compLoop = function (comp, callback) {
    if (this.isArray(comp)) {
      for (var i = 0; i < comp.length; ++i) {
        if (!callback(comp, comp[i])) {
          return;
        }
      }
    } else {
      if (!callback(comp, comp)) {
      }
    }
  };
  this.getCompMainObject = function (comp) {
    if (comp) {
      if (this.isArray(comp)) {
        if (comp.length > 0) {
          return comp[0];
        }
      } else {
        return comp;
      }
    }
  };
  this.isEnum = function (o) {
    return (
      o &&
      o.hasOwnProperty("selected") &&
      o.hasOwnProperty("all") &&
      o.all.length > 0
    );
  };
  this.isHiddenField = function (key) {
    return /^(id|_.*)$/.test(key);
  };
  var jQuery = this;
  this.class2type = {
    "[object Boolean]": "boolean",
    "[object Number]": "number",
    "[object String]": "string",
    "[object Function]": "function",
    "[object Array]": "array",
    "[object Date]": "date",
    "[object RegExp]": "regexp",
    "[object Object]": "object",
    "[object Error]": "error",
  };
  var toString = this.class2type.toString;

  var hasOwn = this.class2type.hasOwnProperty;

  var fnToString = hasOwn.toString;
  var support = {};
  this.extend = function () {
    var src;
    var copyIsArray;
    var copy;
    var name;
    var options;
    var clone;
    var target = arguments[0] || {};
    var i = 1;
    var length = arguments.length;
    var deep = false;

    // Handle a deep copy situation
    if (typeof target === "boolean") {
      deep = target;

      // skip the boolean and the target
      target = arguments[i] || {};
      i++;
    }

    // Handle case when target is a string or something (possible in deep copy)
    if (typeof target !== "object" && !jQuery.isFunction(target)) {
      target = {};
    }

    // extend jQuery itself if only one argument is passed
    if (i === length) {
      target = this;
      i--;
    }

    for (; i < length; i++) {
      // Only deal with non-null/undefined values
      if ((options = arguments[i]) != null) {
        // Extend the base object
        for (name in options) {
          src = target[name];
          copy = options[name];

          // Prevent never-ending loop
          if (target === copy) {
            continue;
          }

          // Recurse if we're merging plain objects or arrays
          if (
            deep &&
            copy &&
            (jQuery.isPlainObject(copy) || (copyIsArray = jQuery.isArray(copy)))
          ) {
            if (copyIsArray) {
              copyIsArray = false;
              clone = src && jQuery.isArray(src) ? src : [];
            } else {
              clone = src && jQuery.isPlainObject(src) ? src : {};
            }

            // Never move original objects, clone them
            target[name] = jQuery.extend(deep, clone, copy);

            // Don't bring in undefined values
          } else if (copy !== undefined) {
            target[name] = copy;
          }
        }
      }
    }

    // Return the modified object
    return target;
  };

  jQuery.extend({
    // Unique for each copy of jQuery on the page
    // expando: "jQuery" + ( version + Math.random() ).replace( /\D/g, "" ),

    // Assume jQuery is ready without the ready module
    isReady: true,

    error: function (msg) {
      throw new Error(msg);
    },

    noop: function () {},

    // See test/unit/core.js for details concerning isFunction.
    // Since version 1.3, DOM methods and functions like alert
    // aren't supported. They return false on IE (#2968).
    isFunction: function (obj) {
      return jQuery.type(obj) === "function";
    },

    isArray:
      Array.isArray ||
      function (obj) {
        return jQuery.type(obj) === "array";
      },

    isWindow: function (obj) {
      /* jshint eqeqeq: false */
      return obj != null && obj == obj.window;
    },

    isNumeric: function (obj) {
      // parseFloat NaNs numeric-cast false positives (null|true|false|"")
      // ...but misinterprets leading-number strings, particularly hex literals ("0x...")
      // subtraction forces infinities to NaN
      // adding 1 corrects loss of precision from parseFloat (#15100)
      return !jQuery.isArray(obj) && obj - parseFloat(obj) + 1 >= 0;
    },

    isEmptyObject: function (obj) {
      var name;
      for (name in obj) {
        return false;
      }
      return true;
    },

    isPlainObject: function (obj) {
      var key;

      // Must be an Object.
      // Because of IE, we also have to check the presence of the constructor property.
      // Make sure that DOM nodes and window objects don't pass through, as well
      if (
        !obj ||
        jQuery.type(obj) !== "object" ||
        obj.nodeType ||
        jQuery.isWindow(obj)
      ) {
        return false;
      }

      try {
        // Not own constructor property must be Object
        if (
          obj.constructor &&
          !hasOwn.call(obj, "constructor") &&
          !hasOwn.call(obj.constructor.prototype, "isPrototypeOf")
        ) {
          return false;
        }
      } catch (e) {
        // IE8,9 Will throw exceptions on certain host objects #9897
        return false;
      }

      // Support: IE<9
      // Handle iteration over inherited properties before own properties.
      if (support.ownLast) {
        for (key in obj) {
          return hasOwn.call(obj, key);
        }
      }

      // Own properties are enumerated firstly, so to speed up,
      // if last one is own, then all properties are own.
      for (key in obj) {
      }

      return key === undefined || hasOwn.call(obj, key);
    },

    type: function (obj) {
      if (obj == null) {
        return obj + "";
      }
      return typeof obj === "object" || typeof obj === "function"
        ? jQuery.class2type[toString.call(obj)] || "object"
        : typeof obj;
    },

    // Evaluates a script in a global context
    // Workarounds based on findings by Jim Driscoll
    // http://weblogs.java.net/blog/driscoll/archive/2009/09/08/eval-javascript-global-context
    globalEval: function (data) {
      if (data && jQuery.trim(data)) {
        // We use execScript on Internet Explorer
        // We use an anonymous function so that context is window
        // rather than jQuery in Firefox
        (
          window.execScript ||
          function (data) {
            window["eval"].call(window, data);
          }
        )(data);
      }
    },

    // Convert dashed to camelCase; used by the css and data modules
    // Microsoft forgot to hump their vendor prefix (#9572)
    camelCase: function (string) {
      return string.replace(rmsPrefix, "ms-").replace(rdashAlpha, fcamelCase);
    },

    nodeName: function (elem, name) {
      return (
        elem.nodeName && elem.nodeName.toLowerCase() === name.toLowerCase()
      );
    },

    // args is for internal usage only
    each: function (obj, callback, args) {
      var value;
      var i = 0;
      var length = obj.length;
      var isArray = isArraylike(obj);

      if (args) {
        if (isArray) {
          for (; i < length; i++) {
            value = callback.apply(obj[i], args);

            if (value === false) {
              break;
            }
          }
        } else {
          for (i in obj) {
            value = callback.apply(obj[i], args);

            if (value === false) {
              break;
            }
          }
        }

        // A special, fast, case for the most common use of each
      } else {
        if (isArray) {
          for (; i < length; i++) {
            value = callback.call(obj[i], i, obj[i]);

            if (value === false) {
              break;
            }
          }
        } else {
          for (i in obj) {
            value = callback.call(obj[i], i, obj[i]);

            if (value === false) {
              break;
            }
          }
        }
      }

      return obj;
    },

    // Support: Android<4.1, IE<9
    trim: function (text) {
      return text == null ? "" : (text + "").replace(rtrim, "");
    },

    // results is for internal usage only
    makeArray: function (arr, results) {
      var ret = results || [];

      if (arr != null) {
        if (isArraylike(Object(arr))) {
          jQuery.merge(ret, typeof arr === "string" ? [arr] : arr);
        } else {
          push.call(ret, arr);
        }
      }

      return ret;
    },

    inArray: function (elem, arr, i) {
      var len;

      if (arr) {
        if (indexOf) {
          return indexOf.call(arr, elem, i);
        }

        len = arr.length;
        i = i ? (i < 0 ? Math.max(0, len + i) : i) : 0;

        for (; i < len; i++) {
          // Skip accessing in sparse arrays
          if (i in arr && arr[i] === elem) {
            return i;
          }
        }
      }

      return -1;
    },

    merge: function (first, second) {
      var len = +second.length;
      var j = 0;
      var i = first.length;

      while (j < len) {
        first[i++] = second[j++];
      }

      // Support: IE<9
      // Workaround casting of .length to NaN on otherwise arraylike objects (e.g., NodeLists)
      if (len !== len) {
        while (second[j] !== undefined) {
          first[i++] = second[j++];
        }
      }

      first.length = i;

      return first;
    },

    grep: function (elems, callback, invert) {
      var callbackInverse;
      var matches = [];
      var i = 0;
      var length = elems.length;
      var callbackExpect = !invert;

      // Go through the array, only saving the items
      // that pass the validator function
      for (; i < length; i++) {
        callbackInverse = !callback(elems[i], i);
        if (callbackInverse !== callbackExpect) {
          matches.push(elems[i]);
        }
      }

      return matches;
    },

    // arg is for internal usage only
    map: function (elems, callback, arg) {
      var value;
      var i = 0;
      var length = elems.length;
      var isArray = isArraylike(elems);
      var ret = [];

      // Go through the array, translating each of the items to their new values
      if (isArray) {
        for (; i < length; i++) {
          value = callback(elems[i], i, arg);

          if (value != null) {
            ret.push(value);
          }
        }

        // Go through every key on the object,
      } else {
        for (i in elems) {
          value = callback(elems[i], i, arg);

          if (value != null) {
            ret.push(value);
          }
        }
      }

      // Flatten any nested arrays
      return concat.apply([], ret);
    },

    // A global GUID counter for objects
    guid: 1,

    now: function () {
      return +new Date();
    },

    // jQuery.support is not used in Core but other projects attach their
    // properties to it so it needs to exist.
    support: support,
  });
  this.__init__(options);
};

export default FT_FormBuilderCompiler;
