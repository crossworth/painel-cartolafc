import React, { useEffect, useState } from 'react'
import { Button, Form, Input, Space, Spin, Switch, Table, Typography } from 'antd'
import { getTopicSearch } from '../api'
import { getGlobalPageSize, parseIntWithDefault, setGlobalPageSize, stringWithDefault, timeStampToDate } from '../util'
import { withRouter } from 'react-router-dom'
import { VK_GROUP_ID } from '../config'

const { Title, Text } = Typography

const typeToPTBR = type => {
  if (type === 'topic') {
    return 'tópico'
  }

  return 'comentário'
}

const columns = [
  {
    title: 'Tipo',
    dataIndex: 'type',
    key: 'type',
    render: (text, data) => <div>{typeToPTBR(text)}</div>
  },
  {
    title: 'Texto',
    dataIndex: 'highlighted_part',
    key: 'highlighted_part',
    render: (text, data) => <div dangerouslySetInnerHTML={{ __html: text }}/>
  },
  {
    title: 'Data',
    dataIndex: 'date',
    key: 'date',
    render: (text, data) => timeStampToDate(text)
  },
  {
    title: '',
    dataIndex: 'topic_id',
    key: 'topic_id',
    render: (text, data) => <div>
      {
        data.type === 'topic' &&
        <Button type="primary" block target="_blank" rel="noopener noreferrer"
                href={`https://vk.com/topic-${VK_GROUP_ID}_${data.topic_id}`}>
          Link original
        </Button>
      }
      {
        data.type === 'comment' &&
        <Button type="primary" block target="_blank" rel="noopener noreferrer"
                href={`https://vk.com/topic-${VK_GROUP_ID}_${data.topic_id}?post=${data.comment_id}`}>
          Link original
        </Button>
      }
    </div>

  },
]

const TopicSearch = (props) => {
  const { history, location } = props

  const searchParams = new URLSearchParams(location.search)

  const [pagination, setPagination] = useState({
    current: parseIntWithDefault(searchParams.get('page'), 1),
    pageSize: parseIntWithDefault(searchParams.get('limit'), getGlobalPageSize(10)),
    position: ['topLeft'],
    showSizeChanger: true,
    term: stringWithDefault(searchParams.get('term'), ''),
    pageSizeOptions: [10, 20, 50, 100, 1000]
  })

  const [loading, setLoading] = useState(true)
  const [termInput, setTermInput] = useState(pagination.term)
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

    if (!searchParams.has('term') || pagination.term !== searchParams.get('term')) {
      shouldUpdate = true
      searchParams.set('term', pagination.term)
    }

    if (shouldUpdate) {
      history.push({
        pathname: location.pathname,
        search: searchParams.toString()
      })
    }
  }, [history, location, pagination])

  const { current, pageSize, term } = pagination

  const setPaginationTotal = (total, page = null, pageSize = null) => {
    let pag = Object.assign({}, pagination, {
      total: total,
      current: page ? page : pagination.current,
      pageSize: pageSize ? pageSize : pagination.pageSize,
    })

    setGlobalPageSize(pag.pageSize)
    setPagination(pag)
  }

  const setSearch = (term) => {
    let pag = Object.assign({}, pagination, {
      term: term,
    })
    setPagination(pag)
  }

  const handleTableChange = pag => {
    setLoading(true)

    getTopicSearch(pag.term, pag.current, pag.pageSize).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
      setPaginationTotal(data.meta.total, pag.current, pag.pageSize)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }

  const onSearch = () => {
    if (termInput.length > 0) {
      setSearch(termInput)
    }
  }

  useEffect(() => {
    if (term === '' || term.length === 0) {
      setLoading(false)
      return
    }

    setLoading(true)
    getTopicSearch(term, current, pageSize).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
      setPaginationTotal(data.meta.total)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }, [term, current, pageSize])

  return (
    <Spin tip="Carregando..." spinning={loading}>
      <Title level={4}>
        {
          !loading ?
            (
              term === '' || term.length === 0 ? <div>Faça uma pesquisa</div>
                : <div>Resultado da pesquisa por {term} - {tableMeta.total} resultados </div>
            )
            : <div>Carregando dados</div>
        }
      </Title>
      <div>
        {!loading &&
        <div>
          <Space direction="vertical" style={{ width: '100%', textAlign: 'left' }}>
            <Form onFinish={onSearch}>
              <Form.Item label="Termo" rules={[{ required: true }]}>
                <Input name="term" placeholder="Messi" defaultValue={term} onChange={event => setTermInput(event.target.value)} />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit">
                  Pesquisar
                </Button>
              </Form.Item>
            </Form>

          </Space>
        </div>
        }
      </div>
      <Table
        bordered={true}
        dataSource={tableData}
        columns={columns}
        rowKey={record => `${record.topic_id}_${record.comment_id}`}
        scroll={{ x: true }}
        pagination={pagination}
        onChange={handleTableChange}
      />
    </Spin>
  )
}

export default withRouter(TopicSearch)
