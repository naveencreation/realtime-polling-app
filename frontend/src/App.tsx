import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { HomePage } from './pages/HomePage'
import { AuthPage } from './pages/AuthPage'
import { DashboardPage } from './pages/DashboardPage'
import { CreatePage } from './pages/CreatePage'
import { PollPage } from './pages/PollPage'

export function App(){return <BrowserRouter><Routes><Route path="/" element={<HomePage/>}/><Route path="/login" element={<AuthPage mode="login"/>}/><Route path="/signup" element={<AuthPage mode="signup"/>}/><Route path="/dashboard" element={<DashboardPage/>}/><Route path="/create" element={<CreatePage/>}/><Route path="/poll/:id" element={<PollPage/>}/><Route path="*" element={<Navigate to="/" replace/>}/></Routes></BrowserRouter>}
