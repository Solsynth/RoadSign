export function createAuthHeader(password: string, username: string = "RoadSign CLI") {
  const credentials = Buffer.from(`${username}:${password}`).toString("base64")
  return `Basic ${credentials}`
}