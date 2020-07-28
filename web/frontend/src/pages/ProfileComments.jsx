import React from 'react'
import { withRouter } from 'react-router-dom'

import { Button } from 'antd'
import { normalizeQuote, timeStampToDate } from '../util'
import ProfileTabularData from '../components/ProfileTabularData'
import { getCommentsFromProfile } from '../api'


const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: 'Comentário',
    dataIndex: 'text',
    key: 'text',
    render: (text, data) => <div>{normalizeQuote(text)}</div>
  },
  {
    title: 'Likes',
    dataIndex: 'likes',
    key: 'likes',
  },
  {
    title: 'Data',
    dataIndex: 'date',
    key: 'date',
    render: (text, data) => timeStampToDate(text)
  },
  {
    title: '',
    dataIndex: 'id',
    key: 'id',
    render: (text, data) => <div>
      <Button type="primary" block target="_blank" rel="noopener noreferrer"
              href={`https://vk.com/topic-73721457_${data.topic_id}?post=${data.id}`}>
        Link original
      </Button>
      {/*<Button style={{ marginTop: 5 }} block target="_blank" rel="noopener noreferrer"*/}
      {/*        href={`/topico/${data.topic_id}?post=${data.id}`}>*/}
      {/*  Reconstituído*/}
      {/*</Button>*/}
    </div>

  },
]

const ProfileComments = (props) => {
  let newProps = Object.assign({}, props, {
    type: 'Comentários',
    columns: columns,
    dataFunc: getCommentsFromProfile
  })

  return (<div>
    <ProfileTabularData {...newProps} />
  </div>)
}

export default withRouter(ProfileComments)
