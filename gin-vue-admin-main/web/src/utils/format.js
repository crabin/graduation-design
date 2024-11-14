import { formatTimeToStr } from '@/utils/date'
import { getDict } from '@/utils/dictionary'

export const formatBoolean = (bool) => {
  if (bool !== null) {
    return bool ? '是' : '否'
  } else {
    return ''
  }
}
export const formatDate = (time) => {
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
  } else {
    return ''
  }
}

export const filterDict = (value, options) => {
  const rowLabel = options && options.filter(item => item.value === value)
  return rowLabel && rowLabel[0] && rowLabel[0].label
}

export const getDictFunc = async(type) => {
  const dicts = await getDict(type)
  return dicts
}

/**
 * 校验IP
 * @param str
 * @returns {boolean}
 */
export const validateIP = (str) => {
  const reg = /^([1-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.([0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.([0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.([0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])(:([1-9](\d{0,3})|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5]))?$/g
  return reg.test(str)
}

