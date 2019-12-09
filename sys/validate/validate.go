package validate

import (
	"fmt"
	r "reflect"
	"strconv"
	"strings"

	//"github.com/c2h5oh/datasize"
	"log"
	"regexp"
)

const msgBadDefinitionOf = "bad definition of"

var allMessages = []string{
	msgRequired,
	msgMinLowly,
	msgErrWhenParsingNumber,
	msgBadDefinitionForMatches,
	msgBadDefinitionOfMax,
	msgBadDefinitionOfMin,
	msgEmailInvalid,
	msgMaxExceeded,
	msgNotMatchingRegex,
	msgNumberInvalid,
	msgPhoneNrInvalid,
	msgUrlInvalid,
	msgDatePatternInvalid,
	msgBadDefinitionOf,
}

type structRules struct {
	name  string
	rules Rules
}

type structHandler struct {
	ErrorMap ErrorMap
}

//Struct validates the interface by the struct tag `validate`
//Example
/*
type Customer struct {
	Age   int    `validate:"required=true,min=5,max=10"`
	Email string `validate:"required=true,email=true"`
}

*/
func Struct(v interface{}) error {
	if v != nil {
		val := _indirect(r.ValueOf(v))
		strType := val.Type()
		if strType.Kind() != r.Struct {
			return fmt.Errorf("not a struct! type:%v val:%v", strType, val.Interface())
		}
		me := &structHandler{ErrorMap: make(ErrorMap)}
		err := me.validate(val, "")
		if err != nil {
			return err
		}
		if len(me.ErrorMap) > 0 {
			return me.ErrorMap
		}
	}
	return nil
}

func _indirect(val r.Value) r.Value {
	for i := 0; i < 10; i++ {
		if val.Kind() == r.Ptr {
			val = r.Indirect(val)
		} else {
			break
		}
	}
	return val
}

func AllMessages() []string {
	return allMessages
}

func getRules(f *r.StructField) *structRules {
	validateRules := f.Tag.Get("validate")
	if strings.TrimSpace(validateRules) != "" {
		sc := structRules{rules: makeRules(validateRules)}
		sc.name = f.Name
		setJsonFieldNameIfExists(f, &sc.name)
		return &sc
	}
	return nil
}

func setJsonFieldNameIfExists(f *r.StructField, name *string) {
	//check json tag as we serialize outgoing data to json
	jsonTagVal := f.Tag.Get("json")
	if jsonTagVal != "" {
		*name = strings.Split(jsonTagVal, ",")[0]
	}
}

func (me *structHandler) validate(val r.Value, structFieldPath string) error {
	val = _indirect(val)
	strType := val.Type()
	if strType.Kind() != r.Struct {
		return fmt.Errorf("not a struct! type:%v val:%v path:%s", strType, val.Interface(), structFieldPath)
	}
	size := val.NumField()
	for i := 0; i < size; i++ {
		tf := strType.Field(i)
		sRules := getRules(&tf)
		name := tf.Name
		f := val.Field(i)
		if sRules != nil {
			intrf := f.Interface()
			name = sRules.name
			errs := Field(intrf, sRules.rules)
			if len(errs) > 0 {
				me.ErrorMap[structFieldPath+sRules.name] = errs
			}
		} else {
			setJsonFieldNameIfExists(&tf, &name)
		}
		f = _indirect(f)

		//deep dive only if we have errors
		switch f.Kind() {
		case r.Struct:
			err := me.validate(f, structFieldPath+name+".")
			if err != nil {
				return err
			}
			break

		case r.Slice, r.Array:
			if f.IsNil() || f.Len() == 0 {
				break
			}
			for i := 0; i < f.Len(); i++ {
				n := f.Index(i)
				k := n.Kind()
				if k == r.Struct || k == r.Ptr {
					err := me.validate(n, fmt.Sprintf("%s.%d.", structFieldPath+name, i))
					if err != nil {
						return err
					}
				} else {
					if sRules != nil {
						if children, ok := sRules.rules[children].(Rules); ok {
							errs := Field(n.Interface(), children)
							if len(errs) > 0 {
								for _, er := range errs {
									er.I = i
								}
								me.ErrorMap[fmt.Sprintf("%s.%d", structFieldPath+name, i)] = errs
							}
						}
					}
				}
			}
		case r.Map:
			if f.IsNil() || f.Len() == 0 {
				break
			}
			for _, mk := range f.MapKeys() {
				n := f.MapIndex(mk)
				k := n.Kind()
				if k == r.Struct || k == r.Ptr {
					err := me.validate(n, fmt.Sprintf("%s.%v.", structFieldPath+name, mk))
					if err != nil {
						return err
					}
				} else {
					if sRules != nil {
						if children, ok := sRules.rules[children].(Rules); ok {
							errs := Field(n.Interface(), children)
							if len(errs) > 0 {
								for _, er := range errs {
									er.I = mk
								}
								me.ErrorMap[fmt.Sprintf("%s.%v", structFieldPath+name, mk)] = errs
							}
						}
					}
				}
			}
		}
	}
	return nil
}

