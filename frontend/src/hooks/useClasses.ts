import { useQuery } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'

export interface ClassOption {
    class_id: number
    grade_level: number
    class_name: string
}

export const useGetClasses = (): UseQueryResult<ClassOption[], Error> =>
    useQuery<ClassOption[], Error>({
        queryKey: ['classes'],
        queryFn: async () => {
            const API_URL = import.meta.env.VITE_API_URL as string + '/api'
            try {
                const { data } = await axiosInstance.get<{ data: ClassOption[] }>(`${API_URL}/classes`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
    })

