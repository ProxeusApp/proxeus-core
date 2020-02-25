package service

type (
	TranslateService interface {
	}
	DefaultTranslateService struct {
	}
)

func NewTranslateService() TranslateService {
	return &DefaultTranslateService{}
}

func (me *DefaultTranslateService) Trans(key string) {

}
