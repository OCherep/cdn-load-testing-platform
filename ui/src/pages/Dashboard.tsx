import { useEffect, useState } from "react";
import { api } from "../api";

export default function Dashboard() {
    const [tests, setTests] = useState([]);

    useEffect(() => {
        api("/tests").then(setTests);
    }, []);

    return (
        <div>
            <h1>CDN Load Tests</h1>
            <ul>
                {tests.map((t:any) => (
                    <li key={t.test_id}>
                        {t.test_id} — {t.status}
                    </li>
                ))}
            </ul>
        </div>
    );
}
