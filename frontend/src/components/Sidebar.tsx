interface SidebarProps {
    activeTab: "upload" | "search" | "logs"
    setActiveTab: (tab: "upload" | "search" | "logs") => void
}

function Sidebar({ activeTab, setActiveTab }: SidebarProps) {
    const navItems = [
        { id: "upload", label: "Upload Center" },
        { id: "search", label: "Application Search" },
        { id: "logs", label: "System Logs" }
    ] as const

    return (
        <div style={{ 
            width: "260px", 
            backgroundColor: "#f9fafb", 
            color: "#111827",
            padding: "24px 16px", 
            display: "flex", 
            flexDirection: "column",
            position: "fixed",
            height: "100vh",
            left: 0,
            top: 0,
            borderRight: "1px solid #e5e7eb" 
        }}>
            <h2 style={{ 
                fontSize: "18px", 
                margin: "0 0 24px 0", 
                letterSpacing: "0.5px", 
                borderBottom: "1px solid #e5e7eb", 
                paddingBottom: "12px", 
                textAlign: "left",
                color: "#111827",
                fontWeight: "bold"
            }}>
                Workflow Dashboard
            </h2>
            
            <nav style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
                {navItems.map((item) => (
                    <button
                        key={item.id}
                        onClick={() => setActiveTab(item.id)}
                        style={{
                            width: "100%",
                            textAlign: "left",
                            padding: "12px 16px",
                            borderRadius: "6px",
                            border: "none",
                            cursor: "pointer",
                            fontWeight: "500",
                            fontSize: "14px",
                            backgroundColor: activeTab === item.id ? "#e5e7eb" : "transparent", // Active state gray subtle
                            color: activeTab === item.id ? "#111827" : "#4b5563", // Text contrast matching request
                            transition: "all 0.2s ease"
                        }}
                    >
                        {item.label}
                    </button>
                ))}
            </nav>
        </div>
    )
}

export default Sidebar