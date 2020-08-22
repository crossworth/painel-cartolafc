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
  if (error.response.status === 300) {
    const url = new URL(`${window.location.protocol}${window.location.host}${error.response.data.to}`)
    const reason = url.searchParams.get('motivo-redirect')
    message.error(reason ? reason : error.response.data.to)
    return Promise.reject(error)
  }

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

const getSearch = (term, page, limit, searchType, fullText) => {
  return api.get(`/search?term=${term}&page=${page}&limit=${limit}&searchType=${searchType}&fullText=${fullText}`)
}

const getAdministratorsProfiles = () => {
  return api.get(`/administrators-profiles`)
}

const setAdministratorsProfiles = profilesIDs => {
  return api.post(`/set-administrators-profiles`, profilesIDs)
}

const getMyProfileStats = () => {
  return api.get(`/my-profile`)
}

const getMembersRules = () => {
  return api.get(`/members-rule`)
}

const setMembersRules = rules => {
  return api.post(`/set-members-rule`, {
    value: rules
  })
}

const getHomePage = () => {
  return api.get(`/home-page`)
}

const setHomePage = content => {
  return api.post(`/set-home-page`, {
    value: content
  })
}

const isAdmin = () => {
  return window.User.type === 'admin' || window.User.type === 'super_admin'
}

const isSuperAdmin = () => {
  return window.User.type === 'super_admin'
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
  getSearch,
  getAdministratorsProfiles,
  setAdministratorsProfiles,
  getMyProfileStats,
  getMembersRules,
  setMembersRules,
  getHomePage,
  setHomePage,
  isAdmin,
  isSuperAdmin
}
