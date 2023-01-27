import React from 'react'

import {Alert} from '@grafana/ui'

export const QueryHelp = () => (
  <Alert title="Query Help" severity="info">
    <div>
      <h3>Format Options:</h3>
      <h4>Table(default):</h4>
      <p>Return any set of columns</p>
      <h3>Time series:</h3>
      <p>Return column named time (UTC in seconds or timestamp) return column(s) with numeric datatype as values</p>
      <p>Result sets of time series queries need to be sorted by time.</p>
    </div>
    <div>
      <h3>Supported Macros:</h3>
      <li>$__dateBin(time) -&gt; date_bin(interval '30 second', time, interval '1970-01-01T00:00:00Z')</li>
      <li>
        $__dateBinAlias(time) -&gt; date_bin(interval '30 second', time, interval '1970-01-01T00:00:00Z') as time_binned
      </li>
      <li>$__interval -&gt; interval '30 second'</li>
      <li>$__timeFilter() -&gt; time &gt;= '2023-01-26T16:24:39Z' AND time &lt;= '2023-01-26T17:24:39Z'</li>
      <li>$__timeFrom -&gt; cast('2023-01-01T00:00:00Z' as timestamp)</li>
      <li>
        $__timeGroup(time, hour) -&gt; datepart('minute', time),datepart('hour', time),datepart('day',
        time),datepart('month', time),datepart('year', time);
      </li>
      <li>
        $__timeGroupAlias(time, minute) -&gt; datepart('minute', time) as time_minute,datepart('hour', time) as
        time_hour,datepart('day', time) as time_day,datepart('month', time) as time_month, datepart('year', time) as
        time_year
      </li>
      <li>$__timeRange -&gt; time &gt;= '2023-01 01T00:00:00Z' and time &lt;= '2023-01-01T01:00:00Z'</li>
      <li>$__timeRangeFrom(time) -&gt; time &gt;= '2023-01-01T00:00:00Z'</li>
      <li> $__timeRangeTo(time) -&gt; time &lt;= '2023-01-01T01:00:00Z'</li>
      <li>$__timeTo(time) -&gt; cast(time as timestamp)</li>
    </div>
  </Alert>
)
