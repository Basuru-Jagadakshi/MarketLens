'use server'

import { thunderid } from '@thunderid/nextjs/server'


export async function getAccessToken(): Promise<string | undefined> {
  try {
    const sdk = await thunderid()
    const sessionId = await sdk.getSessionId()
    console.log('[getAccessToken] Session ID:', sessionId)
    if (!sessionId) return undefined
    const token = await sdk.getAccessToken(sessionId)
    console.log('[getAccessToken] Access Token:', token)

    return token
  } catch {
    return undefined
  }
}