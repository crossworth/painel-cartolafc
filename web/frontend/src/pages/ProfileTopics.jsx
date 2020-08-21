import React from 'react'
import { withRouter } from 'react-router-dom'

import { Button } from 'antd'
import { getTopicsFromProfile } from '../api'
import { timeStampToDate } from '../util'
import ProfileTabularData from '../components/ProfileTabularData'
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
    title: 'Soma likes',
    dataIndex: 'likes',
    key: 'likes',
  },
  {
    title: 'Data',
    dataIndex: 'created_at',
    key: 'created_at',
    render: (text, data) => timeStampToDate(text)
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
      {/*        href={`/topico/${data.id}`}>*/}
      {/*  Reconstituído*/}
      {/*</Button>*/}
    </div>

  },
]

const ProfileTopics = (props) => {
  let newProps = Object.assign({}, props, {
    type: 'Tópicos',
    columns: columns,
    dataFunc: getTopicsFromProfile
  })

  return (<div>
    <ProfileTabularData {...newProps} />
  </div>)
}

export default withRouter(ProfileTopics)
