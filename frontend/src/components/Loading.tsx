export function Loading({ label = 'Loading the room' }: { label?: string }) { return <div className="loading-state" role="status"><span className="loading-line" /><span>{label}</span></div> }
