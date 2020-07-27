import React, { useEffect, useState } from 'react'

import { Alert, Avatar, Button, Col, Divider, List, Row, Spin, Statistic, Typography } from 'antd'
import { getProfileInfo, getProfileHistory, getProfileStats } from '../api'
import { timeStampToDate } from '../util'

const { Title } = Typography

const Profile = (props) => {
  const id = props.match.params.id

  const { history } = props

  const isPossibleScreenName = !(/^-?\d+$/.test(id))

  const [user, setUser] = useState({})
  const [userStats, setUserStats] = useState({})
  const [userProfileHistory, setUserProfileHistory] = useState([])

  useEffect(() => {
    if (isPossibleScreenName) {
      return
    }

    Promise.all([getProfileInfo(id), getProfileStats(id), getProfileHistory(id)]).then(values => {
      setUser(values[0])
      setUserStats(values[1])
      setUserProfileHistory(values[2])
    }).catch(err => {
      if (err.response && err.response.status && err.response.status === 404) {
        history.push('/perfil/nao-encontrado')
      }
    })
  }, [isPossibleScreenName, id, history])

  if (isPossibleScreenName) {
    history.push(`/resolver/${id}`)
    return <div/>
  }

  return (
    <div>
      <Spin tip="Carregando..." spinning={!user.first_name}>
        <Title level={4}>
          {
            user.first_name ?
              <div><a href={`https://vk.com/id${user.id}`} target="_blank"
                      rel="noopener noreferrer">{user.first_name} {user.last_name}</a> (@{user.screen_name}) -
                ID {user.id}
              </div>
              : <div>Carregando dados</div>
          }
        </Title>
        <Divider/>

        <Row gutter={16}>
          <Col span={6}>
            <Statistic title="Tópicos" value={userStats.total_topics ? userStats.total_topics : 0}/>
          </Col>
          <Col span={6}>
            <Statistic title="Comentários" value={userStats.total_comments ? userStats.total_comments : 0}/>
          </Col>
          <Col span={6}>
            <Statistic title="Likes" value={userStats.total_likes ? userStats.total_likes : 0}/>
          </Col>
          <Col span={6}>
            <Statistic title="Alterações do perfil"
                       value={userStats.total_profile_changes ? userStats.total_profile_changes : 0}/>
          </Col>
        </Row>

        <Divider/>

        <Button type="primary" rel="noopener noreferrer"
                href={`/perfil/${user.id}/topicos`}>
          Tópicos
        </Button>
        <Button style={{ marginLeft: 10 }} type="primary" rel="noopener noreferrer"
                href={`/perfil/${user.id}/comentarios`}>
          Comentários
        </Button>

        <Divider/>
        <Title level={4}>Histórico de alteração de perfil ({userProfileHistory.length})</Title>

        <Alert
          message="Porque existe registros duplicados?"
          description="O VK algumas vezes retorna o link da foto de outro servidor do CDN, isso efetivamente significa que o link é diferente, por isso é registrado como uma alteração de foto, mesmo sendo a mesma foto."
          type="info"
          showIcon
        />
        <br/>
        <List
          itemLayout="horizontal"
          dataSource={userProfileHistory}
          renderItem={item => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar src={item.photo}/>}
                title={`${item.first_name} ${item.last_name} - ${item.screen_name} em ${timeStampToDate(item.date)}`}
              />
            </List.Item>
          )}
        />
      </Spin>
    </div>
  )
}

export default Profile
