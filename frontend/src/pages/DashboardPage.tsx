import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import type { Poll } from '../lib/types'
import { PollCard } from '../components/PollCard'
import { Loading } from '../components/Loading'
import { Shell } from '../components/Shell'

export function DashboardPage(){const navigate=useNavigate();const [polls,setPolls]=useState<Poll[]>([]);const [loading,setLoading]=useState(true);const [error,setError]=useState('');useEffect(()=>{api.listMine().then(setPolls).catch(()=>{setError('Your creator session has ended.');navigate('/login')}).finally(()=>setLoading(false))},[navigate]);return <Shell><div className="page-shell"><div className="page-heading"><div><span className="section-label">Creator desk</span><h1>Your rooms.</h1><p>Make a question worth gathering around.</p></div><Link className="primary-button" to="/create">New poll <span aria-hidden="true">↗</span></Link></div>{error&&<div className="error-banner">{error}</div>}{loading?<Loading label="Loading your rooms"/>:polls.length===0?<div className="empty-state"><span className="empty-number">01</span><h2>The room is quiet.</h2><p>You have not published a poll yet. Start with the question your group keeps circling.</p><Link className="primary-button" to="/create">Publish your first poll <span aria-hidden="true">↗</span></Link></div>:<div className="poll-grid">{polls.map((poll)=><PollCard key={poll.id} poll={poll}/>)}</div>}</div></Shell>}
