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

function checkNested(obj, level, ...rest) {
  if (obj === undefined) {
    return false
  }

  if (rest.length === 0 && obj.hasOwnProperty(level)) {
    return true
  }

  return checkNested(obj[level], ...rest)
}

const debounce = (func, wait = 250) => {
  let inDebounce
  return function () {
    const context = this
    const args = arguments
    clearTimeout(inDebounce)
    inDebounce = setTimeout(() => func.apply(context, args), wait)
  }
}

export {
  timeStampToDate,
  getBeforeFromURL,
  normalizeQuote,
  checkNested,
  debounce
}
