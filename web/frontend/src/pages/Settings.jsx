import React, { useEffect, useRef, useState } from 'react'
import { Avatar, Button, Divider, List, Select, Spin, Tag, Typography } from 'antd'
import { autoCompleteProfileNames, getAdministratorsProfiles, setAdministratorsProfiles } from '../api'
import { debounce } from '../util'
import { Link } from 'react-router-dom'

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
        const currentAdminsIds = currentAdmins.map(admin => admin.id)
        setUsers(data.filter(el => {
          return currentAdminsIds.indexOf(el.id) === -1
        }))
      }).catch(err => {
        setUsers([])
      }).finally(() => {
        setLoadingUsers(false)
      })
    })
  }

  const [open, setOpen] = useState(false)
  const adminSelect = useRef(null)

  const handleChange = value => {
    value = value[0]

    users.forEach(user => {
      if (user.id === parseInt(value)) {
        setCurrentAdmins([...currentAdmins, user])
        return false
      }
    })

    setOpen(false)
    adminSelect.current.blur()
  }

  useEffect(() => {
    const currentAdminsIds = currentAdmins.map(admin => admin.id)
    setUsers(users.filter(el => {
      return currentAdminsIds.indexOf(el.id) === -1
    }))

  }, [currentAdmins])

  const deleteUser = userID => {
    setCurrentAdmins(currentAdmins.filter(user => user.id !== userID))
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

  const [saving, setSaving] = useState(false)
  const saveAdmins = () => {
    setSaving(true)
    setAdministratorsProfiles(currentAdmins.map(user => user.id)).finally(() => {
      setSaving(false)
    })
  }

  return (<div>
    <Spin tip="Salvando..." spinning={saving}>
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
          <List.Item.Meta
            avatar={
              <Avatar src={user.photo}/>}
            title={<Link to={`/perfil/${user.id}`}>{user.first_name} {user.last_name}</Link>}
            description={`@${user.screen_name}`}
          />
          <div><Button type="primary" onClick={() => deleteUser(user.id)} danger>
            Remover
          </Button></div>
        </List.Item>}
      />

      <br/>
      <strong>Adicionar novo</strong>
      <br/>
      <Select
        mode="multiple"
        style={{ width: '100%' }}
        value={admins}
        placeholder="Selecione os membros"
        notFoundContent={loadingUsers ? <Spin size="small"/> : <div>Nenhum resultado</div>}
        onSearch={onSearch}
        onChange={handleChange}
        filterOption={false}
        open={open}
        ref={adminSelect}
        onBlur={() => setOpen(false)}
        onFocus={() => setOpen(true)}
        tagRender={tagRender}>
        {users.map(user => (
          <Option key={user.id}>
            <img src={user.photo} alt={user.screen_name} width="40" style={{ marginRight: 5 }}/>
            {user.first_name} {user.last_name} (@{user.screen_name})
          </Option>
        ))}
      </Select>
      <br/>
      <br/>
      <Button type="primary" onClick={saveAdmins}>
        Salvar
      </Button>
    </Spin>
  </div>)
}
export default Settings
