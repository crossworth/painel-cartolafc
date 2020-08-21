import React from 'react'
import { Typography } from 'antd'
import { changeLog } from '../changelog'

const { Title } = Typography

export default () => {
  return (<div>
    <Title level={4}>
      ChangeLog
    </Title>
    <div style={{ whiteSpace: 'pre' }} dangerouslySetInnerHTML={{ __html: changeLog }}/>
  </div>)
}
