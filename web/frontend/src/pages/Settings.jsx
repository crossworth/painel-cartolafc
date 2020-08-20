import React from 'react'
import { Divider, Select, Typography } from 'antd'

const { Option } = Select
const { Title } = Typography

export default () => {

  const children = []
  for (let i = 10; i < 36; i++) {
    children.push(<Option key={i.toString(36) + i}>{i.toString(36) + i}</Option>)
  }

  const handleChange = () => {

  }

  return (<div>
    <Title level={4}>
      Configurações
    </Title>
    <Divider plain>Membros com acesso administrador</Divider>
    <Select
      mode="multiple"
      style={{ width: '100%' }}
      placeholder="Selecione os membros"
      defaultValue={['a10', 'c12']}
      onChange={handleChange}>
      {children}
    </Select>
  </div>)
}



