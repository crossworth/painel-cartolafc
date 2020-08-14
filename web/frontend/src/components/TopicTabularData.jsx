import React, { useEffect, useState } from 'react'
import { withRouter } from 'react-router-dom'
import { Alert, Spin, Table, Typography } from 'antd'
import { getBeforeFromURL, getGlobalPageSize, parseIntWithDefault, setGlobalPageSize } from '../util'
import { unixNow } from '../api'

const { Title } = Typography

const TopicTabularData = (props) => {
  const { history, location } = props

  const searchParams = new URLSearchParams(location.search)

  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: parseIntWithDefault(searchParams.get('limit'), getGlobalPageSize(10)),
    currentTimestamp: parseIntWithDefault(searchParams.get('current'), unixNow()),
    position: ['topLeft'],
    showSizeChanger: true
  })

  const [loading, setLoading] = useState(true)
  const [tableData, setTableData] = useState([])
  const [tableMeta, setTableMeta] = useState({})

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search)
    let shouldUpdate = false

    if (!searchParams.has('limit') || pagination.pageSize !== parseInt(searchParams.get('limit'))) {
      shouldUpdate = true
      searchParams.set('limit', pagination.pageSize)
    }

    if (!searchParams.has('current') || parseInt(pagination.currentTimestamp) !== parseInt(searchParams.get('current'))) {
      shouldUpdate = true
      searchParams.set('current', pagination.currentTimestamp)
    }

    if (shouldUpdate) {
      history.push({
        pathname: location.pathname,
        search: searchParams.toString()
      })
    }
  }, [history, location, pagination])

  const { dataFunc } = props

  const { currentTimestamp, pageSize } = pagination

  const setPaginationCurrentTimeStampAndPage = (currentTimestamp, page, total, pageSize = null) => {
    let pag = Object.assign({}, pagination, {
      currentTimestamp: currentTimestamp,
      current: page,
      total: total,
      pageSize: pageSize ? pageSize : pagination.pageSize
    })

    setGlobalPageSize(pag.pageSize)
    setPagination(pag)
  }

  useEffect(() => {
    dataFunc(currentTimestamp, pageSize).then(result => {
      // NOTE(Pedro): Hack to enable cursor pagination
      // when page is refreshed
      let page = 1
      if (result.meta.prev) {
        page = 2
      }

      if (!result.meta.next) {
        page = Math.ceil(result.meta.total / pageSize)
      }

      setPaginationCurrentTimeStampAndPage(parseInt(getBeforeFromURL(result.meta.current)), page, result.meta.total)

      setTableData(result.data)
      setTableMeta(result.meta)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [history, dataFunc])

  const handleTableChange = pag => {
    let beforeTimestamp = tableMeta.current

    if (pag.current !== pagination.current) {
      beforeTimestamp = pag.current > pagination.current ? tableMeta.next : tableMeta.prev
    }

    setLoading(true)

    dataFunc(getBeforeFromURL(beforeTimestamp), pag.pageSize).then(data => {
      let page = 1
      if (data.meta.prev) {
        page = 2
      }

      if (!data.meta.next) {
        page = Math.ceil(data.meta.total / pag.pageSize)
      }

      setPaginationCurrentTimeStampAndPage(getBeforeFromURL(data.meta.current), page, data.meta.total, pag.pageSize)

      setTableData(data.data)
      setTableMeta(data.meta)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }

  return (<div>
    <Spin tip="Carregando..." spinning={loading}>
      <Title level={4}>
        {
          !loading ?
            <div>Tópicos - {tableMeta ? tableMeta.total : '∅'} {props.type ? props.type.toLowerCase() : ''}
            </div>
            : <div>Carregando dados</div>
        }
      </Title>

      <Alert
        message="Velocidade de atualização dos tópicos"
        description={<div>
          Atualmente a forma como é indexado os tópicos é por <a href={`https://i.imgur.com/iX3VUxP.png`}
                                                                 target="_blank"
                                                                 rel="noopener noreferrer">Workers</a>, que não é a
          forma mais rápida e possui muita overhead em uma comunidade
          com tópicos frequentemente atualizados. Podendo levar minutos para ser verificado todos os tópicos da
          comunidade.
          Futuramente a ideia é alterar para uma forma Pub/Sub com WebHooks fornecida pela API do VK que deve tornar o
          processo
          em tempo real para maioria dos tópicos.
        </div>}
        type="info"
        showIcon
      />

      <div className="cursor-pagination">
        <Table
          bordered={true}
          dataSource={tableData}
          columns={props.columns}
          rowKey='id'
          scroll={{ x: true }}
          pagination={pagination}
          onChange={handleTableChange}
        />
      </div>
    </Spin>
  </div>)
}

export default withRouter(TopicTabularData)
