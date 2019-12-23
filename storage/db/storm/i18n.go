package storm

import (
	//"fmt"
	//"github.com/ProxeusApp/proxeus-core/sys/model"
	//"github.com/asdine/storm/codec/msgpack"
	//"github.com/asdine/storm/q"
	//"path/filepath"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/i18n"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

//
type I18nDB struct {
	db           *storm.DB
	resolver     *i18n.I18nResolver
	langReg      *regexp.Regexp
	allCache     map[string]map[string]string
	allCacheLock sync.RWMutex
	langs        map[string]*model.Lang
	langSlice    []*model.Lang
	fallbackLang string
}

type i18nInternal struct {
	ID   string //=Lang+Code
	Code string //the lang internal unique text identifier like 'max.exceeded'
	Lang string //like en, de
	Text string //the Text of the Code 'max.exceeded' with the Lang 'en' could be 'max length exceeded'
}

func NewI18nDB(dir string) (*I18nDB, error) {
	var err error
	var msgpackDb *storm.DB
	err = ensureDir(dir)
	if err != nil {
		return nil, err
	}
	msgpackDb, err = storm.Open(filepath.Join(dir, "i18n"), storm.Codec(msgpack.Codec))
	if err != nil {
		return nil, err
	}
	udb := &I18nDB{langs: map[string]*model.Lang{}, resolver: &i18n.I18nResolver{}, allCache: map[string]map[string]string{}, db: msgpackDb, langReg: regexp.MustCompile(`^[A-Za-z_]{2,6}$`)}
	example := &i18nInternal{}
	err = udb.db.Init(example)
	if err != nil {
		return nil, err
	}
	udb.langSlice, err = udb.GetAllLangs()
	if err == storm.ErrNotFound {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	for _, lang := range udb.langSlice {
		udb.langs[lang.Code] = lang
	}
	udb.fallbackLang, err = udb.GetFallback()
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func (me *I18nDB) Find(keyContains string, valueContains string, options storage.Options) (map[string]map[string]string, error) {
	sQuery := makeSimpleQuery(options)
	tx, err := me.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var query storm.Query
	m := make(map[string]map[string]string)
	if keyContains != "" {
		keyContains = containsCaseInsensitiveReg(keyContains)
		query = tx.Select(q.Re("Code", keyContains))
	} else if keyContains != "" && valueContains != "" {
		keyContains = containsCaseInsensitiveReg(keyContains)
		valueContains = containsCaseInsensitiveReg(valueContains)
		query = tx.Select(q.Or(q.Re("Code", keyContains), q.Re("Text", valueContains)))
	} else if valueContains != "" {
		valueContains = regexp.QuoteMeta(valueContains)
		query = tx.Select(q.Re("Text", valueContains))
	} else {
		query = tx.Select()
	}
	err = query.
		Limit(sQuery.limit).
		Skip(sQuery.index).
		OrderBy("Code", "Text").
		Each(new(i18nInternal), func(record interface{}) error {
			item := record.(*i18nInternal)
			if m[item.Code] == nil {
				m[item.Code] = make(map[string]string)
			}
			m[item.Code][item.Lang] = item.Text
			return nil
		})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (me *I18nDB) Get(lang string, key string, args ...string) (string, error) {
	var i18nInt i18nInternal
	err := me.db.One("ID", lang+"_"+key, &i18nInt)
	if err != nil {
		return key, nil
	}
	return me.resolver.Resolve(i18nInt.Text, args...), nil
}

func (me *I18nDB) GetAll(lang string) (map[string]string, error) {
	me.allCacheLock.RLock()
	if m, ok := me.allCache[lang]; ok {
		me.allCacheLock.RUnlock()
		return m, nil
	}
	me.allCacheLock.RUnlock()
	m := make(map[string]string)
	err := me.db.Select(q.Eq("Lang", lang)).Each(new(i18nInternal), func(i interface{}) error {
		item := i.(*i18nInternal)
		m[item.Code] = item.Text
		return nil
	})
	if err != nil {
		return nil, err
	}
	me.allCacheLock.Lock()
	me.allCache[lang] = m
	me.allCacheLock.Unlock()
	return m, nil
}

func (me *I18nDB) PutAll(lang string, translations map[string]string) error {
	for key, text := range translations {
		//prevent from any javascript or global css hacks
		translations[key] = i18n.Escape(text)
	}
	me.allCacheLock.Lock()
	if me.fallbackLang != lang {
		if _, ok := me.langs[lang]; !ok {
			me.allCacheLock.Unlock()
			return fmt.Errorf("lang does not exist")
		}
	}
	if m, ok := me.allCache[lang]; ok {
		for key, text := range translations {
			if key == "" {
				continue
			}
			m[key] = text
		}
	}
	me.allCacheLock.Unlock()

	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	tx = tx.WithBatch(true)
	defer tx.Rollback()
	for key, text := range translations {
		if key == "" {
			continue
		}
		err = tx.Save(&i18nInternal{
			ID:   lang + "_" + key,
			Lang: lang,
			Code: key,
			Text: text,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (me *I18nDB) Put(lang string, key string, text string) error {
	tx, err := me.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = me.put(&lang, &key, &text, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (me *I18nDB) put(lang *string, key *string, text *string, tx storm.Node) error {
	if lang == nil || key == nil || text == nil || len(*key) == 0 {
		return fmt.Errorf("invalid arguments: lang and key must be provided")
	}
	//prevent from any javascript or global css hacks
	*text = i18n.Escape(*text)
	me.allCacheLock.Lock()
	if me.fallbackLang != *lang {
		if _, ok := me.langs[*lang]; !ok {
			me.allCacheLock.Unlock()
			return fmt.Errorf("lang does not exist")
		}
	}
	if m, ok := me.allCache[*lang]; ok {
		m[*key] = *text
	}
	me.allCacheLock.Unlock()
	return tx.Save(&i18nInternal{
		ID:   *lang + "_" + *key,
		Lang: *lang,
		Code: *key,
		Text: *text,
	})
}

func (me *I18nDB) PutLang(lang string, enabled bool) error {
	if !me.langReg.MatchString(lang) {
		return os.ErrInvalid
	}
	me.allCacheLock.Lock()
	l := &model.Lang{Code: lang, Enabled: enabled}
	me.langs[lang] = l
	me.langSlice = make([]*model.Lang, len(me.langs))
	i := 0
	for _, lan := range me.langs {
		me.langSlice[i] = lan
		i++
	}
	sort.Slice(me.langSlice, func(i, j int) bool {
		return me.langSlice[i].Code < me.langSlice[j].Code
	})
	//insert lang code as key to make it translatable
	insertLangCodeAsKey := false
	if m, ok := me.allCache[me.fallbackLang]; ok {
		if _, exists := m[lang]; !exists {
			m[lang] = lang
			insertLangCodeAsKey = true
		}
	} else {
		me.allCache[me.fallbackLang] = map[string]string{lang: lang}
		insertLangCodeAsKey = true
	}
	me.allCacheLock.Unlock()
	if insertLangCodeAsKey {
		err := me.db.Save(&i18nInternal{
			ID:   me.fallbackLang + "_" + lang,
			Lang: me.fallbackLang,
			Code: lang,
			Text: lang,
		})
		if err != nil {
			return err
		}
	}
	return me.db.Save(l)
}

func (me *I18nDB) GetLangs(enabled bool) ([]*model.Lang, error) {
	var langs []*model.Lang
	err := me.db.Select(q.Eq("Enabled", enabled)).OrderBy("Code").Find(&langs)
	if err != nil {
		return nil, err
	}
	return langs, nil
}

func (me *I18nDB) HasLang(lang string) bool {
	me.allCacheLock.RLock()
	var exists bool
	if me.langs != nil {
		_, exists = me.langs[lang]
	}
	me.allCacheLock.RUnlock()
	return exists
}

func (me *I18nDB) GetAllLangs() ([]*model.Lang, error) {
	me.allCacheLock.RLock()
	if len(me.langSlice) > 0 {
		me.allCacheLock.RUnlock()
		return me.langSlice, nil
	}
	me.allCacheLock.RUnlock()
	var langs []*model.Lang
	err := me.db.Select().OrderBy("Code").Find(&langs)
	if err != nil {
		return nil, err
	}
	return langs, nil
}

const fallback = "fallback"

func (me *I18nDB) PutFallback(lang string) error {
	if !me.langReg.MatchString(lang) {
		return os.ErrInvalid
	}
	me.fallbackLang = lang
	return me.db.Set(fallback, fallback, lang)
}

const imexI18n = "I18n"

//to provide an error for the entire region
const imexAllKeysMarker = "_all"

func (me *I18nDB) Import(imex storage.ImexIF) error {
	var err error
	err = me.init(imex)
	if err != nil {
		return err
	}
	var langs []*model.Lang
	langs, err = imex.DB().I18n.GetAllLangs()
	if err == nil && len(langs) > 0 {
		for _, lang := range langs {
			if imex.SkipExistingOnImport() {
				if !imex.SysDB().I18n.HasLang(lang.Code) {
					err = imex.SysDB().I18n.PutLang(lang.Code, lang.Enabled)
					imex.ProcessedEntry(imexI18n, lang.Code, err)
				}
			} else {
				err = imex.SysDB().I18n.PutLang(lang.Code, lang.Enabled)
				imex.ProcessedEntry(imexI18n, lang.Code, err)
			}
			var trans map[string]string
			trans, err = imex.DB().I18n.GetAll(lang.Code)
			if err != nil {
				imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
				continue
			}
			if imex.SkipExistingOnImport() {
				var transExisting map[string]string
				transExisting, err = imex.SysDB().I18n.GetAll(lang.Code)
				if err != nil {
					imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
					continue
				}
				tx, err := imex.SysDB().I18n.(*I18nDB).db.Begin(true)
				if err != nil {
					return err
				}
				for k, v := range trans {
					if _, exists := transExisting[k]; exists {
						continue
					} else {
						err = imex.SysDB().I18n.(*I18nDB).put(&lang.Code, &k, &v, tx)
						imex.ProcessedEntry(imexI18n, lang.Code+"_"+k, err)
					}
				}
				if err != nil {
					_ = tx.Rollback()
					imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
				} else {
					err = tx.Commit()
					if err != nil {
						_ = tx.Rollback()
						imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
					}
				}
			} else {
				err = imex.SysDB().I18n.PutAll(lang.Code, trans)
				for k, t := range trans {
					if k == "" {
						imex.ProcessedEntry(imexI18n, lang.Code+"_"+k, fmt.Errorf("invalid argument: key can not be empty {key:'%s',text:'%s'}", k, t))
						continue
					}
					imex.ProcessedEntry(imexI18n, lang.Code+"_"+k, err)
				}
			}
		}
	}
	return nil
}

func (me *I18nDB) init(imex storage.ImexIF) error {
	var err error
	if imex.DB().I18n == nil {
		imex.DB().I18n, err = NewI18nDB(imex.Dir())
	}
	return err
}

func (me *I18nDB) Export(imex storage.ImexIF, id ...string) error {
	var err error
	err = me.init(imex)
	if err != nil {
		return err
	}
	var langs []*model.Lang

	specificIds := len(id) > 0
	langs, err = imex.SysDB().I18n.GetAllLangs()
	if err == nil && len(langs) > 0 {
		for _, lang := range langs {
			if !imex.IsProcessed(imexI18n, lang.Code) {
				err = imex.DB().I18n.PutLang(lang.Code, lang.Enabled)
				imex.ProcessedEntry(imexI18n, lang.Code, err)
			}
			if specificIds {
				var trans map[string]string
				trans, err = imex.SysDB().I18n.GetAll(lang.Code)
				if err != nil {
					imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
					continue
				}
				tx, err := imex.DB().I18n.(*I18nDB).db.Begin(true)
				if err != nil {
					return err
				}
				for _, v := range id {
					if !imex.IsProcessed(imexI18n, lang.Code+"_"+v) {
						if text, ok := trans[v]; ok {
							err = imex.DB().I18n.(*I18nDB).put(&lang.Code, &v, &text, tx)
							imex.ProcessedEntry(imexI18n, lang.Code+"_"+v, err)
						}
					}
				}
				if err != nil {
					_ = tx.Rollback()
					imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
				} else {
					err = tx.Commit()
					if err != nil {
						_ = tx.Rollback()
						imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
					}
				}
			} else if !imex.IsProcessed(imexI18n, lang.Code+imexAllKeysMarker) {
				var trans map[string]string
				trans, err = imex.SysDB().I18n.GetAll(lang.Code)
				if err != nil {
					imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
					continue
				}
				err = imex.DB().I18n.PutAll(lang.Code, trans)
				if err != nil {
					imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, err)
					continue
				}
				for k := range trans {
					imex.ProcessedEntry(imexI18n, lang.Code+"_"+k, nil)
				}
				imex.ProcessedEntry(imexI18n, lang.Code+imexAllKeysMarker, nil)
			}
		}
	}
	return nil
}

func (me *I18nDB) GetFallback() (string, error) {
	if me.fallbackLang != "" {
		return me.fallbackLang, nil
	}
	var l string
	err := me.db.Get(fallback, fallback, &l)
	if err != nil {
		return "en", nil
	}
	return l, nil
}

func (me *I18nDB) Close() error {
	return me.db.Close()
}
