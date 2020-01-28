package form

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	r "reflect"

	"github.com/ProxeusApp/proxeus-core/sys/model/compatability"

	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

var (
	ErrMapInterfaceOrJSON = errors.New("param must be either a json string or a map[string]interface{}")
	ErrParamsMissing      = errors.New("param cannot be nil")
)

func Validate(formInput map[string]interface{}, spec interface{}, submit bool) (validate.ErrorMap, error) {
	if submit {
		if spec == nil {
			return nil, ErrParamsMissing
		}
		mSpec, err := toMapStringInterface(spec)
		if err != nil {
			return nil, err
		}
		return _validate(formInput, ComponentsFrom(mSpec))
	}
	if spec != nil {
		mSpec, err := toMapStringInterface(spec)
		if err != nil {
			return nil, err
		}
		return validateProvidedOnly(formInput, ComponentsFrom(mSpec))
	}
	return nil, nil
}

func validateProvidedOnly(formInput map[string]interface{}, compsSpec map[string]interface{}) (validate.ErrorMap, error) {
	errs := make(validate.ErrorMap)
	whitelistedFieldNames := make(map[string]bool)
	LoopComponents(compsSpec, func(compID, compInstID string, compMain interface{}, comp map[string]interface{}) bool {
		fieldName, ok := CompName(comp)
		if !ok {
			return true
		}
		whitelistedFieldNames[fieldName] = true
		if val, ok := formInput[fieldName]; ok {
			validateSpec, ok := CompValidate(comp)
			if ok {
				ferrs := validate.Field(val, validateSpec)
				if len(ferrs) > 0 {
					errs[fieldName] = ferrs
				}
			}
		}
		return true
	})
	for key := range formInput {
		if !whitelistedFieldNames[key] {
			return errs, fmt.Errorf("variable `%s` is not allowed in this form", key)
		}
	}
	return errs, nil
}

func _validate(formInput map[string]interface{}, compsSpec map[string]interface{}) (validate.ErrorMap, error) {
	errs := make(validate.ErrorMap)
	whitelistedFieldNames := make(map[string]bool)
	LoopComponents(compsSpec, func(compID, compInstID string, compMain interface{}, comp map[string]interface{}) bool {
		fieldName, ok := CompName(comp)
		if !ok {
			return true
		}
		whitelistedFieldNames[fieldName] = true
		validateSpec, ok := CompValidate(comp)
		if ok && IsCompVisible(formInput, compsSpec, comp, compInstID) {
			//TODO use validate children here
			GenericLoop(formInput[fieldName], func(i int, val interface{}) bool {
				ferrs := validate.Field(val, validateSpec)
				if len(ferrs) > 0 {
					for _, er := range ferrs {
						if i >= 0 {
							er.I = i
						}
					}
					errs[fieldName] = ferrs
				}
				return true
			})

		}
		return true
	})
	return errs, nil
}

func ValidateFile(src io.Reader, formSrc interface{}, fieldName string) ([]byte, error) {
	return validate.File(src, RulesOf(formSrc, fieldName))
}

func RulesOf(formSrc interface{}, fieldName string) (vRules validate.Rules) {
	LoopComponents(GetFormSrc(formSrc), func(compID, compInstID string, compMain interface{}, comp map[string]interface{}) bool {
		fName, ok := CompName(comp)
		if !ok {
			return true
		}
		if fName == fieldName {
			vRules, ok = CompValidate(comp)
			return false
		}
		return true
	})
	return
}

func toMapStringInterface(d interface{}) (map[string]interface{}, error) {
	m, ok := compatability.ToMapStringIF(d)
	if !ok {
		sSpec, ok := d.(string)
		if ok {
			err := json.Unmarshal([]byte(sSpec), &m)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, ErrMapInterfaceOrJSON
		}
	}
	return m, nil
}

func GetFormSrc(formSrc interface{}) map[string]interface{} {
	if formSrc != nil {
		formSrcMap, ok := compatability.ToMapStringIF(formSrc)
		if ok {
			if formSrcInner, ok := formSrcMap["formSrc"]; ok {
				formSrcInnerMap, ok := formSrcInner.(map[string]interface{})
				if ok {
					return formSrcInnerMap
				}
			}
			return formSrcMap
		}
	}
	return nil
}

func LoopComponents(formSrc interface{}, looper func(compId, compInstId string, compMain interface{}, comp map[string]interface{}) bool) {
	if formSrc != nil {
		formSrcMap, ok := compatability.ToMapStringIF(formSrc)
		if !ok {
			return
		}
		components := ComponentsFrom(formSrcMap)
		for compInstID, item := range components {
			compId := ""
			if citem, ok := item.(map[string]interface{}); ok {
				compId = citem["_compId"].(string)
			} else if comps, ok := item.([]interface{}); ok {
				if len(comps) > 0 {
					if citem, ok := comps[0].(map[string]interface{}); ok {
						compId = citem["_compId"].(string)
					}
				}
			}
			if !GenericLoop(item, func(i int, c interface{}) bool {
				fcomp, ok := c.(map[string]interface{})
				if ok {
					return looper(compId, compInstID, item, fcomp)
				}
				return true
			}) {
				break
			}
		}
	}
}

func Vars(formSrc interface{}) []string {
	vars := make([]string, 0)
	LoopComponents(formSrc, func(compID, compInstID string, compMain interface{}, comp map[string]interface{}) bool {
		name, ok := CompName(comp)
		if ok {
			vars = append(vars, name)
		}
		return true
	})
	return vars
}

func ComponentsFrom(formSrc map[string]interface{}) map[string]interface{} {
	if formSrcInner, ok := formSrc["formSrc"]; ok {
		formSrcInnerMap, ok := formSrcInner.(map[string]interface{})
		if ok {
			if comps, ok := formSrcInnerMap["components"]; ok {
				compsMap, ok := comps.(map[string]interface{})
				if ok {
					return compsMap
				}
			}
		}
	} else {
		if comps, ok := formSrc["components"]; ok {
			compsMap, ok := comps.(map[string]interface{})
			if ok {
				return compsMap
			}
		}
	}
	return formSrc
}

// GenericLoop executes a function on elements of a slice or array
func GenericLoop(m interface{}, f func(int, interface{}) bool) bool {
	v := r.ValueOf(m)
	switch v.Kind() {
	case r.Slice, r.Array:
		if v.IsNil() {
			return f(-1, nil)
		}
		newLen := v.Len()
		if newLen == 0 {
			return f(-1, nil)
		}
		for i := 0; i < newLen; i++ {
			if !f(i, v.Index(i).Interface()) {
				return false
			}
		}
	default:
		return f(-1, m)
	}
	return true
}
