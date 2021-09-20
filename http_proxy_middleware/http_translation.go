package http_proxy_middleware

import (
	"fmt"
	"github.com/didi/gatekeeper/public"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
			val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})
			val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return matched
			})
			val.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
				matched, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
				return matched
			})
			val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, ms := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(ms, " ")) != 2 {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, ms := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(ms, " ")) != 3 {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_ipportlist", func(fl validator.FieldLevel) bool {
				for _, ms := range strings.Split(fl.Field().String(), ",") {
					if matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !matched {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_iplist", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, item := range strings.Split(fl.Field().String(), ",") {
					matched, _ := regexp.Match(`\S+`, []byte(item)) //ip_addr
					if !matched {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_weightlist", func(fl validator.FieldLevel) bool {
				fmt.Println(fl.Field().String())
				for _, ms := range strings.Split(fl.Field().String(), ",") {
					if matched, _ := regexp.Match(`^\d+$`, []byte(ms)); !matched {
						return false
					}
				}
				return true
			})

			//自定义翻译器
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			val.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
				return ut.Add("valid_username", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_username", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_service_name", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_rule", trans, func(ut ut.Translator) error {
				return ut.Add("valid_rule", "{0} 必须是非空字符", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_rule", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
				return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_url_rewrite", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_header_transfor", trans, func(ut ut.Translator) error {
				return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_header_transfor", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_ipportlist", trans, func(ut ut.Translator) error {
				return ut.Add("valid_ipportlist", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_ipportlist", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_iplist", trans, func(ut ut.Translator) error {
				return ut.Add("valid_iplist", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_iplist", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_weightlist", trans, func(ut ut.Translator) error {
				return ut.Add("valid_weightlist", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_weightlist", fe.Field())
				return t
			})
			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}
