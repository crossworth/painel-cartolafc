import React, { useEffect, useState } from 'react'
import { AutoComplete, Col, Divider, List, Row, Select, Tag, Typography } from 'antd'
import { autoCompleteProfileNames, getAdministratorsProfiles } from '../api'
import { debounce } from '../util'

const { Option } = Select
const { Title } = Typography

const debounceFn = debounce(fn => {
  fn()
}, 300)

const Settings = () => {

  const [admins, setAdmins] = useState([])
  const [currentAdmins, setCurrentAdmins] = useState([])
  const [users, setUsers] = useState([])
  const [loadingUsers, setLoadingUsers] = useState(false)

  const onSearch = value => {
    if (value === '') {
      return
    }

    setUsers([])
    setLoadingUsers(true)

    debounceFn(() => {
      autoCompleteProfileNames(value).then(data => {
        setUsers(data)
      }).catch(err => {
        setUsers([])
      }).finally(() => {
        setLoadingUsers(false)
      })
    })
  }

  const handleChange = value => {
    setAdmins(value)
  }

  const tagRender = props => {
    const { label, value, closable, onClose } = props

    return (
      <Tag value={value} closable={closable} onClose={onClose}
           style={{ marginRight: 1, marginLeft: 1, marginTop: 1, marginBottom: 1 }}>
        {label}
      </Tag>
    )
  }

  useEffect(() => {
    getAdministratorsProfiles().then(data => {
      setCurrentAdmins(data)
    })
  }, [])

  return (<div>
    <Title level={4}>
      Configurações
    </Title>
    <Divider plain>Membros com acesso administrador</Divider>

    <List
      size="small"
      header={<div>Administradores atuais</div>}
      bordered
      dataSource={currentAdmins}
      renderItem={user => <List.Item>
        <Row>
          <Col flex="auto">
            <img src={user.photo} alt={user.screen_name} width="40" style={{ marginRight: 5 }}/>
            {user.first_name} {user.last_name} (@{user.screen_name})
          </Col>
          <Col flex="60px">
            Delete
          </Col>
        </Row>
      </List.Item>}
    />

    <br/>
    <strong>Adicionar novo</strong>
    <br/>
    <AutoComplete
      onSearch={onSearch}
      style={{ width: '100%' }}
      placeholder="Digite um nome">
      {users.map(user => (
        <Option key={user.id} value={user.screen_name}>
          <Row>
            <Col flex="60px">
              <img width="50" height="50" src={user.photo} alt=''/>
            </Col>
            <Col flex="auto">
              {`${user.first_name} ${user.last_name} (@${user.screen_name})`}
            </Col>
          </Row>
        </Option>
      ))}
    </AutoComplete>

    {/*<Select*/}
    {/*  mode="multiple"*/}
    {/*  style={{ width: '100%' }}*/}
    {/*  value={admins}*/}
    {/*  placeholder="Selecione os membros"*/}
    {/*  notFoundContent={loadingUsers ? <Spin size="small"/> : <div>Nenhum resultado</div>}*/}
    {/*  onSearch={onSearch}*/}
    {/*  filterOption={false}*/}
    {/*  tagRender={tagRender}*/}
    {/*  onChange={handleChange}>*/}
    {/*  {users.map(user => (*/}
    {/*    <Option key={user.id}>*/}
    {/*      <img src={user.photo} alt={user.screen_name} width="40" style={{ marginRight: 5 }}/>*/}
    {/*      {user.first_name} {user.last_name} (@{user.screen_name})*/}
    {/*    </Option>*/}
    {/*  ))}*/}
    {/*</Select>*/}
  </div>)
}
export default Settings
