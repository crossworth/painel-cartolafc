import React, { useEffect, useState } from 'react'

import { Button, Col, Divider, List, Row, Spin, Statistic, Tabs, Typography } from 'antd'
import { normalizeComment, normalizeQuote, normalizeScreenName, timeStampToDate } from '../util'
import { getMyLastTopics, getMyProfileBotQuotes, getMyProfileStats } from '../api'
import { VK_GROUP_ID } from '../config'

const { Title } = Typography
const { TabPane } = Tabs

const quotesByBotLimit = 20

const MyProfile = () => {

  const [doneInitialLoading, setDoneInitialLoading] = useState(false)
  const [user, setUser] = useState({})
  const [loading, setLoading] = useState(false)
  const [userStats, setUserStats] = useState({})
  const [topicWithMoreLikes, setTopicWithMoreLikes] = useState([])
  const [topicWithMoreComments, setTopicWithMoreComments] = useState([])
  const [commentsWithMoreLikes, setCommentsWithMoreLikes] = useState([])
  const [quotesByBot, setQuotesByBot] = useState([])
  const [lastTopics, setLastTopics] = useState([])
  const [quotesByBotPage, setQuotesByBotPage] = useState(1)

  useEffect(() => {
    setLoading(true)
    Promise.all([
      getMyProfileStats(),
      getMyProfileBotQuotes(quotesByBotPage, quotesByBotLimit),
      getMyLastTopics()
    ]).then(results => {
      setUser(results[0].user)
      setUserStats(results[0].stats)
      setTopicWithMoreComments(results[0].topic_with_more_comments)
      setCommentsWithMoreLikes(results[0].comments_with_more_likes)
      setTopicWithMoreLikes(results[0].topic_with_more_likes)
      setQuotesByBot(results[1])
      setLastTopics(results[2])
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
      setDoneInitialLoading(true)
    })

  }, [])

  useEffect(() => {
    if (!doneInitialLoading) {
      return
    }

    setLoading(true)
    getMyProfileBotQuotes(quotesByBotPage, quotesByBotLimit).then(values => {
      console.log(values)
      setQuotesByBot(values)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }, [quotesByBotPage])

  return (
    <div>
      <Spin tip="Carregando..." spinning={loading}>
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

        <Tabs type="card">
          <TabPane tab="Estatísticas" key="1">
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
                  title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${topic.id}`} target="_blank"
                            rel="noopener noreferrer">{topic.title}</a>}
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
                  title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${topic.id}`} target="_blank"
                            rel="noopener noreferrer">{topic.title}</a>}
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
                  title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${comment.topic_id}?post=${comment.id}`}
                            target="_blank"
                            rel="noopener noreferrer">{normalizeComment(normalizeQuote(comment.text))}</a>}
                  description={`${comment.likes} likes`}
                />
              </List.Item>}
            />
          </TabPane>
          <TabPane tab={`Citações Bot (${quotesByBot.meta ? quotesByBot.meta.total : 0})`} key="2">
            <Row gutter={16}>
              <Col>
                <Button
                  disabled={quotesByBotPage <= 1}
                  onClick={() => {
                    setQuotesByBotPage(quotesByBotPage - 1)
                  }}>Página anterior</Button>
              </Col>
              <Col>
                <Button
                  disabled={quotesByBot.meta === undefined || quotesByBotPage >= (quotesByBot.meta.total / quotesByBotLimit)}
                  onClick={() => {
                    setQuotesByBotPage(quotesByBotPage + 1)
                  }}>Próxima página</Button>
              </Col>
            </Row>
            <br/>
            <List
              size="small"
              header={<strong>Vezes que você foi citado pelo bot</strong>}
              bordered
              dataSource={quotesByBot.data}
              renderItem={quote => <List.Item>
                <List.Item.Meta
                  title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${quote.topic_id}?post=${quote.comment_id}`}
                            target="_blank" rel="noopener noreferrer">{quote.topic_title}</a>}
                  description={`em ${timeStampToDate(quote.date_comment)}`}
                />
              </List.Item>}
            />
          </TabPane>
          <TabPane tab="Últimos tópicos" key="3">
            <List
              size="small"
              header={<strong>Seus últimos tópicos</strong>}
              bordered
              dataSource={lastTopics}
              renderItem={quote => <List.Item>
                <List.Item.Meta
                  title={<a href={`https://vk.com/topic-${VK_GROUP_ID}_${quote.id}`}
                            target="_blank" rel="noopener noreferrer">{quote.title}</a>}
                  description={`em ${timeStampToDate(quote.created_at)}`}
                />
              </List.Item>}
            />
          </TabPane>
        </Tabs>
      </Spin>
    </div>
  )
}

export default MyProfile
