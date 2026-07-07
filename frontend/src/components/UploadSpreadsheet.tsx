import { useState } from "react"
import { api } from "../api"
import { toast } from "react-toastify"

function UploadSpreadsheet() {
    const [serviceGroupID, setServiceGroupID] = useState("")
    const [serviceName, setServiceName] = useState("")
    const [file, setFile] = useState<File | null>(null)

    async function handleUpload() {
        if (!file) {
            toast.dark("Please select a file to proceed.", {
                style: { backgroundColor: "#374151", color: "#f9fafb" }
            })
            return
        }

        const formData = new FormData()

        formData.append(
            "service_group_id",
            serviceGroupID,
        )

        formData.append(
            "service_name",
            serviceName,
        )

        formData.append(
            "file",
            file,
        )

        try {
            const response = await api.post(
                "/spreadsheet",
                formData,
            )
            // console.log("SUCCESS", response)
            // console.log("DATA", response.data)

            toast.success(response.data.message || "Spreadsheet Ingested Successfully", {
                style: { backgroundColor: "#f3f4f6", color: "#111827", borderLeft: "4px solid #10b981" }
            })
        } catch (err: any) {
            // console.log("ERROR", err)
            const errMsg = err?.response?.data?.error || err?.message || "Upload failed"
            toast.error(errMsg, {
                style: { backgroundColor: "#f3f4f6", color: "#111827", borderLeft: "4px solid #ef4444" }
            })
        }
    }

    return (
        <div>
            <h2>
                Upload Spreadsheet
            </h2>

            <input
                type="text"
                placeholder="Service Group ID"
                value={serviceGroupID}
                onChange={(e) =>
                    setServiceGroupID(
                        e.target.value,
                    )
                }
            />

            <br />

            <input
                type="text"
                placeholder="Service Name"
                value={serviceName}
                onChange={(e) =>
                    setServiceName(
                        e.target.value,
                    )
                }
            />

            <br />

            <input
                type="file"
                accept=".xlsx"
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

export default UploadSpreadsheet