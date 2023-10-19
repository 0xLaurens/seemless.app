export function isJSON(str: string) {
  try {
    JSON.parse(str)
  } catch (e) {
    return false
  }

  return true
}
