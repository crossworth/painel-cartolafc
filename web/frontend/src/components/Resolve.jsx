import React, { useEffect, useState } from 'react'
import { Alert, AutoComplete, Button, Col, Divider, message, Row, Spin, Typography } from 'antd'
import { autoCompleteProfileNames, resolveProfile } from '../api'
import { debounce } from '../util'

const { Title } = Typography
const { Option } = AutoComplete

const debounceFn = debounce(fn => {
  fn()
}, 300)

const Resolve = (props) => {
  const name = props.match.params.name

  const [loading, setLoading] = useState(false)
  const [found, setFound] = useState(false)

  const [options, setOptions] = useState([])
  const [selectedOption, setSelectedOption] = useState('')

  const { history } = props

  useEffect(() => {
    if (name === undefined) {
      setFound(true)
      return
    }

    setLoading(true)
    resolveProfile(name).then(result => {
      setFound(true)
      history.push(`/perfil/${result.id}`)
    }).catch(err => {
      setFound(false)
    }).finally(() => {
      setLoading(false)
    })
  }, [history, name])

  const onSearch = text => {
    debounceFn(() => {
      autoCompleteProfileNames(text).then(data => {
        setOptions(data)
      }).catch(err => {
        console.log(err)
      })
    })
  }

  const onClick = event => {
    event.preventDefault()

    if (selectedOption === '') {
      message.error('Você deve informar um nome ou link de perfil')
      return
    }

    setLoading(true)
    history.push(`/resolver/${selectedOption}`)
  }

  return (
    <div>
      <Spin tip="Carregando..." spinning={loading}>
        <Title level={4}>
          Resolver nome ou link de perfil
        </Title>

        <Alert
          message="O que é resolver um nome ou link?"
          description="É o processo de conversão de um nome/link de perfil ou screen name em um link de perfil canônico"
          type="info"
          showIcon
        />

        <Divider/>
        {
          !found && !loading &&
          <div>
            <Alert
              message="Não encontrado"
              description={`Não foi possível encontrar um perfil com relacionado a ${name}`}
              type="error"
              closable
            />
            <Divider/>
          </div>
        }

        <div style={{ textAlign: 'center' }}>
          <form onSubmit={onClick}>
            <AutoComplete
              disabled={loading}
              className="resolveInput"
              onSearch={onSearch}
              onChange={value => setSelectedOption(value)}
              placeholder="Digite um nome/link de perfil aqui">
              {options.map(profile => (
                <Option key={profile.id} value={profile.screen_name}>
                  <Row gutter={16}>
                    <Col span={1}>
                      <img width="50" height="50" src={profile.photo} alt={profile.screen_name}/>
                    </Col>
                    <Col span={16}>
                      {`${profile.first_name} ${profile.last_name} (@${profile.screen_name})`}
                    </Col>
                  </Row>
                </Option>
              ))}
            </AutoComplete>
            <Button type="submit" onClick={onClick} style={{ marginTop: 20 }} type="primary">Pesquisar</Button>
          </form>
        </div>
      </Spin>
    </div>
  )
}

export default Resolve
