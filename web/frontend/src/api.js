import axios from 'axios'
import { message } from 'antd'

const baseURL = 'http://localhost:8080/api'

const api = axios.create({
  baseURL: baseURL,
})

api.interceptors.response.use((response) => {
  return response.data
}, (error) => {

  let errorMessage = error

  if (error.response &&
    error.response.data &&
    error.response.data.message) {
    errorMessage = error.response.data.message
  }

  const errorMessageNormalized = errorMessage.charAt(0).toUpperCase() + errorMessage.slice(1)
  message.error(errorMessageNormalized)

  return Promise.reject(error)
})

const unixNow = () => {
  const time = Date.now()
  return time / 1000
}

const getUserInfo = id => {
  return api.get(`/user/${id}`)
}

const getTopicsFromUser = (id, before = unixNow(), limit = 10) => {
  return api.get(`/topics/${id}?before=${before}&limit=${limit}`)
}

const getCommentsFromUser = (id, before = unixNow(), limit = 10) => {
  return api.get(`/comments/${id}?before=${before}&limit=${limit}`)
}

export {
  getUserInfo,
  getTopicsFromUser,
  getCommentsFromUser
}
