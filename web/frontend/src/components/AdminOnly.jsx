import React from 'react'
import { isAdmin, isSuperAdmin } from '../api'

export const adminOnly = input => {
  if (!isAdmin()) {
    return
  }

  return input
}

export default props => {
  if (!isAdmin()) {
    return <></>
  }

  return <template>
    {props.children}
  </template>
}
