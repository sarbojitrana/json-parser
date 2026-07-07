import axios from "axios"

export const api = axios.create({
    baseURL: "http://localhost:5555/api",
})

export interface FetchApplicationsParams {
    from: string,
    to: string,
    page: number,
    limit: number
}

export const getApplications = async (params: FetchApplicationsParams) => {
    const response = await api.get("/applications", {
        params: {
            from: params.from,
            to: params.to,
            page: params.page,
            limit: params.limit
        }
    })
    return response.data
}