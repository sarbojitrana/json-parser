import { useState } from "react"
import { api } from "../api"
import { toast } from "react-toastify"

function UploadWorkflow() {
    const [file, setFile] = useState<File | null>(null)

    async function handleUpload() {
        if (!file) {
            toast.dark("Select a JSON workflow file first.", {
                style: { backgroundColor: "#374151", color: "#f9fafb" }
            })
            return
        }

        const formData = new FormData()

        formData.append(
            "file",
            file,
        )

        try {
            const response = await api.post(
                "/workflow",
                formData,
            )

            toast.info(
                <div style={{ fontFamily: "sans-serif", fontSize: "13px" }}>
                    <strong style={{ display: "block", marginBottom: "4px" }}>{response.data.message}</strong>
                    <div>• Events Processed: {response.data.events}</div>
                    <div>• Applications Created: {response.data.applications}</div>
                    <div>• Records Skipped: {response.data.skipped}</div>
                </div>,
                { style: { backgroundColor: "#f9fafb", color: "#111827", borderLeft: "4px solid #2563eb" }, autoClose: 6000 }
            )
        } catch (err: any) {
            // console.log(err)

            const errMsg = err?.response?.data?.error || err.message || "Upload failed"
            toast.error(errMsg, {
                style: { backgroundColor: "#f3f4f6", color: "#111827", borderLeft: "4px solid #ef4444" }
            })
        }
    }

    return (
        <div>
            <h2>Upload Workflow</h2>

            <input
                type="file"
                accept=".json"
                onChange={(e) =>
                    setFile(
                        e.target.files?.[0] ??
                            null,
                    )
                }
            />

            <br />

            <button
                onClick={
                    handleUpload
                }
            >
                Upload
            </button>
        </div>
    )
}

export default UploadWorkflow