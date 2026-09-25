import { Link } from 'react-router-dom'
import { useState } from 'react'
import { PollForm } from '../components/PollForm'
import { PollCreated } from '../components/PollCreated'
import { Shell } from '../components/Shell'
import type { Poll } from '../lib/types'

export function CreatePage(){const [createdPoll,setCreatedPoll]=useState<Poll|null>(null);return <Shell><div className="create-layout">{createdPoll?<PollCreated poll={createdPoll}/>:<><div className="create-intro"><Link className="back-link" to="/dashboard">← Back to desk</Link><span className="section-label">New room</span><h1>Put a question<br />in the air.</h1><p>Keep it clear. Give people a few good ways to answer. The room will take care of the rest.</p><div className="create-note"><span className="mono">01</span><span>Every poll gets a private creator link and a public share link.</span></div></div><section className="form-panel"><PollForm onCreated={setCreatedPoll}/></section></>}</div></Shell>}
