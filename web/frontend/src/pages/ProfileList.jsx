import React, { useEffect, useState } from 'react'
import { Button, Col, Radio, Row, Space, Spin, Table, Typography } from 'antd'
import { getProfiles } from '../api'
import {
  getGlobalPageSize,
  normalizeScreenName,
  parseIntWithDefault,
  setGlobalPageSize,
  stringWithDefault
} from '../util'
import { withRouter } from 'react-router-dom'

const { Title, Text } = Typography

const columns = [
  {
    title: 'Rank',
    dataIndex: 'position',
    key: 'position',
  },
  {
    title: 'Nome',
    dataIndex: 'name',
    render: (text, data) => <div>
      <Row>
        <Col flex="60px">
          <img width="50" height="50" src={data.photo} alt=''/>
        </Col>
        <Col flex="auto">
          {`${data.first_name} ${data.last_name} (@${normalizeScreenName(data.screen_name, data.id)})`}
        </Col>
      </Row>
    </div>
  },
  {
    title: 'Tópicos',
    dataIndex: 'topics',
    key: 'topics',
  },
  {
    title: 'Comentários',
    dataIndex: 'comments',
    key: 'comments',
  },
  {
    title: 'Likes',
    dataIndex: 'likes',
    key: 'likes',
  },
  {
    title: 'Tópicos+Comentários',
    dataIndex: 'topics_plus_comments',
    key: 'topics_plus_comments',
  },
  {
    title: '',
    dataIndex: 'id',
    key: 'id',
    render: (text, data) => <div>
      <Button type="primary" block target="_blank" rel="noopener noreferrer"
              href={`https://vk.com/id${data.id}`}>
        VK
      </Button>
      <Button style={{ marginTop: 5 }} block target="_blank" rel="noopener noreferrer"
              href={`/perfil/${data.id}`}>
        Perfil
      </Button>
    </div>

  },
]

const orderByToPTBR = orderBy => {
  switch (orderBy) {
    case 'topics':
      return 'tópicos'
    case 'comments':
      return 'comentários'
    case 'topics_comments':
      return 'tópicos+comentários'
    default:
      return 'likes'
  }
}

const ProfileList = (props) => {

  const { history, location } = props

  const searchParams = new URLSearchParams(location.search)

  const [pagination, setPagination] = useState({
    current: parseIntWithDefault(searchParams.get('page'), 1),
    pageSize: parseIntWithDefault(searchParams.get('limit'), getGlobalPageSize(10)),
    position: ['topLeft'],
    showSizeChanger: true,
    orderBy: stringWithDefault(searchParams.get('orderBy'), 'topics'),
    orderDir: stringWithDefault(searchParams.get('orderDir'), 'desc'),
    period: stringWithDefault(searchParams.get('period'), 'all'),
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

    if (!searchParams.has('page') || pagination.current !== parseInt(searchParams.get('page'))) {
      shouldUpdate = true
      searchParams.set('page', pagination.current)
    }

    if (!searchParams.has('orderBy') || pagination.orderBy !== searchParams.get('orderBy')) {
      shouldUpdate = true
      searchParams.set('orderBy', pagination.orderBy)
    }

    if (!searchParams.has('orderDir') || pagination.orderDir !== searchParams.get('orderDir')) {
      shouldUpdate = true
      searchParams.set('orderDir', pagination.orderDir)
    }

    if (!searchParams.has('period') || pagination.period !== searchParams.get('period')) {
      shouldUpdate = true
      searchParams.set('period', pagination.period)
    }

    if (shouldUpdate) {
      history.push({
        pathname: location.pathname,
        search: searchParams.toString()
      })
    }
  }, [history, location, pagination])

  const { current, pageSize, orderBy, orderDir, period } = pagination

  const setPaginationTotal = (total, page = null, orderBy = null, orderDir = null, pageSize = null) => {
    let pag = Object.assign({}, pagination, {
      total: total,
      current: page ? page : pagination.current,
      orderBy: orderBy ? orderBy : pagination.orderBy,
      orderDir: orderDir ? orderDir : pagination.orderDir,
      pageSize: pageSize ? pageSize : pagination.pageSize,
    })

    setGlobalPageSize(pag.pageSize)
    setPagination(pag)
  }

  const setOrderBy = orderBy => {
    let pag = Object.assign({}, pagination, {
      orderBy: orderBy ? orderBy : pagination.orderBy,
    })
    setPagination(pag)
  }

  const setPeriod = period => {
    let pag = Object.assign({}, pagination, {
      period: period ? period : pagination.period,
    })
    setPagination(pag)
  }

  const handleTableChange = pag => {
    setLoading(true)

    getProfiles(pag.current, pag.pageSize, pag.orderBy, pag.orderDir, pag.period).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
      setPaginationTotal(data.meta.total, pag.current, pag.orderBy, pag.orderDir, pag.pageSize)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }

  useEffect(() => {
    setLoading(true)
    getProfiles(current, pageSize, orderBy, orderDir, period).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
      setPaginationTotal(data.meta.total)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [current, pageSize, orderBy, orderDir, period])

  return (
    <Spin tip="Carregando..." spinning={loading}>
      <Title level={4}>
        {
          !loading ?
            <div>Lista de membros por {orderByToPTBR(pagination.orderBy)} - {tableMeta.total} membros </div>
            : <div>Carregando dados</div>
        }
      </Title>
      <div>
        {!loading &&
        <div>
          <Space direction="vertical">
            <Radio.Group value={pagination.period} onChange={event => setPeriod(event.target.value)}>
              <Radio.Button value="all">Sempre</Radio.Button>
              <Radio.Button value="last_month">Último mês</Radio.Button>
              <Radio.Button value="last_week">Última semana</Radio.Button>
              <Radio.Button value="last_day">Último dia</Radio.Button>
            </Radio.Group>

            <Radio.Group value={pagination.orderBy} onChange={event => setOrderBy(event.target.value)}>
              <Radio.Button value="topics">Por tópicos</Radio.Button>
              <Radio.Button value="comments">Por comentários</Radio.Button>
              <Radio.Button value="likes">Por likes</Radio.Button>
              <Radio.Button value="topics_comments">Por tópicos+comentários</Radio.Button>
            </Radio.Group>
          </Space>
        </div>
        }
        {
          !loading && tableMeta.cached_at && <div>
            <Text type="secondary">Atualizado em: {(new Date(tableMeta.cached_at)).toLocaleString()}</Text>
          </div>
        }
      </div>
      <Table
        bordered={true}
        dataSource={tableData}
        columns={columns}
        rowKey='id'
        scroll={{ x: true }}
        pagination={pagination}
        onChange={handleTableChange}
      />
    </Spin>
  )
}

export default withRouter(ProfileList)
