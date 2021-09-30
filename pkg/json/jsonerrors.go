/**
 * Copyright (2021, ) Institute of Software, Chinese Academy of Sciences
 */
package json

import (
	"fmt"
	"errors"
)

/**
 *      author: wuheng@iscas.ac.cn
 *      date  : 2021/9/29
 */
var (
	IndexOutOfRangeError = errors.New("index out of range")
	ValueNotNumberError = errors.New("value is not number")
)

type KeyNotFoundError struct {
	Key string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("key[%s] not existed", e.Key)
}

type ValueTransformTypeError struct {
	Type string
}

func (e ValueTransformTypeError) Error() string {
	return fmt.Sprintf("cannot transform into %s", e.Type)
}
