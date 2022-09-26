package validator

import (
	"fmt"
	"ginblogtest/routes/errmsg"
	"reflect"

	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

func Validate(data interface{}) (string, int) {
	//中文翻译器
	uni := ut.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN")
	//校验器
	validate := validator.New()
	//将翻译器注册到校验器中
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err:", err)
	}
	//将struct tag里添加的中文名 作为备用名。
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get("label")
		return label
	})
	//将错误信息转换为指定语言
	err = validate.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			// fmt.Println(err.Translate(trans)) //输出中文错误
			return v.Translate(trans), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCSE
}
