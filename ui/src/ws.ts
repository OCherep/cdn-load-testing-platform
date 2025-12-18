export function connectWS(testId: string, onMessage: (d: any) => void) {
    const ws = new WebSocket(`ws://controller:8080/ws/tests/${testId}`);
    ws.onmessage = e => onMessage(JSON.parse(e.data));
    return ws;
}
