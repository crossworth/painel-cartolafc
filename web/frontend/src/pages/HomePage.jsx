import React, { useEffect, useState } from 'react'

import { Spin } from 'antd'
import { getHomePage } from '../api'
import 'react-quill/dist/quill.snow.css'

const HomePage = () => {

  const [homePageContent, setHomePageContent] = useState(null)

  useEffect(() => {
    getHomePage().then(page => {
      setHomePageContent(page.value)
    }).catch(err => {

    })
  }, [])

  return (
    <div>
      <Spin tip="Carregando..." spinning={homePageContent === null}>
        <div className="ql-editor">
          <div dangerouslySetInnerHTML={{ __html: homePageContent }}/>
        </div>
      </Spin>
    </div>
  )
}

export default HomePage
