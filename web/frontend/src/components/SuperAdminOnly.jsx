import React from 'react'
import { isSuperAdmin } from '../api'

export const superAdminOnly = input => {
  if (!isSuperAdmin()) {
    return
  }

  return input
}

export default props => {
  if (!isSuperAdmin()) {
    return <></>
  }

  return props.children
}
