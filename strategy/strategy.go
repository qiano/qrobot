package strategy

import (
	"fmt"
	// . "github.com/qshuai162/qian/common"
	// . "github.com/qshuai162/qian/config"
	// "github.com/qshuai162/qian/db"
	// "github.com/qshuai162/qian/email"
	// "github.com/qshuai162/qian/huobi"
	"github.com/qshuai162/qian/common/logger"
	"strconv"
	// "time"
)

// Strategy is the interface that must be implemented by a strategy driver.
type Strategy interface {
	Tick() bool
}

var strategyMaps = make(map[string]Strategy)

// Register makes a strategy available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(strategyName string, strategy Strategy) {

	if strategy == nil {
		panic("sql: Register strategy is nil")
	}
	if _, dup := strategyMaps[strategyName]; dup {
		panic("sql: Register called twice for strategy " + strategyName)
	}
	strategyMaps[strategyName] = strategy
	logger.Infoln("注册策略：" + strategyName)
}

func getStrategy(name string) Strategy {
	strategy, ok := strategyMaps[name]
	if !ok {
		logger.Errorf("sql: 未知的策略名 %q (forgotten import? private strategy?)", name)
		return nil
	} else {
		return strategy
	}
}

func Tick(strategyName string) bool {
	return getStrategy(strategyName).Tick()
}

func toString(s interface{}) string {
	if v, ok := s.(string); ok {
		return v
	}
	return fmt.Sprintf("%v", s)
}

func toFloat(s interface{}) float64 {
	var ret float64
	switch v := s.(type) {
	case float64:
		ret = v
	case int64:
		ret = float64(v)
	case string:
		ret, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Errorln("convert ", s, " to float failed")
			return ret
		}
	}
	return ret
}

func float2str(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}
