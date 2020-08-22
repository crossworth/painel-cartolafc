import React, { useEffect, useState } from 'react'

import { Button, Col, Divider, List, Row, Spin, Statistic, Typography } from 'antd'
import { normalizeQuote, normalizeScreenName } from '../util'
import { getMyProfileStats } from '../api'
import { VK_GROUP_ID } from '../config'

const { Title } = Typography

const MyProfile = () => {

  const [user, setUser] = useState({})
  const [userStats, setUserStats] = useState({})
  const [topicWithMoreLikes, setTopicWithMoreLikes] = useState([])
  const [topicWithMoreComments, setTopicWithMoreComments] = useState([])
  const [commentsWithMoreLikes, setCommentsWithMoreLikes] = useState([])

  useEffect(() => {
    getMyProfileStats().then(values => {
      setUser(values.user)
      setUserStats(values.stats)
      setTopicWithMoreComments(values.topic_with_more_comments)
      setCommentsWithMoreLikes(values.comments_with_more_likes)
      setTopicWithMoreLikes(values.topic_with_more_likes)
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
        <Divider plain>O melhor de @{normalizeScreenName(user.screen_name, user.id)}</Divider>

        <List
          size="small"
          header={<strong>Tópicos com mais likes</strong>}
          bordered
          dataSource={topicWithMoreLikes}
          renderItem={topic => <List.Item>
            <List.Item.Meta
              title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${topic.id}`} target="_blank" rel="noopener noreferrer">{topic.title}</a>}
              description={`${topic.likes_count} likes`}
            />
          </List.Item>}
        />
        <br/>
        <List
          size="small"
          header={<strong>Tópicos com mais comentários</strong>}
          bordered
          dataSource={topicWithMoreComments}
          renderItem={topic => <List.Item>
            <List.Item.Meta
              title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${topic.id}`} target="_blank" rel="noopener noreferrer">{topic.title}</a>}
              description={`${topic.comments_count} comentários`}
            />
          </List.Item>}
        />

        <br/>
        <List
          size="small"
          header={<strong>Comentários com mais likes</strong>}
          bordered
          dataSource={commentsWithMoreLikes}
          renderItem={comment => <List.Item>
            <List.Item.Meta
              title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${comment.topic_id}?post=${comment.id}`} target="_blank" rel="noopener noreferrer">{normalizeQuote(comment.text)}</a>}
              description={`${comment.likes} likes`}
            />
          </List.Item>}
        />
      </Spin>
    </div>
  )
}

export default MyProfile
