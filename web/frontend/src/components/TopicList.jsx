import React, { useEffect, useState } from 'react'
import { Button, Table } from 'antd'
import { getTopicsFromUser } from '../api'
import { getBeforeFromURL, timeStampToDate } from '../util'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
  },
  {
    title: 'Título',
    dataIndex: 'title',
    key: 'title',
  },
  {
    title: 'Data',
    dataIndex: 'created_at',
    key: 'created_at',
    render: (text, data) => timeStampToDate(text)
  },
  {
    title: '',
    dataIndex: 'id',
    key: 'id',
    render: (text, data) => <div>
      <Button type="primary" block target="_blank" rel="noopener noreferrer"
              href={`https://vk.com/topic-73721457_${data.id}`}>
        Link original
      </Button>
      <Button style={{ marginTop: 5 }} block target="_blank" rel="noopener noreferrer"
              href={`/topico/${data.id}`}>
        Reconstituído
      </Button>
    </div>

  },
]

const TopicList = ({ onTopicsTotal, profileID }) => {
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    position: ['topLeft']
  })

  const [loading, setLoading] = useState(true)
  const [tableData, setTableData] = useState([])
  const [tableMeta, setTableMeta] = useState({})

  useEffect(() => {
    if (onTopicsTotal && tableMeta.total !== undefined) {
      onTopicsTotal(tableMeta.total)
    }

    setPagination(Object.assign({}, pagination, {
      total: tableMeta.total
    }))

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [onTopicsTotal, tableMeta])

  const handleTableChange = pag => {
    let beforeURL = tableMeta.current

    if (pag.current !== pagination.current) {
      beforeURL = pag.current > pagination.current ? tableMeta.next : tableMeta.prev
    }

    setLoading(true)

    getTopicsFromUser(profileID, getBeforeFromURL(beforeURL), pag.pageSize).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
    }).finally(() => {
      setLoading(false)
    })

    setPagination(pag)
  }

  useEffect(() => {
    getTopicsFromUser(profileID).then(data => {
      setTableData(data.data)
      setTableMeta(data.meta)
    }).finally(() => {
      setLoading(false)
    })
  }, [profileID])

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