const children = "children"

var childrenReg = regexp.MustCompile(`(?i)children=\[[^\]]+\],?`)

func makeRules(rules string) Rules {
	if strings.Contains(rules, children) {
		all := childrenReg.FindAllStringSubmatch(rules, -1)
		if len(all) == 1 && len(all[0]) == 1 {
			var childrenRules Rules
			childs := all[0][0]
			rules = strings.Replace(rules, childs, "", 1)
			rules = strings.TrimSpace(strings.Trim(strings.TrimSpace(rules), ","))
			childrenRules = _makeRules(childs[strings.Index(childs, "[")+1 : strings.Index(childs, "]")])
			if rules != "" {
				mainRules := _makeRules(rules)
				if len(mainRules) == 0 {
					mainRules = Rules{}
				}
				if len(childrenRules) > 0 {
					mainRules[children] = childrenRules
					return mainRules
				}
				return nil
			} else {
				if len(childrenRules) > 0 {
					return Rules{children: childrenRules}
				}
				return nil
			}
		}
	}
	return _makeRules(rules)
}

func _makeRules(rules string) Rules {
	commaSplit := strings.Split(rules, ",")
	convRules := make(Rules)
	for _, b := range commaSplit {
		vv := strings.Split(b, "=")
		key := strings.ToLower(strings.TrimSpace(vv[0]))
		if len(vv) == 1 {
			convRules[key] = nil
		} else if len(vv) == 2 {
			convRules[key] = vv[1]
		}
	}
	return convRules
}

// Rules used for validation
type Rules map[string]interface{}

func (r Rules) isNumber() bool {
	number, ok := r["number"].(bool)
	return ok && number
}

func (v *validator) isDefinedAs(key string) bool {
	if rule, k := v.rules[key]; k {
		var req bool
		var ok bool
		if req, ok = rule.(bool); !ok {
			var err error
			req, err = strconv.ParseBool(fmt.Sprintf("%v", rule))
			if err != nil {
				log.Println(err, rule)
				*v.errs = append(*v.errs, &Error{Msg: fmt.Sprintf(msgBadDefinitionOf+" %s", key)})
			}
		}
		return req
	}
	return false
}

func (v *validator) hasStrValueFor(key string) (string, bool) {
	if rule, ok := v.rules[key]; ok {
		return fmt.Sprintf("%v", rule), ok
	}
	return "", false
}

func (r Rules) isFile() bool {
	_, ok := r["file"]
	return ok
}

//FieldByStrRules accepts rules like "required=true,email=true,matches=^(abc|123)$"
func FieldByStrRules(val interface{}, rules string) Errors {
	return Field(val, makeRules(rules))
}

// Field validates a val with the given rules
func Field(val interface{}, rules Rules) Errors {
	validr := newValidator(val, rules)

	if validr.isDefinedAs("required") {
		validr.required()
	}
	if validr.isDefinedAs("email") {
		validr.email()
	}
	if validr.isDefinedAs("number") {
		validr.number()
	}
	if validr.isDefinedAs("url") {
		validr.url()
	}
	if validr.isDefinedAs("phonenr") {
		validr.phoneNr()
	}
	if v, ok := validr.hasStrValueFor("min"); ok {
		validr.min(v)
	}
	if v, ok := validr.hasStrValueFor("max"); ok {
		validr.max(v)
	}
	if v, ok := validr.hasStrValueFor("matches"); ok {
		validr.matches(v)
	}
	if v, ok := validr.hasStrValueFor("datepattern"); ok {
		validr.datePattern(v)
	}
	if len(*validr.errs) == 0 {
		return nil
	}
	return *validr.errs
}

/**
      validate:{
     	//@method:@spec
          required:true|false,
          email:true|false,
          number:true|false,
          matches:'regex',
          max:[0-9]+,
          min:[0-9]+,
          url:true|false
		  date:dd.mm.yyyy HH:ii
      }
 **/
type validator struct {
	val      r.Value
	rules    Rules
	errs     *Errors
	isNumber bool
}

func newValidator(val interface{}, rules Rules) *validator {
	if rules == nil {
		return nil
	}
	ves := make(Errors, 0)
	validr := &validator{}
	validr.rules = rules
	validr.errs = &ves
	validr.val = r.ValueOf(val)
	//ensure keys are lowercase
	for k, v := range rules {
		lk := strings.ToLower(k)
		if k != lk {
			rules[lk] = v
			delete(rules, k)
		}
	}
	return validr
}

func (v *validator) add(error *Error) {
	v.errs.add(error)
}
