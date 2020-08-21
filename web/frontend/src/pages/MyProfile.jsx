import React, { useEffect, useState } from 'react'

import { Col, Divider, Row, Spin, Statistic, Typography } from 'antd'
import { normalizeScreenName } from '../util'
import { getMyProfileStats } from '../api'

const { Title } = Typography

const MyProfile = () => {

  const [user, setUser] = useState({})
  const [userStats, setUserStats] = useState({})

  useEffect(() => {
    getMyProfileStats().then(values => {
      setUser(values.user)
      setUserStats(values.stats)
    }).catch(err => {

    })
  }, [])

  return (
    <div>
      <Spin tip="Carregando..." spinning={!user.first_name}>
        <Title level={4}>
          {
            user.first_name ?
              <div><a href={`https://vk.com/id${user.id}`} target="_blank"
                      rel="noopener noreferrer">{user.first_name} {user.last_name}</a> (@{normalizeScreenName(user.screen_name, user.id)})
                -
                ID {user.id}
              </div>
              : <div>Carregando dados</div>
          }
        </Title>

        <Divider/>
        <Row gutter={16}>
          <Col md={6}>
            <Statistic title="Tópicos" value={userStats.total_topics ? userStats.total_topics : 0}/>
          </Col>
          <Col md={6}>
            <Statistic title="Comentários" value={userStats.total_comments ? userStats.total_comments : 0}/>
          </Col>
          <Col md={6}>
            <Statistic title="Likes" value={userStats.total_likes ? userStats.total_likes : 0}/>
          </Col>
          <Col md={6}>
            <Statistic title="Tópicos + Comentários"
                       value={userStats.total_topics_plus_comments ? userStats.total_topics_plus_comments : 0}/>
          </Col>
        </Row>
      </Spin>
    </div>
  )
}

export default MyProfile
