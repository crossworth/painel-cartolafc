import React, { useEffect, useState } from 'react'
import { Col, Divider, Row, Spin } from 'antd'
import { ResponsiveCalendar } from '@nivo/calendar'
import { months } from '../util'
import { getCommentsGraph, getTopicsGraph } from '../api'

export default props => {
  const [loading, setLoading] = useState(false)
  const [topicsGraph, setTopicsGraph] = useState([])
  const [commentsGraph, setCommentsGraph] = useState([])

  useEffect(() => {
    setLoading(true)
    Promise.all([
      getTopicsGraph(),
      getCommentsGraph()
    ]).then(results => {
      setTopicsGraph(results[0].result.map(record => {
        return {
          day: record.day.substr(0, 10),
          value: record.value
        }
      }))
      setCommentsGraph(results[1].result.map(record => {
        return {
          day: record.day.substr(0, 10),
          value: record.value
        }
      }))
    }).catch(err => {

    }).finally(() => {
      setLoading(false)
    })
  }, [])

  const getTopicsGraphFrom = () => {
    if (topicsGraph.length === 0) {
      return (new Date()).toISOString().substr(0, 10)
    }
    return topicsGraph[0].day
  }

  const getTopicsGraphTo = () => {
    if (topicsGraph.length === 0) {
      return (new Date()).toISOString().substr(0, 10)
    }

    return topicsGraph[topicsGraph.length - 1].day
  }

  const getCommentsGraphFrom = () => {
    if (commentsGraph.length === 0) {
      return (new Date()).toISOString().substr(0, 10)
    }
    return commentsGraph[0].day
  }

  const getCommentsGraphTo = () => {
    if (commentsGraph.length === 0) {
      return (new Date()).toISOString().substr(0, 10)
    }

    return commentsGraph[commentsGraph.length - 1].day
  }

  return <Spin tip="Carregando..." spinning={loading}>
    <Divider plain>Tópicos por dia</Divider>
    <Row>
      <Col md={24} style={{
        height: '300px',
        maxWidth: '800px',
        margin: '0 auto',
      }}>
        <ResponsiveCalendar
          data={topicsGraph}
          from={getTopicsGraphFrom()}
          to={getTopicsGraphTo()}
          monthLegend={(year, month) => months[month]}
          emptyColor="#eeeeee"
          colors={['#bbdefb', '#64b5f6', '#1e88e5', '#0d47a1']}
          margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
          yearSpacing={40}
          monthBorderColor="#ffffff"
          dayBorderWidth={2}
          dayBorderColor="#ffffff"
          legends={[
            {
              anchor: 'bottom-right',
              direction: 'row',
              translateY: 36,
              itemCount: 4,
              itemWidth: 42,
              itemHeight: 36,
              itemsSpacing: 14,
              itemDirection: 'right-to-left'
            }
          ]}
        />
      </Col>
    </Row>

    <Divider plain>Comentários por dia</Divider>

    <Row>
      <Col md={24} style={{
        height: '300px',
        maxWidth: '800px',
        margin: '0 auto',
      }}>
        <ResponsiveCalendar
          data={commentsGraph}
          from={getCommentsGraphFrom()}
          to={getCommentsGraphTo()}
          monthLegend={(year, month) => months[month]}
          emptyColor="#eeeeee"
          colors={['#bbdefb', '#64b5f6', '#1e88e5', '#0d47a1']}
          margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
          yearSpacing={40}
          monthBorderColor="#ffffff"
          dayBorderWidth={2}
          dayBorderColor="#ffffff"
          legends={[
            {
              anchor: 'bottom-right',
              direction: 'row',
              translateY: 36,
              itemCount: 4,
              itemWidth: 42,
              itemHeight: 36,
              itemsSpacing: 14,
              itemDirection: 'right-to-left'
            }
          ]}
        />
      </Col>
    </Row>
  </Spin>
}
