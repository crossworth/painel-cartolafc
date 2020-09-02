import React, { useEffect, useState } from 'react'
import { List, Spin } from 'antd'
import { VK_GROUP_ID } from '../config'
import { timeStampToDate } from '../util'
import { getMyLastTopics } from '../api'

export default props => {
  const [loading, setLoading] = useState(false)
  const [lastTopics, setLastTopics] = useState([])

  useEffect(() => {
    setLoading(true)
    getMyLastTopics().then(result => setLastTopics(result)).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }, [])

  return <Spin tip="Carregando..." spinning={loading}>
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
  </Spin>
}
