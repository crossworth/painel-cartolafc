import React, { useEffect, useState } from 'react'
import { Link, withRouter } from 'react-router-dom'
import { Spin, Table, Typography } from 'antd'
import { getBeforeFromURL, parseIntWithDefault } from '../util'
import { getUserInfo, unixNow } from '../api'

const { Title } = Typography

const ProfileTabularData = (props) => {
  const id = props.match.params.id

  const { history, location } = props

  const searchParams = new URLSearchParams(location.search)

  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: parseIntWithDefault(searchParams.get('limit'), 10),
    currentTimestamp: parseIntWithDefault(searchParams.get('current'), unixNow()),
    position: ['topLeft']
  })

  const [loading, setLoading] = useState(true)
  const [user, setUser] = useState({})
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
    setPagination(pag)
  }

  useEffect(() => {
    Promise.all([getUserInfo(id), dataFunc(id, currentTimestamp, pageSize)]).then(results => {
      setUser(results[0])
      // NOTE(Pedro): Hack to enable cursor pagination
      // when page refreshed
      let page = 1
      if (results[1].meta.prev) {
        page = 2
      }

      if (!results[1].meta.next) {
        page = Math.ceil(results[1].meta.total / pageSize)
      }

      setPaginationCurrentTimeStampAndPage(parseInt(getBeforeFromURL(results[1].meta.current)), page, results[1].meta.total)

      setTableData(results[1].data)
      setTableMeta(results[1].meta)
    }).catch(err => {
      if (err.response && err.response.status && err.response.status === 404) {
        history.push('/perfil/nao-encontrado')
      }
    }).finally(() => {
      setLoading(false)
    })

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, history, dataFunc])

  const handleTableChange = pag => {
    let beforeTimestamp = tableMeta.current

    if (pag.current !== pagination.current) {
      beforeTimestamp = pag.current > pagination.current ? tableMeta.next : tableMeta.prev
    }

    setLoading(true)

    dataFunc(id, getBeforeFromURL(beforeTimestamp), pag.pageSize).then(data => {
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
    }).finally(() => {
      setLoading(false)
    })
  }

  return (<div>
    <Spin tip="Carregando..." spinning={loading}>
      <Title level={4}>
        {
          !loading ?
            <div>{props.type} de <Link
              to={`/perfil/${user.id}`}>{user.first_name} {user.last_name}</Link> - {tableMeta ? tableMeta.total : '∅'} {props.type ? props.type.toLowerCase() : ''}
            </div>
            : <div>Carregando dados</div>
        }
      </Title>
      <Table
        bordered={true}
        dataSource={tableData}
        columns={props.columns}
        rowKey='id'
        pagination={pagination}
        onChange={handleTableChange}
      />
    </Spin>
  </div>)
}

export default withRouter(ProfileTabularData)
