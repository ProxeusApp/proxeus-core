package form

import (
	"fmt"

	"github.com/dop251/goja"

	"git.proxeus.com/core/central/sys/validate"
)

const (
	jsRegex = `
            var createRegex = function(r){
                r = r + "";
                var regexStart = new RegExp('^\/');
                var regexEnd = new RegExp('\/(\w*)$');
                var rStart = regexStart.exec(r), rEnd = regexEnd.exec(r);
                if(rStart && rEnd && rStart.length>0 && rEnd.length>0){
                    r = r.substring(1, r.length-rEnd[0].length);
                    if(rEnd.length>1 && rEnd[1]){
                        return new RegExp(r,rEnd[1]);
                    }else{
                        return new RegExp(r);
                    }
                }else{
                    return new RegExp(r);
                }
            };
`
)

//TODO improve/replace js engine
type JSRegexEval struct {
	vm *goja.Runtime
}

//var tmpVm *goja.Runtime
func NewJSRegexEval() *JSRegexEval {
	js := &JSRegexEval{}
	js.vm = goja.New()
	_, err := js.vm.RunString(jsRegex)
	if err != nil {
		fmt.Println(err)
	}
	return js
}

func (js *JSRegexEval) Test(regex, val string) bool {
	//return false
	if js.vm != nil {
		v, _ := js.vm.RunString("createRegex('" + regex + "').test('" + val + "');")
		return v.ToBoolean()
	}
	return false

}

func (js *JSRegexEval) Close() {

}

func IsCompVisible(formInput map[string]interface{}, comps, compMain map[string]interface{}, destCompID string) bool {
	regexTester := NewJSRegexEval()
	defer regexTester.Close()
	return isCompVisible(regexTester, formInput, comps, compMain, destCompID)
}

func isCompVisible(regexTester *JSRegexEval, formInput map[string]interface{}, comps, fcomp map[string]interface{}, destCompID string) bool {
	destCompID = hasDestination(comps, fcomp, destCompID) //check if it's a group
	if destCompID != "" {
		visible := false
		LoopComponents(comps, func(compID, compInstID string, compMain interface{}, comp map[string]interface{}) bool {
			if compInstID != destCompID { //need to check if the parent is activated
				if ac, ok := comp["action"]; ok {
					sa, ok := ac.(map[string]interface{})
					if ok {
						if src, ok := sa["source"]; ok { //get what activates the destCompID
							source, ok := src.([]interface{})
							if ok {
								for _, sit := range source {
									if sit != nil {
										if sitem, ok := sit.(map[string]interface{}); ok {
											if sDestCompID, ok := sitem["_destCompId"].(string); ok {
												if sDestCompID == destCompID {
													name, ok := CompName(comp)
													if ok && formInput != nil {
														val, ok := formInput[name]
														if ok && val != nil { // PC-499 incomplete formInput will always fail here, by definition a non submit comes with only the changed form input and therefor incomplete
															if isCompVisible(regexTester, formInput, comps, comp, compInstID) {
																GenericLoop(val, func(vi int, v interface{}) bool {
																	if v != nil {
																		regexVal := fmt.Sprintf("%v", v)
																		if regexVal != "nil" && len(regexVal) > 0 {
																			jsRegexStr := fmt.Sprintf("%v", sitem["regex"])
																			visible = regexTester.Test(jsRegexStr, regexVal)
																			if visible {
																				return false
																			}
																		}
																	}
																	return true
																})
															}
															if visible {
																return false
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			if visible {
				return false
			}
			return true
		})
		return visible
	}
	return true
}

func hasDestination(comps, fcomp map[string]interface{}, compId string) string {
	if fcomp != nil {
		action, ok := fcomp["action"].(map[string]interface{})
		if ok {
			if _, ok := action["destination"]; ok {
				return compId
			}
		}
		if b, ok := fcomp["_grouped"]; ok {
			if bo, ok := b.(bool); ok && bo {
				gcid, groupComp := getCompThatImports(compId, comps)
				if groupComp != nil {
					return hasDestination(comps, groupComp, gcid)
				}
			}
		}
	}
	return ""
}

func getCompThatImports(compId string, comps map[string]interface{}) (resCompId string, gCompId map[string]interface{}) {
	if len(comps) > 0 && compId != "" {
		LoopComponents(comps, func(cId, compInstId string, compMain interface{}, comp map[string]interface{}) bool {
			if comp != nil {
				if imp, ok := comp["_import"]; ok {
					if imp != nil {
						for _, impMapped := range imp.(map[string]interface{}) {
							for _, impSliceStr := range impMapped.([]interface{}) {
								if impSliceStr == compId {
									gCompId = comp
									resCompId = compInstId
									return false
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return
}

func CompName(comp map[string]interface{}) (string, bool) {
	if comp != nil {
		name, ok := comp["name"].(string)
		if ok && name != "" {
			return name, true
		}
	}
	return "", false
}

func CompValidateFile(comp map[string]interface{}) (map[string]interface{}, bool) {
	validateMap, ok := CompValidate(comp)
	if ok {
		fileMap, ok := validateMap["file"].(map[string]interface{})
		if ok {
			return fileMap, true
		}
	}
	return nil, false
}

func CompValidate(comp map[string]interface{}) (validate.Rules, bool) {
	if comp != nil {
		validMap, ok := comp["validate"].(map[string]interface{})
		if ok && len(validMap) > 0 {
			return validMap, true
		}
	}
	return nil, false
}
