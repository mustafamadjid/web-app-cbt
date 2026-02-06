import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

import { AuthProvider } from './contexts/AuthContext'
// Router
import { RouterProvider } from 'react-router'
import {router} from './routes/AppRouter'

import { Toaster } from "react-hot-toast";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider >
      <Toaster position="top-center" />
      <RouterProvider router={router} >
      </RouterProvider>
    </AuthProvider>
  
  </StrictMode>,
)
