import { thunderIDProxy, createRouteMatcher, } from '@thunderid/nextjs/server'

const isProtectedRoute = createRouteMatcher(['/kpis'])

export default thunderIDProxy(async (thunderid, request) => {
  if (isProtectedRoute(request)) {
    console.log('isSignedIn:', thunderid.isSignedIn())
    return await thunderid.protectRoute()
  }
})

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico).*)'],
}