import { useState } from "react"
import { getApplications } from "../api"
import ResultsLedger from "./ResultsLedger"
import { toast } from "react-toastify"

function SearchApplication() {
    const [fromDate, setFromDate] = useState("")
    const [toDate, setToDate] = useState("")
    const [page, setPage] = useState(1)
    const [limit, setLimit] = useState(10)
    const [results, setResults] = useState<any[] | null>(null)
    const [loading, setLoading] = useState(false)
    const [viewMode, setViewMode] = useState<"search" | "ledger">("search")

    async function handleSearch(targetPage = page) {
        if (!fromDate || !toDate) {
            toast.dark("From Date and To Date filters are required parameters.", {
                style: { backgroundColor: "#374151", color: "#f9fafb" }
            })
            return
        }

        try {
            setLoading(true)
            const response = await getApplications({
                from: fromDate,
                to: toDate,
                page: targetPage,
                limit: limit
            })
            setResults(response.data ?? [])
            setViewMode("ledger")
        } catch (err: any) {
            // console.error(err)
            const errMsg = err?.response?.data?.error || "Registry query parsing execution failed."
            toast.error(errMsg, {
                style: { backgroundColor: "#f3f4f6", color: "#111827", borderLeft: "4px solid #ef4444" }
            })
            setResults(null)
        } finally {
            setLoading(false)
        }
    }

    function handlePrevious() {
        if (page > 1) {
            const prevPage = page - 1
            setPage(prevPage)
            handleSearch(prevPage)
        }
    }

    function handleNext() {
        if (results && results.length === limit) {
            const nextPage = page + 1
            setPage(nextPage)
            handleSearch(nextPage)
        }
    }

    function handleBackToSearch() {
        setViewMode("search") 
    }

    if (viewMode === "ledger" && results !== null) {
        return (
            <ResultsLedger 
                results={results}
                page={page}
                limit={limit}
                loading={loading}
                fromDate={fromDate}
                toDate={toDate}
                onBack={handleBackToSearch}
                onNext={handleNext}
                onPrevious={handlePrevious}
            />
        )
    }

    return (
        <div style={{ textAlign: "left" }}>
            <h2 style={{ color: "#111827", marginTop: 0 }}>Search Applications By Date Range</h2>
            <p style={{ color: "#4b5563", fontSize: "14px", marginBottom: "20px" }}>Select specific registry bounds to isolate application records.</p>

            <div style={{ marginBottom: "12px" }}>
                <label style={{ display: "block", marginBottom: "6px", fontWeight: "500", color: "#374151" }}>From Date: </label>
                <input
                    type="date"
                    value={fromDate}
                    onChange={(e) => setFromDate(e.target.value)}
                    style={{ width: "100%", maxWidth: "300px", padding: "8px", border: "1px solid #d1d5db", borderRadius: "4px" }}
                />
            </div>

            <div style={{ marginBottom: "12px" }}>
                <label style={{ display: "block", marginBottom: "6px", fontWeight: "500", color: "#374151" }}>To Date: </label>
                <input
                    type="date"
                    value={toDate}
                    onChange={(e) => setToDate(e.target.value)}
                    style={{ width: "100%", maxWidth: "300px", padding: "8px", border: "1px solid #d1d5db", borderRadius: "4px" }}
                />
            </div>

            <div style={{ marginBottom: "12px" }}>
                <label style={{ display: "block", marginBottom: "6px", fontWeight: "500", color: "#374151" }}>Page: </label>
                <input
                    type="number"
                    min="1"
                    value={page}
                    onChange={(e) => setPage(Math.max(1, Number(e.target.value)))}
                    style={{ width: "80px", padding: "8px", border: "1px solid #d1d5db", borderRadius: "4px" }}
                />
            </div>

            <div style={{ marginBottom: "20px" }}>
                <label style={{ display: "block", marginBottom: "6px", fontWeight: "500", color: "#374151" }}>Limit: </label>
                <input
                    type="number"
                    min="1"
                    max="100"
                    value={limit}
                    onChange={(e) => setLimit(Number(e.target.value))}
                    style={{ width: "80px", padding: "8px", border: "1px solid #d1d5db", borderRadius: "4px" }}
                />
            </div>

            <button 
                onClick={() => { setPage(1); handleSearch(1); }} 
                disabled={loading}
                style={{ padding: "10px 20px", backgroundColor: "#111827", color: "#fff", border: "none", borderRadius: "4px", cursor: "pointer", fontWeight: "bold" }}
            >
                {loading ? "Searching Repository..." : "Execute Search"}
            </button>
        </div>
    )
}

export default SearchApplication