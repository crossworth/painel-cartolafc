import React from 'react'

import { Empty } from 'antd'

export default () => {
  return (<div>
    <Empty description={
      <span>
        Nenhum perfil encontrado
      </span>
    }/>
  </div>)
}
