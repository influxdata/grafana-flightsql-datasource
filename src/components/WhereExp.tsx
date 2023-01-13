import React from 'react'

import {Input} from '@grafana/ui'

export const WhereExp = ({index, onChange, whereValue, setWhere}: any) => (
  <Input
    autoFocus
    type="text"
    spellCheck={false}
    onBlur={() => {}}
    onKeyDown={(e: any) => {
      if (e.key === 'Enter') {
        onChange(whereValue)
      }
    }}
    onChange={(e: any) => {
      setWhere({value: e.currentTarget.value, index: index})
    }}
    value={whereValue}
    width={15}
    placeholder="value = value"
  />
)
