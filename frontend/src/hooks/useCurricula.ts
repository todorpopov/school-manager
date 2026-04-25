import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Curriculum, CreateCurriculumPayload } from '../types/curricula'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetCurricula = (): UseQueryResult<Curriculum[], Error> =>
    useQuery<Curriculum[], Error>({
        queryKey: ['curricula'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Curriculum[] }>(`${API_URL}/curricula`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateCurriculum = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateCurriculumPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Curriculum }>(`${API_URL}/curriculum`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['curricula'] }),
    })
}

export const useDeleteCurriculum = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/curriculum/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['curricula'] }),
    })
}

