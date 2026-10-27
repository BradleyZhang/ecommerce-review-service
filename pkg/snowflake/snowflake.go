package snowflake

import (
	"errors"
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var (
	ErrInvalidInitParam  = errors.New("snowflake 初始化失败，无效的 startTime 或 machineID")
	ErrInvalidTimeFormat = errors.New("snowflake 初始化失败，无效的 startTime 格式")
)
var node *sf.Node

func Init(startTime string, machineID int64) error {
	if len(startTime) == 0 || machineID <= 0 {
		return ErrInvalidInitParam
	}
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return ErrInvalidTimeFormat
	}
	sf.Epoch = st.UnixNano() / 1000_000
	node, err = sf.NewNode(machineID)
	if err != nil {
		return fmt.Errorf("snowflake err:%w", err)
	}
	return nil
}
func GenID() int64 {
	return node.Generate().Int64()
}
