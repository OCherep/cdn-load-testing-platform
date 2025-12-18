import { useEffect, useState } from "react";
import { connectWS } from "../ws";

export default function TestView({ testId }: any) {
    const [metrics, setMetrics] = useState<any[]>([]);

    useEffect(() => {
        const ws = connectWS(testId, m => {
            setMetrics(prev => [...prev.slice(-50), m]);
        });
        return () => ws.close();
    }, []);

    return <pre>{JSON.stringify(metrics, null, 2)}</pre>;
}
