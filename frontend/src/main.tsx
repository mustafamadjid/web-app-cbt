import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'

// Router
import { RouterProvider } from 'react-router'
import {router} from './routes/AppRouter'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router}/>
  </StrictMode>,
)
