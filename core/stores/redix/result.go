package redix

// Result 查询单个结果
type Result struct {
	Value interface{}
	Error error
}

// NewResult 新建查询结果
func NewResult(value interface{}, err error) *Result {
	return &Result{Value: value, Error: err}
}

// Result 返回查询结果的值,defaultValue 表示默认值
func (r *Result) Result(defaultValue ...interface{}) interface{} {
	if r.Error != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		panic(r.Error)
	}
	return r.Value
}

// MResult 列表查询结果
type MResult struct {
	Value []interface{}
	Error error
}

func NewMResult(value []interface{}, err error) *MResult {
	return &MResult{Value: value, Error: err}
}

func (m *MResult) Result(defaultValue ...[]interface{}) []interface{} {
	if m.Error != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		panic(m.Error)
	}
	return m.Value
}

func (m *MResult) Iterator() *Iterator {
	return NewIterator(m.Value)
}
