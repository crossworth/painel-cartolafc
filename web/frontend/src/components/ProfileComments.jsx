import React, { useEffect, useState } from 'react'
import { withRouter } from 'react-router-dom'

import { Spin, Typography } from 'antd'
import { getUserInfo } from '../api'
import CommentList from './CommentList'

const { Title } = Typography

const ProfileComments = (props) => {
  const id = props.match.params.id

  const [user, setUser] = useState({})
  const [total, setTotal] = useState('∅')

  useEffect(() => {
    getUserInfo(id).then(data => {
      setUser(data)
    }).catch(err => {
      if (err.response.status === 404) {
        props.history.push('/perfil/nao-encontrado')
      }
    })
  }, [])

  return (<div>
    <Spin tip="Carregando..." spinning={!user.first_name}>
      <Title level={3}>
        {user.first_name ? `Comentários de ${user.first_name} ${user.last_name} - ${total} comentários` : 'Carregando dados'}
      </Title>
      <CommentList profileID={id} onCommentsTotal={total => setTotal(total)}/>
    </Spin>
  </div>)
}

export default withRouter(ProfileComments)
