import { useState } from "react"
import { api } from "../api"

function UploadSpreadsheet() {
    const [serviceID, setServiceID] = useState("")
    const [serviceName, setServiceName] = useState("")
    const [file, setFile] = useState<File | null>(null)

    async function handleUpload() {
        if (!file) {
            alert("Select a file")
            return
        }

        const formData = new FormData()

        formData.append(
            "service_id",
            serviceID,
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

            alert(
                response.data.message,
            )
        } catch (err) {
            console.error(err)
            alert("Upload failed")
        }
    }

    return (
        <div>
            <h2>Upload Spreadsheet</h2>

            <input
                type="text"
                placeholder="Service ID"
                value={serviceID}
                onChange={(e) =>
                    setServiceID(
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