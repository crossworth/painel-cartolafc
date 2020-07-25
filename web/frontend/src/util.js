const timeStampToDate = timestamp => {
  const date = new Date(timestamp * 1000)
  return date.toLocaleString()
}

const normalizeQuote = text => {
  return text.replace(/\[(id|club)\d+:?(bp-\d+_\d+)?\|([@\d\w\s_\-áàâãéèêíïóôõöúçñ]+)]/gim, '$3')
}

// exemplo de input
// http://localhost:8080/api/comments/259592548?limit=10&before=1595645176
const getBeforeFromURL = url => {
  const u = new URL(url)
  const urlParams = new URLSearchParams(u.search)
  return urlParams.get('before')
}

export {
  timeStampToDate,
  getBeforeFromURL,
  normalizeQuote
}
