export function Loading({ label = 'Loading...' }: { label?: string }) { return <div className="loading-state" role="status"><span className="loading-line" /><span>{label}</span></div> }
