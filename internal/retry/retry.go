package retry

import (
	"errors"
	"fmt"
	"github.com/anatoly32322/metriccollector/internal/logger"
	"reflect"
	"time"
)

func Retry(
	fn interface{}, args []interface{}, maxRetry int,
	startBackoff, maxBackoff time.Duration) ([]reflect.Value, error) {

	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		return nil, errors.New("retry: function type required")
	}

	argVals := make([]reflect.Value, len(args))
	for i, arg := range args {
		argVals[i] = reflect.ValueOf(arg)
	}

	for attempt := 0; attempt < maxRetry; attempt++ {
		result := fnVal.Call(argVals)
		errVal := result[len(result)-1]

		if errVal.IsNil() {
			return result, nil
		}
		time.Sleep(time.Second * startBackoff)
		if startBackoff < maxBackoff {
			startBackoff += 2
		}
		logger.Sugar.Infof(
			"Retrying function call, attempt: %d, error: %v\n",
			attempt+1, errVal,
		)
	}
	return nil, fmt.Errorf("retry: max retries reached without success")
}
