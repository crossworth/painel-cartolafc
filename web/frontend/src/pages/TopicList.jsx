import React from 'react'
import { withRouter } from 'react-router-dom'

import { Button } from 'antd'
import { timeStampToDate } from '../util'
import { getTopics } from '../api'
import TopicTabularData from '../components/TopicTabularData'
import { VK_GROUP_ID } from '../config'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: 'Título',
    dataIndex: 'title',
    key: 'title',
  },
  {
    title: 'Data criação',
    dataIndex: 'created_at',
    key: 'created_at',
    render: (text, data) => timeStampToDate(text)
  },
  {
    title: 'Última atualização',
    dataIndex: 'updated_at',
    key: 'updated_at',
    render: (text, data) => timeStampToDate(text)
  },
  {
    title: 'Comentários',
    dataIndex: 'comments_count',
    key: 'comments_count',
  },
  {
    title: '',
    dataIndex: 'id',
    key: 'id',
    render: (text, data) => <div>
      <Button type="primary" block target="_blank" rel="noopener noreferrer"
              href={`https://vk.com/topic-${VK_GROUP_ID}_${data.id}`}>
        Link original
      </Button>
      {/*<Button style={{ marginTop: 5 }} block target="_blank" rel="noopener noreferrer"*/}
      {/*        href={`/topico/${data.topic_id}?post=${data.id}`}>*/}
      {/*  Reconstituído*/}
      {/*</Button>*/}
    </div>

  },
]

const TopicList = (props) => {
  let newProps = Object.assign({}, props, {
    columns: columns,
    dataFunc: getTopics
  })

  return (<div>
    <TopicTabularData {...newProps} />
  </div>)
}

export default withRouter(TopicList)
