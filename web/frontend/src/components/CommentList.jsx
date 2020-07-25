import React, { useEffect, useState } from 'react'
import { Table } from 'antd'
import { getCommentsFromUser } from '../api'
import { getBeforeFromURL, normalizeQuote, timeStampToDate } from '../util'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    render: (text, data) => <a href={`https://vk.com/topic-73721457_${data.topic_id}?post=${data.id}`}
                               rel="noopener noreferrer"
                               target="_blank">{text}</a>
  },
  {
    title: 'Comentário',
    dataIndex: 'text',
    key: 'text',
    render: (text, data) => <div>{normalizeQuote(text)}</div>
  },
  {
    title: 'Likes',
    dataIndex: 'likes',
    key: 'likes',
  },
  {
    title: 'Data',
    dataIndex: 'date',
    key: 'date',
    render: (text, data) => timeStampToDate(text)
  },
]

const CommentList = (props) => {
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
  })

  const [loading, setLoading] = useState(true)
  const [tableData, setTableData] = useState([])
  const [tableMeta, setTableMeta] = useState({})

  useEffect(() => {
    if (props.onCommentsTotal && tableMeta.total) {
      props.onCommentsTotal(tableMeta.total)
    }

    setPagination(Object.assign({}, pagination, {
      total: tableMeta.total
    }))

  }, [tableMeta])


  const handleTableChange = pag => {
    let beforeURL = tableMeta.current

    if (pag.current !== pagination.current) {
      beforeURL = pag.current > pagination.current ? tableMeta.next : tableMeta.prev
    }

    setLoading(true)

    getCommentsFromUser(props.profileID, getBeforeFromURL(beforeURL), pag.pageSize).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
    }).finally(() => {
      setLoading(false)
    })

    setPagination(pag)
  }

  useEffect(() => {
    getCommentsFromUser(props.profileID).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
    }).finally(() => {
      setLoading(false)
    })
  }, [])

  return (
    <div>
      <Table
        bordered={true}
        dataSource={tableData}
        columns={columns}
        loading={loading}
        rowKey='id'
        pagination={pagination}
        onChange={handleTableChange}
      />
    </div>
  )
}

export default CommentList
