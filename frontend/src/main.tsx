import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

import { AuthProvider } from './contexts/AuthContext'
// Router
import { RouterProvider } from 'react-router'
import {router} from './routes/AppRouter'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider >
      <RouterProvider router={router}/>
    </AuthProvider>
  
  </StrictMode>,
)
