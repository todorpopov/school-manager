import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Class, CreateClassPayload } from '../types/classes'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export interface ClassOption {
    class_id: number
    grade_level: number
    class_name: string
}

export const useGetClasses = (): UseQueryResult<Class[], Error> =>
    useQuery<Class[], Error>({
        queryKey: ['classes'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Class[] }>(`${API_URL}/classes`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useGetClassesBySchoolId = (schoolId: number | undefined): UseQueryResult<Class[], Error> =>
    useQuery<Class[], Error>({
        queryKey: ['classes-by-school', schoolId],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Class[] }>(`${API_URL}/school/${schoolId}/classes`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        enabled: !!schoolId,
        refetchInterval: 5000,
    })

export const useCreateClass = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateClassPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Class }>(`${API_URL}/class`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['classes'] }),
    })
}

export const useDeleteClass = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/class/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['classes'] }),
    })
}

