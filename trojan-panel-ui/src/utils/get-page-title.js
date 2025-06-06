import defaultSettings from '@/settings'

const title = defaultSettings.title || 'Trojan Panel'

// export default function getPageTitle(pageTitle) {
//   if (pageTitle) {
//     return `${pageTitle} - ${title}`
//   }
//   return `${title}`
// }

export default function getPageTitle(pageTitle) {
  if (pageTitle) {
    return `${pageTitle}`
  }
  return ``
}
