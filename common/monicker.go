package common

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/interfaces"
	"github.com/go-faker/faker/v4/pkg/options"
)

type Sex = int32

const (
	Male Sex = iota
	Female
)

type Culture = int32

const (
	AMERICAN Culture = iota
	ARGENTINIAN
	AUSTRALIAN
	BRAZILIAN
	BRITISH
	BULGARIAN
	CANADIAN
	CHINESE
	DANISH
	FINNISH
	FRENCH
	GERMAN
	KAZAKH
	MEXICAN
	NORWEGIAN
	POLISH
	PORTUGUESE
	RUSSIAN
	SPANISH
	SWEDISH
	TURKISH
	UKRAINIAN
)

type MonickerGenerator interface {
	Generate(Sex, Culture) string
}

type GoFakerNameGenerator struct{}

func (g *GoFakerNameGenerator) Generate(sex Sex, culture Culture) string {
	var opt options.OptionFunc
	switch culture {
	case RUSSIAN:
		opt = options.WithStringLanguage(interfaces.LangRUS)
	case CHINESE:
		opt = options.WithStringLanguage(interfaces.LangCHI)
	default:
		opt = options.WithStringLanguage(interfaces.LangENG)
	}

	switch sex {
	case Male:
		return fmt.Sprintf("%s %s", faker.FirstNameMale(opt), faker.LastName(opt))
	case Female:
		return fmt.Sprintf("%s %s", faker.FirstNameFemale(opt), faker.LastName(opt))
	default:
		return faker.Name(opt)
	}
}
