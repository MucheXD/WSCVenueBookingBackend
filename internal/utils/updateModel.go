package utils

import "reflect"

// UpdateField 如果 update 非零，则将其值覆盖到 target 中
// target 必须是指向某个类型的指针，update 必须是该类型的值
func UpdateField[T comparable](target *T, update T) {
	var zero T // 声明一个类型为 T 的零值变量（创建时自动具有类型零值）
	if update != zero {
		*target = update
	}
}

// UpdateFieldPtr 如果 update 非 nil，则将其值覆盖到 target 中
// 适用于区分“未传值”与“传零值”的更新场景
func UpdateFieldPtr[T any](target *T, update *T) {
	if update != nil {
		*target = *update
	}
}

// [谨慎使用] 可能引发潜在问题，建议使用 UpdateField 来更新单个字段
// UpdateModel 将 src 结构体中非零值的字段覆盖到 dest 结构体中
// 操作粒度为一层，不进行递归查询
// dest 和 src 必须是指向同一类型结构体的指针
func UpdateModel(dest any, src any) {
	// 获取指针指向的实际变量值
	// .Elem() 用于“解引用”，获取指针指向的结构体实体
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	// 开始遍历 src 结构体的所有字段
	// .NumField() 返回该结构体中定义的字段总数
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)
		// 检查 src 字段是否为其类型的“零值”
		// .IsZero() 会判断字段是否为“零值”，即 go 语境下的“未定义”
		// 如果 !IsZero() 为真，说明 src 中这个字段有实际的数据，覆盖到 dest
		if !srcField.IsZero() {
			// 检查 dest 的该字段是否可以被修改（导出字段/Public 字段才可修改）
			// 如果字段名是小写开头的（private），CanSet() 会返回 false
			if destField.CanSet() {
				destField.Set(srcField)
			}
		}
	}
}
