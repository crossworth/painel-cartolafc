import axios from 'axios'
import { message } from 'antd'
import { getErrorMessage } from './util'
import { API_BASE_URL } from './config'

const api = axios.create({
  baseURL: API_BASE_URL,
})

api.interceptors.response.use(response => {
  if (response.request.responseURL.indexOf('/fazer-login') !== -1) {
    return window.location.reload()
  }
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
  return api.get(`/profiles/${id}`)
}

const getProfileStats = id => {
  return api.get(`/profiles/${id}/stats`)
}

const getProfileHistory = id => {
  return api.get(`/profiles/${id}/history`)
}

const getTopicsFromProfile = (id, before = unixNow(), limit = 10) => {
  return api.get(`/profiles/${id}/topics?before=${before}&limit=${limit}`)
}

const getCommentsFromProfile = (id, before = unixNow(), limit = 10) => {
  return api.get(`/profiles/${id}/comments?before=${before}&limit=${limit}`)
}

const getProfiles = (page, limit, orderBy = 'topics', orderDir = 'desc', period = 'all') => {
  return api.get(`/profiles?orderBy=${orderBy}&orderDir=${orderDir}&page=${page}&limit=${limit}&period=${period}`)
}

const getTopics = (before, limit, orderBy = 'updated_at') => {
  return api.get(`/topics?orderBy=${orderBy}&before=${before}&limit=${limit}`)
}

const getTopicsRanking = (page, limit, orderBy = 'comments', orderDir = 'desc', period = 'all', showOlderTopics = 'true') => {
  return api.get(`/topics-ranking?orderBy=${orderBy}&orderDir=${orderDir}&page=${page}&limit=${limit}&period=${period}&showOlderTopics=${showOlderTopics}`)
}

const getTopicSearch = (term, page, limit, fromID = 0, createdAfter = 0, createdBefore = 0) => {
  return api.get(`/search?term=${term}&page=${page}&limit=${limit}&fromID=${fromID}&createdAfter=${createdAfter}&createdBefore=${createdBefore}`)
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
  getProfiles,
  getTopics,
  getTopicsRanking,
  getTopicSearch
}
