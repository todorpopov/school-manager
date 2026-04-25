import { useQuery } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'

export interface Term {
    term_id: number
    name: string
}

export const useGetTerms = (): UseQueryResult<Term[], Error> =>
    useQuery<Term[], Error>({
        queryKey: ['terms'],
        queryFn: async () => {
            const API_URL = import.meta.env.VITE_API_URL as string + '/api'
            try {
                const { data } = await axiosInstance.get<{ data: Term[] }>(`${API_URL}/terms`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
    })

