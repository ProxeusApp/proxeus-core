package file

import (
	"encoding/json"
	r "reflect"
	"regexp"
	"strconv"
	"strings"
)

type MapIO map[string]interface{}

//MakeFileInfos makes sure all file map patterns are replaced by actual Info's
func (me MapIO) MakeFileInfos(baseDir string) {
	for k, v := range me {
		check(r.ValueOf(v), r.ValueOf(me), r.ValueOf(k), baseDir, nil)
	}
}

func (me MapIO) GetAllFileInfos(baseDir string) []*IO {
	items := make([]*IO, 0)
	for k, v := range me {
		check(r.ValueOf(v), r.ValueOf(me), r.ValueOf(k), baseDir, func(fio *IO) bool {
			items = append(items, fio)
			return false
		})
	}
	return items
}

func (me MapIO) GetAllDataAndFiles(baseDir string) (dat map[string]interface{}, files []string) {
	bts, err := json.Marshal(me)
	if err != nil {
		return nil, nil
	}
	dat = map[string]interface{}{}
	err = json.Unmarshal(bts, &dat)
	if err != nil {
		return nil, nil
	}
	z := r.Value{}
	files = []string{}
	check(r.ValueOf(dat), z, z, baseDir, func(fio *IO) bool {
		files = append(files, fio.Path())
		return false
	})
	return
}

func (me MapIO) MergeWith(v map[string]interface{}) {
	mergeData(me, "", v)
}

var dataPathReg = regexp.MustCompile(`\[|\]|\.`)

func dataPathUnmarshal(dataPath string) []string {
	if dataPath == "" {
		return []string{""}
	}
	arr := dataPathReg.Split(strings.TrimSpace(dataPath), -1)
	cleanPath := make([]string, 0)
	for _, a := range arr {
		if a != "" {
			cleanPath = append(cleanPath, a)
		}
	}
	return cleanPath
}

func mergeData(all map[string]interface{}, path string, toMerge interface{}) (merged interface{}) {
	if toMerge == nil {
		return nil
	}

	path = strings.TrimPrefix(path, "$")
	if all != nil {
		mergeInner(all, dataPathUnmarshal(path), 0, toMerge)
		return all
	}
	newMap := map[string]interface{}{}
	mergeInner(newMap, dataPathUnmarshal(path), 0, toMerge)
	return newMap
}

func mergeInner(target map[string]interface{}, path []string, i int, toMerge interface{}) {
	var t interface{}
	if path[i] == "" {
		t = target
	} else {
		t = target[path[i]]
	}
	if t != nil {
		if len(path)-1 == i {
			if tm, ok := t.(map[string]interface{}); ok {
				if tmm, ok := toMerge.(map[string]interface{}); ok {
					merge(tm, path, i, tmm)
				} else {
					target[path[i]] = toMerge
				}
			} else {
				target[path[i]] = toMerge
			}
		} else {
			if tmm, ok := t.(map[string]interface{}); ok {
				mergeInner(tmm, path, i+1, toMerge)
			} else {
				newMap := map[string]interface{}{}
				target[path[i]] = newMap
				mergeInner(newMap, path, i+1, toMerge)
			}
		}
	} else {
		if tmm, ok := toMerge.(map[string]interface{}); ok {
			merge(target, path, i, tmm)
		} else {
			if len(path)-1 == i {
				target[path[i]] = toMerge
			} else {
				newMap := map[string]interface{}{}
				target[path[i]] = newMap
				mergeInner(newMap, path, i+1, toMerge)
			}
		}

	}
}

func merge(target map[string]interface{}, path []string, i int, toMerge map[string]interface{}) {
	for k, v := range toMerge {
		if k == "" {
			continue
		}
		if target[k] != nil {
			if m, ok := target[k].(map[string]interface{}); ok {
				if om, ok := v.(map[string]interface{}); ok {
					path = append(path, k)
					mergeInner(m, path, i+1, om)
					continue
				}
			}
		}
		target[k] = v
	}
}

//Get retrieves the value by a path like key.key2.key3
func (me MapIO) Get(path string) interface{} {
	if path == "" {
		return me
	}
	return access(dataPathUnmarshal(path), 0, r.ValueOf(me))
}

//GetFileInfo retrieves the value by a path like key.key2.key3
func (me MapIO) GetFileInfo(baseDir, path string) *IO {
	maybeFi := me.Get(path)
	if fi, ok := maybeFi.(*IO); ok {
		return fi
	} else if m, ok := maybeFi.(map[string]interface{}); ok && IsFileInfo(m) {
		return FromMap(baseDir, m)
	}
	return nil
}

func access(path []string, i int, v r.Value) interface{} {
	for i := 0; i < 10; i++ {
		if v.Kind() == r.Ptr {
			v = r.Indirect(v)
		} else {
			break
		}
	}
	if !v.IsValid() {
		return nil
	}
	if v.CanInterface() {
		v = r.ValueOf(v.Interface())
	}
	if len(path)-1 >= i {
		switch v.Kind() {
		case r.Slice, r.Array:
			if v.IsNil() {
				return nil
			}
			length := v.Len()
			if length == 0 {
				return nil
			}
			ii, err := strconv.Atoi(path[i])
			if err != nil || length-1 < ii || ii < 0 {
				return nil
			}
			return access(path, i+1, v.Index(ii))
		case r.Map:
			if v.IsNil() {
				return nil
			}
			if v.Len() == 0 {
				return nil
			}
			return access(path, i+1, v.MapIndex(r.ValueOf(path[i])))
		case r.Struct:
			return access(path, i+1, v.FieldByName(path[i]))
		}
	}
	return v.Interface()
}

func check(v r.Value, parent r.Value, key r.Value, baseDir string, cb func(fio *IO) bool) {
	for i := 0; i < 10; i++ {
		if v.Kind() == r.Ptr {
			v = r.Indirect(v)
		} else {
			break
		}
	}
	if v.CanInterface() {
		v = r.ValueOf(v.Interface())
	}
	switch v.Kind() {
	case r.Slice, r.Array:
		if v.IsNil() {
			return
		}
		newLen := v.Len()
		if newLen == 0 {
			return
		}
		for i := 0; i < newLen; i++ {
			check(v.Index(i), v, r.ValueOf(i), baseDir, cb)
		}
	case r.Map:
		if v.IsNil() {
			return
		}
		newLen := v.Len()
		if newLen == 0 {
			return
		}
		if f, ok := v.Interface().(map[string]interface{}); ok && IsFileInfo(f) {
			switch parent.Kind() {
			case r.Map:
				fi := FromMap(baseDir, f)
				if cb == nil || cb(fi) {
					parent.SetMapIndex(key, r.ValueOf(fi))
				}
			case r.Slice, r.Array:
				fi := FromMap(baseDir, f)
				if cb == nil || cb(fi) {
					parent.Index(int(key.Int())).Set(r.ValueOf(fi))
				}
			}
		} else {
			for _, mk := range v.MapKeys() {
				check(v.MapIndex(mk), v, mk, baseDir, cb)
			}
		}
	}
}
