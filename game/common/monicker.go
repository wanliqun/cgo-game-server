package common

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/interfaces"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/wanliqun/cgo-game-server/proto"
)

type MonickerGenerator interface {
	Generate(proto.Sex, proto.Culture) string
}

type GoFakerNameGenerator struct{}

func (g *GoFakerNameGenerator) Generate(sex proto.Sex, culture proto.Culture) string {
	var opt options.OptionFunc
	switch culture {
	case proto.Culture_RUSSIAN:
		opt = options.WithStringLanguage(interfaces.LangRUS)
	case proto.Culture_CHINESE:
		opt = options.WithStringLanguage(interfaces.LangCHI)
	default:
		opt = options.WithStringLanguage(interfaces.LangENG)
	}

	switch sex {
	case proto.Sex_MALE:
		return fmt.Sprintf("%s %s", faker.FirstNameMale(opt), faker.LastName(opt))
	case proto.Sex_FEMALE:
		return fmt.Sprintf("%s %s", faker.FirstNameFemale(opt), faker.LastName(opt))
	default:
		return faker.Name(opt)
	}
}
