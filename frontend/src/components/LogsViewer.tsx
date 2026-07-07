import { useEffect, useState } from "react"
import { api } from "../api"
import { toast } from "react-toastify"

function LogsViewer() {
    const [logs, setLogs] = useState<any[]>([])
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        async function fetchLogs() {
            try {
                setLoading(true)
                const response = await api.get("/logs")
                setLogs(response.data ?? [])
            } catch (err) {
                toast.error("Audit Logs fetch routine failed. Check connection parameter roots.", {
                    style: { backgroundColor: "#f3f4f6", color: "#111827", borderLeft: "4px solid #ef4444" }
                })
            } finally {
                setLoading(false)
            }
        }
        fetchLogs()
    }, [])

    return (
        <div style={{ textAlign: "left" }}>
            <h2 style={{ color: "#111827", marginBottom: "5px" }}>System Execution Logs</h2>
            <p style={{ color: "#4b5563", fontSize: "14px", marginBottom: "20px" }}>
                Real-time backend audit prints and system validation routines.
            </p>

            {loading ? (
                <p>Loading registry logs...</p>
            ) : logs.length === 0 ? (
                <p style={{ color: "#6b7280", padding: "20px", backgroundColor: "#f3f4f6", borderRadius: "6px" }}>
                    No system log records available in the database registry.
                </p>
            ) : (
                <div style={{ display: "flex", flexDirection: "column", gap: "10px", maxHeight: "75vh", overflowY: "auto" }}>
                    {logs.map((log: any, index: number) => (
                        <div 
                            key={log.id || index} 
                            style={{ 
                                padding: "12px 16px", 
                                borderLeft: `4px solid ${log.level === "ERROR" ? "#dc2626" : "#2563eb"}`,
                                backgroundColor: "#f9fafb", 
                                borderRadius: "4px",
                                borderTop: "1px solid #e5e7eb",
                                borderRight: "1px solid #e5e7eb",
                                borderBottom: "1px solid #e5e7eb"
                            }}
                        >
                            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "6px", fontSize: "13px" }}>
                                <span style={{ fontWeight: "bold", color: log.level === "ERROR" ? "#dc2626" : "#2563eb" }}>
                                    [{log.level || "INFO"}]
                                </span>
                                <span style={{ fontWeight: "bold", color: "#4b5563" }}>Source: {log.source || "N/A"}</span>
                                <span style={{ color: "#9ca3af" }}>{log.created_at || ""}</span>
                            </div>
                            <p style={{ margin: 0, fontFamily: "monospace", color: "#1f2937", wordBreak: "break-all" }}>
                                {log.message}
                            </p>
                        </div>
                    ))}
                </div>
            )}
        </div>
    )
}

export default LogsViewer