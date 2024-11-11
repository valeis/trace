export function formatDate(ts) {
    const now = new Date();
    const diffInMs = now-ts;
    const oneDayInMs = 24*60*60*1000;

    if (diffInMs < oneDayInMs){
        return ts.toLocaleTimeString('en-GB', {hour: '2-digit', minute: '2-digit'});
    } else {
        return ts.toLocaleDateString('en-GB', {
            day: '2-digit',
            month: 'short',
            year: 'numeric',
        })
    }
}