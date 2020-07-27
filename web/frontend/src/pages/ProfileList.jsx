import React, { useEffect, useState } from 'react'
import { getBeforeFromURL, timeStampToDate } from '../util'
import { Button, Spin, Table, Typography } from 'antd'
import { getTopicsFromUser } from '../api'
import { Link } from 'react-router-dom'

const { Title } = Typography

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

const ProfileList = () => {
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    position: ['topLeft']
  })

  const [loading, setLoading] = useState(true)
  const [tableData, setTableData] = useState([])
  const [tableMeta, setTableMeta] = useState({})


  const handleTableChange = pag => {
    let beforeURL = tableMeta.current

    if (pag.current !== pagination.current) {
      beforeURL = pag.current > pagination.current ? tableMeta.next : tableMeta.prev
    }

    setLoading(true)

    // getTopicsFromUser(profileID, getBeforeFromURL(beforeURL), pag.pageSize).then(data => {
    //   setTableData(data.data)
    //   setTableMeta(data.meta)
    // }).finally(() => {
    //   setLoading(false)
    // })

    setPagination(pag)
  }

  useEffect(() => {
    // getTopicsFromUser(profileID).then(data => {
    //   setTableData(data.data)
    //   setTableMeta(data.meta)
    // }).finally(() => {
    //   setLoading(false)
    // })
  }, [])

  return (
    <Spin tip="Carregando..." spinning={loading}>
      <Title level={4}>
        {
          loading ?
            <div>Lista de membros - {tableMeta.total} membros
            </div>
            : <div>Carregando dados</div>
        }
      </Title>
      <Table
        bordered={true}
        dataSource={tableData}
        columns={columns}
        loading={loading}
        rowKey='id'
        pagination={pagination}
        onChange={handleTableChange}
      />
    </Spin>
  )
}

export default ProfileList
