import axios from 'axios'
import { message } from 'antd'
import { checkNested } from './util'

const baseURL = 'http://localhost:8080/api'

const api = axios.create({
  baseURL: baseURL,
})

api.interceptors.response.use((response) => {
  return response.data
}, (error) => {

  let errorMessage = error

  if (checkNested(error, 'response', 'data', 'message')) {
    errorMessage = error.response.data.message
  }

  if (errorMessage instanceof String) {
    const errorMessageNormalized = errorMessage.charAt(0).toUpperCase() + errorMessage.slice(1)
    message.error(errorMessageNormalized)
  }

  return Promise.reject(error)
})

const unixNow = () => {
  const time = Date.now()
  return time / 1000
}

const resolveProfile = name => {
  if (name.indexOf('vk.com') === -1) {
    name = `vk.com/${name}`
  }

  return api.get(`/resolve-profile?link=${name}`)
}

const autoCompleteProfileNames = name => {
  return api.get(`/auto-complete/profile/${name}`)
}

const getUserInfo = id => {
  return api.get(`/user/${id}`)
}

const getUserStats = id => {
  return api.get(`/user/${id}/stats`)
}

const getUserProfileHistory = id => {
  return api.get(`/user/${id}/history`)
}

const getTopicsFromUser = (id, before = unixNow(), limit = 10) => {
  return api.get(`/topics/${id}?before=${before}&limit=${limit}`)
}

const getCommentsFromUser = (id, before = unixNow(), limit = 10) => {
  return api.get(`/comments/${id}?before=${before}&limit=${limit}`)
}

export {
  resolveProfile,
  autoCompleteProfileNames,
  getUserInfo,
  getUserStats,
  getUserProfileHistory,
  getTopicsFromUser,
  getCommentsFromUser
}
