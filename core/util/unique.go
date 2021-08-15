package util

import (
	"hash/crc32"
	"strings"

	"github.com/twinj/uuid"
)

// NewUUID 创建UUID唯一编号，去分隔符
func NewUUID() string {
	u := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}

// NewSequenceID 创建数字唯一编号，有重复几率
func NewSequenceID() int64 {
	id := int64(1)
	u := NewUUID()
	v := int(crc32.ChecksumIEEE([]byte(u)))
	if v >= 0 {
		id += int64(v)
	}
	if -v >= 0 {
		id += int64(-v)
	}
	return id
}
