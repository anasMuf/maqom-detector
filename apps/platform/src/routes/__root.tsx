import { Outlet, createRootRouteWithContext, Link } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { Music, History, Library, Home } from 'lucide-react'

import '../styles.css'
import type { useAuth } from '../features/auth/AuthContext'

interface MyRouterContext {
  auth: ReturnType<typeof useAuth>
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootComponent,
})

function RootComponent() {
  return (
    <div className="min-h-screen bg-background font-sans text-foreground">
      {/* Navbar */}
      <header className="sticky top-0 z-50 w-full border-b border-border bg-background/80 backdrop-blur">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <Link to="/" className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-primary text-white">
              <Music size={18} />
            </div>
            <span className="text-xl font-bold tracking-tight text-brand-primary">Maqam<span className="text-foreground">Detector</span></span>
          </Link>

          <nav className="hidden md:flex items-center gap-6">
            <Link to="/" className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground [&.active]:text-brand-primary [&.active]:font-semibold">Beranda</Link>
            <Link to="/analyze" className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground [&.active]:text-brand-primary [&.active]:font-semibold">Deteksi</Link>
            <Link to="/maqamat" className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground [&.active]:text-brand-primary [&.active]:font-semibold">Kamus Maqamat</Link>
            <Link to="/history" className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground [&.active]:text-brand-primary [&.active]:font-semibold">Riwayat</Link>
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto flex-1 p-4 md:p-8">
        <Outlet />
      </main>

      {/* Mobile Bottom Nav */}
      <div className="fixed bottom-0 left-0 z-50 w-full border-t border-border bg-background md:hidden pb-safe">
        <div className="flex h-16 items-center justify-around px-2">
          <Link to="/" className="flex flex-col items-center gap-1 text-muted-foreground [&.active]:text-brand-primary">
            <Home size={20} />
            <span className="text-[10px] font-medium">Beranda</span>
          </Link>
          <Link to="/analyze" className="flex flex-col items-center gap-1 text-muted-foreground [&.active]:text-brand-primary">
            <Music size={20} />
            <span className="text-[10px] font-medium">Deteksi</span>
          </Link>
          <Link to="/maqamat" className="flex flex-col items-center gap-1 text-muted-foreground [&.active]:text-brand-primary">
            <Library size={20} />
            <span className="text-[10px] font-medium">Kamus</span>
          </Link>
          <Link to="/history" className="flex flex-col items-center gap-1 text-muted-foreground [&.active]:text-brand-primary">
            <History size={20} />
            <span className="text-[10px] font-medium">Riwayat</span>
          </Link>
        </div>
      </div>

      <TanStackDevtools
        config={{ position: 'bottom-right' }}
        plugins={[{ name: 'TanStack Router', render: <TanStackRouterDevtoolsPanel /> }]}
      />
    </div>
  )
}
