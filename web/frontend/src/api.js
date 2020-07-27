import axios from 'axios'
import { message } from 'antd'
import { getErrorMessage } from './util'

const baseURL = process.env.REACT_APP_API_URL

const api = axios.create({
  baseURL: baseURL,
})

api.interceptors.response.use((response) => {
  return response.data
}, (error) => {
  message.error(getErrorMessage(error))
  return Promise.reject(error)
})

const unixNow = () => {
  const time = Date.now()
  return Math.floor(time / 1000)
}

const resolveProfile = name => {
  if (name.indexOf('vk.com') === -1) {
    name = `vk.com/${name}`
  }

  return api.get(`/resolve-profile?link=${name}`)
}

const autoCompleteProfileNames = name => {
  name = name.replace('https://', '')
  name = name.replace('http://', '')
  name = name.replace('vk.com/', '')

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
  unixNow,
  resolveProfile,
  autoCompleteProfileNames,
  getUserInfo,
  getUserStats,
  getUserProfileHistory,
  getTopicsFromUser,
  getCommentsFromUser
}
