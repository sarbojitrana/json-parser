import { useState } from "react"
import "../App.css"
import Sidebar from "../components/Sidebar"
import UploadSpreadsheet from "../components/UploadSpreadsheet"
import UploadWorkflow from "../components/UploadWorkflow"
import SearchApplication from "../components/SearchApplication"
import LogsViewer from "../components/LogsViewer"

function Home() {
    const [activeTab, setActiveTab] = useState<"upload" | "search" | "logs">("search")

    return (
        <div style={{ display: "flex", minHeight: "100vh", backgroundColor: "#f3f4f6" }}>
            
            <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />

            <div style={{ 
                flex: 1, 
                marginLeft: "260px", 
                padding: "40px",
                boxSizing: "border-box"
            }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "30px" }}>
                    <h1 style={{ margin: 0, color: "#111827", fontSize: "28px", fontWeight: "bold" }}>Official Registry Framework</h1>
                </div>

                {/* Upload View */}
                {activeTab === "upload" && (
                    <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
                        <div style={{ background: "white", padding: "24px", borderRadius: "8px", border: "1px solid #e5e7eb" }}>
                            <UploadSpreadsheet />
                        </div>
                        <div style={{ background: "white", padding: "24px", borderRadius: "8px", border: "1px solid #e5e7eb" }}>
                            <UploadWorkflow />
                        </div>
                    </div>
                )}

                {/* Application Search View */}
                {activeTab === "search" && (
                    <div style={{ background: "white", padding: "24px", borderRadius: "8px", border: "1px solid #e5e7eb" }}>
                        <SearchApplication />
                    </div>
                )}

                {/* Log Viewer */}
                {activeTab === "logs" && (
                    <div style={{ background: "white", padding: "24px", borderRadius: "8px", border: "1px solid #e5e7eb" }}>
                        <LogsViewer />
                    </div>
                )}
            </div>
        </div>
    )
}

export default Home