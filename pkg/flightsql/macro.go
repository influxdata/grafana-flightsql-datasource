package flightsql

import (
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

var macros = sqlutil.Macros{
	"dateBin":      macroDateBin(""),
	"dateBinAlias": macroDateBin("_binned"),
	"interval":     macroInterval,
}

func macroInterval(query *sqlutil.Query, _ []string) (string, error) {
	return fmt.Sprintf("interval '%d second'", int64(query.Interval.Seconds())), nil
}

func macroFrom(query *sqlutil.Query, _ []string) (string, error) {
	return query.TimeRange.From.Format(time.RFC3339), nil
}

func macroTo(query *sqlutil.Query, _ []string) (string, error) {
	return query.TimeRange.To.Format(time.RFC3339), nil
}

func macroDateBin(suffix string) sqlutil.MacroFunc {
	return func(query *sqlutil.Query, args []string) (string, error) {
		if len(args) != 1 {
			return "", fmt.Errorf("%w: expected 1 argument, received %d", sqlutil.ErrorBadArgumentCount, len(args))
		}
		column := args[0]
		aliasing := func() string {
			if suffix == "" {
				return ""
			}
			return fmt.Sprintf(" as %s%s", column, suffix)
		}()
		return fmt.Sprintf("date_bin(interval '%d second', %s, timestamp '1970-01-01T00:00:00Z')%s", int64(query.Interval.Seconds()), column, aliasing), nil
	}
}
