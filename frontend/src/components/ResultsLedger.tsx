interface ResultsLedgerProps {
    results: any[]
    page: number
    limit: number
    loading: boolean
    onBack: () => void
    onNext: () => void
    onPrevious: () => void
    fromDate: string
    toDate: string
}

interface LedgerFieldProps {
    label: string
    value: any
}

function LedgerField({ label, value }: LedgerFieldProps) {
    if (!value) return null // If null, undefined, or empty string, ignore 
    return (
        <p style={{ margin: 0 }}>
            <strong>{label}:</strong> {value}
        </p>
    )
}

function ResultsLedger({
    results,
    page,
    limit,
    loading,
    onBack,
    onNext,
    onPrevious,
    fromDate,
    toDate
}: ResultsLedgerProps) {
    return (
        <div style={{ padding: "30px", backgroundColor: "#fff", minHeight: "100vh", textAlign: "left" }}>
            
            {/* Top Navigation Control */}
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", borderBottom: "2px solid #0B2545", paddingBottom: "15px", marginBottom: "20px" }}>
                <button 
                    onClick={onBack}
                    style={{ padding: "8px 16px", backgroundColor: "#0B2545", color: "#fff", border: "none", borderRadius: "4px", cursor: "pointer", fontWeight: "bold" }}
                >
                    ← BACK TO SEARCH DASHBOARD
                </button>
                <div style={{ fontSize: "14px", color: "#555" }}>
                    <strong>Active Range:</strong> {fromDate} to {toDate} | <strong>Page:</strong> {page}
                </div>
            </div>

            <h2>Official Records Ledger ({results.length} items on this page)</h2>
            <hr style={{ borderColor: "#eee", margin: "15px 0" }} />

            {results.length === 0 ? (
                <p style={{ textAlign: "center", color: "#666", padding: "40px" }}>No records found for this registry filter execution.</p>
            ) : (
                <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
                    {results.map((app: any, index: number) => (
                        <div 
                            key={app.id || index} 
                            style={{ 
                                padding: "20px", 
                                border: "1px solid #ccc", 
                                borderRadius: "4px", 
                                backgroundColor: "#fafafa",
                                boxShadow: "0 1px 3px rgba(0,0,0,0.05)"
                            }}
                        >
                            <h3 style={{ marginTop: 0, color: "#0B2545", borderBottom: "1px style solid #eee", paddingBottom: "8px" }}>
                                Entry #{((page - 1) * limit) + index + 1}
                            </h3>

                            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "10px 20px" }}>
                                <LedgerField label="Application Ref No" value={app.appl_ref_no} />
                                <LedgerField label="Application ID" value={app.appl_id} />
                                <LedgerField label="Applicant Name" value={app.applicant_name} />
                                <LedgerField label="District" value={app.district} />
                                <LedgerField label="District LGD Code" value={app.district_lgd_code} />
                                <LedgerField label="Sub Division" value={app.sub_division} />
                                <LedgerField label="Subdivision LGD Code" value={app.subdivision_lgd_code} />
                                <LedgerField label="Block" value={app.block} />
                                <LedgerField label="Block LGD Code" value={app.block_lgd_code} />
                                <LedgerField label="Pincode" value={app.pincode} />
                                <LedgerField label="Submission Location" value={app.submission_location} />
                                <LedgerField label="Submitted By" value={app.submitted_by} />
                                <LedgerField label="Submission Date" value={app.submission_date} />

                                {/* Status of application */}
                                {app.status && (
                                    <p style={{ margin: 0 }}>
                                        <strong>Status:</strong>{" "}
                                        <span style={{ 
                                            fontWeight: "bold",
                                            color: app.status === "Y" ? "green" : 
                                                   app.status === "N" ? "red" : 
                                                   app.status === "P" ? "orange" : "#333" 
                                        }}>
                                            {app.status}
                                        </span>
                                    </p>
                                )}
                                
                                <LedgerField label="Service ID" value={app.service_id} />
                                <LedgerField label="Service Name" value={app.service_name} />
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Pagination */}
            {results.length > 0 && (
                <div style={{ marginTop: "30px", display: "flex", justifyContent: "center", gap: "20px", alignItems: "center", borderTop: "1px solid #ccc", paddingTop: "20px" }}>
                    <button 
                        onClick={onPrevious} 
                        disabled={page === 1 || loading}
                        style={{ padding: "8px 16px", cursor: page === 1 ? "not-allowed" : "pointer" }}
                    >
                        ← Previous Page
                    </button>
                    <span style={{ fontWeight: "bold" }}>Page {page}</span>
                    <button 
                        onClick={onNext} 
                        disabled={results.length < limit || loading}
                        style={{ padding: "8px 16px", cursor: results.length < limit ? "not-allowed" : "pointer" }}
                    >
                        Next Page →
                    </button>
                </div>
            )}
        </div>
    )
}

export default ResultsLedger