const API = "http://controller:8080";

export async function api(path: string, method = "GET", body?: any) {
    const token = localStorage.getItem("jwt");
    const res = await fetch(API + path, {
        method,
        headers: {
            "Authorization": token || "",
            "Content-Type": "application/json",
        },
        body: body ? JSON.stringify(body) : undefined,
    });
    return res.json();
}
