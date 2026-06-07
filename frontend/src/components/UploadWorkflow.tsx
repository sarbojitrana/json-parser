import { useState } from "react"
import { api } from "../api"

function UploadWorkflow() {
    const [rootType, setRootType] = useState("")
    const [file, setFile] = useState<File | null>(null)

    async function handleUpload() {
        if (!file) {
            alert("Select a file")
            return
        }

        const formData = new FormData()

        formData.append(
            "root_type",
            rootType,
        )

        formData.append(
            "file",
            file,
        )

        try {
            const response = await api.post(
                "/workflow",
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
        <>
            <h2>Upload Workflow</h2>

            <select
                value={rootType}
                onChange={(e) =>
                    setRootType(
                        e.target.value,
                    )
                }
            >
                <option value="">
                    Select Root Type
                </option>

                <option value="initiated">
                    initiated
                </option>

                <option value="execution">
                    execution
                </option>
            </select>

            <br />

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
        </>
    )
}

export default UploadWorkflow