import React, { useEffect, useState } from 'react'
import { Button, Radio, Space, Spin, Switch, Table, Typography } from 'antd'
import { getTopicsRanking } from '../api'
import { getGlobalPageSize, parseIntWithDefault, setGlobalPageSize, stringWithDefault, timeStampToDate } from '../util'
import { withRouter } from 'react-router-dom'
import { VK_GROUP_ID } from '../config'

const { Title, Text } = Typography

const columns = [
  {
    title: 'Rank',
    dataIndex: 'position',
    key: 'position',
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
    dataIndex: 'comments',
    key: 'comments',
  },
  {
    title: 'Likes',
    dataIndex: 'likes',
    key: 'likes',
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
    </div>

  },
]

const orderByToPTBR = orderBy => {
  switch (orderBy) {
    case 'comments':
      return 'comentários'
    default:
      return 'likes'
  }
}

const TopicRankingList = (props) => {

  const { history, location } = props

  const searchParams = new URLSearchParams(location.search)

  const [pagination, setPagination] = useState({
    current: parseIntWithDefault(searchParams.get('page'), 1),
    pageSize: parseIntWithDefault(searchParams.get('limit'), getGlobalPageSize(10)),
    position: ['topLeft'],
    showSizeChanger: true,
    orderBy: stringWithDefault(searchParams.get('orderBy'), 'comments'),
    orderDir: stringWithDefault(searchParams.get('orderDir'), 'desc'),
    period: stringWithDefault(searchParams.get('period'), 'all'),
    showOlderTopics: stringWithDefault(searchParams.get('showOlderTopics'), 'true'),
    pageSizeOptions: [10, 20, 50, 100, 1000]
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

    if (!searchParams.has('showOlderTopics') || pagination.showOlderTopics !== searchParams.get('showOlderTopics')) {
      shouldUpdate = true
      searchParams.set('showOlderTopics', pagination.showOlderTopics)
    }

    if (shouldUpdate) {
      history.push({
        pathname: location.pathname,
        search: searchParams.toString()
      })
    }
  }, [history, location, pagination])

  const { current, pageSize, orderBy, orderDir, period, showOlderTopics } = pagination

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

  const toggleShowOlderTopics = () => {
    let pag = Object.assign({}, pagination, {
      showOlderTopics: pagination.showOlderTopics === 'true' ? 'false' : 'true',
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

    getTopicsRanking(pag.current, pag.pageSize, pag.orderBy, pag.orderDir, pag.period, pag.showOlderTopics).then(data => {
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
    getTopicsRanking(current, pageSize, orderBy, orderDir, period, showOlderTopics).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
      setPaginationTotal(data.meta.total)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [current, pageSize, orderBy, orderDir, period, showOlderTopics])

  return (
    <Spin tip="Carregando..." spinning={loading}>
      <Title level={4}>
        {
          !loading ?
            <div>Ranking de tópicos por número de {orderByToPTBR(pagination.orderBy)} - {tableMeta.total} tópicos </div>
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
              <Radio.Button value="comments">Por comentários</Radio.Button>
              <Radio.Button value="likes">Por likes</Radio.Button>
            </Radio.Group>

            {
              pagination.period !== 'all' && <Switch
                checkedChildren="Exibir tópicos antigos" unCheckedChildren="Exibir tópicos antigos"
                checked={pagination.showOlderTopics === 'true'}
                onChange={toggleShowOlderTopics}/>
            }
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

export default withRouter(TopicRankingList)
