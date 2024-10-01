export function createAuthHeader(username: string, password: string) {
  const credentials = Buffer.from(`${username}:${password}`).toString("base64")
  return `Basic ${credentials}`
}