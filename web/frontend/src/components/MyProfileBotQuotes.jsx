import React, { useEffect, useState } from 'react'
import { Button, Col, List, Row, Spin } from 'antd'
import { VK_GROUP_ID } from '../config'
import { timeStampToDate } from '../util'
import { getMyProfileBotQuotes } from '../api'

const quotesByBotLimit = 20

export default props => {
  const [loading, setLoading] = useState(false)
  const [quotesByBot, setQuotesByBot] = useState([])
  const [quotesByBotPage, setQuotesByBotPage] = useState(1)

  useEffect(() => {
    setLoading(true)
    getMyProfileBotQuotes(quotesByBotPage, quotesByBotLimit).then(result => {
      setQuotesByBot(result)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }, [quotesByBotPage])

  return <Spin tip="Carregando..." spinning={loading}>
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
      header={<strong>Vezes que você foi citado pelo bot {`(${quotesByBot.meta ? quotesByBot.meta.total : 0})`}</strong>}
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
  </Spin>
}
