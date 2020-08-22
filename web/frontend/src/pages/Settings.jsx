import React, { useEffect, useRef, useState } from 'react'
import { Alert, Avatar, Button, Divider, Input, List, Select, Spin, Tabs, Tag, Typography } from 'antd'
import {
  autoCompleteProfileNames,
  getAdministratorsProfiles,
  getHomePage,
  getMembersRules,
  isSuperAdmin,
  setAdministratorsProfiles,
  setHomePage,
  setMembersRules
} from '../api'
import { debounce } from '../util'
import { Link } from 'react-router-dom'
import ReactQuill from 'react-quill'
import 'react-quill/dist/quill.snow.css'
import SuperAdminOnly from '../components/SuperAdminOnly'

const { Option } = Select
const { Title, Text } = Typography
const { TextArea } = Input
const { TabPane } = Tabs

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

  const [membersRule, setMembersRule] = useState('')
  const [homePageContent, setHomePageContent] = useState('')

  useEffect(() => {
    if (isSuperAdmin()) {
      getAdministratorsProfiles().then(data => {
        setCurrentAdmins(data)
      })
    }
  }, [])

  useEffect(() => {
    getMembersRules().then(data => {
      setMembersRule(data.value)
    })
  }, [])

  useEffect(() => {
    getHomePage().then(data => {
      setHomePageContent(data.value)
    })
  }, [])

  const [saving, setSaving] = useState(false)
  const saveAdminsAndMembersRules = () => {
    setSaving(true)

    const promises = []
    promises.push(setMembersRules(membersRule))

    if (isSuperAdmin()) {
      promises.push(setAdministratorsProfiles(currentAdmins.map(user => user.id)))
    }

    Promise.all(promises).finally(() => {
      setSaving(false)
    })
  }

  const saveHomePage = () => {
    setSaving(true)
    setHomePage(homePageContent).finally(() => {
      setSaving(false)
    })
  }

  const modules = {
    'toolbar': [
      [{ 'header': [1, 2, 3, 4, 5, 6] }],
      ['bold', 'italic', 'underline', 'strike'],
      [{ 'color': [] }, { 'background': [] }],
      [{ 'script': 'super' }, { 'script': 'sub' }],
      ['blockquote'],
      [{ 'list': 'ordered' }, { 'list': 'bullet' }, { 'indent': '-1' }, { 'indent': '+1' }],
      [{ 'align': [] }],
      ['link', 'image', 'video'],
      ['clean']
    ]
  }

  const formats = [
    'header',
    'bold', 'italic', 'underline', 'strike', 'blockquote',
    'color', 'background', 'script',
    'list', 'bullet', 'indent', 'align',
    'link', 'image', 'video'
  ]

  return (<div className="settings-page">
    <Spin tip="Salvando..." spinning={saving}>
      <Title level={4}>
        Configurações
      </Title>
      <Tabs type="card">
        <TabPane tab="Membros" key="1">
          <SuperAdminOnly>
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
          </SuperAdminOnly>

          <Divider plain>Regras de acesso membros</Divider>
          <Alert
            message="Como funciona as regras?"
            description={
              <div>
                As regras seguem o padrão de cima para baixo, então são aplicadas conforme a ordem.<br/>
                <Text code> </Text>- (em branco) bloqueado todos membros<br/>
                <Text code>* </Text>- Permitir todos os membros<br/>
                <Text code>admin </Text>- Permitir todos os administradores<br/>
                <Text code>789123 </Text>- Permitido o membro com ID 789123<br/>
                <Text code>!123456 </Text>- Bloqueado o membro com ID 123456<br/>
                Com isso é possível fazer regras como:
                <TextArea value={'*\n!123456'} readOnly={true}/>
                Que permitem todos os membros, menos o membro com ID 123456 (uma regra por linha).<br/>
                Super administradores sempre podem acessar o sistema.
              </div>
            }
            type="info"
            showIcon
          />

          <br/>
          <TextArea
            rows={4}
            onChange={(e) => setMembersRule(e.target.value)}
            value={membersRule}
          />
          <br/>
          <br/>
          <Button type="primary" onClick={saveAdminsAndMembersRules}>
            Salvar
          </Button>
        </TabPane>
        <TabPane tab="Home" key="2">
          <ReactQuill
            theme="snow"
            value={homePageContent}
            modules={modules}
            formats={formats}
            onChange={value => setHomePageContent(value)}/>
          <br/>
          <br/>
          <Button type="primary" onClick={saveHomePage}>
            Salvar
          </Button>
        </TabPane>
      </Tabs>
    </Spin>
  </div>)
}
export default Settings
