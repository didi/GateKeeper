package public

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	if err := c.ShouldBind(params); err != nil {
		return err
	}
	//获取验证器
	valid, err := GetValidator(c)
	if err != nil {
		return err
	}
	//获取翻译器
	trans, err := GetTranslation(c)
	if err != nil {
		return err
	}
	err = valid.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	return nil
}

func GetValidator(c *gin.Context) (*validator.Validate, error) {
	val, ok := c.Get(ValidatorKey)
	if !ok {
		return nil, errors.New("未设置验证器")
	}
	validator, ok := val.(*validator.Validate)
	if !ok {
		return nil, errors.New("获取验证器失败")
	}
	return validator, nil
}

func GetTranslation(c *gin.Context) (ut.Translator, error) {
	trans, ok := c.Get(TranslatorKey)
	if !ok {
		return nil, errors.New("未设置翻译器")
	}
	translator, ok := trans.(ut.Translator)
	if !ok {
		return nil, errors.New("获取翻译器失败")
	}
	return translator, nil
}

func ServiceNameValidate(serviceName string) error {
	if serviceName == "" {
		return errors.New("服务名称不能为空")
	} else {
		reg, _ := regexp.MatchString(`^[0-9a-zA-Z_]{1,}$`, serviceName)
		if !reg { //解释失败，返回false
			return errors.New("服务名称格式错误")
		}
	}
	return nil
}

func HTTPPathsValidate(httpPaths string) error {
	if httpPaths == "" {
		return errors.New("服务地址不能为空")
	} else {
		reg, _ := regexp.MatchString(`^(/[\w\-]+)+`, httpPaths)
		if !reg { //解释失败，返回false
			return errors.New("服务地址格式错误")
		}
	}
	return nil
}

func UpstreamListValidate(upstreamList string) error {
	if upstreamList == "" {
		return errors.New("下游服务器ip和权重不能为空")
	} else {
		tmpLine := strings.Split(upstreamList, "\n")
		for _, tmp := range tmpLine {
			r, _ := regexp.Compile(`^(\S*\:\/\/)(\S*?)\s(.*?)$`)
			submatch := r.FindStringSubmatch(tmp)
			if len(submatch) != 4 {
				return errors.New("下游服务器ip和权重 format error")
			}
		}
	}
	return nil
}
