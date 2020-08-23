import React, { useEffect, useState } from 'react'
import { Button, Checkbox, Form, Input, Radio, Space, Spin, Table, Typography } from 'antd'
import { getSearch } from '../api'
import {
  getGlobalPageSize, normalizeComment,
  normalizeQuote,
  parseIntWithDefault,
  setGlobalPageSize,
  stringWithDefault,
  timeStampToDate
} from '../util'
import { withRouter } from 'react-router-dom'
import { VK_GROUP_ID } from '../config'

const { Title } = Typography

const columns = [
  {
    title: 'Texto',
    dataIndex: 'headline',
    key: 'headline',
    render: (text, data) => <div dangerouslySetInnerHTML={{ __html: normalizeComment(normalizeQuote(text)) }}/>
  },
  {
    title: 'Data',
    dataIndex: 'date',
    key: 'date',
    render: (text, data) => timeStampToDate(text)
  },
  {
    title: '',
    dataIndex: 'type',
    key: 'type',
    render: (text, data) => {
      if (data.type === 'topic') {
        return <div>{data.comments_count} comentários</div>
      } else {
        return <div>{data.likes_count} likes</div>
      }
    }
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

const Search = (props) => {
  const { history, location } = props

  const searchParams = new URLSearchParams(location.search)

  const [pagination, setPagination] = useState({
    current: parseIntWithDefault(searchParams.get('page'), 1),
    pageSize: parseIntWithDefault(searchParams.get('limit'), getGlobalPageSize(10)),
    position: ['topLeft'],
    showSizeChanger: true,
    term: stringWithDefault(searchParams.get('term'), ''),
    searchType: stringWithDefault(searchParams.get('searchType'), 'title'),
    fullText: stringWithDefault(searchParams.get('fullText'), 'true'),
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

    if (!searchParams.has('searchType') || pagination.searchType !== searchParams.get('searchType')) {
      shouldUpdate = true
      searchParams.set('searchType', pagination.searchType)
    }

    if (!searchParams.has('fullText') || pagination.fullText !== searchParams.get('fullText')) {
      shouldUpdate = true
      searchParams.set('fullText', pagination.fullText)
    }

    if (shouldUpdate) {
      history.push({
        pathname: location.pathname,
        search: searchParams.toString()
      })
    }
  }, [history, location, pagination])

  const { current, pageSize, term, searchType, fullText } = pagination

  const setPaginationTotal = (total, page = null, pageSize = null) => {
    let pag = Object.assign({}, pagination, {
      total: total,
      current: page ? page : pagination.current,
      pageSize: pageSize ? pageSize : pagination.pageSize,
    })

    setGlobalPageSize(pag.pageSize)
    setPagination(pag)
  }

  const setSearch = term => {
    let pag = Object.assign({}, pagination, {
      term: term,
    })
    setPagination(pag)
  }

  const setSearchType = searchType => {
    let pag = Object.assign({}, pagination, {
      searchType: searchType,
    })
    setPagination(pag)
  }

  const setFullText = fullText => {
    let pag = Object.assign({}, pagination, {
      fullText: fullText.toString(),
    })
    setPagination(pag)
  }

  const handleTableChange = pag => {
    setLoading(true)

    getSearch(pag.term, pag.current, pag.pageSize, pag.searchType, pag.fullText).then(data => {
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

  const changeSearchType = event => {
    const type = event.target.value
    setSearchType(type)
  }

  useEffect(() => {
    if (term === '' || term.length === 0) {
      setLoading(false)
      return
    }

    setLoading(true)
    getSearch(term, current, pageSize, searchType, fullText).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
      setPaginationTotal(data.meta.total)
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }, [term, current, pageSize, searchType, fullText])

  const searchTypeOptions = [
    { label: 'Título tópico', value: 'title' },
    { label: 'Conteúdo tópico', value: 'text' },
  ]

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
                <Input name="term" placeholder="Messi" defaultValue={term}
                       onChange={event => setTermInput(event.target.value)}/>
              </Form.Item>

              <Form.Item>
                <Radio.Group
                  options={searchTypeOptions}
                  onChange={changeSearchType}
                  value={pagination.searchType}
                  optionType="button"
                  buttonStyle="solid"
                />
              </Form.Item>

              <Form.Item>
                <Checkbox checked={fullText === 'false'} onChange={e => setFullText(!e.target.checked)}>Pesquisa exata</Checkbox>
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

export default withRouter(Search)
