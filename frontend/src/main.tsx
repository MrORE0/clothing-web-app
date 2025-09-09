import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import MainPage from './pages/MainPage.tsx'
import { CustomizationPage } from './pages/CustomizationPage.tsx'
import {createBrowserRouter, RouterProvider} from "react-router-dom"

const router = createBrowserRouter([{path: "/", element: <MainPage/>}, {path: "/customize", element: <CustomizationPage/>}])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router}/>
  </StrictMode>,
)
