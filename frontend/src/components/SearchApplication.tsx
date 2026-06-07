import { useState } from "react"
import { api } from "../api"

function SearchApplication() {
    const [applID, setApplID] = useState("")
    const [serviceID, setServiceID] = useState("")
    const [rootType, setRootType] = useState("")

    async function handleSearch() {
        try {
            const response = await api.get(
                `/applications/${applID}`,
                {
                    params: {
                        service_id:
                            serviceID,
                        root_type:
                            rootType,
                    },
                },
            )

            console.log(
                response.data,
            )
        } catch (err) {
            console.error(err)
            alert(
                "Application not found",
            )
        }
    }

    return (
        <>
            <h2>
                Search Application
            </h2>

            <input
                placeholder="Application ID"
                value={applID}
                onChange={(e) =>
                    setApplID(
                        e.target.value,
                    )
                }
            />

            <br />

            <input
                placeholder="Service ID"
                value={serviceID}
                onChange={(e) =>
                    setServiceID(
                        e.target.value,
                    )
                }
            />

            <br />

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

            <button
                onClick={
                    handleSearch
                }
            >
                Search
            </button>
        </>
    )
}

export default SearchApplication