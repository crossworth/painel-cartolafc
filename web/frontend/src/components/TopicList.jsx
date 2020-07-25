import React, { useEffect, useState } from 'react'
import { Table } from 'antd'
import { getTopicsFromUser } from '../api'
import { getBeforeFromURL, timeStampToDate } from '../util'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    render: (text, data) => <a href={`https://vk.com/topic-73721457_${data.id}`} rel="noopener noreferrer"
                               target="_blank">{text}</a>
  },
  {
    title: 'Título',
    dataIndex: 'title',
    key: 'title',
    render: (text, data) => <div>{text}</div>
  },
  {
    title: 'Data',
    dataIndex: 'created_at',
    key: 'created_at',
    render: (text, data) => timeStampToDate(text)
  },
]

const TopicList = (props) => {
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
  })

  const [loading, setLoading] = useState(true)
  const [tableData, setTableData] = useState([])
  const [tableMeta, setTableMeta] = useState({})

  useEffect(() => {
    if (props.onTopicsTotal && tableMeta.total) {
      props.onTopicsTotal(tableMeta.total)
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

    getTopicsFromUser(props.profileID, getBeforeFromURL(beforeURL), pag.pageSize).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
    }).finally(() => {
      setLoading(false)
    })

    setPagination(pag)
  }

  useEffect(() => {
    getTopicsFromUser(props.profileID).then(data => {
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

export default TopicList
