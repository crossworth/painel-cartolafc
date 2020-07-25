import React, { useEffect, useState } from 'react'
import { Link, withRouter } from 'react-router-dom'

import { Spin, Typography } from 'antd'
import { getUserInfo } from '../api'
import CommentList from './CommentList'

const { Title } = Typography

const ProfileComments = (props) => {
  const id = props.match.params.id

  const [user, setUser] = useState({})
  const [total, setTotal] = useState('∅')

  const { history } = props

  useEffect(() => {
    getUserInfo(id).then(data => {
      setUser(data)
    }).catch(err => {
      if (err.response && err.response.status && err.response.status === 404) {
        history.push('/perfil/nao-encontrado')
      }
    })
  }, [id, history])

  return (<div>
    <Spin tip="Carregando..." spinning={!user.first_name}>
      <Title level={4}>
        {
          user.first_name ?
            <div>Comentários de <Link to={`/perfil/${user.id}`}>{user.first_name} {user.last_name}</Link> - {total} comentários
            </div>
            : <div>Carregando dados</div>
        }
      </Title>
      <CommentList profileID={id} onCommentsTotal={total => setTotal(total)}/>
    </Spin>
  </div>)
}

export default withRouter(ProfileComments)
