import axios from 'axios'
import { message } from 'antd'
import { getErrorMessage } from './util'

const baseURL = process.env.REACT_APP_API_URL

const api = axios.create({
  baseURL: baseURL,
})

api.interceptors.response.use(response => {
  return response.data
}, error => {
  message.error(getErrorMessage(error).toString())
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

const getProfileInfo = id => {
  return api.get(`/profile/${id}`)
}

const getProfileStats = id => {
  return api.get(`/profile/${id}/stats`)
}

const getProfileHistory = id => {
  return api.get(`/profile/${id}/history`)
}

const getTopicsFromProfile = (id, before = unixNow(), limit = 10) => {
  return api.get(`/topics/${id}?before=${before}&limit=${limit}`)
}

const getCommentsFromProfile = (id, before = unixNow(), limit = 10) => {
  return api.get(`/comments/${id}?before=${before}&limit=${limit}`)
}

const getProfiles = (page, limit, orderBy = 'topics', orderDir = 'desc', period = 'all') => {
  return api.get(`/profiles?orderBy=${orderBy}&orderDir=${orderDir}&page=${page}&limit=${limit}&period=${period}`)
}

export {
  unixNow,
  resolveProfile,
  autoCompleteProfileNames,
  getProfileInfo,
  getProfileStats,
  getProfileHistory,
  getTopicsFromProfile,
  getCommentsFromProfile,
  getProfiles
}
